package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/lsp"
	"github.com/microsoft/typescript-go/internal/vfs/osvfs"
)

func runLSP(args []string) int {
	lspFlags := flag.NewFlagSet("lsp", flag.ContinueOnError)
	stdio := lspFlags.Bool("stdio", false, "use stdio for communication")
	pipe := lspFlags.String("pipe", "", "use named pipe for communication")
	_ = pipe
	socket := lspFlags.String("socket", "", "use socket for communication")
	_ = socket
	if err := lspFlags.Parse(args); err != nil {
		return 2
	}

	if !*stdio {
		fmt.Fprintln(os.Stderr, "only stdio is supported")
		return 1
	}

	fs := bundled.WrapFS(osvfs.FS())
	defaultLibraryPath := bundled.LibPath()

	s := lsp.NewServer(&lsp.ServerOptions{
		In:                 os.Stdin,
		Out:                os.Stdout,
		Err:                os.Stderr,
		Cwd:                core.Must(os.Getwd()),
		FS:                 fs,
		DefaultLibraryPath: defaultLibraryPath,
	})

	if err := s.Run(); err != nil && !errors.Is(err, io.EOF) {
		return 1
	}
	return 0
}
