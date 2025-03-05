package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/microsoft/typescript-go/internal/api"
	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/lsp"
	"github.com/microsoft/typescript-go/internal/vfs/osvfs"
)

func runAPI(args []string) int {
	flag := flag.NewFlagSet("api", flag.ContinueOnError)
	useLSP := flag.Bool("lsp", false, "use API over LSP server")
	cwd := flag.String("cwd", core.Must(os.Getwd()), "current working directory")
	if err := flag.Parse(args); err != nil {
		return 2
	}

	fs := bundled.WrapFS(osvfs.FS())
	defaultLibraryPath := bundled.LibPath()

	if *useLSP {
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
	} else {
		s := api.NewServer(&api.ServerOptions{
			In:                 os.Stdin,
			Out:                os.Stdout,
			Err:                os.Stderr,
			Cwd:                *cwd,
			NewLine:            "\n",
			DefaultLibraryPath: defaultLibraryPath,
		})

		if err := s.Run(); err != nil && !errors.Is(err, io.EOF) {
			fmt.Println(err)
			return 1
		}
		return 0
	}
}
