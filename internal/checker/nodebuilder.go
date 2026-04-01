package checker

import (
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/nodebuilder"
	"github.com/microsoft/typescript-go/internal/printer"
)

type NodeBuilder struct {
	ctxStack []*NodeBuilderContext
	host     Host
	impl     *NodeBuilderImpl
}

// EmitContext implements NodeBuilderInterface.
func (b *NodeBuilder) EmitContext() *printer.EmitContext {
	return b.impl.e
}

func (b *NodeBuilder) enterContext(enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) {
	b.enterContextEx(enclosingDeclaration, flags, internalFlags, tracker, -1, nil, 0)
}

func (b *NodeBuilder) enterContextEx(enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker, verbosityLevel int, out *WriterContextOut, maxTruncationLength int) {
	b.ctxStack = append(b.ctxStack, b.impl.ctx)
	b.impl.ctx = &NodeBuilderContext{
		host:                     b.host,
		tracker:                  tracker,
		flags:                    flags,
		internalFlags:            internalFlags,
		maxExpansionDepth:        verbosityLevel,
		maxTruncationLength:      maxTruncationLength,
		enclosingDeclaration:     enclosingDeclaration,
		enclosingFile:            ast.GetSourceFileOfNode(enclosingDeclaration),
		inferTypeParameters:      make([]*Type, 0),
		symbolDepth:              make(map[CompositeSymbolIdentity]int),
		trackedSymbols:           make([]*TrackedSymbolArgs, 0),
		reverseMappedStack:       make([]*ast.Symbol, 0),
		enclosingSymbolTypes:     make(map[ast.SymbolId]*Type),
		remappedSymbolReferences: make(map[ast.SymbolId]*ast.Symbol),
	}
	tracker = NewSymbolTrackerImpl(b.impl.ctx, tracker)
	b.impl.ctx.tracker = tracker
}

func (b *NodeBuilder) popContext() {
	stackSize := len(b.ctxStack)
	if stackSize == 0 {
		b.impl.ctx = nil
	} else {
		b.impl.ctx = b.ctxStack[stackSize-1]
		b.ctxStack = b.ctxStack[:stackSize-1]
	}
}

func (b *NodeBuilder) exitContext(result *ast.Node) *ast.Node {
	b.exitContextCheck()
	defer b.popContext()
	if b.impl.ctx.encounteredError {
		return nil
	}
	return result
}

func (b *NodeBuilder) exitContextSlice(result []*ast.Node) []*ast.Node {
	b.exitContextCheck()
	defer b.popContext()
	if b.impl.ctx.encounteredError {
		return nil
	}
	return result
}

func (b *NodeBuilder) exitContextCheck() {
	if b.impl.ctx.truncating && b.impl.ctx.flags&nodebuilder.FlagsNoTruncation != 0 {
		b.impl.ctx.tracker.ReportTruncationError()
	}
}

