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

// Point something
type Point struct {
	X float64
	Y float64
}

// Segment something
type Segment struct {
	Src Point
	Dst Point
}

func (p Point) isInSegment(s Segment) bool {
	ap := Point{X: p.X - s.Src.X, Y: p.Y - s.Src.Y}
	ab := Point{X: s.Dst.X - s.Src.X, Y: s.Dst.Y - s.Src.Y}

	kap := ab.X*ap.X + ab.Y*ap.Y
	kab := ab.X*ab.X + ab.Y*ab.Y

	return kap == 0 || kap == kab || (0 < kap && kap < kab)
}

func extractSegments(directions []string) []Segment {
	lines := make([]Segment, 0)

	src := Point{X: 0.0, Y: 0.0}
	var dst Point

	for _, direction := range directions {
		units, err := strconv.ParseFloat(direction[1:], 64)
		dst = src

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

		newLine := Segment{Src: src, Dst: dst}
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

// Source: https://en.wikipedia.org/wiki/Line%E2%80%93line_intersection
func tvalue(p1 Point, p2 Point, p3 Point, p4 Point) float64 {
	numerator := (p1.X-p3.X)*(p3.Y-p4.Y) - (p1.Y-p3.Y)*(p3.X-p4.X)
	denominator := (p1.X-p2.X)*(p3.Y-p4.Y) - (p1.Y-p2.Y)*(p3.X-p4.X)

	if denominator == 0.0 {
		return math.MaxFloat64
	}

	return numerator / denominator
}

// Source: https://en.wikipedia.org/wiki/Line%E2%80%93line_intersection
func uvalue(p1 Point, p2 Point, p3 Point, p4 Point) float64 {
	numerator := (p1.X-p2.X)*(p1.Y-p3.Y) - (p1.Y-p2.Y)*(p1.X-p3.X)
	denominator := (p1.X-p2.X)*(p3.Y-p4.Y) - (p1.Y-p2.Y)*(p3.X-p4.X)

	if denominator == 0 {
		return math.MaxFloat64
	}

	return -numerator / denominator
}

func findStepsToIntersection(segments []Segment, intersection Point) int {
	steps := 0

	for _, s := range segments {
		steps += int(math.Abs(s.Dst.X-s.Src.X) + math.Abs(s.Dst.Y-s.Src.Y))

		if intersection.isInSegment(s) {
			break
		}
	}

	return steps
}

func main() {
	file, err := os.Open("test.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	wire1, err := getDirections(scanner)
	wire1Segments := extractSegments(wire1)

	wire2, err := getDirections(scanner)
	wire2Segments := extractSegments(wire2)

	intersections := make([]Point, 0)

	for _, segment1 := range wire1Segments {
		for _, segment2 := range wire2Segments {
			// Source: https://en.wikipedia.org/wiki/Line%E2%80%93line_intersection
			p1 := segment1.Src
			p2 := segment1.Dst
			p3 := segment2.Src
			p4 := segment2.Dst

			t := tvalue(p1, p2, p3, p4)
			u := uvalue(p1, p2, p3, p4)

			if t >= 0.0 && t <= 1.0 && u >= 0 && u <= 1 {
				intersection := Point{X: p1.X + t*(p2.X-p1.X), Y: p1.Y + t*(p2.Y-p1.Y)}
				intersections = append(intersections, intersection)
			}
		}
	}

	minDistance := math.MaxFloat64

	stepsToIntersection := make(map[Point]int)

	for _, intersection := range intersections {
		// Skip the central port since both wires start there
		if intersection.X == 0 && intersection.Y == 0 {
			continue
		}

		// Source: https://en.wikipedia.org/wiki/Taxicab_geometry
		distance := math.Abs(0.0-intersection.X) + math.Abs(0.0-intersection.Y)

		if distance < minDistance {
			minDistance = distance
		}

		wire1Steps := findStepsToIntersection(wire1Segments, intersection)
		wire2Steps := findStepsToIntersection(wire2Segments, intersection)
		stepsToIntersection[intersection] = wire1Steps + wire2Steps
	}

	fmt.Printf("Found %d intersections\n", len(intersections))
	fmt.Printf("Min distance %f\n", minDistance)

	minSteps := math.MaxInt32

	for _, steps := range stepsToIntersection {
		if steps < minSteps {
			minSteps = steps
		}
	}

	fmt.Printf("Min steps %d\n", minSteps)
}