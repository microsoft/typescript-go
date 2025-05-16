package ls

import (
	"fmt"
	"slices"
	"strings"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/checker"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/printer"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type completionsFromTypes struct {
	types           []*checker.StringLiteralType
	isNewIdentifier bool
}

type completionsFromProperties struct {
	symbols           []*ast.Symbol
	hasIndexSignature bool
}

type completionsFromPaths = []*pathCompletion

type pathCompletion struct {
	name string
	// ScriptElementKindScriptElement | ScriptElementKindDirectory | ScriptElementKindExternalModuleName
	kind      ScriptElementKind
	extension string
	textRange *core.TextRange
}

// *completionsFromTypes | *completionsFromProperties | completionsFromPaths
type stringLiteralCompletions = any

func (l *LanguageService) getStringLiteralCompletions(
	file *ast.SourceFile,
	position int,
	contextToken *ast.Node,
	compilerOptions *core.CompilerOptions,
	program *compiler.Program,
	preferences *UserPreferences,
	clientOptions *lsproto.CompletionClientCapabilities,
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
		return l.convertStringLiteralCompletions(
			entries,
			contextToken,
			file,
			position,
			program,
			compilerOptions,
			preferences,
			clientOptions,
		)
	}
	return nil
}

func (l *LanguageService) convertStringLiteralCompletions(
	completion stringLiteralCompletions,
	contextToken *ast.StringLiteralLike,
	file *ast.SourceFile,
	position int,
	program *compiler.Program,
	options *core.CompilerOptions,
	preferences *UserPreferences,
	clientOptions *lsproto.CompletionClientCapabilities,
) *lsproto.CompletionList {
	if completion == nil {
		return nil
	}

	optionalReplacementRange := l.createRangeFromStringLiteralLikeContent(file, contextToken, position)
	switch completion := completion.(type) {
	case completionsFromPaths:
		return l.convertPathCompletions(completion, file, position, clientOptions)
	case *completionsFromProperties:
		data := &completionDataData{
			symbols:                 completion.symbols,
			completionKind:          CompletionKindString,
			isNewIdentifierLocation: completion.hasIndexSignature,
			location:                file.AsNode(),
			contextToken:            contextToken,
		}
		_, items := l.getCompletionEntriesFromSymbols(
			data,
			optionalReplacementRange,
			contextToken, /*replacementToken*/
			position,
			file,
			program,
			core.ScriptTargetESNext,
			preferences,
			options,
			clientOptions,
		)
		defaultCommitCharacters := getDefaultCommitCharacters(completion.hasIndexSignature)
		itemDefaults := setCommitCharacters(clientOptions, items, &defaultCommitCharacters)
		return &lsproto.CompletionList{
			IsIncomplete: false,
			ItemDefaults: itemDefaults,
			Items:        items,
		}
	case *completionsFromTypes:
		var quoteChar printer.QuoteChar
		if contextToken.Kind == ast.KindNoSubstitutionTemplateLiteral {
			quoteChar = printer.QuoteCharBacktick
		} else if strings.HasPrefix(contextToken.Text(), "'") {
			quoteChar = printer.QuoteCharSingleQuote
		} else {
			quoteChar = printer.QuoteCharDoubleQuote
		}
		items := core.Map(completion.types, func(t *checker.StringLiteralType) *lsproto.CompletionItem {
			name := printer.EscapeString(t.AsLiteralType().Value().(string), quoteChar)
			return l.createLSPCompletionItem(
				name,
				"", /*insertText*/
				"", /*filterText*/
				SortTextLocationPriority,
				ScriptElementKindString,
				core.Set[ScriptElementKindModifier]{},
				l.getReplacementRangeForContextToken(file, contextToken, position),
				nil, /*optionalReplacementSpan*/
				nil, /*commitCharacters*/
				nil, /*labelDetails*/
				file,
				position,
				clientOptions,
				false, /*isMemberCompletion*/
				false, /*isSnippet*/
				false, /*hasAction*/
				false, /*preselect*/
				"",    /*source*/
			)
		})
		defaultCommitCharacters := getDefaultCommitCharacters(completion.isNewIdentifier)
		itemDefaults := setCommitCharacters(clientOptions, items, &defaultCommitCharacters)
		return &lsproto.CompletionList{
			IsIncomplete: false,
			ItemDefaults: itemDefaults,
			Items:        items,
		}
	default:
		panic(fmt.Sprintf("Unexpected completion type: %T", completion))
	}
}

