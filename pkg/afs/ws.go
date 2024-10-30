package afs

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
)

type WebsocketClient struct {
	conn *websocket.Conn
}

func NewWebsocketClient(addr string) (*WebsocketClient, error) {
	ws, _, err := websocket.DefaultDialer.Dial(addr, nil)
	return &WebsocketClient{
		conn: ws,
	}, err
}

func NewWebsocketHandler(settings agentSettings) http.HandlerFunc {
	agent := &Agent{
		agentSettings: settings,
		files:         make(map[string]*os.File),
		locks:         make(map[string]*sync.Mutex),
	}
	updater := websocket.Upgrader{}
	return func(w http.ResponseWriter, r *http.Request) {
		// Initialize a new websocket connection
		c, err := updater.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		defer c.Close()
		// Mainloop
		for {
			// Read the message from the client
			messageType, message, err := c.ReadMessage()
			if err != nil {
				log.Println("err:", err)
				break
			}
			// Pass non-text messages
			if messageType != websocket.TextMessage {
				continue
			}
			// Decode request
			req := &agentRequest{}
			if err := json.Unmarshal(message, req); err != nil {
				log.Println("err:", err)
				break
			}
			// Handle request
			var res *agentResponse
			switch req.Op {
			case "Stat":
				res = agent.Stat(req)
			case "Open":
				res = agent.Open(req)
			case "OpenFile":
				res = agent.OpenFile(req)
			default:
				log.Println("err: unknown op", req.Op)
			}
			// Encode and send response
			resBytes, err := json.Marshal(res)
			if err != nil {
				log.Println("err:", err)
				break
			}
			if err := c.WriteMessage(websocket.TextMessage, resBytes); err != nil {
				log.Println("err:", err)
				break
			}
		}
	}

}
