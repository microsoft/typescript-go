package ls

import (
	"context"
	"math"
	"path/filepath"
	"slices"
	"strings"

	"github.com/go-json-experiment/json"
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/astnav"
	"github.com/microsoft/typescript-go/internal/checker"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/scanner"
	"github.com/microsoft/typescript-go/internal/sourcemap"
	"github.com/microsoft/typescript-go/internal/tspath"
)

func (l *LanguageService) ProvideDefinition(ctx context.Context, documentURI lsproto.DocumentUri, position lsproto.Position) (lsproto.DefinitionResponse, error) {
	program, file := l.getProgramAndFile(documentURI)
	node := astnav.GetTouchingPropertyName(file, int(l.converters.LineAndCharacterToPosition(file, position)))
	if node.Kind == ast.KindSourceFile {
		return lsproto.LocationOrLocationsOrDefinitionLinksOrNull{}, nil
	}

	c, done := program.GetTypeCheckerForFile(ctx, file)
	defer done()

	if node.Kind == ast.KindOverrideKeyword {
		if sym := getSymbolForOverriddenMember(c, node); sym != nil {
			return l.createLocationsFromDeclarations(sym.Declarations), nil
		}
	}

	if ast.IsJumpStatementTarget(node) {
		if label := getTargetLabel(node.Parent, node.Text()); label != nil {
			return l.createLocationsFromDeclarations([]*ast.Node{label}), nil
		}
	}

	if node.Kind == ast.KindCaseKeyword || node.Kind == ast.KindDefaultKeyword && ast.IsDefaultClause(node.Parent) {
		if stmt := ast.FindAncestor(node.Parent, ast.IsSwitchStatement); stmt != nil {
			file := ast.GetSourceFileOfNode(stmt)
			return l.createLocationFromFileAndRange(file, scanner.GetRangeOfTokenAtPosition(file, stmt.Pos())), nil
		}
	}

	if node.Kind == ast.KindReturnKeyword || node.Kind == ast.KindYieldKeyword || node.Kind == ast.KindAwaitKeyword {
		if fn := ast.FindAncestor(node, ast.IsFunctionLikeDeclaration); fn != nil {
			return l.createLocationsFromDeclarations([]*ast.Node{fn}), nil
		}
	}

	declarations := getDeclarationsFromLocation(c, node)
	calledDeclaration := tryGetSignatureDeclaration(c, node)
	if calledDeclaration != nil {
		// If we can resolve a call signature, remove all function-like declarations and add that signature.
		nonFunctionDeclarations := core.Filter(slices.Clip(declarations), func(node *ast.Node) bool { return !ast.IsFunctionLike(node) })
		declarations = append(nonFunctionDeclarations, calledDeclaration)
	}
	return l.createLocationsFromDeclarations(declarations), nil
}

func (l *LanguageService) ProvideTypeDefinition(ctx context.Context, documentURI lsproto.DocumentUri, position lsproto.Position) (lsproto.DefinitionResponse, error) {
	program, file := l.getProgramAndFile(documentURI)
	node := astnav.GetTouchingPropertyName(file, int(l.converters.LineAndCharacterToPosition(file, position)))
	if node.Kind == ast.KindSourceFile {
		return lsproto.LocationOrLocationsOrDefinitionLinksOrNull{}, nil
	}

	c, done := program.GetTypeCheckerForFile(ctx, file)
	defer done()

	node = getDeclarationNameForKeyword(node)

	if symbol := c.GetSymbolAtLocation(node); symbol != nil {
		symbolType := getTypeOfSymbolAtLocation(c, symbol, node)
		declarations := getDeclarationsFromType(symbolType)
		if typeArgument := c.GetFirstTypeArgumentFromKnownType(symbolType); typeArgument != nil {
			declarations = core.Concatenate(getDeclarationsFromType(typeArgument), declarations)
		}
		if len(declarations) != 0 {
			return l.createLocationsFromDeclarations(declarations), nil
		}
		if symbol.Flags&ast.SymbolFlagsValue == 0 && symbol.Flags&ast.SymbolFlagsType != 0 {
			return l.createLocationsFromDeclarations(symbol.Declarations), nil
		}
	}

	return lsproto.LocationOrLocationsOrDefinitionLinksOrNull{}, nil
}

