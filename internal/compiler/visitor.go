package compiler

// NodeVisitor

type NodeVisitor interface {
	VisitNode(node *Node) *Node
	VisitNodes(nodes []*Node) []*Node
	// Tokens
	VisitToken(node *Token) *Node
	// Literals
	VisitNumericLiteral(node *NumericLiteral) *Node
	VisitBigIntLiteral(node *BigIntLiteral) *Node
	VisitStringLiteral(node *StringLiteral) *Node
	VisitJsxText(node *JsxText) *Node
	VisitRegularExpressionLiteral(node *RegularExpressionLiteral) *Node
	VisitNoSubstitutionTemplateLiteral(node *NoSubstitutionTemplateLiteral) *Node
	// Pseudo-literals
	VisitTemplateHead(node *TemplateHead) *Node
	VisitTemplateMiddle(node *TemplateMiddle) *Node
	VisitTemplateTail(node *TemplateTail) *Node
	// Identifiers
	VisitIdentifier(node *Identifier) *Node
	VisitPrivateIdentifier(node *PrivateIdentifier) *Node
	// Names
	VisitQualifiedName(node *QualifiedName) *Node
	VisitComputedPropertyName(node *ComputedPropertyName) *Node
	// Lists
	VisitModifierList(node *ModifierList) *Node
	VisitTypeParameterList(node *TypeParameterList) *Node
	VisitTypeArgumentList(node *TypeArgumentList) *Node
	// Signature elements
	VisitTypeParameterDeclaration(node *TypeParameterDeclaration) *Node
	VisitParameterDeclaration(node *ParameterDeclaration) *Node
	VisitDecorator(node *Decorator) *Node
	// Type members
	VisitPropertySignatureDeclaration(node *PropertySignatureDeclaration) *Node
	VisitPropertyDeclaration(node *PropertyDeclaration) *Node
	VisitMethodSignatureDeclaration(node *MethodSignatureDeclaration) *Node
	VisitMethodDeclaration(node *MethodDeclaration) *Node
	VisitClassStaticBlockDeclaration(node *ClassStaticBlockDeclaration) *Node
	VisitConstructorDeclaration(node *ConstructorDeclaration) *Node
	VisitGetAccessorDeclaration(node *GetAccessorDeclaration) *Node
	VisitSetAccessorDeclaration(node *SetAccessorDeclaration) *Node
	VisitCallSignatureDeclaration(node *CallSignatureDeclaration) *Node
	VisitConstructSignatureDeclaration(node *ConstructSignatureDeclaration) *Node
	VisitIndexSignatureDeclaration(node *IndexSignatureDeclaration) *Node
	// Types
	VisitTypePredicateNode(node *TypePredicateNode) *Node
	VisitTypeReferenceNode(node *TypeReferenceNode) *Node
	VisitFunctionTypeNode(node *FunctionTypeNode) *Node
	VisitConstructorTypeNode(node *ConstructorTypeNode) *Node
	VisitTypeQueryNode(node *TypeQueryNode) *Node
	VisitTypeLiteralNode(node *TypeLiteralNode) *Node
	VisitArrayTypeNode(node *ArrayTypeNode) *Node
	VisitTupleTypeNode(node *TupleTypeNode) *Node
	VisitOptionalTypeNode(node *OptionalTypeNode) *Node
	VisitRestTypeNode(node *RestTypeNode) *Node
	VisitUnionTypeNode(node *UnionTypeNode) *Node
	VisitIntersectionTypeNode(node *IntersectionTypeNode) *Node
	VisitConditionalTypeNode(node *ConditionalTypeNode) *Node
	VisitInferTypeNode(node *InferTypeNode) *Node
	VisitParenthesizedTypeNode(node *ParenthesizedTypeNode) *Node
	VisitThisTypeNode(node *ThisTypeNode) *Node
	VisitTypeOperatorNode(node *TypeOperatorNode) *Node
	VisitIndexedAccessTypeNode(node *IndexedAccessTypeNode) *Node
	VisitMappedTypeNode(node *MappedTypeNode) *Node
	VisitLiteralTypeNode(node *LiteralTypeNode) *Node
	VisitNamedTupleMember(node *NamedTupleMember) *Node
	VisitTemplateLiteralTypeNode(node *TemplateLiteralTypeNode) *Node
	VisitTemplateLiteralTypeSpan(node *TemplateLiteralTypeSpan) *Node
	VisitImportTypeNode(node *ImportTypeNode) *Node
	// Binding patterns
	VisitObjectBindingPattern(node *BindingPattern) *Node
	VisitArrayBindingPattern(node *BindingPattern) *Node
	VisitBindingElement(node *BindingElement) *Node
	// Expression
	VisitArrayLiteralExpression(node *ArrayLiteralExpression) *Node
	VisitObjectLiteralExpression(node *ObjectLiteralExpression) *Node
	VisitPropertyAccessExpression(node *PropertyAccessExpression) *Node
	VisitElementAccessExpression(node *ElementAccessExpression) *Node
	VisitCallExpression(node *CallExpression) *Node
	VisitNewExpression(node *NewExpression) *Node
	VisitTaggedTemplateExpression(node *TaggedTemplateExpression) *Node
	VisitTypeAssertion(node *TypeAssertion) *Node
	VisitParenthesizedExpression(node *ParenthesizedExpression) *Node
	VisitFunctionExpression(node *FunctionExpression) *Node
	VisitArrowFunction(node *ArrowFunction) *Node
	VisitDeleteExpression(node *DeleteExpression) *Node
	VisitTypeOfExpression(node *TypeOfExpression) *Node
	VisitVoidExpression(node *VoidExpression) *Node
	VisitAwaitExpression(node *AwaitExpression) *Node
	VisitPrefixUnaryExpression(node *PrefixUnaryExpression) *Node
	VisitPostfixUnaryExpression(node *PostfixUnaryExpression) *Node
	VisitBinaryExpression(node *BinaryExpression) *Node
	VisitConditionalExpression(node *ConditionalExpression) *Node
	VisitTemplateExpression(node *TemplateExpression) *Node
	VisitYieldExpression(node *YieldExpression) *Node
	VisitSpreadElement(node *SpreadElement) *Node
	VisitClassExpression(node *ClassExpression) *Node
	VisitOmittedExpression(node *OmittedExpression) *Node
	VisitExpressionWithTypeArguments(node *ExpressionWithTypeArguments) *Node
	VisitAsExpression(node *AsExpression) *Node
	VisitNonNullExpression(node *NonNullExpression) *Node
	VisitMetaProperty(node *MetaProperty) *Node
	// VisitSyntheticExpression(node *SyntheticExpression) *Node
	VisitSatisfiesExpression(node *SatisfiesExpression) *Node
	// Misc
	VisitTemplateSpan(node *TemplateSpan) *Node
	VisitSemicolonClassElement(node *SemicolonClassElement) *Node
	// Element
	VisitBlock(node *Block) *Node
	VisitEmptyStatement(node *EmptyStatement) *Node
	VisitVariableStatement(node *VariableStatement) *Node
	VisitExpressionStatement(node *ExpressionStatement) *Node
	VisitIfStatement(node *IfStatement) *Node
	VisitDoStatement(node *DoStatement) *Node
	VisitWhileStatement(node *WhileStatement) *Node
	VisitForStatement(node *ForStatement) *Node
	VisitForInStatement(node *ForInOrOfStatement) *Node
	VisitForOfStatement(node *ForInOrOfStatement) *Node
	VisitContinueStatement(node *ContinueStatement) *Node
	VisitBreakStatement(node *BreakStatement) *Node
	VisitReturnStatement(node *ReturnStatement) *Node
	VisitWithStatement(node *WithStatement) *Node
	VisitSwitchStatement(node *SwitchStatement) *Node
	VisitLabeledStatement(node *LabeledStatement) *Node
	VisitThrowStatement(node *ThrowStatement) *Node
	VisitTryStatement(node *TryStatement) *Node
	VisitDebuggerStatement(node *DebuggerStatement) *Node
	VisitVariableDeclaration(node *VariableDeclaration) *Node
	VisitVariableDeclarationList(node *VariableDeclarationList) *Node
	VisitFunctionDeclaration(node *FunctionDeclaration) *Node
	VisitClassDeclaration(node *ClassDeclaration) *Node
	VisitInterfaceDeclaration(node *InterfaceDeclaration) *Node
	VisitTypeAliasDeclaration(node *TypeAliasDeclaration) *Node
	VisitEnumDeclaration(node *EnumDeclaration) *Node
	VisitModuleDeclaration(node *ModuleDeclaration) *Node
	VisitModuleBlock(node *ModuleBlock) *Node
	VisitCaseBlock(node *CaseBlock) *Node
	VisitNamespaceExportDeclaration(node *NamespaceExportDeclaration) *Node
	VisitImportEqualsDeclaration(node *ImportEqualsDeclaration) *Node
	VisitImportDeclaration(node *ImportDeclaration) *Node
	VisitImportClause(node *ImportClause) *Node
	VisitNamespaceImport(node *NamespaceImport) *Node
	VisitNamedImports(node *NamedImports) *Node
	VisitImportSpecifier(node *ImportSpecifier) *Node
	VisitExportAssignment(node *ExportAssignment) *Node
	VisitExportDeclaration(node *ExportDeclaration) *Node
	VisitNamedExports(node *NamedExports) *Node
	VisitNamespaceExport(node *NamespaceExport) *Node
	VisitExportSpecifier(node *ExportSpecifier) *Node
	VisitMissingDeclaration(node *MissingDeclaration) *Node
	// Module references
	VisitExternalModuleReference(node *ExternalModuleReference) *Node
	// JSX
	VisitJsxElement(node *JsxElement) *Node
	VisitJsxSelfClosingElement(node *JsxSelfClosingElement) *Node
	VisitJsxOpeningElement(node *JsxOpeningElement) *Node
	VisitJsxClosingElement(node *JsxClosingElement) *Node
	VisitJsxFragment(node *JsxFragment) *Node
	VisitJsxOpeningFragment(node *JsxOpeningFragment) *Node
	VisitJsxClosingFragment(node *JsxClosingFragment) *Node
	VisitJsxAttribute(node *JsxAttribute) *Node
	VisitJsxAttributes(node *JsxAttributes) *Node
	VisitJsxSpreadAttribute(node *JsxSpreadAttribute) *Node
	VisitJsxExpression(node *JsxExpression) *Node
	VisitJsxNamespacedName(node *JsxNamespacedName) *Node
	// Clauses
	VisitCaseClause(node *CaseOrDefaultClause) *Node
	VisitDefaultClause(node *CaseOrDefaultClause) *Node
	VisitHeritageClause(node *HeritageClause) *Node
	VisitCatchClause(node *CatchClause) *Node
	// Import attributes
	VisitImportAttributes(node *ImportAttributes) *Node
	VisitImportAttribute(node *ImportAttribute) *Node
	// Property assignments
	VisitPropertyAssignment(node *PropertyAssignment) *Node
	VisitShorthandPropertyAssignment(node *ShorthandPropertyAssignment) *Node
	VisitSpreadAssignment(node *SpreadAssignment) *Node
	// Enum
	VisitEnumMember(node *EnumMember) *Node
	// Top-level nodes
	VisitSourceFile(node *SourceFile) *Node
	// VisitBundle(node *Bundle) *Node
	// JSDoc nodes
	// VisitJSDocTypeExpression(node *JSDocTypeExpression) *Node
	// VisitJSDocNameReference(node *JSDocNameReference) *Node
	// VisitJSDocMemberName(node *JSDocMemberName) *Node
	// VisitJSDocAllType(node *JSDocAllType) *Node
	// VisitJSDocUnknownType(node *JSDocUnknownType) *Node
	VisitJSDocNullableType(node *JSDocNullableType) *Node
	VisitJSDocNonNullableType(node *JSDocNonNullableType) *Node
	// VisitJSDocOptionalType(node *JSDocOptionalType) *Node
	// VisitJSDocFunctionType(node *JSDocFunctionType) *Node
	// VisitJSDocVariadicType(node *JSDocVariadicType) *Node
	// VisitJSDocNamepathType(node *JSDocNamepathType) *Node
	// VisitJSDoc(node *JSDoc) *Node
	// VisitJSDocText(node *JSDocText) *Node
	// VisitJSDocTypeLiteral(node *JSDocTypeLiteral) *Node
	// VisitJSDocSignature(node *JSDocSignature) *Node
	// VisitJSDocLink(node *JSDocLink) *Node
	// VisitJSDocLinkCode(node *JSDocLinkCode) *Node
	// VisitJSDocLinkPlain(node *JSDocLinkPlain) *Node
	// VisitJSDocTag(node *JSDocTag) *Node
	// VisitJSDocAugmentsTag(node *JSDocAugmentsTag) *Node
	// VisitJSDocImplementsTag(node *JSDocImplementsTag) *Node
	// VisitJSDocAuthorTag(node *JSDocAuthorTag) *Node
	// VisitJSDocDeprecatedTag(node *JSDocDeprecatedTag) *Node
	// VisitJSDocImmediateTag(node *JSDocImmediateTag) *Node
	// VisitJSDocClassTag(node *JSDocClassTag) *Node
	// VisitJSDocPublicTag(node *JSDocPublicTag) *Node
	// VisitJSDocPrivateTag(node *JSDocPrivateTag) *Node
	// VisitJSDocProtectedTag(node *JSDocProtectedTag) *Node
	// VisitJSDocReadonlyTag(node *JSDocReadonlyTag) *Node
	// VisitJSDocOverrideTag(node *JSDocOverrideTag) *Node
	// VisitJSDocCallbackTag(node *JSDocCallbackTag) *Node
	// VisitJSDocOverloadTag(node *JSDocOverloadTag) *Node
	// VisitJSDocEnumTag(node *JSDocEnumTag) *Node
	// VisitJSDocParameterTag(node *JSDocParameterTag) *Node
	// VisitJSDocReturnTag(node *JSDocReturnTag) *Node
	// VisitJSDocThisTag(node *JSDocThisTag) *Node
	// VisitJSDocTypeTag(node *JSDocTypeTag) *Node
	// VisitJSDocTemplateTag(node *JSDocTemplateTag) *Node
	// VisitJSDocTypedefTag(node *JSDocTypedefTag) *Node
	// VisitJSDocSeeTag(node *JSDocSeeTag) *Node
	// VisitJSDocPropertyTag(node *JSDocPropertyTag) *Node
	// VisitJSDocThrowsTag(node *JSDocThrowsTag) *Node
	// VisitJSDocSatisfiesTag(node *JSDocSatisfiesTag) *Node
	// VisitJSDocImportTag(node *JSDocImportTag) *Node
	// // Transformation nodes
	// VisitNotEmittedStatement(node *NotEmittedStatement) *Node
	// VisitPartiallyEmittedExpression(node *PartiallyEmittedExpression) *Node
	// VisitCommaListExpression(node *CommaListExpression) *Node
	// VisitSyntheticReferenceExpression(node *SyntheticReferenceExpression) *Node

	// Unhandled
	VisitOther(node *Node) *Node
}

