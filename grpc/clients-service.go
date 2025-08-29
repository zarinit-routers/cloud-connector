package grpc

import (
	"context"

	pb "github.com/zarinit-routers/connector-rpc/gen/connector"
)

type clientsService struct {
	pb.UnimplementedClientsServiceServer
}

func newClientsService() *clientsService {
	return &clientsService{}
}

func (s *clientsService) GetClients(ctx context.Context, req *pb.GetClientsRequest) (*pb.GetClientsResponse, error) {
	response := &pb.GetClientsResponse{}
	response.Clients = append(response.Clients,
		&pb.Client{
			Id:   "1",
			Name: "Dummy Client 1",
			Tags: []string{"dummy", "cool", "first"},
		})
	response.Clients = append(response.Clients,
		&pb.Client{
			Id:   "2",
			Name: "Dummy Client 2",
			Tags: []string{"dummy", "metal", "pop"},
		})
	return response, nil
}
