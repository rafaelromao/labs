package labs.filesystem;

import java.util.List;

/** An in-memory tree, not a host-filesystem sandbox or durable store. */
public final class MemoryFileSystem {
    public MemoryFileSystem() {
        // TODO: create the root directory "/".
    }

    public void mkdir(String path) {
        throw new UnsupportedOperationException("TODO: implement");
    }

    public void write(String path, byte[] data) {
        throw new UnsupportedOperationException("TODO: implement");
    }

    public byte[] read(String path) {
        throw new UnsupportedOperationException("TODO: implement");
    }

    public List<String> list(String path) {
        throw new UnsupportedOperationException("TODO: implement");
    }

    public void move(String from, String to) {
        throw new UnsupportedOperationException("TODO: implement");
    }
}
