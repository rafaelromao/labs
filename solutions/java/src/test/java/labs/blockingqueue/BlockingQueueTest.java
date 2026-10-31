package labs.blockingqueue;

import static org.junit.jupiter.api.Assertions.*;

import java.util.concurrent.Executors;
import java.util.concurrent.TimeUnit;
import org.junit.jupiter.api.Test;

class BlockingQueueTest {
    @Test void drainsInOrderAfterClose() throws Exception {
        var queue = new BlockingQueue<Integer>(2);
        queue.put(1); queue.put(2); queue.close(); queue.close();
        assertEquals(1, queue.take().orElseThrow());
        assertEquals(2, queue.take().orElseThrow());
        assertTrue(queue.take().isEmpty());
        assertThrows(IllegalStateException.class, () -> queue.put(3));
    }

    @Test void closeReleasesProducer() throws Exception {
        var queue = new BlockingQueue<Integer>(1);
        queue.put(1);
        try (var executor = Executors.newSingleThreadExecutor()) {
            var result = executor.submit(() -> assertThrows(IllegalStateException.class, () -> queue.put(2)));
            queue.close();
            result.get(1, TimeUnit.SECONDS);
        }
    }

    @Test void rejectsInvalidCapacityAndNull() {
        assertThrows(IllegalArgumentException.class, () -> new BlockingQueue<Integer>(0));
        var queue = new BlockingQueue<Integer>(1);
        assertThrows(NullPointerException.class, () -> queue.put(null));
    }
}
