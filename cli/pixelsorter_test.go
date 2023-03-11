package main

import (
	"testing"

	"github.com/Krzysztofz01/pixel-sorter/cmd"
)

// NOTE: This test suite is used to perform profiling. Command:
// go test -cpuprofile cpu.prof -memprofile mem.prof -bench .

func TestProfilingBenchmarkSortBrightnessIntervalBrightness(t *testing.T) {
	args := []string{
		"brightness",
		"--image-file-path",
		"image-test.png",
		"--output-format",
		"jpg",
		"--direction",
		"ascending",
		"--interval-determinant",
		"brightness",
		"--interval-lower-threshold",
		"0.15",
		"--interval-upper-threshold",
		"0.85",
		"--interval-max-length",
		"0",
		"--order",
		"horizontal-vertical",
		"--angle",
		"0",
		"cycles",
		"1",
		"--verbose",
		"true",
	}

	cmd.Execute(args)
}
