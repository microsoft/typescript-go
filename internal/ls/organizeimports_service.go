package ls

import (
	"context"
	"slices"
	"strings"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/scanner"
)

// OrganizeImportsMode defines different modes for organizing imports
type OrganizeImportsMode int

const (
	// OrganizeImportsModeAll removes unused imports, combines, and sorts imports
	OrganizeImportsModeAll OrganizeImportsMode = iota
	// OrganizeImportsModeSortAndCombine only sorts and combines imports without removing unused ones
	OrganizeImportsModeSortAndCombine
	// OrganizeImportsModeRemoveUnused only removes unused imports
	OrganizeImportsModeRemoveUnused
)

// OrganizeImports organizes the imports in a TypeScript/JavaScript file
// Currently supports sorting imports while preserving blank line groupings
func (l *LanguageService) OrganizeImports(ctx context.Context, documentURI lsproto.DocumentUri, mode OrganizeImportsMode) ([]lsproto.TextEdit, error) {
	program, file := l.getProgramAndFile(documentURI)
	if file == nil {
		return nil, nil
	}

	edits := l.organizeImportsCore(file, program, mode)
	return edits, nil
}

// organizeImportsCore sorts imports within groups separated by blank lines
func (l *LanguageService) organizeImportsCore(file *ast.SourceFile, program *compiler.Program, mode OrganizeImportsMode) []lsproto.TextEdit {
	_ = program
	_ = mode

	var imports []*ast.Statement
	for _, stmt := range file.Statements.Nodes {
		if stmt.Kind == ast.KindImportDeclaration ||
			stmt.Kind == ast.KindJSImportDeclaration ||
			stmt.Kind == ast.KindImportEqualsDeclaration {
			imports = append(imports, stmt)
		}
	}

	if len(imports) == 0 {
		return nil
	}

	// Group imports by blank lines to preserve intentional grouping
	importGroups := l.groupImportsByNewline(file, imports)

	comparer := func(a, b string) int {
		return strings.Compare(a, b)
	}

	var allEdits []lsproto.TextEdit

	// Process each group independently
	for _, group := range importGroups {
		if len(group) == 0 {
			continue
		}

		sortedGroup := make([]*ast.Statement, len(group))
		copy(sortedGroup, group)

		slices.SortFunc(sortedGroup, func(a, b *ast.Statement) int {
			return CompareImportsOrRequireStatements(a, b, comparer)
		})

		// Check if already sorted
		isSorted := true
		for i := range group {
			if group[i] != sortedGroup[i] {
				isSorted = false
				break
			}
		}

		if isSorted {
			continue
		}

		// Create text edit for this group
		firstImport := group[0]
		lastImport := group[len(group)-1]

		startPos := scanner.SkipTrivia(file.Text(), firstImport.Pos())
		endPos := lastImport.End()

		var newText strings.Builder
		for i, sortedImp := range sortedGroup {
			if i > 0 {
				newText.WriteString("\n")
			}
			importText := scanner.GetTextOfNodeFromSourceText(file.Text(), sortedImp.AsNode(), false)
			newText.WriteString(importText)
		}

		startPosition := l.converters.PositionToLineAndCharacter(file, core.TextPos(startPos))
		endPosition := l.converters.PositionToLineAndCharacter(file, core.TextPos(endPos))

		edit := lsproto.TextEdit{
			Range: lsproto.Range{
				Start: startPosition,
				End:   endPosition,
			},
			NewText: newText.String(),
		}

		allEdits = append(allEdits, edit)
	}

	return allEdits
}

// groupImportsByNewline groups consecutive imports separated by blank lines
func (l *LanguageService) groupImportsByNewline(file *ast.SourceFile, imports []*ast.Statement) [][]*ast.Statement {
	if len(imports) == 0 {
		return nil
	}

	var groups [][]*ast.Statement
	currentGroup := []*ast.Statement{imports[0]}

	for i := 1; i < len(imports); i++ {
		if l.hasBlankLineBeforeStatement(file, imports[i]) {
			groups = append(groups, currentGroup)
			currentGroup = []*ast.Statement{imports[i]}
		} else {
			currentGroup = append(currentGroup, imports[i])
		}
	}

	if len(currentGroup) > 0 {
		groups = append(groups, currentGroup)
	}

	return groups
}

// hasBlankLineBeforeStatement checks if a statement has a blank line before it
// A blank line is indicated by 2+ consecutive newlines in the leading trivia
func (l *LanguageService) hasBlankLineBeforeStatement(file *ast.SourceFile, stmt *ast.Statement) bool {
	text := file.Text()
	pos := stmt.Pos()
	newlineCount := 0

	i := pos - 1
	for i >= 0 {
		ch := text[i]
		if ch == '\n' {
			newlineCount++
			i--
			if i >= 0 && text[i] == '\r' {
				i--
			}
		} else if ch == ' ' || ch == '\t' || ch == '\r' {
			i--
		} else {
			break
		}
	}

	return newlineCount >= 2
}
