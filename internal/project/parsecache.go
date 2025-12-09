package project

import (
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
)

type parseCacheKey struct {
	ast.SourceFileParseOptions
	scriptKind core.ScriptKind
}

func newParseCacheKey(
	options ast.SourceFileParseOptions,
	scriptKind core.ScriptKind,
) parseCacheKey {
	return parseCacheKey{
		SourceFileParseOptions: options,
		scriptKind:             scriptKind,
	}
}

type ParseCache = VersionedCache[parseCacheKey, *ast.SourceFile, FileHandle]
