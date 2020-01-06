package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	intcode "github.com/albertowar/adventofcode/intcode"
)

func sendInput(inputs []int, inputChannel chan int) {
	for _, input := range inputs {
		inputChannel <- input
	}
}

func allDifferent(a int, b int, c int, d int, e int) bool {
	return a != b && a != c && a != d && a != e && b != c && b != d && b != e && c != d && c != e && d != e
}

func sleepAndRunIntcodeProgram(intcodes []int, inputChannel chan int, outputChannel chan int, wg *sync.WaitGroup) {
	time.Sleep(500 * time.Millisecond)
	intcode.RunIntcodeProgram(intcodes, inputChannel, outputChannel, wg)
}

func main() {
	intcodes, err := intcode.ReadIntCodes("input.txt")

	if err != nil {
		log.Fatal(err)
	}

	minPhase := 5
	maxPhase := 9

	maxThrusterSignal := 0

	inputAmpA, inputAmpB, inputAmpC, inputAmpD, inputAmpE := make(chan int), make(chan int), make(chan int), make(chan int), make(chan int)

	for phaseA := minPhase; phaseA <= maxPhase; phaseA++ {
		for phaseB := minPhase; phaseB <= maxPhase; phaseB++ {
			for phaseC := minPhase; phaseC <= maxPhase; phaseC++ {
				for phaseD := minPhase; phaseD <= maxPhase; phaseD++ {
					for phaseE := minPhase; phaseE <= maxPhase; phaseE++ {
						if allDifferent(phaseA, phaseB, phaseC, phaseD, phaseE) {
							intcodesAmpA, intcodesAmpB, intcodesAmpC, intcodesAmpD, intcodesAmpE := make([]int, len(intcodes)), make([]int, len(intcodes)), make([]int, len(intcodes)), make([]int, len(intcodes)), make([]int, len(intcodes))
							copy(intcodesAmpA, intcodes)
							copy(intcodesAmpB, intcodes)
							copy(intcodesAmpC, intcodes)
							copy(intcodesAmpD, intcodes)
							copy(intcodesAmpE, intcodes)

							go sendInput([]int{phaseA, 0}, inputAmpA)
							go sendInput([]int{phaseB}, inputAmpB)
							go sendInput([]int{phaseC}, inputAmpC)
							go sendInput([]int{phaseD}, inputAmpD)
							go sendInput([]int{phaseE}, inputAmpE)

							var wg sync.WaitGroup
							wg.Add(4)
							go sleepAndRunIntcodeProgram(intcodesAmpA, inputAmpA, inputAmpB, &wg)
							go sleepAndRunIntcodeProgram(intcodesAmpB, inputAmpB, inputAmpC, &wg)
							go sleepAndRunIntcodeProgram(intcodesAmpC, inputAmpC, inputAmpD, &wg)
							go sleepAndRunIntcodeProgram(intcodesAmpD, inputAmpD, inputAmpE, &wg)
							go sleepAndRunIntcodeProgram(intcodesAmpE, inputAmpE, inputAmpA, nil)
							wg.Wait()

							output := <-inputAmpA

							if output > maxThrusterSignal {
								fmt.Printf("Found new MaxThrusterSignal with value %d. (%d,%d,%d,%d,%d)\n", output, phaseA, phaseB, phaseC, phaseD, phaseE)
								maxThrusterSignal = output
							}
						}
					}
				}
			}
		}
	}

	fmt.Printf("Max Thruster Signal: %d\n", maxThrusterSignal)
}
