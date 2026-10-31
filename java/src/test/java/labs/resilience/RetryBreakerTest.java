package labs.resilience;
import static org.junit.jupiter.api.Assertions.*;
import java.util.concurrent.atomic.*;
import org.junit.jupiter.api.Test;
class RetryBreakerTest{
    @Test void boundedRetryAndRecovery()throws Exception{
        var now=new AtomicLong();var calls=new AtomicInteger();var b=new RetryBreaker(2,100,now::get);
        assertThrows(RetryBreaker.OpenCircuit.class,()->b.execute(5,()->{calls.incrementAndGet();throw new IllegalStateException();},e->true,n->0,()->false,n->{}));
        assertEquals(2,calls.get());now.set(100);
        assertEquals(7,b.execute(1,()->7,e->false,n->0,()->false,n->{}));
    }
    @Test void cancellationSkipsAttempt(){
        var b=new RetryBreaker(1,100,()->0L);
        assertThrows(java.util.concurrent.CancellationException.class,()->b.execute(1,()->{fail("ran");return 0;},e->false,n->0,()->true,n->{}));
    }
}
