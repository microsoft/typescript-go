package checker

import (
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/nodebuilder"
)

type SymbolTrackerImpl struct {
	context            NodeBuilderContext
	inner              nodebuilder.SymbolTracker
	DisableTrackSymbol bool
}

func NewSymbolTrackerImpl(context NodeBuilderContext, tracker nodebuilder.SymbolTracker) *SymbolTrackerImpl {
	// TODO: unwrap `tracker` before setting `inner`
	return &SymbolTrackerImpl{context, tracker, false}
}

func (this *SymbolTrackerImpl) GetModuleSpecifierGenerationHost() any {
	return this.inner.GetModuleSpecifierGenerationHost()
}

func (this *SymbolTrackerImpl) TrackSymbol(symbol *ast.Symbol, enclosingDeclaration *ast.Node, meaning ast.SymbolFlags) bool {
	if !this.DisableTrackSymbol {
		if this.inner.TrackSymbol(symbol, enclosingDeclaration, meaning) {
			this.onDiagnosticReported()
			return true
		}
		// Skip recording type parameters as they dont contribute to late painted statements
		if symbol.Flags&ast.SymbolFlagsTypeParameter == 0 {
			this.context.trackedSymbols = append(this.context.trackedSymbols, &TrackedSymbolArgs{symbol, enclosingDeclaration, meaning})
		}
	}
	return false
}

func (this *SymbolTrackerImpl) ReportInaccessibleThisError() {
	this.onDiagnosticReported()
	this.inner.ReportInaccessibleThisError()
}

func (this *SymbolTrackerImpl) ReportPrivateInBaseOfClassExpression(propertyName string) {
	this.onDiagnosticReported()
	this.inner.ReportPrivateInBaseOfClassExpression(propertyName)
}

func (this *SymbolTrackerImpl) ReportInaccessibleUniqueSymbolError() {
	this.onDiagnosticReported()
	this.inner.ReportInaccessibleUniqueSymbolError()
}

func (this *SymbolTrackerImpl) ReportCyclicStructureError() {
	this.onDiagnosticReported()
	this.inner.ReportCyclicStructureError()
}

func (this *SymbolTrackerImpl) ReportLikelyUnsafeImportRequiredError(specifier string) {
	this.onDiagnosticReported()
	this.inner.ReportLikelyUnsafeImportRequiredError(specifier)
}

func (this *SymbolTrackerImpl) ReportTruncationError() {
	this.onDiagnosticReported()
	this.inner.ReportTruncationError()
}

func (this *SymbolTrackerImpl) ReportNonlocalAugmentation(containingFile *ast.SourceFile, parentSymbol *ast.Symbol, augmentingSymbol *ast.Symbol) {
	this.onDiagnosticReported()
	this.inner.ReportNonlocalAugmentation(containingFile, parentSymbol, augmentingSymbol)
}

func (this *SymbolTrackerImpl) ReportNonSerializableProperty(propertyName string) {
	this.onDiagnosticReported()
	this.inner.ReportNonSerializableProperty(propertyName)
}

func (this *SymbolTrackerImpl) onDiagnosticReported() {
	this.context.reportedDiagnostic = true
}

func (this *SymbolTrackerImpl) ReportInferenceFallback(node *ast.Node) {
	this.inner.ReportInferenceFallback(node)
}

func (this *SymbolTrackerImpl) PushErrorFallbackNode(node *ast.Node) {
	this.inner.PushErrorFallbackNode(node)
}

func (this *SymbolTrackerImpl) PopErrorFallbackNode() {
	this.inner.PopErrorFallbackNode()
}
