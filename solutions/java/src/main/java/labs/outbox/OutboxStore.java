package labs.outbox;
import java.nio.file.Path;
import java.sql.*;
import java.util.*;

/** Real SQLite transactions; producer and consumer normally use separate files. */
public final class OutboxStore{
    public record Event(String id,long amount){}
    private final String url;
    public OutboxStore(Path file)throws SQLException{
        url="jdbc:sqlite:"+file.toAbsolutePath();
        try(var c=connect();var s=c.createStatement()){
            s.execute("CREATE TABLE IF NOT EXISTS commands(id TEXT PRIMARY KEY,amount INTEGER NOT NULL)");
            s.execute("CREATE TABLE IF NOT EXISTS outgoing(id TEXT PRIMARY KEY,amount INTEGER NOT NULL,sent INTEGER NOT NULL DEFAULT 0)");
            s.execute("CREATE TABLE IF NOT EXISTS inbox(id TEXT PRIMARY KEY,amount INTEGER NOT NULL)");
            s.execute("CREATE TABLE IF NOT EXISTS effects(id TEXT PRIMARY KEY,amount INTEGER NOT NULL)");
        }
    }
    private Connection connect()throws SQLException{
        var c=DriverManager.getConnection(url);
        try(var s=c.createStatement()){s.execute("PRAGMA busy_timeout=5000");s.execute("PRAGMA synchronous=FULL");}
        catch(SQLException e){c.close();throw e;}return c;
    }
    // Table names are private fixed constants, never user input.
    private boolean transact(Event event,String first,String second)throws SQLException{
        Objects.requireNonNull(event);Objects.requireNonNull(event.id());if(event.id().isEmpty()||event.amount()<=0)throw new IllegalArgumentException("invalid event");
        try(var c=connect();var control=c.createStatement()){
            control.execute("BEGIN IMMEDIATE");boolean committed=false;
            try{
                try(var q=c.prepareStatement("SELECT amount FROM "+first+" WHERE id=?")){
                    q.setString(1,event.id());try(var rows=q.executeQuery()){if(rows.next()){
                        if(rows.getLong(1)!=event.amount())throw new IllegalArgumentException("payload conflict");
                        control.execute("COMMIT");committed=true;return false;
                    }}
                }
                for(var table:List.of(first,second))try(var q=c.prepareStatement("INSERT INTO "+table+"(id,amount) VALUES(?,?)")){
                    q.setString(1,event.id());q.setLong(2,event.amount());q.executeUpdate();
                }
                control.execute("COMMIT");committed=true;return true;
            }finally{if(!committed)control.execute("ROLLBACK");}
        }
    }
    public boolean produce(Event e)throws SQLException{return transact(e,"commands","outgoing");}
    public boolean consume(Event e)throws SQLException{return transact(e,"inbox","effects");}
    public List<Event> pending(int limit)throws SQLException{
        if(limit<=0)throw new IllegalArgumentException("positive batch size required");
        try(var c=connect();var q=c.prepareStatement("SELECT id,amount FROM outgoing WHERE sent=0 ORDER BY id LIMIT ?")){
            q.setInt(1,limit);var events=new ArrayList<Event>();try(var rows=q.executeQuery()){while(rows.next())events.add(new Event(rows.getString(1),rows.getLong(2)));}return List.copyOf(events);
        }
    }
    public void ack(String id)throws SQLException{try(var c=connect();var q=c.prepareStatement("UPDATE outgoing SET sent=1 WHERE id=?")){q.setString(1,Objects.requireNonNull(id));q.executeUpdate();}}
    public long effectCount()throws SQLException{try(var c=connect();var q=c.createStatement();var rows=q.executeQuery("SELECT count(*) FROM effects")){rows.next();return rows.getLong(1);}}
}
