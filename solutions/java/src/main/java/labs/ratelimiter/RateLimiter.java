package labs.ratelimiter;

import java.time.Duration;
import java.util.HashMap;
import java.util.Map;
import java.util.Objects;
import java.util.function.LongSupplier;

/** Local per-client token buckets with bounded client cardinality. */
public final class RateLimiter {
    public enum Result { ALLOWED, LIMITED, CLIENT_LIMIT }
    private static final class Bucket {
        double tokens;
        long last;
        Bucket(double tokens, long last) { this.tokens = tokens; this.last = last; }
    }
    private final double rate;
    private final int burst;
    private final int maxClients;
    private final long idleNanos;
    private final LongSupplier now;
    private final Map<String, Bucket> clients = new HashMap<>();

    public RateLimiter(double rate, int burst, int maxClients, Duration idle, LongSupplier now) {
        Objects.requireNonNull(idle);
        if (!Double.isFinite(rate) || rate <= 0 || burst <= 0 || maxClients <= 0 || idle.isNegative() || idle.isZero())
            throw new IllegalArgumentException("invalid configuration");
        this.idleNanos = idle.toNanos();
        this.rate = rate; this.burst = burst; this.maxClients = maxClients;
        this.now = Objects.requireNonNull(now);
    }

    public synchronized Result allow(String client) {
        Objects.requireNonNull(client);
        long current = now.getAsLong();
        clients.entrySet().removeIf(e -> current - e.getValue().last >= idleNanos);
        Bucket bucket = clients.get(client);
        if (bucket == null) {
            if (clients.size() == maxClients) return Result.CLIENT_LIMIT;
            bucket = new Bucket(burst, current);
            clients.put(client, bucket);
        }
        long elapsed = current - bucket.last;
        if (elapsed < 0) throw new IllegalStateException("clock moved backwards or exceeded comparison range");
        bucket.tokens = Math.min(burst, bucket.tokens + (elapsed / 1_000_000_000.0) * rate);
        bucket.last = current;
        if (bucket.tokens < 1) return Result.LIMITED;
        bucket.tokens--;
        return Result.ALLOWED;
    }
}
