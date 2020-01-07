package main

import (
	"log"

	"github.com/albertowar/adventofcode/intcode"
)

func main() {
	intcodes, err := intcode.ReadIntCodes("input.txt")

	if err != nil {
		log.Fatal(err)
	}
}
