package main

import (
	"ray_tracing/camera"
	"ray_tracing/hittable"
	"ray_tracing/vector"
	"runtime/debug"
)

func main() {

	debug.SetGCPercent(-1)

	//World
	world := hittable.NewWorld(
		&hittable.Sphere{
			Center: vector.Point{0, 0, -1},
			Radius: 0.5},
		&hittable.Sphere{
			Center: vector.Point{0, -100.5, -1},
			Radius: 100},
		&hittable.Sphere{
			Center: vector.Point{-0.5, 0, -1},
			Radius: 0.5},
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
