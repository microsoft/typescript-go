package main

import "github.com/microsoft/typescript-go/internal/execute"

func watchTime(sys execute.System, commandLineArgs []string) {
	execute.WatchTest(sys, commandLineArgs)
}
