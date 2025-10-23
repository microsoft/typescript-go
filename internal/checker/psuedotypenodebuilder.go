package checker

import (
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/debug"
	"github.com/microsoft/typescript-go/internal/nodebuilder"
	"github.com/microsoft/typescript-go/internal/printer"
	"github.com/microsoft/typescript-go/internal/psuedochecker"
)

// Maps a psuedochecker's psuedotypes into ast nodes and reports any inference fallback errors the psuedotype structure implies
func (b *nodeBuilderImpl) psuedoTypeToNode(t *psuedochecker.PsuedoType) *ast.Node {
	debug.Assert(t != nil, "Attempted to serialize nil psuedotype")
	switch t.Kind {
	case psuedochecker.PsuedoTypeKindDirect:
		return b.reuseTypeNode(t.Data.(*psuedochecker.PsuedoTypeDirect).TypeNode)
	case psuedochecker.PsuedoTypeKindInferred:
		node := t.Data.(*psuedochecker.PsuedoTypeInferred).Expression
		b.ctx.tracker.ReportInferenceFallback(node)
		ty := b.ch.getTypeOfExpression(node)
		return b.typeToTypeNode(ty)
	case psuedochecker.PsuedoTypeKindNoResult:
		node := t.Data.(*psuedochecker.PsuedoTypeNoResult).Declaration
		b.ctx.tracker.ReportInferenceFallback(node)
		return b.serializeTypeForDeclaration(node, nil, nil, false)
	case psuedochecker.PsuedoTypeKindUnion:
		var res []*ast.Node
		var hasElidedType bool
		members := t.Data.(*psuedochecker.PsuedoTypeUnion).Types
		for _, m := range members {
			if !b.ch.strictNullChecks {
				if m.Kind == psuedochecker.PsuedoTypeKindUndefined || m.Kind == psuedochecker.PsuedoTypeKindNull {
					hasElidedType = true
					continue
				}
			}
			res = append(res, b.psuedoTypeToNode(m))
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
	case psuedochecker.PsuedoTypeKindUndefined:
		if !b.ch.strictNullChecks {
			return b.f.NewKeywordTypeNode(ast.KindAnyKeyword)
		}
		return b.f.NewKeywordTypeNode(ast.KindUndefinedKeyword)
	case psuedochecker.PsuedoTypeKindNull:
		if !b.ch.strictNullChecks {
			return b.f.NewKeywordTypeNode(ast.KindAnyKeyword)
		}
		return b.f.NewLiteralTypeNode(b.f.NewKeywordExpression(ast.KindNullKeyword))
	case psuedochecker.PsuedoTypeKindAny:
		return b.f.NewKeywordTypeNode(ast.KindAnyKeyword)
	case psuedochecker.PsuedoTypeKindString:
		return b.f.NewKeywordTypeNode(ast.KindStringKeyword)
	case psuedochecker.PsuedoTypeKindNumber:
		return b.f.NewKeywordTypeNode(ast.KindNumberKeyword)
	case psuedochecker.PsuedoTypeKindBigInt:
		return b.f.NewKeywordTypeNode(ast.KindBigIntKeyword)
	case psuedochecker.PsuedoTypeKindBoolean:
		return b.f.NewKeywordTypeNode(ast.KindBooleanKeyword)
	case psuedochecker.PsuedoTypeKindFalse:
		return b.f.NewLiteralTypeNode(b.f.NewKeywordExpression(ast.KindFalseKeyword))
	case psuedochecker.PsuedoTypeKindTrue:
		return b.f.NewLiteralTypeNode(b.f.NewKeywordExpression(ast.KindTrueKeyword))
	case psuedochecker.PsuedoTypeKindSingleCallSignature:
		d := t.Data.(*psuedochecker.PsuedoTypeSingleCallSignature)
		var typeParams *ast.NodeList
		if len(d.TypeParameters) > 0 {
			res := make([]*ast.Node, 0, len(d.TypeParameters))
			for _, tp := range d.TypeParameters {
				res = append(res, b.reuseNode(tp.AsNode()))
			}
			typeParams = b.f.NewNodeList(res)
		}
		params := b.psuedoParametersToNodeList(d.Parameters)
		returnType := b.psuedoTypeToNode(d.ReturnType)
		return b.f.NewFunctionTypeNode(typeParams, params, returnType)
	case psuedochecker.PsuedoTypeKindTuple:
		var res []*ast.Node
		elements := t.Data.(*psuedochecker.PsuedoTypeTuple).Elements
		for _, e := range elements {
			res = append(res, b.psuedoTypeToNode(e))
		}
		// !!! TODO: psuedo-tuples are implicitly `readonly` since they originate from `as const` contexts
		// but strada fails to add the `readonly` modifier to the generated node. We replicate that bug here.
		// return b.f.NewTypeOperatorNode(ast.KindReadonlyKeyword, b.f.NewTupleTypeNode(b.f.NewNodeList(res)))
		result := b.f.NewTupleTypeNode(b.f.NewNodeList(res))
		b.e.AddEmitFlags(result, printer.EFSingleLine)
		return result
	case psuedochecker.PsuedoTypeKindObjectLiteral:
		elements := t.Data.(*psuedochecker.PsuedoTypeObjectLiteral).Elements
		if len(elements) == 0 {
			result := b.f.NewTypeLiteralNode(b.f.NewNodeList(nil))
			b.e.AddEmitFlags(result, printer.EFSingleLine)
			return result
		}
		// NOTE: using the checker's `isConstContext` instead of the psuedochecker's `isInConstContext`
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
			if isConst || (e.Kind == psuedochecker.PsuedoObjectElementKindPropertyAssignment && e.Data.(*psuedochecker.PsuedoPropertyAssignment).Readonly) {
				modifiers = b.f.NewModifierList([]*ast.Node{b.f.NewModifier(ast.KindReadonlyKeyword)})
			}
			var newProp *ast.Node
			switch e.Kind {
			case psuedochecker.PsuedoObjectElementKindMethod:
				d := e.Data.(*psuedochecker.PsuedoObjectMethod)
				newProp = b.f.NewMethodSignatureDeclaration(
					modifiers,
					b.reuseNode(e.Name),
					nil,
					nil,
					b.psuedoParametersToNodeList(d.Parameters),
					b.psuedoTypeToNode(d.ReturnType),
				)
			case psuedochecker.PsuedoObjectElementKindPropertyAssignment:
				d := e.Data.(*psuedochecker.PsuedoPropertyAssignment)
				newProp = b.f.NewPropertySignatureDeclaration(
					modifiers,
					b.reuseNode(e.Name),
					nil,
					b.psuedoTypeToNode(d.Type),
					nil,
				)
			case psuedochecker.PsuedoObjectElementKindSetAccessor:
				d := e.Data.(*psuedochecker.PsuedoSetAccessor)
				newProp = b.f.NewSetAccessorDeclaration(
					nil,
					b.reuseNode(e.Name),
					nil,
					b.f.NewNodeList([]*ast.Node{b.psuedoParameterToNode(d.Parameter)}),
					nil,
					nil,
					nil,
				)
			case psuedochecker.PsuedoObjectElementKindGetAccessor:
				d := e.Data.(*psuedochecker.PsuedoGetAccessor)
				newProp = b.f.NewSetAccessorDeclaration(
					nil,
					b.reuseNode(e.Name),
					nil,
					nil,
					b.psuedoTypeToNode(d.Type),
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
	case psuedochecker.PsuedoTypeKindStringLiteral, psuedochecker.PsuedoTypeKindNumericLiteral, psuedochecker.PsuedoTypeKindBigIntLiteral:
		source := t.Data.(*psuedochecker.PsuedoTypeLiteral).Node
		return b.f.NewLiteralTypeNode(b.reuseNode(source))
	default:
		debug.AssertNever(t.Kind, "Unhandled psuedotype kind in psuedotype node construction")
		return nil
	}
}

func (b *nodeBuilderImpl) psuedoParametersToNodeList(params []*psuedochecker.PsuedoParameter) *ast.NodeList {
	res := make([]*ast.Node, 0, len(params))
	for _, p := range params {
		res = append(res, b.psuedoParameterToNode(p))
	}
	return b.f.NewNodeList(res)
}

func (b *nodeBuilderImpl) psuedoParameterToNode(p *psuedochecker.PsuedoParameter) *ast.Node {
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
		b.psuedoTypeToNode(p.Type),
		nil,
	)
}
