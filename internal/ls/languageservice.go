package ls

import (
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/vfs"
)

var _ compiler.CompilerHost = (*LanguageService)(nil)

type LanguageService struct {
	host               Host
	program            *compiler.Program
	lastProjectVersion int
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
func (l *LanguageService) GetSourceFile(fileName string, languageVersion core.ScriptTarget) *ast.SourceFile {
	return l.host.GetSourceFile(fileName, languageVersion)
}

// GetProgram updates the program if the project version has changed.
func (l *LanguageService) GetProgram() *compiler.Program {
	hostVersion := l.host.GetProjectVersion()
	if l.program != nil && hostVersion == l.lastProjectVersion {
		return l.program
	}

	l.lastProjectVersion = hostVersion
	rootFileNames := l.host.GetRootFileNames()
	compilerOptions := l.host.GetCompilerOptions()

	l.program = compiler.NewProgram(compiler.ProgramOptions{
		RootFiles:          rootFileNames,
		Host:               l,
		Options:            compilerOptions,
		DefaultLibraryPath: l.host.GetDefaultLibraryPath(),
	})

	l.program.BindSourceFiles()
	return l.program
}
