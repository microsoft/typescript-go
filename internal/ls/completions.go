package ls

import (
	"fmt"
	"maps"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/astnav"
	"github.com/microsoft/typescript-go/internal/checker"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/scanner"
	"github.com/microsoft/typescript-go/internal/stringutil"
	"github.com/microsoft/typescript-go/internal/tspath"
)

func (l *LanguageService) ProvideCompletion(
	fileName string,
	position int,
	context *lsproto.CompletionContext,
	clientOptions *lsproto.ClientCompletionItemOptions) *lsproto.CompletionList {
	program, file := l.getProgramAndFile(fileName)
	node := astnav.GetTouchingPropertyName(file, position)
	if node.Kind == ast.KindSourceFile {
		return nil
	}
	return l.getCompletionsAtPosition(program, file, position, context, nil /*preferences*/, clientOptions) // !!! get preferences
}

// !!! figure out other kinds of completion data return
type completionData struct {
	// !!!
	symbols          []*ast.Symbol
	completionKind   CompletionKind
	isInSnippetScope bool
	// Note that the presence of this alone doesn't mean that we need a conversion. Only do that if the completion is not an ordinary identifier.
	propertyAccessToConvert      *ast.PropertyAccessExpressionNode
	isNewIdentifierLocation      bool
	location                     *ast.Node
	keywordFilters               KeywordCompletionFilters
	literals                     []any
	symbolToOriginInfoMap        map[ast.SymbolId]*symbolOriginInfo
	symbolToSortTextMap          map[ast.SymbolId]sortText
	recommendedCompletion        *ast.Symbol
	previousToken                *ast.Node
	contextToken                 *ast.Node
	jsxInitializer               jsxInitializer
	insideJsDocTagTypeExpression bool
	isTypeOnlyLocation           bool
	// In JSX tag name and attribute names, identifiers like "my-tag" or "aria-name" is valid identifier.
	isJsxIdentifierExpected   bool
	isRightOfOpenTag          bool
	isRightOfDotOrQuestionDot bool
	importStatementCompletion any  // !!!
	hasUnresolvedAutoImports  bool // !!!
	// flags CompletionInfoFlags // !!!
	defaultCommitCharacters []string
}

type importStatementCompletionInfo struct {
	// !!!
}

// If we're after the `=` sign but no identifier has been typed yet,
// value will be `true` but initializer will be `nil`.
type jsxInitializer struct {
	isInitializer bool
	initializer   *ast.IdentifierNode
}

type KeywordCompletionFilters int

const (
	KeywordCompletionFiltersNone                         KeywordCompletionFilters = iota // No keywords
	KeywordCompletionFiltersAll                                                          // Every possible kewyord
	KeywordCompletionFiltersClassElementKeywords                                         // Keywords inside class body
	KeywordCompletionFiltersInterfaceElementKeywords                                     // Keywords inside interface body
	KeywordCompletionFiltersConstructorParameterKeywords                                 // Keywords at constructor parameter
	KeywordCompletionFiltersFunctionLikeBodyKeywords                                     // Keywords at function like body
	KeywordCompletionFiltersTypeAssertionKeywords
	KeywordCompletionFiltersTypeKeywords
	KeywordCompletionFiltersTypeKeyword // Literally just `type`
	KeywordCompletionFiltersLast        = KeywordCompletionFiltersTypeKeyword
)

type CompletionKind int

const (
	CompletionKindNone CompletionKind = iota
	CompletionKindObjectPropertyDeclaration
	CompletionKindGlobal
	CompletionKindPropertyAccess
	CompletionKindMemberLike
	CompletionKindString
)

// All commit characters, valid when `isNewIdentifierLocation` is false.
var allCommitCharacters = []string{".", ",", ";"}

// Commit characters valid at expression positions where we could be inside a parameter list.
var noCommaCommitCharacters = []string{".", ";"}

type sortText string

const (
	SortTextLocalDeclarationPriority         sortText = "10"
	sortTextLocationPriority                 sortText = "11"
	sortTextOptionalMember                   sortText = "12"
	sortTextMemberDeclaredBySpreadAssignment sortText = "13"
	sortTextSuggestedClassMembers            sortText = "14"
	sortTextGlobalsOrKeywords                sortText = "15"
	sortTextAutoImportSuggestions            sortText = "16"
	sortTextClassMemberSnippets              sortText = "17"
	sortTextJavascriptIdentifiers            sortText = "18"
)

func deprecateSortText(original sortText) sortText {
	return "z" + original
}

func sortBelow(original sortText) sortText {
	return original + "1"
}

// !!! sort text transformations

type symbolOriginInfoKind int

const (
	symbolOriginInfoKindThisType symbolOriginInfoKind = 1 << iota
	symbolOriginInfoKindSymbolMember
	symbolOriginInfoKindExport
	symbolOriginInfoKindPromise
	symbolOriginInfoKindNullable
	symbolOriginInfoKindResolvedExport
	symbolOriginInfoKindTypeOnlyAlias
	symbolOriginInfoKindObjectLiteralMethod
	symbolOriginInfoKindIgnore
	symbolOriginInfoKindComputedPropertyName

	symbolOriginInfoKindSymbolMemberNoExport symbolOriginInfoKind = symbolOriginInfoKindSymbolMember
	symbolOriginInfoKindSymbolMemberExport                        = symbolOriginInfoKindSymbolMember | symbolOriginInfoKindExport
)

type symbolOriginInfo struct {
	kind              symbolOriginInfoKind
	isDefaultExport   bool
	isFromPackageJson bool
	fileName          string
	data              any
}

// type originData interface {
// 	symbolName() string
// }

// !!! origin info
func (s *symbolOriginInfo) symbolName() string {
	switch s.data.(type) {
	case *symbolOriginInfoExport:
		return s.data.(*symbolOriginInfoExport).symbolName
	case *symbolOriginInfoResolvedExport:
		return s.data.(*symbolOriginInfoResolvedExport).symbolName
	default:
		panic(fmt.Sprintf("symbolOriginInfo: unknown data type for symbolName(): %T", s.data))
	}
}

type symbolOriginInfoExport struct {
	symbolName      string
	moduleSymbol    *ast.Symbol
	isDefaultExport bool
	exporName       string
	// exportMapKey  ExportMapInfoKey // !!!
}

func (s *symbolOriginInfo) asExport() *symbolOriginInfoExport {
	return s.data.(*symbolOriginInfoExport)
}

type symbolOriginInfoResolvedExport struct {
	symbolName   string
	moduleSymbol *ast.Symbol
	exportName   string
	// exportMapKey ExportMapInfoKey // !!!
	moduleSpecifier string
}

func (s *symbolOriginInfo) asResolvedExport() *symbolOriginInfoResolvedExport {
	return s.data.(*symbolOriginInfoResolvedExport)
}

type symbolOriginInfoObjectLiteralMethod struct {
	insertText   string
	labelDetails *lsproto.CompletionItemLabelDetails
	isSnippet    bool
}

func (s *symbolOriginInfo) asObjectLiteralMethod() *symbolOriginInfoObjectLiteralMethod {
	return s.data.(*symbolOriginInfoObjectLiteralMethod)
}

// Special values for `CompletionInfo['source']` used to disambiguate
// completion items with the same `name`. (Each completion item must
// have a unique name/source combination, because those two fields
// comprise `CompletionEntryIdentifier` in `getCompletionEntryDetails`.
//
// When the completion item is an auto-import suggestion, the source
// is the module specifier of the suggestion. To avoid collisions,
// the values here should not be a module specifier we would ever
// generate for an auto-import.
type completionSource string

const (
	// Completions that require `this.` insertion text.
	completionSourceThisProperty completionSource = "ThisProperty/"
	// Auto-import that comes attached to a class member snippet.
	completionSourceClassMemberSnippet completionSource = "ClassMemberSnippet/"
	// A type-only import that needs to be promoted in order to be used at the completion location.
	completionSourceTypeOnlyAlias completionSource = "TypeOnlyAlias/"
	// Auto-import that comes attached to an object literal method snippet.
	completionSourceObjectLiteralMethodSnippet completionSource = "ObjectLiteralMethodSnippet/"
	// Case completions for switch statements.
	completionSourceSwitchCases completionSource = "SwitchCases/"
	// Completions for an object literal expression.
	completionSourceObjectLiteralMemberWithComma completionSource = "ObjectLiteralMemberWithComma/"
)

// Value is set to false for global variables or completions from external module exports,
// true otherwise.
type uniqueNameSet = map[string]bool

func (l *LanguageService) getCompletionsAtPosition(
	program *compiler.Program,
	file *ast.SourceFile,
	position int,
	context *lsproto.CompletionContext,
	preferences *UserPreferences,
	clientOptions *lsproto.ClientCompletionItemOptions,
) *lsproto.CompletionList {
	previousToken, _ := getRelevantTokens(position, file)
	if context.TriggerCharacter != nil && !isInString(file, position, previousToken) && !isValidTrigger(file, *context.TriggerCharacter, previousToken, position) {
		return nil
	}

	if *context.TriggerCharacter == " " {
		// `isValidTrigger` ensures we are at `import |`
		if preferences.includeCompletionsForImportStatements {
			// !!! isMemberCompletion
			return &lsproto.CompletionList{
				IsIncomplete: true,
				ItemDefaults: &lsproto.CompletionItemDefaults{ // !!! do we need this if no entries? also, check if client supports item defaults
					CommitCharacters: getDefaultCommitCharacters(true /*isNewIdentifierLocation*/),
				},
			}
		}
		return nil
	}

	compilerOptions := program.GetCompilerOptions()

	// !!! see if incomplete completion list and continue or clean

	// !!! string literal completions

	// !!! label completions

	completionData := getCompletionData(program, file, position, preferences)
	if completionData == nil {
		return nil
	}

	// switch completionData.Kind  // !!! other data cases
	// !!! transform data into completion list

	response := completionInfoFromData(
		file,
		program,
		compilerOptions,
		completionData,
		preferences,
		position,
		clientOptions,
	)
	// !!! check if response is incomplete
	return response
}

