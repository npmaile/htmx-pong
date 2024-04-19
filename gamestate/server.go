package gamestate

import (
	"errors"
	"math/rand"
	"strings"
)

type GameStateSingleton struct {
	games                  map[string]*Game
	WaitingRooms           map[string]*Game
	matchMakingWaitingRoom *Game
	MatchMakingWaiting     chan GameUpdateRequest
	NewMatchMakingRequests chan MatchMakingRequest
	NewWaitingRoomRequests chan WaitingRoomRequest
	FriendJoinRequests     chan FriendJoinRequest
	GameUpdateRequests     chan GameUpdateRequest
	CancelRequests         chan CancelRequest
	SinglePlayerRequests   chan SinglePlayerRequest
}

func Init() GameStateSingleton {
	return GameStateSingleton{
		games:                  make(map[string]*Game),
		WaitingRooms:           make(map[string]*Game),
		MatchMakingWaiting:     make(chan GameUpdateRequest),
		NewMatchMakingRequests: make(chan MatchMakingRequest),
		NewWaitingRoomRequests: make(chan WaitingRoomRequest),
		FriendJoinRequests:     make(chan FriendJoinRequest),
		GameUpdateRequests:     make(chan GameUpdateRequest),
		CancelRequests:         make(chan CancelRequest),
		SinglePlayerRequests:   make(chan SinglePlayerRequest),
	}
}

type SinglePlayerRequest struct {
	Res chan GameResponse
}

type CancelRequest struct {
	ID string
}

type MatchMakingRequest struct {
	Res chan GameResponse
}

type WaitingRoomRequest struct {
	Res chan GameResponse
}

type FriendJoinRequest struct {
	Res chan GameResponse
	ID  string
}

type GameResponse struct {
	Error           error
	PlayerID        string
	G               Game
	Ready           bool
	UsingFriendCode bool
	Message         string
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

var ErrNOTREADYYET = errors.New("still waiting")

func (gss *GameStateSingleton) StartProcessing() {
	for {
		select {
		case req := <-gss.SinglePlayerRequests:
			{
				playerID := uuid()
				gameID := uuid()
				g := NewGame(gameID)
				g.LeftPlayerID = playerID
				g.RightPlayerID = "ROBOT"
				gss.games[gameID] = g
				g.start()
				g.LeftPlayerID = playerID
				req.Res <- GameResponse{
					Error:    nil,
					PlayerID: playerID,
					G:        *g,
					Ready:    true,
				}
			}
		case req := <-gss.CancelRequests:
			{
				if gss.matchMakingWaitingRoom != nil && gss.matchMakingWaitingRoom.ID == req.ID {
					gss.matchMakingWaitingRoom = nil
				}
				_, ok := gss.games[req.ID]
				if !ok {
					goto end
				}
				delete(gss.games, req.ID)
			}
		case req := <-gss.NewMatchMakingRequests:
			if gss.matchMakingWaitingRoom == nil {
				gameID := uuid()
				g := NewGame(gameID)
				gss.games[gameID] = g
				g.LeftPlayerID = uuid()
				req.Res <- GameResponse{
					G:        *g,
					Ready:    false,
					PlayerID: g.LeftPlayerID,
				}
				gss.matchMakingWaitingRoom = g
			} else {
				rightPlayerID := uuid()
				rightPlayerIDString := rightPlayerID
				gss.matchMakingWaitingRoom.RightPlayerID = rightPlayerIDString
				gss.matchMakingWaitingRoom.start()
				req.Res <- GameResponse{
					G:        *gss.matchMakingWaitingRoom,
					Ready:    true,
					PlayerID: rightPlayerIDString,
				}
				gss.matchMakingWaitingRoom = nil
			}
		case req := <-gss.NewWaitingRoomRequests:
			{
				idStruct := uuid()
				gameID := idStruct
				g := NewGame(gameID)
				gss.games[gameID] = g
				g.LeftPlayerID = uuid()
				req.Res <- GameResponse{
					G:               *g,
					Ready:           false,
					PlayerID:        g.LeftPlayerID,
					UsingFriendCode: true,
				}
			}
		case req := <-gss.FriendJoinRequests:
			{
				g, ok := gss.games[req.ID]
				if !ok {
					req.Res <- GameResponse{
						Error: errors.New("unable to find game of that id"),
					}
					goto end
				}
				if g.GameState != WAITING {
					req.Res <- GameResponse{
						Error: errors.New("game already full"),
					}
					goto end
				}

				rightPlayerIDString := uuid()
				g.RightPlayerID = rightPlayerIDString
				g.start()
				req.Res <- GameResponse{
					G:        *g,
					Ready:    true,
					PlayerID: rightPlayerIDString,
				}

			}

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
							Error:           nil,
							PlayerID:        req.PlayerID,
							G:               *g,
							Ready:           false,
							UsingFriendCode: false,
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
					inFriendCodeMode := false
					if gss.matchMakingWaitingRoom == nil || gss.matchMakingWaitingRoom.ID != req.ID {
						inFriendCodeMode = true

					}
					req.Res <- GameResponse{
						Error:           ErrNOTREADYYET,
						PlayerID:        req.PlayerID,
						G:               *g,
						Ready:           false,
						UsingFriendCode: inFriendCodeMode,
					}
					goto end
				}
			default:
				{
					//do nothing and allow app to continue
				}
			}
			g.play(req.A, req.PlayerID)
			message := ""
			switch g.GameState {
			case RIGHT_WIN:
				if req.PlayerID == g.RightPlayerID {
					message = "YOU WIN!"
				} else {
					message = "YOU LOSE!"
				}
			case LEFT_WIN:
				if req.PlayerID == g.LeftPlayerID {
					message = "YOU WIN! (reload the page to go again)"
				} else {
					message = "YOU LOSE!(reload the page to go again)"
				}
			}
			req.Res <- GameResponse{
				PlayerID: req.PlayerID,
				G:        *g,
				Ready:    true,
				Error:    nil,
				Message:  message,
			}
		}
	end:
	}
}

var uuidalphabet []string = []string{
	"a",
	"b",
	"c",
	"d",
	"e",
	"f",
	"g",
	"h",
	"i",
	"j",
	"k",
	"m",
	"n",
	"p",
	"q",
	"r",
	"s",
	"t",
	"u",
	"v",
	"w",
	"x",
	"y",
	"z",
	"2",
	"3",
	"4",
	"5",
	"6",
	"7",
	"8",
	"9",
}

const lenUUID = 6

func uuid() string {
	ret := []string{}
	for i := 0; i < lenUUID; i++ {
		index := rand.Intn(len(uuidalphabet))
		ret = append(ret, uuidalphabet[index])
	}

	return strings.Join(ret, "")
}
