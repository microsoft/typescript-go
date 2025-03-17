package ls

import (
	"strings"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/compiler/packagejson"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
)

var _ compiler.CompilerHost = (*LanguageService)(nil)

type LanguageService struct {
	host Host
}

func NewLanguageService(host Host) *LanguageService {
	return &LanguageService{
		host: host,
	}
}

// FS implements compiler.CompilerHost.
func (l *LanguageService) FS() vfs.FS {
	return l.host.FS()
}

// DefaultLibraryPath implements compiler.CompilerHost.
func (l *LanguageService) DefaultLibraryPath() string {
	return l.host.DefaultLibraryPath()
}

// GetCurrentDirectory implements compiler.CompilerHost.
func (l *LanguageService) GetCurrentDirectory() string {
	return l.host.GetCurrentDirectory()
}

// NewLine implements compiler.CompilerHost.
func (l *LanguageService) NewLine() string {
	return l.host.NewLine()
}

// Trace implements compiler.CompilerHost.
func (l *LanguageService) Trace(msg string) {
	l.host.Trace(msg)
}

// GetSourceFile implements compiler.CompilerHost.
func (l *LanguageService) GetSourceFile(fileName string, path tspath.Path, languageVersion core.ScriptTarget, packageJsonScope *packagejson.InfoCacheEntry) *ast.SourceFile {
	return l.host.GetSourceFile(fileName, path, languageVersion, packageJsonScope)
}

func (l *LanguageService) GetImpliedNodeFormatForFileWorker(path string, packageJsonScope *packagejson.InfoCacheEntry) core.ResolutionMode {
	moduleResolution := l.GetProgram().Options().GetModuleResolutionKind()
	shouldLookupFromPackageJson := core.ModuleResolutionKindNode16 <= moduleResolution && moduleResolution <= core.ModuleResolutionKindNodeNext || strings.Contains(path, "/node_modules/")

	if tspath.FileExtensionIsOneOf(path, []string{tspath.ExtensionDmts, tspath.ExtensionMts, tspath.ExtensionMjs}) {
		return core.ResolutionModeESM
	}
	if tspath.FileExtensionIsOneOf(path, []string{tspath.ExtensionDcts, tspath.ExtensionCts, tspath.ExtensionCjs}) {
		return core.ResolutionModeCommonJS
	}
	if shouldLookupFromPackageJson && tspath.FileExtensionIsOneOf(path, []string{tspath.ExtensionDts, tspath.ExtensionTs, tspath.ExtensionTsx, tspath.ExtensionJs, tspath.ExtensionJsx}) {
		return core.IfElse(packageJsonScope.Contents.Type.Value == "module", core.ResolutionModeESM, core.ResolutionModeCommonJS)
	}

	return core.ResolutionModeNone
}

func (l *LanguageService) GetImpliedNodeFormat(fileName string, packageJsonScope *packagejson.InfoCacheEntry) core.ResolutionMode {
	return l.GetImpliedNodeFormatForFileWorker(fileName, packageJsonScope)
}

// GetProgram updates the program if the project version has changed.
func (l *LanguageService) GetProgram() *compiler.Program {
	return l.host.GetProgram()
}

func (l *LanguageService) getProgramAndFile(fileName string) (*compiler.Program, *ast.SourceFile) {
	program := l.GetProgram()
	file := program.GetSourceFile(fileName)
	if file == nil {
		panic("file not found")
	}
	return program, file
}
