package labs.algorithms;

import java.nio.ByteBuffer;
import java.nio.charset.StandardCharsets;
import java.security.*;
import java.util.*;

/** Sequential teaching algorithms. Range sums require representable long intermediates. */
public final class Structures {
  private Structures() {}

  public record Priority(String key, long value) {}

  public static final class IndexedQueue {
    private final ArrayList<Priority> items = new ArrayList<>();
    private final Map<String, Integer> positions = new HashMap<>();

    public int size() {
      return items.size();
    }

    private boolean less(int a, int b) {
      var x = items.get(a);
      var y = items.get(b);
      return x.value() < y.value() || (x.value() == y.value() && x.key().compareTo(y.key()) < 0);
    }

    private void swap(int a, int b) {
      var x = items.get(a);
      items.set(a, items.get(b));
      items.set(b, x);
      positions.put(items.get(a).key(), a);
      positions.put(items.get(b).key(), b);
    }

    private void fix(int i) {
      if (i > 0 && less(i, (i - 1) / 2)) {
        while (i > 0 && less(i, (i - 1) / 2)) {
          int p = (i - 1) / 2;
          swap(i, p);
          i = p;
        }
        return;
      }
      while (true) {
        int c = 2 * i + 1;
        if (c >= size()) return;
        if (c + 1 < size() && less(c + 1, c)) c++;
        if (!less(c, i)) return;
        swap(c, i);
        i = c;
      }
    }

    // region: queue-put
    public void put(String key, long value) {
      Objects.requireNonNull(key);
      Integer i = positions.get(key);
      if (i != null) {
        items.set(i, new Priority(key, value));
        fix(i);
        return;
      }
      positions.put(key, size());
      items.add(new Priority(key, value));
      fix(size() - 1);
    }

    // endregion: queue-put
    public Optional<Priority> peek() {
      return items.isEmpty() ? Optional.empty() : Optional.of(items.getFirst());
    }

    public boolean remove(String key) {
      Integer i = positions.get(key);
      if (i == null) return false;
      int last = size() - 1;
      swap(i, last);
      items.removeLast();
      positions.remove(key);
      if (i < size()) fix(i);
      return true;
    }

    public Optional<Priority> pop() {
      var v = peek();
      v.ifPresent(p -> remove(p.key()));
      return v;
    }

    public boolean valid() {
      if (size() != positions.size()) return false;
      for (int i = 0; i < size(); i++)
        if (!Objects.equals(positions.get(items.get(i).key()), i)
            || (i > 0 && less(i, (i - 1) / 2))) return false;
      return true;
    }
  }

  public static int lower(long[] a, long x) {
    int l = 0, r = a.length;
    while (l < r) {
      int m = l + (r - l) / 2;
      if (a[m] < x) l = m + 1;
      else r = m;
    }
    return l;
  }

  public static int upper(long[] a, long x) {
    int l = 0, r = a.length;
    while (l < r) {
      int m = l + (r - l) / 2;
      if (a[m] <= x) l = m + 1;
      else r = m;
    }
    return l;
  }

  public static final class Trie {
    private static final class Node {
      boolean terminal;
      TreeMap<Character, Node> children = new TreeMap<>();
    }

    private final Node root = new Node();

    private static void ascii(String s) {
      for (char c : s.toCharArray())
        if (c > 127) throw new IllegalArgumentException("ASCII required");
    }

    public void insert(String s) {
      ascii(s);
      var n = root;
      for (char c : s.toCharArray()) n = n.children.computeIfAbsent(c, k -> new Node());
      n.terminal = true;
    }

    public boolean has(String s) {
      ascii(s);
      var n = root;
      for (char c : s.toCharArray()) {
        n = n.children.get(c);
        if (n == null) return false;
      }
      return n.terminal;
    }

    public boolean delete(String s) {
      ascii(s);
      return remove(root, s, 0);
    }

    private boolean remove(Node n, String s, int i) {
      if (i == s.length()) {
        boolean old = n.terminal;
        n.terminal = false;
        return old;
      }
      var child = n.children.get(s.charAt(i));
      if (child == null) return false;
      boolean ok = remove(child, s, i + 1);
      if (!child.terminal && child.children.isEmpty()) n.children.remove(s.charAt(i));
      return ok;
    }

