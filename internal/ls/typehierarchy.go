package ls

import (
	"context"
	"slices"
	"strings"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/astnav"
	"github.com/microsoft/typescript-go/internal/checker"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls/lsconv"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/scanner"
)

// TypeHierarchyDeclaration represents a node that can appear in the type hierarchy.
// This includes classes, interfaces, type aliases, and mixin variables.
type TypeHierarchyDeclaration = *ast.Node

// isPossibleTypeHierarchyDeclaration indicates whether a node could possibly be a type hierarchy declaration.
func isPossibleTypeHierarchyDeclaration(node *ast.Node) bool {
	if node == nil {
		return false
	}
	return ast.IsClassDeclaration(node) ||
		ast.IsClassExpression(node) ||
		ast.IsInterfaceDeclaration(node) ||
		ast.IsTypeAliasDeclaration(node) ||
		ast.IsTypeParameterDeclaration(node) ||
		isTypeHierarchyMixinVariable(node)
}

// isValidTypeHierarchyDeclaration indicates whether a node is a valid type hierarchy declaration.
func isValidTypeHierarchyDeclaration(node *ast.Node) bool {
	if node == nil {
		return false
	}

	// Classes - both named and expressions with names
	if ast.IsClassDeclaration(node) || ast.IsClassExpression(node) {
		return true
	}

	// Interfaces
	if ast.IsInterfaceDeclaration(node) {
		return true
	}

	// Type aliases
	if ast.IsTypeAliasDeclaration(node) {
		return true
	}

	// Type parameters with constraints
	if ast.IsTypeParameterDeclaration(node) {
		return node.AsTypeParameter().Constraint != nil
	}

	// Mixin variables (const Mixed = Mixin(Base))
	if isTypeHierarchyMixinVariable(node) {
		return true
	}

	return false
}

// isTypeHierarchyMixinVariable checks if a node is a mixin variable pattern like `const Mixed = Mixin(Base)`.
func isTypeHierarchyMixinVariable(node *ast.Node) bool {
	if node == nil {
		return false
	}

	// Must be a variable declaration
	if !ast.IsVariableDeclaration(node) {
		return false
	}

	// Must have an initializer
	initializer := node.Initializer()
	if initializer == nil {
		return false
	}

	// Must be const or readonly
	if !((ast.GetCombinedNodeFlags(node)&ast.NodeFlagsConst) != 0 || ast.IsPropertyDeclaration(node)) {
		return false
	}

	// Initializer should be a call expression (mixin function call)
	return ast.IsCallExpression(initializer) && isMixinLikeReturnType(initializer)
}

// isMixinLikeReturnType checks if a call expression returns a class-like type (mixin pattern).
func isMixinLikeReturnType(node *ast.Node) bool {
	// Simplified check: if it's a call expression that takes a class as argument
	// In a full implementation, we'd check the return type
	if !ast.IsCallExpression(node) {
		return false
	}
	callExpr := node.AsCallExpression()
	// A mixin usually has at least one argument (the base class)
	return len(callExpr.Arguments.Nodes) > 0
}

// getTypeHierarchyDeclarationReferenceNode gets the node that can be used as a reference to a type hierarchy declaration.
func getTypeHierarchyDeclarationReferenceNode(node *ast.Node) *ast.Node {
	if node == nil {
		return nil
	}

	if name := node.Name(); name != nil {
		return name
	}

	// For mixin variables
	if ast.IsVariableDeclaration(node) {
		return node.Name()
	}

	return node
}

