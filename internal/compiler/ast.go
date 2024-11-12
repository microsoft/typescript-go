package compiler

import (
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/compiler/textpos"
	"github.com/microsoft/typescript-go/internal/core"
)

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

// NodeFactory

type NodeFactory struct {
	identifierPool core.Pool[Identifier]
}

func (f *NodeFactory) NewNode(kind ast.Kind, data NodeData) *Node {
	n := data.AsNode()
	n.Kind = kind
	n.Data = data
	return n
}

// AST Node
// Interface values stored in AST nodes are never typed nil values. Construction code must ensure that
// interface valued properties either store a true nil or a reference to a non-nil struct.

type Node struct {
	Kind   ast.Kind
	Flags  ast.NodeFlags
	Loc    TextRange
	Id     NodeId
	Parent *Node
	Data   NodeData
}

// Node accessors. Some accessors are implemented as methods on NodeData, others are implemented though
// type switches. Either approach is fine. Interface methods are likely more performant, but have higher
// code size costs because we have hundreds of implementations of the NodeData interface.

func (n *Node) Pos() int                                  { return n.Loc.Pos() }
func (n *Node) End() int                                  { return n.Loc.End() }
func (n *Node) ForEachChild(v Visitor) bool               { return n.Data.ForEachChild(v) }
func (n *Node) Name() *DeclarationName                    { return n.Data.Name() }
func (n *Node) Modifiers() *ModifierListNode              { return n.Data.Modifiers() }
func (n *Node) TypeParameters() *TypeParameterListNode    { return n.Data.TypeParameters() }
func (n *Node) FlowNodeData() *FlowNodeBase               { return n.Data.FlowNodeData() }
func (n *Node) DeclarationData() *DeclarationBase         { return n.Data.DeclarationData() }
func (n *Node) Symbol() *Symbol                           { return n.Data.DeclarationData().Symbol }
func (n *Node) ExportableData() *ExportableBase           { return n.Data.ExportableData() }
func (n *Node) LocalSymbol() *Symbol                      { return n.Data.ExportableData().LocalSymbol }
func (n *Node) LocalsContainerData() *LocalsContainerBase { return n.Data.LocalsContainerData() }
func (n *Node) Locals() SymbolTable                       { return n.Data.LocalsContainerData().Locals_ }
func (n *Node) FunctionLikeData() *FunctionLikeBase       { return n.Data.FunctionLikeData() }
func (n *Node) Parameters() []*ParameterDeclarationNode   { return n.Data.FunctionLikeData().Parameters }
func (n *Node) ReturnType() *TypeNode                     { return n.Data.FunctionLikeData().ReturnType }
func (n *Node) ClassLikeData() *ClassLikeBase             { return n.Data.ClassLikeData() }
func (n *Node) BodyData() *BodyBase                       { return n.Data.BodyData() }

func (n *Node) Text() string {
	switch n.Kind {
	case ast.KindIdentifier:
		return n.AsIdentifier().Text
	case ast.KindPrivateIdentifier:
		return n.AsPrivateIdentifier().Text
	case ast.KindStringLiteral:
		return n.AsStringLiteral().Text
	case ast.KindNumericLiteral:
		return n.AsNumericLiteral().Text
	case ast.KindBigIntLiteral:
		return n.AsBigIntLiteral().Text
	case ast.KindNoSubstitutionTemplateLiteral:
		return n.AsNoSubstitutionTemplateLiteral().Text
	case ast.KindJsxNamespacedName:
		return n.AsJsxNamespacedName().Namespace.Text() + ":" + n.AsJsxNamespacedName().Name_.Text()
	}
	panic("Unhandled case in Node.Text")
}

func (node *Node) Expression() *Node {
	switch node.Kind {
	case ast.KindPropertyAccessExpression:
		return node.AsPropertyAccessExpression().Expression
	case ast.KindElementAccessExpression:
		return node.AsElementAccessExpression().Expression
	case ast.KindParenthesizedExpression:
		return node.AsParenthesizedExpression().Expression
	case ast.KindCallExpression:
		return node.AsCallExpression().Expression
	case ast.KindNewExpression:
		return node.AsNewExpression().Expression
	case ast.KindExpressionWithTypeArguments:
		return node.AsExpressionWithTypeArguments().Expression
	case ast.KindNonNullExpression:
		return node.AsNonNullExpression().Expression
	case ast.KindTypeAssertionExpression:
		return node.AsTypeAssertion().Expression
	case ast.KindAsExpression:
		return node.AsAsExpression().Expression
	case ast.KindSatisfiesExpression:
		return node.AsSatisfiesExpression().Expression
	case ast.KindSpreadAssignment:
		return node.AsSpreadAssignment().Expression
	}
	panic("Unhandled case in Node.Expression")
}

func (node *Node) Arguments() []*Node {
	switch node.Kind {
	case ast.KindCallExpression:
		return node.AsCallExpression().Arguments
	case ast.KindNewExpression:
		return node.AsNewExpression().Arguments
	}
	panic("Unhandled case in Node.Arguments")
}

// Node casts

func (n *Node) AsIdentifier() *Identifier {
	return n.Data.(*Identifier)
}
func (n *Node) AsPrivateIdentifier() *PrivateIdentifier {
	return n.Data.(*PrivateIdentifier)
}
func (n *Node) AsQualifiedName() *QualifiedName {
	return n.Data.(*QualifiedName)
}
func (n *Node) AsModifierList() *ModifierList {
	return n.Data.(*ModifierList)
}
func (n *Node) AsSourceFile() *SourceFile {
	return n.Data.(*SourceFile)
}
func (n *Node) AsPrefixUnaryExpression() *PrefixUnaryExpression {
	return n.Data.(*PrefixUnaryExpression)
}
func (n *Node) AsPostfixUnaryExpression() *PostfixUnaryExpression {
	return n.Data.(*PostfixUnaryExpression)
}
func (n *Node) AsParenthesizedExpression() *ParenthesizedExpression {
	return n.Data.(*ParenthesizedExpression)
}
func (n *Node) AsTypeAssertion() *TypeAssertion {
	return n.Data.(*TypeAssertion)
}
func (n *Node) AsAsExpression() *AsExpression {
	return n.Data.(*AsExpression)
}
func (n *Node) AsSatisfiesExpression() *SatisfiesExpression {
	return n.Data.(*SatisfiesExpression)
}
func (n *Node) AsExpressionWithTypeArguments() *ExpressionWithTypeArguments {
	return n.Data.(*ExpressionWithTypeArguments)
}
func (n *Node) AsNonNullExpression() *NonNullExpression {
	return n.Data.(*NonNullExpression)
}
func (n *Node) AsBindingElement() *BindingElement {
	return n.Data.(*BindingElement)
}
func (n *Node) AsImportSpecifier() *ImportSpecifier {
	return n.Data.(*ImportSpecifier)
}
func (n *Node) AsArrowFunction() *ArrowFunction {
	return n.Data.(*ArrowFunction)
}
func (n *Node) AsCallExpression() *CallExpression {
	return n.Data.(*CallExpression)
}
func (n *Node) AsPropertyAccessExpression() *PropertyAccessExpression {
	return n.Data.(*PropertyAccessExpression)
}
func (n *Node) AsElementAccessExpression() *ElementAccessExpression {
	return n.Data.(*ElementAccessExpression)
}
func (n *Node) AsComputedPropertyName() *ComputedPropertyName {
	return n.Data.(*ComputedPropertyName)
}
func (n *Node) AsBinaryExpression() *BinaryExpression {
	return n.Data.(*BinaryExpression)
}
func (n *Node) AsModuleDeclaration() *ModuleDeclaration {
	return n.Data.(*ModuleDeclaration)
}
func (n *Node) AsStringLiteral() *StringLiteral {
	return n.Data.(*StringLiteral)
}
func (n *Node) AsNumericLiteral() *NumericLiteral {
	return n.Data.(*NumericLiteral)
}
func (n *Node) AsBigIntLiteral() *BigIntLiteral {
	return n.Data.(*BigIntLiteral)
}
func (n *Node) AsNoSubstitutionTemplateLiteral() *NoSubstitutionTemplateLiteral {
	return n.Data.(*NoSubstitutionTemplateLiteral)
}
func (n *Node) AsVariableDeclaration() *VariableDeclaration {
	return n.Data.(*VariableDeclaration)
}
func (n *Node) AsExportAssignment() *ExportAssignment {
	return n.Data.(*ExportAssignment)
}
func (n *Node) AsObjectLiteralExpression() *ObjectLiteralExpression {
	return n.Data.(*ObjectLiteralExpression)
}
func (n *Node) AsIfStatement() *IfStatement {
	return n.Data.(*IfStatement)
}
func (n *Node) AsWhileStatement() *WhileStatement {
	return n.Data.(*WhileStatement)
}
func (n *Node) AsDoStatement() *DoStatement {
	return n.Data.(*DoStatement)
}
func (n *Node) AsForStatement() *ForStatement {
	return n.Data.(*ForStatement)
}
func (n *Node) AsConditionalExpression() *ConditionalExpression {
	return n.Data.(*ConditionalExpression)
}
func (n *Node) AsForInOrOfStatement() *ForInOrOfStatement {
	return n.Data.(*ForInOrOfStatement)
}
func (n *Node) AsShorthandPropertyAssignment() *ShorthandPropertyAssignment {
	return n.Data.(*ShorthandPropertyAssignment)
}
func (n *Node) AsPropertyAssignment() *PropertyAssignment {
	return n.Data.(*PropertyAssignment)
}
func (n *Node) AsExpressionStatement() *ExpressionStatement {
	return n.Data.(*ExpressionStatement)
}
func (n *Node) AsBlock() *Block {
	return n.Data.(*Block)
}
func (n *Node) AsModuleBlock() *ModuleBlock {
	return n.Data.(*ModuleBlock)
}
func (n *Node) AsVariableStatement() *VariableStatement {
	return n.Data.(*VariableStatement)
}
func (n *Node) AsVariableDeclarationList() *VariableDeclarationList {
	return n.Data.(*VariableDeclarationList)
}
func (n *Node) AsMetaProperty() *MetaProperty {
	return n.Data.(*MetaProperty)
}
func (n *Node) AsTypeReference() *TypeReferenceNode {
	return n.Data.(*TypeReferenceNode)
}
func (n *Node) AsConstructorDeclaration() *ConstructorDeclaration {
	return n.Data.(*ConstructorDeclaration)
}
func (n *Node) AsConditionalTypeNode() *ConditionalTypeNode {
	return n.Data.(*ConditionalTypeNode)
}
func (n *Node) AsClassExpression() *ClassExpression {
	return n.Data.(*ClassExpression)
}
func (n *Node) AsHeritageClause() *HeritageClause {
	return n.Data.(*HeritageClause)
}
func (n *Node) AsFunctionExpression() *FunctionExpression {
	return n.Data.(*FunctionExpression)
}
func (n *Node) AsParameterDeclaration() *ParameterDeclaration {
	return n.Data.(*ParameterDeclaration)
}
func (n *Node) AsInferTypeNode() *InferTypeNode {
	return n.Data.(*InferTypeNode)
}
func (n *Node) AsTypeParameter() *TypeParameterDeclaration {
	return n.Data.(*TypeParameterDeclaration)
}
func (n *Node) AsExportSpecifier() *ExportSpecifier {
	return n.Data.(*ExportSpecifier)
}
func (n *Node) AsExportDeclaration() *ExportDeclaration {
	return n.Data.(*ExportDeclaration)
}
func (n *Node) AsPropertyDeclaration() *PropertyDeclaration {
	return n.Data.(*PropertyDeclaration)
}
func (n *Node) AsImportClause() *ImportClause {
	return n.Data.(*ImportClause)
}
func (n *Node) AsImportEqualsDeclaration() *ImportEqualsDeclaration {
	return n.Data.(*ImportEqualsDeclaration)
}
func (n *Node) AsNamespaceImport() *NamespaceImport {
	return n.Data.(*NamespaceImport)
}
func (n *Node) AsPropertySignatureDeclaration() *PropertySignatureDeclaration {
	return n.Data.(*PropertySignatureDeclaration)
}
func (n *Node) AsEnumMember() *EnumMember {
	return n.Data.(*EnumMember)
}
func (n *Node) AsReturnStatement() *ReturnStatement {
	return n.Data.(*ReturnStatement)
}
func (n *Node) AsWithStatement() *WithStatement {
	return n.Data.(*WithStatement)
}
func (n *Node) AsSwitchStatement() *SwitchStatement {
	return n.Data.(*SwitchStatement)
}
func (n *Node) AsCaseOrDefaultClause() *CaseOrDefaultClause {
	return n.Data.(*CaseOrDefaultClause)
}
func (n *Node) AsThrowStatement() *ThrowStatement {
	return n.Data.(*ThrowStatement)
}
func (n *Node) AsTemplateSpan() *TemplateSpan {
	return n.Data.(*TemplateSpan)
}
func (n *Node) AsImportTypeNode() *ImportTypeNode {
	return n.Data.(*ImportTypeNode)
}
func (n *Node) AsNewExpression() *NewExpression {
	return n.Data.(*NewExpression)
}
func (n *Node) AsTaggedTemplateExpression() *TaggedTemplateExpression {
	return n.Data.(*TaggedTemplateExpression)
}
func (n *Node) AsTypeArgumentList() *TypeArgumentList {
	return n.Data.(*TypeArgumentList)
}
func (n *Node) AsJsxOpeningElement() *JsxOpeningElement {
	return n.Data.(*JsxOpeningElement)
}
func (n *Node) AsJsxSelfClosingElement() *JsxSelfClosingElement {
	return n.Data.(*JsxSelfClosingElement)
}
func (n *Node) AsJsxClosingElement() *JsxClosingElement {
	return n.Data.(*JsxClosingElement)
}
func (n *Node) AsImportDeclaration() *ImportDeclaration {
	return n.Data.(*ImportDeclaration)
}
func (n *Node) AsExternalModuleReference() *ExternalModuleReference {
	return n.Data.(*ExternalModuleReference)
}
func (n *Node) AsLiteralTypeNode() *LiteralTypeNode {
	return n.Data.(*LiteralTypeNode)
}
func (n *Node) AsJsxNamespacedName() *JsxNamespacedName {
	return n.Data.(*JsxNamespacedName)
}
func (n *Node) AsTypeParameterList() *TypeParameterList {
	return n.Data.(*TypeParameterList)
}
func (n *Node) AsClassDeclaration() *ClassDeclaration {
	return n.Data.(*ClassDeclaration)
}
func (n *Node) AsInterfaceDeclaration() *InterfaceDeclaration {
	return n.Data.(*InterfaceDeclaration)
}
func (n *Node) AsTypeAliasDeclaration() *TypeAliasDeclaration {
	return n.Data.(*TypeAliasDeclaration)
}
func (n *Node) AsJsxAttribute() *JsxAttribute {
	return n.Data.(*JsxAttribute)
}
func (n *Node) AsParenthesizedTypeNode() *ParenthesizedTypeNode {
	return n.Data.(*ParenthesizedTypeNode)
}
func (n *Node) AsTypePredicateNode() *TypePredicateNode {
	return n.Data.(*TypePredicateNode)
}
func (n *Node) AsTypeOperatorNode() *TypeOperatorNode {
	return n.Data.(*TypeOperatorNode)
}
func (n *Node) AsMappedTypeNode() *MappedTypeNode {
	return n.Data.(*MappedTypeNode)
}
func (n *Node) AsArrayLiteralExpression() *ArrayLiteralExpression {
	return n.Data.(*ArrayLiteralExpression)
}
func (n *Node) AsMethodDeclaration() *MethodDeclaration {
	return n.Data.(*MethodDeclaration)
}
func (n *Node) AsMethodSignatureDeclaration() *MethodSignatureDeclaration {
	return n.Data.(*MethodSignatureDeclaration)
}
func (n *Node) AsTemplateLiteralTypeSpan() *TemplateLiteralTypeSpan {
	return n.Data.(*TemplateLiteralTypeSpan)
}
func (n *Node) AsJsxElement() *JsxElement {
	return n.Data.(*JsxElement)
}
func (n *Node) AsKeywordExpression() *KeywordExpression {
	return n.Data.(*KeywordExpression)
}
func (n *Node) AsCatchClause() *CatchClause {
	return n.Data.(*CatchClause)
}
func (n *Node) AsDeleteExpression() *DeleteExpression {
	return n.Data.(*DeleteExpression)
}
func (n *Node) AsLabeledStatement() *LabeledStatement {
	return n.Data.(*LabeledStatement)
}
func (n *Node) AsNamespaceExportDeclaration() *NamespaceExportDeclaration {
	return n.Data.(*NamespaceExportDeclaration)
}
func (n *Node) AsNamedExports() *NamedExports {
	return n.Data.(*NamedExports)
}
func (n *Node) AsBreakStatement() *BreakStatement {
	return n.Data.(*BreakStatement)
}
func (n *Node) AsContinueStatement() *ContinueStatement {
	return n.Data.(*ContinueStatement)
}
func (n *Node) AsCaseBlock() *CaseBlock {
	return n.Data.(*CaseBlock)
}
func (n *Node) AsTryStatement() *TryStatement {
	return n.Data.(*TryStatement)
}
func (n *Node) AsBindingPattern() *BindingPattern {
	return n.Data.(*BindingPattern)
}
func (n *Node) AsFunctionDeclaration() *FunctionDeclaration {
	return n.Data.(*FunctionDeclaration)
}
func (n *Node) AsTypeOfExpression() *TypeOfExpression {
	return n.Data.(*TypeOfExpression)
}
func (n *Node) AsSpreadElement() *SpreadElement {
	return n.Data.(*SpreadElement)
}
func (n *Node) AsSpreadAssignment() *SpreadAssignment {
	return n.Data.(*SpreadAssignment)
}
func (n *Node) AsArrayTypeNode() *ArrayTypeNode {
	return n.Data.(*ArrayTypeNode)
}
func (n *Node) AsTupleTypeNode() *TupleTypeNode {
	return n.Data.(*TupleTypeNode)
}
func (n *Node) AsUnionTypeNode() *UnionTypeNode {
	return n.Data.(*UnionTypeNode)
}
func (n *Node) AsIntersectionTypeNode() *IntersectionTypeNode {
	return n.Data.(*IntersectionTypeNode)
}
func (n *Node) AsRestTypeNode() *RestTypeNode {
	return n.Data.(*RestTypeNode)
}
func (n *Node) AsNamedTupleMember() *NamedTupleMember {
	return n.Data.(*NamedTupleMember)
}
func (n *Node) AsOptionalTypeNode() *OptionalTypeNode {
	return n.Data.(*OptionalTypeNode)
}
func (n *Node) AsTypeReferenceNode() *TypeReferenceNode {
	return n.Data.(*TypeReferenceNode)
}
func (n *Node) AsTypeQueryNode() *TypeQueryNode {
	return n.Data.(*TypeQueryNode)
}
func (n *Node) AsIndexedAccessTypeNode() *IndexedAccessTypeNode {
	return n.Data.(*IndexedAccessTypeNode)
}
func (n *Node) AsGetAccessorDeclaration() *GetAccessorDeclaration {
	return n.Data.(*GetAccessorDeclaration)
}
func (n *Node) AsSetAccessorDeclaration() *SetAccessorDeclaration {
	return n.Data.(*SetAccessorDeclaration)
}

