package execute

import (
	"io"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/diagnosticwriter"
	"github.com/microsoft/typescript-go/internal/vfs"
)

type System interface {
	Writer() io.Writer
	EndWrite()
	FS() vfs.FS
	Host() compiler.CompilerHost
	SetReportDiagnostics(r DiagnosticReporter)
	ReportDiagnostic(d *ast.Diagnostic)
	GetFormatOpts() *diagnosticwriter.FormattingOptions // todo: should this be part of Host?
}

type ExitStatus int

const (
	ExitStatusSuccess                              ExitStatus = iota
	ExitStatusDiagnosticsPresent_OutputsSkipped    ExitStatus = 1
	ExitStatusDiagnosticsPresent_OutputsGenerated  ExitStatus = 2
	ExitStatusInvalidProject_OutputsSkipped        ExitStatus = 3
	ExitStatusProjectReferenceCycle_OutputsSkipped ExitStatus = 4
	ExitStatusNotImplemented                       ExitStatus = 5
	ExitStatusNotImplementedWatch                  ExitStatus = 6
	ExitStatusNotImplementedIncremental            ExitStatus = 7
)
