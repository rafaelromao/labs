package fencing

import "testing"

func TestOldOwnerRejectedAtResource(t *testing.T) {
	var r Resource
	if !r.Write(7, 0, "old") || !r.Write(8, 1, "new") {
		t.Fatal("setup")
	}
	if r.Write(7, 2, "stale") {
		t.Fatal("lease expiration did not fence")
	}
	if r.Write(8, 1, "duplicate") {
		t.Fatal("revision")
	}
	if value, revision := r.Read(); value != "new" || revision != 2 {
		t.Fatal(value, revision)
	}
}
