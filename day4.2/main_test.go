package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDigitsNeverDecrease_ReturnsTrue_WhenNoDigits(t *testing.T) {
	result := DigitsNeverDecrease([]int{})

	assert.True(t, result)
}
