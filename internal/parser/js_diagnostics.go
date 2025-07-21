package parser

import (
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/scanner"
)

// jsDiagnosticsVisitor is used to find TypeScript-only constructs in JavaScript files
type jsDiagnosticsVisitor struct {
	sourceFile               *ast.SourceFile
	diagnostics              []*ast.Diagnostic
	walkNodeForJSDiagnostics func(node *ast.Node) bool
}

// getJSSyntacticDiagnosticsForFile returns diagnostics for TypeScript-only constructs in JavaScript files
func getJSSyntacticDiagnosticsForFile(sourceFile *ast.SourceFile) []*ast.Diagnostic {
	visitor := &jsDiagnosticsVisitor{
		sourceFile: sourceFile,
	}
	visitor.walkNodeForJSDiagnostics = visitor.walkNodeForJSDiagnosticsWorker

	// Walk the entire AST to find TypeScript-only constructs
	sourceFile.ForEachChild(visitor.walkNodeForJSDiagnostics)

	return visitor.diagnostics
}

// walkNodeForJSDiagnostics walks the AST and collects diagnostics for TypeScript-only constructs
func (v *jsDiagnosticsVisitor) walkNodeForJSDiagnosticsWorker(node *ast.Node) bool {
	if node == nil {
		return false
	}

	parent := node.Parent

	// Bail out early if this node has NodeFlagsReparsed, as they are synthesized type annotations
	if node.Flags&ast.NodeFlagsReparsed != 0 {
		return false
	}

	// Check for question tokens (optional markers) - these are always illegal in JS
	if node.Kind == ast.KindQuestionToken {
		// Skip if this is in an object literal method which already gets a grammar error
		if parent.Kind == ast.KindMethodDeclaration && parent.Parent.Kind == ast.KindObjectLiteralExpression {
			return false
		}
		v.diagnostics = append(v.diagnostics, v.createDiagnosticForNode(node, diagnostics.The_0_modifier_can_only_be_used_in_TypeScript_files, scanner.TokenToString(node.Kind)))
		return false
	}

	// Check for exclamation tokens (definite assignment assertions) - these are always illegal in JS
	if node.Kind == ast.KindExclamationToken {
		// Skip if this is in an object literal method which already gets a grammar error
		if parent.Kind == ast.KindMethodDeclaration && parent.Parent.Kind == ast.KindObjectLiteralExpression {
			return false
		}
		v.diagnostics = append(v.diagnostics, v.createDiagnosticForNode(node, diagnostics.The_0_modifier_can_only_be_used_in_TypeScript_files, scanner.TokenToString(node.Kind)))
		return false
	}

	// Check for type nodes - these are always illegal in JS
	if ast.IsTypeNode(node) {
		v.diagnostics = append(v.diagnostics, v.createDiagnosticForNode(node, diagnostics.Type_annotations_can_only_be_used_in_TypeScript_files))
		return false
	}

	// Check node-specific constructs
	switch node.Kind {
	case ast.KindImportClause:
		if node.AsImportClause() != nil && node.AsImportClause().IsTypeOnly {
			v.diagnostics = append(v.diagnostics, v.createDiagnosticForNode(parent, diagnostics.X_0_declarations_can_only_be_used_in_TypeScript_files, "import type"))
			return false
		}

	case ast.KindExportDeclaration:
		if node.AsExportDeclaration() != nil && node.AsExportDeclaration().IsTypeOnly {
			v.diagnostics = append(v.diagnostics, v.createDiagnosticForNode(node, diagnostics.X_0_declarations_can_only_be_used_in_TypeScript_files, "export type"))
			return false
		}

	case ast.KindImportSpecifier:
		if node.AsImportSpecifier() != nil && node.AsImportSpecifier().IsTypeOnly {
			v.diagnostics = append(v.diagnostics, v.createDiagnosticForNode(node, diagnostics.X_0_declarations_can_only_be_used_in_TypeScript_files, "import...type"))
			return false
		}

	case ast.KindExportSpecifier:
		if node.AsExportSpecifier() != nil && node.AsExportSpecifier().IsTypeOnly {
			v.diagnostics = append(v.diagnostics, v.createDiagnosticForNode(node, diagnostics.X_0_declarations_can_only_be_used_in_TypeScript_files, "export...type"))
			return false
		}

	case ast.KindImportEqualsDeclaration:
		v.diagnostics = append(v.diagnostics, v.createDiagnosticForNode(node, diagnostics.X_import_can_only_be_used_in_TypeScript_files))
		return false

	case ast.KindExportAssignment:
		if node.AsExportAssignment() != nil && node.AsExportAssignment().IsExportEquals {
			v.diagnostics = append(v.diagnostics, v.createDiagnosticForNode(node, diagnostics.X_export_can_only_be_used_in_TypeScript_files))
			return false
		}

	case ast.KindHeritageClause:
		if node.AsHeritageClause() != nil && node.AsHeritageClause().Token == ast.KindImplementsKeyword {
			v.diagnostics = append(v.diagnostics, v.createDiagnosticForNode(node, diagnostics.X_implements_clauses_can_only_be_used_in_TypeScript_files))
			return false
		}

	case ast.KindInterfaceDeclaration:
		v.diagnostics = append(v.diagnostics, v.createDiagnosticForNode(node, diagnostics.X_0_declarations_can_only_be_used_in_TypeScript_files, "interface"))
		return false

	case ast.KindModuleDeclaration:
		moduleKeyword := "module"
		if node.AsModuleDeclaration() != nil {
			switch node.AsModuleDeclaration().Keyword {
			case ast.KindNamespaceKeyword:
				moduleKeyword = "namespace"
			case ast.KindModuleKeyword:
				moduleKeyword = "module"
			case ast.KindGlobalKeyword:
				moduleKeyword = "global"
			}
		}
		v.diagnostics = append(v.diagnostics, v.createDiagnosticForNode(node, diagnostics.X_0_declarations_can_only_be_used_in_TypeScript_files, moduleKeyword))
		return false

	case ast.KindTypeAliasDeclaration:
		v.diagnostics = append(v.diagnostics, v.createDiagnosticForNode(node, diagnostics.Type_aliases_can_only_be_used_in_TypeScript_files))
		return false

	case ast.KindEnumDeclaration:
		v.diagnostics = append(v.diagnostics, v.createDiagnosticForNode(node, diagnostics.X_0_declarations_can_only_be_used_in_TypeScript_files, "enum"))
		return false

	case ast.KindNonNullExpression:
		v.diagnostics = append(v.diagnostics, v.createDiagnosticForNode(node, diagnostics.Non_null_assertions_can_only_be_used_in_TypeScript_files))
		return false

	case ast.KindAsExpression:
		if node.AsAsExpression() != nil {
			v.diagnostics = append(v.diagnostics, v.createDiagnosticForNode(node.AsAsExpression().Type, diagnostics.Type_assertion_expressions_can_only_be_used_in_TypeScript_files))
			return false
		}

	case ast.KindSatisfiesExpression:
		if node.AsSatisfiesExpression() != nil {
			v.diagnostics = append(v.diagnostics, v.createDiagnosticForNode(node.AsSatisfiesExpression().Type, diagnostics.Type_satisfaction_expressions_can_only_be_used_in_TypeScript_files))
			return false
		}

	case ast.KindConstructor, ast.KindMethodDeclaration, ast.KindFunctionDeclaration:
		// Check for signature declarations (functions without bodies)
		if v.isSignatureDeclaration(node) {
			v.diagnostics = append(v.diagnostics, v.createDiagnosticForNode(node, diagnostics.Signature_declarations_can_only_be_used_in_TypeScript_files))
			return false
		}
	}

	// Check for type parameters, type arguments, and modifiers
	v.checkTypeParametersAndModifiers(node)

	// Recursively walk children
	node.ForEachChild(v.walkNodeForJSDiagnostics)

	return false
}


