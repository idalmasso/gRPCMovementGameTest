package main

import (
	"log"
	"net"

	pb "github.com/idalmasso/world-changers/grpc/player-movement"
	"google.golang.org/grpc"
)

const (
	port=":50000"
)

type server struct{}



//Create a gRPC server and listen for incoming requests
func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on port [%s]: %v", port, err)
	}
	s := grpc.NewServer()
	pb.RegisterColorGeneratorServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