// resolveTypeHierarchyDeclaration resolves the type hierarchy declaration at the given node.
func resolveTypeHierarchyDeclaration(program *compiler.Program, node *ast.Node) TypeHierarchyDeclaration {
	if node == nil {
		return nil
	}

	// Walk up to find the containing declaration
	for node != nil {
		if isValidTypeHierarchyDeclaration(node) {
			return node
		}

		// Check if we're on an identifier that references a type
		if ast.IsIdentifier(node) || ast.IsPropertyAccessExpression(node) {
			parent := node.Parent
			if parent != nil {
				// If we're in a heritage clause, resolve the referenced type
				if ast.IsExpressionWithTypeArguments(parent) {
					c, done := program.GetTypeChecker(context.Background())
					defer done()
					symbol := c.GetSymbolAtLocation(node)
					if symbol != nil {
						decl := getTypeDeclarationFromSymbol(symbol)
						if decl != nil && isValidTypeHierarchyDeclaration(decl) {
							return decl
						}
					}
				}
			}
		}

		// If we're on the name of a declaration, return the declaration
		if node.Parent != nil && isValidTypeHierarchyDeclaration(node.Parent) {
			if name := node.Parent.Name(); name == node {
				return node.Parent
			}
		}

		node = node.Parent
	}

	return nil
}

// getTypeDeclarationFromSymbol gets the type declaration from a symbol.
func getTypeDeclarationFromSymbol(symbol *ast.Symbol) *ast.Node {
	if symbol == nil {
		return nil
	}

	// Handle aliased symbols
	if (symbol.Flags & ast.SymbolFlagsAlias) != 0 {
		return nil // Let the caller resolve aliases if needed
	}

	decls := symbol.Declarations
	if len(decls) == 0 {
		return nil
	}

	// Return the first declaration that's a type hierarchy declaration
	for _, decl := range decls {
		if isValidTypeHierarchyDeclaration(decl) {
			return decl
		}
	}

	return decls[0]
}

// getSymbolKindForTypeHierarchy determines the LSP SymbolKind for a type hierarchy item.
// Note: LSP doesn't have a dedicated TypeAlias kind, so we use Struct which is a better
// representation than TypeParameter (used by some implementations). TypeParameter (26)
// is specifically for generic type parameters like T, U, not for type aliases.
func getSymbolKindForTypeHierarchy(node *ast.Node) lsproto.SymbolKind {
	if node == nil {
		return lsproto.SymbolKindClass
	}

	switch {
	case ast.IsClassDeclaration(node) || ast.IsClassExpression(node):
		return lsproto.SymbolKindClass
	case ast.IsInterfaceDeclaration(node):
		return lsproto.SymbolKindInterface
	case ast.IsTypeAliasDeclaration(node):
		// LSP doesn't have TypeAlias, use Struct as it's closer semantically than TypeParameter.
		// TypeParameter (26) is for generic type params (T, U), not type aliases.
		// Struct (23) represents a compound type which is what type aliases often define.
		return lsproto.SymbolKindStruct
	case ast.IsTypeParameterDeclaration(node):
		return lsproto.SymbolKindTypeParameter
	case isTypeHierarchyMixinVariable(node):
		return lsproto.SymbolKindClass // Mixins are class-like
	default:
		return lsproto.SymbolKindClass
	}
}

// getTypeHierarchyItemName gets the name for a type hierarchy item.
func getTypeHierarchyItemName(node *ast.Node) string {
	if node == nil {
		return "<anonymous>"
	}

	if name := node.Name(); name != nil && ast.IsIdentifier(name) {
		return name.Text()
	}

	if ast.IsTypeAliasDeclaration(node) {
		return node.AsTypeAliasDeclaration().Name().Text()
	}

	if ast.IsClassDeclaration(node) && node.AsClassDeclaration().Name() != nil {
		return node.AsClassDeclaration().Name().Text()
	}

	if ast.IsInterfaceDeclaration(node) {
		return node.AsInterfaceDeclaration().Name().Text()
	}

	return "<anonymous>"
}

