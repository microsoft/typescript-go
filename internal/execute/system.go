package execute

import (
	"io"
	"time"

	"github.com/microsoft/typescript-go/internal/vfs"
)

type System interface {
	Writer() io.Writer
	FS() vfs.FS
	DefaultLibraryPath() string
	GetCurrentDirectory() string

	Now() time.Time
	SinceStart() time.Duration
}

type ExitStatus int

const (
	ExitStatusSuccess                              ExitStatus = 0
	ExitStatusDiagnosticsPresent_OutputsGenerated  ExitStatus = 1
	ExitStatusDiagnosticsPresent_OutputsSkipped    ExitStatus = 2
	ExitStatusInvalidProject_OutputsSkipped        ExitStatus = 3
	ExitStatusProjectReferenceCycle_OutputsSkipped ExitStatus = 4
	ExitStatusNotImplemented                       ExitStatus = 5
)
