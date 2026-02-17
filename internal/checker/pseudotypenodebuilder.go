package checker

import (
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/debug"
	"github.com/microsoft/typescript-go/internal/nodebuilder"
	"github.com/microsoft/typescript-go/internal/printer"
	"github.com/microsoft/typescript-go/internal/pseudochecker"
)

// Maps a pseudochecker's pseudotypes into ast nodes and reports any inference fallback errors the pseudotype structure implies
func (b *NodeBuilderImpl) pseudoTypeToNode(t *pseudochecker.PseudoType) *ast.Node {
	debug.Assert(t != nil, "Attempted to serialize nil pseudotype")
	switch t.Kind {
	case pseudochecker.PseudoTypeKindDirect:
		return b.reuseTypeNode(t.AsPseudoTypeDirect().TypeNode)
	case pseudochecker.PseudoTypeKindInferred:
		node := t.AsPseudoTypeInferred().Expression
		b.ctx.tracker.ReportInferenceFallback(node)
		ty := b.ch.getTypeOfExpression(node)
		return b.typeToTypeNode(ty)
	case pseudochecker.PseudoTypeKindNoResult:
		node := t.AsPseudoTypeNoResult().Declaration
		b.ctx.tracker.ReportInferenceFallback(node)
		if ast.IsFunctionLike(node) && !ast.IsAccessor(node) {
			return b.serializeReturnTypeForSignature(b.ch.getSignatureFromDeclaration(node), false)
		}
		return b.serializeTypeForDeclaration(node, nil, nil, false)
	case pseudochecker.PseudoTypeKindMaybeConstLocation:
		d := t.AsPseudoTypeMaybeConstLocation()
		// see checkExpressionWithContextualType for general literal widening rules which need to be emulated here, plus
		// checkTemplateLiteralExpression for template literal widening rules if the pseudochecker ever supports literalized templates
		isInConstContext := b.ch.isConstContext(d.Node)
		if !isInConstContext {
			contextualType := b.ch.getContextualType(d.Node, ContextFlagsNone)
			t := b.pseudoTypeToType(d.ConstType)
			if t != nil && b.ch.isLiteralOfContextualType(t, b.ch.instantiateContextualType(contextualType, d.Node, ContextFlagsNone)) {
				isInConstContext = true
			}
		}
		if isInConstContext {
			return b.pseudoTypeToNode(d.ConstType)
		} else {
			return b.pseudoTypeToNode(d.RegularType)
		}
	case pseudochecker.PseudoTypeKindUnion:
		var res []*ast.Node
		var hasElidedType bool
		members := t.AsPseudoTypeUnion().Types
		for _, m := range members {
			if !b.ch.strictNullChecks {
				if m.Kind == pseudochecker.PseudoTypeKindUndefined || m.Kind == pseudochecker.PseudoTypeKindNull {
					hasElidedType = true
					continue
				}
			}
			res = append(res, b.pseudoTypeToNode(m))
		}
		if len(res) == 1 {
			return res[0]
		}
		if len(res) == 0 {
			if hasElidedType {
				return b.f.NewKeywordTypeNode(ast.KindAnyKeyword)
			}
			return b.f.NewKeywordTypeNode(ast.KindNeverKeyword)
		}
		return b.f.NewUnionTypeNode(b.f.NewNodeList(res))
	case pseudochecker.PseudoTypeKindUndefined:
		if !b.ch.strictNullChecks {
			return b.f.NewKeywordTypeNode(ast.KindAnyKeyword)
		}
		return b.f.NewKeywordTypeNode(ast.KindUndefinedKeyword)
	case pseudochecker.PseudoTypeKindNull:
		if !b.ch.strictNullChecks {
			return b.f.NewKeywordTypeNode(ast.KindAnyKeyword)
		}
		return b.f.NewLiteralTypeNode(b.f.NewKeywordExpression(ast.KindNullKeyword))
	case pseudochecker.PseudoTypeKindAny:
		return b.f.NewKeywordTypeNode(ast.KindAnyKeyword)
	case pseudochecker.PseudoTypeKindString:
		return b.f.NewKeywordTypeNode(ast.KindStringKeyword)
	case pseudochecker.PseudoTypeKindNumber:
		return b.f.NewKeywordTypeNode(ast.KindNumberKeyword)
	case pseudochecker.PseudoTypeKindBigInt:
		return b.f.NewKeywordTypeNode(ast.KindBigIntKeyword)
	case pseudochecker.PseudoTypeKindBoolean:
		return b.f.NewKeywordTypeNode(ast.KindBooleanKeyword)
	case pseudochecker.PseudoTypeKindFalse:
		return b.f.NewLiteralTypeNode(b.f.NewKeywordExpression(ast.KindFalseKeyword))
	case pseudochecker.PseudoTypeKindTrue:
		return b.f.NewLiteralTypeNode(b.f.NewKeywordExpression(ast.KindTrueKeyword))
	case pseudochecker.PseudoTypeKindSingleCallSignature:
		d := t.AsPseudoTypeSingleCallSignature()
		var typeParams *ast.NodeList
		if len(d.TypeParameters) > 0 {
			res := make([]*ast.Node, 0, len(d.TypeParameters))
			for _, tp := range d.TypeParameters {
				res = append(res, b.reuseNode(tp.AsNode()))
			}
			typeParams = b.f.NewNodeList(res)
		}
		params := b.pseudoParametersToNodeList(d.Parameters)
		returnType := b.pseudoTypeToNode(d.ReturnType)
		return b.f.NewFunctionTypeNode(typeParams, params, returnType)
	case pseudochecker.PseudoTypeKindTuple:
		var res []*ast.Node
		elements := t.AsPseudoTypeTuple().Elements
		for _, e := range elements {
			res = append(res, b.pseudoTypeToNode(e))
		}
		// !!! TODO: pseudo-tuples are implicitly `readonly` since they originate from `as const` contexts
		// but strada fails to add the `readonly` modifier to the generated node. We replicate that bug here.
		// return b.f.NewTypeOperatorNode(ast.KindReadonlyKeyword, b.f.NewTupleTypeNode(b.f.NewNodeList(res)))
		result := b.f.NewTupleTypeNode(b.f.NewNodeList(res))
		b.e.AddEmitFlags(result, printer.EFSingleLine)
		return result
	case pseudochecker.PseudoTypeKindObjectLiteral:
		elements := t.AsPseudoTypeObjectLiteral().Elements
		if len(elements) == 0 {
			result := b.f.NewTypeLiteralNode(b.f.NewNodeList(nil))
			b.e.AddEmitFlags(result, printer.EFSingleLine)
			return result
		}
		// NOTE: using the checker's `isConstContext` instead of the pseudochecker's `isInConstContext`
		// results in different results here. The checker one is more "correct" but means we'll mark
		// objects in parameter positions contextually typed by const type parameters as readonly -
		// something a true syntactic ID emitter couldn't possibly know (since the signature could
		// be from across files). This can't *really* happen in any cases ID doesn't already error on, though.
		// Just something to keep in mind if the ID checker keeps growing.
		isConst := b.ch.isConstContext(elements[0].Name)
		newElements := make([]*ast.Node, 0, len(elements))
		// TODO: strada's ID logic is piecemeal in `name` reuse validation - only methods remap `new` to `"new"`
		// we should have a unified `reuseName` codepath that remaps keyword ID names to string literal names
		for _, e := range elements {
			var modifiers *ast.ModifierList
			if isConst || (e.Kind == pseudochecker.PseudoObjectElementKindPropertyAssignment && e.AsPseudoPropertyAssignment().Readonly) {
				modifiers = b.f.NewModifierList([]*ast.Node{b.f.NewModifier(ast.KindReadonlyKeyword)})
			}
			var newProp *ast.Node
			switch e.Kind {
			case pseudochecker.PseudoObjectElementKindMethod:
				d := e.AsPseudoObjectMethod()
				newProp = b.f.NewMethodSignatureDeclaration(
					modifiers,
					b.reuseNode(e.Name),
					nil,
					nil,
					b.pseudoParametersToNodeList(d.Parameters),
					b.pseudoTypeToNode(d.ReturnType),
				)
			case pseudochecker.PseudoObjectElementKindPropertyAssignment:
				d := e.AsPseudoPropertyAssignment()
				newProp = b.f.NewPropertySignatureDeclaration(
					modifiers,
					b.reuseNode(e.Name),
					nil,
					b.pseudoTypeToNode(d.Type),
					nil,
				)
			case pseudochecker.PseudoObjectElementKindSetAccessor:
				d := e.AsPseudoSetAccessor()
				newProp = b.f.NewSetAccessorDeclaration(
					nil,
					b.reuseNode(e.Name),
					nil,
					b.f.NewNodeList([]*ast.Node{b.pseudoParameterToNode(d.Parameter)}),
					nil,
					nil,
					nil,
				)
			case pseudochecker.PseudoObjectElementKindGetAccessor:
				d := e.AsPseudoGetAccessor()
				newProp = b.f.NewGetAccessorDeclaration(
					nil,
					b.reuseNode(e.Name),
					nil,
					nil,
					b.pseudoTypeToNode(d.Type),
					nil,
					nil,
				)
			}
			if b.ctx.enclosingFile == ast.GetSourceFileOfNode(e.Name) {
				b.e.SetCommentRange(newProp, e.Name.Parent.Loc)
			}
			newElements = append(newElements, newProp)
		}
		result := b.f.NewTypeLiteralNode(b.f.NewNodeList(newElements))
		if b.ctx.flags&nodebuilder.FlagsMultilineObjectLiterals == 0 {
			b.e.AddEmitFlags(result, printer.EFSingleLine)
		}
		return result
	case pseudochecker.PseudoTypeKindStringLiteral, pseudochecker.PseudoTypeKindNumericLiteral, pseudochecker.PseudoTypeKindBigIntLiteral:
		source := t.AsPseudoTypeLiteral().Node
		return b.f.NewLiteralTypeNode(b.reuseNode(source))
	default:
		debug.AssertNever(t.Kind, "Unhandled pseudotype kind in pseudotype node construction")
		return nil
	}
}

