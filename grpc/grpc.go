package grpc

import (
	"context"
	"net"

	"github.com/charmbracelet/log"
	pb "github.com/zarinit-routers/connector-rpc/gen/connector"
	googleRPC "google.golang.org/grpc"
)

func Serve() error {
	listener, err := net.Listen("tcp", getListenAddr())
	if err != nil {
		return err
	}

	var opts []googleRPC.ServerOption
	srv := googleRPC.NewServer(opts...)
	pb.RegisterClientsServiceServer(srv, clientsService{})

	log.Info("Starting gRPC server", "address", getListenAddr())
	return srv.Serve(listener)
}
func getListenAddr() string {
	return "localhost:9080"
}

type clientsService struct {
	pb.UnimplementedClientsServiceServer
}

func (s *clientsService) GetClients(ctx context.Context, req *pb.GetClientsRequest) (*pb.GetClientsResponse, error) {
	response := &pb.GetClientsResponse{}
	response.Clients = append(response.Clients,
		&pb.Client{
			Id:   "1",
			Name: "Dummy Client 1",
		})
	response.Clients = append(response.Clients,
		&pb.Client{
			Id:   "2",
			Name: "Dummy Client 2",
		})
	return response, nil
}