// getTypeHierarchyKindModifiers returns additional kind modifiers for a type hierarchy declaration.
// These help distinguish different kinds of type relationships in the UI.
// Modifiers are returned as comma-separated strings (e.g., "abstract", "conditional,extends").
func getTypeHierarchyKindModifiers(node *ast.Node) string {
	if node == nil {
		return ""
	}

	var modifiers []string

	// Check for abstract classes
	if ast.IsClassDeclaration(node) || ast.IsClassExpression(node) {
		if (ast.GetCombinedModifierFlags(node) & ast.ModifierFlagsAbstract) != 0 {
			modifiers = append(modifiers, "abstract")
		}
	}

	// Check for mixin variables
	if isTypeHierarchyMixinVariable(node) {
		modifiers = append(modifiers, "mixin")
	}

	// Check for type alias specific modifiers
	if ast.IsTypeAliasDeclaration(node) {
		typeNode := node.AsTypeAliasDeclaration().Type
		if typeNode != nil {
			switch typeNode.Kind {
			case ast.KindConditionalType:
				modifiers = append(modifiers, "conditional")
				// Check if it uses infer keyword
				if containsInferType(typeNode) {
					modifiers = append(modifiers, "infer")
				} else {
					modifiers = append(modifiers, "extends")
				}
			case ast.KindIntersectionType:
				modifiers = append(modifiers, "intersection")
			case ast.KindUnionType:
				modifiers = append(modifiers, "union")
			case ast.KindMappedType:
				modifiers = append(modifiers, "mapped")
			case ast.KindTupleType:
				modifiers = append(modifiers, "tuple")
			case ast.KindTemplateLiteralType:
				modifiers = append(modifiers, "template")
			case ast.KindIndexedAccessType:
				modifiers = append(modifiers, "indexed")
			case ast.KindTypeOperator:
				typeOp := typeNode.AsTypeOperatorNode()
				switch typeOp.Operator {
				case ast.KindKeyOfKeyword:
					modifiers = append(modifiers, "keyof")
				case ast.KindReadonlyKeyword:
					modifiers = append(modifiers, "readonly")
				default:
					modifiers = append(modifiers, "alias")
				}
			case ast.KindTypeReference:
				// Simple type alias (type Foo = Bar)
				modifiers = append(modifiers, "alias")
			}
		}
	}

	return strings.Join(modifiers, ",")
}

// containsInferType checks if a type node contains an infer type.
func containsInferType(node *ast.Node) bool {
	if node == nil {
		return false
	}

	if node.Kind == ast.KindInferType {
		return true
	}

	// Recursively check children
	found := false
	node.ForEachChild(func(child *ast.Node) bool {
		if containsInferType(child) {
			found = true
			return true // Stop iteration
		}
		return false // Continue iteration
	})

	return found
}

// createTypeHierarchyItem creates an LSP TypeHierarchyItem for the given declaration.
func (l *LanguageService) createTypeHierarchyItem(program *compiler.Program, declaration *ast.Node) *lsproto.TypeHierarchyItem {
	if declaration == nil {
		return nil
	}

	sourceFile := ast.GetSourceFileOfNode(declaration)
	if sourceFile == nil {
		return nil
	}

	script := l.getScript(sourceFile.FileName())
	if script == nil {
		return nil
	}

	name := getTypeHierarchyItemName(declaration)
	kind := getSymbolKindForTypeHierarchy(declaration)

	// Get the range of the entire declaration, skipping leading trivia (whitespace, comments)
	// This matches TypeScript's behavior of using skipTrivia with stopAtComments: true
	startPos := scanner.SkipTriviaEx(sourceFile.Text(), declaration.Pos(), &scanner.SkipTriviaOptions{StopAtComments: true})
	endPos := declaration.End()
	range_ := l.converters.ToLSPRange(script, core.NewTextRange(startPos, endPos))

	// Get the selection range (usually just the name), also skipping trivia
	refNode := getTypeHierarchyDeclarationReferenceNode(declaration)
	var selectionRange lsproto.Range
	if refNode != nil {
		nameStart := scanner.SkipTrivia(sourceFile.Text(), refNode.Pos())
		selectionRange = l.converters.ToLSPRange(script, core.NewTextRange(nameStart, refNode.End()))
	} else {
		selectionRange = range_
	}

	// Get detail (type signature)
	var detail *string
	c, done := program.GetTypeChecker(context.Background())
	symbol := c.GetSymbolAtLocation(declaration)
	if symbol != nil {
		t := c.GetDeclaredTypeOfSymbol(symbol)
		if t != nil {
			typeStr := c.TypeToString(t)
			detail = &typeStr
		}
	}
	done()

	return &lsproto.TypeHierarchyItem{
		Name:           name,
		Kind:           kind,
		Detail:         detail,
		Uri:            lsconv.FileNameToDocumentURI(sourceFile.FileName()),
		Range:          range_,
		SelectionRange: selectionRange,
		Data:           &lsproto.TypeHierarchyItemData{
			// Custom data preserved between requests
		},
	}
}

