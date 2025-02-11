package execute

import "github.com/microsoft/typescript-go/internal/tsoptions"

type Watcher interface {
	Sys() System
	CompileAndEmit()
}

type watcher struct {
	performCompilation func()
	sys                System
}

func createWatcher(sys System, configParseResult *tsoptions.ParsedCommandLine, reportDiagnostic diagnosticReporter) Watcher {
	return &watcher{
		performCompilation: func() {
			performCompilation(sys, nil, configParseResult, reportDiagnostic)
		},
		sys: sys,
		// reportWatchStatus: createWatchStatusReporter(sys, configParseResult.CompilerOptions().Pretty),
	}
}

func (w *watcher) CompileAndEmit() {
	w.performCompilation()
}

func (w *watcher) Sys() System {
	return w.sys
}