// NodeData

type NodeData interface {
	AsNode() *Node
	ForEachChild(v Visitor) bool
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
type JsxAttributeValue = Node           // StringLiteral | JsxExpression | JsxElement | JsxSelfClosingElement | JsxFragment
type JsxTagNameExpression = Node        // IdentifierReference | KeywordExpression | JsxTagNamePropertyAccess | JsxNamespacedName
type ClassLikeDeclaration = Node        // ClassDeclaration | ClassExpression
type AccessorDeclaration = Node         // GetAccessorDeclaration | SetAccessorDeclaration
type LiteralLikeNode = Node             // StringLiteral | NumericLiteral | BigIntLiteral | RegularExpressionLiteral | TemplateLiteralLikeNode | JsxText
type LiteralExpression = Node           // StringLiteral | NumericLiteral | BigIntLiteral | RegularExpressionLiteral | NoSubstitutionTemplateLiteral
type UnionOrIntersectionTypeNode = Node // UnionTypeNode | IntersectionTypeNode
type TemplateLiteralLikeNode = Node     // TemplateHead | TemplateMiddle | TemplateTail
type TemplateMiddleOrTail = Node        // TemplateMiddle | TemplateTail
type TemplateLiteral = Node             // TemplateExpression | NoSubstitutionTemplateLiteral
type TypePredicateParameterName = Node  // Identifier | ThisTypeNode
type ImportAttributeName = Node         // Identifier | StringLiteral
type LeftHandSideExpression = Node      // subset of Expression

// Aliases for node singletons

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
	Symbol *Symbol // Symbol declared by node (initialized by binding)
}

func (node *DeclarationBase) DeclarationData() *DeclarationBase { return node }

func IsDeclarationNode(node *Node) bool {
	return node.DeclarationData() != nil
}

// DeclarationBase

type ExportableBase struct {
	LocalSymbol *Symbol // Local symbol declared by node (initialized by binding only for exported nodes)
}

func (node *ExportableBase) ExportableData() *ExportableBase { return node }

// ModifiersBase

type ModifiersBase struct {
	Modifiers_ *ModifierListNode
}

func (node *ModifiersBase) Modifiers() *ModifierListNode { return node.Modifiers_ }

// LocalsContainerBase

type LocalsContainerBase struct {
	Locals_        SymbolTable // Locals associated with node (initialized by binding)
	NextContainer *Node       // Next container in declaration order (initialized by binding)
}

func (node *LocalsContainerBase) LocalsContainerData() *LocalsContainerBase { return node }

func IsLocalsContainer(node *Node) bool {
	return node.LocalsContainerData() != nil
}

// FunctionLikeBase

type FunctionLikeBase struct {
	LocalsContainerBase
	TypeParameters_ *TypeParameterListNode // Optional
	Parameters     []*ParameterDeclarationNode
	ReturnType     *TypeNode // Optional
}

func (node *FunctionLikeBase) TypeParameters() *TypeParameterListNode { return node.TypeParameters_ }
func (node *FunctionLikeBase) LocalsContainerData() *LocalsContainerBase {
	return &node.LocalsContainerBase
}
func (node *FunctionLikeBase) FunctionLikeData() *FunctionLikeBase { return node }
func (node *FunctionLikeBase) BodyData() *BodyBase                 { return nil }

// BodyBase

type BodyBase struct {
	AsteriskToken *TokenNode
	Body          *BlockOrExpression // Optional, can be Expression only in arrow functions
	EndFlowNode   *FlowNode
}

func (node *BodyBase) BodyData() *BodyBase { return node }

// FunctionLikeWithBodyBase

type FunctionLikeWithBodyBase struct {
	FunctionLikeBase
	BodyBase
}

func (node *FunctionLikeWithBodyBase) TypeParameters() *TypeParameterListNode {
	return node.TypeParameters_
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
	FlowNode *FlowNode
}

func (node *FlowNodeBase) FlowNodeData() *FlowNodeBase { return node }

// Token

type Token struct {
	NodeBase
}

func (f *NodeFactory) NewToken(kind ast.Kind) *Node {
	return f.NewNode(kind, &Token{})
}

// Identifier

type Identifier struct {
	ExpressionBase
	FlowNodeBase
	Text string
}

func (f *NodeFactory) NewIdentifier(text string) *Node {
	data := f.identifierPool.New()
	data.Text = text
	return f.NewNode(ast.KindIdentifier, data)
}

func IsIdentifier(node *Node) bool {
	return node.Kind == ast.KindIdentifier
}

// PrivateIdentifier

type PrivateIdentifier struct {
	ExpressionBase
	Text string
}

func (f *NodeFactory) NewPrivateIdentifier(text string) *Node {
	data := &PrivateIdentifier{}
	data.Text = text
	return f.NewNode(ast.KindPrivateIdentifier, data)
}

func IsPrivateIdentifier(node *Node) bool {
	return node.Kind == ast.KindPrivateIdentifier
}

// QualifiedName

type QualifiedName struct {
	NodeBase
	FlowNodeBase
	Left  *EntityName     // EntityName
	Right *IdentifierNode // IdentifierNode
}

func (f *NodeFactory) NewQualifiedName(left *EntityName, right *IdentifierNode) *Node {
	data := &QualifiedName{}
	data.Left = left
	data.Right = right
	return f.NewNode(ast.KindQualifiedName, data)
}

func (node *QualifiedName) ForEachChild(v Visitor) bool {
	return visit(v, node.Left) || visit(v, node.Right)
}

func IsQualifiedName(node *Node) bool {
	return node.Kind == ast.KindQualifiedName
}

// TypeParameterDeclaration

type TypeParameterDeclaration struct {
	NodeBase
	DeclarationBase
	ModifiersBase
	Name_        *IdentifierNode // IdentifierNode
	Constraint  *TypeNode       // TypeNode. Optional
	DefaultType *TypeNode       // TypeNode. Optional
	Expression  *Expression     // Expression. Optional, For error recovery purposes
}

func (f *NodeFactory) NewTypeParameterDeclaration(modifiers *ModifierListNode, name *IdentifierNode, constraint *TypeNode, defaultType *TypeNode) *Node {
	data := &TypeParameterDeclaration{}
	data.Modifiers_ = modifiers
	data.Name_ = name
	data.Constraint = constraint
	data.DefaultType = defaultType
	return f.NewNode(ast.KindTypeParameter, data)
}

func (node *TypeParameterDeclaration) Kind() ast.Kind {
	return ast.KindTypeParameter
}

func (node *TypeParameterDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.Modifiers_) || visit(v, node.Name_) || visit(v, node.Constraint) || visit(v, node.DefaultType)
}

func (node *TypeParameterDeclaration) Name() *DeclarationName {
	return node.Name_
}

func IsTypeParameterDeclaration(node *Node) bool {
	return node.Kind == ast.KindTypeParameter
}

// ComputedPropertyName

type ComputedPropertyName struct {
	NodeBase
	Expression *Expression // Expression
}

func (f *NodeFactory) NewComputedPropertyName(expression *Expression) *Node {
	data := &ComputedPropertyName{}
	data.Expression = expression
	return f.NewNode(ast.KindComputedPropertyName, data)
}

func (node *ComputedPropertyName) ForEachChild(v Visitor) bool {
	return visit(v, node.Expression)
}

func IsComputedPropertyName(node *Node) bool {
	return node.Kind == ast.KindComputedPropertyName
}

// Modifier

func (f *NodeFactory) NewModifier(kind ast.Kind) *Node {
	return f.NewToken(kind)
}

// Decorator

type Decorator struct {
	NodeBase
	Expression *LeftHandSideExpression // LeftHandSideExpression
}

func (f *NodeFactory) NewDecorator(expression *LeftHandSideExpression) *Node {
	data := &Decorator{}
	data.Expression = expression
	return f.NewNode(ast.KindDecorator, data)
}

func (node *Decorator) ForEachChild(v Visitor) bool {
	return visit(v, node.Expression)
}

func IsDecorator(node *Node) bool {
	return node.Kind == ast.KindDecorator
}

// ModifierList

type ModifierList struct {
	NodeBase
	Modifiers_     []*ModifierLike
	ModifierFlags ast.ModifierFlags
}

func (f *NodeFactory) NewModifierList(modifiers []*ModifierLike, modifierFlags ast.ModifierFlags) *Node {
	data := &ModifierList{}
	data.Modifiers_ = modifiers
	data.ModifierFlags = modifierFlags
	return f.NewNode(ast.KindModifierList, data)
}

func (node *ModifierList) ForEachChild(v Visitor) bool {
	return visitNodes(v, node.Modifiers_)
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
	return f.NewNode(ast.KindEmptyStatement, &EmptyStatement{})
}

func IsEmptyStatement(node *Node) bool {
	return node.Kind == ast.KindEmptyStatement
}

// IfStatement

type IfStatement struct {
	StatementBase
	Expression    *Expression // Expression
	ThenStatement *Statement  // Statement
	ElseStatement *Statement  // Statement. Optional
}

func (f *NodeFactory) NewIfStatement(expression *Expression, thenStatement *Statement, elseStatement *Statement) *Node {
	data := &IfStatement{}
	data.Expression = expression
	data.ThenStatement = thenStatement
	data.ElseStatement = elseStatement
	return f.NewNode(ast.KindIfStatement, data)
}

func (node *IfStatement) ForEachChild(v Visitor) bool {
	return visit(v, node.Expression) || visit(v, node.ThenStatement) || visit(v, node.ElseStatement)
}

// DoStatement

type DoStatement struct {
	StatementBase
	Statement  *Statement  // Statement
	Expression *Expression // Expression
}

func (f *NodeFactory) NewDoStatement(statement *Statement, expression *Expression) *Node {
	data := &DoStatement{}
	data.Statement = statement
	data.Expression = expression
	return f.NewNode(ast.KindDoStatement, data)
}

func (node *DoStatement) ForEachChild(v Visitor) bool {
	return visit(v, node.Statement) || visit(v, node.Expression)
}

// WhileStatement

type WhileStatement struct {
	StatementBase
	Expression *Expression // Expression
	Statement  *Statement  // Statement
}

func (f *NodeFactory) NewWhileStatement(expression *Expression, statement *Statement) *Node {
	data := &WhileStatement{}
	data.Expression = expression
	data.Statement = statement
	return f.NewNode(ast.KindWhileStatement, data)
}

func (node *WhileStatement) ForEachChild(v Visitor) bool {
	return visit(v, node.Expression) || visit(v, node.Statement)
}

// ForStatement

type ForStatement struct {
	StatementBase
	LocalsContainerBase
	Initializer *ForInitializer // ForInitializer. Optional
	Condition   *Expression     // Expression. Optional
	Incrementor *Expression     // Expression. Optional
	Statement   *Statement      // Statement
}

func (f *NodeFactory) NewForStatement(initializer *ForInitializer, condition *Expression, incrementor *Expression, statement *Statement) *Node {
	data := &ForStatement{}
	data.Initializer = initializer
	data.Condition = condition
	data.Incrementor = incrementor
	data.Statement = statement
	return f.NewNode(ast.KindForStatement, data)
}

func (node *ForStatement) ForEachChild(v Visitor) bool {
	return visit(v, node.Initializer) || visit(v, node.Condition) || visit(v, node.Incrementor) || visit(v, node.Statement)
}

// ForInOrOfStatement

type ForInOrOfStatement struct {
	StatementBase
	LocalsContainerBase
	AwaitModifier *TokenNode      // TokenNode. Optional
	Initializer   *ForInitializer // ForInitializer
	Expression    *Expression     // Expression
	Statement     *Statement      // Statement
}

func (f *NodeFactory) NewForInOrOfStatement(kind ast.Kind, awaitModifier *TokenNode, initializer *ForInitializer, expression *Expression, statement *Statement) *Node {
	data := &ForInOrOfStatement{}
	data.AwaitModifier = awaitModifier
	data.Initializer = initializer
	data.Expression = expression
	data.Statement = statement
	return f.NewNode(kind, data)
}

func (node *ForInOrOfStatement) ForEachChild(v Visitor) bool {
	return visit(v, node.AwaitModifier) || visit(v, node.Initializer) || visit(v, node.Expression) || visit(v, node.Statement)
}

func IsForInOrOfStatement(node *Node) bool {
	return node.Kind == ast.KindForInStatement || node.Kind == ast.KindForOfStatement
}

// BreakStatement

type BreakStatement struct {
	StatementBase
	Label *IdentifierNode // IdentifierNode. Optional
}

func (f *NodeFactory) NewBreakStatement(label *IdentifierNode) *Node {
	data := &BreakStatement{}
	data.Label = label
	return f.NewNode(ast.KindBreakStatement, data)
}

func (node *BreakStatement) ForEachChild(v Visitor) bool {
	return visit(v, node.Label)
}

// ContinueStatement

type ContinueStatement struct {
	StatementBase
	Label *IdentifierNode // IdentifierNode. Optional
}

func (f *NodeFactory) NewContinueStatement(label *IdentifierNode) *Node {
	data := &ContinueStatement{}
	data.Label = label
	return f.NewNode(ast.KindContinueStatement, data)
}

func (node *ContinueStatement) ForEachChild(v Visitor) bool {
	return visit(v, node.Label)
}

// ReturnStatement

type ReturnStatement struct {
	StatementBase
	Expression *Expression // Expression. Optional
}

func (f *NodeFactory) NewReturnStatement(expression *Expression) *Node {
	data := &ReturnStatement{}
	data.Expression = expression
	return f.NewNode(ast.KindReturnStatement, data)
}

func (node *ReturnStatement) ForEachChild(v Visitor) bool {
	return visit(v, node.Expression)
}

// WithStatement

type WithStatement struct {
	StatementBase
	Expression *Expression // Expression
	Statement  *Statement  // Statement
}

func (f *NodeFactory) NewWithStatement(expression *Expression, statement *Statement) *Node {
	data := &WithStatement{}
	data.Expression = expression
	data.Statement = statement
	return f.NewNode(ast.KindWithStatement, data)
}

func (node *WithStatement) ForEachChild(v Visitor) bool {
	return visit(v, node.Expression) || visit(v, node.Statement)
}

// SwitchStatement

type SwitchStatement struct {
	StatementBase
	Expression *Expression    // Expression
	CaseBlock  *CaseBlockNode // CaseBlockNode
}

func (f *NodeFactory) NewSwitchStatement(expression *Expression, caseBlock *CaseBlockNode) *Node {
	data := &SwitchStatement{}
	data.Expression = expression
	data.CaseBlock = caseBlock
	return f.NewNode(ast.KindSwitchStatement, data)
}

