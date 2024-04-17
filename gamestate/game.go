package gamestate

import (
	"errors"
	"fmt"
	"math"
	"time"
)

const inputScaling float64 = 4.0
const timeScaling float64 = 0.00000001

type vec2 struct {
	X float64
	Y float64
}

type ballstate struct {
	Speed *vec2
	Loc   *vec2
}

type paddle struct {
	Height float64
	Y      float64
	Left   bool
}

type Game struct {
	Updated       time.Time
	Ball          *ballstate
	PaddR         *paddle
	Paddl         *paddle
	ID            string
	LeftPlayerID  string
	RightPlayerID string
	ScoreL        int
	ScoreR        int
	GameState     state
}

type state int

const (
	STARTED state = iota
	WAITING
	LEFT_WIN
	RIGHT_WIN
)

type GameUpdate struct {
	PlayerID string
	Game
}

func NewGame(ID string) *Game {
	return &Game{
		Ball: &ballstate{Speed: &vec2{X: -2, Y: 1}, Loc: &vec2{X: 50, Y: 50}},
		Paddl: &paddle{
			Height: 25,
			Y:      50,
			Left:   true,
		},
		PaddR: &paddle{
			Height: 25,
			Y:      50,
			Left:   false,
		},
		Updated:   time.Now(),
		GameState: WAITING,
		ID:        ID,
	}
}

func (g *Game) start() {
	g.GameState = STARTED
}

type Action int

const (
	Up Action = iota
	Down
	NoAction
)

func (g *Game) play(action Action, playerID string) error {
	switch g.GameState {
	case WAITING:
		return errors.New("game not started")
	case LEFT_WIN:
		return nil
	case RIGHT_WIN:
		return nil
	case STARTED:
		//continue on
	default:
		return errors.New("wat")
	}
	// move paddles
	var targetPaddle *paddle
	switch playerID {
	case g.LeftPlayerID:
		fmt.Println("left player found")
		targetPaddle = g.Paddl
	case g.RightPlayerID:
		fmt.Println("right player found")
		targetPaddle = g.PaddR
	default:
		return errors.New("no user found")

	}

	switch action {
	case Up:
		fmt.Println("moving up")
		targetPaddle.Y -= 1 * inputScaling
	case Down:
		fmt.Println("moving down")
		targetPaddle.Y += 1 * inputScaling
	default:
	}

	// move balls
	delta := float64(time.Since(g.Updated)) * timeScaling
	g.Ball.Loc.X += g.Ball.Speed.X * delta
	g.Ball.Loc.Y += g.Ball.Speed.Y * delta

	// calculate ceiling/floor collisions
	if g.Ball.Loc.Y >= 100 && g.Ball.Speed.Y > 0 {
		g.Ball.Speed.Y *= -1
	}

	if g.Ball.Loc.Y <= 0 && g.Ball.Speed.Y < 0 {
		g.Ball.Speed.Y *= -1
	}

	// calculate paddle collision left
	if math.Abs(g.Ball.Loc.X-50) >= 43 && g.Ball.Speed.X < 0 {
		if math.Abs(g.Ball.Loc.Y-g.Paddl.Y) < g.Paddl.Height/2 {
			g.Ball.Speed.X *= -1
		} else {
			g.ScoreR += 1
			g.resetBall(false)
		}
		// calculate paddle collision right
	} else if math.Abs(g.Ball.Loc.X-50) >= 43 && g.Ball.Speed.X > 0 {
		if math.Abs(g.Ball.Loc.Y-g.PaddR.Y) < g.PaddR.Height/2 {
			g.Ball.Speed.X *= -1
		} else {

			g.ScoreL += 1
			g.resetBall(true)
		}
	}
	g.Updated = time.Now()

	return nil
}

func (g *Game) resetBall(goingLeft bool) {
	g.Ball.Loc.X = 50
	g.Ball.Loc.Y = 50
	g.Ball.Speed.Y = 1
	if !goingLeft {
		g.Ball.Speed.X = 2
	} else {
		g.Ball.Speed.Y = -2
	}
}
