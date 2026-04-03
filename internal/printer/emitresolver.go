package printer

import (
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/binder"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/evaluator"
	"github.com/microsoft/typescript-go/internal/nodebuilder"
)

type SymbolAccessibility int32

const (
	SymbolAccessibilityAccessible SymbolAccessibility = iota
	SymbolAccessibilityNotAccessible
	SymbolAccessibilityCannotBeNamed
	SymbolAccessibilityNotResolved
)

type SymbolAccessibilityResult struct {
	Accessibility        SymbolAccessibility
	AliasesToMakeVisible []*ast.Node // aliases that need to have this symbol visible
	ErrorSymbolName      string      // Optional - symbol name that results in error
	ErrorNode            *ast.Node   // Optional - node that results in error
	ErrorModuleName      string      // Optional - If the symbol is not visible from module, module's name
}

/**
 * Indicates how to serialize the name for a TypeReference when emitting decorator metadata
 *
 * @internal
 */
type TypeReferenceSerializationKind int32

const (
	// The TypeReference could not be resolved.
	// The type name should be emitted using a safe fallback.
	TypeReferenceSerializationKindUnknown = iota

	// The TypeReference resolves to a type with a constructor
	// function that can be reached at runtime (e.g. a `class`
	// declaration or a `var` declaration for the static side
	// of a type, such as the global `Promise` type in lib.d.ts).
	TypeReferenceSerializationKindTypeWithConstructSignatureAndValue

	// The TypeReference resolves to a Void-like, Nullable, or Never type.
	TypeReferenceSerializationKindVoidNullableOrNeverType

	// The TypeReference resolves to a Number-like type.
	TypeReferenceSerializationKindNumberLikeType

	// The TypeReference resolves to a BigInt-like type.
	TypeReferenceSerializationKindBigIntLikeType

	// The TypeReference resolves to a String-like type.
	TypeReferenceSerializationKindStringLikeType

	// The TypeReference resolves to a Boolean-like type.
	TypeReferenceSerializationKindBooleanType

	// The TypeReference resolves to an Array-like type.
	TypeReferenceSerializationKindArrayLikeType

	// The TypeReference resolves to the ESSymbol type.
	TypeReferenceSerializationKindESSymbolType

	// The TypeReference resolved to the global Promise constructor symbol.
	TypeReferenceSerializationKindPromise

	// The TypeReference resolves to a Function type or a type with call signatures.
	TypeReferenceSerializationKindTypeWithCallSignature

	// The TypeReference resolves to any other type.
	TypeReferenceSerializationKindObjectType
)

type EmitResolver interface {
	binder.ReferenceResolver
	IsReferencedAliasDeclaration(node *ast.Node) bool
	IsValueAliasDeclaration(node *ast.Node) bool
	IsTopLevelValueImportEqualsWithEntityName(node *ast.Node) bool
	MarkLinkedReferencesRecursively(file *ast.SourceFile)
	GetExternalModuleFileFromDeclaration(node *ast.Node) *ast.SourceFile
	GetEffectiveDeclarationFlags(node *ast.Node, flags ast.ModifierFlags) ast.ModifierFlags
	GetResolutionModeOverride(node *ast.Node) core.ResolutionMode

	// decorator metadata
	GetTypeReferenceSerializationKind(name *ast.EntityName, serialScope *ast.Node) TypeReferenceSerializationKind

	// const enum inlining
	GetConstantValue(node *ast.Node) any

	// JSX Emit
	GetJsxFactoryEntity(location *ast.Node) *ast.Node
	GetJsxFragmentFactoryEntity(location *ast.Node) *ast.Node
	SetReferencedImportDeclaration(node *ast.IdentifierNode, ref *ast.Declaration) // for overriding the reference resolver behavior for generated identifiers

	// declaration emit checker functionality projections
	PrecalculateDeclarationEmitVisibility(file *ast.SourceFile)
	IsSymbolAccessible(symbol *ast.Symbol, enclosingDeclaration *ast.Node, meaning ast.SymbolFlags, shouldComputeAliasToMarkVisible bool) SymbolAccessibilityResult
	IsEntityNameVisible(entityName *ast.Node, enclosingDeclaration *ast.Node) SymbolAccessibilityResult // previously SymbolVisibilityResult in strada - ErrorModuleName never set
	IsExpandoFunctionDeclaration(node *ast.Node) bool
	IsExpandoFunctionDeclarationUnsafe(node *ast.Node) bool
	IsLiteralConstDeclaration(node *ast.Node) bool
	RequiresAddingImplicitUndefined(node *ast.Node, symbol *ast.Symbol, enclosingDeclaration *ast.Node) bool
	IsDeclarationVisible(node *ast.Node) bool
	IsImportRequiredByAugmentation(decl *ast.ImportDeclaration) bool
	IsDefinitelyReferenceToGlobalSymbolObject(node *ast.Node) bool
	IsImplementationOfOverload(node *ast.SignatureDeclaration) bool
	GetEnumMemberValue(node *ast.Node) evaluator.Result
	IsLateBound(node *ast.Node) bool
	IsOptionalParameter(node *ast.Node) bool

	// isolatedDeclarations-specific declaration emit
	GetPropertiesOfContainerFunction(node *ast.Node) []*ast.Symbol
	RequiresAddingImplicitUndefinedUnsafe(node *ast.Node, symbol *ast.Symbol, enclosingDeclaration *ast.Node) bool

	// Node construction for declaration emit
	CreateTypeOfDeclaration(emitContext *EmitContext, declaration *ast.Node, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node
	CreateReturnTypeOfSignatureDeclaration(emitContext *EmitContext, signatureDeclaration *ast.Node, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node
	CreateTypeParametersOfSignatureDeclaration(emitContext *EmitContext, signatureDeclaration *ast.Node, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) []*ast.Node
	CreateLiteralConstValue(emitContext *EmitContext, node *ast.Node, tracker nodebuilder.SymbolTracker) *ast.Node
	CreateTypeOfExpression(emitContext *EmitContext, expression *ast.Node, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node
	CreateLateBoundIndexSignatures(emitContext *EmitContext, container *ast.Node, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) []*ast.Node
}
