package glm

import "math"

type Quaternion[T float] struct {
	V Vec3[T]
	S T
}

func QuaternionFromAxisAngle[T float](axis Vec3[T], angleRad T) Quaternion[T] {
	sin, cos := math.Sincos(float64(angleRad) * 0.5)
	return Quaternion[T]{
		S: T(cos),
		V: axis.MulScalar(T(sin)),
	}
}

func (lhs Quaternion[T]) Mul(rhs Quaternion[T]) Quaternion[T] {
	return Quaternion[T]{
		S: lhs.S*rhs.S - lhs.V[0]*rhs.V[0] - lhs.V[1]*rhs.V[1] - lhs.V[2]*rhs.V[2],
		V: Vec3[T]{
			lhs.S*rhs.V[0] + lhs.V[0]*rhs.S + lhs.V[1]*rhs.V[2] - lhs.V[2]*rhs.V[1],
			lhs.S*rhs.V[1] + lhs.V[1]*rhs.S + lhs.V[2]*rhs.V[0] - lhs.V[0]*rhs.V[2],
			lhs.S*rhs.V[2] + lhs.V[2]*rhs.S + lhs.V[0]*rhs.V[1] - lhs.V[1]*rhs.V[0],
		},
	}
}
