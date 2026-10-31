// Package filesystem implements a local, synchronized in-memory file tree.
package filesystem

import ()

type FS struct{}

func New() *FS {
	panic("TODO: implement")
}

func (f *FS) Mkdir(p string) error {
	panic("TODO: implement")
}

// Write creates or replaces a file, copying bytes. Its parent must exist.
func (f *FS) Write(p string, data []byte) error {
	panic("TODO: implement")
}

func (f *FS) Read(p string) ([]byte, error) {
	panic("TODO: implement")
}

func (f *FS) List(p string) ([]string, error) {
	panic("TODO: implement")
}

func (f *FS) Move(from, to string) error {
	panic("TODO: implement")
}
