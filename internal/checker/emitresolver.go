package checker

import (
	"sync"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/binder"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/evaluator"
	"github.com/microsoft/typescript-go/internal/nodebuilder"
	"github.com/microsoft/typescript-go/internal/printer"
)

var _ printer.EmitResolver = &emitResolver{}

type emitResolver struct {
	checker                 *Checker
	checkerMu               sync.Mutex
	isValueAliasDeclaration func(node *ast.Node) bool
	referenceResolver       binder.ReferenceResolver
}

func (r *emitResolver) GetEnumMemberValue(node *ast.Node) evaluator.Result {
	panic("unimplemented") // !!!
}

func (r *emitResolver) IsDeclarationVisible(node *ast.Node) bool {
	panic("unimplemented") // !!!
}

func (r *emitResolver) IsEntityNameVisible(entityName *ast.Node, enclosingDeclaration *ast.Node) printer.SymbolAccessibilityResult {
	panic("unimplemented") // !!!
}

func (r *emitResolver) IsImplementationOfOverload(node *ast.SignatureDeclaration) bool {
	panic("unimplemented") // !!!
}

func (r *emitResolver) IsImportRequiredByAugmentation(decl *ast.ImportDeclaration) bool {
	panic("unimplemented") // !!!
}

// TODO: the emit resolver being respoinsible for some amount of node construction is a very leaky abstraction,
// and requires giving it access to a lot of context it's otherwise not required to have, which also further complicates the API
// and likely reduces performance. There's probably some refactoring that could be done here to simplify this.