func getDeclarationNameForKeyword(node *ast.Node) *ast.Node {
	if node.Kind >= ast.KindFirstKeyword && node.Kind <= ast.KindLastKeyword {
		if ast.IsVariableDeclarationList(node.Parent) {
			if decl := core.FirstOrNil(node.Parent.AsVariableDeclarationList().Declarations.Nodes); decl != nil && decl.Name() != nil {
				return decl.Name()
			}
		} else if node.Parent.DeclarationData() != nil && node.Parent.Name() != nil && node.Pos() < node.Parent.Name().Pos() {
			return node.Parent.Name()
		}
	}
	return node
}

func (l *LanguageService) createLocationsFromDeclarations(declarations []*ast.Node) lsproto.DefinitionResponse {
	locations := make([]lsproto.Location, 0, len(declarations))
	for _, decl := range declarations {
		file := ast.GetSourceFileOfNode(decl)
		name := core.OrElse(ast.GetNameOfDeclaration(decl), decl)

		fileName := file.FileName()
		if strings.HasSuffix(fileName, ".d.ts") {
			if mappedLocation := l.tryMapToOriginalSource(file, name); mappedLocation != nil {
				locations = core.AppendIfUnique(locations, *mappedLocation)
				continue
			}
		}
		locations = core.AppendIfUnique(locations, lsproto.Location{
			Uri:   FileNameToDocumentURI(fileName),
			Range: *l.createLspRangeFromNode(name, file),
		})
	}
	return lsproto.LocationOrLocationsOrDefinitionLinksOrNull{Locations: &locations}
}

func (l *LanguageService) createLocationFromFileAndRange(file *ast.SourceFile, textRange core.TextRange) lsproto.DefinitionResponse {
	return lsproto.LocationOrLocationsOrDefinitionLinksOrNull{
		Location: &lsproto.Location{
			Uri:   FileNameToDocumentURI(file.FileName()),
			Range: *l.createLspRangeFromBounds(textRange.Pos(), textRange.End(), file),
		},
	}
}

func getDeclarationsFromLocation(c *checker.Checker, node *ast.Node) []*ast.Node {
	if ast.IsIdentifier(node) && ast.IsShorthandPropertyAssignment(node.Parent) {
		return c.GetResolvedSymbol(node).Declarations
	}
	node = getDeclarationNameForKeyword(node)
	if symbol := c.GetSymbolAtLocation(node); symbol != nil {
		if symbol.Flags&ast.SymbolFlagsClass != 0 && symbol.Flags&(ast.SymbolFlagsFunction|ast.SymbolFlagsVariable) == 0 && node.Kind == ast.KindConstructorKeyword {
			if constructor := symbol.Members[ast.InternalSymbolNameConstructor]; constructor != nil {
				symbol = constructor
			}
		}
		if symbol.Flags&ast.SymbolFlagsAlias != 0 {
			if resolved, ok := c.ResolveAlias(symbol); ok {
				symbol = resolved
			}
		}
		if symbol.Flags&(ast.SymbolFlagsProperty|ast.SymbolFlagsMethod|ast.SymbolFlagsAccessor) != 0 && symbol.Parent != nil && symbol.Parent.Flags&ast.SymbolFlagsObjectLiteral != 0 {
			if objectLiteral := core.FirstOrNil(symbol.Parent.Declarations); objectLiteral != nil {
				if declarations := c.GetContextualDeclarationsForObjectLiteralElement(objectLiteral, symbol.Name); len(declarations) != 0 {
					return declarations
				}
			}
		}
		return symbol.Declarations
	}
	if indexInfos := c.GetIndexSignaturesAtLocation(node); len(indexInfos) != 0 {
		return indexInfos
	}
	return nil
}

// Returns a CallLikeExpression where `node` is the target being invoked.
func getAncestorCallLikeExpression(node *ast.Node) *ast.Node {
	target := ast.FindAncestor(node, func(n *ast.Node) bool {
		return !isRightSideOfPropertyAccess(n)
	})
	callLike := target.Parent
	if callLike != nil && ast.IsCallLikeExpression(callLike) && ast.GetInvokedExpression(callLike) == target {
		return callLike
	}
	return nil
}