type NodeVisitorBase struct {
	factory NodeFactory
}

func extractSingleNode(nodes []*Node) *Node {
	_assert(len(nodes) == 1, "Expected only a single node to be written to output")
	return nodes[0]
}

func (v *NodeVisitorBase) VisitNode(node *Node) *Node {
	return v.VisitAndLift(node, extractSingleNode)
}

func (v *NodeVisitorBase) VisitAndLift(node *Node, lift func(nodes []*Node) *Node) *Node {
	if node == nil {
		return nil
	}
	visited := node.Accept(v)
	if visited == nil {
		return nil
	}
	if visited.kind == SyntaxKindSyntaxList {
		visited = lift(visited.AsSyntaxList().children)
		_assert(visited.kind != SyntaxKindSyntaxList, "The result of visiting and lifting a Node may not be SyntaxList")
	}
	return visited
}

func (v *NodeVisitorBase) VisitNodes(nodes []*Node) []*Node {
	if nodes == nil {
		return nil
	}

	var updated []*Node = nil
	for i, node := range nodes {
		visited := node.Accept(v)
		if updated != nil || visited == nil || visited != node {
			if updated == nil {
				updated = make([]*Node, i)
				copy(updated, nodes[:i])
			}
			switch {
			case visited == nil: // do nothing
			case visited.kind == SyntaxKindSyntaxList:
				updated = append(updated, visited.AsSyntaxList().children...)
			default:
				updated = append(updated, visited)
			}
		}
	}
	if updated != nil {
		return updated
	}
	return nodes
}

