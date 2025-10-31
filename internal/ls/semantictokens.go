package ls

import (
	"context"
	"slices"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/checker"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls/lsconv"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/tspath"
)

// Token types according to LSP specification
const (
	TokenTypeNamespace = iota
	TokenTypeClass
	TokenTypeEnum
	TokenTypeInterface
	TokenTypeStruct
	TokenTypeTypeParameter
	TokenTypeType
	TokenTypeParameter
	TokenTypeVariable
	TokenTypeProperty
	TokenTypeEnumMember
	TokenTypeDecorator
	TokenTypeEvent
	TokenTypeFunction
	TokenTypeMethod
	TokenTypeMacro
	TokenTypeLabel
	TokenTypeComment
	TokenTypeString
	TokenTypeKeyword
	TokenTypeNumber
	TokenTypeRegexp
	TokenTypeOperator
)

// Token modifiers according to LSP specification
const (
	TokenModifierDeclaration = 1 << iota
	TokenModifierDefinition
	TokenModifierReadonly
	TokenModifierStatic
	TokenModifierDeprecated
	TokenModifierAbstract
	TokenModifierAsync
	TokenModifierModification
	TokenModifierDocumentation
	TokenModifierDefaultLibrary
)

// SemanticTokensLegend returns the legend describing the token types and modifiers
func SemanticTokensLegend() *lsproto.SemanticTokensLegend {
	return &lsproto.SemanticTokensLegend{
		TokenTypes: []string{
			"namespace",
			"class",
			"enum",
			"interface",
			"struct",
			"typeParameter",
			"type",
			"parameter",
			"variable",
			"property",
			"enumMember",
			"decorator",
			"event",
			"function",
			"method",
			"macro",
			"label",
			"comment",
			"string",
			"keyword",
			"number",
			"regexp",
			"operator",
		},
		TokenModifiers: []string{
			"declaration",
			"definition",
			"readonly",
			"static",
			"deprecated",
			"abstract",
			"async",
			"modification",
			"documentation",
			"defaultLibrary",
		},
	}
}

func (l *LanguageService) ProvideSemanticTokens(ctx context.Context, documentURI lsproto.DocumentUri) (lsproto.SemanticTokensResponse, error) {
	program, file := l.getProgramAndFile(documentURI)

	c, done := program.GetTypeCheckerForFile(ctx, file)
	defer done()

	tokens := l.collectSemanticTokens(c, file, program)

	if len(tokens) == 0 {
		return lsproto.SemanticTokensOrNull{}, nil
	}

	// Convert to LSP format (relative encoding)
	encoded := encodeSemanticTokens(tokens, file, l.converters)

	return lsproto.SemanticTokensOrNull{
		SemanticTokens: &lsproto.SemanticTokens{
			Data: encoded,
		},
	}, nil
}

func (l *LanguageService) ProvideSemanticTokensRange(ctx context.Context, documentURI lsproto.DocumentUri, rang lsproto.Range) (lsproto.SemanticTokensRangeResponse, error) {
	program, file := l.getProgramAndFile(documentURI)

	c, done := program.GetTypeCheckerForFile(ctx, file)
	defer done()

	start := int(l.converters.LineAndCharacterToPosition(file, rang.Start))
	end := int(l.converters.LineAndCharacterToPosition(file, rang.End))

	tokens := l.collectSemanticTokensInRange(c, file, program, start, end)

	if len(tokens) == 0 {
		return lsproto.SemanticTokensOrNull{}, nil
	}

	// Convert to LSP format (relative encoding)
	encoded := encodeSemanticTokens(tokens, file, l.converters)

	return lsproto.SemanticTokensOrNull{
		SemanticTokens: &lsproto.SemanticTokens{
			Data: encoded,
		},
	}, nil
}

type semanticToken struct {
	pos           int
	length        int
	tokenType     int
	tokenModifier int
}

func (l *LanguageService) collectSemanticTokens(c *checker.Checker, file *ast.SourceFile, program *compiler.Program) []semanticToken {
	return l.collectSemanticTokensInRange(c, file, program, file.Pos(), file.End())
}

