package labs.workerpool;

import java.util.Objects;
import java.util.concurrent.ArrayBlockingQueue;
import java.util.concurrent.Callable;
import java.util.concurrent.Future;
import java.util.concurrent.ThreadPoolExecutor;
import java.util.concurrent.TimeUnit;

/** Fixed workers, bounded waiting queue, immediate rejection and drain-on-close. */
public final class WorkerPool {
  private final ThreadPoolExecutor executor;

  public WorkerPool(int workers, int capacity) {
    if (workers <= 0 || capacity <= 0)
      throw new IllegalArgumentException("positive workers and capacity required");
    executor =
        new ThreadPoolExecutor(
            workers,
            workers,
            0,
            TimeUnit.MILLISECONDS,
            new ArrayBlockingQueue<>(capacity),
            new ThreadPoolExecutor.AbortPolicy());
  }

  /** Throws RejectedExecutionException after close or when saturated. */
  public <T> Future<T> submit(Callable<T> task) {
    return executor.submit(Objects.requireNonNull(task));
  }

  /** Stops admission immediately; accepted tasks drain. Does not wait. */
  public void close() {
    executor.shutdown();
  }

  /**
   * Reject admission, cancel queued futures, interrupt running cooperative tasks. Await termination
   * separately; interruption does not prove task completion.
   */
  public int abort() {
    var pending = executor.shutdownNow();
    for (var task : pending) if (task instanceof Future<?> future) future.cancel(false);
    return pending.size();
  }

  /** False means the timeout elapsed, not that unfinished tasks were stopped. */
  public boolean awaitTermination(long timeout, TimeUnit unit) throws InterruptedException {
    return executor.awaitTermination(timeout, unit);
  }
}
