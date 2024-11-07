package compiler

import "slices"

// Visitor

type Visitor func(*Node) bool

func visit(v Visitor, node *Node) bool {
	if node != nil {
		return v(node)
	}
	return false
}

func visitNodes(v Visitor, nodes []*Node) bool {
	for _, node := range nodes {
		if v(node) {
			return true
		}
	}
	return false
}

// NodeTest

type NodeTest func(*Node) bool

// NodeFactory

type NodeFactory struct {
	parenthesizerRules ParenthesizerRules
	identifierPool     Pool[Identifier]
}

func (f *NodeFactory) ParenthesizerRules() ParenthesizerRules {
	if f.parenthesizerRules == nil {
		f.parenthesizerRules = NewParenthesizerRules(f)
	}
	return f.parenthesizerRules
}

func (f *NodeFactory) NewNode(kind SyntaxKind, data NodeData) *Node {
	n := data.AsNode()
	n.kind = kind
	n.data = data
	return n
}

func (f *NodeFactory) UpdateNode(updated *Node, original *Node) *Node {
	if updated != original {
		f.SetOriginalNode(updated, original)
		updated.loc = original.loc
	}
	return updated
}

func (f *NodeFactory) SetOriginalNode(node *Node, original *Node) *Node {
	if node.original != original {
		if node.original != nil {
			panic("original node already set")
		}
		node.original = original
		// TODO: copy emitNode data
	}
	return node
}

func (f *NodeFactory) updateOuterExpression(outerExpression *Expression, expression *Expression) *Expression {
	switch outerExpression.kind {
	case SyntaxKindParenthesizedExpression:
		return f.UpdateParenthesizedExpression(outerExpression, expression)
	case SyntaxKindTypeAssertionExpression:
		return f.UpdateTypeAssertion(outerExpression, outerExpression.AsTypeAssertion().typeNode, expression)
	case SyntaxKindAsExpression:
		return f.UpdateAsExpression(outerExpression, expression, outerExpression.AsAsExpression().typeNode)
	case SyntaxKindSatisfiesExpression:
		return f.UpdateSatisfiesExpression(outerExpression, expression, outerExpression.AsSatisfiesExpression().typeNode)
	case SyntaxKindNonNullExpression:
		return f.UpdateNonNullExpression(outerExpression, expression)
	case SyntaxKindExpressionWithTypeArguments:
		return f.UpdateExpressionWithTypeArguments(outerExpression, expression, outerExpression.AsExpressionWithTypeArguments().typeArguments)
	// TODO(rbuckton): Implement PartiallyEmittedExpression if needed
	// case SyntaxKindPartiallyEmittedExpression:
	// 	return f.UpdatePartiallyEmittedExpression(outerExpression, expression)
	default:
		panic("Not a supported outer expression")
	}
}

/**
 * Determines whether a node is a parenthesized expression that can be ignored when recreating outer expressions.
 *
 * A parenthesized expression can be ignored when all of the following are true:
 *
 * - It's `pos` and `end` are not -1
 * - It does not have a custom source map range
 * - It does not have a custom comment range
 * - It does not have synthetic leading or trailing comments
 *
 * If an outermost parenthesized expression is ignored, but the containing expression requires a parentheses around
 * the expression to maintain precedence, a new parenthesized expression should be created automatically when
 * the containing expression is created/updated.
 */
func (f *NodeFactory) isIgnorableParen(node *Expression) bool {
	// TODO(rbuckton): EmitNode-related functions are not yet implemented
	return isParenthesizedExpression(node) &&
		nodeIsSynthesized(node) // &&
	//	nodeIsSynthesized(getSourceMapRange(node)) &&
	//	nodeIsSynthesized(getCommentRange(node)) &&
	//	!some(getSyntheticLeadingComments(node)) &&
	//	!some(getSyntheticTrailingComments(node))
}

func (f *NodeFactory) RestoreOuterExpressions(outerExpression *Expression, innerExpression *Expression, kinds OuterExpressionKinds) *Expression {
	if outerExpression != nil && isOuterExpression(outerExpression, kinds) && !f.isIgnorableParen(outerExpression) {
		return f.updateOuterExpression(
			outerExpression,
			f.RestoreOuterExpressions(outerExpression.Expression(), innerExpression, OEKAll))
	}
	return innerExpression
}

// AST Node
// Interface values stored in AST nodes are never typed nil values. Construction code must ensure that
// interface valued properties either store a true nil or a reference to a non-nil struct.

type Node struct {
	kind     SyntaxKind
	flags    NodeFlags
	loc      TextRange
	id       NodeId
	parent   *Node
	original *Node
	data     NodeData
}

// Node accessors. Some accessors are implemented as methods on NodeData, others are implemented though
// type switches. Either approach is fine. Interface methods are likely more performant, but have higher
// code size costs because we have hundreds of implementations of the NodeData interface.

func (n *Node) Pos() int                                  { return n.loc.Pos() }
func (n *Node) End() int                                  { return n.loc.End() }
func (n *Node) ForEachChild(v Visitor) bool               { return n.data.ForEachChild(v) }
func (n *Node) Accept(v NodeVisitor) *Node                { return n.data.Accept(v) }
func (n *Node) Name() *DeclarationName                    { return n.data.Name() }
func (n *Node) Modifiers() *ModifierListNode              { return n.data.Modifiers() }
func (n *Node) TypeParameters() *TypeParameterListNode    { return n.data.TypeParameters() }
func (n *Node) FlowNodeData() *FlowNodeBase               { return n.data.FlowNodeData() }
func (n *Node) DeclarationData() *DeclarationBase         { return n.data.DeclarationData() }
func (n *Node) Symbol() *Symbol                           { return n.data.DeclarationData().symbol }
func (n *Node) ExportableData() *ExportableBase           { return n.data.ExportableData() }
func (n *Node) LocalSymbol() *Symbol                      { return n.data.ExportableData().localSymbol }
func (n *Node) LocalsContainerData() *LocalsContainerBase { return n.data.LocalsContainerData() }
func (n *Node) Locals() SymbolTable                       { return n.data.LocalsContainerData().locals }
func (n *Node) FunctionLikeData() *FunctionLikeBase       { return n.data.FunctionLikeData() }
func (n *Node) Parameters() []*ParameterDeclarationNode   { return n.data.FunctionLikeData().parameters }
func (n *Node) ReturnType() *TypeNode                     { return n.data.FunctionLikeData().returnType }
func (n *Node) ClassLikeData() *ClassLikeBase             { return n.data.ClassLikeData() }
func (n *Node) BodyData() *BodyBase                       { return n.data.BodyData() }

func (n *Node) Text() string {
	switch n.kind {
	case SyntaxKindIdentifier:
		return n.AsIdentifier().text
	case SyntaxKindPrivateIdentifier:
		return n.AsPrivateIdentifier().text
	case SyntaxKindStringLiteral:
		return n.AsStringLiteral().text
	case SyntaxKindNumericLiteral:
		return n.AsNumericLiteral().text
	case SyntaxKindBigIntLiteral:
		return n.AsBigIntLiteral().text
	case SyntaxKindNoSubstitutionTemplateLiteral:
		return n.AsNoSubstitutionTemplateLiteral().text
	case SyntaxKindJsxNamespacedName:
		return n.AsJsxNamespacedName().namespace.Text() + ":" + n.AsJsxNamespacedName().name.Text()
	}
	panic("Unhandled case in Node.Text")
}

func (node *Node) Expression() *Node {
	switch node.kind {
	case SyntaxKindPropertyAccessExpression:
		return node.AsPropertyAccessExpression().expression
	case SyntaxKindElementAccessExpression:
		return node.AsElementAccessExpression().expression
	case SyntaxKindParenthesizedExpression:
		return node.AsParenthesizedExpression().expression
	case SyntaxKindCallExpression:
		return node.AsCallExpression().expression
	case SyntaxKindNewExpression:
		return node.AsNewExpression().expression
	case SyntaxKindExpressionWithTypeArguments:
		return node.AsExpressionWithTypeArguments().expression
	case SyntaxKindNonNullExpression:
		return node.AsNonNullExpression().expression
	case SyntaxKindTypeAssertionExpression:
		return node.AsTypeAssertion().expression
	case SyntaxKindAsExpression:
		return node.AsAsExpression().expression
	case SyntaxKindSatisfiesExpression:
		return node.AsSatisfiesExpression().expression
	}
	panic("Unhandled case in Node.Expression")
}

func (node *Node) Arguments() []*Node {
	switch node.kind {
	case SyntaxKindCallExpression:
		return node.AsCallExpression().arguments
	case SyntaxKindNewExpression:
		return node.AsNewExpression().arguments
	}
	panic("Unhandled case in Node.Arguments")
}

// Node casts

func (n *Node) AsIdentifier() *Identifier {
	return n.data.(*Identifier)
}
func (n *Node) AsPrivateIdentifier() *PrivateIdentifier {
	return n.data.(*PrivateIdentifier)
}
func (n *Node) AsQualifiedName() *QualifiedName {
	return n.data.(*QualifiedName)
}
func (n *Node) AsModifierList() *ModifierList {
	return n.data.(*ModifierList)
}
func (n *Node) AsSourceFile() *SourceFile {
	return n.data.(*SourceFile)
}
func (n *Node) AsPrefixUnaryExpression() *PrefixUnaryExpression {
	return n.data.(*PrefixUnaryExpression)
}
func (n *Node) AsPostfixUnaryExpression() *PostfixUnaryExpression {
	return n.data.(*PostfixUnaryExpression)
}
func (n *Node) AsYieldExpression() *YieldExpression {
	return n.data.(*YieldExpression)
}
func (n *Node) AsParenthesizedExpression() *ParenthesizedExpression {
	return n.data.(*ParenthesizedExpression)
}
func (n *Node) AsTypeAssertion() *TypeAssertion {
	return n.data.(*TypeAssertion)
}
func (n *Node) AsAsExpression() *AsExpression {
	return n.data.(*AsExpression)
}
func (n *Node) AsSatisfiesExpression() *SatisfiesExpression {
	return n.data.(*SatisfiesExpression)
}
func (n *Node) AsExpressionWithTypeArguments() *ExpressionWithTypeArguments {
	return n.data.(*ExpressionWithTypeArguments)
}
func (n *Node) AsNonNullExpression() *NonNullExpression {
	return n.data.(*NonNullExpression)
}
func (n *Node) AsBindingElement() *BindingElement {
	return n.data.(*BindingElement)
}
func (n *Node) AsMissingDeclaration() *MissingDeclaration {
	return n.data.(*MissingDeclaration)
}
func (n *Node) AsImportSpecifier() *ImportSpecifier {
	return n.data.(*ImportSpecifier)
}
func (n *Node) AsArrowFunction() *ArrowFunction {
	return n.data.(*ArrowFunction)
}
func (n *Node) AsCallExpression() *CallExpression {
	return n.data.(*CallExpression)
}
func (n *Node) AsPropertyAccessExpression() *PropertyAccessExpression {
	return n.data.(*PropertyAccessExpression)
}
func (n *Node) AsElementAccessExpression() *ElementAccessExpression {
	return n.data.(*ElementAccessExpression)
}
func (n *Node) AsComputedPropertyName() *ComputedPropertyName {
	return n.data.(*ComputedPropertyName)
}
func (n *Node) AsDecorator() *Decorator {
	return n.data.(*Decorator)
}
func (n *Node) AsBinaryExpression() *BinaryExpression {
	return n.data.(*BinaryExpression)
}
func (n *Node) AsModuleDeclaration() *ModuleDeclaration {
	return n.data.(*ModuleDeclaration)
}
func (n *Node) AsStringLiteral() *StringLiteral {
	return n.data.(*StringLiteral)
}
func (n *Node) AsNumericLiteral() *NumericLiteral {
	return n.data.(*NumericLiteral)
}
func (n *Node) AsBigIntLiteral() *BigIntLiteral {
	return n.data.(*BigIntLiteral)
}
func (n *Node) AsNoSubstitutionTemplateLiteral() *NoSubstitutionTemplateLiteral {
	return n.data.(*NoSubstitutionTemplateLiteral)
}
func (n *Node) AsVariableDeclaration() *VariableDeclaration {
	return n.data.(*VariableDeclaration)
}
func (n *Node) AsExportAssignment() *ExportAssignment {
	return n.data.(*ExportAssignment)
}
func (n *Node) AsObjectLiteralExpression() *ObjectLiteralExpression {
	return n.data.(*ObjectLiteralExpression)
}
func (n *Node) AsIfStatement() *IfStatement {
	return n.data.(*IfStatement)
}
func (n *Node) AsWhileStatement() *WhileStatement {
	return n.data.(*WhileStatement)
}
func (n *Node) AsDoStatement() *DoStatement {
	return n.data.(*DoStatement)
}
func (n *Node) AsForStatement() *ForStatement {
	return n.data.(*ForStatement)
}
func (n *Node) AsConditionalExpression() *ConditionalExpression {
	return n.data.(*ConditionalExpression)
}
func (n *Node) AsForInOrOfStatement() *ForInOrOfStatement {
	return n.data.(*ForInOrOfStatement)
}
func (n *Node) AsShorthandPropertyAssignment() *ShorthandPropertyAssignment {
	return n.data.(*ShorthandPropertyAssignment)
}
func (n *Node) AsPropertyAssignment() *PropertyAssignment {
	return n.data.(*PropertyAssignment)
}
func (n *Node) AsExpressionStatement() *ExpressionStatement {
	return n.data.(*ExpressionStatement)
}
func (n *Node) AsBlock() *Block {
	return n.data.(*Block)
}
func (n *Node) AsModuleBlock() *ModuleBlock {
	return n.data.(*ModuleBlock)
}
func (n *Node) AsVariableStatement() *VariableStatement {
	return n.data.(*VariableStatement)
}
func (n *Node) AsVariableDeclarationList() *VariableDeclarationList {
	return n.data.(*VariableDeclarationList)
}
func (n *Node) AsMetaProperty() *MetaProperty {
	return n.data.(*MetaProperty)
}
func (n *Node) AsTypeReference() *TypeReferenceNode {
	return n.data.(*TypeReferenceNode)
}
func (n *Node) AsCallSignatureDeclaration() *CallSignatureDeclaration {
	return n.data.(*CallSignatureDeclaration)
}
func (n *Node) AsConstructSignatureDeclaration() *ConstructSignatureDeclaration {
	return n.data.(*ConstructSignatureDeclaration)
}
func (n *Node) AsConstructorDeclaration() *ConstructorDeclaration {
	return n.data.(*ConstructorDeclaration)
}
func (n *Node) AsConditionalTypeNode() *ConditionalTypeNode {
	return n.data.(*ConditionalTypeNode)
}
func (n *Node) AsClassExpression() *ClassExpression {
	return n.data.(*ClassExpression)
}
func (n *Node) AsHeritageClause() *HeritageClause {
	return n.data.(*HeritageClause)
}
func (n *Node) AsFunctionExpression() *FunctionExpression {
	return n.data.(*FunctionExpression)
}
func (n *Node) AsParameterDeclaration() *ParameterDeclaration {
	return n.data.(*ParameterDeclaration)
}
func (n *Node) AsInferTypeNode() *InferTypeNode {
	return n.data.(*InferTypeNode)
}
func (n *Node) AsTypeParameter() *TypeParameterDeclaration {
	return n.data.(*TypeParameterDeclaration)
}
func (n *Node) AsExportSpecifier() *ExportSpecifier {
	return n.data.(*ExportSpecifier)
}
func (n *Node) AsExportDeclaration() *ExportDeclaration {
	return n.data.(*ExportDeclaration)
}
func (n *Node) AsNamespaceExport() *NamespaceExport {
	return n.data.(*NamespaceExport)
}
func (n *Node) AsPropertyDeclaration() *PropertyDeclaration {
	return n.data.(*PropertyDeclaration)
}
func (n *Node) AsImportClause() *ImportClause {
	return n.data.(*ImportClause)
}
func (n *Node) AsImportEqualsDeclaration() *ImportEqualsDeclaration {
	return n.data.(*ImportEqualsDeclaration)
}
func (n *Node) AsNamespaceImport() *NamespaceImport {
	return n.data.(*NamespaceImport)
}
func (n *Node) AsNamedImports() *NamedImports {
	return n.data.(*NamedImports)
}
func (n *Node) AsPropertySignatureDeclaration() *PropertySignatureDeclaration {
	return n.data.(*PropertySignatureDeclaration)
}
func (n *Node) AsEnumMember() *EnumMember {
	return n.data.(*EnumMember)
}
func (n *Node) AsEnumDeclaration() *EnumDeclaration {
	return n.data.(*EnumDeclaration)
}
func (n *Node) AsReturnStatement() *ReturnStatement {
	return n.data.(*ReturnStatement)
}
func (n *Node) AsWithStatement() *WithStatement {
	return n.data.(*WithStatement)
}
func (n *Node) AsSwitchStatement() *SwitchStatement {
	return n.data.(*SwitchStatement)
}
func (n *Node) AsCaseOrDefaultClause() *CaseOrDefaultClause {
	return n.data.(*CaseOrDefaultClause)
}
func (n *Node) AsThrowStatement() *ThrowStatement {
	return n.data.(*ThrowStatement)
}
func (n *Node) AsTemplateExpression() *TemplateExpression {
	return n.data.(*TemplateExpression)
}
func (n *Node) AsTemplateSpan() *TemplateSpan {
	return n.data.(*TemplateSpan)
}
func (n *Node) AsImportAttribute() *ImportAttribute {
	return n.data.(*ImportAttribute)
}
func (n *Node) AsImportAttributes() *ImportAttributes {
	return n.data.(*ImportAttributes)
}
func (n *Node) AsImportTypeNode() *ImportTypeNode {
	return n.data.(*ImportTypeNode)
}
func (n *Node) AsNewExpression() *NewExpression {
	return n.data.(*NewExpression)
}
func (n *Node) AsTaggedTemplateExpression() *TaggedTemplateExpression {
	return n.data.(*TaggedTemplateExpression)
}
func (n *Node) AsTypeArgumentList() *TypeArgumentList {
	return n.data.(*TypeArgumentList)
}
func (n *Node) AsJsxOpeningElement() *JsxOpeningElement {
	return n.data.(*JsxOpeningElement)
}
func (n *Node) AsJsxSelfClosingElement() *JsxSelfClosingElement {
	return n.data.(*JsxSelfClosingElement)
}
func (n *Node) AsJsxClosingElement() *JsxClosingElement {
	return n.data.(*JsxClosingElement)
}
func (n *Node) AsImportDeclaration() *ImportDeclaration {
	return n.data.(*ImportDeclaration)
}
func (n *Node) AsExternalModuleReference() *ExternalModuleReference {
	return n.data.(*ExternalModuleReference)
}
func (n *Node) AsLiteralTypeNode() *LiteralTypeNode {
	return n.data.(*LiteralTypeNode)
}
func (n *Node) AsJsxNamespacedName() *JsxNamespacedName {
	return n.data.(*JsxNamespacedName)
}
func (n *Node) AsTypeParameterList() *TypeParameterList {
	return n.data.(*TypeParameterList)
}
func (n *Node) AsClassDeclaration() *ClassDeclaration {
	return n.data.(*ClassDeclaration)
}
func (n *Node) AsInterfaceDeclaration() *InterfaceDeclaration {
	return n.data.(*InterfaceDeclaration)
}
func (n *Node) AsTypeAliasDeclaration() *TypeAliasDeclaration {
	return n.data.(*TypeAliasDeclaration)
}
func (n *Node) AsJsxAttribute() *JsxAttribute {
	return n.data.(*JsxAttribute)
}
func (n *Node) AsParenthesizedTypeNode() *ParenthesizedTypeNode {
	return n.data.(*ParenthesizedTypeNode)
}
func (n *Node) AsFunctionTypeNode() *FunctionTypeNode {
	return n.data.(*FunctionTypeNode)
}
func (n *Node) AsConstructorTypeNode() *ConstructorTypeNode {
	return n.data.(*ConstructorTypeNode)
}
func (n *Node) AsTypePredicateNode() *TypePredicateNode {
	return n.data.(*TypePredicateNode)
}
func (n *Node) AsTypeOperatorNode() *TypeOperatorNode {
	return n.data.(*TypeOperatorNode)
}
func (n *Node) AsMappedTypeNode() *MappedTypeNode {
	return n.data.(*MappedTypeNode)
}
func (n *Node) AsTypeLiteralNode() *TypeLiteralNode {
	return n.data.(*TypeLiteralNode)
}
func (n *Node) AsArrayLiteralExpression() *ArrayLiteralExpression {
	return n.data.(*ArrayLiteralExpression)
}
func (n *Node) AsMethodDeclaration() *MethodDeclaration {
	return n.data.(*MethodDeclaration)
}
func (n *Node) AsMethodSignatureDeclaration() *MethodSignatureDeclaration {
	return n.data.(*MethodSignatureDeclaration)
}
func (n *Node) AsTemplateLiteralTypeNode() *TemplateLiteralTypeNode {
	return n.data.(*TemplateLiteralTypeNode)
}
func (n *Node) AsTemplateLiteralTypeSpan() *TemplateLiteralTypeSpan {
	return n.data.(*TemplateLiteralTypeSpan)
}
func (n *Node) AsJsxElement() *JsxElement {
	return n.data.(*JsxElement)
}
func (n *Node) AsJsxFragment() *JsxFragment {
	return n.data.(*JsxFragment)
}
func (n *Node) AsJsxAttributes() *JsxAttributes {
	return n.data.(*JsxAttributes)
}
func (n *Node) AsJsxSpreadAttribute() *JsxSpreadAttribute {
	return n.data.(*JsxSpreadAttribute)
}
func (n *Node) AsJsxExpression() *JsxExpression {
	return n.data.(*JsxExpression)
}
func (n *Node) AsJsxText() *JsxText {
	return n.data.(*JsxText)
}
func (n *Node) AsKeywordExpression() *KeywordExpression {
	return n.data.(*KeywordExpression)
}
func (n *Node) AsCatchClause() *CatchClause {
	return n.data.(*CatchClause)
}
func (n *Node) AsDeleteExpression() *DeleteExpression {
	return n.data.(*DeleteExpression)
}
func (n *Node) AsLabeledStatement() *LabeledStatement {
	return n.data.(*LabeledStatement)
}
func (n *Node) AsNamespaceExportDeclaration() *NamespaceExportDeclaration {
	return n.data.(*NamespaceExportDeclaration)
}
func (n *Node) AsNamedExports() *NamedExports {
	return n.data.(*NamedExports)
}
func (n *Node) AsBreakStatement() *BreakStatement {
	return n.data.(*BreakStatement)
}
func (n *Node) AsContinueStatement() *ContinueStatement {
	return n.data.(*ContinueStatement)
}
func (n *Node) AsCaseBlock() *CaseBlock {
	return n.data.(*CaseBlock)
}
func (n *Node) AsTryStatement() *TryStatement {
	return n.data.(*TryStatement)
}
func (n *Node) AsBindingPattern() *BindingPattern {
	return n.data.(*BindingPattern)
}
func (n *Node) AsFunctionDeclaration() *FunctionDeclaration {
	return n.data.(*FunctionDeclaration)
}
func (n *Node) AsTypeOfExpression() *TypeOfExpression {
	return n.data.(*TypeOfExpression)
}
func (n *Node) AsVoidExpression() *VoidExpression {
	return n.data.(*VoidExpression)
}
func (n *Node) AsAwaitExpression() *AwaitExpression {
	return n.data.(*AwaitExpression)
}
func (n *Node) AsSpreadElement() *SpreadElement {
	return n.data.(*SpreadElement)
}
func (n *Node) AsSpreadAssignment() *SpreadAssignment {
	return n.data.(*SpreadAssignment)
}
func (n *Node) AsArrayTypeNode() *ArrayTypeNode {
	return n.data.(*ArrayTypeNode)
}
func (n *Node) AsTupleTypeNode() *TupleTypeNode {
	return n.data.(*TupleTypeNode)
}
func (n *Node) AsUnionTypeNode() *UnionTypeNode {
	return n.data.(*UnionTypeNode)
}
func (n *Node) AsIntersectionTypeNode() *IntersectionTypeNode {
	return n.data.(*IntersectionTypeNode)
}
func (n *Node) AsRestTypeNode() *RestTypeNode {
	return n.data.(*RestTypeNode)
}
func (n *Node) AsNamedTupleMember() *NamedTupleMember {
	return n.data.(*NamedTupleMember)
}
func (n *Node) AsOptionalTypeNode() *OptionalTypeNode {
	return n.data.(*OptionalTypeNode)
}
func (n *Node) AsTypeReferenceNode() *TypeReferenceNode {
	return n.data.(*TypeReferenceNode)
}
func (n *Node) AsTypeQueryNode() *TypeQueryNode {
	return n.data.(*TypeQueryNode)
}
func (n *Node) AsIndexedAccessTypeNode() *IndexedAccessTypeNode {
	return n.data.(*IndexedAccessTypeNode)
}
func (n *Node) AsGetAccessorDeclaration() *GetAccessorDeclaration {
	return n.data.(*GetAccessorDeclaration)
}
func (n *Node) AsSetAccessorDeclaration() *SetAccessorDeclaration {
	return n.data.(*SetAccessorDeclaration)
}
func (n *Node) AsIndexSignatureDeclaration() *IndexSignatureDeclaration {
	return n.data.(*IndexSignatureDeclaration)
}
func (n *Node) AsClassStaticBlockDeclaration() *ClassStaticBlockDeclaration {
	return n.data.(*ClassStaticBlockDeclaration)
}
func (n *Node) AsJSDocNonNullableType() *JSDocNonNullableType {
	return n.data.(*JSDocNonNullableType)
}
func (n *Node) AsJSDocNullableType() *JSDocNullableType {
	return n.data.(*JSDocNullableType)
}
func (n *Node) AsSyntaxList() *SyntaxList {
	return n.data.(*SyntaxList)
}

// NodeData

type NodeData interface {
	AsNode() *Node
	ForEachChild(v Visitor) bool
	Accept(v NodeVisitor) *Node
	Name() *DeclarationName
	Modifiers() *ModifierListNode
	TypeParameters() *TypeParameterListNode
	FlowNodeData() *FlowNodeBase
	DeclarationData() *DeclarationBase
	ExportableData() *ExportableBase
	LocalsContainerData() *LocalsContainerBase
	FunctionLikeData() *FunctionLikeBase
	ClassLikeData() *ClassLikeBase
	BodyData() *BodyBase
}

// NodeDefault

