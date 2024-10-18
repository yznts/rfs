package fusex

import (
	"sync/atomic"

	libfs "github.com/presotto/fuse/fs"
)

var inodeCount uint64

// FS is a fuse file system.
// We are using it as a wrapper around the underlying file system.
type FS struct {
	// Underlying file system.
	backend any
	// Inode counter, used to provide unique inode numbers.
	inoc uint64
}

// GetInoc returns the next unique inode number.
// Used to provide unique inode numbers for all nodes in the file system.
func (f *FS) GetInoc() uint64 {
	atomic.AddUint64(&f.inoc, 1)
	return f.inoc
}

// Root is called by fuse to get the root node of the file system.
func (f *FS) Root() (libfs.Node, error) {
	return NewDir(f, ".", f.GetInoc()), nil
}

// NewFS creates a new fuse file system
// with the given backend file system.
func NewFS(fs any) *FS {
	return &FS{
		backend: fs,
		inoc:    0,
	}
}
