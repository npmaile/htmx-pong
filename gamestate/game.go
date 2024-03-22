package gamestate

import (
	"math"
	"time"
)

const inputScaling float64 = 1.0
const timeScaling float64 = 0.1

type vec2 struct {
	X float64
	Y float64
}

type ballstate struct {
	Speed vec2
	Loc   vec2
}

type paddle struct {
	Height float64
	Y      float64
	Left   bool
}

type Game struct {
	Updated   time.Time
	Ball      ballstate
	PaddR     paddle
	Paddl     paddle
	ScoreL    int
	ScoreR    int
	Completed bool
}

func NewGame() *Game {
	return &Game{
		Ball: ballstate{Speed: vec2{X: 2, Y: 1}, Loc: vec2{X: 50, Y: 50}},
		Paddl: paddle{
			Height: 20,
			Y:      50,
			Left:   true,
		},
		PaddR: paddle{
			Height: 20,
			Y:      50,
			Left:   false,
		},
		Updated:   time.Now(),
		Completed: false,
	}
}

type Action int

const (
	lup Action = iota
	ldown
	rup
	rdown
)

func (g *Game) play(action Action) {
	// move paddles
	switch action {
	case lup:
		g.Paddl.Y -= 1 * inputScaling
	case ldown:
		g.Paddl.Y += 1 * inputScaling
	case rup:
		g.PaddR.Y -= 1 * inputScaling
	case rdown:
		g.PaddR.Y += 1 * inputScaling
	default:
	}

	// move balls
	delta := float64(time.Since(g.Updated)) * timeScaling
	g.Ball.Loc.X += g.Ball.Speed.X * delta
	g.Ball.Loc.Y += g.Ball.Speed.Y * delta

	// calculate ceiling/floor collisions
	if math.Abs(g.Ball.Loc.Y) >= 100.0 {
		g.Ball.Speed.Y *= -1
	}

	// calculate paddle collisions
	if math.Abs(g.Ball.Loc.X) >= 99.0 {
	}
}
