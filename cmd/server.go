package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"

	pb "github.com/idalmasso/gRPCMovementGameTest/grpc/gamesrv"
	"github.com/idalmasso/gRPCMovementGameTest/pkg/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	port=":3000"
)

type server struct{
	Instance *models.GameInstance
}
func (s *server) PlayerPositionSync(stream pb.GameService_StreamSyncServer) error{
	headers, ok :=metadata.FromIncomingContext(stream.Context())
	if !ok{
		return errors.New("Cannot decode metadata")
	}
	token:=headers["Authorization"]
	if len(token) == 0 {
		return errors.New("no token provided")
	}
	 for {
    in, err := stream.Recv()
    if err == io.EOF {
      return nil
    }
    if err != nil {
      return err
    }
		log.Println("Received "+in.GetName())
    
  }
}

func (s *server) Login(c context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error){
	var response *pb.LoginResponse
	response.Token="000"
	return response,nil
}
func (s *server) getClientFromContext(ctx context.Context) (*models.Player,error){
	headers, ok:=metadata.FromIncomingContext(ctx)
	if !ok{
		return nil, fmt.Errorf("Not found metadata headers")
	}
	name:=headers["authorization"]
	if len(name)==0{
		return nil, fmt.Errorf("No authorization token passed")
	}
	s.Instance.GSMutex.RLock()
	p, ok:=s.Instance.Players[name[0]]
	s.Instance.GSMutex.RUnlock()
	if !ok{//TODO: Change this next
		return &models.Player{Name: name[0]}, nil
	}
	return p, nil
}
func (s *server) StreamSync(srv pb.GameService_StreamSyncServer) error{
	ctx := srv.Context()
	currentClient, err := s.getClientFromContext(ctx)
	if err != nil {
		return err
	}

	currentClient.StreamServer = srv
	currentClient.PosX=30
	currentClient.PosY=50
	s.Instance.AddPlayer(currentClient)
	log.Println("start new server")

	// Wait for stream requests.

		for {
			req, err := srv.Recv()
			if err != nil {
				log.Printf("receive error %v", err)
				//currentClient.done <- errors.New("failed to receive request")
				return err
			}
			//currentClient.lastMessage = time.Now()
			var movement models.PlayerMovement 
			movement.GetFromInput(req)
			s.Instance.InputsChannel<-&movement
		}


	
}
//Create a gRPC server and listen for incoming requests
func main() {
	var gi models.GameInstance
	gi.DimX=1000
	gi.DimY=1000
	gi.Name="TEST"
	gi.RunInstance()
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on port [%s]: %v", port, err)
	}
	s := grpc.NewServer()
	pb.RegisterGameServiceServer(s, &server{Instance: &gi})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
	
}
