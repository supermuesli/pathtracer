package main

import (
	"github.com/supermuesli/pathtracer/vec3"
	// "github.com/pkg/profile"
	"github.com/veandco/go-sdl2/sdl"
    "sync"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"fmt"
)

const (
	window_width = 1366
	window_height = 768
)

var input_images [][][]vec3.Vec3
var input_images_eyes [][][]vec3.Vec3
var frame_buffer [][]vec3.Vec3

func save_frame_buffer_to_png(output_name string) {
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

func read_image(path string) [][]vec3.Vec3 {
	_frame_buffer := make([][]vec3.Vec3, window_width)

	for x := 0; x < window_width; x++ {
		_frame_buffer[x] = make([]vec3.Vec3, window_height)

		for y := 0; y < window_height; y++ {
			_frame_buffer[x][y] = vec3.Vec3{0, 0, 0}
		}
	}

    infile, err := os.Open(path)
    if err != nil {
		log.Fatal(err)
    }
    defer infile.Close()

    // Decode will figure out what type of image is in the file on its own.
    // We just have to be sure all the image packages we want are imported.
    src, _, err := image.Decode(infile)
    if err != nil {
		log.Fatal(err)
    }

    // Create a new  image
    bounds := src.Bounds()
    w, h := bounds.Max.X, bounds.Max.Y
    for x := 0; x < w; x++ {
        for y := 0; y < h; y++ {
        	R, G, B, _ := src.At(x, y).RGBA()
            _frame_buffer[x][y] = vec3.Vec3{float64(R%256), float64(G%256), float64(B%256)}
        }
    }

    return _frame_buffer
}

func main() {
	fmt.Println("starty print :)")

	// init window

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Println("initializing SDL", err)
		return 
	}

	window, err := sdl.CreateWindow (
		"morphy boi",
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

	input_images = make([][][]vec3.Vec3, 2)

	for a := 0; a < 2; a++ {
		input_images[a] = make([][]vec3.Vec3, window_width)

		for x := 0; x < window_width; x++ {
			input_images[a][x] = make([]vec3.Vec3, window_height)

			for y := 0; y < window_height; y++ {
				input_images[a][x][y] = vec3.Vec3{0, 0, 0}
			}
		}
	}

	input_images_eyes = make([][][]vec3.Vec3, 2)

	for a := 0; a < 2; a++ {
		input_images_eyes[a] = make([][]vec3.Vec3, window_width)

		for x := 0; x < window_width; x++ {
			input_images_eyes[a][x] = make([]vec3.Vec3, window_height)

			for y := 0; y < window_height; y++ {
				input_images_eyes[a][x][y] = vec3.Vec3{0, 0, 0}
			}
		}
	}

	frame_buffer = make([][]vec3.Vec3, window_width)

	for x := 0; x < window_width; x++ {
		frame_buffer[x] = make([]vec3.Vec3, window_height)

		for y := 0; y < window_height; y++ {
			frame_buffer[x][y] = vec3.Vec3{0, 0, 0}
		}
	}

	input_images[0] = read_image("cat1.png")
	input_images[1] = read_image("cat2.png")
	input_images_eyes[0] = read_image("cat1_eyes.png")
	input_images_eyes[1] = read_image("cat2_eyes.png")

	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.Clear()

	frames := 100
	i := 0
	render_frame(i, frames)
	// game loop
	for {
		i = (i+1)%frames

		for x := 0; x < len(frame_buffer); x++ {
			for y := 0; y < len(frame_buffer[0]); y++ {
				// draw pixels
				renderer.SetDrawColor(uint8(frame_buffer[x][y].X), uint8(frame_buffer[x][y].Y), uint8(frame_buffer[x][y].Z), 255)
				renderer.DrawPoint(int32(x), int32(y))
			} 
		}

		render_frame(i, frames)

		// show pixels on window
		renderer.Present()
	}
}

func render_frame_thread(start_x int, end_x int, start_y int, end_y int, frame_number int, frame_cap int, wg **sync.WaitGroup) {
	// multithreading magic, don't touch this
	_wg := *wg
	defer _wg.Done()

	a := (float64(frame_cap) - float64(frame_number))/float64(frame_cap)
	b := 1.0 - a
	fmt.Println(a, b)

	for x := start_x; x < end_x; x++ {
		for y := start_y; y < end_y; y++ {
			frame_buffer[x][y] = vec3.Vec3 {
				((a*(input_images[0][x][y].X)) + (b*(input_images[1][x][y].X))), 
				((a*(input_images[0][x][y].Y)) + (b*(input_images[1][x][y].Y))), 
				((a*(input_images[0][x][y].Z)) + (b*(input_images[1][x][y].Z))),
			}
		}
	}

	return
}

// renders a frame and generates an output png
func render_frame(frame_number int, frame_cap int) {
	// multithreading using n cpu-cores
	cores := 4
	wg := new(sync.WaitGroup)
	for c := 0; c < cores; c++ {
		wg.Add(1)
		go render_frame_thread (
			int(0), 
			int(window_width), 
			
			int(c*(window_height/(cores))), 
			int((c+1)*(window_height/(cores))), 
			
			frame_number, frame_cap, &wg)
	}
	
	wg.Wait()
}