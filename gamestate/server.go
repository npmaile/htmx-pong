package gamestate

import (
	"fmt"

	"github.com/google/uuid"
)

type GameStateSingleton struct {
	games              map[string]*Game
	NewGameRequests    chan NewGameRequest
	GameUpdateRequests chan GameUpdateRequest
}

func Init() GameStateSingleton {
	return GameStateSingleton{
		games:              make(map[string]*Game),
		NewGameRequests:    make(chan NewGameRequest),
		GameUpdateRequests: make(chan GameUpdateRequest),
	}
}

type NewGameRequest struct {
	Res chan Game
}

type GameUpdateRequest struct {
	Ch chan Game
	ID string
	A  Action
}

func (gss *GameStateSingleton) StartProcessing() {
	for {
		select {
		case req := <-gss.NewGameRequests:
			fmt.Println("got new game request")
			uid, _ := uuid.NewRandom()
			id := uid.String()
			game := NewGame(id)
			gss.games[id] = game
			req.Res <- *game
		case req := <-gss.GameUpdateRequests:
			g := gss.games[req.ID]
			g.play(req.A)
			req.Ch <- *g

		}
	}
}
