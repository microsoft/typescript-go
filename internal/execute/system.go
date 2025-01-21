package execute

import (
	"fmt"
	"os"

	ts "github.com/microsoft/typescript-go/internal/compiler"
	dw "github.com/microsoft/typescript-go/internal/diagnosticwriter"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
)

// todo: implement system? figure out where it is implemented?
// is it the same as fs?
// from looking at go and main, it looks like they are different, host extends vfs sorta.
// compilerhost in compiler/host.go does not have a sys.exit or sys.write
type System interface {
	Exit(status ExitStatus) ExitStatus
	Write(p []byte) (n int, err error)
	EndWrite()
	FS() vfs.FS
	Host() ts.CompilerHost
	GetFormatOpts() *dw.FormattingOptions // todo: should this be part of Host?
}

type osSys struct {
	exit       func(int) // os.Exit
	write      func([]byte) (int, error)
	formatOpts dw.FormattingOptions
	host       ts.CompilerHost
}

func (s *osSys) FS() vfs.FS {
	return s.Host().FS()
}

func (s *osSys) Host() ts.CompilerHost {
	return s.host
}

func (s *osSys) Exit(e ExitStatus) ExitStatus {
	s.exit(int(e))
	return e
}

func (s *osSys) Write(p []byte) (n int, err error) {
	return s.write(p)
}

func (s *osSys) EndWrite() {
	// do nothing, this is needed in the interface for testing
	// todo: revisit if improving tsc/build/watch unittest baselines
}

func (s *osSys) GetFormatOpts() *dw.FormattingOptions {
	return &s.formatOpts
}

func NewSystem() *osSys {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
		os.Exit(int(ExitStatusInvalidProject_OutputsSkipped))
	}
	newHost := ts.NewCompilerHost(nil, cwd, vfs.FromOS())
	return &osSys{
		host:  newHost,
		exit:  os.Exit,
		write: os.Stdout.Write,
		formatOpts: dw.FormattingOptions{
			NewLine: "\n",
			ComparePathsOptions: tspath.ComparePathsOptions{
				CurrentDirectory:          cwd,
				UseCaseSensitiveFileNames: newHost.FS().UseCaseSensitiveFileNames(),
			},
		},
	}
}

type ExitStatus int

const (
	ExitStatusSuccess                              ExitStatus = iota
	ExitStatusDiagnosticsPresent_OutputsSkipped    ExitStatus = 1
	ExitStatusDiagnosticsPresent_OutputsGenerated  ExitStatus = 2
	ExitStatusInvalidProject_OutputsSkipped        ExitStatus = 3
	ExitStatusProjectReferenceCycle_OutputsSkipped ExitStatus = 4
)
