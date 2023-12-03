package ray

import (
	"ray_tracing/vector"
	"sync"
)

var pool = sync.Pool{
	New: func() interface{} { return &Ray{} },
}

func putRay(r *Ray) {

}

type Ray struct {
	Origin    vector.Point
	Direction vector.Vector
	Time      float64
}

func Get() *Ray {

	// fmt.Print(pool.)
	return pool.Get().(*Ray)
}

func Put(r *Ray) {
	r.Origin = vector.Vector{0, 0, 0}
	r.Direction = vector.Vector{0, 0, 0}
	pool.Put(r)
}

func (r *Ray) At(t float64) vector.Point {
	return r.Origin.Add(r.Direction.Multiply(t))
}
