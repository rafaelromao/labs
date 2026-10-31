package labs.windowmetrics;

import static org.junit.jupiter.api.Assertions.*;
import java.util.List;
import org.junit.jupiter.api.Test;

class EventWindowTest {
    @Test void boundsCapacityAndLateData() {
        var window=new EventWindow(10,2);
        window.add(new EventWindow.Event(0,3)); window.add(new EventWindow.Event(5,9));
        assertThrows(IllegalStateException.class,()->window.add(new EventWindow.Event(5,1)));
        assertEquals(List.of(9L,3L),window.query(9,2).top());
        assertEquals(1,window.query(10,2).count());
        assertThrows(IllegalArgumentException.class,()->window.add(new EventWindow.Event(9,7)));
        window.add(new EventWindow.Event(10,7));
        assertTrue(window.query(20,2).maximum().isEmpty());
    }
}
