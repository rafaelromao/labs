package labs.nearby;

import static org.junit.jupiter.api.Assertions.*;

import org.junit.jupiter.api.Test;

class NearbyIndexTest {
  @Test
  void indexedMovesMatchExhaustiveModel() {
    var r = new java.util.Random(37);
    var index = new NearbyIndex();
    var model = new java.util.HashMap<String, NearbyIndex.Point>();
    for (int n = 0; n < 1000; n++) {
      String id = "" + r.nextInt(100);
      var p = new NearbyIndex.Point(r.nextDouble() * 180 - 90, r.nextDouble() * 360 - 180);
      if (n % 7 == 0) {
        index.delete(id);
        model.remove(id);
      } else {
        index.put(id, p);
        model.put(id, p);
      }
      var center = new NearbyIndex.Point(r.nextDouble() * 180 - 90, r.nextDouble() * 360 - 180);
      double radius = r.nextDouble() * 2e7;
      var expected = new java.util.ArrayList<NearbyIndex.Hit>();
      model.forEach(
          (key, point) -> {
            double d = NearbyIndex.distance(center, point);
            if (d <= radius) expected.add(new NearbyIndex.Hit(key, d));
          });
      expected.sort(
          java.util.Comparator.comparingDouble(NearbyIndex.Hit::distanceMeters)
              .thenComparing(NearbyIndex.Hit::id));
      assertEquals(expected, index.radius(center, radius));
    }
  }

  @Test
  void antimeridianMovementAndTies() {
    var index = new NearbyIndex();
    index.put("east", new NearbyIndex.Point(0, 179.9));
    index.put("west", new NearbyIndex.Point(0, -179.9));
    index.put("far", new NearbyIndex.Point(0, 0));
    var center = new NearbyIndex.Point(0, 180);
    assertEquals(2, index.radius(center, 12000).size());
    index.put("far", center);
    assertEquals("far", index.nearest(center, 1).getFirst().id());
    assertTrue(index.delete("far"));
    assertFalse(index.delete("far"));
    index.put("a", center);
    index.put("b", center);
    assertEquals("a", index.radius(center, 0).getFirst().id());
  }

  @Test
  void invalidAndAntipodal() {
    var index = new NearbyIndex();
    assertThrows(
        IllegalArgumentException.class, () -> index.put("a", new NearbyIndex.Point(91, 0)));
    assertThrows(
        IllegalArgumentException.class,
        () -> index.radius(new NearbyIndex.Point(0, 0), Double.NaN));
    assertEquals(
        Math.PI * NearbyIndex.EARTH_RADIUS_METERS,
        NearbyIndex.distance(new NearbyIndex.Point(0, 0), new NearbyIndex.Point(0, 180)),
        0.01);
  }
}