// resolveTypeHierarchyDeclarationAtPosition resolves the type hierarchy declaration
// at the given position in a source file. This helper function extracts common logic
// used by supertypes and subtypes handlers to avoid code duplication.
func (l *LanguageService) resolveTypeHierarchyDeclarationAtPosition(
	program *compiler.Program,
	file *ast.SourceFile,
	pos int,
) TypeHierarchyDeclaration {
	var node *ast.Node
	if pos == 0 {
		node = file.AsNode()
	} else {
		node = astnav.GetTouchingPropertyName(file, pos)
	}

	if node == nil {
		return nil
	}

	return resolveTypeHierarchyDeclaration(program, node)
}

// ProvidePrepareTypeHierarchy prepares the type hierarchy at the given position.
func (l *LanguageService) ProvidePrepareTypeHierarchy(
	ctx context.Context,
	documentURI lsproto.DocumentUri,
	position lsproto.Position,
) (lsproto.TypeHierarchyPrepareResponse, error) {
	program, file := l.getProgramAndFile(documentURI)
	pos := int(l.converters.LineAndCharacterToPosition(file, position))
	declaration := l.resolveTypeHierarchyDeclarationAtPosition(program, file, pos)
	if declaration == nil || declaration.Kind == ast.KindSourceFile {
		return lsproto.TypeHierarchyItemsOrNull{}, nil
	}

	hierarchyItem := l.createTypeHierarchyItem(program, declaration)
	if hierarchyItem == nil {
		return lsproto.TypeHierarchyItemsOrNull{}, nil
	}

	items := []*lsproto.TypeHierarchyItem{hierarchyItem}
	return lsproto.TypeHierarchyItemsOrNull{TypeHierarchyItems: &items}, nil
}

// ProvideTypeHierarchySupertypes gets the supertypes of a type hierarchy item.
func (l *LanguageService) ProvideTypeHierarchySupertypes(
	ctx context.Context,
	item *lsproto.TypeHierarchyItem,
) (lsproto.TypeHierarchySupertypesResponse, error) {
	program := l.GetProgram()
	fileName := item.Uri.FileName()
	file := program.GetSourceFile(fileName)
	if file == nil {
		return lsproto.TypeHierarchyItemsOrNull{}, nil
	}

	pos := int(l.converters.LineAndCharacterToPosition(file, item.SelectionRange.Start))
	declaration := l.resolveTypeHierarchyDeclarationAtPosition(program, file, pos)
	if declaration == nil {
		return lsproto.TypeHierarchyItemsOrNull{}, nil
	}

	supertypes := l.getSupertypes(program, declaration)
	if len(supertypes) == 0 {
		return lsproto.TypeHierarchyItemsOrNull{}, nil
	}

	return lsproto.TypeHierarchyItemsOrNull{TypeHierarchyItems: &supertypes}, nil
}

// ProvideTypeHierarchySubtypes gets the subtypes of a type hierarchy item.
func (l *LanguageService) ProvideTypeHierarchySubtypes(
	ctx context.Context,
	item *lsproto.TypeHierarchyItem,
	orchestrator CrossProjectOrchestrator,
) (lsproto.TypeHierarchySubtypesResponse, error) {
	program := l.GetProgram()
	fileName := item.Uri.FileName()
	file := program.GetSourceFile(fileName)
	if file == nil {
		return lsproto.TypeHierarchyItemsOrNull{}, nil
	}

	pos := int(l.converters.LineAndCharacterToPosition(file, item.SelectionRange.Start))
	declaration := l.resolveTypeHierarchyDeclarationAtPosition(program, file, pos)
	if declaration == nil {
		return lsproto.TypeHierarchyItemsOrNull{}, nil
	}

	subtypes := l.getSubtypes(ctx, program, declaration, orchestrator)
	if len(subtypes) == 0 {
		return lsproto.TypeHierarchyItemsOrNull{}, nil
	}

	return lsproto.TypeHierarchyItemsOrNull{TypeHierarchyItems: &subtypes}, nil
}

