package labs.postgres;

import java.sql.*;
import java.util.*;

/** READ COMMITTED adapter; counter row serializes all producers before key checks. */
public final class PostgresStore {
  public record Event(String scope, String key, long amount) {
    public Event {
      Objects.requireNonNull(scope);
      Objects.requireNonNull(key);
      if (scope.isEmpty() || key.isEmpty() || amount <= 0) throw new IllegalArgumentException();
    }
  }

  @FunctionalInterface
  public interface Delivery {
    void accept(Event e) throws Exception;
  }

  private final String url;

  public PostgresStore(String url) {
    this.url = Objects.requireNonNull(url);
  }

  public long add(Event event, Runnable beforeCommit) throws SQLException {
    try (var c = DriverManager.getConnection(url)) {
      c.setAutoCommit(false);
      c.setTransactionIsolation(Connection.TRANSACTION_READ_COMMITTED);
      try {
        try (var s = c.createStatement()) {
          s.execute("SET LOCAL lock_timeout='3s'");
          s.execute("SET LOCAL statement_timeout='5s'");
        }
        long current;
        try (var q = c.createStatement();
            var rows = q.executeQuery("SELECT value FROM business WHERE id=1 FOR UPDATE")) {
          if (!rows.next()) throw new SQLException("missing counter");
          current = rows.getLong(1);
        }
        try (var q =
            c.prepareStatement("SELECT amount,response FROM requests WHERE scope=? AND key=?")) {
          q.setString(1, event.scope());
          q.setString(2, event.key());
          try (var rows = q.executeQuery()) {
            if (rows.next()) {
              if (rows.getLong(1) != event.amount())
                throw new IllegalArgumentException("payload conflict");
              long response = rows.getLong(2);
              c.commit();
              return response;
            }
          }
        }
        long response = Math.addExact(current, event.amount());
        try (var q = c.prepareStatement("UPDATE business SET value=? WHERE id=1")) {
          q.setLong(1, response);
          q.executeUpdate();
        }
        try (var q = c.prepareStatement("INSERT INTO requests VALUES(?,?,?,?)")) {
          q.setString(1, event.scope());
          q.setString(2, event.key());
          q.setLong(3, event.amount());
          q.setLong(4, response);
          q.executeUpdate();
        }
        try (var q = c.prepareStatement("INSERT INTO outgoing(scope,key,amount) VALUES(?,?,?)")) {
          bind(q, event);
          q.executeUpdate();
        }
        if (beforeCommit != null) beforeCommit.run();
        c.commit();
        return response;
      } catch (SQLException | RuntimeException | Error e) {
        try {
          c.rollback();
        } catch (SQLException rollback) {
          e.addSuppressed(rollback);
        }
        throw e;
      }
    }
  }

  private static void bind(PreparedStatement q, Event e) throws SQLException {
    q.setString(1, e.scope());
    q.setString(2, e.key());
    q.setLong(3, e.amount());
  }

  public boolean consume(Event event) throws SQLException {
    try (var c = DriverManager.getConnection(url)) {
      c.setAutoCommit(false);
      c.setTransactionIsolation(Connection.TRANSACTION_READ_COMMITTED);
      try {
        try (var s = c.createStatement()) {
          s.execute("SET LOCAL lock_timeout='3s'");
          s.execute("SET LOCAL statement_timeout='5s'");
        }
        boolean added;
        try (var q =
            c.prepareStatement(
                "INSERT INTO inbox VALUES(?,?,?) ON CONFLICT DO NOTHING RETURNING amount")) {
          bind(q, event);
          try (var rows = q.executeQuery()) {
            added = rows.next();
          }
        }
        if (!added) {
          try (var q = c.prepareStatement("SELECT amount FROM inbox WHERE scope=? AND key=?")) {
            q.setString(1, event.scope());
            q.setString(2, event.key());
            try (var rows = q.executeQuery()) {
              if (!rows.next()) throw new SQLException("missing receipt");
              if (rows.getLong(1) != event.amount())
                throw new IllegalArgumentException("payload conflict");
            }
          }
          c.commit();
          return false;
        }
        try (var q = c.prepareStatement("INSERT INTO effects VALUES(?,?,?)")) {
          bind(q, event);
          q.executeUpdate();
        }
        c.commit();
        return true;
      } catch (SQLException | RuntimeException | Error e) {
        try {
          c.rollback();
        } catch (SQLException rollback) {
          e.addSuppressed(rollback);
        }
        throw e;
      }
    }
  }

  public void relay(int limit, Delivery deliver) throws Exception {
    if (limit <= 0) throw new IllegalArgumentException();
    Objects.requireNonNull(deliver);
    var events = new ArrayList<Event>();
    try (var c = DriverManager.getConnection(url);
        var q =
            c.prepareStatement(
                "SELECT scope,key,amount FROM outgoing WHERE NOT sent ORDER BY scope,key LIMIT"
                    + " ?")) {
      q.setInt(1, limit);
      try (var rows = q.executeQuery()) {
        while (rows.next())
          events.add(new Event(rows.getString(1), rows.getString(2), rows.getLong(3)));
      }
    }
    for (var e : events) {
      deliver.accept(e);
      try (var c = DriverManager.getConnection(url);
          var q = c.prepareStatement("UPDATE outgoing SET sent=true WHERE scope=? AND key=?")) {
        q.setString(1, e.scope());
        q.setString(2, e.key());
        q.executeUpdate();
      }
    }
  }
}