// isSignatureDeclaration checks if a node is a signature declaration (function without body)
func (v *jsDiagnosticsVisitor) isSignatureDeclaration(node *ast.Node) bool {
	switch node.Kind {
	case ast.KindFunctionDeclaration:
		return node.AsFunctionDeclaration() != nil && node.AsFunctionDeclaration().Body == nil
	case ast.KindMethodDeclaration:
		return node.AsMethodDeclaration() != nil && node.AsMethodDeclaration().Body == nil
	case ast.KindConstructor:
		return node.AsConstructorDeclaration() != nil && node.AsConstructorDeclaration().Body == nil
	}
	return false
}

// checkTypeParametersAndModifiers checks for type parameters, type arguments, and modifiers
func (v *jsDiagnosticsVisitor) checkTypeParametersAndModifiers(node *ast.Node) {
	// Check type parameters
	if typeParams := v.getTypeParameters(node); typeParams != nil {
		v.diagnostics = append(v.diagnostics, v.createDiagnosticForNodeList(typeParams, diagnostics.Type_parameter_declarations_can_only_be_used_in_TypeScript_files))
	}

	// Check type arguments
	if typeArgs := v.getTypeArguments(node); typeArgs != nil {
		v.diagnostics = append(v.diagnostics, v.createDiagnosticForNodeList(typeArgs, diagnostics.Type_arguments_can_only_be_used_in_TypeScript_files))
	}

	// Check modifiers
	v.checkModifiers(node)
}

