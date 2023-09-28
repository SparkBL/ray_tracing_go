package hittable

import (
	"math"
	"ray_tracing/interval"
	"ray_tracing/ray"
	"ray_tracing/vector"
)

type Hittable interface {
	Hit(r ray.Ray, rayT interval.Interval, rec *HitRecord) bool
}

type HitRecord struct {
	Point       vector.Point
	Normal      vector.Vector
	T           float64
	IsFrontFace bool
}

func (hr *HitRecord) SetFaceNormal(r ray.Ray, outwardNormal vector.Vector) {
	// Sets the hit record normal vector.
	// NOTE: the parameter `outward_normal` is assumed to have unit length.

	hr.IsFrontFace = vector.Dot(r.Direction, outwardNormal) < 0
	if hr.IsFrontFace {
		hr.Normal = outwardNormal
	} else {
		hr.Normal = outwardNormal.Negative()
	}

}

type Sphere struct {
	Center vector.Point
	Radius float64
}

func (s *Sphere) Hit(r ray.Ray, rayT interval.Interval, rec *HitRecord) bool {
	ocDistance := r.Origin.Add(s.Center.Negative())
	a := r.Direction.LengthSquared()
	halfB := vector.Dot(ocDistance, r.Direction)
	c := ocDistance.LengthSquared() - s.Radius*s.Radius
	discriminant := halfB*halfB - a*c
	if discriminant < 0 {
		return false
	}
	// Find the nearest root that lies in the acceptable range.
	sqrtd := math.Sqrt(discriminant)
	root := (-halfB - sqrtd) / a

	if !rayT.Surrounds(root) {
		root = (-halfB + sqrtd) / a
		if !rayT.Surrounds(root) {
			return false
		}
	}
	rec.T = root
	rec.Point = r.At(rec.T)
	outwardNormal := rec.Point.Add(s.Center.Negative()).Divide(s.Radius)
	rec.SetFaceNormal(r, outwardNormal)
	return true
}

type Hittables struct {
	objects []Hittable
}

func (hl *Hittables) Append(o ...Hittable) {
	hl.objects = append(hl.objects, o...)
}

func NewWorld(o ...Hittable) *Hittables {
	return &Hittables{
		objects: o,
	}
}

func (hl *Hittables) Hit(r ray.Ray, rayT interval.Interval, rec *HitRecord) bool {
	hitAnything := false
	closestSoFar := rayT.Max
	for _, h := range hl.objects {
		if h.Hit(r, interval.Interval{rayT.Min, closestSoFar}, rec) {
			hitAnything = true
			closestSoFar = rec.T
		}
	}
	return hitAnything
}
