package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var theme image.Image

func run() {

	cfg := pixelgl.WindowConfig{
		Title:  "Clock",
		Bounds: pixel.R(0, 0, float64(theme.Bounds().Dy()), float64(theme.Bounds().Dy())),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	win.Clear(color.Black)

	for !win.Closed() {
		win.Update()
	}
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
