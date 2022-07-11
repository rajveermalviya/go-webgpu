package glm

import "math"

func PerspectiveRH[T float](fovYradians, aspectRatio, zNear, zFar T) Mat4[T] {
	sinFov, cosFov := math.Sincos(float64(0.5) * float64(fovYradians))
	h := T(cosFov) / T(sinFov)
	w := h / aspectRatio
	r := zFar / (zNear - zFar)

	return Mat4[T]{
		w, 0, 0, 0,
		0, h, 0, 0,
		0, 0, r, -1,
		0, 0, r * zNear, 0,
	}
}

func LookToLH[T float](eye, dir, up Vec3[T]) Mat4[T] {
	f := dir.Normalize()
	s := up.Cross(f).Normalize()
	u := f.Cross(s)
	return Mat4[T]{
		s[0], u[0], f[0], 0,
		s[1], u[1], f[1], 0,
		s[2], u[2], f[2], 0,
		-s.Dot(eye), -u.Dot(eye), -f.Dot(eye), 1,
	}
}

func LookAtRH[T float](eye, center, up Vec3[T]) Mat4[T] {
	return LookToLH(eye, eye.Sub(center), up)
}
