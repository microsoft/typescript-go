package printer

import "github.com/microsoft/typescript-go/internal/ast"

type PrivateIdentifierKind int

const (
	PrivateIdentifierKindMethod PrivateIdentifierKind = iota
	PrivateIdentifierKindField
	PrivateIdentifierKindAccessor
	PrivateIdentifierKindUntransformed
)

func (k PrivateIdentifierKind) String() string {
	switch k {
	case PrivateIdentifierKindMethod:
		return "m"
	case PrivateIdentifierKindField:
		return "f"
	case PrivateIdentifierKindAccessor:
		return "a"
	case PrivateIdentifierKindUntransformed:
		return "untransformed"
	default:
		panic("Unhandled PrivateIdentifierKind")
	}
}

type PrivateIdentifierInfo interface {
	Kind() PrivateIdentifierKind
	IsValid() bool
	IsStatic() bool
	BrandCheckIdentifier() *ast.Identifier
}

type PrivateIdentifierInfoBase struct {
	/**
	 * brandCheckIdentifier can contain:
	 *  - For instance field: The WeakMap that will be the storage for the field.
	 *  - For instance methods or accessors: The WeakSet that will be used for brand checking.
	 *  - For static members: The constructor that will be used for brand checking.
	 */
	brandCheckIdentifier *ast.Identifier
	// Stores if the identifier is static or not
	isStatic bool
	// Stores if the identifier declaration is valid or not. Reserved names (e.g. #constructor)
	// or duplicate identifiers are considered invalid.
	isValid bool
}

type PrivateIdentifierAccessorInfo struct {
	PrivateIdentifierInfoBase
	kind PrivateIdentifierKind
	// Identifier for a variable that will contain the private get accessor implementation, if any.
	GetterName *ast.Identifier
	// Identifier for a variable that will contain the private set accessor implementation, if any.
	SetterName *ast.Identifier
}

func (p *PrivateIdentifierAccessorInfo) Kind() PrivateIdentifierKind {
	return p.kind
}

func (p *PrivateIdentifierAccessorInfo) IsValid() bool {
	return p.PrivateIdentifierInfoBase.isValid
}

func (p *PrivateIdentifierAccessorInfo) IsStatic() bool {
	return p.PrivateIdentifierInfoBase.isStatic
}

func (p *PrivateIdentifierAccessorInfo) BrandCheckIdentifier() *ast.Identifier {
	return p.PrivateIdentifierInfoBase.brandCheckIdentifier
}

func NewPrivateIdentifierAccessorInfo(brandCheckIdentifier *ast.Identifier, getterName *ast.Identifier, setterName *ast.Identifier, isValid bool, isStatic bool) *PrivateIdentifierAccessorInfo {
	return &PrivateIdentifierAccessorInfo{
		kind: PrivateIdentifierKindAccessor,
		PrivateIdentifierInfoBase: PrivateIdentifierInfoBase{
			brandCheckIdentifier: brandCheckIdentifier,
			isStatic:             isStatic,
			isValid:              isValid,
		},
		GetterName: getterName,
		SetterName: setterName,
	}
}

type PrivateIdentifierMethodInfo struct {
	PrivateIdentifierInfoBase
	kind PrivateIdentifierKind
	// Identifier for a variable that will contain the private method implementation.
	MethodName *ast.Identifier
}

func (p *PrivateIdentifierMethodInfo) Kind() PrivateIdentifierKind {
	return p.kind
}

func (p *PrivateIdentifierMethodInfo) IsValid() bool {
	return p.PrivateIdentifierInfoBase.isValid
}

func (p *PrivateIdentifierMethodInfo) IsStatic() bool {
	return p.PrivateIdentifierInfoBase.isStatic
}

func (p *PrivateIdentifierMethodInfo) BrandCheckIdentifier() *ast.Identifier {
	return p.PrivateIdentifierInfoBase.brandCheckIdentifier
}

func NewPrivateIdentifierMethodInfo(brandCheckIdentifier *ast.Identifier, methodName *ast.Identifier, isValid bool, isStatic bool) *PrivateIdentifierMethodInfo {
	return &PrivateIdentifierMethodInfo{
		kind: PrivateIdentifierKindMethod,
		PrivateIdentifierInfoBase: PrivateIdentifierInfoBase{
			brandCheckIdentifier: brandCheckIdentifier,
			isStatic:             isStatic,
			isValid:              isValid,
		},
		MethodName: methodName,
	}
}

