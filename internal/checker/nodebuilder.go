package checker

import (
	"fmt"
	"slices"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/jsnum"
	"github.com/microsoft/typescript-go/internal/nodebuilder"
	"github.com/microsoft/typescript-go/internal/printer"
	"github.com/microsoft/typescript-go/internal/scanner"
)

type CompositeSymbolIdentity struct {
	isConstructorNode bool
	symbolId          ast.SymbolId
	nodeId            ast.NodeId
}

type TrackedSymbolArgs struct {
	symbol               *ast.Symbol
	enclosingDeclaration *ast.Node
	meaning              ast.SymbolFlags
}

type SerializedTypeEntry struct {
	node           *ast.Node
	truncating     bool
	addedLength    int
	trackedSymbols []*TrackedSymbolArgs
}

type CompositeTypeCacheIdentity struct {
	typeId        TypeId
	flags         nodebuilder.Flags
	internalFlags nodebuilder.InternalFlags
}

type NodeBuilderLinks struct {
	serializedTypes                  map[CompositeTypeCacheIdentity]*SerializedTypeEntry // Collection of types serialized at this location
	fakeScopeForSignatureDeclaration *string                                             // If present, this is a fake scope injected into an enclosing declaration chain.
}

type NodeBuilderContext struct {
	tracker                         nodebuilder.SymbolTracker
	approximateLength               int
	encounteredError                bool
	truncating                      bool
	reportedDiagnostic              bool
	flags                           nodebuilder.Flags
	internalFlags                   nodebuilder.InternalFlags
	depth                           int
	enclosingDeclaration            *ast.Node
	enclosingFile                   *ast.SourceFile
	inferTypeParameters             []*Type
	visitedTypes                    map[TypeId]bool
	symbolDepth                     map[CompositeSymbolIdentity]int
	trackedSymbols                  []*TrackedSymbolArgs
	mapper                          *TypeMapper
	reverseMappedStack              []*ast.Symbol
	enclosingSymbolTypes            map[ast.SymbolId]*Type
	suppressReportInferenceFallback bool

	// per signature scope state
	mustCreateTypeParameterSymbolList     bool
	mustCreateTypeParametersNamesLookups  bool
	typeParameterNames                    any
	typeParameterNamesByText              any
	typeParameterNamesByTextNextNameCount any
	typeParameterSymbolList               any
}

type NodeBuilder struct {
	// host members
	f  *ast.NodeFactory
	ch *Checker
	e  *printer.EmitContext

	// cache
	links core.LinkStore[*ast.Node, NodeBuilderLinks]

	// closures
	typeToTypeNodeClosure               func(t *Type) *ast.TypeNode
	typeReferenceToTypeNodeClosure      func(t *Type) *ast.TypeNode
	conditionalTypeToTypeNodeClosure    func(t *Type) *ast.TypeNode
	createTypeNodeFromObjectTypeClosure func(t *Type) *ast.TypeNode
	isStringNamedClosure                func(d *ast.Declaration) bool
	isSingleQuotedStringNamedClosure    func(d *ast.Declaration) bool

	// state
	ctx *NodeBuilderContext
}

const defaultMaximumTruncationLength = 160
const noTruncationMaximumTruncationLength = 1_000_000

// Node builder utility functions

// You probably don't mean to use this - use `NewNodeBuilderAPI` instead
func NewNodeBuilder(ch *Checker, e *printer.EmitContext) NodeBuilder {
	result := NodeBuilder{f: e.Factory, ch: ch, e: e, typeToTypeNodeClosure: nil, typeReferenceToTypeNodeClosure: nil, conditionalTypeToTypeNodeClosure: nil, ctx: nil}
	result.initializeClosures()
	return result
}

func (b *NodeBuilder) initializeClosures() {
	b.typeToTypeNodeClosure = b.typeToTypeNode
	b.typeReferenceToTypeNodeClosure = b.typeReferenceToTypeNode
	b.conditionalTypeToTypeNodeClosure = b.conditionalTypeToTypeNode
	b.createTypeNodeFromObjectTypeClosure = b.createTypeNodeFromObjectType
	b.isStringNamedClosure = b.isStringNamed
	b.isSingleQuotedStringNamedClosure = b.isSingleQuotedStringNamed
}

func (b *NodeBuilder) saveRestoreFlags() func() {
	flags := b.ctx.flags
	internalFlags := b.ctx.internalFlags
	depth := b.ctx.depth

	return func() {
		b.ctx.flags = flags
		b.ctx.internalFlags = internalFlags
		b.ctx.depth = depth
	}
}

func (b *NodeBuilder) checkTruncationLength() bool {
	if b.ctx.truncating {
		return b.ctx.truncating
	}
	b.ctx.truncating = b.ctx.approximateLength > (core.IfElse((b.ctx.flags&nodebuilder.FlagsNoTruncation != 0), noTruncationMaximumTruncationLength, defaultMaximumTruncationLength))
	return b.ctx.truncating
}

func (b *NodeBuilder) appendReferenceToType(root *ast.TypeNode, ref *ast.TypeNode) *ast.TypeNode {
	if ast.IsImportTypeNode(root) {
		// first shift type arguments

		// !!! In the old emitter, an Identifier could have type arguments for use with quickinfo:
		// typeArguments := root.TypeArguments
		// qualifier := root.AsImportTypeNode().Qualifier
		// if qualifier != nil {
		// 	if ast.IsIdentifier(qualifier) {
		// 		if typeArguments != getIdentifierTypeArguments(qualifier) {
		// 			qualifier = setIdentifierTypeArguments(b.f.CloneNode(qualifier), typeArguments)
		// 		}
		// 	} else {
		// 		if typeArguments != getIdentifierTypeArguments(qualifier.Right) {
		// 			qualifier = b.f.UpdateQualifiedName(qualifier, qualifier.Left, setIdentifierTypeArguments(b.f.cloneNode(qualifier.Right), typeArguments))
		// 		}
		// 	}
		// }
		// !!! Without the above, nested type args are silently elided
		imprt := root.AsImportTypeNode()
		// then move qualifiers
		ids := getAccessStack(ref)
		var qualifier *ast.Node
		for _, id := range ids {
			if qualifier != nil {
				qualifier = b.f.NewQualifiedName(qualifier, id)
			} else {
				qualifier = id
			}
		}
		return b.f.UpdateImportTypeNode(imprt, imprt.IsTypeOf, imprt.Argument, imprt.Attributes, qualifier, ref.AsTypeReferenceNode().TypeArguments)
	} else {
		// first shift type arguments
		// !!! In the old emitter, an Identifier could have type arguments for use with quickinfo:
		// typeArguments := root.TypeArguments
		// typeName := root.AsTypeReferenceNode().TypeName
		// if ast.IsIdentifier(typeName) {
		// 	if typeArguments != getIdentifierTypeArguments(typeName) {
		// 		typeName = setIdentifierTypeArguments(b.f.cloneNode(typeName), typeArguments)
		// 	}
		// } else {
		// 	if typeArguments != getIdentifierTypeArguments(typeName.Right) {
		// 		typeName = b.f.UpdateQualifiedName(typeName, typeName.Left, setIdentifierTypeArguments(b.f.cloneNode(typeName.Right), typeArguments))
		// 	}
		// }
		// !!! Without the above, nested type args are silently elided
		// then move qualifiers
		ids := getAccessStack(ref)
		var typeName *ast.Node = root.AsTypeReferenceNode().TypeName
		for _, id := range ids {
			typeName = b.f.NewQualifiedName(typeName, id)
		}
		return b.f.UpdateTypeReferenceNode(root.AsTypeReferenceNode(), typeName, ref.AsTypeReferenceNode().TypeArguments)
	}
}

func getAccessStack(ref *ast.Node) []*ast.Node {
	var state *ast.Node = ref.AsTypeReferenceNode().TypeName
	ids := []*ast.Node{}
	for !ast.IsIdentifier(state) {
		entity := state.AsQualifiedName()
		ids = append([]*ast.Node{entity.Right}, ids...)
		state = entity.Left
	}
	ids = append([]*ast.Node{state}, ids...)
	return ids
}

func isClassInstanceSide(c *Checker, t *Type) bool {
	return t.symbol != nil && t.symbol.Flags&ast.SymbolFlagsClass != 0 && (t == c.getDeclaredTypeOfClassOrInterface(t.symbol) || (t.flags&TypeFlagsObject != 0 && t.objectFlags&ObjectFlagsIsClassInstanceClone != 0))
}

func (b *NodeBuilder) createElidedInformationPlaceholder() *ast.TypeNode {
	b.ctx.approximateLength += 3
	if b.ctx.flags&nodebuilder.FlagsNoTruncation == 0 {
		return b.f.NewTypeReferenceNode(b.f.NewIdentifier("..."), nil /*typeArguments*/)
	}
	// addSyntheticLeadingComment(b.f.NewKeywordTypeNode(ast.KindAnyKeyword), ast.KindMultiLineCommentTrivia, "elided") // !!!
	return b.f.NewKeywordTypeNode(ast.KindAnyKeyword)
}

func (b *NodeBuilder) mapToTypeNodes(list []*Type) *ast.NodeList {
	if len(list) == 0 {
		return nil
	}
	contents := core.Map(list, b.typeToTypeNodeClosure)
	return b.f.NewNodeList(contents)
}

func (b *NodeBuilder) setCommentRange(node *ast.Node, range_ *ast.Node) {
	if range_ != nil && b.ctx.enclosingFile != nil && b.ctx.enclosingFile == ast.GetSourceFileOfNode(range_) {
		// Copy comments to node for declaration emit
		b.e.AssignCommentRange(node, range_)
	}
}

func (b *NodeBuilder) tryReuseExistingTypeNodeHelper(existing *ast.TypeNode) *ast.TypeNode {
	return nil // !!!
}

func (b *NodeBuilder) tryReuseExistingTypeNode(typeNode *ast.TypeNode, t *Type, host *ast.Node, addUndefined bool) *ast.TypeNode {
	originalType := t
	if addUndefined {
		t = b.ch.getOptionalType(t, !ast.IsParameter(host))
	}
	clone := b.tryReuseExistingNonParameterTypeNode(typeNode, t, host, nil)
	if clone != nil {
		// explicitly add `| undefined` if it's missing from the input type nodes and the type contains `undefined` (and not the missing type)
		if addUndefined && containsNonMissingUndefinedType(b.ch, t) && !someType(b.getTypeFromTypeNode(typeNode, false), func(t *Type) bool {
			return t.flags&TypeFlagsUndefined != 0
		}) {
			return b.f.NewUnionTypeNode(b.f.NewNodeList([]*ast.TypeNode{clone, b.f.NewKeywordTypeNode(ast.KindUndefinedKeyword)}))
		}
		return clone
	}
	if addUndefined && originalType != t {
		cloneMissingUndefined := b.tryReuseExistingNonParameterTypeNode(typeNode, originalType, host, nil)
		if cloneMissingUndefined != nil {
			return b.f.NewUnionTypeNode(b.f.NewNodeList([]*ast.TypeNode{cloneMissingUndefined, b.f.NewKeywordTypeNode(ast.KindUndefinedKeyword)}))
		}
	}
	return nil
}

func (b *NodeBuilder) typeNodeIsEquivalentToType(annotatedDeclaration *ast.Node, t *Type, typeFromTypeNode *Type) bool {
	if typeFromTypeNode == t {
		return true
	}
	if annotatedDeclaration == nil {
		return false
	}
	// !!!
	// used to be hasEffectiveQuestionToken for JSDoc
	if isOptionalDeclaration(annotatedDeclaration) {
		return b.ch.getTypeWithFacts(t, TypeFactsNEUndefined) == typeFromTypeNode
	}
	return false
}

