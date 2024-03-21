package gamestate

import (
	"math"
	"time"
)

const inputScaling float64 = 1.0
const timeScaling float64 = 1.0

type vec2 struct {
	x float64
	y float64
}

type ballstate struct {
	speed vec2
	loc   vec2
}

type paddle struct {
	height float64
	y      float64
	left   bool
}

type Game struct {
	paddR     paddle
	paddl     paddle
	ball      ballstate
	updated   time.Time
	scoreL    int
	scoreR    int
	completed bool
}

func NewGame() *Game {
	return &Game{
		ball: ballstate{speed: vec2{x: 2, y: 1}, loc: vec2{x: 50, y: 50}},
		paddl: paddle{
			height: 20,
			y:      50,
			left:   true,
		},
		paddR: paddle{
			height: 20,
			y:      50,
			left:   false,
		},
		updated:   time.Now(),
		completed: false,
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
		g.paddl.y -= 1 * inputScaling
	case ldown:
		g.paddl.y += 1 * inputScaling
	case rup:
		g.paddR.y -= 1 * inputScaling
	case rdown:
		g.paddR.y += 1 * inputScaling
	default:
	}

	// move balls
	delta := float64(time.Now().Sub(g.updated)) * timeScaling
	g.ball.loc.x += g.ball.speed.x * delta
	g.ball.loc.y += g.ball.speed.y * delta

	// calculate ceiling/floor collisions
	if math.Abs(g.ball.loc.y) >= 100.0 {
		g.ball.speed.y *= -1
	}

	// calculate paddle collisions
	if math.Abs(g.ball.loc.x) >= 99.0 {
	}
}
