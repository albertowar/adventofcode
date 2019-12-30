package main

import (
	"fmt"
)

func extractDigits(number int) []int {
	digits := make([]int, 0)

	for ok := true; ok; ok = number > 0 {
		digit := number % 10
		number = number / 10

		digits = append([]int{digit}, digits...)
	}

	return digits
}

func hasTwoAdjacentDigits(digits []int) bool {
	hasTwoAdjacentDigits := false

	for i := 0; i < len(digits)-1; i++ {
		hasTwoAdjacentDigits = hasTwoAdjacentDigits || digits[i] == digits[i+1]
	}

	return hasTwoAdjacentDigits
}

func digitsNeverDecrease(digits []int) bool {
	digitsNeverDecrease := true

	for i := 0; i < len(digits)-1; i++ {
		digitsNeverDecrease = digitsNeverDecrease && digits[i] <= digits[i+1]
	}

	return digitsNeverDecrease
}

func main() {
	potentialPasswords := 0

	for i := 236491; i <= 713787; i++ {
		digits := extractDigits(i)

		if hasTwoAdjacentDigits(digits) && digitsNeverDecrease(digits) {
			fmt.Printf("%d is a valid password\n", i)
			potentialPasswords++
		}
	}

	fmt.Printf("There are %d potential passwords\n", potentialPasswords)
}
