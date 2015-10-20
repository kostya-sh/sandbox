package main

import (
	"log"
	"net"

	"github.com/kostya-sh/sandbox/grpc-challenge/grpcch"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

// server is used to implement hellowrld.GreeterServer.
type server struct{}

func (s *server) Call(ctx context.Context, in *grpcch.Request) (*grpcch.Reply, error) {
	return &grpcch.Reply{Message: "Hello " + in.Name}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	grpcch.RegisterServiceServer(s, &server{})
	s.Serve(lis)
}
