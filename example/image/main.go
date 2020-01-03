package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math/rand"
	"net/http"
)

func getRandomColor() color.RGBA {
	return color.RGBA{
		R: uint8(rand.Uint32() % 255),
		G: uint8(rand.Uint32() % 255),
		B: uint8(rand.Uint32() % 255),
		A: 255,
	}
}

func getImage(w int, h int, c color.RGBA) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 150, 150))

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			img.Set(x, y, c)
		}
	}

	return img
}

func main() {


	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")

		c := getRandomColor()
		img := getImage(150, 150, c)

		if err := png.Encode(w, img); err != nil {
			fmt.Println(err)
		}
	})

	if err := http.ListenAndServe(":3001", nil); err != nil {
		log.Fatal(err)
	}
}