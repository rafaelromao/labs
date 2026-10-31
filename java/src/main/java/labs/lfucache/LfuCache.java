package labs.lfucache;

import java.util.Optional;

/** Least-frequently-used eviction with least-recently-used ties; O(1) get and put. */
public final class LfuCache<K, V> {
    public LfuCache(int capacity) {
        // TODO: validate and store the capacity.
    }

    public Optional<V> get(K key) {
        throw new UnsupportedOperationException("TODO: implement");
    }

    public void put(K key, V value) {
        throw new UnsupportedOperationException("TODO: implement");
    }
}
