package project

import (
	"context"

	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
)

type Client interface {
	WatchFiles(ctx context.Context, id WatcherID, watchers []*lsproto.FileSystemWatcher) error
	UnwatchFiles(ctx context.Context, id WatcherID) error
	RefreshDiagnostics(ctx context.Context) error
	PublishDiagnostics(ctx context.Context, params *lsproto.PublishDiagnosticsParams) error
}

type noopClient struct{}

func (n noopClient) WatchFiles(ctx context.Context, id WatcherID, watchers []*lsproto.FileSystemWatcher) error {
	return nil
}

func (n noopClient) UnwatchFiles(ctx context.Context, id WatcherID) error {
	return nil
}

func (n noopClient) RefreshDiagnostics(ctx context.Context) error {
	return nil
}

func (n noopClient) PublishDiagnostics(ctx context.Context, params *lsproto.PublishDiagnosticsParams) error {
	return nil
}
