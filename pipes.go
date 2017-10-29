package main

import (
	"math"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type pipe struct {
	sprite   *pixel.Sprite
	position pixel.Vec
	reverse  bool
}

func newPipe(sceneBounds pixel.Rect) *pipe {
	pic, err := loadPicture("images/pipe.png")
	if err != nil {
		panic(err)
	}

	sprite := pixel.NewSprite(pic, pic.Bounds())

	reverse := rand.Intn(2) == 1

	var startY float64
	if reverse {
		startY = sceneBounds.Max.Y - pic.Bounds().Center().Y
	} else {
		startY = pic.Bounds().Center().Y
	}

	startPosition := pixel.V(sceneBounds.Max.X, startY)

	return &pipe{
		sprite:   sprite,
		position: startPosition,
		reverse:  reverse,
	}
}

func (p *pipe) rect() pixel.Rect {
	rect := p.sprite.Picture().Bounds()
	if p.reverse {
		rect = rect.Moved(pixel.V(p.position.X-rect.Center().X, p.position.Y-rect.Center().Y))
	} else {
		rect = rect.Moved(pixel.V(p.position.X-rect.Center().X, 0))
	}
	return rect
}

type pipes struct {
	pipes        []*pipe
	nextCreation time.Time
	sceneBounds  pixel.Rect
}

func newPipes(sceneBounds pixel.Rect) *pipes {
	return &pipes{
		pipes:        []*pipe{newPipe(sceneBounds)},
		nextCreation: time.Now(),
		sceneBounds:  sceneBounds,
	}
}

func (ps *pipes) clear() {
	ps.pipes = []*pipe{}
}

func (ps *pipes) update(g *gopher) {
	for _, p := range ps.pipes {
		p.position.X -= 5
	}

	now := time.Now()
	if now.After(ps.nextCreation) {
		ps.pipes = append(ps.pipes, newPipe(ps.sceneBounds))
		ps.nextCreation = now.Add(time.Duration(rand.Intn(2)+1) * time.Second)
	}

	for _, p := range ps.pipes {
		g.touch(p)
	}
}

func (ps *pipes) draw(win *pixelgl.Window) {
	for _, pipe := range ps.pipes {
		mat := pixel.IM.Moved(pipe.position)
		if pipe.reverse {
			mat = mat.Rotated(pipe.position, math.Pi)
		}
		pipe.sprite.Draw(win, mat)
	}
}
