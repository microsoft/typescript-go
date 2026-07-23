package project

import (
	"encoding/binary"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/parser"
	"github.com/zeebo/xxh3"
)

type ParseCacheKey struct {
	ast.SourceFileParseOptions
	ScriptKind core.ScriptKind
	Hash       xxh3.Uint128
}

func NewParseCacheKey(
	options ast.SourceFileParseOptions,
	hash xxh3.Uint128,
	scriptKind core.ScriptKind,
) ParseCacheKey {
	return ParseCacheKey{
		SourceFileParseOptions: options,
		Hash:                   hash,
		ScriptKind:             scriptKind,
	}
}

// contentMappedParseCacheKey builds the parse-cache key for a content-mapped file. The transform output
// (text and script kind) is fully determined by the file's raw content and the mapper's transform
// identity, so the key folds both into the hash and uses a fixed placeholder script kind. That lets the
// key be reconstructed from the produced source file (see parseCacheKeyForFile) without knowing the
// output script kind up front.
func contentMappedParseCacheKey(options ast.SourceFileParseOptions, rawHash, transformIdentity xxh3.Uint128) ParseCacheKey {
	var buf [32]byte
	binary.LittleEndian.PutUint64(buf[0:8], rawHash.Hi)
	binary.LittleEndian.PutUint64(buf[8:16], rawHash.Lo)
	binary.LittleEndian.PutUint64(buf[16:24], transformIdentity.Hi)
	binary.LittleEndian.PutUint64(buf[24:32], transformIdentity.Lo)
	return NewParseCacheKey(options, xxh3.Hash128(buf[:]), core.ScriptKindUnknown)
}

// parseCacheKeyForFile reconstructs the parse-cache key for a source file held by a program. Content-
// mapped files store a composite hash (raw content + transform identity) on their Hash field and are
// keyed by a placeholder script kind, matching contentMapperParseCacheKey; all other files key on their
// own hash and script kind.
func parseCacheKeyForFile(file *ast.SourceFile) ParseCacheKey {
	scriptKind := file.ScriptKind
	if file.ContentMapper() != "" {
		scriptKind = core.ScriptKindUnknown
	}
	return NewParseCacheKey(file.ParseOptions(), file.Hash, scriptKind)
}

// parseCacheKeyForDuplicate reconstructs the parse-cache key for a deduplicated source file, mirroring
// parseCacheKeyForFile.
func parseCacheKeyForDuplicate(file *compiler.DuplicateSourceFile) ParseCacheKey {
	scriptKind := file.ScriptKind
	if file.ContentMapper != "" {
		scriptKind = core.ScriptKindUnknown
	}
	return NewParseCacheKey(file.ParseOptions, file.Hash, scriptKind)
}

type ParseCache = RefCountCache[ParseCacheKey, *ast.SourceFile, FileHandle]

func NewParseCache(options RefCountCacheOptions) *ParseCache {
	return NewRefCountCache(
		options,
		func(key ParseCacheKey, fh FileHandle) *ast.SourceFile {
			file := parser.ParseSourceFile(key.SourceFileParseOptions, fh.Content(), key.ScriptKind)
			file.Hash = fh.Hash()
			return file
		},
	)
}
