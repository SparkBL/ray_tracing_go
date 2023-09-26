package main

import (
	"bufio"
	"fmt"
	"os"
	"ray_tracing/vector"
)

func main() {
	OutputImage(256, 256)

	aspectRatio := 16.0 / 9.0
	imageWidth := 400
	imageHeight := int(float64(imageWidth) / aspectRatio)
	viewportHeight := 2.0
	viewportWidth := viewportHeight * float64(imageWidth) / float64(imageHeight)
	camerCenter := vector.Point{0, 0, 0}
	focalLength := 1.0

	viewportU := vector.Vector{viewportWidth, 0, 0}
	viewportV := vector.Vector{0, -viewportHeight, 0}

	pixelDeltaU := viewportU.Divide(float64(imageWidth))
	pixelDeltaV := viewportU.Divide(float64(imageHeight))
}

func OutputImage(width, height int) {
	out, _ := os.Create("test.ppm")
	logger := bufio.NewWriter(os.Stdout)
	out.WriteString(fmt.Sprintf("P3\n%d %d\n255\n", width, height))
	for i := 0; i < height; i++ {
		logger.WriteString(fmt.Sprintf("\rScanlines remaining: %d ", height-i))
		logger.Flush()
		for j := 0; j < width; j++ {
			c := vector.Vector{
				float64(j) / (float64(width) - 1),
				float64(i) / (float64(height) - 1),
				float64(0)}

			out.WriteString(vector.ColorString(&c))
		}
	}
	out.Close()
	logger.Flush()
}
