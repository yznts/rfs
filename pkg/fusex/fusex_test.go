package fusex_test

import (
	"os"
	"testing"

	libfs "github.com/hanwen/go-fuse/v2/fs"
	"github.com/yznts/rfs/pkg/fusex"
	"github.com/yznts/rfs/pkg/sshfs"
	"github.com/yznts/rfs/pkg/testfs"
)

func TestFusex(t *testing.T) {
	// Use the ssh connection from testfs to initialize the file system.
	_sshfs := sshfs.NewFS(testfs.REMOTE_SSHC)
	// Create a root fuse node using fusex.
	root := fusex.NewDir(_sshfs, ".")
	// Mount the file system.
	wd, _ := os.Getwd()
	t.Log("Mounting the file system at ", wd+"/tmp")
	options := &libfs.Options{}
	options.Debug = true
	server, err := libfs.Mount("./tmp", root, options)
	if err != nil {
		t.Fatal(err)
	}
	server.Wait()
}
