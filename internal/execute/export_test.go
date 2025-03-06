package execute

import (
	"github.com/microsoft/typescript-go/internal/tsoptions"
)

func CommandLineTest(sys System, cb cbType, commandLineArgs []string) (*tsoptions.ParsedCommandLine, ExitStatus) {
	parsedCommandLine := tsoptions.ParseCommandLine(commandLineArgs, sys)
	e, _ := executeCommandLineWorker(sys, cb, parsedCommandLine)
	return parsedCommandLine, e
}

func CommandLineTestWatch(sys System, cb cbType, commandLineArgs []string) (*tsoptions.ParsedCommandLine, *watcher) {
	parsedCommandLine := tsoptions.ParseCommandLine(commandLineArgs, sys)
	_, w := executeCommandLineWorker(sys, cb, parsedCommandLine)
	return parsedCommandLine, w
}