type NodeDefault struct {
	Node
}

func (node *NodeDefault) AsNode() *Node                             { return &node.Node }
func (node *NodeDefault) ForEachChild(v Visitor) bool               { return false }
func (node *NodeDefault) Accept(v NodeVisitor) *Node                { return v.VisitOther(node.AsNode()) }
func (node *NodeDefault) Name() *DeclarationName                    { return nil }
func (node *NodeDefault) Modifiers() *ModifierListNode              { return nil }
func (node *NodeDefault) TypeParameters() *TypeParameterListNode    { return nil }
func (node *NodeDefault) FlowNodeData() *FlowNodeBase               { return nil }
func (node *NodeDefault) DeclarationData() *DeclarationBase         { return nil }
func (node *NodeDefault) ExportableData() *ExportableBase           { return nil }
func (node *NodeDefault) LocalsContainerData() *LocalsContainerBase { return nil }
func (node *NodeDefault) FunctionLikeData() *FunctionLikeBase       { return nil }
func (node *NodeDefault) ClassLikeData() *ClassLikeBase             { return nil }
func (node *NodeDefault) BodyData() *BodyBase                       { return nil }

// NodeBase

type NodeBase struct {
	NodeDefault
}

// Aliases for Node unions

type Statement = Node                   // Node with StatementBase
type Declaration = Node                 // Node with DeclarationBase
type Expression = Node                  // Node with ExpressionBase
type TypeNode = Node                    // Node with TypeNodeBase
type TypeElement = Node                 // Node with TypeElementBase
type ClassElement = Node                // Node with ClassElementBase
type NamedMember = Node                 // Node with NamedMemberBase
type ObjectLiteralElement = Node        // Node with ObjectLiteralElementBase
type BlockOrExpression = Node           // Block | Expression
type AccessExpression = Node            // PropertyAccessExpression | ElementAccessExpression
type DeclarationName = Node             // Identifier | PrivateIdentifier | StringLiteral | NumericLiteral | BigIntLiteral | NoSubstitutionTemplateLiteral | ComputedPropertyName | BindingPattern | ElementAccessExpression
type ModuleName = Node                  // Identifier | StringLiteral
type ModuleExportName = Node            // Identifier | StringLiteral
type PropertyName = Node                // Identifier | StringLiteral | NoSubstitutionTemplateLiteral | NumericLiteral | ComputedPropertyName | PrivateIdentifier | BigIntLiteral
type ModuleBody = Node                  // ModuleBlock | ModuleDeclaration
type ForInitializer = Node              // Expression | MissingDeclaration | VariableDeclarationList
type ModuleReference = Node             // Identifier | QualifiedName | ExternalModuleReference
type NamedImportBindings = Node         // NamespaceImport | NamedImports
type NamedExportBindings = Node         // NamespaceExport | NamedExports
type MemberName = Node                  // Identifier | PrivateIdentifier
type EntityName = Node                  // Identifier | QualifiedName
type BindingName = Node                 // Identifier | BindingPattern
type ModifierLike = Node                // Modifier | Decorator
type JsxChild = Node                    // JsxText | JsxExpression | JsxElement | JsxSelfClosingElement | JsxFragment
type JsxAttributeLike = Node            // JsxAttribute | JsxSpreadAttribute
type JsxAttributeName = Node            // Identifier | JsxNamespacedName
type JsxTagNameExpression = Node        // Identifier | KeywordExpression | JsxTagNamePropertyAccess | JsxNamespacedName
type ClassLikeDeclaration = Node        // ClassDeclaration | ClassExpression
type AccessorDeclaration = Node         // GetAccessorDeclaration | SetAccessorDeclaration
type LiteralLikeNode = Node             // StringLiteral | NumericLiteral | BigIntLiteral | RegularExpressionLiteral | TemplateLiteralLikeNode | JsxText
type LiteralExpression = Node           // StringLiteral | NumericLiteral | BigIntLiteral | RegularExpressionLiteral | NoSubstitutionTemplateLiteral
type UnionOrIntersectionTypeNode = Node // UnionTypeNode | IntersectionTypeNode
type TemplateLiteralLikeNode = Node     // TemplateHead | TemplateMiddle | TemplateTail
type TemplateMiddleOrTail = Node        // TemplateMiddle | TemplateTail

// Aliases for node signletons

type IdentifierNode = Node
type ModifierListNode = Node
type TokenNode = Node
type TemplateHeadNode = Node
type TemplateMiddleNode = Node
type TemplateTailNode = Node
type TemplateSpanNode = Node
type TemplateLiteralTypeSpanNode = Node
type BlockNode = Node
type CatchClauseNode = Node
type CaseBlockNode = Node
type CaseOrDefaultClauseNode = Node
type VariableDeclarationNode = Node
type VariableDeclarationListNode = Node
type BindingElementNode = Node
type TypeParameterListNode = Node
type TypeArgumentListNode = Node
type TypeParameterDeclarationNode = Node
type ParameterDeclarationNode = Node
type HeritageClauseNode = Node
type ExpressionWithTypeArgumentsNode = Node
type EnumMemberNode = Node
type ImportClauseNode = Node
type ImportAttributesNode = Node
type ImportAttributeNode = Node
type ImportSpecifierNode = Node
type ExportSpecifierNode = Node
type JsxAttributesNode = Node
type JsxOpeningElementNode = Node
type JsxClosingElementNode = Node
type JsxOpeningFragmentNode = Node
type JsxClosingFragmentNode = Node

// DeclarationBase

type DeclarationBase struct {
	symbol *Symbol // Symbol declared by node (initialized by binding)
}

func (node *DeclarationBase) DeclarationData() *DeclarationBase { return node }

func isDeclarationNode(node *Node) bool {
	return node.DeclarationData() != nil
}

// DeclarationBase

type ExportableBase struct {
	localSymbol *Symbol // Local symbol declared by node (initialized by binding only for exported nodes)
}

func (node *ExportableBase) ExportableData() *ExportableBase { return node }

// ModifiersBase

type ModifiersBase struct {
	modifiers *ModifierListNode
}

func (node *ModifiersBase) Modifiers() *ModifierListNode { return node.modifiers }

// LocalsContainerBase

type LocalsContainerBase struct {
	locals        SymbolTable // Locals associated with node (initialized by binding)
	nextContainer *Node       // Next container in declaration order (initialized by binding)
}

func (node *LocalsContainerBase) LocalsContainerData() *LocalsContainerBase { return node }

func isLocalsContainer(node *Node) bool {
	return node.LocalsContainerData() != nil
}

// FunctionLikeBase

type FunctionLikeBase struct {
	LocalsContainerBase
	typeParameters *TypeParameterListNode // Optional
	parameters     []*ParameterDeclarationNode
	returnType     *TypeNode // Optional
}

func (node *FunctionLikeBase) TypeParameters() *TypeParameterListNode { return node.typeParameters }
func (node *FunctionLikeBase) LocalsContainerData() *LocalsContainerBase {
	return &node.LocalsContainerBase
}
func (node *FunctionLikeBase) FunctionLikeData() *FunctionLikeBase { return node }
func (node *FunctionLikeBase) BodyData() *BodyBase                 { return nil }

// BodyBase

type BodyBase struct {
	asteriskToken *TokenNode
	body          *BlockOrExpression // Optional, can be Expression only in arrow functions
	endFlowNode   *FlowNode
}

func (node *BodyBase) BodyData() *BodyBase { return node }

// FunctionLikeWithBodyBase

type FunctionLikeWithBodyBase struct {
	FunctionLikeBase
	BodyBase
}

func (node *FunctionLikeWithBodyBase) TypeParameters() *TypeParameterListNode {
	return node.typeParameters
}
func (node *FunctionLikeWithBodyBase) LocalsContainerData() *LocalsContainerBase {
	return &node.LocalsContainerBase
}
func (node *FunctionLikeWithBodyBase) FunctionLikeData() *FunctionLikeBase {
	return &node.FunctionLikeBase
}
func (node *FunctionLikeWithBodyBase) BodyData() *BodyBase { return &node.BodyBase }

// FlowNodeBase

type FlowNodeBase struct {
	flowNode *FlowNode
}

func (node *FlowNodeBase) FlowNodeData() *FlowNodeBase { return node }

// Token

type Token struct {
	NodeBase
}

func (f *NodeFactory) NewToken(kind SyntaxKind) *Node {
	return f.NewNode(kind, &Token{})
}

func (node *Token) Accept(v NodeVisitor) *Node {
	return v.VisitToken(node)
}

// Identifier

type Identifier struct {
	ExpressionBase
	FlowNodeBase
	text string
}

func (f *NodeFactory) NewIdentifier(text string) *Node {
	data := f.identifierPool.New()
	data.text = text
	return f.NewNode(SyntaxKindIdentifier, data)
}

func (node *Identifier) Accept(v NodeVisitor) *Node {
	return v.VisitIdentifier(node)
}

func isIdentifier(node *Node) bool {
	return node.kind == SyntaxKindIdentifier
}

// PrivateIdentifier

type PrivateIdentifier struct {
	ExpressionBase
	text string
}

func (f *NodeFactory) NewPrivateIdentifier(text string) *Node {
	data := &PrivateIdentifier{}
	data.text = text
	return f.NewNode(SyntaxKindPrivateIdentifier, data)
}

func (node *PrivateIdentifier) Accept(v NodeVisitor) *Node {
	return v.VisitPrivateIdentifier(node)
}

func isPrivateIdentifier(node *Node) bool {
	return node.kind == SyntaxKindPrivateIdentifier
}

// QualifiedName

type QualifiedName struct {
	NodeBase
	FlowNodeBase
	left  *EntityName
	right *IdentifierNode
}

func (f *NodeFactory) NewQualifiedName(left *EntityName, right *IdentifierNode) *Node {
	data := &QualifiedName{}
	data.left = left
	data.right = right
	return f.NewNode(SyntaxKindQualifiedName, data)
}

func (f *NodeFactory) UpdateQualifiedName(node *Node, left *EntityName, right *IdentifierNode) *Node {
	if n := node.AsQualifiedName(); left != n.left || right != n.right {
		return f.UpdateNode(f.NewQualifiedName(left, right), node)
	}
	return node
}

func (node *QualifiedName) ForEachChild(v Visitor) bool {
	return visit(v, node.left) || visit(v, node.right)
}

func (node *QualifiedName) Accept(v NodeVisitor) *Node {
	return v.VisitQualifiedName(node)
}

func isQualifiedName(node *Node) bool {
	return node.kind == SyntaxKindQualifiedName
}

// TypeParameterDeclaration

type TypeParameterDeclaration struct {
	NodeBase
	DeclarationBase
	ModifiersBase
	name        *IdentifierNode // Identifier
	constraint  *TypeNode       // Optional
	defaultType *TypeNode       // Optional
	expression  *Node           // For error recovery purposes
}

func (f *NodeFactory) NewTypeParameterDeclaration(modifiers *Node, name *IdentifierNode, constraint *TypeNode, defaultType *TypeNode) *Node {
	data := &TypeParameterDeclaration{}
	data.modifiers = modifiers
	data.name = name
	data.constraint = constraint
	data.defaultType = defaultType
	return f.NewNode(SyntaxKindTypeParameter, data)
}

func (f *NodeFactory) UpdateTypeParameterDeclaration(node *Node, modifiers *Node, name *IdentifierNode, constraint *TypeNode, defaultType *TypeNode) *Node {
	if n := node.AsTypeParameter(); modifiers != n.modifiers || name != n.name || constraint != n.constraint || defaultType != n.defaultType {
		return f.UpdateNode(f.NewTypeParameterDeclaration(modifiers, name, constraint, defaultType), node)
	}
	return node
}

func (node *TypeParameterDeclaration) Kind() SyntaxKind {
	return SyntaxKindTypeParameter
}

func (node *TypeParameterDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.modifiers) || visit(v, node.name) || visit(v, node.constraint) || visit(v, node.defaultType)
}

func (node *TypeParameterDeclaration) Accept(v NodeVisitor) *Node {
	return v.VisitTypeParameterDeclaration(node)
}

func (node *TypeParameterDeclaration) Name() *DeclarationName {
	return node.name
}

func isTypeParameterDeclaration(node *Node) bool {
	return node.kind == SyntaxKindTypeParameter
}

// ComputedPropertyName

type ComputedPropertyName struct {
	NodeBase
	expression *Node
}

func (f *NodeFactory) NewComputedPropertyName(expression *Node) *Node {
	data := &ComputedPropertyName{}
	data.expression = f.ParenthesizerRules().ParenthesizeExpressionOfComputedPropertyName(expression)
	return f.NewNode(SyntaxKindComputedPropertyName, data)
}

func (f *NodeFactory) UpdateComputedPropertyName(node *Node, expression *Node) *Node {
	if n := node.AsComputedPropertyName(); expression != n.expression {
		return f.UpdateNode(f.NewComputedPropertyName(expression), node)
	}
	return node
}

func (node *ComputedPropertyName) ForEachChild(v Visitor) bool {
	return visit(v, node.expression)
}

func (node *ComputedPropertyName) Accept(v NodeVisitor) *Node {
	return v.VisitComputedPropertyName(node)
}

func isComputedPropertyName(node *Node) bool {
	return node.kind == SyntaxKindComputedPropertyName
}

// Modifier

func (f *NodeFactory) NewModifier(kind SyntaxKind) *Node {
	return f.NewToken(kind)
}

// Decorator

type Decorator struct {
	NodeBase
	expression *Node
}

func (f *NodeFactory) NewDecorator(expression *Node) *Node {
	data := &Decorator{}
	data.expression = f.ParenthesizerRules().ParenthesizeLeftSideOfAccess(expression, false /*optionalChain*/)
	return f.NewNode(SyntaxKindDecorator, data)
}

func (f *NodeFactory) UpdateDecorator(node *Node, expression *Node) *Node {
	if n := node.AsDecorator(); expression != n.expression {
		return f.UpdateNode(f.NewDecorator(expression), node)
	}
	return node
}

func (node *Decorator) ForEachChild(v Visitor) bool {
	return visit(v, node.expression)
}

func (node *Decorator) Accept(v NodeVisitor) *Node {
	return v.VisitDecorator(node)
}

func isDecorator(node *Node) bool {
	return node.kind == SyntaxKindDecorator
}

// ModifierList

type ModifierList struct {
	NodeBase
	modifiers     []*ModifierLike
	modifierFlags ModifierFlags
}

func (f *NodeFactory) NewModifierList(modifiers []*ModifierLike) *Node {
	data := &ModifierList{}
	data.modifiers = modifiers
	return f.NewNode(SyntaxKindModifierList, data)
}

func (f *NodeFactory) UpdateModifierList(node *Node, modifiers []*ModifierLike) *Node {
	if n := node.AsModifierList(); !slices.Equal(n.modifiers, modifiers) {
		return f.UpdateNode(f.NewModifierList(modifiers), node)
	}
	return node
}

func (node *ModifierList) ForEachChild(v Visitor) bool {
	return visitNodes(v, node.modifiers)
}

func (node *ModifierList) Accept(v NodeVisitor) *Node {
	return v.VisitModifierList(node)
}

// StatementBase

type StatementBase struct {
	NodeBase
	FlowNodeBase
}

// EmptyStatement

type EmptyStatement struct {
	StatementBase
}

func (f *NodeFactory) NewEmptyStatement() *Node {
	return f.NewNode(SyntaxKindEmptyStatement, &EmptyStatement{})
}

func isEmptyStatement(node *Node) bool {
	return node.kind == SyntaxKindEmptyStatement
}

// IfStatement

type IfStatement struct {
	StatementBase
	expression    *Node
	thenStatement *Statement
	elseStatement *Statement // Optional
}

func (f *NodeFactory) NewIfStatement(expression *Node, thenStatement *Statement, elseStatement *Statement) *Node {
	data := &IfStatement{}
	data.expression = expression
	data.thenStatement = thenStatement
	data.elseStatement = elseStatement
	return f.NewNode(SyntaxKindIfStatement, data)
}

func (f *NodeFactory) UpdateIfStatement(node *Node, expression *Node, thenStatement *Statement, elseStatement *Statement) *Node {
	if n := node.AsIfStatement(); expression != n.expression || thenStatement != n.thenStatement || elseStatement != n.elseStatement {
		return f.UpdateNode(f.NewIfStatement(expression, thenStatement, elseStatement), node)
	}
	return node
}

func (node *IfStatement) ForEachChild(v Visitor) bool {
	return visit(v, node.expression) || visit(v, node.thenStatement) || visit(v, node.elseStatement)
}

func (node *IfStatement) Accept(v NodeVisitor) *Node {
	return v.VisitIfStatement(node)
}

// DoStatement

type DoStatement struct {
	StatementBase
	statement  *Statement
	expression *Node
}

func (f *NodeFactory) NewDoStatement(statement *Statement, expression *Node) *Node {
	data := &DoStatement{}
	data.statement = statement
	data.expression = expression
	return f.NewNode(SyntaxKindDoStatement, data)
}

func (f *NodeFactory) UpdateDoStatement(node *Node, statement *Statement, expression *Node) *Node {
	if n := node.AsDoStatement(); statement != n.statement || expression != n.expression {
		return f.UpdateNode(f.NewDoStatement(statement, expression), node)
	}
	return node
}

func (node *DoStatement) ForEachChild(v Visitor) bool {
	return visit(v, node.statement) || visit(v, node.expression)
}

func (node *DoStatement) Accept(v NodeVisitor) *Node {
	return v.VisitDoStatement(node)
}

// WhileStatement

type WhileStatement struct {
	StatementBase
	expression *Node
	statement  *Statement
}

func (f *NodeFactory) NewWhileStatement(expression *Node, statement *Statement) *Node {
	data := &WhileStatement{}
	data.expression = expression
	data.statement = statement
	return f.NewNode(SyntaxKindWhileStatement, data)
}

func (f *NodeFactory) UpdateWhileStatement(node *Node, expression *Node, statement *Statement) *Node {
	if n := node.AsWhileStatement(); expression != n.expression || statement != n.statement {
		return f.UpdateNode(f.NewWhileStatement(expression, statement), node)
	}
	return node
}

func (node *WhileStatement) ForEachChild(v Visitor) bool {
	return visit(v, node.expression) || visit(v, node.statement)
}

func (node *WhileStatement) Accept(v NodeVisitor) *Node {
	return v.VisitWhileStatement(node)
}

// ForStatement

type ForStatement struct {
	StatementBase
	LocalsContainerBase
	initializer *ForInitializer // Optional
	condition   *Node           // Optional
	incrementor *Node           // Optional
	statement   *Statement
}

func (f *NodeFactory) NewForStatement(initializer *ForInitializer, condition *Node, incrementor *Node, statement *Statement) *Node {
	data := &ForStatement{}
	data.initializer = initializer
	data.condition = condition
	data.incrementor = incrementor
	data.statement = statement
	return f.NewNode(SyntaxKindForStatement, data)
}

func (f *NodeFactory) UpdateForStatement(node *Node, initializer *ForInitializer, condition *Node, incrementor *Node, statement *Statement) *Node {
	if n := node.AsForStatement(); initializer != n.initializer || condition != n.condition || incrementor != n.incrementor || statement != n.statement {
		return f.UpdateNode(f.NewForStatement(initializer, condition, incrementor, statement), node)
	}
	return node
}

func (node *ForStatement) ForEachChild(v Visitor) bool {
	return visit(v, node.initializer) || visit(v, node.condition) || visit(v, node.incrementor) || visit(v, node.statement)
}

func (node *ForStatement) Accept(v NodeVisitor) *Node {
	return v.VisitForStatement(node)
}

// ForInOrOfStatement

type ForInOrOfStatement struct {
	StatementBase
	LocalsContainerBase
	kind          SyntaxKind // SyntaxKindForInStatement | SyntaxKindForOfStatement
	awaitModifier *Node      // Optional
	initializer   *ForInitializer
	expression    *Node
	statement     *Statement
}

func (f *NodeFactory) NewForInOrOfStatement(kind SyntaxKind, awaitModifier *Node, initializer *ForInitializer, expression *Node, statement *Statement) *Node {
	data := &ForInOrOfStatement{}
	data.kind = kind
	data.awaitModifier = awaitModifier
	data.initializer = initializer
	data.expression = f.ParenthesizerRules().ParenthesizeExpressionForDisallowedComma(expression)
	data.statement = statement
	return f.NewNode(kind, data)
}

func (f *NodeFactory) UpdateForInOrOfStatement(node *Node, awaitModifier *Node, initializer *ForInitializer, expression *Node, statement *Statement) *Node {
	if n := node.AsForInOrOfStatement(); awaitModifier != n.awaitModifier || initializer != n.initializer || expression != n.expression || statement != n.statement {
		return f.UpdateNode(f.NewForInOrOfStatement(node.kind, awaitModifier, initializer, expression, statement), node)
	}
	return node
}

func (node *ForInOrOfStatement) ForEachChild(v Visitor) bool {
	return visit(v, node.awaitModifier) || visit(v, node.initializer) || visit(v, node.expression) || visit(v, node.statement)
}

func (node *ForInOrOfStatement) Accept(v NodeVisitor) *Node {
	switch node.kind {
	case SyntaxKindForInStatement:
		return v.VisitForInStatement(node)
	case SyntaxKindForOfStatement:
		return v.VisitForOfStatement(node)
	default:
		return v.VisitOther(node.AsNode())
	}
}

func isForInOrOfStatement(node *Node) bool {
	return node.kind == SyntaxKindForInStatement || node.kind == SyntaxKindForOfStatement
}

// BreakStatement

type BreakStatement struct {
	StatementBase
	label *IdentifierNode // Optional
}

func (f *NodeFactory) NewBreakStatement(label *IdentifierNode) *Node {
	data := &BreakStatement{}
	data.label = label
	return f.NewNode(SyntaxKindBreakStatement, data)
}

func (f *NodeFactory) UpdateBreakStatement(node *Node, label *IdentifierNode) *Node {
	if n := node.AsBreakStatement(); label != n.label {
		return f.UpdateNode(f.NewBreakStatement(label), node)
	}
	return node
}

func (node *BreakStatement) ForEachChild(v Visitor) bool {
	return visit(v, node.label)
}

func (node *BreakStatement) Accept(v NodeVisitor) *Node {
	return v.VisitBreakStatement(node)
}

// ContinueStatement

type ContinueStatement struct {
	StatementBase
	label *IdentifierNode // Optional
}

func (f *NodeFactory) NewContinueStatement(label *IdentifierNode) *Node {
	data := &ContinueStatement{}
	data.label = label
	return f.NewNode(SyntaxKindContinueStatement, data)
}

func (f *NodeFactory) UpdateContinueStatement(node *Node, label *IdentifierNode) *Node {
	if n := node.AsContinueStatement(); label != n.label {
		return f.UpdateNode(f.NewContinueStatement(label), node)
	}
	return node
}

func (node *ContinueStatement) ForEachChild(v Visitor) bool {
	return visit(v, node.label)
}

func (node *ContinueStatement) Accept(v NodeVisitor) *Node {
	return v.VisitContinueStatement(node)
}

// ReturnStatement

type ReturnStatement struct {
	StatementBase
	expression *Node // Optional
}

func (f *NodeFactory) NewReturnStatement(expression *Node) *Node {
	data := &ReturnStatement{}
	data.expression = expression
	return f.NewNode(SyntaxKindReturnStatement, data)
}

func (f *NodeFactory) UpdateReturnStatement(node *Node, expression *Node) *Node {
	if n := node.AsReturnStatement(); expression != n.expression {
		return f.UpdateNode(f.NewReturnStatement(expression), node)
	}
	return node
}

func (node *ReturnStatement) ForEachChild(v Visitor) bool {
	return visit(v, node.expression)
}

func (node *ReturnStatement) Accept(v NodeVisitor) *Node {
	return v.VisitReturnStatement(node)
}

// WithStatement

type WithStatement struct {
	StatementBase
	expression *Node
	statement  *Statement
}

func (f *NodeFactory) NewWithStatement(expression *Node, statement *Statement) *Node {
	data := &WithStatement{}
	data.expression = expression
	data.statement = statement
	return f.NewNode(SyntaxKindWithStatement, data)
}

func (f *NodeFactory) UpdateWithStatement(node *Node, expression *Node, statement *Statement) *Node {
	if n := node.AsWithStatement(); expression != n.expression || statement != n.statement {
		return f.UpdateNode(f.NewWithStatement(expression, statement), node)
	}
	return node
}

func (node *WithStatement) ForEachChild(v Visitor) bool {
	return visit(v, node.expression) || visit(v, node.statement)
}

func (node *WithStatement) Accept(v NodeVisitor) *Node {
	return v.VisitWithStatement(node)
}

// SwitchStatement

type SwitchStatement struct {
	StatementBase
	expression *Node
	caseBlock  *CaseBlockNode
}

func (f *NodeFactory) NewSwitchStatement(expression *Node, caseBlock *CaseBlockNode) *Node {
	data := &SwitchStatement{}
	data.expression = f.ParenthesizerRules().ParenthesizeExpressionForDisallowedComma(expression)
	data.caseBlock = caseBlock
	return f.NewNode(SyntaxKindSwitchStatement, data)
}

func (f *NodeFactory) UpdateSwitchStatement(node *Node, expression *Node, caseBlock *CaseBlockNode) *Node {
	if n := node.AsSwitchStatement(); expression != n.expression || caseBlock != n.caseBlock {
		return f.UpdateNode(f.NewSwitchStatement(expression, caseBlock), node)
	}
	return node
}

func (node *SwitchStatement) ForEachChild(v Visitor) bool {
	return visit(v, node.expression) || visit(v, node.caseBlock)
}

func (node *SwitchStatement) Accept(v NodeVisitor) *Node {
	return v.VisitSwitchStatement(node)
}

// CaseBlock

type CaseBlock struct {
	NodeBase
	LocalsContainerBase
	clauses []*CaseOrDefaultClauseNode
}

func (f *NodeFactory) NewCaseBlock(clauses []*CaseOrDefaultClauseNode) *Node {
	data := &CaseBlock{}
	data.clauses = clauses
	return f.NewNode(SyntaxKindCaseBlock, data)
}

func (f *NodeFactory) UpdateCaseBlock(node *Node, clauses []*CaseOrDefaultClauseNode) *Node {
	if n := node.AsCaseBlock(); !slices.Equal(clauses, n.clauses) {
		return f.UpdateNode(f.NewCaseBlock(clauses), node)
	}
	return node
}