func (l *LanguageService) convertPathCompletions(
	pathCompletions completionsFromPaths,
	file *ast.SourceFile,
	position int,
	clientOptions *lsproto.CompletionClientCapabilities,
) *lsproto.CompletionList {
	isNewIdentifierLocation := true // The user may type in a path that doesn't yet exist, creating a "new identifier" with respect to the collection of identifiers the server is aware of.
	defaultCommitCharacters := getDefaultCommitCharacters(isNewIdentifierLocation)
	items := core.Map(pathCompletions, func(pathCompletion *pathCompletion) *lsproto.CompletionItem {
		replacementSpan := l.createLspRangeFromBounds(pathCompletion.textRange.Pos(), pathCompletion.textRange.End(), file)
		return l.createLSPCompletionItem(
			pathCompletion.name,
			"", /*insertText*/
			"", /*filterText*/
			SortTextLocationPriority,
			pathCompletion.kind,
			*core.NewSetFromItems(kindModifiersFromExtension(pathCompletion.extension)),
			replacementSpan,
			nil, /*optionalReplacementSpan*/
			nil, /*commitCharacters*/
			nil, /*labelDetails*/
			file,
			position,
			clientOptions,
			false, /*isMemberCompletion*/
			false, /*isSnippet*/
			false, /*hasAction*/
			false, /*preselect*/
			"",    /*source*/
		)
	})
	itemDefaults := setCommitCharacters(clientOptions, items, &defaultCommitCharacters)
	return &lsproto.CompletionList{
		IsIncomplete: false,
		ItemDefaults: itemDefaults,
		Items:        items,
	}
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
	case ast.KindCallExpression, ast.KindNewExpression, ast.KindJsxAttribute:
		if !isRequireCallArgument(node) && !ast.IsImportCall(parent) {
			// !!! signature help
			// const argumentInfo = SignatureHelp.getArgumentInfoForCompletions(parent.kind === SyntaxKind.JsxAttribute ? parent.parent : node, position, sourceFile, typeChecker);
			// // Get string literal completions from specialized signatures of the target
			// // i.e. declare function f(a: 'A');
			// // f("/*completion position*/")
			// return argumentInfo && getStringLiteralCompletionsFromSignature(argumentInfo.invocation, node, argumentInfo, typeChecker) || fromContextualType(ContextFlags.None);
			return nil
		}
		fallthrough // is `require("")` or `require(""` or `import("")`
	case ast.KindImportDeclaration, ast.KindExportDeclaration, ast.KindExternalModuleReference, ast.KindJSDocImportTag:
		// Get all known external module names or complete a path to a module
		// i.e. import * as ns from "/*completion position*/";
		//      var y = import("/*completion position*/");
		//      import x = require("/*completion position*/");
		//      var y = require("/*completion position*/");
		//      export * from "/*completion position*/";
		return getStringLiteralCompletionsFromModuleNames(file, node, program, preferences)
	case ast.KindCaseClause:
		tracker := newCaseClauseTracker(typeChecker, parent.Parent.AsCaseBlock().Clauses.Nodes)
		contextualTypes := fromContextualType(checker.ContextFlagsCompletions, node, typeChecker)
		if contextualTypes == nil {
			return nil
		}
		literals := core.Filter(contextualTypes.types, func(t *checker.StringLiteralType) bool {
			return !tracker.hasValue(t.AsLiteralType().Value())
		})
		return &completionsFromTypes{
			types:           literals,
			isNewIdentifier: false,
		}
	case ast.KindImportSpecifier, ast.KindExportSpecifier:
		// Complete string aliases in `import { "|" } from` and `export { "|" } from`
		specifier := parent
		if propertyName := specifier.PropertyName(); propertyName != nil && node != propertyName {
			return nil // Don't complete in `export { "..." as "|" } from`
		}
		namedImportsOrExports := specifier.Parent
		var moduleSpecifier *ast.Node
		if namedImportsOrExports.Kind == ast.KindNamedImports {
			moduleSpecifier = namedImportsOrExports.Parent.Parent
		} else {
			moduleSpecifier = namedImportsOrExports.Parent
		}
		if moduleSpecifier == nil {
			return nil
		}
		moduleSpecifierSymbol := typeChecker.GetSymbolAtLocation(moduleSpecifier)
		if moduleSpecifierSymbol == nil {
			return nil
		}
		exports := typeChecker.GetExportsAndPropertiesOfModule(moduleSpecifierSymbol)
		existing := core.NewSetFromItems(core.Map(namedImportsOrExports.Elements(), func(n *ast.Node) string {
			if n.PropertyName() != nil {
				return n.PropertyName().Text()
			}
			return n.Name().Text()
		})...)
		uniques := core.Filter(exports, func(e *ast.Symbol) bool {
			return e.Name != ast.InternalSymbolNameDefault && !existing.Has(e.Name)
		})
		return &completionsFromProperties{
			symbols:           uniques,
			hasIndexSignature: false,
		}
	default:
		result := fromContextualType(checker.ContextFlagsCompletions, node, typeChecker)
		if result != nil {
			return result
		}
		return fromContextualType(checker.ContextFlagsNone, node, typeChecker)
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
	objectLiteralExpression *ast.ObjectLiteralExpressionNode,
) *completionsFromProperties {
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
	preferences *UserPreferences,
) stringLiteralCompletions {
	// !!! needs `getModeForUsageLocationWorker`
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

// Matches
//
//	require(""
//	require("")
func isRequireCallArgument(node *ast.Node) bool {
	return ast.IsCallExpression(node.Parent) && len(node.Parent.Arguments()) > 0 && node.Parent.Arguments()[0] == node &&
		ast.IsIdentifier(node.Parent.Expression()) && node.Parent.Expression().Text() == "require"
}

func kindModifiersFromExtension(extension string) ScriptElementKindModifier {
	switch extension {
	case tspath.ExtensionDts:
		return ScriptElementKindModifierDts
	case tspath.ExtensionJs:
		return ScriptElementKindModifierJs
	case tspath.ExtensionJson:
		return ScriptElementKindModifierJson
	case tspath.ExtensionJsx:
		return ScriptElementKindModifierJsx
	case tspath.ExtensionTs:
		return ScriptElementKindModifierTs
	case tspath.ExtensionTsx:
		return ScriptElementKindModifierTsx
	case tspath.ExtensionDmts:
		return ScriptElementKindModifierDmts
	case tspath.ExtensionMjs:
		return ScriptElementKindModifierMjs
	case tspath.ExtensionMts:
		return ScriptElementKindModifierMts
	case tspath.ExtensionDcts:
		return ScriptElementKindModifierDcts
	case tspath.ExtensionCjs:
		return ScriptElementKindModifierCjs
	case tspath.ExtensionCts:
		return ScriptElementKindModifierCts
	case tspath.ExtensionTsBuildInfo:
		panic(fmt.Sprintf("Extension %v is unsupported.", tspath.ExtensionTsBuildInfo))
	case "":
		return ScriptElementKindModifierNone
	default:
		panic(fmt.Sprintf("Unexpected extension: %v", extension))
	}
}
