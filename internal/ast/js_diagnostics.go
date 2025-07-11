package ast

import (
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/diagnostics"
)

// JS syntactic diagnostics functions

func getJSSyntacticDiagnosticsForFile(sourceFile *SourceFile) []*Diagnostic {
	diagnostics := []*Diagnostic{}

	// Walk the entire AST to find TypeScript-only constructs
	walkNodeForJSDiagnostics(sourceFile, sourceFile.AsNode(), sourceFile.AsNode(), &diagnostics)

	return diagnostics
}

func walkNodeForJSDiagnostics(sourceFile *SourceFile, node *Node, parent *Node, diags *[]*Diagnostic) {
	if node == nil {
		return
	}

	// Bail out early if this node has NodeFlagsReparsed, as they are synthesized type annotations
	if node.Flags&NodeFlagsReparsed != 0 {
		return
	}

	// Handle specific parent-child relationships first
	switch parent.Kind {
	case KindParameter, KindPropertyDeclaration, KindMethodDeclaration:
		// Check for question token (optional markers) - only parameters have question tokens
		if parent.Kind == KindParameter && parent.AsParameterDeclaration() != nil && parent.AsParameterDeclaration().QuestionToken == node {
			*diags = append(*diags, createDiagnosticForNode(sourceFile, node, diagnostics.The_0_modifier_can_only_be_used_in_TypeScript_files, "?"))
			return
		}
		fallthrough
	case KindMethodSignature, KindConstructor, KindGetAccessor, KindSetAccessor,
		KindFunctionExpression, KindFunctionDeclaration, KindArrowFunction, KindVariableDeclaration:
		// Check for type annotations
		if isTypeAnnotation(parent, node) {
			*diags = append(*diags, createDiagnosticForNode(sourceFile, node, diagnostics.Type_annotations_can_only_be_used_in_TypeScript_files))
			return
		}
	}

	// Check node-specific constructs
	switch node.Kind {
	case KindImportClause:
		if node.AsImportClause() != nil && node.AsImportClause().IsTypeOnly {
			*diags = append(*diags, createDiagnosticForNode(sourceFile, parent, diagnostics.X_0_declarations_can_only_be_used_in_TypeScript_files, "import type"))
			return
		}

	case KindExportDeclaration:
		if node.AsExportDeclaration() != nil && node.AsExportDeclaration().IsTypeOnly {
			*diags = append(*diags, createDiagnosticForNode(sourceFile, node, diagnostics.X_0_declarations_can_only_be_used_in_TypeScript_files, "export type"))
			return
		}

	case KindImportSpecifier:
		if node.AsImportSpecifier() != nil && node.AsImportSpecifier().IsTypeOnly {
			*diags = append(*diags, createDiagnosticForNode(sourceFile, node, diagnostics.X_0_declarations_can_only_be_used_in_TypeScript_files, "import...type"))
			return
		}

	case KindExportSpecifier:
		if node.AsExportSpecifier() != nil && node.AsExportSpecifier().IsTypeOnly {
			*diags = append(*diags, createDiagnosticForNode(sourceFile, node, diagnostics.X_0_declarations_can_only_be_used_in_TypeScript_files, "export...type"))
			return
		}

	case KindImportEqualsDeclaration:
		*diags = append(*diags, createDiagnosticForNode(sourceFile, node, diagnostics.X_import_can_only_be_used_in_TypeScript_files))
		return

	case KindExportAssignment:
		if node.AsExportAssignment() != nil && node.AsExportAssignment().IsExportEquals {
			*diags = append(*diags, createDiagnosticForNode(sourceFile, node, diagnostics.X_export_can_only_be_used_in_TypeScript_files))
			return
		}

	case KindHeritageClause:
		if node.AsHeritageClause() != nil && node.AsHeritageClause().Token == KindImplementsKeyword {
			*diags = append(*diags, createDiagnosticForNode(sourceFile, node, diagnostics.X_implements_clauses_can_only_be_used_in_TypeScript_files))
			return
		}

	case KindInterfaceDeclaration:
		*diags = append(*diags, createDiagnosticForNode(sourceFile, node, diagnostics.X_0_declarations_can_only_be_used_in_TypeScript_files, "interface"))
		return

	case KindModuleDeclaration:
		moduleKeyword := "module"
		// For now, we'll just use "module" as the default since we don't have a flag to distinguish namespace vs module
		*diags = append(*diags, createDiagnosticForNode(sourceFile, node, diagnostics.X_0_declarations_can_only_be_used_in_TypeScript_files, moduleKeyword))
		return

	case KindTypeAliasDeclaration:
		*diags = append(*diags, createDiagnosticForNode(sourceFile, node, diagnostics.Type_aliases_can_only_be_used_in_TypeScript_files))
		return

	case KindEnumDeclaration:
		*diags = append(*diags, createDiagnosticForNode(sourceFile, node, diagnostics.X_0_declarations_can_only_be_used_in_TypeScript_files, "enum"))
		return

	case KindNonNullExpression:
		*diags = append(*diags, createDiagnosticForNode(sourceFile, node, diagnostics.Non_null_assertions_can_only_be_used_in_TypeScript_files))
		return

	case KindAsExpression:
		if node.AsAsExpression() != nil {
			*diags = append(*diags, createDiagnosticForNode(sourceFile, node.AsAsExpression().Type, diagnostics.Type_assertion_expressions_can_only_be_used_in_TypeScript_files))
			return
		}

	case KindSatisfiesExpression:
		if node.AsSatisfiesExpression() != nil {
			*diags = append(*diags, createDiagnosticForNode(sourceFile, node.AsSatisfiesExpression().Type, diagnostics.Type_satisfaction_expressions_can_only_be_used_in_TypeScript_files))
			return
		}

	case KindConstructor, KindMethodDeclaration, KindFunctionDeclaration:
		// Check for signature declarations (functions without bodies)
		if isSignatureDeclaration(node) {
			*diags = append(*diags, createDiagnosticForNode(sourceFile, node, diagnostics.Signature_declarations_can_only_be_used_in_TypeScript_files))
			return
		}
	}

	// Check for type parameters, type arguments, and modifiers
	checkTypeParametersAndModifiers(sourceFile, node, diags)

	// Recursively walk children
	node.ForEachChild(func(child *Node) bool {
		walkNodeForJSDiagnostics(sourceFile, child, node, diags)
		return false
	})
}