    public List<String> prefix(String s) {
      ascii(s);
      var n = root;
      for (char c : s.toCharArray()) {
        n = n.children.get(c);
        if (n == null) return List.of();
      }
      var out = new ArrayList<String>();
      walk(n, s, out);
      return List.copyOf(out);
    }

    private void walk(Node n, String p, List<String> out) {
      if (n.terminal) out.add(p);
      n.children.forEach((c, child) -> walk(child, p + c, out));
    }
  }

  public static final class Disjoint {
    private final int[] parent, size;
    private int count;

    public Disjoint(int n) {
      if (n < 0) throw new IllegalArgumentException();
      parent = new int[n];
      size = new int[n];
      count = n;
      for (int i = 0; i < n; i++) {
        parent[i] = i;
        size[i] = 1;
      }
    }

    public int find(int i) {
      Objects.checkIndex(i, parent.length);
      while (i != parent[i]) {
        parent[i] = parent[parent[i]];
        i = parent[i];
      }
      return i;
    }

    public boolean union(int a, int b) {
      a = find(a);
      b = find(b);
      if (a == b) return false;
      if (size[a] < size[b]) {
        int t = a;
        a = b;
        b = t;
      }
      parent[b] = a;
      size[a] += size[b];
      count--;
      return true;
    }

    public int count() {
      return count;
    }

    public int size(int i) {
      return size[find(i)];
    }
  }

  public static final class Fenwick {
    private final long[] tree;

    public Fenwick(int n) {
      if (n < 0) throw new IllegalArgumentException();
      tree = new long[n + 1];
    }

    public void add(int i, long delta) {
      Objects.checkIndex(i, tree.length - 1);
      for (i++; i < tree.length; i += i & -i) tree[i] += delta;
    }

    public long prefix(int end) {
      if (end < 0 || end >= tree.length) throw new IndexOutOfBoundsException();
      long v = 0;
      for (; end > 0; end -= end & -end) v += tree[end];
      return v;
    }

    public long sum(int start, int end) {
      if (start < 0 || start > end) throw new IndexOutOfBoundsException();
      return prefix(end) - prefix(start);
    }
  }

  public static final class SegmentTree {
    private final int n;
    private final long[] tree;

    public SegmentTree(int n) {
      if (n < 0) throw new IllegalArgumentException();
      this.n = n;
      tree = new long[Math.multiplyExact(2, n)];
    }

    public void set(int i, long v) {
      Objects.checkIndex(i, n);
      i += n;
      tree[i] = v;
      while (i > 1) {
        i /= 2;
        tree[i] = tree[2 * i] + tree[2 * i + 1];
      }
    }

    public long sum(int a, int b) {
      if (a < 0 || b < a || b > n) throw new IndexOutOfBoundsException();
      long v = 0;
      for (a += n, b += n; a < b; a /= 2, b /= 2) {
        if ((a & 1) != 0) v += tree[a++];
        if ((b & 1) != 0) v += tree[--b];
      }
      return v;
    }
  }

  public static final class Bloom {
    private final byte[] bits;
    private final int hashes;

    public Bloom(int bytes, int hashes) {
      if (bytes <= 0 || hashes <= 0 || hashes > 8) throw new IllegalArgumentException();
      bits = new byte[bytes];
      this.hashes = hashes;
    }

    private int[] indices(String key) {
      try {
        byte[] digest =
            MessageDigest.getInstance("SHA-256").digest(key.getBytes(StandardCharsets.UTF_8));
        int[] out = new int[hashes];
        for (int i = 0; i < hashes; i++)
          out[i] =
              (int)
                  (Integer.toUnsignedLong(ByteBuffer.wrap(digest, 4 * i, 4).getInt())
                      / 1
                      % ((long) bits.length * 8));
        return out;
      } catch (NoSuchAlgorithmException e) {
        throw new AssertionError(e);
      }
    }

    public void add(String key) {
      for (int i : indices(key)) bits[i / 8] |= (byte) (1 << (i % 8));
    }

