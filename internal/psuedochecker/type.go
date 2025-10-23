package psuedochecker

import (
	"github.com/microsoft/typescript-go/internal/ast"
)

// `PsuedoType`s are skeletons of types - partially interpreted expressions and type nodes
// composed to represent how you *should* construct a type out of them. They can be trivially
// mapped into actual types by a real `Checker`, or into a tree of `Node`s directly, without
// needing to make any intermediate types, by a `NodeBuilder`. Unlike checker `Type`s, these are
// never normalized, and multiple psuedo-types may refer to the same underlying `Type`.

// In strada, these were implicit in the AST nodes constructed in `expressionToTypeNode.ts`, which
// repurposed AST nodes for this purpose, but in so doing, often confused weather or not it had validated
// nested nodes for use at a given use-site. By keeping the mapping deferred like this, we can know we haven't
// done any use-site checks until we're ready to map the `PsuedoType` into a `Node`, and can cache
// `PsuedoType`s across multiple target positions.

type PsuedoTypeKind int16

const (
	PsuedoTypeKindDirect PsuedoTypeKind = iota
	PsuedoTypeKindInferred
	PsuedoTypeKindNoResult
	PsuedoTypeKindMaybeConstLocation
	PsuedoTypeKindUnion
	PsuedoTypeKindUndefined
	PsuedoTypeKindNull
	PsuedoTypeKindAny
	PsuedoTypeKindString
	PsuedoTypeKindNumber
	PsuedoTypeKindBigInt
	PsuedoTypeKindBoolean
	PsuedoTypeKindFalse
	PsuedoTypeKindTrue
	PsuedoTypeKindSingleCallSignature
	PsuedoTypeKindTuple
	PsuedoTypeKindObjectLiteral
	PsuedoTypeKindStringLiteral
	PsuedoTypeKindNumericLiteral
	PsuedoTypeKindBigIntLiteral
)

type PsuedoType struct {
	Kind PsuedoTypeKind
	Data psuedoTypeData
}

func NewPsuedoType(kind PsuedoTypeKind, data psuedoTypeData) *PsuedoType {
	n := data.AsPsuedoType()
	n.Kind = kind
	n.Data = data
	return n
}

type psuedoTypeData interface {
	AsPsuedoType() *PsuedoType
}

type PsuedoTypeDefault struct {
	PsuedoType
}

func (b *PsuedoTypeDefault) AsPsuedoType() *PsuedoType { return &b.PsuedoType }

type PsuedoTypeBase struct {
	PsuedoTypeDefault
}

var (
	PsuedoTypeUndefined = NewPsuedoType(PsuedoTypeKindUndefined, &PsuedoTypeBase{})
	PsuedoTypeNull      = NewPsuedoType(PsuedoTypeKindNull, &PsuedoTypeBase{})
	PsuedoTypeAny       = NewPsuedoType(PsuedoTypeKindAny, &PsuedoTypeBase{})
	PsuedoTypeString    = NewPsuedoType(PsuedoTypeKindString, &PsuedoTypeBase{})
	PsuedoTypeNumber    = NewPsuedoType(PsuedoTypeKindNumber, &PsuedoTypeBase{})
	PsuedoTypeBigInt    = NewPsuedoType(PsuedoTypeKindBigInt, &PsuedoTypeBase{})
	PsuedoTypeBoolean   = NewPsuedoType(PsuedoTypeKindBoolean, &PsuedoTypeBase{})
	PsuedoTypeFalse     = NewPsuedoType(PsuedoTypeKindFalse, &PsuedoTypeBase{})
	PsuedoTypeTrue      = NewPsuedoType(PsuedoTypeKindTrue, &PsuedoTypeBase{})
)

// PsuedoTypeDirect directly encodes the type referred to by a given TypeNode
type PsuedoTypeDirect struct {
	PsuedoTypeBase
	TypeNode *ast.Node
}

func NewPsuedoTypeDirect(typeNode *ast.Node) *PsuedoType {
	return NewPsuedoType(PsuedoTypeKindDirect, &PsuedoTypeDirect{TypeNode: typeNode})
}

