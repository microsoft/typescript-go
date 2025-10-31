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
	"github.com/microsoft/typescript-go/internal/scanner"
	"github.com/microsoft/typescript-go/internal/tspath"
)

// tokenTypes defines the order of token types for encoding
var tokenTypes = []lsproto.SemanticTokenTypes{
	lsproto.SemanticTokenTypesnamespace,
	lsproto.SemanticTokenTypesclass,
	lsproto.SemanticTokenTypesenum,
	lsproto.SemanticTokenTypesinterface,
	lsproto.SemanticTokenTypesstruct,
	lsproto.SemanticTokenTypestypeParameter,
	lsproto.SemanticTokenTypestype,
	lsproto.SemanticTokenTypesparameter,
	lsproto.SemanticTokenTypesvariable,
	lsproto.SemanticTokenTypesproperty,
	lsproto.SemanticTokenTypesenumMember,
	lsproto.SemanticTokenTypesdecorator,
	lsproto.SemanticTokenTypesevent,
	lsproto.SemanticTokenTypesfunction,
	lsproto.SemanticTokenTypesmethod,
	lsproto.SemanticTokenTypesmacro,
	lsproto.SemanticTokenTypeslabel,
	lsproto.SemanticTokenTypescomment,
	lsproto.SemanticTokenTypesstring,
	lsproto.SemanticTokenTypeskeyword,
	lsproto.SemanticTokenTypesnumber,
	lsproto.SemanticTokenTypesregexp,
	lsproto.SemanticTokenTypesoperator,
}

// tokenModifiers defines the order of token modifiers for encoding
var tokenModifiers = []lsproto.SemanticTokenModifiers{
	lsproto.SemanticTokenModifiersdeclaration,
	lsproto.SemanticTokenModifiersdefinition,
	lsproto.SemanticTokenModifiersreadonly,
	lsproto.SemanticTokenModifiersstatic,
	lsproto.SemanticTokenModifiersdeprecated,
	lsproto.SemanticTokenModifiersabstract,
	lsproto.SemanticTokenModifiersasync,
	lsproto.SemanticTokenModifiersmodification,
	lsproto.SemanticTokenModifiersdocumentation,
	lsproto.SemanticTokenModifiersdefaultLibrary,
}

// tokenType represents a semantic token type index
type tokenType int

// Token type indices
const (
	tokenTypeNamespace tokenType = iota
	tokenTypeClass
	tokenTypeEnum
	tokenTypeInterface
	tokenTypeStruct
	tokenTypeTypeParameter
	tokenTypeType
	tokenTypeParameter
	tokenTypeVariable
	tokenTypeProperty
	tokenTypeEnumMember
	tokenTypeDecorator
	tokenTypeEvent
	tokenTypeFunction
	tokenTypeMethod
	tokenTypeMacro
	tokenTypeLabel
	tokenTypeComment
	tokenTypeString
	tokenTypeKeyword
	tokenTypeNumber
	tokenTypeRegexp
	tokenTypeOperator
)

// tokenModifier represents a semantic token modifier bit mask
type tokenModifier int

// Token modifier bit masks
const (
	tokenModifierDeclaration tokenModifier = 1 << iota
	tokenModifierDefinition
	tokenModifierReadonly
	tokenModifierStatic
	tokenModifierDeprecated
	tokenModifierAbstract
	tokenModifierAsync
	tokenModifierModification
	tokenModifierDocumentation
	tokenModifierDefaultLibrary
)

// SemanticTokensLegend returns the legend describing the token types and modifiers.
// If clientCapabilities is provided, it filters the legend to only include types and modifiers
// that the client supports.
func SemanticTokensLegend(clientCapabilities *lsproto.SemanticTokensClientCapabilities) *lsproto.SemanticTokensLegend {
	types := make([]string, 0, len(tokenTypes))
	for _, t := range tokenTypes {
		if clientCapabilities == nil || slices.Contains(clientCapabilities.TokenTypes, string(t)) {
			types = append(types, string(t))
		}
	}
	modifiers := make([]string, 0, len(tokenModifiers))
	for _, m := range tokenModifiers {
		if clientCapabilities == nil || slices.Contains(clientCapabilities.TokenModifiers, string(m)) {
			modifiers = append(modifiers, string(m))
		}
	}
	return &lsproto.SemanticTokensLegend{
		TokenTypes:     types,
		TokenModifiers: modifiers,
	}
}

func (l *LanguageService) ProvideSemanticTokens(ctx context.Context, documentURI lsproto.DocumentUri, clientCapabilities *lsproto.SemanticTokensClientCapabilities) (lsproto.SemanticTokensResponse, error) {
	program, file := l.getProgramAndFile(documentURI)

	c, done := program.GetTypeCheckerForFile(ctx, file)
	defer done()

	tokens := l.collectSemanticTokens(c, file, program)

	if len(tokens) == 0 {
		return lsproto.SemanticTokensOrNull{}, nil
	}

	// Convert to LSP format (relative encoding)
	encoded := encodeSemanticTokens(tokens, file, l.converters, clientCapabilities)

	return lsproto.SemanticTokensOrNull{
		SemanticTokens: &lsproto.SemanticTokens{
			Data: encoded,
		},
	}, nil
}

