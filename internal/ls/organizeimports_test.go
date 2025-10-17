package ls_test

import (
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls"
	"github.com/microsoft/typescript-go/internal/parser"
	"gotest.tools/v3/assert"
)

func parseImports(t *testing.T, source string) []*ast.Statement {
	t.Helper()

	opts := ast.SourceFileParseOptions{
		FileName: "/test.ts",
	}
	sourceFile := parser.ParseSourceFile(opts, source, core.ScriptKindTS)
	assert.Assert(t, sourceFile != nil, "Failed to parse source")

	var imports []*ast.Statement
	for _, stmt := range sourceFile.Statements.Nodes {
		if stmt.Kind == ast.KindImportDeclaration || stmt.Kind == ast.KindJSImportDeclaration || stmt.Kind == ast.KindImportEqualsDeclaration {
			imports = append(imports, stmt)
		}
	}
	return imports
}

func TestCompareImportsOrRequireStatements_ModuleSpecifiers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		import1  string
		import2  string
		expected int
	}{
		{
			name:     "absolute imports come before relative imports",
			import1:  `import "lodash";`,
			import2:  `import "./local";`,
			expected: -1,
		},
		{
			name:     "relative imports come after absolute imports",
			import1:  `import "./local";`,
			import2:  `import "react";`,
			expected: 1,
		},
		{
			name:     "alphabetical order for absolute imports",
			import1:  `import "react";`,
			import2:  `import "lodash";`,
			expected: 1,
		},
		{
			name:     "alphabetical order for relative imports",
			import1:  `import "./utils";`,
			import2:  `import "./api";`,
			expected: 1,
		},
		{
			name:     "same module specifier",
			import1:  `import { foo } from "react";`,
			import2:  `import { bar } from "react";`,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			imports1 := parseImports(t, tt.import1)
			imports2 := parseImports(t, tt.import2)

			assert.Assert(t, len(imports1) == 1, "Expected 1 import in import1")
			assert.Assert(t, len(imports2) == 1, "Expected 1 import in import2")

			comparer := strings.Compare
			result := ls.CompareImportsOrRequireStatements(imports1[0], imports2[0], comparer)

			if tt.expected == 0 {
				assert.Equal(t, 0, result, "Expected imports to be equal")
			} else if tt.expected < 0 {
				assert.Assert(t, result < 0, "Expected import1 < import2, got %d", result)
			} else {
				assert.Assert(t, result > 0, "Expected import1 > import2, got %d", result)
			}
		})
	}
}

func TestCompareImportsOrRequireStatements_ImportKind(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		import1  string
		import2  string
		expected int
	}{
		{
			name:     "side-effect import comes before type-only import",
			import1:  `import "react";`,
			import2:  `import type { FC } from "react";`,
			expected: -1,
		},
		{
			name:     "type-only import comes before namespace import",
			import1:  `import type { FC } from "react";`,
			import2:  `import * as React from "react";`,
			expected: -1,
		},
		{
			name:     "namespace import comes before default import",
			import1:  `import * as React from "react";`,
			import2:  `import React from "react";`,
			expected: -1,
		},
		{
			name:     "default import comes before named import",
			import1:  `import React from "react";`,
			import2:  `import { useState } from "react";`,
			expected: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			imports1 := parseImports(t, tt.import1)
			imports2 := parseImports(t, tt.import2)

			assert.Assert(t, len(imports1) == 1, "Expected 1 import in import1")
			assert.Assert(t, len(imports2) == 1, "Expected 1 import in import2")

			comparer := strings.Compare
			result := ls.CompareImportsOrRequireStatements(imports1[0], imports2[0], comparer)

			if tt.expected < 0 {
				assert.Assert(t, result < 0, "Expected import1 < import2, got %d", result)
			} else if tt.expected > 0 {
				assert.Assert(t, result > 0, "Expected import1 > import2, got %d", result)
			} else {
				assert.Equal(t, 0, result, "Expected imports to be equal")
			}
		})
	}
}

func TestGetImportDeclarationInsertIndex(t *testing.T) {
	t.Parallel()

	source := `import "side-effect";
import type { Type } from "library";
import { namedImport } from "another";
import defaultImport from "yet-another";
import * as namespace from "namespace-lib";
`

	imports := parseImports(t, source)
	assert.Assert(t, len(imports) > 0, "Expected to parse imports")

	newImportSource := `import { newImport } from "library";`
	newImports := parseImports(t, newImportSource)
	assert.Assert(t, len(newImports) == 1, "Expected 1 new import")

	comparer := func(a, b *ast.Statement) int {
		return ls.CompareImportsOrRequireStatements(a, b, strings.Compare)
	}

	index := ls.GetImportDeclarationInsertIndex(imports, newImports[0], comparer)
	assert.Assert(t, index >= 0 && index <= len(imports),
		"Insert index %d out of range [0, %d]", index, len(imports))
}

