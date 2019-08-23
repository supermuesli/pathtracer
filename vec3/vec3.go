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

func (a *Vec3) Clamp() {
	if a.X > 255.0 {
		a.X = 255.0
	}
	if a.Y > 255.0 {
		a.Y = 255.0
	}
	if a.Z > 255.0 {
		a.Z = 255.0
	}
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

func (a *Vec3) Rotate_x(theta float64) {
	temp := a
	a.Y = math.Cos(theta)*temp.Y - math.Sin(theta)*temp.Z
	a.Z = math.Sin(theta)*temp.Y + math.Cos(theta)*temp.Z
}

func (a *Vec3) Rotate_y(theta float64) {
	temp := a
	a.X = math.Cos(theta)*temp.X + math.Sin(theta)*temp.Z
	a.Z = -math.Sin(theta)*temp.X + math.Cos(theta)*temp.Z
}

func (a *Vec3) Rotate_z(theta float64) {
	temp := a
	a.X = math.Cos(theta)*temp.X - math.Sin(theta)*temp.Y
	a.Y = math.Sin(theta)*temp.X + math.Cos(theta)*temp.Y
}

func (a *Vec3) Rotate_around_normal(theta float64, normal Vec3) {
	temp := a
	a.X = (normal.X*normal.X * (1-math.Cos(theta)) + math.Cos(theta)) * temp.X + (normal.X*normal.Y * (1-math.Cos(theta)) - normal.Z*math.Sin(theta)) * temp.Y + (normal.X*normal.Y * (1-math.Cos(theta)) + normal.Y*math.Sin(theta)) * temp.Z
	a.Y = (normal.Y*normal.X * (1-math.Cos(theta)) + normal.Z*math.Sin(theta)) * temp.X + (normal.Y*normal.Y * (1-math.Cos(theta)) + math.Cos(theta)) * temp.Y + (normal.Y*normal.Z * (1-math.Cos(theta)) - normal.X*math.Sin(theta)) * temp.Z
	a.Z = (normal.Z*normal.X * (1-math.Cos(theta)) - normal.Y*math.Sin(theta)) * temp.X + (normal.Z*normal.Y * (1-math.Cos(theta)) + normal.X*math.Sin(theta)) * temp.Y + (normal.Z*normal.Z * (1-math.Cos(theta)) + math.Cos(theta)) * temp.Z
}