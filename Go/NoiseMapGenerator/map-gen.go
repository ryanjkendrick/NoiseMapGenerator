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

const (
 	alpha       = 2.
 	beta        = 2.
 	n           = 3
 	seed  int64 = 100
 )

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please specify the path to save the generated image.")
		os.Exit(1)
	}
	imgPath := os.Args[1]

	out, err := os.Create(imgPath)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	width, height := 1920, 1080
	background := color.RGBA{0, 0xFF, 0, 0xCC}

	if len(os.Args) == 4 {
		if n, err := strconv.Atoi(os.Args[2]); err == nil {
			width = n
		} else {
			log.Fatal(os.Args[2], "is not an integer.")
			os.Exit(1)
		}

		if n, err := strconv.Atoi(os.Args[3]); err == nil {
			height = n
		} else {
			log.Fatal(os.Args[3], "is not an integer.")
			os.Exit(1)
		}
	}

	img := createImage(width, height, background)
	log.Print("Image created.")

	img = perlinizeImage(width, height, img)
	log.Print("Image randomized.")

	if strings.HasSuffix(strings.ToLower(imgPath), ".jpg") {
		var opt jpeg.Options
		opt.Quality = 100
		err = jpeg.Encode(out, img, &opt)
	} else {
		err = png.Encode(out, img)
	}

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	log.Printf("Image saved to %s.\n", imgPath)
}

func createImage(width int, height int, background color.RGBA) *image.RGBA {
	rect := image.Rect(0, 0, width, height)
	img := image.NewRGBA(rect)
	draw.Draw(img, img.Bounds(), &image.Uniform{background}, image.ZP, draw.Src)
	return img
}

func randomizeImage(width int, height int, img *image.RGBA) *image.RGBA {
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			pixelVal := rand.Intn(16777215 * 2)
			r, g, b, a := calcColor(pixelVal)

			img.SetRGBA(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
		}
	}

	return img
}

func perlinizeImage(width int, height int, img *image.RGBA) *image.RGBA {
        for x := 0; x < width; x++ {
                for y := 0; y < height; y++ {
			randVal := rand.Intn(16777215 * 2)
			p := perlin.NewPerlin(alpha, beta, randVal, seed)
                        pixelVal := int(p.Noise2D(float64(x), float64(y))*200)

			log.Print(pixelVal)

                        r, g, b, a := calcColor(pixelVal)

                        img.SetRGBA(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
                }
        }
        return img
}

func calcColor(color int) (red, green, blue, alpha int) {
	alpha = color & 0xFF
	blue = (color >> 8) & 0xFF
	green = (color >> 16) & 0xFF
	red = (color >> 24) & 0xFF

	return red, green, blue, alpha
}