// getTypeParameters returns the type parameters for a node if it has any non-reparsed ones
func (v *jsDiagnosticsVisitor) getTypeParameters(node *ast.Node) *ast.NodeList {
	// Bail out early if this node has NodeFlagsReparsed
	if node.Flags&ast.NodeFlagsReparsed != 0 {
		return nil
	}

	var typeParameters *ast.NodeList
	
	// Try class-like nodes first
	if classLike := node.ClassLikeData(); classLike != nil {
		typeParameters = classLike.TypeParameters
	} else if funcLike := node.FunctionLikeData(); funcLike != nil {
		// Try function-like nodes
		typeParameters = funcLike.TypeParameters
	}

	if typeParameters == nil {
		return nil
	}

	// Check if all type parameters are reparsed (JSDoc originated)
	// Only check the first one since the reparser only sets type parameters if there are none already
	if len(typeParameters.Nodes) > 0 && typeParameters.Nodes[0].Flags&ast.NodeFlagsReparsed != 0 {
		return nil // All type parameters are reparsed (JSDoc originated), so this is valid in JS
	}

	return typeParameters // Found non-reparsed type parameters
}

// hasTypeParameters checks if a node has type parameters
func (v *jsDiagnosticsVisitor) hasTypeParameters(node *ast.Node) bool {
	return v.getTypeParameters(node) != nil
}

// getTypeArguments returns the type arguments for a node if it has any
func (v *jsDiagnosticsVisitor) getTypeArguments(node *ast.Node) *ast.NodeList {
	// Bail out early if this node has NodeFlagsReparsed
	if node.Flags&ast.NodeFlagsReparsed != 0 {
		return nil
	}

	var typeArguments *ast.NodeList
	
	// Use the built-in TypeArgumentList() method for supported nodes
	if ast.IsCallLikeExpression(node) || node.Kind == ast.KindExpressionWithTypeArguments {
		typeArguments = node.TypeArgumentList()
	}

	if typeArguments == nil {
		return nil
	}

	// Check if all type arguments are reparsed (JSDoc originated)
	// Only check the first one since the reparser only sets type arguments if there are none already
	if len(typeArguments.Nodes) > 0 && typeArguments.Nodes[0].Flags&ast.NodeFlagsReparsed != 0 {
		return nil // All type arguments are reparsed (JSDoc originated), so this is valid in JS
	}

	return typeArguments // Found non-reparsed type arguments
}

// checkModifiers checks for TypeScript-only modifiers on various declaration types
func (v *jsDiagnosticsVisitor) checkModifiers(node *ast.Node) {
	modifiers := node.Modifiers()
	if modifiers == nil {
		return
	}

	// Handle different types of nodes with different modifier rules
	switch node.Kind {
	case ast.KindVariableStatement:
		v.checkModifierList(modifiers, true) // const is valid for variable statements
	case ast.KindPropertyDeclaration:
		v.checkPropertyModifiers(modifiers)
	case ast.KindParameter:
		v.checkParameterModifiers(modifiers)
	default:
		v.checkModifierList(modifiers, false) // const is not valid for other declarations
	}
}

// checkModifierList checks a list of modifiers for TypeScript-only constructs
func (v *jsDiagnosticsVisitor) checkModifierList(modifiers *ast.ModifierList, isConstValid bool) {
	if modifiers == nil {
		return
	}

	for _, modifier := range modifiers.Nodes {
		// Skip reparsed modifiers (from JSDoc) but continue checking others
		if modifier.Flags&ast.NodeFlagsReparsed != 0 {
			continue
		}
		v.checkModifier(modifier, isConstValid)
	}
}

