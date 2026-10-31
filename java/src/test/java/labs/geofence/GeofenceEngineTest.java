package labs.geofence;
import static org.junit.jupiter.api.Assertions.*;
import labs.nearby.NearbyIndex;
import org.junit.jupiter.api.Test;
class GeofenceEngineTest{
    @Test void transitionsHysteresisAndStale(){
        var e=new GeofenceEngine(new NearbyIndex.Point(0,0),1000,500,1);
        assertTrue(e.observe("a",1,new NearbyIndex.Point(0,0)).changed());
        var middle=e.observe("a",2,new NearbyIndex.Point(0,0.01));assertTrue(middle.inside());assertFalse(middle.changed());
        var stale=e.observe("a",1,new NearbyIndex.Point(0,1));assertTrue(stale.stale());assertTrue(stale.inside());
        var exit=e.observe("a",3,new NearbyIndex.Point(0,0.02));assertTrue(exit.changed());assertFalse(exit.inside());
        assertThrows(IllegalStateException.class,()->e.observe("b",1,new NearbyIndex.Point(0,0)));
        e.forget("a");assertTrue(e.observe("b",1,new NearbyIndex.Point(0,0)).inside());
    }
}
