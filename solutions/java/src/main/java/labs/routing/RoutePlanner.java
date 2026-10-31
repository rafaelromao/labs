package labs.routing;

import java.util.*;
import java.util.function.ToDoubleFunction;

/** Graph inputs must remain unchanged throughout a call. */
public final class RoutePlanner {
    public record Edge(String to, double cost) {}
    public record Route(List<String> nodes, double cost, boolean found) {}
    private record Item(String node, double cost, double priority) {}
    private RoutePlanner() {}

    /** Null heuristic selects Dijkstra; A* requires an admissible heuristic. */
    public static Route find(Map<String, List<Edge>> graph, String start, String goal, ToDoubleFunction<String> heuristic) {
        if (!graph.containsKey(start) || !graph.containsKey(goal)) throw new IllegalArgumentException("unknown endpoint");
        var h = new HashMap<String, Double>();
        for (var entry : graph.entrySet()) {
            Objects.requireNonNull(entry.getKey());
            double estimate = heuristic == null ? 0 : heuristic.applyAsDouble(entry.getKey());
            if (!Double.isFinite(estimate) || estimate < 0) throw new IllegalArgumentException("invalid heuristic");
            h.put(entry.getKey(), estimate);
            for (var edge : entry.getValue()) {
                if (!graph.containsKey(edge.to()) || !Double.isFinite(edge.cost()) || edge.cost() < 0)
                    throw new IllegalArgumentException("invalid edge");
            }
        }
        if (h.get(goal) != 0) throw new IllegalArgumentException("goal heuristic must be zero");
        var distance = new HashMap<String, Double>();
        var previous = new HashMap<String, String>();
        var queue = new PriorityQueue<Item>(Comparator.comparingDouble(Item::priority).thenComparing(Item::node));
        distance.put(start, 0.0); queue.add(new Item(start, 0, h.get(start)));
        while (!queue.isEmpty()) {
            var current = queue.remove();
            if (current.cost() != distance.get(current.node())) continue;
            if (current.node().equals(goal)) {
                var nodes = new ArrayList<String>();
                for (String node = goal; node != null; node = previous.get(node)) nodes.add(node);
                Collections.reverse(nodes);
                return new Route(List.copyOf(nodes), current.cost(), true);
            }
            for (var edge : graph.get(current.node())) {
                double candidate = current.cost() + edge.cost();
                if (!Double.isFinite(candidate) || !Double.isFinite(candidate + h.get(edge.to())))
                    throw new IllegalArgumentException("numeric overflow");
                if (!distance.containsKey(edge.to()) || candidate < distance.get(edge.to())) {
                    distance.put(edge.to(), candidate); previous.put(edge.to(), current.node());
                    queue.add(new Item(edge.to(), candidate, candidate + h.get(edge.to())));
                }
            }
        }
        return new Route(List.of(), 0, false);
    }
}
