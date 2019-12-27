package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func parseIntArray(input []string) ([]int, error) {
	output := make([]int, 0, len(input))

	for _, element := range input {
		i, err := strconv.Atoi(element)
		if err != nil {
			return output, err
		}
		output = append(output, i)
	}
	return output, nil
}

func main() {
	file, err := os.Open("input.txt")
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

	tokens := strings.Split(text, ",")

	intcodes, err := parseIntArray(tokens)

	fmt.Printf("Tokens: %d Intcodes: %d\n", len(tokens), len(intcodes))

	if err != nil {
		log.Fatal(err)
	}

	intcodes[1] = 12
	intcodes[2] = 2

	i := 0
	opcode := intcodes[i]

	for opcode != 99 {
		position1 := intcodes[i+1]
		position2 := intcodes[i+2]
		position3 := intcodes[i+3]

		if opcode == 1 {
			intcodes[position3] = intcodes[position1] + intcodes[position2]
		} else if opcode == 2 {
			intcodes[position3] = intcodes[position1] * intcodes[position2]
		}

		i += 4
		opcode = intcodes[i]
	}

	fmt.Println(intcodes[0])
}
