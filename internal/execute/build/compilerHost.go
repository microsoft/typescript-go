package build

import (
	"errors"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/contentmapper"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
)

type compilerHost struct {
	host                 *host
	trace                func(msg *diagnostics.Message, args ...any)
	contentMapperProject contentmapper.Project
}

var _ compiler.CompilerHost = (*compilerHost)(nil)

func (h *compilerHost) FS() vfs.FS {
	return h.host.FS()
}

func (h *compilerHost) DefaultLibraryPath() string {
	return h.host.DefaultLibraryPath()
}

func (h *compilerHost) GetCurrentDirectory() string {
	return h.host.GetCurrentDirectory()
}

func (h *compilerHost) Trace(msg *diagnostics.Message, args ...any) {
	h.trace(msg, args...)
}

func (h *compilerHost) GetSourceFile(opts ast.SourceFileParseOptions) *ast.SourceFile {
	return h.host.GetSourceFile(opts)
}

func (h *compilerHost) GetContentMappedSourceFile(parseOptions ast.SourceFileParseOptions, mapper *contentmapper.Mapper, options *core.CompilerOptions) (*ast.SourceFile, error) {
	if h.contentMapperProject == nil {
		return nil, errors.New("content mapper project is unavailable")
	}
	content, ok := h.FS().ReadFile(parseOptions.FileName)
	if !ok {
		return nil, nil
	}
	return contentmapper.TransformAndParse(parseOptions, content, mapper, options, h.contentMapperProject)
}

func (h *compilerHost) ContentMapperProject() contentmapper.Project {
	return h.contentMapperProject
}

func (h *compilerHost) GetResolvedProjectReference(fileName string, path tspath.Path) *tsoptions.ParsedCommandLine {
	return h.host.GetResolvedProjectReference(fileName, path)
}
