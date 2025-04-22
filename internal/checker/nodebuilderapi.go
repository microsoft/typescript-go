package checker

import (
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/nodebuilder"
	"github.com/microsoft/typescript-go/internal/printer"
)

type NodeBuilderInterface interface {
	typeToTypeNode(typ *Type, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node
	typePredicateToTypePredicateNode(predicate *TypePredicate, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node
	serializeTypeForDeclaration(declaration *ast.Node, symbol *ast.Symbol, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node
	serializeReturnTypeForSignature(signatureDeclaration *ast.Node, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node
	serializeTypeForExpression(expr *ast.Node, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node
	indexInfoToIndexSignatureDeclaration(info *IndexInfo, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node
	signatureToSignatureDeclaration(signature *Signature, kind ast.Kind, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node
	symbolToEntityName(symbol *ast.Symbol, meaning ast.SymbolFlags, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node
	symbolToExpression(symbol *ast.Symbol, meaning ast.SymbolFlags, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node
	symbolToTypeParameterDeclarations(symbol *ast.Symbol, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) []*ast.Node
	symbolToParameterDeclaration(symbol *ast.Symbol, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node
	typeParameterToDeclaration(parameter *Type, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node
	symbolTableToDeclarationStatements(symbolTable *ast.SymbolTable, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) []*ast.Node
	symbolToNode(symbol *ast.Symbol, meaning ast.SymbolFlags, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node
}

type NodeBuilderAPI struct {
	ctxStack []*NodeBuilderContext
	impl     *NodeBuilder
}

func (b *NodeBuilderAPI) enterContext(enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) {
	b.impl.ctx = &NodeBuilderContext{
		tracker:              tracker,
		approximateLength:    0,
		encounteredError:     false,
		truncating:           false,
		reportedDiagnostic:   false,
		flags:                flags,
		internalFlags:        internalFlags,
		depth:                0,
		enclosingDeclaration: enclosingDeclaration,
		enclosingFile:        ast.GetSourceFileOfNode(enclosingDeclaration),
		inferTypeParameters:  make([]*Type, 0),
		visitedTypes:         make(map[TypeId]bool),
		symbolDepth:          make(map[CompositeSymbolIdentity]int),
		trackedSymbols:       make([]*TrackedSymbolArgs, 0),
		mapper:               nil,
		reverseMappedStack:   make([]*ast.Symbol, 0),
	}
	if tracker == nil {
		tracker = NewSymbolTrackerImpl(b.impl.ctx, nil)
		b.impl.ctx.tracker = tracker
	}
	b.impl.initializeClosures() // recapture ctx
	b.ctxStack = append(b.ctxStack, b.impl.ctx)
}

func (b *NodeBuilderAPI) popContext() {
	b.impl.ctx = nil
	if len(b.ctxStack) > 1 {
		b.impl.ctx = b.ctxStack[len(b.ctxStack)-1]
	}
	b.ctxStack = b.ctxStack[0 : len(b.ctxStack)-1]
}

func (b *NodeBuilderAPI) exitContext(result *ast.Node) *ast.Node {
	b.exitContextCheck()
	defer b.popContext()
	if b.impl.ctx.encounteredError {
		return nil
	}
	return result
}

func (b *NodeBuilderAPI) exitContextSlice(result []*ast.Node) []*ast.Node {
	b.exitContextCheck()
	defer b.popContext()
	if b.impl.ctx.encounteredError {
		return nil
	}
	return result
}

func (b *NodeBuilderAPI) exitContextCheck() {
	if b.impl.ctx.truncating && b.impl.ctx.flags&nodebuilder.FlagsNoTruncation != 0 {
		b.impl.ctx.tracker.ReportTruncationError()
	}
}

// indexInfoToIndexSignatureDeclaration implements NodeBuilderInterface.
func (b *NodeBuilderAPI) indexInfoToIndexSignatureDeclaration(info *IndexInfo, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node {
	b.enterContext(enclosingDeclaration, flags, internalFlags, tracker)
	return b.exitContext(b.impl.indexInfoToIndexSignatureDeclarationHelper(info, nil))
}

// serializeReturnTypeForSignature implements NodeBuilderInterface.
func (b *NodeBuilderAPI) serializeReturnTypeForSignature(signatureDeclaration *ast.Node, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node {
	b.enterContext(enclosingDeclaration, flags, internalFlags, tracker)
	signature := b.impl.ch.getSignatureFromDeclaration(signatureDeclaration)
	symbol := b.impl.ch.getSymbolOfDeclaration(signatureDeclaration)
	returnType, ok := b.impl.ctx.enclosingSymbolTypes[ast.GetSymbolId(symbol)]
	if !ok || returnType == nil {
		returnType = b.impl.ch.instantiateType(b.impl.ch.getReturnTypeOfSignature(signature), b.impl.ctx.mapper)
	}
	return b.exitContext(b.impl.serializeInferredReturnTypeForSignature(signature, returnType))
}

// serializeTypeForDeclaration implements NodeBuilderInterface.
func (b *NodeBuilderAPI) serializeTypeForDeclaration(declaration *ast.Node, symbol *ast.Symbol, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node {
	b.enterContext(enclosingDeclaration, flags, internalFlags, tracker)
	return b.exitContext(b.impl.serializeTypeForDeclaration(declaration, nil, symbol))
}

// serializeTypeForExpression implements NodeBuilderInterface.
func (b *NodeBuilderAPI) serializeTypeForExpression(expr *ast.Node, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node {
	b.enterContext(enclosingDeclaration, flags, internalFlags, tracker)
	return b.exitContext(b.impl.serializeTypeForExpression(expr))
}

// signatureToSignatureDeclaration implements NodeBuilderInterface.
func (b *NodeBuilderAPI) signatureToSignatureDeclaration(signature *Signature, kind ast.Kind, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node {
	b.enterContext(enclosingDeclaration, flags, internalFlags, tracker)
	return b.exitContext(b.impl.signatureToSignatureDeclarationHelper(signature, kind, nil))
}

// symbolTableToDeclarationStatements implements NodeBuilderInterface.
func (b *NodeBuilderAPI) symbolTableToDeclarationStatements(symbolTable *ast.SymbolTable, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) []*ast.Node {
	b.enterContext(enclosingDeclaration, flags, internalFlags, tracker)
	return b.exitContextSlice(b.impl.symbolTableToDeclarationStatements(symbolTable))
}

// symbolToEntityName implements NodeBuilderInterface.
func (b *NodeBuilderAPI) symbolToEntityName(symbol *ast.Symbol, meaning ast.SymbolFlags, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node {
	b.enterContext(enclosingDeclaration, flags, internalFlags, tracker)
	return b.exitContext(b.impl.symbolToName(symbol, meaning, false))
}

// symbolToExpression implements NodeBuilderInterface.
func (b *NodeBuilderAPI) symbolToExpression(symbol *ast.Symbol, meaning ast.SymbolFlags, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node {
	b.enterContext(enclosingDeclaration, flags, internalFlags, tracker)
	return b.exitContext(b.impl.symbolToExpression(symbol, meaning))
}

// symbolToNode implements NodeBuilderInterface.
func (b *NodeBuilderAPI) symbolToNode(symbol *ast.Symbol, meaning ast.SymbolFlags, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node {
	b.enterContext(enclosingDeclaration, flags, internalFlags, tracker)
	return b.exitContext(b.impl.symbolToNode(symbol, meaning))
}

// symbolToParameterDeclaration implements NodeBuilderInterface.
func (b NodeBuilderAPI) symbolToParameterDeclaration(symbol *ast.Symbol, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node {
	b.enterContext(enclosingDeclaration, flags, internalFlags, tracker)
	return b.exitContext(b.impl.symbolToParameterDeclaration(symbol, false))
}

// symbolToTypeParameterDeclarations implements NodeBuilderInterface.
func (b *NodeBuilderAPI) symbolToTypeParameterDeclarations(symbol *ast.Symbol, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) []*ast.Node {
	b.enterContext(enclosingDeclaration, flags, internalFlags, tracker)
	return b.exitContextSlice(b.impl.symbolToTypeParameterDeclarations(symbol))
}

// typeParameterToDeclaration implements NodeBuilderInterface.
func (b *NodeBuilderAPI) typeParameterToDeclaration(parameter *Type, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node {
	b.enterContext(enclosingDeclaration, flags, internalFlags, tracker)
	return b.exitContext(b.impl.typeParameterToDeclaration(parameter))
}

// typePredicateToTypePredicateNode implements NodeBuilderInterface.
func (b *NodeBuilderAPI) typePredicateToTypePredicateNode(predicate *TypePredicate, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node {
	b.enterContext(enclosingDeclaration, flags, internalFlags, tracker)
	return b.exitContext(b.impl.typePredicateToTypePredicateNode(predicate))
}

// typeToTypeNode implements NodeBuilderInterface.
func (b *NodeBuilderAPI) typeToTypeNode(typ *Type, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node {
	b.enterContext(enclosingDeclaration, flags, internalFlags, tracker)
	return b.exitContext(b.impl.typeToTypeNode(typ))
}

// var _ NodeBuilderInterface = NewNodeBuilderAPI(nil, nil)

func NewNodeBuilderAPI(ch *Checker, e *printer.EmitContext) *NodeBuilderAPI {
	impl := NewNodeBuilder(ch, e)
	return &NodeBuilderAPI{impl: &impl, ctxStack: make([]*NodeBuilderContext, 0, 1)}
}
