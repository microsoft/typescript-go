package project

import (
	"strings"

	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
)

func isDynamicFileName(fileName string) bool {
	return strings.HasPrefix(fileName, "^")
}

type fileSystemWatcherKey struct {
	pattern string
	kind    lsproto.WatchKind
}

func toFileSystemWatcherKey(w *lsproto.FileSystemWatcher) fileSystemWatcherKey {
	if w.GlobPattern.RelativePattern != nil {
		panic("relative globs not implemented")
	}
	kind := w.Kind
	if kind == nil {
		kind = ptrTo(lsproto.WatchKindCreate | lsproto.WatchKindChange | lsproto.WatchKindDelete)
	}
	return fileSystemWatcherKey{pattern: *w.GlobPattern.Pattern, kind: *kind}
}