func (node *SwitchStatement) ForEachChild(v Visitor) bool {
	return visit(v, node.Expression) || visit(v, node.CaseBlock)
}

// CaseBlock

type CaseBlock struct {
	NodeBase
	LocalsContainerBase
	Clauses []*CaseOrDefaultClauseNode
}

func (f *NodeFactory) NewCaseBlock(clauses []*CaseOrDefaultClauseNode) *Node {
	data := &CaseBlock{}
	data.Clauses = clauses
	return f.NewNode(ast.KindCaseBlock, data)
}

func (node *CaseBlock) ForEachChild(v Visitor) bool {
	return visitNodes(v, node.Clauses)
}

// CaseOrDefaultClause

type CaseOrDefaultClause struct {
	NodeBase
	Expression          *Expression  // Expression. nil in default clause
	Statements          []*Statement // []Statement
	FallthroughFlowNode *FlowNode
}

func (f *NodeFactory) NewCaseOrDefaultClause(kind ast.Kind, expression *Expression, statements []*Statement) *Node {
	data := &CaseOrDefaultClause{}
	data.Expression = expression
	data.Statements = statements
	return f.NewNode(kind, data)
}

func (node *CaseOrDefaultClause) ForEachChild(v Visitor) bool {
	return visit(v, node.Expression) || visitNodes(v, node.Statements)
}

// ThrowStatement

type ThrowStatement struct {
	StatementBase
	Expression *Expression // Expression
}

func (f *NodeFactory) NewThrowStatement(expression *Expression) *Node {
	data := &ThrowStatement{}
	data.Expression = expression
	return f.NewNode(ast.KindThrowStatement, data)
}

func (node *ThrowStatement) ForEachChild(v Visitor) bool {
	return visit(v, node.Expression)
}

// TryStatement

type TryStatement struct {
	StatementBase
	TryBlock     *BlockNode       // BlockNode
	CatchClause  *CatchClauseNode // CatchClauseNode. Optional
	FinallyBlock *BlockNode       // BlockNode. Optional
}

func (f *NodeFactory) NewTryStatement(tryBlock *BlockNode, catchClause *CatchClauseNode, finallyBlock *BlockNode) *Node {
	data := &TryStatement{}
	data.TryBlock = tryBlock
	data.CatchClause = catchClause
	data.FinallyBlock = finallyBlock
	return f.NewNode(ast.KindTryStatement, data)
}

func (node *TryStatement) Kind() ast.Kind {
	return ast.KindTryStatement
}

func (node *TryStatement) ForEachChild(v Visitor) bool {
	return visit(v, node.TryBlock) || visit(v, node.CatchClause) || visit(v, node.FinallyBlock)
}

// CatchClause

type CatchClause struct {
	NodeBase
	LocalsContainerBase
	VariableDeclaration *VariableDeclarationNode // Optional
	Block               *BlockNode
}

func (f *NodeFactory) NewCatchClause(variableDeclaration *VariableDeclarationNode, block *BlockNode) *Node {
	data := &CatchClause{}
	data.VariableDeclaration = variableDeclaration
	data.Block = block
	return f.NewNode(ast.KindCatchClause, data)
}

func (node *CatchClause) ForEachChild(v Visitor) bool {
	return visit(v, node.VariableDeclaration) || visit(v, node.Block)
}

// DebuggerStatement

type DebuggerStatement struct {
	StatementBase
}

func (f *NodeFactory) NewDebuggerStatement() *Node {
	return f.NewNode(ast.KindDebuggerStatement, &DebuggerStatement{})
}

// LabeledStatement

type LabeledStatement struct {
	StatementBase
	Label     *IdentifierNode // IdentifierNode
	Statement *Statement      // Statement
}

func (f *NodeFactory) NewLabeledStatement(label *IdentifierNode, statement *Statement) *Node {
	data := &LabeledStatement{}
	data.Label = label
	data.Statement = statement
	return f.NewNode(ast.KindLabeledStatement, data)
}

func (node *LabeledStatement) ForEachChild(v Visitor) bool {
	return visit(v, node.Label) || visit(v, node.Statement)
}

// ExpressionStatement

type ExpressionStatement struct {
	StatementBase
	Expression *Expression // Expression
}

func (f *NodeFactory) NewExpressionStatement(expression *Expression) *Node {
	data := &ExpressionStatement{}
	data.Expression = expression
	return f.NewNode(ast.KindExpressionStatement, data)
}

func (node *ExpressionStatement) ForEachChild(v Visitor) bool {
	return visit(v, node.Expression)
}

func IsExpressionStatement(node *Node) bool {
	return node.Kind == ast.KindExpressionStatement
}

// Block

type Block struct {
	StatementBase
	LocalsContainerBase
	Statements []*Statement
	Multiline  bool
}

func (f *NodeFactory) NewBlock(statements []*Statement, multiline bool) *Node {
	data := &Block{}
	data.Statements = statements
	data.Multiline = multiline
	return f.NewNode(ast.KindBlock, data)
}

func (node *Block) ForEachChild(v Visitor) bool {
	return visitNodes(v, node.Statements)
}

func IsBlock(node *Node) bool {
	return node.Kind == ast.KindBlock
}

// VariableStatement

type VariableStatement struct {
	StatementBase
	ModifiersBase
	DeclarationList *VariableDeclarationListNode // VariableDeclarationListNode
}

func (f *NodeFactory) NewVariableStatement(modifiers *ModifierListNode, declarationList *VariableDeclarationListNode) *Node {
	data := &VariableStatement{}
	data.Modifiers_ = modifiers
	data.DeclarationList = declarationList
	return f.NewNode(ast.KindVariableStatement, data)
}

func (node *VariableStatement) ForEachChild(v Visitor) bool {
	return visit(v, node.Modifiers_) || visit(v, node.DeclarationList)
}

func IsVariableStatement(node *Node) bool {
	return node.Kind == ast.KindVariableStatement
}

// VariableDeclaration

type VariableDeclaration struct {
	NodeBase
	DeclarationBase
	ExportableBase
	Name_             *BindingName // BindingName
	ExclamationToken *TokenNode   // TokenNode. Optional
	TypeNode         *TypeNode    // TypeNode. Optional
	Initializer      *Expression  // Expression. Optional
}

func (f *NodeFactory) NewVariableDeclaration(name *BindingName, exclamationToken *TokenNode, typeNode *TypeNode, initializer *Expression) *Node {
	data := &VariableDeclaration{}
	data.Name_ = name
	data.ExclamationToken = exclamationToken
	data.TypeNode = typeNode
	data.Initializer = initializer
	return f.NewNode(ast.KindVariableDeclaration, data)
}

func (node *VariableDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.Name_) || visit(v, node.ExclamationToken) || visit(v, node.TypeNode) || visit(v, node.Initializer)
}

func (node *VariableDeclaration) Name() *DeclarationName {
	return node.Name_
}

func IsVariableDeclaration(node *Node) bool {
	return node.Kind == ast.KindVariableDeclaration
}

// VariableDeclarationList

type VariableDeclarationList struct {
	NodeBase
	Declarations []*VariableDeclarationNode
}

func (f *NodeFactory) NewVariableDeclarationList(flags ast.NodeFlags, declarations []*VariableDeclarationNode) *Node {
	data := &VariableDeclarationList{}
	data.Declarations = declarations
	node := f.NewNode(ast.KindVariableDeclarationList, data)
	node.Flags = flags
	return node
}

func (node *VariableDeclarationList) ForEachChild(v Visitor) bool {
	return visitNodes(v, node.Declarations)
}

func IsVariableDeclarationList(node *Node) bool {
	return node.Kind == ast.KindVariableDeclarationList
}

// BindingPattern (SyntaxBindObjectBindingPattern | ast.KindArrayBindingPattern)

type BindingPattern struct {
	NodeBase
	Elements []*BindingElementNode
}

func (f *NodeFactory) NewBindingPattern(kind ast.Kind, elements []*BindingElementNode) *Node {
	data := &BindingPattern{}
	data.Elements = elements
	return f.NewNode(kind, data)
}

func (node *BindingPattern) ForEachChild(v Visitor) bool {
	return visitNodes(v, node.Elements)
}

func IsObjectBindingPattern(node *Node) bool {
	return node.Kind == ast.KindObjectBindingPattern
}

func IsArrayBindingPattern(node *Node) bool {
	return node.Kind == ast.KindArrayBindingPattern
}

// ParameterDeclaration

type ParameterDeclaration struct {
	NodeBase
	DeclarationBase
	ModifiersBase
	DotDotDotToken *TokenNode   // TokenNode. Present on rest parameter
	Name_           *BindingName // BindingName. Declared parameter name
	QuestionToken  *TokenNode   // TokenNode. Present on optional parameter
	TypeNode       *TypeNode    // TypeNode. Optional
	Initializer    *Expression  // Expression. Optional
}

func (f *NodeFactory) NewParameterDeclaration(modifiers *ModifierListNode, dotDotDotToken *TokenNode, name *BindingName, questionToken *TokenNode, typeNode *TypeNode, initializer *Expression) *Node {
	data := &ParameterDeclaration{}
	data.Modifiers_ = modifiers
	data.DotDotDotToken = dotDotDotToken
	data.Name_ = name
	data.QuestionToken = questionToken
	data.TypeNode = typeNode
	data.Initializer = initializer
	return f.NewNode(ast.KindParameter, data)
}

func (node *ParameterDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.Modifiers_) || visit(v, node.DotDotDotToken) || visit(v, node.Name_) ||
		visit(v, node.QuestionToken) || visit(v, node.TypeNode) || visit(v, node.Initializer)
}

func (node *ParameterDeclaration) Name() *DeclarationName {
	return node.Name_
}

func IsParameter(node *Node) bool {
	return node.Kind == ast.KindParameter
}

// BindingElement

type BindingElement struct {
	NodeBase
	DeclarationBase
	ExportableBase
	FlowNodeBase
	DotDotDotToken *TokenNode    // TokenNode. Present on rest element (in object binding pattern)
	PropertyName   *PropertyName // PropertyName. Optional binding property name in object binding pattern
	Name_           *BindingName  // BindingName. Optional (nil for missing element)
	Initializer    *Expression   // Expression. Optional
}

func (f *NodeFactory) NewBindingElement(dotDotDotToken *TokenNode, propertyName *PropertyName, name *BindingName, initializer *Expression) *Node {
	data := &BindingElement{}
	data.DotDotDotToken = dotDotDotToken
	data.PropertyName = propertyName
	data.Name_ = name
	data.Initializer = initializer
	return f.NewNode(ast.KindBindingElement, data)
}

func (node *BindingElement) ForEachChild(v Visitor) bool {
	return visit(v, node.PropertyName) || visit(v, node.DotDotDotToken) || visit(v, node.Name_) || visit(v, node.Initializer)
}

func (node *BindingElement) Name() *DeclarationName {
	return node.Name_
}

func IsBindingElement(node *Node) bool {
	return node.Kind == ast.KindBindingElement
}

// MissingDeclaration

type MissingDeclaration struct {
	StatementBase
	DeclarationBase
	ModifiersBase
}

func (f *NodeFactory) NewMissingDeclaration(modifiers *ModifierListNode) *Node {
	data := &MissingDeclaration{}
	data.Modifiers_ = modifiers
	return f.NewNode(ast.KindMissingDeclaration, data)
}

func (node *MissingDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.Modifiers_)
}

// FunctionDeclaration

type FunctionDeclaration struct {
	StatementBase
	DeclarationBase
	ExportableBase
	ModifiersBase
	FunctionLikeWithBodyBase
	Name_           *IdentifierNode // IdentifierNode
	ReturnFlowNode *FlowNode
}

func (f *NodeFactory) NewFunctionDeclaration(modifiers *ModifierListNode, asteriskToken *TokenNode, name *IdentifierNode, typeParameters *TypeParameterListNode, parameters []*ParameterDeclarationNode, returnType *TypeNode, body *BlockNode) *Node {
	data := &FunctionDeclaration{}
	data.Modifiers_ = modifiers
	data.AsteriskToken = asteriskToken
	data.Name_ = name
	data.TypeParameters_ = typeParameters
	data.Parameters = parameters
	data.ReturnType = returnType
	data.Body = body
	return f.NewNode(ast.KindFunctionDeclaration, data)
}

func (node *FunctionDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.Modifiers_) || visit(v, node.AsteriskToken) || visit(v, node.Name_) || visit(v, node.TypeParameters_) ||
		visitNodes(v, node.Parameters) || visit(v, node.ReturnType) || visit(v, node.Body)
}

func (node *FunctionDeclaration) Name() *DeclarationName {
	return node.Name_
}

func (node *FunctionDeclaration) BodyData() *BodyBase { return &node.BodyBase }

func IsFunctionDeclaration(node *Node) bool {
	return node.Kind == ast.KindFunctionDeclaration
}

// ClassLikeDeclarationBase

type ClassLikeBase struct {
	DeclarationBase
	ExportableBase
	ModifiersBase
	Name_            *IdentifierNode        // IdentifierNode
	TypeParameters_  *TypeParameterListNode // TypeParameterListNode
	HeritageClauses []*HeritageClauseNode  // []HeritageClauseNode
	Members         []*ClassElement        // []ClassElement
}

func (node *ClassLikeBase) ForEachChild(v Visitor) bool {
	return visit(v, node.Modifiers_) || visit(v, node.Name_) || visit(v, node.TypeParameters_) ||
		visitNodes(v, node.HeritageClauses) || visitNodes(v, node.Members)
}

func (node *ClassLikeBase) Name() *DeclarationName                 { return node.Name_ }
func (node *ClassLikeBase) TypeParameters() *TypeParameterListNode { return node.TypeParameters_ }
func (node *ClassLikeBase) ClassLikeData() *ClassLikeBase          { return node }

// ClassDeclaration

type ClassDeclaration struct {
	StatementBase
	ClassLikeBase
}

func (f *NodeFactory) NewClassDeclaration(modifiers *ModifierListNode, name *IdentifierNode, typeParameters *TypeParameterListNode, heritageClauses []*HeritageClauseNode, members []*ClassElement) *Node {
	data := &ClassDeclaration{}
	data.Modifiers_ = modifiers
	data.Name_ = name
	data.TypeParameters_ = typeParameters
	data.HeritageClauses = heritageClauses
	data.Members = members
	return f.NewNode(ast.KindClassDeclaration, data)
}

func IsClassDeclaration(node *Node) bool {
	return node.Kind == ast.KindClassDeclaration
}

// ClassExpression

type ClassExpression struct {
	ExpressionBase
	ClassLikeBase
}

func (f *NodeFactory) NewClassExpression(modifiers *ModifierListNode, name *IdentifierNode, typeParameters *TypeParameterListNode, heritageClauses []*HeritageClauseNode, members []*ClassElement) *Node {
	data := &ClassExpression{}
	data.Modifiers_ = modifiers
	data.Name_ = name
	data.TypeParameters_ = typeParameters
	data.HeritageClauses = heritageClauses
	data.Members = members
	return f.NewNode(ast.KindClassExpression, data)
}

func (node *ClassExpression) Kind() ast.Kind { return ast.KindClassExpression }

func IsClassExpression(node *Node) bool {
	return node.Kind == ast.KindClassExpression
}

// HeritageClause

type HeritageClause struct {
	NodeBase
	Token ast.Kind
	Types []*ExpressionWithTypeArgumentsNode
}

func (f *NodeFactory) NewHeritageClause(token ast.Kind, types []*ExpressionWithTypeArgumentsNode) *Node {
	data := &HeritageClause{}
	data.Token = token
	data.Types = types
	return f.NewNode(ast.KindHeritageClause, data)
}

func (node *HeritageClause) ForEachChild(v Visitor) bool {
	return visitNodes(v, node.Types)
}

func IsHeritageClause(node *Node) bool {
	return node.Kind == ast.KindHeritageClause
}

// InterfaceDeclaration

type InterfaceDeclaration struct {
	StatementBase
	DeclarationBase
	ExportableBase
	ModifiersBase
	Name_            *IdentifierNode
	TypeParameters_  *TypeParameterListNode
	HeritageClauses []*HeritageClauseNode
	Members         []*TypeElement
}

func (f *NodeFactory) NewInterfaceDeclaration(modifiers *ModifierListNode, name *IdentifierNode, typeParameters *TypeParameterListNode, heritageClauses []*HeritageClauseNode, members []*TypeElement) *Node {
	data := &InterfaceDeclaration{}
	data.Modifiers_ = modifiers
	data.Name_ = name
	data.TypeParameters_ = typeParameters
	data.HeritageClauses = heritageClauses
	data.Members = members
	return f.NewNode(ast.KindInterfaceDeclaration, data)
}

func (node *InterfaceDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.Modifiers_) || visit(v, node.Name_) || visit(v, node.TypeParameters_) ||
		visitNodes(v, node.HeritageClauses) || visitNodes(v, node.Members)
}

