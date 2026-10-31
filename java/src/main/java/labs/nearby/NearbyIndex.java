package labs.nearby;

import java.util.List;
import java.util.Objects;

/** Spherical search over latitude/longitude points in degrees. */
public final class NearbyIndex {
    public static final double EARTH_RADIUS_METERS = 6371008.8;

    public record Point(double latitude, double longitude) {}

    public record Hit(String id, double distanceMeters) {}

    public void put(String id, Point p) {
        throw new UnsupportedOperationException("TODO: implement");
    }

    public boolean delete(String id) {
        throw new UnsupportedOperationException("TODO: implement");
    }

    /** Provided (geofence and trajectory use it): haversine distance in metres. */
    public static double distance(Point a, Point b) {
        validate(a);
        validate(b);
        double dlat = Math.toRadians(b.latitude() - a.latitude());
        double dlon = Math.toRadians(b.longitude() - a.longitude());
        double h = Math.pow(Math.sin(dlat / 2), 2)
                + Math.cos(Math.toRadians(a.latitude())) * Math.cos(Math.toRadians(b.latitude()))
                * Math.pow(Math.sin(dlon / 2), 2);
        return 2 * EARTH_RADIUS_METERS * Math.asin(Math.sqrt(Math.max(0, Math.min(1, h))));
    }

    /** Provided: rejects null, non-finite or out-of-range coordinates. */
    static void validate(Point p) {
        Objects.requireNonNull(p);
        if (!Double.isFinite(p.latitude()) || !Double.isFinite(p.longitude())
                || p.latitude() < -90 || p.latitude() > 90
                || p.longitude() < -180 || p.longitude() > 180) {
            throw new IllegalArgumentException("invalid coordinates");
        }
    }

    public List<Hit> radius(Point center, double meters) {
        throw new UnsupportedOperationException("TODO: implement");
    }

    public List<Hit> nearest(Point center, int k) {
        throw new UnsupportedOperationException("TODO: implement");
    }
}
