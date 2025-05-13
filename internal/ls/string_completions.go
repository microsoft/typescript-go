package ls

import (
	"fmt"
	"slices"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/checker"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
)

type completionsFromTypes struct {
	types           []*checker.StringLiteralType
	isNewIdentifier bool
}

type completionsFromProperties struct {
	symbols           []*ast.Symbol
	hasIndexSignature bool
}

// !!!
type stringLiteralCompletions = any

func (l *LanguageService) getStringLiteralCompletions(
	file *ast.SourceFile,
	position int,
	contextToken *ast.Node,
	compilerOptions *core.CompilerOptions,
	program *compiler.Program,
	preferences *UserPreferences,
) *lsproto.CompletionList {
	// !!! reference comment
	if IsInString(file, position, contextToken) {
		if contextToken == nil || !ast.IsStringOrNumericLiteralLike(contextToken) {
			return nil
		}
		entries := l.getStringLiteralCompletionEntries(
			file,
			contextToken,
			position,
			program,
			preferences)
		// return l.convertStringLiteralCompletions(entries) // !!! HERE
	}
	return nil
}

func (l *LanguageService) getStringLiteralCompletionEntries(
	file *ast.SourceFile,
	node *ast.StringLiteralLike,
	position int,
	program *compiler.Program,
	preferences *UserPreferences,
) stringLiteralCompletions {
	typeChecker := program.GetTypeChecker()
	parent := walkUpParentheses(node.Parent)
	switch parent.Kind {
	case ast.KindLiteralType:
		grandparent := walkUpParentheses(parent.Parent)
		if grandparent.Kind == ast.KindImportType {
			return getStringLiteralCompletionsFromModuleNames(
				file,
				node,
				program,
				preferences,
			)
		}
		return fromUnionableLiteralType(grandparent, parent, position, typeChecker)
	case ast.KindPropertyAssignment:
		if ast.IsObjectLiteralExpression(parent.Parent) && parent.Name() == node {
			// Get quoted name of properties of the object literal expression
			// i.e. interface ConfigFiles {
			//          'jspm:dev': string
			//      }
			//      let files: ConfigFiles = {
			//          '/*completion position*/'
			//      }
			//
			//      function foo(c: ConfigFiles) {}
			//      foo({
			//          '/*completion position*/'
			//      });
			return stringLiteralCompletionsForObjectLiteral(typeChecker, parent.Parent)
		}
		result := fromContextualType(checker.ContextFlagsCompletions, node, typeChecker)
		if result != nil {
			return result
		}
		return fromContextualType(checker.ContextFlagsNone, node, typeChecker)
	case ast.KindElementAccessExpression:
		expression := parent.Expression()
		argumentExpression := parent.AsElementAccessExpression().ArgumentExpression
		if node == ast.SkipParentheses(argumentExpression) {
			// Get all names of properties on the expression
			// i.e. interface A {
			//      'prop1': string
			// }
			// let a: A;
			// a['/*completion position*/']
			t := typeChecker.GetTypeAtLocation(expression)
			return stringLiteralCompletionsFromProperties(t, typeChecker)
		}
		return nil
		// !!! HERE
	}
}

func fromContextualType(contextFlags checker.ContextFlags, node *ast.Node, typeChecker *checker.Checker) *completionsFromTypes {
	// Get completion for string literal from string literal type
	// i.e. var x: "hi" | "hello" = "/*completion position*/"
	types := getStringLiteralTypes(getContextualTypeFromParent(node, typeChecker, contextFlags), nil, typeChecker)
	if len(types) == 0 {
		return nil
	}
	return &completionsFromTypes{
		types:           types,
		isNewIdentifier: false,
	}
}

