package fusex

import (
	"context"
	"io/fs"
	"sync"
	"syscall"
	"time"

	libfuse "github.com/presotto/fuse"
	libfs "github.com/presotto/fuse/fs"
	"github.com/yznts/rfs/pkg/fsx"
)

// File is a wrapper around go-fuse's Inode struct,
// which is the base struct for all file system nodes.
// File is used to provide a common implementation for all file nodes
// in the fuse file system.
type File struct {
	// File system pointer to access top-level instances/methods.
	fs *FS
	// Mutex to lock the file.
	*sync.Mutex
	// Current node inode number.
	inode uint64
	// Current node path in the file system.
	path string
	// Current node file descriptor on open.
	file fsx.File
}

// Attr is called by fuse to get the attributes of a file.
func (f *File) Attr(ctx context.Context, a *libfuse.Attr) error {
	// Log and measure.
	defer flog(time.Now(), "getattr "+f.path)
	// Get the file's attributes.
	attr, err := f.fs.backend.(fs.StatFS).Stat(f.path)
	if err != nil {
		return err
	}
	// Set the attributes to FUSE.
	a.Mode = attr.Mode()
	a.Size = uint64(attr.Size())
	a.Inode = f.inode
	a.Mtime = attr.ModTime()
	a.Atime = a.Mtime
	a.Ctime = a.Mtime
	// Return success.
	return nil
}

// Setattr is called by fuse to set the attributes of a file.
func (f *File) Setattr(ctx context.Context, req *libfuse.SetattrRequest, resp *libfuse.SetattrResponse) error {
	// Log and measure.
	defer flog(time.Now(), "setattr "+f.path)
	// TODO: Implement the Setattr method.
	// ...
	// Return success.
	return nil
}

// Open is called by fuse to open a file descriptor.
func (f *File) Open(ctx context.Context, req *libfuse.OpenRequest, resp *libfuse.OpenResponse) (libfs.Handle, error) {
	// Log and measure.
	defer flog(time.Now(), "open "+f.path)
	// Unsupported flags (taken from sshfs-go).
	if req.Flags&libfuse.OpenAppend == libfuse.OpenAppend {
		return nil, syscall.ENOTSUP
	}
	// Open the file.
	var err error
	f.file, err = f.fs.backend.(fsx.OpenFileFS).OpenFile(f.path, int(req.Flags), 0)
	if err != nil {
		return nil, err
	}
	// Lock (will be unlocked in the Release method)
	f.Lock()
	// Return with managed flags.
	if req.Flags.IsReadOnly() {
		return f, err
	}
	if req.Flags.IsWriteOnly() {
		// Somehow the OpenTruncate flag is not actually truncating the file.
		// Let's do this manually.
		// resp.Flags = libfuse.OpenResponseFlags(libfuse.OpenWriteOnly | libfuse.OpenTruncate)
		f.file.Truncate(0)
		f.file.Seek(0, 0)
		resp.Flags = libfuse.OpenResponseFlags(libfuse.OpenWriteOnly)
		return f, nil
	}
	// Default return (cannot determine mode).
	return nil, syscall.ENOTSUP
}

// Read is called by fuse to read a file.
func (f *File) Read(ctx context.Context, req *libfuse.ReadRequest, resp *libfuse.ReadResponse) error {
	// Log and measure.
	defer flog(time.Now(), "read "+f.path)
	// Seek to the offset.
	f.file.Seek(req.Offset, 0)
	// Read the data.
	resp.Data = make([]byte, req.Size)
	f.file.Read(resp.Data)
	// Return success.
	return nil
}

// // Write is called by fuse to write a file.
func (f *File) Write(ctx context.Context, req *libfuse.WriteRequest, resp *libfuse.WriteResponse) error {
	// Log and measure.
	defer flog(time.Now(), "write "+f.path)
	// Seek to the offset.
	f.file.Seek(req.Offset, 0)
	// Write the data.
	_, err := f.file.Write(req.Data)
	resp.Size = len(req.Data)
	if err != nil {
		panic(err)
	}
	return err
}

// // Realease is called by fuse to release a file.
func (f *File) Release(ctx context.Context, req *libfuse.ReleaseRequest) error {
	// Log and measure.
	defer flog(time.Now(), "release "+f.path)
	// Unlock.
	defer f.Unlock()
	// Close the file descriptor.
	if f.file != nil {
		return f.file.Close()
	}
	return nil
}

// NewFile creates a new file node.
func NewFile(fs *FS, path string, ino uint64) *File {
	return &File{
		fs:    fs,
		path:  path,
		Mutex: &sync.Mutex{},
		inode: ino,
	}
}
