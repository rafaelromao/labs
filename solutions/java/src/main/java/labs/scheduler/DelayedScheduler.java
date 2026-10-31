package labs.scheduler;

import labs.algorithms.Structures;
import java.util.*;
import java.util.function.LongSupplier;

/** Bounded indexed scheduling with manual drive or one interruptible owner. */
public final class DelayedScheduler {
  private record Task(long due, Runnable run) {}

  private final int capacity;
  private final Map<String, Task> tasks = new HashMap<>();
  private final Structures.IndexedQueue queue = new Structures.IndexedQueue();

  public DelayedScheduler(int capacity) {
    if (capacity <= 0) throw new IllegalArgumentException("positive capacity required");
    this.capacity = capacity;
  }

  public synchronized void schedule(String id, long due, Runnable run) {
    Objects.requireNonNull(id);
    Objects.requireNonNull(run);
    if (!tasks.containsKey(id) && tasks.size() == capacity)
      throw new IllegalStateException("capacity");
    tasks.put(id, new Task(due, run));
    queue.put(id, due);
    notifyAll();
  }

  public synchronized boolean cancel(String id) {
    Objects.requireNonNull(id);
    if (!queue.remove(id)) return false;
    tasks.remove(id);
    notifyAll();
    return true;
  }

  public synchronized OptionalLong next() {
    var p = queue.peek();
    return p.isEmpty() ? OptionalLong.empty() : OptionalLong.of(p.get().value());
  }

  /** Claims a snapshot; cancellation cannot recall claimed callbacks. */
  public int runDue(long now) {
    var work = new ArrayList<Task>();
    synchronized (this) {
      while (queue.peek().isPresent() && queue.peek().get().value() <= now) {
        var p = queue.pop().orElseThrow();
        work.add(tasks.remove(p.key()));
      }
    }
    for (var task : work) task.run().run();
    return work.size();
  }

  /**
   * now supplies monotonic nanoseconds. Interrupt and join the owner to stop. No other thread may
   * drive callbacks concurrently if serial execution is required.
   */
  public void run(LongSupplier now) throws InterruptedException {
    while (true) {
      if (Thread.interrupted()) throw new InterruptedException();
      runDue(now.getAsLong());
      synchronized (this) {
        var deadline = next();
        if (deadline.isEmpty()) {
          wait();
          continue;
        }
        long current = now.getAsLong();
        if (deadline.getAsLong() <= current) continue;
        long delay;
        try {
          delay = Math.subtractExact(deadline.getAsLong(), current);
        } catch (ArithmeticException e) {
          delay = Long.MAX_VALUE;
        }
        wait(delay / 1_000_000, (int) (delay % 1_000_000));
      }
    }
  }
}
