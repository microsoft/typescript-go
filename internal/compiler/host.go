package compiler

import (
	"strings"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/compiler/packagejson"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/parser"
	"github.com/microsoft/typescript-go/internal/scanner"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
)

type CompilerHost interface {
	FS() vfs.FS
	DefaultLibraryPath() string
	GetCurrentDirectory() string
	NewLine() string
	Trace(msg string)
	GetSourceFile(fileName string, path tspath.Path, languageVersion core.ScriptTarget, packageJsonScope *packagejson.InfoCacheEntry) *ast.SourceFile
	GetImpliedNodeFormat(fileName string, packageJsonScope *packagejson.InfoCacheEntry) core.ResolutionMode
}

type FileInfo struct {
	Name string
	Size int64
}

var _ CompilerHost = (*compilerHost)(nil)

type compilerHost struct {
	options            *core.CompilerOptions
	currentDirectory   string
	fs                 vfs.FS
	defaultLibraryPath string
}

func NewCompilerHost(options *core.CompilerOptions, currentDirectory string, fs vfs.FS, defaultLibraryPath string) CompilerHost {
	h := &compilerHost{}
	h.options = options
	h.currentDirectory = currentDirectory
	h.fs = fs
	h.defaultLibraryPath = defaultLibraryPath
	return h
}

func (h *compilerHost) FS() vfs.FS {
	return h.fs
}

func (h *compilerHost) DefaultLibraryPath() string {
	return h.defaultLibraryPath
}

func (h *compilerHost) SetOptions(options *core.CompilerOptions) {
	h.options = options
}

func (h *compilerHost) GetCurrentDirectory() string {
	return h.currentDirectory
}

func (h *compilerHost) NewLine() string {
	if h.options == nil {
		return "\n"
	}
	return h.options.NewLine.GetNewLineCharacter()
}

func (h *compilerHost) Trace(msg string) {
	//!!! TODO: implement
}

func (h *compilerHost) GetSourceFile(fileName string, path tspath.Path, languageVersion core.ScriptTarget, packageJsonScope *packagejson.InfoCacheEntry) *ast.SourceFile {
	text, _ := h.FS().ReadFile(fileName)
	if tspath.FileExtensionIs(fileName, tspath.ExtensionJson) {
		return parser.ParseJSONText(fileName, path, text)
	}
	return parser.ParseSourceFile(fileName, path, text, languageVersion, scanner.JSDocParsingModeParseForTypeErrors, h.GetImpliedNodeFormat(fileName, packageJsonScope), packageJsonScope)
}

func (h *compilerHost) GetImpliedNodeFormatForFileWorker(path string, packageJsonScope *packagejson.InfoCacheEntry) core.ResolutionMode {
	var moduleResolution core.ModuleResolutionKind
	if h.options != nil {
		moduleResolution = h.options.GetModuleResolutionKind()
	}

	shouldLookupFromPackageJson := core.ModuleResolutionKindNode16 <= moduleResolution && moduleResolution <= core.ModuleResolutionKindNodeNext || strings.Contains(path, "/node_modules/")

	if tspath.FileExtensionIsOneOf(path, []string{tspath.ExtensionDmts, tspath.ExtensionMts, tspath.ExtensionMjs}) {
		return core.ResolutionModeESM
	}
	if tspath.FileExtensionIsOneOf(path, []string{tspath.ExtensionDcts, tspath.ExtensionCts, tspath.ExtensionCjs}) {
		return core.ResolutionModeCommonJS
	}
	if shouldLookupFromPackageJson && packageJsonScope != nil && tspath.FileExtensionIsOneOf(path, []string{tspath.ExtensionDts, tspath.ExtensionTs, tspath.ExtensionTsx, tspath.ExtensionJs, tspath.ExtensionJsx}) {
		return core.IfElse(packageJsonScope.Contents.Type.Value == "module", core.ResolutionModeESM, core.ResolutionModeCommonJS)
	}

	return core.ResolutionModeNone
}

func (h *compilerHost) GetImpliedNodeFormat(fileName string, packageJsonScope *packagejson.InfoCacheEntry) core.ResolutionMode {
	return h.GetImpliedNodeFormatForFileWorker(fileName, packageJsonScope)
}
