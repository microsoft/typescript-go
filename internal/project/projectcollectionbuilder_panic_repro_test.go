package project

import (
	"context"
	"testing"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/project/logging"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
)

func TestProjectCollectionBuilder_HandlesStaleConfigRetainer(t *testing.T) {
	t.Parallel()

	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}

	files := map[string]any{
		"/src/tsconfig.json": `{"compilerOptions": {"strict": true}}`,
		"/src/index.ts":      `export const x = 1;`,
	}

	fs := bundled.WrapFS(vfstest.FromMap(files, false /*useCaseSensitiveFileNames*/))
	session := NewSession(&SessionInit{
		BackgroundCtx: context.Background(),
		Options: &SessionOptions{
			CurrentDirectory:   "/",
			DefaultLibraryPath: bundled.LibPath(),
			TypingsLocation:    "/home/src/Library/Caches/typescript",
			PositionEncoding:   lsproto.PositionEncodingKindUTF8,
			WatchEnabled:       false,
			LoggingEnabled:     false,
		},
		FS:          fs,
		Client:      noopClient{},
		Logger:      logging.NewTestLogger(),
		NpmExecutor: nil,
	})

	session.DidOpenFile(
		context.Background(),
		lsproto.DocumentUri("file:///src/index.ts"),
		1,
		files["/src/index.ts"].(string),
		lsproto.LanguageKindTypeScript,
	)

	snapshot := session.Snapshot()
	configPath := session.toPath("/src/tsconfig.json")
	entry := snapshot.ConfigFileRegistry.configs[configPath]
	if entry == nil {
		t.Fatal("expected /src/tsconfig.json config entry to exist")
	}
	if entry.retainingProjects == nil {
		entry.retainingProjects = make(map[tspath.Path]struct{})
	}
	entry.retainingProjects[session.toPath("/stale/tsconfig.json")] = struct{}{}

	session.DidChangeWatchedFiles(context.Background(), []*lsproto.FileEvent{
		{
			Uri:  lsproto.DocumentUri("file:///src/tsconfig.json"),
			Type: lsproto.FileChangeTypeChanged,
		},
	})

	_, _ = session.GetLanguageService(context.Background(), lsproto.DocumentUri("file:///src/index.ts"))
}
