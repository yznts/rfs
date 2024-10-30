package sshfs

import (
	"io"
	"io/fs"

	"github.com/pkg/sftp"
	"github.com/yznts/rfs/pkg/fsx"
	"github.com/yznts/zen/v3/slice"
	"golang.org/x/crypto/ssh"
)

// FS implements filesystem operations over sftp client.
type FS struct {
	*sftp.Client
}

// Stat implements the Stat method of fs.StatFS.
func (f *FS) Stat(name string) (fs.FileInfo, error) {
	return f.Client.Stat(name)
}

// Open implements the Open method of fs.FS.
func (f *FS) Open(name string) (fs.File, error) {
	return f.Client.Open(name)
}

// OpenFile implements the OpenFile method of fsx.OpenFileFS.
func (f *FS) OpenFile(name string, flag int, _ fs.FileMode) (fsx.File, error) {
	return f.Client.OpenFile(name, flag)
}

// ReadFile implements the ReadFile method of fs.ReadFileFS.
func (f *FS) ReadFile(name string) ([]byte, error) {
	file, err := f.Client.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return io.ReadAll(file)
}

// WriteFile implements the WriteFile method of fsx.WriteFileFS.
func (f *FS) WriteFile(name string, data []byte, _ fs.FileMode) error {
	file, err := f.Client.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(data)
	return err
}

// ReadDir implements the ReadDir method of fs.ReadDirFS.
func (f *FS) ReadDir(name string) ([]fs.DirEntry, error) {
	filesinfo, err := f.Client.ReadDir(name)
	if err != nil {
		return nil, err
	}
	return slice.Map(filesinfo, func(v fs.FileInfo) fs.DirEntry {
		return fsx.NewEntryFromFileInfo(v)
	}), nil
}

// New creates a new FS instance wrapping the given ssh client.
func New(client *ssh.Client) *FS {
	// Wrap ssh client into sftp client.
	sftpc, err := sftp.NewClient(client)
	if err != nil {
		panic(err)
	}
	// Return the wrapped into fs sftp client.
	return &FS{sftpc}
}

// Validate if implementation satisfies required fs interfaces
var _ fs.FS = &FS{}
var _ fs.StatFS = &FS{}
var _ fs.ReadDirFS = &FS{}
var _ fs.ReadFileFS = &FS{}
var _ fsx.OpenFileFS = &FS{}
var _ fsx.WriteFileFS = &FS{}