// getSupertypes collects all supertypes of a declaration.
func (l *LanguageService) getSupertypes(program *compiler.Program, declaration *ast.Node) []*lsproto.TypeHierarchyItem {
	if declaration == nil {
		return nil
	}

	var results []*lsproto.TypeHierarchyItem
	seen := make(map[*ast.Node]bool)

	c, done := program.GetTypeChecker(context.Background())
	defer done()

	switch {
	case ast.IsClassDeclaration(declaration) || ast.IsClassExpression(declaration):
		// Get base class
		if baseType := getEffectiveBaseTypeNode(declaration); baseType != nil {
			if baseDecl := resolveTypeReferenceToDeclaration(c, baseType); baseDecl != nil && !seen[baseDecl] {
				seen[baseDecl] = true
				if item := l.createTypeHierarchyItem(program, baseDecl); item != nil {
					results = append(results, item)
				}
			}
		}

		// Get implemented interfaces
		for _, heritage := range getHeritageClausesWithKind(declaration, ast.KindImplementsKeyword) {
			for _, typeRef := range heritage.Types.Nodes {
				if decl := resolveTypeReferenceToDeclaration(c, typeRef); decl != nil && !seen[decl] {
					seen[decl] = true
					if item := l.createTypeHierarchyItem(program, decl); item != nil {
						results = append(results, item)
					}
				}
			}
		}

	case ast.IsInterfaceDeclaration(declaration):
		// Get extended interfaces
		for _, heritage := range getHeritageClausesWithKind(declaration, ast.KindExtendsKeyword) {
			for _, typeRef := range heritage.Types.Nodes {
				if decl := resolveTypeReferenceToDeclaration(c, typeRef); decl != nil && !seen[decl] {
					seen[decl] = true
					if item := l.createTypeHierarchyItem(program, decl); item != nil {
						results = append(results, item)
					}
				}
			}
		}

	case ast.IsTypeAliasDeclaration(declaration):
		// Analyze the type alias body for referenced types
		typeNode := declaration.AsTypeAliasDeclaration().Type
		if typeNode != nil {
			collectReferencedTypesFromTypeNode(c, typeNode, seen, func(decl *ast.Node) {
				if item := l.createTypeHierarchyItem(program, decl); item != nil {
					results = append(results, item)
				}
			})
		}

	case ast.IsTypeParameterDeclaration(declaration):
		// Get constraint
		if constraint := declaration.AsTypeParameter().Constraint; constraint != nil {
			if decl := resolveTypeReferenceToDeclaration(c, constraint); decl != nil && !seen[decl] {
				seen[decl] = true
				if item := l.createTypeHierarchyItem(program, decl); item != nil {
					results = append(results, item)
				}
			}
		}

	case isTypeHierarchyMixinVariable(declaration):
		// Get the mixin chain
		collectMixinChain(c, declaration, seen, func(decl *ast.Node) {
			if item := l.createTypeHierarchyItem(program, decl); item != nil {
				results = append(results, item)
			}
		})
	}

	// Sort results by name for consistent ordering
	slices.SortFunc(results, func(a, b *lsproto.TypeHierarchyItem) int {
		return strings.Compare(a.Name, b.Name)
	})

	return results
}