func (b *NodeBuilder) existingTypeNodeIsNotReferenceOrIsReferenceWithCompatibleTypeArgumentCount(existing *ast.TypeNode, t *Type) bool {
	// In JS, you can say something like `Foo` and get a `Foo<any>` implicitly - we don't want to preserve that original `Foo` in these cases, though.
	if t.objectFlags&ObjectFlagsReference == 0 {
		return true
	}
	if !ast.IsTypeReferenceNode(existing) {
		return true
	}
	// `type` is a reference type, and `existing` is a type reference node, but we still need to make sure they refer to the _same_ target type
	// before we go comparing their type argument counts.
	b.ch.getTypeFromTypeReference(existing)
	// call to ensure symbol is resolved
	links := b.ch.symbolNodeLinks.TryGet(existing)
	if links == nil {
		return true
	}
	symbol := links.resolvedSymbol
	if symbol == nil {
		return true
	}
	existingTarget := b.ch.getDeclaredTypeOfSymbol(symbol)
	if existingTarget == nil || existingTarget != t.AsTypeReference().target {
		return true
	}
	return len(existing.TypeArguments()) >= b.ch.getMinTypeArgumentCount(t.AsTypeReference().target.AsInterfaceType().TypeParameters())
}

func (b *NodeBuilder) tryReuseExistingNonParameterTypeNode(existing *ast.TypeNode, t *Type, host *ast.Node, annotationType *Type) *ast.TypeNode {
	if host == nil {
		host = b.ctx.enclosingDeclaration
	}
	if annotationType == nil {
		annotationType = b.getTypeFromTypeNode(existing, true)
	}
	if annotationType != nil && b.typeNodeIsEquivalentToType(host, t, annotationType) && b.existingTypeNodeIsNotReferenceOrIsReferenceWithCompatibleTypeArgumentCount(existing, t) {
		result := b.tryReuseExistingTypeNodeHelper(existing)
		if result != nil {
			return result
		}
	}
	return nil
}

func (b *NodeBuilder) getResolvedTypeWithoutAbstractConstructSignatures(t *StructuredType) *Type {
	if len(t.ConstructSignatures()) == 0 {
		return t.AsType()
	}
	if t.objectTypeWithoutAbstractConstructSignatures != nil {
		return t.objectTypeWithoutAbstractConstructSignatures
	}
	constructSignatures := core.Filter(t.ConstructSignatures(), func(signature *Signature) bool {
		return signature.flags&SignatureFlagsAbstract == 0
	})
	if len(constructSignatures) == len(t.ConstructSignatures()) {
		t.objectTypeWithoutAbstractConstructSignatures = t.AsType()
		return t.AsType()
	}
	typeCopy := b.ch.newAnonymousType(t.symbol, t.members, t.CallSignatures(), core.IfElse(len(constructSignatures) > 0, constructSignatures, []*Signature{}), t.indexInfos)
	t.objectTypeWithoutAbstractConstructSignatures = typeCopy
	typeCopy.AsStructuredType().objectTypeWithoutAbstractConstructSignatures = typeCopy
	return typeCopy
}

func (b *NodeBuilder) symbolToNode(symbol *ast.Symbol, meaning ast.SymbolFlags) *ast.Node {
	if b.ctx.internalFlags&nodebuilder.InternalFlagsWriteComputedProps != 0 {
		if symbol.ValueDeclaration != nil {
			name := ast.GetNameOfDeclaration(symbol.ValueDeclaration)
			if name != nil && ast.IsComputedPropertyName(name) {
				return name
			}
			if b.ch.valueSymbolLinks.Has(symbol) {
				nameType := b.ch.valueSymbolLinks.Get(symbol).nameType
				if nameType != nil && nameType.flags&(TypeFlagsEnumLiteral|TypeFlagsUniqueESSymbol) != 0 {
					oldEnclosing := b.ctx.enclosingDeclaration
					b.ctx.enclosingDeclaration = nameType.symbol.ValueDeclaration
					result := b.f.NewComputedPropertyName(b.symbolToExpression(nameType.symbol, meaning))
					b.ctx.enclosingDeclaration = oldEnclosing
					return result
				}
			}
		}
	}
	return b.symbolToExpression(symbol, meaning)
}

func (b *NodeBuilder) symbolToName(symbol *ast.Symbol, meaning ast.SymbolFlags, expectsIdentifier bool) *ast.Node {
	panic("unimplemented") // !!!
}

func (b *NodeBuilder) symbolToEntityNameNode(symbol *ast.Symbol) *ast.EntityName {
	panic("unimplemented") // !!!
}

func (b *NodeBuilder) symbolToTypeNode(symbol *ast.Symbol, mask ast.SymbolFlags, typeArguments *ast.NodeList) *ast.TypeNode {
	panic("unimplemented") // !!!
}

func (b *NodeBuilder) symbolToExpression(symbol *ast.Symbol, mask ast.SymbolFlags) *ast.Expression {
	// chain := b.lookupSymbolChain(symbol, meaning)
	panic("unimplemented") // !!!
}

func (b *NodeBuilder) typeParameterToDeclarationWithConstraint(typeParameter *Type, constraintNode *ast.TypeNode) *ast.TypeParameterDeclarationNode {
	panic("unimplemented") // !!!
}

func (b *NodeBuilder) typeParameterToName(typeParameter *Type) *ast.Identifier {
	panic("unimplemented") // !!!
}

func (b *NodeBuilder) createMappedTypeNodeFromType(type_ *Type) *ast.TypeNode {
	panic("unimplemented") // !!!
}

func (b *NodeBuilder) typePredicateToTypePredicateNode(predicate *TypePredicate) *ast.Node {
	panic("unimplemented") // !!!
}

func (b *NodeBuilder) typeParameterToDeclaration(parameter *Type) *ast.Node {
	panic("unimplemented") // !!!
}

func (b *NodeBuilder) symbolToTypeParameterDeclarations(symbol *ast.Symbol) *ast.Node {
	panic("unimplemented") // !!!
}

func (b *NodeBuilder) symbolToParameterDeclaration(symbol *ast.Symbol, preserveModifierFlags bool) *ast.Node {
	panic("unimplemented") // !!!
}

func (b *NodeBuilder) symbolTableToDeclarationStatements(symbolTable *ast.SymbolTable) []*ast.Node {
	panic("unimplemented") // !!!
}

func (b *NodeBuilder) serializeTypeForExpression(expr *ast.Node) *ast.Node {
	panic("unimplemented") // !!!
}

func (b *NodeBuilder) serializeInferredReturnTypeForSignature(signature *Signature, returnType *Type) *ast.Node {
	oldSuppressReportInferenceFallback := b.ctx.suppressReportInferenceFallback
	b.ctx.suppressReportInferenceFallback = true
	typePredicate := b.ch.getTypePredicateOfSignature(signature)
	var returnTypeNode *ast.Node
	if typePredicate != nil {
		var predicate *TypePredicate
		if b.ctx.mapper != nil {
			predicate = b.ch.instantiateTypePredicate(typePredicate, b.ctx.mapper)
		} else {
			predicate = typePredicate
		}
		returnTypeNode = b.typePredicateToTypePredicateNodeHelper(predicate)
	} else {
		returnTypeNode = b.typeToTypeNodeClosure(returnType)
	}
	b.ctx.suppressReportInferenceFallback = oldSuppressReportInferenceFallback
	return returnTypeNode
}

func (b *NodeBuilder) typePredicateToTypePredicateNodeHelper(typePredicate *TypePredicate) *ast.Node {
	var assertsModifier *ast.Node
	if typePredicate.kind == TypePredicateKindAssertsThis || typePredicate.kind == TypePredicateKindAssertsIdentifier {
		assertsModifier = b.f.NewToken(ast.KindAssertsKeyword)
	} else {
		assertsModifier = nil
	}
	var parameterName *ast.Node
	if typePredicate.kind == TypePredicateKindIdentifier || typePredicate.kind == TypePredicateKindAssertsIdentifier {
		parameterName = b.f.NewIdentifier(typePredicate.parameterName)
		b.e.SetEmitFlags(parameterName, printer.EFNoAsciiEscaping)
	} else {
		parameterName = b.f.NewThisTypeNode()
	}
	var typeNode *ast.Node
	if typePredicate.t != nil {
		typeNode = b.typeToTypeNodeClosure(typePredicate.t)
	}
	return b.f.NewTypePredicateNode(assertsModifier, parameterName, typeNode)
}

type SignatureToSignatureDeclarationOptions struct {
	modifiers     []*ast.Node
	name          *ast.PropertyName
	questionToken *ast.Node
}

