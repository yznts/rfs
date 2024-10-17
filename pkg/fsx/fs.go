package fsx

import "io/fs"

// OpenFileFS defines how we are opening file with FS instance.
// Unfortunately, io/fs have only Open method interface,
// which doesn't allow to write.
//
// Repeats os.OpenFile signature.
type OpenFileFS interface {
	OpenFile(name string, flag int, perm fs.FileMode) (File, error)
}

// WriteFileFS defines how we are writing with FS instance.
// Unfortunately, io/fs is missing writing interfaces.
//
// Repeats os.WriteFile signature.
// Still, perm is not obligatory and depends on implementation.
type WriteFileFS interface {
	WriteFile(name string, data []byte, perm fs.FileMode) error
}
