package labs.postgres;

import static org.junit.jupiter.api.Assertions.*;

import java.nio.file.*;
import java.sql.*;
import java.util.*;
import java.util.concurrent.*;
import org.junit.jupiter.api.Test;

class PostgresIntegrationTest {
  @Test
  void duplicatesCrashRecoveryLostAckAndPostGIS() throws Exception {
    String base = System.getenv("DEPTH_POSTGRES_JDBC");
    assertNotNull(base, "integration requires DEPTH_POSTGRES_JDBC");
    String schema = "depth_java_" + UUID.randomUUID().toString().replace("-", "");
    try (var admin = DriverManager.getConnection(base);
        var control = admin.createStatement()) {
      control.execute("CREATE SCHEMA " + schema);
      try {
        String url =
            base
                + (base.contains("?") ? "&" : "?")
                + "currentSchema="
                + schema
                + ",public&connectTimeout=5&socketTimeout=15";
        var root = Path.of(System.getProperty("labs.root", "."));
        try (var c = DriverManager.getConnection(url);
            var s = c.createStatement()) {
          for (var sql :
              Files.readString(root.resolve("fixtures/postgres.sql"))
                  .split(";")) if (!sql.isBlank()) s.execute(sql);
        }
        var store = new PostgresStore(url);
        var e = new PostgresStore.Event("tenant", "duplicate", 7);
        try (var pool = Executors.newFixedThreadPool(4)) {
          var results = new ArrayList<Future<Long>>();
          for (int i = 0; i < 12; i++) results.add(pool.submit(() -> store.add(e, null)));
          for (var result : results) assertEquals(7, result.get(10, TimeUnit.SECONDS));
        }
        assertThrows(
            IllegalArgumentException.class,
            () -> store.add(new PostgresStore.Event("tenant", "duplicate", 8), null));
        assertThrows(
            IllegalStateException.class,
            () ->
                store.add(
                    new PostgresStore.Event("tenant", "rollback", 4),
                    () -> {
                      throw new IllegalStateException("injected");
                    }));
        for (var phase : List.of("before", "after")) {
          var builder =
              new ProcessBuilder(
                  Path.of(System.getProperty("java.home"), "bin/java").toString(),
                  "-cp",
                  System.getProperty(
                      "surefire.test.class.path", System.getProperty("java.class.path")),
                  CrashClient.class.getName(),
                  phase);
          builder.environment().put("DEPTH_CRASH_JDBC", url);
          builder.redirectErrorStream(true);
          builder.redirectOutput(ProcessBuilder.Redirect.DISCARD);
          var child = builder.start();
          try {
            assertTrue(child.waitFor(15, TimeUnit.SECONDS));
            assertEquals(phase.equals("before") ? 17 : 18, child.exitValue());
          } finally {
            if (child.isAlive()) {
              child.destroyForcibly();
              child.waitFor();
            }
          }
          store.add(new PostgresStore.Event("tenant", phase, 3), null);
        }
        try (var c = DriverManager.getConnection(url);
            var s = c.createStatement();
            var rows = s.executeQuery("SELECT value FROM business WHERE id=1")) {
          assertTrue(rows.next());
          assertEquals(13, rows.getLong(1));
        }
        assertThrows(
            IllegalStateException.class,
            () ->
                store.relay(
                    10,
                    event -> {
                      store.consume(event);
                      throw new IllegalStateException("lost ack");
                    }));
        store.relay(10, store::consume);
        try (var c = DriverManager.getConnection(url);
            var s = c.createStatement();
            var rows = s.executeQuery("SELECT count(*) FROM effects")) {
          assertTrue(rows.next());
          assertEquals(3, rows.getLong(1));
        }
        try (var c = DriverManager.getConnection(url);
            var s = c.createStatement();
            var rows =
                s.executeQuery(
                    Files.readString(
                        root.resolve("fixtures/postgis.sql")))) {
          assertTrue(rows.next());
          assertTrue(rows.getBoolean(1));
        }
      } finally {
        control.execute("DROP SCHEMA " + schema + " CASCADE");
      }
    }
  }
}