func (b *NodeBuilder) signatureToSignatureDeclarationHelper(signature *Signature, kind ast.Kind, options *SignatureToSignatureDeclarationOptions) *ast.Node {
	var typeParameters *[]*ast.Node
	var typeArguments *[]*ast.Node

	expandedParams := b.ch.getExpandedParameters(signature, true /*skipUnionExpanding*/)[0]
	cleanup := b.enterNewScope(signature.declaration, &expandedParams, &signature.typeParameters, &signature.parameters, signature.mapper)
	b.ctx.approximateLength += 3
	// Usually a signature contributes a few more characters than this, but 3 is the minimum

	if b.ctx.flags&nodebuilder.FlagsWriteTypeArgumentsOfSignature != 0 && signature.target != nil && signature.mapper != nil && signature.target.typeParameters != nil {
		for _, parameter := range signature.target.typeParameters {
			node := b.typeToTypeNodeClosure(b.ch.instantiateType(parameter, signature.mapper))
			if typeArguments == nil {
				typeArguments = &[]*ast.Node{}
			}
			args := append(*typeArguments, node)
			typeArguments = &args
		}
	} else if signature.typeParameters != nil {
		for _, parameter := range signature.typeParameters {
			node := b.typeParameterToDeclaration(parameter)
			if typeParameters == nil {
				typeParameters = &[]*ast.Node{}
			}
			args := append(*typeParameters, node)
			typeParameters = &args
		}
	}

	restoreFlags := b.saveRestoreFlags()
	b.ctx.flags &^= nodebuilder.FlagsSuppressAnyReturnType
	// If the expanded parameter list had a variadic in a non-trailing position, don't expand it
	parameters := core.Map(core.IfElse(core.Some(expandedParams, func(p *ast.Symbol) bool {
		return p != expandedParams[len(expandedParams)-1] && p.CheckFlags&ast.CheckFlagsRestParameter != 0
	}), signature.parameters, expandedParams), func(parameter *ast.Symbol) *ast.Node {
		return b.symbolToParameterDeclaration(parameter, kind == ast.KindConstructor)
	})
	var thisParameter *ast.Node
	if b.ctx.flags&nodebuilder.FlagsOmitThisParameter != 0 {
		thisParameter = nil
	} else {
		thisParameter = b.tryGetThisParameterDeclaration(signature)
	}
	if thisParameter != nil {
		parameters = append([]*ast.Node{thisParameter}, parameters...)
	}
	restoreFlags()

	returnTypeNode := b.serializeReturnTypeForSignature(signature)

	var modifiers []*ast.Node
	if options != nil {
		modifiers = options.modifiers
	}
	if (kind == ast.KindConstructorType) && signature.flags&SignatureFlagsAbstract != 0 {
		flags := ast.ModifiersToFlags(modifiers)
		modifiers = ast.CreateModifiersFromModifierFlags(flags|ast.ModifierFlagsAbstract, b.f.NewModifier)
	}

	paramList := b.f.NewNodeList(parameters)
	var typeParamList *ast.NodeList
	if typeParameters != nil {
		typeParamList = b.f.NewNodeList(*typeParameters)
	}
	var modifierList *ast.ModifierList
	if modifiers != nil && len(modifiers) > 0 {
		modifierList = b.f.NewModifierList(modifiers)
	}
	var name *ast.Node
	if options != nil {
		name = options.name
	}
	if name == nil {
		name = b.f.NewIdentifier("")
	}

	var node *ast.Node
	switch {
	case kind == ast.KindCallSignature:
		node = b.f.NewCallSignatureDeclaration(typeParamList, paramList, returnTypeNode)
	case kind == ast.KindConstructSignature:
		node = b.f.NewConstructSignatureDeclaration(typeParamList, paramList, returnTypeNode)
	case kind == ast.KindMethodSignature:
		var questionToken *ast.Node
		if options != nil {
			questionToken = options.questionToken
		}
		node = b.f.NewMethodSignatureDeclaration(modifierList, name, questionToken, typeParamList, paramList, returnTypeNode)
	case kind == ast.KindMethodDeclaration:
		node = b.f.NewMethodDeclaration(modifierList, nil /*asteriskToken*/, name, nil /*questionToken*/, typeParamList, paramList, returnTypeNode, nil /*body*/)
	case kind == ast.KindConstructor:
		node = b.f.NewConstructorDeclaration(modifierList, nil /*typeParamList*/, paramList, nil /*returnTypeNode*/, nil /*body*/)
	case kind == ast.KindGetAccessor:
		node = b.f.NewGetAccessorDeclaration(modifierList, name, nil /*typeParamList*/, paramList, returnTypeNode, nil /*body*/)
	case kind == ast.KindSetAccessor:
		node = b.f.NewSetAccessorDeclaration(modifierList, name, nil /*typeParamList*/, paramList, nil /*returnTypeNode*/, nil /*body*/)
	case kind == ast.KindIndexSignature:
		node = b.f.NewIndexSignatureDeclaration(modifierList, paramList, returnTypeNode)
	// !!! JSDoc Support
	// case kind == ast.KindJSDocFunctionType:
	// 	node = b.f.NewJSDocFunctionType(parameters, returnTypeNode)
	case kind == ast.KindFunctionType:
		if returnTypeNode == nil {
			returnTypeNode = b.f.NewTypeReferenceNode(b.f.NewIdentifier(""), nil)
		}
		node = b.f.NewFunctionTypeNode(typeParamList, paramList, returnTypeNode)
	case kind == ast.KindConstructorType:
		if returnTypeNode == nil {
			returnTypeNode = b.f.NewTypeReferenceNode(b.f.NewIdentifier(""), nil)
		}
		node = b.f.NewConstructorTypeNode(modifierList, typeParamList, paramList, returnTypeNode)
	case kind == ast.KindFunctionDeclaration:
		// TODO: assert name is Identifier
		node = b.f.NewFunctionDeclaration(modifierList, nil /*asteriskToken*/, name, typeParamList, paramList, returnTypeNode, nil /*body*/)
	case kind == ast.KindFunctionExpression:
		// TODO: assert name is Identifier
		node = b.f.NewFunctionExpression(modifierList, nil /*asteriskToken*/, name, typeParamList, paramList, returnTypeNode, b.f.NewBlock(b.f.NewNodeList([]*ast.Node{}), false))
	case kind == ast.KindArrowFunction:
		node = b.f.NewArrowFunction(modifierList, typeParamList, paramList, returnTypeNode, nil /*equalsGreaterThanToken*/, b.f.NewBlock(b.f.NewNodeList([]*ast.Node{}), false))
	default:
		panic("Unhandled kind in signatureToSignatureDeclarationHelper")
	}

	// !!! TODO: Smuggle type arguments of signatures out for quickinfo
	// if typeArguments != nil {
	// 	node.TypeArguments = b.f.NewNodeList(typeArguments)
	// }
	// !!! TODO: synthetic comment support
	// if signature.declaration. /* ? */ kind == ast.KindJSDocSignature && signature.declaration.Parent.Kind == ast.KindJSDocOverloadTag {
	// 	comment := getTextOfNode(signature.declaration.Parent.Parent, true /*includeTrivia*/).slice(2, -2).split(regexp.MustParse(`\r\n|\n|\r`)).map_(func(line string) string {
	// 		return line.replace(regexp.MustParse(`^\s+`), " ")
	// 	}).join("\n")
	// 	addSyntheticLeadingComment(node, ast.KindMultiLineCommentTrivia, comment, true /*hasTrailingNewLine*/)
	// }

	cleanup()
	return node
}

func (c *Checker) getExpandedParameters(sig *Signature, skipUnionExpanding bool) [][]*ast.Symbol {
	if signatureHasRestParameter(sig) {
		restIndex := len(sig.parameters) - 1
		restSymbol := sig.parameters[restIndex]
		restType := c.getTypeOfSymbol(restSymbol)
		getUniqAssociatedNamesFromTupleType := func(t *Type, restSymbol *ast.Symbol) []string {
			names := core.MapIndex(t.AsTupleType().elementInfos, func(info TupleElementInfo, i int) string {
				return c.getTupleElementLabel(info, restSymbol, i)
			})
			if len(names) > 0 {
				duplicates := []int{}
				uniqueNames := make(map[string]bool)
				for i, name := range names {
					_, ok := uniqueNames[name]
					if ok {
						duplicates = append(duplicates, i)
					} else {
						uniqueNames[name] = true
					}
				}
				counters := make(map[string]int)
				for _, i := range duplicates {
					counter, ok := counters[names[i]]
					if !ok {
						counter = 1
					}
					var name string
					for true {
						name = fmt.Sprintf("%s_%d", names[i], counter)
						_, ok := uniqueNames[name]
						if ok {
							counter++
							continue
						} else {
							uniqueNames[name] = true
							break
						}
					}
					names[i] = name
					counters[names[i]] = counter + 1
				}
			}
			return names
		}
		expandSignatureParametersWithTupleMembers := func(restType *Type, restIndex int, restSymbol *ast.Symbol) []*ast.Symbol {
			elementTypes := c.getTypeArguments(restType)
			associatedNames := getUniqAssociatedNamesFromTupleType(restType, restSymbol)
			restParams := core.MapIndex(elementTypes, func(t *Type, i int) *ast.Symbol {
				// Lookup the label from the individual tuple passed in before falling back to the signature `rest` parameter name
				// TODO: getTupleElementLabel can no longer fail, investigate if this lack of falliability meaningfully changes output
				// var name *string
				// if associatedNames != nil && associatedNames[i] != nil {
				// 	name = associatedNames[i]
				// } else {
				// 	name = c.getParameterNameAtPosition(sig, restIndex+i, restType)
				// }
				name := associatedNames[i]
				flags := restType.AsTupleType().elementInfos[i].flags
				var checkFlags ast.CheckFlags
				switch {
				case flags&ElementFlagsVariable != 0:
					checkFlags = ast.CheckFlagsRestParameter
				case flags&ElementFlagsOptional != 0:
					checkFlags = ast.CheckFlagsOptionalParameter
				default:
					checkFlags = 0
				}
				symbol := c.newSymbolEx(ast.SymbolFlagsFunctionScopedVariable, name, checkFlags)
				links := c.valueSymbolLinks.Get(symbol)
				if flags&ElementFlagsRest != 0 {
					links.resolvedType = c.createArrayType(t)
				} else {
					links.resolvedType = t
				}
				return symbol
			})
			return core.Concatenate(sig.parameters[0:restIndex], restParams)
		}

		if isTupleType(restType) {
			return [][]*ast.Symbol{expandSignatureParametersWithTupleMembers(restType, restIndex, restSymbol)}
		} else if !skipUnionExpanding && restType.flags&TypeFlagsUnion != 0 && core.Every(restType.AsUnionType().types, isTupleType) {
			return core.Map(restType.AsUnionType().types, func(t *Type) []*ast.Symbol {
				return expandSignatureParametersWithTupleMembers(t, restIndex, restSymbol)
			})
		}
	}
	return [][]*ast.Symbol{sig.parameters}

}

func (b *NodeBuilder) tryGetThisParameterDeclaration(signature *Signature) *ast.Node {
	if signature.thisParameter != nil {
		return b.symbolToParameterDeclaration(signature.thisParameter, false)
	}
	if signature.declaration != nil && ast.IsInJSFile(signature.declaration) {
		// !!! JSDoc Support
		// thisTag := getJSDocThisTag(signature.declaration)
		// if (thisTag && thisTag.typeExpression) {
		// 	return factory.createParameterDeclaration(
		// 		/*modifiers*/ undefined,
		// 		/*dotDotDotToken*/ undefined,
		// 		"this",
		// 		/*questionToken*/ undefined,
		// 		typeToTypeNodeHelper(getTypeFromTypeNode(context, thisTag.typeExpression), context),
		// 	);
		// }
	}
	return nil
}

/**
* Serializes the return type of the signature by first trying to use the syntactic printer if possible and falling back to the checker type if not.
 */
func (b *NodeBuilder) serializeReturnTypeForSignature(signature *Signature) *ast.Node {
	suppressAny := b.ctx.flags&nodebuilder.FlagsSuppressAnyReturnType != 0
	restoreFlags := b.saveRestoreFlags()
	if suppressAny {
		b.ctx.flags &= ^nodebuilder.FlagsSuppressAnyReturnType // suppress only toplevel `any`s
	}
	var returnTypeNode *ast.Node

	returnType := b.ch.getReturnTypeOfSignature(signature)
	if !(suppressAny && IsTypeAny(returnType)) {
		// !!! IsolatedDeclaration support
		// if signature.declaration != nil && !ast.NodeIsSynthesized(signature.declaration) {
		// 	declarationSymbol := b.ch.getSymbolOfDeclaration(signature.declaration)
		// 	restore := addSymbolTypeToContext(declarationSymbol, returnType)
		// 	returnTypeNode = syntacticNodeBuilder.serializeReturnTypeForSignature(signature.declaration, declarationSymbol)
		// 	restore()
		// }
		if returnTypeNode == nil {
			returnTypeNode = b.serializeInferredReturnTypeForSignature(signature, returnType)
		}
	}

	if returnTypeNode == nil && !suppressAny {
		returnTypeNode = b.f.NewKeywordTypeNode(ast.KindAnyKeyword)
	}
	restoreFlags()
	return returnTypeNode
}

func (b *NodeBuilder) indexInfoToIndexSignatureDeclarationHelper(indexInfo *IndexInfo, typeNode *ast.TypeNode) *ast.Node {
	name := getNameFromIndexInfo(indexInfo)
	indexerTypeNode := b.typeToTypeNodeClosure(indexInfo.keyType)

	indexingParameter := b.f.NewParameterDeclaration(nil, nil, b.f.NewIdentifier(name), nil, indexerTypeNode, nil)
	if typeNode == nil {
		if indexInfo.valueType == nil {
			typeNode = b.f.NewKeywordTypeNode(ast.KindAnyKeyword)
		} else {
			typeNode = b.typeToTypeNodeClosure(indexInfo.valueType)
		}
	}
	if indexInfo.valueType == nil && b.ctx.flags&nodebuilder.FlagsAllowEmptyIndexInfoType == 0 {
		b.ctx.encounteredError = true
	}
	b.ctx.approximateLength += len(name) + 4
	var modifiers *ast.ModifierList
	if indexInfo.isReadonly {
		b.ctx.approximateLength += 9
		modifiers = b.f.NewModifierList([]*ast.Node{b.f.NewModifier(ast.KindReadonlyKeyword)})
	}
	return b.f.NewIndexSignatureDeclaration(modifiers, b.f.NewNodeList([]*ast.Node{indexingParameter}), typeNode)
}

/**
* Unlike `typeToTypeNodeHelper`, this handles setting up the `AllowUniqueESSymbolType` flag
* so a `unique symbol` is returned when appropriate for the input symbol, rather than `typeof sym`
* @param declaration - The preferred declaration to pull existing type nodes from (the symbol will be used as a fallback to find any annotated declaration)
* @param type - The type to write; an existing annotation must match this type if it's used, otherwise this is the type serialized as a new type node
* @param symbol - The symbol is used both to find an existing annotation if declaration is not provided, and to determine if `unique symbol` should be printed
 */
