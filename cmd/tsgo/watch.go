package main

import (
	"fmt"
	"time"

	"github.com/microsoft/typescript-go/internal/execute"
)

func start(watcher execute.Watcher) int {
	watchInterval := 100 * time.Millisecond
	for {
		fmt.Fprint(watcher.Sys().Writer(), "build starting at ", time.Now(), "\n")
		watcher.CompileAndEmit()
		fmt.Fprint(watcher.Sys().Writer(), "build finished ", time.Now(), "\n")
		time.Sleep(watchInterval)
	}
}