// getSubtypes collects all subtypes of a declaration.
// This is a simplified implementation that doesn't use cross-project references.
// A full implementation would use handleCrossProject similar to callhierarchy.go.
func (l *LanguageService) getSubtypes(
	ctx context.Context,
	program *compiler.Program,
	declaration *ast.Node,
	orchestrator CrossProjectOrchestrator,
) []*lsproto.TypeHierarchyItem {
	if declaration == nil {
		return nil
	}

	var results []*lsproto.TypeHierarchyItem
	seen := make(map[*ast.Node]bool)

	c, done := program.GetTypeChecker(context.Background())
	defer done()

	// Get the symbol for the declaration - use the name node for GetSymbolAtLocation
	nameNode := declaration.Name()
	if nameNode == nil {
		return nil
	}
	symbol := c.GetSymbolAtLocation(nameNode)
	if symbol == nil {
		return nil
	}

	// For now, scan all source files for heritage clauses
	// A full implementation would use FindAllReferences with implementations flag
	for _, sourceFile := range program.GetSourceFiles() {
		scanForSubtypes(sourceFile, declaration, symbol, seen, func(decl *ast.Node) {
			if item := l.createTypeHierarchyItem(program, decl); item != nil {
				results = append(results, item)
			}
		})
	}

	// Sort results
	slices.SortFunc(results, func(a, b *lsproto.TypeHierarchyItem) int {
		if cmp := strings.Compare(string(a.Uri), string(b.Uri)); cmp != 0 {
			return cmp
		}
		return strings.Compare(a.Name, b.Name)
	})

	return results
}

// scanForSubtypes scans a source file for types that extend/implement the target type.
func scanForSubtypes(sourceFile *ast.SourceFile, target *ast.Node, targetSymbol *ast.Symbol, seen map[*ast.Node]bool, callback func(*ast.Node)) {
	sourceFile.AsNode().ForEachChild(func(node *ast.Node) bool {
		return scanNodeForSubtypes(node, target, targetSymbol, seen, callback)
	})
}

// scanNodeForSubtypes recursively scans nodes for subtype relationships.
func scanNodeForSubtypes(node *ast.Node, target *ast.Node, targetSymbol *ast.Symbol, seen map[*ast.Node]bool, callback func(*ast.Node)) bool {
	if node == nil {
		return false
	}

	// Check if this is a class or interface declaration
	if isValidTypeHierarchyDeclaration(node) && node != target {
		// Check heritage clauses
		if hasHeritageReferenceToSymbol(node, targetSymbol) {
			if !seen[node] {
				seen[node] = true
				callback(node)
			}
		}
	}

	// Continue scanning children
	node.ForEachChild(func(child *ast.Node) bool {
		return scanNodeForSubtypes(child, target, targetSymbol, seen, callback)
	})

	return false
}

// hasHeritageReferenceToSymbol checks if a node has a heritage clause referencing the target symbol.
func hasHeritageReferenceToSymbol(node *ast.Node, targetSymbol *ast.Symbol) bool {
	// Get heritage clauses
	var heritageClauses *ast.NodeList
	switch {
	case ast.IsClassDeclaration(node):
		heritageClauses = node.AsClassDeclaration().HeritageClauses
	case ast.IsClassExpression(node):
		heritageClauses = node.AsClassExpression().HeritageClauses
	case ast.IsInterfaceDeclaration(node):
		heritageClauses = node.AsInterfaceDeclaration().HeritageClauses
	}

	if heritageClauses == nil {
		return false
	}

	for _, clause := range heritageClauses.Nodes {
		if !ast.IsHeritageClause(clause) {
			continue
		}
		heritageClause := clause.AsHeritageClause()
		if heritageClause.Types == nil {
			continue
		}
		for _, typeRef := range heritageClause.Types.Nodes {
			// Get the identifier being referenced
			var expr *ast.Node
			if ast.IsExpressionWithTypeArguments(typeRef) {
				expr = typeRef.AsExpressionWithTypeArguments().Expression
			}
			if expr == nil {
				continue
			}
			// Simple name comparison (a full implementation would resolve symbols)
			if ast.IsIdentifier(expr) && targetSymbol != nil {
				targetName := getSymbolName(targetSymbol)
				if expr.Text() == targetName {
					return true
				}
			}
		}
	}

	return false
}

// getSymbolName gets the name of a symbol.
func getSymbolName(symbol *ast.Symbol) string {
	if symbol == nil {
		return ""
	}
	return symbol.Name
}

