package labs.pubsub;

import static org.junit.jupiter.api.Assertions.*;
import org.junit.jupiter.api.Test;

class BrokerTest {
    @Test void slowConsumerAndDrain() {
        var b=new Broker<Integer>(2,1); var a=b.subscribe(); var c=b.subscribe();
        assertThrows(IllegalStateException.class,b::subscribe);
        assertEquals(new Broker.Delivery(2,0),b.publish(1));
        assertEquals(1,a.poll().orElseThrow());
        assertEquals(new Broker.Delivery(1,1),b.publish(2));
        c.unsubscribe(); c.unsubscribe(); assertTrue(c.isClosed());
        assertEquals(1,c.poll().orElseThrow()); assertTrue(c.poll().isEmpty());
        b.close(); b.close(); assertTrue(a.isClosed());
        assertEquals(2,a.poll().orElseThrow()); assertTrue(a.poll().isEmpty());
        assertThrows(IllegalStateException.class,()->b.publish(3));
    }
}
