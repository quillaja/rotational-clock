package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"math"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var theme image.Image

func run() {

	textureSize := float64(theme.Bounds().Dy())
	numFaces := theme.Bounds().Dx() / theme.Bounds().Dy()
	fmt.Println(numFaces)

	cfg := pixelgl.WindowConfig{
		Title:  "Clock",
		Bounds: pixel.R(0, 0, textureSize, textureSize),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	picdata := pixel.PictureDataFromImage(theme)
	batch := pixel.NewBatch(&pixel.TrianglesData{}, picdata)

	clock := clock{}
	{ // initialize clock sprites and positions
		sprites := make([]*pixel.Sprite, len(files))
		clip := pixel.R(0, 0, textureSize, textureSize)
		for i := range sprites {
			sprites[i] = pixel.NewSprite(picdata, clip)
			clip.Min.X += textureSize
			clip.Max.X += textureSize
		}

		spritePos := make(positions, len(files))

		clock.faces = sprites
		clock.positions = spritePos
	}

	for !win.Closed() {

		// recalculate sprite positions and rotations
		clock.positions.calculate(textureSize / 2)

		// draw each sprite
		batch.Clear()
		for i := range clock.faces {
			clock.faces[i].Draw(batch, pixel.IM.
				Rotated(pixel.ZV, -clock.positions[i].angle).
				Scaled(pixel.ZV, powHalf(i)).
				Moved(clock.positions[i].position))
		}

		// position clock within window
		win.SetMatrix(pixel.IM.
			Moved(win.Bounds().Center()).
			Moved(clock.positions[Month].position.Scaled(-1)).
			Scaled(win.Bounds().Center(), pow2(Month)))

		// draw to window
		win.Clear(color.Black)
		batch.Draw(win)
		win.Update()
	}
}

// powHalf calculates 0.5^(i)
func powHalf(i int) float64 {
	return math.Pow(2, float64(-i))
}

// pow2 calculates 2^i
func pow2(i int) float64 {
	return math.Pow(2, float64(i))
}

// clock is a convenience struct to contain all
// data needed to represent the clock
type clock struct {
	faces []*pixel.Sprite
	positions
}

// pos contains the positional data for a sprite.
type pos struct {
	position pixel.Vec
	angle    float64
}

// positions is a slice of the positional data for sprites.
type positions []pos

// calculate uses the current time to determine the position and rotation
// for each of the len(p) sprites.
func (p positions) calculate(radius float64) {
	now := time.Now()
	n := len(p)

	for i := 0; i < n; i++ {
		angle := 0.0

		switch i {
		case Year:
			angle = float64(now.YearDay()-1) / 365.0
		case Month:
			// uses idiosyncrasy of time package to figure out the number
			// of days in the current month.
			// see: https://yourbasic.org/golang/last-day-month-date/
			daysInMonth := time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, time.Local).Day()
			angle = float64(now.Day()-1) / float64(daysInMonth)

		case Day:
			angle = float64(now.Hour()) / 24.0 // Hour() returns [0,23]

		case Hour:
			angle = float64(now.Minute()) / 60.0 // Minute() [0,59]

		case Minute:
			angle = float64(now.Second()) / 60.0 // Second() [0,59]

		case Second:
			angle = float64(now.Nanosecond()) / 1e9
		}

		p[i].angle = angle * TwoPi
		p[i].position = p.locForIndex(i, radius)
	}

}

// locForIndex calculates the position of the component at i given the initial
// radius.
//    loc = prev_loc + radius * 2^(-i) * Vec2((cos&sin(prev_angle))
// except when i==0, loc = (0,0)
func (p positions) locForIndex(i int, radius float64) pixel.Vec {
	if i == 0 {
		return pixel.ZV
	}

	dirUnitCircle := pixel.V(math.Cos(p[i-1].angle), math.Sin(p[i-1].angle))
	return p[i-1].position.Add(dirUnitCircle.Scaled(radius * powHalf(i)))
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
