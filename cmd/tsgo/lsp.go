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

func runLSP(args []string) (exitCode int) {
	flags := flag.NewFlagSet("lsp", flag.ContinueOnError)
	stdio := flags.Bool("stdio", false, "use stdio for communication")
	pipe := flags.String("pipe", "", "use named pipe for communication")
	_ = pipe
	socket := flags.String("socket", "", "use socket for communication")
	_ = socket
	if err := flags.Parse(args); err != nil {
		return 2
	}

	if !*stdio {
		errStdioSupport := errors.New("only stdio supported")
		if _, err := fmt.Fprintln(os.Stderr); err != nil {
			fmt.Printf("stderr unavailable, exiting with error %v\n", errStdioSupport)
		}
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
