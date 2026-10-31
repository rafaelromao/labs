package labs.routing;

import java.util.List;
import java.util.Set;

/** Searches over incoming roads so turn restrictions can be enforced; route nodes are road ids. */
public final class TurnPlanner {
    public record Road(String id, String from, String to, double cost) {}

    public record Turn(String incoming, String outgoing) {}

    private TurnPlanner() {}

    public static RoutePlanner.Route find(List<Road> roads, Set<Turn> blocked, String start, String goal) {
        throw new UnsupportedOperationException("TODO: implement");
    }
}