func isTypeAnnotation(parent *Node, node *Node) bool {
	switch parent.Kind {
	case KindFunctionDeclaration:
		return parent.AsFunctionDeclaration() != nil && parent.AsFunctionDeclaration().Type == node
	case KindFunctionExpression:
		return parent.AsFunctionExpression() != nil && parent.AsFunctionExpression().Type == node
	case KindArrowFunction:
		return parent.AsArrowFunction() != nil && parent.AsArrowFunction().Type == node
	case KindMethodDeclaration:
		return parent.AsMethodDeclaration() != nil && parent.AsMethodDeclaration().Type == node
	case KindGetAccessor:
		return parent.AsGetAccessorDeclaration() != nil && parent.AsGetAccessorDeclaration().Type == node
	case KindSetAccessor:
		return parent.AsSetAccessorDeclaration() != nil && parent.AsSetAccessorDeclaration().Type == node
	case KindConstructor:
		return parent.AsConstructorDeclaration() != nil && parent.AsConstructorDeclaration().Type == node
	case KindVariableDeclaration:
		return parent.AsVariableDeclaration() != nil && parent.AsVariableDeclaration().Type == node
	case KindParameter:
		return parent.AsParameterDeclaration() != nil && parent.AsParameterDeclaration().Type == node
	case KindPropertyDeclaration:
		return parent.AsPropertyDeclaration() != nil && parent.AsPropertyDeclaration().Type == node
	}
	return false
}

func isSignatureDeclaration(node *Node) bool {
	switch node.Kind {
	case KindFunctionDeclaration:
		return node.AsFunctionDeclaration() != nil && node.AsFunctionDeclaration().Body == nil
	case KindMethodDeclaration:
		return node.AsMethodDeclaration() != nil && node.AsMethodDeclaration().Body == nil
	case KindConstructor:
		return node.AsConstructorDeclaration() != nil && node.AsConstructorDeclaration().Body == nil
	}
	return false
}

