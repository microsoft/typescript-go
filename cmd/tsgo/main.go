package main

import (
	"os"
)

func main() {
	os.Exit(runMain())
}

func runMain() int {
	args := os.Args[1:]
	if len(args) > 0 {
		switch args[0] {
		case "--api":
			return runAPI(args[1:])
		}
	}
	return 1
}
