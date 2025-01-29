package main

import (
	"fmt"
	"io"
	"os"

	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/diagnosticwriter"
	"github.com/microsoft/typescript-go/internal/execute"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
)

type osSys struct {
	writer     io.Writer
	formatOpts *diagnosticwriter.FormattingOptions
	host       compiler.CompilerHost
}

func (s *osSys) FS() vfs.FS {
	return s.Host().FS()
}

func (s *osSys) Host() compiler.CompilerHost {
	return s.host
}

func (s *osSys) Exit(e execute.ExitStatus) execute.ExitStatus {
	return e
}

func (s *osSys) Writer() io.Writer {
	return s.writer
}

func (s *osSys) EndWrite() {
	// do nothing, this is needed in the interface for testing
	// todo: revisit if improving tsc/build/watch unittest baselines
}

func (s *osSys) GetFormatOpts() *diagnosticwriter.FormattingOptions {
	return s.formatOpts
}

func NewSystem() *osSys {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
		os.Exit(int(execute.ExitStatusInvalidProject_OutputsSkipped))
	}
	newHost := compiler.NewCompilerHost(nil, tspath.NormalizePath(cwd), vfs.FromOS())
	return &osSys{
		host:       newHost,
		writer:     os.Stdout,
		formatOpts: getFormatOpts(newHost),
	}
}
