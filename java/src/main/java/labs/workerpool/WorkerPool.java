package labs.workerpool;

import java.util.concurrent.Callable;
import java.util.concurrent.Future;
import java.util.concurrent.TimeUnit;

/** Fixed workers, bounded waiting queue, immediate rejection and drain-on-close. */
public final class WorkerPool {
    public WorkerPool(int workers, int capacity) {
        // TODO: validate and start the workers.
    }

    /** Throws RejectedExecutionException after close or when saturated. */
    public <T> Future<T> submit(Callable<T> task) {
        throw new UnsupportedOperationException("TODO: implement");
    }

    /** Stops admission immediately; accepted tasks drain. Does not wait. */
    public void close() {
        throw new UnsupportedOperationException("TODO: implement");
    }

    /** Rejects admission, cancels queued futures, interrupts running tasks; returns the number cancelled. */
    public int abort() {
        throw new UnsupportedOperationException("TODO: implement");
    }

    /** False means the timeout elapsed, not that unfinished tasks were stopped. */
    public boolean awaitTermination(long timeout, TimeUnit unit) throws InterruptedException {
        throw new UnsupportedOperationException("TODO: implement");
    }
}
