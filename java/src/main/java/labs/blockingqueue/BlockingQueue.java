package labs.blockingqueue;

import java.util.Optional;

/** Bounded FIFO; closing rejects puts and lets consumers drain accepted items. */
public final class BlockingQueue<T> {
    public BlockingQueue(int capacity) {
        // TODO: validate and store the capacity.
    }

    /** Waits for space; interruption cancels the wait. */
    public void put(T value) throws InterruptedException {
        throw new UnsupportedOperationException("TODO: implement");
    }

    /** Returns empty only when closed and drained. */
    public Optional<T> take() throws InterruptedException {
        throw new UnsupportedOperationException("TODO: implement");
    }

    public void close() {
        throw new UnsupportedOperationException("TODO: implement");
    }
}