func (l *LanguageService) collectSemanticTokensInRange(c *checker.Checker, file *ast.SourceFile, program *compiler.Program, spanStart, spanEnd int) []semanticToken {
	tokens := []semanticToken{}

	inJSXElement := false

	var visit func(*ast.Node) bool
	visit = func(node *ast.Node) bool {
		// Note: cancellation is handled at the handler level, not here

		if node == nil {
			return false
		}
		nodeEnd := node.End()
		if node.Pos() >= spanEnd || nodeEnd <= spanStart {
			return false
		}

		prevInJSXElement := inJSXElement
		if ast.IsJsxElement(node) || ast.IsJsxSelfClosingElement(node) {
			inJSXElement = true
		}
		if ast.IsJsxExpression(node) {
			inJSXElement = false
		}

		if ast.IsIdentifier(node) && !inJSXElement && !isInImportClause(node) && !isInfinityOrNaNString(node.Text()) {
			symbol := c.GetSymbolAtLocation(node)
			if symbol != nil {
				// Resolve aliases
				if symbol.Flags&ast.SymbolFlagsAlias != 0 {
					symbol = c.GetAliasedSymbol(symbol)
				}

				tokenType, ok := classifySymbol(symbol, getMeaningFromLocation(node))
				if ok {
					tokenModifier := 0

					// Check if this is a declaration
					parent := node.Parent
					if parent != nil {
						parentIsDeclaration := ast.IsBindingElement(parent) || tokenFromDeclarationMapping(parent.Kind) == tokenType
						if parentIsDeclaration && parent.Name() == node {
							tokenModifier |= TokenModifierDeclaration
						}
					}

					// Reclassify parameters as properties in property access context
					if tokenType == TokenTypeParameter && ast.IsRightSideOfQualifiedNameOrPropertyAccess(node) {
						tokenType = TokenTypeProperty
					}

					// Reclassify based on type information
					tokenType = reclassifyByType(c, node, tokenType)

					// Get the value declaration to check modifiers
					if decl := symbol.ValueDeclaration; decl != nil {
						modifiers := ast.GetCombinedModifierFlags(decl)
						nodeFlags := ast.GetCombinedNodeFlags(decl)

						if modifiers&ast.ModifierFlagsStatic != 0 {
							tokenModifier |= TokenModifierStatic
						}
						if modifiers&ast.ModifierFlagsAsync != 0 {
							tokenModifier |= TokenModifierAsync
						}
						if tokenType != TokenTypeClass && tokenType != TokenTypeInterface {
							if (modifiers&ast.ModifierFlagsReadonly != 0) || (nodeFlags&ast.NodeFlagsConst != 0) || (symbol.Flags&ast.SymbolFlagsEnumMember != 0) {
								tokenModifier |= TokenModifierReadonly
							}
						}
						if (tokenType == TokenTypeVariable || tokenType == TokenTypeFunction) && isLocalDeclaration(decl, file) {
							// Local variables get no special modifier in LSP, but we track it internally
						}
						declSourceFile := ast.GetSourceFileOfNode(decl)
						if declSourceFile != nil && program.IsSourceFileDefaultLibrary(tspath.Path(declSourceFile.FileName())) {
							tokenModifier |= TokenModifierDefaultLibrary
						}
					} else if symbol.Declarations != nil {
						for _, decl := range symbol.Declarations {
							declSourceFile := ast.GetSourceFileOfNode(decl)
							if declSourceFile != nil && program.IsSourceFileDefaultLibrary(tspath.Path(declSourceFile.FileName())) {
								tokenModifier |= TokenModifierDefaultLibrary
								break
							}
						}
					}

					tokens = append(tokens, semanticToken{
						pos:           node.Pos(),
						length:        node.End() - node.Pos(),
						tokenType:     tokenType,
						tokenModifier: tokenModifier,
					})
				}
			}
		}

		node.ForEachChild(visit)
		inJSXElement = prevInJSXElement
		return false
	}

	visit(&file.Node)
	return tokens
}

func classifySymbol(symbol *ast.Symbol, meaning ast.SemanticMeaning) (int, bool) {
	flags := symbol.Flags
	if flags&ast.SymbolFlagsClass != 0 {
		return TokenTypeClass, true
	}
	if flags&ast.SymbolFlagsEnum != 0 {
		return TokenTypeEnum, true
	}
	if flags&ast.SymbolFlagsTypeAlias != 0 {
		return TokenTypeType, true
	}
	if flags&ast.SymbolFlagsInterface != 0 {
		if meaning&ast.SemanticMeaningType != 0 {
			return TokenTypeInterface, true
		}
	}
	if flags&ast.SymbolFlagsTypeParameter != 0 {
		return TokenTypeTypeParameter, true
	}

	// Check the value declaration
	decl := symbol.ValueDeclaration
	if decl == nil && symbol.Declarations != nil && len(symbol.Declarations) > 0 {
		decl = symbol.Declarations[0]
	}
	if decl != nil && ast.IsBindingElement(decl) {
		decl = getDeclarationForBindingElement(decl)
	}
	if decl != nil {
		if tokenType := tokenFromDeclarationMapping(decl.Kind); tokenType >= 0 {
			return tokenType, true
		}
	}

	return 0, false
}

func tokenFromDeclarationMapping(kind ast.Kind) int {
	switch kind {
	case ast.KindVariableDeclaration:
		return TokenTypeVariable
	case ast.KindParameter:
		return TokenTypeParameter
	case ast.KindPropertyDeclaration:
		return TokenTypeProperty
	case ast.KindModuleDeclaration:
		return TokenTypeNamespace
	case ast.KindEnumDeclaration:
		return TokenTypeEnum
	case ast.KindEnumMember:
		return TokenTypeEnumMember
	case ast.KindClassDeclaration:
		return TokenTypeClass
	case ast.KindMethodDeclaration:
		return TokenTypeMethod
	case ast.KindFunctionDeclaration, ast.KindFunctionExpression:
		return TokenTypeFunction
	case ast.KindMethodSignature:
		return TokenTypeMethod
	case ast.KindGetAccessor, ast.KindSetAccessor:
		return TokenTypeProperty
	case ast.KindPropertySignature:
		return TokenTypeProperty
	case ast.KindInterfaceDeclaration:
		return TokenTypeInterface
	case ast.KindTypeAliasDeclaration:
		return TokenTypeType
	case ast.KindTypeParameter:
		return TokenTypeTypeParameter
	case ast.KindPropertyAssignment, ast.KindShorthandPropertyAssignment:
		return TokenTypeProperty
	default:
		return -1
	}
}

