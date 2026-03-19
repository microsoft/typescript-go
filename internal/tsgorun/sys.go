package tsgorun

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/execute/tsc"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/microsoft/typescript-go/internal/vfs/osvfs"
	"golang.org/x/term"
)

type osSys struct {
	writer             io.Writer
	fs                 vfs.FS
	defaultLibraryPath string
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

// NewSystem creates a new OS-backed system for the compiler.
// If defaultLibraryPath is empty, the bundled library path is used.
func NewSystem(defaultLibraryPath string) tsc.System {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
		os.Exit(int(tsc.ExitStatusInvalidProject_OutputsSkipped))
	}

	if defaultLibraryPath == "" {
		defaultLibraryPath = bundled.LibPath()
	}

	return &osSys{
		cwd:                tspath.NormalizePath(cwd),
		fs:                 bundled.WrapFS(osvfs.FS()),
		defaultLibraryPath: defaultLibraryPath,
		writer:             os.Stdout,
		start:              time.Now(),
	}
}
