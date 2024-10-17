package fusex

import (
	"context"
	"io/fs"
	"log"
	"syscall"
	"time"

	libfs "github.com/hanwen/go-fuse/v2/fs"
	libfuse "github.com/hanwen/go-fuse/v2/fuse"
)

// File is a wrapper around go-fuse's Inode struct,
// which is the base struct for all file system nodes.
// File is used to provide a common implementation for all file nodes
// in the FUSE file system.
type File struct {
	// Nest Inode, the go-fuse struct that Node wraps.
	libfs.Inode
	// Nest the remote file system.
	fs any
	// Nest the node's path in the file system.
	path string
}

func (f *File) Getattr(ctx context.Context, h libfs.FileHandle, out *libfuse.AttrOut) syscall.Errno {
	// Log and measure.
	now := time.Now()
	defer func() {
		log.Printf("FUSEX: getattr: %s, %s", f.path, time.Since(now))
	}()
	// Get the file's attributes.
	attr, err := f.fs.(fs.StatFS).Stat(f.path)
	if err != nil {
		return syscall.ENOENT
	}
	// Set the attributes to FUSE.
	out.Mode = uint32(attr.Mode())
	out.Size = uint64(attr.Size())
	out.Mtime = uint64(attr.ModTime().Unix())
	out.Atime = out.Mtime
	out.Ctime = out.Mtime
	// Return success.
	return 0
}

// NewFile creates a new file node.
func NewFile(fs any, path string) *File {
	return &File{
		fs:   fs,
		path: path,
	}
}