func TestGetImportDeclarationInsertIndex_EmptyList(t *testing.T) {
	t.Parallel()

	newImportSource := `import { foo } from "bar";`
	newImports := parseImports(t, newImportSource)
	assert.Assert(t, len(newImports) == 1, "Expected 1 new import")

	comparer := func(a, b *ast.Statement) int {
		return ls.CompareImportsOrRequireStatements(a, b, strings.Compare)
	}

	index := ls.GetImportDeclarationInsertIndex([]*ast.Statement{}, newImports[0], comparer)
	assert.Equal(t, 0, index, "Expected index 0 for empty list")
}

func TestCompareImports_RelativeVsAbsolute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		imports []string
		want    []string
	}{
		{
			name: "mix of relative and absolute",
			imports: []string{
				`import "./utils";`,
				`import "react";`,
				`import "../parent";`,
				`import "lodash";`,
			},
			want: []string{
				`import "lodash";`,
				`import "react";`,
				`import "../parent";`,
				`import "./utils";`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var imports []*ast.Statement
			for _, imp := range tt.imports {
				parsed := parseImports(t, imp)
				assert.Assert(t, len(parsed) == 1, "Expected 1 import")
				imports = append(imports, parsed[0])
			}

			comparer := func(a, b *ast.Statement) int {
				return ls.CompareImportsOrRequireStatements(a, b, strings.Compare)
			}

			for i := 0; i < len(imports); i++ {
				for j := i + 1; j < len(imports); j++ {
					if comparer(imports[i], imports[j]) > 0 {
						imports[i], imports[j] = imports[j], imports[i]
					}
				}
			}

			for i, want := range tt.want {
				wantImports := parseImports(t, want)
				assert.Assert(t, len(wantImports) == 1, "Expected 1 wanted import")

				gotSpec := ls.GetModuleSpecifierExpression(imports[i])
				wantSpec := ls.GetModuleSpecifierExpression(wantImports[0])

				if gotSpec != nil && wantSpec != nil {
					assert.Assert(t, ast.IsStringLiteral(gotSpec), "Expected string literal in got")
					assert.Assert(t, ast.IsStringLiteral(wantSpec), "Expected string literal in want")

					gotText := gotSpec.AsStringLiteral().Text
					wantText := wantSpec.AsStringLiteral().Text
					assert.Equal(t, wantText, gotText, "Import at position %d doesn't match", i)
				}
			}
		})
	}
}

func TestGetModuleSpecifierExpression(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		source   string
		expected string
	}{
		{
			name:     "import declaration",
			source:   `import { foo } from "react";`,
			expected: "react",
		},
		{
			name:     "import equals declaration with external module",
			source:   `import foo = require("lodash");`,
			expected: "lodash",
		},
		{
			name:     "variable statement with require",
			source:   `const foo = require("express");`,
			expected: "express",
		},
		{
			name:     "side-effect import",
			source:   `import "./styles.css";`,
			expected: "./styles.css",
		},
		{
			name:     "namespace import",
			source:   `import * as React from "react";`,
			expected: "react",
		},
		{
			name:     "type-only import",
			source:   `import type { FC } from "react";`,
			expected: "react",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			opts := ast.SourceFileParseOptions{
				FileName: "/test.ts",
			}
			sourceFile := parser.ParseSourceFile(opts, tt.source, core.ScriptKindTS)
			assert.Assert(t, sourceFile != nil, "Failed to parse source")
			assert.Assert(t, len(sourceFile.Statements.Nodes) > 0, "Expected at least one statement")

			stmt := sourceFile.Statements.Nodes[0]
			expr := ls.GetModuleSpecifierExpression(stmt)

			assert.Assert(t, expr != nil, "Expected non-nil module specifier")
			assert.Assert(t, ast.IsStringLiteral(expr), "Expected string literal")

			gotText := expr.AsStringLiteral().Text
			assert.Equal(t, tt.expected, gotText)
		})
	}
}

func TestCompareImportsWithDifferentTypes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		import1  string
		import2  string
		expected string
	}{
		{
			name:     "import equals vs import declaration",
			import1:  `import foo = require("library");`,
			import2:  `import { bar } from "library";`,
			expected: "import2 first",
		},
		{
			name:     "namespace import with named import",
			import1:  `import * as NS from "library";`,
			import2:  `import { foo } from "library";`,
			expected: "import1 first",
		},
		{
			name:     "default and named combined vs named only",
			import1:  `import React, { useState } from "react";`,
			import2:  `import { useEffect } from "react";`,
			expected: "import1 first",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			imports1 := parseImports(t, tt.import1)
			imports2 := parseImports(t, tt.import2)

			if len(imports1) == 0 || len(imports2) == 0 {
				t.Skip("Failed to parse one of the imports")
			}

			comparer := strings.Compare
			result := ls.CompareImportsOrRequireStatements(imports1[0], imports2[0], comparer)

			switch tt.expected {
			case "import1 first":
				assert.Assert(t, result < 0, "Expected import1 < import2, got %d", result)
			case "import2 first":
				assert.Assert(t, result > 0, "Expected import1 > import2, got %d", result)
			case "equal":
				assert.Equal(t, 0, result)
			}
		})
	}
}
