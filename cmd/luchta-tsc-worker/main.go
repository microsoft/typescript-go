package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/microsoft/typescript-go/internal/luchta"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	if err := luchta.Serve(ctx, os.Stdin, os.Stdout, os.Stderr); err != nil {
		os.Exit(1)
	}
}
