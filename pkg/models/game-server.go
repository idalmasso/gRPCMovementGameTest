package models

import (
	"fmt"
	"log"
	"sync"

	"github.com/idalmasso/world-changers/grpc/gamesrv"
)

//GameInstance is a struct with all the data of an instance of the game
type GameInstance struct{
	Name string
	DimX float32
	DimY float32
	SomethingHappened bool
	InputsChannel chan InputEvent
	EventsChannel chan InputEvent
	Players map[string]*Player
	PlayersChanged map[string]bool
	GSMutex sync.RWMutex
}
//AddPlayer adds a player to the instance
func (g *GameInstance) AddPlayer(p *Player) error{
	fmt.Println("Add Player: ", p.Name)
	g.GSMutex.Lock()
	defer	g.GSMutex.Unlock()
	if _, ok:= g.Players[p.Name]; !ok{
		g.Players[p.Name]=p;
		g.PlayersChanged[p.Name]=true
	} else{
		return fmt.Errorf("User already exists")
	}
	g.SomethingHappened=true
	return nil
}
//RemovePlayer removes a player to the instance
func (g *GameInstance) RemovePlayer(name string) {
	g.GSMutex.Lock()
	defer	g.GSMutex.Unlock()
	delete(g.Players, name)
	delete(g.PlayersChanged, name)
}

//RunInstance launch the actual game loop
func (g* GameInstance) RunInstance(){
	g.InputsChannel=make(chan InputEvent, 1)
	g.EventsChannel = make(chan InputEvent, 1)
	g.Players = make(map[string]*Player)
	g.PlayersChanged = make(map[string]bool)
	go func(){
		for{
			g.manageInputs()
			g.update()
			g.GSMutex.Lock()
			if g.SomethingHappened{
				log.Println("WOW")
				g.broadcastStatus()
				g.SomethingHappened=false
			}
			g.GSMutex.Unlock()
			
		}
	}()
}

func (g *GameInstance) manageInputs(){
	
	for{
		select {
		case value:=<-g.InputsChannel:
			if value.CanPerform(g){
				value.Perform(g)
				
				g.GSMutex.Lock()
				g.SomethingHappened=true
				g.GSMutex.Unlock()
			}
		default:
			return
		}
	}
}


func (g *GameInstance) update(){
	return
}


func (g *GameInstance) broadcastStatus(){
	
	positions:=g.getPlayerPositionsObject()
	for _,p:=range(g.Players){
		p.PlayerMutex.Lock()
		p.StreamServer.Send(positions)
		p.PlayerMutex.Unlock()
	}
}

func (g *GameInstance) getPlayerPositionsObject() *gamesrv.PlayerPositionArray{
	var positions gamesrv.PlayerPositionArray
	positions.PlayerPositions=make([]*gamesrv.PlayerPosition, len(g.Players))
	idx:=0
	for _, p:=range(g.Players){
		if g.PlayersChanged[p.Name]{
			var pos gamesrv.PlayerPosition
			pos.Name=p.Name
			pos.X=p.PosX
			pos.Y=p.PosY
			positions.PlayerPositions[idx]=&pos
			g.PlayersChanged[p.Name]=false
		}
	}
	return &positions
}


