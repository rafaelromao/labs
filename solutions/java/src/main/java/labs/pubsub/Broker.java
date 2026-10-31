package labs.pubsub;

import java.util.*;
import java.util.concurrent.ArrayBlockingQueue;

/** Bounded local broadcast, with drop-newest per full subscriber queue. */
public final class Broker<T> {
    public record Delivery(int delivered,int dropped) {}
    public final class Subscription {
        private final ArrayBlockingQueue<T> messages=new ArrayBlockingQueue<>(capacity);
        private volatile boolean closed;
        private Subscription() {}
        public Optional<T> poll() { return Optional.ofNullable(messages.poll()); }
        public boolean isClosed() { return closed; }
        public void unsubscribe() {
            synchronized (Broker.this) { if (subscribers.remove(this)) closed=true; }
        }
    }
    private final int limit,capacity;
    private boolean closed;
    private final Set<Subscription> subscribers=new HashSet<>();
    public Broker(int limit,int capacity) {
        if (limit<=0 || capacity<=0) throw new IllegalArgumentException("positive limits required");
        this.limit=limit; this.capacity=capacity;
    }
    public synchronized Subscription subscribe() {
        if (closed || subscribers.size()==limit) throw new IllegalStateException("closed or subscriber limit reached");
        var s=new Subscription(); subscribers.add(s); return s;
    }
    public synchronized Delivery publish(T value) {
        Objects.requireNonNull(value);
        if (closed) throw new IllegalStateException("closed");
        int delivered=0,dropped=0;
        for (var s:subscribers) { if (s.messages.offer(value)) delivered++; else dropped++; }
        return new Delivery(delivered,dropped);
    }
    public synchronized void close() {
        closed=true; for (var s:subscribers) s.closed=true; subscribers.clear();
    }
}
