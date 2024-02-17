package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCIntnShouldPanicForInvalidArgument(t *testing.T) {
	assert.Panics(t, func() {
		CIntn(-2)
	})

	assert.Panics(t, func() {
		CIntn(0)
	})
}

func TestCIntnShouldReturnRandomValueInCorrectRange(t *testing.T) {
	const (
		limit      int = 25
		iterations int = 100_000_000
	)

	for i := 0; i < iterations; i += 1 {
		value := CIntn(limit)
		if value < 0 || value >= limit {
			assert.FailNow(t, "CIntn generated a value out of the expected range")
		}
	}
}
