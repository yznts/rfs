package fusex

import (
	"sync/atomic"

	libfs "github.com/presotto/fuse/fs"
)

var inodeCount uint64

// FS is a fuse filesystem.
// We are using it as a wrapper around the underlying filesystem.
type FS struct {
	// Underlying filesystem.
	backend any
	// Inode counter, used to provide unique inode numbers.
	inoc uint64
}

// GetInoc returns the next unique inode number.
// Used to provide unique inode numbers for all nodes in the filesystem.
func (f *FS) GetInoc() uint64 {
	atomic.AddUint64(&f.inoc, 1)
	return f.inoc
}

// Root is called by fuse to get the root node of the filesystem.
func (f *FS) Root() (libfs.Node, error) {
	return NewDir(f, ".", f.GetInoc()), nil
}

// NewFS creates a new fuse filesystem
// with the given backend filesystem.
func NewFS(fs any) *FS {
	return &FS{
		backend: fs,
		inoc:    0,
	}
}
