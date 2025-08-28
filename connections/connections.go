package connections

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/gorilla/websocket"
	"github.com/zarinit-routers/cloud-connector/globals"
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
)

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

		globals.AppendConnection(auth.NodeID, conn)

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

func serveConnection(conn *websocket.Conn) {
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Error("Failed to read message", "error", err)
			continue
		}
		log.Info("Received message", "messageType", messageType, "message", string(message))

		var request models.FromNodeResponse
		err = json.Unmarshal(message, &request)
		if err != nil {
			log.Error("Failed to unmarshal message", "error", err)
			continue
		}
		if err := globals.SendResponse(request.RequestID, request.ToCloud()); err != nil {
			log.Error("Failed to send response", "error", err)
			continue
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
