package labs.fencing;

import static org.junit.jupiter.api.Assertions.*;

import org.junit.jupiter.api.Test;

class FencedResourceTest {
  @Test
  void rejectsOldOwner() {
    var r = new FencedResource();
    assertTrue(r.write(7, 0, "old"));
    assertTrue(r.write(8, 1, "new"));
    assertFalse(r.write(7, 2, "stale"));
    assertFalse(r.write(8, 1, "duplicate"));
    assertEquals(new FencedResource.Snapshot("new", 2), r.read());
  }
}
