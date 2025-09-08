package checker

import (
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/modulespecifiers"
	"github.com/microsoft/typescript-go/internal/nodebuilder"
)

func (b *nodeBuilderImpl) reuseNode(node *ast.Node) *ast.Node {
	// !!!
	return node
}

type recoveryBoundary struct {
	ctx                 *NodeBuilderContext
	hadError            bool
	deferredReports     []func()
	oldTracker          nodebuilder.SymbolTracker
	oldTrackedSymbols   []*TrackedSymbolArgs
	oldEncounteredError bool
}

func (b *recoveryBoundary) markError(f func()) {
	b.hadError = true
	if f != nil {
		b.deferredReports = append(b.deferredReports, f)
	}
}

type originalRecoveryScopeState struct {
	trackedSymbolsTop   int
	unreportedErrorsTop int
	hadError            bool
}

func (b *recoveryBoundary) startRecoveryScope() originalRecoveryScopeState {
	trackedSymbolsTop := len(b.ctx.trackedSymbols)
	unreportedErrorsTop := len(b.deferredReports)
	return originalRecoveryScopeState{trackedSymbolsTop: trackedSymbolsTop, unreportedErrorsTop: unreportedErrorsTop, hadError: b.hadError}
}

func (b *recoveryBoundary) endRecoveryScope(state originalRecoveryScopeState) {
	b.hadError = state.hadError
	b.ctx.trackedSymbols = b.ctx.trackedSymbols[0:state.trackedSymbolsTop]
	b.deferredReports = b.deferredReports[0:state.unreportedErrorsTop]
}

type wrappingTracker struct {
	wrapped nodebuilder.SymbolTracker
	bound   *recoveryBoundary
}

func (w *wrappingTracker) GetModuleSpecifierGenerationHost() modulespecifiers.ModuleSpecifierGenerationHost {
	return w.wrapped.GetModuleSpecifierGenerationHost()
}

func (w *wrappingTracker) PopErrorFallbackNode() {
	w.wrapped.PopErrorFallbackNode()
}

func (w *wrappingTracker) PushErrorFallbackNode(node *ast.Node) {
	w.wrapped.PushErrorFallbackNode(node)
}

func (w *wrappingTracker) ReportCyclicStructureError() {
	w.bound.markError(w.wrapped.ReportCyclicStructureError)
}

func (w *wrappingTracker) ReportInaccessibleThisError() {
	w.bound.markError(w.wrapped.ReportInaccessibleThisError)
}

func (w *wrappingTracker) ReportInaccessibleUniqueSymbolError() {
	w.bound.markError(w.wrapped.ReportInaccessibleUniqueSymbolError)
}

func (w *wrappingTracker) ReportInferenceFallback(node *ast.Node) {
	w.wrapped.ReportInferenceFallback(node) // Should this also be deferred?
}

func (w *wrappingTracker) ReportLikelyUnsafeImportRequiredError(specifier string) {
	w.bound.markError(func() { w.wrapped.ReportLikelyUnsafeImportRequiredError(specifier) })
}

func (w *wrappingTracker) ReportNonSerializableProperty(propertyName string) {
	w.bound.markError(func() { w.wrapped.ReportNonSerializableProperty(propertyName) })
}

func (w *wrappingTracker) ReportNonlocalAugmentation(containingFile *ast.SourceFile, parentSymbol *ast.Symbol, augmentingSymbol *ast.Symbol) {
	w.wrapped.ReportNonlocalAugmentation(containingFile, parentSymbol, augmentingSymbol) // Should this also be deferred?
}

func (w *wrappingTracker) ReportPrivateInBaseOfClassExpression(propertyName string) {
	w.bound.markError(func() { w.wrapped.ReportPrivateInBaseOfClassExpression(propertyName) })
}

func (w *wrappingTracker) ReportTruncationError() {
	w.wrapped.ReportTruncationError() // Should this also be deferred?
}

