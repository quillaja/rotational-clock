package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/lucasb-eyer/go-colorful"
)

var theme image.Image
var winSize pixel.Vec
var config Configuration

func run() {

	// window and texture sizing info
	textureSize := pixel.V(float64(theme.Bounds().Dx()/6), float64(theme.Bounds().Dy()))

	if winSize == pixel.ZV {
		maxtex := math.Max(textureSize.XY())
		winSize = pixel.V(maxtex, maxtex)
	}
	textureScale := math.Min(winSize.XY()) / math.Max(textureSize.XY())
	numFaces := 6

	// configure and create window
	cfg := pixelgl.WindowConfig{
		Title:  "Clock",
		Bounds: pixel.R(0, 0, winSize.X, winSize.Y),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.SetSmooth(true)

	// create graphical resources
	picdata := pixel.PictureDataFromImage(theme)
	batch := pixel.NewBatch(&pixel.TrianglesData{}, picdata)

	// set up clock data
	clock := clock{}
	{ // initialize clock sprites and positions
		sprites := make([]*pixel.Sprite, numFaces)
		clip := pixel.R(0, 0, textureSize.X, textureSize.Y)
		for i := range sprites {
			sprites[i] = pixel.NewSprite(picdata, clip)
			clip.Min.X += textureSize.X
			clip.Max.X += textureSize.X
		}

		spritePos := make(positions, numFaces)

		clock.faces = sprites
		clock.positions = spritePos
	}

	// vars used in loop
	targetFPS := 1.0 / float64(config.TargetFPS)
	focus := Year

	var bgcolor color.Color
	bgcolor, err = colorful.Hex(config.BackgroundColor)
	if err != nil {
		bgcolor = color.White
	}

	for !win.Closed() {
		start := time.Now()

		// recalculate sprite positions and rotations
		clock.positions.calculate(start, textureScale*textureSize.Y/2, config.RotationMode)

		// draw each sprite
		batch.Clear()
		for i := range clock.faces {
			clock.faces[i].Draw(batch, pixel.IM.
				Rotated(pixel.ZV, clock.positions[i].angle*float64(config.RotationDirection)).
				Scaled(pixel.ZV, textureScale*powHalf(i)).
				Moved(clock.positions[i].position))
		}

		// position clock within window
		getFocus(&focus, win)
		win.SetMatrix(pixel.IM.
			Moved(win.Bounds().Center()).
			Moved(clock.positions[focus].position.Scaled(-1)).
			Scaled(win.Bounds().Center(), pow2(focus)))

		// draw to window
		win.Clear(bgcolor)
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

// getFocus gets the clock component to 'zoom' from user keypresses
// using the buttons defined in Focus[].
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
	size := flag.String("size", "", "width and height of window (eg 800x600)")
	flag.Parse()

	if *folder == "" {
		fmt.Println("'folder' flag must be specified.")
		return
	}

	// parse size
	if *size != "" {
		dims := strings.Split(*size, "x")
		var w, h int
		var err error
		if len(dims) == 2 {
			w, err = strconv.Atoi(strings.TrimSpace(dims[0]))
			h, err = strconv.Atoi(strings.TrimSpace(dims[1]))
			if err != nil {
				fmt.Println("error with window size of", *size)
			}
		}
		winSize = pixel.V(float64(w), float64(h))
	}

	// prepare theme
	theme, err = MakeTile(*folder)
	if err != nil {
		panic(err)
	}

	config = DefaultConfig()

	pixelgl.Run(run)
}
