package object

import (
	"github.com/supermuesli/pathtracer/vec3"
	"math"
)

type Line struct {
	Origin, Dir vec3.Vec3
}

type Triangle struct {
	A, B, C vec3.Vec3
	Pdf func(vec3.Vec3, vec3.Vec3) vec3.Vec3 
	Mterial Material
}

type Sphere struct {
	Origin vec3.Vec3
	Radius float64
	Pdf func(vec3.Vec3, vec3.Vec3) vec3.Vec3 
	Mterial Material
}

type Object struct {
	Mesh []Triangle
}

type Material struct {
	Diffuse_color vec3.Vec3
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

// rotate object in 3d space
func (o *Object) Rotate_x(x float64) {
	for i := 0; i < len(o.Mesh); i++ {
		o.Mesh[i].A.Rotate_x(x)
		o.Mesh[i].C.Rotate_x(x)
		o.Mesh[i].B.Rotate_x(x)
	}
}

// rotate object in 3d space
func (o *Object) Rotate_y(x float64) {
	for i := 0; i < len(o.Mesh); i++ {
		o.Mesh[i].A.Rotate_y(x)
		o.Mesh[i].C.Rotate_y(x)
		o.Mesh[i].B.Rotate_y(x)
	}
}

// rotate object in 3d space
func (o *Object) Rotate_z(x float64) {
	for i := 0; i < len(o.Mesh); i++ {
		o.Mesh[i].A.Rotate_z(x)
		o.Mesh[i].C.Rotate_z(x)
		o.Mesh[i].B.Rotate_z(x)
	}
}

func min (a float64, b float64) float64 {
	if a < b {
		return a
	}

	return b
}

func (s Sphere) Intersection(ray *Line) (bool, float64) {
	ro_so := ray.Origin
	ro_so.Sub(s.Origin)
	rd_so_dot := ray.Dir.Dot(ro_so)
	t := math.Pow(rd_so_dot, 2) - math.Pow(ro_so.Euclidean_norm(), 2) + math.Pow(s.Radius, 2)
	d1 := -rd_so_dot + math.Sqrt(t)
	d2 := -rd_so_dot - math.Sqrt(t)

	if d1 < 0 && d2 < 0 {
		if math.IsNaN(d1) && math.IsNaN(d2) {
			return false, math.Inf(1)
		}
	}

	if math.IsNaN(d1) {
		d1 = math.Inf(1)
	}

	if math.IsNaN(d2) {
		d2 = math.Inf(1)
	}

	_min := min(d1, d2)
	if _min > 0 {
		return true, _min
	}

	return false, math.Inf(1)
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
