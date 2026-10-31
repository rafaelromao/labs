package labs.fencing;

/** In-memory destination fencing; tokens must come from an authoritative issuer. */
public final class FencedResource {
  public record Snapshot(String value, long revision) {}

  private long token, revision;
  private String value = "";

  public synchronized boolean write(long token, long expected, String value) {
    if (token < 0 || token < this.token || expected != revision || revision == Long.MAX_VALUE)
      return false;
    this.token = token;
    this.value = value;
    revision++;
    return true;
  }

  public synchronized Snapshot read() {
    return new Snapshot(value, revision);
  }
}
