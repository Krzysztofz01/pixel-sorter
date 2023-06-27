package main

import (
	"os"
	"strings"
	"testing"

	"github.com/Krzysztofz01/pixel-sorter/cmd"
)

func BenchmarkPixelSorterCliFromEnvArgs(b *testing.B) {
	args := os.Getenv("PIXEL_SORTER_CLI_ARGS")

	cmd.Execute(strings.Split(args, " "))
}
