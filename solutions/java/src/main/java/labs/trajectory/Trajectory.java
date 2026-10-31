package labs.trajectory;

import labs.nearby.NearbyIndex;
import java.util.*;

/** One bounded track, with timestamps in nonnegative milliseconds. */
public final class Trajectory {
    public record Observation(long millis,NearbyIndex.Point position) {}
    public record Segment(long fromMillis,long toMillis,double meters,double speedMetersPerSecond,boolean implausible) {}
    private final int capacity;
    private final NavigableMap<Long,Observation> samples=new TreeMap<>();
    public Trajectory(int capacity) { if(capacity<=0)throw new IllegalArgumentException("positive capacity required");this.capacity=capacity; }
    public synchronized void put(Observation o) {
        Objects.requireNonNull(o);
        if(o.millis()<0)throw new IllegalArgumentException("negative time");
        NearbyIndex.distance(o.position(),o.position());
        samples.put(o.millis(),o);
        while(samples.size()>capacity)samples.pollFirstEntry();
    }
    public synchronized List<Segment> segments(double maxSpeed) {
        if(!Double.isFinite(maxSpeed)||maxSpeed<=0)throw new IllegalArgumentException("invalid speed threshold");
        var result=new ArrayList<Segment>(); Observation previous=null;
        for(var current:samples.values()) {
            if(previous!=null) {
                double meters=NearbyIndex.distance(previous.position(),current.position());
                double speed=meters/((current.millis()-previous.millis())/1000.0);
                result.add(new Segment(previous.millis(),current.millis(),meters,speed,speed>maxSpeed));
            }
            previous=current;
        }
        return List.copyOf(result);
    }
}
