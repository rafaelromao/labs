package labs.lfucache;

import static org.junit.jupiter.api.Assertions.*;
import org.junit.jupiter.api.Test;

class LfuCacheTest {
    @Test void recencyBreaksFrequencyTies(){
        var cache=new LfuCache<String,Integer>(2);
        cache.put("a",1);cache.put("b",2);cache.get("a");cache.get("b");cache.put("c",3);
        assertTrue(cache.get("a").isEmpty());assertEquals(2,cache.get("b").orElseThrow());
    }
    @Test void frequencyAndRecency() {
        var cache = new LfuCache<String, Integer>(2);
        cache.put("a", 1); cache.put("b", 2); cache.get("a"); cache.put("c", 3);
        assertTrue(cache.get("b").isEmpty());
        cache.put("d", 4);
        assertTrue(cache.get("c").isEmpty());
        assertEquals(1, cache.get("a").orElseThrow());
    }
    @Test void updateCountsAsAccess() {
        var cache = new LfuCache<String, Integer>(2);
        cache.put("a", 1); cache.put("b", 2); cache.put("a", 9); cache.put("c", 3);
        assertTrue(cache.get("b").isEmpty());
        assertEquals(9, cache.get("a").orElseThrow());
    }
}
