package labs.postgres;

/** Subprocess-only fixture. Halts at a named transaction boundary. */
public final class CrashClient {
  public static void main(String[] args) throws Exception {
    var store = new PostgresStore(System.getenv("DEPTH_CRASH_JDBC"));
    String phase = args[0];
    store.add(
        new PostgresStore.Event("tenant", phase, 3),
        phase.equals("before") ? () -> Runtime.getRuntime().halt(17) : null);
    Runtime.getRuntime().halt(18);
  }
}