func (node *CaseBlock) ForEachChild(v Visitor) bool {
	return visitNodes(v, node.clauses)
}

func (node *CaseBlock) Accept(v NodeVisitor) *Node {
	return v.VisitCaseBlock(node)
}

// CaseOrDefaultClause

type CaseOrDefaultClause struct {
	NodeBase
	expression          *Node // nil in default clause
	statements          []*Statement
	fallthroughFlowNode *FlowNode
}

func (f *NodeFactory) NewCaseOrDefaultClause(kind SyntaxKind, expression *Node, statements []*Statement) *Node {
	data := &CaseOrDefaultClause{}
	if expression != nil {
		data.expression = f.ParenthesizerRules().ParenthesizeExpressionForDisallowedComma(expression)
	}
	data.statements = statements
	return f.NewNode(kind, data)
}

func (f *NodeFactory) UpdateCaseOrDefaultClause(node *Node, expression *Node, statements []*Statement) *Node {
	if n := node.AsCaseOrDefaultClause(); expression != n.expression || !slices.Equal(statements, n.statements) {
		return f.UpdateNode(f.NewCaseOrDefaultClause(n.kind, expression, statements), node)
	}
	return node
}

func (node *CaseOrDefaultClause) ForEachChild(v Visitor) bool {
	return visit(v, node.expression) || visitNodes(v, node.statements)
}

func (node *CaseOrDefaultClause) Accept(v NodeVisitor) *Node {
	switch node.kind {
	case SyntaxKindCaseClause:
		return v.VisitCaseClause(node)
	case SyntaxKindDefaultClause:
		return v.VisitDefaultClause(node)
	default:
		return v.VisitOther(node.AsNode())
	}
}

// ThrowStatement

type ThrowStatement struct {
	StatementBase
	expression *Node
}

func (f *NodeFactory) NewThrowStatement(expression *Node) *Node {
	data := &ThrowStatement{}
	data.expression = expression
	return f.NewNode(SyntaxKindThrowStatement, data)
}

func (f *NodeFactory) UpdateThrowStatement(node *Node, expression *Node) *Node {
	if n := node.AsThrowStatement(); expression != n.expression {
		return f.UpdateNode(f.NewThrowStatement(expression), node)
	}
	return node
}

func (node *ThrowStatement) ForEachChild(v Visitor) bool {
	return visit(v, node.expression)
}

func (node *ThrowStatement) Accept(v NodeVisitor) *Node {
	return v.VisitThrowStatement(node)
}

// TryStatement

type TryStatement struct {
	StatementBase
	tryBlock     *BlockNode
	catchClause  *CatchClauseNode // Optional
	finallyBlock *BlockNode       // Optional
}

func (f *NodeFactory) NewTryStatement(tryBlock *BlockNode, catchClause *CatchClauseNode, finallyBlock *BlockNode) *Node {
	data := &TryStatement{}
	data.tryBlock = tryBlock
	data.catchClause = catchClause
	data.finallyBlock = finallyBlock
	return f.NewNode(SyntaxKindTryStatement, data)
}

func (f *NodeFactory) UpdateTryStatement(node *Node, tryBlock *BlockNode, catchClause *CatchClauseNode, finallyBlock *BlockNode) *Node {
	if n := node.AsTryStatement(); tryBlock != n.tryBlock || catchClause != n.catchClause || finallyBlock != n.finallyBlock {
		return f.UpdateNode(f.NewTryStatement(tryBlock, catchClause, finallyBlock), node)
	}
	return node
}

func (node *TryStatement) Kind() SyntaxKind {
	return SyntaxKindTryStatement
}

func (node *TryStatement) ForEachChild(v Visitor) bool {
	return visit(v, node.tryBlock) || visit(v, node.catchClause) || visit(v, node.finallyBlock)
}

func (node *TryStatement) Accept(v NodeVisitor) *Node {
	return v.VisitTryStatement(node)
}

// CatchClause

type CatchClause struct {
	NodeBase
	LocalsContainerBase
	variableDeclaration *VariableDeclarationNode // Optional
	block               *BlockNode
}

func (f *NodeFactory) NewCatchClause(variableDeclaration *VariableDeclarationNode, block *BlockNode) *Node {
	data := &CatchClause{}
	data.variableDeclaration = variableDeclaration
	data.block = block
	return f.NewNode(SyntaxKindCatchClause, data)
}

func (f *NodeFactory) UpdateCatchClause(node *Node, variableDeclaration *VariableDeclarationNode, block *BlockNode) *Node {
	if n := node.AsCatchClause(); variableDeclaration != n.variableDeclaration || block != n.block {
		return f.UpdateNode(f.NewCatchClause(variableDeclaration, block), node)
	}
	return node
}

func (node *CatchClause) ForEachChild(v Visitor) bool {
	return visit(v, node.variableDeclaration) || visit(v, node.block)
}

func (node *CatchClause) Accept(v NodeVisitor) *Node {
	return v.VisitCatchClause(node)
}

// DebuggerStatement

type DebuggerStatement struct {
	StatementBase
}

func (f *NodeFactory) NewDebuggerStatement() *Node {
	return f.NewNode(SyntaxKindDebuggerStatement, &DebuggerStatement{})
}

// LabeledStatement

type LabeledStatement struct {
	StatementBase
	label     *IdentifierNode
	statement *Statement
}

func (f *NodeFactory) NewLabeledStatement(label *IdentifierNode, statement *Statement) *Node {
	data := &LabeledStatement{}
	data.label = label
	data.statement = statement
	return f.NewNode(SyntaxKindLabeledStatement, data)
}

func (f *NodeFactory) UpdateLabeledStatement(node *Node, label *IdentifierNode, statement *Statement) *Node {
	if n := node.AsLabeledStatement(); label != n.label || statement != n.statement {
		return f.UpdateNode(f.NewLabeledStatement(label, statement), node)
	}
	return node
}

func (node *LabeledStatement) ForEachChild(v Visitor) bool {
	return visit(v, node.label) || visit(v, node.statement)
}

func (node *LabeledStatement) Accept(v NodeVisitor) *Node {
	return v.VisitLabeledStatement(node)
}

// ExpressionStatement

type ExpressionStatement struct {
	StatementBase
	expression *Node
}

func (f *NodeFactory) NewExpressionStatement(expression *Node) *Node {
	data := &ExpressionStatement{}
	data.expression = f.ParenthesizerRules().ParenthesizeExpressionOfExpressionStatement(expression)
	return f.NewNode(SyntaxKindExpressionStatement, data)
}

func (f *NodeFactory) UpdateExpressionStatement(node *Node, expression *Node) *Node {
	if n := node.AsExpressionStatement(); expression != n.expression {
		return f.UpdateNode(f.NewExpressionStatement(expression), node)
	}
	return node
}

func (node *ExpressionStatement) ForEachChild(v Visitor) bool {
	return visit(v, node.expression)
}

func (node *ExpressionStatement) Accept(v NodeVisitor) *Node {
	return v.VisitExpressionStatement(node)
}

func isExpressionStatement(node *Node) bool {
	return node.kind == SyntaxKindExpressionStatement
}

// Block

type Block struct {
	StatementBase
	LocalsContainerBase
	statements []*Statement
	multiline  bool
}

func (f *NodeFactory) NewBlock(statements []*Statement, multiline bool) *Node {
	data := &Block{}
	data.statements = statements
	data.multiline = multiline
	return f.NewNode(SyntaxKindBlock, data)
}

func (f *NodeFactory) UpdateBlock(node *Node, statements []*Statement) *Node {
	if n := node.AsBlock(); !slices.Equal(statements, n.statements) {
		return f.UpdateNode(f.NewBlock(statements, n.multiline), node)
	}
	return node
}

func (node *Block) ForEachChild(v Visitor) bool {
	return visitNodes(v, node.statements)
}

func (node *Block) Accept(v NodeVisitor) *Node {
	return v.VisitBlock(node)
}

func isBlock(node *Node) bool {
	return node.kind == SyntaxKindBlock
}

// VariableStatement

type VariableStatement struct {
	StatementBase
	ModifiersBase
	declarationList *VariableDeclarationListNode
}

func (f *NodeFactory) NewVariableStatement(modifiers *ModifierListNode, declarationList *VariableDeclarationListNode) *Node {
	data := &VariableStatement{}
	data.modifiers = modifiers
	data.declarationList = declarationList
	return f.NewNode(SyntaxKindVariableStatement, data)
}

func (f *NodeFactory) UpdateVariableStatement(node *Node, modifiers *ModifierListNode, declarationList *VariableDeclarationListNode) *Node {
	if n := node.AsVariableStatement(); modifiers != n.modifiers || declarationList != n.declarationList {
		return f.UpdateNode(f.NewVariableStatement(modifiers, declarationList), node)
	}
	return node
}

func (node *VariableStatement) ForEachChild(v Visitor) bool {
	return visit(v, node.modifiers) || visit(v, node.declarationList)
}

func (node *VariableStatement) Accept(v NodeVisitor) *Node {
	return v.VisitVariableStatement(node)
}

func isVariableStatement(node *Node) bool {
	return node.kind == SyntaxKindVariableStatement
}

// VariableDeclaration

type VariableDeclaration struct {
	NodeBase
	DeclarationBase
	ExportableBase
	name             *BindingName
	exclamationToken *TokenNode // Optional
	typeNode         *TypeNode  // Optional
	initializer      *Node      // Optional
}

func (f *NodeFactory) NewVariableDeclaration(name *BindingName, exclamationToken *TokenNode, typeNode *TypeNode, initializer *Node) *Node {
	data := &VariableDeclaration{}
	data.name = name
	data.exclamationToken = exclamationToken
	data.typeNode = typeNode
	if initializer != nil {
		data.initializer = f.ParenthesizerRules().ParenthesizeExpressionForDisallowedComma(initializer)
	}
	return f.NewNode(SyntaxKindVariableDeclaration, data)
}

func (f *NodeFactory) UpdateVariableDeclaration(node *Node, name *BindingName, exclamationToken *TokenNode, typeNode *TypeNode, initializer *Node) *Node {
	if n := node.AsVariableDeclaration(); name != n.name || exclamationToken != n.exclamationToken || typeNode != n.typeNode || initializer != n.initializer {
		return f.UpdateNode(f.NewVariableDeclaration(name, exclamationToken, typeNode, initializer), node)
	}
	return node
}

func (node *VariableDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.name) || visit(v, node.exclamationToken) || visit(v, node.typeNode) || visit(v, node.initializer)
}

func (node *VariableDeclaration) Accept(v NodeVisitor) *Node {
	return v.VisitVariableDeclaration(node)
}

func (node *VariableDeclaration) Name() *DeclarationName {
	return node.name
}

func isVariableDeclaration(node *Node) bool {
	return node.kind == SyntaxKindVariableDeclaration
}

// VariableDeclarationList

type VariableDeclarationList struct {
	NodeBase
	declarations []*VariableDeclarationNode
}

func (f *NodeFactory) NewVariableDeclarationList(flags NodeFlags, declarations []*VariableDeclarationNode) *Node {
	data := &VariableDeclarationList{}
	data.declarations = declarations
	node := f.NewNode(SyntaxKindVariableDeclarationList, data)
	node.flags = flags
	return node
}

func (f *NodeFactory) UpdateVariableDeclarationList(node *Node, declarations []*VariableDeclarationNode) *Node {
	if n := node.AsVariableDeclarationList(); !slices.Equal(declarations, n.declarations) {
		return f.UpdateNode(f.NewVariableDeclarationList(n.flags, declarations), node)
	}
	return node
}

func (node *VariableDeclarationList) ForEachChild(v Visitor) bool {
	return visitNodes(v, node.declarations)
}

func (node *VariableDeclarationList) Accept(v NodeVisitor) *Node {
	return v.VisitVariableDeclarationList(node)
}

func isVariableDeclarationList(node *Node) bool {
	return node.kind == SyntaxKindVariableDeclarationList
}

// BindingPattern (SyntaxBindObjectBindingPattern | SyntaxKindArrayBindingPattern)

type BindingPattern struct {
	NodeBase
	elements []*BindingElementNode
}

func (f *NodeFactory) NewBindingPattern(kind SyntaxKind, elements []*BindingElementNode) *Node {
	data := &BindingPattern{}
	data.elements = elements
	return f.NewNode(kind, data)
}

func (f *NodeFactory) UpdateBindingPattern(node *Node, elements []*BindingElementNode) *Node {
	if n := node.AsBindingPattern(); !slices.Equal(elements, n.elements) {
		return f.UpdateNode(f.NewBindingPattern(n.kind, elements), node)
	}
	return node
}

func (node *BindingPattern) ForEachChild(v Visitor) bool {
	return visitNodes(v, node.elements)
}

func (node *BindingPattern) Accept(v NodeVisitor) *Node {
	switch node.kind {
	case SyntaxKindObjectBindingPattern:
		return v.VisitObjectBindingPattern(node)
	case SyntaxKindArrayBindingPattern:
		return v.VisitArrayBindingPattern(node)
	}
	return v.VisitOther(node.AsNode())
}

func isObjectBindingPattern(node *Node) bool {
	return node.kind == SyntaxKindObjectBindingPattern
}

func isArrayBindingPattern(node *Node) bool {
	return node.kind == SyntaxKindArrayBindingPattern
}

func isBindingPattern(node *Node) bool {
	return isObjectBindingPattern(node) || isArrayBindingPattern(node)
}

// ParameterDeclaration

type ParameterDeclaration struct {
	NodeBase
	DeclarationBase
	ModifiersBase
	dotDotDotToken *TokenNode   // Present on rest parameter
	name           *BindingName // Declared parameter name
	questionToken  *TokenNode   // Present on optional parameter
	typeNode       *TypeNode    // Optional
	initializer    *Node        // Optional
}

func (f *NodeFactory) NewParameterDeclaration(modifiers *ModifierListNode, dotDotDotToken *TokenNode, name *BindingName, questionToken *TokenNode, typeNode *TypeNode, initializer *Node) *Node {
	data := &ParameterDeclaration{}
	data.modifiers = modifiers
	data.dotDotDotToken = dotDotDotToken
	data.name = name
	data.questionToken = questionToken
	data.typeNode = typeNode
	if initializer != nil {
		data.initializer = f.ParenthesizerRules().ParenthesizeExpressionForDisallowedComma(initializer)
	}
	return f.NewNode(SyntaxKindParameter, data)
}

func (f *NodeFactory) UpdateParameterDeclaration(node *Node, modifiers *ModifierListNode, dotDotDotToken *TokenNode, name *BindingName, questionToken *TokenNode, typeNode *TypeNode, initializer *Node) *Node {
	if n := node.AsParameterDeclaration(); modifiers != n.modifiers || dotDotDotToken != n.dotDotDotToken || name != n.name || questionToken != n.questionToken || typeNode != n.typeNode || initializer != n.initializer {
		return f.UpdateNode(f.NewParameterDeclaration(modifiers, dotDotDotToken, name, questionToken, typeNode, initializer), node)
	}
	return node
}

func (node *ParameterDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.modifiers) || visit(v, node.dotDotDotToken) || visit(v, node.name) ||
		visit(v, node.questionToken) || visit(v, node.typeNode) || visit(v, node.initializer)
}

func (node *ParameterDeclaration) Accept(v NodeVisitor) *Node {
	return v.VisitParameterDeclaration(node)
}

func (node *ParameterDeclaration) Name() *DeclarationName {
	return node.name
}

func isParameter(node *Node) bool {
	return node.kind == SyntaxKindParameter
}

// BindingElement

type BindingElement struct {
	NodeBase
	DeclarationBase
	ExportableBase
	FlowNodeBase
	dotDotDotToken *TokenNode    // Present on rest element (in object binding pattern)
	propertyName   *PropertyName // Optional binding property name in object binding pattern
	name           *BindingName  // Optional (nil for missing element)
	initializer    *Node         // Optional
}

func (f *NodeFactory) NewBindingElement(dotDotDotToken *TokenNode, propertyName *PropertyName, name *BindingName, initializer *Node) *Node {
	data := &BindingElement{}
	data.dotDotDotToken = dotDotDotToken
	data.propertyName = propertyName
	data.name = name
	if initializer != nil {
		data.initializer = f.ParenthesizerRules().ParenthesizeExpressionForDisallowedComma(initializer)
	}
	return f.NewNode(SyntaxKindBindingElement, data)
}

func (f *NodeFactory) UpdateBindingElement(node *Node, dotDotDotToken *TokenNode, propertyName *PropertyName, name *BindingName, initializer *Node) *Node {
	if n := node.AsBindingElement(); dotDotDotToken != n.dotDotDotToken || propertyName != n.propertyName || name != n.name || initializer != n.initializer {
		return f.UpdateNode(f.NewBindingElement(dotDotDotToken, propertyName, name, initializer), node)
	}
	return node
}

func (node *BindingElement) ForEachChild(v Visitor) bool {
	return visit(v, node.propertyName) || visit(v, node.dotDotDotToken) || visit(v, node.name) || visit(v, node.initializer)
}

func (node *BindingElement) Accept(v NodeVisitor) *Node {
	return v.VisitBindingElement(node)
}

func (node *BindingElement) Name() *DeclarationName {
	return node.name
}

func isBindingElement(node *Node) bool {
	return node.kind == SyntaxKindBindingElement
}

// MissingDeclaration

type MissingDeclaration struct {
	StatementBase
	DeclarationBase
	ModifiersBase
}

func (f *NodeFactory) NewMissingDeclaration(modifiers *ModifierListNode) *Node {
	data := &MissingDeclaration{}
	data.modifiers = modifiers
	return f.NewNode(SyntaxKindMissingDeclaration, data)
}

func (f *NodeFactory) UpdateMissingDeclaration(node *Node, modifiers *ModifierListNode) *Node {
	if n := node.AsMissingDeclaration(); modifiers != n.modifiers {
		return f.UpdateNode(f.NewMissingDeclaration(modifiers), node)
	}
	return node
}

func (node *MissingDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.modifiers)
}

func (node *MissingDeclaration) Accept(v NodeVisitor) *Node {
	return v.VisitMissingDeclaration(node)
}

// FunctionDeclaration

type FunctionDeclaration struct {
	StatementBase
	DeclarationBase
	ExportableBase
	ModifiersBase
	FunctionLikeWithBodyBase
	name           *IdentifierNode
	returnFlowNode *FlowNode
}

func (f *NodeFactory) NewFunctionDeclaration(modifiers *ModifierListNode, asteriskToken *TokenNode, name *IdentifierNode, typeParameters *TypeParameterListNode, parameters []*ParameterDeclarationNode, returnType *TypeNode, body *BlockNode) *Node {
	data := &FunctionDeclaration{}
	data.modifiers = modifiers
	data.asteriskToken = asteriskToken
	data.name = name
	data.typeParameters = typeParameters
	data.parameters = parameters
	data.returnType = returnType
	data.body = body
	return f.NewNode(SyntaxKindFunctionDeclaration, data)
}

func (f *NodeFactory) UpdateFunctionDeclaration(node *Node, modifiers *ModifierListNode, asteriskToken *TokenNode, name *IdentifierNode, typeParameters *TypeParameterListNode, parameters []*ParameterDeclarationNode, returnType *TypeNode, body *BlockNode) *Node {
	if n := node.AsFunctionDeclaration(); modifiers != n.modifiers || asteriskToken != n.asteriskToken || name != n.name || typeParameters != n.typeParameters || !slices.Equal(parameters, n.parameters) || returnType != n.returnType || body != n.body {
		return f.UpdateNode(f.NewFunctionDeclaration(modifiers, asteriskToken, name, typeParameters, parameters, returnType, body), node)
	}
	return node
}

func (node *FunctionDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.modifiers) || visit(v, node.asteriskToken) || visit(v, node.name) || visit(v, node.typeParameters) ||
		visitNodes(v, node.parameters) || visit(v, node.returnType) || visit(v, node.body)
}

func (node *FunctionDeclaration) Accept(v NodeVisitor) *Node {
	return v.VisitFunctionDeclaration(node)
}

func (node *FunctionDeclaration) Name() *DeclarationName {
	return node.name
}

func (node *FunctionDeclaration) BodyData() *BodyBase { return &node.BodyBase }

func isFunctionDeclaration(node *Node) bool {
	return node.kind == SyntaxKindFunctionDeclaration
}

// ClassLikeDeclarationBase

type ClassLikeBase struct {
	DeclarationBase
	ExportableBase
	ModifiersBase
	name            *IdentifierNode
	typeParameters  *TypeParameterListNode
	heritageClauses []*HeritageClauseNode
	members         []*ClassElement
}

func (node *ClassLikeBase) ForEachChild(v Visitor) bool {
	return visit(v, node.modifiers) || visit(v, node.name) || visit(v, node.typeParameters) ||
		visitNodes(v, node.heritageClauses) || visitNodes(v, node.members)
}

func (node *ClassLikeBase) Name() *DeclarationName                 { return node.name }
func (node *ClassLikeBase) TypeParameters() *TypeParameterListNode { return node.typeParameters }
func (node *ClassLikeBase) ClassLikeData() *ClassLikeBase          { return node }

// ClassDeclaration

type ClassDeclaration struct {
	StatementBase
	ClassLikeBase
}

func (f *NodeFactory) NewClassDeclaration(modifiers *ModifierListNode, name *IdentifierNode, typeParameters *TypeParameterListNode, heritageClauses []*HeritageClauseNode, members []*ClassElement) *Node {
	data := &ClassDeclaration{}
	data.modifiers = modifiers
	data.name = name
	data.typeParameters = typeParameters
	data.heritageClauses = heritageClauses
	data.members = members
	return f.NewNode(SyntaxKindClassDeclaration, data)
}

func (f *NodeFactory) UpdateClassDeclaration(node *Node, modifiers *ModifierListNode, name *IdentifierNode, typeParameters *TypeParameterListNode, heritageClauses []*HeritageClauseNode, members []*ClassElement) *Node {
	if n := node.AsClassDeclaration(); modifiers != n.modifiers || name != n.name || typeParameters != n.typeParameters || !slices.Equal(heritageClauses, n.heritageClauses) || !slices.Equal(members, n.members) {
		return f.UpdateNode(f.NewClassDeclaration(modifiers, name, typeParameters, heritageClauses, members), node)
	}
	return node
}

func (node *ClassDeclaration) Accept(v NodeVisitor) *Node {
	return v.VisitClassDeclaration(node)
}

func isClassDeclaration(node *Node) bool {
	return node.kind == SyntaxKindClassDeclaration
}

// ClassExpression

type ClassExpression struct {
	ExpressionBase
	ClassLikeBase
}

func (f *NodeFactory) NewClassExpression(modifiers *ModifierListNode, name *IdentifierNode, typeParameters *TypeParameterListNode, heritageClauses []*HeritageClauseNode, members []*ClassElement) *Node {
	data := &ClassExpression{}
	data.modifiers = modifiers
	data.name = name
	data.typeParameters = typeParameters
	data.heritageClauses = heritageClauses
	data.members = members
	return f.NewNode(SyntaxKindClassExpression, data)
}

func (f *NodeFactory) UpdateClassExpression(node *Node, modifiers *ModifierListNode, name *IdentifierNode, typeParameters *TypeParameterListNode, heritageClauses []*HeritageClauseNode, members []*ClassElement) *Node {
	if n := node.AsClassExpression(); modifiers != n.modifiers || name != n.name || typeParameters != n.typeParameters || !slices.Equal(heritageClauses, n.heritageClauses) || !slices.Equal(members, n.members) {
		return f.UpdateNode(f.NewClassExpression(modifiers, name, typeParameters, heritageClauses, members), node)
	}
	return node
}

func (node *ClassExpression) Kind() SyntaxKind { return SyntaxKindClassExpression }

func (node *ClassExpression) Accept(v NodeVisitor) *Node {
	return v.VisitClassExpression(node)
}

func isClassExpression(node *Node) bool {
	return node.kind == SyntaxKindClassExpression
}

// HeritageClause

type HeritageClause struct {
	NodeBase
	token SyntaxKind
	types []*ExpressionWithTypeArgumentsNode
}

func (f *NodeFactory) NewHeritageClause(token SyntaxKind, types []*ExpressionWithTypeArgumentsNode) *Node {
	data := &HeritageClause{}
	data.token = token
	data.types = types
	return f.NewNode(SyntaxKindHeritageClause, data)
}

func (f *NodeFactory) UpdateHeritageClause(node *Node, types []*ExpressionWithTypeArgumentsNode) *Node {
	if n := node.AsHeritageClause(); !slices.Equal(types, n.types) {
		return f.UpdateNode(f.NewHeritageClause(n.token, types), node)
	}
	return node
}

func (node *HeritageClause) ForEachChild(v Visitor) bool {
	return visitNodes(v, node.types)
}

func (node *HeritageClause) Accept(v NodeVisitor) *Node {
	return v.VisitHeritageClause(node)
}

func isHeritageClause(node *Node) bool {
	return node.kind == SyntaxKindHeritageClause
}

// InterfaceDeclaration

type InterfaceDeclaration struct {
	StatementBase
	DeclarationBase
	ExportableBase
	ModifiersBase
	name            *IdentifierNode
	typeParameters  *TypeParameterListNode
	heritageClauses []*HeritageClauseNode
	members         []*TypeElement
}

func (f *NodeFactory) NewInterfaceDeclaration(modifiers *ModifierListNode, name *IdentifierNode, typeParameters *TypeParameterListNode, heritageClauses []*HeritageClauseNode, members []*TypeElement) *Node {
	data := &InterfaceDeclaration{}
	data.modifiers = modifiers
	data.name = name
	data.typeParameters = typeParameters
	data.heritageClauses = heritageClauses
	data.members = members
	return f.NewNode(SyntaxKindInterfaceDeclaration, data)
}

func (f *NodeFactory) UpdateInterfaceDeclaration(node *Node, modifiers *ModifierListNode, name *IdentifierNode, typeParameters *TypeParameterListNode, heritageClauses []*HeritageClauseNode, members []*TypeElement) *Node {
	if n := node.AsInterfaceDeclaration(); modifiers != n.modifiers || name != n.name || typeParameters != n.typeParameters || !slices.Equal(heritageClauses, n.heritageClauses) || !slices.Equal(members, n.members) {
		return f.UpdateNode(f.NewInterfaceDeclaration(modifiers, name, typeParameters, heritageClauses, members), node)
	}
	return node
}

func (node *InterfaceDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.modifiers) || visit(v, node.name) || visit(v, node.typeParameters) ||
		visitNodes(v, node.heritageClauses) || visitNodes(v, node.members)
}

