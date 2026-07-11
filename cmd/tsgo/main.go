package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/execute"
	"github.com/microsoft/typescript-go/internal/execute/tsc"
)

func main() {
	os.Exit(runMain())
}

func runMain() int {
	core.ApplyDebugStackLimit()
	args := os.Args[1:]
	if len(args) > 0 {
		switch args[0] {
		case "--lsp":
			return runLSP(args[1:])
		case "--api":
			return runAPI(args[1:])
		}
	}

	// Notify on our own channel rather than using signal.NotifyContext: we need the
	// actual signal so a canceled run can exit with the conventional 128+signum code,
	// matching the JS tsc (which installs no handler and lets node's default fire).
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigCh)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var receivedSignal os.Signal
	go func() {
		select {
		case receivedSignal = <-sigCh:
			cancel()
		case <-ctx.Done():
		}
	}()

	result := execute.CommandLine(ctx, newSystem(), args, nil)

	if result.Status == tsc.ExitStatusCanceled && receivedSignal != nil {
		// A signal interrupted the run. Restore the default disposition and re-raise so
		// the process terminates via the signal itself: this yields the conventional
		// exit code (130 for SIGINT, 143 for SIGTERM) and lets the runtime reset the
		// terminal, exactly as an unhandled signal would.
		signal.Reset(receivedSignal)
		if sig, ok := receivedSignal.(syscall.Signal); ok {
			_ = syscall.Kill(os.Getpid(), sig)
			// Block until the re-raised signal is delivered; do not fall through to a
			// normal return, which would exit 6 and mask the signal.
			select {}
		}
	}
	return int(result.Status)
}
