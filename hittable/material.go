package hittable

import (
	"ray_tracing/ray"
	"ray_tracing/vector"
)

type Material interface {
	Scatter(rIn ray.Ray, rec *HitRecord) (bool, ray.Ray, vector.Color)
}

type Lambertian struct {
	Albedo vector.Color
}

//	func (l *Lambertian) Scatter(rIn, rScattered *ray.Ray, rec *HitRecord, attenuation *vector.Color) bool {
//		scatterDirection := rec.Normal.Add(vector.RandomUnitVector())
//		if scatterDirection.IsCloseToZero() {
//			scatterDirection = rec.Normal
//		}
//		rScattered = &ray.Ray{Origin: rec.Point, Direction: scatterDirection}
//		attenuation = &l.Albedo
//		return true
//	}

func (l *Lambertian) Scatter(rIn ray.Ray, rec *HitRecord) (bool, ray.Ray, vector.Color) {
	scatterDirection := rec.Normal.Add(vector.RandomUnitVector())
	if scatterDirection.IsCloseToZero() {
		scatterDirection = rec.Normal
	}
	return true, ray.Ray{Origin: rec.Point, Direction: scatterDirection}, l.Albedo

}

type Metal struct {
	Albedo    vector.Color
	Fuzziness float64 //0 <= x < 1
}

// func (l *Metal) Scatter(rIn, rScattered *ray.Ray, rec *HitRecord, attenuation *vector.Color) bool {
// 	reflected := vector.Reflect(vector.UnitVector(rIn.Direction), rec.Normal)
// 	rScattered = &ray.Ray{Origin: rec.Point, Direction: reflected}
// 	attenuation = &l.Albedo
// 	return true

// }

func (l *Metal) Scatter(rIn ray.Ray, rec *HitRecord) (bool, ray.Ray, vector.Color) {
	reflected := vector.Reflect(vector.UnitVector(rIn.Direction), rec.Normal)
	scattered := ray.Ray{Origin: rec.Point, Direction: reflected.Add(vector.RandomUnitVector().Multiply(l.Fuzziness))}
	return vector.Dot(scattered.Direction, rec.Normal) > 0.0, scattered, l.Albedo

}
