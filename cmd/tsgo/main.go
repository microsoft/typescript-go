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

	// Not signal.NotifyContext: we need to know which signal fired so we can exit by
	// re-raising it, the way the JS tsc terminates under node's default handler.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigCh)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Written before cancel(), so a canceled CommandLine always finds the signal here.
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
		// Terminate via the signal itself, for the conventional exit code (130 for
		// SIGINT, 143 for SIGTERM) and the terminal reset an unhandled signal produces.
		select {
		case sig := <-canceledBy:
			// Does not return if the signal is re-delivered; otherwise (e.g. Windows)
			// fall through to the same exit code numerically.
			reRaiseSignal(sig)
			if signo := signalNumber(sig); signo != 0 {
				return 128 + signo
			}
		default:
		}
	}
	return int(result.Status)
}

// signalNumber returns the platform signal number for sig, or 0 if it has none.
func signalNumber(sig os.Signal) int {
	if s, ok := sig.(syscall.Signal); ok {
		return int(s)
	}
	return 0
}
