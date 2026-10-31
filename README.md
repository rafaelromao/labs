# Rafael Romão's Coding Labs

Twenty-two small, production-style backend components to implement yourself, each with a written
contract, a diagram, a starter stub in **Java** and in **Go**, a test suite that checks the contract,
and a reference solution to compare against afterwards.

**Lab pages:** <https://rafaelromao.github.io/labs/> — each lab's task, diagram, API, rules, worked example, steps and checklist.

## What is in the repository

| Path | Contents |
| --- | --- |
| [`java/`](java) | IntelliJ IDEA project (Gradle, Kotlin DSL, JDK 21). Stubs in `src/main/java/labs/<lab>`, tests in `src/test/java/labs/<lab>`. |
| [`go/`](go) | IntelliJ IDEA / GoLand project (Go module). One package per lab: stub files plus `*_test.go`. |
| [`solutions/`](solutions) | Reference implementations for both languages, with the full test suites and some supporting examples. |
| [`site/`](site) | The lab pages, published to GitHub Pages by [`.github/workflows/pages.yml`](.github/workflows/pages.yml). |
| [`fixtures/`](fixtures), [`compose.yaml`](compose.yaml) | Shared test fixtures and a disposable PostGIS for the optional integration tests. |
| [`labs.json`](labs.json) | The lab catalogue: id, slug, title, page and starter locations. |

Every stub keeps the public API the tests need — class and method signatures, records, enums,
exceptions — and every method body is `throw new UnsupportedOperationException("TODO: implement")`
in Java or `panic("TODO: implement")` in Go. Out of the box, every lab test fails; your job is to
make them pass.

Two helpers ship fully implemented because other labs depend on them: `labs.algorithms.Structures`
(Java) / `go/algorithms` (indexed heap, trie, union-find, top-K, Fenwick and segment trees, Bloom
filter and more), and the haversine `distance` in the nearby-search lab, which the geofence and
trajectory labs reuse.

## Getting started

### Java

1. Install a JDK 17 or newer to run Gradle. The build uses a JDK 21 toolchain and downloads it if
   none is installed.
2. In IntelliJ IDEA, **File → Open…** and choose the `java/` directory. It imports as a Gradle
   project.
3. Run a single lab's tests from the gutter icon in its test class, or run the shared **All labs**
   configuration.

From a terminal:

```sh
cd java
./gradlew test                                # every lab
./gradlew test --tests 'labs.ttlcache.*'      # one lab
```

### Go

1. Install Go 1.25 or newer, and a C compiler: the SQLite labs (L10, L11) use CGO.
2. Open the `go/` directory in GoLand, or in IntelliJ IDEA with the Go plugin.
3. Run a package's tests from the gutter icon, or run the shared **All labs** configuration
   (`go test -race ./...`).

From a terminal:

```sh
cd go
go test -race -count=1 ./...                  # every lab
go test -race -count=1 ./ttlcache             # one lab
```

Start with **L05 (TTL cache)** if you are new to the set. It combines capacity, recency, time and
shared state in one small class, and its tests show the shape the other labs expect: an injected
clock, an inclusive deadline and a concurrent schedule.

## How to work a lab

1. **Read the contract and the tests.** Open the lab page and the test file. List inputs, outputs,
   boundaries, invariants, failure behaviour and resource limits before writing code.
2. **Implement the starter.** Replace every `TODO`. Run the lab's tests and fix what fails before
   reading anything else.
3. **Compare.** Open the reference in `solutions/` and explain one difference in representation or
   synchronization, and whether your version still keeps the contract.
4. **Change something.** Add a test for an edge case you chose, or implement the variation
   suggested on the lab page.

Tests you write yourself mostly confirm what you already believed; a supplied suite fails on the
boundary you skipped, the tie you broke differently, or the interleaving you never imagined. That
is why the order above matters more than the code.

## The labs

### L01 · Bounded Blocking Queue

Build the hand-off buffer that sits between threads producing work and threads consuming it. The queue holds at most `capacity` items: `put` blocks while it is full, `take` blocks while it is empty, and `close()` stops new puts while letting consumers finish the items already accepted before they see “end of stream”.

