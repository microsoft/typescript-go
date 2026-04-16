//go:build tinygo

package main

import "fmt"

func runLSP(args []string) int {
	fmt.Println("LSP mode is not supported in this build")
	return 1
}