// Tokens
func (v *NodeVisitorBase) VisitToken(node *Token) *Node {
	return node.AsNode()
}

// Literals
func (v *NodeVisitorBase) VisitNumericLiteral(node *NumericLiteral) *Node {
	return node.AsNode()
}
func (v *NodeVisitorBase) VisitBigIntLiteral(node *BigIntLiteral) *Node {
	return node.AsNode()
}
func (v *NodeVisitorBase) VisitStringLiteral(node *StringLiteral) *Node {
	return node.AsNode()
}
func (v *NodeVisitorBase) VisitJsxText(node *JsxText) *Node {
	return node.AsNode()
}
func (v *NodeVisitorBase) VisitRegularExpressionLiteral(node *RegularExpressionLiteral) *Node {
	return node.AsNode()
}
func (v *NodeVisitorBase) VisitNoSubstitutionTemplateLiteral(node *NoSubstitutionTemplateLiteral) *Node {
	return node.AsNode()
}

// Pseudo-literals
func (v *NodeVisitorBase) VisitTemplateHead(node *TemplateHead) *Node {
	return node.AsNode()
}
func (v *NodeVisitorBase) VisitTemplateMiddle(node *TemplateMiddle) *Node {
	return node.AsNode()
}
func (v *NodeVisitorBase) VisitTemplateTail(node *TemplateTail) *Node {
	return node.AsNode()
}

// Identifiers
func (v *NodeVisitorBase) VisitIdentifier(node *Identifier) *Node {
	return node.AsNode()
}
func (v *NodeVisitorBase) VisitPrivateIdentifier(node *PrivateIdentifier) *Node {
	return node.AsNode()
}

// Names
func (v *NodeVisitorBase) VisitQualifiedName(node *QualifiedName) *Node {
	return v.factory.UpdateQualifiedName(
		node.AsNode(),
		v.VisitNode(node.left),
		v.VisitNode(node.right),
	)
}
func (v *NodeVisitorBase) VisitComputedPropertyName(node *ComputedPropertyName) *Node {
	return v.factory.UpdateComputedPropertyName(
		node.AsNode(),
		v.VisitNode(node.expression),
	)
}

// Lists
func (v *NodeVisitorBase) VisitModifierList(node *ModifierList) *Node {
	return v.factory.UpdateModifierList(
		node.AsNode(),
		v.VisitNodes(node.modifiers),
	)
}
func (v *NodeVisitorBase) VisitTypeParameterList(node *TypeParameterList) *Node {
	return v.factory.UpdateTypeParameterList(
		node.AsNode(),
		v.VisitNodes(node.parameters),
	)
}
func (v *NodeVisitorBase) VisitTypeArgumentList(node *TypeArgumentList) *Node {
	return v.factory.UpdateTypeArgumentList(
		node.AsNode(),
		v.VisitNodes(node.arguments),
	)
}

// Signature elements
func (v *NodeVisitorBase) VisitTypeParameterDeclaration(node *TypeParameterDeclaration) *Node {
	return v.factory.UpdateTypeParameterDeclaration(
		node.AsNode(),
		v.VisitNode(node.modifiers),
		v.VisitNode(node.name),
		v.VisitNode(node.constraint),
		v.VisitNode(node.defaultType),
	)
}
func (v *NodeVisitorBase) VisitParameterDeclaration(node *ParameterDeclaration) *Node {
	return v.factory.UpdateParameterDeclaration(
		node.AsNode(),
		v.VisitNode(node.modifiers),
		v.VisitNode(node.dotDotDotToken),
		v.VisitNode(node.name),
		v.VisitNode(node.questionToken),
		v.VisitNode(node.typeNode),
		v.VisitNode(node.initializer),
	)
}
func (v *NodeVisitorBase) VisitDecorator(node *Decorator) *Node {
	return v.factory.UpdateDecorator(
		node.AsNode(),
		v.VisitNode(node.expression),
	)
}

// Type members
func (v *NodeVisitorBase) VisitPropertySignatureDeclaration(node *PropertySignatureDeclaration) *Node {
	return v.factory.UpdatePropertySignatureDeclaration(
		node.AsNode(),
		v.VisitNode(node.modifiers),
		v.VisitNode(node.name),
		v.VisitNode(node.postfixToken),
		v.VisitNode(node.typeNode),
		v.VisitNode(node.initializer),
	)
}
func (v *NodeVisitorBase) VisitPropertyDeclaration(node *PropertyDeclaration) *Node {
	return v.factory.UpdatePropertyDeclaration(
		node.AsNode(),
		v.VisitNode(node.modifiers),
		v.VisitNode(node.name),
		v.VisitNode(node.postfixToken),
		v.VisitNode(node.typeNode),
		v.VisitNode(node.initializer),
	)
}
func (v *NodeVisitorBase) VisitMethodSignatureDeclaration(node *MethodSignatureDeclaration) *Node {
	return v.factory.UpdateMethodSignatureDeclaration(
		node.AsNode(),
		v.VisitNode(node.modifiers),
		v.VisitNode(node.name),
		v.VisitNode(node.postfixToken),
		v.VisitNode(node.typeParameters),
		v.VisitNodes(node.parameters),
		v.VisitNode(node.returnType),
	)
}
func (v *NodeVisitorBase) VisitMethodDeclaration(node *MethodDeclaration) *Node {
	return v.factory.UpdateMethodDeclaration(
		node.AsNode(),
		v.VisitNode(node.modifiers),
		v.VisitNode(node.asteriskToken),
		v.VisitNode(node.name),
		v.VisitNode(node.postfixToken),
		v.VisitNode(node.typeParameters),
		v.VisitParameters(node.parameters),
		v.VisitNode(node.returnType),
		v.VisitFunctionBody(node.body),
	)
}
func (v *NodeVisitorBase) VisitClassStaticBlockDeclaration(node *ClassStaticBlockDeclaration) *Node {
	// A `static {}` Block does not have parameters, but we must still ensure we enter the lexical scope
	v.VisitParameters(nil)
	return v.factory.UpdateClassStaticBlockDeclaration(
		node.AsNode(),
		v.VisitNode(node.modifiers),
		v.VisitFunctionBody(node.body),
	)
}
func (v *NodeVisitorBase) VisitConstructorDeclaration(node *ConstructorDeclaration) *Node {
	return v.factory.UpdateConstructorDeclaration(
		node.AsNode(),
		v.VisitNode(node.modifiers),
		v.VisitNode(node.typeParameters),
		v.VisitParameters(node.parameters),
		v.VisitNode(node.returnType),
		v.VisitFunctionBody(node.body),
	)
}
func (v *NodeVisitorBase) VisitGetAccessorDeclaration(node *GetAccessorDeclaration) *Node {
	return v.factory.UpdateGetAccessorDeclaration(
		node.AsNode(),
		v.VisitNode(node.modifiers),
		v.VisitNode(node.name),
		v.VisitNode(node.typeParameters),
		v.VisitParameters(node.parameters),
		v.VisitNode(node.returnType),
		v.VisitFunctionBody(node.body),
	)
}
func (v *NodeVisitorBase) VisitSetAccessorDeclaration(node *SetAccessorDeclaration) *Node {
	return v.factory.UpdateSetAccessorDeclaration(
		node.AsNode(),
		v.VisitNode(node.modifiers),
		v.VisitNode(node.name),
		v.VisitNode(node.typeParameters),
		v.VisitParameters(node.parameters),
		v.VisitNode(node.returnType),
		v.VisitFunctionBody(node.body),
	)
}
func (v *NodeVisitorBase) VisitCallSignatureDeclaration(node *CallSignatureDeclaration) *Node {
	return v.factory.UpdateCallSignatureDeclaration(
		node.AsNode(),
		v.VisitNode(node.typeParameters),
		v.VisitNodes(node.parameters),
		v.VisitNode(node.returnType),
	)
}
func (v *NodeVisitorBase) VisitConstructSignatureDeclaration(node *ConstructSignatureDeclaration) *Node {
	return v.factory.UpdateConstructSignatureDeclaration(
		node.AsNode(),
		v.VisitNode(node.typeParameters),
		v.VisitNodes(node.parameters),
		v.VisitNode(node.returnType),
	)
}
func (v *NodeVisitorBase) VisitIndexSignatureDeclaration(node *IndexSignatureDeclaration) *Node {
	return v.factory.UpdateIndexSignatureDeclaration(
		node.AsNode(),
		v.VisitNode(node.typeParameters),
		v.VisitNodes(node.parameters),
		v.VisitNode(node.returnType),
	)
}

