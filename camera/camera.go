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
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/seehuhn/mt19937"

	"ray_tracing/hittable"
)

var randGenerator *rand.Rand = rand.New(mt19937.New())
var mu sync.Mutex

type pair struct {
	f int
	s string
}

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

	//	focalLength    float64
	viewportHeight float64
	viewportWidth  float64

	verticalFieldOfView float64
	lookFrom            vector.Point
	lookAt              vector.Point
	vUp                 vector.Vector
	u, v, w             vector.Vector

	defocusAngle  float64       // Variation angle of rays through each pixel
	focusDistance float64       // Distance from camera lookfrom point to plane of perfect focus
	defocusDiskU  vector.Vector // Defocus disk horizontal radius
	defocusDiskV  vector.Vector // Defocus disk vertical radius

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

	c.defocusAngle = 0
	c.focusDistance = 10
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

func WithFocus(defocusAngle, focusDistance float64) CameraOption {
	return func(c *Camera) *Camera {
		c.defocusAngle = defocusAngle
		c.focusDistance = focusDistance
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

	// Determine viewport dimensions.
	c.center = c.lookFrom
	//We assume focal_length == focus_distance
	theta := util.DegressToRadians(c.verticalFieldOfView)
	h := math.Tan(theta / 2)
	c.viewportHeight = 2.0 * h * c.focusDistance //think as matrix
	c.viewportWidth = c.viewportHeight * (float64(c.imageWidth) / float64(c.imageHeight))

	// Calculate the u,v,w unit basis vectors for the camera coordinate frame.
	c.w = vector.UnitVector(c.lookFrom.Add(c.lookAt.Negative()))
	c.u = vector.UnitVector(vector.Cross(c.vUp, c.w))
	c.v = vector.Cross(c.w, c.u)

	// Calculate the horizontal and vertical delta vectors to the next pixel.
	viewportU := c.u.Multiply(c.viewportWidth)
	viewportV := c.v.Negative().Multiply(c.viewportHeight)
	c.pixelDeltaU = viewportU.Divide(float64(c.imageWidth))
	c.pixelDeltaV = viewportV.Divide(float64(c.imageHeight))

	// Calculate the location of the upper left pixel.
	viewPortUpleft :=
		c.center.Add(c.w.Multiply(c.focusDistance).Negative()).
			Add(viewportU.Divide(2).Negative()).
			Add(viewportV.Divide(2).Negative())
	c.pixelZeroLocation = viewPortUpleft.Add(c.pixelDeltaU.Add(c.pixelDeltaV).Multiply(0.5))

	// Calculate the camera defocus disk basis vectors.
	defocusRadius := c.focusDistance * math.Tan(util.DegressToRadians(c.defocusAngle/2))
	c.defocusDiskU = c.u.Multiply(defocusRadius)
	c.defocusDiskV = c.v.Multiply(defocusRadius)

	c.logger = bufio.NewWriter(os.Stdout)

	fmt.Printf("%#v\n", c)
}

func (c *Camera) Render(filename string, world hittable.Hittable) {
	start := time.Now()
	err := os.Remove(filename)
	if err != nil {
		log.Println(err)
	}
	var output chan pair = make(chan pair, 500)
	var quit chan bool = make(chan bool)
	//var mux sync.Mutex
	wg := sync.WaitGroup{}
	out, _ := os.Create(filename)
	var buf strings.Builder = strings.Builder{}
	keys := make([]int, 0, c.imageHeight*c.imageWidth)
	pixels := make(map[int]string, c.imageHeight*c.imageWidth)
	go func(chan pair, chan bool) {
		out.WriteString(fmt.Sprintf("P3\n%d %d\n255\n", c.imageWidth, c.imageHeight))
		for {
			select {
			case s := <-output:
				keys = append(keys, s.f)
				pixels[s.f] = s.s
				wg.Done()
			case <-quit:
				break
			}
		}
	}(output, quit)
	count := 0
	for i := 0; i < c.imageHeight; i++ {
		c.logger.WriteString(fmt.Sprintf("\rremaining: %.2f%%", 100.0*float64(c.imageHeight-i)/float64(c.imageHeight)))
		c.logger.Flush()
		for j := 0; j < c.imageWidth; j++ {
			wg.Add(1)
			go func(k, w, cnt int) {
				pixelColor := vector.Color{0, 0, 0}
				for sample := 0; sample < c.samplesPerPixel; sample++ {
					r := c.getRay(w, k)
					pixelColor = pixelColor.Add(c.rayColor(r, c.maxRayDepth, world)) //performance boost if pointer
				}
				output <- pair{cnt, ColorString(&pixelColor, c.samplesPerPixel)}
			}(i, j, count)
			count++
		}
	}
	wg.Wait()
	quit <- true

	sort.Ints(keys)
	for _, k := range keys {
		buf.WriteString(pixels[k])
	}
	out.WriteString(buf.String())
	buf.Reset()

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
	// Get a randomly-sampled camera ray for the pixel at location i,j, originating from
	// the camera defocus disk.
	pixelCenter := c.pixelZeroLocation.
		Add(c.pixelDeltaU.Multiply(float64(i))).
		Add(c.pixelDeltaV.Multiply(float64(j)))

	pixelSample := pixelCenter.Add(c.pixelSampleSquare())
	return ray.Ray{
		Origin:    c.defocusDiskSample(),
		Direction: pixelSample.Add(c.center.Negative()),
	}
}

func (c *Camera) pixelSampleSquare() vector.Vector {
	mu.Lock()
	px, py := -0.5+randGenerator.Float64(), -0.5+randGenerator.Float64()
	mu.Unlock()
	return c.pixelDeltaU.Multiply(px).Add(c.pixelDeltaV.Multiply(py))
}

func (c *Camera) defocusDiskSample() vector.Point {
	p := vector.RandomInUnitDisk()
	if c.defocusAngle <= 0 {
		return c.center
	}
	return c.center.Add(c.defocusDiskU.Multiply(p[0])).Add(c.defocusDiskV.Multiply(p[1]))
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