// IndexInfoToIndexSignatureDeclaration implements NodeBuilderInterface.
func (b *NodeBuilder) IndexInfoToIndexSignatureDeclaration(info *IndexInfo, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node {
	b.enterContext(enclosingDeclaration, flags, internalFlags, tracker)
	return b.exitContext(b.impl.indexInfoToIndexSignatureDeclarationHelper(info, nil))
}

// SerializeReturnTypeForSignature implements NodeBuilderInterface.
func (b *NodeBuilder) SerializeReturnTypeForSignature(signatureDeclaration *ast.Node, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node {
	b.enterContext(enclosingDeclaration, flags, internalFlags, tracker)
	signature := b.impl.ch.getSignatureFromDeclaration(signatureDeclaration)
	return b.exitContext(b.impl.serializeReturnTypeForSignature(signature, true))
}

func (b *NodeBuilder) SerializeTypeParametersForSignature(signatureDeclaration *ast.Node, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) []*ast.Node {
	b.enterContext(enclosingDeclaration, flags, internalFlags, tracker)
	symbol := b.impl.ch.getSymbolOfDeclaration(signatureDeclaration)
	typeParams := b.SymbolToTypeParameterDeclarations(symbol, enclosingDeclaration, flags, internalFlags, tracker)
	return b.exitContextSlice(typeParams)
}

// SerializeTypeForDeclaration implements NodeBuilderInterface.
func (b *NodeBuilder) SerializeTypeForDeclaration(declaration *ast.Node, symbol *ast.Symbol, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node {
	b.enterContext(enclosingDeclaration, flags, internalFlags, tracker)
	return b.exitContext(b.impl.serializeTypeForDeclaration(declaration, nil, symbol, true))
}

// SerializeTypeForExpression implements NodeBuilderInterface.
func (b *NodeBuilder) SerializeTypeForExpression(expr *ast.Node, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node {
	b.enterContext(enclosingDeclaration, flags, internalFlags, tracker)
	return b.exitContext(b.impl.serializeTypeForExpression(expr))
}

// SignatureToSignatureDeclaration implements NodeBuilderInterface.
func (b *NodeBuilder) SignatureToSignatureDeclaration(signature *Signature, kind ast.Kind, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node {
	b.enterContext(enclosingDeclaration, flags, internalFlags, tracker)
	return b.exitContext(b.impl.signatureToSignatureDeclarationHelper(signature, kind, nil))
}

// SymbolTableToDeclarationStatements implements NodeBuilderInterface.
func (b *NodeBuilder) SymbolTableToDeclarationStatements(symbolTable *ast.SymbolTable, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) []*ast.Node {
	b.enterContext(enclosingDeclaration, flags, internalFlags, tracker)
	return b.exitContextSlice(b.impl.symbolTableToDeclarationStatements(symbolTable))
}

// SymbolToDeclarationsWithVerbosity produces declaration nodes for a symbol with verbosity level support.
func (b *NodeBuilder) SymbolToDeclarationsWithVerbosity(symbol *ast.Symbol, meaning ast.SymbolFlags, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker, verbosityLevel int, out *WriterContextOut, maxTruncationLength int) []*ast.Node {
	b.enterContextEx(enclosingDeclaration, flags, internalFlags, tracker, verbosityLevel, out, maxTruncationLength)

	// Push the declared type onto the type stack to prevent re-expansion.
	// We push a sentinel 0 after the real id so that the cycle-detection loops
	// in shouldExpandType/canPossiblyExpandType (which skip the last element
	// via `range len(typeStack)-1`) still check declaredType.id.
	declaredType := b.impl.ch.getDeclaredTypeOfSymbol(symbol)
	b.impl.ctx.typeStack = append(b.impl.ctx.typeStack, declaredType.id)
	b.impl.ctx.typeStack = append(b.impl.ctx.typeStack, 0)

	table := createSymbolTable([]*ast.Symbol{symbol})
	nodes := b.impl.symbolTableToDeclarationStatements(&table)

	b.impl.ctx.typeStack = b.impl.ctx.typeStack[:len(b.impl.ctx.typeStack)-2]

	if out != nil {
		out.CanIncreaseExpansionDepth = b.impl.ctx.out.CanIncreaseExpansionDepth
		out.Truncated = b.impl.ctx.out.Truncated
	}

	// Simplify declarations by applying original modifiers
	result := make([]*ast.Node, 0, len(nodes))
	for _, node := range nodes {
		switch node.Kind {
		case ast.KindClassDeclaration:
			result = append(result, b.simplifyClassDeclaration(node, symbol))
		case ast.KindEnumDeclaration:
			result = append(result, simplifyModifiers(b.impl.f, node, ast.IsEnumDeclaration, symbol))
		case ast.KindInterfaceDeclaration:
			if meaning&ast.SymbolFlagsInterface != 0 {
				result = append(result, simplifyModifiers(b.impl.f, node, ast.IsInterfaceDeclaration, symbol))
			}
		case ast.KindModuleDeclaration:
			result = append(result, simplifyModifiers(b.impl.f, node, ast.IsModuleDeclaration, symbol))
		}
	}

	return b.exitContextSlice(result)
}

func (b *NodeBuilder) simplifyClassDeclaration(classDecl *ast.Node, symbol *ast.Symbol) *ast.Node {
	classDeclarations := core.Filter(symbol.Declarations, ast.IsClassLike)
	var originalClassDecl *ast.Node
	if len(classDeclarations) > 0 {
		originalClassDecl = classDeclarations[0]
	} else {
		originalClassDecl = classDecl
	}
	modifiers := originalClassDecl.ModifierFlags() & ^(ast.ModifierFlagsExport | ast.ModifierFlagsAmbient)
	isAnonymous := ast.IsClassExpression(originalClassDecl)
	if isAnonymous {
		cd := classDecl.AsClassDeclaration()
		classDecl = b.impl.f.UpdateClassDeclaration(
			cd,
			classDecl.Modifiers(),
			nil,
			cd.TypeParameters,
			cd.HeritageClauses,
			cd.Members,
		)
	}
	return ast.ReplaceModifiers(b.impl.f, classDecl, b.impl.f.NewModifierList(ast.CreateModifiersFromModifierFlags(modifiers, b.impl.f.NewModifier)))
}

func simplifyModifiers(f *ast.NodeFactory, newDecl *ast.Node, isDeclKind func(*ast.Node) bool, symbol *ast.Symbol) *ast.Node {
	decls := core.Filter(symbol.Declarations, isDeclKind)
	var declWithModifiers *ast.Node
	if len(decls) > 0 {
		declWithModifiers = decls[0]
	} else {
		declWithModifiers = newDecl
	}
	modifiers := declWithModifiers.ModifierFlags() & ^(ast.ModifierFlagsExport | ast.ModifierFlagsAmbient)
	return ast.ReplaceModifiers(f, newDecl, f.NewModifierList(ast.CreateModifiersFromModifierFlags(modifiers, f.NewModifier)))
}

// SymbolToEntityName implements NodeBuilderInterface.
func (b *NodeBuilder) SymbolToEntityName(symbol *ast.Symbol, meaning ast.SymbolFlags, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node {
	b.enterContext(enclosingDeclaration, flags, internalFlags, tracker)
	return b.exitContext(b.impl.symbolToName(symbol, meaning, false))
}

// SymbolToExpression implements NodeBuilderInterface.
func (b *NodeBuilder) SymbolToExpression(symbol *ast.Symbol, meaning ast.SymbolFlags, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node {
	b.enterContext(enclosingDeclaration, flags, internalFlags, tracker)
	return b.exitContext(b.impl.symbolToExpression(symbol, meaning))
}

// SymbolToNode implements NodeBuilderInterface.
func (b *NodeBuilder) SymbolToNode(symbol *ast.Symbol, meaning ast.SymbolFlags, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node {
	b.enterContext(enclosingDeclaration, flags, internalFlags, tracker)
	return b.exitContext(b.impl.symbolToNode(symbol, meaning))
}

// SymbolToParameterDeclaration implements NodeBuilderInterface.
func (b NodeBuilder) SymbolToParameterDeclaration(symbol *ast.Symbol, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node {
	b.enterContext(enclosingDeclaration, flags, internalFlags, tracker)
	return b.exitContext(b.impl.symbolToParameterDeclaration(symbol, false))
}

// SymbolToTypeParameterDeclarations implements NodeBuilderInterface.
func (b *NodeBuilder) SymbolToTypeParameterDeclarations(symbol *ast.Symbol, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) []*ast.Node {
	b.enterContext(enclosingDeclaration, flags, internalFlags, tracker)
	return b.exitContextSlice(b.impl.symbolToTypeParameterDeclarations(symbol))
}

// TypeParameterToDeclaration implements NodeBuilderInterface.
func (b *NodeBuilder) TypeParameterToDeclaration(parameter *Type, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node {
	b.enterContext(enclosingDeclaration, flags, internalFlags, tracker)
	return b.exitContext(b.impl.typeParameterToDeclaration(parameter))
}

// TypeParameterToDeclarationWithVerbosity is like TypeParameterToDeclaration but with verbosity level support.
func (b *NodeBuilder) TypeParameterToDeclarationWithVerbosity(parameter *Type, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker, verbosityLevel int, out *WriterContextOut) *ast.Node {
	b.enterContextEx(enclosingDeclaration, flags, internalFlags, tracker, verbosityLevel, out, 0)
	result := b.impl.typeParameterToDeclaration(parameter)
	if out != nil {
		out.CanIncreaseExpansionDepth = b.impl.ctx.out.CanIncreaseExpansionDepth
		out.Truncated = b.impl.ctx.out.Truncated
	}
	return b.exitContext(result)
}

// TypePredicateToTypePredicateNode implements NodeBuilderInterface.
func (b *NodeBuilder) TypePredicateToTypePredicateNode(predicate *TypePredicate, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node {
	b.enterContext(enclosingDeclaration, flags, internalFlags, tracker)
	return b.exitContext(b.impl.typePredicateToTypePredicateNode(predicate))
}

// TypeToTypeNode implements NodeBuilderInterface.
func (b *NodeBuilder) TypeToTypeNode(typ *Type, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node {
	b.enterContext(enclosingDeclaration, flags, internalFlags, tracker)
	return b.exitContext(b.impl.typeToTypeNode(typ))
}

// TypeToTypeNodeWithVerbosity is like TypeToTypeNode but with verbosity level support for expandable hover.
func (b *NodeBuilder) TypeToTypeNodeWithVerbosity(typ *Type, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker, verbosityLevel int, out *WriterContextOut) *ast.Node {
	b.enterContextEx(enclosingDeclaration, flags, internalFlags, tracker, verbosityLevel, out, 0)
	result := b.impl.typeToTypeNode(typ)
	if out != nil {
		out.CanIncreaseExpansionDepth = b.impl.ctx.out.CanIncreaseExpansionDepth
		out.Truncated = b.impl.ctx.out.Truncated
	}
	return b.exitContext(result)
}

// SignatureToSignatureDeclarationWithVerbosity is like SignatureToSignatureDeclaration but with verbosity level support.
func (b *NodeBuilder) SignatureToSignatureDeclarationWithVerbosity(signature *Signature, kind ast.Kind, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker, verbosityLevel int, out *WriterContextOut) *ast.Node {
	b.enterContextEx(enclosingDeclaration, flags, internalFlags, tracker, verbosityLevel, out, 0)
	result := b.impl.signatureToSignatureDeclarationHelper(signature, kind, nil)
	if out != nil {
		out.CanIncreaseExpansionDepth = b.impl.ctx.out.CanIncreaseExpansionDepth
		out.Truncated = b.impl.ctx.out.Truncated
	}
	return b.exitContext(result)
}

// var _ NodeBuilderInterface = NewNodeBuilderAPI(nil, nil)

func NewNodeBuilder(ch *Checker, e *printer.EmitContext) *NodeBuilder {
	return NewNodeBuilderEx(ch, e, nil /*idToSymbol*/)
}

func NewNodeBuilderEx(ch *Checker, e *printer.EmitContext, idToSymbol map[*ast.IdentifierNode]*ast.Symbol) *NodeBuilder {
	impl := newNodeBuilderImpl(ch, e, idToSymbol)
	return &NodeBuilder{impl: impl, ctxStack: make([]*NodeBuilderContext, 0, 1), host: ch.program}
}

func (c *Checker) getNodeBuilder() *NodeBuilder {
	return c.getNodeBuilderEx(nil /*idToSymbol*/)
}

func (c *Checker) getNodeBuilderEx(idToSymbol map[*ast.IdentifierNode]*ast.Symbol) *NodeBuilder {
	b := NewNodeBuilderEx(c, printer.NewEmitContext(), idToSymbol)
	return b
}
