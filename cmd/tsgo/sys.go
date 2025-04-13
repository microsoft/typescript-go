package main

import (
    "fmt"
    "io"
    "os"
    "runtime"
    "time"

    "github.com/microsoft/typescript-go/internal/bundled"
    "github.com/microsoft/typescript-go/internal/core"
    "github.com/microsoft/typescript-go/internal/execute"
    "github.com/microsoft/typescript-go/internal/tspath"
    "github.com/microsoft/typescript-go/internal/vfs"
    "github.com/microsoft/typescript-go/internal/vfs/osvfs"
)

type osSys struct {
    writer             io.Writer
    fs                 vfs.FS
    defaultLibraryPath string
    newLine            string
    cwd                string
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

func (s *osSys) EndWrite() {}

func newSystem() *osSys {
    cwd, err := os.Getwd()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
        os.Exit(int(execute.ExitStatusInvalidProject_OutputsSkipped))
    }
    return &osSys{
        cwd:                tspath.NormalizePath(cwd),
        fs:                 bundled.WrapFS(osvfs.FS()),
        defaultLibraryPath: bundled.LibPath(),
        writer:             os.Stdout,
        newLine:            map[string]string{"windows": "\r\n"}[runtime.GOOS],
    }
}
