package labs.idempotency;
import java.nio.file.Path;
import java.sql.*;
import java.util.Objects;

/** SQLite transactional counter plus retained scoped responses. */
public final class DurableCounter{
    private final String url;
    public DurableCounter(Path filename)throws SQLException{
        url="jdbc:sqlite:"+filename.toAbsolutePath();
        try(var c=connect();var s=c.createStatement()){
            s.execute("CREATE TABLE IF NOT EXISTS business(id INTEGER PRIMARY KEY CHECK(id=1),value INTEGER NOT NULL)");
            s.execute("INSERT OR IGNORE INTO business VALUES(1,0)");
            s.execute("CREATE TABLE IF NOT EXISTS requests(scope TEXT NOT NULL,key TEXT NOT NULL,amount INTEGER NOT NULL,response INTEGER NOT NULL,PRIMARY KEY(scope,key))");
        }
    }
    private Connection connect()throws SQLException{
        var c=DriverManager.getConnection(url);
        try(var s=c.createStatement()){s.execute("PRAGMA busy_timeout=5000");s.execute("PRAGMA synchronous=FULL");}
        catch(SQLException e){c.close();throw e;}return c;
    }
    public long add(String scope,String key,long amount)throws SQLException{
        Objects.requireNonNull(scope);Objects.requireNonNull(key);
        if(scope.isEmpty()||key.isEmpty()||amount<=0)throw new IllegalArgumentException("invalid input");
        try(var c=connect();var control=c.createStatement()){
            control.execute("BEGIN IMMEDIATE");boolean committed=false;
            try{
                try(var q=c.prepareStatement("SELECT amount,response FROM requests WHERE scope=? AND key=?")){
                    q.setString(1,scope);q.setString(2,key);
                    try(var rows=q.executeQuery()){if(rows.next()){
                        if(rows.getLong(1)!=amount)throw new IllegalArgumentException("key conflict");
                        long response=rows.getLong(2);control.execute("COMMIT");committed=true;return response;
                    }}
                }
                long current;try(var q=control.executeQuery("SELECT value FROM business WHERE id=1")){if(!q.next())throw new SQLException("missing business row");current=q.getLong(1);}
                long response=Math.addExact(current,amount);
                try(var q=c.prepareStatement("UPDATE business SET value=? WHERE id=1")){q.setLong(1,response);q.executeUpdate();}
                try(var q=c.prepareStatement("INSERT INTO requests(scope,key,amount,response) VALUES(?,?,?,?)")){
                    q.setString(1,scope);q.setString(2,key);q.setLong(3,amount);q.setLong(4,response);q.executeUpdate();
                }
                control.execute("COMMIT");committed=true;return response;
            }finally{if(!committed)control.execute("ROLLBACK");}
        }
    }
}
