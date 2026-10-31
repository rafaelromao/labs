package labs.trajectory;

import static org.junit.jupiter.api.Assertions.*;

import java.util.*;
import org.junit.jupiter.api.Test;

class SegmentMatcherTest {
  @Test
  void projectionAndAmbiguity() {
    assertEquals(
        new SegmentMatcher.XY(5, 0),
        SegmentMatcher.project(
            new SegmentMatcher.XY(5, 4),
            new SegmentMatcher.XY(0, 0),
            new SegmentMatcher.XY(10, 0)));
    var m =
        SegmentMatcher.nearest(
            new SegmentMatcher.XY(5, 1),
            List.of(
                new SegmentMatcher.Road(
                    "a", new SegmentMatcher.XY(0, 0), new SegmentMatcher.XY(10, 0)),
                new SegmentMatcher.Road(
                    "b", new SegmentMatcher.XY(0, 2), new SegmentMatcher.XY(10, 2))),
            0.01);
    assertEquals("a", m.id());
    assertTrue(m.ambiguous());
  }
}
