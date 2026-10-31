package labs.cells;

import static org.junit.jupiter.api.Assertions.*;

import com.google.common.geometry.*;
import org.junit.jupiter.api.Test;

class CellsTest {
  @Test
  void cellRoundTrips() throws Exception {
    var cells = new Cells();
    for (var p : new double[][] {{0, 0}, {89.9, 179.9}, {-89.9, -179.9}, {0, 180}}) {
      var ids = cells.ids(p[0], p[1], 7, 12);
      assertTrue(cells.h3CenterRoundTrip(ids.h3(), 7));
      assertEquals(12, ids.s2().level());
      assertTrue(new S2Cell(ids.s2()).contains(S2LatLng.fromDegrees(p[0], p[1]).toPoint()));
    }
  }
}
