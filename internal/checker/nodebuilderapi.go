package checker

import (
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/nodebuilder"
	"github.com/microsoft/typescript-go/internal/printer"
)

type NodeBuilderInterface interface {
	EmitContext() *printer.EmitContext

	TypeToTypeNode(typ *Type, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node
	TypePredicateToTypePredicateNode(predicate *TypePredicate, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node
	SerializeTypeForDeclaration(declaration *ast.Node, symbol *ast.Symbol, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node
	SerializeReturnTypeForSignature(signatureDeclaration *ast.Node, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node
	SerializeTypeForExpression(expr *ast.Node, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node
	IndexInfoToIndexSignatureDeclaration(info *IndexInfo, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node
	SignatureToSignatureDeclaration(signature *Signature, kind ast.Kind, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node
	SymbolToEntityName(symbol *ast.Symbol, meaning ast.SymbolFlags, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node
	SymbolToExpression(symbol *ast.Symbol, meaning ast.SymbolFlags, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node
	SymbolToTypeParameterDeclarations(symbol *ast.Symbol, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) []*ast.Node
	SymbolToParameterDeclaration(symbol *ast.Symbol, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node
	TypeParameterToDeclaration(parameter *Type, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node
	SymbolTableToDeclarationStatements(symbolTable *ast.SymbolTable, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) []*ast.Node
	SymbolToNode(symbol *ast.Symbol, meaning ast.SymbolFlags, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node
}

type NodeBuilderAPI struct {
	ctxStack []*NodeBuilderContext
	impl     *NodeBuilder
}

// EmitContext implements NodeBuilderInterface.
func (b *NodeBuilderAPI) EmitContext() *printer.EmitContext {
	return b.impl.e
}

func (b *NodeBuilderAPI) enterContext(enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) {
	b.impl.ctx = &NodeBuilderContext{
		tracker:                  tracker,
		approximateLength:        0,
		encounteredError:         false,
		truncating:               false,
		reportedDiagnostic:       false,
		flags:                    flags,
		internalFlags:            internalFlags,
		depth:                    0,
		enclosingDeclaration:     enclosingDeclaration,
		enclosingFile:            ast.GetSourceFileOfNode(enclosingDeclaration),
		inferTypeParameters:      make([]*Type, 0),
		visitedTypes:             make(map[TypeId]bool),
		symbolDepth:              make(map[CompositeSymbolIdentity]int),
		trackedSymbols:           make([]*TrackedSymbolArgs, 0),
		mapper:                   nil,
		reverseMappedStack:       make([]*ast.Symbol, 0),
		enclosingSymbolTypes:     make(map[ast.SymbolId]*Type),
		remappedSymbolReferences: make(map[ast.SymbolId]*ast.Symbol),
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

// IndexInfoToIndexSignatureDeclaration implements NodeBuilderInterface.
func (b *NodeBuilderAPI) IndexInfoToIndexSignatureDeclaration(info *IndexInfo, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node {
	b.enterContext(enclosingDeclaration, flags, internalFlags, tracker)
	return b.exitContext(b.impl.indexInfoToIndexSignatureDeclarationHelper(info, nil))
}

// SerializeReturnTypeForSignature implements NodeBuilderInterface.
func (b *NodeBuilderAPI) SerializeReturnTypeForSignature(signatureDeclaration *ast.Node, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node {
	b.enterContext(enclosingDeclaration, flags, internalFlags, tracker)
	signature := b.impl.ch.getSignatureFromDeclaration(signatureDeclaration)
	symbol := b.impl.ch.getSymbolOfDeclaration(signatureDeclaration)
	returnType, ok := b.impl.ctx.enclosingSymbolTypes[ast.GetSymbolId(symbol)]
	if !ok || returnType == nil {
		returnType = b.impl.ch.instantiateType(b.impl.ch.getReturnTypeOfSignature(signature), b.impl.ctx.mapper)
	}
	return b.exitContext(b.impl.serializeInferredReturnTypeForSignature(signature, returnType))
}

// SerializeTypeForDeclaration implements NodeBuilderInterface.
func (b *NodeBuilderAPI) SerializeTypeForDeclaration(declaration *ast.Node, symbol *ast.Symbol, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node {
	b.enterContext(enclosingDeclaration, flags, internalFlags, tracker)
	return b.exitContext(b.impl.serializeTypeForDeclaration(declaration, nil, symbol))
}

// SerializeTypeForExpression implements NodeBuilderInterface.
func (b *NodeBuilderAPI) SerializeTypeForExpression(expr *ast.Node, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node {
	b.enterContext(enclosingDeclaration, flags, internalFlags, tracker)
	return b.exitContext(b.impl.serializeTypeForExpression(expr))
}

// SignatureToSignatureDeclaration implements NodeBuilderInterface.
func (b *NodeBuilderAPI) SignatureToSignatureDeclaration(signature *Signature, kind ast.Kind, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node {
	b.enterContext(enclosingDeclaration, flags, internalFlags, tracker)
	return b.exitContext(b.impl.signatureToSignatureDeclarationHelper(signature, kind, nil))
}

// SymbolTableToDeclarationStatements implements NodeBuilderInterface.
func (b *NodeBuilderAPI) SymbolTableToDeclarationStatements(symbolTable *ast.SymbolTable, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) []*ast.Node {
	b.enterContext(enclosingDeclaration, flags, internalFlags, tracker)
	return b.exitContextSlice(b.impl.symbolTableToDeclarationStatements(symbolTable))
}

// SymbolToEntityName implements NodeBuilderInterface.
func (b *NodeBuilderAPI) SymbolToEntityName(symbol *ast.Symbol, meaning ast.SymbolFlags, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node {
	b.enterContext(enclosingDeclaration, flags, internalFlags, tracker)
	return b.exitContext(b.impl.symbolToName(symbol, meaning, false))
}

// SymbolToExpression implements NodeBuilderInterface.
func (b *NodeBuilderAPI) SymbolToExpression(symbol *ast.Symbol, meaning ast.SymbolFlags, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node {
	b.enterContext(enclosingDeclaration, flags, internalFlags, tracker)
	return b.exitContext(b.impl.symbolToExpression(symbol, meaning))
}

// SymbolToNode implements NodeBuilderInterface.
func (b *NodeBuilderAPI) SymbolToNode(symbol *ast.Symbol, meaning ast.SymbolFlags, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node {
	b.enterContext(enclosingDeclaration, flags, internalFlags, tracker)
	return b.exitContext(b.impl.symbolToNode(symbol, meaning))
}

// SymbolToParameterDeclaration implements NodeBuilderInterface.
func (b NodeBuilderAPI) SymbolToParameterDeclaration(symbol *ast.Symbol, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node {
	b.enterContext(enclosingDeclaration, flags, internalFlags, tracker)
	return b.exitContext(b.impl.symbolToParameterDeclaration(symbol, false))
}

// SymbolToTypeParameterDeclarations implements NodeBuilderInterface.
func (b *NodeBuilderAPI) SymbolToTypeParameterDeclarations(symbol *ast.Symbol, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) []*ast.Node {
	b.enterContext(enclosingDeclaration, flags, internalFlags, tracker)
	return b.exitContextSlice(b.impl.symbolToTypeParameterDeclarations(symbol))
}

// TypeParameterToDeclaration implements NodeBuilderInterface.
func (b *NodeBuilderAPI) TypeParameterToDeclaration(parameter *Type, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node {
	b.enterContext(enclosingDeclaration, flags, internalFlags, tracker)
	return b.exitContext(b.impl.typeParameterToDeclaration(parameter))
}

// TypePredicateToTypePredicateNode implements NodeBuilderInterface.
func (b *NodeBuilderAPI) TypePredicateToTypePredicateNode(predicate *TypePredicate, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node {
	b.enterContext(enclosingDeclaration, flags, internalFlags, tracker)
	return b.exitContext(b.impl.typePredicateToTypePredicateNode(predicate))
}

// TypeToTypeNode implements NodeBuilderInterface.
func (b *NodeBuilderAPI) TypeToTypeNode(typ *Type, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node {
	b.enterContext(enclosingDeclaration, flags, internalFlags, tracker)
	return b.exitContext(b.impl.typeToTypeNode(typ))
}

// var _ NodeBuilderInterface = NewNodeBuilderAPI(nil, nil)

func NewNodeBuilderAPI(ch *Checker, e *printer.EmitContext) *NodeBuilderAPI {
	impl := NewNodeBuilder(ch, e)
	return &NodeBuilderAPI{impl: &impl, ctxStack: make([]*NodeBuilderContext, 0, 1)}
}

func (c *Checker) GetDiagnosticNodeBuilder() NodeBuilderInterface {
	return c.nodeBuilder
}
