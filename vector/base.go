package vector

import (
	"fmt"
	"math"
)

type Vector [3]float64

func (v *Vector) X() float64 {
	return v[0]
}

func (v *Vector) Y() float64 {
	return v[1]
}

func (v *Vector) Z() float64 {
	return v[2]
}

func (v Vector) Negative() Vector {
	return [3]float64{-v[0], -v[1], -v[2]}
}

func (v Vector) Add(addend Vector) Vector {
	return Vector{v[0] + addend[0], v[1] + addend[1], v[2] + addend[2]}
}

func (v Vector) Multiply(t float64) Vector {
	return Vector{v[0] * t, v[1] * t, v[2] * t}
}

func (v Vector) Divide(t float64) Vector {
	return Vector{v[0] * (1 / t), v[1] * (1 / t), v[2] * (1 / t)}
}

func (v *Vector) LengthSquared() float64 {
	return v[0]*v[0] + v[1]*v[1] + v[2]*v[2]
}

func (v *Vector) Length() float64 {
	return math.Sqrt(v.LengthSquared())
}

func (v *Vector) IsCloseToZero() bool {
	return v[0] < 1e-8 && v[1] < 1e-8 && v[2] < 1e-8
}

func (v Vector) String() string {
	return fmt.Sprintf("%f %f %f", v[0], v[1], v[2])
}

//Type aliases

type Color = Vector
type Point = Vector
