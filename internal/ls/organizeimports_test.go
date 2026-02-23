package ls

import (
	"fmt"
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls/lsutil"
	"github.com/microsoft/typescript-go/internal/parser"
	"github.com/microsoft/typescript-go/internal/stringutil"
	"gotest.tools/v3/assert"
)

func parseImports(importStrings ...string) ([]*ast.Statement, *ast.SourceFile) {
	content := strings.Join(importStrings, "\n")
	sourceFile := parser.ParseSourceFile(ast.SourceFileParseOptions{
		FileName: "/a.ts",
		Path:     "/a.ts",
	}, content, core.ScriptKindTS)

	var imports []*ast.Statement
	for _, stmt := range sourceFile.Statements.Nodes {
		if stmt.Kind == ast.KindImportDeclaration {
			imports = append(imports, stmt)
		}
	}
	if len(imports) != len(importStrings) {
		panic(fmt.Sprintf("parseImports: expected %d imports, got %d", len(importStrings), len(imports)))
	}
	return imports, sourceFile
}

func parseExports(exportStrings ...string) ([]*ast.Statement, *ast.SourceFile) {
	content := strings.Join(exportStrings, "\n")
	sourceFile := parser.ParseSourceFile(ast.SourceFileParseOptions{
		FileName: "/a.ts",
		Path:     "/a.ts",
	}, content, core.ScriptKindTS)

	var exports []*ast.Statement
	for _, stmt := range sourceFile.Statements.Nodes {
		if stmt.Kind == ast.KindExportDeclaration {
			exports = append(exports, stmt)
		}
	}
	if len(exports) != len(exportStrings) {
		panic(fmt.Sprintf("parseExports: expected %d exports, got %d", len(exportStrings), len(exports)))
	}
	return exports, sourceFile
}

func assertListEqual(t *testing.T, actual, expected []*ast.Statement) {
	t.Helper()
	assert.Equal(t, len(actual), len(expected), "declaration count mismatch\nActual: %s\nExpected: %s",
		formatStatements(actual), formatStatements(expected))
	for i := range actual {
		assertEqual(t, actual[i].AsNode(), expected[i].AsNode())
	}
}

func assertEqual(t *testing.T, node1, node2 *ast.Node) {
	t.Helper()
	if node1 == nil && node2 == nil {
		return
	}
	if node1 == nil {
		t.Fatal("node mismatch: first is nil, second is not")
		return
	}
	if node2 == nil {
		t.Fatal("node mismatch: first is not nil, second is nil")
		return
	}

	assert.Equal(t, node1.Kind, node2.Kind, "kind mismatch")

	switch node1.Kind {
	case ast.KindImportDeclaration:
		decl1 := node1.AsImportDeclaration()
		decl2 := node2.AsImportDeclaration()
		assertEqual(t, decl1.ImportClause, decl2.ImportClause)
		assertEqual(t, decl1.ModuleSpecifier.AsNode(), decl2.ModuleSpecifier.AsNode())
	case ast.KindImportClause:
		clause1 := node1.AsImportClause()
		clause2 := node2.AsImportClause()
		assertEqual(t, clause1.Name(), clause2.Name())
		assertEqual(t, clause1.NamedBindings, clause2.NamedBindings)
	case ast.KindNamespaceImport:
		nsi1 := node1.AsNamespaceImport()
		nsi2 := node2.AsNamespaceImport()
		assertEqual(t, nsi1.Name().AsNode(), nsi2.Name().AsNode())
	case ast.KindNamedImports:
		ni1 := node1.AsNamedImports()
		ni2 := node2.AsNamedImports()
		assertNodeListEqual(t, ni1.Elements.Nodes, ni2.Elements.Nodes)
	case ast.KindImportSpecifier:
		is1 := node1.AsImportSpecifier()
		is2 := node2.AsImportSpecifier()
		assertEqual(t, is1.Name().AsNode(), is2.Name().AsNode())
		assertEqual(t, is1.PropertyName, is2.PropertyName)
	case ast.KindExportDeclaration:
		ed1 := node1.AsExportDeclaration()
		ed2 := node2.AsExportDeclaration()
		assertEqual(t, ed1.ExportClause, ed2.ExportClause)
		assertEqual(t, ed1.ModuleSpecifier, ed2.ModuleSpecifier)
	case ast.KindNamedExports:
		ne1 := node1.AsNamedExports()
		ne2 := node2.AsNamedExports()
		assertNodeListEqual(t, ne1.Elements.Nodes, ne2.Elements.Nodes)
	case ast.KindExportSpecifier:
		es1 := node1.AsExportSpecifier()
		es2 := node2.AsExportSpecifier()
		assertEqual(t, es1.Name().AsNode(), es2.Name().AsNode())
		assertEqual(t, es1.PropertyName, es2.PropertyName)
	case ast.KindIdentifier:
		id1 := node1.AsIdentifier()
		id2 := node2.AsIdentifier()
		assert.Equal(t, id1.Text, id2.Text, "identifier text mismatch")
	case ast.KindStringLiteral, ast.KindNoSubstitutionTemplateLiteral:
		assert.Equal(t, node1.Text(), node2.Text(), "string literal text mismatch")
	default:
		t.Fatalf("assertEqual: unhandled node kind: %v", node1.Kind)
	}
}

