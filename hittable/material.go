package hittable

import (
	"math"
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

type Dielectric struct {
	IR float64 //Refraction Index
}

func (d *Dielectric) Scatter(rIn ray.Ray, rec *HitRecord) (bool, ray.Ray, vector.Color) {

	attenuation := vector.Color{1, 1, 1}
	refractionRatio := d.IR
	if rec.IsFrontFace {
		refractionRatio = 1.0 / d.IR
	}
	unitDirection := vector.UnitVector(rIn.Direction)
	cosTheta := math.Min(vector.Dot(unitDirection.Negative(), rec.Normal), 1.0)
	sinTheta := math.Sqrt(1.0 - cosTheta*cosTheta)
	direction := vector.Vector{}

	cannotRefract := refractionRatio*sinTheta > 1.0
	mu.Lock()
	randFloat := randGenerator.Float64()
	mu.Unlock()
	if cannotRefract || d.reflectance(cosTheta, refractionRatio) > randFloat {
		direction = vector.Reflect(unitDirection, rec.Normal)
	} else {
		direction = vector.Refract(unitDirection, rec.Normal, refractionRatio)
	}
	scattered := ray.Ray{rec.Point, direction}
	return true, scattered, attenuation

}

func (d *Dielectric) reflectance(cosine, refIdx float64) float64 {
	//Use Shlick's approx for reflectance
	r0 := (1 - refIdx) / (1 + refIdx)
	r0 = r0 * r0
	return r0 + (1-r0)*math.Pow((1-cosine), 5)
}
