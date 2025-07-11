package compiler

import (
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/binder"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/scanner"
)

// jsDiagnosticsVisitor is used to find TypeScript-only constructs in JavaScript files
type jsDiagnosticsVisitor struct {
	sourceFile  *ast.SourceFile
	diagnostics []*ast.Diagnostic
}

// getJSSyntacticDiagnosticsForFile returns diagnostics for TypeScript-only constructs in JavaScript files
func getJSSyntacticDiagnosticsForFile(sourceFile *ast.SourceFile) []*ast.Diagnostic {
	visitor := &jsDiagnosticsVisitor{
		sourceFile:  sourceFile,
		diagnostics: []*ast.Diagnostic{},
	}

	// Walk the entire AST to find TypeScript-only constructs
	visitor.walkNodeForJSDiagnostics(sourceFile.AsNode(), sourceFile.AsNode())

	return visitor.diagnostics
}

// getJSSyntacticDiagnosticsForFile returns diagnostics for TypeScript-only constructs in JavaScript files
func (p *Program) getJSSyntacticDiagnosticsForFile(sourceFile *ast.SourceFile) []*ast.Diagnostic {
	if cached, ok := p.jsDiagnosticCache.Load(sourceFile); ok {
		return cached
	}
	diagnostics, _ := p.jsDiagnosticCache.LoadOrStore(sourceFile, getJSSyntacticDiagnosticsForFile(sourceFile))
	return diagnostics
}

