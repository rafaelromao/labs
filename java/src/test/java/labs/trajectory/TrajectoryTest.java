package labs.trajectory;

import static org.junit.jupiter.api.Assertions.*;
import labs.nearby.NearbyIndex;
import org.junit.jupiter.api.Test;

class TrajectoryTest {
    @Test void orderingReplacementRetentionAndJumps() {
        var t=new Trajectory(3);
        t.put(new Trajectory.Observation(2000,new NearbyIndex.Point(0,0.01)));
        t.put(new Trajectory.Observation(0,new NearbyIndex.Point(0,0)));
        t.put(new Trajectory.Observation(1000,new NearbyIndex.Point(0,0)));
        t.put(new Trajectory.Observation(1000,new NearbyIndex.Point(0,0.005)));
        assertEquals(2,t.segments(100).size());assertTrue(t.segments(100).getFirst().implausible());
        t.put(new Trajectory.Observation(3000,new NearbyIndex.Point(0,0.01)));
        assertEquals(1000,t.segments(100).getFirst().fromMillis());assertEquals(0,t.segments(100).getLast().speedMetersPerSecond());
        t.put(new Trajectory.Observation(0,new NearbyIndex.Point(0,1)));
        assertEquals(1000,t.segments(100).getFirst().fromMillis());
    }
}
