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

func parseIntArray(input []string) ([]int64, error) {
	output := make([]int64, 0, 10*len(input))

	for _, element := range input {
		i, err := strconv.Atoi(element)
		if err != nil {
			return output, err
		}
		output = append(output, i)
	}
	return output, nil
}

func ReadIntCodes(filename string) ([]int64, error) {
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

func readInput() (int64, error) {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)
	value, err := strconv.Atoi(text)
	return value, err
}

func extractOpcodeDigits(rawOpcode int64) []int64 {
	digits := make([]int64, 0)

	for rawOpcode > 0 {
		digit := rawOpcode % 10
		rawOpcode = rawOpcode / 10

		digits = append([]int64{digit}, digits...)
	}

	return digits
}

var relativeBase int = 0

func extractArgument(intcodes []int64, start int64, argNumber int64, mode int64) int64 {
	var arg int64

	switch mode {
	case 0:
		arg = intcodes[intcodes[start+argNumber]]
	case 1:
		arg = intcodes[start+argNumber]
	case 2:
		arg = intcodes[relativeBase+intcodes[start+argNumber]]
	}

	return arg
}

func parseOpcode(opcodeDigits []int64) int64 {
	last := len(opcodeDigits) - 1

	var opcode int64
	if len(opcodeDigits) == 1 {
		opcode = opcodeDigits[last]
	} else {
		opcode = 10*opcodeDigits[last-1] + opcodeDigits[last]
	}

	return opcode
}

func parseArgMode(opcodeDigits []int64, argNumber int64) int64 {
	var mode int64

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

func RunIntcodeProgram(intcodes []int64, input chan int64, output chan int64, wg *sync.WaitGroup) {
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

			var value int64
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
		case 9:
			arg1 := extractArgument(intcodes, i, 1, parseArgMode(opcodeDigits, 1))

			relativeBase += arg1

			i += 2
		}
	}
}