func (node *InterfaceDeclaration) Accept(v NodeVisitor) *Node {
	return v.VisitInterfaceDeclaration(node)
}

func (node *InterfaceDeclaration) Name() *DeclarationName                 { return node.name }
func (node *InterfaceDeclaration) TypeParameters() *TypeParameterListNode { return node.typeParameters }

func isInterfaceDeclaration(node *Node) bool {
	return node.kind == SyntaxKindInterfaceDeclaration
}

// TypeAliasDeclaration

type TypeAliasDeclaration struct {
	StatementBase
	DeclarationBase
	ExportableBase
	ModifiersBase
	LocalsContainerBase
	name           *IdentifierNode
	typeParameters *TypeParameterListNode
	typeNode       *TypeNode
}

func (f *NodeFactory) NewTypeAliasDeclaration(modifiers *ModifierListNode, name *IdentifierNode, typeParameters *TypeParameterListNode, typeNode *TypeNode) *Node {
	data := &TypeAliasDeclaration{}
	data.modifiers = modifiers
	data.name = name
	data.typeParameters = typeParameters
	data.typeNode = typeNode
	return f.NewNode(SyntaxKindTypeAliasDeclaration, data)
}

func (f *NodeFactory) UpdateTypeAliasDeclaration(node *Node, modifiers *ModifierListNode, name *IdentifierNode, typeParameters *TypeParameterListNode, typeNode *TypeNode) *Node {
	if n := node.AsTypeAliasDeclaration(); modifiers != n.modifiers || name != n.name || typeParameters != n.typeParameters || typeNode != n.typeNode {
		return f.UpdateNode(f.NewTypeAliasDeclaration(modifiers, name, typeParameters, typeNode), node)
	}
	return node
}

func (node *TypeAliasDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.modifiers) || visit(v, node.name) || visit(v, node.typeParameters) || visit(v, node.typeNode)
}

func (node *TypeAliasDeclaration) Accept(v NodeVisitor) *Node {
	return v.VisitTypeAliasDeclaration(node)
}

func (node *TypeAliasDeclaration) Name() *DeclarationName                 { return node.name }
func (node *TypeAliasDeclaration) TypeParameters() *TypeParameterListNode { return node.typeParameters }

func isTypeAliasDeclaration(node *Node) bool {
	return node.kind == SyntaxKindTypeAliasDeclaration
}

// EnumMember

type EnumMember struct {
	NodeBase
	NamedMemberBase
	initializer *Node // Optional
}

func (f *NodeFactory) NewEnumMember(name *PropertyName, initializer *Node) *Node {
	data := &EnumMember{}
	data.name = name
	if initializer != nil {
		data.initializer = f.ParenthesizerRules().ParenthesizeExpressionForDisallowedComma(initializer)
	}
	return f.NewNode(SyntaxKindEnumMember, data)
}

func (f *NodeFactory) UpdateEnumMember(node *Node, name *PropertyName, initializer *Node) *Node {
	if n := node.AsEnumMember(); name != n.name || initializer != n.initializer {
		return f.UpdateNode(f.NewEnumMember(name, initializer), node)
	}
	return node
}

func (node *EnumMember) ForEachChild(v Visitor) bool {
	return visit(v, node.name) || visit(v, node.initializer)
}

func (node *EnumMember) Accept(v NodeVisitor) *Node {
	return v.VisitEnumMember(node)
}

func (node *EnumMember) Name() *DeclarationName {
	return node.name
}

// EnumDeclaration

type EnumDeclaration struct {
	StatementBase
	DeclarationBase
	ExportableBase
	ModifiersBase
	name    *IdentifierNode
	members []*EnumMemberNode
}

func (f *NodeFactory) NewEnumDeclaration(modifiers *ModifierListNode, name *IdentifierNode, members []*EnumMemberNode) *Node {
	data := &EnumDeclaration{}
	data.modifiers = modifiers
	data.name = name
	data.members = members
	return f.NewNode(SyntaxKindEnumDeclaration, data)
}

func (f *NodeFactory) UpdateEnumDeclaration(node *Node, modifiers *ModifierListNode, name *IdentifierNode, members []*EnumMemberNode) *Node {
	if n := node.AsEnumDeclaration(); modifiers != n.modifiers || name != n.name || !slices.Equal(members, n.members) {
		return f.UpdateNode(f.NewEnumDeclaration(modifiers, name, members), node)
	}
	return node
}

func (node *EnumDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.modifiers) || visit(v, node.name) || visitNodes(v, node.members)
}

func (node *EnumDeclaration) Accept(v NodeVisitor) *Node {
	return v.VisitEnumDeclaration(node)
}

func (node *EnumDeclaration) Name() *DeclarationName {
	return node.name
}

func isEnumDeclaration(node *Node) bool {
	return node.kind == SyntaxKindEnumDeclaration
}

// ModuleBlock

type ModuleBlock struct {
	StatementBase
	statements []*Statement
}

func (f *NodeFactory) NewModuleBlock(statements []*Statement) *Node {
	data := &ModuleBlock{}
	data.statements = statements
	return f.NewNode(SyntaxKindModuleBlock, data)
}

func (f *NodeFactory) UpdateModuleBlock(node *Node, statements []*Statement) *Node {
	if n := node.AsModuleBlock(); !slices.Equal(statements, n.statements) {
		return f.UpdateNode(f.NewModuleBlock(statements), node)
	}
	return node
}

func (node *ModuleBlock) ForEachChild(v Visitor) bool {
	return visitNodes(v, node.statements)
}

func (node *ModuleBlock) Accept(v NodeVisitor) *Node {
	return v.VisitModuleBlock(node)
}

func isModuleBlock(node *Node) bool {
	return node.kind == SyntaxKindModuleBlock
}

// ModuleDeclaration

type ModuleDeclaration struct {
	StatementBase
	DeclarationBase
	ExportableBase
	ModifiersBase
	LocalsContainerBase
	name *ModuleName
	body *ModuleBody // Optional (may be nil in ambient module declaration)
}

func (f *NodeFactory) NewModuleDeclaration(modifiers *ModifierListNode, name *ModuleName, body *ModuleBody, flags NodeFlags) *Node {
	data := &ModuleDeclaration{}
	data.modifiers = modifiers
	data.name = name
	data.body = body
	node := f.NewNode(SyntaxKindModuleDeclaration, data)
	node.flags |= flags & (NodeFlagsNamespace | NodeFlagsNestedNamespace | NodeFlagsGlobalAugmentation)
	return node
}

func (f *NodeFactory) UpdateModuleDeclaration(node *Node, modifiers *ModifierListNode, name *ModuleName, body *ModuleBody) *Node {
	if n := node.AsModuleDeclaration(); modifiers != n.modifiers || name != n.name || body != n.body {
		return f.UpdateNode(f.NewModuleDeclaration(modifiers, name, body, n.flags), node)
	}
	return node
}

func (node *ModuleDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.modifiers) || visit(v, node.name) || visit(v, node.body)
}

func (node *ModuleDeclaration) Accept(v NodeVisitor) *Node {
	return v.VisitModuleDeclaration(node)
}

func (node *ModuleDeclaration) Name() *DeclarationName {
	return node.name
}

func isModuleDeclaration(node *Node) bool {
	return node.kind == SyntaxKindModuleDeclaration
}

// ModuleEqualsDeclaration

type ImportEqualsDeclaration struct {
	StatementBase
	DeclarationBase
	ExportableBase
	ModifiersBase
	modifiers  *ModifierListNode
	isTypeOnly bool
	name       *IdentifierNode
	// 'EntityName' for an internal module reference, 'ExternalModuleReference' for an external
	// module reference.
	moduleReference *ModuleReference
}

func (f *NodeFactory) NewImportEqualsDeclaration(modifiers *ModifierListNode, isTypeOnly bool, name *IdentifierNode, moduleReference *ModuleReference) *Node {
	data := &ImportEqualsDeclaration{}
	data.modifiers = modifiers
	data.isTypeOnly = isTypeOnly
	data.name = name
	data.moduleReference = moduleReference
	return f.NewNode(SyntaxKindImportEqualsDeclaration, data)
}

func (f *NodeFactory) UpdateImportEqualsDeclaration(node *Node, modifiers *ModifierListNode, isTypeOnly bool, name *IdentifierNode, moduleReference *ModuleReference) *Node {
	if n := node.AsImportEqualsDeclaration(); modifiers != n.modifiers || isTypeOnly != n.isTypeOnly || name != n.name || moduleReference != n.moduleReference {
		return f.UpdateNode(f.NewImportEqualsDeclaration(modifiers, isTypeOnly, name, moduleReference), node)
	}
	return node
}

func (node *ImportEqualsDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.modifiers) || visit(v, node.name) || visit(v, node.moduleReference)
}

func (node *ImportEqualsDeclaration) Accept(v NodeVisitor) *Node {
	return v.VisitImportEqualsDeclaration(node)
}

func (node *ImportEqualsDeclaration) Name() *DeclarationName {
	return node.name
}

func isImportEqualsDeclaration(node *Node) bool {
	return node.kind == SyntaxKindImportEqualsDeclaration
}

// ImportDeclaration

type ImportDeclaration struct {
	StatementBase
	ModifiersBase
	importClause    *ImportClauseNode
	moduleSpecifier *Node
	attributes      *ImportAttributesNode
}

func (f *NodeFactory) NewImportDeclaration(modifiers *ModifierListNode, importClause *ImportClauseNode, moduleSpecifier *Node, attributes *ImportAttributesNode) *Node {
	data := &ImportDeclaration{}
	data.modifiers = modifiers
	data.importClause = importClause
	data.moduleSpecifier = moduleSpecifier
	data.attributes = attributes
	return f.NewNode(SyntaxKindImportDeclaration, data)
}

func (f *NodeFactory) UpdateImportDeclaration(node *Node, modifiers *ModifierListNode, importClause *ImportClauseNode, moduleSpecifier *Node, attributes *ImportAttributesNode) *Node {
	if n := node.AsImportDeclaration(); modifiers != n.modifiers || importClause != n.importClause || moduleSpecifier != n.moduleSpecifier || attributes != n.attributes {
		return f.UpdateNode(f.NewImportDeclaration(modifiers, importClause, moduleSpecifier, attributes), node)
	}
	return node
}

func (node *ImportDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.modifiers) || visit(v, node.importClause) || visit(v, node.moduleSpecifier) || visit(v, node.attributes)
}

func (node *ImportDeclaration) Accept(v NodeVisitor) *Node {
	return v.VisitImportDeclaration(node)
}

func isImportDeclaration(node *Node) bool {
	return node.kind == SyntaxKindImportDeclaration
}

// ImportSpecifier

type ImportSpecifier struct {
	NodeBase
	DeclarationBase
	ExportableBase
	isTypeOnly   bool
	propertyName *ModuleExportName
	name         *IdentifierNode
}

func (f *NodeFactory) NewImportSpecifier(isTypeOnly bool, propertyName *ModuleExportName, name *IdentifierNode) *Node {
	data := &ImportSpecifier{}
	data.isTypeOnly = isTypeOnly
	data.propertyName = propertyName
	data.name = name
	return f.NewNode(SyntaxKindImportSpecifier, data)
}

func (f *NodeFactory) UpdateImportSpecifier(node *Node, isTypeOnly bool, propertyName *ModuleExportName, name *IdentifierNode) *Node {
	if n := node.AsImportSpecifier(); isTypeOnly != n.isTypeOnly || propertyName != n.propertyName || name != n.name {
		return f.UpdateNode(f.NewImportSpecifier(isTypeOnly, propertyName, name), node)
	}
	return node
}

func (node *ImportSpecifier) ForEachChild(v Visitor) bool {
	return visit(v, node.propertyName) || visit(v, node.name)
}

func (node *ImportSpecifier) Accept(v NodeVisitor) *Node {
	return v.VisitImportSpecifier(node)
}

func (node *ImportSpecifier) Name() *DeclarationName {
	return node.name
}

func isImportSpecifier(node *Node) bool {
	return node.kind == SyntaxKindImportSpecifier
}

// ExternalModuleReference

type ExternalModuleReference struct {
	NodeBase
	expression *Node
}

func (f *NodeFactory) NewExternalModuleReference(expression *Node) *Node {
	data := &ExternalModuleReference{}
	data.expression = expression
	return f.NewNode(SyntaxKindExternalModuleReference, data)
}

func (f *NodeFactory) UpdateExternalModuleReference(node *Node, expression *Node) *Node {
	if n := node.AsExternalModuleReference(); expression != n.expression {
		return f.UpdateNode(f.NewExternalModuleReference(expression), node)
	}
	return node
}

func (node *ExternalModuleReference) ForEachChild(v Visitor) bool {
	return visit(v, node.expression)
}

func (node *ExternalModuleReference) Accept(v NodeVisitor) *Node {
	return v.VisitExternalModuleReference(node)
}

func isExternalModuleReference(node *Node) bool {
	return node.kind == SyntaxKindExternalModuleReference
}

// ImportClause

type ImportClause struct {
	NodeBase
	DeclarationBase
	ExportableBase
	isTypeOnly    bool
	namedBindings *NamedImportBindings // Optional, named bindings
	name          *IdentifierNode      // Optional, default binding
}

func (f *NodeFactory) NewImportClause(isTypeOnly bool, name *IdentifierNode, namedBindings *NamedImportBindings) *Node {
	data := &ImportClause{}
	data.isTypeOnly = isTypeOnly
	data.name = name
	data.namedBindings = namedBindings
	return f.NewNode(SyntaxKindImportClause, data)
}

func (f *NodeFactory) UpdateImportClause(node *Node, isTypeOnly bool, name *IdentifierNode, namedBindings *NamedImportBindings) *Node {
	if n := node.AsImportClause(); isTypeOnly != n.isTypeOnly || name != n.name || namedBindings != n.namedBindings {
		return f.UpdateNode(f.NewImportClause(isTypeOnly, name, namedBindings), node)
	}
	return node
}

func (node *ImportClause) ForEachChild(v Visitor) bool {
	return visit(v, node.name) || visit(v, node.namedBindings)
}

func (node *ImportClause) Accept(v NodeVisitor) *Node {
	return v.VisitImportClause(node)
}

func (node *ImportClause) Name() *DeclarationName {
	return node.name
}

func isImportClause(node *Node) bool {
	return node.kind == SyntaxKindImportClause
}

// NamespaceImport

type NamespaceImport struct {
	NodeBase
	DeclarationBase
	ExportableBase
	name *IdentifierNode
}

func (f *NodeFactory) NewNamespaceImport(name *IdentifierNode) *Node {
	data := &NamespaceImport{}
	data.name = name
	return f.NewNode(SyntaxKindNamespaceImport, data)
}

func (f *NodeFactory) UpdateNamespaceImport(node *Node, name *IdentifierNode) *Node {
	if n := node.AsNamespaceImport(); name != n.name {
		return f.UpdateNode(f.NewNamespaceImport(name), node)
	}
	return node
}

func (node *NamespaceImport) ForEachChild(v Visitor) bool {
	return visit(v, node.name)
}

func (node *NamespaceImport) Accept(v NodeVisitor) *Node {
	return v.VisitNamespaceImport(node)
}

func (node *NamespaceImport) Name() *DeclarationName {
	return node.name
}

func isNamespaceImport(node *Node) bool {
	return node.kind == SyntaxKindNamespaceImport
}

// NamedImports

type NamedImports struct {
	NodeBase
	elements []*ImportSpecifierNode
}

func (f *NodeFactory) NewNamedImports(elements []*ImportSpecifierNode) *Node {
	data := &NamedImports{}
	data.elements = elements
	return f.NewNode(SyntaxKindNamedImports, data)
}

func (f *NodeFactory) UpdateNamedImports(node *Node, elements []*ImportSpecifierNode) *Node {
	if n := node.AsNamedImports(); !slices.Equal(elements, n.elements) {
		return f.UpdateNode(f.NewNamedImports(elements), node)
	}
	return node
}

func (node *NamedImports) ForEachChild(v Visitor) bool {
	return visitNodes(v, node.elements)
}

func (node *NamedImports) Accept(v NodeVisitor) *Node {
	return v.VisitNamedImports(node)
}

// ExportAssignment

// This is either an `export =` or an `export default` declaration.
// Unless `isExportEquals` is set, this node was parsed as an `export default`.
type ExportAssignment struct {
	StatementBase
	DeclarationBase
	ModifiersBase
	isExportEquals bool
	expression     *Node
}

func (f *NodeFactory) NewExportAssignment(modifiers *ModifierListNode, isExportEquals bool, expression *Node) *Node {
	data := &ExportAssignment{}
	data.modifiers = modifiers
	data.isExportEquals = isExportEquals
	if isExportEquals {
		data.expression = f.ParenthesizerRules().ParenthesizeRightSideOfBinary(SyntaxKindEqualsToken, nil /*leftSide*/, expression)
	} else {
		data.expression = f.ParenthesizerRules().ParenthesizeExpressionOfExportDefault(expression)
	}
	return f.NewNode(SyntaxKindExportAssignment, data)
}

func (f *NodeFactory) UpdateExportAssignment(node *Node, modifiers *ModifierListNode, expression *Node) *Node {
	if n := node.AsExportAssignment(); modifiers != n.modifiers || expression != n.expression {
		return f.UpdateNode(f.NewExportAssignment(modifiers, n.isExportEquals, expression), node)
	}
	return node
}

func (node *ExportAssignment) ForEachChild(v Visitor) bool {
	return visit(v, node.modifiers) || visit(v, node.expression)
}

func (node *ExportAssignment) Accept(v NodeVisitor) *Node {
	return v.VisitExportAssignment(node)
}

func isExportAssignment(node *Node) bool {
	return node.kind == SyntaxKindExportAssignment
}

// NamespaceExportDeclaration

type NamespaceExportDeclaration struct {
	StatementBase
	DeclarationBase
	ModifiersBase
	name *IdentifierNode
}

func (f *NodeFactory) NewNamespaceExportDeclaration(modifiers *ModifierListNode, name *IdentifierNode) *Node {
	data := &NamespaceExportDeclaration{}
	data.modifiers = modifiers
	data.name = name
	return f.NewNode(SyntaxKindNamespaceExportDeclaration, data)
}

func (f *NodeFactory) UpdateNamespaceExportDeclaration(node *Node, modifiers *ModifierListNode, name *IdentifierNode) *Node {
	if n := node.AsNamespaceExportDeclaration(); modifiers != n.modifiers || name != n.name {
		return f.UpdateNode(f.NewNamespaceExportDeclaration(modifiers, name), node)
	}
	return node
}

func (node *NamespaceExportDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.modifiers) || visit(v, node.name)
}

func (node *NamespaceExportDeclaration) Accept(v NodeVisitor) *Node {
	return v.VisitNamespaceExportDeclaration(node)
}

func (node *NamespaceExportDeclaration) Name() *DeclarationName {
	return node.name
}

func isNamespaceExportDeclaration(node *Node) bool {
	return node.kind == SyntaxKindNamespaceExportDeclaration
}

// ExportDeclaration

type ExportDeclaration struct {
	StatementBase
	DeclarationBase
	ModifiersBase
	isTypeOnly      bool
	exportClause    *NamedExportBindings  // Optional
	moduleSpecifier *Node                 // Optional
	attributes      *ImportAttributesNode // Optional
}

func (f *NodeFactory) NewExportDeclaration(modifiers *ModifierListNode, isTypeOnly bool, exportClause *NamedExportBindings, moduleSpecifier *Node, attributes *ImportAttributesNode) *Node {
	data := &ExportDeclaration{}
	data.modifiers = modifiers
	data.isTypeOnly = isTypeOnly
	data.exportClause = exportClause
	data.moduleSpecifier = moduleSpecifier
	data.attributes = attributes
	return f.NewNode(SyntaxKindExportDeclaration, data)
}

func (f *NodeFactory) UpdateExportDeclaration(node *Node, modifiers *ModifierListNode, isTypeOnly bool, exportClause *NamedExportBindings, moduleSpecifier *Node, attributes *ImportAttributesNode) *Node {
	if n := node.AsExportDeclaration(); modifiers != n.modifiers || exportClause != n.exportClause || moduleSpecifier != n.moduleSpecifier || attributes != n.attributes {
		return f.UpdateNode(f.NewExportDeclaration(modifiers, isTypeOnly, exportClause, moduleSpecifier, attributes), node)
	}
	return node
}

func (node *ExportDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.modifiers) || visit(v, node.exportClause) || visit(v, node.moduleSpecifier) || visit(v, node.attributes)
}

func (node *ExportDeclaration) Accept(v NodeVisitor) *Node {
	return v.VisitExportDeclaration(node)
}

func isExportDeclaration(node *Node) bool {
	return node.kind == SyntaxKindExportDeclaration
}

// NamespaceExport

type NamespaceExport struct {
	NodeBase
	DeclarationBase
	name *ModuleExportName
}

func (f *NodeFactory) NewNamespaceExport(name *ModuleExportName) *Node {
	data := &NamespaceExport{}
	data.name = name
	return f.NewNode(SyntaxKindNamespaceExport, data)
}

func (f *NodeFactory) UpdateNamespaceExport(node *Node, name *ModuleExportName) *Node {
	if n := node.AsNamespaceExport(); name != n.name {
		return f.UpdateNode(f.NewNamespaceExport(name), node)
	}
	return node
}

func (node *NamespaceExport) ForEachChild(v Visitor) bool {
	return visit(v, node.name)
}

func (node *NamespaceExport) Accept(v NodeVisitor) *Node {
	return v.VisitNamespaceExport(node)
}

func (node *NamespaceExport) Name() *DeclarationName {
	return node.name
}

func isNamespaceExport(node *Node) bool {
	return node.kind == SyntaxKindNamespaceExport
}

// NamedExports

type NamedExports struct {
	NodeBase
	elements []*ExportSpecifierNode
}

func (f *NodeFactory) NewNamedExports(elements []*ExportSpecifierNode) *Node {
	data := &NamedExports{}
	data.elements = elements
	return f.NewNode(SyntaxKindNamedExports, data)
}

func (f *NodeFactory) UpdateNamedExports(node *Node, elements []*ExportSpecifierNode) *Node {
	if n := node.AsNamedExports(); !slices.Equal(elements, n.elements) {
		return f.UpdateNode(f.NewNamedExports(elements), node)
	}
	return node
}

func (node *NamedExports) ForEachChild(v Visitor) bool {
	return visitNodes(v, node.elements)
}

func (node *NamedExports) Accept(v NodeVisitor) *Node {
	return v.VisitNamedExports(node)
}

// ExportSpecifier

type ExportSpecifier struct {
	NodeBase
	DeclarationBase
	ExportableBase
	isTypeOnly   bool
	propertyName *ModuleExportName // Optional, name preceding 'as' keyword
	name         *ModuleExportName
}

func (f *NodeFactory) NewExportSpecifier(isTypeOnly bool, propertyName *ModuleExportName, name *ModuleExportName) *Node {
	data := &ExportSpecifier{}
	data.isTypeOnly = isTypeOnly
	data.propertyName = propertyName
	data.name = name
	return f.NewNode(SyntaxKindExportSpecifier, data)
}

func (f *NodeFactory) UpdateExportSpecifier(node *Node, isTypeOnly bool, propertyName *ModuleExportName, name *ModuleExportName) *Node {
	if n := node.AsExportSpecifier(); isTypeOnly != n.isTypeOnly || propertyName != n.propertyName || name != n.name {
		return f.UpdateNode(f.NewExportSpecifier(isTypeOnly, propertyName, name), node)
	}
	return node
}

func (node *ExportSpecifier) ForEachChild(v Visitor) bool {
	return visit(v, node.propertyName) || visit(v, node.name)
}

func (node *ExportSpecifier) Accept(v NodeVisitor) *Node {
	return v.VisitExportSpecifier(node)
}

func (node *ExportSpecifier) Name() *DeclarationName {
	return node.name
}

func isExportSpecifier(node *Node) bool {
	return node.kind == SyntaxKindExportSpecifier
}

// TypeElementBase

type TypeElementBase struct{}

// ClassElementBase

type ClassElementBase struct{}

// NamedMemberBase

type NamedMemberBase struct {
	DeclarationBase
	ModifiersBase
	name         *PropertyName
	postfixToken *TokenNode
}

func (node *NamedMemberBase) DeclarationData() *DeclarationBase { return &node.DeclarationBase }
func (node *NamedMemberBase) Modifiers() *ModifierListNode      { return node.modifiers }
func (node *NamedMemberBase) Name() *DeclarationName            { return node.name }

// CallSignatureDeclaration

type CallSignatureDeclaration struct {
	NodeBase
	DeclarationBase
	FunctionLikeBase
	TypeElementBase
}

func (f *NodeFactory) NewCallSignatureDeclaration(typeParameters *TypeParameterListNode, parameters []*ParameterDeclarationNode, returnType *Node) *Node {
	data := &CallSignatureDeclaration{}
	data.typeParameters = typeParameters
	data.parameters = parameters
	data.returnType = returnType
	return f.NewNode(SyntaxKindCallSignature, data)
}

func (f *NodeFactory) UpdateCallSignatureDeclaration(node *Node, typeParameters *TypeParameterListNode, parameters []*ParameterDeclarationNode, returnType *Node) *Node {
	if n := node.AsCallSignatureDeclaration(); typeParameters != n.typeParameters || !slices.Equal(parameters, n.parameters) || returnType != n.returnType {
		return f.UpdateNode(f.NewCallSignatureDeclaration(typeParameters, parameters, returnType), node)
	}
	return node
}

func (node *CallSignatureDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.typeParameters) || visitNodes(v, node.parameters) || visit(v, node.returnType)
}

func (node *CallSignatureDeclaration) Accept(v NodeVisitor) *Node {
	return v.VisitCallSignatureDeclaration(node)
}

func isCallSignatureDeclaration(node *Node) bool {
	return node.kind == SyntaxKindCallSignature
}

// ConstructSignatureDeclaration

type ConstructSignatureDeclaration struct {
	NodeBase
	DeclarationBase
	FunctionLikeBase
	TypeElementBase
}

func (f *NodeFactory) NewConstructSignatureDeclaration(typeParameters *TypeParameterListNode, parameters []*ParameterDeclarationNode, returnType *Node) *Node {
	data := &ConstructSignatureDeclaration{}
	data.typeParameters = typeParameters
	data.parameters = parameters
	data.returnType = returnType
	return f.NewNode(SyntaxKindConstructSignature, data)
}