func getCompletionData(program *compiler.Program, file *ast.SourceFile, position int, preferences *UserPreferences) *completionData {
	typeChecker := program.GetTypeChecker()
	inCheckedFile := isCheckedFile(file, program.GetCompilerOptions())

	currentToken := astnav.GetTokenAtPosition(file, position)

	insideComment := isInComment(file, position, currentToken)

	insideJsDocTagTypeExpression := false
	insideJsDocImportTag := false
	isInSnippetScope := false
	if insideComment != nil {
		// !!! jsdoc
	}

	// The decision to provide completion depends on the contextToken, which is determined through the previousToken.
	// Note: 'previousToken' (and thus 'contextToken') can be undefined if we are the beginning of the file
	isJSOnlyLocation := !insideJsDocTagTypeExpression && !insideJsDocImportTag && ast.IsSourceFileJs(file)
	previousToken, contextToken := getRelevantTokens(position, file)

	// Find the node where completion is requested on.
	// Also determine whether we are trying to complete with members of that node
	// or attributes of a JSX tag.
	node := currentToken
	var propertyAccessToConvert *ast.PropertyAccessExpressionNode
	isRightOfDot := false
	isRightOfQuestionDot := false
	isRightOfOpenTag := false
	isStartingCloseTag := false
	var jsxInitializer jsxInitializer
	isJsxIdentifierExpected := false
	var importStatementCompletion *importStatementCompletionInfo
	location := astnav.GetTouchingPropertyName(file, position)
	keywordFilters := KeywordCompletionFiltersNone
	isNewIdentifierLocation := false
	// !!!
	// flags := CompletionInfoFlagsNone
	var defaultCommitCharacters []string

	if contextToken != nil {
		// !!! import completions

		parent := contextToken.Parent
		if contextToken.Kind == ast.KindDotToken || contextToken.Kind == ast.KindQuestionDotToken {
			isRightOfDot = contextToken.Kind == ast.KindDotToken
			isRightOfQuestionDot = contextToken.Kind == ast.KindQuestionDotToken
			switch parent.Kind {
			case ast.KindPropertyAccessExpression:
				propertyAccessToConvert = parent
				node = propertyAccessToConvert.Expression()
				leftMostAccessExpression := ast.GetLeftmostAccessExpression(parent)
				if ast.NodeIsMissing(leftMostAccessExpression) ||
					((ast.IsCallExpression(node) || ast.IsFunctionLike(node)) &&
						node.End() == contextToken.Pos() &&
						getLastChild(node, file).Kind != ast.KindCloseParenToken) {
					// This is likely dot from incorrectly parsed expression and user is starting to write spread
					// eg: Math.min(./**/)
					// const x = function (./**/) {}
					// ({./**/})
					return nil
				}
			case ast.KindQualifiedName:
				node = parent.AsQualifiedName().Left
			case ast.KindModuleDeclaration:
				node = parent.Name()
			case ast.KindImportType:
				node = parent
			case ast.KindMetaProperty:
				node = getFirstToken(parent, file)
				if node.Kind != ast.KindImportKeyword || node.Kind != ast.KindNewKeyword {
					panic("Unexpected token kind: " + node.Kind.String())
				}
			default:
				// There is nothing that precedes the dot, so this likely just a stray character
				// or leading into a '...' token. Just bail out instead.
				return nil
			}
		} else { // !!! else if (!importStatementCompletion)
			// <UI.Test /* completion position */ />
			// If the tagname is a property access expression, we will then walk up to the top most of property access expression.
			// Then, try to get a JSX container and its associated attributes type.
			if parent != nil && parent.Kind == ast.KindPropertyAccessExpression {
				contextToken = parent
				parent = parent.Parent
			}

			// Fix location
			if parent == location {
				switch currentToken.Kind {
				case ast.KindGreaterThanToken:
					if parent.Kind == ast.KindJsxElement || parent.Kind == ast.KindJsxOpeningElement {
						location = currentToken
					}
				case ast.KindSlashToken:
					if parent.Kind == ast.KindJsxSelfClosingElement {
						location = currentToken
					}
				}
			}

			switch parent.Kind {
			case ast.KindJsxClosingElement:
				if contextToken.Kind == ast.KindSlashToken {
					isStartingCloseTag = true
					location = contextToken
				}
			case ast.KindBinaryExpression:
				if !binaryExpressionMayBeOpenTag(parent.AsBinaryExpression()) {
					break
				}
				fallthrough
			case ast.KindJsxSelfClosingElement, ast.KindJsxElement, ast.KindJsxOpeningElement:
				isJsxIdentifierExpected = true
				if contextToken.Kind == ast.KindLessThanToken {
					isRightOfOpenTag = true
					location = contextToken
				}
			case ast.KindJsxExpression, ast.KindJsxSpreadAttribute:
				// First case is for `<div foo={true} [||] />` or `<div foo={true} [||] ></div>`,
				// `parent` will be `{true}` and `previousToken` will be `}`.
				// Second case is for `<div foo={true} t[||] ></div>`.
				// Second case must not match for `<div foo={undefine[||]}></div>`.
				if previousToken.Kind == ast.KindCloseBraceToken ||
					previousToken.Kind == ast.KindIdentifier && previousToken.Parent.Kind == ast.KindJsxAttribute {
					isJsxIdentifierExpected = true
				}
			case ast.KindJsxAttribute:
				// For `<div className="x" [||] ></div>`, `parent` will be JsxAttribute and `previousToken` will be its initializer.
				if parent.Initializer() == previousToken && previousToken.End() < position {
					isJsxIdentifierExpected = true
				} else {
					switch previousToken.Kind {
					case ast.KindEqualsToken:
						jsxInitializer.isInitializer = true
					case ast.KindIdentifier:
						isJsxIdentifierExpected = true
						// For `<div x=[|f/**/|]`, `parent` will be `x` and `previousToken.parent` will be `f` (which is its own JsxAttribute).
						// Note for `<div someBool f>` we don't want to treat this as a jsx inializer, instead it's the attribute name.
						if parent != previousToken.Parent &&
							parent.Initializer() == nil &&
							findChildOfKind(parent, ast.KindEqualsToken, file) != nil {
							jsxInitializer.initializer = previousToken
						}
					}
				}
			}
		}
	}

	completionKind := CompletionKindNone
	hasUnresolvedAutoImports := false
	// This also gets mutated in nested-functions after the return
	var symbols []*ast.Symbol
	var symbolToOriginInfoMap map[ast.SymbolId]*symbolOriginInfo
	var symbolToSortTextMap map[ast.SymbolId]sortText
	var importSpecifierResolver any // !!! auto import
	var seenPropertySymbols core.Set[ast.SymbolId]
	isTypeOnlyLocation := insideJsDocTagTypeExpression || insideJsDocImportTag ||
		importStatementCompletion != nil && ast.IsTypeOnlyImportOrExportDeclaration(location.Parent) ||
		!isContextTokenValueLocation(contextToken) &&
			(isPossiblyTypeArgumentPosition(contextToken, file, typeChecker) ||
				ast.IsPartOfTypeNode(location) ||
				isContextTokenTypeLocation(contextToken))
	// var getModuleSpecifierResolutionHost any // !!! auto import

	addSymbolOriginInfo := func(symbol *ast.Symbol, insertQuestionDot bool, insertAwait bool) {
		symbolId := ast.GetSymbolId(symbol)
		if insertAwait && core.AddIfAbsent(seenPropertySymbols, symbolId) {
			symbolToOriginInfoMap[symbolId] = &symbolOriginInfo{kind: getNullableSymbolOriginInfoKind(symbolOriginInfoKindPromise, insertQuestionDot)}
		} else if insertQuestionDot {
			symbolToOriginInfoMap[symbolId] = &symbolOriginInfo{kind: symbolOriginInfoKindNullable}
		}
	}

	addSymbolSortInfo := func(symbol *ast.Symbol) {
		symbolId := ast.GetSymbolId(symbol)
		if isStaticProperty(symbol) {
			symbolToSortTextMap[symbolId] = SortTextLocalDeclarationPriority
		}
	}

	addPropertySymbol := func(symbol *ast.Symbol, insertAwait bool, insertQuestionDot bool) {
		// For a computed property with an accessible name like `Symbol.iterator`,
		// we'll add a completion for the *name* `Symbol` instead of for the property.
		// If this is e.g. [Symbol.iterator], add a completion for `Symbol`.
		computedPropertyName := core.FirstNonNil(symbol.Declarations, func(decl *ast.Node) *ast.Node {
			name := ast.GetNameOfDeclaration(decl)
			if name != nil && name.Kind == ast.KindComputedPropertyName {
				return name
			}
			return nil
		})

		if computedPropertyName != nil {
			leftMostName := getLeftMostName(computedPropertyName.Expression()) // The completion is for `Symbol`, not `iterator`.
			var nameSymbol *ast.Symbol
			if leftMostName != nil {
				nameSymbol = typeChecker.GetSymbolAtLocation(leftMostName)
			}
			// If this is nested like for `namespace N { export const sym = Symbol(); }`, we'll add the completion for `N`.
			var firstAccessibleSymbol *ast.Symbol
			if nameSymbol != nil {
				firstAccessibleSymbol = getFirstSymbolInChain(nameSymbol, contextToken, typeChecker)
			}
			var firstAccessibleSymbolId ast.SymbolId
			if firstAccessibleSymbol != nil {
				firstAccessibleSymbolId = ast.GetSymbolId(firstAccessibleSymbol)
			}
			if firstAccessibleSymbolId != 0 && core.AddIfAbsent(seenPropertySymbols, firstAccessibleSymbolId) {
				index := len(symbols)
				symbols = append(symbols, firstAccessibleSymbol)
				moduleSymbol := firstAccessibleSymbol.Parent
				if moduleSymbol == nil ||
					!checker.IsExternalModuleSymbol(moduleSymbol) ||
					typeChecker.TryGetMemberInModuleExportsAndProperties(firstAccessibleSymbol.Name, moduleSymbol) != firstAccessibleSymbol {
					symbolToOriginInfoMap[ast.GetSymbolId(symbol)] = &symbolOriginInfo{kind: getNullableSymbolOriginInfoKind(symbolOriginInfoKindSymbolMemberNoExport, insertQuestionDot)}
				} else {
					var fileName string
					if tspath.IsExternalModuleNameRelative(core.StripQuotes(moduleSymbol.Name)) {
						fileName = ast.GetSourceFileOfModule(moduleSymbol).FileName()
					}
					if importSpecifierResolver == nil { // !!! verify if this is right, depending on the type of importSpecifierResolver
						// !!!
						// importSpecifierResolver ||= codefix.createImportSpecifierResolver(sourceFile, program, host, preferences))
					}
					// !!!
					// const { moduleSpecifier } = importSpecifier.getModuleSpecifierForBestExportInfo(
					// 	[{
					// 		exportKind: ExportKind.Named,
					// 		moduleFileName: fileName,
					// 		isFromPackageJson: false,
					// 		moduleSymbol,
					// 		symbol: firstAccessibleSymbol,
					// 		targetFlags: skipAlias(firstAccessibleSymbol, typeChecker).flags,
					// 	}],
					// 	position,
					// 	isValidTypeOnlyAliasUseSite(location),
					// ) || {};
					// if (moduleSpecifier) {
					// 	const origin: SymbolOriginInfoResolvedExport = {
					// 		kind: getNullableSymbolOriginInfoKind(SymbolOriginInfoKind.SymbolMemberExport),
					// 		moduleSymbol,
					// 		isDefaultExport: false,
					// 		symbolName: firstAccessibleSymbol.name,
					// 		exportName: firstAccessibleSymbol.name,
					// 		fileName,
					// 		moduleSpecifier,
					// 	};
					// 	symbolToOriginInfoMap[index] = origin;
					// }
				}
			} else if firstAccessibleSymbolId == 0 || !seenPropertySymbols.Has(firstAccessibleSymbolId) {
				symbols = append(symbols, symbol)
				addSymbolOriginInfo(symbol, insertQuestionDot, insertAwait)
				addSymbolSortInfo(symbol)
			}
		} else {
			symbols = append(symbols, symbol)
			addSymbolOriginInfo(symbol, insertQuestionDot, insertAwait)
			addSymbolSortInfo(symbol)
		}
	}

	addTypeProperties := func(t *checker.Type, insertAwait bool, insertQuestionDot bool) {
		if typeChecker.GetStringIndexType(t) != nil {
			isNewIdentifierLocation = true
			defaultCommitCharacters = make([]string, 0)
		}
		if isRightOfQuestionDot && len(typeChecker.GetCallSignatures(t)) != 0 {
			isNewIdentifierLocation = true
			if len(defaultCommitCharacters) == 0 {
				defaultCommitCharacters = slices.Clone(allCommitCharacters) // Only invalid commit character here would be `(`.
			}
		}

		var propertyAccess *ast.Node
		if node.Kind == ast.KindImportType {
			propertyAccess = node
		} else {
			propertyAccess = node.Parent
		}

		if inCheckedFile {
			for _, symbol := range typeChecker.GetApparentProperties(t) {
				if typeChecker.IsValidPropertyAccessForCompletions(propertyAccess, t, symbol) {
					addPropertySymbol(symbol, false /*insertAwait*/, insertQuestionDot)
				}
			}
		} else {
			// In javascript files, for union types, we don't just get the members that
			// the individual types have in common, we also include all the members that
			// each individual type has. This is because we're going to add all identifiers
			// anyways. So we might as well elevate the members that were at least part
			// of the individual types to a higher status since we know what they are.
			for _, symbol := range getPropertiesForCompletion(t, typeChecker) {
				if typeChecker.IsValidPropertyAccessForCompletions(propertyAccess, t, symbol) {
					symbols = append(symbols, symbol)
				}
			}
		}

		if insertAwait {
			promiseType := typeChecker.GetPromisedTypeOfPromise(t)
			if promiseType != nil {
				for _, symbol := range typeChecker.GetApparentProperties(promiseType) {
					if typeChecker.IsValidPropertyAccessForCompletions(propertyAccess, promiseType, symbol) {
						addPropertySymbol(symbol, true /*insertAwait*/, insertQuestionDot)
					}
				}
			}
		}
	}

	getTypeScriptMemberSymbols := func() {
		// Right of dot member completion list
		completionKind = CompletionKindPropertyAccess

		// Since this is qualified name check it's a type node location
		isImportType := ast.IsLiteralImportTypeNode(node)
		isTypeLocation := (isImportType && !node.AsImportTypeNode().IsTypeOf) ||
			ast.IsPartOfTypeNode(node.Parent) ||
			isPossiblyTypeArgumentPosition(contextToken, file, typeChecker)
		isRhsOfImportDeclaration := isInRightSideOfInternalImportEqualsDeclaration(node)
		if ast.IsEntityName(node) || isImportType || ast.IsPropertyAccessExpression(node) {
			isNamespaceName := ast.IsModuleDeclaration(node.Parent)
			if isNamespaceName {
				isNewIdentifierLocation = true
				defaultCommitCharacters = make([]string, 0)
			}
			symbol := typeChecker.GetSymbolAtLocation(node)
			if symbol != nil {
				symbol := checker.SkipAlias(symbol, typeChecker)
				if symbol.Flags&(ast.SymbolFlagsModule|ast.SymbolFlagsEnum) != 0 {
					var valueAccessNode *ast.Node
					if isImportType {
						valueAccessNode = node
					} else {
						valueAccessNode = node.Parent
					}
					// Extract module or enum members
					exportedSymbols := typeChecker.GetExportsOfModule(symbol)
					for _, exportedSymbol := range exportedSymbols {
						if exportedSymbol == nil {
							panic("getExporsOfModule() should all be defined")
						}
						isValidValueAccess := func(s *ast.Symbol) bool {
							return typeChecker.IsValidPropertyAccess(valueAccessNode, s.Name)
						}
						isValidTypeAccess := func(s *ast.Symbol) bool {
							return symbolCanBeReferencedAtTypeLocation(s, typeChecker, core.Set[ast.SymbolId]{})
						}
						var isValidAccess bool
						if isNamespaceName {
							// At `namespace N.M/**/`, if this is the only declaration of `M`, don't include `M` as a completion.
							isValidAccess = exportedSymbol.Flags&ast.SymbolFlagsNamespace != 0 &&
								!core.Every(exportedSymbol.Declarations, func(declaration *ast.Declaration) bool {
									return declaration.Parent == node.Parent
								})
						} else if isRhsOfImportDeclaration {
							// Any kind is allowed when dotting off namespace in internal import equals declaration
							isValidAccess = isValidTypeAccess(exportedSymbol) || isValidValueAccess(exportedSymbol)
						} else if isTypeLocation || insideJsDocTagTypeExpression {
							isValidAccess = isValidTypeAccess(exportedSymbol)
						} else {
							isValidAccess = isValidValueAccess(exportedSymbol)
						}
						if isValidAccess {
							symbols = append(symbols, exportedSymbol)
						}
					}

					// If the module is merged with a value, we must get the type of the class and add its properties (for inherited static methods).
					if !isTypeLocation && !insideJsDocTagTypeExpression &&
						core.Some(
							symbol.Declarations,
							func(decl *ast.Declaration) bool {
								return decl.Kind != ast.KindSourceFile && decl.Kind != ast.KindModuleDeclaration && decl.Kind != ast.KindEnumDeclaration
							}) {
						t := typeChecker.GetNonOptionalType(typeChecker.GetTypeOfSymbolAtLocation(symbol, node))
						insertQuestionDot := false
						if typeChecker.IsNullableType(t) {
							canCorrectToQuestionDot := isRightOfDot && !isRightOfQuestionDot &&
								preferences.includeAutomaticOptionalChainCompletions
							if canCorrectToQuestionDot || isRightOfQuestionDot {
								t = typeChecker.GetNonNullableType(t)
								if canCorrectToQuestionDot {
									insertQuestionDot = true
								}
							}
						}
						addTypeProperties(t, node.Flags&ast.NodeFlagsAwaitContext != 0, insertQuestionDot)
					}
				}
			}
		}
	}

	if isRightOfDot || isRightOfQuestionDot {
		getTypeScriptMemberSymbols()
	} else if isRightOfOpenTag {
		// !!! jsx completions
	} else if isStartingCloseTag {
		// !!! jsx completions
	} else {
		// For JavaScript or TypeScript, if we're not after a dot, then just try to get the
		// global symbols in scope.  These results should be valid for either language as
		// the set of symbols that can be referenced from this location.
		// !!! global completions
	}

	var contextualType *checker.Type
	if previousToken != nil {
		contextualType = getContextualType(previousToken, position, file, typeChecker)
	}

	// exclude literal suggestions after <input type="text" [||] /> (#51667) and after closing quote (#52675)
	// for strings getStringLiteralCompletions handles completions
	isLiteralExpected := !ast.IsStringLiteralLike(previousToken) && !isJsxIdentifierExpected
	var literals []any
	if isLiteralExpected {
		var types []*checker.Type
		if contextualType != nil && checker.IsUnion(contextualType) {
			types = contextualType.Types()
		} else if contextualType != nil {
			types = []*checker.Type{contextualType}
		}
		literals = core.MapNonNil(types, func(t *checker.Type) any {
			if isLiteral(t) && !checker.IsEnumLiteral(t) {
				return t.AsLiteralType().Value()
			}
			return nil
		})
	}

	var recommendedCompletion *ast.Symbol
	if previousToken != nil && contextualType != nil {
		recommendedCompletion = getRecommendedCompletion(previousToken, contextualType, typeChecker)
	}

	return &completionData{
		symbols:                      symbols,
		completionKind:               completionKind,
		isInSnippetScope:             isInSnippetScope,
		propertyAccessToConvert:      propertyAccessToConvert,
		isNewIdentifierLocation:      isNewIdentifierLocation,
		location:                     location,
		keywordFilters:               keywordFilters,
		literals:                     literals,
		symbolToOriginInfoMap:        symbolToOriginInfoMap,
		symbolToSortTextMap:          symbolToSortTextMap,
		recommendedCompletion:        recommendedCompletion,
		previousToken:                previousToken,
		contextToken:                 contextToken,
		jsxInitializer:               jsxInitializer,
		insideJsDocTagTypeExpression: insideJsDocTagTypeExpression,
		isTypeOnlyLocation:           isTypeOnlyLocation,
		isJsxIdentifierExpected:      isJsxIdentifierExpected,
		isRightOfOpenTag:             isRightOfOpenTag,
		isRightOfDotOrQuestionDot:    isRightOfDot || isRightOfQuestionDot,
		importStatementCompletion:    importStatementCompletion,
		hasUnresolvedAutoImports:     hasUnresolvedAutoImports,
		defaultCommitCharacters:      defaultCommitCharacters,
	}
}

