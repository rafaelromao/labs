package labs.ttlcache;

import java.time.Duration;
import java.util.Optional;
import java.util.function.LongSupplier;

/** A bounded, linearizable LRU cache with a fixed time-to-live. */
public final class TtlCache<K, V> {
    public TtlCache(int capacity, Duration ttl) {
        this(capacity, ttl, System::nanoTime);
    }

    /** Injects monotonic nanoseconds for deterministic tests. */
    public TtlCache(int capacity, Duration ttl, LongSupplier nanoTime) {
        // TODO: validate and store the configuration.
    }

    public Optional<V> get(K key) {
        throw new UnsupportedOperationException("TODO: implement");
    }

    public void put(K key, V value) {
        throw new UnsupportedOperationException("TODO: implement");
    }

    public boolean delete(K key) {
        throw new UnsupportedOperationException("TODO: implement");
    }

    public int size() {
        throw new UnsupportedOperationException("TODO: implement");
    }
}
