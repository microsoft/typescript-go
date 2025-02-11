package execute

import (
	"fmt"
	"time"

	"github.com/microsoft/typescript-go/internal/tsoptions"
)

type watcher struct {
	sys                System
	config 		   *tsoptions.ParsedCommandLine
	reportDiagnostic diagnosticReporter
}

func createWatcher(sys System, configParseResult *tsoptions.ParsedCommandLine, reportDiagnostic diagnosticReporter) *watcher {
	return &watcher{
		sys: sys,
		config: configParseResult,
		reportDiagnostic: reportDiagnostic,
		// reportWatchStatus: createWatchStatusReporter(sys, configParseResult.CompilerOptions().Pretty),
	}
}

func (w *watcher) Start() ExitStatus {
	watchInterval := 100 * time.Millisecond
	for {
		// todo: for non interval based watch, wait for file change event here
		fmt.Fprint(w.sys.Writer(), "build starting at ", time.Now(), "\n")
		w.compileAndEmit()
		fmt.Fprint(w.sys.Writer(), "build finished ", time.Now(), "\n")
		time.Sleep(watchInterval)
	}
}

func (w *watcher) compileAndEmit() {
	performCompilation(w.sys, nil, w.config, w.reportDiagnostic)
}
