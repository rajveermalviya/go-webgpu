package glm

import "math"

type Vec3[T float] [3]T

func (v1 Vec3[T]) Dot(v2 Vec3[T]) T {
	return (v1[0] * v2[0]) + (v1[1] * v2[1]) + (v1[2] * v2[2])
}

func (v1 Vec3[T]) Length() T {
	return T(math.Sqrt(float64(v1.Dot(v1))))
}

func (v1 Vec3[T]) LengthRecip() T {
	return (1 / v1.Length())
}

func (v1 Vec3[T]) MulScalar(s T) Vec3[T] {
	return Vec3[T]{
		v1[0] * s,
		v1[1] * s,
		v1[2] * s,
	}
}

func (v1 Vec3[T]) Normalize() Vec3[T] {
	return v1.MulScalar(v1.LengthRecip())
}

func (v1 Vec3[T]) Cross(v2 Vec3[T]) Vec3[T] {
	return Vec3[T]{
		v1[1]*v2[2] - v2[1]*v1[2],
		v1[2]*v2[0] - v2[2]*v1[0],
		v1[0]*v2[1] - v2[0]*v1[1],
	}
}

func (v1 Vec3[T]) Sub(v2 Vec3[T]) Vec3[T] {
	return Vec3[T]{
		v1[0] - v2[0],
		v1[1] - v2[1],
		v1[2] - v2[2],
	}
}
