package ls

import (
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/sourcemap"
)

type LanguageService struct {
	host       Host
	program    *compiler.Program
	converters *Converters
}

func NewLanguageService(
	program *compiler.Program,
	host Host,
) *LanguageService {
	return &LanguageService{
		host:       host,
		program:    program,
		converters: host.Converters(),
	}
}

func (l *LanguageService) GetProgram() *compiler.Program {
	return l.program
}

func (l *LanguageService) tryGetProgramAndFile(fileName string) (*compiler.Program, *ast.SourceFile) {
	program := l.GetProgram()
	file := program.GetSourceFile(fileName)
	return program, file
}

func (l *LanguageService) getProgramAndFile(documentURI lsproto.DocumentUri) (*compiler.Program, *ast.SourceFile) {
	fileName := documentURI.FileName()
	program, file := l.tryGetProgramAndFile(fileName)
	if file == nil {
		panic("file not found: " + fileName)
	}
	return program, file
}

func (l *LanguageService) GetDocumentPositionMapper(fileName string) *sourcemap.DocumentPositionMapper {
	// d, ok := l.documentPositionMappers[fileName]
	// if !ok {
	// 	d = sourcemap.GetDocumentPositionMapper(l, fileName)
	// 	l.documentPositionMappers[fileName] = d
	// }
	// return d
	// !!!
	return nil
}
