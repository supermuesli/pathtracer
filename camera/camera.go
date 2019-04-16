package camera

import (
	"github.com/supermuesli/pathtracer/vec3"
)

type Camera struct {
	Width, Height int
	Origin vec3.Vec3
}

func (c *Camera) Move(x float64, y float64, z float64) {
	c.Origin.X += x
	c.Origin.Y += y
	c.Origin.Z += z
}

func (c *Camera) Rotate_x(x float64) {
	c.Origin.Rotate_x(x)
}

func (c *Camera) Rotate_y(x float64) {
	c.Origin.Rotate_y(x)
}

func (c *Camera) Rotate_z(x float64) {
	c.Origin.Rotate_z(x)
}