// Types
func (v *NodeVisitorBase) VisitTypePredicateNode(node *TypePredicateNode) *Node {
	return v.factory.UpdateTypePredicateNode(
		node.AsNode(),
		v.VisitNode(node.assertsModifier),
		v.VisitNode(node.parameterName),
		v.VisitNode(node.typeNode),
	)
}
func (v *NodeVisitorBase) VisitTypeReferenceNode(node *TypeReferenceNode) *Node {
	return v.factory.UpdateTypeReferenceNode(
		node.AsNode(),
		v.VisitNode(node.typeName),
		v.VisitNode(node.typeArguments),
	)
}
func (v *NodeVisitorBase) VisitFunctionTypeNode(node *FunctionTypeNode) *Node {
	return v.factory.UpdateFunctionTypeNode(
		node.AsNode(),
		v.VisitNode(node.typeParameters),
		v.VisitNodes(node.parameters),
		v.VisitNode(node.returnType),
	)
}
func (v *NodeVisitorBase) VisitConstructorTypeNode(node *ConstructorTypeNode) *Node {
	return v.factory.UpdateConstructorTypeNode(
		node.AsNode(),
		v.VisitNode(node.modifiers),
		v.VisitNode(node.typeParameters),
		v.VisitNodes(node.parameters),
		v.VisitNode(node.returnType),
	)
}
func (v *NodeVisitorBase) VisitTypeQueryNode(node *TypeQueryNode) *Node {
	return v.factory.UpdateTypeQueryNode(
		node.AsNode(),
		v.VisitNode(node.exprName),
		v.VisitNode(node.typeArguments),
	)
}
func (v *NodeVisitorBase) VisitTypeLiteralNode(node *TypeLiteralNode) *Node {
	return v.factory.UpdateTypeLiteralNode(
		node.AsNode(),
		v.VisitNodes(node.members),
	)
}
func (v *NodeVisitorBase) VisitArrayTypeNode(node *ArrayTypeNode) *Node {
	return v.factory.UpdateArrayTypeNode(
		node.AsNode(),
		v.VisitNode(node.elementType),
	)
}
func (v *NodeVisitorBase) VisitTupleTypeNode(node *TupleTypeNode) *Node {
	return v.factory.UpdateTupleTypeNode(
		node.AsNode(),
		v.VisitNodes(node.elements),
	)
}
func (v *NodeVisitorBase) VisitOptionalTypeNode(node *OptionalTypeNode) *Node {
	return v.factory.UpdateOptionalTypeNode(
		node.AsNode(),
		v.VisitNode(node.typeNode),
	)
}
func (v *NodeVisitorBase) VisitRestTypeNode(node *RestTypeNode) *Node {
	return v.factory.UpdateRestTypeNode(
		node.AsNode(),
		v.VisitNode(node.typeNode),
	)
}
func (v *NodeVisitorBase) VisitUnionTypeNode(node *UnionTypeNode) *Node {
	return v.factory.UpdateUnionTypeNode(
		node.AsNode(),
		v.VisitNodes(node.types),
	)
}
func (v *NodeVisitorBase) VisitIntersectionTypeNode(node *IntersectionTypeNode) *Node {
	return v.factory.UpdateIntersectionTypeNode(
		node.AsNode(),
		v.VisitNodes(node.types),
	)
}
func (v *NodeVisitorBase) VisitConditionalTypeNode(node *ConditionalTypeNode) *Node {
	return v.factory.UpdateConditionalTypeNode(
		node.AsNode(),
		v.VisitNode(node.checkType),
		v.VisitNode(node.extendsType),
		v.VisitNode(node.trueType),
		v.VisitNode(node.falseType),
	)
}
func (v *NodeVisitorBase) VisitInferTypeNode(node *InferTypeNode) *Node {
	return v.factory.UpdateInferTypeNode(
		node.AsNode(),
		v.VisitNode(node.typeParameter),
	)
}
func (v *NodeVisitorBase) VisitParenthesizedTypeNode(node *ParenthesizedTypeNode) *Node {
	return v.factory.UpdateParenthesizedTypeNode(
		node.AsNode(),
		v.VisitNode(node.typeNode),
	)
}
func (v *NodeVisitorBase) VisitThisTypeNode(node *ThisTypeNode) *Node {
	return node.AsNode()
}
func (v *NodeVisitorBase) VisitTypeOperatorNode(node *TypeOperatorNode) *Node {
	return v.factory.UpdateTypeOperatorNode(
		node.AsNode(),
		v.VisitNode(node.typeNode),
	)
}
func (v *NodeVisitorBase) VisitIndexedAccessTypeNode(node *IndexedAccessTypeNode) *Node {
	return v.factory.UpdateIndexedAccessTypeNode(
		node.AsNode(),
		v.VisitNode(node.objectType),
		v.VisitNode(node.indexType),
	)
}
func (v *NodeVisitorBase) VisitMappedTypeNode(node *MappedTypeNode) *Node {
	return v.factory.UpdateMappedTypeNode(
		node.AsNode(),
		v.VisitNode(node.readonlyToken),
		v.VisitNode(node.typeParameter),
		v.VisitNode(node.nameType),
		v.VisitNode(node.questionToken),
		v.VisitNode(node.typeNode),
		v.VisitNodes(node.members),
	)
}
func (v *NodeVisitorBase) VisitLiteralTypeNode(node *LiteralTypeNode) *Node {
	return v.factory.UpdateLiteralTypeNode(
		node.AsNode(),
		v.VisitNode(node.literal),
	)
}
func (v *NodeVisitorBase) VisitNamedTupleMember(node *NamedTupleMember) *Node {
	return v.factory.UpdateNamedTupleMember(
		node.AsNode(),
		v.VisitNode(node.dotDotDotToken),
		v.VisitNode(node.name),
		v.VisitNode(node.questionToken),
		v.VisitNode(node.typeNode),
	)
}
func (v *NodeVisitorBase) VisitTemplateLiteralTypeNode(node *TemplateLiteralTypeNode) *Node {
	return v.factory.UpdateTemplateLiteralTypeNode(
		node.AsNode(),
		v.VisitNode(node.head),
		v.VisitNodes(node.templateSpans),
	)
}
func (v *NodeVisitorBase) VisitTemplateLiteralTypeSpan(node *TemplateLiteralTypeSpan) *Node {
	return v.factory.UpdateTemplateLiteralTypeSpan(
		node.AsNode(),
		v.VisitNode(node.typeNode),
		v.VisitNode(node.literal),
	)
}
func (v *NodeVisitorBase) VisitImportTypeNode(node *ImportTypeNode) *Node {
	return v.factory.UpdateImportTypeNode(
		node.AsNode(),
		node.isTypeOf,
		v.VisitNode(node.argument),
		v.VisitNode(node.attributes),
		v.VisitNode(node.qualifier),
		v.VisitNode(node.typeArguments),
	)
}