func (node *InterfaceDeclaration) Name() *DeclarationName                 { return node.Name_ }
func (node *InterfaceDeclaration) TypeParameters() *TypeParameterListNode { return node.TypeParameters_ }

func IsInterfaceDeclaration(node *Node) bool {
	return node.Kind == ast.KindInterfaceDeclaration
}

// TypeAliasDeclaration

type TypeAliasDeclaration struct {
	StatementBase
	DeclarationBase
	ExportableBase
	ModifiersBase
	LocalsContainerBase
	Name_           *IdentifierNode        // IdentifierNode
	TypeParameters_ *TypeParameterListNode // TypeParameterListNode. Optional
	TypeNode       *TypeNode              // TypeNode
}

func (f *NodeFactory) NewTypeAliasDeclaration(modifiers *ModifierListNode, name *IdentifierNode, typeParameters *TypeParameterListNode, typeNode *TypeNode) *Node {
	data := &TypeAliasDeclaration{}
	data.Modifiers_ = modifiers
	data.Name_ = name
	data.TypeParameters_ = typeParameters
	data.TypeNode = typeNode
	return f.NewNode(ast.KindTypeAliasDeclaration, data)
}

func (node *TypeAliasDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.Modifiers_) || visit(v, node.Name_) || visit(v, node.TypeParameters_) || visit(v, node.TypeNode)
}

func (node *TypeAliasDeclaration) Name() *DeclarationName                 { return node.Name_ }
func (node *TypeAliasDeclaration) TypeParameters() *TypeParameterListNode { return node.TypeParameters_ }

func IsTypeAliasDeclaration(node *Node) bool {
	return node.Kind == ast.KindTypeAliasDeclaration
}

// EnumMember

type EnumMember struct {
	NodeBase
	NamedMemberBase
	Initializer *Expression // Expression. Optional
}

func (f *NodeFactory) NewEnumMember(name *PropertyName, initializer *Expression) *Node {
	data := &EnumMember{}
	data.Name_ = name
	data.Initializer = initializer
	return f.NewNode(ast.KindEnumMember, data)
}

func (node *EnumMember) ForEachChild(v Visitor) bool {
	return visit(v, node.Name_) || visit(v, node.Initializer)
}

func (node *EnumMember) Name() *DeclarationName {
	return node.Name_
}

// EnumDeclaration

type EnumDeclaration struct {
	StatementBase
	DeclarationBase
	ExportableBase
	ModifiersBase
	Name_    *IdentifierNode   // IdentifierNode
	Members []*EnumMemberNode // []EnumMemberNode
}

func (f *NodeFactory) NewEnumDeclaration(modifiers *ModifierListNode, name *IdentifierNode, members []*EnumMemberNode) *Node {
	data := &EnumDeclaration{}
	data.Modifiers_ = modifiers
	data.Name_ = name
	data.Members = members
	return f.NewNode(ast.KindEnumDeclaration, data)
}

func (node *EnumDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.Modifiers_) || visit(v, node.Name_) || visitNodes(v, node.Members)
}

func (node *EnumDeclaration) Name() *DeclarationName {
	return node.Name_
}

func IsEnumDeclaration(node *Node) bool {
	return node.Kind == ast.KindEnumDeclaration
}

// ModuleBlock

type ModuleBlock struct {
	StatementBase
	Statements []*Statement
}

func (f *NodeFactory) NewModuleBlock(statements []*Statement) *Node {
	data := &ModuleBlock{}
	data.Statements = statements
	return f.NewNode(ast.KindModuleBlock, data)
}

func (node *ModuleBlock) ForEachChild(v Visitor) bool {
	return visitNodes(v, node.Statements)
}

func IsModuleBlock(node *Node) bool {
	return node.Kind == ast.KindModuleBlock
}

// ModuleDeclaration

type ModuleDeclaration struct {
	StatementBase
	DeclarationBase
	ExportableBase
	ModifiersBase
	LocalsContainerBase
	Name_ *ModuleName // ModuleName
	Body *ModuleBody // ModuleBody. Optional (may be nil in ambient module declaration)
}

func (f *NodeFactory) NewModuleDeclaration(modifiers *ModifierListNode, name *ModuleName, body *ModuleBody, flags ast.NodeFlags) *Node {
	data := &ModuleDeclaration{}
	data.Modifiers_ = modifiers
	data.Name_ = name
	data.Body = body
	node := f.NewNode(ast.KindModuleDeclaration, data)
	node.Flags |= flags & (ast.NodeFlagsNamespace | ast.NodeFlagsNestedNamespace | ast.NodeFlagsGlobalAugmentation)
	return node
}

func (node *ModuleDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.Modifiers_) || visit(v, node.Name_) || visit(v, node.Body)
}

func (node *ModuleDeclaration) Name() *DeclarationName {
	return node.Name_
}

func IsModuleDeclaration(node *Node) bool {
	return node.Kind == ast.KindModuleDeclaration
}

// ModuleEqualsDeclaration

type ImportEqualsDeclaration struct {
	StatementBase
	DeclarationBase
	ExportableBase
	ModifiersBase
	Modifiers_  *ModifierListNode // ModifierListNode
	IsTypeOnly bool
	Name_       *IdentifierNode // IdentifierNode
	// 'EntityName' for an internal module reference, 'ExternalModuleReference' for an external
	// module reference.
	ModuleReference *ModuleReference // ModuleReference
}

func (f *NodeFactory) NewImportEqualsDeclaration(modifiers *ModifierListNode, isTypeOnly bool, name *IdentifierNode, moduleReference *ModuleReference) *Node {
	data := &ImportEqualsDeclaration{}
	data.Modifiers_ = modifiers
	data.IsTypeOnly = isTypeOnly
	data.Name_ = name
	data.ModuleReference = moduleReference
	return f.NewNode(ast.KindImportEqualsDeclaration, data)
}

func (node *ImportEqualsDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.Modifiers_) || visit(v, node.Name_) || visit(v, node.ModuleReference)
}

func (node *ImportEqualsDeclaration) Name() *DeclarationName {
	return node.Name_
}

func IsImportEqualsDeclaration(node *Node) bool {
	return node.Kind == ast.KindImportEqualsDeclaration
}

// ImportDeclaration

type ImportDeclaration struct {
	StatementBase
	ModifiersBase
	ImportClause    *ImportClauseNode     // ImportClauseNode. Optional
	ModuleSpecifier *Expression           // Expression
	Attributes      *ImportAttributesNode // ImportAttributesNode. Optional
}

func (f *NodeFactory) NewImportDeclaration(modifiers *ModifierListNode, importClause *ImportClauseNode, moduleSpecifier *Expression, attributes *ImportAttributesNode) *Node {
	data := &ImportDeclaration{}
	data.Modifiers_ = modifiers
	data.ImportClause = importClause
	data.ModuleSpecifier = moduleSpecifier
	data.Attributes = attributes
	return f.NewNode(ast.KindImportDeclaration, data)
}

func (node *ImportDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.Modifiers_) || visit(v, node.ImportClause) || visit(v, node.ModuleSpecifier) || visit(v, node.Attributes)
}

func IsImportDeclaration(node *Node) bool {
	return node.Kind == ast.KindImportDeclaration
}

// ImportSpecifier

type ImportSpecifier struct {
	NodeBase
	DeclarationBase
	ExportableBase
	IsTypeOnly   bool
	PropertyName *ModuleExportName // ModuleExportName. Optional
	Name_         *IdentifierNode   // IdentifierNode
}

func (f *NodeFactory) NewImportSpecifier(isTypeOnly bool, propertyName *ModuleExportName, name *IdentifierNode) *Node {
	data := &ImportSpecifier{}
	data.IsTypeOnly = isTypeOnly
	data.PropertyName = propertyName
	data.Name_ = name
	return f.NewNode(ast.KindImportSpecifier, data)
}

func (node *ImportSpecifier) ForEachChild(v Visitor) bool {
	return visit(v, node.PropertyName) || visit(v, node.Name_)
}

func (node *ImportSpecifier) Name() *DeclarationName {
	return node.Name_
}

func IsImportSpecifier(node *Node) bool {
	return node.Kind == ast.KindImportSpecifier
}

// ExternalModuleReference

type ExternalModuleReference struct {
	NodeBase
	expression *Expression // Expression
}

func (f *NodeFactory) NewExternalModuleReference(expression *Expression) *Node {
	data := &ExternalModuleReference{}
	data.expression = expression
	return f.NewNode(ast.KindExternalModuleReference, data)
}

func (node *ExternalModuleReference) ForEachChild(v Visitor) bool {
	return visit(v, node.expression)
}

func IsExternalModuleReference(node *Node) bool {
	return node.Kind == ast.KindExternalModuleReference
}

// ImportClause

type ImportClause struct {
	NodeBase
	DeclarationBase
	ExportableBase
	IsTypeOnly    bool
	NamedBindings *NamedImportBindings // NamedImportBindings. Optional, named bindings
	Name_          *IdentifierNode      // IdentifierNode. Optional, default binding
}

func (f *NodeFactory) NewImportClause(isTypeOnly bool, name *IdentifierNode, namedBindings *NamedImportBindings) *Node {
	data := &ImportClause{}
	data.IsTypeOnly = isTypeOnly
	data.Name_ = name
	data.NamedBindings = namedBindings
	return f.NewNode(ast.KindImportClause, data)
}

func (node *ImportClause) ForEachChild(v Visitor) bool {
	return visit(v, node.Name_) || visit(v, node.NamedBindings)
}

func (node *ImportClause) Name() *DeclarationName {
	return node.Name_
}

// NamespaceImport

type NamespaceImport struct {
	NodeBase
	DeclarationBase
	ExportableBase
	Name_ *IdentifierNode // IdentifierNode
}

func (f *NodeFactory) NewNamespaceImport(name *IdentifierNode) *Node {
	data := &NamespaceImport{}
	data.Name_ = name
	return f.NewNode(ast.KindNamespaceImport, data)
}

func (node *NamespaceImport) ForEachChild(v Visitor) bool {
	return visit(v, node.Name_)
}

func (node *NamespaceImport) Name() *DeclarationName {
	return node.Name_
}

func IsNamespaceImport(node *Node) bool {
	return node.Kind == ast.KindNamespaceImport
}

// NamedImports

type NamedImports struct {
	NodeBase
	Elements []*ImportSpecifierNode
}

func (f *NodeFactory) NewNamedImports(elements []*ImportSpecifierNode) *Node {
	data := &NamedImports{}
	data.Elements = elements
	return f.NewNode(ast.KindNamedImports, data)
}

func (node *NamedImports) ForEachChild(v Visitor) bool {
	return visitNodes(v, node.Elements)
}

func IsNamedImports(node *Node) bool {
	return node.Kind == ast.KindNamedImports
}

// ExportAssignment

// This is either an `export =` or an `export default` declaration.
// Unless `isExportEquals` is set, this node was parsed as an `export default`.
type ExportAssignment struct {
	StatementBase
	DeclarationBase
	ModifiersBase
	IsExportEquals bool
	Expression     *Expression // Expression
}

func (f *NodeFactory) NewExportAssignment(modifiers *ModifierListNode, isExportEquals bool, expression *Expression) *Node {
	data := &ExportAssignment{}
	data.Modifiers_ = modifiers
	data.IsExportEquals = isExportEquals
	data.Expression = expression
	return f.NewNode(ast.KindExportAssignment, data)
}

func (node *ExportAssignment) ForEachChild(v Visitor) bool {
	return visit(v, node.Modifiers_) || visit(v, node.Expression)
}

func IsExportAssignment(node *Node) bool {
	return node.Kind == ast.KindExportAssignment
}

// NamespaceExportDeclaration

type NamespaceExportDeclaration struct {
	StatementBase
	DeclarationBase
	ModifiersBase
	Name_ *IdentifierNode // IdentifierNode
}

func (f *NodeFactory) NewNamespaceExportDeclaration(modifiers *ModifierListNode, name *IdentifierNode) *Node {
	data := &NamespaceExportDeclaration{}
	data.Modifiers_ = modifiers
	data.Name_ = name
	return f.NewNode(ast.KindNamespaceExportDeclaration, data)
}

func (node *NamespaceExportDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.Modifiers_) || visit(v, node.Name_)
}

func (node *NamespaceExportDeclaration) Name() *DeclarationName {
	return node.Name_
}

func IsNamespaceExportDeclaration(node *Node) bool {
	return node.Kind == ast.KindNamespaceExportDeclaration
}

// ExportDeclaration

type ExportDeclaration struct {
	StatementBase
	DeclarationBase
	ModifiersBase
	IsTypeOnly      bool
	ExportClause    *NamedExportBindings  // NamedExportBindings. Optional
	ModuleSpecifier *Expression           // Expression. Optional
	Attributes      *ImportAttributesNode // ImportAttributesNode. Optional
}

func (f *NodeFactory) NewExportDeclaration(modifiers *ModifierListNode, isTypeOnly bool, exportClause *NamedExportBindings, moduleSpecifier *Expression, attributes *ImportAttributesNode) *Node {
	data := &ExportDeclaration{}
	data.Modifiers_ = modifiers
	data.IsTypeOnly = isTypeOnly
	data.ExportClause = exportClause
	data.ModuleSpecifier = moduleSpecifier
	data.Attributes = attributes
	return f.NewNode(ast.KindExportDeclaration, data)
}

func (node *ExportDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.Modifiers_) || visit(v, node.ExportClause) || visit(v, node.ModuleSpecifier) || visit(v, node.Attributes)
}

func IsExportDeclaration(node *Node) bool {
	return node.Kind == ast.KindExportDeclaration
}

// NamespaceExport

type NamespaceExport struct {
	NodeBase
	DeclarationBase
	Name_ *ModuleExportName // ModuleExportName
}

func (f *NodeFactory) NewNamespaceExport(name *ModuleExportName) *Node {
	data := &NamespaceExport{}
	data.Name_ = name
	return f.NewNode(ast.KindNamespaceExport, data)
}

func (node *NamespaceExport) ForEachChild(v Visitor) bool {
	return visit(v, node.Name_)
}

func (node *NamespaceExport) Name() *DeclarationName {
	return node.Name_
}

func IsNamespaceExport(node *Node) bool {
	return node.Kind == ast.KindNamespaceExport
}

// NamedExports

type NamedExports struct {
	NodeBase
	Elements []*ExportSpecifierNode // []ExportSpecifierNode
}

func (f *NodeFactory) NewNamedExports(elements []*ExportSpecifierNode) *Node {
	data := &NamedExports{}
	data.Elements = elements
	return f.NewNode(ast.KindNamedExports, data)
}

func (node *NamedExports) ForEachChild(v Visitor) bool {
	return visitNodes(v, node.Elements)
}

// ExportSpecifier

type ExportSpecifier struct {
	NodeBase
	DeclarationBase
	ExportableBase
	IsTypeOnly   bool
	PropertyName *ModuleExportName // ModuleExportName. Optional, name preceding 'as' keyword
	Name_         *ModuleExportName // ModuleExportName
}

func (f *NodeFactory) NewExportSpecifier(isTypeOnly bool, propertyName *ModuleExportName, name *ModuleExportName) *Node {
	data := &ExportSpecifier{}
	data.IsTypeOnly = isTypeOnly
	data.PropertyName = propertyName
	data.Name_ = name
	return f.NewNode(ast.KindExportSpecifier, data)
}

func (node *ExportSpecifier) ForEachChild(v Visitor) bool {
	return visit(v, node.PropertyName) || visit(v, node.Name_)
}

func (node *ExportSpecifier) Name() *DeclarationName {
	return node.Name_
}

func IsExportSpecifier(node *Node) bool {
	return node.Kind == ast.KindExportSpecifier
}

// TypeElementBase

type TypeElementBase struct{}

// ClassElementBase

type ClassElementBase struct{}

// NamedMemberBase

type NamedMemberBase struct {
	DeclarationBase
	ModifiersBase
	Name_         *PropertyName // PropertyName
	PostfixToken *TokenNode    // TokenNode. Optional
}

func (node *NamedMemberBase) DeclarationData() *DeclarationBase { return &node.DeclarationBase }
func (node *NamedMemberBase) Modifiers() *ModifierListNode      { return node.Modifiers_ }
func (node *NamedMemberBase) Name() *DeclarationName            { return node.Name_ }

// CallSignatureDeclaration

type CallSignatureDeclaration struct {
	NodeBase
	DeclarationBase
	FunctionLikeBase
	TypeElementBase
}

func (f *NodeFactory) NewCallSignatureDeclaration(typeParameters *TypeParameterListNode, parameters []*ParameterDeclarationNode, returnType *TypeNode) *Node {
	data := &CallSignatureDeclaration{}
	data.TypeParameters_ = typeParameters
	data.Parameters = parameters
	data.ReturnType = returnType
	return f.NewNode(ast.KindCallSignature, data)
}

func (node *CallSignatureDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.TypeParameters_) || visitNodes(v, node.Parameters) || visit(v, node.ReturnType)
}

