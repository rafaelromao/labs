package labs.reservation;

/** Atomic weighted half-open reservations with compatible one-unit booking. */
public final class ReservationCalendar {
    public record Interval(long start, long end) {}

    public ReservationCalendar(int capacity) {
        // TODO: validate and store the capacity.
    }

    public boolean book(String id, Interval interval) {
        return bookUnits(id, interval, 1);
    }

    public boolean bookUnits(String id, Interval interval, int units) {
        throw new UnsupportedOperationException("TODO: implement");
    }

    public boolean reschedule(String id, Interval interval, int units) {
        throw new UnsupportedOperationException("TODO: implement");
    }

    public boolean cancel(String id) {
        throw new UnsupportedOperationException("TODO: implement");
    }

    public int expire(long now) {
        throw new UnsupportedOperationException("TODO: implement");
    }
}