// PsuedoTypeInferred directly encodes the type referred to by a given Expression
// These represent cases where the expression was too complex for the psuedochecker.
// Most of the time, these locations will produce an error under ID.
type PsuedoTypeInferred struct {
	PsuedoTypeBase
	Expression *ast.Node
}

func NewPsuedoTypeInferred(expr *ast.Node) *PsuedoType {
	return NewPsuedoType(PsuedoTypeKindInferred, &PsuedoTypeInferred{Expression: expr})
}

// PsuedoTypeNoResult is anlogous to PsuedoTypeInferred in that it references a case
// where the type was too complex for the psuedochecker. Rather than an expression, however,
// it is referring to the return type of a signature or declaration.
type PsuedoTypeNoResult struct {
	PsuedoTypeBase
	Declaration *ast.Node
}

func NewPsuedoTypeNoResult(decl *ast.Node) *PsuedoType {
	return NewPsuedoType(PsuedoTypeKindNoResult, &PsuedoTypeNoResult{Declaration: decl})
}

// PsuedoTypeMaybeConstLocation encodes the const/regular types of a location so the builder
// can later select the appropriate psuedotype based on the location's context. This is used
// to ensure accuracy in nested expressions without exposing type-based functionality to the psuedochecker.
// A nodebuilder that doesn't do contextual typing would need to, as policy, reject these types if they
// are in a contextually typed position! (Otherwise they could pick one, but either type could be wrong, depending on context!)
// At the top-level, which is generally what ID is concerned with, nothing is contextually typed, so these cases don't generally
// cause problems. Once you get into reused nodes in nested expressions, however, this becomes important.
// In strada, checker `isConstContext` functionality exposed to the psuedochecker + type comparison sanity checking
// on nested results masks the need for this abstraction, but with it present it clearly highlights a shortcoming
// of the ID infernce model and how "standalone" it can(n't) truly be without substantial restrictions on expression inference.
type PsuedoTypeMaybeConstLocation struct {
	PsuedoTypeBase
	Node        *ast.Node
	ConstType   *PsuedoType
	RegularType *PsuedoType
}

func NewPsuedoTypeMaybeConstLocation(loc *ast.Node, ct *PsuedoType, reg *PsuedoType) *PsuedoType {
	return NewPsuedoType(PsuedoTypeKindMaybeConstLocation, &PsuedoTypeMaybeConstLocation{Node: loc, ConstType: ct, RegularType: reg})
}

// PsuedoTypeUnion is a collection of psudotypes joined into a union
type PsuedoTypeUnion struct {
	PsuedoTypeBase
	Types []*PsuedoType
}

func NewPsuedoTypeUnion(types []*PsuedoType) *PsuedoType {
	return NewPsuedoType(PsuedoTypeKindUnion, &PsuedoTypeUnion{Types: types})
}

type PsuedoParameter struct {
	Rest     bool
	Name     *ast.Node
	Optional bool
	Type     *PsuedoType
}

func NewPsuedoParameter(isRest bool, name *ast.Node, isOptional bool, t *PsuedoType) *PsuedoParameter {
	return &PsuedoParameter{Rest: isRest, Name: name, Optional: isOptional, Type: t}
}

// PsuedoTypeSingleCallSignature represents an object type with a single call signature, like an arrow or function expression
type PsuedoTypeSingleCallSignature struct {
	PsuedoTypeBase
	Parameters     []*PsuedoParameter
	TypeParameters []*ast.TypeParameterDeclaration
	ReturnType     *PsuedoType
}

func NewPsuedoTypeSingleCallSignature(parameters []*PsuedoParameter, typeParameters []*ast.TypeParameterDeclaration, returnType *PsuedoType) *PsuedoType {
	return NewPsuedoType(PsuedoTypeKindSingleCallSignature, &PsuedoTypeSingleCallSignature{
		Parameters:     parameters,
		TypeParameters: typeParameters,
		ReturnType:     returnType,
	})
}

// PsuedoTypeTuple represents a tuple originaing from an `as const` array literal
type PsuedoTypeTuple struct {
	PsuedoTypeBase
	Elements []*PsuedoType
}

func NewPsuedoTypeTuple(elements []*PsuedoType) *PsuedoType {
	return NewPsuedoType(PsuedoTypeKindTuple, &PsuedoTypeTuple{
		Elements: elements,
	})
}

