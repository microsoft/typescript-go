package ls

import (
	"cmp"
	"strings"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/modulespecifiers"
	"github.com/microsoft/typescript-go/internal/tspath"
)

// GetImportDeclarationInsertIndex determines where to insert a new import in a sorted list
// It uses binary search to find the appropriate insertion index
// statement = anyImportOrRequireStatement
func GetImportDeclarationInsertIndex(sortedImports []*ast.Statement, newImport *ast.Statement, comparer func(a, b *ast.Statement) int) int {
	n := len(sortedImports)
	if n == 0 {
		return 0
	}

	low, high := 0, n
	for low < high {
		mid := low + (high-low)/2
		if comparer(sortedImports[mid], newImport) < 0 {
			low = mid + 1
		} else {
			high = mid
		}
	}
	return low
}

// CompareImportsOrRequireStatements compares two import statements for sorting purposes
// Returns:
//   - negative if s1 should come before s2
//   - 0 if they are equivalent
//   - positive if s1 should come after s2
func CompareImportsOrRequireStatements(s1, s2 *ast.Statement, comparer func(a, b string) int) int {
	// First compare module specifiers
	comparison := compareModuleSpecifiersWorker(
		GetModuleSpecifierExpression(s1),
		GetModuleSpecifierExpression(s2),
		comparer,
	)
	if comparison != 0 {
		return comparison
	}
	// If module specifiers are equal, compare by import kind
	return compareImportKind(s1, s2)
}

// GetModuleSpecifierExpression extracts the module specifier expression from an import/export statement
func GetModuleSpecifierExpression(declaration *ast.Statement) *ast.Expression {
	switch declaration.Kind {
	case ast.KindImportDeclaration, ast.KindJSImportDeclaration:
		return declaration.AsImportDeclaration().ModuleSpecifier
	case ast.KindImportEqualsDeclaration:
		moduleRef := declaration.AsImportEqualsDeclaration().ModuleReference
		if moduleRef != nil && moduleRef.Kind == ast.KindExternalModuleReference {
			return moduleRef.AsExternalModuleReference().Expression
		}
		return nil
	case ast.KindVariableStatement:
		// Handle require statements: const x = require("...")
		declList := declaration.AsVariableStatement().DeclarationList
		if declList != nil && declList.Kind == ast.KindVariableDeclarationList {
			decls := declList.AsVariableDeclarationList().Declarations.Nodes
			if len(decls) > 0 {
				varDecl := decls[0].AsVariableDeclaration()
				if varDecl.Initializer != nil && ast.IsCallExpression(varDecl.Initializer) {
					call := varDecl.Initializer.AsCallExpression()
					if len(call.Arguments.Nodes) > 0 {
						return call.Arguments.Nodes[0]
					}
				}
			}
		}
		return nil
	}
	return nil
}

// compareModuleSpecifiersWorker compares two module specifier expressions.
// Ordering:
//  1. undefined module specifiers come last
//  2. Relative imports come after absolute imports
//  3. Otherwise, compare by the comparer function
func compareModuleSpecifiersWorker(m1, m2 *ast.Expression, comparer func(a, b string) int) int {
	name1 := getExternalModuleNameText(m1)
	name2 := getExternalModuleNameText(m2)

	// undefined names come last
	if comparison := compareBooleans(name1 == "", name2 == ""); comparison != 0 {
		return comparison
	}

	// If both are defined, absolute imports come before relative imports
	if name1 != "" && name2 != "" {
		isRelative1 := tspath.IsExternalModuleNameRelative(name1)
		isRelative2 := tspath.IsExternalModuleNameRelative(name2)
		// Reverse parameter order because we want absolute imports (isRelative=false) before relative imports (isRelative=true)
		if comparison := compareBooleans(isRelative2, isRelative1); comparison != 0 {
			return comparison
		}

		// Finally, compare by the provided comparer
		return comparer(name1, name2)
	}

	return 0
}

// getExternalModuleNameText extracts the text of a module name from an expression
func getExternalModuleNameText(expr *ast.Expression) string {
	if expr == nil {
		return ""
	}

	if ast.IsStringLiteral(expr) {
		return expr.AsStringLiteral().Text
	}

	return ""
}

// compareImportKind compares imports by their kind/type for sorting
// Import order:
//  1. Side-effect imports (import "foo")
//  2. Type-only imports (import type { Foo } from "foo")
//  3. Namespace imports (import * as foo from "foo")
//  4. Default imports (import foo from "foo")
//  5. Named imports (import { foo } from "foo")
//  6. ImportEquals declarations (import foo = require("foo"))
//  7. Require variable statements (const foo = require("foo"))
func compareImportKind(s1, s2 *ast.Statement) int {
	order1 := getImportKindOrder(s1)
	order2 := getImportKindOrder(s2)
	return cmp.Compare(order1, order2)
}

// getImportKindOrder returns the sorting order for different import kinds
func getImportKindOrder(statement *ast.Statement) int {
	switch statement.Kind {
	case ast.KindImportDeclaration, ast.KindJSImportDeclaration:
		importDecl := statement.AsImportDeclaration()
		// Side-effect import (no import clause)
		if importDecl.ImportClause == nil {
			return 0
		}
		importClause := importDecl.ImportClause.AsImportClause()
		// Type-only import
		if importClause.IsTypeOnly {
			return 1
		}
		// Namespace import (import * as foo)
		if importClause.NamedBindings != nil && importClause.NamedBindings.Kind == ast.KindNamespaceImport {
			return 2
		}
		// Default import (import foo)
		if importClause.Name() != nil {
			return 3
		}
		// Named import (import { foo })
		return 4
	case ast.KindImportEqualsDeclaration:
		return 5
	case ast.KindVariableStatement:
		return 6
	}
	return 999
}

