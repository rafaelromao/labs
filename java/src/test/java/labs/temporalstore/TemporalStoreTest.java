package labs.temporalstore;

import static org.junit.jupiter.api.Assertions.*;

import java.util.ArrayList;
import java.util.concurrent.Executors;
import org.junit.jupiter.api.Test;

class TemporalStoreTest {
    @Test void predecessorAndReplacement() {
        var store = new TemporalStore<String, String>();
        store.put("a", 20, "twenty");
        store.put("a", 10, "ten");
        store.put("a", 20, "updated");
        assertTrue(store.get("a", 9).isEmpty());
        assertEquals("ten", store.get("a", 10).orElseThrow());
        assertEquals("ten", store.get("a", 19).orElseThrow());
        assertEquals("updated", store.get("a", 20).orElseThrow());
        assertEquals("updated", store.get("a", 21).orElseThrow());
        assertTrue(store.get("missing", 100).isEmpty());
    }

    @Test void concurrentWrites() throws Exception {
        var store = new TemporalStore<String, Integer>();
        try (var executor = Executors.newFixedThreadPool(4)) {
            var futures = new ArrayList<java.util.concurrent.Future<?>>();
            for (int i = 0; i < 100; i++) {
                int timestamp = i;
                futures.add(executor.submit(() -> store.put("a", timestamp, timestamp)));
            }
            for (var future : futures) future.get();
        }
        assertEquals(99, store.get("a", 99).orElseThrow());
    }
}