func (b *NodeBuilder) serializeTypeForDeclaration(declaration *ast.Declaration, t *Type, symbol *ast.Symbol) *ast.Node {
	// !!! node reuse logic
	restoreFlags := b.saveRestoreFlags()
	if t.flags&TypeFlagsUniqueESSymbol != 0 && t.symbol == symbol && (b.ctx.enclosingDeclaration == nil || core.Some(symbol.Declarations, func(d *ast.Declaration) bool {
		return ast.GetSourceFileOfNode(d) == b.ctx.enclosingFile
	})) {
		b.ctx.flags |= nodebuilder.FlagsAllowUniqueESSymbolType
	}
	result := b.typeToTypeNodeClosure(t) // !!! expressionOrTypeToTypeNode
	restoreFlags()
	return result
}

const MAX_REVERSE_MAPPED_NESTING_INSPECTION_DEPTH = 3

func (b *NodeBuilder) shouldUsePlaceholderForProperty(propertySymbol *ast.Symbol) bool {
	// Use placeholders for reverse mapped types we've either
	// (1) already descended into, or
	// (2) are nested reverse mappings within a mapping over a non-anonymous type, or
	// (3) are deeply nested properties that originate from the same mapped type.
	// Condition (2) is a restriction mostly just to
	// reduce the blowup in printback size from doing, eg, a deep reverse mapping over `Window`.
	// Since anonymous types usually come from expressions, this allows us to preserve the output
	// for deep mappings which likely come from expressions, while truncating those parts which
	// come from mappings over library functions.
	// Condition (3) limits printing of possibly infinitely deep reverse mapped types.
	if propertySymbol.CheckFlags&ast.CheckFlagsReverseMapped == 0 {
		return false
	}
	// (1)
	for _, elem := range b.ctx.reverseMappedStack {
		if elem == propertySymbol {
			return true
		}
	}
	// (2)
	if len(b.ctx.reverseMappedStack) > 0 {
		last := b.ctx.reverseMappedStack[len(b.ctx.reverseMappedStack)-1]
		if b.ch.ReverseMappedSymbolLinks.Has(last) {
			links := b.ch.ReverseMappedSymbolLinks.TryGet(last)
			propertyType := links.propertyType
			if propertyType != nil && propertyType.objectFlags&ObjectFlagsAnonymous == 0 {
				return true
			}
		}
	}
	// (3) - we only inspect the last MAX_REVERSE_MAPPED_NESTING_INSPECTION_DEPTH elements of the
	// stack for approximate matches to catch tight infinite loops
	// TODO: Why? Reasoning lost to time. this could probably stand to be improved?
	if len(b.ctx.reverseMappedStack) < MAX_REVERSE_MAPPED_NESTING_INSPECTION_DEPTH {
		return false
	}
	if !b.ch.ReverseMappedSymbolLinks.Has(propertySymbol) {
		return false
	}
	propertyLinks := b.ch.ReverseMappedSymbolLinks.TryGet(propertySymbol)
	propMappedType := propertyLinks.mappedType
	if propMappedType == nil || propMappedType.symbol == nil {
		return false
	}
	for i := range b.ctx.reverseMappedStack {
		if i > MAX_REVERSE_MAPPED_NESTING_INSPECTION_DEPTH {
			break
		}
		prop := b.ctx.reverseMappedStack[len(b.ctx.reverseMappedStack)-1-i]
		if b.ch.ReverseMappedSymbolLinks.Has(prop) {
			links := b.ch.ReverseMappedSymbolLinks.TryGet(prop)
			mappedType := links.mappedType
			if mappedType != nil && mappedType.symbol == propMappedType.symbol {
				return true
			}
		}
	}
	return false
}

func (b *NodeBuilder) trackComputedName(accessExpression *ast.Node, enclosingDeclaration *ast.Node) {
	// get symbol of the first identifier of the entityName
	firstIdentifier := ast.GetFirstIdentifier(accessExpression)
	name := b.ch.resolveName(firstIdentifier, firstIdentifier.Text(), ast.SymbolFlagsValue|ast.SymbolFlagsExportValue, nil /*nameNotFoundMessage*/, true /*isUse*/, false)
	if name != nil {
		b.ctx.tracker.TrackSymbol(name, enclosingDeclaration, ast.SymbolFlagsValue)
	}
}

func (b *NodeBuilder) createPropertyNameNodeForIdentifierOrLiteral(name string, target core.ScriptTarget, _singleQuote bool, stringNamed bool, isMethod bool) *ast.Node {
	isMethodNamedNew := isMethod && name == "new"
	if !isMethodNamedNew && scanner.IsIdentifierText(name, target) {
		return b.f.NewIdentifier(name)
	}
	if !stringNamed && !isMethodNamedNew && isNumericLiteralName(name) && jsnum.FromString(name) >= 0 {
		return b.f.NewNumericLiteral(name)
	}
	result := b.f.NewStringLiteral(name)
	// !!! TODO: set singleQuote
	return result
}

func (b *NodeBuilder) isStringNamed(d *ast.Declaration) bool {
	name := ast.GetNameOfDeclaration(d)
	if name == nil {
		return false
	}
	if ast.IsComputedPropertyName(name) {
		t := b.ch.checkExpression(name.AsComputedPropertyName().Expression)
		return t.flags&TypeFlagsStringLike != 0
	}
	if ast.IsElementAccessExpression(name) {
		t := b.ch.checkExpression(name.AsElementAccessExpression().ArgumentExpression)
		return t.flags&TypeFlagsStringLike != 0
	}
	return ast.IsStringLiteral(name)
}

func (b *NodeBuilder) isSingleQuotedStringNamed(d *ast.Declaration) bool {
	return false // !!!
	// TODO: actually support single-quote-style-maintenance
	// name := ast.GetNameOfDeclaration(d)
	// return name != nil && ast.IsStringLiteral(name) && (name.AsStringLiteral().SingleQuote || !nodeIsSynthesized(name) && startsWith(getTextOfNode(name, false /*includeTrivia*/), "'"))
}

func (b *NodeBuilder) getPropertyNameNodeForSymbol(symbol *ast.Symbol) *ast.Node {
	stringNamed := len(symbol.Declarations) != 0 && core.Every(symbol.Declarations, b.isStringNamedClosure)
	singleQuote := len(symbol.Declarations) != 0 && core.Every(symbol.Declarations, b.isSingleQuotedStringNamedClosure)
	isMethod := symbol.Flags&ast.SymbolFlagsMethod != 0
	fromNameType := b.getPropertyNameNodeForSymbolFromNameType(symbol, singleQuote, stringNamed, isMethod)
	if fromNameType != nil {
		return fromNameType
	}
	return b.createPropertyNameNodeForIdentifierOrLiteral(symbol.Name, b.ch.compilerOptions.GetEmitScriptTarget(), singleQuote, stringNamed, isMethod)
}

// See getNameForSymbolFromNameType for a stringy equivalent
func (b *NodeBuilder) getPropertyNameNodeForSymbolFromNameType(symbol *ast.Symbol, singleQuote bool, stringNamed bool, isMethod bool) *ast.Node {
	if !b.ch.valueSymbolLinks.Has(symbol) {
		return nil
	}
	nameType := b.ch.valueSymbolLinks.TryGet(symbol).nameType
	if nameType == nil {
		return nil
	}
	if nameType.flags&TypeFlagsStringOrNumberLiteral != 0 {
		name := nameType.AsLiteralType().value.(string)
		if !scanner.IsIdentifierText(name, b.ch.compilerOptions.GetEmitScriptTarget()) && (stringNamed || !isNumericLiteralName(name)) {
			// !!! TODO: set singleQuote
			return b.f.NewStringLiteral(name)
		}
		if isNumericLiteralName(name) && name[0] == '-' {
			return b.f.NewComputedPropertyName(b.f.NewPrefixUnaryExpression(ast.KindMinusToken, b.f.NewNumericLiteral(name[1:])))
		}
		return b.createPropertyNameNodeForIdentifierOrLiteral(name, b.ch.compilerOptions.GetEmitScriptTarget(), singleQuote, stringNamed, isMethod)
	}
	if nameType.flags&TypeFlagsUniqueESSymbol != 0 {
		return b.f.NewComputedPropertyName(b.symbolToExpression(nameType.AsUniqueESSymbolType().symbol, ast.SymbolFlagsValue))
	}
	return nil
}

func (b *NodeBuilder) addPropertyToElementList(propertySymbol *ast.Symbol, typeElements []*ast.TypeElement) []*ast.TypeElement {
	propertyIsReverseMapped := propertySymbol.CheckFlags&ast.CheckFlagsReverseMapped != 0
	var propertyType *Type
	if b.shouldUsePlaceholderForProperty(propertySymbol) {
		propertyType = b.ch.anyType
	} else {
		propertyType = b.ch.getNonMissingTypeOfSymbol(propertySymbol)
	}
	saveEnclosingDeclaration := b.ctx.enclosingDeclaration
	b.ctx.enclosingDeclaration = nil
	if isLateBoundName(propertySymbol.Name) {
		if len(propertySymbol.Declarations) > 0 {
			decl := propertySymbol.Declarations[0]
			if b.ch.hasLateBindableName(decl) {
				if ast.IsBinaryExpression(decl) {
					name := ast.GetNameOfDeclaration(decl)
					if name != nil && ast.IsElementAccessExpression(name) && ast.IsPropertyAccessEntityNameExpression(name.AsElementAccessExpression().ArgumentExpression) {
						b.trackComputedName(name.AsElementAccessExpression().ArgumentExpression, saveEnclosingDeclaration)
					}
				} else {
					b.trackComputedName(decl.Name().Expression(), saveEnclosingDeclaration)
				}
			}
		} else {
			b.ctx.tracker.ReportNonSerializableProperty(b.ch.symbolToString(propertySymbol))
		}
	}
	if propertySymbol.ValueDeclaration != nil {
		b.ctx.enclosingDeclaration = propertySymbol.ValueDeclaration
	} else if len(propertySymbol.Declarations) > 0 && propertySymbol.Declarations[0] != nil {
		b.ctx.enclosingDeclaration = propertySymbol.Declarations[0]
	} else {
		b.ctx.enclosingDeclaration = saveEnclosingDeclaration
	}
	propertyName := b.getPropertyNameNodeForSymbol(propertySymbol)
	b.ctx.enclosingDeclaration = saveEnclosingDeclaration
	b.ctx.approximateLength += len(ast.SymbolName(propertySymbol)) + 1

	if propertySymbol.Flags&ast.SymbolFlagsAccessor != 0 {
		writeType := b.ch.getWriteTypeOfSymbol(propertySymbol)
		if propertyType != writeType && !b.ch.isErrorType(propertyType) && !b.ch.isErrorType(writeType) {
			getterDeclaration := ast.GetDeclarationOfKind(propertySymbol, ast.KindGetAccessor)
			getterSignature := b.ch.getSignatureFromDeclaration(getterDeclaration)
			getter := b.signatureToSignatureDeclarationHelper(getterSignature, ast.KindGetAccessor, &SignatureToSignatureDeclarationOptions{
				name: propertyName,
			})
			b.setCommentRange(getter, getterDeclaration)
			typeElements = append(typeElements, getter)
			setterDeclaration := ast.GetDeclarationOfKind(propertySymbol, ast.KindSetAccessor)
			setterSignature := b.ch.getSignatureFromDeclaration(setterDeclaration)
			setter := b.signatureToSignatureDeclarationHelper(setterSignature, ast.KindSetAccessor, &SignatureToSignatureDeclarationOptions{
				name: propertyName,
			})
			b.setCommentRange(setter, setterDeclaration)
			typeElements = append(typeElements, setter)
			return typeElements
		}
	}

	var optionalToken *ast.Node
	if propertySymbol.Flags&ast.SymbolFlagsOptional != 0 {
		optionalToken = b.f.NewToken(ast.KindQuestionToken)
	} else {
		optionalToken = nil
	}
	if propertySymbol.Flags&(ast.SymbolFlagsFunction|ast.SymbolFlagsMethod) != 0 && len(b.ch.getPropertiesOfObjectType(propertyType)) == 0 && !b.ch.isReadonlySymbol(propertySymbol) {
		signatures := b.ch.getSignaturesOfType(b.ch.filterType(propertyType, func(t *Type) bool {
			return t.flags&TypeFlagsUndefined == 0
		}), SignatureKindCall)
		for _, signature := range signatures {
			methodDeclaration := b.signatureToSignatureDeclarationHelper(signature, ast.KindMethodSignature, &SignatureToSignatureDeclarationOptions{
				name:          propertyName,
				questionToken: optionalToken,
			})
			b.setCommentRange(methodDeclaration, propertySymbol.ValueDeclaration) // !!! missing JSDoc support formerly provided by preserveCommentsOn
			typeElements = append(typeElements, methodDeclaration)
		}
		if len(signatures) != 0 || optionalToken == nil {
			return typeElements
		}
	}
	var propertyTypeNode *ast.TypeNode
	if b.shouldUsePlaceholderForProperty(propertySymbol) {
		propertyTypeNode = b.createElidedInformationPlaceholder()
	} else {
		if propertyIsReverseMapped {
			b.ctx.reverseMappedStack = append(b.ctx.reverseMappedStack, propertySymbol)
		}
		if propertyType != nil {
			propertyTypeNode = b.serializeTypeForDeclaration(nil /*declaration*/, propertyType, propertySymbol)
		} else {
			propertyTypeNode = b.f.NewKeywordTypeNode(ast.KindAnyKeyword)
		}
		if propertyIsReverseMapped {
			b.ctx.reverseMappedStack = b.ctx.reverseMappedStack[:len(b.ctx.reverseMappedStack)-1]
		}
	}

	var modifiers *ast.ModifierList
	if b.ch.isReadonlySymbol(propertySymbol) {
		modifiers = b.f.NewModifierList([]*ast.Node{b.f.NewModifier(ast.KindReadonlyKeyword)})
		b.ctx.approximateLength += 9
	}
	propertySignature := b.f.NewPropertySignatureDeclaration(modifiers, propertyName, optionalToken, propertyTypeNode, nil)

	b.setCommentRange(propertySignature, propertySymbol.ValueDeclaration) // !!! missing JSDoc support formerly provided by preserveCommentsOn
	typeElements = append(typeElements, propertySignature)

	return typeElements
}

