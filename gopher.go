package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type gopher struct {
	sceneBounds pixel.Rect
	sprite      *pixel.Sprite
	position    pixel.Vec
	dead        bool
	speed       float64
}

func newGopher(sceneBounds pixel.Rect) *gopher {
	pic, err := loadPicture("images/gopher.png")
	if err != nil {
		panic(err)
	}

	sprite := pixel.NewSprite(pic, pic.Bounds())

	startPosition := pixel.V(pic.Bounds().Center().X, sceneBounds.Center().Y)

	return &gopher{
		sceneBounds: sceneBounds,
		sprite:      sprite,
		position:    startPosition,
		speed:       5,
	}
}

func (g *gopher) update() {
	g.position.Y -= g.speed
	if g.rect().Min.Y <= 0 || g.rect().Max.Y > g.sceneBounds.Max.Y {
		g.dead = true
	}
}

func (g *gopher) reborn() {
	g.position.Y = g.sceneBounds.Center().Y
	g.dead = false
	g.speed = 5
}

func (g *gopher) rect() pixel.Rect {
	rect := g.sprite.Picture().Bounds()
	return rect.Moved(pixel.V(0, g.position.Y-rect.Center().Y))
}

func (g *gopher) touch(p *pipe) {
	if p.rect().Intersect(g.rect()) != pixel.R(0, 0, 0, 0) {
		g.dead = true
	}
}

func (g *gopher) draw(win *pixelgl.Window) {
	g.sprite.Draw(win, pixel.IM.Moved(g.position))
}
