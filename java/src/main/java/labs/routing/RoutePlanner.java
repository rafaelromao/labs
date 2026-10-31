package labs.routing;

import java.util.List;
import java.util.Map;
import java.util.function.ToDoubleFunction;

/** Graph inputs must remain unchanged throughout a call. */
public final class RoutePlanner {
    public record Edge(String to, double cost) {}

    public record Route(List<String> nodes, double cost, boolean found) {}

    private RoutePlanner() {}

    /** A null heuristic selects Dijkstra; A* requires an admissible heuristic. */
    public static Route find(Map<String, List<Edge>> graph, String start, String goal,
            ToDoubleFunction<String> heuristic) {
        throw new UnsupportedOperationException("TODO: implement");
    }
}
