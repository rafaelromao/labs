package labs.streamio;

import java.io.BufferedInputStream;
import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.nio.ByteBuffer;
import java.nio.charset.CodingErrorAction;
import java.nio.charset.StandardCharsets;
import java.util.function.BooleanSupplier;

/** Bounded UTF-8, LF-delimited signed-decimal records. The caller owns input. */
public final class RecordProcessor {
    public record Result(long count, long sum) {}
    private RecordProcessor() {}

    public static Result process(InputStream input, int maxBytes, BooleanSupplier cancelled) throws IOException {
        if (maxBytes <= 0) throw new IllegalArgumentException("positive record limit required");
        var reader = new BufferedInputStream(input);
        var line = new ByteArrayOutputStream(Math.min(maxBytes, 4096));
        Result result = new Result(0, 0);
        while (true) {
            if (cancelled.getAsBoolean()) throw new IOException("cancelled");
            int b = reader.read();
            if (b == -1) return line.size() == 0 ? result : consume(line, result);
            if (b == '\n') { result = consume(line, result); line.reset(); }
            else {
                if (line.size() == maxBytes) throw new IOException("record too large");
                line.write(b);
            }
        }
    }

    private static Result consume(ByteArrayOutputStream line, Result result) throws IOException {
        String text = StandardCharsets.UTF_8.newDecoder()
                .onMalformedInput(CodingErrorAction.REPORT).onUnmappableCharacter(CodingErrorAction.REPORT)
                .decode(ByteBuffer.wrap(line.toByteArray())).toString();
        if (text.endsWith("\r")) text = text.substring(0, text.length() - 1);
        if (!text.matches("[+-]?[0-9]+")) throw new IOException("decimal record required");
        try {
            long value = Long.parseLong(text);
            return new Result(Math.addExact(result.count(), 1), Math.addExact(result.sum(), value));
        } catch (NumberFormatException | ArithmeticException e) { throw new IOException("numeric overflow", e); }
    }
}
