package labs.ratelimiter;

import java.time.Duration;
import java.util.function.LongSupplier;

/** Local per-client token buckets with bounded client cardinality. */
public final class RateLimiter {
    public enum Result { ALLOWED, LIMITED, CLIENT_LIMIT }

    public RateLimiter(double rate, int burst, int maxClients, Duration idle, LongSupplier now) {
        // TODO: validate and store the configuration.
    }

    public Result allow(String client) {
        throw new UnsupportedOperationException("TODO: implement");
    }
}
