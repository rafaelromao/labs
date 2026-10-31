package labs.windowmetrics;

import java.util.List;
import java.util.OptionalLong;

/** Exact bounded window with a nondecreasing, nonnegative caller watermark. */
public final class EventWindow {
    public record Event(long time, long value) {}

    public record Snapshot(int count, OptionalLong maximum, List<Long> top) {}

    public EventWindow(long width, int capacity) {
        // TODO: validate and store the configuration.
    }

    /** Advances time even if the bounded collection is still full. */
    public void add(Event event) {
        throw new UnsupportedOperationException("TODO: implement");
    }

    public Snapshot query(long now, int k) {
        throw new UnsupportedOperationException("TODO: implement");
    }
}
