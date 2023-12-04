package main

import (
	"math"
	"math/rand"
	"ray_tracing/camera"
	"ray_tracing/hittable"
	"ray_tracing/vector"
	"runtime/debug"
)

func Scene1() {
	// World
	materialGround := hittable.Lambertian{Albedo: vector.Color{0.8, 0.8, 0.0}}
	materialCenter := hittable.Lambertian{Albedo: vector.Color{0.1, 0.2, 0.5}}
	//materialLeft := hittable.Metal{Albedo: vector.Color{1, 1, 1}, Fuzziness: 0.2}
	materialLeft := hittable.Dielectric{1.5}
	materialRight := hittable.Metal{Albedo: vector.Color{0.8, 0.6, 0.2}, Fuzziness: 0.0}

	world := hittable.NewWorld(
		hittable.NewSphere(
			vector.Point{0, -100.5, -1.0},
			100,
			&materialGround),
		hittable.NewSphere(
			vector.Point{0, 0, -1},
			0.5,
			&materialCenter),
		hittable.NewSphere(
			vector.Point{-1.0, 0, -1.0},
			0.5,
			&materialLeft),
		hittable.NewSphere(
			vector.Point{-1.0, 0, -1.0},
			-0.4,
			&materialLeft),
		hittable.NewSphere(
			vector.Point{1.0, 0, -1.0},
			0.5,
			&materialRight),
	)

	c := camera.Camera{}
	c.Init(
		// camera.WithSamplesPerPixel(100),
		// camera.WithMaxRayDepth(50),
		camera.WithVFOV(20),
		camera.WithPosition(vector.Vector{0, 1, 0},
			vector.Vector{-2, 2, 1},
			vector.Vector{0, 0, -1},
		),
		camera.WithFocus(10.0, 3.4),
		camera.WithImageWidth(600),
	)
	c.Render("test_ray.ppm", world, 12)
}

func Scene2() {
	materialLeft := hittable.Lambertian{Albedo: vector.Color{0, 0, 1}}
	materialRight := hittable.Lambertian{Albedo: vector.Color{1, 0, 0}}

	r := math.Cos(math.Pi / 4)

	world := hittable.NewWorld(
		hittable.NewSphere(
			vector.Point{-r, 0, -1.0},
			r,
			&materialLeft),
		hittable.NewSphere(
			vector.Point{r, 0, -1.0},
			r,
			&materialRight),
	)

	c := camera.Camera{}
	c.Init(camera.WithImageWidth(500))
	c.Render("test_ray.ppm", world, 16)
}

func Scene3() {
	// World
	materialGround := hittable.Lambertian{Albedo: vector.Color{0.5, 0.5, 0.5}}

	world := hittable.NewWorld(
		hittable.NewSphere(
			vector.Point{0, -1000, 0},
			1000,
			&materialGround),
	)

	for a := -11; a < 11; a++ {
		for b := -11; b < 11; b++ {
			chooseMaterial := rand.Float64()
			center := vector.Point{float64(a) + 0.9*rand.Float64(), 0.2, float64(b) + 0.9*rand.Float64()}
			if center.Add(vector.Point{4, 0.2, 0}.Negative()).Length() > 0.9 {
				var sphereMaterial hittable.Material

				if chooseMaterial < 0.8 {
					//diffuse
					albedo := vector.Multiply(vector.Random(), vector.Random())
					sphereMaterial = &hittable.Lambertian{Albedo: albedo}
				} else if chooseMaterial < 0.95 {
					//metal
					albedo := vector.RandomBounded(0.5, 1)
					fuzz := rand.Float64() / 2
					sphereMaterial = &hittable.Metal{Albedo: albedo, Fuzziness: fuzz}
				} else {
					// glass
					sphereMaterial = &hittable.Dielectric{IR: 1.5}
				}
				lilSphere := hittable.NewSphere(center, 0.2, sphereMaterial)
				lilSphere.MoveTo(center.Add(vector.RandomBounded(0.0, 0.5)))
				world.Append(lilSphere)
			}
		}
	}

	world.Append(
		hittable.NewSphere(
			vector.Point{0, 1, 0},
			1.0,
			&hittable.Dielectric{IR: 1.5},
		),

		hittable.NewSphere(
			vector.Point{-4, 1, 0},
			1.0,
			&hittable.Lambertian{Albedo: vector.Color{0.4, 0.2, 0.1}},
		),

		hittable.NewSphere(
			vector.Point{4, 1, 0},
			1.0,
			&hittable.Metal{Albedo: vector.Color{0.7, 0.6, 0.5}, Fuzziness: 0},
		),
	)

	c := camera.Camera{}
	c.Init(
		camera.WithVFOV(20),
		camera.WithPosition(vector.Vector{0, 1, 0},
			vector.Vector{13, 2, 3},
			vector.Vector{0, 0, 0},
		),
		camera.WithFocus(0.02, 10.0),
		camera.WithImageWidth(1600),
		camera.WithSamplesPerPixel(100),
		camera.WithMaxRayDepth(50),
	)
	c.Render("test_ray.ppm", world.ToBVHTree(), 12)
}

func main() {
	debug.SetGCPercent(1000)
	Scene3()
}