// Binding patterns
func (v *NodeVisitorBase) VisitObjectBindingPattern(node *BindingPattern) *Node {
	return v.factory.UpdateBindingPattern(
		node.AsNode(),
		v.VisitNodes(node.elements),
	)
}
func (v *NodeVisitorBase) VisitArrayBindingPattern(node *BindingPattern) *Node {
	return v.factory.UpdateBindingPattern(
		node.AsNode(),
		v.VisitNodes(node.elements),
	)
}
func (v *NodeVisitorBase) VisitBindingElement(node *BindingElement) *Node {
	return v.factory.UpdateBindingElement(
		node.AsNode(),
		v.VisitNode(node.dotDotDotToken),
		v.VisitNode(node.propertyName),
		v.VisitNode(node.name),
		v.VisitNode(node.initializer),
	)
}

// Expression
func (v *NodeVisitorBase) VisitArrayLiteralExpression(node *ArrayLiteralExpression) *Node {
	return v.factory.UpdateArrayLiteralExpression(
		node.AsNode(),
		v.VisitNodes(node.elements),
	)
}
func (v *NodeVisitorBase) VisitObjectLiteralExpression(node *ObjectLiteralExpression) *Node {
	return v.factory.UpdateObjectLiteralExpression(
		node.AsNode(),
		v.VisitNodes(node.properties),
	)
}
func (v *NodeVisitorBase) VisitPropertyAccessExpression(node *PropertyAccessExpression) *Node {
	return v.factory.UpdatePropertyAccessExpression(
		node.AsNode(),
		v.VisitNode(node.expression),
		v.VisitNode(node.questionDotToken),
		v.VisitNode(node.name),
	)
}
func (v *NodeVisitorBase) VisitElementAccessExpression(node *ElementAccessExpression) *Node {
	return v.factory.UpdateElementAccessExpression(
		node.AsNode(),
		v.VisitNode(node.expression),
		v.VisitNode(node.questionDotToken),
		v.VisitNode(node.argumentExpression),
	)
}
func (v *NodeVisitorBase) VisitCallExpression(node *CallExpression) *Node {
	return v.factory.UpdateCallExpression(
		node.AsNode(),
		v.VisitNode(node.expression),
		v.VisitNode(node.questionDotToken),
		v.VisitNode(node.typeArguments),
		v.VisitNodes(node.arguments),
	)
}
func (v *NodeVisitorBase) VisitNewExpression(node *NewExpression) *Node {
	return v.factory.UpdateNewExpression(
		node.AsNode(),
		v.VisitNode(node.expression),
		v.VisitNode(node.typeArguments),
		v.VisitNodes(node.arguments),
	)
}
func (v *NodeVisitorBase) VisitTaggedTemplateExpression(node *TaggedTemplateExpression) *Node {
	return v.factory.UpdateTaggedTemplateExpression(
		node.AsNode(),
		v.VisitNode(node.tag),
		v.VisitNode(node.questionDotToken),
		v.VisitNode(node.typeArguments),
		v.VisitNode(node.template),
	)
}
func (v *NodeVisitorBase) VisitTypeAssertion(node *TypeAssertion) *Node {
	return v.factory.UpdateTypeAssertion(
		node.AsNode(),
		v.VisitNode(node.typeNode),
		v.VisitNode(node.expression),
	)
}
func (v *NodeVisitorBase) VisitParenthesizedExpression(node *ParenthesizedExpression) *Node {
	return v.factory.UpdateParenthesizedExpression(
		node.AsNode(),
		v.VisitNode(node.expression),
	)
}
func (v *NodeVisitorBase) VisitFunctionExpression(node *FunctionExpression) *Node {
	return v.factory.UpdateFunctionExpression(
		node.AsNode(),
		v.VisitNode(node.modifiers),
		v.VisitNode(node.asteriskToken),
		v.VisitNode(node.name),
		v.VisitNode(node.typeParameters),
		v.VisitParameters(node.parameters),
		v.VisitNode(node.returnType),
		v.VisitFunctionBody(node.body),
	)
}
func (v *NodeVisitorBase) VisitArrowFunction(node *ArrowFunction) *Node {
	return v.factory.UpdateArrowFunction(
		node.AsNode(),
		v.VisitNode(node.modifiers),
		v.VisitNode(node.typeParameters),
		v.VisitParameters(node.parameters),
		v.VisitNode(node.returnType),
		v.VisitNode(node.equalsGreaterThanToken),
		v.VisitFunctionBody(node.body),
	)
}
func (v *NodeVisitorBase) VisitDeleteExpression(node *DeleteExpression) *Node {
	return v.factory.UpdateDeleteExpression(
		node.AsNode(),
		v.VisitNode(node.expression),
	)
}
func (v *NodeVisitorBase) VisitTypeOfExpression(node *TypeOfExpression) *Node {
	return v.factory.UpdateTypeOfExpression(
		node.AsNode(),
		v.VisitNode(node.expression),
	)
}
func (v *NodeVisitorBase) VisitVoidExpression(node *VoidExpression) *Node {
	return v.factory.UpdateVoidExpression(
		node.AsNode(),
		v.VisitNode(node.expression),
	)
}
func (v *NodeVisitorBase) VisitAwaitExpression(node *AwaitExpression) *Node {
	return v.factory.UpdateAwaitExpression(
		node.AsNode(),
		v.VisitNode(node.expression),
	)
}
func (v *NodeVisitorBase) VisitPrefixUnaryExpression(node *PrefixUnaryExpression) *Node {
	return v.factory.UpdatePrefixUnaryExpression(
		node.AsNode(),
		v.VisitNode(node.operand),
	)
}
func (v *NodeVisitorBase) VisitPostfixUnaryExpression(node *PostfixUnaryExpression) *Node {
	return v.factory.UpdatePostfixUnaryExpression(
		node.AsNode(),
		v.VisitNode(node.operand),
	)
}
func (v *NodeVisitorBase) VisitBinaryExpression(node *BinaryExpression) *Node {
	return v.factory.UpdateBinaryExpression(
		node.AsNode(),
		v.VisitNode(node.left),
		v.VisitNode(node.operatorToken),
		v.VisitNode(node.right),
	)
}
func (v *NodeVisitorBase) VisitConditionalExpression(node *ConditionalExpression) *Node {
	return v.factory.UpdateConditionalExpression(
		node.AsNode(),
		v.VisitNode(node.condition),
		v.VisitNode(node.questionToken),
		v.VisitNode(node.whenTrue),
		v.VisitNode(node.colonToken),
		v.VisitNode(node.whenFalse),
	)
}
func (v *NodeVisitorBase) VisitTemplateExpression(node *TemplateExpression) *Node {
	return v.factory.UpdateTemplateExpression(
		node.AsNode(),
		v.VisitNode(node.head),
		v.VisitNodes(node.templateSpans),
	)
}
func (v *NodeVisitorBase) VisitYieldExpression(node *YieldExpression) *Node {
	return v.factory.UpdateYieldExpression(
		node.AsNode(),
		v.VisitNode(node.asteriskToken),
		v.VisitNode(node.expression),
	)
}
func (v *NodeVisitorBase) VisitSpreadElement(node *SpreadElement) *Node {
	return v.factory.UpdateSpreadElement(
		node.AsNode(),
		v.VisitNode(node.expression),
	)
}
func (v *NodeVisitorBase) VisitClassExpression(node *ClassExpression) *Node {
	return v.factory.UpdateClassExpression(
		node.AsNode(),
		v.VisitNode(node.modifiers),
		v.VisitNode(node.name),
		v.VisitNode(node.typeParameters),
		v.VisitNodes(node.heritageClauses),
		v.VisitNodes(node.members),
	)
}
func (v *NodeVisitorBase) VisitOmittedExpression(node *OmittedExpression) *Node {
	return node.AsNode()
}
func (v *NodeVisitorBase) VisitExpressionWithTypeArguments(node *ExpressionWithTypeArguments) *Node {
	return v.factory.UpdateExpressionWithTypeArguments(
		node.AsNode(),
		v.VisitNode(node.expression),
		v.VisitNode(node.typeArguments),
	)
}
func (v *NodeVisitorBase) VisitAsExpression(node *AsExpression) *Node {
	return v.factory.UpdateAsExpression(
		node.AsNode(),
		v.VisitNode(node.expression),
		v.VisitNode(node.typeNode),
	)
}
func (v *NodeVisitorBase) VisitNonNullExpression(node *NonNullExpression) *Node {
	return v.factory.UpdateNonNullExpression(
		node.AsNode(),
		v.VisitNode(node.expression),
	)
}
func (v *NodeVisitorBase) VisitMetaProperty(node *MetaProperty) *Node {
	return v.factory.UpdateMetaProperty(
		node.AsNode(),
		v.VisitNode(node.name),
	)
}

