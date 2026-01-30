package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/microsoft/typescript-go/internal/api"
	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/core"
)

func runAPI(args []string) int {
	flag := flag.NewFlagSet("api", flag.ContinueOnError)
	cwd := flag.String("cwd", core.Must(os.Getwd()), "current working directory")
	callbacks := flag.String("callbacks", "", "comma-separated list of FS callbacks to enable (readFile,fileExists,directoryExists,getAccessibleEntries,realpath)")
	if err := flag.Parse(args); err != nil {
		return 2
	}

	defaultLibraryPath := bundled.LibPath()

	// Parse callbacks list
	var callbacksList []string
	if *callbacks != "" {
		callbacksList = strings.Split(*callbacks, ",")
	}

	s := api.NewStdioServer(&api.StdioServerOptions{
		In:                 os.Stdin,
		Out:                os.Stdout,
		Err:                os.Stderr,
		Cwd:                *cwd,
		DefaultLibraryPath: defaultLibraryPath,
		Callbacks:          callbacksList,
	})
	defer s.Close()

	if err := s.Run(context.Background()); err != nil && !errors.Is(err, io.EOF) {
		fmt.Println(err)
		return 1
	}
	return 0
}
