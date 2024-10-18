package fsx

import (
	"encoding/gob"
	"io"
	"io/fs"
	"time"
)

func init() {
	gob.Register(Entry{})
}

// File interface includes basic interfaces
// file must to implement.
// Most of the methods are from os.File,
// but we don't include all of them.
type File interface {
	io.Reader
	io.Writer
	io.Closer
	io.Seeker

	Truncate(size int64) error
}

// Entry is a universal struct with multiple purposes.
// - It implements both fs.FileInfo and fs.DirEntry interfaces.
// - It can be serialized and deserialized.
type Entry struct {
	EntryName    string
	EntrySize    int64
	EntryMode    fs.FileMode
	EntryModTime time.Time
}

func (e Entry) Name() string {
	return e.EntryName
}

func (e Entry) Size() int64 {
	return e.EntrySize
}

func (e Entry) Mode() fs.FileMode {
	return e.EntryMode
}

func (e Entry) ModTime() time.Time {
	return e.EntryModTime
}

func (e Entry) IsDir() bool {
	return e.EntryMode.IsDir()
}

func (e Entry) Sys() any {
	return nil
}

// DirEntry compatibility

func (e Entry) Type() fs.FileMode {
	return e.EntryMode
}

func (e Entry) Info() (fs.FileInfo, error) {
	return e, nil
}

// Entry constructors

func NewEntryFromFileInfo(f fs.FileInfo) Entry {
	return Entry{
		EntryName:    f.Name(),
		EntrySize:    f.Size(),
		EntryMode:    f.Mode(),
		EntryModTime: f.ModTime(),
	}
}

func NewEntryFromDirEntry(d fs.DirEntry) Entry {
	var (
		size    int64     = 0
		modtime time.Time = time.Now()
	)
	if i, err := d.Info(); err == nil {
		size = i.Size()
		modtime = i.ModTime()
	}
	return Entry{
		EntryName:    d.Name(),
		EntrySize:    size,
		EntryMode:    d.Type(),
		EntryModTime: modtime,
	}
}
