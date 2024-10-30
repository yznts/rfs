package afs

type Client interface {
	Stat(req *agentRequest) *agentResponse
	Open(req *agentRequest) *agentResponse
	OpenFile(req *agentRequest) *agentResponse
}