func (b *NodeBuilderImpl) pseudoParametersToNodeList(params []*pseudochecker.PseudoParameter) *ast.NodeList {
	res := make([]*ast.Node, 0, len(params))
	for _, p := range params {
		res = append(res, b.pseudoParameterToNode(p))
	}
	return b.f.NewNodeList(res)
}

func (b *NodeBuilderImpl) pseudoParameterToNode(p *pseudochecker.PseudoParameter) *ast.Node {
	var dotDotDot *ast.Node
	var questionMark *ast.Node
	if p.Rest {
		dotDotDot = b.f.NewToken(ast.KindDotDotDotToken)
	}
	if p.Optional {
		questionMark = b.f.NewToken(ast.KindQuestionToken)
	}
	return b.f.NewParameterDeclaration(
		nil,
		dotDotDot,
		b.reuseNode(p.Name),
		questionMark,
		b.pseudoTypeToNode(p.Type),
		nil,
	)
}

func (b *NodeBuilderImpl) pseudoTypeEquivalentToType(t *pseudochecker.PseudoType, type_ *Type) bool {
	if type_ != nil && b.ch.isErrorType(type_) {
		return true
	}
	typeFromPseudo := b.pseudoTypeToType(t) // note: cannot convert complex types like objects, which must be validated separately
	if typeFromPseudo == type_ {
		return true
	}
	if t.Kind == pseudochecker.PseudoTypeKindObjectLiteral || t.Kind == pseudochecker.PseudoTypeKindTuple {
		return true // !!! TODO: validate. This relies on the psuedochecker tossing out all complex literals with insufficient/incorrect annotations/assertions
	}
	return false
}

