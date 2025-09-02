package models

import (
	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5/pgtype"
)

type JsonMap = map[string]any

type UUID = pgtype.UUID

func ParseUUID(str string) (UUID, error) {

	id := pgtype.UUID{}
	err := id.Scan(str)
	if err != nil {
		log.Warn("Bad UUID string", "string", str, "error", err)
		return id, err
	}
	return id, nil
}

type FromCloudRequest struct {
	NodeID  string  `json:"nodeId"`
	Command string  `json:"command"`
	Args    JsonMap `json:"args"`
}
type ToCloudResponse struct {
	RequestError string  `json:"requestError"` // Connector error
	CommandError string  `json:"commandError"` // Node error
	Data         JsonMap `json:"data"`
}

type ToNodeRequest struct {
	RequestID string  `json:"requestId"`
	Command   string  `json:"command"`
	Args      JsonMap `json:"args"`
}
type FromNodeResponse struct {
	RequestID string  `json:"requestId"`
	Data      JsonMap `json:"data"`
	Error     string  `json:"error"`
}

func (r *FromCloudRequest) ToNode(requestId string) *ToNodeRequest {
	return &ToNodeRequest{
		RequestID: requestId,
		Command:   r.Command,
		Args:      r.Args,
	}
}

func (r *FromNodeResponse) ToCloud() *ToCloudResponse {
	return &ToCloudResponse{
		CommandError: r.Error,
		Data:         r.Data,
	}
}
