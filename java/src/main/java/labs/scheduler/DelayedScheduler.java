package labs.scheduler;

import java.util.OptionalLong;
import java.util.function.LongSupplier;

/**
 * Bounded scheduling by id with manual drive or one interruptible owner.
 * {@code labs.algorithms.Structures.IndexedQueue} is available if you want an indexed heap.
 */
public final class DelayedScheduler {
    public DelayedScheduler(int capacity) {
        // TODO: validate and store the capacity.
    }

    public void schedule(String id, long due, Runnable run) {
        throw new UnsupportedOperationException("TODO: implement");
    }

    public boolean cancel(String id) {
        throw new UnsupportedOperationException("TODO: implement");
    }

    public OptionalLong next() {
        throw new UnsupportedOperationException("TODO: implement");
    }

    /** Claims due tasks under the lock, then runs them outside it. */
    public int runDue(long now) {
        throw new UnsupportedOperationException("TODO: implement");
    }

    /** Runs until interrupted; {@code now} supplies monotonic nanoseconds. */
    public void run(LongSupplier now) throws InterruptedException {
        throw new UnsupportedOperationException("TODO: implement");
    }
}
