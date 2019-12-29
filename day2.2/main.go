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

func findOpcode(noun int, verb int, intcodes []int) {
	defer func() {
		if r := recover(); r != nil {
		}
	}()

	intcodes[1] = noun
	intcodes[2] = verb

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

	if intcodes[0] == 19690720 {
		fmt.Printf("Noun: %d Verb: %d Result: %d\n", noun, verb, intcodes[0])
		fmt.Println(100*noun + verb)
	}
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

	intcodesOriginal, err := parseIntArray(tokens)

	fmt.Printf("Tokens: %d Intcodes: %d\n", len(tokens), len(intcodesOriginal))

	if err != nil {
		log.Fatal(err)
	}

	for noun := 0; noun <= len(intcodesOriginal); noun++ {
		for verb := 0; verb <= len(intcodesOriginal); verb++ {
			intcodes := make([]int, len(intcodesOriginal))
			copy(intcodes, intcodesOriginal)

			go findOpcode(noun, verb, intcodes)
		}
	}
}
