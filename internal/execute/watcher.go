package execute

import (
	"time"

	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/tsoptions"
)

type watcher struct {
	sys              System
	configFileName   string
	options          *tsoptions.ParsedCommandLine
	reportDiagnostic diagnosticReporter

	host         compiler.CompilerHost
	program      *compiler.Program
	prevModified map[string]time.Time
}

func createWatcher(sys System, configParseResult *tsoptions.ParsedCommandLine, reportDiagnostic diagnosticReporter) *watcher {
	return &watcher{
		sys:              sys,
		configFileName:   configParseResult.ConfigFile.SourceFile.FileName(),
		options:          configParseResult,
		reportDiagnostic: reportDiagnostic,
		// reportWatchStatus: createWatchStatusReporter(sys, configParseResult.CompilerOptions().Pretty),
	}
}

func (w *watcher) compileAndEmit() {
	// !!! output/error reporting is currently the same as non-watch mode
	// diagnostics, emitResult, exitStatus :=
	compileAndEmit(w.sys, w.program, w.reportDiagnostic)
}
