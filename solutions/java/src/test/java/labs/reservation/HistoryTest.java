package labs.reservation;

import static org.junit.jupiter.api.Assertions.*;

import java.util.*;
import java.util.concurrent.*;
import java.util.concurrent.atomic.*;
import org.junit.jupiter.api.Test;

class HistoryTest {
  private record Call(String kind, long start, long finish, boolean accepted) {}

  private boolean legal(List<Call> h, boolean[] used, Set<String> state, int count) {
    if (count == h.size()) return true;
    for (int i = 0; i < h.size(); i++) {
      if (used[i]) continue;
      var c = h.get(i);
      boolean ready = true;
      for (int j = 0; j < h.size(); j++)
        if (!used[j] && h.get(j).finish() < c.start()) ready = false;
      if (!ready) continue;
      var next = new HashSet<>(state);
      boolean answer;
      if (c.kind().equals("cancel")) answer = next.remove("a");
      else {
        answer = next.contains(c.kind()) || next.isEmpty();
        if (answer) next.add(c.kind());
      }
      if (answer != c.accepted()) continue;
      used[i] = true;
      if (legal(h, used, next, count + 1)) return true;
      used[i] = false;
    }
    return false;
  }

  @Test
  void boundedHistories() throws Exception {
    try (var pool = Executors.newFixedThreadPool(3)) {
      for (int round = 0; round < 40; round++) {
        var c = new ReservationCalendar(1);
        var clock = new AtomicLong();
        var futures = new ArrayList<Future<Call>>();
        for (var kind : List.of("a", "b", "cancel"))
          futures.add(
              pool.submit(
                  () -> {
                    long start = clock.incrementAndGet();
                    boolean answer =
                        kind.equals("cancel")
                            ? c.cancel("a")
                            : c.book(kind, new ReservationCalendar.Interval(0, 10));
                    return new Call(kind, start, clock.incrementAndGet(), answer);
                  }));
        var history = new ArrayList<Call>();
        for (var f : futures) history.add(f.get(2, TimeUnit.SECONDS));
        assertTrue(legal(history, new boolean[3], Set.of(), 0));
      }
    }
    assertFalse(
        legal(
            List.of(new Call("a", 1, 4, true), new Call("b", 2, 3, true)),
            new boolean[2],
            Set.of(),
            0));
  }
}
