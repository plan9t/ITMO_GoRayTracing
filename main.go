package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

type Vec3f struct {
	x, y, z float64
}

func render() {
	const width, height = 1024, 768
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for j := 0; j < height; j++ {
		for i := 0; i < width; i++ {
			img.Set(i, j, color.RGBA{
				R: uint8(math.Max(0, math.Min(255, 255*float64(j)/float64(height)))),
				G: uint8(math.Max(0, math.Min(255, 255*float64(i)/float64(width)))),
				B: 0,
				A: 255,
			})
		}
	}

	file, err := os.Create("out/out.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	png.Encode(file, img)
}

func main() {
	render()
}
