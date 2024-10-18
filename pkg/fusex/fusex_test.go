package fusex_test

import (
	"testing"

	libfuse "github.com/presotto/fuse"
	libfs "github.com/presotto/fuse/fs"
	"github.com/yznts/rfs/pkg/fusex"
	"github.com/yznts/rfs/pkg/sshfs"
	"github.com/yznts/rfs/pkg/testfs"
)

func TestFusex(t *testing.T) {
	// Use the ssh connection from testfs to initialize the file system.
	_sshfs := sshfs.NewFS(testfs.REMOTE_SSHC)
	// Create a root fuse node using fusex.
	root := fusex.NewFS(_sshfs)
	// Mount the file system.
	c, err := libfuse.Mount("/tmp/rfs", libfuse.FSName("rfs"), libfuse.Subtype("rfs"))
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	// Serve the file system.
	if err := libfs.Serve(c, root); err != nil {
		t.Fatal(err)
	}
}
