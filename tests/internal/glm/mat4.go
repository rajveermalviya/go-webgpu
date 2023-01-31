package glm

import "math"

type Mat4[T float] [16]T

func Mat4FromQuaternion[T float](quat Quaternion[T]) Mat4[T] {
	x2 := quat.V[0] + quat.V[0]
	y2 := quat.V[1] + quat.V[1]
	z2 := quat.V[2] + quat.V[2]

	xx2 := x2 * quat.V[0]
	xy2 := x2 * quat.V[1]
	xz2 := x2 * quat.V[2]

	yy2 := y2 * quat.V[1]
	yz2 := y2 * quat.V[2]
	zz2 := z2 * quat.V[2]

	sy2 := y2 * quat.S
	sz2 := z2 * quat.S
	sx2 := x2 * quat.S

	return Mat4[T]{
		1 - yy2 - zz2, xy2 + sz2, xz2 - sy2, 0,
		xy2 - sz2, 1 - xx2 - zz2, yz2 + sx2, 0,
		xz2 + sy2, yz2 - sx2, 1 - xx2 - yy2, 0,
		0, 0, 0, 1,
	}
}

func Mat4FromTranslation[T float](v Vec3[T]) Mat4[T] {
	return Mat4[T]{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		v[0], v[1], v[2], 1,
	}
}

func Mat4FromAngleZ[T float](thetaRad T) Mat4[T] {
	s, c := math.Sincos(float64(thetaRad))

	return Mat4[T]{
		T(c), T(s), 0, 0,
		-T(s), T(c), 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
}

func (lhs Mat4[T]) Mul4(rhs Mat4[T]) Mat4[T] {
	return Mat4[T]{
		lhs[0]*rhs[0] + lhs[4]*rhs[1] + lhs[8]*rhs[2] + lhs[12]*rhs[3],
		lhs[1]*rhs[0] + lhs[5]*rhs[1] + lhs[9]*rhs[2] + lhs[13]*rhs[3],
		lhs[2]*rhs[0] + lhs[6]*rhs[1] + lhs[10]*rhs[2] + lhs[14]*rhs[3],
		lhs[3]*rhs[0] + lhs[7]*rhs[1] + lhs[11]*rhs[2] + lhs[15]*rhs[3],
		lhs[0]*rhs[4] + lhs[4]*rhs[5] + lhs[8]*rhs[6] + lhs[12]*rhs[7],
		lhs[1]*rhs[4] + lhs[5]*rhs[5] + lhs[9]*rhs[6] + lhs[13]*rhs[7],
		lhs[2]*rhs[4] + lhs[6]*rhs[5] + lhs[10]*rhs[6] + lhs[14]*rhs[7],
		lhs[3]*rhs[4] + lhs[7]*rhs[5] + lhs[11]*rhs[6] + lhs[15]*rhs[7],
		lhs[0]*rhs[8] + lhs[4]*rhs[9] + lhs[8]*rhs[10] + lhs[12]*rhs[11],
		lhs[1]*rhs[8] + lhs[5]*rhs[9] + lhs[9]*rhs[10] + lhs[13]*rhs[11],
		lhs[2]*rhs[8] + lhs[6]*rhs[9] + lhs[10]*rhs[10] + lhs[14]*rhs[11],
		lhs[3]*rhs[8] + lhs[7]*rhs[9] + lhs[11]*rhs[10] + lhs[15]*rhs[11],
		lhs[0]*rhs[12] + lhs[4]*rhs[13] + lhs[8]*rhs[14] + lhs[12]*rhs[15],
		lhs[1]*rhs[12] + lhs[5]*rhs[13] + lhs[9]*rhs[14] + lhs[13]*rhs[15],
		lhs[2]*rhs[12] + lhs[6]*rhs[13] + lhs[10]*rhs[14] + lhs[14]*rhs[15],
		lhs[3]*rhs[12] + lhs[7]*rhs[13] + lhs[11]*rhs[14] + lhs[15]*rhs[15],
	}
}
