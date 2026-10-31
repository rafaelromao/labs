package labs.spatial;

import java.util.*;

/** Revision updates reclassify retained planar observations atomically. */
public final class VersionedFence {
  public record Transition(String id, boolean inside, long revision, long sequence) {}

  private record State(long sequence, Polygon.Point point, boolean inside) {}

  private Polygon shape;
  private long revision;
  private final int capacity;
  private final Map<String, State> objects = new HashMap<>();

  public VersionedFence(Polygon shape, long revision, int capacity) {
    this.shape = Objects.requireNonNull(shape);
    if (revision < 0 || capacity <= 0) throw new IllegalArgumentException();
    this.revision = revision;
    this.capacity = capacity;
  }

  public synchronized Optional<Transition> observe(String id, long sequence, Polygon.Point point) {
    Objects.requireNonNull(id);
    if (sequence < 0 || !Polygon.valid(point)) throw new IllegalArgumentException();
    var old = objects.get(id);
    if (old != null && sequence <= old.sequence()) return Optional.empty();
    if (old == null && objects.size() == capacity) throw new IllegalStateException("capacity");
    boolean inside = shape.covers(point);
    objects.put(id, new State(sequence, point, inside));
    return inside == (old != null && old.inside())
        ? Optional.empty()
        : Optional.of(new Transition(id, inside, revision, sequence));
  }

  public synchronized List<Transition> replace(Polygon shape, long revision) {
    Objects.requireNonNull(shape);
    if (revision <= this.revision) throw new IllegalArgumentException("stale revision");
    var changes = new ArrayList<Transition>();
    for (var entry : objects.entrySet()) {
      var old = entry.getValue();
      boolean inside = shape.covers(old.point());
      if (inside != old.inside())
        changes.add(new Transition(entry.getKey(), inside, revision, old.sequence()));
      entry.setValue(new State(old.sequence(), old.point(), inside));
    }
    this.shape = shape;
    this.revision = revision;
    changes.sort(Comparator.comparing(Transition::id));
    return List.copyOf(changes);
  }
}
