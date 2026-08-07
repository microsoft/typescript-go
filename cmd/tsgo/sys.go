package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/execute/tsc"
	"github.com/microsoft/typescript-go/internal/pnp"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/microsoft/typescript-go/internal/vfs/osvfs"
	"github.com/microsoft/typescript-go/internal/vfs/pnpvfs"
	"golang.org/x/term"
)

type osSys struct {
	writer             io.Writer
	fs                 vfs.FS
	defaultLibraryPath string
	cwd                string
	start              time.Time
	pnpApi             *pnp.PnpApi
}

func (s *osSys) SinceStart() time.Duration {
	return time.Since(s.start)
}

func (s *osSys) Now() time.Time {
	return time.Now()
}

func (s *osSys) FS() vfs.FS {
	return s.fs
}

func (s *osSys) DefaultLibraryPath() string {
	return s.defaultLibraryPath
}

func (s *osSys) GetCurrentDirectory() string {
	return s.cwd
}

func (s *osSys) Writer() io.Writer {
	return s.writer
}

func (s *osSys) WriteOutputIsTTY() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

func (s *osSys) GetWidthOfTerminal() int {
	width, _, _ := term.GetSize(int(os.Stdout.Fd()))
	return width
}

func (s *osSys) GetEnvironmentVariable(name string) string {
	return os.Getenv(name)
}

func (s *osSys) PnpApi() *pnp.PnpApi {
	return s.pnpApi
}

func newSystem() *osSys {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
		os.Exit(int(tsc.ExitStatusInvalidProject_OutputsSkipped))
	}

	var fs vfs.FS = osvfs.FS()

	pnpApi := pnp.InitPnpApi(fs, tspath.NormalizePath(cwd))
	if pnpApi != nil {
		fs = pnpvfs.From(fs)
	}

	return &osSys{
		cwd:                tspath.NormalizePath(cwd),
		fs:                 bundled.WrapFS(fs),
		defaultLibraryPath: bundled.LibPath(),
		writer:             os.Stdout,
		start:              time.Now(),
		pnpApi:             pnpApi,
	}
}
