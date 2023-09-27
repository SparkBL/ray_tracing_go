package vector

func Add(v1, v2 Vector) Vector {
	return Vector{v1[0] + v2[0], v1[1] + v2[1], v1[2] + v2[2]}
}

func Substract(v1, v2 Vector) Vector {
	return Vector{v1[0] - v2[0], v1[1] - v2[1], v1[2] - v2[2]}
}

func Multiply(v1, v2 Vector) Vector {
	return Vector{v1[0] * v2[0], v1[1] * v2[1], v1[2] * v2[2]}
}

func MultiplyByValue(v Vector, t float64) Vector {
	return Vector{v[0] * t, v[1] * t, v[2] * t}
}

func DivideByValue(v Vector, t float64) Vector {
	return Vector{v[0] * (1 / t), v[1] * (1 / t), v[2] * (1 / t)}
}

func Dot(v1, v2 Vector) float64 {
	return v1[0]*v2[0] +
		v1[1]*v2[1] +
		v1[2]*v2[2]
}

func Cross(v1, v2 Vector) Vector {
	return Vector{
		v1[1]*v2[2] - v1[2]*v2[1],
		v1[2]*v2[0] - v1[0]*v2[2],
		v1[0]*v2[1] - v1[1]*v2[0],
	}
}

func UnitVector(v Vector) Vector {
	return v.Divide(v.Length())

}
