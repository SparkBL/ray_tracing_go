package main

import (
	"math"
	"ray_tracing/ray"
	"ray_tracing/vector"
)

type Hittable interface {
	Hit(r ray.Ray, tMin, tMax float64, rec *HitRecord) bool
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

	if vector.Dot(r.Direction, outwardNormal) < 0 {
		hr.Normal = outwardNormal
	} else {
		hr.Normal = outwardNormal.Negative()
	}

}

type Sphere struct {
	Center vector.Point
	Radius float64
}

func (s *Sphere) Hit(r ray.Ray, tMin, tMax float64, rec *HitRecord) bool {
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

	if root <= tMin || tMax <= root {
		root = (-halfB + sqrtd) / a
		if root <= tMin || tMax <= root {
			return false
		}
	}
	rec.T = root
	rec.Point = r.At(rec.T)
	outwardNormal := rec.Point.Add(s.Center.Negative()).Divide(s.Radius)
	rec.SetFaceNormal(r, outwardNormal)
	return true
}

type HittableList []Hittable

func (hl HittableList) Hit(r ray.Ray, tMin, tMax float64, rec *HitRecord) bool {
	hitAnything := false
	closestSoFar := tMax
	for _, h := range hl {
		if h.Hit(r, tMin, closestSoFar, rec) {
			hitAnything = true
			closestSoFar = rec.T
		}
	}
	return hitAnything
}
