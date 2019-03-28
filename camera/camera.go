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