package labs.spatial;

import java.util.*;

/** Sequential bounded planar quadtree; depth cap prevents duplicate-point recursion. */
public final class QuadTree {
  public record Box(Polygon.Point min, Polygon.Point max) {
    boolean contains(Polygon.Point p) {
      return p.x() >= min.x() && p.x() <= max.x() && p.y() >= min.y() && p.y() <= max.y();
    }

    boolean intersects(Box b) {
      return min.x() <= b.max.x()
          && b.min.x() <= max.x()
          && min.y() <= b.max.y()
          && b.min.y() <= max.y();
    }
  }

  private static final class Node {
    final Box box;
    Map<String, Polygon.Point> ids = new HashMap<>();
    Node[] children;

    Node(Box box) {
      this.box = box;
    }
  }

  private final Node root =
      new Node(new Box(new Polygon.Point(-1000000, -1000000), new Polygon.Point(1000000, 1000000)));
  private final Map<String, Polygon.Point> points = new HashMap<>();

  private static int index(Box b, Polygon.Point p) {
    return (p.x() > (b.min.x() + b.max.x()) / 2 ? 1 : 0)
        + (p.y() > (b.min.y() + b.max.y()) / 2 ? 2 : 0);
  }

  private static Box child(Box b, int i) {
    long x = (b.min.x() + b.max.x()) / 2, y = (b.min.y() + b.max.y()) / 2;
    return new Box(
        new Polygon.Point((i & 1) == 0 ? b.min.x() : x + 1, (i & 2) == 0 ? b.min.y() : y + 1),
        new Polygon.Point((i & 1) == 0 ? x : b.max.x(), (i & 2) == 0 ? y : b.max.y()));
  }

  private static void insert(Node n, String id, Polygon.Point p, int depth) {
    if (n.children == null && (n.ids.size() < 8 || depth == 16)) {
      n.ids.put(id, p);
      return;
    }
    if (n.children == null) {
      n.children = new Node[4];
      for (int i = 0; i < 4; i++) n.children[i] = new Node(child(n.box, i));
      for (var e : n.ids.entrySet())
        insert(n.children[index(n.box, e.getValue())], e.getKey(), e.getValue(), depth + 1);
      n.ids.clear();
    }
    insert(n.children[index(n.box, p)], id, p, depth + 1);
  }

  private static boolean remove(Node n, String id, Polygon.Point p) {
    if (n.children == null) {
      n.ids.remove(id);
      return n.ids.isEmpty();
    }
    remove(n.children[index(n.box, p)], id, p);
    for (var c : n.children) if (c.children != null || !c.ids.isEmpty()) return false;
    n.children = null;
    return true;
  }

  public void put(String id, Polygon.Point p) {
    Objects.requireNonNull(id);
    if (!Polygon.valid(p)) throw new IllegalArgumentException();
    delete(id);
    points.put(id, p);
    insert(root, id, p, 0);
  }

  public boolean delete(String id) {
    var p = points.remove(id);
    if (p == null) return false;
    remove(root, id, p);
    return true;
  }

  public Map<String, Polygon.Point> range(Box b) {
    if (!Polygon.valid(b.min())
        || !Polygon.valid(b.max())
        || b.min().x() > b.max().x()
        || b.min().y() > b.max().y()) throw new IllegalArgumentException();
    var out = new HashMap<String, Polygon.Point>();
    visit(root, b, out);
    return Map.copyOf(out);
  }

  private static void visit(Node n, Box b, Map<String, Polygon.Point> out) {
    if (!n.box.intersects(b)) return;
    if (n.children != null) for (var c : n.children) visit(c, b, out);
    else
      n.ids.forEach(
          (id, p) -> {
            if (b.contains(p)) out.put(id, p);
          });
  }
}