func (b *NodeBuilder) createTypeNodesFromResolvedType(resolvedType *StructuredType) *ast.NodeList {
	if b.checkTruncationLength() {
		if b.ctx.flags&nodebuilder.FlagsNoTruncation != 0 {
			panic("NotEmittedTypeElement not implemented") // !!!
			// elem := b.f.NewNotEmittedTypeElement()
			// b.e.addSyntheticTrailingComment(elem, ast.KindMultiLineCommentTrivia, "elided")
			// return b.f.NewNodeList([]*ast.TypeElement{elem})
		}
		return b.f.NewNodeList([]*ast.Node{b.f.NewPropertySignatureDeclaration(nil, b.f.NewIdentifier("..."), nil, nil, nil)})
	}
	var typeElements []*ast.TypeElement
	for _, signature := range resolvedType.CallSignatures() {
		typeElements = append(typeElements, b.signatureToSignatureDeclarationHelper(signature, ast.KindCallSignature, nil))
	}
	for _, signature := range resolvedType.ConstructSignatures() {
		if signature.flags&SignatureFlagsAbstract != 0 {
			continue
		}
		typeElements = append(typeElements, b.signatureToSignatureDeclarationHelper(signature, ast.KindConstructSignature, nil))
	}
	for _, info := range resolvedType.indexInfos {
		typeElements = append(typeElements, b.indexInfoToIndexSignatureDeclarationHelper(info, core.IfElse(resolvedType.objectFlags&ObjectFlagsReverseMapped != 0, b.createElidedInformationPlaceholder(), nil)))
	}

	properties := resolvedType.properties
	if len(properties) == 0 {
		return b.f.NewNodeList(typeElements)
	}

	i := 0
	for _, propertySymbol := range properties {
		i++
		if b.ctx.flags&nodebuilder.FlagsWriteClassExpressionAsTypeLiteral != 0 {
			if propertySymbol.Flags&ast.SymbolFlagsPrototype != 0 {
				continue
			}
			if getDeclarationModifierFlagsFromSymbol(propertySymbol)&(ast.ModifierFlagsPrivate|ast.ModifierFlagsProtected) != 0 {
				b.ctx.tracker.ReportPrivateInBaseOfClassExpression(propertySymbol.Name)
			}
		}
		if b.checkTruncationLength() && (i+2 < len(properties)-1) {
			if b.ctx.flags&nodebuilder.FlagsNoTruncation != 0 {
				// !!! synthetic comment support - missing middle silently elided without
				// typeElement := typeElements[len(typeElements) - 1].Clone()
				// typeElements = typeElements[0:len(typeElements)-1]
				// b.e.addSyntheticTrailingComment(typeElement, ast.KindMultiLineCommentTrivia, __TEMPLATE__("... ", properties.length-i, " more elided ..."))
				// typeElements = append(typeElements, typeElement)
			} else {
				text := fmt.Sprintf("... %d more ...", len(properties)-i)
				typeElements = append(typeElements, b.f.NewPropertySignatureDeclaration(nil, b.f.NewIdentifier(text), nil, nil, nil))
			}
			typeElements = b.addPropertyToElementList(properties[len(properties)-1], typeElements)
			break
		}
		typeElements = b.addPropertyToElementList(propertySymbol, typeElements)
	}
	if len(typeElements) != 0 {
		return b.f.NewNodeList(typeElements)
	} else {
		return nil
	}
}

func (b *NodeBuilder) createTypeNodeFromObjectType(t *Type) *ast.TypeNode {
	if b.ch.isGenericMappedType(t) || (t.objectFlags&ObjectFlagsMapped != 0 && t.AsMappedType().containsError) {
		return b.createMappedTypeNodeFromType(t)
	}

	resolved := b.ch.resolveStructuredTypeMembers(t)
	callSigs := resolved.CallSignatures()
	ctorSigs := resolved.ConstructSignatures()
	if len(resolved.properties) == 0 && len(resolved.indexInfos) == 0 {
		if len(callSigs) == 0 && len(ctorSigs) == 0 {
			b.ctx.approximateLength += 2
			result := b.f.NewTypeLiteralNode(b.f.NewNodeList([]*ast.Node{}))
			b.e.SetEmitFlags(result, printer.EFSingleLine)
			return result
		}

		if len(callSigs) == 1 && len(ctorSigs) == 0 {
			signature := callSigs[0]
			signatureNode := b.signatureToSignatureDeclarationHelper(signature, ast.KindFunctionType, nil)
			return signatureNode
		}

		if len(ctorSigs) == 1 && len(callSigs) == 0 {
			signature := ctorSigs[0]
			signatureNode := b.signatureToSignatureDeclarationHelper(signature, ast.KindConstructorType, nil)
			return signatureNode
		}
	}

	abstractSignatures := core.Filter(ctorSigs, func(signature *Signature) bool {
		return signature.flags&SignatureFlagsAbstract != 0
	})
	if len(callSigs) > 0 {
		types := core.Map(abstractSignatures, func(s *Signature) *Type {
			return b.ch.getOrCreateTypeFromSignature(s, nil)
		})
		// count the number of type elements excluding abstract constructors
		typeElementCount := len(callSigs) + (len(ctorSigs) - len(abstractSignatures)) + len(resolved.indexInfos) + (core.IfElse(b.ctx.flags&nodebuilder.FlagsWriteClassExpressionAsTypeLiteral != 0, core.CountWhere(resolved.properties, func(p *ast.Symbol) bool {
			return p.Flags&ast.SymbolFlagsPrototype == 0
		}), len(resolved.properties)))
		// don't include an empty object literal if there were no other static-side
		// properties to write, i.e. `abstract class C { }` becomes `abstract new () => {}`
		// and not `(abstract new () => {}) & {}`
		if typeElementCount != 0 {
			// create a copy of the object type without any abstract construct signatures.
			types = append(types, b.getResolvedTypeWithoutAbstractConstructSignatures(resolved))
		}
		return b.typeToTypeNodeClosure(b.ch.getIntersectionType(types))
	}

	restoreFlags := b.saveRestoreFlags()
	b.ctx.flags |= nodebuilder.FlagsInObjectTypeLiteral
	members := b.createTypeNodesFromResolvedType(resolved)
	restoreFlags()
	typeLiteralNode := b.f.NewTypeLiteralNode(members)
	b.ctx.approximateLength += 2
	b.e.SetEmitFlags(typeLiteralNode, core.IfElse((b.ctx.flags&nodebuilder.FlagsMultilineObjectLiterals != 0), 0, printer.EFSingleLine))
	return typeLiteralNode
}

func getTypeAliasForTypeLiteral(c *Checker, t *Type) *ast.Symbol {
	if t.symbol != nil && t.symbol.Flags&ast.SymbolFlagsTypeLiteral != 0 && t.symbol.Declarations != nil {
		node := ast.WalkUpParenthesizedTypes(t.symbol.Declarations[0].Parent)
		if ast.IsTypeAliasDeclaration(node) {
			return c.getSymbolOfDeclaration(node)
		}
	}
	return nil
}

func (b *NodeBuilder) shouldWriteTypeOfFunctionSymbol(symbol *ast.Symbol, typeId TypeId) bool {
	isStaticMethodSymbol := symbol.Flags&ast.SymbolFlagsMethod != 0 && core.Some(symbol.Declarations, func(declaration *ast.Node) bool {
		return ast.IsStatic(declaration)
	})
	isNonLocalFunctionSymbol := false
	if symbol.Flags&ast.SymbolFlagsFunction != 0 {
		if symbol.Parent != nil {
			isNonLocalFunctionSymbol = true
		} else {
			for _, declaration := range symbol.Declarations {
				if declaration.Parent.Kind == ast.KindSourceFile || declaration.Parent.Kind == ast.KindModuleBlock {
					isNonLocalFunctionSymbol = true
					break
				}
			}
		}
	}
	if isStaticMethodSymbol || isNonLocalFunctionSymbol {
		// typeof is allowed only for static/non local functions
		_, visited := b.ctx.visitedTypes[typeId]
		return (b.ctx.flags&nodebuilder.FlagsUseTypeOfFunction != 0 || visited) && (b.ctx.flags&nodebuilder.FlagsUseStructuralFallback == 0 || b.ch.IsValueSymbolAccessible(symbol, b.ctx.enclosingDeclaration))
		// And the build is going to succeed without visibility error or there is no structural fallback allowed
	}
	return false
}

