package main

import (
	"fmt"
	"io"
	"os"

	"github.com/microsoft/typescript-go/internal/diagnosticwriter"
	"github.com/microsoft/typescript-go/internal/execute"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
)

type osSys struct {
	writer     io.Writer
	formatOpts *diagnosticwriter.FormattingOptions
	fs         vfs.FS
	cwd        string
}

func (s *osSys) FS() vfs.FS {
	return s.fs
}

func (s *osSys) GetCurrentDirectory() string {
	return s.cwd
}

func (s *osSys) GetFormatOpts() *diagnosticwriter.FormattingOptions {
	return s.formatOpts
}

func (s *osSys) Writer() io.Writer {
	return s.writer
}

func (s *osSys) EndWrite() {
	// do nothing, this is needed in the interface for testing
	// todo: revisit if improving tsc/build/watch unittest baselines
}

func NewSystem() *osSys {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
		os.Exit(int(execute.ExitStatusInvalidProject_OutputsSkipped))
	}
	cwd = tspath.NormalizePath(cwd)
	fs := vfs.FromOS()
	return &osSys{
		cwd:    cwd,
		fs:     fs,
		writer: os.Stdout,
		formatOpts: &diagnosticwriter.FormattingOptions{
			NewLine: "\n",
			ComparePathsOptions: tspath.ComparePathsOptions{
				CurrentDirectory:          cwd,
				UseCaseSensitiveFileNames: fs.UseCaseSensitiveFileNames(),
			},
		},
	}
}