// getEffectiveBaseTypeNode gets the effective base type node from a class declaration.
func getEffectiveBaseTypeNode(node *ast.Node) *ast.Node {
	if node == nil {
		return nil
	}

	heritageClauses := getHeritageClausesWithKind(node, ast.KindExtendsKeyword)
	if len(heritageClauses) == 0 {
		return nil
	}

	types := heritageClauses[0].Types
	if types == nil || len(types.Nodes) == 0 {
		return nil
	}

	return types.Nodes[0]
}

// getHeritageClausesWithKind gets heritage clauses of a specific kind.
func getHeritageClausesWithKind(node *ast.Node, kind ast.Kind) []*ast.HeritageClause {
	if node == nil {
		return nil
	}

	var heritageClauses *ast.NodeList
	switch {
	case ast.IsClassDeclaration(node):
		heritageClauses = node.AsClassDeclaration().HeritageClauses
	case ast.IsClassExpression(node):
		heritageClauses = node.AsClassExpression().HeritageClauses
	case ast.IsInterfaceDeclaration(node):
		heritageClauses = node.AsInterfaceDeclaration().HeritageClauses
	}

	if heritageClauses == nil {
		return nil
	}

	var result []*ast.HeritageClause
	for _, clause := range heritageClauses.Nodes {
		if ast.IsHeritageClause(clause) && clause.AsHeritageClause().Token == kind {
			result = append(result, clause.AsHeritageClause())
		}
	}
	return result
}

// resolveTypeReferenceToDeclaration resolves a type reference to its declaration.
func resolveTypeReferenceToDeclaration(c *checker.Checker, typeRef *ast.Node) *ast.Node {
	if typeRef == nil {
		return nil
	}

	// Get the expression from ExpressionWithTypeArguments
	var expr *ast.Node
	if ast.IsExpressionWithTypeArguments(typeRef) {
		expr = typeRef.AsExpressionWithTypeArguments().Expression
	} else if ast.IsTypeReferenceNode(typeRef) {
		expr = typeRef.AsTypeReferenceNode().TypeName
	} else {
		expr = typeRef
	}

	if expr == nil {
		return nil
	}

	symbol := c.GetSymbolAtLocation(expr)
	if symbol == nil {
		return nil
	}

	// Resolve aliases
	if (symbol.Flags & ast.SymbolFlagsAlias) != 0 {
		symbol = c.GetAliasedSymbol(symbol)
	}

	return getTypeDeclarationFromSymbol(symbol)
}

// collectReferencedTypesFromTypeNode collects type declarations referenced in a type node.
// For type hierarchy purposes, we only want concrete types (classes, interfaces, type aliases),
// not type parameters, as those represent the structural relationship.
func collectReferencedTypesFromTypeNode(c *checker.Checker, typeNode *ast.Node, seen map[*ast.Node]bool, callback func(*ast.Node)) {
	if typeNode == nil {
		return
	}

	switch typeNode.Kind {
	case ast.KindTypeReference:
		if decl := resolveTypeReferenceToDeclaration(c, typeNode); decl != nil && !seen[decl] {
			// Filter out type parameters - they're not concrete types we want in the hierarchy
			if !ast.IsTypeParameterDeclaration(decl) {
				seen[decl] = true
				callback(decl)
			}
		}

	case ast.KindIntersectionType:
		// For intersection types, collect all member types
		for _, member := range typeNode.AsIntersectionTypeNode().Types.Nodes {
			collectReferencedTypesFromTypeNode(c, member, seen, callback)
		}

	case ast.KindUnionType:
		// For union types, collect all member types
		for _, member := range typeNode.AsUnionTypeNode().Types.Nodes {
			collectReferencedTypesFromTypeNode(c, member, seen, callback)
		}

	case ast.KindConditionalType:
		// For conditional types like `T extends Dog ? true : false`:
		// - CheckType (T) is often a type parameter, which we skip
		// - ExtendsType (Dog) is the constraint/base type we want to show
		// We only collect the ExtendsType as it represents the structural relationship
		condType := typeNode.AsConditionalTypeNode()
		// Only collect ExtendsType, not CheckType (which is often just a type parameter T)
		collectReferencedTypesFromTypeNode(c, condType.ExtendsType, seen, callback)

	case ast.KindMappedType:
		// For mapped types, collect the constraint type
		mappedType := typeNode.AsMappedTypeNode()
		if mappedType.TypeParameter != nil {
			if constraint := mappedType.TypeParameter.AsTypeParameter().Constraint; constraint != nil {
				collectReferencedTypesFromTypeNode(c, constraint, seen, callback)
			}
		}
	}
}

