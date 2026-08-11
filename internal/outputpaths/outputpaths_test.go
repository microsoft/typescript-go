package outputpaths_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/outputpaths"
	"github.com/microsoft/typescript-go/internal/parser"
	"gotest.tools/v3/assert"
)

type contentMapperOutputPathsHost struct{}

func (contentMapperOutputPathsHost) CommonSourceDirectory() string { return "/" }

func (contentMapperOutputPathsHost) ContentMapperExtensionRewrites() []core.ExtensionRewrite {
	return []core.ExtensionRewrite{{Source: ".mdx", Target: ".js"}}
}
func (contentMapperOutputPathsHost) GetCurrentDirectory() string     { return "/" }
func (contentMapperOutputPathsHost) UseCaseSensitiveFileNames() bool { return true }

func TestContentMappedJSOutputSuppressesSourceMapPath(t *testing.T) {
	t.Parallel()
	sourceFile := parser.ParseSourceFile(ast.SourceFileParseOptions{FileName: "/component.mdx"}, "export {};", core.ScriptKindTS)
	sourceFile.SetContentMapper("mapper@1.0.0")

	paths := outputpaths.GetOutputPathsFor(sourceFile, &core.CompilerOptions{SourceMap: core.TSTrue}, contentMapperOutputPathsHost{}, outputpaths.ForceEmitPaths{})
	assert.Equal(t, paths.JsFilePath(), "/component.mdx.js")
	assert.Equal(t, paths.SourceMapFilePath(), "")
}
