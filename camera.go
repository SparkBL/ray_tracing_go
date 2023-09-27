package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"ray_tracing/ray"
	"ray_tracing/vector"
	"time"
)

type Camera struct {
	aspectRatio       float64
	imageWidth        int
	imageHeight       int
	center            vector.Point
	pixelZeroLocation vector.Vector
	pixelDeltaU       vector.Vector
	pixelDeltaV       vector.Vector

	focalLength    float64
	viewportHeight float64
	viewportWidth  float64
	logger         *bufio.Writer
}

type CameraOption func(c *Camera)

func (c *Camera) Init(opts ...CameraOption) {
	c.aspectRatio = 16.0 / 9.0
	c.imageWidth = 600
	c.imageHeight = int(float64(c.imageWidth) / c.aspectRatio)

	c.focalLength = 1.0
	c.viewportHeight = 2.0 //think as matrix
	c.viewportWidth = c.viewportHeight * (float64(c.imageWidth) / float64(c.imageHeight))

	c.center = vector.Point{0, 0, 0}

	for _, o := range opts {
		o(c)
	}

	if c.imageHeight < 1 {
		c.imageHeight = 1
	}
	viewportU := vector.Vector{c.viewportWidth, 0, 0}
	viewportV := vector.Vector{0, -c.viewportHeight, 0}

	c.pixelDeltaU = viewportU.Divide(float64(c.imageWidth))
	c.pixelDeltaV = viewportV.Divide(float64(c.imageHeight))

	viewPortUpleft :=
		c.center.Add(vector.Vector{0, 0, c.focalLength}.Negative()).
			Add(viewportU.Divide(2).Negative()).
			Add(viewportV.Divide(2).Negative())

	c.pixelZeroLocation = viewPortUpleft.Add(c.pixelDeltaU.Add(c.pixelDeltaV).Multiply(0.5))
	c.logger = bufio.NewWriter(os.Stdout)
}

func (c *Camera) Render(filename string, world Hittable) {
	start := time.Now()
	err := os.Remove(filename)
	if err != nil {
		log.Println(err)
	}
	out, _ := os.Create(filename)

	out.WriteString(fmt.Sprintf("P3\n%d %d\n255\n", c.imageWidth, c.imageHeight))

	for i := 0; i < c.imageHeight; i++ {

		c.logger.WriteString(fmt.Sprintf("\rremaining: %.2f%%", 100.0*float64(c.imageHeight-i)/float64(c.imageHeight)))
		c.logger.Flush()

		for j := 0; j < c.imageWidth; j++ {
			pixelCenter := c.pixelZeroLocation.Add(c.pixelDeltaU.Multiply(float64(j)).Add(c.pixelDeltaV.Multiply(float64(i))))
			rayDirection := pixelCenter.Add(c.center.Negative())
			r := ray.Ray{Origin: c.center, Direction: rayDirection}
			color := c.rayColor(r, world)
			out.WriteString(vector.ColorString(&color))
		}
	}
	out.Close()
	c.logger.WriteString(fmt.Sprintf("\relapsed: %v\n", time.Since(start)))
	c.logger.Flush()
}

func (c *Camera) rayColor(r ray.Ray, world Hittable) vector.Color {
	rec := HitRecord{}

	if world.Hit(r, Interval{0, math.Inf(1)}, &rec) {
		return rec.Normal.Add(vector.Color{1, 1, 1}).Multiply(0.5)
	}
	unitDirection := vector.UnitVector(r.Direction)
	a := 0.5 * (unitDirection.Y() + 1.0)
	return vector.Color{1.0, 1.0, 1.0}.
		Multiply(1.0 - a).
		Add(vector.Color{0.5, 0.7, 1.0}.Multiply(a))
}