// checkPropertyModifiers checks property modifiers for TypeScript-only constructs
func (v *jsDiagnosticsVisitor) checkPropertyModifiers(modifiers *ast.ModifierList) {
	if modifiers == nil {
		return
	}

	for _, modifier := range modifiers.Nodes {
		// Skip reparsed modifiers (from JSDoc) but continue checking others
		if modifier.Flags&ast.NodeFlagsReparsed != 0 {
			continue
		}
		// Property modifiers allow static and accessor, but all other modifiers are invalid
		switch modifier.Kind {
		case ast.KindStaticKeyword, ast.KindAccessorKeyword:
			// These are valid in JavaScript
			continue
		default:
			// All other modifiers are invalid on properties in JavaScript
			v.diagnostics = append(v.diagnostics, v.createDiagnosticForNode(modifier, diagnostics.The_0_modifier_can_only_be_used_in_TypeScript_files, scanner.TokenToString(modifier.Kind)))
		}
	}
}

// checkParameterModifiers checks parameter modifiers for TypeScript-only constructs
func (v *jsDiagnosticsVisitor) checkParameterModifiers(modifiers *ast.ModifierList) {
	if modifiers == nil {
		return
	}

	for _, modifier := range modifiers.Nodes {
		// Skip reparsed modifiers (from JSDoc) but continue checking others
		if modifier.Flags&ast.NodeFlagsReparsed != 0 {
			continue
		}
		// All parameter modifiers are invalid in JavaScript
		v.diagnostics = append(v.diagnostics, v.createDiagnosticForNode(modifier, diagnostics.Parameter_modifiers_can_only_be_used_in_TypeScript_files))
	}
}

// checkModifier checks a single modifier for TypeScript-only constructs
func (v *jsDiagnosticsVisitor) checkModifier(modifier *ast.Node, isConstValid bool) {
	switch modifier.Kind {
	case ast.KindConstKeyword:
		if !isConstValid {
			v.diagnostics = append(v.diagnostics, v.createDiagnosticForNode(modifier, diagnostics.The_0_modifier_can_only_be_used_in_TypeScript_files, "const"))
		}
	case ast.KindPublicKeyword, ast.KindPrivateKeyword, ast.KindProtectedKeyword, ast.KindReadonlyKeyword,
		ast.KindDeclareKeyword, ast.KindAbstractKeyword, ast.KindOverrideKeyword, ast.KindInKeyword, ast.KindOutKeyword:
		v.diagnostics = append(v.diagnostics, v.createDiagnosticForNode(modifier, diagnostics.The_0_modifier_can_only_be_used_in_TypeScript_files, scanner.TokenToString(modifier.Kind)))
	case ast.KindStaticKeyword, ast.KindExportKeyword, ast.KindDefaultKeyword, ast.KindAccessorKeyword:
		// These are valid in JavaScript
	}
}

// createDiagnosticForNode creates a diagnostic for a specific node
func (v *jsDiagnosticsVisitor) createDiagnosticForNode(node *ast.Node, message *diagnostics.Message, args ...any) *ast.Diagnostic {
	return ast.NewDiagnostic(v.sourceFile, scanner.GetErrorRangeForNode(v.sourceFile, node), message, args...)
}

// createDiagnosticForNodeList creates a diagnostic for a NodeList
func (v *jsDiagnosticsVisitor) createDiagnosticForNodeList(nodeList *ast.NodeList, message *diagnostics.Message, args ...any) *ast.Diagnostic {
	// If the NodeList's location is not set correctly, fall back to using the first node
	if nodeList.Loc.Pos() == nodeList.Loc.End() && len(nodeList.Nodes) > 0 {
		// Calculate the range from the first to last node
		start := nodeList.Nodes[0].Pos()
		end := nodeList.Nodes[len(nodeList.Nodes)-1].End()
		return ast.NewDiagnostic(v.sourceFile, core.NewTextRange(start, end), message, args...)
	}
	return ast.NewDiagnostic(v.sourceFile, nodeList.Loc, message, args...)
}
