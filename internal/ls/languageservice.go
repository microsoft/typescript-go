package ls

import (
	"context"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/checker"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
)

type LanguageService struct {
	ctx         context.Context
	host        Host
	converters  *Converters
	disposables []func()
}

func NewLanguageService(ctx context.Context, host Host) *LanguageService {
	return &LanguageService{
		ctx:        ctx,
		host:       host,
		converters: NewConverters(host.GetPositionEncoding(), host.GetLineMap),
	}
}

// GetProgram updates the program if the project version has changed.
func (l *LanguageService) GetProgram() *compiler.Program {
	return l.host.GetProgram()
}

func (l *LanguageService) GetTypeChecker(file *ast.SourceFile) *checker.Checker {
	var checker *checker.Checker
	var done func()
	if file == nil {
		checker, done = l.GetProgram().GetTypeChecker(l.ctx)
	} else {
		checker, done = l.GetProgram().GetTypeCheckerForFile(l.ctx, file)
	}
	l.disposables = append(l.disposables, done)
	return checker
}

func (l *LanguageService) Dispose() {
	for _, dispose := range l.disposables {
		dispose()
	}
	l.disposables = nil
}

func (l *LanguageService) tryGetProgramAndFile(fileName string) (*compiler.Program, *ast.SourceFile) {
	program := l.GetProgram()
	file := program.GetSourceFile(fileName)
	return program, file
}

func (l *LanguageService) getSourceFile(documentURI lsproto.DocumentUri) *ast.SourceFile {
	fileName := DocumentURIToFileName(documentURI)
	_, file := l.tryGetProgramAndFile(fileName)
	if file == nil {
		return nil
	}
	return file
}

func (l *LanguageService) getProgramAndFile(documentURI lsproto.DocumentUri) (*compiler.Program, *ast.SourceFile) {
	fileName := DocumentURIToFileName(documentURI)
	program, file := l.tryGetProgramAndFile(fileName)
	if file == nil {
		panic("file not found: " + fileName)
	}
	return program, file
}