func reclassifyByType(c *checker.Checker, node *ast.Node, tokenType int) int {
	// Type-based reclassification for variables, properties, and parameters
	if tokenType == TokenTypeVariable || tokenType == TokenTypeProperty || tokenType == TokenTypeParameter {
		typ := c.GetTypeAtLocation(node)
		if typ != nil {
			test := func(condition func(*checker.Type) bool) bool {
				if condition(typ) {
					return true
				}
				if typ.Flags()&checker.TypeFlagsUnion != 0 {
					for _, t := range typ.AsUnionType().Types() {
						if condition(t) {
							return true
						}
					}
				}
				return false
			}

			// Check for constructor signatures (class-like)
			if tokenType != TokenTypeParameter && test(func(t *checker.Type) bool {
				return len(c.GetSignaturesOfType(t, checker.SignatureKindConstruct)) > 0
			}) {
				return TokenTypeClass
			}

			// Check for call signatures (function-like)
			if test(func(t *checker.Type) bool {
				callSigs := c.GetSignaturesOfType(t, checker.SignatureKindCall)
				if len(callSigs) == 0 {
					return false
				}
				// Must have call signatures and no properties (or be used in call context)
				return len(t.AsObjectType().Properties()) == 0 || isExpressionInCallExpression(node)
			}) {
				if tokenType == TokenTypeProperty {
					return TokenTypeMethod
				}
				return TokenTypeFunction
			}
		}
	}
	return tokenType
}

func isLocalDeclaration(decl *ast.Node, sourceFile *ast.SourceFile) bool {
	if ast.IsBindingElement(decl) {
		decl = getDeclarationForBindingElement(decl)
	}
	if ast.IsVariableDeclaration(decl) {
		parent := decl.Parent
		if parent != nil && ast.IsVariableDeclarationList(parent) {
			grandparent := parent.Parent
			if grandparent != nil {
				return (!ast.IsSourceFile(grandparent) || ast.IsCatchClause(grandparent)) &&
					ast.GetSourceFileOfNode(decl) == sourceFile
			}
		}
	} else if ast.IsFunctionDeclaration(decl) {
		parent := decl.Parent
		return parent != nil && !ast.IsSourceFile(parent) && ast.GetSourceFileOfNode(decl) == sourceFile
	}
	return false
}

func getDeclarationForBindingElement(element *ast.Node) *ast.Node {
	for {
		parent := element.Parent
		if parent != nil && ast.IsBindingPattern(parent) {
			grandparent := parent.Parent
			if grandparent != nil && ast.IsBindingElement(grandparent) {
				element = grandparent
				continue
			}
		}
		if parent != nil && ast.IsBindingPattern(parent) {
			return parent.Parent
		}
		return element
	}
}

func isInImportClause(node *ast.Node) bool {
	parent := node.Parent
	if parent == nil {
		return false
	}
	return ast.IsImportClause(parent) || ast.IsImportSpecifier(parent) || ast.IsNamespaceImport(parent)
}

func isExpressionInCallExpression(node *ast.Node) bool {
	for ast.IsRightSideOfQualifiedNameOrPropertyAccess(node) {
		node = node.Parent
	}
	parent := node.Parent
	return parent != nil && ast.IsCallExpression(parent) && parent.Expression() == node
}

func isInfinityOrNaNString(text string) bool {
	return text == "Infinity" || text == "NaN"
}

// encodeSemanticTokens encodes tokens into the LSP format using relative positioning
func encodeSemanticTokens(tokens []semanticToken, file *ast.SourceFile, converters *lsconv.Converters) []uint32 {
	// Sort tokens by position
	slices.SortFunc(tokens, func(a, b semanticToken) int {
		return a.pos - b.pos
	})

	encoded := []uint32{}
	prevLine := uint32(0)
	prevChar := uint32(0)

	for _, token := range tokens {
		pos := converters.PositionToLineAndCharacter(file, core.TextPos(token.pos))
		line := pos.Line
		char := pos.Character

		// Encode as: [deltaLine, deltaChar, length, tokenType, tokenModifiers]
		deltaLine := line - prevLine
		var deltaChar uint32
		if deltaLine == 0 {
			deltaChar = char - prevChar
		} else {
			deltaChar = char
		}

		encoded = append(encoded,
			deltaLine,
			deltaChar,
			uint32(token.length),
			uint32(token.tokenType),
			uint32(token.tokenModifier),
		)

		prevLine = line
		prevChar = char
	}

	return encoded
}
