package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

// Coordinate something
type Point struct {
	X int
	Y int
}

// Line something
type Line struct {
	Src Point
	Dst Point
}

func intAbs(number int) int {
	if number < 0 {
		return -number
	}

	return number
}

func extractLines(directions []string) []Line {
	lines := make([]Line, 0)

	src := Point{X: 0, Y: 0}
	dst := src

	for _, direction := range directions {
		units, err := strconv.Atoi(direction[1:])

		if err != nil {
			panic(err)
		}

		if direction[0] == 'D' {
			dst.Y -= units
		} else if direction[0] == 'U' {
			dst.Y += units
		} else if direction[0] == 'R' {
			dst.X += units
		} else if direction[0] == 'L' {
			dst.X -= units
		}

		newLine := Line{Src: src, Dst: dst}
		lines = append(lines, newLine)
		src = dst
	}

	return lines
}

func getDirections(scanner *bufio.Scanner) ([]string, error) {
	scanner.Scan()

	text := scanner.Text()

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return strings.Split(text, ","), nil
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	wire1, err := getDirections(scanner)
	wireLines1 := extractLines(wire1)

	wire2, err := getDirections(scanner)
	wireLines2 := extractLines(wire2)

	intersections := make([]Point, 0)

	for _, wire1Line := range wireLines1 {
		for _, wire2Line := range wireLines2 {
			// Source: https://en.wikipedia.org/wiki/Line%E2%80%93line_intersection
			numerator := (wire1Line.Src.X-wire2Line.Src.X)*(wire2Line.Src.Y-wire2Line.Dst.Y) - (wire1Line.Src.Y-wire2Line.Src.Y)*(wire2Line.Src.X-wire2Line.Dst.X)
			denominator := (wire1Line.Src.X-wire1Line.Dst.X)*(wire2Line.Src.Y-wire2Line.Dst.Y) - (wire1Line.Src.Y-wire1Line.Dst.Y)*(wire2Line.Src.X-wire2Line.Dst.X)

			if denominator != 0 {
				t := numerator / denominator
				x := wire1Line.Src.X + t*(wire1Line.Dst.X-wire1Line.Src.X)
				y := wire1Line.Src.Y + t*(wire1Line.Dst.Y-wire1Line.Src.Y)

				intersection := Point{X: x, Y: y}
				intersections = append(intersections, intersection)
			}
		}
	}

	minDistance := math.MaxInt32

	for _, intersection := range intersections {
		// Skip the central port since both wires start there
		if intersection.X == 0 && intersection.Y == 0 {
			continue
		}

		// Source: https://en.wikipedia.org/wiki/Taxicab_geometry
		distance := intAbs(0-intersection.X) + intAbs(0-intersection.Y)

		if distance < minDistance {
			minDistance = distance
		}
	}

	fmt.Printf("Found %d intersections\n", len(intersections))
	fmt.Printf("Min distance %d\n", minDistance)
}