func IsCallSignatureDeclaration(node *Node) bool {
	return node.Kind == ast.KindCallSignature
}

// ConstructSignatureDeclaration

type ConstructSignatureDeclaration struct {
	NodeBase
	DeclarationBase
	FunctionLikeBase
	TypeElementBase
}

func (f *NodeFactory) NewConstructSignatureDeclaration(typeParameters *TypeParameterListNode, parameters []*ParameterDeclarationNode, returnType *TypeNode) *Node {
	data := &ConstructSignatureDeclaration{}
	data.TypeParameters_ = typeParameters
	data.Parameters = parameters
	data.ReturnType = returnType
	return f.NewNode(ast.KindConstructSignature, data)
}

func (node *ConstructSignatureDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.TypeParameters_) || visitNodes(v, node.Parameters) || visit(v, node.ReturnType)
}

func IsConstructSignatureDeclaration(node *Node) bool {
	return node.Kind == ast.KindConstructSignature
}

// ConstructorDeclaration

type ConstructorDeclaration struct {
	NodeBase
	DeclarationBase
	ModifiersBase
	FunctionLikeWithBodyBase
	ClassElementBase
	ReturnFlowNode *FlowNode
}

func (f *NodeFactory) NewConstructorDeclaration(modifiers *ModifierListNode, typeParameters *TypeParameterListNode, parameters []*ParameterDeclarationNode, returnType *TypeNode, body *BlockNode) *Node {
	data := &ConstructorDeclaration{}
	data.Modifiers_ = modifiers
	data.TypeParameters_ = typeParameters
	data.Parameters = parameters
	data.ReturnType = returnType
	data.Body = body
	return f.NewNode(ast.KindConstructor, data)
}

func (node *ConstructorDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.Modifiers_) || visit(v, node.TypeParameters_) || visitNodes(v, node.Parameters) || visit(v, node.ReturnType) || visit(v, node.Body)
}

func IsConstructorDeclaration(node *Node) bool {
	return node.Kind == ast.KindConstructor
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
	return visit(v, node.Modifiers_) || visit(v, node.Name_) || visit(v, node.TypeParameters_) || visitNodes(v, node.Parameters) ||
		visit(v, node.ReturnType) || visit(v, node.Body)
}

func (node *AccessorDeclarationBase) IsAccessorDeclaration() {}

// GetAccessorDeclaration

type GetAccessorDeclaration struct {
	AccessorDeclarationBase
}

func (f *NodeFactory) NewGetAccessorDeclaration(modifiers *ModifierListNode, name *PropertyName, typeParameters *TypeParameterListNode, parameters []*ParameterDeclarationNode, returnType *TypeNode, body *BlockNode) *Node {
	data := &GetAccessorDeclaration{}
	data.Modifiers_ = modifiers
	data.Name_ = name
	data.TypeParameters_ = typeParameters
	data.Parameters = parameters
	data.ReturnType = returnType
	data.Body = body
	return f.NewNode(ast.KindGetAccessor, data)
}

func IsGetAccessorDeclaration(node *Node) bool {
	return node.Kind == ast.KindGetAccessor
}

// SetAccessorDeclaration

type SetAccessorDeclaration struct {
	AccessorDeclarationBase
}

func (f *NodeFactory) NewSetAccessorDeclaration(modifiers *ModifierListNode, name *PropertyName, typeParameters *TypeParameterListNode, parameters []*ParameterDeclarationNode, returnType *TypeNode, body *BlockNode) *Node {
	data := &SetAccessorDeclaration{}
	data.Modifiers_ = modifiers
	data.Name_ = name
	data.TypeParameters_ = typeParameters
	data.Parameters = parameters
	data.ReturnType = returnType
	data.Body = body
	return f.NewNode(ast.KindSetAccessor, data)
}

func IsSetAccessorDeclaration(node *Node) bool {
	return node.Kind == ast.KindSetAccessor
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

func (f *NodeFactory) NewIndexSignatureDeclaration(modifiers *ModifierListNode, parameters []*ParameterDeclarationNode, returnType *TypeNode) *Node {
	data := &IndexSignatureDeclaration{}
	data.Modifiers_ = modifiers
	data.Parameters = parameters
	data.ReturnType = returnType
	return f.NewNode(ast.KindIndexSignature, data)
}

func (node *IndexSignatureDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.Modifiers_) || visitNodes(v, node.Parameters) || visit(v, node.ReturnType)
}

func IsIndexSignatureDeclaration(node *Node) bool {
	return node.Kind == ast.KindIndexSignature
}

// MethodSignatureDeclaration

type MethodSignatureDeclaration struct {
	NodeBase
	NamedMemberBase
	FunctionLikeBase
	TypeElementBase
}

func (f *NodeFactory) NewMethodSignatureDeclaration(modifiers *ModifierListNode, name *PropertyName, postfixToken *TokenNode, typeParameters *TypeParameterListNode, parameters []*ParameterDeclarationNode, returnType *TypeNode) *Node {
	data := &MethodSignatureDeclaration{}
	data.Modifiers_ = modifiers
	data.Name_ = name
	data.PostfixToken = postfixToken
	data.TypeParameters_ = typeParameters
	data.Parameters = parameters
	data.ReturnType = returnType
	return f.NewNode(ast.KindMethodSignature, data)
}

func (node *MethodSignatureDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.Modifiers_) || visit(v, node.Name_) || visit(v, node.PostfixToken) || visit(v, node.TypeParameters_) ||
		visitNodes(v, node.Parameters) || visit(v, node.ReturnType)
}

func IsMethodSignatureDeclaration(node *Node) bool {
	return node.Kind == ast.KindMethodSignature
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

func (f *NodeFactory) NewMethodDeclaration(modifiers *ModifierListNode, asteriskToken *TokenNode, name *PropertyName, postfixToken *TokenNode, typeParameters *TypeParameterListNode, parameters []*ParameterDeclarationNode, returnType *TypeNode, body *BlockNode) *Node {
	data := &MethodDeclaration{}
	data.Modifiers_ = modifiers
	data.AsteriskToken = asteriskToken
	data.Name_ = name
	data.PostfixToken = postfixToken
	data.TypeParameters_ = typeParameters
	data.Parameters = parameters
	data.ReturnType = returnType
	data.Body = body
	return f.NewNode(ast.KindMethodDeclaration, data)
}

func (node *MethodDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.Modifiers_) || visit(v, node.AsteriskToken) || visit(v, node.Name_) || visit(v, node.PostfixToken) ||
		visit(v, node.TypeParameters_) || visitNodes(v, node.Parameters) || visit(v, node.ReturnType) || visit(v, node.Body)
}

func IsMethodDeclaration(node *Node) bool {
	return node.Kind == ast.KindMethodDeclaration
}

// PropertySignatureDeclaration

type PropertySignatureDeclaration struct {
	NodeBase
	NamedMemberBase
	TypeElementBase
	TypeNode    *TypeNode   // TypeNode
	Initializer *Expression // Expression. For error reporting purposes
}

func (f *NodeFactory) NewPropertySignatureDeclaration(modifiers *ModifierListNode, name *PropertyName, postfixToken *TokenNode, typeNode *TypeNode, initializer *Expression) *Node {
	data := &PropertySignatureDeclaration{}
	data.Modifiers_ = modifiers
	data.Name_ = name
	data.PostfixToken = postfixToken
	data.TypeNode = typeNode
	data.Initializer = initializer
	return f.NewNode(ast.KindPropertySignature, data)
}

func (node *PropertySignatureDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.Modifiers_) || visit(v, node.Name_) || visit(v, node.PostfixToken) || visit(v, node.TypeNode) || visit(v, node.Initializer)
}

func IsPropertySignatureDeclaration(node *Node) bool {
	return node.Kind == ast.KindPropertySignature
}

// PropertyDeclaration

type PropertyDeclaration struct {
	NodeBase
	NamedMemberBase
	ClassElementBase
	TypeNode    *TypeNode   // TypeNode. Optional
	Initializer *Expression // Expression. Optional
}

func (f *NodeFactory) NewPropertyDeclaration(modifiers *ModifierListNode, name *PropertyName, postfixToken *TokenNode, typeNode *TypeNode, initializer *Expression) *Node {
	data := &PropertyDeclaration{}
	data.Modifiers_ = modifiers
	data.Name_ = name
	data.PostfixToken = postfixToken
	data.TypeNode = typeNode
	data.Initializer = initializer
	return f.NewNode(ast.KindPropertyDeclaration, data)
}

func (node *PropertyDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.Modifiers_) || visit(v, node.Name_) || visit(v, node.PostfixToken) || visit(v, node.TypeNode) || visit(v, node.Initializer)
}

func IsPropertyDeclaration(node *Node) bool {
	return node.Kind == ast.KindPropertyDeclaration
}

// SemicolonClassElement

type SemicolonClassElement struct {
	NodeBase
	DeclarationBase
	ClassElementBase
}

func (f *NodeFactory) NewSemicolonClassElement() *Node {
	return f.NewNode(ast.KindSemicolonClassElement, &SemicolonClassElement{})
}

// ClassStaticBlockDeclaration

type ClassStaticBlockDeclaration struct {
	NodeBase
	DeclarationBase
	ModifiersBase
	LocalsContainerBase
	ClassElementBase
	Body *BlockNode // BlockNode
}

func (f *NodeFactory) NewClassStaticBlockDeclaration(modifiers *ModifierListNode, body *BlockNode) *Node {
	data := &ClassStaticBlockDeclaration{}
	data.Modifiers_ = modifiers
	data.Body = body
	return f.NewNode(ast.KindClassStaticBlockDeclaration, data)
}

func (node *ClassStaticBlockDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.Modifiers_) || visit(v, node.Body)
}

func IsClassStaticBlockDeclaration(node *Node) bool {
	return node.Kind == ast.KindClassStaticBlockDeclaration
}

// TypeParameterList

type TypeParameterList struct {
	NodeBase
	Parameters []*TypeParameterDeclarationNode // []TypeParameterDeclarationNode
}

func (f *NodeFactory) NewTypeParameterList(parameters []*TypeParameterDeclarationNode) *Node {
	data := &TypeParameterList{}
	data.Parameters = parameters
	return f.NewNode(ast.KindTypeParameterList, data)
}

func (node *TypeParameterList) ForEachChild(v Visitor) bool {
	return visitNodes(v, node.Parameters)
}

func IsTypeParameterList(node *Node) bool {
	return node.Kind == ast.KindTypeParameterList
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
	return f.NewNode(ast.KindOmittedExpression, &OmittedExpression{})
}

// KeywordExpression

type KeywordExpression struct {
	ExpressionBase
	FlowNodeBase // For 'this' and 'super' expressions
}

func (f *NodeFactory) NewKeywordExpression(kind ast.Kind) *Node {
	return f.NewNode(kind, &KeywordExpression{})
}

// LiteralLikeBase

type LiteralLikeBase struct {
	Text string
}

// StringLiteral

type StringLiteral struct {
	ExpressionBase
	LiteralLikeBase
}

func (f *NodeFactory) NewStringLiteral(text string) *Node {
	data := &StringLiteral{}
	data.Text = text
	return f.NewNode(ast.KindStringLiteral, data)
}

func IsStringLiteral(node *Node) bool {
	return node.Kind == ast.KindStringLiteral
}

// NumericLiteral

type NumericLiteral struct {
	ExpressionBase
	LiteralLikeBase
}

func (f *NodeFactory) NewNumericLiteral(text string) *Node {
	data := &NumericLiteral{}
	data.Text = text
	return f.NewNode(ast.KindNumericLiteral, data)
}

// BigIntLiteral

type BigIntLiteral struct {
	ExpressionBase
	LiteralLikeBase
}

func (f *NodeFactory) NewBigIntLiteral(text string) *Node {
	data := &BigIntLiteral{}
	data.Text = text
	return f.NewNode(ast.KindBigIntLiteral, data)
}

// RegularExpressionLiteral

type RegularExpressionLiteral struct {
	ExpressionBase
	LiteralLikeBase
}

func (f *NodeFactory) NewRegularExpressionLiteral(text string) *Node {
	data := &RegularExpressionLiteral{}
	data.Text = text
	return f.NewNode(ast.KindRegularExpressionLiteral, data)
}

// NoSubstitutionTemplateLiteral

type NoSubstitutionTemplateLiteral struct {
	ExpressionBase
	TemplateLiteralLikeBase
}

func (f *NodeFactory) NewNoSubstitutionTemplateLiteral(text string) *Node {
	data := &NoSubstitutionTemplateLiteral{}
	data.Text = text
	return f.NewNode(ast.KindNoSubstitutionTemplateLiteral, data)
}

// BinaryExpression

type BinaryExpression struct {
	ExpressionBase
	DeclarationBase
	Left          *Expression // Expression
	OperatorToken *TokenNode  // TokenNode
	Right         *Expression // Expression
}

func (f *NodeFactory) NewBinaryExpression(left *Expression, operatorToken *TokenNode, right *Expression) *Node {
	data := &BinaryExpression{}
	data.Left = left
	data.OperatorToken = operatorToken
	data.Right = right
	return f.NewNode(ast.KindBinaryExpression, data)
}

func (node *BinaryExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.Left) || visit(v, node.OperatorToken) || visit(v, node.Right)
}

// PrefixUnaryExpression

type PrefixUnaryExpression struct {
	ExpressionBase
	Operator ast.Kind
	Operand  *Expression // Expression
}

func (f *NodeFactory) NewPrefixUnaryExpression(operator ast.Kind, operand *Expression) *Node {
	data := &PrefixUnaryExpression{}
	data.Operator = operator
	data.Operand = operand
	return f.NewNode(ast.KindPrefixUnaryExpression, data)
}

func (node *PrefixUnaryExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.Operand)
}

func IsPrefixUnaryExpression(node *Node) bool {
	return node.Kind == ast.KindPrefixUnaryExpression
}

// PostfixUnaryExpression

type PostfixUnaryExpression struct {
	ExpressionBase
	Operand  *Expression // Expression
	Operator ast.Kind
}

func (f *NodeFactory) NewPostfixUnaryExpression(operand *Expression, operator ast.Kind) *Node {
	data := &PostfixUnaryExpression{}
	data.Operand = operand
	data.Operator = operator
	return f.NewNode(ast.KindPostfixUnaryExpression, data)
}

func (node *PostfixUnaryExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.Operand)
}

// YieldExpression

type YieldExpression struct {
	ExpressionBase
	AsteriskToken *TokenNode  // TokenNode
	Expression    *Expression // Expression. Optional
}

func (f *NodeFactory) NewYieldExpression(asteriskToken *TokenNode, expression *Expression) *Node {
	data := &YieldExpression{}
	data.AsteriskToken = asteriskToken
	data.Expression = expression
	return f.NewNode(ast.KindYieldExpression, data)
}

func (node *YieldExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.AsteriskToken) || visit(v, node.Expression)
}

// ArrowFunction

type ArrowFunction struct {
	ExpressionBase
	DeclarationBase
	ModifiersBase
	FunctionLikeWithBodyBase
	FlowNodeBase
	EqualsGreaterThanToken *TokenNode // TokenNode
}

func (f *NodeFactory) NewArrowFunction(modifiers *ModifierListNode, typeParameters *TypeParameterListNode, parameters []*ParameterDeclarationNode, returnType *TypeNode, equalsGreaterThanToken *TokenNode, body *BlockOrExpression) *Node {
	data := &ArrowFunction{}
	data.Modifiers_ = modifiers
	data.TypeParameters_ = typeParameters
	data.Parameters = parameters
	data.ReturnType = returnType
	data.EqualsGreaterThanToken = equalsGreaterThanToken
	data.Body = body
	return f.NewNode(ast.KindArrowFunction, data)
}

func (node *ArrowFunction) ForEachChild(v Visitor) bool {
	return visit(v, node.Modifiers_) || visit(v, node.TypeParameters_) || visitNodes(v, node.Parameters) ||
		visit(v, node.ReturnType) || visit(v, node.EqualsGreaterThanToken) || visit(v, node.Body)
}

func (node *ArrowFunction) Name() *DeclarationName {
	return nil
}

func IsArrowFunction(node *Node) bool {
	return node.Kind == ast.KindArrowFunction
}

// FunctionExpression

type FunctionExpression struct {
	ExpressionBase
	DeclarationBase
	ModifiersBase
	FunctionLikeWithBodyBase
	FlowNodeBase
	Name_           *IdentifierNode // IdentifierNode. Optional
	ReturnFlowNode *FlowNode
}

func (f *NodeFactory) NewFunctionExpression(modifiers *ModifierListNode, asteriskToken *TokenNode, name *IdentifierNode, typeParameters *TypeParameterListNode, parameters []*ParameterDeclarationNode, returnType *TypeNode, body *BlockNode) *Node {
	data := &FunctionExpression{}
	data.Modifiers_ = modifiers
	data.AsteriskToken = asteriskToken
	data.Name_ = name
	data.TypeParameters_ = typeParameters
	data.Parameters = parameters
	data.ReturnType = returnType
	data.Body = body
	return f.NewNode(ast.KindFunctionExpression, data)
}

