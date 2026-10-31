package labs.resilience;

import java.util.concurrent.Callable;
import java.util.function.BooleanSupplier;
import java.util.function.IntToLongFunction;
import java.util.function.LongSupplier;
import java.util.function.Predicate;

/** Each attempt counts as a breaker observation; one half-open probe at a time. */
public final class RetryBreaker {
    public static final class OpenCircuit extends Exception {}

    @FunctionalInterface
    public interface Sleeper {
        void sleep(long nanos) throws InterruptedException;
    }

    public RetryBreaker(int threshold, long cooldownNanos, LongSupplier clock) {
        // TODO: validate and store the configuration.
    }

    public <V> V execute(int attempts, Callable<V> task, Predicate<Exception> retryable,
            IntToLongFunction backoff, BooleanSupplier cancelled, Sleeper sleeper) throws Exception {
        throw new UnsupportedOperationException("TODO: implement");
    }
}
