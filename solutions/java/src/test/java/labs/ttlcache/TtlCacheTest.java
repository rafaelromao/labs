package labs.ttlcache;

import static org.junit.jupiter.api.Assertions.*;

import java.time.Duration;
import java.util.ArrayList;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.Executors;
import java.util.concurrent.Future;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicLong;
import org.junit.jupiter.api.Test;

class TtlCacheTest {
    @Test
    void expiresAtInclusiveDeadline() {
        AtomicLong now = new AtomicLong(10);
        TtlCache<String, Integer> cache = new TtlCache<>(2, Duration.ofNanos(5), now::get);
        cache.put("k", 7);
        now.set(14);
        assertEquals(7, cache.get("k").orElseThrow());
        now.set(15);
        assertTrue(cache.get("k").isEmpty());
    }

    @Test
    void getPromotesAndEvictsLeastRecentlyUsed() {
        AtomicLong now = new AtomicLong(0);
        TtlCache<String, Integer> cache = new TtlCache<>(2, Duration.ofHours(1), now::get);
        cache.put("a", 1);
        cache.put("b", 2);
        assertEquals(1, cache.get("a").orElseThrow());
        cache.put("c", 3);
        assertTrue(cache.get("b").isEmpty());
        assertEquals(2, cache.size());
    }

    @Test
    void updateRefreshesValueAndDeadline() {
        AtomicLong now = new AtomicLong(0);
        TtlCache<String, Integer> cache = new TtlCache<>(2, Duration.ofNanos(10), now::get);
        cache.put("a", 1);
        now.set(5);
        cache.put("a", 9);
        now.set(10);
        assertEquals(9, cache.get("a").orElseThrow());
        now.set(15);
        assertTrue(cache.get("a").isEmpty());
    }

    @Test
    void putSweepsExpiredEntriesBeforeEviction() {
        AtomicLong now = new AtomicLong(0);
        TtlCache<String, Integer> cache = new TtlCache<>(2, Duration.ofNanos(10), now::get);
        cache.put("a", 1);
        cache.put("b", 2);
        now.set(10);
        cache.put("c", 3);
        cache.put("d", 4);
        assertEquals(2, cache.size());
        assertEquals(3, cache.get("c").orElseThrow());
    }

    @Test
    void rejectsInvalidConfiguration() {
        assertThrows(IllegalArgumentException.class,
                () -> new TtlCache<String, Integer>(0, Duration.ofSeconds(1)));
        assertThrows(IllegalArgumentException.class,
                () -> new TtlCache<String, Integer>(1, Duration.ZERO));
    }

    @Test
    void rejectsNullsAndDeletesExpiredResident() {
        AtomicLong now = new AtomicLong(0);
        TtlCache<String, Integer> cache = new TtlCache<>(2, Duration.ofNanos(10), now::get);
        assertThrows(NullPointerException.class, () -> cache.put(null, 1));
        assertThrows(NullPointerException.class, () -> cache.put("a", null));
        cache.put("a", 1);
        now.set(10);
        assertTrue(cache.delete("a"));
        assertFalse(cache.delete("a"));
    }

    @Test
    void concurrentOperationsStayWithinCapacity() throws Exception {
        TtlCache<Integer, Integer> cache = new TtlCache<>(64, Duration.ofHours(1));
        int workers = 8;
        int iterations = 500;
        CountDownLatch start = new CountDownLatch(1);
        try (var executor = Executors.newFixedThreadPool(workers)) {
            var futures = new ArrayList<Future<?>>();
            for (int worker = 0; worker < workers; worker++) {
                int id = worker;
                futures.add(executor.submit(() -> {
                    start.await();
                    for (int i = 0; i < iterations; i++) {
                        int key = (id + i) % 100;
                        cache.put(key, i);
                        cache.get(key);
                        cache.size();
                    }
                    return null;
                }));
            }
            start.countDown();
            executor.shutdown();
            assertTrue(executor.awaitTermination(5, TimeUnit.SECONDS));
            for (Future<?> future : futures) {
                future.get();
            }
        }
        assertTrue(cache.size() <= 64);
    }
}
