package main

import (
	"bufio"
	"fmt"
	"os"
	"ray_tracing/ray"
	"ray_tracing/vector"

	"log"
)

func main() {

	//Image
	aspectRatio := 16.0 / 9.0
	imageWidth := 400

	//Calc image height
	imageHeight := int(float64(imageWidth) / aspectRatio)

	//Camera
	focalLength := 1.0
	viewportHeight := 2.0 //think as matrix
	viewportWidth := viewportHeight * (float64(imageWidth) / float64(imageHeight))
	camerCenter := vector.Point{0, 0, 0}
	// Calculate the vectors across the horizontal and down the vertical viewport edges.
	viewportU := vector.Vector{viewportWidth, 0, 0}
	viewportV := vector.Vector{0, -viewportHeight, 0}
	// Calculate the horizontal and vertical delta vectors from pixel to pixel.
	pixelDeltaU := viewportU.Divide(float64(imageWidth))
	pixelDeltaV := viewportV.Divide(float64(imageHeight))

	viewPortUpleft :=
		camerCenter.Add(vector.Vector{0, 0, focalLength}.Negative()).
			Add(viewportU.Divide(2).Negative()).
			Add(viewportV.Divide(2).Negative())
	pixel00Location := viewPortUpleft.Add(pixelDeltaU.Add(pixelDeltaV).Multiply(0.5))

	//Render
	//OutputImage(256, 256)
	OutputRenderedImage(RenderImageOption{
		PixelZeroLocation: pixel00Location,
		CameraCenter:      camerCenter,
		PixelDeltaU:       pixelDeltaU,
		PixelDeltaV:       pixelDeltaV,
		Width:             imageWidth,
		Height:            imageHeight,
	})
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

func RayColor(r ray.Ray) vector.Color {
	if HitSphere(vector.Point{0, 0, -1}, 0.5, r) {
		return vector.Color{1, 0, 0}
	}
	unitDirection := vector.UnitVector(r.Direction)
	a := 0.5 * (unitDirection.Y() + 1.0)
	return vector.Color{1.0, 1.0, 1.0}.Multiply(1.0 - a).Add(vector.Color{0.5, 0.7, 1.0}.Multiply(a))
}

func HitSphere(center vector.Point, radius float64, ray ray.Ray) bool {
	ocDistance := ray.Origin.Add(center.Negative())
	a := vector.Dot(ray.Direction, ray.Direction)
	b := 2.0 * vector.Dot(ocDistance, ray.Direction)
	c := vector.Dot(ocDistance, ocDistance) - radius*radius
	discriminant := b*b - 4.0*a*c
	return discriminant >= 0
}

type RenderImageOption struct {
	PixelZeroLocation vector.Vector
	CameraCenter      vector.Vector
	PixelDeltaU       vector.Vector
	PixelDeltaV       vector.Vector
	Width             int
	Height            int
}

func OutputRenderedImage(o RenderImageOption) {
	err := os.Remove("test_ray.ppm")
	if err != nil {
		log.Println(err)
	}

	out, _ := os.Create("test_ray.ppm")
	logger := bufio.NewWriter(os.Stdout)
	out.WriteString(fmt.Sprintf("P3\n%d %d\n255\n", o.Width, o.Height))
	for i := 0; i < o.Height; i++ {
		logger.WriteString(fmt.Sprintf("\rScanlines remaining: %d ", o.Height-i))
		logger.Flush()
		for j := 0; j < o.Width; j++ {
			pixelCenter := o.PixelZeroLocation.Add(o.PixelDeltaU.Multiply(float64(j)).Add(o.PixelDeltaV.Multiply(float64(i))))
			rayDirection := pixelCenter.Add(o.CameraCenter.Negative())
			r := ray.Ray{Origin: o.CameraCenter, Direction: rayDirection}
			color := RayColor(r)
			out.WriteString(vector.ColorString(&color))
		}
	}
	out.Close()
	logger.Flush()
}
