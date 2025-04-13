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
    "github.com/microsoft/typescript-go/internal/pprof"
    "github.com/microsoft/typescript-go/internal/vfs/osvfs"
)

func runLSP(args []string) int {
    // Create a new flag set to parse command-line arguments.
    flagSet := flag.NewFlagSet("lsp", flag.ContinueOnError)
    // Define a flag to enable stdio communication.
    stdio := flagSet.Bool("stdio", false, "Use stdio for communication")
    // Define a flag to specify a directory for pprof profiles.
    pprofDir := flagSet.String("pprofDir", "", "Generate pprof CPU/memory profiles to the given directory.")

    // Parse command-line arguments and handle errors gracefully.
    if err := flagSet.Parse(args); err != nil {
        // Print parsing error and return exit code 2.
        fmt.Fprintln(os.Stderr, "Error parsing flags:", err)
        return 2
    }

    // Check if stdio communication is enabled; if not, exit with an error.
    if !*stdio {
        fmt.Fprintln(os.Stderr, "Only stdio is supported")
        return 1
    }

    // If pprofDir is specified, initialize profiling and handle cleanup using defer.
    if *pprofDir != "" {
        fmt.Fprintf(os.Stderr, "pprof profiles will be written to: %v\n", *pprofDir)
        profileSession := pprof.BeginProfiling(*pprofDir, os.Stderr)
        // Ensure that the profiling session is stopped properly at the end.
        defer profileSession.Stop()
    }

    // Wrap the file system using a helper from bundled and retrieve the default library path.
    fs := bundled.WrapFS(osvfs.FS())
    defaultLibraryPath := bundled.LibPath()

    // Create a new LSP server instance with required options.
    s := lsp.NewServer(&lsp.ServerOptions{
        In:                 os.Stdin,                  // Input stream from standard input.
        Out:                os.Stdout,                 // Output stream to standard output.
        Err:                os.Stderr,                 // Error stream to standard error.
        Cwd:                core.Must(os.Getwd()),     // Current working directory.
        FS:                 fs,                        // Wrapped file system.
        DefaultLibraryPath: defaultLibraryPath,        // Path to default TypeScript library.
    })

    // Run the LSP server and handle errors during execution.
    if err := s.Run(); err != nil {
        // Only log non-EOF errors as critical issues.
        if !errors.Is(err, io.EOF) {
            fmt.Fprintln(os.Stderr, "Server run error:", err)
            return 1
        }
    }
    return 0
}


// Flag initialization:

// Added comments to explain the purpose of each flag (stdio and pprofDir), making the code easier to understand for new developers.

// Error handling:

// Added inline comments where errors are handled, explaining the exit conditions and their significance.

// Profiling setup:

// Added comments to describe the profiling logic and why defer is used to ensure cleanup.

// LSP server setup:

// Explained the configuration options (In, Out, Err, etc.) for the LSP server.

// Error management during server runtime:

// Added a comment explaining the differentiation between EOF errors and critical errors.
