package labs.reservation;

import java.util.*;

/** Atomic weighted half-open reservations with compatible one-unit booking. */
public final class ReservationCalendar {
  public record Interval(long start, long end) {}

  private record Booking(Interval interval, int units) {}

  private record Endpoint(long time, int units, boolean start) {}

  private final int capacity;
  private final Map<String, Booking> entries = new HashMap<>();

  public ReservationCalendar(int capacity) {
    if (capacity <= 0) throw new IllegalArgumentException();
    this.capacity = capacity;
  }

  public boolean book(String id, Interval interval) {
    return bookUnits(id, interval, 1);
  }

  public boolean bookUnits(String id, Interval interval, int units) {
    return change(id, interval, units, false);
  }

  public boolean reschedule(String id, Interval interval, int units) {
    return change(id, interval, units, true);
  }

  private synchronized boolean change(String id, Interval interval, int units, boolean replace) {
    Objects.requireNonNull(id);
    Objects.requireNonNull(interval);
    if (id.isEmpty() || interval.start() >= interval.end() || units <= 0 || units > capacity)
      throw new IllegalArgumentException();
    var candidate = new Booking(interval, units);
    var old = entries.get(id);
    if (replace && old == null) throw new IllegalArgumentException("missing booking");
    if (!replace && old != null) {
      if (!old.equals(candidate)) throw new IllegalArgumentException("ID conflict");
      return true;
    }
    var events = new ArrayList<Endpoint>();
    for (var entry : entries.entrySet())
      if (!entry.getKey().equals(id)) add(events, entry.getValue());
    add(events, candidate);
    events.sort(Comparator.comparingLong(Endpoint::time).thenComparing(Endpoint::start));
    int used = 0;
    for (var e : events) {
      if (e.start()) {
        if (e.units() > capacity - used) return false;
        used += e.units();
      } else used -= e.units();
    }
    entries.put(id, candidate);
    return true;
  }

  private static void add(List<Endpoint> events, Booking b) {
    events.add(new Endpoint(b.interval().start(), b.units(), true));
    events.add(new Endpoint(b.interval().end(), b.units(), false));
  }

  public synchronized boolean cancel(String id) {
    return entries.remove(Objects.requireNonNull(id)) != null;
  }

  public synchronized int expire(long now) {
    int before = entries.size();
    entries.entrySet().removeIf(e -> e.getValue().interval().end() <= now);
    return before - entries.size();
  }
}
