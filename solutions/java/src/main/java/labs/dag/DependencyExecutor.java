package labs.dag;

import labs.fanout.RequestAggregator;
import java.util.*;
import java.util.concurrent.Callable;
import java.util.function.BooleanSupplier;

/** Validates a finite DAG, then executes bounded topological waves. */
public final class DependencyExecutor {
    public record Task<V>(List<String> dependsOn,Callable<V> run) {
        public Task { dependsOn=List.copyOf(dependsOn); Objects.requireNonNull(run); }
    }
    public record Outcome<V>(V value,Throwable failure,boolean skipped) {}
    private DependencyExecutor() {}
    public static <V> Map<String,Outcome<V>> execute(Map<String,Task<V>> input,int concurrency,int maxTasks,
            BooleanSupplier cancelled) throws InterruptedException {
        if (concurrency<=0 || maxTasks<=0 || input.size()>maxTasks) throw new IllegalArgumentException("invalid limits");
        var tasks=new TreeMap<>(input); var degree=new TreeMap<String,Integer>();
        var children=new HashMap<String,List<String>>();
        for (var entry:tasks.entrySet()) {
            var seen=new HashSet<String>();
            for (var dep:entry.getValue().dependsOn()) {
                if (!tasks.containsKey(dep) || !seen.add(dep)) throw new IllegalArgumentException("missing or duplicate dependency");
                children.computeIfAbsent(dep,k->new ArrayList<>()).add(entry.getKey());
            }
            degree.put(entry.getKey(),seen.size());
        }
        var check=new HashMap<>(degree); var ready=new ArrayList<String>();
        check.forEach((id,n)->{if(n==0)ready.add(id);}); int visited=0;
        while (!ready.isEmpty()) {
            String id=ready.removeLast(); visited++;
            for (var child:children.getOrDefault(id,List.of())) if (check.merge(child,-1,Integer::sum)==0) ready.add(child);
        }
        if (visited!=tasks.size()) throw new IllegalArgumentException("cycle");
        var result=new TreeMap<String,Outcome<V>>();
        degree.forEach((id,n)->{if(n==0)ready.add(id);});
        while (!ready.isEmpty()) {
            var calls=new TreeMap<String,Callable<V>>();
            for (var id:ready) {
                boolean blocked=tasks.get(id).dependsOn().stream().anyMatch(dep->result.get(dep).failure()!=null || result.get(dep).skipped());
                if (blocked) result.put(id,new Outcome<>(null,null,true)); else calls.put(id,tasks.get(id).run());
            }
            var outcomes=RequestAggregator.run(calls,concurrency,maxTasks,cancelled);
            outcomes.forEach((id,out)->result.put(id,new Outcome<>(out.value(),out.failure(),false)));
            var next=new ArrayList<String>();
            for (var id:ready) for (var child:children.getOrDefault(id,List.of())) if (degree.merge(child,-1,Integer::sum)==0) next.add(child);
            ready.clear(); ready.addAll(next);
        }
        return Collections.unmodifiableMap(result);
    }
}
