package sorter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultSorterOptionsShouldValidate(t *testing.T) {
	options := GetDefaultSorterOptions()

	valid, msg := options.AreValid()

	assert.True(t, valid)
	assert.Empty(t, msg)
}

func TestSorterOptionsShouldNotValidateInvalidIntervalDeterminantLowerThreshold(t *testing.T) {
	cases := []float64{
		-0.5,
		1.5,
	}

	for _, value := range cases {
		options := GetDefaultSorterOptions()
		options.IntervalDeterminantLowerThreshold = value
		valid, msg := options.AreValid()

		assert.False(t, valid)
		assert.NotEmpty(t, msg)
	}
}

func TestSorterOptionsShouldNotValidateInvalidIntervalDeterminantUpperThreshold(t *testing.T) {
	cases := []float64{
		-0.5,
		1.5,
	}

	for _, value := range cases {
		options := GetDefaultSorterOptions()
		options.IntervalDeterminantUpperThreshold = value
		valid, msg := options.AreValid()

		assert.False(t, valid)
		assert.NotEmpty(t, msg)
	}
}

func TestSorterOptionsShouldNotValidateInvalidIntervalDeterminantThresholds(t *testing.T) {
	options := GetDefaultSorterOptions()
	options.IntervalDeterminantLowerThreshold = 0.6
	options.IntervalDeterminantUpperThreshold = 0.4

	valid, msg := options.AreValid()

	assert.False(t, valid)
	assert.NotEmpty(t, msg)
}

func TestSorterOptionsShouldNotValidateInvalidCycles(t *testing.T) {
	cases := []int{
		-1,
		0,
	}

	for _, value := range cases {
		options := GetDefaultSorterOptions()
		options.Cycles = value

		valid, msg := options.AreValid()

		assert.False(t, valid)
		assert.NotEmpty(t, msg)
	}
}

func TestSorterOptionsShouldNotValidateInvalidScale(t *testing.T) {
	cases := []float64{
		-0.5,
		0.0,
		1.5,
	}

	for _, value := range cases {
		options := GetDefaultSorterOptions()
		options.Scale = value

		valid, msg := options.AreValid()

		assert.False(t, valid)
		assert.NotEmpty(t, msg)
	}
}

func TestSorterOptionsShouldNotValidateIntervalLength(t *testing.T) {
	options := GetDefaultSorterOptions()
	options.Scale = -1

	valid, msg := options.AreValid()

	assert.False(t, valid)
	assert.NotEmpty(t, msg)
}