// walkNodeForJSDiagnostics walks the AST and collects diagnostics for TypeScript-only constructs
func (v *jsDiagnosticsVisitor) walkNodeForJSDiagnostics(node *ast.Node, parent *ast.Node) {
	if node == nil {
		return
	}

	// Bail out early if this node has NodeFlagsReparsed, as they are synthesized type annotations
	if node.Flags&ast.NodeFlagsReparsed != 0 {
		return
	}

	// Handle specific parent-child relationships first
	switch parent.Kind {
	case ast.KindParameter, ast.KindPropertyDeclaration, ast.KindMethodDeclaration:
		// Check for question token (optional markers)
		if parent.Kind == ast.KindParameter && parent.AsParameterDeclaration() != nil && parent.AsParameterDeclaration().QuestionToken == node {
			v.diagnostics = append(v.diagnostics, v.createDiagnosticForNode(node, diagnostics.The_0_modifier_can_only_be_used_in_TypeScript_files, scanner.TokenToString(node.Kind)))
			return
		}
		if parent.Kind == ast.KindPropertyDeclaration && parent.AsPropertyDeclaration() != nil && parent.AsPropertyDeclaration().PostfixToken == node {
			v.diagnostics = append(v.diagnostics, v.createDiagnosticForNode(node, diagnostics.The_0_modifier_can_only_be_used_in_TypeScript_files, scanner.TokenToString(node.Kind)))
			return
		}
		if parent.Kind == ast.KindMethodDeclaration && parent.AsMethodDeclaration() != nil && parent.AsMethodDeclaration().PostfixToken == node {
			if node.Kind == ast.KindExclamationToken && parent.Parent.Kind == ast.KindObjectLiteralExpression {
				// This already gets a grammar error.
				return
			}
			v.diagnostics = append(v.diagnostics, v.createDiagnosticForNode(node, diagnostics.The_0_modifier_can_only_be_used_in_TypeScript_files, scanner.TokenToString(node.Kind)))
			return
		}
		fallthrough
	case ast.KindConstructor, ast.KindGetAccessor, ast.KindSetAccessor,
		ast.KindFunctionExpression, ast.KindFunctionDeclaration, ast.KindArrowFunction, ast.KindVariableDeclaration:
		// Check for type annotations
		if v.isTypeAnnotation(parent, node) {
			v.diagnostics = append(v.diagnostics, v.createDiagnosticForNode(node, diagnostics.Type_annotations_can_only_be_used_in_TypeScript_files))
			return
		}
	}

	// Check node-specific constructs
	switch node.Kind {
	case ast.KindImportClause:
		if node.AsImportClause() != nil && node.AsImportClause().IsTypeOnly {
			v.diagnostics = append(v.diagnostics, v.createDiagnosticForNode(parent, diagnostics.X_0_declarations_can_only_be_used_in_TypeScript_files, "import type"))
			return
		}

	case ast.KindExportDeclaration:
		if node.AsExportDeclaration() != nil && node.AsExportDeclaration().IsTypeOnly {
			v.diagnostics = append(v.diagnostics, v.createDiagnosticForNode(node, diagnostics.X_0_declarations_can_only_be_used_in_TypeScript_files, "export type"))
			return
		}

	case ast.KindImportSpecifier:
		if node.AsImportSpecifier() != nil && node.AsImportSpecifier().IsTypeOnly {
			v.diagnostics = append(v.diagnostics, v.createDiagnosticForNode(node, diagnostics.X_0_declarations_can_only_be_used_in_TypeScript_files, "import...type"))
			return
		}

	case ast.KindExportSpecifier:
		if node.AsExportSpecifier() != nil && node.AsExportSpecifier().IsTypeOnly {
			v.diagnostics = append(v.diagnostics, v.createDiagnosticForNode(node, diagnostics.X_0_declarations_can_only_be_used_in_TypeScript_files, "export...type"))
			return
		}

	case ast.KindImportEqualsDeclaration:
		v.diagnostics = append(v.diagnostics, v.createDiagnosticForNode(node, diagnostics.X_import_can_only_be_used_in_TypeScript_files))
		return

	case ast.KindExportAssignment:
		if node.AsExportAssignment() != nil && node.AsExportAssignment().IsExportEquals {
			v.diagnostics = append(v.diagnostics, v.createDiagnosticForNode(node, diagnostics.X_export_can_only_be_used_in_TypeScript_files))
			return
		}

	case ast.KindHeritageClause:
		if node.AsHeritageClause() != nil && node.AsHeritageClause().Token == ast.KindImplementsKeyword {
			v.diagnostics = append(v.diagnostics, v.createDiagnosticForNode(node, diagnostics.X_implements_clauses_can_only_be_used_in_TypeScript_files))
			return
		}

	case ast.KindInterfaceDeclaration:
		v.diagnostics = append(v.diagnostics, v.createDiagnosticForNode(node, diagnostics.X_0_declarations_can_only_be_used_in_TypeScript_files, "interface"))
		return

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
		return

	case ast.KindTypeAliasDeclaration:
		v.diagnostics = append(v.diagnostics, v.createDiagnosticForNode(node, diagnostics.Type_aliases_can_only_be_used_in_TypeScript_files))
		return

	case ast.KindEnumDeclaration:
		v.diagnostics = append(v.diagnostics, v.createDiagnosticForNode(node, diagnostics.X_0_declarations_can_only_be_used_in_TypeScript_files, "enum"))
		return

	case ast.KindNonNullExpression:
		v.diagnostics = append(v.diagnostics, v.createDiagnosticForNode(node, diagnostics.Non_null_assertions_can_only_be_used_in_TypeScript_files))
		return

	case ast.KindAsExpression:
		if node.AsAsExpression() != nil {
			v.diagnostics = append(v.diagnostics, v.createDiagnosticForNode(node.AsAsExpression().Type, diagnostics.Type_assertion_expressions_can_only_be_used_in_TypeScript_files))
			return
		}

	case ast.KindSatisfiesExpression:
		if node.AsSatisfiesExpression() != nil {
			v.diagnostics = append(v.diagnostics, v.createDiagnosticForNode(node.AsSatisfiesExpression().Type, diagnostics.Type_satisfaction_expressions_can_only_be_used_in_TypeScript_files))
			return
		}

	case ast.KindConstructor, ast.KindMethodDeclaration, ast.KindFunctionDeclaration:
		// Check for signature declarations (functions without bodies)
		if v.isSignatureDeclaration(node) {
			v.diagnostics = append(v.diagnostics, v.createDiagnosticForNode(node, diagnostics.Signature_declarations_can_only_be_used_in_TypeScript_files))
			return
		}
	}

	// Check for type parameters, type arguments, and modifiers
	v.checkTypeParametersAndModifiers(node)

	// Recursively walk children
	node.ForEachChild(func(child *ast.Node) bool {
		v.walkNodeForJSDiagnostics(child, node)
		return false
	})
}

