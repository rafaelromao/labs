package labs.fanout;

import java.util.Map;
import java.util.concurrent.Callable;
import java.util.function.BooleanSupplier;

/** All-settled aggregation: failures do not cancel healthy siblings. */
public final class RequestAggregator {
    public record Outcome<V>(V value, Throwable failure) {}

    private RequestAggregator() {}

    public static <V> Map<String, Outcome<V>> run(
            Map<String, Callable<V>> tasks, int concurrency, int maxTasks, BooleanSupplier cancelled)
            throws InterruptedException {
        throw new UnsupportedOperationException("TODO: implement");
    }
}
