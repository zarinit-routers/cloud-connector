package grpc

import (
	"context"
	"net"
	"os"

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

const ENV_GRPC_ADDR = "CONNECTOR_GRPC_ADDR"

func getListenAddr() string {
	addr := os.Getenv(ENV_GRPC_ADDR)
	if addr == "" {
		log.Fatal("GRPC address not set", "envVariable", ENV_GRPC_ADDR)
	}
	return addr
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
