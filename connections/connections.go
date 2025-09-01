package connections

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/zarinit-routers/cloud-connector/models"
	"github.com/zarinit-routers/cloud-connector/storage/repository"
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
	connections = map[models.UUID]*websocket.Conn{}

	ctx = context.Background()
)

func AppendConnection(node *AuthData, conn *websocket.Conn) {
	if connections[node.NodeID] != nil {
		log.Warn("Connection with that node already exists, closing it", "nodeId", node.NodeID)
		closeConn(node.NodeID, connections[node.NodeID])
		repository.GetQueries().ReconnectNode(ctx, repository.ReconnectNodeParams{
			Id:      node.NodeID,
			GroupId: node.GroupID,
		})
	} else {
		repository.GetQueries().NewNode(ctx, repository.NewNodeParams{
			Id:      node.NodeID,
			GroupId: node.GroupID,
			Name: pgtype.Text{
				String: GenNodeName(),
				Valid:  true},
		})
	}

	connections[node.NodeID] = conn
}

func closeConn(nodeId models.UUID, conn *websocket.Conn) {
	addr := fmt.Sprintf("%s %s", conn.RemoteAddr().Network(), conn.RemoteAddr().String())

	if err := conn.Close(); err != nil {
		log.Error("Failed close connection, process will be forced", "address", addr, "error", err)
	}

	repository.GetQueries().UpdateLastConnection(ctx, nodeId)

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

		AppendConnection(auth, conn)

		go serveConnection(auth.NodeID, conn)

		log.Info("Connection established", "nodeId", auth.NodeID, "groupId", auth.GroupID)
	})

	log.Info("Starting connections server", "address", getAddress())
	http.ListenAndServe(getAddress(), srv)
}

type AuthData struct {
	NodeID  models.UUID
	GroupID models.UUID
}
type MessageHandlerFunc func(message []byte) error

var handlers = []MessageHandlerFunc{}

func AddHandler(handler MessageHandlerFunc) {
	handlers = append(handlers, handler)
}

func serveConnection(nodeId models.UUID, conn *websocket.Conn) {
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Error("Failed to read message", "error", err)
			continue
		}

		if messageType == websocket.CloseMessage {
			closeConn(nodeId, conn)
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
	groupIdStr := r.Header.Get(GroupIDHeader)

	if passphrase == "" {
		log.Error("Failed to authenticate passphrase header is empty")
		return nil, fmt.Errorf("missing passphrase header")
	}
	if groupIdStr == "" {
		log.Error("Group ID header is empty")
		return nil, fmt.Errorf("missing group id header %q", GroupIDHeader)
	}

	if err := authenticateInGroup(groupIdStr, passphrase); err != nil {
		log.Error("Failed to authenticate", "error", err, "groupId", groupIdStr, "passphrase", passphrase)
		return nil, err
	}

	routerId, idParseErr := models.UUIDFromString(r.Header.Get(RouterIDHeader))
	groupId, groupParseErr := models.UUIDFromString(groupIdStr)

	if err := errors.Join(idParseErr, groupParseErr); err != nil {
		log.Error("Bad UUID specifications", "error", err)
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
	id, err := models.UUIDFromString(nodeId)
	if err != nil {
		return fmt.Errorf("bad node id %q: %s", nodeId, err)
	}
	conn, ok := connections[id]
	if !ok {
		return fmt.Errorf("node with id %q not connected", nodeId)
	}

	message, err := json.Marshal(r)
	if err != nil {
		return err
	}

	return conn.WriteMessage(websocket.TextMessage, message)
}
