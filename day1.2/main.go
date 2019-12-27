package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
)

func calculateRequiredFuel(mass float64) float64 {
	if mass <= 0.0 {
		return 0.0
	}

	requiredFuel := math.Trunc(mass/3) - 2

	if requiredFuel <= 0.0 {
		return 0.0
	}

	requiredFuelForFuel := calculateRequiredFuel(requiredFuel)

	return requiredFuel + requiredFuelForFuel
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	totalFuel := 0.0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		containerMass, err := strconv.ParseFloat(text, 64)

		if err != nil {
			break
		}

		totalFuel += calculateRequiredFuel(containerMass)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(totalFuel)
}