// VisitSyntheticExpression(node *SyntheticExpression) *Node
func (v *NodeVisitorBase) VisitSatisfiesExpression(node *SatisfiesExpression) *Node {
	return v.factory.UpdateSatisfiesExpression(
		node.AsNode(),
		v.VisitNode(node.expression),
		v.VisitNode(node.typeNode),
	)
}

// Misc
func (v *NodeVisitorBase) VisitTemplateSpan(node *TemplateSpan) *Node {
	return v.factory.UpdateTemplateSpan(
		node.AsNode(),
		v.VisitNode(node.expression),
		v.VisitNode(node.literal),
	)
}
func (v *NodeVisitorBase) VisitSemicolonClassElement(node *SemicolonClassElement) *Node {
	return node.AsNode()
}

// Element
func (v *NodeVisitorBase) VisitBlock(node *Block) *Node {
	return v.factory.UpdateBlock(
		node.AsNode(),
		v.VisitNodes(node.statements),
	)
}
func (v *NodeVisitorBase) VisitEmptyStatement(node *EmptyStatement) *Node {
	return node.AsNode()
}
func (v *NodeVisitorBase) VisitVariableStatement(node *VariableStatement) *Node {
	return v.factory.UpdateVariableStatement(
		node.AsNode(),
		v.VisitNode(node.modifiers),
		v.VisitNode(node.declarationList),
	)
}
func (v *NodeVisitorBase) VisitExpressionStatement(node *ExpressionStatement) *Node {
	return v.factory.UpdateExpressionStatement(
		node.AsNode(),
		v.VisitNode(node.expression),
	)
}
func (v *NodeVisitorBase) VisitIfStatement(node *IfStatement) *Node {
	return v.factory.UpdateIfStatement(
		node.AsNode(),
		v.VisitNode(node.expression),
		v.VisitNode(node.thenStatement),
		v.VisitNode(node.elseStatement),
	)
}
func (v *NodeVisitorBase) VisitDoStatement(node *DoStatement) *Node {
	return v.factory.UpdateDoStatement(
		node.AsNode(),
		v.VisitIterationBody(node.statement),
		v.VisitNode(node.expression),
	)
}
func (v *NodeVisitorBase) VisitWhileStatement(node *WhileStatement) *Node {
	return v.factory.UpdateWhileStatement(
		node.AsNode(),
		v.VisitNode(node.expression),
		v.VisitIterationBody(node.statement),
	)
}
func (v *NodeVisitorBase) VisitForStatement(node *ForStatement) *Node {
	return v.factory.UpdateForStatement(
		node.AsNode(),
		v.VisitNode(node.initializer),
		v.VisitNode(node.condition),
		v.VisitNode(node.incrementor),
		v.VisitIterationBody(node.statement),
	)
}
func (v *NodeVisitorBase) VisitForInStatement(node *ForInOrOfStatement) *Node {
	return v.factory.UpdateForInOrOfStatement(
		node.AsNode(),
		nil,
		v.VisitNode(node.initializer),
		v.VisitNode(node.expression),
		v.VisitIterationBody(node.statement),
	)
}
func (v *NodeVisitorBase) VisitForOfStatement(node *ForInOrOfStatement) *Node {
	return v.factory.UpdateForInOrOfStatement(
		node.AsNode(),
		v.VisitNode(node.awaitModifier),
		v.VisitNode(node.initializer),
		v.VisitNode(node.expression),
		v.VisitIterationBody(node.statement),
	)
}
func (v *NodeVisitorBase) VisitContinueStatement(node *ContinueStatement) *Node {
	return v.factory.UpdateContinueStatement(
		node.AsNode(),
		v.VisitNode(node.label),
	)
}
func (v *NodeVisitorBase) VisitBreakStatement(node *BreakStatement) *Node {
	return v.factory.UpdateBreakStatement(
		node.AsNode(),
		v.VisitNode(node.label),
	)
}
func (v *NodeVisitorBase) VisitReturnStatement(node *ReturnStatement) *Node {
	return v.factory.UpdateReturnStatement(
		node.AsNode(),
		v.VisitNode(node.expression),
	)
}
func (v *NodeVisitorBase) VisitWithStatement(node *WithStatement) *Node {
	return v.factory.UpdateWithStatement(
		node.AsNode(),
		v.VisitNode(node.expression),
		v.VisitNode(node.statement),
	)
}
func (v *NodeVisitorBase) VisitSwitchStatement(node *SwitchStatement) *Node {
	return v.factory.UpdateSwitchStatement(
		node.AsNode(),
		v.VisitNode(node.expression),
		v.VisitNode(node.caseBlock),
	)
}
func (v *NodeVisitorBase) VisitLabeledStatement(node *LabeledStatement) *Node {
	return v.factory.UpdateLabeledStatement(
		node.AsNode(),
		v.VisitNode(node.label),
		v.VisitNode(node.statement),
	)
}
func (v *NodeVisitorBase) VisitThrowStatement(node *ThrowStatement) *Node {
	return v.factory.UpdateThrowStatement(
		node.AsNode(),
		v.VisitNode(node.expression),
	)
}
func (v *NodeVisitorBase) VisitTryStatement(node *TryStatement) *Node {
	return v.factory.UpdateTryStatement(
		node.AsNode(),
		v.VisitNode(node.tryBlock),
		v.VisitNode(node.catchClause),
		v.VisitNode(node.finallyBlock),
	)
}
func (v *NodeVisitorBase) VisitDebuggerStatement(node *DebuggerStatement) *Node {
	return node.AsNode()
}
func (v *NodeVisitorBase) VisitVariableDeclaration(node *VariableDeclaration) *Node {
	return v.factory.UpdateVariableDeclaration(
		node.AsNode(),
		v.VisitNode(node.name),
		v.VisitNode(node.exclamationToken),
		v.VisitNode(node.typeNode),
		v.VisitNode(node.initializer),
	)
}
func (v *NodeVisitorBase) VisitVariableDeclarationList(node *VariableDeclarationList) *Node {
	return v.factory.UpdateVariableDeclarationList(
		node.AsNode(),
		v.VisitNodes(node.declarations),
	)
}
func (v *NodeVisitorBase) VisitFunctionDeclaration(node *FunctionDeclaration) *Node {
	return v.factory.UpdateFunctionDeclaration(
		node.AsNode(),
		v.VisitNode(node.modifiers),
		v.VisitNode(node.asteriskToken),
		v.VisitNode(node.name),
		v.VisitNode(node.typeParameters),
		v.VisitParameters(node.parameters),
		v.VisitNode(node.returnType),
		v.VisitFunctionBody(node.body),
	)
}
func (v *NodeVisitorBase) VisitClassDeclaration(node *ClassDeclaration) *Node {
	return v.factory.UpdateClassDeclaration(
		node.AsNode(),
		v.VisitNode(node.modifiers),
		v.VisitNode(node.name),
		v.VisitNode(node.typeParameters),
		v.VisitNodes(node.heritageClauses),
		v.VisitNodes(node.members),
	)
}
func (v *NodeVisitorBase) VisitInterfaceDeclaration(node *InterfaceDeclaration) *Node {
	return v.factory.UpdateInterfaceDeclaration(
		node.AsNode(),
		v.VisitNode(node.modifiers),
		v.VisitNode(node.name),
		v.VisitNode(node.typeParameters),
		v.VisitNodes(node.heritageClauses),
		v.VisitNodes(node.members),
	)
}
func (v *NodeVisitorBase) VisitTypeAliasDeclaration(node *TypeAliasDeclaration) *Node {
	return v.factory.UpdateTypeAliasDeclaration(
		node.AsNode(),
		v.VisitNode(node.modifiers),
		v.VisitNode(node.name),
		v.VisitNode(node.typeParameters),
		v.VisitNode(node.typeNode),
	)
}
func (v *NodeVisitorBase) VisitEnumDeclaration(node *EnumDeclaration) *Node {
	return v.factory.UpdateEnumDeclaration(
		node.AsNode(),
		v.VisitNode(node.modifiers),
		v.VisitNode(node.name),
		v.VisitNodes(node.members),
	)
}
func (v *NodeVisitorBase) VisitModuleDeclaration(node *ModuleDeclaration) *Node {
	return v.factory.UpdateModuleDeclaration(
		node.AsNode(),
		v.VisitNode(node.modifiers),
		v.VisitNode(node.name),
		v.VisitNode(node.body),
	)
}
func (v *NodeVisitorBase) VisitModuleBlock(node *ModuleBlock) *Node {
	return v.factory.UpdateModuleBlock(
		node.AsNode(),
		v.VisitNodes(node.statements),
	)
}
func (v *NodeVisitorBase) VisitCaseBlock(node *CaseBlock) *Node {
	return v.factory.UpdateCaseBlock(
		node.AsNode(),
		v.VisitNodes(node.clauses),
	)
}
func (v *NodeVisitorBase) VisitNamespaceExportDeclaration(node *NamespaceExportDeclaration) *Node {
	return v.factory.UpdateNamespaceExportDeclaration(
		node.AsNode(),
		v.VisitNode(node.modifiers),
		v.VisitNode(node.name),
	)
}
func (v *NodeVisitorBase) VisitImportEqualsDeclaration(node *ImportEqualsDeclaration) *Node {
	return v.factory.UpdateImportEqualsDeclaration(
		node.AsNode(),
		v.VisitNode(node.modifiers),
		node.isTypeOnly,
		v.VisitNode(node.name),
		v.VisitNode(node.moduleReference),
	)
}
func (v *NodeVisitorBase) VisitImportDeclaration(node *ImportDeclaration) *Node {
	return v.factory.UpdateImportDeclaration(
		node.AsNode(),
		v.VisitNode(node.modifiers),
		v.VisitNode(node.importClause),
		v.VisitNode(node.moduleSpecifier),
		v.VisitNode(node.attributes),
	)
}
func (v *NodeVisitorBase) VisitImportClause(node *ImportClause) *Node {
	return v.factory.UpdateImportClause(
		node.AsNode(),
		node.isTypeOnly,
		v.VisitNode(node.name),
		v.VisitNode(node.namedBindings),
	)
}
func (v *NodeVisitorBase) VisitNamespaceImport(node *NamespaceImport) *Node {
	return v.factory.UpdateNamespaceImport(
		node.AsNode(),
		v.VisitNode(node.name),
	)
}
func (v *NodeVisitorBase) VisitNamedImports(node *NamedImports) *Node {
	return v.factory.UpdateNamedImports(
		node.AsNode(),
		v.VisitNodes(node.elements),
	)
}
func (v *NodeVisitorBase) VisitImportSpecifier(node *ImportSpecifier) *Node {
	return v.factory.UpdateImportSpecifier(
		node.AsNode(),
		node.isTypeOnly,
		v.VisitNode(node.propertyName),
		v.VisitNode(node.name),
	)
}
func (v *NodeVisitorBase) VisitExportAssignment(node *ExportAssignment) *Node {
	return v.factory.UpdateExportAssignment(
		node.AsNode(),
		v.VisitNode(node.modifiers),
		v.VisitNode(node.expression),
	)
}
func (v *NodeVisitorBase) VisitExportDeclaration(node *ExportDeclaration) *Node {
	return v.factory.UpdateExportDeclaration(
		node.AsNode(),
		v.VisitNode(node.modifiers),
		node.isTypeOnly,
		v.VisitNode(node.exportClause),
		v.VisitNode(node.moduleSpecifier),
		v.VisitNode(node.attributes),
	)
}
func (v *NodeVisitorBase) VisitNamedExports(node *NamedExports) *Node {
	return v.factory.UpdateNamedExports(
		node.AsNode(),
		v.VisitNodes(node.elements),
	)
}
func (v *NodeVisitorBase) VisitNamespaceExport(node *NamespaceExport) *Node {
	return v.factory.UpdateNamespaceExport(
		node.AsNode(),
		v.VisitNode(node.name),
	)
}
func (v *NodeVisitorBase) VisitExportSpecifier(node *ExportSpecifier) *Node {
	return v.factory.UpdateExportSpecifier(
		node.AsNode(),
		node.isTypeOnly,
		v.VisitNode(node.propertyName),
		v.VisitNode(node.name),
	)
}
func (v *NodeVisitorBase) VisitMissingDeclaration(node *MissingDeclaration) *Node {
	return v.factory.UpdateMissingDeclaration(
		node.AsNode(),
		v.VisitNode(node.modifiers),
	)
}

