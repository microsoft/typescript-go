package ls

import (
	"context"
	"errors"
	"fmt"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/astnav"
	"github.com/microsoft/typescript-go/internal/checker"
	"github.com/microsoft/typescript-go/internal/compiler"
)

var (
	ErrNoSourceFile      = errors.New("source file not found")
	ErrNoTokenAtPosition = errors.New("no token found at position")
)

func (l *LanguageService) GetSymbolAtPosition(ctx context.Context, fileName string, position int) (*ast.Symbol, error) {
	program, file := l.tryGetProgramAndFile(fileName)
	if file == nil {
		return nil, fmt.Errorf("%w: %s", ErrNoSourceFile, fileName)
	}
	node := astnav.GetTokenAtPosition(file, position)
	if node == nil {
		return nil, fmt.Errorf("%w: %s:%d", ErrNoTokenAtPosition, fileName, position)
	}
	checker, done := program.GetTypeCheckerForFile(ctx, file)
	defer done()
	return checker.GetSymbolAtLocation(node), nil
}

func (l *LanguageService) GetSymbolAtLocation(ctx context.Context, node *ast.Node) *ast.Symbol {
	program := l.GetProgram()
	checker, done := program.GetTypeCheckerForFile(ctx, ast.GetSourceFileOfNode(node))
	defer done()
	return checker.GetSymbolAtLocation(node)
}

func (l *LanguageService) GetTypeOfSymbol(ctx context.Context, symbol *ast.Symbol) *checker.Type {
	program := l.GetProgram()
	checker, done := program.GetTypeChecker(ctx)
	defer done()
	return checker.GetTypeOfSymbolAtLocation(symbol, nil)
}

type Position struct {
	Line      int64 `json:"line"`
	Character int64 `json:"character"`
}

func getPosition(file *ast.SourceFile, position int, ls *LanguageService) Position {
	pos := ls.createLspPosition(position, file)
	return Position{
		Line:      int64(pos.Line),
		Character: int64(pos.Character),
	}
}

type DiagnosticId uint32

type Diagnostic struct {
	FileName           string        `json:"fileName"`
	Start              Position      `json:"start"`
	End                Position      `json:"end"`
	StartPos           int           `json:"startPos"`
	EndPos             int           `json:"endPos"`
	Code               int32         `json:"code"`
	Category           string        `json:"category"`
	Message            string        `json:"message"`
	MessageChain       []*Diagnostic `json:"messageChain"`
	RelatedInformation []*Diagnostic `json:"relatedInformation"`
	ReportsUnnecessary bool          `json:"reportsUnnecessary"`
	ReportsDeprecated  bool          `json:"reportsDeprecated"`
	SkippedOnNoEmit    bool          `json:"skippedOnNoEmit"`
	SourceLine         string        `json:"sourceLine"`
}

type diagnosticList struct {
	diagnostics []*Diagnostic
}

func (d *diagnosticList) addDiagnostic(diagnostic *ast.Diagnostic, ls *LanguageService, shouldAdd bool) *Diagnostic {
	startPos := diagnostic.Loc().Pos()
	startPosLineCol := getPosition(diagnostic.File(), startPos, ls)
	lineMap := ls.converters.getLineMap(diagnostic.File().FileName())
	lineStartPos := lineMap.LineStarts[startPosLineCol.Line]
	var lineEndPos int
	if int(startPosLineCol.Line+1) >= len(lineMap.LineStarts) {
		lineEndPos = len(diagnostic.File().Text())
	} else {
		lineEndPos = int(lineMap.LineStarts[startPosLineCol.Line+1]) - 1
	}
	sourceLine := diagnostic.File().Text()[lineStartPos:lineEndPos]

	diag := &Diagnostic{
		FileName:           diagnostic.File().FileName(),
		Start:              startPosLineCol,
		End:                getPosition(diagnostic.File(), diagnostic.Loc().End(), ls),
		StartPos:           startPos,
		EndPos:             diagnostic.Loc().End(),
		SourceLine:         sourceLine,
		Code:               diagnostic.Code(),
		Category:           diagnostic.Category().Name(),
		Message:            diagnostic.Message(),
		MessageChain:       make([]*Diagnostic, 0, len(diagnostic.MessageChain())),
		RelatedInformation: make([]*Diagnostic, 0, len(diagnostic.RelatedInformation())),
	}

	for _, messageChain := range diagnostic.MessageChain() {
		diag.MessageChain = append(diag.MessageChain, d.addDiagnostic(messageChain, ls, false))
	}
	for _, relatedInformation := range diagnostic.RelatedInformation() {
		diag.RelatedInformation = append(diag.RelatedInformation, d.addDiagnostic(relatedInformation, ls, false))
	}

	if shouldAdd {
		d.diagnostics = append(d.diagnostics, diag)
	}
	return diag
}

func (d *diagnosticList) getDiagnostics() []*Diagnostic {
	return d.diagnostics
}

func (l *LanguageService) GetDiagnostics(ctx context.Context, fileNames []string) []*Diagnostic {
	program := l.GetProgram()
	var sourceFiles []*ast.SourceFile
	if len(fileNames) > 0 {
		sourceFiles = make([]*ast.SourceFile, 0, len(fileNames))
		for _, fileName := range fileNames {
			sourceFiles = append(sourceFiles, program.GetSourceFile(fileName))
		}
	} else {
		sourceFiles = program.GetSourceFiles()
	}
	diagnosticList := &diagnosticList{
		diagnostics: make([]*Diagnostic, 0),
	}
	diagnostics := make([]*ast.Diagnostic, 0, len(sourceFiles))
	for _, sourceFile := range sourceFiles {
		diagnostics = append(diagnostics, program.GetSyntacticDiagnostics(ctx, sourceFile)...)
		diagnostics = append(diagnostics, program.GetSemanticDiagnostics(ctx, sourceFile)...)
	}
	diagnostics = compiler.SortAndDeduplicateDiagnostics(diagnostics)
	for _, diagnostic := range diagnostics {
		diagnosticList.addDiagnostic(diagnostic, l, true)
	}
	return diagnosticList.getDiagnostics()
}