func (b *NodeBuilder) createAnonymousTypeNode(t *Type) *ast.TypeNode {
	typeId := t.id
	symbol := t.symbol
	if symbol != nil {
		isInstantiationExpressionType := t.objectFlags&ObjectFlagsInstantiationExpressionType != 0
		if isInstantiationExpressionType {
			instantiationExpressionType := t.AsInstantiationExpressionType()
			existing := instantiationExpressionType.node
			if ast.IsTypeQueryNode(existing) {
				typeNode := b.tryReuseExistingNonParameterTypeNode(existing, t, nil, nil)
				if typeNode != nil {
					return typeNode
				}
			}
			if _, ok := b.ctx.visitedTypes[typeId]; ok {
				return b.createElidedInformationPlaceholder()
			}
			return b.visitAndTransformType(t, b.createTypeNodeFromObjectTypeClosure)
		}
		var isInstanceType ast.SymbolFlags
		if isClassInstanceSide(b.ch, t) {
			isInstanceType = ast.SymbolFlagsType
		} else {
			isInstanceType = ast.SymbolFlagsValue
		}

		// !!! JS support
		// if c.isJSConstructor(symbol.ValueDeclaration) {
		// 	// Instance and static types share the same symbol; only add 'typeof' for the static side.
		// 	return b.symbolToTypeNode(symbol, isInstanceType, nil)
		// } else
		if symbol.Flags&ast.SymbolFlagsClass != 0 && b.ch.getBaseTypeVariableOfClass(symbol) == nil && !(symbol.ValueDeclaration != nil && ast.IsClassLike(symbol.ValueDeclaration) && b.ctx.flags&nodebuilder.FlagsWriteClassExpressionAsTypeLiteral != 0 && (!ast.IsClassDeclaration(symbol.ValueDeclaration) || b.ch.IsSymbolAccessible(symbol, b.ctx.enclosingDeclaration, isInstanceType, false /*shouldComputeAliasesToMakeVisible*/).Accessibility != printer.SymbolAccessibilityAccessible)) || symbol.Flags&(ast.SymbolFlagsEnum|ast.SymbolFlagsValueModule) != 0 || b.shouldWriteTypeOfFunctionSymbol(symbol, typeId) {
			return b.symbolToTypeNode(symbol, isInstanceType, nil)
		} else if _, ok := b.ctx.visitedTypes[typeId]; ok {
			// If type is an anonymous type literal in a type alias declaration, use type alias name
			typeAlias := getTypeAliasForTypeLiteral(b.ch, t)
			if typeAlias != nil {
				// The specified symbol flags need to be reinterpreted as type flags
				return b.symbolToTypeNode(typeAlias, ast.SymbolFlagsType, nil)
			} else {
				return b.createElidedInformationPlaceholder()
			}
		} else {
			return b.visitAndTransformType(t, b.createTypeNodeFromObjectTypeClosure)
		}
	} else {
		// Anonymous types without a symbol are never circular.
		return b.createTypeNodeFromObjectTypeClosure(t)
	}

}

func (b *NodeBuilder) getTypeFromTypeNode(node *ast.TypeNode, noMappedTypes bool) *Type {
	// !!! noMappedTypes optional param support
	t := b.ch.getTypeFromTypeNode(node)
	if b.ctx.mapper == nil {
		return t
	}

	instantiated := b.ch.instantiateType(t, b.ctx.mapper)
	if noMappedTypes && instantiated != t {
		return nil
	}
	return instantiated
}

func (b *NodeBuilder) typeToTypeNodeOrCircularityElision(t *Type) *ast.TypeNode {
	if t.flags&TypeFlagsUnion != 0 {
		if _, ok := b.ctx.visitedTypes[t.id]; ok {
			if b.ctx.flags&nodebuilder.FlagsAllowAnonymousIdentifier == 0 {
				b.ctx.encounteredError = true
				b.ctx.tracker.ReportCyclicStructureError()
			}
			return b.createElidedInformationPlaceholder()
		}
		return b.visitAndTransformType(t, b.typeToTypeNodeClosure)
	}
	return b.typeToTypeNodeClosure(t)
}

func (b *NodeBuilder) conditionalTypeToTypeNode(_t *Type) *ast.TypeNode {
	t := _t.AsConditionalType()
	checkTypeNode := b.typeToTypeNodeClosure(t.checkType)
	b.ctx.approximateLength += 15
	if b.ctx.flags&nodebuilder.FlagsGenerateNamesForShadowedTypeParams != 0 && t.root.isDistributive && t.checkType.flags&TypeFlagsTypeParameter == 0 {
		newParam := b.ch.newTypeParameter(b.ch.newSymbol(ast.SymbolFlagsTypeParameter, "T" /* as __String */))
		name := b.typeParameterToName(newParam)
		newTypeVariable := b.f.NewTypeReferenceNode(name.AsNode(), nil)
		b.ctx.approximateLength += 37
		// 15 each for two added conditionals, 7 for an added infer type
		newMapper := prependTypeMapping(t.root.checkType, newParam, t.mapper)
		saveInferTypeParameters := b.ctx.inferTypeParameters
		b.ctx.inferTypeParameters = t.root.inferTypeParameters
		extendsTypeNode := b.typeToTypeNodeClosure(b.ch.instantiateType(t.root.extendsType, newMapper))
		b.ctx.inferTypeParameters = saveInferTypeParameters
		trueTypeNode := b.typeToTypeNodeOrCircularityElision(b.ch.instantiateType(b.getTypeFromTypeNode(t.root.node.TrueType, false), newMapper))
		falseTypeNode := b.typeToTypeNodeOrCircularityElision(b.ch.instantiateType(b.getTypeFromTypeNode(t.root.node.FalseType, false), newMapper))

		// outermost conditional makes `T` a type parameter, allowing the inner conditionals to be distributive
		// second conditional makes `T` have `T & checkType` substitution, so it is correctly usable as the checkType
		// inner conditional runs the check the user provided on the check type (distributively) and returns the result
		// checkType extends infer T ? T extends checkType ? T extends extendsType<T> ? trueType<T> : falseType<T> : never : never;
		// this is potentially simplifiable to
		// checkType extends infer T ? T extends checkType & extendsType<T> ? trueType<T> : falseType<T> : never;
		// but that may confuse users who read the output more.
		// On the other hand,
		// checkType extends infer T extends checkType ? T extends extendsType<T> ? trueType<T> : falseType<T> : never;
		// may also work with `infer ... extends ...` in, but would produce declarations only compatible with the latest TS.
		newId := newTypeVariable.AsTypeReferenceNode().TypeName.AsIdentifier().Clone(b.f)
		syntheticExtendsNode := b.f.NewInferTypeNode(b.f.NewTypeParameterDeclaration(nil, newId, nil, nil))
		innerCheckConditionalNode := b.f.NewConditionalTypeNode(newTypeVariable, extendsTypeNode, trueTypeNode, falseTypeNode)
		syntheticTrueNode := b.f.NewConditionalTypeNode(b.f.NewTypeReferenceNode(name.Clone(b.f), nil), b.f.DeepCloneNode(checkTypeNode), innerCheckConditionalNode, b.f.NewKeywordTypeNode(ast.KindNeverKeyword))
		return b.f.NewConditionalTypeNode(checkTypeNode, syntheticExtendsNode, syntheticTrueNode, b.f.NewKeywordTypeNode(ast.KindNeverKeyword))
	}
	saveInferTypeParameters := b.ctx.inferTypeParameters
	b.ctx.inferTypeParameters = t.root.inferTypeParameters
	extendsTypeNode := b.typeToTypeNodeClosure(t.extendsType)
	b.ctx.inferTypeParameters = saveInferTypeParameters
	trueTypeNode := b.typeToTypeNodeOrCircularityElision(b.ch.getTrueTypeFromConditionalType(_t))
	falseTypeNode := b.typeToTypeNodeOrCircularityElision(b.ch.getFalseTypeFromConditionalType(_t))
	return b.f.NewConditionalTypeNode(checkTypeNode, extendsTypeNode, trueTypeNode, falseTypeNode)
}

func (b *NodeBuilder) getParentSymbolOfTypeParameter(typeParameter *TypeParameter) *ast.Symbol {
	tp := ast.GetDeclarationOfKind(typeParameter.symbol, ast.KindTypeParameter)
	var host *ast.Node
	// !!! JSDoc support
	// if ast.IsJSDocTemplateTag(tp.Parent) {
	// 	host = getEffectiveContainerForJSDocTemplateTag(tp.Parent)
	// } else {
	host = tp.Parent
	// }
	if host == nil {
		return nil
	}
	return b.ch.getSymbolOfNode(host)
}

func (b *NodeBuilder) typeReferenceToTypeNode(t *Type) *ast.TypeNode {
	var typeArguments []*Type = b.ch.getTypeArguments(t)
	if t.Target() == b.ch.globalArrayType || t.Target() == b.ch.globalReadonlyArrayType {
		if b.ctx.flags&nodebuilder.FlagsWriteArrayAsGenericType != 0 {
			typeArgumentNode := b.typeToTypeNodeClosure(typeArguments[0])
			return b.f.NewTypeReferenceNode(b.f.NewIdentifier(core.IfElse(t.Target() == b.ch.globalArrayType, "Array", "ReadonlyArray")), b.f.NewNodeList([]*ast.TypeNode{typeArgumentNode}))
		}
		elementType := b.typeToTypeNodeClosure(typeArguments[0])
		arrayType := b.f.NewArrayTypeNode(elementType)
		if t.Target() == b.ch.globalArrayType {
			return arrayType
		} else {
			return b.f.NewTypeOperatorNode(ast.KindReadonlyKeyword, arrayType)
		}
	} else if t.Target().objectFlags&ObjectFlagsTuple != 0 {
		typeArguments = core.SameMapIndex(typeArguments, func(t *Type, i int) *Type {
			return b.ch.removeMissingType(t, t.Target().AsTupleType().elementInfos[i].flags&ElementFlagsOptional != 0)
		})
		if len(typeArguments) > 0 {
			arity := b.ch.getTypeReferenceArity(t)
			tupleConstituentNodes := b.mapToTypeNodes(typeArguments[0:arity])
			if tupleConstituentNodes != nil {
				for i := 0; i < len(tupleConstituentNodes.Nodes); i++ {
					flags := t.Target().AsTupleType().elementInfos[i].flags
					labeledElementDeclaration := t.Target().AsTupleType().elementInfos[i].labeledDeclaration

					if labeledElementDeclaration != nil {
						tupleConstituentNodes.Nodes[i] = b.f.NewNamedTupleMember(core.IfElse(flags&ElementFlagsVariable != 0, b.f.NewToken(ast.KindDotDotDotToken), nil), b.f.NewIdentifier(b.ch.getTupleElementLabel(t.Target().AsTupleType().elementInfos[i], nil, i)), core.IfElse(flags&ElementFlagsOptional != 0, b.f.NewToken(ast.KindQuestionToken), nil), core.IfElse(flags&ElementFlagsRest != 0, b.f.NewArrayTypeNode(tupleConstituentNodes.Nodes[i]), tupleConstituentNodes.Nodes[i]))
					} else {
						switch {
						case flags&ElementFlagsVariable != 0:
							tupleConstituentNodes.Nodes[i] = b.f.NewRestTypeNode(core.IfElse(flags&ElementFlagsRest != 0, b.f.NewArrayTypeNode(tupleConstituentNodes.Nodes[i]), tupleConstituentNodes.Nodes[i]))
						case flags&ElementFlagsOptional != 0:
							tupleConstituentNodes.Nodes[i] = b.f.NewOptionalTypeNode(tupleConstituentNodes.Nodes[i])
						}
					}
				}
				tupleTypeNode := b.f.NewTupleTypeNode(tupleConstituentNodes)
				b.e.SetEmitFlags(tupleTypeNode, printer.EFSingleLine)
				if t.Target().AsTupleType().readonly {
					return b.f.NewTypeOperatorNode(ast.KindReadonlyKeyword, tupleTypeNode)
				} else {
					return tupleTypeNode
				}
			}
		}
		if b.ctx.encounteredError || (b.ctx.flags&nodebuilder.FlagsAllowEmptyTuple != 0) {
			tupleTypeNode := b.f.NewTupleTypeNode(b.f.NewNodeList([]*ast.TypeNode{}))
			b.e.SetEmitFlags(tupleTypeNode, printer.EFSingleLine)
			if t.Target().AsTupleType().readonly {
				return b.f.NewTypeOperatorNode(ast.KindReadonlyKeyword, tupleTypeNode)
			} else {
				return tupleTypeNode
			}
		}
		b.ctx.encounteredError = true
		return nil
		// TODO: GH#18217
	} else if b.ctx.flags&nodebuilder.FlagsWriteClassExpressionAsTypeLiteral != 0 && t.symbol.ValueDeclaration != nil && ast.IsClassLike(t.symbol.ValueDeclaration) && !b.ch.IsValueSymbolAccessible(t.symbol, b.ctx.enclosingDeclaration) {
		return b.createAnonymousTypeNode(t)
	} else {
		outerTypeParameters := t.Target().AsInterfaceType().OuterTypeParameters()
		i := 0
		var resultType *ast.TypeNode
		if outerTypeParameters != nil {
			length := len(outerTypeParameters)
			for i < length {
				// Find group of type arguments for type parameters with the same declaring container.
				start := i
				parent := b.getParentSymbolOfTypeParameter(outerTypeParameters[i].AsTypeParameter())
				for ok := true; ok; ok = i < length && b.getParentSymbolOfTypeParameter(outerTypeParameters[i].AsTypeParameter()) == parent { // do-while loop
					i++
				}
				// When type parameters are their own type arguments for the whole group (i.e. we have
				// the default outer type arguments), we don't show the group.

				if !slices.Equal(outerTypeParameters[start:i], typeArguments[start:i]) {
					typeArgumentSlice := b.mapToTypeNodes(typeArguments[start:i])
					restoreFlags := b.saveRestoreFlags()
					b.ctx.flags |= nodebuilder.FlagsForbidIndexedAccessSymbolReferences
					ref := b.symbolToTypeNode(parent, ast.SymbolFlagsType, typeArgumentSlice)
					restoreFlags()
					if resultType == nil {
						resultType = ref
					} else {
						resultType = b.appendReferenceToType(resultType, ref)
					}
				}
			}
		}
		var typeArgumentNodes *ast.NodeList
		if len(typeArguments) > 0 {
			typeParameterCount := 0
			typeParams := t.Target().AsInterfaceType().TypeParameters()
			if typeParams != nil {
				typeParameterCount = min(len(typeParams), len(typeArguments))

				// Maybe we should do this for more types, but for now we only elide type arguments that are
				// identical to their associated type parameters' defaults for `Iterable`, `IterableIterator`,
				// `AsyncIterable`, and `AsyncIterableIterator` to provide backwards-compatible .d.ts emit due
				// to each now having three type parameters instead of only one.
				if b.ch.isReferenceToType(t, b.ch.getGlobalIterableType()) || b.ch.isReferenceToType(t, b.ch.getGlobalIterableIteratorType()) || b.ch.isReferenceToType(t, b.ch.getGlobalAsyncIterableType()) || b.ch.isReferenceToType(t, b.ch.getGlobalAsyncIterableIteratorType()) {
					if t.AsInterfaceType().node == nil || !ast.IsTypeReferenceNode(t.AsInterfaceType().node) || t.AsInterfaceType().node.TypeArguments() == nil || len(t.AsInterfaceType().node.TypeArguments()) < typeParameterCount {
						for typeParameterCount > 0 {
							typeArgument := typeArguments[typeParameterCount-1]
							typeParameter := t.Target().AsInterfaceType().TypeParameters()[typeParameterCount-1]
							defaultType := b.ch.getDefaultFromTypeParameter(typeParameter)
							if defaultType == nil || !b.ch.isTypeIdenticalTo(typeArgument, defaultType) {
								break
							}
							typeParameterCount--
						}
					}
				}
			}

			typeArgumentNodes = b.mapToTypeNodes(typeArguments[i:typeParameterCount])
		}
		restoreFlags := b.saveRestoreFlags()
		b.ctx.flags |= nodebuilder.FlagsForbidIndexedAccessSymbolReferences
		finalRef := b.symbolToTypeNode(t.symbol, ast.SymbolFlagsType, typeArgumentNodes)
		restoreFlags()
		if resultType == nil {
			return finalRef
		} else {
			return b.appendReferenceToType(resultType, finalRef)
		}
	}
}

