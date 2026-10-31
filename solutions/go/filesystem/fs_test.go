package filesystem

import "testing"

func TestCopiesAndAtomicSubtreeMove(t *testing.T) {
	f := New()
	if err := f.Mkdir("/a"); err != nil {
		t.Fatal(err)
	}
	data := []byte("hello")
	if err := f.Write("/a/file", data); err != nil {
		t.Fatal(err)
	}
	data[0] = 'x'
	got, _ := f.Read("/a/file")
	if string(got) != "hello" {
		t.Fatal(string(got))
	}
	got[0] = 'x'
	if err := f.Move("/a", "/a/child"); err == nil {
		t.Fatal("cycle accepted")
	}
	if err := f.Move("/a", "/b"); err != nil {
		t.Fatal(err)
	}
	got, err := f.Read("/b/file")
	if err != nil || string(got) != "hello" {
		t.Fatal(string(got), err)
	}
	if _, err := f.Read("/a/file"); err == nil {
		t.Fatal("old subtree remains")
	}
	names, err := f.List("/")
	if err != nil || len(names) != 1 || names[0] != "b" {
		t.Fatal(names, err)
	}
	if err := f.Write("/b/../escape", nil); err == nil {
		t.Fatal("dot traversal accepted")
	}
}
