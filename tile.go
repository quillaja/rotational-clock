package main

import (
	"image"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
)

var files = []string{"year.png", "month.png", "day.png", "hour.png", "minute.png", "second.png"}

// MakeTile looks in folder for PNG files named "year, month, day, hour,
// minute, and second" (with .png extension) and combines them into one
// larger file. Each file must be the same dimensions. Error is returned
// if there is an error opening or decoding the any of the files.
func MakeTile(folder string) (image.Image, error) {

	// open all images and place into array
	images := make([]image.Image, len(files))
	for i, name := range files {
		fpath := filepath.Join(folder, name)
		file, err := os.Open(fpath)
		if err != nil {
			return nil, err
		}
		images[i], err = png.Decode(file)
		if err != nil {
			return nil, err
		}
		file.Close()
	}

	// copy each image into a region of a larger single image
	width := images[0].Bounds().Dx()
	height := images[0].Bounds().Dy()
	dest := image.NewRGBA(image.Rect(0, 0, width*len(images), height))
	rect := images[0].Bounds()
	for _, img := range images {
		draw.Draw(dest, rect, img, img.Bounds().Min, draw.Src)
		rect.Min.X += width
		rect.Max.X += width
	}

	return dest, nil
}