// isTypeAnnotation checks if a node is a type annotation in relation to its parent
func (v *jsDiagnosticsVisitor) isTypeAnnotation(parent *ast.Node, node *ast.Node) bool {
	switch parent.Kind {
	case ast.KindFunctionDeclaration:
		return parent.AsFunctionDeclaration() != nil && parent.AsFunctionDeclaration().Type == node
	case ast.KindFunctionExpression:
		return parent.AsFunctionExpression() != nil && parent.AsFunctionExpression().Type == node
	case ast.KindArrowFunction:
		return parent.AsArrowFunction() != nil && parent.AsArrowFunction().Type == node
	case ast.KindMethodDeclaration:
		return parent.AsMethodDeclaration() != nil && parent.AsMethodDeclaration().Type == node
	case ast.KindGetAccessor:
		return parent.AsGetAccessorDeclaration() != nil && parent.AsGetAccessorDeclaration().Type == node
	case ast.KindSetAccessor:
		return parent.AsSetAccessorDeclaration() != nil && parent.AsSetAccessorDeclaration().Type == node
	case ast.KindConstructor:
		return parent.AsConstructorDeclaration() != nil && parent.AsConstructorDeclaration().Type == node
	case ast.KindVariableDeclaration:
		return parent.AsVariableDeclaration() != nil && parent.AsVariableDeclaration().Type == node
	case ast.KindParameter:
		return parent.AsParameterDeclaration() != nil && parent.AsParameterDeclaration().Type == node
	case ast.KindPropertyDeclaration:
		return parent.AsPropertyDeclaration() != nil && parent.AsPropertyDeclaration().Type == node
	}
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
	// Bail out early if this node has NodeFlagsReparsed
	if node.Flags&ast.NodeFlagsReparsed != 0 {
		return
	}

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
	switch node.Kind {
	case ast.KindClassDeclaration:
		if node.AsClassDeclaration() != nil {
			typeParameters = node.AsClassDeclaration().TypeParameters
		}
	case ast.KindClassExpression:
		if node.AsClassExpression() != nil {
			typeParameters = node.AsClassExpression().TypeParameters
		}
	case ast.KindMethodDeclaration:
		if node.AsMethodDeclaration() != nil {
			typeParameters = node.AsMethodDeclaration().TypeParameters
		}
	case ast.KindConstructor:
		if node.AsConstructorDeclaration() != nil {
			typeParameters = node.AsConstructorDeclaration().TypeParameters
		}
	case ast.KindGetAccessor:
		if node.AsGetAccessorDeclaration() != nil {
			typeParameters = node.AsGetAccessorDeclaration().TypeParameters
		}
	case ast.KindSetAccessor:
		if node.AsSetAccessorDeclaration() != nil {
			typeParameters = node.AsSetAccessorDeclaration().TypeParameters
		}
	case ast.KindFunctionExpression:
		if node.AsFunctionExpression() != nil {
			typeParameters = node.AsFunctionExpression().TypeParameters
		}
	case ast.KindFunctionDeclaration:
		if node.AsFunctionDeclaration() != nil {
			typeParameters = node.AsFunctionDeclaration().TypeParameters
		}
	case ast.KindArrowFunction:
		if node.AsArrowFunction() != nil {
			typeParameters = node.AsArrowFunction().TypeParameters
		}
	default:
		return nil
	}

	if typeParameters == nil {
		return nil
	}

	// Check if all type parameters are reparsed (JSDoc originated)
	for _, tp := range typeParameters.Nodes {
		if tp.Flags&ast.NodeFlagsReparsed == 0 {
			return typeParameters // Found a non-reparsed type parameter, so return the type parameters
		}
	}

	return nil // All type parameters are reparsed (JSDoc originated), so this is valid in JS
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
	switch node.Kind {
	case ast.KindCallExpression:
		if node.AsCallExpression() != nil {
			typeArguments = node.AsCallExpression().TypeArguments
		}
	case ast.KindNewExpression:
		if node.AsNewExpression() != nil {
			typeArguments = node.AsNewExpression().TypeArguments
		}
	case ast.KindExpressionWithTypeArguments:
		if node.AsExpressionWithTypeArguments() != nil {
			typeArguments = node.AsExpressionWithTypeArguments().TypeArguments
		}
	case ast.KindJsxSelfClosingElement:
		if node.AsJsxSelfClosingElement() != nil {
			typeArguments = node.AsJsxSelfClosingElement().TypeArguments
		}
	case ast.KindJsxOpeningElement:
		if node.AsJsxOpeningElement() != nil {
			typeArguments = node.AsJsxOpeningElement().TypeArguments
		}
	case ast.KindTaggedTemplateExpression:
		if node.AsTaggedTemplateExpression() != nil {
			typeArguments = node.AsTaggedTemplateExpression().TypeArguments
		}
	}

	if typeArguments == nil {
		return nil
	}

	// Check if all type arguments are reparsed (JSDoc originated)
	for _, ta := range typeArguments.Nodes {
		if ta.Flags&ast.NodeFlagsReparsed == 0 {
			return typeArguments // Found a non-reparsed type argument, so return the type arguments
		}
	}

	return nil // All type arguments are reparsed (JSDoc originated), so this is valid in JS
}

