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

func readIntCodes(filename string) ([]int, error) {
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

	tokens := strings.Split(text, ",")

	return parseIntArray(tokens)
}

func readInput() (int, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Introduce the ID of the system to test: ")
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)
	value, err := strconv.Atoi(text)
	fmt.Printf("Entered %d\n", value)
	return value, err
}

func extractOpcodeDigits(rawOpcode int) []int {
	digits := make([]int, 0)

	for rawOpcode > 0 {
		digit := rawOpcode % 10
		rawOpcode = rawOpcode / 10

		digits = append([]int{digit}, digits...)
	}

	return digits
}

func extractArgument(intcodes []int, start int, argNumber int, mode int) int {
	var arg int

	if mode == 0 {
		arg = intcodes[intcodes[start+argNumber]]
	} else {
		arg = intcodes[start+argNumber]
	}

	return arg
}

func parseOpcode(opcodeDigits []int) int {
	last := len(opcodeDigits) - 1

	var opcode int
	if len(opcodeDigits) == 1 {
		opcode = opcodeDigits[last]
	} else {
		opcode = 10*opcodeDigits[last-1] + opcodeDigits[last]
	}

	return opcode
}

func parseArgMode(opcodeDigits []int, argNumber int) int {
	var mode int

	last := len(opcodeDigits) - 1

	switch argNumber {
	case 1:
		if len(opcodeDigits) > 2 {
			mode = opcodeDigits[last-2]
		} else {
			mode = 0
		}
	case 2:
		if len(opcodeDigits) > 3 {
			mode = opcodeDigits[last-3]
		} else {
			mode = 0
		}
	default:
		mode = 0
	}

	return mode
}

func runIntcodeProgram(intcodes []int) {
	i := 0
	rawOpcode := intcodes[i]
	opcodeDigits := extractOpcodeDigits(rawOpcode)
	opcode := parseOpcode(opcodeDigits)

	for opcode != 99 {
		//time.Sleep(1 * time.Second)
		fmt.Println(opcodeDigits)
		if opcode == 1 {
			mode1 := parseArgMode(opcodeDigits, 1)
			mode2 := parseArgMode(opcodeDigits, 2)

			arg1 := extractArgument(intcodes, i, 1, mode1)
			arg2 := extractArgument(intcodes, i, 2, mode2)

			// arg3 is always in position mode
			arg3 := intcodes[i+3]
			fmt.Printf("Memory[%d] <- %d + %d\n", arg3, arg1, arg2)

			intcodes[arg3] = arg1 + arg2
			i += 4
		} else if opcode == 2 {
			mode1 := parseArgMode(opcodeDigits, 1)
			mode2 := parseArgMode(opcodeDigits, 2)

			arg1 := extractArgument(intcodes, i, 1, mode1)
			arg2 := extractArgument(intcodes, i, 2, mode2)

			// arg3 is always in position mode
			arg3 := intcodes[i+3]
			fmt.Printf("Memory[%d] <- %d * %d\n", arg3, arg1, arg2)

			intcodes[arg3] = arg1 * arg2
			i += 4
		} else if opcode == 3 {
			// Instruction 3 argument is always in position mode
			arg1 := extractArgument(intcodes, i, 1, 1)

			value, _ := readInput()
			fmt.Printf("Memory[%d] <- %d\n", arg1, value)

			intcodes[arg1] = value
			i += 2
		} else if opcode == 4 {
			mode1 := parseArgMode(opcodeDigits, 1)

			arg1 := extractArgument(intcodes, i, 1, mode1)

			fmt.Println(arg1)
			i += 2
		}

		rawOpcode = intcodes[i]
		opcodeDigits = extractOpcodeDigits(rawOpcode)
		opcode = parseOpcode(opcodeDigits)
	}
}

func main() {
	intcodes, err := readIntCodes("input.txt")

	if err != nil {
		log.Fatal(err)
	}

	runIntcodeProgram(intcodes)
}
