package gamestate

import (
	"fmt"

	"github.com/google/uuid"
)

type GameStateSingleton struct {
	games                  map[string]*Game
	NewMatchMakingRequests chan MatchMakingRequest
	NewWaitingRoomRequests chan WaitingRoomRequest
	GameUpdateRequests     chan GameUpdateRequest
}

func Init() GameStateSingleton {
	return GameStateSingleton{
		games:                  make(map[string]*Game),
		NewMatchMakingRequests: make(chan MatchMakingRequest),
		NewWaitingRoomRequests: make(chan WaitingRoomRequest),
		GameUpdateRequests:     make(chan GameUpdateRequest),
	}
}

type MatchMakingRequest struct {
	res chan MatchMakingResponse
}

type WaitingRoomRequest struct {
	res chan WaitingRoomResponse
}

type WaitingRoomResponse struct {
	g     Game
	ready bool
}

type MatchMakingResponse struct {
	g     Game
	ready bool
}

type GameUpdateRequest struct {
	Ch       chan Game
	ID       string
	PlayerID string
	A        Action
}

func (gss *GameStateSingleton) StartProcessing() {
	for {
		select {
		case req := <-gss.NewMatchMakingRequests:
		case req := <-gss.NewWaitingRoomRequests:
		case req := <-gss.GameUpdateRequests:
			g := gss.games[req.ID]
			g.play(req.A, req.PlayerID)
			req.Ch <- *g
		}
	}
}
