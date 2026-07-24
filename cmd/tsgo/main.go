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

	// Use signal.Notify with our own channel rather than signal.NotifyContext: the
	// latter's context can't tell us which signal fired, and we need it to exit like
	// the JS tsc, which installs no handler and lets node terminate via the signal.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigCh)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// canceledBy carries the interrupting signal (if any) to the code below. It is
	// written before cancel() and thus before CommandLine can observe cancellation.
	canceledBy := make(chan os.Signal, 1)
	go func() {
		select {
		case sig := <-sigCh:
			canceledBy <- sig
			cancel()
		case <-ctx.Done():
		}
	}()

	result := execute.CommandLine(ctx, newSystem(), args, nil)

	if result.Status == tsc.ExitStatusCanceled {
		// A signal canceled the run. Re-raise it so we terminate via the signal itself,
		// yielding the conventional exit code (130 for SIGINT, 143 for SIGTERM) and the
		// terminal reset an unhandled signal would produce. On platforms that cannot
		// re-deliver the signal (e.g. Windows), reRaiseSignal returns 0 and we exit with
		// the numeric status instead.
		select {
		case sig := <-canceledBy:
			if signo := reRaiseSignal(sig); signo != 0 {
				return 128 + signo // fallback in case the signal doesn't land promptly
			}
		default:
		}
	}
	return int(result.Status)
}
