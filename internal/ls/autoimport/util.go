package autoimport

import (
	"unicode/utf8"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/stringutil"
)

func getModuleIDOfModuleSymbol(symbol *ast.Symbol) ModuleID {
	if !symbol.IsExternalModule() {
		panic("symbol is not an external module")
	}
	if sourceFile := ast.GetSourceFileOfModule(symbol); sourceFile != nil {
		return ModuleID(sourceFile.Path())
	}
	if ast.IsModuleWithStringLiteralName(symbol.ValueDeclaration) {
		return ModuleID(stringutil.StripQuotes(symbol.Name))
	}
	panic("could not determine module ID of module symbol")
}

// wordIndices splits an identifier into its constituent words based on camelCase and snake_case conventions
// by returning the starting byte indices of each word.
//   - CamelCase
//     ^    ^
//   - snake_case
//     ^     ^
//   - ParseURL
//     ^    ^
//   - __proto__
//     ^
func wordIndices(s string) []int {
	var indices []int
	for byteIndex, runeValue := range s {
		if byteIndex == 0 {
			indices = append(indices, byteIndex)
			continue
		}
		if runeValue == '_' {
			if byteIndex+1 < len(s) && s[byteIndex+1] != '_' {
				indices = append(indices, byteIndex+1)
			}
			continue
		}
		if isUpper(runeValue) && (isLower(core.FirstResult(utf8.DecodeLastRuneInString(s[:byteIndex]))) || (byteIndex+1 < len(s) && isLower(core.FirstResult(utf8.DecodeRuneInString(s[byteIndex+1:]))))) {
			indices = append(indices, byteIndex)
		}
	}
	return indices
}

func isUpper(c rune) bool {
	return c >= 'A' && c <= 'Z'
}

func isLower(c rune) bool {
	return c >= 'a' && c <= 'z'
}
