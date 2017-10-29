package main

import (
	"fmt"
	"image"
	_ "image/png"
	"os"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return pixel.PictureDataFromImage(img), nil
}

func run() {
	pic, err := loadPicture("images/bg.png")
	if err != nil {
		panic(err)
	}

	cfg := pixelgl.WindowConfig{
		Title:  "Test PixelGL",
		Bounds: pic.Bounds(),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	deadTxt := text.New(pixel.V(400, 400), text.NewAtlas(basicfont.Face7x13, text.ASCII))
	deadTxt.Color = colornames.Red
	fmt.Fprintln(deadTxt, "Game Over !")
	fmt.Fprintln(deadTxt, "Press enter")

	bg := pixel.NewSprite(pic, pic.Bounds())
	bg.Draw(win, pixel.IM.Moved(win.Bounds().Center()))

	sceneBounds := win.Bounds()
	gopher := newGopher(sceneBounds)
	pipes := newPipes(sceneBounds)

	for !win.Closed() {
		win.Clear(colornames.White)
		bg.Draw(win, pixel.IM.Moved(win.Bounds().Center()))

		gopher.update()
		if gopher.dead {
			pipes.clear()
			deadTxt.Draw(win, pixel.IM.Scaled(deadTxt.Orig, 10))
			if win.JustPressed(pixelgl.KeyEnter) {
				gopher.reborn()
			}

		} else {
			pipes.update(gopher)
			pipes.draw(win)

			if win.JustPressed(pixelgl.KeyUp) {
				gopher.speed = -5
			} else {
				gopher.speed += 0.1
			}

			// Space debug :)
			if win.JustPressed(pixelgl.KeySpace) {
				fmt.Println("=== Not Dead ===")
				fmt.Printf("gopher: %v\n", gopher.rect())
				var closerPipe *pipe
				for _, p := range pipes.pipes {
					if p.position.X > 0 {
						if closerPipe == nil {
							closerPipe = p
						} else if p.position.X < closerPipe.position.X {
							closerPipe = p
						}
					}
				}
				fmt.Printf("pipe: %v\n", closerPipe.rect())
				fmt.Printf("Intersect: %v\n", closerPipe.rect().Intersect(gopher.rect()))
				time.Sleep(20 * time.Second)
			}

			gopher.draw(win)
		}

		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
