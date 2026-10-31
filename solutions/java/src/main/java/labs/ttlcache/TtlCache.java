package labs.ttlcache;

import java.time.Duration;
import java.util.LinkedHashMap;
import java.util.Objects;
import java.util.Optional;
import java.util.function.LongSupplier;

/** A bounded, linearizable LRU cache with a fixed time-to-live. */
public final class TtlCache<K, V> {
    private record Entry<V>(V value, long expiresAtNanos) {}

    private final int capacity;
    private final long ttlNanos;
    private final LongSupplier nanoTime;
    private final LinkedHashMap<K, Entry<V>> entries;

    public TtlCache(int capacity, Duration ttl) {
        this(capacity, ttl, System::nanoTime);
    }

    /** Injects monotonic nanoseconds for deterministic tests. */
    public TtlCache(int capacity, Duration ttl, LongSupplier nanoTime) {
        Objects.requireNonNull(ttl, "ttl");
        this.nanoTime = Objects.requireNonNull(nanoTime, "nanoTime");
        if (capacity <= 0 || ttl.isZero() || ttl.isNegative()) {
            throw new IllegalArgumentException("capacity and ttl must be positive");
        }
        try {
            this.ttlNanos = ttl.toNanos();
        } catch (ArithmeticException e) {
            throw new IllegalArgumentException("ttl is too large", e);
        }
        this.capacity = capacity;
        this.entries = new LinkedHashMap<>(capacity, 0.75f, true);
    }

    public synchronized Optional<V> get(K key) {
        Objects.requireNonNull(key, "key");
        Entry<V> entry = entries.get(key);
        if (entry == null) {
            return Optional.empty();
        }
        if (expired(nanoTime.getAsLong(), entry.expiresAtNanos())) {
            entries.remove(key);
            return Optional.empty();
        }
        return Optional.of(entry.value());
    }

    public synchronized void put(K key, V value) {
        Objects.requireNonNull(key, "key");
        Objects.requireNonNull(value, "value");
        long now = nanoTime.getAsLong();
        removeExpired(now);
        entries.put(key, new Entry<>(value, now + ttlNanos));
        if (entries.size() > capacity) {
            K leastRecentlyUsed = entries.keySet().iterator().next();
            entries.remove(leastRecentlyUsed);
        }
    }

    public synchronized boolean delete(K key) {
        Objects.requireNonNull(key, "key");
        return entries.remove(key) != null;
    }

    public synchronized int size() {
        removeExpired(nanoTime.getAsLong());
        return entries.size();
    }

    private void removeExpired(long now) {
        entries.entrySet().removeIf(entry -> expired(now, entry.getValue().expiresAtNanos()));
    }

    // Wrap-safe while the actual difference is less than 2^63 nanoseconds.
    private static boolean expired(long now, long deadline) {
        return now - deadline >= 0;
    }
}
