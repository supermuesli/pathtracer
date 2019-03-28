package vec3

import (
	"math"
)

type Vec3 struct {
	X, Y, Z float64
}

func (a *Vec3) Add(b Vec3) {
	a.X += b.X
	a.Y += b.Y
	a.Z += b.Z
}

func (a *Vec3) Sub(b Vec3) {
	a.X -= b.X
	a.Y -= b.Y
	a.Z -= b.Z
}

func (a *Vec3) Scale(b float64) {
	a.X *= b
	a.Y *= b
	a.Z *= b
}

func (a *Vec3) Component_wise_mul(b Vec3) {
	a.X *= b.X
	a.Y *= b.Y
	a.Z *= b.Z
}

func (a Vec3) Dot(b Vec3) float64 {
	return a.X * b.X + a.Y * b.Y + a.Z * b.Z
}

func (a *Vec3) Cross(b Vec3) {
	temp_x := a.X
	a.X = a.Y * b.Z - a.Z * b.Y
	temp_y := a.Y 
	a.Y = a.Z * b.X - temp_x * b.Z
	a.Z = temp_x * b.Y - temp_y * b.X
}
	
func (a Vec3) Euclidean_norm() float64 {
	return math.Sqrt(a.Dot(a))
}

func (a *Vec3) Normalize() {
	a.Scale(1/a.Euclidean_norm())
}