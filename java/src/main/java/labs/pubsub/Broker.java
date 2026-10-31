package labs.pubsub;

import java.util.Optional;

/** Bounded local broadcast, with drop-newest per full subscriber queue. */
public final class Broker<T> {
    public record Delivery(int delivered, int dropped) {}

    public final class Subscription {
        private Subscription() {}

        public Optional<T> poll() {
            throw new UnsupportedOperationException("TODO: implement");
        }

        public boolean isClosed() {
            throw new UnsupportedOperationException("TODO: implement");
        }

        public void unsubscribe() {
            throw new UnsupportedOperationException("TODO: implement");
        }
    }

    public Broker(int limit, int capacity) {
        // TODO: validate and store the limits.
    }

    public Subscription subscribe() {
        throw new UnsupportedOperationException("TODO: implement");
    }

    public Delivery publish(T value) {
        throw new UnsupportedOperationException("TODO: implement");
    }

    public void close() {
        throw new UnsupportedOperationException("TODO: implement");
    }
}
