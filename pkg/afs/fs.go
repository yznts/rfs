/*
AFS is a remote filesystem implementation with agent.
It requires a server part (agent) to be running on the remote machine.

Package consists of multiple parts:
- Agent: a backend that serves filesystem operations over a simple request/response schema.
- Server: implements a protocol server that listens for incoming connections and forwards requests to the agent.
- Client: implements a protocol client that connects to the server and sends requests to the agent.
- FS: a wrapper around a protocol client that acts as a bridge between the client and the filesystem interfaces.

agent <-> server <-> client <-> fs

It designed to be able to work with different protocols, like rpc, ws, etc.
*/
package afs

import (
	"io/fs"
	"math/rand/v2"
)

// FS is a wrapper around a protocol client.
// Acts as a bridge between the client and the filesystem interfaces.
type FS struct {
	client Client
}

// Stat implements the Stat method of fs.StatFS.
func (f *FS) Stat(name string) (fs.FileInfo, error) {
	// Send a Stat request to the agent
	res := f.client.Stat(&agentRequest{
		Id:   rand.IntN(9999999999),
		Op:   "Stat",
		Args: []any{name},
	})
	_ = res
	return nil, nil
}

func NewFS(client Client) *FS {
	return &FS{
		client: client,
	}
}