func (f *NodeFactory) UpdateConstructSignatureDeclaration(node *Node, typeParameters *TypeParameterListNode, parameters []*ParameterDeclarationNode, returnType *Node) *Node {
	if n := node.AsConstructSignatureDeclaration(); typeParameters != n.typeParameters || !slices.Equal(parameters, n.parameters) || returnType != n.returnType {
		return f.UpdateNode(f.NewConstructSignatureDeclaration(typeParameters, parameters, returnType), node)
	}
	return node
}

func (node *ConstructSignatureDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.typeParameters) || visitNodes(v, node.parameters) || visit(v, node.returnType)
}

func (node *ConstructSignatureDeclaration) Accept(v NodeVisitor) *Node {
	return v.VisitConstructSignatureDeclaration(node)
}

// ConstructorDeclaration

type ConstructorDeclaration struct {
	NodeBase
	DeclarationBase
	ModifiersBase
	FunctionLikeWithBodyBase
	ClassElementBase
	returnFlowNode *FlowNode
}

func (f *NodeFactory) NewConstructorDeclaration(modifiers *ModifierListNode, typeParameters *TypeParameterListNode, parameters []*ParameterDeclarationNode, returnType *Node, body *BlockNode) *Node {
	data := &ConstructorDeclaration{}
	data.modifiers = modifiers
	data.typeParameters = typeParameters
	data.parameters = parameters
	data.returnType = returnType
	data.body = body
	return f.NewNode(SyntaxKindConstructor, data)
}

func (f *NodeFactory) UpdateConstructorDeclaration(node *Node, modifiers *ModifierListNode, typeParameters *TypeParameterListNode, parameters []*ParameterDeclarationNode, returnType *Node, body *BlockNode) *Node {
	if n := node.AsConstructorDeclaration(); modifiers != n.modifiers || typeParameters != n.typeParameters || !slices.Equal(parameters, n.parameters) || returnType != n.returnType || body != n.body {
		return f.UpdateNode(f.NewConstructorDeclaration(modifiers, typeParameters, parameters, returnType, body), node)
	}
	return node
}

func (node *ConstructorDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.modifiers) || visit(v, node.typeParameters) || visitNodes(v, node.parameters) || visit(v, node.returnType) || visit(v, node.body)
}

func (node *ConstructorDeclaration) Accept(v NodeVisitor) *Node {
	return v.VisitConstructorDeclaration(node)
}

func isConstructorDeclaration(node *Node) bool {
	return node.kind == SyntaxKindConstructor
}

// AccessorDeclarationBase

type AccessorDeclarationBase struct {
	NodeBase
	NamedMemberBase
	FunctionLikeWithBodyBase
	FlowNodeBase
	TypeElementBase
	ClassElementBase
	ObjectLiteralElementBase
}

func (node *AccessorDeclarationBase) ForEachChild(v Visitor) bool {
	return visit(v, node.modifiers) || visit(v, node.name) || visit(v, node.typeParameters) || visitNodes(v, node.parameters) ||
		visit(v, node.returnType) || visit(v, node.body)
}

func (node *AccessorDeclarationBase) isAccessorDeclaration() {}

// GetAccessorDeclaration

type GetAccessorDeclaration struct {
	AccessorDeclarationBase
}

func (f *NodeFactory) NewGetAccessorDeclaration(modifiers *ModifierListNode, name *PropertyName, typeParameters *TypeParameterListNode, parameters []*ParameterDeclarationNode, returnType *Node, body *BlockNode) *Node {
	data := &GetAccessorDeclaration{}
	data.modifiers = modifiers
	data.name = name
	data.typeParameters = typeParameters
	data.parameters = parameters
	data.returnType = returnType
	data.body = body
	return f.NewNode(SyntaxKindGetAccessor, data)
}

func (f *NodeFactory) UpdateGetAccessorDeclaration(node *Node, modifiers *ModifierListNode, name *PropertyName, typeParameters *TypeParameterListNode, parameters []*ParameterDeclarationNode, returnType *Node, body *BlockNode) *Node {
	if n := node.AsGetAccessorDeclaration(); modifiers != n.modifiers || name != n.name || typeParameters != n.typeParameters || !slices.Equal(parameters, n.parameters) || returnType != n.returnType || body != n.body {
		return f.UpdateNode(f.NewGetAccessorDeclaration(modifiers, name, typeParameters, parameters, returnType, body), node)
	}
	return node
}

func (node *GetAccessorDeclaration) Accept(v NodeVisitor) *Node {
	return v.VisitGetAccessorDeclaration(node)
}

func isGetAccessorDeclaration(node *Node) bool {
	return node.kind == SyntaxKindGetAccessor
}

// SetAccessorDeclaration

type SetAccessorDeclaration struct {
	AccessorDeclarationBase
}

func (f *NodeFactory) NewSetAccessorDeclaration(modifiers *ModifierListNode, name *PropertyName, typeParameters *TypeParameterListNode, parameters []*ParameterDeclarationNode, returnType *Node, body *BlockNode) *Node {
	data := &SetAccessorDeclaration{}
	data.modifiers = modifiers
	data.name = name
	data.typeParameters = typeParameters
	data.parameters = parameters
	data.returnType = returnType
	data.body = body
	return f.NewNode(SyntaxKindSetAccessor, data)
}

func (f *NodeFactory) UpdateSetAccessorDeclaration(node *Node, modifiers *ModifierListNode, name *PropertyName, typeParameters *TypeParameterListNode, parameters []*ParameterDeclarationNode, returnType *Node, body *BlockNode) *Node {
	if n := node.AsSetAccessorDeclaration(); modifiers != n.modifiers || name != n.name || typeParameters != n.typeParameters || !slices.Equal(parameters, n.parameters) || returnType != n.returnType || body != n.body {
		return f.UpdateNode(f.NewSetAccessorDeclaration(modifiers, name, typeParameters, parameters, returnType, body), node)
	}
	return node
}

func (node *SetAccessorDeclaration) Accept(v NodeVisitor) *Node {
	return v.VisitSetAccessorDeclaration(node)
}

func isSetAccessorDeclaration(node *Node) bool {
	return node.kind == SyntaxKindSetAccessor
}

// IndexSignatureDeclaration

type IndexSignatureDeclaration struct {
	NodeBase
	DeclarationBase
	ModifiersBase
	FunctionLikeBase
	TypeElementBase
	ClassElementBase
}

func (f *NodeFactory) NewIndexSignatureDeclaration(modifiers *Node, parameters []*Node, returnType *Node) *Node {
	data := &IndexSignatureDeclaration{}
	data.modifiers = modifiers
	data.parameters = parameters
	data.returnType = returnType
	return f.NewNode(SyntaxKindIndexSignature, data)
}

func (f *NodeFactory) UpdateIndexSignatureDeclaration(node *Node, modifiers *Node, parameters []*Node, returnType *Node) *Node {
	if n := node.AsIndexSignatureDeclaration(); modifiers != n.modifiers || !slices.Equal(parameters, n.parameters) || returnType != n.returnType {
		return f.UpdateNode(f.NewIndexSignatureDeclaration(modifiers, parameters, returnType), node)
	}
	return node
}

func (node *IndexSignatureDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.modifiers) || visitNodes(v, node.parameters) || visit(v, node.returnType)
}

func (node *IndexSignatureDeclaration) Accept(v NodeVisitor) *Node {
	return v.VisitIndexSignatureDeclaration(node)
}

func isIndexSignatureDeclaration(node *Node) bool {
	return node.kind == SyntaxKindIndexSignature
}

// MethodSignatureDeclaration

type MethodSignatureDeclaration struct {
	NodeBase
	NamedMemberBase
	FunctionLikeBase
	TypeElementBase
}

func (f *NodeFactory) NewMethodSignatureDeclaration(modifiers *Node, name *Node, postfixToken *Node, typeParameters *Node, parameters []*Node, returnType *Node) *Node {
	data := &MethodSignatureDeclaration{}
	data.modifiers = modifiers
	data.name = name
	data.postfixToken = postfixToken
	data.typeParameters = typeParameters
	data.parameters = parameters
	data.returnType = returnType
	return f.NewNode(SyntaxKindMethodSignature, data)
}

func (f *NodeFactory) UpdateMethodSignatureDeclaration(node *Node, modifiers *Node, name *Node, postfixToken *Node, typeParameters *Node, parameters []*Node, returnType *Node) *Node {
	if n := node.AsMethodSignatureDeclaration(); modifiers != n.modifiers || name != n.name || postfixToken != n.postfixToken || typeParameters != n.typeParameters || !slices.Equal(parameters, n.parameters) || returnType != n.returnType {
		return f.UpdateNode(f.NewMethodSignatureDeclaration(modifiers, name, postfixToken, typeParameters, parameters, returnType), node)
	}
	return node
}

func (node *MethodSignatureDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.modifiers) || visit(v, node.name) || visit(v, node.postfixToken) || visit(v, node.typeParameters) ||
		visitNodes(v, node.parameters) || visit(v, node.returnType)
}

func (node *MethodSignatureDeclaration) Accept(v NodeVisitor) *Node {
	return v.VisitMethodSignatureDeclaration(node)
}

func isMethodSignatureDeclaration(node *Node) bool {
	return node.kind == SyntaxKindMethodSignature
}

// MethodSignatureDeclaration

type MethodDeclaration struct {
	NodeBase
	NamedMemberBase
	FunctionLikeWithBodyBase
	FlowNodeBase
	ClassElementBase
	ObjectLiteralElementBase
}

func (f *NodeFactory) NewMethodDeclaration(modifiers *Node, asteriskToken *Node, name *Node, postfixToken *Node, typeParameters *Node, parameters []*Node, returnType *Node, body *Node) *Node {
	data := &MethodDeclaration{}
	data.modifiers = modifiers
	data.asteriskToken = asteriskToken
	data.name = name
	data.postfixToken = postfixToken
	data.typeParameters = typeParameters
	data.parameters = parameters
	data.returnType = returnType
	data.body = body
	return f.NewNode(SyntaxKindMethodDeclaration, data)
}

func (f *NodeFactory) UpdateMethodDeclaration(node *Node, modifiers *Node, asteriskToken *Node, name *Node, postfixToken *Node, typeParameters *Node, parameters []*Node, returnType *Node, body *Node) *Node {
	if n := node.AsMethodDeclaration(); modifiers != n.modifiers || asteriskToken != n.asteriskToken || name != n.name || postfixToken != n.postfixToken || typeParameters != n.typeParameters || !slices.Equal(parameters, n.parameters) || returnType != n.returnType || body != n.body {
		return f.UpdateNode(f.NewMethodDeclaration(modifiers, asteriskToken, name, postfixToken, typeParameters, parameters, returnType, body), node)
	}
	return node
}

func (node *MethodDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.modifiers) || visit(v, node.asteriskToken) || visit(v, node.name) || visit(v, node.postfixToken) ||
		visit(v, node.typeParameters) || visitNodes(v, node.parameters) || visit(v, node.returnType) || visit(v, node.body)
}

func (node *MethodDeclaration) Accept(v NodeVisitor) *Node {
	return v.VisitMethodDeclaration(node)
}

func isMethodDeclaration(node *Node) bool {
	return node.kind == SyntaxKindMethodDeclaration
}

// PropertySignatureDeclaration

type PropertySignatureDeclaration struct {
	NodeBase
	NamedMemberBase
	TypeElementBase
	typeNode    *Node
	initializer *Node // For error reporting purposes
}

func (f *NodeFactory) NewPropertySignatureDeclaration(modifiers *Node, name *Node, postfixToken *Node, typeNode *Node, initializer *Node) *Node {
	data := &PropertySignatureDeclaration{}
	data.modifiers = modifiers
	data.name = name
	data.postfixToken = postfixToken
	data.typeNode = typeNode
	data.initializer = initializer
	return f.NewNode(SyntaxKindPropertySignature, data)
}

func (f *NodeFactory) UpdatePropertySignatureDeclaration(node *Node, modifiers *Node, name *Node, postfixToken *Node, typeNode *Node, initializer *Node) *Node {
	if n := node.AsPropertySignatureDeclaration(); modifiers != n.modifiers || name != n.name || postfixToken != n.postfixToken || typeNode != n.typeNode || initializer != n.initializer {
		return f.UpdateNode(f.NewPropertySignatureDeclaration(modifiers, name, postfixToken, typeNode, initializer), node)
	}
	return node
}

func (node *PropertySignatureDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.modifiers) || visit(v, node.name) || visit(v, node.postfixToken) || visit(v, node.typeNode) || visit(v, node.initializer)
}

func (node *PropertySignatureDeclaration) Accept(v NodeVisitor) *Node {
	return v.VisitPropertySignatureDeclaration(node)
}

func isPropertySignatureDeclaration(node *Node) bool {
	return node.kind == SyntaxKindPropertySignature
}

// PropertyDeclaration

type PropertyDeclaration struct {
	NodeBase
	NamedMemberBase
	ClassElementBase
	typeNode    *Node // Optional
	initializer *Node // Optional
}

func (f *NodeFactory) NewPropertyDeclaration(modifiers *Node, name *Node, postfixToken *Node, typeNode *Node, initializer *Node) *Node {
	data := &PropertyDeclaration{}
	data.modifiers = modifiers
	data.name = name
	data.postfixToken = postfixToken
	data.typeNode = typeNode
	if initializer != nil {
		data.initializer = f.ParenthesizerRules().ParenthesizeExpressionForDisallowedComma(initializer)
	}
	return f.NewNode(SyntaxKindPropertyDeclaration, data)
}

func (f *NodeFactory) UpdatePropertyDeclaration(node *Node, modifiers *Node, name *Node, postfixToken *Node, typeNode *Node, initializer *Node) *Node {
	if n := node.AsPropertyDeclaration(); modifiers != n.modifiers || name != n.name || postfixToken != n.postfixToken || typeNode != n.typeNode || initializer != n.initializer {
		return f.UpdateNode(f.NewPropertyDeclaration(modifiers, name, postfixToken, typeNode, initializer), node)
	}
	return node
}

func (node *PropertyDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.modifiers) || visit(v, node.name) || visit(v, node.postfixToken) || visit(v, node.typeNode) || visit(v, node.initializer)
}

func (node *PropertyDeclaration) Accept(v NodeVisitor) *Node {
	return v.VisitPropertyDeclaration(node)
}

func isPropertyDeclaration(node *Node) bool {
	return node.kind == SyntaxKindPropertyDeclaration
}

// SemicolonClassElement

type SemicolonClassElement struct {
	NodeBase
	DeclarationBase
	ClassElementBase
}

func (f *NodeFactory) NewSemicolonClassElement() *Node {
	return f.NewNode(SyntaxKindSemicolonClassElement, &SemicolonClassElement{})
}

// ClassStaticBlockDeclaration

type ClassStaticBlockDeclaration struct {
	NodeBase
	DeclarationBase
	ModifiersBase
	LocalsContainerBase
	ClassElementBase
	body *Node
}

func (f *NodeFactory) NewClassStaticBlockDeclaration(modifiers *Node, body *Node) *Node {
	data := &ClassStaticBlockDeclaration{}
	data.modifiers = modifiers
	data.body = body
	return f.NewNode(SyntaxKindClassStaticBlockDeclaration, data)
}

func (f *NodeFactory) UpdateClassStaticBlockDeclaration(node *Node, modifiers *Node, body *Node) *Node {
	if n := node.AsClassStaticBlockDeclaration(); modifiers != n.modifiers || body != n.body {
		return f.UpdateNode(f.NewClassStaticBlockDeclaration(modifiers, body), node)
	}
	return node
}

func (node *ClassStaticBlockDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.modifiers) || visit(v, node.body)
}

func (node *ClassStaticBlockDeclaration) Accept(v NodeVisitor) *Node {
	return v.VisitClassStaticBlockDeclaration(node)
}

func isClassStaticBlockDeclaration(node *Node) bool {
	return node.kind == SyntaxKindClassStaticBlockDeclaration
}

// TypeParameterList

type TypeParameterList struct {
	NodeBase
	parameters []*Node
}

func (f *NodeFactory) NewTypeParameterList(parameters []*Node) *Node {
	data := &TypeParameterList{}
	data.parameters = parameters
	return f.NewNode(SyntaxKindTypeParameterList, data)
}

func (f *NodeFactory) UpdateTypeParameterList(node *Node, parameters []*Node) *Node {
	if n := node.AsTypeParameterList(); !slices.Equal(parameters, n.parameters) {
		return f.UpdateNode(f.NewTypeParameterList(parameters), node)
	}
	return node
}

func (node *TypeParameterList) ForEachChild(v Visitor) bool {
	return visitNodes(v, node.parameters)
}

func (node *TypeParameterList) Accept(v NodeVisitor) *Node {
	return v.VisitTypeParameterList(node)
}

func isTypeParameterList(node *Node) bool {
	return node.kind == SyntaxKindTypeParameterList
}

// ExpressionBase

type ExpressionBase struct {
	NodeBase
}

// OmittedExpression

type OmittedExpression struct {
	ExpressionBase
}

func (f *NodeFactory) NewOmittedExpression() *Node {
	return f.NewNode(SyntaxKindOmittedExpression, &OmittedExpression{})
}

// KeywordExpression

type KeywordExpression struct {
	ExpressionBase
	FlowNodeBase // For 'this' and 'super' expressions
}

func (f *NodeFactory) NewKeywordExpression(kind SyntaxKind) *Node {
	return f.NewNode(kind, &KeywordExpression{})
}

// LiteralLikeBase

type LiteralLikeBase struct {
	text string
}

// StringLiteral

type StringLiteral struct {
	ExpressionBase
	LiteralLikeBase
}

func (f *NodeFactory) NewStringLiteral(text string) *Node {
	data := &StringLiteral{}
	data.text = text
	return f.NewNode(SyntaxKindStringLiteral, data)
}

func isStringLiteral(node *Node) bool {
	return node.kind == SyntaxKindStringLiteral
}

// NumericLiteral

type NumericLiteral struct {
	ExpressionBase
	LiteralLikeBase
}

func (f *NodeFactory) NewNumericLiteral(text string) *Node {
	data := &NumericLiteral{}
	data.text = text
	return f.NewNode(SyntaxKindNumericLiteral, data)
}

// BigIntLiteral

type BigIntLiteral struct {
	ExpressionBase
	LiteralLikeBase
}

func (f *NodeFactory) NewBigIntLiteral(text string) *Node {
	data := &BigIntLiteral{}
	data.text = text
	return f.NewNode(SyntaxKindBigIntLiteral, data)
}

// RegularExpressionLiteral

type RegularExpressionLiteral struct {
	ExpressionBase
	LiteralLikeBase
}

func (f *NodeFactory) NewRegularExpressionLiteral(text string) *Node {
	data := &RegularExpressionLiteral{}
	data.text = text
	return f.NewNode(SyntaxKindRegularExpressionLiteral, data)
}

// NoSubstitutionTemplateLiteral

type NoSubstitutionTemplateLiteral struct {
	ExpressionBase
	TemplateLiteralLikeBase
}

func (f *NodeFactory) NewNoSubstitutionTemplateLiteral(text string) *Node {
	data := &NoSubstitutionTemplateLiteral{}
	data.text = text
	return f.NewNode(SyntaxKindNoSubstitutionTemplateLiteral, data)
}

// BinaryExpression

type BinaryExpression struct {
	ExpressionBase
	DeclarationBase
	left          *Node
	operatorToken *Node
	right         *Node
}

func (f *NodeFactory) NewBinaryExpression(left *Node, operatorToken *Node, right *Node) *Node {
	data := &BinaryExpression{}
	data.left = f.ParenthesizerRules().ParenthesizeLeftSideOfBinary(operatorToken.kind, left)
	data.operatorToken = operatorToken
	data.right = f.ParenthesizerRules().ParenthesizeRightSideOfBinary(operatorToken.kind, data.left, right)
	return f.NewNode(SyntaxKindBinaryExpression, data)
}

func (f *NodeFactory) UpdateBinaryExpression(node *Node, left *Node, operatorToken *Node, right *Node) *Node {
	if n := node.AsBinaryExpression(); left != n.left || operatorToken != n.operatorToken || right != n.right {
		return f.UpdateNode(f.NewBinaryExpression(left, operatorToken, right), node)
	}
	return node
}

func (node *BinaryExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.left) || visit(v, node.operatorToken) || visit(v, node.right)
}

func (node *BinaryExpression) Accept(v NodeVisitor) *Node {
	return v.VisitBinaryExpression(node)
}

// PrefixUnaryExpression

type PrefixUnaryExpression struct {
	ExpressionBase
	operator SyntaxKind
	operand  *Node
}

func (f *NodeFactory) NewPrefixUnaryExpression(operator SyntaxKind, operand *Node) *Node {
	data := &PrefixUnaryExpression{}
	data.operator = operator
	data.operand = f.ParenthesizerRules().ParenthesizeOperandOfPrefixUnary(operand)
	return f.NewNode(SyntaxKindPrefixUnaryExpression, data)
}

func (f *NodeFactory) UpdatePrefixUnaryExpression(node *Node, operand *Node) *Node {
	if n := node.AsPrefixUnaryExpression(); operand != n.operand {
		return f.UpdateNode(f.NewPrefixUnaryExpression(n.operator, operand), node)
	}
	return node
}

func (node *PrefixUnaryExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.operand)
}

func (node *PrefixUnaryExpression) Accept(v NodeVisitor) *Node {
	return v.VisitPrefixUnaryExpression(node)
}

func isPrefixUnaryExpression(node *Node) bool {
	return node.kind == SyntaxKindPrefixUnaryExpression
}

// PostfixUnaryExpression

type PostfixUnaryExpression struct {
	ExpressionBase
	operand  *Node
	operator SyntaxKind
}

func (f *NodeFactory) NewPostfixUnaryExpression(operand *Node, operator SyntaxKind) *Node {
	data := &PostfixUnaryExpression{}
	data.operand = f.ParenthesizerRules().ParenthesizeOperandOfPostfixUnary(operand)
	data.operator = operator
	return f.NewNode(SyntaxKindPostfixUnaryExpression, data)
}

func (f *NodeFactory) UpdatePostfixUnaryExpression(node *Node, operand *Node) *Node {
	if n := node.AsPostfixUnaryExpression(); operand != n.operand {
		return f.UpdateNode(f.NewPostfixUnaryExpression(operand, n.operator), node)
	}
	return node
}

func (node *PostfixUnaryExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.operand)
}

func (node *PostfixUnaryExpression) Accept(v NodeVisitor) *Node {
	return v.VisitPostfixUnaryExpression(node)
}

// YieldExpression

type YieldExpression struct {
	ExpressionBase
	asteriskToken *Node
	expression    *Node // Optional
}

func (f *NodeFactory) NewYieldExpression(asteriskToken *Node, expression *Node) *Node {
	data := &YieldExpression{}
	data.asteriskToken = asteriskToken
	data.expression = f.ParenthesizerRules().ParenthesizeExpressionForDisallowedComma(expression)
	return f.NewNode(SyntaxKindYieldExpression, data)
}

func (f *NodeFactory) UpdateYieldExpression(node *Node, asteriskToken *Node, expression *Node) *Node {
	if n := node.AsYieldExpression(); asteriskToken != n.asteriskToken || expression != n.expression {
		return f.UpdateNode(f.NewYieldExpression(asteriskToken, expression), node)
	}
	return node
}

func (node *YieldExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.asteriskToken) || visit(v, node.expression)
}

func (node *YieldExpression) Accept(v NodeVisitor) *Node {
	return v.VisitYieldExpression(node)
}

// ArrowFunction

type ArrowFunction struct {
	ExpressionBase
	DeclarationBase
	ModifiersBase
	FunctionLikeWithBodyBase
	FlowNodeBase
	equalsGreaterThanToken *Node
}

func (f *NodeFactory) NewArrowFunction(modifiers *Node, typeParameters *Node, parameters []*Node, returnType *Node, equalsGreaterThanToken *Node, body *Node) *Node {
	data := &ArrowFunction{}
	data.modifiers = modifiers
	data.typeParameters = typeParameters
	data.parameters = parameters
	data.returnType = returnType
	data.equalsGreaterThanToken = equalsGreaterThanToken
	data.body = f.ParenthesizerRules().ParenthesizeConciseBodyOfArrowFunction(body)
	return f.NewNode(SyntaxKindArrowFunction, data)
}

func (f *NodeFactory) UpdateArrowFunction(node *Node, modifiers *Node, typeParameters *Node, parameters []*Node, returnType *Node, equalsGreaterThanToken *Node, body *Node) *Node {
	if n := node.AsArrowFunction(); modifiers != n.modifiers || typeParameters != n.typeParameters || !slices.Equal(parameters, n.parameters) || returnType != n.returnType || equalsGreaterThanToken != n.equalsGreaterThanToken || body != n.body {
		return f.UpdateNode(f.NewArrowFunction(modifiers, typeParameters, parameters, returnType, equalsGreaterThanToken, body), node)
	}
	return node
}

func (node *ArrowFunction) ForEachChild(v Visitor) bool {
	return visit(v, node.modifiers) || visit(v, node.typeParameters) || visitNodes(v, node.parameters) ||
		visit(v, node.returnType) || visit(v, node.equalsGreaterThanToken) || visit(v, node.body)
}

func (node *ArrowFunction) Accept(v NodeVisitor) *Node {
	return v.VisitArrowFunction(node)
}

func (node *ArrowFunction) Name() *DeclarationName {
	return nil
}

func isArrowFunction(node *Node) bool {
	return node.kind == SyntaxKindArrowFunction
}

// FunctionExpression

type FunctionExpression struct {
	ExpressionBase
	DeclarationBase
	ModifiersBase
	FunctionLikeWithBodyBase
	FlowNodeBase
	name           *Node // Optional
	returnFlowNode *FlowNode
}

func (f *NodeFactory) NewFunctionExpression(modifiers *Node, asteriskToken *Node, name *Node, typeParameters *Node, parameters []*Node, returnType *Node, body *Node) *Node {
	data := &FunctionExpression{}
	data.modifiers = modifiers
	data.asteriskToken = asteriskToken
	data.name = name
	data.typeParameters = typeParameters
	data.parameters = parameters
	data.returnType = returnType
	data.body = body
	return f.NewNode(SyntaxKindFunctionExpression, data)
}