func completionInfoFromData(
	file *ast.SourceFile,
	program *compiler.Program,
	compilerOptions *core.CompilerOptions,
	data *completionData,
	preferences *UserPreferences,
	position int,
	clientOptions *lsproto.ClientCompletionItemOptions,
) *lsproto.CompletionList {
	keywordFilters := data.keywordFilters
	symbols := data.symbols
	isNewIdentifierLocation := data.isNewIdentifierLocation
	contextToken := data.contextToken
	literals := data.literals
	typeChecker := program.GetTypeChecker()

	// Verify if the file is JSX language variant
	if ast.GetLanguageVariant(file.ScriptKind) == core.LanguageVariantJSX {
		// !!! jsx
		return nil
	}

	// When the completion is for the expression of a case clause (e.g. `case |`),
	// filter literals & enum symbols whose values are already present in existing case clauses.
	caseClause := ast.FindAncestor(contextToken, ast.IsCaseClause)
	if caseClause != nil &&
		(contextToken.Kind == ast.KindCaseKeyword ||
			ast.IsNodeDescendantOf(contextToken, caseClause.Expression())) {
		// !!! switch completions
	}

	isChecked := isCheckedFile(file, compilerOptions)
	if isChecked && !isNewIdentifierLocation && len(symbols) == 0 && keywordFilters == KeywordCompletionFiltersNone {
		return nil
	}

	uniqueNames, sortedEntries := getCompletionEntriesFromSymbols(
		data,
		nil, /*replacementToken*/
		position,
		file,
		program,
		compilerOptions.GetEmitScriptTarget(),
		preferences,
		compilerOptions,
		clientOptions,
	)

	if data.keywordFilters != KeywordCompletionFiltersNone {
		keywordCompletions := getKeywordCompletions(
			data.keywordFilters,
			!data.insideJsDocTagTypeExpression && ast.IsSourceFileJs(sourceFile))
		for _, keywordEntry := keywordCompletions {
			if data.isTypeOnlyLocation && isTypeKeyword(scanner.StringToToken(keywordEntry.name)) ||
				false { // !!! HERE HERE

			}
		}	
	}

	
	// !!! exhaustive case completions

	// !!! here
	return nil
}

