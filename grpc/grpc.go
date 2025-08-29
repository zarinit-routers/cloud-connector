package grpc

import (
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
	pb.RegisterNodesServer(srv, newNodesService())

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
