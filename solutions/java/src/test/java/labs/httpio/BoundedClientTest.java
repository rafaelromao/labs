package labs.httpio;

import static org.junit.jupiter.api.Assertions.*;

import com.sun.net.httpserver.HttpServer;
import java.net.*;
import java.net.http.*;
import java.time.*;
import java.util.concurrent.*;
import org.junit.jupiter.api.Test;

class BoundedClientTest {
  @Test
  void boundsAndDeadline() throws Exception {
    var release = new CountDownLatch(1);
    var server = HttpServer.create(new InetSocketAddress("127.0.0.1", 0), 0);
    var executor = Executors.newFixedThreadPool(2);
    server.setExecutor(executor);
    server.createContext(
        "/",
        exchange -> {
          try {
            exchange.sendResponseHeaders(200, 0);
            exchange.getResponseBody().write(new byte[] {'1', '\n'});
            exchange.getResponseBody().flush();
            exchange.getResponseBody().write(new byte[] {'2', '\n'});
          } finally {
            exchange.close();
          }
        });
    server.createContext(
        "/blocked",
        exchange -> {
          try {
            release.await(2, TimeUnit.SECONDS);
          } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
          } finally {
            exchange.close();
          }
        });
    server.start();
    try (var client = HttpClient.newHttpClient()) {
      var uri = URI.create("http://127.0.0.1:" + server.getAddress().getPort());
      assertArrayEquals(
          new byte[] {'1', '\n', '2', '\n'},
          BoundedClient.fetch(client, uri, 4, Duration.ofSeconds(2)));
      assertThrows(
          Exception.class, () -> BoundedClient.fetch(client, uri, 3, Duration.ofSeconds(2)));
      assertThrows(
          Exception.class,
          () -> BoundedClient.fetch(client, uri.resolve("/blocked"), 4, Duration.ofMillis(30)));
    } finally {
      release.countDown();
      server.stop(0);
      executor.shutdownNow();
      assertTrue(executor.awaitTermination(3, TimeUnit.SECONDS));
    }
  }
}
