package connections

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/gorilla/websocket"
	"github.com/zarinit-routers/cloud-connector/models"
)

const (
	AuthorizationHeader = "Authorization"
	RouterIDHeader      = "X-Router-ID"
	GroupIDHeader       = "X-Group-ID"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	connections = map[string]*websocket.Conn{}
)

func AppendConnection(nodeId string, conn *websocket.Conn) {
	if connections[nodeId] != nil {
		log.Warn("Connection with that node already exists, closing it", "nodeId", nodeId)
		closeConn(connections[nodeId])
	}
	connections[nodeId] = conn
}

func closeConn(conn *websocket.Conn) {
	addr := fmt.Sprintf("%s %s", conn.RemoteAddr().Network(), conn.RemoteAddr().String())

	if err := conn.Close(); err != nil {
		log.Error("Failed close connection, process will be forced", "address", addr, "error", err)
	}

	log.Warn("Connection closed", "address", addr)
}

func Serve() {
	srv := http.NewServeMux()
	srv.HandleFunc("/api/ipc/connect", func(w http.ResponseWriter, r *http.Request) {
		auth, err := checkAuth(r)
		if err != nil {
			log.Error("Failed authenticate connection", "error", err)
			return
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Error("Failed to upgrade connection", "error", err)
			return
		}

		AppendConnection(auth.NodeID, conn)

		go serveConnection(conn)

		log.Info("Connection established", "nodeId", auth.NodeID, "groupId", auth.GroupID)
	})

	log.Info("Starting connections server", "address", getAddress())
	http.ListenAndServe(getAddress(), srv)
}

type AuthData struct {
	NodeID  string
	GroupID string
}
type MessageHandlerFunc func(message []byte) error

var handlers = []MessageHandlerFunc{}

func AddHandler(handler MessageHandlerFunc) {
	handlers = append(handlers, handler)
}

func serveConnection(conn *websocket.Conn) {
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Error("Failed to read message", "error", err)
			continue
		}

		if messageType == websocket.CloseMessage {
			closeConn(conn)
			return
		}

		for _, handler := range handlers {
			go func() {
				if err := handler(message); err != nil {
					log.Error("Failed to handle message", "error", err)
				}
			}()
		}

	}
}

func checkAuth(r *http.Request) (*AuthData, error) {
	passphrase := r.Header.Get(AuthorizationHeader)
	routerId := r.Header.Get(RouterIDHeader)
	groupId := r.Header.Get(GroupIDHeader)

	if passphrase == "" || routerId == "" || groupId == "" {
		log.Error("Failed to authenticate", "routerId", routerId, "groupId", groupId, "passphrase", passphrase)
		return nil, errors.New("missing required headers")
	}

	if err := authenticateInGroup(groupId, passphrase); err != nil {
		log.Error("Failed to authenticate", "error", err, "routerId", routerId, "groupId", groupId, "passphrase", passphrase)
		return nil, err
	}

	return &AuthData{
		NodeID:  routerId,
		GroupID: groupId,
	}, nil
}

// Uses gRPC to connect to auth service
//
// TODO: properly implement
func authenticateInGroup(groupId string, passphrase string) error {
	return nil
}

func getAddress() string {
	return ":8080"
}

func SendRequest(nodeId string, r *models.ToNodeRequest) error {
	conn, ok := connections[nodeId]
	if !ok {
		return fmt.Errorf("node with id %q not connected", nodeId)
	}

	message, err := json.Marshal(r)
	if err != nil {
		return err
	}

	return conn.WriteMessage(websocket.TextMessage, message)
}