type PrivateIdentifierInstanceFieldInfo struct {
	PrivateIdentifierInfoBase
	kind PrivateIdentifierKind
}

func NewPrivateIdentifierInstanceFieldInfo(brandCheckIdentifier *ast.Identifier, isValid bool) *PrivateIdentifierInstanceFieldInfo {
	return &PrivateIdentifierInstanceFieldInfo{
		kind: PrivateIdentifierKindField,
		PrivateIdentifierInfoBase: PrivateIdentifierInfoBase{
			brandCheckIdentifier: brandCheckIdentifier,
			isStatic:             false,
			isValid:              isValid,
		},
	}
}

func (p *PrivateIdentifierInstanceFieldInfo) Kind() PrivateIdentifierKind {
	return p.kind
}

func (p *PrivateIdentifierInstanceFieldInfo) IsValid() bool {
	return p.PrivateIdentifierInfoBase.isValid
}

func (p *PrivateIdentifierInstanceFieldInfo) IsStatic() bool {
	return p.PrivateIdentifierInfoBase.isStatic
}

func (p *PrivateIdentifierInstanceFieldInfo) BrandCheckIdentifier() *ast.Identifier {
	return p.PrivateIdentifierInfoBase.brandCheckIdentifier
}

type PrivateIdentifierStaticFieldInfo struct {
	PrivateIdentifierInfoBase
	kind PrivateIdentifierKind
	// Contains the variable that will serve as the storage for the field.
	VariableName *ast.Identifier
}

func NewPrivateIdentifierStaticFieldInfo(brandCheckIdentifier *ast.Identifier, variableName *ast.Identifier, isValid bool) *PrivateIdentifierStaticFieldInfo {
	return &PrivateIdentifierStaticFieldInfo{
		kind: PrivateIdentifierKindField,
		PrivateIdentifierInfoBase: PrivateIdentifierInfoBase{
			brandCheckIdentifier: brandCheckIdentifier,
			isStatic:             true,
			isValid:              isValid,
		},
		VariableName: variableName,
	}
}

func (p *PrivateIdentifierStaticFieldInfo) Kind() PrivateIdentifierKind {
	return p.kind
}

func (p *PrivateIdentifierStaticFieldInfo) IsValid() bool {
	return p.PrivateIdentifierInfoBase.isValid
}

func (p *PrivateIdentifierStaticFieldInfo) IsStatic() bool {
	return p.PrivateIdentifierInfoBase.isStatic
}

func (p *PrivateIdentifierStaticFieldInfo) BrandCheckIdentifier() *ast.Identifier {
	return p.PrivateIdentifierInfoBase.brandCheckIdentifier
}

type PrivateIdentifierUntransformedInfo struct {
	PrivateIdentifierInfoBase
	kind PrivateIdentifierKind
}

func NewPrivateIdentifierUntransformedInfo(brandCheckIdentifier *ast.Identifier, isValid bool, isStatic bool) *PrivateIdentifierUntransformedInfo {
	return &PrivateIdentifierUntransformedInfo{
		kind: PrivateIdentifierKindUntransformed,
		PrivateIdentifierInfoBase: PrivateIdentifierInfoBase{
			brandCheckIdentifier: brandCheckIdentifier,
			isStatic:             isStatic,
			isValid:              isValid,
		},
	}
}

func (p *PrivateIdentifierUntransformedInfo) Kind() PrivateIdentifierKind {
	return p.kind
}

func (p *PrivateIdentifierUntransformedInfo) IsValid() bool {
	return p.PrivateIdentifierInfoBase.isValid
}

func (p *PrivateIdentifierUntransformedInfo) IsStatic() bool {
	return p.PrivateIdentifierInfoBase.isStatic
}

func (p *PrivateIdentifierUntransformedInfo) BrandCheckIdentifier() *ast.Identifier {
	return p.PrivateIdentifierInfoBase.brandCheckIdentifier
}
