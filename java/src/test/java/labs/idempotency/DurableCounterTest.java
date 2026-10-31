package labs.idempotency;
import static org.junit.jupiter.api.Assertions.*;
import java.nio.file.Path;
import java.sql.*;
import java.util.ArrayList;
import java.util.concurrent.Executors;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.io.TempDir;
class DurableCounterTest{
    @TempDir Path directory;
    @Test void duplicatesConflictAndReopen()throws Exception{
        Path file=directory.resolve("counter.db");var counter=new DurableCounter(file);var other=new DurableCounter(file);
        try(var executor=Executors.newFixedThreadPool(4)){
            var futures=new ArrayList<java.util.concurrent.Future<Long>>();
            for(int i=0;i<12;i++){var c=i%2==0?counter:other;futures.add(executor.submit(()->c.add("tenant","key",7)));}
            for(var result:futures)assertEquals(7,result.get());
        }
        assertThrows(IllegalArgumentException.class,()->counter.add("tenant","key",8));
        var reopened=new DurableCounter(file);assertEquals(7,reopened.add("tenant","key",7));assertEquals(8,reopened.add("tenant","next",1));
    }
    @Test void failureRollsBackMutation()throws Exception{
        Path file=directory.resolve("counter.db");var counter=new DurableCounter(file);
        try(var c=DriverManager.getConnection("jdbc:sqlite:"+file);var s=c.createStatement()){
            s.execute("CREATE TRIGGER fail_insert BEFORE INSERT ON requests WHEN NEW.key='fail' BEGIN SELECT RAISE(ABORT,'injected'); END");
        }
        assertThrows(SQLException.class,()->counter.add("tenant","fail",7));assertEquals(1,counter.add("tenant","ok",1));
    }
}