// collectMixinChain collects the mixin chain from a mixin variable.
func collectMixinChain(c *checker.Checker, node *ast.Node, seen map[*ast.Node]bool, callback func(*ast.Node)) {
	if node == nil || !ast.IsVariableDeclaration(node) {
		return
	}

	initializer := node.Initializer()
	if initializer == nil || !ast.IsCallExpression(initializer) {
		return
	}

	callExpr := initializer.AsCallExpression()

	// Process arguments (base classes)
	for _, arg := range callExpr.Arguments.Nodes {
		if decl := resolveExpressionToDeclaration(c, arg); decl != nil && !seen[decl] {
			seen[decl] = true
			callback(decl)
		}

		// Recursively process nested mixin calls
		if ast.IsCallExpression(arg) {
			collectMixinChainFromCall(c, arg.AsCallExpression(), seen, callback)
		}
	}
}

// collectMixinChainFromCall collects mixin chain from a nested call expression.
func collectMixinChainFromCall(c *checker.Checker, callExpr *ast.CallExpression, seen map[*ast.Node]bool, callback func(*ast.Node)) {
	for _, arg := range callExpr.Arguments.Nodes {
		if decl := resolveExpressionToDeclaration(c, arg); decl != nil && !seen[decl] {
			seen[decl] = true
			callback(decl)
		}

		if ast.IsCallExpression(arg) {
			collectMixinChainFromCall(c, arg.AsCallExpression(), seen, callback)
		}
	}
}

// resolveExpressionToDeclaration resolves an expression to its declaration.
func resolveExpressionToDeclaration(c *checker.Checker, expr *ast.Node) *ast.Node {
	if expr == nil {
		return nil
	}

	symbol := c.GetSymbolAtLocation(expr)
	if symbol == nil {
		return nil
	}

	// Skip aliases for now - a full implementation would resolve them
	if (symbol.Flags & ast.SymbolFlagsAlias) != 0 {
		return nil
	}

	return getTypeDeclarationFromSymbol(symbol)
}

// findContainingTypeDeclaration finds the containing type declaration for a node.
func findContainingTypeDeclaration(node *ast.Node) *ast.Node {
	for node != nil {
		if isValidTypeHierarchyDeclaration(node) {
			return node
		}
		node = node.Parent
	}
	return nil
}

// isSubtypeRelationship checks if a reference establishes a subtype relationship.
func isSubtypeRelationship(refNode *ast.Node, supertype *ast.Node) bool {
	if refNode == nil {
		return false
	}

	// Check if reference is in a heritage clause
	parent := refNode.Parent
	for parent != nil {
		if ast.IsHeritageClause(parent) {
			return true
		}
		if ast.IsExpressionWithTypeArguments(parent) {
			parent = parent.Parent
			continue
		}
		// Check if we're in an intersection type (which creates a subtype)
		if ast.IsIntersectionTypeNode(parent) {
			return true
		}
		break
	}

	return false
}

// collectIntersectionSubtypes finds type aliases that are intersection types including the given type.
func collectIntersectionSubtypes(c *checker.Checker, program *compiler.Program, declaration *ast.Node, seen map[*ast.Node]bool, callback func(*ast.Node)) {
	// This is a simplified implementation
	// A full implementation would scan all type aliases in the program
	// looking for intersection types that include the given type

	symbol := c.GetSymbolAtLocation(declaration)
	if symbol == nil {
		return
	}

	// For now, we rely on FindAllReferences to find these
	// This function can be expanded later for more comprehensive coverage
}
