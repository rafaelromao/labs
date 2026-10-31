package labs.dag;

import static org.junit.jupiter.api.Assertions.*;
import java.util.List;
import java.util.Map;
import org.junit.jupiter.api.Test;

class DependencyExecutorTest {
    @Test void skipsFailedPrerequisites() throws Exception {
        var tasks=Map.of("a",new DependencyExecutor.Task<Integer>(List.of(),()->1),
                "b",new DependencyExecutor.Task<Integer>(List.of("a"),()->{throw new IllegalStateException();}),
                "c",new DependencyExecutor.Task<Integer>(List.of("b"),()->{fail("blocked task ran");return 0;}));
        var result=DependencyExecutor.execute(tasks,2,3,()->false);
        assertEquals(1,result.get("a").value()); assertNotNull(result.get("b").failure()); assertTrue(result.get("c").skipped());
    }
    @Test void rejectsCycleBeforeEffects() {
        var tasks=Map.of("a",new DependencyExecutor.Task<Integer>(List.of("b"),()->{fail("ran");return 0;}),
                "b",new DependencyExecutor.Task<Integer>(List.of("a"),()->{fail("ran");return 0;}));
        assertThrows(IllegalArgumentException.class,()->DependencyExecutor.execute(tasks,1,2,()->false));
    }
}
