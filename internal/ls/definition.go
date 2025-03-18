package ls

import (
	"fmt"
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/astnav"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/scanner"
)

func (l *LanguageService) ProvideDefinitions(fileName string, position int) ([]Location, error) {
	program, file, err := l.getProgramAndFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to get program and file: %w", err)
	}

	node := astnav.GetTouchingPropertyName(file, position)
	if node.Kind == ast.KindSourceFile {
		return nil, nil
	}

	checker := program.GetTypeChecker()
	if symbol := checker.GetSymbolAtLocation(node); symbol != nil {
		if symbol.Flags&ast.SymbolFlagsAlias != 0 {
			if resolved, ok := checker.ResolveAlias(symbol); ok {
				symbol = resolved
			}
		}

		locations := make([]Location, 0, len(symbol.Declarations))
		for _, decl := range symbol.Declarations {
			file := ast.GetSourceFileOfNode(decl)
			loc := decl.Loc
			pos := scanner.GetTokenPosOfNode(decl, file, false /*includeJsDoc*/)

			locations = append(locations, Location{
				FileName: file.FileName(),
				Range:    core.NewTextRange(pos, loc.End()),
			})
		}
		return locations, nil
	}
	return nil, nil
}
