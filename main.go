package main

import (
	"ray_tracing/camera"
	"ray_tracing/hittable"
	"ray_tracing/vector"
)

func main() {

	// debug.SetGCPercent(-1)

	//World
	materialGround := hittable.Lambertian{Albedo: vector.Color{0.8, 0.8, 0.0}}
	materialCenter := hittable.Lambertian{Albedo: vector.Color{0.1, 0.2, 0.5}}
	// materialLeft := hittable.Metal{Albedo: vector.Color{1, 1, 1}, Fuzziness: 0.2}
	materialLeft := hittable.Dielectric{1.5}
	materialRight := hittable.Metal{Albedo: vector.Color{0.8, 0.6, 0.2}, Fuzziness: 0.0}

	world := hittable.NewWorld(
		&hittable.Sphere{
			Center:   vector.Point{0, -100.5, -1.0},
			Material: &materialGround,
			Radius:   100},
		&hittable.Sphere{
			Center:   vector.Point{0, 0, -1.0},
			Material: &materialCenter,
			Radius:   0.5},
		&hittable.Sphere{
			Center:   vector.Point{-1.0, 0, -1.0},
			Material: &materialLeft,
			Radius:   0.5},
		&hittable.Sphere{
			Center:   vector.Point{1.0, 0, -1.0},
			Material: &materialRight,
			Radius:   0.5},
		// &hittable.Plane{
		// 	Center:   vector.Point{0, 0, -1.5},
		// 	Material: &materialCenter,
		// 	Normal:   vector.Vector{0, 0.5, 0.3}},
	)

	camera := camera.Camera{}
	camera.Init()
	camera.Render("test_ray.ppm", world)
}

// func OutputImage(width, height int) {
// 	out, _ := os.Create("test.ppm")
// 	logger := bufio.NewWriter(os.Stdout)
// 	out.WriteString(fmt.Sprintf("P3\n%d %d\n255\n", width, height))
// 	for i := 0; i < height; i++ {
// 		logger.WriteString(fmt.Sprintf("\rScanlines remaining: %d ", height-i))
// 		logger.Flush()
// 		for j := 0; j < width; j++ {
// 			c := vector.Vector{
// 				float64(j) / (float64(width) - 1),
// 				float64(i) / (float64(height) - 1),
// 				float64(0)}

// 			out.WriteString(vector.ColorString(&c))
// 		}
// 	}
// 	out.Close()
// 	logger.Flush()
// }
