package labs.algorithms;

import static org.junit.jupiter.api.Assertions.*;

import java.nio.file.*;
import java.util.*;
import org.junit.jupiter.api.Test;

class StructuresTest {
  @Test
  void indexedQueueModel() {
    var r = new Random(42);
    var q = new Structures.IndexedQueue();
    var model = new HashMap<String, Long>();
    for (int step = 0; step < 3000; step++) {
      String key = "" + r.nextInt(40);
      if (r.nextInt(3) == 0) {
        boolean exists = model.containsKey(key);
        assertEquals(exists, q.remove(key));
        model.remove(key);
      } else {
        long value = r.nextInt(100) - 50;
        q.put(key, value);
        model.put(key, value);
      }
      assertTrue(q.valid());
      assertEquals(model.size(), q.size());
      var want =
          model.entrySet().stream()
              .map(e -> new Structures.Priority(e.getKey(), e.getValue()))
              .min(
                  Comparator.comparingLong(Structures.Priority::value)
                      .thenComparing(Structures.Priority::key));
      assertEquals(want, q.peek());
    }
    long previous = Long.MIN_VALUE;
    while (q.size() > 0) {
      long current = q.pop().orElseThrow().value();
      assertTrue(current >= previous);
      previous = current;
      assertTrue(q.valid());
    }
  }

  @Test
  void rangesAndSharedFixture() throws Exception {
    var root = Path.of(System.getProperty("labs.root", "."));
    var file = root.resolve("fixtures/structures.tsv");
    var f = new Structures.Fenwick(3);
    var s = new Structures.SegmentTree(3);
    long[] a = new long[3];
    for (var line : Files.readAllLines(file)) {
      if (line.startsWith("#") || line.isBlank()) continue;
      var fields = line.split("\\s+");
      int i = Integer.parseInt(fields[0]),
          l = Integer.parseInt(fields[2]),
          h = Integer.parseInt(fields[3]);
      long delta = Long.parseLong(fields[1]), want = Long.parseLong(fields[4]);
      a[i] += delta;
      f.add(i, delta);
      s.set(i, a[i]);
      assertEquals(want, f.sum(l, h));
      assertEquals(want, s.sum(l, h));
    }
    var r = new Random(7);
    a = new long[37];
    f = new Structures.Fenwick(a.length);
    s = new Structures.SegmentTree(a.length);
    for (int step = 0; step < 2000; step++) {
      int i = r.nextInt(a.length);
      long delta = r.nextInt(41) - 20;
      a[i] += delta;
      f.add(i, delta);
      s.set(i, a[i]);
      int l = r.nextInt(a.length + 1), h = l + r.nextInt(a.length - l + 1);
      long want = 0;
      for (int j = l; j < h; j++) want += a[j];
      assertEquals(want, f.sum(l, h));
      assertEquals(want, s.sum(l, h));
    }
  }

  @Test
  void trieUnionBloomAndBounds() {
    var trie = new Structures.Trie();
    for (String s : List.of("", "car", "cart", "cat")) trie.insert(s);
    assertTrue(trie.delete("car"));
    assertFalse(trie.has("car"));
    assertTrue(trie.has("cart"));
    assertEquals(List.of("cart", "cat"), trie.prefix("ca"));
    var d = new Structures.Disjoint(6);
    assertTrue(d.union(0, 1));
    assertTrue(d.union(2, 3));
    assertTrue(d.union(1, 2));
    assertFalse(d.union(0, 3));
    assertEquals(4, d.size(0));
    assertEquals(3, d.count());
    var b = new Structures.Bloom(128, 4);
    for (int i = 0; i < 100; i++) b.add("" + i);
    for (int i = 0; i < 100; i++) assertTrue(b.has("" + i));
    assertEquals(0, Structures.lower(new long[] {2, 2, 7}, 2));
    assertEquals(2, Structures.upper(new long[] {2, 2, 7}, 2));
  }

  @Test
  void streamsGraphsAndIntervals() {
    var r = new Random(9);
    for (int n = 1; n < 60; n++) {
      long[] a = new long[n];
      for (int i = 0; i < n; i++) a[i] = r.nextInt(100) - 50;
      int width = 1 + r.nextInt(n);
      var maxima = Structures.maxima(a, width);
      for (int i = 0; i < maxima.size(); i++) {
        long want = Long.MIN_VALUE;
        for (int j = i; j < i + width; j++) want = Math.max(want, a[j]);
        assertEquals(want, maxima.get(i));
      }
      int k = r.nextInt(n + 1);
      assertEquals(
          Arrays.stream(a).boxed().sorted(Comparator.reverseOrder()).limit(k).toList(),
          Structures.topK(a, k));
    }
    assertEquals(
        List.of(1L, 2L, 3L, 4L, 5L),
        Structures.mergeSorted(List.of(new long[] {1, 4}, new long[] {}, new long[] {2, 3, 5})));
    assertEquals(
        List.of(new Structures.Interval(0, 8), new Structures.Interval(9, 10)),
        Structures.mergeIntervals(
            List.of(
                new Structures.Interval(5, 8),
                new Structures.Interval(0, 5),
                new Structures.Interval(9, 10))));
    int[][] g = {{1, 2}, {3}, {3}, {}};
    assertEquals(4, Structures.traverse(g, 0, false).size());
    assertEquals(4, Structures.traverse(g, 0, true).size());
    assertEquals(3, Structures.topological(g).orElseThrow().getLast());
    assertTrue(Structures.topological(new int[][] {{1}, {0}}).isEmpty());
  }
}