func (b *NodeBuilderImpl) pseudoTypeToType(t *pseudochecker.PseudoType) *Type {
	// !!! TODO: only literal types currently mapped because this is only used to determine if literal contextual typing need apply to the pseudotype
	// If this is used more broadly, the implementation needs to be filled out more to handle the structural pseudotypes - signatures, objects, tuples, etc
	debug.Assert(t != nil, "Attempted to realize nil pseudotype")
	switch t.Kind {
	case pseudochecker.PseudoTypeKindDirect:
		return b.ch.getTypeFromTypeNode(t.AsPseudoTypeDirect().TypeNode)
	case pseudochecker.PseudoTypeKindInferred:
		node := t.AsPseudoTypeInferred().Expression
		ty := b.ch.getTypeOfExpression(node)
		return ty
	case pseudochecker.PseudoTypeKindNoResult:
		return nil // TODO: extract type selection logic from `serializeTypeForDeclaration`, not needed for current usecases but needed if completeness becomes required
	case pseudochecker.PseudoTypeKindMaybeConstLocation:
		d := t.AsPseudoTypeMaybeConstLocation()
		if b.ch.isConstContext(d.Node) {
			return b.pseudoTypeToType(d.ConstType)
		}
		return b.pseudoTypeToType(d.RegularType)
	case pseudochecker.PseudoTypeKindUnion:
		var res []*Type
		var hasElidedType bool
		members := t.AsPseudoTypeUnion().Types
		for _, m := range members {
			if !b.ch.strictNullChecks {
				if m.Kind == pseudochecker.PseudoTypeKindUndefined || m.Kind == pseudochecker.PseudoTypeKindNull {
					hasElidedType = true
					continue
				}
			}
			t := b.pseudoTypeToType(m)
			if t == nil {
				return nil // propagate failure
			}
			res = append(res, t)
		}
		if len(res) == 1 {
			return res[0]
		}
		if len(res) == 0 {
			if hasElidedType {
				return b.ch.anyType
			}
			return b.ch.neverType
		}
		return b.ch.newUnionType(ObjectFlagsNone, res)
	case pseudochecker.PseudoTypeKindUndefined:
		return b.ch.undefinedWideningType
	case pseudochecker.PseudoTypeKindNull:
		return b.ch.nullWideningType
	case pseudochecker.PseudoTypeKindAny:
		return b.ch.anyType
	case pseudochecker.PseudoTypeKindString:
		return b.ch.stringType
	case pseudochecker.PseudoTypeKindNumber:
		return b.ch.numberType
	case pseudochecker.PseudoTypeKindBigInt:
		return b.ch.bigintType
	case pseudochecker.PseudoTypeKindBoolean:
		return b.ch.booleanType
	case pseudochecker.PseudoTypeKindFalse:
		return b.ch.falseType
	case pseudochecker.PseudoTypeKindTrue:
		return b.ch.trueType
	case pseudochecker.PseudoTypeKindStringLiteral, pseudochecker.PseudoTypeKindNumericLiteral, pseudochecker.PseudoTypeKindBigIntLiteral:
		source := t.AsPseudoTypeLiteral().Node
		return b.ch.getTypeOfExpression(source) // big shortcut, uses cached expression types where possible
	default:
		return nil
	}
}