func (f *NodeFactory) UpdateFunctionExpression(node *Node, modifiers *Node, asteriskToken *Node, name *Node, typeParameters *Node, parameters []*Node, returnType *Node, body *Node) *Node {
	if n := node.AsFunctionExpression(); modifiers != n.modifiers || asteriskToken != n.asteriskToken || name != n.name || typeParameters != n.typeParameters || !slices.Equal(parameters, n.parameters) || returnType != n.returnType || body != n.body {
		return f.UpdateNode(f.NewFunctionExpression(modifiers, asteriskToken, name, typeParameters, parameters, returnType, body), node)
	}
	return node
}

func (node *FunctionExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.modifiers) || visit(v, node.asteriskToken) || visit(v, node.name) || visit(v, node.typeParameters) ||
		visitNodes(v, node.parameters) || visit(v, node.returnType) || visit(v, node.body)
}

func (node *FunctionExpression) Accept(v NodeVisitor) *Node {
	return v.VisitFunctionExpression(node)
}

func (node *FunctionExpression) Name() *DeclarationName {
	return node.name
}

func isFunctionExpression(node *Node) bool {
	return node.kind == SyntaxKindFunctionExpression
}

// AsExpression

type AsExpression struct {
	ExpressionBase
	expression *Node
	typeNode   *Node
}

func (f *NodeFactory) NewAsExpression(expression *Node, typeNode *Node) *Node {
	data := &AsExpression{}
	data.expression = expression
	data.typeNode = typeNode
	return f.NewNode(SyntaxKindAsExpression, data)
}

func (f *NodeFactory) UpdateAsExpression(node *Node, expression *Node, typeNode *Node) *Node {
	if n := node.AsAsExpression(); expression != n.expression || typeNode != n.typeNode {
		return f.UpdateNode(f.NewAsExpression(expression, typeNode), node)
	}
	return node
}

func (node *AsExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.expression) || visit(v, node.typeNode)
}

func (node *AsExpression) Accept(v NodeVisitor) *Node {
	return v.VisitAsExpression(node)
}

// SatisfiesExpression

type SatisfiesExpression struct {
	ExpressionBase
	expression *Node
	typeNode   *Node
}

func (f *NodeFactory) NewSatisfiesExpression(expression *Node, typeNode *Node) *Node {
	data := &SatisfiesExpression{}
	data.expression = expression
	data.typeNode = typeNode
	return f.NewNode(SyntaxKindSatisfiesExpression, data)
}

func (f *NodeFactory) UpdateSatisfiesExpression(node *Node, expression *Node, typeNode *Node) *Node {
	if n := node.AsSatisfiesExpression(); expression != n.expression || typeNode != n.typeNode {
		return f.UpdateNode(f.NewSatisfiesExpression(expression, typeNode), node)
	}
	return node
}

func (node *SatisfiesExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.expression) || visit(v, node.typeNode)
}

func (node *SatisfiesExpression) Accept(v NodeVisitor) *Node {
	return v.VisitSatisfiesExpression(node)
}

// ConditionalExpression

type ConditionalExpression struct {
	ExpressionBase
	condition     *Node
	questionToken *Node
	whenTrue      *Node
	colonToken    *Node
	whenFalse     *Node
}

func (f *NodeFactory) NewConditionalExpression(condition *Node, questionToken *Node, whenTrue *Node, colonToken *Node, whenFalse *Node) *Node {
	data := &ConditionalExpression{}
	data.condition = f.ParenthesizerRules().ParenthesizeConditionOfConditionalExpression(condition)
	data.questionToken = questionToken
	data.whenTrue = f.ParenthesizerRules().ParenthesizeBranchOfConditionalExpression(whenTrue)
	data.colonToken = colonToken
	data.whenFalse = f.ParenthesizerRules().ParenthesizeBranchOfConditionalExpression(whenFalse)
	return f.NewNode(SyntaxKindConditionalExpression, data)
}

func (f *NodeFactory) UpdateConditionalExpression(node *Node, condition *Node, questionToken *Node, whenTrue *Node, colonToken *Node, whenFalse *Node) *Node {
	if n := node.AsConditionalExpression(); condition != n.condition || questionToken != n.questionToken || whenTrue != n.whenTrue || colonToken != n.colonToken || whenFalse != n.whenFalse {
		return f.UpdateNode(f.NewConditionalExpression(condition, questionToken, whenTrue, colonToken, whenFalse), node)
	}
	return node
}

func (node *ConditionalExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.condition) || visit(v, node.questionToken) || visit(v, node.whenTrue) ||
		visit(v, node.colonToken) || visit(v, node.whenFalse)
}

func (node *ConditionalExpression) Accept(v NodeVisitor) *Node {
	return v.VisitConditionalExpression(node)
}

// PropertyAccessExpression

type PropertyAccessExpression struct {
	ExpressionBase
	FlowNodeBase
	expression       *Node
	questionDotToken *Node
	name             *Node
}

func (f *NodeFactory) NewPropertyAccessExpression(expression *Node, questionDotToken *Node, name *Node, flags NodeFlags) *Node {
	data := &PropertyAccessExpression{}
	data.expression = f.ParenthesizerRules().ParenthesizeLeftSideOfAccess(expression, flags&NodeFlagsOptionalChain != 0)
	data.questionDotToken = questionDotToken
	data.name = name
	node := f.NewNode(SyntaxKindPropertyAccessExpression, data)
	node.flags |= flags & NodeFlagsOptionalChain
	return node
}

func (f *NodeFactory) UpdatePropertyAccessExpression(node *Node, expression *Node, questionDotToken *Node, name *Node) *Node {
	if n := node.AsPropertyAccessExpression(); expression != n.expression || questionDotToken != n.questionDotToken || name != n.name {
		return f.UpdateNode(f.NewPropertyAccessExpression(expression, questionDotToken, name, n.flags), node)
	}
	return node
}

func (node *PropertyAccessExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.expression) || visit(v, node.questionDotToken) || visit(v, node.name)
}

func (node *PropertyAccessExpression) Accept(v NodeVisitor) *Node {
	return v.VisitPropertyAccessExpression(node)
}

func (node *PropertyAccessExpression) Name() *DeclarationName { return node.name }

func isPropertyAccessExpression(node *Node) bool {
	return node.kind == SyntaxKindPropertyAccessExpression
}

// ElementAccessExpression

type ElementAccessExpression struct {
	ExpressionBase
	FlowNodeBase
	expression         *Node
	questionDotToken   *Node
	argumentExpression *Node
}

func (f *NodeFactory) NewElementAccessExpression(expression *Node, questionDotToken *Node, argumentExpression *Node, flags NodeFlags) *Node {
	data := &ElementAccessExpression{}
	data.expression = f.ParenthesizerRules().ParenthesizeLeftSideOfAccess(expression, flags&NodeFlagsOptionalChain != 0)
	data.questionDotToken = questionDotToken
	data.argumentExpression = argumentExpression
	node := f.NewNode(SyntaxKindElementAccessExpression, data)
	node.flags |= flags & NodeFlagsOptionalChain
	return node
}

func (f *NodeFactory) UpdateElementAccessExpression(node *Node, expression *Node, questionDotToken *Node, argumentExpression *Node) *Node {
	if n := node.AsElementAccessExpression(); expression != n.expression || questionDotToken != n.questionDotToken || argumentExpression != n.argumentExpression {
		return f.UpdateNode(f.NewElementAccessExpression(expression, questionDotToken, argumentExpression, n.flags), node)
	}
	return node
}

func (node *ElementAccessExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.expression) || visit(v, node.questionDotToken) || visit(v, node.argumentExpression)
}

func (node *ElementAccessExpression) Accept(v NodeVisitor) *Node {
	return v.VisitElementAccessExpression(node)
}

func isElementAccessExpression(node *Node) bool {
	return node.kind == SyntaxKindElementAccessExpression
}

// CallExpression

type CallExpression struct {
	ExpressionBase
	expression       *Node
	questionDotToken *Node
	typeArguments    *Node
	arguments        []*Node
}

func (f *NodeFactory) NewCallExpression(expression *Node, questionDotToken *Node, typeArguments *Node, arguments []*Node, flags NodeFlags) *Node {
	data := &CallExpression{}
	data.expression = f.ParenthesizerRules().ParenthesizeLeftSideOfAccess(expression, flags&NodeFlagsOptionalChain != 0)
	data.questionDotToken = questionDotToken
	data.typeArguments = typeArguments
	data.arguments = f.ParenthesizerRules().ParenthesizeExpressionsOfCommaDelimitedList(arguments)
	node := f.NewNode(SyntaxKindCallExpression, data)
	node.flags |= flags & NodeFlagsOptionalChain
	return node
}

func (f *NodeFactory) UpdateCallExpression(node *Node, expression *Node, questionDotToken *Node, typeArguments *Node, arguments []*Node) *Node {
	if n := node.AsCallExpression(); expression != n.expression || questionDotToken != n.questionDotToken || typeArguments != n.typeArguments || !slices.Equal(arguments, n.arguments) {
		return f.UpdateNode(f.NewCallExpression(expression, questionDotToken, typeArguments, arguments, n.flags), node)
	}
	return node
}

func (node *CallExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.expression) || visit(v, node.questionDotToken) || visit(v, node.typeArguments) || visitNodes(v, node.arguments)
}

func (node *CallExpression) Accept(v NodeVisitor) *Node {
	return v.VisitCallExpression(node)
}

func isCallExpression(node *Node) bool {
	return node.kind == SyntaxKindCallExpression
}

// NewExpression

type NewExpression struct {
	ExpressionBase
	expression    *Node
	typeArguments *Node
	arguments     []*Node
}

func (f *NodeFactory) NewNewExpression(expression *Node, typeArguments *Node, arguments []*Node) *Node {
	data := &NewExpression{}
	data.expression = f.ParenthesizerRules().ParenthesizeExpressionOfNew(expression)
	data.typeArguments = typeArguments
	data.arguments = f.ParenthesizerRules().ParenthesizeExpressionsOfCommaDelimitedList(arguments)
	return f.NewNode(SyntaxKindNewExpression, data)
}

func (f *NodeFactory) UpdateNewExpression(node *Node, expression *Node, typeArguments *Node, arguments []*Node) *Node {
	if n := node.AsNewExpression(); expression != n.expression || typeArguments != n.typeArguments || !slices.Equal(arguments, n.arguments) {
		return f.UpdateNode(f.NewNewExpression(expression, typeArguments, arguments), node)
	}
	return node
}

func (node *NewExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.expression) || visit(v, node.typeArguments) || visitNodes(v, node.arguments)
}

func (node *NewExpression) Accept(v NodeVisitor) *Node {
	return v.VisitNewExpression(node)
}

func isNewExpression(node *Node) bool {
	return node.kind == SyntaxKindNewExpression
}

// MetaProperty

type MetaProperty struct {
	ExpressionBase
	FlowNodeBase
	keywordToken SyntaxKind
	name         *Node
}

func (f *NodeFactory) NewMetaProperty(keywordToken SyntaxKind, name *Node) *Node {
	data := &MetaProperty{}
	data.keywordToken = keywordToken
	data.name = name
	return f.NewNode(SyntaxKindMetaProperty, data)
}

func (f *NodeFactory) UpdateMetaProperty(node *Node, name *Node) *Node {
	if n := node.AsMetaProperty(); name != n.name {
		return f.UpdateNode(f.NewMetaProperty(n.keywordToken, name), node)
	}
	return node
}

func (node *MetaProperty) ForEachChild(v Visitor) bool {
	return visit(v, node.name)
}

func (node *MetaProperty) Accept(v NodeVisitor) *Node {
	return v.VisitMetaProperty(node)
}

func isMetaProperty(node *Node) bool {
	return node.kind == SyntaxKindMetaProperty
}

// NonNullExpression

type NonNullExpression struct {
	ExpressionBase
	expression *Node
}

func (f *NodeFactory) NewNonNullExpression(expression *Node, flags NodeFlags) *Node {
	data := &NonNullExpression{}
	data.expression = f.ParenthesizerRules().ParenthesizeLeftSideOfAccess(expression /*optionalChain*/, flags&NodeFlagsOptionalChain != 0)
	node := f.NewNode(SyntaxKindNonNullExpression, data)
	node.flags |= flags & NodeFlagsOptionalChain
	return node
}

func (f *NodeFactory) UpdateNonNullExpression(node *Node, expression *Node) *Node {
	if n := node.AsNonNullExpression(); expression != n.expression {
		return f.UpdateNode(f.NewNonNullExpression(expression, n.flags), node)
	}
	return node
}

func (node *NonNullExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.expression)
}

func (node *NonNullExpression) Accept(v NodeVisitor) *Node {
	return v.VisitNonNullExpression(node)
}

// SpreadElement

type SpreadElement struct {
	ExpressionBase
	expression *Node
}

func (f *NodeFactory) NewSpreadElement(expression *Node) *Node {
	data := &SpreadElement{}
	data.expression = f.ParenthesizerRules().ParenthesizeExpressionForDisallowedComma(expression)
	return f.NewNode(SyntaxKindSpreadElement, data)
}

func (f *NodeFactory) UpdateSpreadElement(node *Node, expression *Node) *Node {
	if n := node.AsSpreadElement(); expression != n.expression {
		return f.UpdateNode(f.NewSpreadElement(expression), node)
	}
	return node
}

func (node *SpreadElement) ForEachChild(v Visitor) bool {
	return visit(v, node.expression)
}

func (node *SpreadElement) Accept(v NodeVisitor) *Node {
	return v.VisitSpreadElement(node)
}

// TemplateExpression

type TemplateExpression struct {
	ExpressionBase
	head          *Node
	templateSpans []*Node
}

func (f *NodeFactory) NewTemplateExpression(head *Node, templateSpans []*Node) *Node {
	data := &TemplateExpression{}
	data.head = head
	data.templateSpans = templateSpans
	return f.NewNode(SyntaxKindTemplateExpression, data)
}

func (f *NodeFactory) UpdateTemplateExpression(node *Node, head *Node, templateSpans []*Node) *Node {
	if n := node.AsTemplateExpression(); head != n.head || !slices.Equal(templateSpans, n.templateSpans) {
		return f.UpdateNode(f.NewTemplateExpression(head, templateSpans), node)
	}
	return node
}

func (node *TemplateExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.head) || visitNodes(v, node.templateSpans)
}

func (node *TemplateExpression) Accept(v NodeVisitor) *Node {
	return v.VisitTemplateExpression(node)
}

// TemplateLiteralTypeSpan

type TemplateSpan struct {
	NodeBase
	expression *Node
	literal    *Node
}

func (f *NodeFactory) NewTemplateSpan(expression *Node, literal *Node) *Node {
	data := &TemplateSpan{}
	data.expression = expression
	data.literal = literal
	return f.NewNode(SyntaxKindTemplateSpan, data)
}

func (f *NodeFactory) UpdateTemplateSpan(node *Node, expression *Node, literal *Node) *Node {
	if n := node.AsTemplateSpan(); expression != n.expression || literal != n.literal {
		return f.UpdateNode(f.NewTemplateSpan(expression, literal), node)
	}
	return node
}

func (node *TemplateSpan) ForEachChild(v Visitor) bool {
	return visit(v, node.expression) || visit(v, node.literal)
}

func (node *TemplateSpan) Accept(v NodeVisitor) *Node {
	return v.VisitTemplateSpan(node)
}

// TaggedTemplateExpression

type TaggedTemplateExpression struct {
	ExpressionBase
	tag              *Node
	questionDotToken *Node // For error reporting purposes only
	typeArguments    *Node
	template         *Node
}

func (f *NodeFactory) NewTaggedTemplateExpression(tag *Node, questionDotToken *Node, typeArguments *Node, template *Node, flags NodeFlags) *Node {
	data := &TaggedTemplateExpression{}
	data.tag = f.ParenthesizerRules().ParenthesizeLeftSideOfAccess(tag, false /*optionalChain*/)
	data.questionDotToken = questionDotToken
	data.typeArguments = typeArguments
	data.template = template
	node := f.NewNode(SyntaxKindTaggedTemplateExpression, data)
	node.flags |= flags & NodeFlagsOptionalChain
	return node
}

func (f *NodeFactory) UpdateTaggedTemplateExpression(node *Node, tag *Node, questionDotToken *Node, typeArguments *Node, template *Node) *Node {
	if n := node.AsTaggedTemplateExpression(); tag != n.tag || questionDotToken != n.questionDotToken || typeArguments != n.typeArguments || template != n.template {
		return f.UpdateNode(f.NewTaggedTemplateExpression(tag, questionDotToken, typeArguments, template, n.flags), node)
	}
	return node
}

func (node *TaggedTemplateExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.tag) || visit(v, node.questionDotToken) || visit(v, node.typeArguments) || visit(v, node.template)
}

func (node *TaggedTemplateExpression) Accept(v NodeVisitor) *Node {
	return v.VisitTaggedTemplateExpression(node)
}

// ParenthesizedExpression

type ParenthesizedExpression struct {
	ExpressionBase
	expression *Node
}

func (f *NodeFactory) NewParenthesizedExpression(expression *Node) *Node {
	data := &ParenthesizedExpression{}
	data.expression = expression
	return f.NewNode(SyntaxKindParenthesizedExpression, data)
}

func (f *NodeFactory) UpdateParenthesizedExpression(node *Node, expression *Node) *Node {
	if n := node.AsParenthesizedExpression(); expression != n.expression {
		return f.UpdateNode(f.NewParenthesizedExpression(expression), node)
	}
	return node
}

func (node *ParenthesizedExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.expression)
}

func (node *ParenthesizedExpression) Accept(v NodeVisitor) *Node {
	return v.VisitParenthesizedExpression(node)
}

func isParenthesizedExpression(node *Node) bool {
	return node.kind == SyntaxKindParenthesizedExpression
}

// ArrayLiteralExpression

type ArrayLiteralExpression struct {
	ExpressionBase
	elements  []*Node
	multiLine bool
}

func (f *NodeFactory) NewArrayLiteralExpression(elements []*Node, multiLine bool) *Node {
	data := &ArrayLiteralExpression{}
	data.elements = f.ParenthesizerRules().ParenthesizeExpressionsOfCommaDelimitedList(elements)
	data.multiLine = multiLine
	return f.NewNode(SyntaxKindArrayLiteralExpression, data)
}

func (f *NodeFactory) UpdateArrayLiteralExpression(node *Node, elements []*Node) *Node {
	if n := node.AsArrayLiteralExpression(); !slices.Equal(elements, n.elements) {
		return f.UpdateNode(f.NewArrayLiteralExpression(elements, n.multiLine), node)
	}
	return node
}

func (node *ArrayLiteralExpression) ForEachChild(v Visitor) bool {
	return visitNodes(v, node.elements)
}

func (node *ArrayLiteralExpression) Accept(v NodeVisitor) *Node {
	return v.VisitArrayLiteralExpression(node)
}

// ObjectLiteralExpression

type ObjectLiteralExpression struct {
	ExpressionBase
	DeclarationBase
	properties []*Node
	multiLine  bool
}

func (f *NodeFactory) NewObjectLiteralExpression(properties []*Node, multiLine bool) *Node {
	data := &ObjectLiteralExpression{}
	data.properties = properties
	data.multiLine = multiLine
	return f.NewNode(SyntaxKindObjectLiteralExpression, data)

}

func (f *NodeFactory) UpdateObjectLiteralExpression(node *Node, properties []*Node) *Node {
	if n := node.AsObjectLiteralExpression(); !slices.Equal(properties, n.properties) {
		return f.UpdateNode(f.NewObjectLiteralExpression(properties, n.multiLine), node)
	}
	return node
}

func (node *ObjectLiteralExpression) ForEachChild(v Visitor) bool {
	return visitNodes(v, node.properties)
}

func (node *ObjectLiteralExpression) Accept(v NodeVisitor) *Node {
	return v.VisitObjectLiteralExpression(node)
}

func isObjectLiteralExpression(node *Node) bool {
	return node.kind == SyntaxKindObjectLiteralExpression
}

// ObjectLiteralElementBase

type ObjectLiteralElementBase struct{}

// SpreadAssignment

type SpreadAssignment struct {
	NodeBase
	ObjectLiteralElementBase
	expression *Node
}

func (f *NodeFactory) NewSpreadAssignment(expression *Node) *Node {
	data := &SpreadAssignment{}
	data.expression = f.ParenthesizerRules().ParenthesizeExpressionForDisallowedComma(expression)
	return f.NewNode(SyntaxKindSpreadAssignment, data)
}

func (f *NodeFactory) UpdateSpreadAssignment(node *Node, expression *Node) *Node {
	if n := node.AsSpreadAssignment(); expression != n.expression {
		return f.UpdateNode(f.NewSpreadAssignment(expression), node)
	}
	return node
}

func (node *SpreadAssignment) ForEachChild(v Visitor) bool {
	return visit(v, node.expression)
}

func (node *SpreadAssignment) Accept(v NodeVisitor) *Node {
	return v.VisitSpreadAssignment(node)
}

// PropertyAssignment

type PropertyAssignment struct {
	NodeBase
	NamedMemberBase
	ObjectLiteralElementBase
	initializer *Node
}

func (f *NodeFactory) NewPropertyAssignment(modifiers *Node, name *Node, postfixToken *Node, initializer *Node) *Node {
	data := &PropertyAssignment{}
	data.modifiers = modifiers
	data.name = name
	data.postfixToken = postfixToken
	data.initializer = f.ParenthesizerRules().ParenthesizeExpressionForDisallowedComma(initializer)
	return f.NewNode(SyntaxKindPropertyAssignment, data)
}

func (f *NodeFactory) UpdatePropertyAssignment(node *Node, modifiers *Node, name *Node, postfixToken *Node, initializer *Node) *Node {
	if n := node.AsPropertyAssignment(); modifiers != n.modifiers || name != n.name || postfixToken != n.postfixToken || initializer != n.initializer {
		return f.UpdateNode(f.NewPropertyAssignment(modifiers, name, postfixToken, initializer), node)
	}
	return node
}

func (node *PropertyAssignment) ForEachChild(v Visitor) bool {
	return visit(v, node.modifiers) || visit(v, node.name) || visit(v, node.postfixToken) || visit(v, node.initializer)
}

func (node *PropertyAssignment) Accept(v NodeVisitor) *Node {
	return v.VisitPropertyAssignment(node)
}

func isPropertyAssignment(node *Node) bool {
	return node.kind == SyntaxKindPropertyAssignment
}

// ShorthandPropertyAssignment

type ShorthandPropertyAssignment struct {
	NodeBase
	NamedMemberBase
	ObjectLiteralElementBase
	objectAssignmentInitializer *Node // Optional
}

func (f *NodeFactory) NewShorthandPropertyAssignment(modifiers *Node, name *Node, postfixToken *Node, objectAssignmentInitializer *Node) *Node {
	data := &ShorthandPropertyAssignment{}
	data.modifiers = modifiers
	data.name = name
	data.postfixToken = postfixToken
	if objectAssignmentInitializer != nil {
		data.objectAssignmentInitializer = f.ParenthesizerRules().ParenthesizeExpressionForDisallowedComma(objectAssignmentInitializer)
	}
	return f.NewNode(SyntaxKindShorthandPropertyAssignment, data)
}

func (f *NodeFactory) UpdateShorthandPropertyAssignment(node *Node, modifiers *Node, name *Node, postfixToken *Node, objectAssignmentInitializer *Node) *Node {
	if n := node.AsShorthandPropertyAssignment(); modifiers != n.modifiers || name != n.name || postfixToken != n.postfixToken || objectAssignmentInitializer != n.objectAssignmentInitializer {
		return f.UpdateNode(f.NewShorthandPropertyAssignment(modifiers, name, postfixToken, objectAssignmentInitializer), node)
	}
	return node
}

func (node *ShorthandPropertyAssignment) ForEachChild(v Visitor) bool {
	return visit(v, node.modifiers) || visit(v, node.name) || visit(v, node.postfixToken) || visit(v, node.objectAssignmentInitializer)
}

func (node *ShorthandPropertyAssignment) Accept(v NodeVisitor) *Node {
	return v.VisitShorthandPropertyAssignment(node)
}

func isShorthandPropertyAssignment(node *Node) bool {
	return node.kind == SyntaxKindShorthandPropertyAssignment
}

// DeleteExpression

type DeleteExpression struct {
	ExpressionBase
	expression *Node
}

func (f *NodeFactory) NewDeleteExpression(expression *Node) *Node {
	data := &DeleteExpression{}
	data.expression = f.ParenthesizerRules().ParenthesizeOperandOfPrefixUnary(expression)
	return f.NewNode(SyntaxKindDeleteExpression, data)

}

func (f *NodeFactory) UpdateDeleteExpression(node *Node, expression *Node) *Node {
	if n := node.AsDeleteExpression(); expression != n.expression {
		return f.UpdateNode(f.NewDeleteExpression(expression), node)
	}
	return node
}

func (node *DeleteExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.expression)
}

func (node *DeleteExpression) Accept(v NodeVisitor) *Node {
	return v.VisitDeleteExpression(node)
}

// TypeOfExpression

type TypeOfExpression struct {
	ExpressionBase
	expression *Node
}

func (f *NodeFactory) NewTypeOfExpression(expression *Node) *Node {
	data := &TypeOfExpression{}
	data.expression = f.ParenthesizerRules().ParenthesizeOperandOfPrefixUnary(expression)
	return f.NewNode(SyntaxKindTypeOfExpression, data)
}

func (f *NodeFactory) UpdateTypeOfExpression(node *Node, expression *Node) *Node {
	if n := node.AsTypeOfExpression(); expression != n.expression {
		return f.UpdateNode(f.NewTypeOfExpression(expression), node)
	}
	return node
}

func (node *TypeOfExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.expression)
}

func (node *TypeOfExpression) Accept(v NodeVisitor) *Node {
	return v.VisitTypeOfExpression(node)
}

func isTypeOfExpression(node *Node) bool {
	return node.kind == SyntaxKindTypeOfExpression
}

// VoidExpression

type VoidExpression struct {
	ExpressionBase
	expression *Node
}

func (f *NodeFactory) NewVoidExpression(expression *Node) *Node {
	data := &VoidExpression{}
	data.expression = f.ParenthesizerRules().ParenthesizeOperandOfPrefixUnary(expression)
	return f.NewNode(SyntaxKindVoidExpression, data)
}

func (f *NodeFactory) UpdateVoidExpression(node *Node, expression *Node) *Node {
	if n := node.AsVoidExpression(); expression != n.expression {
		return f.UpdateNode(f.NewVoidExpression(expression), node)
	}
	return node
}

func (node *VoidExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.expression)
}

func (node *VoidExpression) Accept(v NodeVisitor) *Node {
	return v.VisitVoidExpression(node)
}

