package labs.cells;

import com.google.common.geometry.*;
import com.uber.h3core.H3Core;
import java.io.IOException;

/** Pinned API examples, not a complete region-query adapter. */
public final class Cells {
  public record IDs(long h3, S2CellId s2) {}

  private final H3Core h3;

  public Cells() throws IOException {
    h3 = H3Core.newInstance();
  }

  public IDs ids(double latitude, double longitude, int resolution, int level) {
    if (!Double.isFinite(latitude)
        || !Double.isFinite(longitude)
        || latitude < -90
        || latitude > 90
        || longitude < -180
        || longitude > 180
        || resolution < 0
        || resolution > 15
        || level < 0
        || level > 30) throw new IllegalArgumentException();
    return new IDs(
        h3.latLngToCell(latitude, longitude, resolution),
        S2CellId.fromLatLng(S2LatLng.fromDegrees(latitude, longitude)).parent(level));
  }

  public boolean h3CenterRoundTrip(long cell, int resolution) {
    var center = h3.cellToLatLng(cell);
    return h3.latLngToCell(center.lat, center.lng, resolution) == cell;
  }
}
