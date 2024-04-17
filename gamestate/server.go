package gamestate

import (
	"errors"

	"github.com/google/uuid"
)

type GameStateSingleton struct {
	games                  map[string]*Game
	WaitingRooms           map[string]*Game
	matchMakingWaitingRoom *Game
	MatchMakingWaiting     chan GameUpdateRequest
	NewMatchMakingRequests chan MatchMakingRequest
	NewWaitingRoomRequests chan WaitingRoomRequest
	GameUpdateRequests     chan GameUpdateRequest
}

func Init() GameStateSingleton {
	return GameStateSingleton{
		games:                  make(map[string]*Game),
		WaitingRooms:           make(map[string]*Game),
		MatchMakingWaiting:     make(chan GameUpdateRequest),
		NewMatchMakingRequests: make(chan MatchMakingRequest),
		NewWaitingRoomRequests: make(chan WaitingRoomRequest),
		GameUpdateRequests:     make(chan GameUpdateRequest),
	}
}

type MatchMakingRequest struct {
	Res chan GameResponse
}

type WaitingRoomRequest struct {
	res chan GameResponse
}

type FriendJoinRequest struct {
	res chan FriendJoinResponse
}

type GameResponse struct {
	Error    error
	PlayerID string
	G        Game
	Ready    bool
}

type FriendJoinResponse struct {
	G     Game
	Found bool
}

type GameUpdateRequest struct {
	Res      chan GameResponse
	ID       string
	PlayerID string
	A        Action
}

func (gss *GameStateSingleton) StartProcessing() {
	for {
		select {
		case req := <-gss.NewMatchMakingRequests:
			if gss.matchMakingWaitingRoom == nil {
				idStruct, err := uuid.NewV6()
				if err != nil {
					panic(err)
				}
				gameID := idStruct.String()
				g := NewGame(gameID)
				gss.games[gameID] = g
				leftPlayerID, err := uuid.NewV6()
				if err != nil {
					panic(err)
				}
				g.LeftPlayerID = leftPlayerID.String()
				req.Res <- GameResponse{
					G:        *g,
					Ready:    false,
					PlayerID: g.LeftPlayerID,
				}
				gss.matchMakingWaitingRoom = g
			} else {
				rightPlayerID, err := uuid.NewV6()
				if err != nil {
					panic(err)
				}
				rightPlayerIDString := rightPlayerID.String()
				gss.matchMakingWaitingRoom.RightPlayerID = rightPlayerIDString
				gss.matchMakingWaitingRoom.start()
				req.Res <- GameResponse{
					G:        *gss.matchMakingWaitingRoom,
					Ready:    true,
					PlayerID: rightPlayerIDString,
				}
				gss.matchMakingWaitingRoom.start()
				gss.matchMakingWaitingRoom = nil
			}
		//case req := <-gss.NewWaitingRoomRequests:

		//case req := <-gss.FriendJoinRequests:
		case req := <-gss.MatchMakingWaiting:
			{
				g, ok := gss.games[req.ID]
				if !ok {
					req.Res <- GameResponse{
						Error: errors.New("no game found"),
					}
					goto end
				}
				switch g.GameState {
				case WAITING:
					{
						req.Res <- GameResponse{
							PlayerID: req.PlayerID,
							G:        *g,
							Ready:    false,
							Error:    nil,
						}
						goto end
					}
				default:
					{
						req.Res <- GameResponse{
							G:        *gss.matchMakingWaitingRoom,
							Ready:    true,
							PlayerID: req.PlayerID,
						} //do nothing and allow app to continue
					}
				}

			}

		case req := <-gss.GameUpdateRequests:
			g, ok := gss.games[req.ID]
			if !ok {
				req.Res <- GameResponse{
					Error: errors.New("no game found"),
				}
				goto end
			}
			switch g.GameState {
			case WAITING:
				{
					req.Res <- GameResponse{
						PlayerID: req.PlayerID,
						G:        *g,
						Ready:    false,
						Error:    nil,
					}
					goto end
				}
			default:
				{
					//do nothing and allow app to continue
				}
			}
			g.play(req.A, req.PlayerID)
			req.Res <- GameResponse{
				PlayerID: req.PlayerID,
				G:        *g,
				Ready:    true,
				Error:    nil,
			}
		}
	end:
	}
}