// AwaitExpression

type AwaitExpression struct {
	ExpressionBase
	expression *Node
}

func (f *NodeFactory) NewAwaitExpression(expression *Node) *Node {
	data := &AwaitExpression{}
	data.expression = f.ParenthesizerRules().ParenthesizeOperandOfPrefixUnary(expression)
	return f.NewNode(SyntaxKindAwaitExpression, data)
}

func (f *NodeFactory) UpdateAwaitExpression(node *Node, expression *Node) *Node {
	if n := node.AsAwaitExpression(); expression != n.expression {
		return f.UpdateNode(f.NewAwaitExpression(expression), node)
	}
	return node
}

func (node *AwaitExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.expression)
}

func (node *AwaitExpression) Accept(v NodeVisitor) *Node {
	return v.VisitAwaitExpression(node)
}

// TypeAssertion

type TypeAssertion struct {
	ExpressionBase
	typeNode   *Node
	expression *Node
}

func (f *NodeFactory) NewTypeAssertion(typeNode *Node, expression *Node) *Node {
	data := &TypeAssertion{}
	data.typeNode = typeNode
	data.expression = f.ParenthesizerRules().ParenthesizeOperandOfPrefixUnary(expression)
	return f.NewNode(SyntaxKindTypeAssertionExpression, data)
}

func (f *NodeFactory) UpdateTypeAssertion(node *Node, typeNode *Node, expression *Node) *Node {
	if n := node.AsTypeAssertion(); typeNode != n.typeNode || expression != n.expression {
		return f.UpdateNode(f.NewTypeAssertion(typeNode, expression), node)
	}
	return node
}

func (node *TypeAssertion) ForEachChild(v Visitor) bool {
	return visit(v, node.typeNode) || visit(v, node.expression)
}

func (node *TypeAssertion) Accept(v NodeVisitor) *Node {
	return v.VisitTypeAssertion(node)
}

// TypeNodeBase

type TypeNodeBase struct {
	NodeBase
}

// KeywordTypeNode

type KeywordTypeNode struct {
	TypeNodeBase
}

func (f *NodeFactory) NewKeywordTypeNode(kind SyntaxKind) *Node {
	return f.NewNode(kind, &KeywordTypeNode{})
}

// UnionOrIntersectionTypeBase

type UnionOrIntersectionTypeNodeBase struct {
	TypeNodeBase
	types []*Node
}

func (node *UnionOrIntersectionTypeNodeBase) ForEachChild(v Visitor) bool {
	return visitNodes(v, node.types)
}

// UnionTypeNode

type UnionTypeNode struct {
	UnionOrIntersectionTypeNodeBase
}

func (f *NodeFactory) NewUnionTypeNode(types []*Node) *Node {
	data := &UnionTypeNode{}
	data.types = f.ParenthesizerRules().ParenthesizeConstituentTypesOfUnionType(types)
	return f.NewNode(SyntaxKindUnionType, data)
}

func (f *NodeFactory) UpdateUnionTypeNode(node *Node, types []*Node) *Node {
	if n := node.AsUnionTypeNode(); !slices.Equal(types, n.types) {
		return f.UpdateNode(f.NewUnionTypeNode(types), node)
	}
	return node
}

func (node *UnionTypeNode) Accept(v NodeVisitor) *Node {
	return v.VisitUnionTypeNode(node)
}

// IntersectionTypeNode

type IntersectionTypeNode struct {
	UnionOrIntersectionTypeNodeBase
}

func (f *NodeFactory) NewIntersectionTypeNode(types []*Node) *Node {
	data := &IntersectionTypeNode{}
	data.types = f.ParenthesizerRules().ParenthesizeConstituentTypesOfIntersectionType(types)
	return f.NewNode(SyntaxKindIntersectionType, data)
}

func (f *NodeFactory) UpdateIntersectionTypeNode(node *Node, types []*Node) *Node {
	if n := node.AsIntersectionTypeNode(); !slices.Equal(types, n.types) {
		return f.UpdateNode(f.NewIntersectionTypeNode(types), node)
	}
	return node
}

func (node *IntersectionTypeNode) Accept(v NodeVisitor) *Node {
	return v.VisitIntersectionTypeNode(node)
}

// ConditionalTypeNode

type ConditionalTypeNode struct {
	TypeNodeBase
	LocalsContainerBase
	checkType   *Node
	extendsType *Node
	trueType    *Node
	falseType   *Node
}

func (f *NodeFactory) NewConditionalTypeNode(checkType *Node, extendsType *Node, trueType *Node, falseType *Node) *Node {
	data := &ConditionalTypeNode{}
	data.checkType = f.ParenthesizerRules().ParenthesizeCheckTypeOfConditionalType(checkType)
	data.extendsType = f.ParenthesizerRules().ParenthesizeExtendsTypeOfConditionalType(extendsType)
	data.trueType = trueType
	data.falseType = falseType
	return f.NewNode(SyntaxKindConditionalType, data)
}

func (f *NodeFactory) UpdateConditionalTypeNode(node *Node, checkType *Node, extendsType *Node, trueType *Node, falseType *Node) *Node {
	if n := node.AsConditionalTypeNode(); checkType != n.checkType || extendsType != n.extendsType || trueType != n.trueType || falseType != n.falseType {
		return f.UpdateNode(f.NewConditionalTypeNode(checkType, extendsType, trueType, falseType), node)
	}
	return node
}

func (node *ConditionalTypeNode) ForEachChild(v Visitor) bool {
	return visit(v, node.checkType) || visit(v, node.extendsType) || visit(v, node.trueType) || visit(v, node.falseType)
}

func (node *ConditionalTypeNode) Accept(v NodeVisitor) *Node {
	return v.VisitConditionalTypeNode(node)
}

func isConditionalTypeNode(node *Node) bool {
	return node.kind == SyntaxKindConditionalType
}

// TypeOperatorNode

type TypeOperatorNode struct {
	TypeNodeBase
	operator SyntaxKind // SyntaxKindKeyOfKeyword | SyntaxKindUniqueKeyword | SyntaxKindReadonlyKeyword
	typeNode *Node
}

func (f *NodeFactory) NewTypeOperatorNode(operator SyntaxKind, typeNode *Node) *Node {
	data := &TypeOperatorNode{}
	data.operator = operator
	switch operator {
	case SyntaxKindReadonlyKeyword:
		data.typeNode = f.ParenthesizerRules().ParenthesizeOperandOfReadonlyTypeOperator(typeNode)
	default:
		data.typeNode = f.ParenthesizerRules().ParenthesizeOperandOfTypeOperator(typeNode)
	}
	return f.NewNode(SyntaxKindTypeOperator, data)
}

func (f *NodeFactory) UpdateTypeOperatorNode(node *Node, typeNode *Node) *Node {
	if n := node.AsTypeOperatorNode(); typeNode != n.typeNode {
		return f.UpdateNode(f.NewTypeOperatorNode(n.operator, typeNode), node)
	}
	return node
}

func (node *TypeOperatorNode) ForEachChild(v Visitor) bool {
	return visit(v, node.typeNode)
}

func (node *TypeOperatorNode) Accept(v NodeVisitor) *Node {
	return v.VisitTypeOperatorNode(node)
}

func isTypeOperatorNode(node *Node) bool {
	return node.kind == SyntaxKindTypeOperator
}

// InferTypeNode

type InferTypeNode struct {
	TypeNodeBase
	typeParameter *Node
}

func (f *NodeFactory) NewInferTypeNode(typeParameter *Node) *Node {
	data := &InferTypeNode{}
	data.typeParameter = typeParameter
	return f.NewNode(SyntaxKindInferType, data)
}

func (f *NodeFactory) UpdateInferTypeNode(node *Node, typeParameter *Node) *Node {
	if n := node.AsInferTypeNode(); typeParameter != n.typeParameter {
		return f.UpdateNode(f.NewInferTypeNode(typeParameter), node)
	}
	return node
}

func (node *InferTypeNode) ForEachChild(v Visitor) bool {
	return visit(v, node.typeParameter)
}

func (node *InferTypeNode) Accept(v NodeVisitor) *Node {
	return v.VisitInferTypeNode(node)
}

// ArrayTypeNode

type ArrayTypeNode struct {
	TypeNodeBase
	elementType *Node
}

func (f *NodeFactory) NewArrayTypeNode(elementType *Node) *Node {
	data := &ArrayTypeNode{}
	data.elementType = f.ParenthesizerRules().ParenthesizeNonArrayTypeOfPostfixType(elementType)
	return f.NewNode(SyntaxKindArrayType, data)
}

func (f *NodeFactory) UpdateArrayTypeNode(node *Node, elementType *Node) *Node {
	if n := node.AsArrayTypeNode(); elementType != n.elementType {
		return f.UpdateNode(f.NewArrayTypeNode(elementType), node)
	}
	return node
}

func (node *ArrayTypeNode) ForEachChild(v Visitor) bool {
	return visit(v, node.elementType)
}

func (node *ArrayTypeNode) Accept(v NodeVisitor) *Node {
	return v.VisitArrayTypeNode(node)
}

// IndexedAccessTypeNode

type IndexedAccessTypeNode struct {
	TypeNodeBase
	objectType *Node
	indexType  *Node
}

func (f *NodeFactory) NewIndexedAccessTypeNode(objectType *Node, indexType *Node) *Node {
	data := &IndexedAccessTypeNode{}
	data.objectType = f.ParenthesizerRules().ParenthesizeNonArrayTypeOfPostfixType(objectType)
	data.indexType = indexType
	return f.NewNode(SyntaxKindIndexedAccessType, data)
}

func (f *NodeFactory) UpdateIndexedAccessTypeNode(node *Node, objectType *Node, indexType *Node) *Node {
	if n := node.AsIndexedAccessTypeNode(); objectType != n.objectType || indexType != n.indexType {
		return f.UpdateNode(f.NewIndexedAccessTypeNode(objectType, indexType), node)
	}
	return node
}

func (node *IndexedAccessTypeNode) ForEachChild(v Visitor) bool {
	return visit(v, node.objectType) || visit(v, node.indexType)
}

func (node *IndexedAccessTypeNode) Accept(v NodeVisitor) *Node {
	return v.VisitIndexedAccessTypeNode(node)
}

func isIndexedAccessTypeNode(node *Node) bool {
	return node.kind == SyntaxKindIndexedAccessType
}

// TypeArgumentList

type TypeArgumentList struct {
	NodeBase
	arguments []*Node
}

func (f *NodeFactory) NewTypeArgumentList(arguments []*Node) *Node {
	data := &TypeArgumentList{}
	data.arguments = arguments
	return f.NewNode(SyntaxKindTypeArgumentList, data)
}

func (f *NodeFactory) UpdateTypeArgumentList(node *Node, arguments []*Node) *Node {
	if n := node.AsTypeArgumentList(); !slices.Equal(arguments, n.arguments) {
		return f.UpdateNode(f.NewTypeArgumentList(arguments), node)
	}
	return node
}

func (node *TypeArgumentList) ForEachChild(v Visitor) bool {
	return visitNodes(v, node.arguments)
}

func (node *TypeArgumentList) Accept(v NodeVisitor) *Node {
	return v.VisitTypeArgumentList(node)
}

// TypeReferenceNode

type TypeReferenceNode struct {
	TypeNodeBase
	typeName      *Node
	typeArguments *Node // TypeArgumentList
}

func (f *NodeFactory) NewTypeReferenceNode(typeName *Node, typeArguments *Node) *Node {
	data := &TypeReferenceNode{}
	data.typeName = typeName
	data.typeArguments = f.ParenthesizerRules().ParenthesizeTypeArguments(typeArguments)
	return f.NewNode(SyntaxKindTypeReference, data)
}

func (f *NodeFactory) UpdateTypeReferenceNode(node *Node, typeName *Node, typeArguments *Node) *Node {
	if n := node.AsTypeReferenceNode(); typeName != n.typeName || typeArguments != n.typeArguments {
		return f.UpdateNode(f.NewTypeReferenceNode(typeName, typeArguments), node)
	}
	return node
}

func (node *TypeReferenceNode) ForEachChild(v Visitor) bool {
	return visit(v, node.typeName) || visit(v, node.typeArguments)
}

func (node *TypeReferenceNode) Accept(v NodeVisitor) *Node {
	return v.VisitTypeReferenceNode(node)
}

func isTypeReferenceNode(node *Node) bool {
	return node.kind == SyntaxKindTypeReference
}

// ExpressionWithTypeArguments

type ExpressionWithTypeArguments struct {
	ExpressionBase
	expression    *Node
	typeArguments *Node
}

func (f *NodeFactory) NewExpressionWithTypeArguments(expression *Node, typeArguments *Node) *Node {
	data := &ExpressionWithTypeArguments{}
	data.expression = f.ParenthesizerRules().ParenthesizeLeftSideOfAccess(expression, false /*optionalChain*/)
	data.typeArguments = f.ParenthesizerRules().ParenthesizeTypeArguments(typeArguments)
	return f.NewNode(SyntaxKindExpressionWithTypeArguments, data)
}

func (f *NodeFactory) UpdateExpressionWithTypeArguments(node *Node, expression *Node, typeArguments *Node) *Node {
	if n := node.AsExpressionWithTypeArguments(); expression != n.expression || typeArguments != n.typeArguments {
		return f.UpdateNode(f.NewExpressionWithTypeArguments(expression, typeArguments), node)
	}
	return node
}

func (node *ExpressionWithTypeArguments) ForEachChild(v Visitor) bool {
	return visit(v, node.expression) || visit(v, node.typeArguments)
}

func (node *ExpressionWithTypeArguments) Accept(v NodeVisitor) *Node {
	return v.VisitExpressionWithTypeArguments(node)
}

// LiteralTypeNode

type LiteralTypeNode struct {
	TypeNodeBase
	literal *Node // KeywordExpression | LiteralExpression | PrefixUnaryExpression
}

func (f *NodeFactory) NewLiteralTypeNode(literal *Node) *Node {
	data := &LiteralTypeNode{}
	data.literal = literal
	return f.NewNode(SyntaxKindLiteralType, data)
}

func (f *NodeFactory) UpdateLiteralTypeNode(node *Node, literal *Node) *Node {
	if n := node.AsLiteralTypeNode(); literal != n.literal {
		return f.UpdateNode(f.NewLiteralTypeNode(literal), node)
	}
	return node
}

func (node *LiteralTypeNode) ForEachChild(v Visitor) bool {
	return visit(v, node.literal)
}

func (node *LiteralTypeNode) Accept(v NodeVisitor) *Node {
	return v.VisitLiteralTypeNode(node)
}

func isLiteralTypeNode(node *Node) bool {
	return node.kind == SyntaxKindLiteralType
}

// ThisTypeNode

type ThisTypeNode struct {
	TypeNodeBase
}

func (f *NodeFactory) NewThisTypeNode() *Node {
	return f.NewNode(SyntaxKindThisType, &ThisTypeNode{})
}

func isThisTypeNode(node *Node) bool {
	return node.kind == SyntaxKindThisType
}

// TypePredicateNode

type TypePredicateNode struct {
	TypeNodeBase
	assertsModifier *Node // Optional
	parameterName   *Node // Identifier | ThisTypeNode
	typeNode        *Node // Optional
}

func (f *NodeFactory) NewTypePredicateNode(assertsModifier *Node, parameterName *Node, typeNode *Node) *Node {
	data := &TypePredicateNode{}
	data.assertsModifier = assertsModifier
	data.parameterName = parameterName
	data.typeNode = typeNode
	return f.NewNode(SyntaxKindTypePredicate, data)
}

func (f *NodeFactory) UpdateTypePredicateNode(node *Node, assertsModifier *Node, parameterName *Node, typeNode *Node) *Node {
	if n := node.AsTypePredicateNode(); assertsModifier != n.assertsModifier || parameterName != n.parameterName || typeNode != n.typeNode {
		return f.UpdateNode(f.NewTypePredicateNode(assertsModifier, parameterName, typeNode), node)
	}
	return node
}

func (node *TypePredicateNode) ForEachChild(v Visitor) bool {
	return visit(v, node.assertsModifier) || visit(v, node.parameterName) || visit(v, node.typeNode)
}

func (node *TypePredicateNode) Accept(v NodeVisitor) *Node {
	return v.VisitTypePredicateNode(node)
}

// ImportTypeNode

type ImportTypeNode struct {
	TypeNodeBase
	isTypeOf      bool
	argument      *Node
	attributes    *Node // Optional
	qualifier     *Node // Optional
	typeArguments *Node // Optional
}

func (f *NodeFactory) NewImportTypeNode(isTypeOf bool, argument *Node, attributes *Node, qualifier *Node, typeArguments *Node) *Node {
	data := &ImportTypeNode{}
	data.isTypeOf = isTypeOf
	data.argument = argument
	data.attributes = attributes
	data.qualifier = qualifier
	data.typeArguments = f.ParenthesizerRules().ParenthesizeTypeArguments(typeArguments)
	return f.NewNode(SyntaxKindImportType, data)
}

func (f *NodeFactory) UpdateImportTypeNode(node *Node, isTypeOf bool, argument *Node, attributes *Node, qualifier *Node, typeArguments *Node) *Node {
	if n := node.AsImportTypeNode(); isTypeOf != n.isTypeOf || argument != n.argument || attributes != n.attributes || qualifier != n.qualifier || typeArguments != n.typeArguments {
		return f.UpdateNode(f.NewImportTypeNode(isTypeOf, argument, attributes, qualifier, typeArguments), node)
	}
	return node
}

func (node *ImportTypeNode) ForEachChild(v Visitor) bool {
	return visit(v, node.argument) || visit(v, node.attributes) || visit(v, node.qualifier) || visit(v, node.typeArguments)
}

func (node *ImportTypeNode) Accept(v NodeVisitor) *Node {
	return v.VisitImportTypeNode(node)
}

func isImportTypeNode(node *Node) bool {
	return node.kind == SyntaxKindImportType
}

// ImportAttribute

type ImportAttribute struct {
	NodeBase
	name  *Node
	value *Node
}

func (f *NodeFactory) NewImportAttribute(name *Node, value *Node) *Node {
	data := &ImportAttribute{}
	data.name = name
	data.value = value
	return f.NewNode(SyntaxKindImportAttribute, data)
}

func (f *NodeFactory) UpdateImportAttribute(node *Node, name *Node, value *Node) *Node {
	if n := node.AsImportAttribute(); name != n.name || value != n.value {
		return f.UpdateNode(f.NewImportAttribute(name, value), node)
	}
	return node
}

func (node *ImportAttribute) ForEachChild(v Visitor) bool {
	return visit(v, node.name) || visit(v, node.value)
}

func (node *ImportAttribute) Accept(v NodeVisitor) *Node {
	return v.VisitImportAttribute(node)
}

// ImportAttributes

type ImportAttributes struct {
	NodeBase
	token      SyntaxKind
	attributes []*Node
	multiLine  bool
}

func (f *NodeFactory) NewImportAttributes(token SyntaxKind, attributes []*Node, multiLine bool) *Node {
	data := &ImportAttributes{}
	data.token = token
	data.attributes = attributes
	data.multiLine = multiLine
	return f.NewNode(SyntaxKindImportAttributes, data)
}

func (f *NodeFactory) UpdateImportAttributes(node *Node, attributes []*Node) *Node {
	if n := node.AsImportAttributes(); !slices.Equal(attributes, n.attributes) {
		return f.UpdateNode(f.NewImportAttributes(n.token, attributes, n.multiLine), node)
	}
	return node
}

func (node *ImportAttributes) ForEachChild(v Visitor) bool {
	return visitNodes(v, node.attributes)
}

func (node *ImportAttributes) Accept(v NodeVisitor) *Node {
	return v.VisitImportAttributes(node)
}

// TypeQueryNode

type TypeQueryNode struct {
	TypeNodeBase
	exprName      *Node
	typeArguments *Node
}

func (f *NodeFactory) NewTypeQueryNode(exprName *Node, typeArguments *Node) *Node {
	data := &TypeQueryNode{}
	data.exprName = exprName
	data.typeArguments = f.ParenthesizerRules().ParenthesizeTypeArguments(typeArguments)
	return f.NewNode(SyntaxKindTypeQuery, data)
}

func (f *NodeFactory) UpdateTypeQueryNode(node *Node, exprName *Node, typeArguments *Node) *Node {
	if n := node.AsTypeQueryNode(); exprName != n.exprName || typeArguments != n.typeArguments {
		return f.UpdateNode(f.NewTypeQueryNode(exprName, typeArguments), node)
	}
	return node
}

func (node *TypeQueryNode) ForEachChild(v Visitor) bool {
	return visit(v, node.exprName) || visit(v, node.typeArguments)
}

func (node *TypeQueryNode) Accept(v NodeVisitor) *Node {
	return v.VisitTypeQueryNode(node)
}

func isTypeQueryNode(node *Node) bool {
	return node.kind == SyntaxKindTypeQuery
}

// MappedTypeNode

type MappedTypeNode struct {
	TypeNodeBase
	DeclarationBase
	LocalsContainerBase
	readonlyToken *Node // Optional
	typeParameter *Node
	nameType      *Node   // Optional
	questionToken *Node   // Optional
	typeNode      *Node   // Optional (error if missing)
	members       []*Node // Used only to produce grammar errors
}

func (f *NodeFactory) NewMappedTypeNode(readonlyToken *Node, typeParameter *Node, nameType *Node, questionToken *Node, typeNode *Node, members []*Node) *Node {
	data := &MappedTypeNode{}
	data.readonlyToken = readonlyToken
	data.typeParameter = typeParameter
	data.nameType = nameType
	data.questionToken = questionToken
	data.typeNode = typeNode
	data.members = members
	return f.NewNode(SyntaxKindMappedType, data)
}

func (f *NodeFactory) UpdateMappedTypeNode(node *Node, readonlyToken *Node, typeParameter *Node, nameType *Node, questionToken *Node, typeNode *Node, members []*Node) *Node {
	if n := node.AsMappedTypeNode(); readonlyToken != n.readonlyToken || typeParameter != n.typeParameter || nameType != n.nameType || questionToken != n.questionToken || typeNode != n.typeNode || !slices.Equal(members, n.members) {
		return f.UpdateNode(f.NewMappedTypeNode(readonlyToken, typeParameter, nameType, questionToken, typeNode, members), node)
	}
	return node
}

func (node *MappedTypeNode) ForEachChild(v Visitor) bool {
	return visit(v, node.readonlyToken) || visit(v, node.typeParameter) || visit(v, node.nameType) ||
		visit(v, node.questionToken) || visit(v, node.typeNode) || visitNodes(v, node.members)
}

func (node *MappedTypeNode) Accept(v NodeVisitor) *Node {
	return v.VisitMappedTypeNode(node)
}

// TypeLiteralNode

type TypeLiteralNode struct {
	TypeNodeBase
	DeclarationBase
	members []*TypeElement
}

func (f *NodeFactory) NewTypeLiteralNode(members []*TypeElement) *Node {
	data := &TypeLiteralNode{}
	data.members = members
	return f.NewNode(SyntaxKindTypeLiteral, data)
}

func (f *NodeFactory) UpdateTypeLiteralNode(node *Node, members []*TypeElement) *Node {
	if n := node.AsTypeLiteralNode(); !slices.Equal(members, n.members) {
		return f.UpdateNode(f.NewTypeLiteralNode(members), node)
	}
	return node
}

func (node *TypeLiteralNode) ForEachChild(v Visitor) bool {
	return visitNodes(v, node.members)
}

func (node *TypeLiteralNode) Accept(v NodeVisitor) *Node {
	return v.VisitTypeLiteralNode(node)
}

// TupleTypeNode

type TupleTypeNode struct {
	TypeNodeBase
	elements []*TypeNode
}

func (f *NodeFactory) NewTupleTypeNode(elements []*TypeNode) *Node {
	data := &TupleTypeNode{}
	data.elements = f.ParenthesizerRules().ParenthesizeElementTypesOfTupleType(elements)
	return f.NewNode(SyntaxKindTupleType, data)
}

func (f *NodeFactory) UpdateTupleTypeNode(node *Node, elements []*TypeNode) *Node {
	if n := node.AsTupleTypeNode(); !slices.Equal(elements, n.elements) {
		return f.UpdateNode(f.NewTupleTypeNode(elements), node)
	}
	return node
}

func (node *TupleTypeNode) Kind() SyntaxKind {
	return SyntaxKindTupleType
}

func (node *TupleTypeNode) ForEachChild(v Visitor) bool {
	return visitNodes(v, node.elements)
}

func (node *TupleTypeNode) Accept(v NodeVisitor) *Node {
	return v.VisitTupleTypeNode(node)
}

// NamedTupleMember

type NamedTupleMember struct {
	TypeNodeBase
	DeclarationBase
	dotDotDotToken *Node
	name           *Node
	questionToken  *Node
	typeNode       *Node
}

func (f *NodeFactory) NewNamedTupleMember(dotDotDotToken *Node, name *Node, questionToken *Node, typeNode *Node) *Node {
	data := &NamedTupleMember{}
	data.dotDotDotToken = dotDotDotToken
	data.name = name
	data.questionToken = questionToken
	data.typeNode = typeNode
	return f.NewNode(SyntaxKindNamedTupleMember, data)
}

func (f *NodeFactory) UpdateNamedTupleMember(node *Node, dotDotDotToken *Node, name *Node, questionToken *Node, typeNode *Node) *Node {
	if n := node.AsNamedTupleMember(); dotDotDotToken != n.dotDotDotToken || name != n.name || questionToken != n.questionToken || typeNode != n.typeNode {
		return f.UpdateNode(f.NewNamedTupleMember(dotDotDotToken, name, questionToken, typeNode), node)
	}
	return node
}

func (node *NamedTupleMember) ForEachChild(v Visitor) bool {
	return visit(v, node.dotDotDotToken) || visit(v, node.name) || visit(v, node.questionToken) || visit(v, node.typeNode)
}

func (node *NamedTupleMember) Accept(v NodeVisitor) *Node {
	return v.VisitNamedTupleMember(node)
}

func isNamedTupleMember(node *Node) bool {
	return node.kind == SyntaxKindNamedTupleMember
}

// OptionalTypeNode

type OptionalTypeNode struct {
	TypeNodeBase
	typeNode *TypeNode
}

func (f *NodeFactory) NewOptionalTypeNode(typeNode *TypeNode) *Node {
	data := &OptionalTypeNode{}
	data.typeNode = f.ParenthesizerRules().ParenthesizeTypeOfOptionalType(typeNode)
	return f.NewNode(SyntaxKindOptionalType, data)
}