func (r *emitResolver) CreateReturnTypeOfSignatureDeclaration(emitContext *printer.EmitContext, signatureDeclaration *ast.Node, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node {
	original := emitContext.ParseNode(signatureDeclaration)
	if original == nil {
		return emitContext.Factory.NewKeywordTypeNode(ast.KindAnyKeyword)
	}
	requestNodeBuilder := NewNodeBuilderAPI(r.checker, emitContext) // TODO: cache per-context
	return requestNodeBuilder.serializeReturnTypeForSignature(original, enclosingDeclaration, flags, internalFlags, tracker)
}

func (r *emitResolver) CreateTypeOfDeclaration(emitContext *printer.EmitContext, declaration *ast.Node, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node {
	original := emitContext.ParseNode(declaration)
	if original == nil {
		return emitContext.Factory.NewKeywordTypeNode(ast.KindAnyKeyword)
	}
	requestNodeBuilder := NewNodeBuilderAPI(r.checker, emitContext) // TODO: cache per-context
	// // Get type of the symbol if this is the valid symbol otherwise get type at location
	symbol := r.checker.getSymbolOfDeclaration(declaration)
	return requestNodeBuilder.serializeTypeForDeclaration(declaration, symbol, enclosingDeclaration, flags|nodebuilder.FlagsMultilineObjectLiterals, internalFlags, tracker)
}

func (r *emitResolver) RequiresAddingImplicitUndefined(declaration *ast.Node, symbol *ast.Symbol, enclosingDeclaration *ast.Node) bool {
	switch declaration.Kind {
	case ast.KindPropertyDeclaration, ast.KindPropertySignature, ast.KindJSDocPropertyTag:
		r.checkerMu.Lock()
		defer r.checkerMu.Unlock()
		if symbol == nil {
			symbol = r.checker.getSymbolOfDeclaration(declaration)
		}
		type_ := r.checker.getTypeOfSymbol(symbol)
		r.checker.mappedSymbolLinks.Has(symbol)
		return !!((symbol.Flags&ast.SymbolFlagsProperty != 0) && (symbol.Flags&ast.SymbolFlagsOptional != 0) && isOptionalDeclaration(declaration) && r.checker.ReverseMappedSymbolLinks.Has(symbol) && r.checker.ReverseMappedSymbolLinks.Get(symbol).mappedType != nil && containsNonMissingUndefinedType(r.checker, type_))
	case ast.KindParameter, ast.KindJSDocParameterTag:
		return r.requiresAddingImplicitUndefined(declaration, enclosingDeclaration)
	default:
		panic("Node cannot possibly require adding undefined")
	}
}

func (r *emitResolver) requiresAddingImplicitUndefined(parameter *ast.Node, enclosingDeclaration *ast.Node) bool {
	return (r.isRequiredInitializedParameter(parameter, enclosingDeclaration) || r.isOptionalUninitializedParameterProperty(parameter)) && !r.declaredParameterTypeContainsUndefined(parameter)
}

func (r *emitResolver) declaredParameterTypeContainsUndefined(parameter *ast.Node) bool {
	// typeNode := getNonlocalEffectiveTypeAnnotationNode(parameter); // !!! JSDoc Support
	typeNode := parameter.Type()
	if typeNode == nil {
		return false
	}
	r.checkerMu.Lock()
	defer r.checkerMu.Unlock()
	type_ := r.checker.getTypeFromTypeNode(typeNode)
	// allow error type here to avoid confusing errors that the annotation has to contain undefined when it does in cases like this:
	//
	// export function fn(x?: Unresolved | undefined): void {}
	return r.checker.isErrorType(type_) || r.checker.containsUndefinedType(type_)
}

func (r *emitResolver) isOptionalUninitializedParameterProperty(parameter *ast.Node) bool {
	return r.checker.strictNullChecks &&
		r.isOptionalParameter(parameter) &&
		( /*isJSDocParameterTag(parameter) ||*/ parameter.Initializer() != nil) && // !!! TODO: JSDoc support
		ast.HasSyntacticModifier(parameter, ast.ModifierFlagsParameterPropertyModifier)
}

func (r *emitResolver) isRequiredInitializedParameter(parameter *ast.Node, enclosingDeclaration *ast.Node) bool {
	if r.checker.strictNullChecks || r.isOptionalParameter(parameter) || /*isJSDocParameterTag(parameter) ||*/ parameter.Initializer() == nil { // !!! TODO: JSDoc Support
		return false
	}
	if ast.HasSyntacticModifier(parameter, ast.ModifierFlagsParameterPropertyModifier) {
		return enclosingDeclaration != nil && ast.IsFunctionLikeDeclaration(enclosingDeclaration)
	}
	return true
}

func (r *emitResolver) isOptionalParameter(node *ast.Node) bool {
	// !!! TODO: JSDoc support
	// if (hasEffectiveQuestionToken(node)) {
	// 	return true;
	// }
	if ast.IsParameter(node) && node.AsParameterDeclaration().QuestionToken != nil {
		return true
	}
	if !ast.IsParameter(node) {
		return false
	}
	if node.Initializer() != nil {
		signature := r.checker.getSignatureFromDeclaration(node.Parent)
		parameterIndex := core.FindIndex(node.Parent.Parameters(), func(p *ast.ParameterDeclarationNode) bool { return p == node })
		// Debug.assert(parameterIndex >= 0); // !!!
		// Only consider syntactic or instantiated parameters as optional, not `void` parameters as this function is used
		// in grammar checks and checking for `void` too early results in parameter types widening too early
		// and causes some noImplicitAny errors to be lost.
		return parameterIndex >= r.checker.getMinArgumentCountEx(signature, MinArgumentCountFlagsStrongArityForUntypedJS|MinArgumentCountFlagsVoidIsNonOptional)
	}
	iife := ast.GetImmediatelyInvokedFunctionExpression(node.Parent)
	if iife != nil {
		parameterIndex := core.FindIndex(node.Parent.Parameters(), func(p *ast.ParameterDeclarationNode) bool { return p == node })
		return node.Type() == nil &&
			node.AsParameterDeclaration().DotDotDotToken == nil &&
			parameterIndex >= len(r.checker.getEffectiveCallArguments(iife))
	}

	return false
}

func (r *emitResolver) IsLiteralConstDeclaration(node *ast.Node) bool {
	if isDeclarationReadonly(node) || ast.IsVariableDeclaration(node) && ast.IsVarConst(node) {
		r.checkerMu.Lock()
		defer r.checkerMu.Unlock()
		return isFreshLiteralType(r.checker.getTypeOfSymbol(r.checker.getSymbolOfDeclaration(node)))
	}
	return false
}

func (r *emitResolver) IsExpandoFunctionDeclaration(node *ast.Node) bool {
	// !!! TODO: expando function support
	return false
}

func (r *emitResolver) IsSymbolAccessible(symbol *ast.Symbol, enclosingDeclaration *ast.Node, meaning ast.SymbolFlags, shouldComputeAliasToMarkVisible bool) printer.SymbolAccessibilityResult {
	r.checkerMu.Lock()
	defer r.checkerMu.Unlock()
	return r.checker.IsSymbolAccessible(symbol, enclosingDeclaration, meaning, shouldComputeAliasToMarkVisible)
}

func isConstEnumOrConstEnumOnlyModule(s *ast.Symbol) bool {
	return isConstEnumSymbol(s) || s.Flags&ast.SymbolFlagsConstEnumOnlyModule != 0
}

func (r *emitResolver) IsReferencedAliasDeclaration(node *ast.Node) bool {
	c := r.checker
	if !c.canCollectSymbolAliasAccessibilityData || !ast.IsParseTreeNode(node) {
		return true
	}

	r.checkerMu.Lock()
	defer r.checkerMu.Unlock()

	if ast.IsAliasSymbolDeclaration(node) {
		if symbol := c.getSymbolOfDeclaration(node); symbol != nil {
			aliasLinks := c.aliasSymbolLinks.Get(symbol)
			if aliasLinks.referenced {
				return true
			}
			target := aliasLinks.aliasTarget
			if target != nil && node.ModifierFlags()&ast.ModifierFlagsExport != 0 &&
				c.getSymbolFlags(target)&ast.SymbolFlagsValue != 0 &&
				(c.compilerOptions.ShouldPreserveConstEnums() || !isConstEnumOrConstEnumOnlyModule(target)) {
				return true
			}
		}
	}
	return false
}

func (r *emitResolver) IsValueAliasDeclaration(node *ast.Node) bool {
	c := r.checker
	if !c.canCollectSymbolAliasAccessibilityData || !ast.IsParseTreeNode(node) {
		return true
	}

	r.checkerMu.Lock()
	defer r.checkerMu.Unlock()

	return r.isValueAliasDeclarationWorker(node)
}

func (r *emitResolver) isValueAliasDeclarationWorker(node *ast.Node) bool {
	c := r.checker

	switch node.Kind {
	case ast.KindImportEqualsDeclaration:
		return r.isAliasResolvedToValue(c.getSymbolOfDeclaration(node), false /*excludeTypeOnlyValues*/)
	case ast.KindImportClause,
		ast.KindNamespaceImport,
		ast.KindImportSpecifier,
		ast.KindExportSpecifier:
		symbol := c.getSymbolOfDeclaration(node)
		return symbol != nil && r.isAliasResolvedToValue(symbol, true /*excludeTypeOnlyValues*/)
	case ast.KindExportDeclaration:
		exportClause := node.AsExportDeclaration().ExportClause
		if r.isValueAliasDeclaration == nil {
			r.isValueAliasDeclaration = r.isValueAliasDeclarationWorker
		}
		return exportClause != nil && (ast.IsNamespaceExport(exportClause) ||
			core.Some(exportClause.AsNamedExports().Elements.Nodes, r.isValueAliasDeclaration))
	case ast.KindExportAssignment:
		if node.AsExportAssignment().Expression != nil && node.AsExportAssignment().Expression.Kind == ast.KindIdentifier {
			return r.isAliasResolvedToValue(c.getSymbolOfDeclaration(node), true /*excludeTypeOnlyValues*/)
		}
		return true
	}
	return false
}

func (r *emitResolver) isAliasResolvedToValue(symbol *ast.Symbol, excludeTypeOnlyValues bool) bool {
	c := r.checker
	if symbol == nil {
		return false
	}
	if symbol.ValueDeclaration != nil {
		if container := ast.GetSourceFileOfNode(symbol.ValueDeclaration); container != nil {
			fileSymbol := c.getSymbolOfDeclaration(container.AsNode())
			// Ensures cjs export assignment is setup, since this symbol may point at, and merge with, the file itself.
			// If we don't, the merge may not have yet occurred, and the flags check below will be missing flags that
			// are added as a result of the merge.
			c.resolveExternalModuleSymbol(fileSymbol, false /*dontResolveAlias*/)
		}
	}
	target := c.getExportSymbolOfValueSymbolIfExported(c.resolveAlias(symbol))
	if target == c.unknownSymbol {
		return !excludeTypeOnlyValues || c.getTypeOnlyAliasDeclaration(symbol) == nil
	}
	// const enums and modules that contain only const enums are not considered values from the emit perspective
	// unless 'preserveConstEnums' option is set to true
	return c.getSymbolFlagsEx(symbol, excludeTypeOnlyValues, true /*excludeLocalMeanings*/)&ast.SymbolFlagsValue != 0 &&
		(c.compilerOptions.ShouldPreserveConstEnums() ||
			!isConstEnumOrConstEnumOnlyModule(target))
}

func (r *emitResolver) IsTopLevelValueImportEqualsWithEntityName(node *ast.Node) bool {
	c := r.checker
	if !c.canCollectSymbolAliasAccessibilityData {
		return true
	}
	if !ast.IsParseTreeNode(node) || node.Kind != ast.KindImportEqualsDeclaration || node.Parent.Kind != ast.KindSourceFile {
		return false
	}
	n := node.AsImportEqualsDeclaration()
	if ast.NodeIsMissing(n.ModuleReference) || n.ModuleReference.Kind != ast.KindExternalModuleReference {
		return false
	}

	r.checkerMu.Lock()
	defer r.checkerMu.Unlock()

	return r.isAliasResolvedToValue(c.getSymbolOfDeclaration(node), false /*excludeTypeOnlyValues*/)
}

func (r *emitResolver) MarkLinkedReferencesRecursively(file *ast.SourceFile) {
	r.checkerMu.Lock()
	defer r.checkerMu.Unlock()

	if file != nil {
		var visit ast.Visitor
		visit = func(n *ast.Node) bool {
			if ast.IsImportEqualsDeclaration(n) && n.ModifierFlags()&ast.ModifierFlagsExport == 0 {
				return false // These are deferred and marked in a chain when referenced
			}
			if ast.IsImportDeclaration(n) {
				return false // likewise, these are ultimately what get marked by calls on other nodes - we want to skip them
			}
			r.checker.markLinkedReferences(n, ReferenceHintUnspecified, nil /*propSymbol*/, nil /*parentType*/)
			n.ForEachChild(visit)
			return false
		}
		file.ForEachChild(visit)
	}
}

func (r *emitResolver) GetExternalModuleFileFromDeclaration(node *ast.Node) *ast.SourceFile {
	r.checkerMu.Lock()
	defer r.checkerMu.Unlock()

	if ast.IsParseTreeNode(node) {
		// !!!
		// return r.checker.getExternalModuleFileFromDeclaration(node)
	}
	return nil
}

func (r *emitResolver) getReferenceResolver() binder.ReferenceResolver {
	if r.referenceResolver == nil {
		r.referenceResolver = binder.NewReferenceResolver(r.checker.compilerOptions, binder.ReferenceResolverHooks{
			ResolveName:                            r.checker.resolveName,
			GetResolvedSymbol:                      r.checker.getResolvedSymbol,
			GetMergedSymbol:                        r.checker.getMergedSymbol,
			GetParentOfSymbol:                      r.checker.getParentOfSymbol,
			GetSymbolOfDeclaration:                 r.checker.getSymbolOfDeclaration,
			GetTypeOnlyAliasDeclaration:            r.checker.getTypeOnlyAliasDeclarationEx,
			GetExportSymbolOfValueSymbolIfExported: r.checker.getExportSymbolOfValueSymbolIfExported,
		})
	}
	return r.referenceResolver
}

func (r *emitResolver) GetReferencedExportContainer(node *ast.IdentifierNode, prefixLocals bool) *ast.Node /*SourceFile|ModuleDeclaration|EnumDeclaration*/ {
	if !ast.IsParseTreeNode(node) {
		return nil
	}

	r.checkerMu.Lock()
	defer r.checkerMu.Unlock()

	return r.getReferenceResolver().GetReferencedExportContainer(node, prefixLocals)
}

func (r *emitResolver) GetReferencedImportDeclaration(node *ast.IdentifierNode) *ast.Declaration {
	if !ast.IsParseTreeNode(node) {
		return nil
	}

	r.checkerMu.Lock()
	defer r.checkerMu.Unlock()

	return r.getReferenceResolver().GetReferencedImportDeclaration(node)
}

func (r *emitResolver) GetReferencedValueDeclaration(node *ast.IdentifierNode) *ast.Declaration {
	if !ast.IsParseTreeNode(node) {
		return nil
	}

	r.checkerMu.Lock()
	defer r.checkerMu.Unlock()

	return r.getReferenceResolver().GetReferencedValueDeclaration(node)
}

func (r *emitResolver) GetReferencedValueDeclarations(node *ast.IdentifierNode) []*ast.Declaration {
	if !ast.IsParseTreeNode(node) {
		return nil
	}

	r.checkerMu.Lock()
	defer r.checkerMu.Unlock()

	return r.getReferenceResolver().GetReferencedValueDeclarations(node)
}