// returns `-1` if `a` is better than `b`
//
//	note: this sorts in descending order of preference; different than convention in other cmp-like functions
func compareModuleSpecifiers(
	a *ImportFix, // !!! ImportFixWithModuleSpecifier
	b *ImportFix, // !!! ImportFixWithModuleSpecifier
	importingFile *ast.SourceFile, // | FutureSourceFile,
	program *compiler.Program,
	preferences UserPreferences,
	allowsImportingSpecifier func(specifier string) bool,
	toPath func(fileName string) tspath.Path,
) int {
	if a.kind == ImportFixKindUseNamespace || b.kind == ImportFixKindUseNamespace {
		return 0
	}
	if comparison := compareBooleans(
		b.moduleSpecifierKind != modulespecifiers.ResultKindNodeModules || allowsImportingSpecifier(b.moduleSpecifier),
		a.moduleSpecifierKind != modulespecifiers.ResultKindNodeModules || allowsImportingSpecifier(a.moduleSpecifier),
	); comparison != 0 {
		return comparison
	}
	if comparison := compareModuleSpecifierRelativity(a, b, preferences); comparison != 0 {
		return comparison
	}
	if comparison := compareNodeCoreModuleSpecifiers(a.moduleSpecifier, b.moduleSpecifier, importingFile, program); comparison != 0 {
		return comparison
	}
	if comparison := compareBooleans(isFixPossiblyReExportingImportingFile(a, importingFile.Path(), toPath), isFixPossiblyReExportingImportingFile(b, importingFile.Path(), toPath)); comparison != 0 {
		return comparison
	}
	if comparison := compareNumberOfDirectorySeparators(a.moduleSpecifier, b.moduleSpecifier); comparison != 0 {
		return comparison
	}
	return 0
}

// True > False
func compareBooleans(a, b bool) int {
	if a && !b {
		return -1
	} else if !a && b {
		return 1
	}
	return 0
}

// returns `-1` if `a` is better than `b`
func compareModuleSpecifierRelativity(a *ImportFix, b *ImportFix, preferences UserPreferences) int {
	switch preferences.ImportModuleSpecifierPreference {
	case modulespecifiers.ImportModuleSpecifierPreferenceNonRelative, modulespecifiers.ImportModuleSpecifierPreferenceProjectRelative:
		return compareBooleans(a.moduleSpecifierKind == modulespecifiers.ResultKindRelative, b.moduleSpecifierKind == modulespecifiers.ResultKindRelative)
	}
	return 0
}

func compareNodeCoreModuleSpecifiers(a, b string, importingFile *ast.SourceFile, program *compiler.Program) int {
	if strings.HasPrefix(a, "node:") && !strings.HasPrefix(b, "node:") {
		if shouldUseUriStyleNodeCoreModules(importingFile, program) {
			return -1
		}
		return 1
	}
	if strings.HasPrefix(b, "node:") && !strings.HasPrefix(a, "node:") {
		if shouldUseUriStyleNodeCoreModules(importingFile, program) {
			return 1
		}
		return -1
	}
	return 0
}

func shouldUseUriStyleNodeCoreModules(file *ast.SourceFile, program *compiler.Program) bool {
	for _, node := range file.Imports() {
		if core.NodeCoreModules()[node.Text()] && !core.ExclusivelyPrefixedNodeCoreModules[node.Text()] {
			if strings.HasPrefix(node.Text(), "node:") {
				return true
			} else {
				return false
			}
		}
	}

	return program.UsesUriStyleNodeCoreModules()
}

// This is a simple heuristic to try to avoid creating an import cycle with a barrel re-export.
// E.g., do not `import { Foo } from ".."` when you could `import { Foo } from "../Foo"`.
// This can produce false positives or negatives if re-exports cross into sibling directories
// (e.g. `export * from "../whatever"`) or are not named "index".
func isFixPossiblyReExportingImportingFile(fix *ImportFix, importingFilePath tspath.Path, toPath func(fileName string) tspath.Path) bool {
	if fix.isReExport != nil && *(fix.isReExport) &&
		fix.exportInfo != nil && fix.exportInfo.moduleFileName != "" && isIndexFileName(fix.exportInfo.moduleFileName) {
		reExportDir := toPath(tspath.GetDirectoryPath(fix.exportInfo.moduleFileName))
		return strings.HasPrefix(string(importingFilePath), string(reExportDir))
	}
	return false
}

func compareNumberOfDirectorySeparators(path1, path2 string) int {
	return cmp.Compare(strings.Count(path1, "/"), strings.Count(path2, "/"))
}

func isIndexFileName(fileName string) bool {
	fileName = tspath.GetBaseFileName(fileName)
	if tspath.FileExtensionIsOneOf(fileName, []string{".js", ".jsx", ".d.ts", ".ts", ".tsx"}) {
		fileName = tspath.RemoveFileExtension(fileName)
	}
	return fileName == "index"
}
