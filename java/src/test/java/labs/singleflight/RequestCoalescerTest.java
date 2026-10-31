package labs.singleflight;

import static org.junit.jupiter.api.Assertions.*;

import java.util.concurrent.CountDownLatch;
import java.util.concurrent.Executors;
import java.util.concurrent.RejectedExecutionException;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicInteger;
import org.junit.jupiter.api.Test;

class RequestCoalescerTest {
  @Test
  void closeJoinsAndRejects() throws Exception {
    var g = new RequestCoalescer<Integer>(1, Runnable::run, key -> 7);
    assertEquals(7, g.request("a").join());
    g.closeAndWait();
    assertThrows(RejectedExecutionException.class, () -> g.request("b"));
  }

  @Test
  void sharesAndIsolatesCancellation() throws Exception {
    var started = new CountDownLatch(1);
    var release = new CountDownLatch(1);
    var calls = new AtomicInteger();
    try (var executor = Executors.newFixedThreadPool(2)) {
      try {
        var group =
            new RequestCoalescer<Integer>(
                1,
                executor,
                key -> {
                  calls.incrementAndGet();
                  started.countDown();
                  try {
                    release.await();
                  } catch (InterruptedException e) {
                    Thread.currentThread().interrupt();
                    throw new RuntimeException(e);
                  }
                  return 7;
                });
        var first = group.request("a");
        assertTrue(started.await(1, TimeUnit.SECONDS));
        var second = group.request("a");
        assertThrows(RejectedExecutionException.class, () -> group.request("b"));
        first.cancel(false);
        release.countDown();
        assertEquals(7, second.get(1, TimeUnit.SECONDS));
        assertEquals(1, calls.get());
      } finally {
        release.countDown();
      }
    }
  }

  @Test
  void failureAndExecutorRejectionReleaseCapacity() {
    var calls = new AtomicInteger();
    var group =
        new RequestCoalescer<Integer>(
            1,
            Runnable::run,
            key -> {
              if (calls.incrementAndGet() == 1) throw new IllegalStateException("failed");
              return 9;
            });
    assertTrue(group.request("a").isCompletedExceptionally());
    assertEquals(9, group.request("a").join());
    var rejecting =
        new RequestCoalescer<Integer>(
            1,
            task -> {
              throw new RejectedExecutionException();
            },
            key -> 1);
    assertTrue(rejecting.request("a").isCompletedExceptionally());
    assertTrue(rejecting.request("b").isCompletedExceptionally());
  }
}