func (b *NodeBuilder) visitAndTransformType(t *Type, transform func(t *Type) *ast.TypeNode) *ast.TypeNode {
	typeId := t.id
	isConstructorObject := t.objectFlags&ObjectFlagsAnonymous != 0 && t.symbol != nil && t.symbol.Flags&ast.SymbolFlagsClass != 0
	var id *CompositeSymbolIdentity
	switch {
	case t.objectFlags&ObjectFlagsReference != 0 && t.AsTypeReference().node != nil:
		id = &CompositeSymbolIdentity{false, 0, ast.GetNodeId(t.AsTypeReference().node)}
	case t.flags&TypeFlagsConditional != 0:
		id = &CompositeSymbolIdentity{false, 0, ast.GetNodeId(t.AsConditionalType().root.node.AsNode())}
	case t.symbol != nil:
		id = &CompositeSymbolIdentity{isConstructorObject, ast.GetSymbolId(t.symbol), 0}
	default:
		id = nil
	}
	// Since instantiations of the same anonymous type have the same symbol, tracking symbols instead
	// of types allows us to catch circular references to instantiations of the same anonymous type

	key := CompositeTypeCacheIdentity{typeId, b.ctx.flags, b.ctx.internalFlags}
	if b.links.Has(b.ctx.enclosingDeclaration) {
		links := b.links.Get(b.ctx.enclosingDeclaration)
		cachedResult, ok := links.serializedTypes[key]
		if ok {
			// TODO:: check if we instead store late painted statements associated with this?
			for _, arg := range cachedResult.trackedSymbols {
				b.ctx.tracker.TrackSymbol(arg.symbol, arg.enclosingDeclaration, arg.meaning)
			}
			if cachedResult.truncating {
				b.ctx.truncating = true
			}
			b.ctx.approximateLength += cachedResult.addedLength
			return b.f.DeepCloneNode(cachedResult.node)
		}
	}

	var depth int
	if id != nil {
		depth = b.ctx.symbolDepth[*id]
		if depth > 10 {
			return b.createElidedInformationPlaceholder()
		}
		b.ctx.symbolDepth[*id] = depth + 1
	}
	b.ctx.visitedTypes[typeId] = true
	prevTrackedSymbols := b.ctx.trackedSymbols
	b.ctx.trackedSymbols = nil
	startLength := b.ctx.approximateLength
	result := transform(t)
	addedLength := b.ctx.approximateLength - startLength
	if !b.ctx.reportedDiagnostic && !b.ctx.encounteredError {
		links := b.links.Get(b.ctx.enclosingDeclaration)
		links.serializedTypes[key] = &SerializedTypeEntry{
			node:           result,
			truncating:     b.ctx.truncating,
			addedLength:    addedLength,
			trackedSymbols: b.ctx.trackedSymbols,
		}
	}
	delete(b.ctx.visitedTypes, typeId)
	if id != nil {
		b.ctx.symbolDepth[*id] = depth
	}
	b.ctx.trackedSymbols = prevTrackedSymbols
	return result

	// !!! TODO: Attempt node reuse or parse nodes to minimize copying once text range setting is set up
	// deepCloneOrReuseNode := func(node T) T {
	// 	if !nodeIsSynthesized(node) && getParseTreeNode(node) == node {
	// 		return node
	// 	}
	// 	return setTextRange(b.ctx, b.f.cloneNode(visitEachChildWorker(node, deepCloneOrReuseNode, nil /*b.ctx*/, deepCloneOrReuseNodes, deepCloneOrReuseNode)), node)
	// }

	// deepCloneOrReuseNodes := func(nodes *NodeArray[*ast.Node], visitor Visitor, test func(node *ast.Node) bool, start number, count number) *NodeArray[*ast.Node] {
	// 	if nodes != nil && nodes.length == 0 {
	// 		// Ensure we explicitly make a copy of an empty array; visitNodes will not do this unless the array has elements,
	// 		// which can lead to us reusing the same empty NodeArray more than once within the same AST during type noding.
	// 		return setTextRangeWorker(b.f.NewNodeArray(nil, nodes.hasTrailingComma), nodes)
	// 	}
	// 	return visitNodes(nodes, visitor, test, start, count)
	// }
}

