package labs.dag;

import java.util.List;
import java.util.Map;
import java.util.Objects;
import java.util.concurrent.Callable;
import java.util.function.BooleanSupplier;

/** Validates a finite DAG, then executes bounded topological waves. */
public final class DependencyExecutor {
    public record Task<V>(List<String> dependsOn, Callable<V> run) {
        public Task {
            dependsOn = List.copyOf(dependsOn);
            Objects.requireNonNull(run);
        }
    }

    public record Outcome<V>(V value, Throwable failure, boolean skipped) {}

    private DependencyExecutor() {}

    public static <V> Map<String, Outcome<V>> execute(
            Map<String, Task<V>> input, int concurrency, int maxTasks, BooleanSupplier cancelled)
            throws InterruptedException {
        throw new UnsupportedOperationException("TODO: implement");
    }
}