// Module references
func (v *NodeVisitorBase) VisitExternalModuleReference(node *ExternalModuleReference) *Node {
	return v.factory.UpdateExternalModuleReference(
		node.AsNode(),
		v.VisitNode(node.expression),
	)
}

// JSX
func (v *NodeVisitorBase) VisitJsxElement(node *JsxElement) *Node {
	return v.factory.UpdateJsxElement(
		node.AsNode(),
		v.VisitNode(node.openingElement),
		v.VisitNodes(node.children),
		v.VisitNode(node.closingElement),
	)
}
func (v *NodeVisitorBase) VisitJsxSelfClosingElement(node *JsxSelfClosingElement) *Node {
	return v.factory.UpdateJsxSelfClosingElement(
		node.AsNode(),
		v.VisitNode(node.tagName),
		v.VisitNode(node.typeArguments),
		v.VisitNode(node.attributes),
	)
}
func (v *NodeVisitorBase) VisitJsxOpeningElement(node *JsxOpeningElement) *Node {
	return v.factory.UpdateJsxOpeningElement(
		node.AsNode(),
		v.VisitNode(node.tagName),
		v.VisitNode(node.typeArguments),
		v.VisitNode(node.attributes),
	)
}
func (v *NodeVisitorBase) VisitJsxClosingElement(node *JsxClosingElement) *Node {
	return v.factory.UpdateJsxClosingElement(
		node.AsNode(),
		v.VisitNode(node.tagName),
	)
}
func (v *NodeVisitorBase) VisitJsxFragment(node *JsxFragment) *Node {
	return v.factory.UpdateJsxFragment(
		node.AsNode(),
		v.VisitNode(node.openingFragment),
		v.VisitNodes(node.children),
		v.VisitNode(node.closingFragment),
	)
}
func (v *NodeVisitorBase) VisitJsxOpeningFragment(node *JsxOpeningFragment) *Node {
	return node.AsNode()
}
func (v *NodeVisitorBase) VisitJsxClosingFragment(node *JsxClosingFragment) *Node {
	return node.AsNode()
}
func (v *NodeVisitorBase) VisitJsxAttribute(node *JsxAttribute) *Node {
	return v.factory.UpdateJsxAttribute(
		node.AsNode(),
		v.VisitNode(node.name),
		v.VisitNode(node.initializer),
	)
}
func (v *NodeVisitorBase) VisitJsxAttributes(node *JsxAttributes) *Node {
	return v.factory.UpdateJsxAttributes(
		node.AsNode(),
		v.VisitNodes(node.properties),
	)
}
func (v *NodeVisitorBase) VisitJsxSpreadAttribute(node *JsxSpreadAttribute) *Node {
	return v.factory.UpdateJsxSpreadAttribute(
		node.AsNode(),
		v.VisitNode(node.expression),
	)
}
func (v *NodeVisitorBase) VisitJsxExpression(node *JsxExpression) *Node {
	return v.factory.UpdateJsxExpression(
		node.AsNode(),
		v.VisitNode(node.dotDotDotToken),
		v.VisitNode(node.expression),
	)
}
func (v *NodeVisitorBase) VisitJsxNamespacedName(node *JsxNamespacedName) *Node {
	return v.factory.UpdateJsxNamespacedName(
		node.AsNode(),
		v.VisitNode(node.name),
		v.VisitNode(node.namespace),
	)
}

