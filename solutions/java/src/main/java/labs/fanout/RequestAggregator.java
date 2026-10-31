package labs.fanout;

import java.util.*;
import java.util.concurrent.*;
import java.util.function.BooleanSupplier;

/** All-settled aggregation: failures do not cancel healthy siblings. */
public final class RequestAggregator {
    public record Outcome<V>(V value, Throwable failure) {}
    private RequestAggregator() {}
    public static <V> Map<String,Outcome<V>> run(Map<String,Callable<V>> tasks, int concurrency, int maxTasks,
            BooleanSupplier cancelled) throws InterruptedException {
        if (concurrency<=0 || maxTasks<=0 || tasks.size()>maxTasks) throw new IllegalArgumentException("invalid limits");
        var copied=new TreeMap<String,Callable<V>>();
        tasks.forEach((key,task)->copied.put(Objects.requireNonNull(key),Objects.requireNonNull(task)));
        var result=new TreeMap<String,Outcome<V>>();
        if (copied.isEmpty()) return Collections.unmodifiableMap(result);
        try (var executor=Executors.newFixedThreadPool(Math.min(concurrency,copied.size()))) {
            var futures=new TreeMap<String,Future<V>>();
            for (var entry:copied.entrySet()) futures.put(entry.getKey(),executor.submit(()-> {
                if (cancelled.getAsBoolean()) throw new CancellationException("cancelled before dispatch");
                return entry.getValue().call();
            }));
            for (var entry:futures.entrySet()) {
                try { result.put(entry.getKey(),new Outcome<>(entry.getValue().get(),null)); }
                catch (ExecutionException e) { result.put(entry.getKey(),new Outcome<>(null,e.getCause())); }
            }
        }
        return Collections.unmodifiableMap(result);
    }
}
