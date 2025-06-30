package execute

import (
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/incremental"
	"github.com/microsoft/typescript-go/internal/tsoptions"
)

func CommandLineTest(sys System, commandLineArgs []string) (ExitStatus, *tsoptions.ParsedCommandLine, *watcher) {
	return commandLineWorker(sys, commandLineArgs)
}

func StartForTest(w *watcher) {
	// this function should perform any initializations before w.doCycle() in `start(watcher)`
	w.initialize()
}

func RunWatchCycle(w *watcher) {
	// this function should perform the same stuff as w.doCycle() without printing time-related output
	if w.hasErrorsInTsConfig() {
		// these are unrecoverable errors--report them and do not build
		return
	}
	// todo: updateProgram()
	w.program = incremental.NewProgram(compiler.NewProgram(compiler.ProgramOptions{
		Config:           w.options,
		Host:             w.host,
		JSDocParsingMode: ast.JSDocParsingModeParseForTypeErrors,
	}), w.program)
	if w.hasBeenModified(w.program.GetProgram()) {
		w.compileAndEmit()
	}
}