func assertNodeListEqual(t *testing.T, list1, list2 []*ast.Node) {
	t.Helper()
	assert.Equal(t, len(list1), len(list2), "list length mismatch")
	for i := range list1 {
		assertEqual(t, list1[i], list2[i])
	}
}

func formatStatements(stmts []*ast.Statement) string {
	var parts []string
	for _, stmt := range stmts {
		switch stmt.Kind {
		case ast.KindImportDeclaration:
			parts = append(parts, formatImportDecl(stmt))
		case ast.KindExportDeclaration:
			parts = append(parts, formatExportDecl(stmt))
		default:
			parts = append(parts, "<unknown>")
		}
	}
	return "[" + strings.Join(parts, ", ") + "]"
}

func formatImportDecl(stmt *ast.Statement) string {
	if stmt.Kind != ast.KindImportDeclaration {
		return "<not-import>"
	}
	decl := stmt.AsImportDeclaration()
	var result strings.Builder
	result.WriteString("import ")
	if decl.ImportClause != nil {
		clause := decl.ImportClause.AsImportClause()
		if clause.IsTypeOnly() {
			result.WriteString("type ")
		}
		if clause.Name() != nil {
			result.WriteString(clause.Name().Text())
			if clause.NamedBindings != nil {
				result.WriteString(", ")
			}
		}
		if clause.NamedBindings != nil {
			switch clause.NamedBindings.Kind {
			case ast.KindNamespaceImport:
				result.WriteString("* as ")
				result.WriteString(clause.NamedBindings.AsNamespaceImport().Name().Text())
			case ast.KindNamedImports:
				result.WriteString("{ ")
				elements := clause.NamedBindings.AsNamedImports().Elements.Nodes
				for i, elem := range elements {
					if i > 0 {
						result.WriteString(", ")
					}
					spec := elem.AsImportSpecifier()
					if spec.IsTypeOnly {
						result.WriteString("type ")
					}
					if spec.PropertyName != nil {
						result.WriteString(spec.PropertyName.Text())
						result.WriteString(" as ")
					}
					result.WriteString(spec.Name().Text())
				}
				result.WriteString(" }")
			}
		}
		result.WriteString(" from ")
	}
	result.WriteString("\"")
	result.WriteString(decl.ModuleSpecifier.Text())
	result.WriteString("\"")
	return result.String()
}

