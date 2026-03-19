package tsgorun

import (
	"os"

	"github.com/microsoft/typescript-go/internal/execute"
)

// RunMain runs the main tsgo command line, dispatching to --lsp or --api
// subcommands as needed. It returns the process exit code.
// If defaultLibraryPath is empty, the bundled library path is used.
func RunMain(args []string, defaultLibraryPath string) int {
	if len(args) > 0 {
		switch args[0] {
		case "--lsp":
			return RunLSP(args[1:])
		case "--api":
			return RunAPI(args[1:])
		}
	}
	result := execute.CommandLine(NewSystem(defaultLibraryPath), args, nil)
	return int(result.Status)
}

// Main calls RunMain with the process arguments and exits.
func Main() {
	os.Exit(RunMain(os.Args[1:], ""))
}
