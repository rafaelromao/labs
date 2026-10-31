package labs.temporalstore;

import java.util.Optional;

/** In-memory time-indexed store. Values are returned without copying. */
public final class TemporalStore<K, V> {
    /** Replaces an existing value at the same timestamp. */
    public void put(K key, long timestamp, V value) {
        throw new UnsupportedOperationException("TODO: implement");
    }

    /** Returns the value at the greatest timestamp not greater than {@code at}. */
    public Optional<V> get(K key, long at) {
        throw new UnsupportedOperationException("TODO: implement");
    }
}
