package labs.outbox;

import java.nio.file.Path;
import java.sql.SQLException;
import java.util.List;

/**
 * Real SQLite transactions; producer and consumer normally use separate files.
 * Tables: {@code commands}, {@code outgoing(sent)}, {@code inbox}, {@code effects}.
 */
public final class OutboxStore {
    public record Event(String id, long amount) {}

    public OutboxStore(Path file) throws SQLException {
        // TODO: remember the JDBC URL and create the tables.
    }

    public boolean produce(Event e) throws SQLException {
        throw new UnsupportedOperationException("TODO: implement");
    }

    public boolean consume(Event e) throws SQLException {
        throw new UnsupportedOperationException("TODO: implement");
    }

    public List<Event> pending(int limit) throws SQLException {
        throw new UnsupportedOperationException("TODO: implement");
    }

    public void ack(String id) throws SQLException {
        throw new UnsupportedOperationException("TODO: implement");
    }

    public long effectCount() throws SQLException {
        throw new UnsupportedOperationException("TODO: implement");
    }
}