func (node *FunctionExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.Modifiers_) || visit(v, node.AsteriskToken) || visit(v, node.Name_) || visit(v, node.TypeParameters_) ||
		visitNodes(v, node.Parameters) || visit(v, node.ReturnType) || visit(v, node.Body)
}

func (node *FunctionExpression) Name() *DeclarationName {
	return node.Name_
}

func IsFunctionExpression(node *Node) bool {
	return node.Kind == ast.KindFunctionExpression
}

// AsExpression

type AsExpression struct {
	ExpressionBase
	Expression *Expression // Expression
	TypeNode   *TypeNode   // TypeNode
}

func (f *NodeFactory) NewAsExpression(expression *Expression, typeNode *TypeNode) *Node {
	data := &AsExpression{}
	data.Expression = expression
	data.TypeNode = typeNode
	return f.NewNode(ast.KindAsExpression, data)
}

func (node *AsExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.Expression) || visit(v, node.TypeNode)
}

// SatisfiesExpression

type SatisfiesExpression struct {
	ExpressionBase
	Expression *Expression // Expression
	TypeNode   *TypeNode   // TypeNode
}

func (f *NodeFactory) NewSatisfiesExpression(expression *Expression, typeNode *TypeNode) *Node {
	data := &SatisfiesExpression{}
	data.Expression = expression
	data.TypeNode = typeNode
	return f.NewNode(ast.KindSatisfiesExpression, data)
}

func (node *SatisfiesExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.Expression) || visit(v, node.TypeNode)
}

// ConditionalExpression

type ConditionalExpression struct {
	ExpressionBase
	Condition     *Expression
	QuestionToken *TokenNode
	WhenTrue      *Expression
	ColonToken    *TokenNode
	WhenFalse     *Expression
}

func (f *NodeFactory) NewConditionalExpression(condition *Expression, questionToken *TokenNode, whenTrue *Expression, colonToken *TokenNode, whenFalse *Expression) *Node {
	data := &ConditionalExpression{}
	data.Condition = condition
	data.QuestionToken = questionToken
	data.WhenTrue = whenTrue
	data.ColonToken = colonToken
	data.WhenFalse = whenFalse
	return f.NewNode(ast.KindConditionalExpression, data)
}

func (node *ConditionalExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.Condition) || visit(v, node.QuestionToken) || visit(v, node.WhenTrue) ||
		visit(v, node.ColonToken) || visit(v, node.WhenFalse)
}

// PropertyAccessExpression

type PropertyAccessExpression struct {
	ExpressionBase
	FlowNodeBase
	Expression       *Expression // Expression
	QuestionDotToken *TokenNode  // TokenNode
	Name_             *MemberName // MemberName
}

func (f *NodeFactory) NewPropertyAccessExpression(expression *Expression, questionDotToken *TokenNode, name *MemberName, flags ast.NodeFlags) *Node {
	data := &PropertyAccessExpression{}
	data.Expression = expression
	data.QuestionDotToken = questionDotToken
	data.Name_ = name
	node := f.NewNode(ast.KindPropertyAccessExpression, data)
	node.Flags |= flags & ast.NodeFlagsOptionalChain
	return node
}

func (node *PropertyAccessExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.Expression) || visit(v, node.QuestionDotToken) || visit(v, node.Name_)
}

func (node *PropertyAccessExpression) Name() *DeclarationName { return node.Name_ }

func IsPropertyAccessExpression(node *Node) bool {
	return node.Kind == ast.KindPropertyAccessExpression
}

// ElementAccessExpression

type ElementAccessExpression struct {
	ExpressionBase
	FlowNodeBase
	Expression         *Expression // Expression
	QuestionDotToken   *TokenNode  // TokenNode
	ArgumentExpression *Expression // Expression
}

func (f *NodeFactory) NewElementAccessExpression(expression *Expression, questionDotToken *TokenNode, argumentExpression *Expression, flags ast.NodeFlags) *Node {
	data := &ElementAccessExpression{}
	data.Expression = expression
	data.QuestionDotToken = questionDotToken
	data.ArgumentExpression = argumentExpression
	node := f.NewNode(ast.KindElementAccessExpression, data)
	node.Flags |= flags & ast.NodeFlagsOptionalChain
	return node
}

func (node *ElementAccessExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.Expression) || visit(v, node.QuestionDotToken) || visit(v, node.ArgumentExpression)
}

func IsElementAccessExpression(node *Node) bool {
	return node.Kind == ast.KindElementAccessExpression
}

// CallExpression

type CallExpression struct {
	ExpressionBase
	Expression       *Expression           // Expression
	QuestionDotToken *TokenNode            // TokenNode
	TypeArguments    *TypeArgumentListNode // TypeArgumentListNode
	Arguments        []*Expression         // []Expression
}

func (f *NodeFactory) NewCallExpression(expression *Expression, questionDotToken *TokenNode, typeArguments *TypeArgumentListNode, arguments []*Expression, flags ast.NodeFlags) *Node {
	data := &CallExpression{}
	data.Expression = expression
	data.QuestionDotToken = questionDotToken
	data.TypeArguments = typeArguments
	data.Arguments = arguments
	node := f.NewNode(ast.KindCallExpression, data)
	node.Flags |= flags & ast.NodeFlagsOptionalChain
	return node
}

func (node *CallExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.Expression) || visit(v, node.QuestionDotToken) || visit(v, node.TypeArguments) || visitNodes(v, node.Arguments)
}

func IsCallExpression(node *Node) bool {
	return node.Kind == ast.KindCallExpression
}

// NewExpression

type NewExpression struct {
	ExpressionBase
	Expression    *Expression           // Expression
	TypeArguments *TypeArgumentListNode // TypeArgumentListNode
	Arguments     []*Expression         // []Expression
}

func (f *NodeFactory) NewNewExpression(expression *Expression, typeArguments *TypeArgumentListNode, arguments []*Expression) *Node {
	data := &NewExpression{}
	data.Expression = expression
	data.TypeArguments = typeArguments
	data.Arguments = arguments
	return f.NewNode(ast.KindNewExpression, data)
}

func (node *NewExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.Expression) || visit(v, node.TypeArguments) || visitNodes(v, node.Arguments)
}

func IsNewExpression(node *Node) bool {
	return node.Kind == ast.KindNewExpression
}

// MetaProperty

type MetaProperty struct {
	ExpressionBase
	FlowNodeBase
	KeywordToken ast.Kind      // NewKeyword | ImportKeyword
	Name_         *IdentifierNode // IdentifierNode
}

func (f *NodeFactory) NewMetaProperty(keywordToken ast.Kind, name *IdentifierNode) *Node {
	data := &MetaProperty{}
	data.KeywordToken = keywordToken
	data.Name_ = name
	return f.NewNode(ast.KindMetaProperty, data)
}

func (node *MetaProperty) ForEachChild(v Visitor) bool {
	return visit(v, node.Name_)
}

func IsMetaProperty(node *Node) bool {
	return node.Kind == ast.KindMetaProperty
}

// NonNullExpression

type NonNullExpression struct {
	ExpressionBase
	Expression *Expression // Expression
}

func (f *NodeFactory) NewNonNullExpression(expression *Expression) *Node {
	data := &NonNullExpression{}
	data.Expression = expression
	return f.NewNode(ast.KindNonNullExpression, data)
}

func (node *NonNullExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.Expression)
}

// SpreadElement

type SpreadElement struct {
	ExpressionBase
	Expression *Expression // Expression
}

func (f *NodeFactory) NewSpreadElement(expression *Expression) *Node {
	data := &SpreadElement{}
	data.Expression = expression
	return f.NewNode(ast.KindSpreadElement, data)
}

func (node *SpreadElement) ForEachChild(v Visitor) bool {
	return visit(v, node.Expression)
}

func IsSpreadElement(node *Node) bool {
	return node.Kind == ast.KindSpreadElement
}

// TemplateExpression

type TemplateExpression struct {
	ExpressionBase
	Head          *TemplateHeadNode   // TemplateHeadNode
	TemplateSpans []*TemplateSpanNode // []TemplateSpanNode
}

func (f *NodeFactory) NewTemplateExpression(head *TemplateHeadNode, templateSpans []*TemplateSpanNode) *Node {
	data := &TemplateExpression{}
	data.Head = head
	data.TemplateSpans = templateSpans
	return f.NewNode(ast.KindTemplateExpression, data)
}

func (node *TemplateExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.Head) || visitNodes(v, node.TemplateSpans)
}

// TemplateLiteralTypeSpan

type TemplateSpan struct {
	NodeBase
	Expression *Expression           // Expression
	Literal    *TemplateMiddleOrTail // TemplateMiddleOrTail
}

func (f *NodeFactory) NewTemplateSpan(expression *Expression, literal *TemplateMiddleOrTail) *Node {
	data := &TemplateSpan{}
	data.Expression = expression
	data.Literal = literal
	return f.NewNode(ast.KindTemplateSpan, data)
}

func (node *TemplateSpan) ForEachChild(v Visitor) bool {
	return visit(v, node.Expression) || visit(v, node.Literal)
}

func IsTemplateSpan(node *Node) bool {
	return node.Kind == ast.KindTemplateSpan
}

// TaggedTemplateExpression

type TaggedTemplateExpression struct {
	ExpressionBase
	Tag              *Expression           // Expression
	QuestionDotToken *TokenNode            // TokenNode. For error reporting purposes only
	TypeArguments    *TypeArgumentListNode // TypeArgumentListNode
	Template         *TemplateLiteral      // TemplateLiteral
}

func (f *NodeFactory) NewTaggedTemplateExpression(tag *Expression, questionDotToken *TokenNode, typeArguments *TypeArgumentListNode, template *TemplateLiteral, flags ast.NodeFlags) *Node {
	data := &TaggedTemplateExpression{}
	data.Tag = tag
	data.QuestionDotToken = questionDotToken
	data.TypeArguments = typeArguments
	data.Template = template
	node := f.NewNode(ast.KindTaggedTemplateExpression, data)
	node.Flags |= flags & ast.NodeFlagsOptionalChain
	return node
}

func (node *TaggedTemplateExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.Tag) || visit(v, node.QuestionDotToken) || visit(v, node.TypeArguments) || visit(v, node.Template)
}

// ParenthesizedExpression

type ParenthesizedExpression struct {
	ExpressionBase
	Expression *Expression // Expression
}

func (f *NodeFactory) NewParenthesizedExpression(expression *Expression) *Node {
	data := &ParenthesizedExpression{}
	data.Expression = expression
	return f.NewNode(ast.KindParenthesizedExpression, data)
}

func (node *ParenthesizedExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.Expression)
}

func IsParenthesizedExpression(node *Node) bool {
	return node.Kind == ast.KindParenthesizedExpression
}

// ArrayLiteralExpression

type ArrayLiteralExpression struct {
	ExpressionBase
	Elements  []*Expression // []Expression
	MultiLine bool
}

func (f *NodeFactory) NewArrayLiteralExpression(elements []*Expression, multiLine bool) *Node {
	data := &ArrayLiteralExpression{}
	data.Elements = elements
	data.MultiLine = multiLine
	return f.NewNode(ast.KindArrayLiteralExpression, data)
}

func (node *ArrayLiteralExpression) ForEachChild(v Visitor) bool {
	return visitNodes(v, node.Elements)
}

func IsArrayLiteralExpression(node *Node) bool {
	return node.Kind == ast.KindArrayLiteralExpression
}

// ObjectLiteralExpression

type ObjectLiteralExpression struct {
	ExpressionBase
	DeclarationBase
	Properties []*ObjectLiteralElement // []ObjectLiteralElement
	MultiLine  bool
}

func (f *NodeFactory) NewObjectLiteralExpression(properties []*ObjectLiteralElement, multiLine bool) *Node {
	data := &ObjectLiteralExpression{}
	data.Properties = properties
	data.MultiLine = multiLine
	return f.NewNode(ast.KindObjectLiteralExpression, data)

}

func (node *ObjectLiteralExpression) ForEachChild(v Visitor) bool {
	return visitNodes(v, node.Properties)
}

func IsObjectLiteralExpression(node *Node) bool {
	return node.Kind == ast.KindObjectLiteralExpression
}

// ObjectLiteralElementBase

type ObjectLiteralElementBase struct{}

// SpreadAssignment

type SpreadAssignment struct {
	NodeBase
	DeclarationBase
	ObjectLiteralElementBase
	Expression *Expression // Expression
}

func (f *NodeFactory) NewSpreadAssignment(expression *Expression) *Node {
	data := &SpreadAssignment{}
	data.Expression = expression
	return f.NewNode(ast.KindSpreadAssignment, data)
}

func (node *SpreadAssignment) ForEachChild(v Visitor) bool {
	return visit(v, node.Expression)
}

// PropertyAssignment

type PropertyAssignment struct {
	NodeBase
	NamedMemberBase
	ObjectLiteralElementBase
	Initializer *Expression // Expression
}

func (f *NodeFactory) NewPropertyAssignment(modifiers *ModifierListNode, name *PropertyName, postfixToken *TokenNode, initializer *Expression) *Node {
	data := &PropertyAssignment{}
	data.Modifiers_ = modifiers
	data.Name_ = name
	data.PostfixToken = postfixToken
	data.Initializer = initializer
	return f.NewNode(ast.KindPropertyAssignment, data)
}

func (node *PropertyAssignment) ForEachChild(v Visitor) bool {
	return visit(v, node.Modifiers_) || visit(v, node.Name_) || visit(v, node.PostfixToken) || visit(v, node.Initializer)
}

func IsPropertyAssignment(node *Node) bool {
	return node.Kind == ast.KindPropertyAssignment
}

// ShorthandPropertyAssignment

type ShorthandPropertyAssignment struct {
	NodeBase
	NamedMemberBase
	ObjectLiteralElementBase
	ObjectAssignmentInitializer *Expression // Optional
}

func (f *NodeFactory) NewShorthandPropertyAssignment(modifiers *ModifierListNode, name *PropertyName, postfixToken *TokenNode, objectAssignmentInitializer *Expression) *Node {
	data := &ShorthandPropertyAssignment{}
	data.Modifiers_ = modifiers
	data.Name_ = name
	data.PostfixToken = postfixToken
	data.ObjectAssignmentInitializer = objectAssignmentInitializer
	return f.NewNode(ast.KindShorthandPropertyAssignment, data)
}

func (node *ShorthandPropertyAssignment) ForEachChild(v Visitor) bool {
	return visit(v, node.Modifiers_) || visit(v, node.Name_) || visit(v, node.PostfixToken) || visit(v, node.ObjectAssignmentInitializer)
}

func IsShorthandPropertyAssignment(node *Node) bool {
	return node.Kind == ast.KindShorthandPropertyAssignment
}

// DeleteExpression

type DeleteExpression struct {
	ExpressionBase
	Expression *Expression // Expression
}

func (f *NodeFactory) NewDeleteExpression(expression *Expression) *Node {
	data := &DeleteExpression{}
	data.Expression = expression
	return f.NewNode(ast.KindDeleteExpression, data)

}

func (node *DeleteExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.Expression)
}

// TypeOfExpression

type TypeOfExpression struct {
	ExpressionBase
	Expression *Expression // Expression
}

func (f *NodeFactory) NewTypeOfExpression(expression *Expression) *Node {
	data := &TypeOfExpression{}
	data.Expression = expression
	return f.NewNode(ast.KindTypeOfExpression, data)
}

func (node *TypeOfExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.Expression)
}

func IsTypeOfExpression(node *Node) bool {
	return node.Kind == ast.KindTypeOfExpression
}

// VoidExpression

type VoidExpression struct {
	ExpressionBase
	Expression *Expression // Expression
}

func (f *NodeFactory) NewVoidExpression(expression *Expression) *Node {
	data := &VoidExpression{}
	data.Expression = expression
	return f.NewNode(ast.KindVoidExpression, data)
}

func (node *VoidExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.Expression)
}

// AwaitExpression

type AwaitExpression struct {
	ExpressionBase
	Expression *Expression // Expression
}

func (f *NodeFactory) NewAwaitExpression(expression *Expression) *Node {
	data := &AwaitExpression{}
	data.Expression = expression
	return f.NewNode(ast.KindAwaitExpression, data)
}

func (node *AwaitExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.Expression)
}

// TypeAssertion

type TypeAssertion struct {
	ExpressionBase
	TypeNode   *TypeNode   // TypeNode
	Expression *Expression // Expression
}

func (f *NodeFactory) NewTypeAssertion(typeNode *TypeNode, expression *Expression) *Node {
	data := &TypeAssertion{}
	data.TypeNode = typeNode
	data.Expression = expression
	return f.NewNode(ast.KindTypeAssertionExpression, data)
}

func (node *TypeAssertion) ForEachChild(v Visitor) bool {
	return visit(v, node.TypeNode) || visit(v, node.Expression)
}

