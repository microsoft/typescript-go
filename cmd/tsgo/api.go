package main

import (
	"errors"
	"flag"
	"io"
	"os"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/lsp"
	"github.com/microsoft/typescript-go/internal/vfs/osvfs"
)

func runAPI(args []string) int {
	flag := flag.NewFlagSet("api", flag.ContinueOnError)
	cwd := flag.String("cwd", core.Must(os.Getwd()), "current working directory")
	if err := flag.Parse(args); err != nil {
		return 2
	}

	fs := bundled.WrapFS(osvfs.FS())
	defaultLibraryPath := bundled.LibPath()

	s := lsp.NewServer(&lsp.ServerOptions{
		In:                 os.Stdin,
		Out:                os.Stdout,
		Err:                os.Stderr,
		API:                true,
		Cwd:                *cwd,
		FS:                 fs,
		DefaultLibraryPath: defaultLibraryPath,
	})

	if err := s.Run(); err != nil && !errors.Is(err, io.EOF) {
		return 1
	}
	return 0
}
