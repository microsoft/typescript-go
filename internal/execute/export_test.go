package execute

import (
	"github.com/microsoft/typescript-go/internal/compiler"
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

func RunWatchCycle(w *watcher) {
	// this function should run watchCycle() without printing time-related output
	if errorsInTsConfig(w) {
		// these are unrecoverable errors--report them and do not build
		return
	}
	// todo: updateProgram()
	w.program = compiler.NewProgramFromParsedCommandLine(w.options, w.host)
	if hasBeenModified(w, w.program) {
		w.compileAndEmit()
	}
}