// TypeNodeBase

type TypeNodeBase struct {
	NodeBase
}

// KeywordTypeNode

type KeywordTypeNode struct {
	TypeNodeBase
}

func (f *NodeFactory) NewKeywordTypeNode(kind ast.Kind) *Node {
	return f.NewNode(kind, &KeywordTypeNode{})
}

// UnionOrIntersectionTypeBase

type UnionOrIntersectionTypeNodeBase struct {
	TypeNodeBase
	Types []*TypeNode // []TypeNode
}

func (node *UnionOrIntersectionTypeNodeBase) ForEachChild(v Visitor) bool {
	return visitNodes(v, node.Types)
}

// UnionTypeNode

type UnionTypeNode struct {
	UnionOrIntersectionTypeNodeBase
}

func (f *NodeFactory) NewUnionTypeNode(types []*TypeNode) *Node {
	data := &UnionTypeNode{}
	data.Types = types
	return f.NewNode(ast.KindUnionType, data)
}

// IntersectionTypeNode

type IntersectionTypeNode struct {
	UnionOrIntersectionTypeNodeBase
}

func (f *NodeFactory) NewIntersectionTypeNode(types []*TypeNode) *Node {
	data := &IntersectionTypeNode{}
	data.Types = types
	return f.NewNode(ast.KindIntersectionType, data)
}

// ConditionalTypeNode

type ConditionalTypeNode struct {
	TypeNodeBase
	LocalsContainerBase
	CheckType   *TypeNode // TypeNode
	ExtendsType *TypeNode // TypeNode
	TrueType    *TypeNode // TypeNode
	FalseType   *TypeNode // TypeNode
}

func (f *NodeFactory) NewConditionalTypeNode(checkType *TypeNode, extendsType *TypeNode, trueType *TypeNode, falseType *TypeNode) *Node {
	data := &ConditionalTypeNode{}
	data.CheckType = checkType
	data.ExtendsType = extendsType
	data.TrueType = trueType
	data.FalseType = falseType
	return f.NewNode(ast.KindConditionalType, data)
}

func (node *ConditionalTypeNode) ForEachChild(v Visitor) bool {
	return visit(v, node.CheckType) || visit(v, node.ExtendsType) || visit(v, node.TrueType) || visit(v, node.FalseType)
}

func IsConditionalTypeNode(node *Node) bool {
	return node.Kind == ast.KindConditionalType
}

// TypeOperatorNode

type TypeOperatorNode struct {
	TypeNodeBase
	Operator ast.Kind // ast.KindKeyOfKeyword | ast.KindUniqueKeyword | ast.KindReadonlyKeyword
	TypeNode *TypeNode  // TypeNode
}

func (f *NodeFactory) NewTypeOperatorNode(operator ast.Kind, typeNode *TypeNode) *Node {
	data := &TypeOperatorNode{}
	data.Operator = operator
	data.TypeNode = typeNode
	return f.NewNode(ast.KindTypeOperator, data)
}

func (node *TypeOperatorNode) ForEachChild(v Visitor) bool {
	return visit(v, node.TypeNode)
}

func IsTypeOperatorNode(node *Node) bool {
	return node.Kind == ast.KindTypeOperator
}

// InferTypeNode

type InferTypeNode struct {
	TypeNodeBase
	TypeParameter *TypeParameterDeclarationNode // TypeParameterDeclarationNode
}

func (f *NodeFactory) NewInferTypeNode(typeParameter *TypeParameterDeclarationNode) *Node {
	data := &InferTypeNode{}
	data.TypeParameter = typeParameter
	return f.NewNode(ast.KindInferType, data)
}

func (node *InferTypeNode) ForEachChild(v Visitor) bool {
	return visit(v, node.TypeParameter)
}

// ArrayTypeNode

type ArrayTypeNode struct {
	TypeNodeBase
	ElementType *TypeNode // TypeNode
}

func (f *NodeFactory) NewArrayTypeNode(elementType *TypeNode) *Node {
	data := &ArrayTypeNode{}
	data.ElementType = elementType
	return f.NewNode(ast.KindArrayType, data)
}

func (node *ArrayTypeNode) ForEachChild(v Visitor) bool {
	return visit(v, node.ElementType)
}

// IndexedAccessTypeNode

type IndexedAccessTypeNode struct {
	TypeNodeBase
	ObjectType *TypeNode // TypeNode
	IndexType  *TypeNode // TypeNode
}

func (f *NodeFactory) NewIndexedAccessTypeNode(objectType *TypeNode, indexType *TypeNode) *Node {
	data := &IndexedAccessTypeNode{}
	data.ObjectType = objectType
	data.IndexType = indexType
	return f.NewNode(ast.KindIndexedAccessType, data)
}

func (node *IndexedAccessTypeNode) ForEachChild(v Visitor) bool {
	return visit(v, node.ObjectType) || visit(v, node.IndexType)
}

func IsIndexedAccessTypeNode(node *Node) bool {
	return node.Kind == ast.KindIndexedAccessType
}

// TypeArgumentList

type TypeArgumentList struct {
	NodeBase
	Arguments []*TypeNode // []TypeNode
}

func (f *NodeFactory) NewTypeArgumentList(arguments []*TypeNode) *Node {
	data := &TypeArgumentList{}
	data.Arguments = arguments
	return f.NewNode(ast.KindTypeArgumentList, data)
}

func (node *TypeArgumentList) ForEachChild(v Visitor) bool {
	return visitNodes(v, node.Arguments)
}

// TypeReferenceNode

type TypeReferenceNode struct {
	TypeNodeBase
	TypeName      *EntityName           // EntityName
	TypeArguments *TypeArgumentListNode // TypeArgumentListNode
}

func (f *NodeFactory) NewTypeReferenceNode(typeName *EntityName, typeArguments *TypeArgumentListNode) *Node {
	data := &TypeReferenceNode{}
	data.TypeName = typeName
	data.TypeArguments = typeArguments
	return f.NewNode(ast.KindTypeReference, data)
}

func (node *TypeReferenceNode) ForEachChild(v Visitor) bool {
	return visit(v, node.TypeName) || visit(v, node.TypeArguments)
}

func IsTypeReferenceNode(node *Node) bool {
	return node.Kind == ast.KindTypeReference
}

// ExpressionWithTypeArguments

type ExpressionWithTypeArguments struct {
	ExpressionBase
	Expression    *Expression           // Expression
	TypeArguments *TypeArgumentListNode // TypeArgumentListNode. Optional
}

func (f *NodeFactory) NewExpressionWithTypeArguments(expression *Expression, typeArguments *TypeArgumentListNode) *Node {
	data := &ExpressionWithTypeArguments{}
	data.Expression = expression
	data.TypeArguments = typeArguments
	return f.NewNode(ast.KindExpressionWithTypeArguments, data)
}

func (node *ExpressionWithTypeArguments) ForEachChild(v Visitor) bool {
	return visit(v, node.Expression) || visit(v, node.TypeArguments)
}

// LiteralTypeNode

type LiteralTypeNode struct {
	TypeNodeBase
	Literal *Node // KeywordExpression | LiteralExpression | PrefixUnaryExpression
}

func (f *NodeFactory) NewLiteralTypeNode(literal *Node) *Node {
	data := &LiteralTypeNode{}
	data.Literal = literal
	return f.NewNode(ast.KindLiteralType, data)
}

func (node *LiteralTypeNode) ForEachChild(v Visitor) bool {
	return visit(v, node.Literal)
}

func IsLiteralTypeNode(node *Node) bool {
	return node.Kind == ast.KindLiteralType
}

// ThisTypeNode

type ThisTypeNode struct {
	TypeNodeBase
}

func (f *NodeFactory) NewThisTypeNode() *Node {
	return f.NewNode(ast.KindThisType, &ThisTypeNode{})
}

func IsThisTypeNode(node *Node) bool {
	return node.Kind == ast.KindThisType
}

// TypePredicateNode

type TypePredicateNode struct {
	TypeNodeBase
	AssertsModifier *TokenNode                  // TokenNode. Optional
	ParameterName   *TypePredicateParameterName // TypePredicateParameterName (Identifier | ThisTypeNode)
	TypeNode        *TypeNode                   // TypeNode. Optional
}

func (f *NodeFactory) NewTypePredicateNode(assertsModifier *TokenNode, parameterName *TypePredicateParameterName, typeNode *TypeNode) *Node {
	data := &TypePredicateNode{}
	data.AssertsModifier = assertsModifier
	data.ParameterName = parameterName
	data.TypeNode = typeNode
	return f.NewNode(ast.KindTypePredicate, data)
}

func (node *TypePredicateNode) ForEachChild(v Visitor) bool {
	return visit(v, node.AssertsModifier) || visit(v, node.ParameterName) || visit(v, node.TypeNode)
}

func IsTypePredicateNode(node *Node) bool {
	return node.Kind == ast.KindTypePredicate
}

// ImportTypeNode

type ImportTypeNode struct {
	TypeNodeBase
	IsTypeOf      bool
	Argument      *TypeNode             // TypeNode
	Attributes    *ImportAttributesNode // ImportAttributesNode. Optional
	Qualifier     *EntityName           // EntityName. Optional
	TypeArguments *TypeArgumentListNode // TypeArgumentListNode. Optional
}

func (f *NodeFactory) NewImportTypeNode(isTypeOf bool, argument *TypeNode, attributes *ImportAttributesNode, qualifier *EntityName, typeArguments *TypeArgumentListNode) *Node {
	data := &ImportTypeNode{}
	data.IsTypeOf = isTypeOf
	data.Argument = argument
	data.Attributes = attributes
	data.Qualifier = qualifier
	data.TypeArguments = typeArguments
	return f.NewNode(ast.KindImportType, data)
}

func (node *ImportTypeNode) ForEachChild(v Visitor) bool {
	return visit(v, node.Argument) || visit(v, node.Attributes) || visit(v, node.Qualifier) || visit(v, node.TypeArguments)
}

func IsImportTypeNode(node *Node) bool {
	return node.Kind == ast.KindImportType
}

// ImportAttribute

type ImportAttribute struct {
	NodeBase
	Name_  *ImportAttributeName // ImportAttributeName
	Value *Expression          // Expression
}

func (f *NodeFactory) NewImportAttribute(name *ImportAttributeName, value *Expression) *Node {
	data := &ImportAttribute{}
	data.Name_ = name
	data.Value = value
	return f.NewNode(ast.KindImportAttribute, data)
}

func (node *ImportAttribute) ForEachChild(v Visitor) bool {
	return visit(v, node.Name_) || visit(v, node.Value)
}

// ImportAttributes

type ImportAttributes struct {
	NodeBase
	Token      ast.Kind
	Attributes []*ImportAttributeNode // []ImportAttributeNode
	MultiLine  bool
}

func (f *NodeFactory) NewImportAttributes(token ast.Kind, attributes []*ImportAttributeNode, multiLine bool) *Node {
	data := &ImportAttributes{}
	data.Token = token
	data.Attributes = attributes
	data.MultiLine = multiLine
	return f.NewNode(ast.KindImportAttributes, data)
}

func (node *ImportAttributes) ForEachChild(v Visitor) bool {
	return visitNodes(v, node.Attributes)
}

// TypeQueryNode

type TypeQueryNode struct {
	TypeNodeBase
	ExprName      *EntityName           // EntityName
	TypeArguments *TypeArgumentListNode // TypeArgumentListNode
}

func (f *NodeFactory) NewTypeQueryNode(exprName *EntityName, typeArguments *TypeArgumentListNode) *Node {
	data := &TypeQueryNode{}
	data.ExprName = exprName
	data.TypeArguments = typeArguments
	return f.NewNode(ast.KindTypeQuery, data)
}

func (node *TypeQueryNode) ForEachChild(v Visitor) bool {
	return visit(v, node.ExprName) || visit(v, node.TypeArguments)
}

func IsTypeQueryNode(node *Node) bool {
	return node.Kind == ast.KindTypeQuery
}

// MappedTypeNode

type MappedTypeNode struct {
	TypeNodeBase
	DeclarationBase
	LocalsContainerBase
	ReadonlyToken *TokenNode                    // TokenNode. Optional
	TypeParameter *TypeParameterDeclarationNode // TypeParameterDeclarationNode
	NameType      *TypeNode                     // TypeNode. Optional
	QuestionToken *TokenNode                    // TokenNode. Optional
	TypeNode      *TypeNode                     // TypeNode. Optional (error if missing)
	Members       []*TypeElement                // []TypeElement. Used only to produce grammar errors
}

func (f *NodeFactory) NewMappedTypeNode(readonlyToken *TokenNode, typeParameter *TypeParameterDeclarationNode, nameType *TypeNode, questionToken *TokenNode, typeNode *TypeNode, members []*TypeElement) *Node {
	data := &MappedTypeNode{}
	data.ReadonlyToken = readonlyToken
	data.TypeParameter = typeParameter
	data.NameType = nameType
	data.QuestionToken = questionToken
	data.TypeNode = typeNode
	data.Members = members
	return f.NewNode(ast.KindMappedType, data)
}

func (node *MappedTypeNode) ForEachChild(v Visitor) bool {
	return visit(v, node.ReadonlyToken) || visit(v, node.TypeParameter) || visit(v, node.NameType) ||
		visit(v, node.QuestionToken) || visit(v, node.TypeNode) || visitNodes(v, node.Members)
}

func IsMappedTypeNode(node *Node) bool {
	return node.Kind == ast.KindMappedType
}

// TypeLiteralNode

type TypeLiteralNode struct {
	TypeNodeBase
	DeclarationBase
	Members []*TypeElement // []TypeElement
}

func (f *NodeFactory) NewTypeLiteralNode(members []*TypeElement) *Node {
	data := &TypeLiteralNode{}
	data.Members = members
	return f.NewNode(ast.KindTypeLiteral, data)
}

func (node *TypeLiteralNode) ForEachChild(v Visitor) bool {
	return visitNodes(v, node.Members)
}

// TupleTypeNode

type TupleTypeNode struct {
	TypeNodeBase
	Elements []*TypeNode // []TypeNode
}

func (f *NodeFactory) NewTupleTypeNode(elements []*TypeNode) *Node {
	data := &TupleTypeNode{}
	data.Elements = elements
	return f.NewNode(ast.KindTupleType, data)
}

func (node *TupleTypeNode) Kind() ast.Kind {
	return ast.KindTupleType
}

func (node *TupleTypeNode) ForEachChild(v Visitor) bool {
	return visitNodes(v, node.Elements)
}

// NamedTupleTypeMember

type NamedTupleMember struct {
	TypeNodeBase
	DeclarationBase
	DotDotDotToken *TokenNode      // TokenNode
	Name_           *IdentifierNode // IdentifierNode
	QuestionToken  *TokenNode      // TokenNode
	TypeNode       *TypeNode       // TypeNode
}

func (f *NodeFactory) NewNamedTupleTypeMember(dotDotDotToken *TokenNode, name *IdentifierNode, questionToken *TokenNode, typeNode *TypeNode) *Node {
	data := &NamedTupleMember{}
	data.DotDotDotToken = dotDotDotToken
	data.Name_ = name
	data.QuestionToken = questionToken
	data.TypeNode = typeNode
	return f.NewNode(ast.KindNamedTupleMember, data)
}

func (node *NamedTupleMember) ForEachChild(v Visitor) bool {
	return visit(v, node.DotDotDotToken) || visit(v, node.Name_) || visit(v, node.QuestionToken) || visit(v, node.TypeNode)
}

func IsNamedTupleMember(node *Node) bool {
	return node.Kind == ast.KindNamedTupleMember
}

// OptionalTypeNode

type OptionalTypeNode struct {
	TypeNodeBase
	TypeNode *TypeNode // TypeNode
}

func (f *NodeFactory) NewOptionalTypeNode(typeNode *TypeNode) *Node {
	data := &OptionalTypeNode{}
	data.TypeNode = typeNode
	return f.NewNode(ast.KindOptionalType, data)
}

func (node *OptionalTypeNode) ForEachChild(v Visitor) bool {
	return visit(v, node.TypeNode)
}

// RestTypeNode

type RestTypeNode struct {
	TypeNodeBase
	TypeNode *TypeNode // TypeNode
}

func (f *NodeFactory) NewRestTypeNode(typeNode *TypeNode) *Node {
	data := &RestTypeNode{}
	data.TypeNode = typeNode
	return f.NewNode(ast.KindRestType, data)
}

func (node *RestTypeNode) ForEachChild(v Visitor) bool {
	return visit(v, node.TypeNode)
}

// ParenthesizedTypeNode

type ParenthesizedTypeNode struct {
	TypeNodeBase
	TypeNode *TypeNode // TypeNode
}

func (f *NodeFactory) NewParenthesizedTypeNode(typeNode *TypeNode) *Node {
	data := &ParenthesizedTypeNode{}
	data.TypeNode = typeNode
	return f.NewNode(ast.KindParenthesizedType, data)
}

