package main

import (
	"github.com/supermuesli/pathtracer/vec3"
	"github.com/supermuesli/pathtracer/object"
	"github.com/supermuesli/pathtracer/camera"
	// "github.com/pkg/profile"
	"github.com/veandco/go-sdl2/sdl"
	"math"
    "sync"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"fmt"
	"math/rand"
	"strconv"
	_ "time"
)

const (
	window_width = 500
	window_height = 500
)

var floats []float64
var float_offset int = -1
var float_amount int = 1000000
var inf float64 = math.Inf(1)
var objects []object.Object
var spheres []object.Sphere

// returns the next random float in sequence
func rand_float() float64 {
	float_offset = (float_offset + 1) % float_amount 
	return floats[float_offset]
}

// returns next random float (possibly negative) in sequence
func rand_neg_float() float64 {
	res := rand_float()
	if rand_float() < 0.5 {
		return -res
	}

	return res
}

// takes a ray and checks for intersections among all objects in world space
// returns color, normal, hit distance and emission
func trace(ray *object.Line) (vec3.Vec3, vec3.Vec3, float64, float64, (func(vec3.Vec3, vec3.Vec3) vec3.Vec3)) {
	min_dist := inf
	closest_hit_color := vec3.Vec3{0, 0, 0}
	normal := vec3.Vec3{0, 0, 0}
	emission := 0.0
	var pdf func(vec3.Vec3, vec3.Vec3) vec3.Vec3

	for i := 0; i < len(objects); i++ {
		// iterate through object mesh (triangles)
		for j := 0; j < len(objects[i].Mesh); j++ {
			// camera ray
			intersection, hit_distance := objects[i].Mesh[j].Intersection(ray)
			if intersection {
				// only keep the closest intersections
				if hit_distance < min_dist {
					min_dist = hit_distance
					closest_hit_color = objects[i].Mesh[j].Mterial.Diffuse_color
					normal = surface_normal(&objects[i].Mesh[j])
					emission = objects[i].Mesh[j].Mterial.Emission
					pdf = objects[i].Mesh[j].Pdf
				}
			}
		}
	}

	for i := 0; i < len(spheres); i++ {
		intersection, hit_distance := spheres[i].Intersection(ray)
		if intersection {
			// only keep the closest intersections
			if hit_distance < min_dist {
				min_dist = hit_distance
				closest_hit_color = spheres[i].Mterial.Diffuse_color
				
				// compute normal
				hit_position := ray.Origin
				d := ray.Dir
				d.Scale(hit_distance)
				hit_position.Add(d)
				normal = hit_position
				normal.Sub(spheres[i].Origin)
				normal.Normalize()

				emission = spheres[i].Mterial.Emission
				pdf = spheres[i].Pdf
			}
		}
	}

	return closest_hit_color, normal, min_dist, emission, pdf
}

func surface_normal(tri *object.Triangle) vec3.Vec3 {
	a := tri.A
	normal := tri.B
	c := tri.C
	normal.Sub(a)
	c.Sub(a)
	normal.Cross(c)
	normal.Normalize()
	return normal
}

func max (a float64, b float64) float64 {
	if a < b {
		return b
	}

	return a
}

func min (a float64, b float64) float64 {
	if a > b {
		return b
	}

	return a
}

func cosine_hemisphere_sample() vec3.Vec3 {
	u1 := rand_float()
    r := math.Sqrt(u1)
    theta := 2 * math.Pi * rand_float()
 
    x := r * math.Cos(theta)
    y := r * math.Sin(theta)
 
    return vec3.Vec3{x, y, math.Sqrt(max(0.0, 1.0 - u1))}
}

