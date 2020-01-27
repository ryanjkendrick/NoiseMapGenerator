package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/aquilax/go-perlin"
)

// Perlin
const (
	alpha       = 1000. // His default was 2
	beta        = 10.
	n           = 3
	seed  int64 = 100
)

const (
	gridSize          = 4
	threadPayloadSize = 10
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please specify the path to save the generated image.")
	}
	imgPath := os.Args[1]

	out, err := os.Create(imgPath)
	if err != nil {
		log.Fatal(err)
	}

	width, height := 1920, 1080
	background := color.RGBA{0xFF, 0, 0, 0xCC}
	randomizeTheImage := false

	if len(os.Args) >= 3 {
		if n, err := strconv.Atoi(os.Args[2]); err == nil {
			width = n
		} else if os.Args[2] == "-r" {
			randomizeTheImage = true
		} else {
			log.Fatal(os.Args[2], "is not an integer.")
		}

		if !randomizeTheImage { // Randomize flag should be last parameter
			if n, err := strconv.Atoi(os.Args[3]); err == nil {
				height = n
			} else if os.Args[3] == "-r" {
				randomizeTheImage = true
			} else {
				log.Fatal(os.Args[3], "is not an integer.")
			}
		}
	}

	if len(os.Args) == 5 && os.Args[4] == "-r" {
		randomizeTheImage = true
	}

	log.Print("Generating Noise Map...")

	var pixels [][]int
	if randomizeTheImage {
		pixels = createRandomMap(width, height)
		log.Print("Map randomized.")
	} else {
		pixels = createPerlinMap(width, height)
		log.Print("Map perlinized.")
	}

	img := createImage(width, height, background)
	log.Print("Image created.")

	img = convertArrayToImage(width, height, img, pixels)
	log.Print("Mapped to image.")

	if path := strings.ToLower(imgPath); strings.HasSuffix(path, ".jpg") || strings.HasSuffix(path, ".jpeg") {
		var opt jpeg.Options
		opt.Quality = 100
		err = jpeg.Encode(out, img, &opt)
	} else {
		err = png.Encode(out, img)
	}

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Image saved to %s.\n", imgPath)
}

func createImage(width int, height int, background color.RGBA) *image.RGBA {
	rect := image.Rect(0, 0, width, height)
	img := image.NewRGBA(rect)
	draw.Draw(img, img.Bounds(), &image.Uniform{background}, image.ZP, draw.Src)
	return img
}

func createRandomMap(width int, height int) [][]int {
	pixels := make([][]int, width)

	for x := 0; x < width; x++ {
		pixels[x] = make([]int, height)

		for y := 0; y < height; y++ {
			pixelVal := rand.Intn(16777215)

			pixels[x][y] = pixelVal
		}
	}

	return pixels
}

func createPerlinMap(width int, height int) [][]int {
	pixels := make([][]int, width)

	for x := 0; x < width; x++ {
		pixels[x] = make([]int, height)

		for y := 0; y < height; y++ {
			p := perlin.NewPerlinRandSource(alpha, beta, n, rand.NewSource(seed))
			pixelVal := p.Noise2D(float64(x)/10, float64(y)/10) * 100000000.0

			if pixelVal < 0 { // Hex color value cannot be negative
				pixelVal *= -1
			}

			pixels[x][y] = int(pixelVal)
		}
	}

	return pixels
}

func calcColor(color int) (red, green, blue, alpha int) {
	//log.Print(color)

	alpha = color & 0xFF
	blue = (color >> 8) & 0xFF
	green = (color >> 16) & 0xFF
	red = (color >> 24) & 0xFF

	return red, green, blue, alpha
}

func convertArrayToImage(width int, height int, img *image.RGBA, pixels [][]int) *image.RGBA {
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			r, g, b, a := calcColor(pixels[x][y])

			img.SetRGBA(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
		}
	}
	return img
}
