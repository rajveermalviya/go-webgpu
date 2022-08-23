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

func Perspective[T float](fovyDeg, aspect, near, far T) Mat4[T] {
	fovyRad := degToRad(fovyDeg)
	f := T(1 / math.Tan(float64(fovyRad*0.5)))

	return Mat4[T]{
		f / aspect, 0, 0, 0,
		0, f, 0, 0,
		0, 0, (far + near) / (near - far), -1,
		0, 0, (2 * far * near) / (near - far), 0,
	}
}

func LookAtRH[T float](eye, center, up Vec3[T]) Mat4[T] {
	f := (center.Sub(eye)).Normalize()
	s := f.Cross(up).Normalize()
	u := s.Cross(f)

	return Mat4[T]{
		s[0], u[0], -f[0], 0,
		s[1], u[1], -f[1], 0,
		s[2], u[2], -f[2], 0,
		-eye.Dot(s), -eye.Dot(u), eye.Dot(f), 1,
	}
}

func degToRad[T float](deg T) (rad T) {
	return deg * (math.Pi / 180)
}