func getCompletionEntriesFromSymbols(
	data *completionData,
	replacementToken *ast.Node,
	position int,
	file *ast.SourceFile,
	program *compiler.Program,
	target core.ScriptTarget,
	preferences *UserPreferences,
	compilerOptions *core.CompilerOptions,
	clientOptions *lsproto.ClientCompletionItemOptions,
) (uniqueNames *core.Set[string], sortedEntries []*lsproto.CompletionItem) {
	closestSymbolDeclaration := getClosestSymbolDeclaration(data.contextToken, data.location)
	useSemicolons := probablyUsesSemicolons(file)
	typeChecker := program.GetTypeChecker()
	// Tracks unique names.
	// Value is set to false for global variables or completions from external module exports, because we can have multiple of those;
	// true otherwise. Based on the order we add things we will always see locals first, then globals, then module exports.
	// So adding a completion for a local will prevent us from adding completions for external module exports sharing the same name.
	uniques := make(uniqueNameSet)
	for _, symbol := range data.symbols {
		symbolId := ast.GetSymbolId(symbol)
		origin := data.symbolToOriginInfoMap[symbolId]
		name, needsConvertPropertyAccess := getCompletionEntryDisplayNameForSymbol(
			symbol,
			target,
			origin,
			data.completionKind,
			data.isJsxIdentifierExpected,
		)
		if name == "" ||
			uniques[name] && (origin == nil || !originIsObjectLiteralMethod(origin)) ||
			data.completionKind == CompletionKindGlobal &&
				!shouldIncludeSymbol(symbol, data, closestSymbolDeclaration, file, typeChecker, compilerOptions) {
			continue
		}

		// When in a value location in a JS file, ignore symbols that definitely seem to be type-only.
		if !data.isTypeOnlyLocation && ast.IsSourceFileJs(file) && symbolAppearsToBeTypeOnly(symbol, typeChecker) {
			continue
		}

		originalSortText := data.symbolToSortTextMap[ast.GetSymbolId(symbol)]
		if originalSortText == "" {
			originalSortText = sortTextLocationPriority
		}
		sortText := core.IfElse(isDeprecated(symbol, typeChecker), deprecateSortText(originalSortText), originalSortText)
		entry := createCompletionItem(
			symbol,
			sortText,
			replacementToken,
			data.contextToken,
			data.location,
			position,
			file,
			program,
			name,
			needsConvertPropertyAccess,
			origin,
			data.recommendedCompletion,
			data.propertyAccessToConvert,
			data.jsxInitializer,
			data.importStatementCompletion,
			useSemicolons,
			compilerOptions,
			preferences,
			clientOptions,
			data.completionKind,
			data.isJsxIdentifierExpected,
			data.isRightOfOpenTag,
		)
		if entry == nil {
			continue
		}

		/** True for locals; false for globals, module exports from other files, `this.` completions. */
		shouldShadowLaterSymbols := (origin == nil || originIsTypeOnlyAlias(origin)) &&
			!(symbol.Parent == nil &&
				!core.Some(symbol.Declarations, func(d *ast.Node) bool { return ast.GetSourceFileOfNode(d) == file }))
		uniques[name] = shouldShadowLaterSymbols
		core.InsertSorted(sortedEntries, entry, compareCompletionEntries)
	}

	uniqueSet := core.NewSetWithSizeHint[string](len(uniques))
	for name := range maps.Keys(uniques) {
		uniqueSet.Add(name)
	}
	return uniqueSet, sortedEntries
}