    public boolean has(String key) {
      for (int i : indices(key)) if ((bits[i / 8] & (1 << (i % 8))) == 0) return false;
      return true;
    }
  }

  public static List<Long> maxima(long[] a, int width) {
    if (width <= 0) throw new IllegalArgumentException();
    var q = new ArrayDeque<Integer>();
    var out = new ArrayList<Long>();
    for (int i = 0; i < a.length; i++) {
      while (!q.isEmpty() && q.getFirst() <= i - width) q.removeFirst();
      while (!q.isEmpty() && a[q.getLast()] <= a[i]) q.removeLast();
      q.addLast(i);
      if (i + 1 >= width) out.add(a[q.getFirst()]);
    }
    return List.copyOf(out);
  }

  public static List<Long> topK(long[] a, int k) {
    if (k < 0) throw new IllegalArgumentException();
    var q = new PriorityQueue<Long>();
    for (long v : a) {
      if (q.size() < k) q.add(v);
      else if (k > 0 && v > q.element()) {
        q.remove();
        q.add(v);
      }
    }
    var out = new ArrayList<>(q);
    out.sort(Comparator.reverseOrder());
    return List.copyOf(out);
  }

  public static List<Long> mergeSorted(List<long[]> streams) {
    record Head(long value, int stream, int index) {}
    var q =
        new PriorityQueue<Head>(
            Comparator.comparingLong(Head::value).thenComparingInt(Head::stream));
    for (int i = 0; i < streams.size(); i++) {
      var a = streams.get(i);
      for (int j = 1; j < a.length; j++)
        if (a[j] < a[j - 1]) throw new IllegalArgumentException("unsorted");
      if (a.length > 0) q.add(new Head(a[0], i, 0));
    }
    var out = new ArrayList<Long>();
    while (!q.isEmpty()) {
      var h = q.remove();
      out.add(h.value());
      var a = streams.get(h.stream());
      if (h.index() + 1 < a.length) q.add(new Head(a[h.index() + 1], h.stream(), h.index() + 1));
    }
    return List.copyOf(out);
  }

  public record Interval(long start, long end) {}

  public static List<Interval> mergeIntervals(List<Interval> input) {
    var a = new ArrayList<>(input);
    for (var v : a) if (v.start() >= v.end()) throw new IllegalArgumentException();
    a.sort(Comparator.comparingLong(Interval::start));
    var out = new ArrayList<Interval>();
    for (var v : a) {
      if (out.isEmpty() || v.start() > out.getLast().end()) out.add(v);
      else if (v.end() > out.getLast().end())
        out.set(out.size() - 1, new Interval(out.getLast().start(), v.end()));
    }
    return List.copyOf(out);
  }

  public static List<Integer> traverse(int[][] graph, int start, boolean depth) {
    Objects.checkIndex(start, graph.length);
    for (var edges : graph) for (int v : edges) Objects.checkIndex(v, graph.length);
    boolean[] seen = new boolean[graph.length];
    var q = new ArrayDeque<Integer>();
    q.add(start);
    var out = new ArrayList<Integer>();
    while (!q.isEmpty()) {
      int v = depth ? q.removeLast() : q.removeFirst();
      if (seen[v]) continue;
      seen[v] = true;
      out.add(v);
      for (int w : graph[v]) q.addLast(w);
    }
    return List.copyOf(out);
  }

  public static Optional<List<Integer>> topological(int[][] g) {
    int[] degree = new int[g.length];
    for (var edges : g)
      for (int v : edges) {
        Objects.checkIndex(v, g.length);
        degree[v]++;
      }
    var q = new ArrayDeque<Integer>();
    for (int i = 0; i < degree.length; i++) if (degree[i] == 0) q.add(i);
    var out = new ArrayList<Integer>();
    while (!q.isEmpty()) {
      int v = q.removeFirst();
      out.add(v);
      for (int w : g[v]) if (--degree[w] == 0) q.add(w);
    }
    return out.size() == g.length ? Optional.of(List.copyOf(out)) : Optional.empty();
  }
}