func tryGetSignatureDeclaration(typeChecker *checker.Checker, node *ast.Node) *ast.Node {
	var signature *checker.Signature
	callLike := getAncestorCallLikeExpression(node)
	if callLike != nil {
		signature = typeChecker.GetResolvedSignature(callLike)
	}
	// Don't go to a function type, go to the value having that type.
	var declaration *ast.Node
	if signature != nil && signature.Declaration() != nil {
		declaration = signature.Declaration()
		if ast.IsFunctionLike(declaration) && !ast.IsFunctionTypeNode(declaration) {
			return declaration
		}
	}
	return nil
}

func getSymbolForOverriddenMember(typeChecker *checker.Checker, node *ast.Node) *ast.Symbol {
	classElement := ast.FindAncestor(node, ast.IsClassElement)
	if classElement == nil || classElement.Name() == nil {
		return nil
	}
	baseDeclaration := ast.FindAncestor(classElement, ast.IsClassLike)
	if baseDeclaration == nil {
		return nil
	}
	baseTypeNode := ast.GetClassExtendsHeritageElement(baseDeclaration)
	if baseTypeNode == nil {
		return nil
	}
	expression := ast.SkipParentheses(baseTypeNode.Expression())
	var base *ast.Symbol
	if ast.IsClassExpression(expression) {
		base = expression.Symbol()
	} else {
		base = typeChecker.GetSymbolAtLocation(expression)
	}
	if base == nil {
		return nil
	}
	name := ast.GetTextOfPropertyName(classElement.Name())
	if ast.HasStaticModifier(classElement) {
		return typeChecker.GetPropertyOfType(typeChecker.GetTypeOfSymbol(base), name)
	}
	return typeChecker.GetPropertyOfType(typeChecker.GetDeclaredTypeOfSymbol(base), name)
}

func getTypeOfSymbolAtLocation(c *checker.Checker, symbol *ast.Symbol, node *ast.Node) *checker.Type {
	t := c.GetTypeOfSymbolAtLocation(symbol, node)
	// If the type is just a function's inferred type, go-to-type should go to the return type instead since
	// go-to-definition takes you to the function anyway.
	if t.Symbol() == symbol || t.Symbol() != nil && symbol.ValueDeclaration != nil && ast.IsVariableDeclaration(symbol.ValueDeclaration) && symbol.ValueDeclaration.Initializer() == t.Symbol().ValueDeclaration {
		sigs := c.GetCallSignatures(t)
		if len(sigs) == 1 {
			return c.GetReturnTypeOfSignature(sigs[0])
		}
	}
	return t
}

func getDeclarationsFromType(t *checker.Type) []*ast.Node {
	var result []*ast.Node
	for _, t := range t.Distributed() {
		if t.Symbol() != nil {
			for _, decl := range t.Symbol().Declarations {
				result = core.AppendIfUnique(result, decl)
			}
		}
	}
	return result
}

