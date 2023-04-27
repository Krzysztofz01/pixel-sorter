package main

import (
	"os"

	"github.com/Krzysztofz01/pixel-sorter/cmd"
)

func main() {
	cmd.Execute(os.Args[1:])
}
