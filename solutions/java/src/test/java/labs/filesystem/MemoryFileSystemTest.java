package labs.filesystem;

import static org.junit.jupiter.api.Assertions.*;
import java.util.List;
import org.junit.jupiter.api.Test;

class MemoryFileSystemTest {
    @Test void copiesAndSubtreeMoves(){
        var fs=new MemoryFileSystem();fs.mkdir("/a");byte[] data={1,2};fs.write("/a/file",data);data[0]=9;
        assertArrayEquals(new byte[]{1,2},fs.read("/a/file"));var returned=fs.read("/a/file");returned[0]=8;
        assertThrows(IllegalArgumentException.class,()->fs.move("/a","/a/child"));fs.move("/a","/b");
        assertArrayEquals(new byte[]{1,2},fs.read("/b/file"));assertThrows(IllegalStateException.class,()->fs.read("/a/file"));
        assertEquals(List.of("b"),fs.list("/"));assertThrows(IllegalArgumentException.class,()->fs.write("/b/../escape",data));
    }
}