func (l *LanguageService) tryMapToOriginalSource(declFile *ast.SourceFile, node *ast.Node) *lsproto.Location {
	fs := l.GetProgram().Host().FS()

	declFileName := declFile.FileName()
	declContent, ok := fs.ReadFile(declFileName)
	if !ok {
		return nil
	}

	lineMap := l.converters.getLineMap(declFileName)
	if lineMap == nil {
		lineMap = ComputeLineStarts(declContent)
	}
	lineInfo := sourcemap.GetLineInfo(declContent, lineMap.LineStarts)

	sourceMappingURL := sourcemap.TryGetSourceMappingURL(lineInfo)
	if sourceMappingURL == "" || strings.HasPrefix(sourceMappingURL, "data:") {
		return nil
	}

	sourceMapPath := tspath.NormalizePath(filepath.Join(filepath.Dir(declFileName), sourceMappingURL))
	sourceMapContent, ok := fs.ReadFile(sourceMapPath)
	if !ok {
		return nil
	}

	var sourceMapData struct {
		SourceRoot string   `json:"sourceRoot"`
		Sources    []string `json:"sources"`
		Mappings   string   `json:"mappings"`
	}
	if err := json.Unmarshal([]byte(sourceMapContent), &sourceMapData); err != nil {
		return nil
	}

	decoder := sourcemap.DecodeMappings(sourceMapData.Mappings)
	if decoder.Error() != nil {
		return nil
	}

	declPosition := l.converters.PositionToLineAndCharacter(declFile, core.TextPos(node.Pos()))

	var bestMapping *sourcemap.Mapping
	for mapping := range decoder.Values() {
		if mapping.GeneratedLine == int(declPosition.Line) &&
			mapping.GeneratedCharacter <= int(declPosition.Character) &&
			mapping.IsSourceMapping() {
			if bestMapping == nil || mapping.GeneratedCharacter > bestMapping.GeneratedCharacter {
				bestMapping = mapping
			}
		}
	}

	if bestMapping == nil || int(bestMapping.SourceIndex) >= len(sourceMapData.Sources) {
		return nil
	}

	sourceFileName := sourceMapData.Sources[bestMapping.SourceIndex]
	if !filepath.IsAbs(sourceFileName) {
		if sourceMapData.SourceRoot != "" {
			sourceFileName = filepath.Join(sourceMapData.SourceRoot, sourceFileName)
		}
		sourceFileName = tspath.NormalizePath(filepath.Join(filepath.Dir(declFileName), sourceFileName))
	}

	if !fs.FileExists(sourceFileName) {
		return nil
	}

	sourceContent, ok := fs.ReadFile(sourceFileName)
	if !ok {
		return nil
	}

	sourceFileScript := &sourceFileScript{
		fileName: sourceFileName,
		text:     sourceContent,
		lineMap:  ComputeLineStarts(sourceContent).LineStarts,
	}

	sourceLspPosition := lsproto.Position{
		Line:      uint32(bestMapping.SourceLine),
		Character: uint32(bestMapping.SourceCharacter),
	}

	var sourceStartLsp, sourceEndLsp lsproto.Position

	var symbolName string
	if node.Kind == ast.KindIdentifier || node.Kind == ast.KindPrivateIdentifier {
		symbolName = node.Text()
	}
	if symbolName != "" {
		sourceBytePos := l.converters.LineAndCharacterToPosition(sourceFileScript, sourceLspPosition)
		if symbolStart := findSymbolNearPosition(sourceContent, symbolName, int(sourceBytePos)); symbolStart != -1 {
			sourceStartPos := core.TextPos(symbolStart)
			sourceEndPos := core.TextPos(symbolStart + len(symbolName))
			sourceStartLsp = l.converters.PositionToLineAndCharacter(sourceFileScript, sourceStartPos)
			sourceEndLsp = l.converters.PositionToLineAndCharacter(sourceFileScript, sourceEndPos)
		}
	}

	if sourceStartLsp == (lsproto.Position{}) {
		sourceStartLsp = sourceLspPosition
		sourceEndLsp = sourceLspPosition
	}

	return &lsproto.Location{
		Uri: FileNameToDocumentURI(sourceFileName),
		Range: lsproto.Range{
			Start: sourceStartLsp,
			End:   sourceEndLsp,
		},
	}
}

type sourceFileScript struct {
	fileName string
	text     string
	lineMap  []core.TextPos
}

func (s *sourceFileScript) FileName() string {
	return s.fileName
}

func (s *sourceFileScript) Text() string {
	return s.text
}

func (s *sourceFileScript) LineMap() []core.TextPos {
	return s.lineMap
}

func findSymbolNearPosition(text, symbolName string, targetPos int) int {
	if symbolName == "" {
		return -1
	}

	symbolLen := len(symbolName)
	textLen := len(text)
	bestMatch := -1
	bestDistance := math.MaxInt

	pos := strings.Index(text, symbolName)
	for pos >= 0 {
		if (pos == 0 || !scanner.IsIdentifierPart(rune(text[pos-1]))) &&
			(pos+symbolLen >= textLen || !scanner.IsIdentifierPart(rune(text[pos+symbolLen]))) {

			distance := targetPos - pos
			if distance < 0 {
				distance = -distance
			}

			if distance < bestDistance {
				bestDistance = distance
				bestMatch = pos
			}
		}

		nextPos := strings.Index(text[pos+1:], symbolName)
		if nextPos >= 0 {
			pos = pos + 1 + nextPos
		} else {
			pos = -1
		}
	}

	return bestMatch
}