func checkTypeParametersAndModifiers(sourceFile *SourceFile, node *Node, diags *[]*Diagnostic) {
	// Check type parameters
	if hasTypeParameters(node) {
		*diags = append(*diags, createDiagnosticForNode(sourceFile, node, diagnostics.Type_parameter_declarations_can_only_be_used_in_TypeScript_files))
	}

	// Check type arguments
	if hasTypeArguments(node) {
		*diags = append(*diags, createDiagnosticForNode(sourceFile, node, diagnostics.Type_arguments_can_only_be_used_in_TypeScript_files))
	}

	// Check modifiers
	checkModifiers(sourceFile, node, diags)
}

func hasTypeParameters(node *Node) bool {
	switch node.Kind {
	case KindClassDeclaration:
		return node.AsClassDeclaration() != nil && node.AsClassDeclaration().TypeParameters != nil
	case KindClassExpression:
		return node.AsClassExpression() != nil && node.AsClassExpression().TypeParameters != nil
	case KindMethodDeclaration:
		return node.AsMethodDeclaration() != nil && node.AsMethodDeclaration().TypeParameters != nil
	case KindConstructor:
		return node.AsConstructorDeclaration() != nil && node.AsConstructorDeclaration().TypeParameters != nil
	case KindGetAccessor:
		return node.AsGetAccessorDeclaration() != nil && node.AsGetAccessorDeclaration().TypeParameters != nil
	case KindSetAccessor:
		return node.AsSetAccessorDeclaration() != nil && node.AsSetAccessorDeclaration().TypeParameters != nil
	case KindFunctionExpression:
		return node.AsFunctionExpression() != nil && node.AsFunctionExpression().TypeParameters != nil
	case KindFunctionDeclaration:
		return node.AsFunctionDeclaration() != nil && node.AsFunctionDeclaration().TypeParameters != nil
	case KindArrowFunction:
		return node.AsArrowFunction() != nil && node.AsArrowFunction().TypeParameters != nil
	}
	return false
}

func hasTypeArguments(node *Node) bool {
	switch node.Kind {
	case KindCallExpression:
		return node.AsCallExpression() != nil && node.AsCallExpression().TypeArguments != nil
	case KindNewExpression:
		return node.AsNewExpression() != nil && node.AsNewExpression().TypeArguments != nil
	case KindExpressionWithTypeArguments:
		return node.AsExpressionWithTypeArguments() != nil && node.AsExpressionWithTypeArguments().TypeArguments != nil
	case KindJsxSelfClosingElement:
		return node.AsJsxSelfClosingElement() != nil && node.AsJsxSelfClosingElement().TypeArguments != nil
	case KindJsxOpeningElement:
		return node.AsJsxOpeningElement() != nil && node.AsJsxOpeningElement().TypeArguments != nil
	case KindTaggedTemplateExpression:
		return node.AsTaggedTemplateExpression() != nil && node.AsTaggedTemplateExpression().TypeArguments != nil
	}
	return false
}

func checkModifiers(sourceFile *SourceFile, node *Node, diags *[]*Diagnostic) {
	// Check for TypeScript-only modifiers on various declaration types
	switch node.Kind {
	case KindVariableStatement:
		if node.AsVariableStatement() != nil && node.AsVariableStatement().Modifiers() != nil {
			checkModifierList(sourceFile, node.AsVariableStatement().Modifiers(), true, diags)
		}
	case KindPropertyDeclaration:
		if node.AsPropertyDeclaration() != nil && node.AsPropertyDeclaration().Modifiers() != nil {
			checkPropertyModifiers(sourceFile, node.AsPropertyDeclaration().Modifiers(), diags)
		}
	case KindParameter:
		if node.AsParameterDeclaration() != nil && node.AsParameterDeclaration().Modifiers() != nil {
			checkParameterModifiers(sourceFile, node.AsParameterDeclaration().Modifiers(), diags)
		}
	}
}

