package gamestate

import (
	"errors"
	"math"
	"time"
)

const inputScaling float64 = 4.0
const timeScaling float64 = 0.00000001
const speedupDuringGameConstant float64 = 0.1
const score2win int = 5

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
	Dir    direction
}

type direction int

const (
	UP direction = iota
	DOWN
	NEUTRAL
)

type Game struct {
	Updated        time.Time
	Ball           *ballstate
	PaddR          *paddle
	Paddl          *paddle
	ID             string
	LeftPlayerID   string
	RightPlayerID  string
	ScoreL         int
	ScoreR         int
	GameState      state
	MatchStartTime time.Time
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
			Dir:    NEUTRAL,
		},
		PaddR: &paddle{
			Height: 25,
			Y:      50,
			Left:   false,
			Dir:    NEUTRAL,
		},
		ScoreL:    0,
		ScoreR:    0,
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
	NotUp
	NotDown
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
		targetPaddle = g.Paddl
	case g.RightPlayerID:
		targetPaddle = g.PaddR
	default:
		return errors.New("no user found")

	}

	switch action {
	case Up:
		targetPaddle.Dir = UP

	case Down:
		targetPaddle.Dir = DOWN

	case NotUp:
		targetPaddle.Dir = NEUTRAL
	case NotDown:
		targetPaddle.Dir = NEUTRAL
	default:
	}

	timeInMatchScaling := 1.0 + (time.Since(g.MatchStartTime).Seconds() * speedupDuringGameConstant)
	// move balls
	delta := float64(time.Since(g.Updated)) * timeScaling * timeInMatchScaling
	g.Ball.Loc.X += g.Ball.Speed.X * delta
	g.Ball.Loc.Y += g.Ball.Speed.Y * delta

	// move paddles
	switch g.Paddl.Dir {
	case UP:
		g.Paddl.Y -= 1 * inputScaling * delta
	case DOWN:
		g.Paddl.Y += 1 * inputScaling * delta
	}

	switch g.PaddR.Dir {
	case UP:
		g.PaddR.Y -= 1 * inputScaling * delta * .001
	case DOWN:
		g.PaddR.Y += 1 * inputScaling * delta * .001
	}

	// calculate ceiling/floor collisions
	if g.Ball.Loc.Y >= 100 && g.Ball.Speed.Y > 0 {
		g.Ball.Speed.Y *= -1
	}

	if g.Ball.Loc.Y <= 0 && g.Ball.Speed.Y < 0 {
		g.Ball.Speed.Y *= -1
	}

	// calculate paddle collision left
	if g.Ball.Loc.X <= 7 && g.Ball.Speed.X < 0 {
		if math.Abs(g.Ball.Loc.Y-g.Paddl.Y) < g.Paddl.Height/2 && g.Ball.Speed.X < 0 {
			g.Ball.Speed.X *= -1
		} else if g.Ball.Speed.X < 0 {
			g.ScoreR += 1
			g.resetBall(false)
		}
		// calculate paddle collision right
	} else if g.Ball.Loc.X >= 93 && g.Ball.Speed.X > 0 {
		if math.Abs(g.Ball.Loc.Y-g.PaddR.Y) < g.PaddR.Height/2 && g.Ball.Speed.X > 0 {
			g.Ball.Speed.X *= -1
		} else if g.Ball.Speed.X > 0 {
			g.ScoreL += 1
			g.resetBall(true)
		}
	}
	g.Updated = time.Now()
	if g.RightPlayerID == "ROBOT" {
		g.PaddR.Y = g.Ball.Loc.Y
	}
	if g.ScoreR > score2win {
		g.GameState = RIGHT_WIN
	}

	if g.ScoreL > score2win {
		g.GameState = LEFT_WIN
	}
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
	g.MatchStartTime = time.Now()
}
