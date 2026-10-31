package labs.streamio;

import static org.junit.jupiter.api.Assertions.*;

import java.io.ByteArrayInputStream;
import java.io.IOException;
import java.nio.charset.StandardCharsets;
import org.junit.jupiter.api.Test;

class RecordProcessorTest {
  @Test
  void shortReadsAndFailure() throws Exception {
    var input =
        new ByteArrayInputStream("1\r\n-2\n+3".getBytes(StandardCharsets.UTF_8)) {
          @Override
          public synchronized int read(byte[] bytes, int offset, int length) {
            return super.read(bytes, offset, Math.min(length, 1));
          }
        };
    assertEquals(new RecordProcessor.Result(3, 2), RecordProcessor.process(input, 4, () -> false));
    var failing =
        new java.io.InputStream() {
          int i;

          @Override
          public int read() throws IOException {
            if (i++ < 2) return i == 1 ? '1' : '\n';
            throw new IOException("injected");
          }
        };
    assertThrows(IOException.class, () -> RecordProcessor.process(failing, 4, () -> false));
  }

  private RecordProcessor.Result read(String text, int max) throws IOException {
    return RecordProcessor.process(
        new ByteArrayInputStream(text.getBytes(StandardCharsets.UTF_8)), max, () -> false);
  }

  @Test
  void validBoundaries() throws Exception {
    assertEquals(new RecordProcessor.Result(0, 0), read("", 4));
    assertEquals(new RecordProcessor.Result(3, 2), read("1\r\n-2\n+3", 4));
    assertEquals(new RecordProcessor.Result(1, 123), read("123\n", 3));
  }

  @Test
  void rejectsBadRecords() {
    for (String text : new String[] {"1234", "\n", "1\nno", "1 2", "+"})
      assertThrows(IOException.class, () -> read(text, 3));
    assertThrows(IOException.class, () -> read("9223372036854775807\n1", 20));
    assertThrows(IOException.class, () -> read("-9223372036854775808\n-1", 20));
    assertThrows(
        IOException.class,
        () ->
            RecordProcessor.process(
                new ByteArrayInputStream(new byte[] {(byte) 0xff}), 3, () -> false));
  }

  @Test
  void cancelledBeforeRead() {
    assertThrows(
        IOException.class,
        () -> RecordProcessor.process(new ByteArrayInputStream(new byte[0]), 3, () -> true));
  }
}
