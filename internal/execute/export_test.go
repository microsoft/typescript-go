package execute

import (
	"github.com/microsoft/typescript-go/internal/tsoptions"
)

func CommandLineTest(sys System, commandLineArgs []string) (ExitStatus, *tsoptions.ParsedCommandLine, *watcher) {
	return commandLineWorker(sys, commandLineArgs)
}

func StartForTest(w *watcher) {
	// this function should perform any initializations before w.doCycle() in `start(watcher)`
	w.initialize()
	w.doCycle()
}

func RunWatchCycle(w *watcher) {
	w.doCycle()
}
