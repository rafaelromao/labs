package labs.trajectory;

import java.util.List;
import labs.nearby.NearbyIndex;

/** One bounded track, with timestamps in nonnegative milliseconds. */
public final class Trajectory {
    public record Observation(long millis, NearbyIndex.Point position) {}

    public record Segment(long fromMillis, long toMillis, double meters, double speedMetersPerSecond,
            boolean implausible) {}

    public Trajectory(int capacity) {
        // TODO: validate and store the capacity.
    }

    public void put(Observation o) {
        throw new UnsupportedOperationException("TODO: implement");
    }

    public List<Segment> segments(double maxSpeed) {
        throw new UnsupportedOperationException("TODO: implement");
    }
}