func save_frame_buffer_to_png(frame_buffer [][]vec3.Vec3, output_name string) {
	// stores output image
	img := image.NewNRGBA(image.Rect(0, 0, len(frame_buffer[0]), len(frame_buffer)))

	for x := 0; x < len(frame_buffer); x++ {
		for y := 0; y < len(frame_buffer[0]); y++ {	
			// paint pixel to output.png
			img.Set(x, y, color.NRGBA {
				R: uint8(frame_buffer[x][y].X),
				G: uint8(frame_buffer[x][y].Y),
				B: uint8(frame_buffer[x][y].Z),
				A: 255,
			})	
		}
	}

	// create output file and catch errors
	f, err := os.Create(output_name + ".png")
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(f, img); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	fmt.Println("starty print :)")

	// init window

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Println("initializing SDL", err)
		return 
	}

	window, err := sdl.CreateWindow (
		"pathtracy boi",
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		window_width, window_height,
		sdl.WINDOW_OPENGL)

	if err != nil {
		fmt.Println("initializing window")
		return
	}

	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)

	if err != nil {
		fmt.Println("initializing renderer", err)
		return
	}

	defer renderer.Destroy()

	//
	// ************************************************************
	//
	// define materials
	green := object.Material {
		Diffuse_color : vec3.Vec3{0, 255, 0},
		Emission      : 0,
	}

	white := object.Material {
		Diffuse_color : vec3.Vec3{255, 255, 255},
		Emission      : 0,
	}

	blue := object.Material {
		Diffuse_color : vec3.Vec3{0, 0, 255},
		Emission      : 0,
	}

	red := object.Material {
		Diffuse_color : vec3.Vec3{255, 0, 0},
		Emission      : 0,
	}

	purple := object.Material {
		Diffuse_color : vec3.Vec3{200, 0, 200},
		Emission      : 0,
	}

	white_light := object.Material {
		Diffuse_color : vec3.Vec3{255, 255, 255},
		Emission      : 1.0,
	}

	_ = blue
	_ = red
	_ = green
	_ = purple
	_ = white

	diffuse_pdf := func(incident vec3.Vec3, n vec3.Vec3) vec3.Vec3 {
		direction := vec3.Vec3{rand_neg_float(), rand_neg_float(), rand_neg_float()}
		for {
			direction.Normalize()
			if direction.Dot(n) >= 0 {
				break
			}
			direction = vec3.Vec3{rand_neg_float(), rand_neg_float(), rand_neg_float()}
		}

		return direction
	}

	specular_pdf := func(incident vec3.Vec3, n vec3.Vec3) vec3.Vec3 {
		n.Scale(2*incident.Dot(n))
		incident.Sub(n)
		return incident
	}

	room_size := 1000.0/2

	// declare objects in 3d space
	room := object.Object {
		[](object.Triangle) {
			// back wall
			object.Triangle{vec3.Vec3{0, 0, room_size}, vec3.Vec3{0, room_size, room_size}, vec3.Vec3{room_size, room_size, room_size}, diffuse_pdf, green},
			object.Triangle{vec3.Vec3{0, 0, room_size}, vec3.Vec3{room_size, room_size, room_size}, vec3.Vec3{room_size, 0, room_size}, diffuse_pdf, green},
			// left wall
			object.Triangle{vec3.Vec3{0, 0, 0}, vec3.Vec3{0, room_size, 0}, vec3.Vec3{0, room_size, room_size}, specular_pdf, white},
			object.Triangle{vec3.Vec3{0, 0, 0}, vec3.Vec3{0, room_size, room_size}, vec3.Vec3{0, 0, room_size}, specular_pdf, white},
			// right wall
			object.Triangle{vec3.Vec3{room_size, room_size, room_size}, vec3.Vec3{room_size, room_size, 0}, vec3.Vec3{room_size, 0, 0}, specular_pdf, white},
			object.Triangle{vec3.Vec3{room_size, 0, room_size}, vec3.Vec3{room_size, room_size, room_size}, vec3.Vec3{room_size, 0, 0}, specular_pdf, white},
			// ceiling
			object.Triangle{vec3.Vec3{0, 0, 0}, vec3.Vec3{0, 0, room_size}, vec3.Vec3{room_size, 0, room_size}, diffuse_pdf, purple},
			object.Triangle{vec3.Vec3{0, 0, 0}, vec3.Vec3{room_size, 0, room_size}, vec3.Vec3{room_size, 0, 0}, diffuse_pdf, purple},
			// floor
			object.Triangle{vec3.Vec3{0, room_size, 0}, vec3.Vec3{room_size, room_size, 0}, vec3.Vec3{0, room_size, room_size}, specular_pdf, white},
			object.Triangle{vec3.Vec3{0, room_size, room_size}, vec3.Vec3{room_size, room_size, 0}, vec3.Vec3{room_size, room_size, room_size}, specular_pdf, white},
		},
	}

	cuboid_size := 10000.0

	cuboid := object.Object {
		[](object.Triangle) {
			// back wall
			object.Triangle{vec3.Vec3{0, 0, cuboid_size}, vec3.Vec3{0, cuboid_size, cuboid_size}, vec3.Vec3{cuboid_size, cuboid_size, cuboid_size}, diffuse_pdf, blue},
			object.Triangle{vec3.Vec3{0, 0, cuboid_size}, vec3.Vec3{cuboid_size, cuboid_size, cuboid_size}, vec3.Vec3{cuboid_size, 0, cuboid_size}, diffuse_pdf, blue},
			// left wall
			object.Triangle{vec3.Vec3{0, 0, 0}, vec3.Vec3{0, cuboid_size, 0}, vec3.Vec3{0, cuboid_size, cuboid_size}, diffuse_pdf, blue},
			object.Triangle{vec3.Vec3{0, 0, 0}, vec3.Vec3{0, cuboid_size, cuboid_size}, vec3.Vec3{0, 0, cuboid_size}, diffuse_pdf, blue},
			// right wall
			object.Triangle{vec3.Vec3{cuboid_size, cuboid_size, cuboid_size}, vec3.Vec3{cuboid_size, cuboid_size, 0}, vec3.Vec3{cuboid_size, 0, 0}, diffuse_pdf, blue},
			object.Triangle{vec3.Vec3{cuboid_size, 0, cuboid_size}, vec3.Vec3{cuboid_size, cuboid_size, cuboid_size}, vec3.Vec3{cuboid_size, 0, 0}, diffuse_pdf, blue},
			// ceiling
			object.Triangle{vec3.Vec3{0, 0, 0}, vec3.Vec3{0, 0, cuboid_size}, vec3.Vec3{cuboid_size, 0, cuboid_size}, diffuse_pdf, blue},
			object.Triangle{vec3.Vec3{0, 0, 0}, vec3.Vec3{cuboid_size, 0, cuboid_size}, vec3.Vec3{cuboid_size, 0, 0}, diffuse_pdf, blue},
			// floor
			object.Triangle{vec3.Vec3{0, cuboid_size, 0}, vec3.Vec3{cuboid_size, cuboid_size, 0}, vec3.Vec3{0, cuboid_size, cuboid_size}, diffuse_pdf, blue},
			object.Triangle{vec3.Vec3{0, cuboid_size, cuboid_size}, vec3.Vec3{cuboid_size, cuboid_size, 0}, vec3.Vec3{cuboid_size, cuboid_size, cuboid_size}, diffuse_pdf, blue},
			// front plane
			object.Triangle{vec3.Vec3{0, 0, 0}, vec3.Vec3{0, cuboid_size, 0}, vec3.Vec3{cuboid_size, cuboid_size, 0}, diffuse_pdf, blue},
			object.Triangle{vec3.Vec3{0, 0, 0}, vec3.Vec3{cuboid_size, cuboid_size, 0}, vec3.Vec3{cuboid_size, 0, 0}, diffuse_pdf, blue},
		},
	}

	cuboid2 := object.Object {
		[](object.Triangle) {
			// back wall
			object.Triangle{vec3.Vec3{0, 0, cuboid_size}, vec3.Vec3{0, cuboid_size, cuboid_size}, vec3.Vec3{cuboid_size, cuboid_size, cuboid_size}, diffuse_pdf, white},
			object.Triangle{vec3.Vec3{0, 0, cuboid_size}, vec3.Vec3{cuboid_size, cuboid_size, cuboid_size}, vec3.Vec3{cuboid_size, 0, cuboid_size}, diffuse_pdf, white},
			// left wall
			object.Triangle{vec3.Vec3{0, 0, 0}, vec3.Vec3{0, cuboid_size, 0}, vec3.Vec3{0, cuboid_size, cuboid_size}, diffuse_pdf, white},
			object.Triangle{vec3.Vec3{0, 0, 0}, vec3.Vec3{0, cuboid_size, cuboid_size}, vec3.Vec3{0, 0, cuboid_size}, diffuse_pdf, white},
			// right wall
			object.Triangle{vec3.Vec3{cuboid_size, cuboid_size, cuboid_size}, vec3.Vec3{cuboid_size, cuboid_size, 0}, vec3.Vec3{cuboid_size, 0, 0}, diffuse_pdf, white},
			object.Triangle{vec3.Vec3{cuboid_size, 0, cuboid_size}, vec3.Vec3{cuboid_size, cuboid_size, cuboid_size}, vec3.Vec3{cuboid_size, 0, 0}, diffuse_pdf, white},
			// ceiling
			object.Triangle{vec3.Vec3{0, 0, 0}, vec3.Vec3{0, 0, cuboid_size}, vec3.Vec3{cuboid_size, 0, cuboid_size}, diffuse_pdf, white},
			object.Triangle{vec3.Vec3{0, 0, 0}, vec3.Vec3{cuboid_size, 0, cuboid_size}, vec3.Vec3{cuboid_size, 0, 0}, diffuse_pdf, white},
			// floor
			object.Triangle{vec3.Vec3{0, cuboid_size, 0}, vec3.Vec3{cuboid_size, cuboid_size, 0}, vec3.Vec3{0, cuboid_size, cuboid_size}, diffuse_pdf, white},
			object.Triangle{vec3.Vec3{0, cuboid_size, cuboid_size}, vec3.Vec3{cuboid_size, cuboid_size, 0}, vec3.Vec3{cuboid_size, cuboid_size, cuboid_size}, diffuse_pdf, white},
			// front plane
			object.Triangle{vec3.Vec3{0, 0, 0}, vec3.Vec3{0, cuboid_size, 0}, vec3.Vec3{cuboid_size, cuboid_size, 0}, diffuse_pdf, white},
			object.Triangle{vec3.Vec3{0, 0, 0}, vec3.Vec3{cuboid_size, cuboid_size, 0}, vec3.Vec3{cuboid_size, 0, 0}, diffuse_pdf, white},
		},
	}

	cuboid3 := object.Object {
		[](object.Triangle) {
			// back wall
			object.Triangle{vec3.Vec3{0, 0, cuboid_size}, vec3.Vec3{0, cuboid_size, cuboid_size}, vec3.Vec3{cuboid_size, cuboid_size, cuboid_size}, diffuse_pdf, white},
			object.Triangle{vec3.Vec3{0, 0, cuboid_size}, vec3.Vec3{cuboid_size, cuboid_size, cuboid_size}, vec3.Vec3{cuboid_size, 0, cuboid_size}, diffuse_pdf, white},
			// left wall
			object.Triangle{vec3.Vec3{0, 0, 0}, vec3.Vec3{0, cuboid_size, 0}, vec3.Vec3{0, cuboid_size, cuboid_size}, diffuse_pdf, white},
			object.Triangle{vec3.Vec3{0, 0, 0}, vec3.Vec3{0, cuboid_size, cuboid_size}, vec3.Vec3{0, 0, cuboid_size}, diffuse_pdf, white},
			// right wall
			object.Triangle{vec3.Vec3{cuboid_size, cuboid_size, cuboid_size}, vec3.Vec3{cuboid_size, cuboid_size, 0}, vec3.Vec3{cuboid_size, 0, 0}, diffuse_pdf, white},
			object.Triangle{vec3.Vec3{cuboid_size, 0, cuboid_size}, vec3.Vec3{cuboid_size, cuboid_size, cuboid_size}, vec3.Vec3{cuboid_size, 0, 0}, diffuse_pdf, white},
			// ceiling
			object.Triangle{vec3.Vec3{0, 0, 0}, vec3.Vec3{0, 0, cuboid_size}, vec3.Vec3{cuboid_size, 0, cuboid_size}, diffuse_pdf, white},
			object.Triangle{vec3.Vec3{0, 0, 0}, vec3.Vec3{cuboid_size, 0, cuboid_size}, vec3.Vec3{cuboid_size, 0, 0}, diffuse_pdf, white},
			// floor
			object.Triangle{vec3.Vec3{0, cuboid_size, 0}, vec3.Vec3{cuboid_size, cuboid_size, 0}, vec3.Vec3{0, cuboid_size, cuboid_size}, diffuse_pdf, white},
			object.Triangle{vec3.Vec3{0, cuboid_size, cuboid_size}, vec3.Vec3{cuboid_size, cuboid_size, 0}, vec3.Vec3{cuboid_size, cuboid_size, cuboid_size}, diffuse_pdf, white},
			// front plane
			object.Triangle{vec3.Vec3{0, 0, 0}, vec3.Vec3{0, cuboid_size, 0}, vec3.Vec3{cuboid_size, cuboid_size, 0}, diffuse_pdf, white},
			object.Triangle{vec3.Vec3{0, 0, 0}, vec3.Vec3{cuboid_size, cuboid_size, 0}, vec3.Vec3{cuboid_size, 0, 0}, diffuse_pdf, white},
		},
	}

	cuboid4 := object.Object {
		[](object.Triangle) {
			// back wall
			object.Triangle{vec3.Vec3{0, 0, cuboid_size}, vec3.Vec3{0, cuboid_size, cuboid_size}, vec3.Vec3{cuboid_size, cuboid_size, cuboid_size}, diffuse_pdf, white},
			object.Triangle{vec3.Vec3{0, 0, cuboid_size}, vec3.Vec3{cuboid_size, cuboid_size, cuboid_size}, vec3.Vec3{cuboid_size, 0, cuboid_size}, diffuse_pdf, white},
			// left wall
			object.Triangle{vec3.Vec3{0, 0, 0}, vec3.Vec3{0, cuboid_size, 0}, vec3.Vec3{0, cuboid_size, cuboid_size}, diffuse_pdf, white},
			object.Triangle{vec3.Vec3{0, 0, 0}, vec3.Vec3{0, cuboid_size, cuboid_size}, vec3.Vec3{0, 0, cuboid_size}, diffuse_pdf, white},
			// right wall
			object.Triangle{vec3.Vec3{cuboid_size, cuboid_size, cuboid_size}, vec3.Vec3{cuboid_size, cuboid_size, 0}, vec3.Vec3{cuboid_size, 0, 0}, diffuse_pdf, white},
			object.Triangle{vec3.Vec3{cuboid_size, 0, cuboid_size}, vec3.Vec3{cuboid_size, cuboid_size, cuboid_size}, vec3.Vec3{cuboid_size, 0, 0}, diffuse_pdf, white},
			// ceiling
			object.Triangle{vec3.Vec3{0, 0, 0}, vec3.Vec3{0, 0, cuboid_size}, vec3.Vec3{cuboid_size, 0, cuboid_size}, diffuse_pdf, white},
			object.Triangle{vec3.Vec3{0, 0, 0}, vec3.Vec3{cuboid_size, 0, cuboid_size}, vec3.Vec3{cuboid_size, 0, 0}, diffuse_pdf, white},
			// floor
			object.Triangle{vec3.Vec3{0, cuboid_size, 0}, vec3.Vec3{cuboid_size, cuboid_size, 0}, vec3.Vec3{0, cuboid_size, cuboid_size}, diffuse_pdf, white},
			object.Triangle{vec3.Vec3{0, cuboid_size, cuboid_size}, vec3.Vec3{cuboid_size, cuboid_size, 0}, vec3.Vec3{cuboid_size, cuboid_size, cuboid_size}, diffuse_pdf, white},
			// front plane
			object.Triangle{vec3.Vec3{0, 0, 0}, vec3.Vec3{0, cuboid_size, 0}, vec3.Vec3{cuboid_size, cuboid_size, 0}, diffuse_pdf, white},
			object.Triangle{vec3.Vec3{0, 0, 0}, vec3.Vec3{cuboid_size, cuboid_size, 0}, vec3.Vec3{cuboid_size, 0, 0}, diffuse_pdf, white},
		},
	}

	// example of how you can move an object
	cuboid.Move(-2500, -2500, -5000)
	cuboid2.Move(1 + cuboid_size, 1 + cuboid_size, 1 + cuboid_size)

	cuboid3.Move(1 + 1.8*cuboid_size, 1 + 1.6*cuboid_size, 1 + 2*cuboid_size)
	cuboid3.Rotate_x(0.5)
	cuboid3.Rotate_y(0.3)
	cuboid3.Rotate_z(0.5)
	cuboid3.Move(-430, -270, -20)

	sphere1 := object.Sphere {
		vec3.Vec3{150, 150, 250},
		120.0,
		specular_pdf,
		white,
	}

	sphere2 := object.Sphere {
		vec3.Vec3{350, 350, 250},
		120.0,
		specular_pdf,
		white,
	}

	sphere3 := object.Sphere {
		vec3.Vec3{400, 100, 350},
		90.0,
		specular_pdf,
		white,
	}

	// output dimensions
	width := 1000/2
	height := 1000/2

	// define camera (and image output dimensions)
	camera := camera.Camera { 
		Width: width,
		Height: height,
		Origin: vec3.Vec3{float64(width/2), float64(height/2), -float64(width)},
	}

	// define light sources
	spotlight1_radius := 200.0/2

	depth := height

	lamp1 := object.Object {
		[]object.Triangle {
			object.Triangle {
				vec3.Vec3{float64(width/2) - float64(spotlight1_radius/2), 0.0000001, float64(depth/2) - float64(spotlight1_radius/2)}, 
				vec3.Vec3{float64(width/2) + float64(spotlight1_radius/2), 0.0000001, float64(width/2) - float64(spotlight1_radius/2)}, 
				vec3.Vec3{float64(width/2) - float64(spotlight1_radius/2), 0.0000001, float64(depth/2) + float64(spotlight1_radius/2)}, diffuse_pdf, white_light},
			object.Triangle {
				vec3.Vec3{float64(width/2) - float64(spotlight1_radius/2), 0.0000001, float64(depth/2) + float64(spotlight1_radius/2),}, 
				vec3.Vec3{float64(width/2) + float64(spotlight1_radius/2), 0.0000001, float64(depth/2) - float64(spotlight1_radius/2),}, 
				vec3.Vec3{float64(width/2) + float64(spotlight1_radius/2), 0.0000001, float64(depth/2) + float64(spotlight1_radius/2)}, diffuse_pdf, white_light},
		},
	}

	_ = cuboid
	_ = cuboid2
	_ = cuboid3
	_ = cuboid4
	_ = sphere1

	objects = append(objects, room, lamp1)
	spheres = append(spheres, sphere1, sphere2, sphere3)

	// CPU profiling by default
	// defer profile.Start().Stop()

	// cache random floats for quicker computation
	floats = make([]float64, float_amount)
	for i := 0; i < float_amount; i++ {
		floats[i] = rand.Float64()
	}
	
	// how many times a single pixel is sampled
	pixel_samples, err := strconv.Atoi(string(os.Args[1]))
	// how many times a ray bounces
	hops               := 7
	
	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.Clear()

	frame_buffer := render_frame(camera, pixel_samples, hops)

	for x := 0; x < len(frame_buffer); x++ {
		for y := 0; y < len(frame_buffer[0]); y++ {
			// draw pixels
			renderer.SetDrawColor(uint8(frame_buffer[x][y].X), uint8(frame_buffer[x][y].Y), uint8(frame_buffer[x][y].Z), 255)
			renderer.DrawPoint(int32(x), int32(y))
		} 
	}

	// show pixels on window
	renderer.Present()
	save_frame_buffer_to_png(frame_buffer, "output@" + strconv.Itoa(pixel_samples) + "_samples")
}

