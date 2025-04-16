package declarations

import (
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/printer"
)

func isPreservedDeclarationStatement(node *ast.Node) bool {
	switch node.Kind {
	case ast.KindFunctionDeclaration,
		ast.KindModuleDeclaration,
		ast.KindImportEqualsDeclaration,
		ast.KindInterfaceDeclaration,
		ast.KindClassDeclaration,
		ast.KindTypeAliasDeclaration,
		ast.KindEnumDeclaration,
		ast.KindVariableStatement,
		ast.KindImportDeclaration,
		ast.KindExportDeclaration,
		ast.KindExportAssignment:
		return true
	}
	return false
}

func needsScopeMarker(result *ast.Node) bool {
	return !ast.IsAnyImportOrReExport(result) && !ast.IsExportAssignment(result) && !ast.HasSyntacticModifier(result, ast.ModifierFlagsExport) && !ast.IsAmbientModule(result)
}

func isLateVisibilityPaintedStatement(node *ast.Node) bool {
	switch node.Kind {
	case ast.KindImportDeclaration,
		ast.KindImportEqualsDeclaration,
		ast.KindVariableStatement,
		ast.KindClassDeclaration,
		ast.KindFunctionDeclaration,
		ast.KindModuleDeclaration,
		ast.KindTypeAliasDeclaration,
		ast.KindInterfaceDeclaration,
		ast.KindEnumDeclaration:
		return true
	default:
		return false
	}
}

func canHaveLiteralInitializer(host DeclarationEmitHost, node *ast.Node) bool {
	switch node.Kind {
	case ast.KindPropertyDeclaration,
		ast.KindPropertySignature:
		return host.GetEffectiveDeclarationFlags(node, ast.ModifierFlagsPrivate) != 0
	case ast.KindParameter,
		ast.KindVariableDeclaration:
		return true
	}
	return false
}

func canProduceDiagnostics(node *ast.Node) bool {
	return ast.IsVariableDeclaration(node) ||
		ast.IsPropertyDeclaration(node) ||
		ast.IsPropertySignatureDeclaration(node) ||
		ast.IsBindingElement(node) ||
		ast.IsSetAccessorDeclaration(node) ||
		ast.IsGetAccessorDeclaration(node) ||
		ast.IsConstructSignatureDeclaration(node) ||
		ast.IsCallSignatureDeclaration(node) ||
		ast.IsMethodDeclaration(node) ||
		ast.IsMethodSignatureDeclaration(node) ||
		ast.IsFunctionDeclaration(node) ||
		ast.IsParameter(node) ||
		ast.IsTypeParameterDeclaration(node) ||
		ast.IsExpressionWithTypeArguments(node) ||
		ast.IsImportEqualsDeclaration(node) ||
		ast.IsTypeAliasDeclaration(node) ||
		ast.IsConstructorDeclaration(node) ||
		ast.IsIndexSignatureDeclaration(node) ||
		ast.IsPropertyAccessExpression(node) ||
		ast.IsElementAccessExpression(node) ||
		ast.IsBinaryExpression(node) // || // !!! TODO: JSDoc support
	/* ast.IsJSDocTypeAlias(node); */
}

func hasInferredType(node *ast.Node) bool {
	// Debug.type<HasInferredType>(node); // !!!
	switch node.Kind {
	case ast.KindParameter,
		ast.KindPropertySignature,
		ast.KindPropertyDeclaration,
		ast.KindBindingElement,
		ast.KindPropertyAccessExpression,
		ast.KindElementAccessExpression,
		ast.KindBinaryExpression,
		ast.KindVariableDeclaration,
		ast.KindExportAssignment,
		ast.KindPropertyAssignment,
		ast.KindShorthandPropertyAssignment,
		ast.KindJSDocParameterTag,
		ast.KindJSDocPropertyTag:
		return true
	default:
		// assertType<never>(node); // !!!
		return false
	}
}

func isDeclarationAndNotVisible(emitContext *printer.EmitContext, resolver printer.EmitResolver, node *ast.Node) bool {
	node = emitContext.ParseNode(node)
	switch node.Kind {
	case ast.KindFunctionDeclaration,
		ast.KindModuleDeclaration,
		ast.KindInterfaceDeclaration,
		ast.KindClassDeclaration,
		ast.KindTypeAliasDeclaration,
		ast.KindEnumDeclaration:
		return !resolver.IsDeclarationVisible(node)
	// The following should be doing their own visibility checks based on filtering their members
	case ast.KindVariableDeclaration:
		return !getBindingNameVisible(resolver, node)
	case ast.KindImportEqualsDeclaration:
	case ast.KindImportDeclaration:
	case ast.KindExportDeclaration:
	case ast.KindExportAssignment:
		return false
	case ast.KindClassStaticBlockDeclaration:
		return true
	}
	return false
}

func getBindingNameVisible(resolver printer.EmitResolver, elem *ast.Node) bool {
	if ast.IsOmittedExpression(elem) {
		return false
	}
	if ast.IsBindingPattern(elem.Name()) {
		// If any child binding pattern element has been marked visible (usually by collect linked aliases), then this is visible
		for _, elem := range elem.Name().AsBindingPattern().Elements.Nodes {
			if getBindingNameVisible(resolver, elem) {
				return true
			}
		}
		return false
	} else {
		return resolver.IsDeclarationVisible(elem)
	}
}

func isEnclosingDeclaration(node *ast.Node) bool {
	return ast.IsSourceFile(node) ||
		ast.IsTypeAliasDeclaration(node) ||
		ast.IsModuleDeclaration(node) ||
		ast.IsClassDeclaration(node) ||
		ast.IsInterfaceDeclaration(node) ||
		ast.IsFunctionLike(node) ||
		ast.IsIndexSignatureDeclaration(node) ||
		ast.IsMappedTypeNode(node)
}

func isAlwaysType(node *ast.Node) bool {
	if node.Kind == ast.KindInterfaceDeclaration {
		return true
	}
	return false
}

func maskModifierFlags(host DeclarationEmitHost, node *ast.Node) ast.ModifierFlags {
	return maskModifierFlagsEx(host, node, ast.ModifierFlagsAll^ast.ModifierFlagsPublic, ast.ModifierFlagsNone)
}

func maskModifierFlagsEx(host DeclarationEmitHost, node *ast.Node, modifierMask ast.ModifierFlags, modifierAdditions ast.ModifierFlags) ast.ModifierFlags {
	flags := host.GetEffectiveDeclarationFlags(node, modifierMask) | modifierAdditions
	if flags&ast.ModifierFlagsDefault != 0 && (flags&ast.ModifierFlagsExport == 0) {
		// A non-exported default is a nonsequitor - we usually try to remove all export modifiers
		// from statements in ambient declarations; but a default export must retain its export modifier to be syntactically valid
		flags ^= ast.ModifierFlagsExport
	}
	if flags&ast.ModifierFlagsDefault != 0 && flags&ast.ModifierFlagsAmbient != 0 {
		flags ^= ast.ModifierFlagsAmbient // `declare` is never required alongside `default` (and would be an error if printed)
	}
	return flags
}
