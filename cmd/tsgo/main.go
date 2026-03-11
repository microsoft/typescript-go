package main

import (
	"os"
	"strings"

	"github.com/microsoft/typescript-go/internal/execute"
)

func init() {
	const tracebackAncestors = "tracebackancestors"
	const tracebackAncestorsSetting = tracebackAncestors + "=5"
	godebug := os.Getenv("GODEBUG")
	if godebug == "" {
		os.Setenv("GODEBUG", tracebackAncestorsSetting)
	} else if !strings.Contains(godebug, tracebackAncestors+"=") {
		os.Setenv("GODEBUG", godebug+","+tracebackAncestorsSetting)
	}
}

func main() {
	os.Exit(runMain())
}

func runMain() int {
	args := os.Args[1:]
	if len(args) > 0 {
		switch args[0] {
		case "--lsp":
			return runLSP(args[1:])
		case "--api":
			return runAPI(args[1:])
		}
	}
	result := execute.CommandLine(newSystem(), args, nil)
	return int(result.Status)
}
