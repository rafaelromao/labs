package labs.benchmarks;

import labs.algorithms.Structures;
import labs.nearby.NearbyIndex;
import labs.ttlcache.TtlCache;
import labs.workerpool.WorkerPool;
import java.time.Duration;
import java.util.concurrent.*;
import org.openjdk.jmh.annotations.*;

@State(Scope.Thread)
@BenchmarkMode(Mode.AverageTime)
@OutputTimeUnit(TimeUnit.MICROSECONDS)
public class ComponentBenchmarks {
  private final long[] values = new long[10000];
  private NearbyIndex index;
  private TtlCache<String, Integer> cache;
  private WorkerPool pool;
  private CountDownLatch release;

  @Setup
  public void setup() throws Exception {
    for (int i = 0; i < values.length; i++) values[i] = (i * 7919L) % 10007;
    index = new NearbyIndex();
    for (int i = 0; i < 10000; i++)
      index.put("" + i, new NearbyIndex.Point(i % 180 - 90, i % 360 - 180));
    cache = new TtlCache<>(1000, Duration.ofHours(1), System::nanoTime);
    cache.put("key", 7);
    pool = new WorkerPool(1, 1);
    release = new CountDownLatch(1);
    var started = new CountDownLatch(1);
    pool.submit(
        () -> {
          started.countDown();
          release.await();
          return 1;
        });
    if (!started.await(2, TimeUnit.SECONDS))
      throw new IllegalStateException("worker did not start");
    pool.submit(() -> 2);
  }

  @TearDown
  public void teardown() throws Exception {
    release.countDown();
    pool.close();
    if (!pool.awaitTermination(2, TimeUnit.SECONDS)) throw new IllegalStateException("pool leak");
  }

  @Benchmark
  public Object topK() {
    return Structures.topK(values, 20);
  }

  @Benchmark
  public Object spatialRadius() {
    return index.radius(new NearbyIndex.Point(0, 0), 10000);
  }

  @Benchmark
  public Object ttlRead() {
    return cache.get("key");
  }

  @Benchmark
  public boolean workerSaturation() {
    try {
      pool.submit(() -> 3);
      throw new IllegalStateException("unexpected admission");
    } catch (RejectedExecutionException expected) {
      return true;
    }
  }
}
