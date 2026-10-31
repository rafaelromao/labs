package labs.outbox;
import static org.junit.jupiter.api.Assertions.*;
import java.nio.file.Path;
import java.sql.*;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.io.TempDir;
class OutboxStoreTest{
    @TempDir Path directory;
    @Test void lostAckReplayAndReopen()throws Exception{
        var producer=new OutboxStore(directory.resolve("producer.db"));var file=directory.resolve("consumer.db");var consumer=new OutboxStore(file);
        var event=new OutboxStore.Event("event-1",7);assertTrue(producer.produce(event));assertFalse(producer.produce(event));
        assertTrue(consumer.consume(producer.pending(10).getFirst()));
        var reopened=new OutboxStore(file);assertFalse(reopened.consume(producer.pending(10).getFirst()));
        producer.ack(event.id());assertEquals(1,reopened.effectCount());assertTrue(producer.pending(10).isEmpty());
        assertThrows(IllegalArgumentException.class,()->reopened.consume(new OutboxStore.Event(event.id(),8)));
    }
    @Test void outboxFailureRollsBackBusiness()throws Exception{
        Path file=directory.resolve("producer.db");var store=new OutboxStore(file);
        try(var c=DriverManager.getConnection("jdbc:sqlite:"+file);var s=c.createStatement()){
            s.execute("CREATE TRIGGER reject_event BEFORE INSERT ON outgoing BEGIN SELECT RAISE(ABORT,'injected'); END");
        }
        assertThrows(SQLException.class,()->store.produce(new OutboxStore.Event("event",1)));
        try(var c=DriverManager.getConnection("jdbc:sqlite:"+file);var s=c.createStatement();var rows=s.executeQuery("SELECT count(*) FROM commands")){assertTrue(rows.next());assertEquals(0,rows.getLong(1));}
    }
}
