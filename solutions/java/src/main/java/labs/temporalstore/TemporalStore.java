package labs.temporalstore;

import java.util.HashMap;
import java.util.Map;
import java.util.NavigableMap;
import java.util.Objects;
import java.util.Optional;
import java.util.TreeMap;

/** In-memory time-indexed store. Values are returned without copying. */
public final class TemporalStore<K, V> {
    private final Map<K, NavigableMap<Long, V>> values = new HashMap<>();

    /** Replaces an existing value at the same timestamp. */
    public synchronized void put(K key, long timestamp, V value) {
        Objects.requireNonNull(key, "key");
        Objects.requireNonNull(value, "value");
        values.computeIfAbsent(key, ignored -> new TreeMap<>()).put(timestamp, value);
    }

    /** Returns the value at the greatest timestamp not greater than at. */
    public synchronized Optional<V> get(K key, long at) {
        Objects.requireNonNull(key, "key");
        NavigableMap<Long, V> history = values.get(key);
        if (history == null) return Optional.empty();
        Map.Entry<Long, V> entry = history.floorEntry(at);
        return entry == null ? Optional.empty() : Optional.of(entry.getValue());
    }
}
