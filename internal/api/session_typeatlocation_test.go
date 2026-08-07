package api

import (
	"context"
	"testing"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/testutil/projecttestutil"
	"gotest.tools/v3/assert"
)

// TestGetTypeAtLocationTupleTypes asks for the type of one node per case and turns
// it into a TypeResponse, covering the shapes an array literal or a tuple
// annotation can produce.
//
// An array literal contextually typed by an empty tuple is the interesting case:
// createTupleTypeEx returns the synthesized tuple target itself at arity zero, and
// createArrayLiteralType then clones it with cloneTypeReference, which keeps
// ObjectFlagsTuple while storing TypeReference data. Such a type cannot be cast
// with AsTupleType, so its tuple shape has to be read off the target.
func TestGetTypeAtLocationTupleTypes(t *testing.T) {
	t.Parallel()
	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}

	testCases := []struct {
		name string
		// content is the whole source file; the type is taken from the first node
		// of kind kind.
		content string
		kind    ast.Kind
		// expectedType is the type the checker prints for that node.
		expectedType string
		// Tuple metadata is only carried by types that have ObjectFlagsTuple set,
		// which in practice means empty tuple targets and clones of them; other
		// tuples are plain references that describe their shape through Target.
		expectedTupleMetadata bool
		expectedFixedLength   int
		expectedReadonly      bool
	}{
		{
			name:                  "array literal for an empty readonly tuple parameter",
			content:               "declare function f(t: readonly []): void;\nexport const g = () => f([]);",
			kind:                  ast.KindArrayLiteralExpression,
			expectedType:          "[]",
			expectedTupleMetadata: true,
		},
		{
			name:                  "array literal for an empty mutable tuple parameter",
			content:               "declare function f(t: []): void;\nexport const g = () => f([]);",
			kind:                  ast.KindArrayLiteralExpression,
			expectedType:          "[]",
			expectedTupleMetadata: true,
		},
		{
			name:                  "array literal for an empty tuple annotation",
			content:               "export const x: readonly [] = [];",
			kind:                  ast.KindArrayLiteralExpression,
			expectedType:          "[]",
			expectedTupleMetadata: true,
		},
		{
			name:                  "array literal for an empty tuple reached through an alias",
			content:               "type E = readonly [];\nexport const x: E = [];",
			kind:                  ast.KindArrayLiteralExpression,
			expectedType:          "[]",
			expectedTupleMetadata: true,
		},
		{
			name:                  "array literal for an empty tuple property",
			content:               "declare function h(o: { items: readonly [] }): void;\nexport const g = () => h({ items: [] });",
			kind:                  ast.KindArrayLiteralExpression,
			expectedType:          "[]",
			expectedTupleMetadata: true,
		},
		{
			name:                  "array literal asserted to an empty tuple",
			content:               "declare function f<T extends readonly unknown[]>(t: T): T;\nexport const g = () => f([] as readonly []);",
			kind:                  ast.KindArrayLiteralExpression,
			expectedType:          "[]",
			expectedTupleMetadata: true,
		},
		{
			name:                  "empty tuple type node",
			content:               "export const x: readonly [] = [];",
			kind:                  ast.KindTupleType,
			expectedType:          "readonly []",
			expectedTupleMetadata: true,
			expectedReadonly:      true,
		},
		{
			name:                  "identifier annotated with an empty tuple",
			content:               "const x: readonly [] = [];\nexport { x };",
			kind:                  ast.KindIdentifier,
			expectedType:          "readonly []",
			expectedTupleMetadata: true,
			expectedReadonly:      true,
		},
		{
			name:         "array literal for a single element tuple parameter",
			content:      "declare function f(t: readonly [number]): void;\nexport const g = () => f([1]);",
			kind:         ast.KindArrayLiteralExpression,
			expectedType: "[number]",
		},
		{
			name:         "array literal for a single element tuple annotation",
			content:      "export const x: readonly [number] = [1];",
			kind:         ast.KindArrayLiteralExpression,
			expectedType: "[number]",
		},
		{
			name:         "single element tuple type node",
			content:      "export const x: readonly [number] = [1];",
			kind:         ast.KindTupleType,
			expectedType: "readonly [number]",
		},
		{
			name:         "array literal for an array parameter",
			content:      "declare function f(t: number[]): void;\nexport const g = () => f([]);",
			kind:         ast.KindArrayLiteralExpression,
			expectedType: "never[]",
		},
		{
			name:         "array literal without a contextual type",
			content:      "export const x = [];",
			kind:         ast.KindArrayLiteralExpression,
			expectedType: "never[]",
		},
	}

	const fileName = "/home/projects/p/src/index.ts"

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			projectSession, _ := projecttestutil.Setup(map[string]any{
				"/home/projects/p/tsconfig.json": `{ "compilerOptions": { "strict": true } }`,
				fileName:                         testCase.content,
			})
			defer projectSession.Close()
			session := NewSession(projectSession, nil)
			defer session.Close()

			ctx := context.Background()
			snapshot, err := session.handleUpdateSnapshot(ctx, &UpdateSnapshotParams{
				OpenFiles: []DocumentIdentifier{{FileName: fileName}},
			})
			assert.NilError(t, err)
			assert.Assert(t, len(snapshot.Projects) > 0, "expected at least one project")

			setup, err := session.setupChecker(ctx, snapshot.Snapshot, snapshot.Projects[0].Id)
			assert.NilError(t, err)
			defer setup.done()

			sourceFile := setup.program.GetSourceFile(fileName)
			assert.Assert(t, sourceFile != nil, "expected %s to be part of the program", fileName)
			node := findFirstNodeOfKind(sourceFile.AsNode(), testCase.kind)
			assert.Assert(t, node != nil, "expected a node of kind %s", testCase.kind.String())

			nodeType := setup.checker.GetTypeAtLocation(node)
			assert.Equal(t, setup.checker.TypeToString(nodeType), testCase.expectedType)

			response := setup.newTypeResponse(nodeType)
			assert.Assert(t, response != nil)
			assert.Assert(t, response.Target != 0, "expected a target type")

			if !testCase.expectedTupleMetadata {
				assert.Assert(t, response.FixedLength == nil, "expected no tuple metadata")
				assert.Assert(t, response.TupleReadonly == nil, "expected no tuple metadata")
				return
			}
			assert.Assert(t, response.FixedLength != nil, "expected tuple metadata")
			assert.Equal(t, *response.FixedLength, testCase.expectedFixedLength)
			assert.Assert(t, response.TupleReadonly != nil, "expected tuple metadata")
			assert.Equal(t, *response.TupleReadonly, testCase.expectedReadonly)
			assert.Equal(t, len(response.ElementFlags), testCase.expectedFixedLength)
		})
	}
}

func findFirstNodeOfKind(node *ast.Node, kind ast.Kind) *ast.Node {
	var found *ast.Node
	var visit ast.Visitor
	visit = func(child *ast.Node) bool {
		if found != nil {
			return true
		}
		if child.Kind == kind {
			found = child
			return true
		}
		return child.ForEachChild(visit)
	}
	visit(node)
	return found
}
