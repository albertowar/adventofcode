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

func runIntcodeProgram(intcodes []int, input chan int, output chan int) {
	i := 0
	rawOpcode := intcodes[i]
	opcodeDigits := extractOpcodeDigits(rawOpcode)
	opcode := parseOpcode(opcodeDigits)

	for opcode != 99 {
		//time.Sleep(1 * time.Second)
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

func runProgram(intcodes []int, inputChannel chan int, outputChannel chan int, phase int, input int) int {
	go runIntcodeProgram(intcodes, inputChannel, outputChannel)
	inputChannel <- phase
	inputChannel <- input

	output := <-outputChannel

	return output
}

func allDifferent(a int, b int, c int, d int, e int) bool {
	return a != b && a != c && a != d && a != e && b != c && b != d && b != e && c != d && c != e && d != e
}

func main() {
	inputChannel := make(chan int)
	outputChannel := make(chan int)

	maxThrusterSignal := 0

	for phaseA := 0; phaseA < 5; phaseA++ {
		for phaseB := 0; phaseB < 5; phaseB++ {
			for phaseC := 0; phaseC < 5; phaseC++ {
				for phaseD := 0; phaseD < 5; phaseD++ {
					for phaseE := 0; phaseE < 5; phaseE++ {
						if allDifferent(phaseA, phaseB, phaseC, phaseD, phaseE) {
							intcodes, err := readIntCodes("input.txt")
							if err != nil {
								log.Fatal(err)
							}

							outputA := runProgram(intcodes, inputChannel, outputChannel, phaseA, 0)
							outputB := runProgram(intcodes, inputChannel, outputChannel, phaseB, outputA)
							outputC := runProgram(intcodes, inputChannel, outputChannel, phaseC, outputB)
							outputD := runProgram(intcodes, inputChannel, outputChannel, phaseD, outputC)
							outputE := runProgram(intcodes, inputChannel, outputChannel, phaseE, outputD)

							if outputE > maxThrusterSignal {
								maxThrusterSignal = outputE
							}
						}
					}
				}
			}
		}
	}

	fmt.Printf("Max Thruster Signal: %d\n", maxThrusterSignal)
}
