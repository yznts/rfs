package afs

import (
	"encoding/json"
	"errors"
	"os"
	"path"
	"sync"

	"github.com/yznts/rfs/pkg/fsx"
	"github.com/yznts/zen/v3/errx"
)

type agentSettings struct {
	User string // Authorization user
	Pass string // Authorization pass
	Root string // Root directory to serve from
}

type agentRequest struct {
	Id   int    `json:"id"`
	Op   string `json:"op"`
	Args []any  `json:"args"`
}

type agentResponse struct {
	Id  int    `json:"id"`
	Err string `json:"err"`
	Res string `json:"res"` // use marshal/unmarshal
}

func (r *agentResponse) Error() error {
	if r.Err != "" {
		return errors.New(r.Err)
	}
	return nil
}

// Agent is a protocol agnostic filesystem agent.
// It implements filesystem operations and serves them over
// a simple request/response schema.
// It is intended to be used as a backend for a specific
// protocol, like websocket or http.
type Agent struct {
	// Agent settings.
	agentSettings
	// Opened file descriptors.
	files map[string]*os.File
	// Opened file locks.
	locks map[string]*sync.Mutex
}

func (a *Agent) marshal(data any) string {
	return string(errx.Must(json.Marshal(data)))
}

func (a *Agent) unmarshal(data string, dest any) {
	errx.Must(0, json.Unmarshal([]byte(data), dest))
}

// Stat handles file stat request.
func (a *Agent) Stat(req *agentRequest) *agentResponse {
	// Unpack args.
	name, _ := req.Args[0].(string)
	// Get file info.
	fileinfo, err := os.Stat(
		path.Join(a.agentSettings.Root, name),
	)
	// Prepare and return response.
	return &agentResponse{
		Id:  req.Id,
		Err: marshalError(err),
		Res: a.marshal(fsx.NewEntryFromFileInfo(fileinfo)),
	}
}

// Open handles file open request.
func (a *Agent) Open(req *agentRequest) *agentResponse {
	// Unpack args.
	name, _ := req.Args[0].(string)
	// Resolve absolute file path.
	filepath := path.Join(a.agentSettings.Root, name)
	// Open file descriptor.
	file, err := os.Open(filepath)
	// Save it and prepare response.
	a.files[filepath] = file
	// Prepare and return response.
	return &agentResponse{
		Id:  req.Id,
		Err: marshalError(err),
		Res: filepath,
	}
}

func (a *Agent) OpenFile(req *agentRequest) *agentResponse {
	// Unpack args.
	name, _ := req.Args[0].(string)
	fileflag, _ := req.Args[1].(int)
	fileperm, _ := req.Args[2].(int)
	// Resolve absolute file path.
	filepath := path.Join(a.agentSettings.Root, name)
	// Open file descriptor.
	file, err := os.OpenFile(filepath, fileflag, os.FileMode(fileperm))
	// Save it and prepare response.
	a.files[filepath] = file
	// Prepare and return response.
	return &agentResponse{
		Id:  req.Id,
		Err: marshalError(err),
		Res: filepath,
	}
}
