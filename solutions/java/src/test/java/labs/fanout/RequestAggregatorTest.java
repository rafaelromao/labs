package labs.fanout;

import static org.junit.jupiter.api.Assertions.*;
import java.util.Map;
import java.util.concurrent.Callable;
import java.util.concurrent.CancellationException;
import org.junit.jupiter.api.Test;

class RequestAggregatorTest {
    @Test void partialFailure() throws Exception {
        Map<String,Callable<Integer>> tasks=Map.of("a",()->7,"b",()->{ throw new IllegalStateException("unavailable"); });
        var result=RequestAggregator.run(tasks,1,2,()->false);
        assertEquals(7,result.get("a").value());
        assertInstanceOf(IllegalStateException.class,result.get("b").failure());
    }
    @Test void cancellationAndInputBound() throws Exception {
        Map<String,Callable<Integer>> tasks=Map.of("a",()->{ fail("task started"); return 0; });
        var result=RequestAggregator.run(tasks,1,1,()->true);
        assertInstanceOf(CancellationException.class,result.get("a").failure());
        assertThrows(IllegalArgumentException.class,()->RequestAggregator.run(tasks,0,1,()->false));
    }
}
