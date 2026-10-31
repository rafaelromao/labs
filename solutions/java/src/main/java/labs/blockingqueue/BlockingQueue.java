package labs.blockingqueue;

import java.util.ArrayDeque;
import java.util.Objects;
import java.util.Optional;

/** Bounded FIFO; closing rejects puts and lets consumers drain accepted items. */
public final class BlockingQueue<T> {
    private final int capacity;
    private final ArrayDeque<T> items = new ArrayDeque<>();
    private boolean closed;

    public BlockingQueue(int capacity) {
        if (capacity <= 0) throw new IllegalArgumentException("capacity must be positive");
        this.capacity = capacity;
    }

    /** Waits for space; interruption cancels the wait. */
    public synchronized void put(T value) throws InterruptedException {
        Objects.requireNonNull(value, "value");
        while (!closed && items.size() == capacity) wait();
        if (closed) throw new IllegalStateException("queue closed");
        items.addLast(value);
        notifyAll();
    }

    /** Returns empty only when closed and drained. */
    public synchronized Optional<T> take() throws InterruptedException {
        while (!closed && items.isEmpty()) wait();
        if (items.isEmpty()) return Optional.empty();
        T result = items.removeFirst();
        notifyAll();
        return Optional.of(result);
    }

    public synchronized void close() {
        if (!closed) { closed = true; notifyAll(); }
    }
}
