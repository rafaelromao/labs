// Package filesystem implements a local, synchronized in-memory file tree.
package filesystem

import (
	"bytes"
	"errors"
	"path"
	"sort"
	"strings"
	"sync"
)

type node struct {
	directory bool
	data      []byte
}
type FS struct {
	mu    sync.Mutex
	nodes map[string]node
}

func New() *FS { return &FS{nodes: map[string]node{"/": {directory: true}}} }
func normalize(p string) (string, error) {
	if !strings.HasPrefix(p, "/") || strings.ContainsRune(p, 0) {
		return "", errors.New("filesystem: absolute path required")
	}
	for _, part := range strings.Split(p, "/") {
		if part == "." || part == ".." {
			return "", errors.New("filesystem: dot components forbidden")
		}
	}
	return path.Clean(p), nil
}
func (f *FS) create(p string, n node) error {
	if _, ok := f.nodes[p]; ok {
		return errors.New("filesystem: already exists")
	}
	parent, ok := f.nodes[path.Dir(p)]
	if !ok || !parent.directory {
		return errors.New("filesystem: parent directory missing")
	}
	f.nodes[p] = n
	return nil
}
func (f *FS) Mkdir(p string) error {
	p, err := normalize(p)
	if err != nil {
		return err
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.create(p, node{directory: true})
}

// Write creates or replaces a file, copying bytes. Its parent must exist.
func (f *FS) Write(p string, data []byte) error {
	p, err := normalize(p)
	if err != nil {
		return err
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if n, ok := f.nodes[p]; ok {
		if n.directory {
			return errors.New("filesystem: is directory")
		}
		f.nodes[p] = node{data: bytes.Clone(data)}
		return nil
	}
	return f.create(p, node{data: bytes.Clone(data)})
}
func (f *FS) Read(p string) ([]byte, error) {
	p, err := normalize(p)
	if err != nil {
		return nil, err
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	n, ok := f.nodes[p]
	if !ok || n.directory {
		return nil, errors.New("filesystem: file missing")
	}
	return bytes.Clone(n.data), nil
}
func (f *FS) List(p string) ([]string, error) {
	p, err := normalize(p)
	if err != nil {
		return nil, err
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	n, ok := f.nodes[p]
	if !ok || !n.directory {
		return nil, errors.New("filesystem: directory missing")
	}
	result := []string{}
	for key := range f.nodes {
		if key != p && path.Dir(key) == p {
			result = append(result, path.Base(key))
		}
	}
	sort.Strings(result)
	return result, nil
}
func (f *FS) Move(from, to string) error {
	from, err := normalize(from)
	if err != nil {
		return err
	}
	to, err = normalize(to)
	if err != nil {
		return err
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if from == "/" {
		return errors.New("filesystem: cannot move root")
	}
	n, ok := f.nodes[from]
	if !ok {
		return errors.New("filesystem: source missing")
	}
	if from == to {
		return nil
	}
	if _, ok := f.nodes[to]; ok {
		return errors.New("filesystem: destination exists")
	}
	if n.directory && strings.HasPrefix(to, from+"/") {
		return errors.New("filesystem: move into descendant")
	}
	parent, ok := f.nodes[path.Dir(to)]
	if !ok || !parent.directory {
		return errors.New("filesystem: destination parent missing")
	}
	moved := map[string]node{}
	for key, value := range f.nodes {
		if key == from || strings.HasPrefix(key, from+"/") {
			moved[to+strings.TrimPrefix(key, from)] = value
		}
	}
	for key := range f.nodes {
		if key == from || strings.HasPrefix(key, from+"/") {
			delete(f.nodes, key)
		}
	}
	for key, value := range moved {
		f.nodes[key] = value
	}
	return nil
}
