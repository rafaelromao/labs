package labs.httpio;

import java.io.*;
import java.net.*;
import java.net.http.*;
import java.nio.*;
import java.time.*;
import java.util.*;
import java.util.concurrent.*;
import java.util.concurrent.Flow;

/** Owner reuses HttpClient. The subscriber enforces the limit before body accumulation. */
public final class BoundedClient {
  private BoundedClient() {}

  private static final class Body implements HttpResponse.BodySubscriber<byte[]> {
    private final int limit;
    private final CompletableFuture<byte[]> result = new CompletableFuture<>();
    private final ByteArrayOutputStream bytes = new ByteArrayOutputStream();
    private Flow.Subscription subscription;

    Body(int limit) {
      this.limit = limit;
    }

    public CompletionStage<byte[]> getBody() {
      return result;
    }

    public void onSubscribe(Flow.Subscription s) {
      subscription = s;
      s.request(1);
    }

    public void onNext(List<ByteBuffer> buffers) {
      if (result.isDone()) return;
      for (var b : buffers) {
        if (b.remaining() > limit - bytes.size()) {
          subscription.cancel();
          result.completeExceptionally(new IOException("body limit"));
          return;
        }
        byte[] part = new byte[b.remaining()];
        b.get(part);
        bytes.writeBytes(part);
      }
      subscription.request(1);
    }

    public void onError(Throwable t) {
      result.completeExceptionally(t);
    }

    public void onComplete() {
      result.complete(bytes.toByteArray());
    }
  }

  public static byte[] fetch(HttpClient client, URI uri, int limit, Duration timeout)
      throws Exception {
    Objects.requireNonNull(client);
    if (limit <= 0 || limit > 1 << 20 || timeout.isNegative() || timeout.isZero())
      throw new IllegalArgumentException();
    var request = HttpRequest.newBuilder(uri).timeout(timeout).GET().build();
    var future = client.sendAsync(request, info -> new Body(limit));
    try {
      var response = future.get(timeout.toNanos(), TimeUnit.NANOSECONDS);
      if (response.statusCode() < 200 || response.statusCode() >= 300)
        throw new IOException("HTTP status rejected");
      return response.body();
    } catch (InterruptedException e) {
      Thread.currentThread().interrupt();
      throw e;
    } finally {
      if (!future.isDone()) future.cancel(true);
    }
  }
}
