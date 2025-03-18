package emittestutil

import (
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/binder"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/printer"
	"github.com/microsoft/typescript-go/internal/testutil/parsetestutil"
	"github.com/microsoft/typescript-go/internal/tspath"
	"gotest.tools/v3/assert"
)

// Checks that pretty-printing the given file matches the expected output.
func CheckEmit(t *testing.T, emitContext *printer.EmitContext, file *ast.SourceFile, expected string) {
	t.Helper()
	printer := printer.NewPrinter(
		printer.PrinterOptions{
			NewLine: core.NewLineKindLF,
		},
		printer.PrintHandlers{},
		emitContext,
	)
	text := printer.EmitSourceFile(file)
	actual := strings.TrimSuffix(text, "\n")
	assert.Equal(t, expected, actual)
	file2 := parsetestutil.ParseTypeScript(text, file.LanguageVariant == core.LanguageVariantJSX)
	parsetestutil.CheckDiagnosticsMessage(t, file2, "error on reparse: ")
}

type fakeProgram struct {
	singleThreaded              bool
	compilerOptions             *core.CompilerOptions
	files                       []*ast.SourceFile
	getEmitModuleFormatOfFile   func(sourceFile *ast.SourceFile) core.ModuleKind
	getImpliedNodeFormatForEmit func(sourceFile *ast.SourceFile) core.ModuleKind
	getResolvedModule           func(currentSourceFile *ast.SourceFile, moduleReference string) *ast.SourceFile
}

func (p *fakeProgram) Options() *core.CompilerOptions {
	return p.compilerOptions
}

func (p *fakeProgram) SourceFiles() []*ast.SourceFile {
	return p.files
}

func (p *fakeProgram) BindSourceFiles() {
	wg := core.NewWorkGroup(p.singleThreaded)
	for _, file := range p.files {
		if !file.IsBound() {
			wg.Queue(func() {
				binder.BindSourceFile(file, p.compilerOptions)
			})
		}
	}
	wg.RunAndWait()
}

func (p *fakeProgram) GetEmitModuleFormatOfFile(sourceFile *ast.SourceFile) core.ModuleKind {
	return p.getEmitModuleFormatOfFile(sourceFile)
}

func (p *fakeProgram) GetImpliedNodeFormatForEmit(sourceFile *ast.SourceFile) core.ModuleKind {
	return p.getImpliedNodeFormatForEmit(sourceFile)
}

func (p *fakeProgram) GetResolvedModule(currentSourceFile *ast.SourceFile, moduleReference string) *ast.SourceFile {
	return p.getResolvedModule(currentSourceFile, moduleReference)
}

func (p *fakeProgram) GetSourceFileMetaData(path tspath.Path) *ast.SourceFileMetaData {
	return nil
}

func NewFakeProgram(singleThreaded bool, compilerOptions *core.CompilerOptions, files []*ast.SourceFile, file, other *ast.SourceFile) *fakeProgram {
	return &fakeProgram{
		singleThreaded:  singleThreaded,
		compilerOptions: compilerOptions,
		files:           files,
		getEmitModuleFormatOfFile: func(sourceFile *ast.SourceFile) core.ModuleKind {
			return core.ModuleKindESNext
		},
		getImpliedNodeFormatForEmit: func(sourceFile *ast.SourceFile) core.ModuleKind {
			return core.ModuleKindESNext
		},
		getResolvedModule: func(currentSourceFile *ast.SourceFile, moduleReference string) *ast.SourceFile {
			if currentSourceFile == file && moduleReference == "other" {
				return other
			}
			return nil
		},
	}
}
