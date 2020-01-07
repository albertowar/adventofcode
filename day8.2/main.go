package main

import (
	"bufio"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
)

const PIXELS_WIDE int = 25
const PIXELS_TALL int = 6

//const PIXELS_WIDE int = 3
//const PIXELS_TALL int = 2

func extractDigits(filename string) []int {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Scan()

	text := scanner.Text()

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	digits := make([]int, 0)

	for _, digitRune := range text {
		digit := int(digitRune - '0')
		digits = append(digits, digit)
	}

	return digits
}

func createImage(imagePixels [PIXELS_TALL][PIXELS_WIDE]int) {
	upLeft := image.Point{0, 0}
	lowRight := image.Point{PIXELS_TALL, PIXELS_WIDE}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	for x := 0; x < PIXELS_TALL; x++ {
		for y := 0; y < PIXELS_WIDE; y++ {
			switch imagePixels[x][y] {
			case 0:
				img.Set(x, y, color.Black)
			case 1:
				img.Set(x, y, color.White)
			case 2:
				img.Set(x, y, color.Transparent)
			}
		}
	}

	// Encode as PNG.
	f, _ := os.Create("image.png")
	png.Encode(f, img)
}

func main() {
	digits := extractDigits("input.txt")

	l, r, c := 0, 0, 0
	var layers [][PIXELS_TALL][PIXELS_WIDE]int
	var layer [PIXELS_TALL][PIXELS_WIDE]int
	layers = append(layers, layer)

	for _, digit := range digits {
		if c == PIXELS_WIDE {
			c = 0
			r++
		}

		if r == PIXELS_TALL {
			r = 0
			l++
			var newLayer [PIXELS_TALL][PIXELS_WIDE]int
			layers = append(layers, newLayer)
		}

		layers[l][r][c] = digit
		c++
	}

	var finalImage [PIXELS_TALL][PIXELS_WIDE]int

	for r := 0; r < PIXELS_TALL; r++ {
		for c := 0; c < PIXELS_WIDE; c++ {
			finalImage[r][c] = layers[0][r][c]
			for l := 0; l < len(layers); l++ {
				if layers[l][r][c] < finalImage[r][c] && finalImage[r][c] == 2 {
					finalImage[r][c] = layers[l][r][c]
				}
			}
		}
	}

	createImage(finalImage)
}