func createCompletionItem(
	symbol *ast.Symbol,
	sortText sortText,
	replacementToken *ast.Node,
	contextToken *ast.Node,
	location *ast.Node,
	position int,
	file *ast.SourceFile,
	program *compiler.Program,
	name string,
	needsConvertPropertyAccess bool,
	origin *symbolOriginInfo,
	recommendedCompletion *ast.Symbol,
	propertyAccessToConvert *ast.Node,
	jsxInitializer jsxInitializer,
	importStatementCompletion any,
	useSemicolons bool,
	compilerOptions *core.CompilerOptions,
	preferences *UserPreferences,
	clientOptions *lsproto.ClientCompletionItemOptions,
	completionKind CompletionKind,
	isJsxIdentifierExpected bool,
	isRightOfOpenTag bool,
) *lsproto.CompletionItem {
	var insertText string
	var filterText string
	replacementSpan := getReplacementRangeForContextToken(file, replacementToken, position)
	var isSnippet, hasAction bool
	source := getSourceFromOrigin(origin)
	var labelDetails *lsproto.CompletionItemLabelDetails

	typeChecker := program.GetTypeChecker()
	insertQuestionDot := originIsNullableMember(origin)
	useBraces := originIsSymbolMember(origin) || needsConvertPropertyAccess
	if originIsThisType(origin) {
		if needsConvertPropertyAccess {
			insertText = fmt.Sprintf(
				"this%s[%s]",
				core.IfElse(insertQuestionDot, "?.", ""),
				quotePropertyName(file, preferences, name))
		} else {
			insertText = fmt.Sprintf(
				"this%s%s",
				core.IfElse(insertQuestionDot, "?.", ""),
				name)
		}
	} else if propertyAccessToConvert != nil && (useBraces || insertQuestionDot) {
		// We should only have needsConvertPropertyAccess if there's a property access to convert. But see #21790.
		// Somehow there was a global with a non-identifier name. Hopefully someone will complain about getting a "foo bar" global completion and provide a repro.
		if useBraces {
			if needsConvertPropertyAccess {
				insertText = fmt.Sprintf("[%s]", quotePropertyName(file, preferences, name))
			} else {
				insertText = fmt.Sprintf("[%s]", name)
			}
		} else {
			insertText = name
		}

		if insertQuestionDot || propertyAccessToConvert.AsPropertyAccessExpression().QuestionDotToken != nil {
			insertText = fmt.Sprintf("?.%s", insertText)
		}

		dot := findChildOfKind(propertyAccessToConvert, ast.KindDotToken, file)
		if dot == nil {
			dot = findChildOfKind(propertyAccessToConvert, ast.KindQuestionDotToken, file)
		}

		if dot == nil {
			return nil
		}

		// If the text after the '.' starts with this name, write over it. Else, add new text.
		var end int
		if strings.HasPrefix(name, propertyAccessToConvert.Name().Text()) {
			end = propertyAccessToConvert.End()
		} else {
			end = dot.End()
		}
		replacementSpan = createLspRangeFromBounds(getStartOfNode(dot, file), end, file)
	}

	if jsxInitializer.isInitializer {
		if insertText == "" {
			insertText = name
		}
		insertText = fmt.Sprintf("{%s}", insertText)
		if jsxInitializer.initializer != nil {
			replacementSpan = createLspRangeFromNode(jsxInitializer.initializer, file)
		}
	}

	if originIsPromise(origin) && propertyAccessToConvert != nil {
		if insertText == "" {
			insertText = name
		}
		precedingToken := astnav.FindPrecedingToken(file, propertyAccessToConvert.Pos())
		var awaitText string
		if precedingToken != nil && positionIsASICandidate(precedingToken.End(), precedingToken.Parent, file) {
			awaitText = ";"
		}

		awaitText += "(await " + scanner.GetTextOfNode(propertyAccessToConvert.Expression()) + ")"
		if needsConvertPropertyAccess {
			insertText = awaitText + insertText
		} else {
			dotStr := core.IfElse(insertQuestionDot, "?.", ".")
			insertText = awaitText + dotStr + insertText
		}
		isInAwaitExpression := ast.IsAwaitExpression(propertyAccessToConvert.Parent)
		wrapNode := core.IfElse(isInAwaitExpression, propertyAccessToConvert.Parent, propertyAccessToConvert.Expression())
		replacementSpan = createLspRangeFromBounds(getStartOfNode(wrapNode, file), propertyAccessToConvert.End(), file)
	}

	if originIsResolvedExport(origin) {
		labelDetails = &lsproto.CompletionItemLabelDetails{
			Description: &origin.asResolvedExport().moduleSpecifier, // !!! vscode @link support
		}
		if importStatementCompletion != nil {
			// !!! auto-imports
		}
	}

	if originIsTypeOnlyAlias(origin) {
		hasAction = true
	}

	// Provide object member completions when missing commas, and insert missing commas.
	// For example:
	//
	//    interface I {
	//        a: string;
	//        b: number
	//     }
	//
	//     const cc: I = { a: "red" | }
	//
	// Completion should add a comma after "red" and provide completions for b
	if completionKind == CompletionKindObjectPropertyDeclaration &&
		contextToken != nil &&
		!ast.NodeHasKind(astnav.FindPrecedingTokenEx(file, contextToken.Pos(), contextToken), ast.KindCommaToken) {
		if ast.IsMethodDeclaration(contextToken.Parent.Parent) ||
			ast.IsGetAccessorDeclaration(contextToken.Parent.Parent) ||
			ast.IsSetAccessorDeclaration(contextToken.Parent.Parent) ||
			ast.IsSpreadAssignment(contextToken.Parent) ||
			getLastToken(ast.FindAncestor(contextToken.Parent, ast.IsPropertyAssignment), file) == contextToken ||
			ast.IsShorthandPropertyAssignment(contextToken.Parent) &&
				getLineOfPosition(file, contextToken.End()) != getLineOfPosition(file, position) {
			source = string(completionSourceObjectLiteralMemberWithComma)
			hasAction = true
		}
	}

	if preferences.includeCompletionsWithClassMemberSnippets &&
		completionKind == CompletionKindMemberLike &&
		isClassLikeMemberCompletion(symbol, location, file) {
		// !!! class member completions
	}

	if originIsObjectLiteralMethod(origin) {
		insertText = origin.asObjectLiteralMethod().insertText
		isSnippet = origin.asObjectLiteralMethod().isSnippet
		labelDetails = origin.asObjectLiteralMethod().labelDetails // !!! check if this can conflict with case above where we set label details
		if !ptrIsTrue(clientOptions.LabelDetailsSupport) {
			name = name + *origin.asObjectLiteralMethod().labelDetails.Detail
			labelDetails = nil
		}
		source = string(completionSourceObjectLiteralMethodSnippet)
		sortText = sortBelow(sortText)
	}

	if isJsxIdentifierExpected &&
		!isRightOfOpenTag &&
		ptrIsTrue(clientOptions.SnippetSupport) &&
		preferences.jsxAttributeCompletionStyle != JsxAttributeCompletionStyleNone &&
		!(ast.IsJsxAttribute(location.Parent) && location.Parent.Initializer() != nil) {
		useBraces := preferences.jsxAttributeCompletionStyle == JsxAttributeCompletionStyleBraces
		t := typeChecker.GetTypeOfSymbolAtLocation(symbol, location)

		// If is boolean like or undefined, don't return a snippet, we want to return just the completion.
		if preferences.jsxAttributeCompletionStyle == JsxAttributeCompletionStyleAuto &&
			t.Flags()&checker.TypeFlagsBooleanLike == 0 &&
			!(t.Flags()&checker.TypeFlagsUnion != 0 && core.Some(t.Types(), func(t *checker.Type) bool { return t.Flags()&checker.TypeFlagsBooleanLike != 0 })) {
			if t.Flags()&checker.TypeFlagsStringLike != 0 ||
				t.Flags()&checker.TypeFlagsUnion != 0 &&
					core.Every(
						t.Types(),
						func(t *checker.Type) bool {
							return t.Flags()&(checker.TypeFlagsStringLike|checker.TypeFlagsUndefined) != 0 ||
								isStringAndEmptyAnonymousObjectIntersection(typeChecker, t)
						}) {
				// If type is string-like or undefined, use quotes.
				insertText = fmt.Sprintf("%s=%s", escapeSnippetText(name), quote(file, preferences, "$1"))
				isSnippet = true
			} else {
				// Use braces for everything else.
				useBraces = true
			}
		}

		if useBraces {
			insertText = fmt.Sprintf("%s={$1}", escapeSnippetText(name))
			isSnippet = true
		}
	}

	if originIsExport(origin) || originIsResolvedExport(origin) {
		// !!! auto-imports
		// data = originToCompletionEntryData(origin)
		// hasAction = importStatementCompletion == nil
	}

	parentNamedImportOrExport := ast.FindAncestor(location, isNamedImportsOrExports)
	if parentNamedImportOrExport != nil {
		languageVersion := compilerOptions.GetEmitScriptTarget()
		if !scanner.IsIdentifierText(name, languageVersion) {
			insertText = quotePropertyName(file, preferences, name)

			if parentNamedImportOrExport.Kind == ast.KindNamedImports {
				// Check if it is `import { ^here as name } from '...'``.
				// We have to access the scanner here to check if it is `{ ^here as name }`` or `{ ^here, as, name }`.
				scanner := scanner.NewScanner()
				scanner.SetText(file.Text)
				scanner.ResetPos(position)
				if !(scanner.Scan() == ast.KindAsKeyword && scanner.Scan() == ast.KindIdentifier) {
					insertText += " as " + generateIdentifierForArbitraryString(name, languageVersion)
				}
			}
		} else if parentNamedImportOrExport.Kind == ast.KindNamedImports {
			possibleToken := scanner.StringToToken(name)
			if possibleToken != ast.KindUnknown &&
				(possibleToken == ast.KindAwaitKeyword || isNonContextualKeyword(possibleToken)) {
				insertText = fmt.Sprintf("%s as %s_")
			}
		}
	}

	elementKind := getSymbolKind(typeChecker, symbol, location)
	kind := getCompletionsSymbolKind(elementKind)
	var commitCharacters []string
	if elementKind == ScriptElementKindWarning || elementKind == ScriptElementKindString {
		commitCharacters = []string{}
	} else {
		commitCharacters = nil // Use the completion list default.
	}

	kindModifiers := getSymbolModifiers(typeChecker, symbol)
	var tags *[]lsproto.CompletionItemTag
	var detail *string
	// Copied from vscode ts extension.
	if kindModifiers.Has(ScriptElementKindModifierOptional) {
		if insertText == "" {
			insertText = name
		}
		if filterText == "" {
			filterText = name
		}
		name = name + "?"
	}
	if kindModifiers.Has(ScriptElementKindModifierDeprecated) {
		tags = &[]lsproto.CompletionItemTag{lsproto.CompletionItemTagDeprecated}
	}
	if kind == lsproto.CompletionItemKindFile {
		for _, extensionModifier := range fileExtensionKindModifiers {
			if kindModifiers.Has(extensionModifier) {
				if strings.HasSuffix(name, string(extensionModifier)) {
					detail = ptrTo(name)
				} else {
					detail = ptrTo(name + string(extensionModifier))
				}
				break
			}
		}
	}

	if hasAction && source != "" {
		// !!! adjust label like vscode does
	}

	var insertTextFormat *lsproto.InsertTextFormat
	if isSnippet {
		insertTextFormat = ptrTo(lsproto.InsertTextFormatSnippet)
	} else {
		insertTextFormat = ptrTo(lsproto.InsertTextFormatPlainText)
	}

	var textEdit *lsproto.TextEditOrInsertReplaceEdit
	if replacementSpan != nil {
		textEdit = &lsproto.TextEditOrInsertReplaceEdit{
			TextEdit: &lsproto.TextEdit{
				NewText: core.IfElse(insertText == "", name, insertText),
				Range:   *replacementSpan,
			},
		}
	}

	return &lsproto.CompletionItem{
		Label:            name,
		LabelDetails:     labelDetails,
		Kind:             &kind,
		Tags:             tags,
		Detail:           detail,
		Preselect:        boolToPtr(isRecommendedCompletionMatch(symbol, recommendedCompletion, typeChecker)),
		SortText:         ptrTo(string(sortText)),
		FilterText:       strPtrTo(filterText),
		InsertText:       strPtrTo(insertText),
		InsertTextFormat: insertTextFormat,
		TextEdit:         textEdit,
		CommitCharacters: slicesPtrTo(commitCharacters),
		Data:             nil, // !!! auto-imports
	}
}

func isRecommendedCompletionMatch(localSymbol *ast.Symbol, recommendedCompletion *ast.Symbol, typeChecker *checker.Checker) bool {
	return localSymbol == recommendedCompletion ||
		localSymbol.Flags&ast.SymbolFlagsExportValue != 0 && typeChecker.GetExportSymbolOfSymbol(localSymbol) == recommendedCompletion
}

func strPtrTo(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}

func slicesPtrTo[T any](s []T) *[]T {
	if s == nil {
		return nil
	}
	return &s
}

func ptrTo[T any](v T) *T {
	return &v
}

func ptrIsTrue(ptr *bool) bool {
	if ptr == nil {
		return false
	}
	return *ptr
}

func boolToPtr(v bool) *bool {
	if v {
		return ptrTo(true)
	}
	return nil
}

func getLineOfPosition(file *ast.SourceFile, pos int) int {
	line, _ := scanner.GetLineAndCharacterOfPosition(file, pos)
	return line
}

func isClassLikeMemberCompletion(symbol *ast.Symbol, location *ast.Node, file *ast.SourceFile) bool {
	// !!! class member completions
	return false
}

func symbolAppearsToBeTypeOnly(symbol *ast.Symbol, typeChecker *checker.Checker) bool {
	flags := checker.GetCombinedLocalAndExportSymbolFlags(checker.SkipAlias(symbol, typeChecker))
	return flags&ast.SymbolFlagsValue == 0 &&
		(len(symbol.Declarations) == 0 || !ast.IsInJSFile(symbol.Declarations[0]) || flags&ast.SymbolFlagsType != 0)
}