func render_frame_thread(start_x int, end_x int, start_y int, end_y int, camera camera.Camera, frame_buffer [][]vec3.Vec3, samples int, hops int, wg **sync.WaitGroup) {
	// multithreading magic, don't touch this
	_wg := *wg
	defer _wg.Done()

	// camera position data: compute this only once
	cam_x := camera.Origin.X - float64(camera.Width/2)
	cam_y := camera.Origin.Y - float64(camera.Height/2)
	zero_vector := vec3.Vec3{0, 0, 0}
	
	// this jumbo wumbo loop solves the rendering equations for path tracing
	for x := start_x; x < end_x; x++ {
		for y := start_y; y < end_y; y++ {
			// generate camera ray
			camera_ray_dir := vec3.Vec3 {
				cam_x + float64(x), 
				cam_y + float64(y), 
				// camera.Height is 1000, which is also the distance from camera to view plane
				// TODO find a nicer way to implement this
				camera.Origin.Z + float64(camera.Height),
			}
			
			camera_ray_dir.Sub(camera.Origin)
			camera_ray_dir.Normalize()
			color := zero_vector

			for s := 0; s < samples; s++ {
				hops_done := 0
				origin := camera.Origin
				direction := camera_ray_dir
				cur_weight := 1.0
				cur_color := zero_vector
				hit_a_light_source := false
				for h := 0; h < hops; h++ {
					hops_done += 1
					pixel_color, n, distance, emission, pdf := trace(&object.Line{origin, direction})

					// no intersection, ray probably left the cornel box
					if distance == inf {
						cur_color = zero_vector
						break
					}

					// ---------------------------------------------------------------------------
					// DO NOT USE closest_triangle BEFORE THIS LINE
					// ONLY USE closest_triangle IF INTERSECTION OCCURRED
					// ---------------------------------------------------------------------------

					// light attentuation (fall-off)
					pixel_color.Scale(math.Pow(0.95, float64(hops_done-1)))

					cur_color.Add(pixel_color)
					cur_weight += emission

					// hit a light source
					if emission > 0.0 {
						hit_a_light_source = true
						break
					}

					// bounce
					// update origin
					incident := direction
					_ = incident
					direction.Scale(distance)
					origin.Add(direction)

					// <update direction>
					direction = pdf(incident, n)
					// </update direction>
				}

				if hit_a_light_source {
					cur_color.Scale(cur_weight)
					color.Add(cur_color)
				}
			}
			color.Scale(1.0/float64(samples))

			// prevent overflow
			_max := color.X
			if color.Y > color.X {
				_max = color.Y
			}
			if color.Z > color.Y {
				_max = color.Z
			}

			if _max > 255 {
				color.Scale(255.0/_max)
			}

			frame_buffer[x][y] = color
			// gamma correction
			frame_buffer[x][y].X = math.Pow(frame_buffer[x][y].X/255, 1.0/2.20)
			frame_buffer[x][y].Y = math.Pow(frame_buffer[x][y].Y/255, 1.0/2.20)
			frame_buffer[x][y].Z = math.Pow(frame_buffer[x][y].Z/255, 1.0/2.20)
			frame_buffer[x][y].Scale(255)
		}
	}

	return
}

// renders a frame and generates an output png
func render_frame(camera camera.Camera, samples int, hops int) [][]vec3.Vec3 {
	// initialise g_buffers
	frame_buffer := make([][]vec3.Vec3, int(camera.Width/1))

	for x := 0; x < camera.Width; x++ {
		frame_buffer[x] = make([]vec3.Vec3, int(camera.Height/1))

		for y := 0; y < camera.Height; y++ {
			frame_buffer[x][y] = vec3.Vec3{0, 0, 0}
		}
	}

	// multithreading using n cpu-cores
	cores := 4
	wg := new(sync.WaitGroup)
	for c := 0; c < cores; c++ {
		wg.Add(1)
		go render_frame_thread (
			int(0), 
			int(camera.Width), 
			
			int(c*(camera.Height/(cores))), 
			int((c+1)*(camera.Height/(cores))), 
			
			camera, frame_buffer, samples, hops, &wg)
	}
	
	wg.Wait()

	return frame_buffer
}