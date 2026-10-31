package labs.reservation;

import static org.junit.jupiter.api.Assertions.*;

import java.util.ArrayList;
import java.util.concurrent.Executors;
import org.junit.jupiter.api.Test;

class ReservationCalendarTest {
  @Test
  void weightedRescheduleRetainsOld() {
    var c = new ReservationCalendar(3);
    assertTrue(c.bookUnits("a", new ReservationCalendar.Interval(0, 10), 2));
    assertTrue(c.book("b", new ReservationCalendar.Interval(10, 20)));
    assertFalse(c.reschedule("a", new ReservationCalendar.Interval(10, 20), 3));
    assertFalse(c.bookUnits("c", new ReservationCalendar.Interval(0, 10), 2));
    assertTrue(c.reschedule("a", new ReservationCalendar.Interval(20, 30), 3));
  }

  @Test
  void boundariesAndLifecycle() {
    var c = new ReservationCalendar(1);
    assertTrue(c.book("a", new ReservationCalendar.Interval(0, 10)));
    assertTrue(c.book("b", new ReservationCalendar.Interval(10, 20)));
    assertFalse(c.book("c", new ReservationCalendar.Interval(5, 15)));
    assertTrue(c.book("a", new ReservationCalendar.Interval(0, 10)));
    assertThrows(
        IllegalArgumentException.class, () -> c.book("a", new ReservationCalendar.Interval(0, 11)));
    assertTrue(c.cancel("a"));
    assertFalse(c.cancel("a"));
    assertEquals(1, c.expire(20));
  }

  @Test
  void concurrentCapacity() throws Exception {
    var c = new ReservationCalendar(3);
    int accepted = 0;
    try (var executor = Executors.newFixedThreadPool(8)) {
      var futures = new ArrayList<java.util.concurrent.Future<Boolean>>();
      for (int i = 0; i < 100; i++) {
        String id = Integer.toString(i);
        futures.add(executor.submit(() -> c.book(id, new ReservationCalendar.Interval(0, 10))));
      }
      for (var future : futures) if (future.get()) accepted++;
    }
    assertEquals(3, accepted);
  }
}
