package models

import (
	"sync"

	"github.com/google/uuid"
	pm "github.com/idalmasso/gRPCMovementGameTest/grpc/gamesrv"
)

//Player contains the actual data of all this game stuff
type Player struct{
	token uuid.UUID
	Name string
	PosX float32
	PosY float32
	StreamServer pm.GameService_StreamSyncServer
	PlayerMutex sync.RWMutex
}



