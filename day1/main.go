package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
)

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

		requiredFuel := math.Trunc(containerMass/3) - 2
		totalFuel += requiredFuel
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(totalFuel)
}