func checkModifierList(sourceFile *SourceFile, modifiers *ModifierList, isConstValid bool, diags *[]*Diagnostic) {
	if modifiers == nil {
		return
	}

	for _, modifier := range modifiers.Nodes {
		checkModifier(sourceFile, modifier, isConstValid, diags)
	}
}

func checkPropertyModifiers(sourceFile *SourceFile, modifiers *ModifierList, diags *[]*Diagnostic) {
	if modifiers == nil {
		return
	}

	for _, modifier := range modifiers.Nodes {
		// Property modifiers allow static and accessor, but not other TypeScript modifiers
		switch modifier.Kind {
		case KindStaticKeyword, KindAccessorKeyword:
			// These are valid in JavaScript
			continue
		default:
			if isTypeScriptOnlyModifier(modifier) {
				*diags = append(*diags, createDiagnosticForNode(sourceFile, modifier, diagnostics.The_0_modifier_can_only_be_used_in_TypeScript_files, getTokenText(modifier)))
			}
		}
	}
}

func checkParameterModifiers(sourceFile *SourceFile, modifiers *ModifierList, diags *[]*Diagnostic) {
	if modifiers == nil {
		return
	}

	for _, modifier := range modifiers.Nodes {
		if isTypeScriptOnlyModifier(modifier) {
			*diags = append(*diags, createDiagnosticForNode(sourceFile, modifier, diagnostics.Parameter_modifiers_can_only_be_used_in_TypeScript_files))
		}
	}
}

func checkModifier(sourceFile *SourceFile, modifier *Node, isConstValid bool, diags *[]*Diagnostic) {
	switch modifier.Kind {
	case KindConstKeyword:
		if !isConstValid {
			*diags = append(*diags, createDiagnosticForNode(sourceFile, modifier, diagnostics.The_0_modifier_can_only_be_used_in_TypeScript_files, "const"))
		}
	case KindPublicKeyword, KindPrivateKeyword, KindProtectedKeyword, KindReadonlyKeyword,
		KindDeclareKeyword, KindAbstractKeyword, KindOverrideKeyword, KindInKeyword, KindOutKeyword:
		*diags = append(*diags, createDiagnosticForNode(sourceFile, modifier, diagnostics.The_0_modifier_can_only_be_used_in_TypeScript_files, getTokenText(modifier)))
	case KindStaticKeyword, KindExportKeyword, KindDefaultKeyword, KindAccessorKeyword:
		// These are valid in JavaScript
	}
}

func isTypeScriptOnlyModifier(modifier *Node) bool {
	switch modifier.Kind {
	case KindPublicKeyword, KindPrivateKeyword, KindProtectedKeyword, KindReadonlyKeyword,
		KindDeclareKeyword, KindAbstractKeyword, KindOverrideKeyword, KindInKeyword, KindOutKeyword:
		return true
	}
	return false
}

func getTokenText(node *Node) string {
	switch node.Kind {
	case KindPublicKeyword:
		return "public"
	case KindPrivateKeyword:
		return "private"
	case KindProtectedKeyword:
		return "protected"
	case KindReadonlyKeyword:
		return "readonly"
	case KindDeclareKeyword:
		return "declare"
	case KindAbstractKeyword:
		return "abstract"
	case KindOverrideKeyword:
		return "override"
	case KindInKeyword:
		return "in"
	case KindOutKeyword:
		return "out"
	case KindConstKeyword:
		return "const"
	case KindStaticKeyword:
		return "static"
	case KindAccessorKeyword:
		return "accessor"
	default:
		return ""
	}
}

func createDiagnosticForNode(sourceFile *SourceFile, node *Node, message *diagnostics.Message, args ...any) *Diagnostic {
	// Find the source file for this node
	nodeSourceFile := GetSourceFileOfNode(node)
	if nodeSourceFile == nil {
		// Fallback to empty range if we can't find the source file
		return NewDiagnostic(nil, core.TextRange{}, message, args...)
	}
	return NewDiagnostic(nodeSourceFile, getErrorRangeForNode(node), message, args...)
}

func getErrorRangeForNode(node *Node) core.TextRange {
	if node == nil {
		return core.TextRange{}
	}
	return core.NewTextRange(node.Pos(), node.End())
}
