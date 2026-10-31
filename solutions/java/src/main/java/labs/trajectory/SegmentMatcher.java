package labs.trajectory;

import java.util.*;

/** Local projected metres; pointwise approximation, not road-path inference. */
public final class SegmentMatcher {
  public record XY(double x, double y) {}

  public record Road(String id, XY a, XY b) {}

  public record Match(String id, XY position, double distance, boolean ambiguous) {}

  private static void validate(XY p) {
    if (!Double.isFinite(p.x())
        || !Double.isFinite(p.y())
        || Math.abs(p.x()) > 1e6
        || Math.abs(p.y()) > 1e6) throw new IllegalArgumentException();
  }

  public static XY project(XY p, XY a, XY b) {
    validate(p);
    validate(a);
    validate(b);
    double dx = b.x() - a.x(), dy = b.y() - a.y(), length = dx * dx + dy * dy;
    if (length == 0) return a;
    double t = Math.max(0, Math.min(1, ((p.x() - a.x()) * dx + (p.y() - a.y()) * dy) / length));
    return new XY(a.x() + t * dx, a.y() + t * dy);
  }

  public static Match nearest(XY p, List<Road> roads, double tolerance) {
    if (roads.isEmpty() || !Double.isFinite(tolerance) || tolerance < 0)
      throw new IllegalArgumentException();
    var matches = new ArrayList<Match>();
    var seen = new HashSet<String>();
    for (var road : roads) {
      Objects.requireNonNull(road.id());
      if (!seen.add(road.id())) throw new IllegalArgumentException("duplicate ID");
      var q = project(p, road.a(), road.b());
      matches.add(new Match(road.id(), q, Math.hypot(p.x() - q.x(), p.y() - q.y()), false));
    }
    matches.sort(Comparator.comparingDouble(Match::distance).thenComparing(Match::id));
    var first = matches.getFirst();
    return new Match(
        first.id(),
        first.position(),
        first.distance(),
        matches.size() > 1 && matches.get(1).distance() - first.distance() <= tolerance);
  }
}
