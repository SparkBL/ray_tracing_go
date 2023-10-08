package camera

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"ray_tracing/interval"
	"ray_tracing/ray"
	"ray_tracing/util"
	"ray_tracing/vector"
	"strings"
	"time"

	"github.com/seehuhn/mt19937"

	"ray_tracing/hittable"
)

var randGenerator *rand.Rand = rand.New(mt19937.New())

type Camera struct {
	aspectRatio       float64
	imageWidth        int
	imageHeight       int
	center            vector.Point
	pixelZeroLocation vector.Vector
	pixelDeltaU       vector.Vector
	pixelDeltaV       vector.Vector

	samplesPerPixel int
	maxRayDepth     int

	focalLength    float64
	viewportHeight float64
	viewportWidth  float64

	verticalFieldOfView float64
	lookFrom            vector.Point
	lookAt              vector.Point
	vUp                 vector.Vector
	u, v, w             vector.Vector

	logger *bufio.Writer
}

type CameraOption func(c *Camera) *Camera

func DefaultOption(c *Camera) *Camera {
	c.aspectRatio = 16.0 / 9.0
	c.imageWidth = 1000

	c.lookFrom = vector.Vector{0, 0, 0}
	c.lookAt = vector.Vector{0, 0, -1}
	c.vUp = vector.Vector{0, 1, 0}

	c.verticalFieldOfView = 90 //degrees

	c.samplesPerPixel = 10
	c.maxRayDepth = 50
	return c
}

func WithPosition(vUp, lookFrom, lookAt vector.Vector) CameraOption {
	return func(c *Camera) *Camera {
		c.lookFrom = lookFrom
		c.lookAt = lookAt
		c.vUp = vUp
		return c
	}
}

func WithAspectRatio(aspectRatio float64) CameraOption {
	return func(c *Camera) *Camera {
		c.aspectRatio = aspectRatio
		return c
	}
}

func WithImageWidth(imageWidth int) CameraOption {
	return func(c *Camera) *Camera {
		c.imageWidth = imageWidth
		return c
	}
}

func WithSamplesPerPixel(samplesPerPixel int) CameraOption {
	return func(c *Camera) *Camera {
		c.samplesPerPixel = samplesPerPixel
		return c
	}
}

func WithMaxRayDepth(maxRayDepth int) CameraOption {
	return func(c *Camera) *Camera {
		c.maxRayDepth = maxRayDepth
		return c
	}
}

func WithVFOV(vFov float64) CameraOption {
	return func(c *Camera) *Camera {
		c.verticalFieldOfView = vFov
		return c
	}
}

func (c *Camera) Init(opts ...CameraOption) {
	c = DefaultOption(c)
	for _, o := range opts {
		c = o(c)
	}

	c.imageHeight = int(float64(c.imageWidth) / c.aspectRatio)
	if c.imageHeight < 1 {
		c.imageHeight = 1
	}
	c.center = c.lookFrom
	c.focalLength = c.lookFrom.Add(c.lookAt.Negative()).Length()
	theta := util.DegressToRadians(c.verticalFieldOfView)
	h := math.Tan(theta / 2)
	c.viewportHeight = 2.0 * h * c.focalLength //think as matrix
	c.viewportWidth = c.viewportHeight * (float64(c.imageWidth) / float64(c.imageHeight))

	// Calculate the u,v,w unit basis vectors for the camera coordinate frame.
	c.w = vector.UnitVector(c.lookFrom.Add(c.lookAt.Negative()))
	c.u = vector.UnitVector(vector.Cross(c.vUp, c.w))
	c.v = vector.Cross(c.w, c.u)

	viewportU := c.u.Multiply(c.viewportWidth)
	viewportV := c.v.Negative().Multiply(c.viewportHeight)
	c.pixelDeltaU = viewportU.Divide(float64(c.imageWidth))
	c.pixelDeltaV = viewportV.Divide(float64(c.imageHeight))

	viewPortUpleft :=
		c.center.Add(c.w.Multiply(c.focalLength).Negative()).
			Add(viewportU.Divide(2).Negative()).
			Add(viewportV.Divide(2).Negative())
	c.pixelZeroLocation = viewPortUpleft.Add(c.pixelDeltaU.Add(c.pixelDeltaV).Multiply(0.5))
	c.logger = bufio.NewWriter(os.Stdout)

	fmt.Printf("%#v\n", c)
}