func (w *wrappingTracker) TrackSymbol(symbol *ast.Symbol, enclosingDeclaration *ast.Node, meaning ast.SymbolFlags) bool {
	w.bound.ctx.trackedSymbols = append(w.bound.ctx.trackedSymbols, &TrackedSymbolArgs{symbol, enclosingDeclaration, meaning})
	return false
}

func newWrappingTracker(inner nodebuilder.SymbolTracker, bound *recoveryBoundary) *wrappingTracker {
	return &wrappingTracker{
		wrapped: inner,
		bound:   bound,
	}
}

func (b *nodeBuilderImpl) createRecoveryBoundary() *recoveryBoundary {
	b.ch.checkNotCanceled()
	bound := &recoveryBoundary{oldTracker: b.ctx.tracker, oldTrackedSymbols: b.ctx.trackedSymbols, oldEncounteredError: b.ctx.encounteredError}
	newTracker := NewSymbolTrackerImpl(b.ctx, newWrappingTracker(b.ctx.tracker, bound), b.ctx.tracker.GetModuleSpecifierGenerationHost())
	b.ctx.tracker = newTracker
	b.ctx.trackedSymbols = nil
	return bound
}

func (b *nodeBuilderImpl) finalizeBoundary(bound *recoveryBoundary) bool {
	b.ctx.tracker = bound.oldTracker
	b.ctx.trackedSymbols = bound.oldTrackedSymbols
	b.ctx.encounteredError = bound.oldEncounteredError

	for _, f := range bound.deferredReports {
		f()
	}
	if bound.hadError {
		return false
	}
	for _, a := range b.ctx.trackedSymbols {
		b.ctx.tracker.TrackSymbol(a.symbol, a.enclosingDeclaration, a.meaning)
	}
	return true
}

func (b *nodeBuilderImpl) tryReuseExistingTypeNodeHelper(existing *ast.TypeNode) *ast.TypeNode {
	bound := b.createRecoveryBoundary()
	var transformed *ast.Node
	// !!!
	if !b.finalizeBoundary(bound) {
		return nil
	}
	b.ctx.approximateLength += existing.Loc.End() - existing.Loc.Pos()
	return transformed
}

func getExistingNodeTreeVisitor(b *nodeBuilderImpl, bound *recoveryBoundary, factory *ast.NodeFactory) *ast.NodeVisitor {
	visitExistingNodeTreeSymbolsWorker := func(node *ast.Node) *ast.Node {
		return node // !!!
	}
	return ast.NewNodeVisitor(func(node *ast.Node) *ast.Node {
		if bound.hadError {
			return node
		}
		recover := bound.startRecoveryScope()
		introducesNewScope := ast.IsFunctionLike(node) || ast.IsMappedTypeNode(node)
		var exit func()
		if introducesNewScope {
			var params []*ast.Symbol
			var typeParams []*Type
			if ast.IsFunctionLike(node) {
				sig := b.ch.getSignatureFromDeclaration(node)
				params = sig.parameters
				typeParams = sig.typeParameters
			} else if ast.IsConditionalTypeNode(node) { // !!! TODO: impossible in combination with the scope start check???
				typeParams = b.ch.getInferTypeParameters(node)
			} else if ast.IsMappedTypeNode(node) {
				typeParams = []*Type{b.ch.getDeclaredTypeOfTypeParameter(b.ch.getSymbolOfDeclaration(node.AsMappedTypeNode().TypeParameter))}
			}
			exit = b.enterNewScope(node, params, typeParams, nil, nil)
		}
		result := visitExistingNodeTreeSymbolsWorker(node)
		if exit != nil {
			exit()
		}

		if bound.hadError {
			if ast.IsTypeNode(node) && !ast.IsTypePredicateNode(node) {
				bound.endRecoveryScope(recover)
				return nil // !!!
			}
			return node
		}

		return result
	}, factory, ast.NodeVisitorHooks{})
}
