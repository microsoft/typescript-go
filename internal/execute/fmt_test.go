package execute

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/repo"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/microsoft/typescript-go/internal/vfs/osvfs"
)

func BenchmarkFormat(b *testing.B) {
	checkerPath := tspath.CombinePaths(repo.TypeScriptSubmodulePath, "src", "compiler", "checker.ts")

	tmp := b.TempDir()
	sys := newSystem()

	var count int

	b.ReportAllocs()
	time.Sleep(5 * time.Second)
	for b.Loop() {
		out := tspath.CombinePaths(tmp, fmt.Sprintf("out%d.ts", count))
		code := fmtMain(sys, checkerPath, out)
		if code != 0 {
			b.Fatalf("Unexpected exit code: %d", code)
		}
	}
}

type osSys struct {
	writer             io.Writer
	fs                 vfs.FS
	defaultLibraryPath string
	newLine            string
	cwd                string
	start              time.Time
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

func (s *osSys) NewLine() string {
	return s.newLine
}

func (s *osSys) Writer() io.Writer {
	return s.writer
}

func (s *osSys) EndWrite() {
	// do nothing, this is needed in the interface for testing
	// todo: revisit if improving tsc/build/watch unittest baselines
}

func newSystem() *osSys {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
		os.Exit(int(ExitStatusInvalidProject_OutputsSkipped))
	}

	return &osSys{
		cwd:                tspath.NormalizePath(cwd),
		fs:                 bundled.WrapFS(osvfs.FS()),
		defaultLibraryPath: bundled.LibPath(),
		writer:             os.Stdout,
		newLine:            core.IfElse(runtime.GOOS == "windows", "\r\n", "\n"),
		start:              time.Now(),
	}
}
