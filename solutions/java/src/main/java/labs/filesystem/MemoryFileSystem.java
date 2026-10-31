package labs.filesystem;

import java.util.*;

/** An in-memory tree, not a host-filesystem sandbox or durable store. */
public final class MemoryFileSystem {
    private record Node(boolean directory,byte[] data) {}
    private final Map<String,Node> nodes=new HashMap<>();
    public MemoryFileSystem(){nodes.put("/",new Node(true,null));}
    private static String normalize(String p){
        Objects.requireNonNull(p);
        if(!p.startsWith("/")||p.indexOf('\0')>=0)throw new IllegalArgumentException("absolute path required");
        var parts=new ArrayList<String>();for(var part:p.split("/")){
            if(part.equals(".")||part.equals(".."))throw new IllegalArgumentException("dot components forbidden");
            if(!part.isEmpty())parts.add(part);
        }
        return "/"+String.join("/",parts);
    }
    private static String parent(String p){int i=p.lastIndexOf('/');return i==0?"/":p.substring(0,i);}
    private void requireParent(String p){var n=nodes.get(parent(p));if(n==null||!n.directory())throw new IllegalStateException("parent missing");}
    public synchronized void mkdir(String p){p=normalize(p);if(nodes.containsKey(p))throw new IllegalStateException("exists");requireParent(p);nodes.put(p,new Node(true,null));}
    public synchronized void write(String p,byte[] data){p=normalize(p);Objects.requireNonNull(data);var old=nodes.get(p);if(old!=null&&old.directory())throw new IllegalStateException("is directory");requireParent(p);nodes.put(p,new Node(false,data.clone()));}
    public synchronized byte[] read(String p){var n=nodes.get(normalize(p));if(n==null||n.directory())throw new IllegalStateException("file missing");return n.data().clone();}
    public synchronized List<String> list(String p){
        p=normalize(p);var n=nodes.get(p);if(n==null||!n.directory())throw new IllegalStateException("directory missing");
        var result=new ArrayList<String>();for(var key:nodes.keySet())if(!key.equals(p)&&parent(key).equals(p))result.add(key.substring(key.lastIndexOf('/')+1));Collections.sort(result);return List.copyOf(result);
    }
    public synchronized void move(String from,String to){
        from=normalize(from);to=normalize(to);
        if(from.equals("/"))throw new IllegalArgumentException("cannot move root");
        var n=nodes.get(from);if(n==null)throw new IllegalStateException("source missing");if(from.equals(to))return;
        if(nodes.containsKey(to))throw new IllegalStateException("destination exists");
        if(n.directory()&&to.startsWith(from+"/"))throw new IllegalArgumentException("move into descendant");
        requireParent(to);var moved=new HashMap<String,Node>();var removed=new ArrayList<String>();
        for(var entry:nodes.entrySet())if(entry.getKey().equals(from)||entry.getKey().startsWith(from+"/")){moved.put(to+entry.getKey().substring(from.length()),entry.getValue());removed.add(entry.getKey());}
        for(var key:removed)nodes.remove(key);nodes.putAll(moved);
    }
}