// hasTypeArguments checks if a node has type arguments
func (v *jsDiagnosticsVisitor) hasTypeArguments(node *ast.Node) bool {
	// Bail out early if this node has NodeFlagsReparsed
	if node.Flags&ast.NodeFlagsReparsed != 0 {
		return false
	}

	var typeArguments *ast.NodeList
	switch node.Kind {
	case ast.KindCallExpression:
		if node.AsCallExpression() != nil {
			typeArguments = node.AsCallExpression().TypeArguments
		}
	case ast.KindNewExpression:
		if node.AsNewExpression() != nil {
			typeArguments = node.AsNewExpression().TypeArguments
		}
	case ast.KindExpressionWithTypeArguments:
		if node.AsExpressionWithTypeArguments() != nil {
			typeArguments = node.AsExpressionWithTypeArguments().TypeArguments
		}
	case ast.KindJsxSelfClosingElement:
		if node.AsJsxSelfClosingElement() != nil {
			typeArguments = node.AsJsxSelfClosingElement().TypeArguments
		}
	case ast.KindJsxOpeningElement:
		if node.AsJsxOpeningElement() != nil {
			typeArguments = node.AsJsxOpeningElement().TypeArguments
		}
	case ast.KindTaggedTemplateExpression:
		if node.AsTaggedTemplateExpression() != nil {
			typeArguments = node.AsTaggedTemplateExpression().TypeArguments
		}
	}

	if typeArguments == nil {
		return false
	}

	// Check if all type arguments are reparsed (JSDoc originated)
	for _, ta := range typeArguments.Nodes {
		if ta.Flags&ast.NodeFlagsReparsed == 0 {
			return true // Found a non-reparsed type argument, so this is a TypeScript-only construct
		}
	}

	return false // All type arguments are reparsed (JSDoc originated), so this is valid in JS
}

// checkModifiers checks for TypeScript-only modifiers on various declaration types
func (v *jsDiagnosticsVisitor) checkModifiers(node *ast.Node) {
	// Bail out early if this node has NodeFlagsReparsed
	if node.Flags&ast.NodeFlagsReparsed != 0 {
		return
	}

	// Check for TypeScript-only modifiers on various declaration types
	switch node.Kind {
	case ast.KindVariableStatement:
		if node.AsVariableStatement() != nil && node.AsVariableStatement().Modifiers() != nil {
			v.checkModifierList(node.AsVariableStatement().Modifiers(), true)
		}
	case ast.KindPropertyDeclaration:
		if node.AsPropertyDeclaration() != nil && node.AsPropertyDeclaration().Modifiers() != nil {
			v.checkPropertyModifiers(node.AsPropertyDeclaration().Modifiers())
		}
	case ast.KindParameter:
		if node.AsParameterDeclaration() != nil && node.AsParameterDeclaration().Modifiers() != nil {
			v.checkParameterModifiers(node.AsParameterDeclaration().Modifiers())
		}
	case ast.KindClassDeclaration:
		if node.AsClassDeclaration() != nil && node.AsClassDeclaration().Modifiers() != nil {
			v.checkModifierList(node.AsClassDeclaration().Modifiers(), false)
		}
	case ast.KindMethodDeclaration:
		if node.AsMethodDeclaration() != nil && node.AsMethodDeclaration().Modifiers() != nil {
			v.checkModifierList(node.AsMethodDeclaration().Modifiers(), false)
		}
	case ast.KindFunctionDeclaration:
		if node.AsFunctionDeclaration() != nil && node.AsFunctionDeclaration().Modifiers() != nil {
			v.checkModifierList(node.AsFunctionDeclaration().Modifiers(), false)
		}
	}
}

// checkModifierList checks a list of modifiers for TypeScript-only constructs
func (v *jsDiagnosticsVisitor) checkModifierList(modifiers *ast.ModifierList, isConstValid bool) {
	if modifiers == nil {
		return
	}

	for _, modifier := range modifiers.Nodes {
		// Bail out early if this modifier has NodeFlagsReparsed
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
		// Bail out early if this modifier has NodeFlagsReparsed
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
		// Bail out early if this modifier has NodeFlagsReparsed
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
	return ast.NewDiagnostic(v.sourceFile, binder.GetErrorRangeForNode(v.sourceFile, node), message, args...)
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