func shouldIncludeSymbol(
	symbol *ast.Symbol,
	data *completionData,
	closestSymbolDeclaration *ast.Declaration,
	file *ast.SourceFile,
	typeChecker *checker.Checker,
	compilerOptions *core.CompilerOptions) bool {
	allFlags := symbol.Flags
	location := data.location
	if !ast.IsSourceFile(location) {
		// export = /**/ here we want to get all meanings, so any symbol is ok
		if ast.IsExportAssignment(location.Parent) {
			return true
		}

		// Filter out variables from their own initializers
		// `const a = /* no 'a' here */`
		if closestSymbolDeclaration != nil &&
			ast.IsVariableDeclaration(closestSymbolDeclaration) &&
			symbol.ValueDeclaration == closestSymbolDeclaration {
			return false
		}

		// Filter out current and latter parameters from defaults
		// `function f(a = /* no 'a' and 'b' here */, b) { }` or
		// `function f<T = /* no 'T' and 'T2' here */>(a: T, b: T2) { }`
		var symbolDeclaration *ast.Declaration
		if symbol.ValueDeclaration != nil {
			symbolDeclaration = symbol.ValueDeclaration
		} else if len(symbol.Declarations) > 0 {
			symbolDeclaration = symbol.Declarations[0]
		}

		if closestSymbolDeclaration != nil && symbolDeclaration != nil {
			if ast.IsParameter(closestSymbolDeclaration) && ast.IsParameter(symbolDeclaration) {
				parameters := closestSymbolDeclaration.Parent.ParameterList()
				if symbolDeclaration.Pos() >= closestSymbolDeclaration.Pos() &&
					symbolDeclaration.Pos() < parameters.End() {
					return false
				}
			} else if ast.IsTypeParameterDeclaration(closestSymbolDeclaration) &&
				ast.IsTypeParameterDeclaration(symbolDeclaration) {
				if closestSymbolDeclaration == symbolDeclaration && data.contextToken != nil && data.contextToken.Kind == ast.KindExtendsKeyword {
					// filter out the directly self-recursive type parameters
					// `type A<K extends /* no 'K' here*/> = K`
					return false
				}
				if isInTypeParameterDefault(data.contextToken) && !ast.IsInferTypeNode(closestSymbolDeclaration.Parent) {
					typeParameters := closestSymbolDeclaration.Parent.TypeParameterList()
					if typeParameters != nil && symbolDeclaration.Pos() >= closestSymbolDeclaration.Pos() &&
						symbolDeclaration.Pos() < typeParameters.End() {
						return false
					}
				}
			}
		}

		// External modules can have global export declarations that will be
		// available as global keywords in all scopes. But if the external module
		// already has an explicit export and user only wants to use explicit
		// module imports then the global keywords will be filtered out so auto
		// import suggestions will win in the completion.
		symbolOrigin := checker.SkipAlias(symbol, typeChecker)
		// We only want to filter out the global keywords.
		// Auto Imports are not available for scripts so this conditional is always false.
		if file.AsSourceFile().ExternalModuleIndicator != nil &&
			compilerOptions.AllowUmdGlobalAccess != core.TSTrue &&
			data.symbolToSortTextMap[ast.GetSymbolId(symbol)] == sortTextGlobalsOrKeywords &&
			(data.symbolToSortTextMap[ast.GetSymbolId(symbolOrigin)] == sortTextAutoImportSuggestions ||
				data.symbolToSortTextMap[ast.GetSymbolId(symbolOrigin)] == sortTextLocationPriority) {
			return false
		}

		allFlags = allFlags | checker.GetCombinedLocalAndExportSymbolFlags(symbolOrigin)

		// import m = /**/ <-- It can only access namespace (if typing import = x. this would get member symbols and not namespace)
		if isInRightSideOfInternalImportEqualsDeclaration(data.location) {
			return allFlags&ast.SymbolFlagsNamespace != 0
		}

		if data.isTypeOnlyLocation {
			// It's a type, but you can reach it by namespace.type as well.
			return symbolCanBeReferencedAtTypeLocation(symbol, typeChecker, core.Set[ast.SymbolId]{})
		}
	}

	// expressions are value space (which includes the value namespaces)
	return allFlags&ast.SymbolFlagsValue != 0
}

func getCompletionEntryDisplayNameForSymbol(
	symbol *ast.Symbol,
	target core.ScriptTarget,
	origin *symbolOriginInfo,
	completionKind CompletionKind,
	isJsxIdentifierExpected bool,
) (displayName string, needsConvertPropertyAccess bool) {
	if originIsIgnore(origin) {
		return "", false
	}

	name := core.IfElse(originIncludesSymbolName(origin), origin.symbolName(), symbol.Name)
	if name == "" ||
		// If the symbol is external module, don't show it in the completion list
		// (i.e declare module "http" { const x; } | // <= request completion here, "http" should not be there)
		symbol.Flags&ast.SymbolFlagsModule != 0 && startsWithQuote(name) ||
		// If the symbol is the internal name of an ES symbol, it is not a valid entry. Internal names for ES symbols start with "__@"
		checker.IsKnownSymbol(symbol) {
		return "", false
	}

	// !!! isIdentifierText should take in identifierVariant language variant
	// name is a valid identifier or private identifier text
	if scanner.IsIdentifierText(name, target) ||
		symbol.ValueDeclaration != nil && ast.IsPrivateIdentifierClassElementDeclaration(symbol.ValueDeclaration) {
		return name, false
	}
	if symbol.Flags&ast.SymbolFlagsAlias != 0 {
		// Allow non-identifier import/export aliases since we can insert them as string literals
		return name, true
	}

	switch completionKind {
	case CompletionKindMemberLike:
		if originIsComputedPropertyName(origin) {
			return origin.symbolName(), false
		}
		return "", false
	case CompletionKindObjectPropertyDeclaration:
		// TODO: GH#18169
		escapedName, _ := core.StringifyJson(name, "", "")
		return escapedName, false
	case CompletionKindPropertyAccess, CompletionKindGlobal:
		// For a 'this.' completion it will be in a global context, but may have a non-identifier name.
		// Don't add a completion for a name starting with a space. See https://github.com/Microsoft/TypeScript/pull/20547
		ch, _ := utf8.DecodeRuneInString(name)
		if ch == ' ' {
			return "", false
		}
		return name, true
	case CompletionKindNone, CompletionKindString:
		return name, false
	default:
		panic(fmt.Sprintf("Unexpected completion kind: %v", completionKind))
	}
}

// !!! refactor symbolOriginInfo so that we can tell the difference between flags and the kind of data it has
func originIsIgnore(origin *symbolOriginInfo) bool {
	return origin != nil && origin.kind&symbolOriginInfoKindIgnore != 0
}

func originIncludesSymbolName(origin *symbolOriginInfo) bool {
	return originIsExport(origin) || originIsResolvedExport(origin) || originIsComputedPropertyName(origin)
}

func originIsExport(origin *symbolOriginInfo) bool {
	return origin != nil && origin.kind&symbolOriginInfoKindExport != 0
}

func originIsResolvedExport(origin *symbolOriginInfo) bool {
	return origin != nil && origin.kind&symbolOriginInfoKindResolvedExport != 0
}

func originIsComputedPropertyName(origin *symbolOriginInfo) bool {
	return origin != nil && origin.kind&symbolOriginInfoKindComputedPropertyName != 0
}

func originIsObjectLiteralMethod(origin *symbolOriginInfo) bool {
	return origin != nil && origin.kind&symbolOriginInfoKindObjectLiteralMethod != 0
}

func originIsThisType(origin *symbolOriginInfo) bool {
	return origin != nil && origin.kind&symbolOriginInfoKindThisType != 0
}

func originIsTypeOnlyAlias(origin *symbolOriginInfo) bool {
	return origin != nil && origin.kind&symbolOriginInfoKindTypeOnlyAlias != 0
}

func originIsSymbolMember(origin *symbolOriginInfo) bool {
	return origin != nil && origin.kind&symbolOriginInfoKindSymbolMember != 0
}

func originIsNullableMember(origin *symbolOriginInfo) bool {
	return origin != nil && origin.kind&symbolOriginInfoKindNullable != 0
}

func originIsPromise(origin *symbolOriginInfo) bool {
	return origin != nil && origin.kind&symbolOriginInfoKindPromise != 0
}

func getSourceFromOrigin(origin *symbolOriginInfo) string {
	if originIsExport(origin) {
		return core.StripQuotes(origin.asExport().moduleSymbol.Name)
	}

	if originIsResolvedExport(origin) {
		return origin.asResolvedExport().moduleSpecifier
	}

	if originIsThisType(origin) {
		return string(completionSourceThisProperty)
	}

	if originIsTypeOnlyAlias(origin) {
		return string(completionSourceTypeOnlyAlias)
	}

	return ""
}

// In a scenarion such as `const x = 1 * |`, the context and previous tokens are both `*`.
// In `const x = 1 * o|`, the context token is *, and the previous token is `o`.
// `contextToken` and `previousToken` can both be nil if we are at the beginning of the file.
func getRelevantTokens(position int, file *ast.SourceFile) (contextToken *ast.Node, previousToken *ast.Node) {
	previousToken = astnav.FindPrecedingToken(file, position)
	if previousToken != nil && position <= previousToken.End() && (ast.IsMemberName(previousToken) || ast.IsKeywordKind(previousToken.Kind)) {
		contextToken := astnav.FindPrecedingToken(file, previousToken.Pos())
		return contextToken, previousToken
	}
	return previousToken, previousToken
}

// "." | '"' | "'" | "`" | "/" | "@" | "<" | "#" | " "
type CompletionsTriggerCharacter = string

func isValidTrigger(file *ast.SourceFile, triggerCharacter CompletionsTriggerCharacter, contextToken *ast.Node, position int) bool {
	switch triggerCharacter {
	case ".", "@":
		return true
	case "\"", "'", "`":
		// Only automatically bring up completions if this is an opening quote.
		return contextToken != nil &&
			isStringLiteralOrTemplate(contextToken) &&
			position == getStartOfNode(contextToken, file)+1
	case "#":
		return contextToken != nil &&
			ast.IsPrivateIdentifier(contextToken) &&
			ast.GetContainingClass(contextToken) != nil
	case "<":
		// Opening JSX tag
		return contextToken != nil &&
			contextToken.Kind == ast.KindLessThanToken &&
			(!ast.IsBinaryExpression(contextToken.Parent) || binaryExpressionMayBeOpenTag(contextToken.Parent.AsBinaryExpression()))
	case "/":
		if contextToken == nil {
			return false
		}
		if ast.IsStringLiteralLike(contextToken) {
			return tryGetImportFromModuleSpecifier(contextToken) != nil
		}
		return contextToken.Kind == ast.KindSlashToken && ast.IsJsxClosingElement(contextToken.Parent)
	case " ":
		return contextToken != nil && contextToken.Kind == ast.KindImportKeyword && contextToken.Parent.Kind == ast.KindSourceFile
	default:
		panic("Unknown trigger character: " + triggerCharacter)
	}
}

func isStringLiteralOrTemplate(node *ast.Node) bool {
	switch node.Kind {
	case ast.KindStringLiteral, ast.KindNoSubstitutionTemplateLiteral, ast.KindTemplateExpression,
		ast.KindTaggedTemplateExpression:
		return true
	}
	return false
}

func binaryExpressionMayBeOpenTag(binaryExpression *ast.BinaryExpression) bool {
	return ast.NodeIsMissing(binaryExpression.Left)
}

func getDefaultCommitCharacters(isNewIdentifierLocation bool) *[]string {
	// !!!
	return nil
}

func isCheckedFile(file *ast.SourceFile, compilerOptions *core.CompilerOptions) bool {
	return !ast.IsSourceFileJs(file) || ast.IsCheckJsEnabledForFile(file, compilerOptions)
}

