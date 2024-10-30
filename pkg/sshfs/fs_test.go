package sshfs_test

import (
	"net/url"
	"os"
	"testing"

	"github.com/pkg/sftp"
	"github.com/yznts/rfs/pkg/sshfs"
	"github.com/yznts/rfs/pkg/testfs"
	"github.com/yznts/zen/v3/errx"
	"github.com/yznts/zen/v3/logic"
	"golang.org/x/crypto/ssh"
)

var (
	_fs *sshfs.FS
)

func TestFS(t *testing.T) {
	// Get remote filesystem url.
	url := errx.Must(url.Parse(os.Getenv("TEST_SSHFS_DSN")))
	// Get private key.
	key := errx.Must(os.ReadFile(
		logic.Or(os.Getenv("TEST_SSHFS_KEY"), os.Getenv("HOME")+"/.ssh/id_rsa"),
	))
	// Open ssh connection.
	_ssh := errx.Must(ssh.Dial("tcp", url.Host, &ssh.ClientConfig{
		User:            url.User.Username(),
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(errx.Must(ssh.ParsePrivateKey(key))),
		},
	}))
	// Open sftp connection.
	_sftp := errx.Must(sftp.NewClient(_ssh))
	// Create filesystem instance.
	_fs = &sshfs.FS{_sftp}
	// Run tests
	t.Run("TestOpen", func(t *testing.T) {
		testfs.RunTestOpen(t, _fs)
	})
	t.Run("TestOpenFile", func(t *testing.T) {
		testfs.RunTestOpenFile(t, _fs)
	})
	t.Run("TestReadFile", func(t *testing.T) {
		testfs.RunTestReadFile(t, _fs)
	})
	t.Run("TestWriteFile", func(t *testing.T) {
		testfs.RunTestWriteFile(t, _fs)
	})
	t.Run("TestReadDir", func(t *testing.T) {
		testfs.RunTestReadDir(t, _fs)
	})
}
