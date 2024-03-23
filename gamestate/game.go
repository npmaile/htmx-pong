package gamestate

import (
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
	Updated   time.Time
	Ball      *ballstate
	PaddR     *paddle
	Paddl     *paddle
	ScoreL    int
	ScoreR    int
	Completed bool
	ID        string
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
		Completed: false,
		ID:        ID,
	}
}

type Action int

const (
	Lup Action = iota
	Ldown
	Rup
	Rdown
	NoAction
)

func (g *Game) play(action Action) {
	// move paddles
	switch action {
	case Lup:
		g.Paddl.Y -= 1 * inputScaling
	case Ldown:
		g.Paddl.Y += 1 * inputScaling
	case Rup:
		g.PaddR.Y -= 1 * inputScaling
	case Rdown:
		g.PaddR.Y += 1 * inputScaling
	default:
	}

	// move balls
	delta := float64(time.Since(g.Updated)) * timeScaling
	g.Ball.Loc.X += g.Ball.Speed.X * delta
	g.Ball.Loc.Y += g.Ball.Speed.Y * delta

	// calculate ceiling/floor collisions
	if math.Abs(g.Ball.Loc.Y-50) >= 50 {
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
	} else if math.Abs(g.Ball.Loc.X-50) >= 43 && g.Ball.Speed.X > 0 {
		if math.Abs(g.Ball.Loc.Y-g.PaddR.Y) < g.PaddR.Height/2 {
			g.Ball.Speed.X *= -1
		} else {

			g.ScoreL += 1
			g.resetBall(true)
		}
	}
	g.Updated = time.Now()
}

func (g *Game) resetBall(goingLeft bool) {
	g.Ball.Loc.X = 50
	g.Ball.Loc.Y = 50
	g.Ball.Speed.Y = 1
	if goingLeft {
		g.Ball.Speed.X = 2
	} else {
		g.Ball.Speed.Y = -2
	}
}
