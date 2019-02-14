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

const TwoPi = math.Pi * 2

var theme image.Image

func run() {

	textureSize := float64(theme.Bounds().Dy())

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

	sprites := make([]*pixel.Sprite, len(files))
	clip := pixel.R(0, 0, textureSize, textureSize)
	for i := range sprites {
		sprites[i] = pixel.NewSprite(picdata, clip)
		clip.Min.X += textureSize
		clip.Max.X += textureSize
	}

	for !win.Closed() {

		now := time.Now().Second()
		rot := float64(now) / 60.0
		for i := range sprites {
			sprites[i].Draw(batch, pixel.IM.
				Rotated(pixel.ZV, -TwoPi*rot).
				Scaled(pixel.ZV, math.Pow(2, float64(-i))).
				Moved(pixel.ZV))
		}

		win.SetMatrix(pixel.IM.Moved(win.Bounds().Center()))

		win.Clear(color.Black)
		batch.Draw(win)
		win.Update()
	}
}

type pos struct {
	pixel.Vec
	angle float64
}

type positions []pos

func (p positions) calculate() {
	// now := time.Now()
	// n := len(p)

}

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
