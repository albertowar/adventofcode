package intcode

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
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

func ReadIntCodes(filename string) ([]int, error) {
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
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)
	value, err := strconv.Atoi(text)
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

func RunIntcodeProgram(intcodes []int, input chan int, output chan int, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}

	i := 0
	rawOpcode := intcodes[i]
	opcodeDigits := extractOpcodeDigits(rawOpcode)
	opcode := parseOpcode(opcodeDigits)

	for opcode != 99 {
		switch opcode {
		case 1:
			arg1 := extractArgument(intcodes, i, 1, parseArgMode(opcodeDigits, 1))
			arg2 := extractArgument(intcodes, i, 2, parseArgMode(opcodeDigits, 2))

			// arg3 is always in position mode
			arg3 := intcodes[i+3]

			intcodes[arg3] = arg1 + arg2
			i += 4
		case 2:
			arg1 := extractArgument(intcodes, i, 1, parseArgMode(opcodeDigits, 1))
			arg2 := extractArgument(intcodes, i, 2, parseArgMode(opcodeDigits, 2))

			// arg3 is always in position mode
			arg3 := intcodes[i+3]

			intcodes[arg3] = arg1 * arg2
			i += 4
		case 3:
			// Instruction 3 argument is always in position mode
			arg1 := extractArgument(intcodes, i, 1, 1)

			var value int
			if input != nil {
				value = <-input
			} else {
				value, _ = readInput()
			}

			intcodes[arg1] = value
			i += 2
		case 4:
			arg1 := extractArgument(intcodes, i, 1, parseArgMode(opcodeDigits, 1))

			if output != nil {
				output <- arg1
			} else {
				fmt.Println(arg1)
			}

			i += 2
		case 5:
			arg1 := extractArgument(intcodes, i, 1, parseArgMode(opcodeDigits, 1))

			if arg1 != 0 {
				arg2 := extractArgument(intcodes, i, 2, parseArgMode(opcodeDigits, 2))
				i = arg2
			} else {
				i += 3
			}
		case 6:
			arg1 := extractArgument(intcodes, i, 1, parseArgMode(opcodeDigits, 1))

			if arg1 == 0 {
				arg2 := extractArgument(intcodes, i, 2, parseArgMode(opcodeDigits, 2))
				i = arg2
			} else {
				i += 3
			}
		case 7:
			arg1 := extractArgument(intcodes, i, 1, parseArgMode(opcodeDigits, 1))
			arg2 := extractArgument(intcodes, i, 2, parseArgMode(opcodeDigits, 2))
			arg3 := extractArgument(intcodes, i, 3, 1)

			if arg1 < arg2 {
				intcodes[arg3] = 1
			} else {
				intcodes[arg3] = 0
			}

			i += 4
		case 8:
			arg1 := extractArgument(intcodes, i, 1, parseArgMode(opcodeDigits, 1))
			arg2 := extractArgument(intcodes, i, 2, parseArgMode(opcodeDigits, 2))
			arg3 := extractArgument(intcodes, i, 3, 1)

			if arg1 == arg2 {
				intcodes[arg3] = 1
			} else {
				intcodes[arg3] = 0
			}

			i += 4
		}

		rawOpcode = intcodes[i]
		opcodeDigits = extractOpcodeDigits(rawOpcode)
		opcode = parseOpcode(opcodeDigits)
	}
}
