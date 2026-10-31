package labs.trajectory;

import java.util.List;

/** Local projected metres; pointwise approximation, not road-path inference. */
public final class SegmentMatcher {
    public record XY(double x, double y) {}

    public record Road(String id, XY a, XY b) {}

    public record Match(String id, XY position, double distance, boolean ambiguous) {}

    private SegmentMatcher() {}

    public static XY project(XY p, XY a, XY b) {
        throw new UnsupportedOperationException("TODO: implement");
    }

    public static Match nearest(XY p, List<Road> roads, double tolerance) {
        throw new UnsupportedOperationException("TODO: implement");
    }
}
