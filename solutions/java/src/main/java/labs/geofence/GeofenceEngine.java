package labs.geofence;
import labs.nearby.NearbyIndex;
import java.util.*;
/** One immutable circular fence; per-object sequences and spatial hysteresis. */
public final class GeofenceEngine{
    public record Result(boolean inside,boolean changed,boolean stale,long sequence){}
    private record State(long sequence,boolean inside){}
    private final NearbyIndex.Point center;private final double radius,margin;private final int limit;
    private final Map<String,State> objects=new HashMap<>();
    public GeofenceEngine(NearbyIndex.Point center,double radius,double margin,int limit){
        NearbyIndex.distance(center,center);
        if(radius<0||margin<0||!Double.isFinite(radius)||!Double.isFinite(margin)||!Double.isFinite(radius+margin)||limit<=0)throw new IllegalArgumentException("invalid configuration");
        this.center=center;this.radius=radius;this.margin=margin;this.limit=limit;
    }
    public synchronized Result observe(String id,long sequence,NearbyIndex.Point position){
        Objects.requireNonNull(id);if(sequence<0)throw new IllegalArgumentException("negative sequence");
        double d=NearbyIndex.distance(center,position);var old=objects.get(id);
        if(old!=null&&sequence<=old.sequence())return new Result(old.inside(),false,true,old.sequence());
        if(old==null&&objects.size()==limit)throw new IllegalStateException("object capacity");
        boolean before=old!=null&&old.inside();boolean inside=d<=radius+(before?margin:0);
        objects.put(id,new State(sequence,inside));return new Result(inside,inside!=before,false,sequence);
    }
    public synchronized void forget(String id){objects.remove(Objects.requireNonNull(id));}
}
