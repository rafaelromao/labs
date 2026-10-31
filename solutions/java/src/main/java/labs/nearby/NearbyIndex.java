package labs.nearby;

import java.util.*;

/** Exhaustive spherical search. Coordinates are latitude, longitude in degrees. */
public final class NearbyIndex {
  public static final double EARTH_RADIUS_METERS = 6371008.8;

  public record Point(double latitude, double longitude) {}

  public record Hit(String id, double distanceMeters) {}

  private final Map<String, Point> points = new HashMap<>();
  private final Map<Integer, Set<String>> bands = new HashMap<>();

  private static int band(Point p) {
    return (int) Math.floor(p.latitude());
  }

  private void removeBand(String id, Point p) {
    var ids = bands.get(band(p));
    ids.remove(id);
    if (ids.isEmpty()) bands.remove(band(p));
  }

  private static void validate(Point p) {
    Objects.requireNonNull(p);
    if (!Double.isFinite(p.latitude())
        || !Double.isFinite(p.longitude())
        || p.latitude() < -90
        || p.latitude() > 90
        || p.longitude() < -180
        || p.longitude() > 180) throw new IllegalArgumentException("invalid coordinates");
  }

  public synchronized void put(String id, Point p) {
    Objects.requireNonNull(id);
    validate(p);
    var old = points.get(id);
    if (old != null) removeBand(id, old);
    points.put(id, p);
    bands.computeIfAbsent(band(p), b -> new HashSet<>()).add(id);
  }

  public synchronized boolean delete(String id) {
    var old = points.remove(Objects.requireNonNull(id));
    if (old == null) return false;
    removeBand(id, old);
    return true;
  }

  public static double distance(Point a, Point b) {
    validate(a);
    validate(b);
    double dlat = Math.toRadians(b.latitude() - a.latitude());
    double dlon = Math.toRadians(b.longitude() - a.longitude());
    double h =
        Math.pow(Math.sin(dlat / 2), 2)
            + Math.cos(Math.toRadians(a.latitude()))
                * Math.cos(Math.toRadians(b.latitude()))
                * Math.pow(Math.sin(dlon / 2), 2);
    return 2 * EARTH_RADIUS_METERS * Math.asin(Math.sqrt(Math.max(0, Math.min(1, h))));
  }

  public synchronized List<Hit> radius(Point center, double meters) {
    validate(center);
    if (!Double.isFinite(meters) || meters < 0)
      throw new IllegalArgumentException("invalid radius");
    var hits = new ArrayList<Hit>();
    double delta = meters / EARTH_RADIUS_METERS * 180 / Math.PI + 1e-9;
    for (var band : bands.entrySet()) {
      if (band.getKey() > center.latitude() + delta
          || band.getKey() + 1 < center.latitude() - delta) continue;
      for (var id : band.getValue()) {
        double d = distance(center, points.get(id));
        if (d <= meters) hits.add(new Hit(id, d));
      }
    }
    hits.sort(Comparator.comparingDouble(Hit::distanceMeters).thenComparing(Hit::id));
    return List.copyOf(hits);
  }

  public synchronized List<Hit> nearest(Point center, int k) {
    if (k < 0) throw new IllegalArgumentException("negative k");
    var hits = radius(center, Math.PI * EARTH_RADIUS_METERS);
    return List.copyOf(hits.subList(0, Math.min(k, hits.size())));
  }
}
