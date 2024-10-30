package fusex

import (
	"context"
	"io/fs"
	"sync/atomic"
	"syscall"
	"time"

	libfuse "github.com/presotto/fuse"
	libfs "github.com/presotto/fuse/fs"
	"github.com/yznts/zen/v3/logic"
	"github.com/yznts/zen/v3/slice"
)

// Dir is a wrapper around go-fuse's Inode struct,
// which is the base struct for all filesystem nodes.
// Dir is used to provide a common implementation for all directory nodes
// in the fuse filesystem.
type Dir struct {
	// filesystem pointer to access top-level instances/methods.
	fs *FS
	// Current node inode number.
	inode uint64
	// Current node path in the filesystem.
	path string
	// Children directory entries.
	children []libfuse.Dirent
}

// Attr is called by fuse to get the attributes of a directory.
func (d *Dir) Attr(ctx context.Context, a *libfuse.Attr) error {
	// Log and measure.
	defer flog(time.Now(), "getattr "+d.path)
	// Get the directory's attributes.
	stat, err := d.fs.backend.(fs.StatFS).Stat(d.path)
	if err != nil {
		return err
	}
	// Set the attributes to FUSE.
	a.Mode = stat.Mode()
	// a.Mode = os.ModeDir | 0o555
	a.Size = 4096
	a.Inode = d.inode
	a.Mtime = stat.ModTime()
	a.Atime = a.Mtime
	a.Ctime = a.Mtime
	// Return success.
	return nil
}

// Lookup is called by fuse to lookup a directory entry by name.
func (d *Dir) Lookup(ctx context.Context, name string) (libfs.Node, error) {
	// Log and measure.
	defer flog(time.Now(), "lookup "+d.path+"/"+name)
	// Get and correctly wrap directory entry as a node.
	for _, e := range d.children {
		// Found entry.
		if e.Name == name {
			// Correctly wrap the entry as a node.
			if e.Type == libfuse.DT_Dir {
				return NewDir(d.fs, d.path+"/"+name, e.Inode), nil
			} else {
				return NewFile(d.fs, d.path+"/"+name, e.Inode), nil
			}
		}
	}
	// If we are here, the entry was not found.
	return nil, syscall.ENOENT
}

// ReadDirAll is a method called by fuse to read the directory entries.
// We are reloading the directory children every time readdir is called
// to ensure that the directory is always up to date,
// and returning the directory entries.
//
// TODO: Probably we should somehow reuse existing children because of inode numbers.
func (d *Dir) ReadDirAll(ctx context.Context) ([]libfuse.Dirent, error) {
	// Log and measure.
	defer flog(time.Now(), "readdir "+d.path)
	// Read the directory entries.
	entries, err := d.fs.backend.(fs.ReadDirFS).ReadDir(d.path)
	if err != nil {
		panic(err)
	}
	// Transform them into libfuse directory entries
	// and store them as children (will be used in lookup).
	d.children = slice.Map(entries, func(e fs.DirEntry) libfuse.Dirent {
		atomic.AddUint64(&inodeCount, 1)
		return libfuse.Dirent{
			Name:  e.Name(),
			Type:  logic.Tr(e.IsDir(), libfuse.DT_Dir, libfuse.DT_File),
			Inode: inodeCount,
		}
	})
	// Return the directory entries.
	return d.children, nil
}

// NewDir creates a new directory node.
func NewDir(fs *FS, path string, ino uint64) *Dir {
	return &Dir{
		fs:    fs,
		path:  path,
		inode: ino,
	}
}
