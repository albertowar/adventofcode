package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"math"
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

func countLayerDigits(layer [PIXELS_TALL][PIXELS_WIDE] int, candidateDigit int) int {
	count := 0

	for _, row := range layer {
		for _, digit := range row {
			if digit == candidateDigit {
				count++
			}
		}
	}

	return count
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

	totalZeros := math.MaxInt32
	var layerWithFewestZeros [PIXELS_TALL][PIXELS_WIDE]int
	for _, layer := range layers {
		count := countLayerDigits(layer, 0)

		if count < totalZeros {
			totalZeros = count
			layerWithFewestZeros = layer
		}
	}

	totalOnes := countLayerDigits(layerWithFewestZeros, 1)
	totalTwos := countLayerDigits(layerWithFewestZeros, 2)
	fmt.Println(layerWithFewestZeros)
	fmt.Println(totalOnes * totalTwos)
}
