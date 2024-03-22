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

type NewGameResponse struct {
	GameID string
	Game   Game
}

type NewGameRequest struct {
	Res chan NewGameResponse
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
			game := NewGame()
			gss.games[id] = game
			req.Res <- NewGameResponse{
				Game:   *game,
				GameID: id,
			}
		case req := <-gss.GameUpdateRequests:
			g := gss.games[req.ID]
			g.play(req.A)
			req.Ch <- *g

		}
	}
}
