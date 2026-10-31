package labs.routing;

import static org.junit.jupiter.api.Assertions.*;

import java.util.*;
import org.junit.jupiter.api.Test;

class TurnPlannerTest {
  @Test
  void alternateArrival() {
    var roads =
        List.of(
            new TurnPlanner.Road("ab", "a", "b", 1),
            new TurnPlanner.Road("ac", "a", "c", 2),
            new TurnPlanner.Road("cb", "c", "b", 1),
            new TurnPlanner.Road("bd", "b", "d", 1));
    assertEquals(
        new RoutePlanner.Route(List.of("ac", "cb", "bd"), 4, true),
        TurnPlanner.find(roads, Set.of(new TurnPlanner.Turn("ab", "bd")), "a", "d"));
  }
}
