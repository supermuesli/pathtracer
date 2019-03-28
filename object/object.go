package object

import (
	"github.com/supermuesli/pathtracer/vec3"
)

type Line struct {
	Origin, Dir vec3.Vec3
}

type Triangle struct {
	A, B, C vec3.Vec3
}

type Sphere struct {
	Origin vec3.Vec3
	Radius float64
}

type Object struct {
	Mesh []Triangle
	Mterial Material
}

type Material struct {
	Ambient_color vec3.Vec3
	Diffuse_color vec3.Vec3
	Specular_color vec3.Vec3
	Emission float64
}

// move object in 3d space
func (o *Object) Move(x float64, y float64, z float64) {
	adder := vec3.Vec3{x, y, z}
	for i := 0; i < len(o.Mesh); i++ {
		o.Mesh[i].A.Add(adder)
		o.Mesh[i].C.Add(adder)
		o.Mesh[i].B.Add(adder)
	}
}

func (t Triangle) Intersection(ray *Line) (bool, float64) {
	const epsilon = 0.0000001 // minimum offset distance (otherwise rays will always intersect the hit_positions they're on)

	ta := t.A
	edge1 := t.B
	edge2 := t.C
	raydir := ray.Dir
	s := ray.Origin

	edge1.Sub(ta)
	edge2.Sub(ta)
	
	h := raydir
	h.Cross(edge2)

	a := edge1.Dot(h)
	
	if a > -epsilon && a < epsilon {
		return false, 0
	}

	f := 1.0/a
	
	s.Sub(ta)

	u := f * (s.Dot(h))
	
	if u < 0.0 || u > 1.0 {
		return false, 0
	}

	s.Cross(edge1)
	
	q := s 
	v := f * raydir.Dot(q)

	if v < 0.0 || u + v > 1.0 {
		return false, 0
	}
	
	// At this stage we can compute d to find out where the intersection point is on the line (hit_vector = ray.origin + d*ray.direction).
	d := f * edge2.Dot(q)

	// ray intersection
	if d > epsilon { 
		return true, d
	}

	// This means that there is a line intersection but not a ray intersection.
	return false, 0
}
