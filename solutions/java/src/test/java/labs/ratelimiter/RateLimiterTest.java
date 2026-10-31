package labs.ratelimiter;

import static org.junit.jupiter.api.Assertions.*;
import java.time.Duration;
import java.util.ArrayList;
import java.util.concurrent.Executors;
import java.util.concurrent.atomic.AtomicLong;
import org.junit.jupiter.api.Test;

class RateLimiterTest {
    @Test void burstRefillAndCardinality() {
        var now = new AtomicLong();
        var limiter = new RateLimiter(2, 2, 1, Duration.ofMinutes(1), now::get);
        assertEquals(RateLimiter.Result.ALLOWED, limiter.allow("a"));
        assertEquals(RateLimiter.Result.ALLOWED, limiter.allow("a"));
        assertEquals(RateLimiter.Result.LIMITED, limiter.allow("a"));
        now.set(500_000_000);
        assertEquals(RateLimiter.Result.ALLOWED, limiter.allow("a"));
        assertEquals(RateLimiter.Result.CLIENT_LIMIT, limiter.allow("b"));
        now.addAndGet(Duration.ofMinutes(1).toNanos());
        assertEquals(RateLimiter.Result.ALLOWED, limiter.allow("b"));
    }

    @Test void concurrentAdmission() throws Exception {
        var limiter = new RateLimiter(1, 5, 1, Duration.ofMinutes(1), () -> 0L);
        int allowed = 0;
        try (var executor = Executors.newFixedThreadPool(8)) {
            var futures = new ArrayList<java.util.concurrent.Future<RateLimiter.Result>>();
            for (int i = 0; i < 100; i++) futures.add(executor.submit(() -> limiter.allow("a")));
            for (var future : futures) if (future.get() == RateLimiter.Result.ALLOWED) allowed++;
        }
        assertEquals(5, allowed);
    }
}
