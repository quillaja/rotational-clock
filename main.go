package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var theme image.Image

func run() {

	textureSize := float64(theme.Bounds().Dy())
	numFaces := theme.Bounds().Dx() / theme.Bounds().Dy()

	cfg := pixelgl.WindowConfig{
		Title:  "Clock",
		Bounds: pixel.R(0, 0, textureSize, textureSize),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.SetSmooth(true)

	picdata := pixel.PictureDataFromImage(theme)
	batch := pixel.NewBatch(&pixel.TrianglesData{}, picdata)

	clock := clock{}
	{ // initialize clock sprites and positions
		sprites := make([]*pixel.Sprite, numFaces)
		clip := pixel.R(0, 0, textureSize, textureSize)
		for i := range sprites {
			sprites[i] = pixel.NewSprite(picdata, clip)
			clip.Min.X += textureSize
			clip.Max.X += textureSize
		}

		spritePos := make(positions, numFaces)

		clock.faces = sprites
		clock.positions = spritePos
	}

	const targetFPS = 1.0 / 60.0
	focus := Year

	for !win.Closed() {
		start := time.Now()

		// recalculate sprite positions and rotations
		clock.positions.calculate(start, textureSize/2)

		// draw each sprite
		batch.Clear()
		for i := range clock.faces {
			clock.faces[i].Draw(batch, pixel.IM.
				Rotated(pixel.ZV, -clock.positions[i].angle).
				Scaled(pixel.ZV, powHalf(i)).
				Moved(clock.positions[i].position))
		}

		// position clock within window
		getFocus(&focus, win)
		win.SetMatrix(pixel.IM.
			Moved(win.Bounds().Center()).
			Moved(clock.positions[focus].position.Scaled(-1)).
			Scaled(win.Bounds().Center(), pow2(focus)))

		// draw to window
		win.Clear(color.Black)
		batch.Draw(win)
		win.Update()

		// control FPS
		elapsed := time.Since(start)
		if elapsed.Seconds() < targetFPS {
			pause := targetFPS - elapsed.Seconds()
			time.Sleep(time.Duration(pause * float64(time.Second)))
		}
	}
}

func getFocus(focus *int, win *pixelgl.Window) {
	for i := Year; i <= Second; i++ {
		if win.JustPressed(Focus[i]) {
			*focus = i
			return
		}
	}
}

// main parses flags, loads the texture/image/sprite/whatever data from
// disk and then runs the pixel window junk
func main() {
	var err error

	folder := flag.String("folder", "", "relative path to folder containing theme to use.")
	flag.Parse()

	if *folder == "" {
		fmt.Println("'folder' flag must be specified.")
		return
	}

	theme, err = MakeTile(*folder)
	if err != nil {
		panic(err)
	}

	pixelgl.Run(run)
}