func isContextTokenValueLocation(contextToken *ast.Node) bool {
	return contextToken != nil && ((contextToken.Kind == ast.KindTypeOfKeyword &&
		(contextToken.Parent.Kind == ast.KindTypeQuery || ast.IsTypeOfExpression(contextToken.Parent))) ||
		(contextToken.Kind == ast.KindAssertsKeyword && contextToken.Parent.Kind == ast.KindTypePredicate))
}

func isPossiblyTypeArgumentPosition(token *ast.Node, sourceFile *ast.SourceFile, typeChecker *checker.Checker) bool {
	info := getPossibleTypeArgumentsInfo(token, sourceFile)
	return info != nil && (ast.IsPartOfTypeNode(info.called) ||
		len(getPossibleGenericSignatures(info.called, info.nTypeArguments, typeChecker)) != 0 ||
		isPossiblyTypeArgumentPosition(info.called, sourceFile, typeChecker))
}

func isContextTokenTypeLocation(contextToken *ast.Node) bool {
	if contextToken != nil {
		parentKind := contextToken.Parent.Kind
		switch contextToken.Kind {
		case ast.KindColonToken:
			return parentKind == ast.KindPropertyDeclaration ||
				parentKind == ast.KindPropertySignature ||
				parentKind == ast.KindParameter ||
				parentKind == ast.KindVariableDeclaration ||
				ast.IsFunctionLikeKind(parentKind)
		case ast.KindEqualsToken:
			return parentKind == ast.KindTypeAliasDeclaration || parentKind == ast.KindTypeParameter
		case ast.KindAsKeyword:
			return parentKind == ast.KindAsExpression
		case ast.KindLessThanToken:
			return parentKind == ast.KindTypeReference || parentKind == ast.KindTypeAssertionExpression
		case ast.KindExtendsKeyword:
			return parentKind == ast.KindTypeParameter
		case ast.KindSatisfiesKeyword:
			return parentKind == ast.KindSatisfiesExpression
		}
	}
	return false
}

// True if symbol is a type or a module containing at least one type.
func symbolCanBeReferencedAtTypeLocation(symbol *ast.Symbol, typeChecker *checker.Checker, seenModules core.Set[ast.SymbolId]) bool {
	// Since an alias can be merged with a local declaration, we need to test both the alias and its target.
	// This code used to just test the result of `skipAlias`, but that would ignore any locally introduced meanings.
	return nonAliasCanBeReferencedAtTypeLocation(symbol, typeChecker, seenModules) ||
		nonAliasCanBeReferencedAtTypeLocation(
			checker.SkipAlias(core.IfElse(symbol.ExportSymbol != nil, symbol.ExportSymbol, symbol), typeChecker),
			typeChecker,
			seenModules,
		)
}

func nonAliasCanBeReferencedAtTypeLocation(symbol *ast.Symbol, typeChecker *checker.Checker, seenModules core.Set[ast.SymbolId]) bool {
	return symbol.Flags&ast.SymbolFlagsType != 0 || typeChecker.IsUnknownSymbol(symbol) ||
		symbol.Flags&ast.SymbolFlagsModule != 0 && core.AddIfAbsent(seenModules, ast.GetSymbolId(symbol)) &&
			core.Some(
				typeChecker.GetExportsOfModule(symbol),
				func(e *ast.Symbol) bool { return symbolCanBeReferencedAtTypeLocation(e, typeChecker, seenModules) })
}

// Gets all properties on a type, but if that type is a union of several types,
// excludes array-like types or callable/constructable types.
func getPropertiesForCompletion(t *checker.Type, typeChecker *checker.Checker) []*ast.Symbol {
	if checker.IsUnion(t) {
		return core.CheckEachDefined(typeChecker.GetAllPossiblePropertiesOfTypes(t.Types()), "getAllPossiblePropertiesOfTypes() should all be defined.")
	} else {
		return core.CheckEachDefined(typeChecker.GetApparentProperties(t), "getApparentProperties() should all be defined.")
	}
}

// Given 'a.b.c', returns 'a'.
func getLeftMostName(e *ast.Expression) *ast.IdentifierNode {
	if ast.IsIdentifier(e) {
		return e
	} else if ast.IsPropertyAccessExpression(e) {
		return getLeftMostName(e.Expression())
	} else {
		return nil
	}
}

func getFirstSymbolInChain(symbol *ast.Symbol, enclosingDeclaration *ast.Node, typeChecker *checker.Checker) *ast.Symbol {
	chain := typeChecker.GetAccessibleSymbolChain(
		symbol,
		enclosingDeclaration,
		ast.SymbolFlagsAll, /*meaning*/
		false /*useOnlyExternalAliasing*/)
	if len(chain) > 0 {
		return chain[0]
	}
	if symbol.Parent != nil {
		if isModuleSymbol(symbol.Parent) {
			return symbol
		}
		return getFirstSymbolInChain(symbol.Parent, enclosingDeclaration, typeChecker)
	}
	return nil
}

func isModuleSymbol(symbol *ast.Symbol) bool {
	return core.Some(symbol.Declarations, func(decl *ast.Declaration) bool { return decl.Kind == ast.KindSourceFile })
}

func getNullableSymbolOriginInfoKind(kind symbolOriginInfoKind, insertQuestionDot bool) symbolOriginInfoKind {
	if insertQuestionDot {
		kind |= symbolOriginInfoKindNullable
	}
	return kind
}

func isStaticProperty(symbol *ast.Symbol) bool {
	return symbol.ValueDeclaration != nil &&
		checker.GetEffectiveModifierFlags(symbol.ValueDeclaration)&ast.ModifierFlagsStatic != 0 &&
		ast.IsClassLike(symbol.ValueDeclaration.Parent)
}

func getContextualType(previousToken *ast.Node, position int, file *ast.SourceFile, typeChecker *checker.Checker) *checker.Type {
	parent := previousToken.Parent
	switch previousToken.Kind {
	case ast.KindIdentifier:
		return getContextualTypeFromParent(previousToken, typeChecker, checker.ContextFlagsNone)
	case ast.KindEqualsToken:
		switch parent.Kind {
		case ast.KindVariableDeclaration:
			return typeChecker.GetContextualType(parent.Initializer(), checker.ContextFlagsNone)
		case ast.KindBinaryExpression:
			return typeChecker.GetTypeAtLocation(parent.AsBinaryExpression().Left)
		case ast.KindJsxAttribute:
			// return typeChecker.GetContextualTypeForJsxAttribute(parent) // !!! jsx
			return nil
		default:
			return nil
		}
	case ast.KindNewKeyword:
		return typeChecker.GetContextualType(parent, checker.ContextFlagsNone)
	case ast.KindCaseKeyword:
		caseClause := core.IfElse(ast.IsCaseClause(parent), parent, nil)
		if caseClause != nil {
			return getSwitchedType(caseClause, typeChecker)
		}
		return nil
	case ast.KindOpenBraceToken:
		if ast.IsJsxExpression(parent) && !ast.IsJsxElement(parent.Parent) && !ast.IsJsxFragment(parent.Parent) {
			// return typeChecker.GetContextualTypeForJsxAttribute(parent.Parent) // !!! jsx
			return nil
		}
		return nil
	default:
		// argInfo := getArgumentInfoForCompletions(previousToken, position, file, typeChecker) // !!! signature help
		var argInfo *struct{} // !!! signature help
		if argInfo != nil {
			// return typeChecker.GetContextualTypeForArgumentAtIndex() // !!! signature help
			return nil
		} else if isEqualityOperatorKind(previousToken.Kind) && ast.IsBinaryExpression(parent) && isEqualityOperatorKind(parent.AsBinaryExpression().OperatorToken.Kind) {
			// completion at `x ===/**/`
			return typeChecker.GetTypeAtLocation(parent.AsBinaryExpression().Left)
		} else {
			contextualType := typeChecker.GetContextualType(previousToken, checker.ContextFlagsCompletions)
			if contextualType != nil {
				return contextualType
			}
			return typeChecker.GetContextualType(previousToken, checker.ContextFlagsNone)
		}
	}
}

// TODO: this is also used by string completions originally, but passing a flag `ContextFlagsCompletions`.
// What difference does it make?
func getContextualTypeFromParent(node *ast.Expression, typeChecker *checker.Checker, contextFlags checker.ContextFlags) *checker.Type {
	parent := ast.WalkUpParenthesizedExpressions(node.Parent)
	switch parent.Kind {
	case ast.KindNewExpression:
		return typeChecker.GetContextualType(parent, contextFlags)
	case ast.KindBinaryExpression:
		if isEqualityOperatorKind(parent.AsBinaryExpression().OperatorToken.Kind) {
			return typeChecker.GetTypeAtLocation(core.IfElse(node == parent.AsBinaryExpression().Right, parent.AsBinaryExpression().Left, parent.AsBinaryExpression().Right))
		}
		return typeChecker.GetContextualType(node, contextFlags)
	case ast.KindCaseClause:
		return getSwitchedType(parent, typeChecker)
	default:
		return typeChecker.GetContextualType(node, contextFlags)
	}
}

func getSwitchedType(caseClause *ast.CaseClauseNode, typeChecker *checker.Checker) *checker.Type {
	return typeChecker.GetTypeAtLocation(caseClause.Parent.Parent.Expression())
}

func isEqualityOperatorKind(kind ast.Kind) bool {
	switch kind {
	case ast.KindEqualsEqualsEqualsToken, ast.KindEqualsEqualsToken,
		ast.KindExclamationEqualsEqualsToken, ast.KindExclamationEqualsToken:
		return true
	default:
		return false
	}
}

func isLiteral(t *checker.Type) bool {
	return checker.IsStringLiteral(t) || checker.IsNumberLiteral(t) || checker.IsBigIntLiteral(t)
}

func getRecommendedCompletion(previousToken *ast.Node, contextualType *checker.Type, typeChecker *checker.Checker) *ast.Symbol {
	// For a union, return the first one with a recommended completion.
	return core.FirstNonNil(
		core.IfElse(checker.IsUnion(contextualType), contextualType.Types(), []*checker.Type{contextualType}),
		func(t *checker.Type) *ast.Symbol {
			symbol := t.Symbol()
			// Don't make a recommended completion for an abstract class.
			if symbol != nil &&
				symbol.Flags&(ast.SymbolFlagsEnumMember|ast.SymbolFlagsEnum|ast.SymbolFlagsClass) != 0 &&
				!isAbstractConstructorSymbol(symbol) {
				return getFirstSymbolInChain(symbol, previousToken, typeChecker)
			}
			return nil
		},
	)
}

