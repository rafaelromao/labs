package labs.routing;

import java.util.*;

/** Expanded incoming-edge states; route nodes contain road IDs. */
public final class TurnPlanner {
  public record Road(String id, String from, String to, double cost) {}

  public record Turn(String incoming, String outgoing) {}

  private TurnPlanner() {}

  public static RoutePlanner.Route find(
      List<Road> roads, Set<Turn> blocked, String start, String goal) {
    var byID = new HashMap<String, Road>();
    var junctions = new HashSet<String>();
    for (var r : roads) {
      Objects.requireNonNull(r.id());
      Objects.requireNonNull(r.from());
      Objects.requireNonNull(r.to());
      if (r.id().isEmpty()
          || !Double.isFinite(r.cost())
          || r.cost() < 0
          || byID.putIfAbsent(r.id(), r) != null)
        throw new IllegalArgumentException("invalid road");
      junctions.add(r.from());
      junctions.add(r.to());
    }
    if (!junctions.contains(start) || !junctions.contains(goal))
      throw new IllegalArgumentException("unknown endpoint");
    for (var t : blocked) {
      var a = byID.get(t.incoming());
      var b = byID.get(t.outgoing());
      if (a == null || b == null || !a.to().equals(b.from()))
        throw new IllegalArgumentException("invalid turn");
    }
    var g = new HashMap<String, List<RoutePlanner.Edge>>();
    g.put("start", new ArrayList<>());
    g.put("goal", new ArrayList<>());
    for (var r : roads) {
      var edges = new ArrayList<RoutePlanner.Edge>();
      g.put("road:" + r.id(), edges);
      if (r.from().equals(start))
        g.get("start").add(new RoutePlanner.Edge("road:" + r.id(), r.cost()));
      if (r.to().equals(goal)) edges.add(new RoutePlanner.Edge("goal", 0));
    }
    if (start.equals(goal)) g.get("start").add(new RoutePlanner.Edge("goal", 0));
    for (var a : roads)
      for (var b : roads)
        if (a.to().equals(b.from()) && !blocked.contains(new Turn(a.id(), b.id())))
          g.get("road:" + a.id()).add(new RoutePlanner.Edge("road:" + b.id(), b.cost()));
    var r = RoutePlanner.find(g, "start", "goal", null);
    return new RoutePlanner.Route(
        r.nodes().stream().filter(n -> n.startsWith("road:")).map(n -> n.substring(5)).toList(),
        r.cost(),
        r.found());
  }
}