// Clauses
func (v *NodeVisitorBase) VisitCaseClause(node *CaseOrDefaultClause) *Node {
	return v.factory.UpdateCaseOrDefaultClause(
		node.AsNode(),
		v.VisitNode(node.expression),
		v.VisitNodes(node.statements),
	)
}
func (v *NodeVisitorBase) VisitDefaultClause(node *CaseOrDefaultClause) *Node {
	return v.factory.UpdateCaseOrDefaultClause(
		node.AsNode(),
		nil,
		v.VisitNodes(node.statements),
	)
}
func (v *NodeVisitorBase) VisitHeritageClause(node *HeritageClause) *Node {
	return v.factory.UpdateHeritageClause(
		node.AsNode(),
		v.VisitNodes(node.types),
	)
}
func (v *NodeVisitorBase) VisitCatchClause(node *CatchClause) *Node {
	return v.factory.UpdateCatchClause(
		node.AsNode(),
		v.VisitNode(node.variableDeclaration),
		v.VisitNode(node.block),
	)
}

// Import attributes
func (v *NodeVisitorBase) VisitImportAttributes(node *ImportAttributes) *Node {
	return v.factory.UpdateImportAttributes(
		node.AsNode(),
		v.VisitNodes(node.attributes),
	)
}
func (v *NodeVisitorBase) VisitImportAttribute(node *ImportAttribute) *Node {
	return v.factory.UpdateImportAttribute(
		node.AsNode(),
		v.VisitNode(node.name),
		v.VisitNode(node.value),
	)
}

// Property assignments
func (v *NodeVisitorBase) VisitPropertyAssignment(node *PropertyAssignment) *Node {
	return v.factory.UpdatePropertyAssignment(
		node.AsNode(),
		v.VisitNode(node.modifiers),
		v.VisitNode(node.name),
		v.VisitNode(node.postfixToken),
		v.VisitNode(node.initializer),
	)
}
func (v *NodeVisitorBase) VisitShorthandPropertyAssignment(node *ShorthandPropertyAssignment) *Node {
	return v.factory.UpdateShorthandPropertyAssignment(
		node.AsNode(),
		v.VisitNode(node.modifiers),
		v.VisitNode(node.name),
		v.VisitNode(node.postfixToken),
		v.VisitNode(node.objectAssignmentInitializer),
	)
}
func (v *NodeVisitorBase) VisitSpreadAssignment(node *SpreadAssignment) *Node {
	return v.factory.UpdateSpreadAssignment(
		node.AsNode(),
		v.VisitNode(node.expression),
	)
}

// Enum
func (v *NodeVisitorBase) VisitEnumMember(node *EnumMember) *Node {
	return v.factory.UpdateEnumMember(
		node.AsNode(),
		v.VisitNode(node.name),
		v.VisitNode(node.initializer),
	)
}

// Top-level nodes
func (v *NodeVisitorBase) VisitSourceFile(node *SourceFile) *Node {
	return v.factory.UpdateSourceFile(
		node.AsNode(),
		v.VisitNodes(node.statements),
	)
}

// JSDoc nodes
func (v *NodeVisitorBase) VisitJSDocNullableType(node *JSDocNullableType) *Node {
	return v.factory.UpdateJSDocNullableType(
		node.AsNode(),
		v.VisitNode(node.typeNode),
	)
}
func (v *NodeVisitorBase) VisitJSDocNonNullableType(node *JSDocNonNullableType) *Node {
	return v.factory.UpdateJSDocNonNullableType(
		node.AsNode(),
		v.VisitNode(node.typeNode),
	)
}

func (v *NodeVisitorBase) VisitOther(node *Node) *Node {
	return node
}

// Starts a new lexical environment and visits a parameter list, suspending the lexical
// environment upon completion.
func (v *NodeVisitorBase) VisitParameters(nodes []*ParameterDeclarationNode) []*ParameterDeclarationNode {
	// TODO(rbuckton): To be implemented
	return v.VisitNodes(nodes)
}

// Resumes a suspended lexical environment and visits a function body, ending the lexical
// environment and merging hoisted declarations upon completion.
func (v *NodeVisitorBase) VisitFunctionBody(node *BlockOrExpression) *Node {
	// TODO(rbuckton): To be implemented
	return v.VisitNode(node)
}

// Visits an iteration body, adding any block-scoped variables required by the transformation.
func (v *NodeVisitorBase) VisitIterationBody(body *Statement) *Statement {
	// TODO(rbuckton): To be implemented
	return v.VisitNode(body)
}
