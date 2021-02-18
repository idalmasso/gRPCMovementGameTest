package models

import (
	"fmt"
	"log"

	pb "github.com/idalmasso/gRPCMovementGameTest/grpc/gamesrv"
)

const (
	//POSITION const INPUT TYPE position event
	POSITION =iota
	//PLAYERSPAWN const player spawn event
	PLAYERSPAWN
	//ENTITYCREATE const entity created event
	ENTITYCREATE
)
//InputEvent is the interface that will manage the events from clients
type InputEvent interface {
	GetType() int32
	CanPerform(instance *GameInstance ) bool
	Perform(instance *GameInstance) error
}
//PlayerMovement is a struct implementing InputEvent for player movement
type PlayerMovement struct{
	playerName  string 
	deltaX float32
	deltaY float32
}
//GetFromInput return a plpayermovement from a grpc playerposition
func (playerMovement *PlayerMovement) GetFromInput(position *pb.PlayerPosition){
	playerMovement.playerName=position.Name
	playerMovement.deltaX=position.X
	playerMovement.deltaY=position.Y
}
//GetType returns POSITION for movement
func (playerMovement *PlayerMovement) GetType() int32{
	return POSITION
}
//CanPerform always trye
func (playerMovement *PlayerMovement) CanPerform(instance *GameInstance) bool{
	log.Println("Can perform X: ",playerMovement.deltaX,"Y: ",playerMovement.deltaY)
	return playerMovement.deltaX!=0 || playerMovement.deltaY!=0
}
//Perform do the actual move for a player
func (playerMovement *PlayerMovement) Perform(instance *GameInstance) error{
	log.Println("Perform X: ",playerMovement.deltaX,"Y: ",playerMovement.deltaY)
	instance.GSMutex.Lock()
	defer instance.GSMutex.Unlock()
	p, ok:= instance.Players[playerMovement.playerName]
	if ok{
		p.PlayerMutex.Lock()
		defer p.PlayerMutex.Unlock()
		p.PosX+=playerMovement.deltaX
		p.PosY+=playerMovement.deltaY
	}
	instance.PlayersChanged[playerMovement.playerName]=true
	return nil
}
//PlayerSpawn is the event for spawning of player
type PlayerSpawn struct{
	player  Player 
}
//GetType returns always PLAYERSPAWN
func (spawn *PlayerSpawn) GetType() int32{
	return PLAYERSPAWN
}
//CanPerform only f the player name does not exists
func (spawn *PlayerSpawn) CanPerform(instance *GameInstance) bool{
	instance.GSMutex.RLock()
	instance.GSMutex.RUnlock()
	_, ok:=instance.Players[spawn.player.Name];
	return !ok
}
//Perform add the player in the instance
func (spawn *PlayerSpawn) Perform(instance *GameInstance) error{
	instance.GSMutex.Lock()
	instance.GSMutex.Unlock()
	if _, ok:=instance.Players[spawn.player.Name]; !ok{
		instance.AddPlayer(&spawn.player)
	} else {
		return fmt.Errorf("User already exists")
	}
	return nil
}