func (f *NodeFactory) UpdateOptionalTypeNode(node *Node, typeNode *TypeNode) *Node {
	if n := node.AsOptionalTypeNode(); typeNode != n.typeNode {
		return f.UpdateNode(f.NewOptionalTypeNode(typeNode), node)
	}
	return node
}

func (node *OptionalTypeNode) ForEachChild(v Visitor) bool {
	return visit(v, node.typeNode)
}

func (node *OptionalTypeNode) Accept(v NodeVisitor) *Node {
	return v.VisitOptionalTypeNode(node)
}

// RestTypeNode

type RestTypeNode struct {
	TypeNodeBase
	typeNode *TypeNode
}

func (f *NodeFactory) NewRestTypeNode(typeNode *TypeNode) *Node {
	data := &RestTypeNode{}
	data.typeNode = typeNode
	return f.NewNode(SyntaxKindRestType, data)
}

func (f *NodeFactory) UpdateRestTypeNode(node *Node, typeNode *TypeNode) *Node {
	if n := node.AsRestTypeNode(); typeNode != n.typeNode {
		return f.UpdateNode(f.NewRestTypeNode(typeNode), node)
	}
	return node
}

func (node *RestTypeNode) ForEachChild(v Visitor) bool {
	return visit(v, node.typeNode)
}

func (node *RestTypeNode) Accept(v NodeVisitor) *Node {
	return v.VisitRestTypeNode(node)
}

// ParenthesizedTypeNode

type ParenthesizedTypeNode struct {
	TypeNodeBase
	typeNode *TypeNode
}

func (f *NodeFactory) NewParenthesizedTypeNode(typeNode *TypeNode) *Node {
	data := &ParenthesizedTypeNode{}
	data.typeNode = typeNode
	return f.NewNode(SyntaxKindParenthesizedType, data)
}

func (f *NodeFactory) UpdateParenthesizedTypeNode(node *Node, typeNode *TypeNode) *Node {
	if n := node.AsParenthesizedTypeNode(); typeNode != n.typeNode {
		return f.UpdateNode(f.NewParenthesizedTypeNode(typeNode), node)
	}
	return node
}

func (node *ParenthesizedTypeNode) ForEachChild(v Visitor) bool {
	return visit(v, node.typeNode)
}

func (node *ParenthesizedTypeNode) Accept(v NodeVisitor) *Node {
	return v.VisitParenthesizedTypeNode(node)
}

func isParenthesizedTypeNode(node *Node) bool {
	return node.kind == SyntaxKindParenthesizedType
}

// FunctionOrConstructorTypeNodeBase

type FunctionOrConstructorTypeNodeBase struct {
	TypeNodeBase
	DeclarationBase
	ModifiersBase
	FunctionLikeBase
}

func (node *FunctionOrConstructorTypeNodeBase) ForEachChild(v Visitor) bool {
	return visit(v, node.modifiers) || visit(v, node.typeParameters) || visitNodes(v, node.parameters) || visit(v, node.returnType)
}

// FunctionTypeNode

type FunctionTypeNode struct {
	FunctionOrConstructorTypeNodeBase
}

func (f *NodeFactory) NewFunctionTypeNode(typeParameters *Node, parameters []*Node, returnType *Node) *Node {
	data := &FunctionTypeNode{}
	data.typeParameters = typeParameters
	data.parameters = parameters
	data.returnType = returnType
	return f.NewNode(SyntaxKindFunctionType, data)
}

func (f *NodeFactory) UpdateFunctionTypeNode(node *Node, typeParameters *Node, parameters []*Node, returnType *Node) *Node {
	if n := node.AsFunctionTypeNode(); typeParameters != n.typeParameters || !slices.Equal(parameters, n.parameters) || returnType != n.returnType {
		return f.UpdateNode(f.NewFunctionTypeNode(typeParameters, parameters, returnType), node)
	}
	return node
}

func (node *FunctionTypeNode) Accept(v NodeVisitor) *Node {
	return v.VisitFunctionTypeNode(node)
}

func isFunctionTypeNode(node *Node) bool {
	return node.kind == SyntaxKindFunctionType
}

// ConstructorTypeNode

type ConstructorTypeNode struct {
	FunctionOrConstructorTypeNodeBase
}

func (f *NodeFactory) NewConstructorTypeNode(modifiers *Node, typeParameters *Node, parameters []*Node, returnType *Node) *Node {
	data := &ConstructorTypeNode{}
	data.modifiers = modifiers
	data.typeParameters = typeParameters
	data.parameters = parameters
	data.returnType = returnType
	return f.NewNode(SyntaxKindConstructorType, data)
}

func (f *NodeFactory) UpdateConstructorTypeNode(node *Node, modifiers *Node, typeParameters *Node, parameters []*Node, returnType *Node) *Node {
	if n := node.AsConstructorTypeNode(); modifiers != n.modifiers || typeParameters != n.typeParameters || !slices.Equal(parameters, n.parameters) || returnType != n.returnType {
		return f.UpdateNode(f.NewConstructorTypeNode(modifiers, typeParameters, parameters, returnType), node)
	}
	return node
}

func (node *ConstructorTypeNode) Accept(v NodeVisitor) *Node {
	return v.VisitConstructorTypeNode(node)
}

func isConstructorTypeNode(node *Node) bool {
	return node.kind == SyntaxKindConstructorType
}

// TemplateLiteralLikeBase

type TemplateLiteralLikeBase struct {
	LiteralLikeBase
	rawText       string
	templateFlags TokenFlags
}

// TemplateHead

type TemplateHead struct {
	NodeBase
	TemplateLiteralLikeBase
}

func (f *NodeFactory) NewTemplateHead(text string, rawText string, templateFlags TokenFlags) *Node {
	data := &TemplateHead{}
	data.text = text
	data.rawText = rawText
	data.templateFlags = templateFlags
	return f.NewNode(SyntaxKindTemplateHead, data)
}

func (node *TemplateHead) Accept(v NodeVisitor) *Node {
	return v.VisitTemplateHead(node)
}

// TemplateMiddle

type TemplateMiddle struct {
	NodeBase
	TemplateLiteralLikeBase
}

func (f *NodeFactory) NewTemplateMiddle(text string, rawText string, templateFlags TokenFlags) *Node {
	data := &TemplateMiddle{}
	data.text = text
	data.rawText = rawText
	data.templateFlags = templateFlags
	return f.NewNode(SyntaxKindTemplateMiddle, data)
}

func (node *TemplateMiddle) Accept(v NodeVisitor) *Node {
	return v.VisitTemplateMiddle(node)
}

// TemplateTail

type TemplateTail struct {
	NodeBase
	TemplateLiteralLikeBase
}

func (f *NodeFactory) NewTemplateTail(text string, rawText string, templateFlags TokenFlags) *Node {
	data := &TemplateTail{}
	data.text = text
	data.rawText = rawText
	data.templateFlags = templateFlags
	return f.NewNode(SyntaxKindTemplateTail, data)
}

func (node *TemplateTail) Accept(v NodeVisitor) *Node {
	return v.VisitTemplateTail(node)
}

// TemplateLiteralTypeNode

type TemplateLiteralTypeNode struct {
	TypeNodeBase
	head          *Node
	templateSpans []*Node
}

func (f *NodeFactory) NewTemplateLiteralTypeNode(head *Node, templateSpans []*Node) *Node {
	data := &TemplateLiteralTypeNode{}
	data.head = head
	data.templateSpans = templateSpans
	return f.NewNode(SyntaxKindTemplateLiteralType, data)
}

func (f *NodeFactory) UpdateTemplateLiteralTypeNode(node *Node, head *Node, templateSpans []*Node) *Node {
	if n := node.AsTemplateLiteralTypeNode(); head != n.head || !slices.Equal(templateSpans, n.templateSpans) {
		return f.UpdateNode(f.NewTemplateLiteralTypeNode(head, templateSpans), node)
	}
	return node
}

func (node *TemplateLiteralTypeNode) ForEachChild(v Visitor) bool {
	return visit(v, node.head) || visitNodes(v, node.templateSpans)
}

func (node *TemplateLiteralTypeNode) Accept(v NodeVisitor) *Node {
	return v.VisitTemplateLiteralTypeNode(node)
}

// TemplateLiteralTypeSpan

type TemplateLiteralTypeSpan struct {
	NodeBase
	typeNode *Node
	literal  *Node
}

func (f *NodeFactory) NewTemplateLiteralTypeSpan(typeNode *Node, literal *Node) *Node {
	data := &TemplateLiteralTypeSpan{}
	data.typeNode = typeNode
	data.literal = literal
	return f.NewNode(SyntaxKindTemplateLiteralTypeSpan, data)
}

func (f *NodeFactory) UpdateTemplateLiteralTypeSpan(node *Node, typeNode *Node, literal *Node) *Node {
	if n := node.AsTemplateLiteralTypeSpan(); typeNode != n.typeNode || literal != n.literal {
		return f.UpdateNode(f.NewTemplateLiteralTypeSpan(typeNode, literal), node)
	}
	return node
}

func (node *TemplateLiteralTypeSpan) ForEachChild(v Visitor) bool {
	return visit(v, node.typeNode) || visit(v, node.literal)
}

func (node *TemplateLiteralTypeSpan) Accept(v NodeVisitor) *Node {
	return v.VisitTemplateLiteralTypeSpan(node)
}

/// A JSX expression of the form <TagName attrs>...</TagName>

type JsxElement struct {
	ExpressionBase
	openingElement *Node
	children       []*Node
	closingElement *Node
}

func (f *NodeFactory) NewJsxElement(openingElement *Node, children []*Node, closingElement *Node) *Node {
	data := &JsxElement{}
	data.openingElement = openingElement
	data.children = children
	data.closingElement = closingElement
	return f.NewNode(SyntaxKindJsxElement, data)
}

func (f *NodeFactory) UpdateJsxElement(node *Node, openingElement *Node, children []*Node, closingElement *Node) *Node {
	if n := node.AsJsxElement(); openingElement != n.openingElement || !slices.Equal(children, n.children) || closingElement != n.closingElement {
		return f.UpdateNode(f.NewJsxElement(openingElement, children, closingElement), node)
	}
	return node
}

func (node *JsxElement) ForEachChild(v Visitor) bool {
	return visit(v, node.openingElement) || visitNodes(v, node.children) || visit(v, node.closingElement)
}

func (node *JsxElement) Accept(v NodeVisitor) *Node {
	return v.VisitJsxElement(node)
}

// JsxAttributes

type JsxAttributes struct {
	ExpressionBase
	DeclarationBase
	properties []*JsxAttributeLike
}

func (f *NodeFactory) NewJsxAttributes(properties []*JsxAttributeLike) *Node {
	data := &JsxAttributes{}
	data.properties = properties
	return f.NewNode(SyntaxKindJsxAttributes, data)
}

func (f *NodeFactory) UpdateJsxAttributes(node *Node, properties []*JsxAttributeLike) *Node {
	if n := node.AsJsxAttributes(); !slices.Equal(properties, n.properties) {
		return f.UpdateNode(f.NewJsxAttributes(properties), node)
	}
	return node
}

func (node *JsxAttributes) ForEachChild(v Visitor) bool {
	return visitNodes(v, node.properties)
}

func (node *JsxAttributes) Accept(v NodeVisitor) *Node {
	return v.VisitJsxAttributes(node)
}

// JsxNamespacedName

type JsxNamespacedName struct {
	ExpressionBase
	name      *Node
	namespace *Node
}

func (f *NodeFactory) NewJsxNamespacedName(name *Node, namespace *Node) *Node {
	data := &JsxNamespacedName{}
	data.name = name
	data.namespace = namespace
	return f.NewNode(SyntaxKindJsxNamespacedName, data)
}

func (f *NodeFactory) UpdateJsxNamespacedName(node *Node, name *Node, namespace *Node) *Node {
	if n := node.AsJsxNamespacedName(); name != n.name || namespace != n.namespace {
		return f.UpdateNode(f.NewJsxNamespacedName(name, namespace), node)
	}
	return node
}

func (node *JsxNamespacedName) ForEachChild(v Visitor) bool {
	return visit(v, node.name) || visit(v, node.namespace)
}

func (node *JsxNamespacedName) Accept(v NodeVisitor) *Node {
	return v.VisitJsxNamespacedName(node)
}

func isJsxNamespacedName(node *Node) bool {
	return node.kind == SyntaxKindJsxNamespacedName
}

/// The opening element of a <Tag>...</Tag> JsxElement

type JsxOpeningElement struct {
	ExpressionBase
	tagName       *Node // Identifier | KeywordExpression | PropertyAccessExpression | JsxNamespacedName
	typeArguments *Node
	attributes    *Node
}

func (f *NodeFactory) NewJsxOpeningElement(tagName *Node, typeArguments *Node, attributes *Node) *Node {
	data := &JsxOpeningElement{}
	data.tagName = tagName
	data.typeArguments = typeArguments
	data.attributes = attributes
	return f.NewNode(SyntaxKindJsxOpeningElement, data)
}

func (f *NodeFactory) UpdateJsxOpeningElement(node *Node, tagName *Node, typeArguments *Node, attributes *Node) *Node {
	if n := node.AsJsxOpeningElement(); tagName != n.tagName || typeArguments != n.typeArguments || attributes != n.attributes {
		return f.UpdateNode(f.NewJsxOpeningElement(tagName, typeArguments, attributes), node)
	}
	return node
}

func (node *JsxOpeningElement) ForEachChild(v Visitor) bool {
	return visit(v, node.tagName) || visit(v, node.typeArguments) || visit(v, node.attributes)
}

func (node *JsxOpeningElement) Accept(v NodeVisitor) *Node {
	return v.VisitJsxOpeningElement(node)
}

func isJsxOpeningElement(node *Node) bool {
	return node.kind == SyntaxKindJsxOpeningElement
}

/// A JSX expression of the form <TagName attrs />

type JsxSelfClosingElement struct {
	ExpressionBase
	tagName       *Node // Identifier | KeywordExpression | PropertyAccessExpression | JsxNamespacedName
	typeArguments *Node
	attributes    *Node
}

func (f *NodeFactory) NewJsxSelfClosingElement(tagName *Node, typeArguments *Node, attributes *Node) *Node {
	data := &JsxSelfClosingElement{}
	data.tagName = tagName
	data.typeArguments = typeArguments
	data.attributes = attributes
	return f.NewNode(SyntaxKindJsxSelfClosingElement, data)
}

func (f *NodeFactory) UpdateJsxSelfClosingElement(node *Node, tagName *Node, typeArguments *Node, attributes *Node) *Node {
	if n := node.AsJsxSelfClosingElement(); tagName != n.tagName || typeArguments != n.typeArguments || attributes != n.attributes {
		return f.UpdateNode(f.NewJsxSelfClosingElement(tagName, typeArguments, attributes), node)
	}
	return node
}

func (node *JsxSelfClosingElement) ForEachChild(v Visitor) bool {
	return visit(v, node.tagName) || visit(v, node.typeArguments) || visit(v, node.attributes)
}

func (node *JsxSelfClosingElement) Accept(v NodeVisitor) *Node {
	return v.VisitJsxSelfClosingElement(node)
}

/// A JSX expression of the form <>...</>

type JsxFragment struct {
	ExpressionBase
	openingFragment *Node
	children        []*Node
	closingFragment *Node
}

func (f *NodeFactory) NewJsxFragment(openingFragment *Node, children []*Node, closingFragment *Node) *Node {
	data := &JsxFragment{}
	data.openingFragment = openingFragment
	data.children = children
	data.closingFragment = closingFragment
	return f.NewNode(SyntaxKindJsxFragment, data)
}

func (f *NodeFactory) UpdateJsxFragment(node *Node, openingFragment *Node, children []*Node, closingFragment *Node) *Node {
	if n := node.AsJsxFragment(); openingFragment != n.openingFragment || !slices.Equal(children, n.children) || closingFragment != n.closingFragment {
		return f.UpdateNode(f.NewJsxFragment(openingFragment, children, closingFragment), node)
	}
	return node
}

func (node *JsxFragment) ForEachChild(v Visitor) bool {
	return visit(v, node.openingFragment) || visitNodes(v, node.children) || visit(v, node.closingFragment)
}

func (node *JsxFragment) Accept(v NodeVisitor) *Node {
	return v.VisitJsxFragment(node)
}

/// The opening element of a <>...</> JsxFragment

type JsxOpeningFragment struct {
	ExpressionBase
}

func (f *NodeFactory) NewJsxOpeningFragment() *Node {
	return f.NewNode(SyntaxKindJsxOpeningFragment, &JsxOpeningFragment{})
}

func (f *NodeFactory) UpdateJsxAttribute(node *Node, name *Node, initializer *Node) *Node {
	if n := node.AsJsxAttribute(); name != n.name || initializer != n.initializer {
		return f.UpdateNode(f.NewJsxAttribute(name, initializer), node)
	}
	return node
}

func (node *JsxOpeningFragment) Accept(v NodeVisitor) *Node {
	return v.VisitJsxOpeningFragment(node)
}

func isJsxOpeningFragment(node *Node) bool {
	return node.kind == SyntaxKindJsxOpeningFragment
}

/// The closing element of a <>...</> JsxFragment

type JsxClosingFragment struct {
	ExpressionBase
}

func (f *NodeFactory) NewJsxClosingFragment() *Node {
	return f.NewNode(SyntaxKindJsxClosingFragment, &JsxClosingFragment{})
}

func (node *JsxClosingFragment) Accept(v NodeVisitor) *Node {
	return v.VisitJsxClosingFragment(node)
}

// JsxAttribute

type JsxAttribute struct {
	NodeBase
	DeclarationBase
	name *Node
	/// JSX attribute initializers are optional; <X y /> is sugar for <X y={true} />
	initializer *Node
}

func (f *NodeFactory) NewJsxAttribute(name *Node, initializer *Node) *Node {
	data := &JsxAttribute{}
	data.name = name
	data.initializer = initializer
	return f.NewNode(SyntaxKindJsxAttribute, data)
}

func (node *JsxAttribute) ForEachChild(v Visitor) bool {
	return visit(v, node.name) || visit(v, node.initializer)
}

func (node *JsxAttribute) Accept(v NodeVisitor) *Node {
	return v.VisitJsxAttribute(node)
}

func isJsxAttribute(node *Node) bool {
	return node.kind == SyntaxKindJsxAttribute
}

// JsxSpreadAttribute

type JsxSpreadAttribute struct {
	NodeBase
	expression *Node
}

func (f *NodeFactory) NewJsxSpreadAttribute(expression *Node) *Node {
	data := &JsxSpreadAttribute{}
	data.expression = expression
	return f.NewNode(SyntaxKindJsxSpreadAttribute, data)
}

func (f *NodeFactory) UpdateJsxSpreadAttribute(node *Node, expression *Node) *Node {
	if n := node.AsJsxSpreadAttribute(); expression != n.expression {
		return f.UpdateNode(f.NewJsxSpreadAttribute(expression), node)
	}
	return node
}

func (node *JsxSpreadAttribute) ForEachChild(v Visitor) bool {
	return visit(v, node.expression)
}

func (node *JsxSpreadAttribute) Accept(v NodeVisitor) *Node {
	return v.VisitJsxSpreadAttribute(node)
}

// JsxClosingElement

type JsxClosingElement struct {
	NodeBase
	tagName *Node // Identifier | KeywordExpression | PropertyAccessExpression | JsxNamespacedName
}

func (f *NodeFactory) NewJsxClosingElement(tagName *Node) *Node {
	data := &JsxClosingElement{}
	data.tagName = tagName
	return f.NewNode(SyntaxKindJsxClosingElement, data)
}

func (f *NodeFactory) UpdateJsxClosingElement(node *Node, tagName *Node) *Node {
	if n := node.AsJsxClosingElement(); tagName != n.tagName {
		return f.UpdateNode(f.NewJsxClosingElement(tagName), node)
	}
	return node
}

func (node *JsxClosingElement) ForEachChild(v Visitor) bool {
	return visit(v, node.tagName)
}

func (node *JsxClosingElement) Accept(v NodeVisitor) *Node {
	return v.VisitJsxClosingElement(node)
}

// JsxExpression

type JsxExpression struct {
	ExpressionBase
	dotDotDotToken *Node
	expression     *Node
}

func (f *NodeFactory) NewJsxExpression(dotDotDotToken *Node, expression *Node) *Node {
	data := &JsxExpression{}
	data.dotDotDotToken = dotDotDotToken
	data.expression = expression
	return f.NewNode(SyntaxKindJsxExpression, data)
}

func (f *NodeFactory) UpdateJsxExpression(node *Node, dotDotDotToken *Node, expression *Node) *Node {
	if n := node.AsJsxExpression(); dotDotDotToken != n.dotDotDotToken || expression != n.expression {
		return f.UpdateNode(f.NewJsxExpression(dotDotDotToken, expression), node)
	}
	return node
}

func (node *JsxExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.dotDotDotToken) || visit(v, node.expression)
}

func (node *JsxExpression) Accept(v NodeVisitor) *Node {
	return v.VisitJsxExpression(node)
}

// JsxText

type JsxText struct {
	ExpressionBase
	LiteralLikeBase
	containsOnlyTriviaWhiteSpaces bool
}

func (f *NodeFactory) NewJsxText(text string, containsOnlyTriviaWhiteSpace bool) *Node {
	data := &JsxText{}
	data.text = text
	data.containsOnlyTriviaWhiteSpaces = containsOnlyTriviaWhiteSpace
	return f.NewNode(SyntaxKindJsxText, data)
}

func (node *JsxText) Accept(v NodeVisitor) *Node {
	return v.VisitJsxText(node)
}

// SyntaxList

type SyntaxList struct {
	NodeBase
	children []*Node
}

func (f *NodeFactory) NewSyntaxList(children []*Node) *Node {
	data := &SyntaxList{}
	data.children = children
	return f.NewNode(SyntaxKindSyntaxList, data)
}

// JSDocNonNullableType

type JSDocNonNullableType struct {
	TypeNodeBase
	typeNode *Node
	postfix  bool
}

func (f *NodeFactory) NewJSDocNonNullableType(typeNode *Node, postfix bool) *Node {
	data := &JSDocNonNullableType{}
	data.typeNode = typeNode
	data.postfix = postfix
	return f.NewNode(SyntaxKindJSDocNonNullableType, data)
}

func (f *NodeFactory) UpdateJSDocNonNullableType(node *Node, typeNode *Node) *Node {
	if n := node.AsJSDocNonNullableType(); typeNode != n.typeNode {
		return f.UpdateNode(f.NewJSDocNonNullableType(typeNode, n.postfix), node)
	}
	return node
}

func (node *JSDocNonNullableType) ForEachChild(v Visitor) bool {
	return visit(v, node.typeNode)
}

func (node *JSDocNonNullableType) Accept(v NodeVisitor) *Node {
	return v.VisitJSDocNonNullableType(node)
}

// JSDocNullableType

type JSDocNullableType struct {
	TypeNodeBase
	typeNode *Node
	postfix  bool
}

func (f *NodeFactory) NewJSDocNullableType(typeNode *Node, postfix bool) *Node {
	data := &JSDocNullableType{}
	data.typeNode = typeNode
	data.postfix = postfix
	return f.NewNode(SyntaxKindJSDocNullableType, data)
}

func (f *NodeFactory) UpdateJSDocNullableType(node *Node, typeNode *Node) *Node {
	if n := node.AsJSDocNullableType(); typeNode != n.typeNode {
		return f.UpdateNode(f.NewJSDocNullableType(typeNode, n.postfix), node)
	}
	return node
}

func (node *JSDocNullableType) ForEachChild(v Visitor) bool {
	return visit(v, node.typeNode)
}

func (node *JSDocNullableType) Accept(v NodeVisitor) *Node {
	return v.VisitJSDocNullableType(node)
}

// PatternAmbientModule

type PatternAmbientModule struct {
	pattern Pattern
	symbol  *Symbol
}

// SourceFile

type SourceFile struct {
	NodeBase
	DeclarationBase
	LocalsContainerBase
	text                        string
	fileName                    string
	path                        string
	statements                  []*Statement
	diagnostics                 []*Diagnostic
	bindDiagnostics             []*Diagnostic
	bindSuggestionDiagnostics   []*Diagnostic
	lineMap                     []TextPos
	languageVersion             ScriptTarget
	languageVariant             LanguageVariant
	scriptKind                  ScriptKind
	externalModuleIndicator     *Node
	endFlowNode                 *FlowNode
	jsGlobalAugmentations       SymbolTable
	isDeclarationFile           bool
	isBound                     bool
	moduleReferencesProcessed   bool
	usesUriStyleNodeCoreModules Tristate
	symbolCount                 int
	classifiableNames           set[string]
	imports                     []*LiteralLikeNode
	moduleAugmentations         []*ModuleName
	patternAmbientModules       []PatternAmbientModule
	ambientModuleNames          []string
}

func (f *NodeFactory) NewSourceFile(text string, fileName string, statements []*Node) *Node {
	data := &SourceFile{}
	data.text = text
	data.fileName = fileName
	data.statements = statements
	data.languageVersion = ScriptTargetLatest
	return f.NewNode(SyntaxKindSourceFile, data)
}

func (f *NodeFactory) UpdateSourceFile(node *Node, statements []*Node) *Node {
	if n := node.AsSourceFile(); !slices.Equal(statements, n.statements) {
		updated := f.NewSourceFile(n.text, n.fileName, statements).AsSourceFile()
		updated.path = n.path
		updated.languageVersion = n.languageVersion
		updated.languageVariant = n.languageVariant
		updated.scriptKind = n.scriptKind
		updated.isDeclarationFile = n.isDeclarationFile
		// TODO: Include other fields or use .original to get to original source file
		return f.UpdateNode(updated.AsNode(), node)
	}
	return node
}

func (node *SourceFile) FileName() string {
	return node.fileName
}

func (node *SourceFile) Path() string {
	return node.path
}

func (node *SourceFile) Diagnostics() []*Diagnostic {
	return node.diagnostics
}

func (node *SourceFile) BindDiagnostics() []*Diagnostic {
	return node.bindDiagnostics
}

func (node *SourceFile) ForEachChild(v Visitor) bool {
	return visitNodes(v, node.statements)
}

func (node *SourceFile) Accept(v NodeVisitor) *Node {
	return v.VisitSourceFile(node)
}

func isSourceFile(node *Node) bool {
	return node.kind == SyntaxKindSourceFile
}
