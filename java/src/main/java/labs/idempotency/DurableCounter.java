package labs.idempotency;

import java.nio.file.Path;
import java.sql.SQLException;

/**
 * SQLite transactional counter plus retained scoped responses. Use the tables
 * {@code business(id, value)} and {@code requests(scope, key, amount, response)}.
 */
public final class DurableCounter {
    public DurableCounter(Path filename) throws SQLException {
        // TODO: remember the JDBC URL ("jdbc:sqlite:" + path) and create the tables.
    }

    public long add(String scope, String key, long amount) throws SQLException {
        throw new UnsupportedOperationException("TODO: implement");
    }
}