func (node *ParenthesizedTypeNode) ForEachChild(v Visitor) bool {
	return visit(v, node.TypeNode)
}

func IsParenthesizedTypeNode(node *Node) bool {
	return node.Kind == ast.KindParenthesizedType
}

// FunctionOrConstructorTypeNodeBase

type FunctionOrConstructorTypeNodeBase struct {
	TypeNodeBase
	DeclarationBase
	ModifiersBase
	FunctionLikeBase
}

func (node *FunctionOrConstructorTypeNodeBase) ForEachChild(v Visitor) bool {
	return visit(v, node.Modifiers_) || visit(v, node.TypeParameters_) || visitNodes(v, node.Parameters) || visit(v, node.ReturnType)
}

// FunctionTypeNode

type FunctionTypeNode struct {
	FunctionOrConstructorTypeNodeBase
}

func (f *NodeFactory) NewFunctionTypeNode(typeParameters *TypeParameterListNode, parameters []*ParameterDeclarationNode, returnType *TypeNode) *Node {
	data := &FunctionTypeNode{}
	data.TypeParameters_ = typeParameters
	data.Parameters = parameters
	data.ReturnType = returnType
	return f.NewNode(ast.KindFunctionType, data)
}

func IsFunctionTypeNode(node *Node) bool {
	return node.Kind == ast.KindFunctionType
}

// ConstructorTypeNode

type ConstructorTypeNode struct {
	FunctionOrConstructorTypeNodeBase
}

func (f *NodeFactory) NewConstructorTypeNode(modifiers *ModifierListNode, typeParameters *TypeParameterListNode, parameters []*ParameterDeclarationNode, returnType *TypeNode) *Node {
	data := &ConstructorTypeNode{}
	data.Modifiers_ = modifiers
	data.TypeParameters_ = typeParameters
	data.Parameters = parameters
	data.ReturnType = returnType
	return f.NewNode(ast.KindConstructorType, data)
}

func IsConstructorTypeNode(node *Node) bool {
	return node.Kind == ast.KindConstructorType
}

// TemplateLiteralLikeBase

type TemplateLiteralLikeBase struct {
	LiteralLikeBase
	RawText       string
	TemplateFlags TokenFlags
}

// TemplateHead

type TemplateHead struct {
	NodeBase
	TemplateLiteralLikeBase
}

func (f *NodeFactory) NewTemplateHead(text string, rawText string, templateFlags TokenFlags) *Node {
	data := &TemplateHead{}
	data.Text = text
	data.RawText = rawText
	data.TemplateFlags = templateFlags
	return f.NewNode(ast.KindTemplateHead, data)
}

// TemplateMiddle

type TemplateMiddle struct {
	NodeBase
	TemplateLiteralLikeBase
}

func (f *NodeFactory) NewTemplateMiddle(text string, rawText string, templateFlags TokenFlags) *Node {
	data := &TemplateMiddle{}
	data.Text = text
	data.RawText = rawText
	data.TemplateFlags = templateFlags
	return f.NewNode(ast.KindTemplateMiddle, data)
}

// TemplateTail

type TemplateTail struct {
	NodeBase
	TemplateLiteralLikeBase
}

func (f *NodeFactory) NewTemplateTail(text string, rawText string, templateFlags TokenFlags) *Node {
	data := &TemplateTail{}
	data.Text = text
	data.RawText = rawText
	data.TemplateFlags = templateFlags
	return f.NewNode(ast.KindTemplateTail, data)
}

// TemplateLiteralTypeNode

type TemplateLiteralTypeNode struct {
	TypeNodeBase
	Head          *TemplateHeadNode              // TemplateHeadNode
	TemplateSpans []*TemplateLiteralTypeSpanNode // []TemplateLiteralTypeSpanNode
}

func (f *NodeFactory) NewTemplateLiteralTypeNode(head *TemplateHeadNode, templateSpans []*TemplateLiteralTypeSpanNode) *Node {
	data := &TemplateLiteralTypeNode{}
	data.Head = head
	data.TemplateSpans = templateSpans
	return f.NewNode(ast.KindTemplateLiteralType, data)
}

func (node *TemplateLiteralTypeNode) ForEachChild(v Visitor) bool {
	return visit(v, node.Head) || visitNodes(v, node.TemplateSpans)
}

// TemplateLiteralTypeSpan

type TemplateLiteralTypeSpan struct {
	NodeBase
	TypeNode *TypeNode             // TypeNode
	Literal  *TemplateMiddleOrTail // TemplateMiddleOrTail
}

func (f *NodeFactory) NewTemplateLiteralTypeSpan(typeNode *TypeNode, literal *TemplateMiddleOrTail) *Node {
	data := &TemplateLiteralTypeSpan{}
	data.TypeNode = typeNode
	data.Literal = literal
	return f.NewNode(ast.KindTemplateLiteralTypeSpan, data)
}

func (node *TemplateLiteralTypeSpan) ForEachChild(v Visitor) bool {
	return visit(v, node.TypeNode) || visit(v, node.Literal)
}

/// A JSX expression of the form <TagName attrs>...</TagName>

type JsxElement struct {
	ExpressionBase
	OpeningElement *JsxOpeningElementNode // JsxOpeningElementNode
	Children       []*JsxChild            // []JsxChild
	ClosingElement *JsxClosingElementNode // JsxClosingElementNode
}

func (f *NodeFactory) NewJsxElement(openingElement *JsxOpeningElementNode, children []*JsxChild, closingElement *JsxClosingElementNode) *Node {
	data := &JsxElement{}
	data.OpeningElement = openingElement
	data.Children = children
	data.ClosingElement = closingElement
	return f.NewNode(ast.KindJsxElement, data)
}

func (node *JsxElement) ForEachChild(v Visitor) bool {
	return visit(v, node.OpeningElement) || visitNodes(v, node.Children) || visit(v, node.ClosingElement)
}

// JsxAttributes
type JsxAttributes struct {
	ExpressionBase
	DeclarationBase
	Properties []*JsxAttributeLike // []JsxAttributeLike
}

func (f *NodeFactory) NewJsxAttributes(properties []*JsxAttributeLike) *Node {
	data := &JsxAttributes{}
	data.Properties = properties
	return f.NewNode(ast.KindJsxAttributes, data)
}

func (node *JsxAttributes) ForEachChild(v Visitor) bool {
	return visitNodes(v, node.Properties)
}

func IsJsxAttributes(node *Node) bool {
	return node.Kind == ast.KindJsxAttributes
}

// JsxNamespacedName

type JsxNamespacedName struct {
	ExpressionBase
	Name_      *IdentifierNode // IdentifierNode
	Namespace *IdentifierNode // IdentifierNode
}

func (f *NodeFactory) NewJsxNamespacedName(name *IdentifierNode, namespace *IdentifierNode) *Node {
	data := &JsxNamespacedName{}
	data.Name_ = name
	data.Namespace = namespace
	return f.NewNode(ast.KindJsxNamespacedName, data)
}

func (node *JsxNamespacedName) ForEachChild(v Visitor) bool {
	return visit(v, node.Name_) || visit(v, node.Namespace)
}

func IsJsxNamespacedName(node *Node) bool {
	return node.Kind == ast.KindJsxNamespacedName
}

/// The opening element of a <Tag>...</Tag> JsxElement

type JsxOpeningElement struct {
	ExpressionBase
	TagName       *JsxTagNameExpression // JsxTagNameExpression (Identifier | KeywordExpression | JsxTagNamePropertyAccess | JsxNamespacedName)
	TypeArguments *TypeArgumentListNode // TypeArgumentListNode
	Attributes    *JsxAttributesNode    // JsxAttributesNode
}

func (f *NodeFactory) NewJsxOpeningElement(tagName *JsxTagNameExpression, typeArguments *TypeArgumentListNode, attributes *JsxAttributesNode) *Node {
	data := &JsxOpeningElement{}
	data.TagName = tagName
	data.TypeArguments = typeArguments
	data.Attributes = attributes
	return f.NewNode(ast.KindJsxOpeningElement, data)
}

func (node *JsxOpeningElement) ForEachChild(v Visitor) bool {
	return visit(v, node.TagName) || visit(v, node.TypeArguments) || visit(v, node.Attributes)
}

func IsJsxOpeningElement(node *Node) bool {
	return node.Kind == ast.KindJsxOpeningElement
}

/// A JSX expression of the form <TagName attrs />

type JsxSelfClosingElement struct {
	ExpressionBase
	TagName       *JsxTagNameExpression // JsxTagNameExpression (IdentifierReference | KeywordExpression | JsxTagNamePropertyAccess | JsxNamespacedName)
	TypeArguments *TypeArgumentListNode // TypeArgumentListNode
	Attributes    *JsxAttributesNode    // JsxAttributesNode
}

func (f *NodeFactory) NewJsxSelfClosingElement(tagName *JsxTagNameExpression, typeArguments *TypeArgumentListNode, attributes *JsxAttributesNode) *Node {
	data := &JsxSelfClosingElement{}
	data.TagName = tagName
	data.TypeArguments = typeArguments
	data.Attributes = attributes
	return f.NewNode(ast.KindJsxSelfClosingElement, data)
}

func (node *JsxSelfClosingElement) ForEachChild(v Visitor) bool {
	return visit(v, node.TagName) || visit(v, node.TypeArguments) || visit(v, node.Attributes)
}

func IsJsxSelfClosingElement(node *Node) bool {
	return node.Kind == ast.KindJsxSelfClosingElement
}

/// A JSX expression of the form <>...</>

type JsxFragment struct {
	ExpressionBase
	OpeningFragment *JsxOpeningFragmentNode // JsxOpeningFragmentNode
	Children        []*JsxChild             // []JsxChild
	ClosingFragment *JsxClosingFragmentNode // JsxClosingFragmentNode
}

func (f *NodeFactory) NewJsxFragment(openingFragment *JsxOpeningFragmentNode, children []*JsxChild, closingFragment *JsxClosingFragmentNode) *Node {
	data := &JsxFragment{}
	data.OpeningFragment = openingFragment
	data.Children = children
	data.ClosingFragment = closingFragment
	return f.NewNode(ast.KindJsxFragment, data)
}

func (node *JsxFragment) ForEachChild(v Visitor) bool {
	return visit(v, node.OpeningFragment) || visitNodes(v, node.Children) || visit(v, node.ClosingFragment)
}

/// The opening element of a <>...</> JsxFragment

type JsxOpeningFragment struct {
	ExpressionBase
}

func (f *NodeFactory) NewJsxOpeningFragment() *Node {
	return f.NewNode(ast.KindJsxOpeningFragment, &JsxOpeningFragment{})
}

func IsJsxOpeningFragment(node *Node) bool {
	return node.Kind == ast.KindJsxOpeningFragment
}

/// The closing element of a <>...</> JsxFragment

type JsxClosingFragment struct {
	ExpressionBase
}

func (f *NodeFactory) NewJsxClosingFragment() *Node {
	return f.NewNode(ast.KindJsxClosingFragment, &JsxClosingFragment{})
}

// JsxAttribute

type JsxAttribute struct {
	NodeBase
	DeclarationBase
	Name_        *JsxAttributeName  // JsxAttributeName
	Initializer *JsxAttributeValue // JsxAttributeValue. Optional, <X y /> is sugar for <X y={true} />
}

func (f *NodeFactory) NewJsxAttribute(name *JsxAttributeName, initializer *JsxAttributeValue) *Node {
	data := &JsxAttribute{}
	data.Name_ = name
	data.Initializer = initializer
	return f.NewNode(ast.KindJsxAttribute, data)
}

func (node *JsxAttribute) ForEachChild(v Visitor) bool {
	return visit(v, node.Name_) || visit(v, node.Initializer)
}

func IsJsxAttribute(node *Node) bool {
	return node.Kind == ast.KindJsxAttribute
}

// JsxSpreadAttribute

type JsxSpreadAttribute struct {
	NodeBase
	Expression *Expression // Expression
}

func (f *NodeFactory) NewJsxSpreadAttribute(expression *Expression) *Node {
	data := &JsxSpreadAttribute{}
	data.Expression = expression
	return f.NewNode(ast.KindJsxSpreadAttribute, data)
}

func (node *JsxSpreadAttribute) ForEachChild(v Visitor) bool {
	return visit(v, node.Expression)
}

// JsxClosingElement

type JsxClosingElement struct {
	NodeBase
	TagName *JsxTagNameExpression // JsxTagNameExpression
}

func (f *NodeFactory) NewJsxClosingElement(tagName *JsxTagNameExpression) *Node {
	data := &JsxClosingElement{}
	data.TagName = tagName
	return f.NewNode(ast.KindJsxClosingElement, data)
}

func (node *JsxClosingElement) ForEachChild(v Visitor) bool {
	return visit(v, node.TagName)
}

// JsxExpression

type JsxExpression struct {
	ExpressionBase
	DotDotDotToken *TokenNode  // TokenNode. Optional
	Expression     *Expression // Expression
}

func (f *NodeFactory) NewJsxExpression(dotDotDotToken *TokenNode, expression *Expression) *Node {
	data := &JsxExpression{}
	data.DotDotDotToken = dotDotDotToken
	data.Expression = expression
	return f.NewNode(ast.KindJsxExpression, data)
}

func (node *JsxExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.DotDotDotToken) || visit(v, node.Expression)
}

// JsxText

type JsxText struct {
	ExpressionBase
	LiteralLikeBase
	ContainsOnlyTriviaWhiteSpaces bool
}

func (f *NodeFactory) NewJsxText(text string, containsOnlyTriviaWhiteSpace bool) *Node {
	data := &JsxText{}
	data.Text = text
	data.ContainsOnlyTriviaWhiteSpaces = containsOnlyTriviaWhiteSpace
	return f.NewNode(ast.KindJsxText, data)
}

// JSDocNonNullableType

type JSDocNonNullableType struct {
	TypeNodeBase
	TypeNode *TypeNode // TypeNode
	Postfix  bool
}

func (f *NodeFactory) NewJSDocNonNullableType(typeNode *TypeNode, postfix bool) *Node {
	data := &JSDocNonNullableType{}
	data.TypeNode = typeNode
	data.Postfix = postfix
	return f.NewNode(ast.KindJSDocNonNullableType, data)
}

func (node *JSDocNonNullableType) ForEachChild(v Visitor) bool {
	return visit(v, node.TypeNode)
}

// JSDocNullableType

type JSDocNullableType struct {
	TypeNodeBase
	TypeNode *TypeNode // TypeNode
	Postfix  bool
}

func (f *NodeFactory) NewJSDocNullableType(typeNode *TypeNode, postfix bool) *Node {
	data := &JSDocNullableType{}
	data.TypeNode = typeNode
	data.Postfix = postfix
	return f.NewNode(ast.KindJSDocNullableType, data)
}

func (node *JSDocNullableType) ForEachChild(v Visitor) bool {
	return visit(v, node.TypeNode)
}

// PatternAmbientModule

type PatternAmbientModule struct {
	Pattern Pattern
	Symbol  *Symbol
}

// SourceFile

type SourceFile struct {
	NodeBase
	DeclarationBase
	LocalsContainerBase
	Text                        string
	FileName_                    string
	Path_                        string
	Statements                  []*Statement // []Statement
	Diagnostics_                 []*Diagnostic
	BindDiagnostics_             []*Diagnostic
	BindSuggestionDiagnostics   []*Diagnostic
	LineMap                     []textpos.TextPos
	LanguageVersion             core.ScriptTarget
	LanguageVariant             core.LanguageVariant
	ScriptKind                  core.ScriptKind
	ExternalModuleIndicator     *Node
	EndFlowNode                 *FlowNode
	JsGlobalAugmentations       SymbolTable
	IsDeclarationFile           bool
	IsBound                     bool
	ModuleReferencesProcessed   bool
	UsesUriStyleNodeCoreModules core.Tristate
	SymbolCount                 int
	ClassifiableNames           set[string]
	Imports                     []*LiteralLikeNode // []LiteralLikeNode
	ModuleAugmentations         []*ModuleName      // []ModuleName
	PatternAmbientModules       []PatternAmbientModule
	AmbientModuleNames          []string
}

func (f *NodeFactory) NewSourceFile(text string, fileName string, statements []*Statement) *Node {
	data := &SourceFile{}
	data.Text = text
	data.FileName_ = fileName
	data.Statements = statements
	data.LanguageVersion = core.ScriptTargetLatest
	return f.NewNode(ast.KindSourceFile, data)
}

func (node *SourceFile) FileName() string {
	return node.FileName_
}

func (node *SourceFile) Path() string {
	return node.Path_
}

func (node *SourceFile) Diagnostics() []*Diagnostic {
	return node.Diagnostics_
}

func (node *SourceFile) BindDiagnostics() []*Diagnostic {
	return node.BindDiagnostics_
}

func (node *SourceFile) ForEachChild(v Visitor) bool {
	return visitNodes(v, node.Statements)
}

func IsSourceFile(node *Node) bool {
	return node.Kind == ast.KindSourceFile
}
