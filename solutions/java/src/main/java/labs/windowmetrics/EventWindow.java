package labs.windowmetrics;

import java.util.*;

/** Exact bounded window with a nondecreasing, nonnegative caller watermark. */
public final class EventWindow {
    public record Event(long time, long value) {}
    public record Snapshot(int count, OptionalLong maximum, List<Long> top) {}
    private final long width;
    private final int capacity;
    private long watermark;
    private boolean initialized;
    private final ArrayDeque<Event> events = new ArrayDeque<>();
    public EventWindow(long width, int capacity) {
        if (width<=0 || capacity<=0) throw new IllegalArgumentException("positive width and capacity required");
        this.width=width; this.capacity=capacity;
    }
    private void advance(long now) {
        if (now<0 || (initialized && now<watermark)) throw new IllegalArgumentException("invalid or late time");
        watermark=now; initialized=true;
        while (!events.isEmpty() && now-events.getFirst().time()>=width) events.removeFirst();
    }
    /** Advances time even if the bounded collection is still full. */
    public synchronized void add(Event event) {
        advance(event.time());
        if (events.size()==capacity) throw new IllegalStateException("capacity reached");
        events.addLast(event);
    }
    public synchronized Snapshot query(long now,int k) {
        if (k<0) throw new IllegalArgumentException("negative k");
        advance(now);
        var values=new ArrayList<Long>(); for (var event:events) values.add(event.value());
        values.sort(Comparator.reverseOrder());
        return new Snapshot(values.size(),values.isEmpty()?OptionalLong.empty():OptionalLong.of(values.getFirst()),
                List.copyOf(values.subList(0,Math.min(k,values.size()))));
    }
}
