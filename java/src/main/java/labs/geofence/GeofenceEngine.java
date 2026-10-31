package labs.geofence;

import labs.nearby.NearbyIndex;

/** One immutable circular fence; per-object sequences and spatial hysteresis. */
public final class GeofenceEngine {
    public record Result(boolean inside, boolean changed, boolean stale, long sequence) {}

    public GeofenceEngine(NearbyIndex.Point center, double radius, double margin, int limit) {
        // TODO: validate and store the configuration.
    }

    public Result observe(String id, long sequence, NearbyIndex.Point position) {
        throw new UnsupportedOperationException("TODO: implement");
    }

    public void forget(String id) {
        throw new UnsupportedOperationException("TODO: implement");
    }
}