func isAbstractConstructorSymbol(symbol *ast.Symbol) bool {
	if symbol.Flags&ast.SymbolFlagsClass != 0 {
		declaration := ast.GetClassLikeDeclarationOfSymbol(symbol)
		return declaration != nil && ast.HasSyntacticModifier(declaration, ast.ModifierFlagsAbstract)
	}
	return false
}

func startsWithQuote(s string) bool {
	r, _ := utf8.DecodeRuneInString(s)
	return r == '"' || r == '\''
}

func getClosestSymbolDeclaration(contextToken *ast.Node, location *ast.Node) *ast.Declaration {
	if contextToken == nil {
		return nil
	}

	closestDeclaration := ast.FindAncestorOrQuit(contextToken, func(node *ast.Node) ast.FindAncestorResult {
		if ast.IsFunctionBlock(node) || isArrowFunctionBody(node) || ast.IsBindingPattern(node) {
			return ast.FindAncestorQuit
		}

		if (ast.IsParameter(node) || ast.IsTypeParameterDeclaration(node)) &&
			!ast.IsIndexSignatureDeclaration(node.Parent) {
			return ast.FindAncestorTrue
		}
		return ast.FindAncestorFalse
	})

	if closestDeclaration == nil {
		closestDeclaration = ast.FindAncestorOrQuit(contextToken, func(node *ast.Node) ast.FindAncestorResult {
			if ast.IsFunctionBlock(node) || isArrowFunctionBody(node) || ast.IsBindingPattern(node) {
				return ast.FindAncestorQuit
			}

			if ast.IsVariableDeclaration(node) {
				return ast.FindAncestorTrue
			}
			return ast.FindAncestorFalse
		})
	}
	return closestDeclaration
}

func isArrowFunctionBody(node *ast.Node) bool {
	return node.Parent != nil && ast.IsArrowFunction(node.Parent) &&
		(node.Parent.Body() == node ||
			// const a = () => /**/;
			node.Kind == ast.KindEqualsGreaterThanToken)
}

func isInTypeParameterDefault(contextToken *ast.Node) bool {
	if contextToken == nil {
		return false
	}

	node := contextToken
	parent := contextToken.Parent
	for parent != nil {
		if ast.IsTypeParameterDeclaration(parent) {
			return parent.AsTypeParameter().DefaultType == node || node.Kind == ast.KindEqualsToken
		}
		node = parent
		parent = parent.Parent
	}

	return false
}

func isDeprecated(symbol *ast.Symbol, typeChecker *checker.Checker) bool {
	declarations := checker.SkipAlias(symbol, typeChecker).Declarations
	return len(declarations) > 0 && core.Every(declarations, func(decl *ast.Declaration) bool { return typeChecker.IsDeprecatedDeclaration(decl) })
}

func getReplacementRangeForContextToken(file *ast.SourceFile, contextToken *ast.Node, position int) *lsproto.Range {
	if contextToken == nil {
		return nil
	}

	// !!! ensure range is single line
	switch contextToken.Kind {
	case ast.KindStringLiteral, ast.KindNoSubstitutionTemplateLiteral:
		return createRangeFromStringLiteralLikeContent(file, contextToken, position)
	default:
		return createLspRangeFromNode(contextToken, file)
	}
}

func createRangeFromStringLiteralLikeContent(file *ast.SourceFile, node *ast.StringLiteralLike, position int) *lsproto.Range {
	replacementEnd := node.End() - 1
	nodeStart := getStartOfNode(node, file)
	if ast.IsUnterminatedLiteral(node) {
		// we return no replacement range only if unterminated string is empty
		if nodeStart == replacementEnd {
			return nil
		}
		replacementEnd = min(position, node.End())
	}
	return createLspRangeFromBounds(nodeStart+1, replacementEnd, file)
}

func quotePropertyName(file *ast.SourceFile, preferences *UserPreferences, name string) string {
	r, _ := utf8.DecodeRuneInString(name)
	if unicode.IsDigit(r) {
		return name
	}
	return quote(file, preferences, name)
}

// Checks whether type is `string & {}`, which is semantically equivalent to string but
// is not reduced by the checker as a special case used for supporting string literal completions
// for string type.
func isStringAndEmptyAnonymousObjectIntersection(typeChecker *checker.Checker, t *checker.Type) bool {
	if t.Flags()&checker.TypeFlagsIntersection == 0 {
		return false
	}

	return len(t.Types()) == 2 &&
		(areIntersectedTypesAvoidingStringReduction(typeChecker, t.Types()[0], t.Types()[1]) ||
			areIntersectedTypesAvoidingStringReduction(typeChecker, t.Types()[1], t.Types()[0]))
}

func areIntersectedTypesAvoidingStringReduction(typeChecker *checker.Checker, t1 *checker.Type, t2 *checker.Type) bool {
	return t1.Flags()&checker.TypeFlagsString != 0 && typeChecker.IsEmptyAnonymousObjectType(t2)
}

func escapeSnippetText(text string) string {
	return strings.ReplaceAll(text, `$`, `\$`)
}

func isNamedImportsOrExports(node *ast.Node) bool {
	return ast.IsNamedImports(node) || ast.IsNamedExports(node)
}

func generateIdentifierForArbitraryString(text string, languageVersion core.ScriptTarget) string {
	needsUnderscore := false
	identifier := ""
	var ch rune
	var size int

	// Convert "(example, text)" into "_example_text_"
	for pos := 0; pos < len(text); pos += size {
		ch, size = utf8.DecodeRuneInString(text[pos:])
		var validChar bool
		if pos == 0 {
			validChar = scanner.IsIdentifierStart(ch, languageVersion)
		} else {
			validChar = scanner.IsIdentifierPart(ch, languageVersion)
		}
		if size > 0 && validChar {
			if needsUnderscore {
				identifier += "_"
				identifier += string(ch)
				needsUnderscore = false
			}
		} else {
			needsUnderscore = true
		}
	}

	if needsUnderscore {
		identifier += "_"
	}

	// Default to "_" if the provided text was empty
	if identifier == "" {
		return "_"
	}

	return identifier
}

// Copied from vscode TS extension.
func getCompletionsSymbolKind(kind ScriptElementKind) lsproto.CompletionItemKind {
	switch kind {
	case ScriptElementKindPrimitiveType:
	case ScriptElementKindKeyword:
		return lsproto.CompletionItemKindKeyword
	case ScriptElementKindConstElement:
	case ScriptElementKindLetElement:
	case ScriptElementKindVariableElement:
	case ScriptElementKindLocalVariableElement:
	case ScriptElementKindAlias:
	case ScriptElementKindParameterElement:
		return lsproto.CompletionItemKindVariable

	case ScriptElementKindMemberVariableElement:
	case ScriptElementKindMemberGetAccessorElement:
	case ScriptElementKindMemberSetAccessorElement:
		return lsproto.CompletionItemKindField

	case ScriptElementKindFunctionElement:
	case ScriptElementKindLocalFunctionElement:
		return lsproto.CompletionItemKindFunction

	case ScriptElementKindMemberFunctionElement:
	case ScriptElementKindConstructSignatureElement:
	case ScriptElementKindCallSignatureElement:
	case ScriptElementKindIndexSignatureElement:
		return lsproto.CompletionItemKindMethod

	case ScriptElementKindEnumElement:
		return lsproto.CompletionItemKindEnum

	case ScriptElementKindEnumMemberElement:
		return lsproto.CompletionItemKindEnumMember

	case ScriptElementKindModuleElement:
	case ScriptElementKindExternalModuleName:
		return lsproto.CompletionItemKindModule

	case ScriptElementKindClassElement:
	case ScriptElementKindTypeElement:
		return lsproto.CompletionItemKindClass

	case ScriptElementKindInterfaceElement:
		return lsproto.CompletionItemKindInterface

	case ScriptElementKindWarning:
		return lsproto.CompletionItemKindText

	case ScriptElementKindScriptElement:
		return lsproto.CompletionItemKindFile

	case ScriptElementKindDirectory:
		return lsproto.CompletionItemKindFolder

	case ScriptElementKindString:
		return lsproto.CompletionItemKindConstant

	default:
		return lsproto.CompletionItemKindProperty
	}
	panic("Unhandled script element kind: " + kind)
}

// Editors will use the `sortText` and then fall back to `name` for sorting, but leave ties in response order.
// So, it's important that we sort those ties in the order we want them displayed if it matters. We don't
// strictly need to sort by name or SortText here since clients are going to do it anyway, but we have to
// do the work of comparing them so we can sort those ties appropriately; plus, it makes the order returned
// by the language service consistent with what TS Server does and what editors typically do. This also makes
// completions tests make more sense. We used to sort only alphabetically and only in the server layer, but
// this made tests really weird, since most fourslash tests don't use the server.
func compareCompletionEntries(entryInSlice *lsproto.CompletionItem, entryToInsert *lsproto.CompletionItem) int {
	// !!! use locale-aware comparison
	result := stringutil.CompareStringsCaseSensitive(*entryInSlice.SortText, *entryToInsert.SortText)
	if result == stringutil.ComparisonEqual {
		result = stringutil.CompareStringsCaseSensitive(entryInSlice.Label, entryToInsert.Label)
	}
	// !!! auto-imports
	// if (result === Comparison.EqualTo && entryInArray.data?.moduleSpecifier && entryToInsert.data?.moduleSpecifier) {
	//     // Sort same-named auto-imports by module specifier
	//     result = compareNumberOfDirectorySeparators(
	//         (entryInArray.data as CompletionEntryDataResolved).moduleSpecifier,
	//         (entryToInsert.data as CompletionEntryDataResolved).moduleSpecifier,
	//     );
	// }
	if result == stringutil.ComparisonEqual {
		// Fall back to symbol order - if we return `EqualTo`, `insertSorted` will put later symbols first.
		return stringutil.ComparisonLessThan
	}

	return result
}


func getKeywordCompletions(keywordFilter KeywordCompletionFilters, filterOutTsOnlyKeywords bool) []*lsproto.CompletionItem {
	if !filterOutTsOnlyKeywords {
		return getTypescriptKeywordCompletions(keywordFilter)
	}

	// !!! cache keyword list per filter
	return core.Filter(
		getTypescriptKeywordCompletions(keywordFilter),
		func(ci *lsproto.CompletionItem) bool {
			return !isTypeScriptOnlyKeyword(scanner.StringToToken(ci.Label))
		})
}

func getTypescriptKeywordCompletions(keywordFilter KeywordCompletionFilters) []*lsproto.CompletionItem {
	// !!! cache keyword list per filter
}
