package ls

import (
	"bufio"
	"context"
	"os"
	"regexp"
	"slices"
	"strings"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/astnav"
	"github.com/microsoft/typescript-go/internal/checker"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/scanner"
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
		
		// Try to find the corresponding source file if this is a declaration file
		// Use the file path directly so it works even for external declaration files
		sourceFile := l.tryFindSourceFileForDeclarationByPath(file.FileName())
		
		// If we found a source file via project references, use it
		if sourceFile != nil {
			// Try to find the corresponding declaration in the source file
			sourceDecl := l.tryFindDeclarationInSourceFile(decl, sourceFile)
			
			if sourceDecl != nil {
				// Use the found declaration in the source file
				sourceName := core.OrElse(ast.GetNameOfDeclaration(sourceDecl), sourceDecl)
				locations = core.AppendIfUnique(locations, lsproto.Location{
					Uri:   FileNameToDocumentURI(sourceFile.FileName()),
					Range: *l.createLspRangeFromNode(sourceName, sourceFile),
				})
			} else {
				// Fallback: point to the beginning of the source file
				// This is not perfect but ensures we point to the right file
				locations = core.AppendIfUnique(locations, lsproto.Location{
					Uri:   FileNameToDocumentURI(sourceFile.FileName()),
					Range: *l.createLspRangeFromBounds(0, 0, sourceFile),
				})
			}
		} else {
			// Fallback: Try heuristic approach for monorepos without project references
			if candidateSourcePath := l.tryFindSourceFilePathHeuristic(file.FileName()); candidateSourcePath != "" {
				// Try to find the approximate location of the symbol in the source file
				if symbolLocation := l.tryFindSymbolLocationInFile(decl, candidateSourcePath); symbolLocation != nil {
					locations = core.AppendIfUnique(locations, lsproto.Location{
						Uri:   FileNameToDocumentURI(candidateSourcePath),
						Range: *symbolLocation,
					})
				} else {
					// Fallback: point to the beginning of the source file
					locations = core.AppendIfUnique(locations, lsproto.Location{
						Uri:   FileNameToDocumentURI(candidateSourcePath),
						Range: lsproto.Range{
							Start: lsproto.Position{Line: 0, Character: 0},
							End:   lsproto.Position{Line: 0, Character: 0},
						},
					})
				}
			} else {
				// Use the original approach for non-declaration files or when no mapping found
				locations = core.AppendIfUnique(locations, lsproto.Location{
					Uri:   FileNameToDocumentURI(file.FileName()),
					Range: *l.createLspRangeFromNode(name, file),
				})
			}
		}
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

// tryFindSourceFileForDeclaration attempts to find the corresponding source file
// for a declaration file using the project reference system.
//
// This handles monorepo scenarios where packages have aliases like 
// @alias/package that resolve to packages/package/dist/*.d.ts files.
// The project reference system maps these back to the original source files.
func (l *LanguageService) tryFindSourceFileForDeclaration(declFile *ast.SourceFile) *ast.SourceFile {
	return l.tryFindSourceFileForDeclarationByPath(declFile.FileName())
}

// tryFindSourceFileForDeclarationByPath attempts to find the corresponding source file
// for a declaration file path using the project reference system.
// This works even for external declaration files that aren't loaded in the current program.
func (l *LanguageService) tryFindSourceFileForDeclarationByPath(declFilePath string) *ast.SourceFile {
	// If not a declaration file, return nil (caller should use original)
	if !strings.HasSuffix(declFilePath, ".d.ts") {
		return nil
	}

	program := l.GetProgram()

	// Convert the file path to a tspath.Path for lookup
	path := tspath.ToPath(declFilePath, program.Host().GetCurrentDirectory(), program.Host().FS().UseCaseSensitiveFileNames())
	
	// Strategy 1: Use the project reference system to find the source file
	// This works when TypeScript project references are properly configured
	if projectRef := program.GetProjectReferenceFromOutputDts(path); projectRef != nil && projectRef.Source != "" {
		if sourceFile := program.GetSourceFile(projectRef.Source); sourceFile != nil {
			return sourceFile
		}
	}
	
	// Strategy 2: Fallback heuristic mapping for monorepos without project references
	// This handles common patterns like: packages/pkg/dist/file.d.ts -> packages/pkg/src/file.ts
	if sourceFile := l.tryHeuristicSourceMapping(declFilePath, program); sourceFile != nil {
		return sourceFile
	}
	
	return nil // No mapping found
}

// tryHeuristicSourceMapping attempts to map declaration files to source files using common patterns
// when project references aren't configured. This handles typical monorepo structures.
func (l *LanguageService) tryHeuristicSourceMapping(declFilePath string, program *compiler.Program) *ast.SourceFile {
	// Common monorepo patterns to try
	mappings := []struct {
		buildPattern string
		sourcePattern string
	}{
		// Most common: packages/pkg/dist/ -> packages/pkg/src/
		{"/dist/", "/src/"},
		// Alternative patterns
		{"/lib/", "/src/"},
		{"/build/", "/src/"},
		{"/out/", "/src/"},
		// Sometimes dist is alongside src
		{"/dist/", "/"},
	}
	
	for _, mapping := range mappings {
		if strings.Contains(declFilePath, mapping.buildPattern) {
			// Replace the build directory with source directory
			sourceDir := strings.Replace(declFilePath, mapping.buildPattern, mapping.sourcePattern, 1)
			
			// Try different source file extensions
			baseName := strings.TrimSuffix(sourceDir, ".d.ts")
			candidates := []string{
				baseName + ".ts",
				baseName + ".tsx",
				baseName + ".js",
				baseName + ".jsx",
			}
			
			for _, candidate := range candidates {
				if sourceFile := program.GetSourceFile(candidate); sourceFile != nil {
					return sourceFile
				}
			}
		}
	}
	
	return nil
}

// tryFindSourceFilePathHeuristic attempts to find source files on the filesystem using common patterns
// This is used when project references aren't configured and source files aren't loaded in the program
func (l *LanguageService) tryFindSourceFilePathHeuristic(declFilePath string) string {
	// Only process declaration files
	if !strings.HasSuffix(declFilePath, ".d.ts") {
		return ""
	}
	
	// Common monorepo patterns to try
	mappings := []struct {
		buildPattern  string
		sourcePattern string
	}{
		// Most common: packages/pkg/dist/ -> packages/pkg/src/
		{"/dist/", "/src/"},
		// Alternative patterns
		{"/lib/", "/src/"},
		{"/build/", "/src/"},
		{"/out/", "/src/"},
		// Sometimes dist is alongside src
		{"/dist/", "/"},
	}
	
	for _, mapping := range mappings {
		if strings.Contains(declFilePath, mapping.buildPattern) {
			// Replace the build directory with source directory
			sourceDir := strings.Replace(declFilePath, mapping.buildPattern, mapping.sourcePattern, 1)
			
			// Try different source file extensions
			baseName := strings.TrimSuffix(sourceDir, ".d.ts")
			candidates := []string{
				baseName + ".ts",
				baseName + ".tsx",
				baseName + ".js",
				baseName + ".jsx",
			}
			
			// Check if any candidate exists on the filesystem
			for _, candidate := range candidates {
				if _, err := os.Stat(candidate); err == nil {
					// File exists! Return this path
					return candidate
				}
			}
		}
	}
	
	return ""
}

// tryLoadSourceFile attempts to load a source file into the TypeScript program
// This allows us to analyze the file and find specific declarations within it
func (l *LanguageService) tryLoadSourceFile(sourceFilePath string) *ast.SourceFile {
	program := l.GetProgram()
	
	// First, try to get the file if it's already loaded in the program
	if existingFile := program.GetSourceFile(sourceFilePath); existingFile != nil {
		return existingFile
	}
	
	// For now, return nil - dynamic loading is complex and may not be necessary
	// since the TypeScript server can handle the source file once we navigate to it
	return nil
}

// tryFindSymbolLocationInFile attempts to find the location of a symbol in a source file
// using text-based search patterns. This provides approximate location when AST analysis isn't available.
func (l *LanguageService) tryFindSymbolLocationInFile(declNode *ast.Node, sourceFilePath string) *lsproto.Range {
	// Get the name of the symbol we're looking for
	symbolName := l.getSymbolName(declNode)
	if symbolName == "" {
		return nil
	}
	
	// Read the source file
	file, err := os.Open(sourceFilePath)
	if err != nil {
		return nil
	}
	defer file.Close()
	
	scanner := bufio.NewScanner(file)
	lineNumber := 0
	
	// Common TypeScript/JavaScript declaration patterns
	patterns := []*regexp.Regexp{
		// export function funcName(
		regexp.MustCompile(`export\s+function\s+` + regexp.QuoteMeta(symbolName) + `\s*\(`),
		// export const funcName = 
		regexp.MustCompile(`export\s+const\s+` + regexp.QuoteMeta(symbolName) + `\s*[=:]`),
		// function funcName(
		regexp.MustCompile(`function\s+` + regexp.QuoteMeta(symbolName) + `\s*\(`),
		// const funcName = 
		regexp.MustCompile(`const\s+` + regexp.QuoteMeta(symbolName) + `\s*[=:]`),
		// let funcName = 
		regexp.MustCompile(`let\s+` + regexp.QuoteMeta(symbolName) + `\s*[=:]`),
		// var funcName = 
		regexp.MustCompile(`var\s+` + regexp.QuoteMeta(symbolName) + `\s*[=:]`),
		// class ClassName
		regexp.MustCompile(`class\s+` + regexp.QuoteMeta(symbolName) + `\s*[{<]`),
		// interface InterfaceName
		regexp.MustCompile(`interface\s+` + regexp.QuoteMeta(symbolName) + `\s*[{<]`),
		// type TypeName
		regexp.MustCompile(`type\s+` + regexp.QuoteMeta(symbolName) + `\s*[=<]`),
	}
	
	for scanner.Scan() {
		line := scanner.Text()
		
		// Check each pattern
		for _, pattern := range patterns {
			if match := pattern.FindStringIndex(line); match != nil {
				// Found a match! Calculate the position
				character := match[0]
				
				// Try to find the exact position of the symbol name within the match
				symbolIndex := strings.Index(line[character:], symbolName)
				if symbolIndex != -1 {
					character += symbolIndex
				}
				
				return &lsproto.Range{
					Start: lsproto.Position{
						Line:      uint32(lineNumber),
						Character: uint32(character),
					},
					End: lsproto.Position{
						Line:      uint32(lineNumber),
						Character: uint32(character + len(symbolName)),
					},
				}
			}
		}
		
		lineNumber++
	}
	
	return nil
}

// getSymbolName extracts the symbol name from a declaration node
func (l *LanguageService) getSymbolName(declNode *ast.Node) string {
	if nameNode := ast.GetNameOfDeclaration(declNode); nameNode != nil {
		return nameNode.Text()
	}
	
	// Fallback: try to extract from the node text
	nodeText := strings.TrimSpace(declNode.Text())
	
	// Simple patterns to extract symbol names
	patterns := []string{
		`function\s+(\w+)`,
		`const\s+(\w+)`,
		`let\s+(\w+)`,
		`var\s+(\w+)`,
		`class\s+(\w+)`,
		`interface\s+(\w+)`,
		`type\s+(\w+)`,
	}
	
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		if matches := re.FindStringSubmatch(nodeText); len(matches) > 1 {
			return matches[1]
		}
	}
	
	return ""
}

// tryFindDeclarationInSourceFile attempts to find a declaration with the same name
// in the source file. This is a simple heuristic that may not always be accurate.
func (l *LanguageService) tryFindDeclarationInSourceFile(declNode *ast.Node, sourceFile *ast.SourceFile) *ast.Node {
	// Get the name of the original declaration
	originalName := ast.GetNameOfDeclaration(declNode)
	if originalName == nil {
		return nil
	}
	
	// For now, return nil to use the fallback approach
	// TODO: Implement proper symbol matching in the source file
	// This would require walking the AST and finding declarations with matching names
	return nil
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