func fromUnionableLiteralType(grandparent *ast.Node, parent *ast.Node, position int, typeChecker *checker.Checker) stringLiteralCompletions {
	switch grandparent.Kind {
	case ast.KindExpressionWithTypeArguments, ast.KindTypeReference:
		typeArgument := ast.FindAncestor(parent, func(n *ast.Node) bool { return n.Parent == grandparent })
		if typeArgument != nil {
			t := typeChecker.GetTypeArgumentConstraint(typeArgument)
			return &completionsFromTypes{
				types:           getStringLiteralTypes(t, nil, typeChecker),
				isNewIdentifier: false,
			}
		}
		return nil
	case ast.KindIndexedAccessType:
		// Get all apparent property names
		// i.e. interface Foo {
		//          foo: string;
		//          bar: string;
		//      }
		//      let x: Foo["/*completion position*/"]
		indexType := grandparent.AsIndexedAccessTypeNode().IndexType
		objectType := grandparent.AsIndexedAccessTypeNode().ObjectType
		if !indexType.Loc.ContainsInclusive(position) {
			return nil
		}
		t := typeChecker.GetTypeFromTypeNode(objectType)
		return stringLiteralCompletionsFromProperties(t, typeChecker)
	case ast.KindUnionType:
		result := fromUnionableLiteralType(
			walkUpParentheses(grandparent.Parent),
			parent,
			position,
			typeChecker)
		if result == nil {
			return nil
		}
		alreadyUsedTypes := getAlreadyUsedTypesInStringLiteralUnion(grandparent, parent)
		switch result := result.(type) {
		case *completionsFromProperties:
			return &completionsFromProperties{
				symbols: core.Filter(
					result.symbols,
					func(s *ast.Symbol) bool { return !slices.Contains(alreadyUsedTypes, s.Name) },
				),
				hasIndexSignature: result.hasIndexSignature,
			}
		case *completionsFromTypes:
			return &completionsFromTypes{
				types: core.Filter(result.types, func(t *checker.StringLiteralType) bool {
					return !slices.Contains(alreadyUsedTypes, t.AsLiteralType().Value().(string))
				}),
				isNewIdentifier: false,
			}
		default:
			panic(fmt.Sprintf("Unexpected result type: %T", result))
		}
	default:
		return nil
	}
}

func stringLiteralCompletionsForObjectLiteral(
	typeChecker *checker.Checker,
	objectLiteralExpression *ast.ObjectLiteralExpressionNode) *completionsFromProperties {
	contextualType := typeChecker.GetContextualType(objectLiteralExpression, checker.ContextFlagsNone)
	if contextualType == nil {
		return nil
	}

	completionsType := typeChecker.GetContextualType(objectLiteralExpression, checker.ContextFlagsCompletions)
	symbols := getPropertiesForObjectExpression(
		contextualType,
		completionsType,
		objectLiteralExpression,
		typeChecker)

	return &completionsFromProperties{
		symbols:           symbols,
		hasIndexSignature: hasIndexSignature(contextualType, typeChecker),
	}
}

func stringLiteralCompletionsFromProperties(t *checker.Type, typeChecker *checker.Checker) *completionsFromProperties {
	return &completionsFromProperties{
		symbols: core.Filter(typeChecker.GetApparentProperties(t), func(s *ast.Symbol) bool {
			return !(s.ValueDeclaration != nil && ast.IsPrivateIdentifierClassElementDeclaration(s.ValueDeclaration))
		}),
		hasIndexSignature: hasIndexSignature(t, typeChecker),
	}
}

func getStringLiteralCompletionsFromModuleNames(
	file *ast.SourceFile,
	node *ast.LiteralExpression,
	program *compiler.Program,
	preferences *UserPreferences) stringLiteralCompletions {
	// !!! here
	return nil
}

func walkUpParentheses(node *ast.Node) *ast.Node {
	switch node.Kind {
	case ast.KindParenthesizedType:
		return ast.WalkUpParenthesizedTypes(node)
	case ast.KindParenthesizedExpression:
		return ast.WalkUpParenthesizedExpressions(node)
	default:
		return node
	}
}

func getStringLiteralTypes(t *checker.Type, uniques *core.Set[string], typeChecker *checker.Checker) []*checker.StringLiteralType {
	if t == nil {
		return nil
	}
	if uniques == nil {
		uniques = &core.Set[string]{}
	}
	t = skipConstraint(t, typeChecker)
	if t.IsUnion() {
		var types []*checker.StringLiteralType
		for _, elementType := range t.Types() {
			types = append(types, getStringLiteralTypes(elementType, uniques, typeChecker)...)
		}
		return types
	}
	if t.IsStringLiteral() && !t.IsEnumLiteral() && uniques.AddIfAbsent(t.AsLiteralType().Value().(string)) {
		return []*checker.StringLiteralType{t}
	}
	return nil
}

func getAlreadyUsedTypesInStringLiteralUnion(union *ast.UnionType, current *ast.LiteralType) []string {
	typesList := union.AsUnionTypeNode().Types
	if typesList == nil {
		return nil
	}
	var values []string
	for _, typeNode := range typesList.Nodes {
		if typeNode != current && ast.IsLiteralTypeNode(typeNode) &&
			ast.IsStringLiteral(typeNode.AsLiteralTypeNode().Literal) {
			values = append(values, typeNode.AsLiteralTypeNode().Literal.Text())
		}
	}
	return values
}

func hasIndexSignature(t *checker.Type, typeChecker *checker.Checker) bool {
	return typeChecker.GetStringIndexType(t) != nil || typeChecker.GetNumberIndexType(t) != nil
}
