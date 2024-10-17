package testfs

import (
	"io"
	"io/fs"
	"net/url"
	"os"
	"testing"

	"github.com/yznts/rfs/pkg/fsx"
	"github.com/yznts/zen/v3/errx"
	"golang.org/x/crypto/ssh"
)

// This test suite requires a remote system with ssh access,
// even for non-ssh filesystems.
// We need to prepare a remote filesystem for testing.
var (
	// REMOTE_SSH is the remote system we are testing against.
	REMOTE_SSH = os.Getenv("TEST_REMOTE_SSH")
	// REMOTE_SSH_KEY is the private key path for ssh access.
	REMOTE_SSH_KEY = os.Getenv("TEST_REMOTE_KEY")

	// REMOTE_SSHC is the ssh client,
	// which will be prepared in the init function.
	REMOTE_SSHC *ssh.Client
)

func init() {
	// Get remote filesystem url.
	REMOTE_SSH_URL := errx.Must(url.Parse(REMOTE_SSH))
	// Open ssh connection.
	REMOTE_SSHC = errx.Must(ssh.Dial("tcp", REMOTE_SSH_URL.Host, &ssh.ClientConfig{
		User:            REMOTE_SSH_URL.User.Username(),
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(errx.Must(ssh.ParsePrivateKey(errx.Must(os.ReadFile(REMOTE_SSH_KEY))))),
		},
	}))
}

// run executes a command on the remote machine.
// This is a helper function to prepare the remote filesystem for testing.
func run(cmd string) ([]byte, error) {
	// Open a new session.
	session := errx.Must(REMOTE_SSHC.NewSession())
	defer session.Close()
	// Run a command.
	return session.Output(cmd)
}

func RunTestStat(t *testing.T, _fs fs.FS) {
	// First, let's create a file to test.
	_, err := run("echo 'Hello, World!' > /tmp/test.txt")
	if err != nil {
		t.Fatal("something wrong with preparations:", err)
	}
	// Defer cleanup.
	defer func() {
		_, err := run("rm /tmp/test.txt")
		if err != nil {
			t.Fatal("something wrong with cleanup:", err)
		}
	}()
	// Test stat.
	info, err := _fs.(fs.StatFS).Stat("/tmp/test.txt")
	if err != nil {
		t.Fatal("stat failed:", err)
	}
	if info.Name() != "test.txt" {
		t.Error("unexpected name:", info.Name())
	}
	if info.Size() != 14 {
		t.Error("unexpected size:", info.Size())
	}
	if !info.Mode().IsRegular() {
		t.Error("unexpected mode:", info.Mode())
	}
	if info.ModTime().IsZero() {
		t.Error("unexpected modtime:", info.ModTime())
	}
}

func RunTestOpen(t *testing.T, _fs fs.FS) {
	// First, let's create a file to test.
	_, err := run("echo 'Hello, World!' > /tmp/test.txt")
	if err != nil {
		t.Fatal("something wrong with preparations:", err)
	}
	// Defer cleanup.
	defer func() {
		_, err := run("rm /tmp/test.txt")
		if err != nil {
			t.Fatal("something wrong with cleanup:", err)
		}
	}()
	// Test open.
	file, err := _fs.Open("/tmp/test.txt")
	if err != nil {
		t.Fatal("open failed:", err)
	}
	defer file.Close()
	// Read file to test content.
	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatal("read failed:", err)
	}
	if string(data) != "Hello, World!\n" {
		t.Error("unexpected content:", string(data))
	}
}

func RunTestOpenFile(t *testing.T, _fs fs.FS) {
	// We will test both write and read operations here.
	// First, let's create a file to test.
	file, err := _fs.(fsx.OpenFileFS).OpenFile("/tmp/test.txt", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		t.Fatal("open failed:", err)
	}
	defer file.Close()
	// Write some data.
	_, err = file.Write([]byte("Hello, World!\n"))
	if err != nil {
		t.Fatal("write failed:", err)
	}
	// Seek to the beginning.
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		t.Fatal("seek failed:", err)
	}
	// Read data.
	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatal("read failed:", err)
	}
	if string(data) != "Hello, World!\n" {
		t.Error("unexpected content:", string(data))
	}
	// Remove the file.
	_, err = run("rm /tmp/test.txt")
	if err != nil {
		t.Fatal("remove failed:", err)
	}
}

func RunTestReadFile(t *testing.T, _fs fs.FS) {
	// First, let's create a file to test.
	_, err := run("echo 'Hello, World!' > /tmp/test.txt")
	if err != nil {
		t.Fatal("something wrong with preparations:", err)
	}
	// Defer cleanup.
	defer func() {
		_, err := run("rm /tmp/test.txt")
		if err != nil {
			t.Fatal("something wrong with cleanup:", err)
		}
	}()
	// Test read file.
	data, err := _fs.(fs.ReadFileFS).ReadFile("/tmp/test.txt")
	if err != nil {
		t.Fatal("read failed:", err)
	}
	if string(data) != "Hello, World!\n" {
		t.Error("unexpected content:", string(data))
	}
}

func RunTestWriteFile(t *testing.T, _fs fs.FS) {
	// Write file using fs.
	err := _fs.(fsx.WriteFileFS).WriteFile("/tmp/test.txt", []byte("Hello, World!\n"), 0644)
	if err != nil {
		t.Fatal("write failed:", err)
	}
	// Read file to test content.
	data, err := run("cat /tmp/test.txt")
	if err != nil {
		t.Fatal("read failed:", err)
	}
	if string(data) != "Hello, World!\n" {
		t.Error("unexpected content:", string(data))
	}
	// Remove the file.
	_, err = run("rm /tmp/test.txt")
	if err != nil {
		t.Fatal("remove failed:", err)
	}
}

func RunTestReadDir(t *testing.T, _fs fs.FS) {
	// First, let's create a directory to test.
	_, err := run("mkdir /tmp/test")
	if err != nil {
		t.Fatal("something wrong with preparations:", err)
	}
	// Defer cleanup.
	defer func() {
		_, err := run("rm -r /tmp/test")
		if err != nil {
			t.Fatal("something wrong with cleanup:", err)
		}
	}()
	// Test read directory.
	entries, err := _fs.(fs.ReadDirFS).ReadDir("/tmp")
	if err != nil {
		t.Fatal("readdir failed:", err)
	}
	if len(entries) == 0 {
		t.Error("unexpected entries:", len(entries))
	}
	// Find the test directory.
	var dir os.DirEntry
	for _, entry := range entries {
		if entry.Name() == "test" {
			dir = entry
			break
		}
	}
	if dir == nil {
		t.Error("test directory not found")
	}
	// Test the directory.
	if !dir.IsDir() {
		t.Error("unexpected type:", dir.Type())
	}
}
