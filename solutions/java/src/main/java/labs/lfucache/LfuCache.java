package labs.lfucache;

import java.util.*;

/** Expected O(1) frequency buckets with LRU ties and saturated frequencies. */
public final class LfuCache<K,V>{
    private static final class Entry<V>{V value;long frequency=1;Entry(V value){this.value=value;}}
    private final int capacity;
    private long minimum;
    private final Map<K,Entry<V>> entries=new HashMap<>();
    private final Map<Long,LinkedHashSet<K>> buckets=new HashMap<>();
    public LfuCache(int capacity){if(capacity<=0)throw new IllegalArgumentException("capacity must be positive");this.capacity=capacity;}
    public synchronized Optional<V> get(K key){Objects.requireNonNull(key);var e=entries.get(key);if(e==null)return Optional.empty();touch(key,e);return Optional.of(e.value);}
    public synchronized void put(K key,V value){
        Objects.requireNonNull(key);Objects.requireNonNull(value);var e=entries.get(key);
        if(e!=null){e.value=value;touch(key,e);return;}
        if(entries.size()==capacity){
            var bucket=buckets.get(minimum);K victim=bucket.iterator().next();bucket.remove(victim);entries.remove(victim);if(bucket.isEmpty())buckets.remove(minimum);
        }
        entries.put(key,new Entry<>(value));buckets.computeIfAbsent(1L,f->new LinkedHashSet<>()).add(key);minimum=1;
    }
    private void touch(K key,Entry<V> e){
        long old=e.frequency;var bucket=buckets.get(old);bucket.remove(key);if(bucket.isEmpty())buckets.remove(old);
        if(e.frequency<Long.MAX_VALUE)e.frequency++;
        buckets.computeIfAbsent(e.frequency,f->new LinkedHashSet<>()).add(key);
        if(minimum==old&&!buckets.containsKey(old))minimum=e.frequency;
    }
}
