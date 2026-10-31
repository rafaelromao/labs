package labs.streamio;

import java.io.IOException;
import java.io.InputStream;
import java.util.function.BooleanSupplier;

/** Bounded UTF-8, LF-delimited signed-decimal records. The caller owns the input. */
public final class RecordProcessor {
    public record Result(long count, long sum) {}

    private RecordProcessor() {}

    public static Result process(InputStream input, int maxBytes, BooleanSupplier cancelled) throws IOException {
        throw new UnsupportedOperationException("TODO: implement");
    }
}