func (l *LanguageService) ProvideSemanticTokensRange(ctx context.Context, documentURI lsproto.DocumentUri, rang lsproto.Range, clientCapabilities *lsproto.SemanticTokensClientCapabilities) (lsproto.SemanticTokensRangeResponse, error) {
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
	encoded := encodeSemanticTokens(tokens, file, l.converters, clientCapabilities)

	return lsproto.SemanticTokensOrNull{
		SemanticTokens: &lsproto.SemanticTokens{
			Data: encoded,
		},
	}, nil
}

type semanticToken struct {
	node          *ast.Node
	tokenType     tokenType
	tokenModifier tokenModifier
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
					tokenModifier := tokenModifier(0)

					// Check if this is a declaration
					parent := node.Parent
					if parent != nil {
						parentIsDeclaration := ast.IsBindingElement(parent) || tokenFromDeclarationMapping(parent.Kind) == tokenType
						if parentIsDeclaration && parent.Name() == node {
							tokenModifier |= tokenModifierDeclaration
						}
					}

					// Reclassify parameters as properties in property access context
					if tokenType == tokenTypeParameter && ast.IsRightSideOfQualifiedNameOrPropertyAccess(node) {
						tokenType = tokenTypeProperty
					}

					// Reclassify based on type information
					tokenType = reclassifyByType(c, node, tokenType)

					// Get the value declaration to check modifiers
					if decl := symbol.ValueDeclaration; decl != nil {
						modifiers := ast.GetCombinedModifierFlags(decl)
						nodeFlags := ast.GetCombinedNodeFlags(decl)

						if modifiers&ast.ModifierFlagsStatic != 0 {
							tokenModifier |= tokenModifierStatic
						}
						if modifiers&ast.ModifierFlagsAsync != 0 {
							tokenModifier |= tokenModifierAsync
						}
						if tokenType != tokenTypeClass && tokenType != tokenTypeInterface {
							if (modifiers&ast.ModifierFlagsReadonly != 0) || (nodeFlags&ast.NodeFlagsConst != 0) || (symbol.Flags&ast.SymbolFlagsEnumMember != 0) {
								tokenModifier |= tokenModifierReadonly
							}
						}
						if (tokenType == tokenTypeVariable || tokenType == tokenTypeFunction) && isLocalDeclaration(decl, file) {
							// Local variables get no special modifier in LSP, but we track it internally
						}
						declSourceFile := ast.GetSourceFileOfNode(decl)
						if declSourceFile != nil && program.IsSourceFileDefaultLibrary(tspath.Path(declSourceFile.FileName())) {
							tokenModifier |= tokenModifierDefaultLibrary
						}
					} else if symbol.Declarations != nil {
						for _, decl := range symbol.Declarations {
							declSourceFile := ast.GetSourceFileOfNode(decl)
							if declSourceFile != nil && program.IsSourceFileDefaultLibrary(tspath.Path(declSourceFile.FileName())) {
								tokenModifier |= tokenModifierDefaultLibrary
								break
							}
						}
					}

					tokens = append(tokens, semanticToken{
						node:          node,
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

func classifySymbol(symbol *ast.Symbol, meaning ast.SemanticMeaning) (tokenType, bool) {
	flags := symbol.Flags
	if flags&ast.SymbolFlagsClass != 0 {
		return tokenTypeClass, true
	}
	if flags&ast.SymbolFlagsEnum != 0 {
		return tokenTypeEnum, true
	}
	if flags&ast.SymbolFlagsTypeAlias != 0 {
		return tokenTypeType, true
	}
	if flags&ast.SymbolFlagsInterface != 0 {
		if meaning&ast.SemanticMeaningType != 0 {
			return tokenTypeInterface, true
		}
	}
	if flags&ast.SymbolFlagsTypeParameter != 0 {
		return tokenTypeTypeParameter, true
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

func tokenFromDeclarationMapping(kind ast.Kind) tokenType {
	switch kind {
	case ast.KindVariableDeclaration:
		return tokenTypeVariable
	case ast.KindParameter:
		return tokenTypeParameter
	case ast.KindPropertyDeclaration:
		return tokenTypeProperty
	case ast.KindModuleDeclaration:
		return tokenTypeNamespace
	case ast.KindEnumDeclaration:
		return tokenTypeEnum
	case ast.KindEnumMember:
		return tokenTypeEnumMember
	case ast.KindClassDeclaration:
		return tokenTypeClass
	case ast.KindMethodDeclaration:
		return tokenTypeMethod
	case ast.KindFunctionDeclaration, ast.KindFunctionExpression:
		return tokenTypeFunction
	case ast.KindMethodSignature:
		return tokenTypeMethod
	case ast.KindGetAccessor, ast.KindSetAccessor:
		return tokenTypeProperty
	case ast.KindPropertySignature:
		return tokenTypeProperty
	case ast.KindInterfaceDeclaration:
		return tokenTypeInterface
	case ast.KindTypeAliasDeclaration:
		return tokenTypeType
	case ast.KindTypeParameter:
		return tokenTypeTypeParameter
	case ast.KindPropertyAssignment, ast.KindShorthandPropertyAssignment:
		return tokenTypeProperty
	default:
		return -1
	}
}

func reclassifyByType(c *checker.Checker, node *ast.Node, tt tokenType) tokenType {
	// Type-based reclassification for variables, properties, and parameters
	if tt == tokenTypeVariable || tt == tokenTypeProperty || tt == tokenTypeParameter {
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
			if tt != tokenTypeParameter && test(func(t *checker.Type) bool {
				return len(c.GetSignaturesOfType(t, checker.SignatureKindConstruct)) > 0
			}) {
				return tokenTypeClass
			}

			// Check for call signatures (function-like)
			if test(func(t *checker.Type) bool {
				callSigs := c.GetSignaturesOfType(t, checker.SignatureKindCall)
				if len(callSigs) == 0 {
					return false
				}
				// Must have call signatures and no properties (or be used in call context)
				objType := t.AsObjectType()
				if objType == nil {
					return true
				}
				return len(objType.Properties()) == 0 || isExpressionInCallExpression(node)
			}) {
				if tt == tokenTypeProperty {
					return tokenTypeMethod
				}
				return tokenTypeFunction
			}
		}
	}
	return tt
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

// encodeSemanticTokens encodes tokens into the LSP format using relative positioning.
// It filters tokens based on client capabilities, only including types and modifiers that the client supports.
func encodeSemanticTokens(tokens []semanticToken, file *ast.SourceFile, converters *lsconv.Converters, clientCapabilities *lsproto.SemanticTokensClientCapabilities) []uint32 {
	// Build mapping from server token types/modifiers to client indices
	typeMapping := make(map[tokenType]int)
	modifierMapping := make(map[lsproto.SemanticTokenModifiers]int)

	if clientCapabilities != nil {
		// Map server token types to client-supported indices
		clientIdx := 0
		for _, serverType := range tokenTypes {
			if slices.Contains(clientCapabilities.TokenTypes, string(serverType)) {
				// Find the server index for this type
				for i, t := range tokenTypes {
					if t == serverType {
						typeMapping[tokenType(i)] = clientIdx
						break
					}
				}
				clientIdx++
			}
		}

		// Map server token modifiers to client-supported bit positions
		clientBit := 0
		for _, serverModifier := range tokenModifiers {
			if slices.Contains(clientCapabilities.TokenModifiers, string(serverModifier)) {
				modifierMapping[serverModifier] = clientBit
				clientBit++
			}
		}
	} else {
		// No filtering - use direct mapping
		for i := range tokenTypes {
			typeMapping[tokenType(i)] = i
		}
		for i, mod := range tokenModifiers {
			modifierMapping[mod] = i
		}
	}

	// Each token encodes 5 uint32 values: deltaLine, deltaChar, length, tokenType, tokenModifiers
	encoded := make([]uint32, 0, len(tokens)*5)
	prevLine := uint32(0)
	prevChar := uint32(0)

	for _, token := range tokens {
		// Skip tokens with types not supported by the client
		clientTypeIdx, typeSupported := typeMapping[token.tokenType]
		if !typeSupported {
			continue
		}

		// Map modifiers to client-supported bit mask
		clientModifierMask := uint32(0)
		for i, serverModifier := range tokenModifiers {
			if token.tokenModifier&(1<<i) != 0 {
				if clientBit, ok := modifierMapping[serverModifier]; ok {
					clientModifierMask |= 1 << clientBit
				}
			}
		}

		// Use GetTokenPosOfNode to skip trivia (comments, whitespace) before the identifier
		tokenStart := scanner.GetTokenPosOfNode(token.node, file, false)
		tokenEnd := token.node.End()

		// Convert both start and end positions to LSP coordinates, then compute length
		startPos := converters.PositionToLineAndCharacter(file, core.TextPos(tokenStart))
		endPos := converters.PositionToLineAndCharacter(file, core.TextPos(tokenEnd))

		// Length is the character difference when on the same line
		var tokenLength uint32
		if startPos.Line == endPos.Line {
			tokenLength = endPos.Character - startPos.Character
		} else {
			// Multi-line tokens shouldn't happen for identifiers, but handle it
			tokenLength = endPos.Character
		}

		line := startPos.Line
		char := startPos.Character

		// Verify that positions are strictly increasing (visitor walks in order)
		if line < prevLine || (line == prevLine && char <= prevChar) {
			panic("semantic tokens: positions must be strictly increasing")
		}

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
			tokenLength,
			uint32(clientTypeIdx),
			clientModifierMask,
		)

		prevLine = line
		prevChar = char
	}

	return encoded
}