[Lab page](https://rafaelromao.github.io/labs/lab/l01-blockingqueue.html) · Java: [`BlockingQueue.java`](java/src/main/java/labs/blockingqueue/BlockingQueue.java) · Go: [`go/blockingqueue`](go/blockingqueue)

### L02 · Bounded Worker Pool

Build an executor that runs submitted tasks on a fixed number of worker threads, with a fixed-size waiting line in front of them. When every worker is busy and the line is full, `submit` rejects the task immediately instead of blocking the caller or growing memory. `close()` stops new submissions but lets every accepted task finish; `abort()` throws away queued work and interrupts running work.

[Lab page](https://rafaelromao.github.io/labs/lab/l02-workerpool.html) · Java: [`WorkerPool.java`](java/src/main/java/labs/workerpool/WorkerPool.java) · Go: [`go/workerpool`](go/workerpool)

### L03 · Concurrent Request Aggregator

Build a helper that runs a set of named calls concurrently — at most `concurrency` at a time — waits for all of them, and returns one outcome per name: either the value it returned or the exception it threw. A failing call must not cancel or hide its siblings (“all-settled”, as opposed to fail-fast).

[Lab page](https://rafaelromao.github.io/labs/lab/l03-fanout.html) · Java: [`RequestAggregator.java`](java/src/main/java/labs/fanout/RequestAggregator.java) · Go: [`go/fanout`](go/fanout)

### L04 · Per-Key Request Coalescing

Build a request coalescer (“single flight”): when many callers ask for the same key while a load for that key is already running, they wait on that one load instead of starting their own. When the load finishes — successfully or not — the key is forgotten, so the next request starts a fresh load. Each caller gets its own future, so one caller giving up never cancels the load for the others.

[Lab page](https://rafaelromao.github.io/labs/lab/l04-singleflight.html) · Java: [`RequestCoalescer.java`](java/src/main/java/labs/singleflight/RequestCoalescer.java) · Go: [`go/singleflight`](go/singleflight)

### L05 · Bounded TTL Cache

Build an in-memory key-value cache with two limits. Size: at most `capacity` entries; inserting into a full cache evicts the least-recently-used entry. Time: every write gives the entry a lifetime of `ttl`, after which it is never returned again. Reads refresh recency but not lifetime. Time comes from an injected clock so tests can move it deterministically.

[Lab page](https://rafaelromao.github.io/labs/lab/l05-ttlcache.html) · Java: [`TtlCache.java`](java/src/main/java/labs/ttlcache/TtlCache.java) · Go: [`go/ttlcache`](go/ttlcache)

### L06 · LFU Cache

Build a fixed-capacity cache that, when full, evicts the key that has been used the *fewest* times; if several keys tie on that count, evict the one among them used least recently. Every hit (read or overwrite) increases a key's count. Both `get` and `put` should run in O(1) — no scanning all entries to find the victim.

[Lab page](https://rafaelromao.github.io/labs/lab/l06-lfucache.html) · Java: [`LfuCache.java`](java/src/main/java/labs/lfucache/LfuCache.java) · Go: [`go/lfucache`](go/lfucache)

### L07 · Per-Client Rate Limiter

Build a per-client rate limiter using token buckets. Each client has a bucket that holds up to `burst` tokens and refills continuously at `rate` tokens per second; each request spends one token or is rejected immediately (no waiting). The limiter also caps how many clients it tracks at once and forgets clients that have been idle for a configured time.

[Lab page](https://rafaelromao.github.io/labs/lab/l07-ratelimiter.html) · Java: [`RateLimiter.java`](java/src/main/java/labs/ratelimiter/RateLimiter.java) · Go: [`go/ratelimiter`](go/ratelimiter)

### L08 · Delayed Task Scheduler

Build a timer queue. Callers register a callback under an id with a due time, can replace or cancel it while it is still pending, and a driver calls `runDue(now)` to execute every callback whose due time has arrived, earliest first. Scheduling an id that is already pending replaces it, so there is at most one pending task per id.

[Lab page](https://rafaelromao.github.io/labs/lab/l08-scheduler.html) · Java: [`DelayedScheduler.java`](java/src/main/java/labs/scheduler/DelayedScheduler.java) · Go: [`go/scheduler`](go/scheduler)

### L09 · Retry and Circuit Breaker Executor

Build an executor that protects calls to a flaky dependency in two ways. **Retries**: if an attempt fails with an exception you declare retryable, wait (backoff) and try again, up to `attempts` times. **Circuit breaker**: after `threshold` consecutive failed attempts, stop calling the dependency for `cooldown`, failing fast with `OpenCircuit`; after the cooldown, let exactly one probe call through to decide whether to close the circuit again.

[Lab page](https://rafaelromao.github.io/labs/lab/l09-resilience.html) · Java: [`RetryBreaker.java`](java/src/main/java/labs/resilience/RetryBreaker.java) · Go: [`go/resilience`](go/resilience)

### L10 · Durable Idempotent Operation

Build a retry-safe `add(scope, key, amount)` on a counter stored in SQLite. The first call with a given (scope, key) increments the counter *and* stores the resulting value as that request's response, in the same transaction. A repeat with the same key and amount returns the stored response without incrementing again; the same key with a different amount is rejected as a conflict.

[Lab page](https://rafaelromao.github.io/labs/lab/l10-idempotency.html) · Java: [`DurableCounter.java`](java/src/main/java/labs/idempotency/DurableCounter.java) · Go: [`go/idempotency`](go/idempotency)

### L11 · Outbox and Idempotent Consumer

Build both halves of the transactional outbox pattern on SQLite. The producer saves a business command *and* the event announcing it in one transaction. A relay API reads unsent events in batches and marks each one sent after delivery. The consumer applies each event's effect at most once by recording the event id in an inbox table in the same transaction as the effect, so a redelivered event is detected and skipped.

[Lab page](https://rafaelromao.github.io/labs/lab/l11-outbox.html) · Java: [`OutboxStore.java`](java/src/main/java/labs/outbox/OutboxStore.java) · Go: [`go/outbox`](go/outbox)

### L12 · Dependency Task Executor

Build an executor for tasks that depend on other tasks. Each task lists the tasks that must succeed before it may run. First validate the whole graph — unknown dependency, duplicate edge or cycle means reject before running anything. Then execute in waves: every task whose prerequisites are done runs in parallel (at most `concurrency` at once), and a task whose prerequisite failed or was skipped is marked skipped instead of run.

[Lab page](https://rafaelromao.github.io/labs/lab/l12-dag.html) · Java: [`DependencyExecutor.java`](java/src/main/java/labs/dag/DependencyExecutor.java) · Go: [`go/dag`](go/dag)

### L13 · In-Memory File System

Build an in-memory file system addressed by absolute paths such as `/docs/a.txt`. It supports creating directories, writing and reading files as byte arrays, listing a directory's children in sorted order, and moving a file or a whole directory subtree to a new path in one atomic step.

[Lab page](https://rafaelromao.github.io/labs/lab/l13-filesystem.html) · Java: [`MemoryFileSystem.java`](java/src/main/java/labs/filesystem/MemoryFileSystem.java) · Go: [`go/filesystem`](go/filesystem)

### L14 · Extending a Concurrent Reservation Component

Build a booking calendar for a resource that has `capacity` parallel units, for example three charging bays. A booking reserves some units over a half-open time interval `[start, end)` and is accepted only if, at every instant, the units of all overlapping bookings stay within capacity. The base contract is one-unit `book` plus cancel and expiry; the extension adds multi-unit `bookUnits` and `reschedule` without changing how `book` behaves.

[Lab page](https://rafaelromao.github.io/labs/lab/l14-reservation.html) · Java: [`ReservationCalendar.java`](java/src/main/java/labs/reservation/ReservationCalendar.java) · Go: [`go/reservation`](go/reservation)

### L15 · Bounded Publish-Subscribe Component

Build an in-process broadcast broker. Every subscriber gets its own bounded mailbox, and `publish(message)` offers the message to each current subscriber's mailbox without ever waiting. If a subscriber's mailbox is full, that subscriber misses this message (it is dropped and counted) while everyone else still receives it. The broker also caps how many subscribers it accepts.

[Lab page](https://rafaelromao.github.io/labs/lab/l15-pubsub.html) · Java: [`Broker.java`](java/src/main/java/labs/pubsub/Broker.java) · Go: [`go/pubsub`](go/pubsub)

### L16 · Time-Indexed Key-Value Store

Build a versioned key-value store. `put(key, timestamp, value)` records what a key's value was from a given time, and `get(key, at)` answers “what was the value as of time `at`?” — the value with the greatest timestamp that is `<= at`. Writes may arrive in any time order.

[Lab page](https://rafaelromao.github.io/labs/lab/l16-temporalstore.html) · Java: [`TemporalStore.java`](java/src/main/java/labs/temporalstore/TemporalStore.java) · Go: [`go/temporalstore`](go/temporalstore)

### L17 · Time-Window Event Aggregator

Build a rolling-window aggregator over timestamped integer events. You add events as time moves forward, and `query(now, k)` returns the count, the maximum and the `k` largest values among events from the last `width` time units. Time is driven only by the timestamps you pass in: each call moves a *watermark* forward, and an event older than the watermark is rejected as late.

[Lab page](https://rafaelromao.github.io/labs/lab/l17-windowmetrics.html) · Java: [`EventWindow.java`](java/src/main/java/labs/windowmetrics/EventWindow.java) · Go: [`go/windowmetrics`](go/windowmetrics)

### L18 · Streaming Record Processor

Build a function that reads a stream of newline-separated integers — one per line — and returns how many records there were and their sum. It reads through a buffer with memory bounded by the longest allowed line, never loading the whole input. It is strict: any malformed line, oversized line or numeric overflow fails the whole run with an `IOException`, and it never closes the caller's stream.

[Lab page](https://rafaelromao.github.io/labs/lab/l18-streamio.html) · Java: [`RecordProcessor.java`](java/src/main/java/labs/streamio/RecordProcessor.java) · Go: [`go/streamio`](go/streamio)

### L19 · Nearby-Object Search

Build an in-memory index of moving objects (id → latitude/longitude) that answers two questions: which objects are within R metres of a point, sorted by distance, and which `k` objects are nearest. Distances use the haversine formula on a sphere. To avoid checking every object, the index groups points into 1-degree latitude bands, and a radius query only examines bands that could possibly contain a match.

[Lab page](https://rafaelromao.github.io/labs/lab/l19-nearby.html) · Java: [`NearbyIndex.java`](java/src/main/java/labs/nearby/NearbyIndex.java) · Go: [`go/nearby`](go/nearby)

### L20 · Geofencing Engine

Build a geofence monitor for one circular zone. Feed it position updates per object; for each update it reports whether the object is now inside and whether that is a change (an enter or exit event). A hysteresis margin stops an object hovering at the boundary from flapping in and out, and per-object sequence numbers make it ignore updates that arrive out of order.

[Lab page](https://rafaelromao.github.io/labs/lab/l20-geofence.html) · Java: [`GeofenceEngine.java`](java/src/main/java/labs/geofence/GeofenceEngine.java) · Go: [`go/geofence`](go/geofence)

### L21 · Route Planner

Build a shortest-path finder for a directed graph with non-negative edge costs: return the cheapest node sequence from `start` to `goal` and its total cost. Implement Dijkstra first; then accept an optional heuristic `h(node)`, a never-too-high estimate of the remaining cost, which turns the same loop into A*. The extension, `TurnPlanner`, supports forbidden turns (“no left from road X onto road Y”) by searching over roads instead of junctions.

[Lab page](https://rafaelromao.github.io/labs/lab/l21-routing.html) · Java: [`RoutePlanner.java`](java/src/main/java/labs/routing/RoutePlanner.java), [`TurnPlanner.java`](java/src/main/java/labs/routing/TurnPlanner.java) · Go: [`go/routing`](go/routing)

### L22 · Trajectory Processor

Build a per-vehicle location history. Store timestamped GPS points (sorted by time even when they arrive out of order, keeping only the newest `capacity`), turn each pair of consecutive points into a segment with distance and speed, and flag segments whose speed is physically implausible — a GPS glitch. Separately, snap a point onto the nearest road segment and report whether the choice was ambiguous.

[Lab page](https://rafaelromao.github.io/labs/lab/l22-trajectory.html) · Java: [`SegmentMatcher.java`](java/src/main/java/labs/trajectory/SegmentMatcher.java), [`Trajectory.java`](java/src/main/java/labs/trajectory/Trajectory.java) · Go: [`go/trajectory`](go/trajectory)

## Reference solutions and verification

```sh
cd solutions/java && ./gradlew test           # reference Java suite
cd solutions/go && go test -race ./...        # reference Go suite
scripts/verify-starters.sh                    # starter tests pass once the reference replaces the stubs
```

CI ([`.github/workflows/ci.yml`](.github/workflows/ci.yml)) runs all three on every push.

A few Go tests in `solutions/go` inspect internal fields, such as a lock or a bucket map. The
starter versions check the same behaviour through the public API so that your design stays your
own.

`solutions/` also holds supporting examples that are not labs: exact planar polygon predicates and
a quadtree (`spatial`), fencing tokens (`fencing`), H3/S2 cell round trips (`cells`), bounded HTTP
bodies with deadlines (`httpio`), and a PostgreSQL/PostGIS adapter (`postgres`). The PostgreSQL
tests are integration-only:

```sh
docker compose up -d
cd solutions/java && ./gradlew integrationTest
```

## Where the tests stop

A green suite is evidence about the cases someone wrote, not proof over every schedule or failure.
The components are local and in memory unless a lab says otherwise; clocks are injected, so real
wall-clock drift is untested; and mutable values are not deep-copied. When you reason about
production readiness, name the limit you would test next rather than claiming the suite covers it.
