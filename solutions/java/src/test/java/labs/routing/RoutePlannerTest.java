package labs.routing;

import static org.junit.jupiter.api.Assertions.*;
import java.util.List;
import java.util.Map;
import java.util.function.ToDoubleFunction;
import org.junit.jupiter.api.Test;

class RoutePlannerTest {
    @Test void shortestPaths() {
        var graph = Map.of("a", List.of(new RoutePlanner.Edge("b", 4), new RoutePlanner.Edge("c", 1)),
                "c", List.of(new RoutePlanner.Edge("b", 1), new RoutePlanner.Edge("z", 10)),
                "b", List.of(new RoutePlanner.Edge("z", 1)), "z", List.<RoutePlanner.Edge>of(), "isolated", List.<RoutePlanner.Edge>of());
        ToDoubleFunction<String> h = node -> Map.of("a", 3.0, "c", 2.0, "b", 1.0).getOrDefault(node, 0.0);
        var expected = new RoutePlanner.Route(List.of("a", "c", "b", "z"), 3, true);
        assertEquals(expected, RoutePlanner.find(graph, "a", "z", null));
        assertEquals(expected, RoutePlanner.find(graph, "a", "z", h));
        assertFalse(RoutePlanner.find(graph, "z", "isolated", null).found());
        assertEquals(List.of("a"), RoutePlanner.find(graph, "a", "a", null).nodes());
    }
    @Test void invalidCosts() {
        for (double cost : new double[]{-1, Double.NaN, Double.POSITIVE_INFINITY}) {
            var graph = Map.of("a", List.of(new RoutePlanner.Edge("b", cost)), "b", List.<RoutePlanner.Edge>of());
            assertThrows(IllegalArgumentException.class, () -> RoutePlanner.find(graph, "a", "b", null));
        }
    }
}
