package labs.spatial;

import static org.junit.jupiter.api.Assertions.*;

import java.nio.file.*;
import java.util.*;
import org.junit.jupiter.api.Test;

class SpatialTest {
  private List<Polygon.Point> square(long a, long b) {
    return List.of(
        new Polygon.Point(a, a),
        new Polygon.Point(b, a),
        new Polygon.Point(b, b),
        new Polygon.Point(a, b));
  }

  @Test
  void sharedPolygonFixtureAndRevision() throws Exception {
    var p = new Polygon(square(0, 10), List.of(square(3, 7)));
    var file =
        Path.of(System.getProperty("labs.root", "."))
            .resolve("fixtures/polygon.tsv");
    for (var line : Files.readAllLines(file)) {
      if (line.startsWith("#") || line.isBlank()) continue;
      var a = line.split("\\s+");
      assertEquals(
          Boolean.parseBoolean(a[2]),
          p.covers(new Polygon.Point(Long.parseLong(a[0]), Long.parseLong(a[1]))));
    }
    var f = new VersionedFence(p, 1, 2);
    assertTrue(f.observe("a", 1, new Polygon.Point(1, 1)).orElseThrow().inside());
    var changes = f.replace(new Polygon(square(2, 8), List.of()), 2);
    assertEquals(1, changes.size());
    assertFalse(changes.getFirst().inside());
    assertTrue(f.observe("a", 1, new Polygon.Point(5, 5)).isEmpty());
  }

  @Test
  void invalidGeometry() {
    assertThrows(
        IllegalArgumentException.class,
        () ->
            new Polygon(
                List.of(
                    new Polygon.Point(0, 0),
                    new Polygon.Point(5, 5),
                    new Polygon.Point(0, 5),
                    new Polygon.Point(5, 0)),
                List.of()));
    assertThrows(
        IllegalArgumentException.class, () -> new Polygon(square(0, 10), List.of(square(0, 2))));
    assertThrows(
        IllegalArgumentException.class,
        () -> new Polygon(square(0, 10), List.of(square(2, 6), square(4, 8))));
  }

  @Test
  void quadtreeModel() {
    var q = new QuadTree();
    var model = new HashMap<String, Polygon.Point>();
    var r = new Random(17);
    for (int step = 0; step < 1200; step++) {
      String id = "" + r.nextInt(100);
      var p = new Polygon.Point(r.nextInt(101) - 50, r.nextInt(101) - 50);
      if (step % 4 == 0) {
        q.delete(id);
        model.remove(id);
      } else {
        q.put(id, p);
        model.put(id, p);
      }
      var box = new QuadTree.Box(new Polygon.Point(-20, -20), new Polygon.Point(20, 20));
      var want = new HashMap<String, Polygon.Point>();
      model.forEach(
          (key, point) -> {
            if (box.contains(point)) want.put(key, point);
          });
      assertEquals(want, q.range(box));
    }
    for (int i = 0; i < 100; i++) q.put("" + i, new Polygon.Point(0, 0));
    assertEquals(
        100, q.range(new QuadTree.Box(new Polygon.Point(0, 0), new Polygon.Point(0, 0))).size());
  }
}
