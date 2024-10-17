package fusex

import (
	"context"
	"io/fs"

	libfs "github.com/hanwen/go-fuse/v2/fs"
	libfuse "github.com/hanwen/go-fuse/v2/fuse"
)

// Dir is a wrapper around go-fuse's Inode struct,
// which is the base struct for all file system nodes.
// Dir is used to provide a common implementation for all directory nodes
// in the FUSE file system.
type Dir struct {
	// Nest Inode, the go-fuse struct that Node wraps.
	libfs.Inode
	// Nest the remote file system.
	fs any
	// Nest the node's path in the file system.
	path string
}

// OnAdd is a method that is called when a directory node is added to the file system.
func (d *Dir) OnAdd(ctx context.Context) {
	// Read the directory entries.
	entries, err := d.fs.(fs.ReadDirFS).ReadDir(d.path)
	if err != nil {
		panic(err)
	}
	// Add entries to the directory node.
	for _, entry := range entries {
		// Create a new node, depending on the entry type.
		var node *libfs.Inode
		if entry.IsDir() {
			node = d.NewInode(ctx, NewDir(d.fs, d.path+"/"+entry.Name()), libfs.StableAttr{
				Mode: libfuse.S_IFDIR | 0755,
			})
		} else {
			node = d.NewInode(ctx, NewFile(d.fs, d.path+"/"+entry.Name()), libfs.StableAttr{})
		}
		// Add the node to the directory.
		if success := d.AddChild(entry.Name(), node, true); !success {
			panic("failed to add child")
		}
	}
}

// NewDir creates a new directory node.
func NewDir(fs any, path string) *Dir {
	return &Dir{
		fs:   fs,
		path: path,
	}
}
