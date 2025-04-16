package checker

import (
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/printer"
)

func (ch *Checker) IsTypeSymbolAccessible(symbol *ast.Symbol, enclosingDeclaration *ast.Node) bool {
	return false // !!!
}

func (ch *Checker) IsValueSymbolAccessible(symbol *ast.Symbol, enclosingDeclaration *ast.Node) bool {
	return false // !!!
}

/**
 * Check if the given symbol in given enclosing declaration is accessible and mark all associated alias to be visible if requested
 *
 * @param symbol a Symbol to check if accessible
 * @param enclosingDeclaration a Node containing reference to the symbol
 * @param meaning a SymbolFlags to check if such meaning of the symbol is accessible
 * @param shouldComputeAliasToMakeVisible a boolean value to indicate whether to return aliases to be mark visible in case the symbol is accessible
 */

func (c *Checker) IsSymbolAccessible(symbol *ast.Symbol, enclosingDeclaration *ast.Node, meaning ast.SymbolFlags, shouldComputeAliasesToMakeVisible bool) printer.SymbolAccessibilityResult {
	return c.isSymbolAccessibleWorker(symbol, enclosingDeclaration, meaning, shouldComputeAliasesToMakeVisible, true /*allowModules*/)
}

func (c *Checker) isSymbolAccessibleWorker(symbol *ast.Symbol, enclosingDeclaration *ast.Node, meaning ast.SymbolFlags, shouldComputeAliasesToMakeVisible bool, allowModules bool) printer.SymbolAccessibilityResult {
	// if symbol != nil && enclosingDeclaration != nil {
	// 	result := c.isAnySymbolAccessible([]*ast.Symbol{symbol}, enclosingDeclaration, symbol, meaning, shouldComputeAliasesToMakeVisible, allowModules)
	// 	if result != nil {
	// 		return result
	// 	}

	// 	// This could be a symbol that is not exported in the external module
	// 	// or it could be a symbol from different external module that is not aliased and hence cannot be named
	// 	symbolExternalModule := forEach(symbol.Declarations, c.getExternalModuleContainer)
	// 	if symbolExternalModule != nil {
	// 		enclosingExternalModule := c.getExternalModuleContainer(enclosingDeclaration)
	// 		if symbolExternalModule != enclosingExternalModule {
	// 			// name from different external module that is not visible
	// 			return SymbolAccessibilityResult{
	// 				accessibility:   SymbolAccessibilityCannotBeNamed,
	// 				errorSymbolName: c.symbolToString(symbol, enclosingDeclaration, meaning),
	// 				errorModuleName: c.symbolToString(symbolExternalModule),
	// 				errorNode:       ifElse(isInJSFile(enclosingDeclaration), enclosingDeclaration, nil),
	// 			}
	// 		}
	// 	}

	// 	// Just a local name that is not accessible
	// 	return SymbolAccessibilityResult{
	// 		accessibility:   SymbolAccessibilityNotAccessible,
	// 		errorSymbolName: c.symbolToString(symbol, enclosingDeclaration, meaning),
	// 	}
	// }

	// return SymbolAccessibilityResult{
	// 	accessibility: SymbolAccessibilityAccessible,
	// }
	return printer.SymbolAccessibilityResult{} // !!!
}