func formatExportDecl(stmt *ast.Statement) string {
	if stmt.Kind != ast.KindExportDeclaration {
		return "<not-export>"
	}
	decl := stmt.AsExportDeclaration()
	var result strings.Builder
	result.WriteString("export ")
	if decl.IsTypeOnly {
		result.WriteString("type ")
	}
	if decl.ExportClause != nil {
		switch decl.ExportClause.Kind {
		case ast.KindNamedExports:
			result.WriteString("{ ")
			elements := decl.ExportClause.AsNamedExports().Elements.Nodes
			for i, elem := range elements {
				if i > 0 {
					result.WriteString(", ")
				}
				spec := elem.AsExportSpecifier()
				if spec.IsTypeOnly {
					result.WriteString("type ")
				}
				if spec.PropertyName != nil {
					result.WriteString(spec.PropertyName.Text())
					result.WriteString(" as ")
				}
				result.WriteString(spec.Name().Text())
			}
			result.WriteString(" }")
		case ast.KindNamespaceExport:
			result.WriteString("* as ")
			result.WriteString(decl.ExportClause.Name().Text())
		}
	} else {
		result.WriteString("*")
	}
	if decl.ModuleSpecifier != nil {
		result.WriteString(" from \"")
		result.WriteString(decl.ModuleSpecifier.Text())
		result.WriteString("\"")
	}
	return result.String()
}

func testCoalesceImports(importGroup []*ast.Statement, ignoreCase bool, preferences *lsutil.UserPreferences) []*ast.Statement {
	comparer := getOrdinalStringComparer(ignoreCase)
	if preferences == nil {
		preferences = &lsutil.UserPreferences{}
	}
	specifierComparer := lsutil.GetNamedImportSpecifierComparer(preferences, comparer)
	return coalesceImportsWorker(importGroup, comparer, specifierComparer, nil /*sourceFile*/, nil /*changeTracker*/)
}

func testCoalesceExports(exportGroup []*ast.Statement, ignoreCase bool, preferences *lsutil.UserPreferences) []*ast.Statement {
	if preferences == nil {
		preferences = &lsutil.UserPreferences{
			OrganizeImportsTypeOrder: lsutil.OrganizeImportsTypeOrderLast,
		}
	}
	comparer := getOrdinalStringComparer(ignoreCase)
	specifierComparer := lsutil.GetNamedImportSpecifierComparer(preferences, comparer)
	moduleSpecifierComparer := comparer
	return coalesceExportsWorker(exportGroup, specifierComparer, moduleSpecifierComparer, nil /*sourceFile*/, nil /*changeTracker*/)
}

func getOrdinalStringComparer(ignoreCase bool) func(a, b string) int {
	if ignoreCase {
		return stringutil.CompareStringsCaseInsensitiveEslintCompatible
	}
	return stringutil.CompareStringsCaseSensitive
}

func assertSortsBefore(t *testing.T, importString1, importString2 string) {
	t.Helper()
	imports, _ := parseImports(importString1, importString2)
	assert.Equal(t, len(imports), 2, "expected 2 imports")
	comparer := getOrdinalStringComparer(true)
	moduleSpec1 := imports[0].AsImportDeclaration().ModuleSpecifier
	moduleSpec2 := imports[1].AsImportDeclaration().ModuleSpecifier
	result := lsutil.CompareModuleSpecifiers(moduleSpec1, moduleSpec2, comparer)
	assert.Assert(t, result < 0, "expected %q < %q, got comparison %d", importString1, importString2, result)
	result2 := lsutil.CompareModuleSpecifiers(moduleSpec2, moduleSpec1, comparer)
	assert.Assert(t, result2 > 0, "expected %q > %q, got comparison %d", importString2, importString1, result2)
}