type PsuedoObjectElement struct {
	Name     *ast.Node
	Optional bool
	Kind     PsuedoObjectElementKind
	Data     psuedoObjectElementData
}

func (e *PsuedoObjectElement) AsPsuedoObjectElement() *PsuedoObjectElement { return e }

type PsuedoObjectElementKind int8

const (
	PsuedoObjectElementKindMethod PsuedoObjectElementKind = iota
	PsuedoObjectElementKindPropertyAssignment
	PsuedoObjectElementKindSetAccessor
	PsuedoObjectElementKindGetAccessor
)

type psuedoObjectElementData interface {
	AsPsuedoObjectElement() *PsuedoObjectElement
}

func NewPsuedoObjectElement(kind PsuedoObjectElementKind, name *ast.Node, optional bool, data psuedoObjectElementData) *PsuedoObjectElement {
	e := data.AsPsuedoObjectElement()
	e.Kind = kind
	e.Name = name
	e.Optional = optional
	e.Data = data
	return e
}

type PsuedoObjectMethod struct {
	PsuedoObjectElement
	Parameters []*PsuedoParameter
	ReturnType *PsuedoType
}

func NewPsuedoObjectMethod(name *ast.Node, optional bool, parameters []*PsuedoParameter, returnType *PsuedoType) *PsuedoObjectElement {
	return NewPsuedoObjectElement(PsuedoObjectElementKindMethod, name, optional, &PsuedoObjectMethod{
		Parameters: parameters,
		ReturnType: returnType,
	})
}

type PsuedoPropertyAssignment struct {
	PsuedoObjectElement
	Readonly bool
	Type     *PsuedoType
}

func NewPsuedoPropertyAssignment(readonly bool, name *ast.Node, optional bool, t *PsuedoType) *PsuedoObjectElement {
	return NewPsuedoObjectElement(PsuedoObjectElementKindPropertyAssignment, name, optional, &PsuedoPropertyAssignment{
		Readonly: readonly,
		Type:     t,
	})
}

type PsuedoSetAccessor struct {
	PsuedoObjectElement
	Parameter *PsuedoParameter
}

func NewPsuedoSetAccessor(name *ast.Node, optional bool, p *PsuedoParameter) *PsuedoObjectElement {
	return NewPsuedoObjectElement(PsuedoObjectElementKindSetAccessor, name, optional, &PsuedoSetAccessor{
		Parameter: p,
	})
}

type PsuedoGetAccessor struct {
	PsuedoObjectElement
	Type *PsuedoType
}

func NewPsuedoGetAccessor(name *ast.Node, optional bool, t *PsuedoType) *PsuedoObjectElement {
	return NewPsuedoObjectElement(PsuedoObjectElementKindGetAccessor, name, optional, &PsuedoGetAccessor{
		Type: t,
	})
}

// PsuedoTypeObjectLiteral represents an object type originaing from an object literal
type PsuedoTypeObjectLiteral struct {
	PsuedoTypeBase
	Elements []*PsuedoObjectElement
}

func NewPsuedoTypeObjectLiteral(elements []*PsuedoObjectElement) *PsuedoType {
	return NewPsuedoType(PsuedoTypeKindObjectLiteral, &PsuedoTypeObjectLiteral{
		Elements: elements,
	})
}

// PsuedoTypeLiteral represents a literal type
type PsuedoTypeLiteral struct {
	PsuedoTypeBase
	Node *ast.Node
}

func NewPsuedoTypeStringLiteral(node *ast.Node) *PsuedoType {
	return NewPsuedoType(PsuedoTypeKindStringLiteral, &PsuedoTypeLiteral{
		Node: node,
	})
}

func NewPsuedoTypeNumericLiteral(node *ast.Node) *PsuedoType {
	return NewPsuedoType(PsuedoTypeKindNumericLiteral, &PsuedoTypeLiteral{
		Node: node,
	})
}

func NewPsuedoTypeBigIntLiteral(node *ast.Node) *PsuedoType {
	return NewPsuedoType(PsuedoTypeKindBigIntLiteral, &PsuedoTypeLiteral{
		Node: node,
	})
}
