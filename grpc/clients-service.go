package grpc

import (
	"context"

	pb "github.com/zarinit-routers/connector-rpc/gen/connector"
)

type nodesService struct {
	pb.UnimplementedNodesServer
}

func newNodesService() *nodesService {
	return &nodesService{}
}

func (s *nodesService) NodesByGroup(ctx context.Context, req *pb.NodesByGroupRequest) (*pb.NodesResponse, error) {
	response := &pb.NodesResponse{}
	response.Clients = append(response.Clients,
		&pb.Node{
			Id:   "1",
			Name: "Dummy Client 1",
			Tags: []string{"dummy", "cool", "first"},
		})
	response.Clients = append(response.Clients,
		&pb.Node{
			Id:   "2",
			Name: "Dummy Client 2",
			Tags: []string{"dummy", "metal", "pop"},
		})
	return response, nil
}

func (s *nodesService) AddTag(ctx context.Context, req *pb.TagRequest) (*pb.Node, error) {
	return &pb.Node{
		Id:   req.ModeId,
		Name: "Dummy Client " + req.ModeId,
		Tags: []string{req.Tag},
	}, nil
}

func (s *nodesService) RemoveTag(ctx context.Context, req *pb.TagRequest) (*pb.Node, error) {
	return &pb.Node{
		Id:   req.ModeId,
		Name: "Dummy Client " + req.ModeId,
		Tags: []string{},
	}, nil
}