func (b *NodeBuilder) typeToTypeNode(t *Type) *ast.TypeNode {

	inTypeAlias := b.ctx.flags & nodebuilder.FlagsInTypeAlias
	b.ctx.flags &^= nodebuilder.FlagsInTypeAlias

	if t == nil {
		if b.ctx.flags&nodebuilder.FlagsAllowEmptyUnionOrIntersection == 0 {
			b.ctx.encounteredError = true
			return nil
			// TODO: GH#18217
		}
		b.ctx.approximateLength += 3
		return b.f.NewKeywordTypeNode(ast.KindAnyKeyword)
	}

	if b.ctx.flags&nodebuilder.FlagsNoTypeReduction == 0 {
		t = b.ch.getReducedType(t)
	}

	if t.flags&TypeFlagsAny != 0 {
		if t.alias != nil {
			return t.alias.ToTypeReferenceNode(b)
		}
		// !!! TODO: add comment once synthetic comment additions to nodes are supported
		// if t == b.ch.unresolvedType {
		// 	return e.AddSyntheticLeadingComment(b.f.NewKeywordTypeNode(ast.KindAnyKeyword), ast.KindMultiLineCommentTrivia, "unresolved")
		// }
		b.ctx.approximateLength += 3
		return b.f.NewLiteralTypeNode(b.f.NewKeywordExpression(core.IfElse(t == b.ch.intrinsicMarkerType, ast.KindIntrinsicKeyword, ast.KindAnyKeyword)))
	}
	if t.flags&TypeFlagsUnknown != 0 {
		return b.f.NewKeywordTypeNode(ast.KindUnknownKeyword)
	}
	if t.flags&TypeFlagsString != 0 {
		b.ctx.approximateLength += 6
		return b.f.NewKeywordTypeNode(ast.KindStringKeyword)
	}
	if t.flags&TypeFlagsNumber != 0 {
		b.ctx.approximateLength += 6
		return b.f.NewKeywordTypeNode(ast.KindNumberKeyword)
	}
	if t.flags&TypeFlagsBigInt != 0 {
		b.ctx.approximateLength += 6
		return b.f.NewKeywordTypeNode(ast.KindBigIntKeyword)
	}
	if t.flags&TypeFlagsBoolean != 0 && t.alias == nil {
		b.ctx.approximateLength += 7
		return b.f.NewKeywordTypeNode(ast.KindBooleanKeyword)
	}
	if t.flags&TypeFlagsEnumLike != 0 {
		if t.symbol.Flags&ast.SymbolFlagsEnumMember != 0 {
			parentSymbol := b.ch.getParentOfSymbol(t.symbol)
			parentName := b.symbolToTypeNode(parentSymbol, ast.SymbolFlagsType, nil)
			if b.ch.getDeclaredTypeOfSymbol(parentSymbol) == t {
				return parentName
			}
			memberName := ast.SymbolName(t.symbol)
			if scanner.IsIdentifierText(memberName, core.ScriptTargetES5) {
				return b.appendReferenceToType(parentName /* as TypeReferenceNode | ImportTypeNode */, b.f.NewTypeReferenceNode(b.f.NewIdentifier(memberName), nil /*typeArguments*/))
			}
			if ast.IsImportTypeNode(parentName) {
				parentName.AsImportTypeNode().IsTypeOf = true
				// mutably update, node is freshly manufactured anyhow
				return b.f.NewIndexedAccessTypeNode(parentName, b.f.NewLiteralTypeNode(b.f.NewStringLiteral(memberName)))
			} else if ast.IsTypeReferenceNode(parentName) {
				return b.f.NewIndexedAccessTypeNode(b.f.NewTypeQueryNode(parentName.AsTypeReferenceNode().TypeName, nil), b.f.NewLiteralTypeNode(b.f.NewStringLiteral(memberName)))
			} else {
				panic("Unhandled type node kind returned from `symbolToTypeNode`.")
			}
		}
		return b.symbolToTypeNode(t.symbol, ast.SymbolFlagsType, nil)
	}
	if t.flags&TypeFlagsStringLiteral != 0 {
		b.ctx.approximateLength += len(t.AsLiteralType().value.(string)) + 2
		return b.f.NewLiteralTypeNode(b.f.NewStringLiteral(t.AsLiteralType().value.(string) /*, b.flags&nodebuilder.FlagsUseSingleQuotesForStringLiteralType != 0*/))
	}
	if t.flags&TypeFlagsNumberLiteral != 0 {
		value := t.AsLiteralType().value.(jsnum.Number)
		b.ctx.approximateLength += len(value.String())
		return b.f.NewLiteralTypeNode(core.IfElse(value < 0, b.f.NewPrefixUnaryExpression(ast.KindMinusToken, b.f.NewNumericLiteral("-"+value.String())), b.f.NewNumericLiteral(value.String())))
	}
	if t.flags&TypeFlagsBigIntLiteral != 0 {
		b.ctx.approximateLength += len(pseudoBigIntToString(getBigIntLiteralValue(t))) + 1
		return b.f.NewLiteralTypeNode(b.f.NewBigIntLiteral(pseudoBigIntToString(getBigIntLiteralValue(t))))
	}
	if t.flags&TypeFlagsBooleanLiteral != 0 {
		b.ctx.approximateLength += len(t.AsIntrinsicType().intrinsicName)
		return b.f.NewLiteralTypeNode(core.IfElse(t.AsIntrinsicType().intrinsicName == "true", b.f.NewKeywordExpression(ast.KindTrueKeyword), b.f.NewKeywordExpression(ast.KindFalseKeyword)))
	}
	if t.flags&TypeFlagsUniqueESSymbol != 0 {
		if b.ctx.flags&nodebuilder.FlagsAllowUniqueESSymbolType == 0 {
			if b.ch.IsValueSymbolAccessible(t.symbol, b.ctx.enclosingDeclaration) {
				b.ctx.approximateLength += 6
				return b.symbolToTypeNode(t.symbol, ast.SymbolFlagsValue, nil)
			}
			b.ctx.tracker.ReportInaccessibleUniqueSymbolError()
		}
		b.ctx.approximateLength += 13
		return b.f.NewTypeOperatorNode(ast.KindUniqueKeyword, b.f.NewKeywordTypeNode(ast.KindSymbolKeyword))
	}
	if t.flags&TypeFlagsVoid != 0 {
		b.ctx.approximateLength += 4
		return b.f.NewKeywordTypeNode(ast.KindVoidKeyword)
	}
	if t.flags&TypeFlagsUndefined != 0 {
		b.ctx.approximateLength += 9
		return b.f.NewKeywordTypeNode(ast.KindUndefinedKeyword)
	}
	if t.flags&TypeFlagsNull != 0 {
		b.ctx.approximateLength += 4
		return b.f.NewLiteralTypeNode(b.f.NewKeywordExpression(ast.KindNullKeyword))
	}
	if t.flags&TypeFlagsNever != 0 {
		b.ctx.approximateLength += 5
		return b.f.NewKeywordTypeNode(ast.KindNeverKeyword)
	}
	if t.flags&TypeFlagsESSymbol != 0 {
		b.ctx.approximateLength += 6
		return b.f.NewKeywordTypeNode(ast.KindSymbolKeyword)
	}
	if t.flags&TypeFlagsNonPrimitive != 0 {
		b.ctx.approximateLength += 6
		return b.f.NewKeywordTypeNode(ast.KindObjectKeyword)
	}
	if isThisTypeParameter(t) {
		if b.ctx.flags&nodebuilder.FlagsInObjectTypeLiteral != 0 {
			if !b.ctx.encounteredError && b.ctx.flags&nodebuilder.FlagsAllowThisInObjectLiteral == 0 {
				b.ctx.encounteredError = true
			}
			b.ctx.tracker.ReportInaccessibleThisError()
		}
		b.ctx.approximateLength += 4
		return b.f.NewThisTypeNode()
	}

	if inTypeAlias == 0 && t.alias != nil && (b.ctx.flags&nodebuilder.FlagsUseAliasDefinedOutsideCurrentScope != 0 || b.ch.IsTypeSymbolAccessible(t.alias.Symbol(), b.ctx.enclosingDeclaration)) {
		sym := t.alias.Symbol()
		typeArgumentNodes := b.mapToTypeNodes(t.alias.TypeArguments())
		if isReservedMemberName(sym.Name) && sym.Flags&ast.SymbolFlagsClass == 0 {
			return b.f.NewTypeReferenceNode(b.f.NewIdentifier(""), typeArgumentNodes)
		}
		if typeArgumentNodes != nil && len(typeArgumentNodes.Nodes) == 1 && sym == b.ch.globalArrayType.symbol {
			return b.f.NewArrayTypeNode(typeArgumentNodes.Nodes[0])
		}
		return b.symbolToTypeNode(sym, ast.SymbolFlagsType, typeArgumentNodes)
	}

	objectFlags := t.objectFlags

	if objectFlags&ObjectFlagsReference != 0 {
		// Debug.assert(t.flags&TypeFlagsObject != 0) // !!!
		if t.AsTypeReference().node != nil {
			return b.visitAndTransformType(t, b.typeReferenceToTypeNodeClosure)
		} else {
			return b.typeReferenceToTypeNodeClosure(t)
		}
	}
	if t.flags&TypeFlagsTypeParameter != 0 || objectFlags&ObjectFlagsClassOrInterface != 0 {
		if t.flags&TypeFlagsTypeParameter != 0 && slices.Contains(b.ctx.inferTypeParameters, t) {
			b.ctx.approximateLength += len(ast.SymbolName(t.symbol)) + 6
			var constraintNode *ast.TypeNode
			constraint := b.ch.getConstraintOfTypeParameter(t)
			if constraint != nil {
				// If the infer type has a constraint that is not the same as the constraint
				// we would have normally inferred based on b, we emit the constraint
				// using `infer T extends ?`. We omit inferred constraints from type references
				// as they may be elided.
				inferredConstraint := b.ch.getInferredTypeParameterConstraint(t, true /*omitTypeReferences*/)
				if !(inferredConstraint != nil && b.ch.isTypeIdenticalTo(constraint, inferredConstraint)) {
					b.ctx.approximateLength += 9
					constraintNode = b.typeToTypeNodeClosure(constraint)
				}
			}
			return b.f.NewInferTypeNode(b.typeParameterToDeclarationWithConstraint(t, constraintNode))
		}
		if b.ctx.flags&nodebuilder.FlagsGenerateNamesForShadowedTypeParams != 0 && t.flags&TypeFlagsTypeParameter != 0 {
			name := b.typeParameterToName(t)
			b.ctx.approximateLength += len(name.Text)
			return b.f.NewTypeReferenceNode(b.f.NewIdentifier(name.Text), nil /*typeArguments*/)
		}
		// Ignore constraint/default when creating a usage (as opposed to declaration) of a type parameter.
		if t.symbol != nil {
			return b.symbolToTypeNode(t.symbol, ast.SymbolFlagsType, nil)
		}
		var name string
		if (t == b.ch.markerSuperTypeForCheck || t == b.ch.markerSubTypeForCheck) && b.ch.varianceTypeParameter != nil && b.ch.varianceTypeParameter.symbol != nil {
			name = (core.IfElse(t == b.ch.markerSubTypeForCheck, "sub-", "super-")) + ast.SymbolName(b.ch.varianceTypeParameter.symbol)
		} else {
			name = "?"
		}
		return b.f.NewTypeReferenceNode(b.f.NewIdentifier(name), nil /*typeArguments*/)
	}
	if t.flags&TypeFlagsUnion != 0 && t.AsUnionType().origin != nil {
		t = t.AsUnionType().origin
	}
	if t.flags&(TypeFlagsUnion|TypeFlagsIntersection) != 0 {
		var types []*Type
		if t.flags&TypeFlagsUnion != 0 {
			types = b.ch.formatUnionTypes(t.AsUnionType().types)
		} else {
			types = t.AsIntersectionType().types
		}
		if len(types) == 1 {
			return b.typeToTypeNodeClosure(types[0])
		}
		typeNodes := b.mapToTypeNodes(types)
		if typeNodes != nil && len(typeNodes.Nodes) > 0 {
			if t.flags&TypeFlagsUnion != 0 {
				return b.f.NewUnionTypeNode(typeNodes)
			} else {
				return b.f.NewIntersectionTypeNode(typeNodes)
			}
		} else {
			if !b.ctx.encounteredError && b.ctx.flags&nodebuilder.FlagsAllowEmptyUnionOrIntersection == 0 {
				b.ctx.encounteredError = true
			}
			return nil
			// TODO: GH#18217
		}
	}
	if objectFlags&(ObjectFlagsAnonymous|ObjectFlagsMapped) != 0 {
		// Debug.assert(t.flags&TypeFlagsObject != 0) // !!!
		// The type is an object literal type.
		return b.createAnonymousTypeNode(t)
	}
	if t.flags&TypeFlagsIndex != 0 {
		indexedType := t.Target()
		b.ctx.approximateLength += 6
		indexTypeNode := b.typeToTypeNodeClosure(indexedType)
		return b.f.NewTypeOperatorNode(ast.KindKeyOfKeyword, indexTypeNode)
	}
	if t.flags&TypeFlagsTemplateLiteral != 0 {
		texts := t.AsTemplateLiteralType().texts
		types := t.AsTemplateLiteralType().types
		templateHead := b.f.NewTemplateHead(texts[0], texts[0], ast.TokenFlagsNone)
		templateSpans := b.f.NewNodeList(core.MapIndex(types, func(t *Type, i int) *ast.Node {
			var res *ast.TemplateMiddleOrTail
			if i < len(types)-1 {
				res = b.f.NewTemplateMiddle(texts[i+1], texts[i+1], ast.TokenFlagsNone)
			} else {
				res = b.f.NewTemplateTail(texts[i+1], texts[i+1], ast.TokenFlagsNone)
			}
			return b.f.NewTemplateLiteralTypeSpan(b.typeToTypeNodeClosure(t), res)
		}))
		b.ctx.approximateLength += 2
		return b.f.NewTemplateLiteralTypeNode(templateHead, templateSpans)
	}
	if t.flags&TypeFlagsStringMapping != 0 {
		typeNode := b.typeToTypeNodeClosure(t.Target())
		return b.symbolToTypeNode(t.AsStringMappingType().symbol, ast.SymbolFlagsType, b.f.NewNodeList([]*ast.Node{typeNode}))
	}
	if t.flags&TypeFlagsIndexedAccess != 0 {
		objectTypeNode := b.typeToTypeNodeClosure(t.AsIndexedAccessType().objectType)
		indexTypeNode := b.typeToTypeNodeClosure(t.AsIndexedAccessType().indexType)
		b.ctx.approximateLength += 2
		return b.f.NewIndexedAccessTypeNode(objectTypeNode, indexTypeNode)
	}
	if t.flags&TypeFlagsConditional != 0 {
		return b.visitAndTransformType(t, b.conditionalTypeToTypeNodeClosure)
	}
	if t.flags&TypeFlagsSubstitution != 0 {
		typeNode := b.typeToTypeNodeClosure(t.AsSubstitutionType().baseType)
		if !b.ch.isNoInferType(t) {
			return typeNode
		}
		noInferSymbol := b.ch.getGlobalTypeAliasSymbol("NoInfer", 1, false)
		if noInferSymbol != nil {
			return b.symbolToTypeNode(noInferSymbol, ast.SymbolFlagsType, b.f.NewNodeList([]*ast.Node{typeNode}))
		} else {
			return typeNode
		}
	}

	panic("Should be unreachable.")
}

// Direct serialization core functions for types, type aliases, and symbols

func (t *TypeAlias) ToTypeReferenceNode(b *NodeBuilder) *ast.Node {
	return b.f.NewTypeReferenceNode(b.symbolToEntityNameNode(t.Symbol()), b.mapToTypeNodes(t.TypeArguments()))
}
