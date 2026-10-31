package labs.resilience;
import java.util.Objects;
import java.util.concurrent.Callable;
import java.util.function.*;

/** Each attempt counts as a breaker observation; one half-open probe at a time. */
public final class RetryBreaker{
    public static final class OpenCircuit extends Exception{}
    private record Permit(long generation,boolean probe){}
    @FunctionalInterface public interface Sleeper{void sleep(long nanos)throws InterruptedException;}
    private final int threshold;private final long cooldown;private final LongSupplier clock;
    private int state,failures;private long until,generation;
    public RetryBreaker(int threshold,long cooldownNanos,LongSupplier clock){if(threshold<=0||cooldownNanos<=0)throw new IllegalArgumentException("invalid configuration");this.threshold=threshold;this.cooldown=cooldownNanos;this.clock=Objects.requireNonNull(clock);}
    private synchronized Permit acquire()throws OpenCircuit{
        if(state==1&&clock.getAsLong()-until>=0){state=2;return new Permit(generation,true);}
        if(state!=0)throw new OpenCircuit();return new Permit(generation,false);
    }
    private synchronized void finish(Permit p,Exception failure){
        if(p.generation()!=generation)return;
        if(p.probe()){if(failure==null){state=0;failures=0;generation++;}else open();return;}
        if(state!=0)return;if(failure==null)failures=0;else if(++failures>=threshold)open();
    }
    private void open(){state=1;until=clock.getAsLong()+cooldown;generation++;}
    public <V>V execute(int attempts,Callable<V> task,Predicate<Exception> retryable,IntToLongFunction backoff,
            BooleanSupplier cancelled,Sleeper sleeper)throws Exception{
        if(attempts<=0)throw new IllegalArgumentException("positive attempts required");
        Objects.requireNonNull(task);Objects.requireNonNull(retryable);Objects.requireNonNull(backoff);Objects.requireNonNull(cancelled);Objects.requireNonNull(sleeper);
        for(int attempt=1;;attempt++){
            if(cancelled.getAsBoolean())throw new java.util.concurrent.CancellationException();
            Permit p=acquire();V value;
            try{value=task.call();}catch(Exception failure){
                finish(p,failure);
                if(failure instanceof InterruptedException){Thread.currentThread().interrupt();throw failure;}
                if(attempt==attempts||!retryable.test(failure))throw failure;
                long delay=backoff.applyAsLong(attempt);if(delay<0)throw new IllegalArgumentException("negative delay");sleeper.sleep(delay);continue;
            }catch(Error fatal){finish(p,new Exception("fatal task failure",fatal));throw fatal;}
            finish(p,null);return value;
        }
    }
}
