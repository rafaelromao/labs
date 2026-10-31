package labs.workerpool;

import static org.junit.jupiter.api.Assertions.*;

import java.util.concurrent.CountDownLatch;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.RejectedExecutionException;
import java.util.concurrent.TimeUnit;
import org.junit.jupiter.api.Test;

class WorkerPoolTest {
  @Test
  void abortAccountsForPending() throws Exception {
    var pool = new WorkerPool(1, 1);
    var started = new CountDownLatch(1);
    try {
      var first =
          pool.submit(
              () -> {
                started.countDown();
                new CountDownLatch(1).await();
                return 1;
              });
      assertTrue(started.await(1, TimeUnit.SECONDS));
      var second = pool.submit(() -> 2);
      assertEquals(1, pool.abort());
      assertTrue(second.isCancelled());
      assertThrows(ExecutionException.class, () -> first.get(1, TimeUnit.SECONDS));
      assertTrue(pool.awaitTermination(1, TimeUnit.SECONDS));
    } finally {
      pool.abort();
    }
  }

  @Test
  void saturationAndDrain() throws Exception {
    var pool = new WorkerPool(1, 1);
    var started = new CountDownLatch(1);
    var release = new CountDownLatch(1);
    try {
      var first =
          pool.submit(
              () -> {
                started.countDown();
                release.await();
                return 1;
              });
      assertTrue(started.await(1, TimeUnit.SECONDS));
      var second = pool.submit(() -> 2);
      assertThrows(RejectedExecutionException.class, () -> pool.submit(() -> 3));
      pool.close();
      pool.close();
      assertThrows(RejectedExecutionException.class, () -> pool.submit(() -> 4));
      release.countDown();
      assertEquals(1, first.get(1, TimeUnit.SECONDS));
      assertEquals(2, second.get(1, TimeUnit.SECONDS));
      assertTrue(pool.awaitTermination(1, TimeUnit.SECONDS));
    } finally {
      release.countDown();
      pool.close();
    }
  }

  @Test
  void failuresAreObserved() throws Exception {
    var pool = new WorkerPool(1, 1);
    try {
      var result =
          pool.submit(
              () -> {
                throw new IllegalStateException("boom");
              });
      assertInstanceOf(
          IllegalStateException.class,
          assertThrows(ExecutionException.class, () -> result.get(1, TimeUnit.SECONDS)).getCause());
      var next = pool.submit(() -> 7);
      assertEquals(7, next.get(1, TimeUnit.SECONDS));
    } finally {
      pool.close();
    }
    assertTrue(pool.awaitTermination(1, TimeUnit.SECONDS));
  }
}