func (c *Camera) Render(filename string, world hittable.Hittable) {
	start := time.Now()
	err := os.Remove(filename)
	if err != nil {
		log.Println(err)
	}
	var output chan string = make(chan string, 3)
	var quit chan bool = make(chan bool)
	go func(chan string, chan bool) {
		out, _ := os.Create(filename)
		var buf strings.Builder = strings.Builder{}
		out.WriteString(fmt.Sprintf("P3\n%d %d\n255\n", c.imageWidth, c.imageHeight))
		for {
			select {
			case s := <-output:
				buf.WriteString(s)
				if buf.Len() > 1000 {
					out.WriteString(buf.String())
					buf.Reset()
				}

			case <-quit:
				out.Close()
			}
		}
	}(output, quit)

	for i := 0; i < c.imageHeight; i++ {

		c.logger.WriteString(fmt.Sprintf("\rremaining: %.2f%%", 100.0*float64(c.imageHeight-i)/float64(c.imageHeight)))
		c.logger.Flush()

		for j := 0; j < c.imageWidth; j++ {
			pixelColor := vector.Color{0, 0, 0}
			for sample := 0; sample < c.samplesPerPixel; sample++ {
				r := c.getRay(j, i)
				pixelColor = pixelColor.Add(c.rayColor(r, c.maxRayDepth, world)) //performance boost if pointer
			}
			output <- ColorString(&pixelColor, c.samplesPerPixel)
		}
	}
	quit <- true

	c.logger.WriteString(fmt.Sprintf("\relapsed: %v\n", time.Since(start)))
	c.logger.Flush()
}

func (c *Camera) rayColor(r ray.Ray, depth int, world hittable.Hittable) vector.Color {
	rec := hittable.HitRecord{}
	// If we've exceeded the ray bounce limit, no more light is gathered.
	if depth <= 0 {
		return vector.Color{0, 0, 0}
	}

	if world.Hit(r, interval.Interval{0.001, math.Inf(1)}, &rec) {
		if ok, scattered, attenuation := rec.Material.Scatter(r, &rec); ok {
			return vector.Multiply(attenuation, c.rayColor(scattered, depth-1, world))
		}
		return vector.Color{0, 0, 0}
		//	direction := vector.RandomOnHemisphere(rec.Normal) //Random
		//direction := rec.Normal.Add(vector.RandomUnitVector()) //Lambertian
		//return c.rayColor(ray.Ray{Origin: rec.Point, Direction: direction}, depth-1, world).Multiply(0.5)

	}
	unitDirection := vector.UnitVector(r.Direction)
	a := 0.5 * (unitDirection.Y() + 1.0)
	return vector.Color{1.0, 1.0, 1.0}.
		Multiply(1.0 - a).
		Add(vector.Color{0.5, 0.7, 1.0}.Multiply(a))
}

func (c *Camera) getRay(i, j int) ray.Ray {
	pixelCenter := c.pixelZeroLocation.
		Add(c.pixelDeltaU.Multiply(float64(i))).
		Add(c.pixelDeltaV.Multiply(float64(j)))

	pixelSample := pixelCenter.Add(c.pixelSampleSquare())
	return ray.Ray{
		Origin:    c.center,
		Direction: pixelSample.Add(c.center.Negative()),
	}
}

func (c *Camera) pixelSampleSquare() vector.Vector {
	px, py := -0.5+randGenerator.Float64(), -0.5+randGenerator.Float64()
	return c.pixelDeltaU.Multiply(px).Add(c.pixelDeltaV.Multiply(py))
}

func ColorString(c *vector.Color, samples int) string {

	scale := 1.0 / float64(samples) //get avg of samples
	r := util.LinearToGamma(c.X() * scale)
	g := util.LinearToGamma(c.Y() * scale)
	b := util.LinearToGamma(c.Z() * scale)
	intensity := interval.Interval{0.000, 0.999}

	return fmt.Sprintf(
		"%d %d %d\n",
		int(intensity.Clamp(r)*float64(256)),
		int(intensity.Clamp(g)*float64(256)),
		int(intensity.Clamp(b)*float64(256)),
	)
}
