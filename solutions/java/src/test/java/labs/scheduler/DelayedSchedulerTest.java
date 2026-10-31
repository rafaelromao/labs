package labs.scheduler;

import static org.junit.jupiter.api.Assertions.*;

import java.util.ArrayList;
import java.util.List;
import org.junit.jupiter.api.Test;

class DelayedSchedulerTest {
  @Test
  void earlierDeadlineWakesAndOwnerJoins() throws Exception {
    var s = new DelayedScheduler(2);
    long origin = System.nanoTime();
    var ran = new java.util.concurrent.CountDownLatch(1);
    s.schedule(
        "later", java.util.concurrent.TimeUnit.HOURS.toNanos(1), () -> fail("late task ran"));
    var failure = new java.util.concurrent.atomic.AtomicReference<Throwable>();
    var owner =
        new Thread(
            () -> {
              try {
                s.run(() -> System.nanoTime() - origin);
              } catch (InterruptedException expected) {
              } catch (Throwable e) {
                failure.set(e);
              }
            });
    owner.start();
    try {
      s.schedule("earlier", 0, ran::countDown);
      assertTrue(ran.await(1, java.util.concurrent.TimeUnit.SECONDS));
    } finally {
      owner.interrupt();
      owner.join(1000);
    }
    assertFalse(owner.isAlive());
    assertNull(failure.get());
  }

  @Test
  void dueRescheduleAndReentrant() {
    var s = new DelayedScheduler(2);
    var order = new ArrayList<Integer>();
    s.schedule("a", 10, () -> order.add(1));
    s.schedule("b", 20, () -> order.add(2));
    s.schedule("a", 30, () -> order.add(3));
    assertEquals(20, s.next().orElseThrow());
    assertEquals(0, s.runDue(19));
    assertTrue(s.cancel("b"));
    assertFalse(s.cancel("b"));
    s.schedule(
        "b",
        30,
        () -> {
          order.add(4);
          s.schedule("c", 30, () -> order.add(5));
        });
    assertEquals(2, s.runDue(30));
    assertEquals(List.of(3, 4), order);
    assertEquals(1, s.runDue(30));
    assertEquals(List.of(3, 4, 5), order);
  }
}
