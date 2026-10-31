package labs.singleflight;

import java.util.HashMap;
import java.util.Map;
import java.util.Objects;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.Executor;
import java.util.concurrent.RejectedExecutionException;
import java.util.function.Function;

/** In-flight results only. The caller owns the executor and loader lifecycle. */
public final class RequestCoalescer<V> {
  private final int capacity;
  private final Executor executor;
  private final Function<String, V> loader;
  private final Map<String, CompletableFuture<V>> flights = new HashMap<>();
  private boolean closed;
  private int active;

  public RequestCoalescer(int capacity, Executor executor, Function<String, V> loader) {
    if (capacity <= 0) throw new IllegalArgumentException("positive capacity required");
    this.capacity = capacity;
    this.executor = Objects.requireNonNull(executor);
    this.loader = Objects.requireNonNull(loader);
  }

  /** Returns an independent future whose cancellation cannot cancel the load. */
  public CompletableFuture<V> request(String key) {
    Objects.requireNonNull(key);
    CompletableFuture<V> shared;
    boolean start = false;
    synchronized (this) {
      if (closed) throw new RejectedExecutionException("admission closed");
      shared = flights.get(key);
      if (shared == null) {
        if (flights.size() == capacity)
          throw new RejectedExecutionException("distinct-key capacity reached");
        shared = new CompletableFuture<>();
        flights.put(key, shared);
        start = true;
        active++;
      }
    }
    if (start) {
      CompletableFuture<V> flight = shared;
      try {
        executor.execute(
            () -> {
              V value = null;
              Throwable failure = null;
              try {
                value = loader.apply(key);
              } catch (Throwable t) {
                failure = t;
              }
              synchronized (this) {
                flights.remove(key, flight);
              }
              if (failure == null) flight.complete(value);
              else flight.completeExceptionally(failure);
              synchronized (this) {
                active--;
                notifyAll();
              }
            });
      } catch (RuntimeException rejected) {
        synchronized (this) {
          flights.remove(key, flight);
        }
        flight.completeExceptionally(rejected);
        synchronized (this) {
          active--;
          notifyAll();
        }
      }
    }
    return shared.copy();
  }

  /**
   * Close admission and join loads; owner must cancel blocked dependencies. Never call this from a
   * loader or its synchronous completion callback.
   */
  public synchronized void closeAndWait() throws InterruptedException {
    closed = true;
    while (active != 0) wait();
  }
}
