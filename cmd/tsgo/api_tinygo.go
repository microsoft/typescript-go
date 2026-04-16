//go:build tinygo

package main

import "fmt"

func runAPI(args []string) int {
	fmt.Println("API mode is not supported in this build")
	return 1
}