func TestOrganizeImports(t *testing.T) {
	t.Parallel()

	t.Run("Sort imports", func(t *testing.T) {
		t.Parallel()

		t.Run("Sort - non-relative vs non-relative", func(t *testing.T) {
			t.Parallel()
			assertSortsBefore(t,
				`import y from "lib1";`,
				`import x from "lib2";`,
			)
		})

		t.Run("Sort - relative vs relative", func(t *testing.T) {
			t.Parallel()
			assertSortsBefore(t,
				`import y from "./lib1";`,
				`import x from "./lib2";`,
			)
		})

		t.Run("Sort - relative vs non-relative", func(t *testing.T) {
			t.Parallel()
			assertSortsBefore(t,
				`import y from "lib";`,
				`import x from "./lib";`,
			)
		})

		t.Run("Sort - case-insensitive", func(t *testing.T) {
			t.Parallel()
			assertSortsBefore(t,
				`import y from "a";`,
				`import x from "Z";`,
			)
			assertSortsBefore(t,
				`import y from "A";`,
				`import x from "z";`,
			)
		})
	})

	t.Run("Coalesce imports", func(t *testing.T) {
		t.Parallel()

		t.Run("No imports", func(t *testing.T) {
			t.Parallel()
			result := testCoalesceImports(nil, true, nil)
			assert.Equal(t, len(result), 0)
		})

		t.Run("Sort specifiers - case-insensitive", func(t *testing.T) {
			t.Parallel()
			sortedImports, _ := parseImports(`import { default as M, a as n, B, y, Z as O } from "lib";`)
			actual := testCoalesceImports(sortedImports, true, nil)
			expected, _ := parseImports(`import { B, default as M, a as n, Z as O, y } from "lib";`)
			assertListEqual(t, actual, expected)
		})

		t.Run("Combine side-effect-only imports", func(t *testing.T) {
			t.Parallel()
			sortedImports, _ := parseImports(
				`import "lib";`,
				`import "lib";`,
			)
			actual := testCoalesceImports(sortedImports, true, nil)
			expected, _ := parseImports(`import "lib";`)
			assertListEqual(t, actual, expected)
		})

		t.Run("Combine namespace imports", func(t *testing.T) {
			t.Parallel()
			sortedImports, _ := parseImports(
				`import * as x from "lib";`,
				`import * as y from "lib";`,
			)
			actual := testCoalesceImports(sortedImports, true, nil)
			expected, _ := parseImports(
				`import * as x from "lib";`,
				`import * as y from "lib";`,
			)
			assertListEqual(t, actual, expected)
		})

		t.Run("Combine default imports", func(t *testing.T) {
			t.Parallel()
			sortedImports, _ := parseImports(
				`import x from "lib";`,
				`import y from "lib";`,
			)
			actual := testCoalesceImports(sortedImports, true, nil)
			expected, _ := parseImports(`import { default as x, default as y } from "lib";`)
			assertListEqual(t, actual, expected)
		})

		t.Run("Combine property imports", func(t *testing.T) {
			t.Parallel()
			sortedImports, _ := parseImports(
				`import { x } from "lib";`,
				`import { y as z } from "lib";`,
			)
			actual := testCoalesceImports(sortedImports, true, nil)
			expected, _ := parseImports(`import { x, y as z } from "lib";`)
			assertListEqual(t, actual, expected)
		})

		t.Run("Combine side-effect-only import with namespace import", func(t *testing.T) {
			t.Parallel()
			sortedImports, _ := parseImports(
				`import "lib";`,
				`import * as x from "lib";`,
			)
			actual := testCoalesceImports(sortedImports, true, nil)
			expected, _ := parseImports(
				`import "lib";`,
				`import * as x from "lib";`,
			)
			assertListEqual(t, actual, expected)
		})

		t.Run("Combine side-effect-only import with default import", func(t *testing.T) {
			t.Parallel()
			sortedImports, _ := parseImports(
				`import "lib";`,
				`import x from "lib";`,
			)
			actual := testCoalesceImports(sortedImports, true, nil)
			expected, _ := parseImports(
				`import "lib";`,
				`import x from "lib";`,
			)
			assertListEqual(t, actual, expected)
		})

		t.Run("Combine side-effect-only import with property import", func(t *testing.T) {
			t.Parallel()
			sortedImports, _ := parseImports(
				`import "lib";`,
				`import { x } from "lib";`,
			)
			actual := testCoalesceImports(sortedImports, true, nil)
			expected, _ := parseImports(
				`import "lib";`,
				`import { x } from "lib";`,
			)
			assertListEqual(t, actual, expected)
		})

		t.Run("Combine namespace import with default import", func(t *testing.T) {
			t.Parallel()
			sortedImports, _ := parseImports(
				`import * as x from "lib";`,
				`import y from "lib";`,
			)
			actual := testCoalesceImports(sortedImports, true, nil)
			expected, _ := parseImports(`import y, * as x from "lib";`)
			assertListEqual(t, actual, expected)
		})

		t.Run("Combine namespace import with property import", func(t *testing.T) {
			t.Parallel()
			sortedImports, _ := parseImports(
				`import * as x from "lib";`,
				`import { y } from "lib";`,
			)
			actual := testCoalesceImports(sortedImports, true, nil)
			expected, _ := parseImports(
				`import * as x from "lib";`,
				`import { y } from "lib";`,
			)
			assertListEqual(t, actual, expected)
		})

		t.Run("Combine default import with property import", func(t *testing.T) {
			t.Parallel()
			sortedImports, _ := parseImports(
				`import x from "lib";`,
				`import { y } from "lib";`,
			)
			actual := testCoalesceImports(sortedImports, true, nil)
			expected, _ := parseImports(`import x, { y } from "lib";`)
			assertListEqual(t, actual, expected)
		})

		t.Run("Combine many imports", func(t *testing.T) {
			t.Parallel()
			sortedImports, _ := parseImports(
				`import "lib";`,
				`import * as y from "lib";`,
				`import w from "lib";`,
				`import { b } from "lib";`,
				`import "lib";`,
				`import * as x from "lib";`,
				`import z from "lib";`,
				`import { a } from "lib";`,
			)
			actual := testCoalesceImports(sortedImports, true, nil)
			expected, _ := parseImports(
				`import "lib";`,
				`import * as x from "lib";`,
				`import * as y from "lib";`,
				`import { a, b, default as w, default as z } from "lib";`,
			)
			assertListEqual(t, actual, expected)
		})

		// This is descriptive, rather than normative
		t.Run("Combine two namespace imports with one default import", func(t *testing.T) {
			t.Parallel()
			sortedImports, _ := parseImports(
				`import * as x from "lib";`,
				`import * as y from "lib";`,
				`import z from "lib";`,
			)
			actual := testCoalesceImports(sortedImports, true, nil)
			expected, _ := parseImports(
				`import * as x from "lib";`,
				`import * as y from "lib";`,
				`import z from "lib";`,
			)
			assertListEqual(t, actual, expected)
		})

		t.Run("Combine type-only imports separately from other imports", func(t *testing.T) {
			t.Parallel()
			sortedImports, _ := parseImports(
				`import type { x } from "lib";`,
				`import type { y } from "lib";`,
				`import { z } from "lib";`,
			)
			actual := testCoalesceImports(sortedImports, true, nil)
			expected, _ := parseImports(
				`import { z } from "lib";`,
				`import type { x, y } from "lib";`,
			)
			assertListEqual(t, actual, expected)
		})

		t.Run("Do not combine type-only default, namespace, or named imports with each other", func(t *testing.T) {
			t.Parallel()
			sortedImports, _ := parseImports(
				`import type { x } from "lib";`,
				`import type * as y from "lib";`,
				`import type z from "lib";`,
			)
			actual := testCoalesceImports(sortedImports, true, nil)
			// Default import could be rewritten as a named import to combine with `x`,
			// but seems of debatable merit.
			// The three type-only kinds (named, namespace, default) should not be combined into one.
			// The Go implementation outputs them in: namespace, default, named order.
			expected, _ := parseImports(
				`import type * as y from "lib";`,
				`import type z from "lib";`,
				`import type { x } from "lib";`,
			)
			assertListEqual(t, actual, expected)
		})
	})

	t.Run("Coalesce exports", func(t *testing.T) {
		t.Parallel()

		t.Run("No exports", func(t *testing.T) {
			t.Parallel()
			result := testCoalesceExports(nil, true, nil)
			assert.Equal(t, len(result), 0)
		})

		t.Run("Sort specifiers - case-insensitive", func(t *testing.T) {
			t.Parallel()
			sortedExports, _ := parseExports(`export { default as M, a as n, B, y, Z as O } from "lib";`)
			actual := testCoalesceExports(sortedExports, true, nil)
			expected, _ := parseExports(`export { B, default as M, a as n, Z as O, y } from "lib";`)
			assertListEqual(t, actual, expected)
		})

		t.Run("Sort specifiers - type-only-inline", func(t *testing.T) {
			t.Parallel()
			sortedImports, _ := parseImports(`import { type z, y, type x, c, type b, a } from "lib";`)
			prefs := &lsutil.UserPreferences{
				OrganizeImportsTypeOrder: lsutil.OrganizeImportsTypeOrderInline,
			}
			actual := testCoalesceImports(sortedImports, true, prefs)
			expected, _ := parseImports(`import { a, type b, c, type x, y, type z } from "lib";`)
			assertListEqual(t, actual, expected)
		})

		t.Run("Combine namespace re-exports", func(t *testing.T) {
			t.Parallel()
			sortedExports, _ := parseExports(
				`export * from "lib";`,
				`export * from "lib";`,
			)
			actual := testCoalesceExports(sortedExports, true, nil)
			expected, _ := parseExports(`export * from "lib";`)
			assertListEqual(t, actual, expected)
		})

		t.Run("Combine property exports", func(t *testing.T) {
			t.Parallel()
			sortedExports, _ := parseExports(
				`export { x };`,
				`export { y as z };`,
			)
			actual := testCoalesceExports(sortedExports, true, nil)
			expected, _ := parseExports(`export { x, y as z };`)
			assertListEqual(t, actual, expected)
		})

		t.Run("Combine property re-exports", func(t *testing.T) {
			t.Parallel()
			sortedExports, _ := parseExports(
				`export { x } from "lib";`,
				`export { y as z } from "lib";`,
			)
			actual := testCoalesceExports(sortedExports, true, nil)
			expected, _ := parseExports(`export { x, y as z } from "lib";`)
			assertListEqual(t, actual, expected)
		})

		t.Run("Combine namespace re-export with property re-export", func(t *testing.T) {
			t.Parallel()
			sortedExports, _ := parseExports(
				`export * from "lib";`,
				`export { y } from "lib";`,
			)
			actual := testCoalesceExports(sortedExports, true, nil)
			expected, _ := parseExports(
				`export * from "lib";`,
				`export { y } from "lib";`,
			)
			assertListEqual(t, actual, expected)
		})

		t.Run("Combine many exports", func(t *testing.T) {
			t.Parallel()
			sortedExports, _ := parseExports(
				`export { x };`,
				`export { y as w, z as default };`,
				`export { w as q };`,
			)
			actual := testCoalesceExports(sortedExports, true, nil)
			expected, _ := parseExports(
				`export { z as default, w as q, y as w, x };`,
			)
			assertListEqual(t, actual, expected)
		})

		t.Run("Combine many re-exports", func(t *testing.T) {
			t.Parallel()
			sortedExports, _ := parseExports(
				`export { x as a, y } from "lib";`,
				`export * from "lib";`,
				`export { z as b } from "lib";`,
			)
			actual := testCoalesceExports(sortedExports, true, nil)
			expected, _ := parseExports(
				`export * from "lib";`,
				`export { x as a, z as b, y } from "lib";`,
			)
			assertListEqual(t, actual, expected)
		})

		t.Run("Keep type-only exports separate", func(t *testing.T) {
			t.Parallel()
			sortedExports, _ := parseExports(
				`export { x };`,
				`export type { y };`,
			)
			actual := testCoalesceExports(sortedExports, true, nil)
			expected, _ := parseExports(
				`export { x };`,
				`export type { y };`,
			)
			assertListEqual(t, actual, expected)
		})

		t.Run("Combine type-only exports", func(t *testing.T) {
			t.Parallel()
			sortedExports, _ := parseExports(
				`export type { x };`,
				`export type { y };`,
			)
			actual := testCoalesceExports(sortedExports, true, nil)
			expected, _ := parseExports(
				`export type { x, y };`,
			)
			assertListEqual(t, actual, expected)
		})
	})
}
