package labs.singleflight;

import java.util.concurrent.CompletableFuture;
import java.util.concurrent.Executor;
import java.util.function.Function;

/** In-flight results only. The caller owns the executor and loader lifecycle. */
public final class RequestCoalescer<V> {
    public RequestCoalescer(int capacity, Executor executor, Function<String, V> loader) {
        // TODO: validate and store the configuration.
    }

    /** Returns an independent future whose cancellation cannot cancel the load. */
    public CompletableFuture<V> request(String key) {
        throw new UnsupportedOperationException("TODO: implement");
    }

    /** Closes admission and waits for running loads. Never call it from a loader. */
    public void closeAndWait() throws InterruptedException {
        throw new UnsupportedOperationException("TODO: implement");
    }
}
