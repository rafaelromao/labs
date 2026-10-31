package labs.spatial;

import java.util.*;

/** Exact planar grid predicates on [-1_000_000,1_000_000], not degrees. */
public final class Polygon {
  public record Point(long x, long y) {}

  private final List<Point> shell;
  private final List<List<Point>> holes;

  public static boolean valid(Point p) {
    return p != null
        && p.x() >= -1000000
        && p.x() <= 1000000
        && p.y() >= -1000000
        && p.y() <= 1000000;
  }

  private static long cross(Point a, Point b, Point c) {
    return (b.x() - a.x()) * (c.y() - a.y()) - (b.y() - a.y()) * (c.x() - a.x());
  }

  private static boolean on(Point a, Point b, Point p) {
    return cross(a, b, p) == 0
        && p.x() >= Math.min(a.x(), b.x())
        && p.x() <= Math.max(a.x(), b.x())
        && p.y() >= Math.min(a.y(), b.y())
        && p.y() <= Math.max(a.y(), b.y());
  }

  private static boolean intersects(Point a, Point b, Point c, Point d) {
    return on(a, b, c)
        || on(a, b, d)
        || on(c, d, a)
        || on(c, d, b)
        || (Long.signum(cross(a, b, c)) * Long.signum(cross(a, b, d)) < 0
            && Long.signum(cross(c, d, a)) * Long.signum(cross(c, d, b)) < 0);
  }

  private static int relation(List<Point> r, Point p) {
    boolean inside = false;
    for (int i = 0; i < r.size(); i++) {
      var a = r.get(i);
      var b = r.get((i + 1) % r.size());
      if (on(a, b, p)) return 0;
      if ((a.y() > p.y()) != (b.y() > p.y()) && (cross(a, b, p) > 0) == (b.y() > a.y()))
        inside = !inside;
    }
    return inside ? 1 : -1;
  }

  private static boolean ringValid(List<Point> r) {
    if (r.size() < 3 || r.size() > 4096) return false;
    long area = 0;
    for (int i = 0; i < r.size(); i++) {
      var a = r.get(i);
      var b = r.get((i + 1) % r.size());
      if (!valid(a) || !valid(b) || a.equals(b)) return false;
      area += a.x() * b.y() - a.y() * b.x();
      for (int j = i + 1; j < r.size(); j++) {
        if (j == i + 1 || (i == 0 && j == r.size() - 1)) continue;
        if (intersects(a, b, r.get(j), r.get((j + 1) % r.size()))) return false;
      }
      if (cross(a, b, r.get((i + 2) % r.size())) == 0) return false;
    }
    return area != 0;
  }

  private static boolean ringsIntersect(List<Point> a, List<Point> b) {
    for (int i = 0; i < a.size(); i++)
      for (int j = 0; j < b.size(); j++)
        if (intersects(a.get(i), a.get((i + 1) % a.size()), b.get(j), b.get((j + 1) % b.size())))
          return true;
    return false;
  }

  public Polygon(List<Point> shell, List<List<Point>> holes) {
    this.shell = List.copyOf(shell);
    if (!ringValid(this.shell)) throw new IllegalArgumentException("invalid shell");
    var copied = new ArrayList<List<Point>>();
    for (var input : holes) {
      var h = List.copyOf(input);
      if (!ringValid(h) || relation(this.shell, h.getFirst()) != 1 || ringsIntersect(this.shell, h))
        throw new IllegalArgumentException("invalid hole");
      for (var other : copied)
        if (ringsIntersect(other, h)
            || relation(other, h.getFirst()) >= 0
            || relation(h, other.getFirst()) >= 0)
          throw new IllegalArgumentException("overlapping holes");
      copied.add(h);
    }
    this.holes = List.copyOf(copied);
  }

  public boolean covers(Point p) {
    if (!valid(p)) throw new IllegalArgumentException("coordinate domain");
    int r = relation(shell, p);
    if (r < 0) return false;
    if (r == 0) return true;
    for (var h : holes) {
      r = relation(h, p);
      if (r == 0) return true;
      if (r > 0) return false;
    }
    return true;
  }
}
