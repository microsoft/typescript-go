package execute

import (
	"io"

	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/diagnosticwriter"
	"github.com/microsoft/typescript-go/internal/vfs"
)

// todo: implement system? figure out where it is implemented?
// is it the same as fs?
// from looking at go and main, it looks like they are different, host extends vfs sorta.
// compilerhost in compiler/host.go does not have a sys.exit or sys.write
type System interface {
	Exit(status ExitStatus) ExitStatus
	Writer() io.Writer
	EndWrite()
	FS() vfs.FS
	Host() compiler.CompilerHost
	GetFormatOpts() *diagnosticwriter.FormattingOptions // todo: should this be part of Host?
}

type ExitStatus int

const (
	ExitStatusSuccess                              ExitStatus = iota
	ExitStatusDiagnosticsPresent_OutputsSkipped    ExitStatus = 1
	ExitStatusDiagnosticsPresent_OutputsGenerated  ExitStatus = 2
	ExitStatusInvalidProject_OutputsSkipped        ExitStatus = 3
	ExitStatusProjectReferenceCycle_OutputsSkipped ExitStatus = 4
)
