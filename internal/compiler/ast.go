package compiler

import (
	"github.com/microsoft/typescript-go/internal/compiler/textpos"
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
	identifierPool Pool[Identifier]
}

func (f *NodeFactory) NewNode(kind SyntaxKind, data NodeData) *Node {
	n := data.AsNode()
	n.Kind = kind
	n.Data = data
	return n
}

// AST Node
// Interface values stored in AST nodes are never typed nil values. Construction code must ensure that
// interface valued properties either store a true nil or a reference to a non-nil struct.

type Node struct {
	Kind   SyntaxKind
	Flags  NodeFlags
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
func (n *Node) GetName() *DeclarationName                 { return n.Data.GetName() }
func (n *Node) GetModifiers() *ModifierListNode           { return n.Data.GetModifiers() }
func (n *Node) GetTypeParameters() *TypeParameterListNode { return n.Data.GetTypeParameters() }
func (n *Node) FlowNodeData() *FlowNodeBase               { return n.Data.FlowNodeData() }
func (n *Node) DeclarationData() *DeclarationBase         { return n.Data.DeclarationData() }
func (n *Node) Symbol() *Symbol                           { return n.Data.DeclarationData().Symbol }
func (n *Node) ExportableData() *ExportableBase           { return n.Data.ExportableData() }
func (n *Node) LocalSymbol() *Symbol                      { return n.Data.ExportableData().LocalSymbol }
func (n *Node) LocalsContainerData() *LocalsContainerBase { return n.Data.LocalsContainerData() }
func (n *Node) Locals() SymbolTable                       { return n.Data.LocalsContainerData().locals }
func (n *Node) FunctionLikeData() *FunctionLikeBase       { return n.Data.FunctionLikeData() }
func (n *Node) Parameters() []*ParameterDeclarationNode   { return n.Data.FunctionLikeData().Parameters }
func (n *Node) ReturnType() *TypeNode                     { return n.Data.FunctionLikeData().ReturnType }
func (n *Node) ClassLikeData() *ClassLikeBase             { return n.Data.ClassLikeData() }
func (n *Node) BodyData() *BodyBase                       { return n.Data.BodyData() }

func (n *Node) Text() string {
	switch n.Kind {
	case SyntaxKindIdentifier:
		return n.AsIdentifier().Text
	case SyntaxKindPrivateIdentifier:
		return n.AsPrivateIdentifier().Text
	case SyntaxKindStringLiteral:
		return n.AsStringLiteral().Text
	case SyntaxKindNumericLiteral:
		return n.AsNumericLiteral().Text
	case SyntaxKindBigIntLiteral:
		return n.AsBigIntLiteral().Text
	case SyntaxKindNoSubstitutionTemplateLiteral:
		return n.AsNoSubstitutionTemplateLiteral().Text
	case SyntaxKindJsxNamespacedName:
		return n.AsJsxNamespacedName().Namespace.Text() + ":" + n.AsJsxNamespacedName().Name.Text()
	}
	panic("Unhandled case in Node.Text")
}

func (node *Node) Expression() *Node {
	switch node.Kind {
	case SyntaxKindPropertyAccessExpression:
		return node.AsPropertyAccessExpression().Expression
	case SyntaxKindElementAccessExpression:
		return node.AsElementAccessExpression().Expression
	case SyntaxKindParenthesizedExpression:
		return node.AsParenthesizedExpression().Expression
	case SyntaxKindCallExpression:
		return node.AsCallExpression().Expression
	case SyntaxKindNewExpression:
		return node.AsNewExpression().Expression
	case SyntaxKindExpressionWithTypeArguments:
		return node.AsExpressionWithTypeArguments().Expression
	case SyntaxKindNonNullExpression:
		return node.AsNonNullExpression().Expression
	case SyntaxKindTypeAssertionExpression:
		return node.AsTypeAssertion().Expression
	case SyntaxKindAsExpression:
		return node.AsAsExpression().Expression
	case SyntaxKindSatisfiesExpression:
		return node.AsSatisfiesExpression().Expression
	case SyntaxKindSpreadAssignment:
		return node.AsSpreadAssignment().Expression
	}
	panic("Unhandled case in Node.Expression")
}

func (node *Node) Arguments() []*Node {
	switch node.Kind {
	case SyntaxKindCallExpression:
		return node.AsCallExpression().Arguments
	case SyntaxKindNewExpression:
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
	GetName() *DeclarationName
	GetModifiers() *ModifierListNode
	GetTypeParameters() *TypeParameterListNode
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
func (node *NodeDefault) GetName() *DeclarationName                 { return nil }
func (node *NodeDefault) GetModifiers() *ModifierListNode           { return nil }
func (node *NodeDefault) GetTypeParameters() *TypeParameterListNode { return nil }
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
type JsxAttributeLike = Node            // JsxAttribute | JsxSpreadAttribute
type JsxAttributeName = Node            // Identifier | JsxNamespacedName
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
type BlockNode = Node
type CatchClauseNode = Node
type CaseBlockNode = Node
type CaseOrDefaultClauseNode = Node
type VariableDeclarationNode = Node
type VariableDeclarationListNode = Node
type BindingElementNode = Node
type TypeParameterListNode = Node
type ParameterDeclarationNode = Node
type HeritageClauseNode = Node
type ExpressionWithTypeArgumentsNode = Node
type EnumMemberNode = Node
type ImportClauseNode = Node
type ImportAttributesNode = Node
type ImportSpecifierNode = Node
type ExportSpecifierNode = Node

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
	Modifiers *ModifierListNode
}

func (node *ModifiersBase) GetModifiers() *ModifierListNode { return node.Modifiers }

// LocalsContainerBase

type LocalsContainerBase struct {
	locals        SymbolTable // Locals associated with node (initialized by binding)
	NextContainer *Node       // Next container in declaration order (initialized by binding)
}

func (node *LocalsContainerBase) LocalsContainerData() *LocalsContainerBase { return node }

func IsLocalsContainer(node *Node) bool {
	return node.LocalsContainerData() != nil
}

// FunctionLikeBase

type FunctionLikeBase struct {
	LocalsContainerBase
	TypeParameters *TypeParameterListNode // Optional
	Parameters     []*ParameterDeclarationNode
	ReturnType     *TypeNode // Optional
}

func (node *FunctionLikeBase) GetTypeParameters() *TypeParameterListNode { return node.TypeParameters }
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

func (node *FunctionLikeWithBodyBase) GetTypeParameters() *TypeParameterListNode {
	return node.TypeParameters
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

func (f *NodeFactory) NewToken(kind SyntaxKind) *Node {
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
	return f.NewNode(SyntaxKindIdentifier, data)
}

func IsIdentifier(node *Node) bool {
	return node.Kind == SyntaxKindIdentifier
}

// PrivateIdentifier

type PrivateIdentifier struct {
	ExpressionBase
	Text string
}

func (f *NodeFactory) NewPrivateIdentifier(text string) *Node {
	data := &PrivateIdentifier{}
	data.Text = text
	return f.NewNode(SyntaxKindPrivateIdentifier, data)
}

func IsPrivateIdentifier(node *Node) bool {
	return node.Kind == SyntaxKindPrivateIdentifier
}

// QualifiedName

type QualifiedName struct {
	NodeBase
	FlowNodeBase
	Left  *EntityName
	Right *IdentifierNode
}

func (f *NodeFactory) NewQualifiedName(left *EntityName, right *IdentifierNode) *Node {
	data := &QualifiedName{}
	data.Left = left
	data.Right = right
	return f.NewNode(SyntaxKindQualifiedName, data)
}

func (node *QualifiedName) ForEachChild(v Visitor) bool {
	return visit(v, node.Left) || visit(v, node.Right)
}

func IsQualifiedName(node *Node) bool {
	return node.Kind == SyntaxKindQualifiedName
}

// TypeParameterDeclaration

type TypeParameterDeclaration struct {
	NodeBase
	DeclarationBase
	ModifiersBase
	Name        *IdentifierNode // Identifier
	Constraint  *TypeNode       // Optional
	DefaultType *TypeNode       // Optional
	Expression  *Node           // For error recovery purposes
}

func (f *NodeFactory) NewTypeParameterDeclaration(modifiers *Node, name *IdentifierNode, constraint *TypeNode, defaultType *TypeNode) *Node {
	data := &TypeParameterDeclaration{}
	data.Modifiers = modifiers
	data.Name = name
	data.Constraint = constraint
	data.DefaultType = defaultType
	return f.NewNode(SyntaxKindTypeParameter, data)
}

func (node *TypeParameterDeclaration) Kind() SyntaxKind {
	return SyntaxKindTypeParameter
}

func (node *TypeParameterDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.Modifiers) || visit(v, node.Name) || visit(v, node.Constraint) || visit(v, node.DefaultType)
}

func (node *TypeParameterDeclaration) GetName() *DeclarationName {
	return node.Name
}

func IsTypeParameterDeclaration(node *Node) bool {
	return node.Kind == SyntaxKindTypeParameter
}

// ComputedPropertyName

type ComputedPropertyName struct {
	NodeBase
	Expression *Node
}

func (f *NodeFactory) NewComputedPropertyName(expression *Node) *Node {
	data := &ComputedPropertyName{}
	data.Expression = expression
	return f.NewNode(SyntaxKindComputedPropertyName, data)
}

func (node *ComputedPropertyName) ForEachChild(v Visitor) bool {
	return visit(v, node.Expression)
}

func IsComputedPropertyName(node *Node) bool {
	return node.Kind == SyntaxKindComputedPropertyName
}

// Modifier

func (f *NodeFactory) NewModifier(kind SyntaxKind) *Node {
	return f.NewToken(kind)
}

// Decorator

type Decorator struct {
	NodeBase
	Expression *Node
}

func (f *NodeFactory) NewDecorator(expression *Node) *Node {
	data := &Decorator{}
	data.Expression = expression
	return f.NewNode(SyntaxKindDecorator, data)
}

func (node *Decorator) ForEachChild(v Visitor) bool {
	return visit(v, node.Expression)
}

func IsDecorator(node *Node) bool {
	return node.Kind == SyntaxKindDecorator
}

// ModifierList

type ModifierList struct {
	NodeBase
	Modifiers     []*ModifierLike
	ModifierFlags ModifierFlags
}

func (f *NodeFactory) NewModifierList(modifiers []*ModifierLike, modifierFlags ModifierFlags) *Node {
	data := &ModifierList{}
	data.Modifiers = modifiers
	data.ModifierFlags = modifierFlags
	return f.NewNode(SyntaxKindModifierList, data)
}

func (node *ModifierList) ForEachChild(v Visitor) bool {
	return visitNodes(v, node.Modifiers)
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

func IsEmptyStatement(node *Node) bool {
	return node.Kind == SyntaxKindEmptyStatement
}

// IfStatement

type IfStatement struct {
	StatementBase
	Expression    *Node
	ThenStatement *Statement
	ElseStatement *Statement // Optional
}

func (f *NodeFactory) NewIfStatement(expression *Node, thenStatement *Statement, elseStatement *Statement) *Node {
	data := &IfStatement{}
	data.Expression = expression
	data.ThenStatement = thenStatement
	data.ElseStatement = elseStatement
	return f.NewNode(SyntaxKindIfStatement, data)
}

func (node *IfStatement) ForEachChild(v Visitor) bool {
	return visit(v, node.Expression) || visit(v, node.ThenStatement) || visit(v, node.ElseStatement)
}

// DoStatement

type DoStatement struct {
	StatementBase
	Statement  *Statement
	Expression *Node
}

func (f *NodeFactory) NewDoStatement(statement *Statement, expression *Node) *Node {
	data := &DoStatement{}
	data.Statement = statement
	data.Expression = expression
	return f.NewNode(SyntaxKindDoStatement, data)
}

func (node *DoStatement) ForEachChild(v Visitor) bool {
	return visit(v, node.Statement) || visit(v, node.Expression)
}

// WhileStatement

type WhileStatement struct {
	StatementBase
	Expression *Node
	Statement  *Statement
}

func (f *NodeFactory) NewWhileStatement(expression *Node, statement *Statement) *Node {
	data := &WhileStatement{}
	data.Expression = expression
	data.Statement = statement
	return f.NewNode(SyntaxKindWhileStatement, data)
}

func (node *WhileStatement) ForEachChild(v Visitor) bool {
	return visit(v, node.Expression) || visit(v, node.Statement)
}

// ForStatement

type ForStatement struct {
	StatementBase
	LocalsContainerBase
	Initializer *ForInitializer // Optional
	Condition   *Node           // Optional
	Incrementor *Node           // Optional
	Statement   *Statement
}

func (f *NodeFactory) NewForStatement(initializer *ForInitializer, condition *Node, incrementor *Node, statement *Statement) *Node {
	data := &ForStatement{}
	data.Initializer = initializer
	data.Condition = condition
	data.Incrementor = incrementor
	data.Statement = statement
	return f.NewNode(SyntaxKindForStatement, data)
}

func (node *ForStatement) ForEachChild(v Visitor) bool {
	return visit(v, node.Initializer) || visit(v, node.Condition) || visit(v, node.Incrementor) || visit(v, node.Statement)
}

// ForInOrOfStatement

type ForInOrOfStatement struct {
	StatementBase
	LocalsContainerBase
	Kind          SyntaxKind // SyntaxKindForInStatement | SyntaxKindForOfStatement
	AwaitModifier *Node      // Optional
	Initializer   *ForInitializer
	Expression    *Node
	Statement     *Statement
}

func (f *NodeFactory) NewForInOrOfStatement(kind SyntaxKind, awaitModifier *Node, initializer *ForInitializer, expression *Node, statement *Statement) *Node {
	data := &ForInOrOfStatement{}
	data.Kind = kind
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
	return node.Kind == SyntaxKindForInStatement || node.Kind == SyntaxKindForOfStatement
}

// BreakStatement

type BreakStatement struct {
	StatementBase
	Label *IdentifierNode // Optional
}

func (f *NodeFactory) NewBreakStatement(label *IdentifierNode) *Node {
	data := &BreakStatement{}
	data.Label = label
	return f.NewNode(SyntaxKindBreakStatement, data)
}

func (node *BreakStatement) ForEachChild(v Visitor) bool {
	return visit(v, node.Label)
}

// ContinueStatement

type ContinueStatement struct {
	StatementBase
	Label *IdentifierNode // Optional
}

func (f *NodeFactory) NewContinueStatement(label *IdentifierNode) *Node {
	data := &ContinueStatement{}
	data.Label = label
	return f.NewNode(SyntaxKindContinueStatement, data)
}

func (node *ContinueStatement) ForEachChild(v Visitor) bool {
	return visit(v, node.Label)
}

// ReturnStatement

type ReturnStatement struct {
	StatementBase
	Expression *Node // Optional
}

func (f *NodeFactory) NewReturnStatement(expression *Node) *Node {
	data := &ReturnStatement{}
	data.Expression = expression
	return f.NewNode(SyntaxKindReturnStatement, data)
}

func (node *ReturnStatement) ForEachChild(v Visitor) bool {
	return visit(v, node.Expression)
}

// WithStatement

type WithStatement struct {
	StatementBase
	Expression *Node
	Statement  *Statement
}

func (f *NodeFactory) NewWithStatement(expression *Node, statement *Statement) *Node {
	data := &WithStatement{}
	data.Expression = expression
	data.Statement = statement
	return f.NewNode(SyntaxKindWithStatement, data)
}

func (node *WithStatement) ForEachChild(v Visitor) bool {
	return visit(v, node.Expression) || visit(v, node.Statement)
}

// SwitchStatement

type SwitchStatement struct {
	StatementBase
	Expression *Node
	CaseBlock  *CaseBlockNode
}

func (f *NodeFactory) NewSwitchStatement(expression *Node, caseBlock *CaseBlockNode) *Node {
	data := &SwitchStatement{}
	data.Expression = expression
	data.CaseBlock = caseBlock
	return f.NewNode(SyntaxKindSwitchStatement, data)
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
	return f.NewNode(SyntaxKindCaseBlock, data)
}

func (node *CaseBlock) ForEachChild(v Visitor) bool {
	return visitNodes(v, node.Clauses)
}

// CaseOrDefaultClause

type CaseOrDefaultClause struct {
	NodeBase
	Expression          *Node // nil in default clause
	Statements          []*Statement
	FallthroughFlowNode *FlowNode
}

func (f *NodeFactory) NewCaseOrDefaultClause(kind SyntaxKind, expression *Node, statements []*Statement) *Node {
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
	Expression *Node
}

func (f *NodeFactory) NewThrowStatement(expression *Node) *Node {
	data := &ThrowStatement{}
	data.Expression = expression
	return f.NewNode(SyntaxKindThrowStatement, data)
}

func (node *ThrowStatement) ForEachChild(v Visitor) bool {
	return visit(v, node.Expression)
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

func (node *TryStatement) Kind() SyntaxKind {
	return SyntaxKindTryStatement
}

func (node *TryStatement) ForEachChild(v Visitor) bool {
	return visit(v, node.tryBlock) || visit(v, node.catchClause) || visit(v, node.finallyBlock)
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
	return f.NewNode(SyntaxKindCatchClause, data)
}

func (node *CatchClause) ForEachChild(v Visitor) bool {
	return visit(v, node.VariableDeclaration) || visit(v, node.Block)
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
	Label     *IdentifierNode
	Statement *Statement
}

func (f *NodeFactory) NewLabeledStatement(label *IdentifierNode, statement *Statement) *Node {
	data := &LabeledStatement{}
	data.Label = label
	data.Statement = statement
	return f.NewNode(SyntaxKindLabeledStatement, data)
}

func (node *LabeledStatement) ForEachChild(v Visitor) bool {
	return visit(v, node.Label) || visit(v, node.Statement)
}

// ExpressionStatement

type ExpressionStatement struct {
	StatementBase
	Expression *Node
}

func (f *NodeFactory) NewExpressionStatement(expression *Node) *Node {
	data := &ExpressionStatement{}
	data.Expression = expression
	return f.NewNode(SyntaxKindExpressionStatement, data)
}

func (node *ExpressionStatement) ForEachChild(v Visitor) bool {
	return visit(v, node.Expression)
}

func IsExpressionStatement(node *Node) bool {
	return node.Kind == SyntaxKindExpressionStatement
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
	return f.NewNode(SyntaxKindBlock, data)
}

func (node *Block) ForEachChild(v Visitor) bool {
	return visitNodes(v, node.Statements)
}

func IsBlock(node *Node) bool {
	return node.Kind == SyntaxKindBlock
}

// VariableStatement

type VariableStatement struct {
	StatementBase
	ModifiersBase
	DeclarationList *VariableDeclarationListNode
}

func (f *NodeFactory) NewVariableStatement(modifiers *ModifierListNode, declarationList *VariableDeclarationListNode) *Node {
	data := &VariableStatement{}
	data.Modifiers = modifiers
	data.DeclarationList = declarationList
	return f.NewNode(SyntaxKindVariableStatement, data)
}

func (node *VariableStatement) ForEachChild(v Visitor) bool {
	return visit(v, node.Modifiers) || visit(v, node.DeclarationList)
}

func IsVariableStatement(node *Node) bool {
	return node.Kind == SyntaxKindVariableStatement
}

// VariableDeclaration

type VariableDeclaration struct {
	NodeBase
	DeclarationBase
	ExportableBase
	Name             *BindingName
	ExclamationToken *TokenNode // Optional
	TypeNode         *TypeNode  // Optional
	Initializer      *Node      // Optional
}

func (f *NodeFactory) NewVariableDeclaration(name *BindingName, exclamationToken *TokenNode, typeNode *TypeNode, initializer *Node) *Node {
	data := &VariableDeclaration{}
	data.Name = name
	data.ExclamationToken = exclamationToken
	data.TypeNode = typeNode
	data.Initializer = initializer
	return f.NewNode(SyntaxKindVariableDeclaration, data)
}

func (node *VariableDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.Name) || visit(v, node.ExclamationToken) || visit(v, node.TypeNode) || visit(v, node.Initializer)
}

func (node *VariableDeclaration) GetName() *DeclarationName {
	return node.Name
}

func IsVariableDeclaration(node *Node) bool {
	return node.Kind == SyntaxKindVariableDeclaration
}

// VariableDeclarationList

type VariableDeclarationList struct {
	NodeBase
	Declarations []*VariableDeclarationNode
}

func (f *NodeFactory) NewVariableDeclarationList(flags NodeFlags, declarations []*VariableDeclarationNode) *Node {
	data := &VariableDeclarationList{}
	data.Declarations = declarations
	node := f.NewNode(SyntaxKindVariableDeclarationList, data)
	node.Flags = flags
	return node
}

func (node *VariableDeclarationList) ForEachChild(v Visitor) bool {
	return visitNodes(v, node.Declarations)
}

func IsVariableDeclarationList(node *Node) bool {
	return node.Kind == SyntaxKindVariableDeclarationList
}

// BindingPattern (SyntaxBindObjectBindingPattern | SyntaxKindArrayBindingPattern)

type BindingPattern struct {
	NodeBase
	Elements []*BindingElementNode
}

func (f *NodeFactory) NewBindingPattern(kind SyntaxKind, elements []*BindingElementNode) *Node {
	data := &BindingPattern{}
	data.Elements = elements
	return f.NewNode(kind, data)
}

func (node *BindingPattern) ForEachChild(v Visitor) bool {
	return visitNodes(v, node.Elements)
}

func IsObjectBindingPattern(node *Node) bool {
	return node.Kind == SyntaxKindObjectBindingPattern
}

func IsArrayBindingPattern(node *Node) bool {
	return node.Kind == SyntaxKindArrayBindingPattern
}

// ParameterDeclaration

type ParameterDeclaration struct {
	NodeBase
	DeclarationBase
	ModifiersBase
	DotDotDotToken *TokenNode   // Present on rest parameter
	Name           *BindingName // Declared parameter name
	QuestionToken  *TokenNode   // Present on optional parameter
	TypeNode       *TypeNode    // Optional
	Initializer    *Node        // Optional
}

func (f *NodeFactory) NewParameterDeclaration(modifiers *ModifierListNode, dotDotDotToken *TokenNode, name *BindingName, questionToken *TokenNode, typeNode *TypeNode, initializer *Node) *Node {
	data := &ParameterDeclaration{}
	data.Modifiers = modifiers
	data.DotDotDotToken = dotDotDotToken
	data.Name = name
	data.QuestionToken = questionToken
	data.TypeNode = typeNode
	data.Initializer = initializer
	return f.NewNode(SyntaxKindParameter, data)
}

func (node *ParameterDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.Modifiers) || visit(v, node.DotDotDotToken) || visit(v, node.Name) ||
		visit(v, node.QuestionToken) || visit(v, node.TypeNode) || visit(v, node.Initializer)
}

func (node *ParameterDeclaration) GetName() *DeclarationName {
	return node.Name
}

func IsParameter(node *Node) bool {
	return node.Kind == SyntaxKindParameter
}

// BindingElement

type BindingElement struct {
	NodeBase
	DeclarationBase
	ExportableBase
	FlowNodeBase
	DotDotDotToken *TokenNode    // Present on rest element (in object binding pattern)
	PropertyName   *PropertyName // Optional binding property name in object binding pattern
	Name           *BindingName  // Optional (nil for missing element)
	Initializer    *Node         // Optional
}

func (f *NodeFactory) NewBindingElement(dotDotDotToken *TokenNode, propertyName *PropertyName, name *BindingName, initializer *Node) *Node {
	data := &BindingElement{}
	data.DotDotDotToken = dotDotDotToken
	data.PropertyName = propertyName
	data.Name = name
	data.Initializer = initializer
	return f.NewNode(SyntaxKindBindingElement, data)
}

func (node *BindingElement) ForEachChild(v Visitor) bool {
	return visit(v, node.PropertyName) || visit(v, node.DotDotDotToken) || visit(v, node.Name) || visit(v, node.Initializer)
}

func (node *BindingElement) GetName() *DeclarationName {
	return node.Name
}

func IsBindingElement(node *Node) bool {
	return node.Kind == SyntaxKindBindingElement
}

// MissingDeclaration

type MissingDeclaration struct {
	StatementBase
	DeclarationBase
	ModifiersBase
}

func (f *NodeFactory) NewMissingDeclaration(modifiers *ModifierListNode) *Node {
	data := &MissingDeclaration{}
	data.Modifiers = modifiers
	return f.NewNode(SyntaxKindMissingDeclaration, data)
}

func (node *MissingDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.Modifiers)
}

// FunctionDeclaration

type FunctionDeclaration struct {
	StatementBase
	DeclarationBase
	ExportableBase
	ModifiersBase
	FunctionLikeWithBodyBase
	Name           *IdentifierNode
	ReturnFlowNode *FlowNode
}

func (f *NodeFactory) NewFunctionDeclaration(modifiers *ModifierListNode, asteriskToken *TokenNode, name *IdentifierNode, typeParameters *TypeParameterListNode, parameters []*ParameterDeclarationNode, returnType *TypeNode, body *BlockNode) *Node {
	data := &FunctionDeclaration{}
	data.Modifiers = modifiers
	data.AsteriskToken = asteriskToken
	data.Name = name
	data.TypeParameters = typeParameters
	data.Parameters = parameters
	data.ReturnType = returnType
	data.Body = body
	return f.NewNode(SyntaxKindFunctionDeclaration, data)
}

func (node *FunctionDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.Modifiers) || visit(v, node.AsteriskToken) || visit(v, node.Name) || visit(v, node.TypeParameters) ||
		visitNodes(v, node.Parameters) || visit(v, node.ReturnType) || visit(v, node.Body)
}

func (node *FunctionDeclaration) GetName() *DeclarationName {
	return node.Name
}

func (node *FunctionDeclaration) BodyData() *BodyBase { return &node.BodyBase }

func IsFunctionDeclaration(node *Node) bool {
	return node.Kind == SyntaxKindFunctionDeclaration
}

// ClassLikeDeclarationBase

type ClassLikeBase struct {
	DeclarationBase
	ExportableBase
	ModifiersBase
	Name            *IdentifierNode
	TypeParameters  *TypeParameterListNode
	HeritageClauses []*HeritageClauseNode
	Members         []*ClassElement
}

func (node *ClassLikeBase) ForEachChild(v Visitor) bool {
	return visit(v, node.Modifiers) || visit(v, node.Name) || visit(v, node.TypeParameters) ||
		visitNodes(v, node.HeritageClauses) || visitNodes(v, node.Members)
}

func (node *ClassLikeBase) GetName() *DeclarationName                 { return node.Name }
func (node *ClassLikeBase) GetTypeParameters() *TypeParameterListNode { return node.TypeParameters }
func (node *ClassLikeBase) ClassLikeData() *ClassLikeBase             { return node }

// ClassDeclaration

type ClassDeclaration struct {
	StatementBase
	ClassLikeBase
}

func (f *NodeFactory) NewClassDeclaration(modifiers *ModifierListNode, name *IdentifierNode, typeParameters *TypeParameterListNode, heritageClauses []*HeritageClauseNode, members []*ClassElement) *Node {
	data := &ClassDeclaration{}
	data.Modifiers = modifiers
	data.Name = name
	data.TypeParameters = typeParameters
	data.HeritageClauses = heritageClauses
	data.Members = members
	return f.NewNode(SyntaxKindClassDeclaration, data)
}

func IsClassDeclaration(node *Node) bool {
	return node.Kind == SyntaxKindClassDeclaration
}

// ClassExpression

type ClassExpression struct {
	ExpressionBase
	ClassLikeBase
}

func (f *NodeFactory) NewClassExpression(modifiers *ModifierListNode, name *IdentifierNode, typeParameters *TypeParameterListNode, heritageClauses []*HeritageClauseNode, members []*ClassElement) *Node {
	data := &ClassExpression{}
	data.Modifiers = modifiers
	data.Name = name
	data.TypeParameters = typeParameters
	data.HeritageClauses = heritageClauses
	data.Members = members
	return f.NewNode(SyntaxKindClassExpression, data)
}

func (node *ClassExpression) Kind() SyntaxKind { return SyntaxKindClassExpression }

func IsClassExpression(node *Node) bool {
	return node.Kind == SyntaxKindClassExpression
}

// HeritageClause

type HeritageClause struct {
	NodeBase
	Token SyntaxKind
	Types []*ExpressionWithTypeArgumentsNode
}

func (f *NodeFactory) NewHeritageClause(token SyntaxKind, types []*ExpressionWithTypeArgumentsNode) *Node {
	data := &HeritageClause{}
	data.Token = token
	data.Types = types
	return f.NewNode(SyntaxKindHeritageClause, data)
}

func (node *HeritageClause) ForEachChild(v Visitor) bool {
	return visitNodes(v, node.Types)
}

func IsHeritageClause(node *Node) bool {
	return node.Kind == SyntaxKindHeritageClause
}

// InterfaceDeclaration

type InterfaceDeclaration struct {
	StatementBase
	DeclarationBase
	ExportableBase
	ModifiersBase
	Name            *IdentifierNode
	TypeParameters  *TypeParameterListNode
	HeritageClauses []*HeritageClauseNode
	Members         []*TypeElement
}

func (f *NodeFactory) NewInterfaceDeclaration(modifiers *ModifierListNode, name *IdentifierNode, typeParameters *TypeParameterListNode, heritageClauses []*HeritageClauseNode, members []*TypeElement) *Node {
	data := &InterfaceDeclaration{}
	data.Modifiers = modifiers
	data.Name = name
	data.TypeParameters = typeParameters
	data.HeritageClauses = heritageClauses
	data.Members = members
	return f.NewNode(SyntaxKindInterfaceDeclaration, data)
}

func (node *InterfaceDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.Modifiers) || visit(v, node.Name) || visit(v, node.TypeParameters) ||
		visitNodes(v, node.HeritageClauses) || visitNodes(v, node.Members)
}

func (node *InterfaceDeclaration) GetName() *DeclarationName { return node.Name }
func (node *InterfaceDeclaration) GetTypeParameters() *TypeParameterListNode {
	return node.TypeParameters
}

func IsInterfaceDeclaration(node *Node) bool {
	return node.Kind == SyntaxKindInterfaceDeclaration
}

// TypeAliasDeclaration

type TypeAliasDeclaration struct {
	StatementBase
	DeclarationBase
	ExportableBase
	ModifiersBase
	LocalsContainerBase
	Name           *IdentifierNode
	TypeParameters *TypeParameterListNode
	TypeNode       *TypeNode
}

func (f *NodeFactory) NewTypeAliasDeclaration(modifiers *ModifierListNode, name *IdentifierNode, typeParameters *TypeParameterListNode, typeNode *TypeNode) *Node {
	data := &TypeAliasDeclaration{}
	data.Modifiers = modifiers
	data.Name = name
	data.TypeParameters = typeParameters
	data.TypeNode = typeNode
	return f.NewNode(SyntaxKindTypeAliasDeclaration, data)
}

func (node *TypeAliasDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.Modifiers) || visit(v, node.Name) || visit(v, node.TypeParameters) || visit(v, node.TypeNode)
}

func (node *TypeAliasDeclaration) GetName() *DeclarationName { return node.Name }
func (node *TypeAliasDeclaration) GetTypeParameters() *TypeParameterListNode {
	return node.TypeParameters
}

func IsTypeAliasDeclaration(node *Node) bool {
	return node.Kind == SyntaxKindTypeAliasDeclaration
}

// EnumMember

type EnumMember struct {
	NodeBase
	NamedMemberBase
	Initializer *Node // Optional
}

func (f *NodeFactory) NewEnumMember(name *PropertyName, initializer *Node) *Node {
	data := &EnumMember{}
	data.Name = name
	data.Initializer = initializer
	return f.NewNode(SyntaxKindEnumMember, data)
}

func (node *EnumMember) ForEachChild(v Visitor) bool {
	return visit(v, node.Name) || visit(v, node.Initializer)
}

func (node *EnumMember) GetName() *DeclarationName {
	return node.Name
}

// EnumDeclaration

type EnumDeclaration struct {
	StatementBase
	DeclarationBase
	ExportableBase
	ModifiersBase
	Name    *IdentifierNode
	Members []*EnumMemberNode
}

func (f *NodeFactory) NewEnumDeclaration(modifiers *ModifierListNode, name *IdentifierNode, members []*EnumMemberNode) *Node {
	data := &EnumDeclaration{}
	data.Modifiers = modifiers
	data.Name = name
	data.Members = members
	return f.NewNode(SyntaxKindEnumDeclaration, data)
}

func (node *EnumDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.Modifiers) || visit(v, node.Name) || visitNodes(v, node.Members)
}

func (node *EnumDeclaration) GetName() *DeclarationName {
	return node.Name
}

func IsEnumDeclaration(node *Node) bool {
	return node.Kind == SyntaxKindEnumDeclaration
}

// ModuleBlock

type ModuleBlock struct {
	StatementBase
	Statements []*Statement
}

func (f *NodeFactory) NewModuleBlock(statements []*Statement) *Node {
	data := &ModuleBlock{}
	data.Statements = statements
	return f.NewNode(SyntaxKindModuleBlock, data)
}

func (node *ModuleBlock) ForEachChild(v Visitor) bool {
	return visitNodes(v, node.Statements)
}

func IsModuleBlock(node *Node) bool {
	return node.Kind == SyntaxKindModuleBlock
}

// ModuleDeclaration

type ModuleDeclaration struct {
	StatementBase
	DeclarationBase
	ExportableBase
	ModifiersBase
	LocalsContainerBase
	Name *ModuleName
	Body *ModuleBody // Optional (may be nil in ambient module declaration)
}

func (f *NodeFactory) NewModuleDeclaration(modifiers *ModifierListNode, name *ModuleName, body *ModuleBody, flags NodeFlags) *Node {
	data := &ModuleDeclaration{}
	data.Modifiers = modifiers
	data.Name = name
	data.Body = body
	node := f.NewNode(SyntaxKindModuleDeclaration, data)
	node.Flags |= flags & (NodeFlagsNamespace | NodeFlagsNestedNamespace | NodeFlagsGlobalAugmentation)
	return node
}

func (node *ModuleDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.Modifiers) || visit(v, node.Name) || visit(v, node.Body)
}

func (node *ModuleDeclaration) GetName() *DeclarationName {
	return node.Name
}

func IsModuleDeclaration(node *Node) bool {
	return node.Kind == SyntaxKindModuleDeclaration
}

// ModuleEqualsDeclaration

type ImportEqualsDeclaration struct {
	StatementBase
	DeclarationBase
	ExportableBase
	ModifiersBase
	Modifiers  *ModifierListNode
	IsTypeOnly bool
	Name       *IdentifierNode
	// 'EntityName' for an internal module reference, 'ExternalModuleReference' for an external
	// module reference.
	ModuleReference *ModuleReference
}

func (f *NodeFactory) NewImportEqualsDeclaration(modifiers *ModifierListNode, isTypeOnly bool, name *IdentifierNode, moduleReference *ModuleReference) *Node {
	data := &ImportEqualsDeclaration{}
	data.Modifiers = modifiers
	data.IsTypeOnly = isTypeOnly
	data.Name = name
	data.ModuleReference = moduleReference
	return f.NewNode(SyntaxKindImportEqualsDeclaration, data)
}

func (node *ImportEqualsDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.Modifiers) || visit(v, node.Name) || visit(v, node.ModuleReference)
}

func (node *ImportEqualsDeclaration) GetName() *DeclarationName {
	return node.Name
}

func IsImportEqualsDeclaration(node *Node) bool {
	return node.Kind == SyntaxKindImportEqualsDeclaration
}

// ImportDeclaration

type ImportDeclaration struct {
	StatementBase
	ModifiersBase
	ImportClause    *ImportClauseNode
	ModuleSpecifier *Node
	Attributes      *ImportAttributesNode
}

func (f *NodeFactory) NewImportDeclaration(modifiers *ModifierListNode, importClause *ImportClauseNode, moduleSpecifier *Node, attributes *ImportAttributesNode) *Node {
	data := &ImportDeclaration{}
	data.Modifiers = modifiers
	data.ImportClause = importClause
	data.ModuleSpecifier = moduleSpecifier
	data.Attributes = attributes
	return f.NewNode(SyntaxKindImportDeclaration, data)
}

func (node *ImportDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.Modifiers) || visit(v, node.ImportClause) || visit(v, node.ModuleSpecifier) || visit(v, node.Attributes)
}

func IsImportDeclaration(node *Node) bool {
	return node.Kind == SyntaxKindImportDeclaration
}

// ImportSpecifier

type ImportSpecifier struct {
	NodeBase
	DeclarationBase
	ExportableBase
	IsTypeOnly   bool
	PropertyName *ModuleExportName
	Name         *IdentifierNode
}

func (f *NodeFactory) NewImportSpecifier(isTypeOnly bool, propertyName *ModuleExportName, name *IdentifierNode) *Node {
	data := &ImportSpecifier{}
	data.IsTypeOnly = isTypeOnly
	data.PropertyName = propertyName
	data.Name = name
	return f.NewNode(SyntaxKindImportSpecifier, data)
}

func (node *ImportSpecifier) ForEachChild(v Visitor) bool {
	return visit(v, node.PropertyName) || visit(v, node.Name)
}

func (node *ImportSpecifier) GetName() *DeclarationName {
	return node.Name
}

func IsImportSpecifier(node *Node) bool {
	return node.Kind == SyntaxKindImportSpecifier
}

// ExternalModuleReference

type ExternalModuleReference struct {
	NodeBase
	Expression *Node
}

func (f *NodeFactory) NewExternalModuleReference(expression *Node) *Node {
	data := &ExternalModuleReference{}
	data.Expression = expression
	return f.NewNode(SyntaxKindExternalModuleReference, data)
}

func (node *ExternalModuleReference) ForEachChild(v Visitor) bool {
	return visit(v, node.Expression)
}

func IsExternalModuleReference(node *Node) bool {
	return node.Kind == SyntaxKindExternalModuleReference
}

// ImportClause

type ImportClause struct {
	NodeBase
	DeclarationBase
	ExportableBase
	IsTypeOnly    bool
	NamedBindings *NamedImportBindings // Optional, named bindings
	Name          *IdentifierNode      // Optional, default binding
}

func (f *NodeFactory) NewImportClause(isTypeOnly bool, name *IdentifierNode, namedBindings *NamedImportBindings) *Node {
	data := &ImportClause{}
	data.IsTypeOnly = isTypeOnly
	data.Name = name
	data.NamedBindings = namedBindings
	return f.NewNode(SyntaxKindImportClause, data)
}

func (node *ImportClause) ForEachChild(v Visitor) bool {
	return visit(v, node.Name) || visit(v, node.NamedBindings)
}

func (node *ImportClause) GetName() *DeclarationName {
	return node.Name
}

// NamespaceImport

type NamespaceImport struct {
	NodeBase
	DeclarationBase
	ExportableBase
	Name *IdentifierNode
}

func (f *NodeFactory) NewNamespaceImport(name *IdentifierNode) *Node {
	data := &NamespaceImport{}
	data.Name = name
	return f.NewNode(SyntaxKindNamespaceImport, data)
}

func (node *NamespaceImport) ForEachChild(v Visitor) bool {
	return visit(v, node.Name)
}

func (node *NamespaceImport) GetName() *DeclarationName {
	return node.Name
}

func IsNamespaceImport(node *Node) bool {
	return node.Kind == SyntaxKindNamespaceImport
}

// NamedImports

type NamedImports struct {
	NodeBase
	Elements []*ImportSpecifierNode
}

func (f *NodeFactory) NewNamedImports(elements []*ImportSpecifierNode) *Node {
	data := &NamedImports{}
	data.Elements = elements
	return f.NewNode(SyntaxKindNamedImports, data)
}

func (node *NamedImports) ForEachChild(v Visitor) bool {
	return visitNodes(v, node.Elements)
}

// ExportAssignment

// This is either an `export =` or an `export default` declaration.
// Unless `isExportEquals` is set, this node was parsed as an `export default`.
type ExportAssignment struct {
	StatementBase
	DeclarationBase
	ModifiersBase
	IsExportEquals bool
	Expression     *Node
}

func (f *NodeFactory) NewExportAssignment(modifiers *ModifierListNode, isExportEquals bool, expression *Node) *Node {
	data := &ExportAssignment{}
	data.Modifiers = modifiers
	data.IsExportEquals = isExportEquals
	data.Expression = expression
	return f.NewNode(SyntaxKindExportAssignment, data)
}

func (node *ExportAssignment) ForEachChild(v Visitor) bool {
	return visit(v, node.Modifiers) || visit(v, node.Expression)
}

func IsExportAssignment(node *Node) bool {
	return node.Kind == SyntaxKindExportAssignment
}

// NamespaceExportDeclaration

type NamespaceExportDeclaration struct {
	StatementBase
	DeclarationBase
	ModifiersBase
	Name *IdentifierNode
}

func (f *NodeFactory) NewNamespaceExportDeclaration(modifiers *ModifierListNode, name *IdentifierNode) *Node {
	data := &NamespaceExportDeclaration{}
	data.Modifiers = modifiers
	data.Name = name
	return f.NewNode(SyntaxKindNamespaceExportDeclaration, data)
}

func (node *NamespaceExportDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.Modifiers) || visit(v, node.Name)
}

func (node *NamespaceExportDeclaration) GetName() *DeclarationName {
	return node.Name
}

func IsNamespaceExportDeclaration(node *Node) bool {
	return node.Kind == SyntaxKindNamespaceExportDeclaration
}

// ExportDeclaration

type ExportDeclaration struct {
	StatementBase
	DeclarationBase
	ModifiersBase
	IsTypeOnly      bool
	ExportClause    *NamedExportBindings  // Optional
	ModuleSpecifier *Node                 // Optional
	Attributes      *ImportAttributesNode // Optional
}

func (f *NodeFactory) NewExportDeclaration(modifiers *ModifierListNode, isTypeOnly bool, exportClause *NamedExportBindings, moduleSpecifier *Node, attributes *ImportAttributesNode) *Node {
	data := &ExportDeclaration{}
	data.Modifiers = modifiers
	data.IsTypeOnly = isTypeOnly
	data.ExportClause = exportClause
	data.ModuleSpecifier = moduleSpecifier
	data.Attributes = attributes
	return f.NewNode(SyntaxKindExportDeclaration, data)
}

func (node *ExportDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.Modifiers) || visit(v, node.ExportClause) || visit(v, node.ModuleSpecifier) || visit(v, node.Attributes)
}

func IsExportDeclaration(node *Node) bool {
	return node.Kind == SyntaxKindExportDeclaration
}

// NamespaceExport

type NamespaceExport struct {
	NodeBase
	DeclarationBase
	Name *ModuleExportName
}

func (f *NodeFactory) NewNamespaceExport(name *ModuleExportName) *Node {
	data := &NamespaceExport{}
	data.Name = name
	return f.NewNode(SyntaxKindNamespaceExport, data)
}

func (node *NamespaceExport) ForEachChild(v Visitor) bool {
	return visit(v, node.Name)
}

func (node *NamespaceExport) GetName() *DeclarationName {
	return node.Name
}

func IsNamespaceExport(node *Node) bool {
	return node.Kind == SyntaxKindNamespaceExport
}

// NamedExports

type NamedExports struct {
	NodeBase
	Elements []*ExportSpecifierNode
}

func (f *NodeFactory) NewNamedExports(elements []*ExportSpecifierNode) *Node {
	data := &NamedExports{}
	data.Elements = elements
	return f.NewNode(SyntaxKindNamedExports, data)
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
	PropertyName *ModuleExportName // Optional, name preceding 'as' keyword
	Name         *ModuleExportName
}

func (f *NodeFactory) NewExportSpecifier(isTypeOnly bool, propertyName *ModuleExportName, name *ModuleExportName) *Node {
	data := &ExportSpecifier{}
	data.IsTypeOnly = isTypeOnly
	data.PropertyName = propertyName
	data.Name = name
	return f.NewNode(SyntaxKindExportSpecifier, data)
}

func (node *ExportSpecifier) ForEachChild(v Visitor) bool {
	return visit(v, node.PropertyName) || visit(v, node.Name)
}

func (node *ExportSpecifier) GetName() *DeclarationName {
	return node.Name
}

func IsExportSpecifier(node *Node) bool {
	return node.Kind == SyntaxKindExportSpecifier
}

// TypeElementBase

type TypeElementBase struct{}

// ClassElementBase

type ClassElementBase struct{}

// NamedMemberBase

type NamedMemberBase struct {
	DeclarationBase
	ModifiersBase
	Name         *PropertyName
	PostfixToken *TokenNode
}

func (node *NamedMemberBase) DeclarationData() *DeclarationBase { return &node.DeclarationBase }
func (node *NamedMemberBase) GetModifiers() *ModifierListNode   { return node.Modifiers }
func (node *NamedMemberBase) GetName() *DeclarationName         { return node.Name }

// CallSignatureDeclaration

type CallSignatureDeclaration struct {
	NodeBase
	DeclarationBase
	FunctionLikeBase
	TypeElementBase
}

func (f *NodeFactory) NewCallSignatureDeclaration(typeParameters *TypeParameterListNode, parameters []*ParameterDeclarationNode, returnType *Node) *Node {
	data := &CallSignatureDeclaration{}
	data.TypeParameters = typeParameters
	data.Parameters = parameters
	data.ReturnType = returnType
	return f.NewNode(SyntaxKindCallSignature, data)
}

func (node *CallSignatureDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.TypeParameters) || visitNodes(v, node.Parameters) || visit(v, node.ReturnType)
}

func IsCallSignatureDeclaration(node *Node) bool {
	return node.Kind == SyntaxKindCallSignature
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
	data.TypeParameters = typeParameters
	data.Parameters = parameters
	data.ReturnType = returnType
	return f.NewNode(SyntaxKindConstructSignature, data)
}

func (node *ConstructSignatureDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.TypeParameters) || visitNodes(v, node.Parameters) || visit(v, node.ReturnType)
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

func (f *NodeFactory) NewConstructorDeclaration(modifiers *ModifierListNode, typeParameters *TypeParameterListNode, parameters []*ParameterDeclarationNode, returnType *Node, body *BlockNode) *Node {
	data := &ConstructorDeclaration{}
	data.Modifiers = modifiers
	data.TypeParameters = typeParameters
	data.Parameters = parameters
	data.ReturnType = returnType
	data.Body = body
	return f.NewNode(SyntaxKindConstructor, data)
}

func (node *ConstructorDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.Modifiers) || visit(v, node.TypeParameters) || visitNodes(v, node.Parameters) || visit(v, node.ReturnType) || visit(v, node.Body)
}

func IsConstructorDeclaration(node *Node) bool {
	return node.Kind == SyntaxKindConstructor
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
	return visit(v, node.Modifiers) || visit(v, node.Name) || visit(v, node.TypeParameters) || visitNodes(v, node.Parameters) ||
		visit(v, node.ReturnType) || visit(v, node.Body)
}

func (node *AccessorDeclarationBase) isAccessorDeclaration() {}

// GetAccessorDeclaration

type GetAccessorDeclaration struct {
	AccessorDeclarationBase
}

func (f *NodeFactory) NewGetAccessorDeclaration(modifiers *ModifierListNode, name *PropertyName, typeParameters *TypeParameterListNode, parameters []*ParameterDeclarationNode, returnType *Node, body *BlockNode) *Node {
	data := &GetAccessorDeclaration{}
	data.Modifiers = modifiers
	data.Name = name
	data.TypeParameters = typeParameters
	data.Parameters = parameters
	data.ReturnType = returnType
	data.Body = body
	return f.NewNode(SyntaxKindGetAccessor, data)
}

func IsGetAccessorDeclaration(node *Node) bool {
	return node.Kind == SyntaxKindGetAccessor
}

// SetAccessorDeclaration

type SetAccessorDeclaration struct {
	AccessorDeclarationBase
}

func (f *NodeFactory) NewSetAccessorDeclaration(modifiers *ModifierListNode, name *PropertyName, typeParameters *TypeParameterListNode, parameters []*ParameterDeclarationNode, returnType *Node, body *BlockNode) *Node {
	data := &SetAccessorDeclaration{}
	data.Modifiers = modifiers
	data.Name = name
	data.TypeParameters = typeParameters
	data.Parameters = parameters
	data.ReturnType = returnType
	data.Body = body
	return f.NewNode(SyntaxKindSetAccessor, data)
}

func IsSetAccessorDeclaration(node *Node) bool {
	return node.Kind == SyntaxKindSetAccessor
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
	data.Modifiers = modifiers
	data.Parameters = parameters
	data.ReturnType = returnType
	return f.NewNode(SyntaxKindIndexSignature, data)
}

func (node *IndexSignatureDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.Modifiers) || visitNodes(v, node.Parameters) || visit(v, node.ReturnType)
}

func IsIndexSignatureDeclaration(node *Node) bool {
	return node.Kind == SyntaxKindIndexSignature
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
	data.Modifiers = modifiers
	data.Name = name
	data.PostfixToken = postfixToken
	data.TypeParameters = typeParameters
	data.Parameters = parameters
	data.ReturnType = returnType
	return f.NewNode(SyntaxKindMethodSignature, data)
}

func (node *MethodSignatureDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.Modifiers) || visit(v, node.Name) || visit(v, node.PostfixToken) || visit(v, node.TypeParameters) ||
		visitNodes(v, node.Parameters) || visit(v, node.ReturnType)
}

func IsMethodSignatureDeclaration(node *Node) bool {
	return node.Kind == SyntaxKindMethodSignature
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
	data.Modifiers = modifiers
	data.AsteriskToken = asteriskToken
	data.Name = name
	data.PostfixToken = postfixToken
	data.TypeParameters = typeParameters
	data.Parameters = parameters
	data.ReturnType = returnType
	data.Body = body
	return f.NewNode(SyntaxKindMethodDeclaration, data)
}

func (node *MethodDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.Modifiers) || visit(v, node.AsteriskToken) || visit(v, node.Name) || visit(v, node.PostfixToken) ||
		visit(v, node.TypeParameters) || visitNodes(v, node.Parameters) || visit(v, node.ReturnType) || visit(v, node.Body)
}

func IsMethodDeclaration(node *Node) bool {
	return node.Kind == SyntaxKindMethodDeclaration
}

// PropertySignatureDeclaration

type PropertySignatureDeclaration struct {
	NodeBase
	NamedMemberBase
	TypeElementBase
	TypeNode    *Node
	Initializer *Node // For error reporting purposes
}

func (f *NodeFactory) NewPropertySignatureDeclaration(modifiers *Node, name *Node, postfixToken *Node, typeNode *Node, initializer *Node) *Node {
	data := &PropertySignatureDeclaration{}
	data.Modifiers = modifiers
	data.Name = name
	data.PostfixToken = postfixToken
	data.TypeNode = typeNode
	data.Initializer = initializer
	return f.NewNode(SyntaxKindPropertySignature, data)
}

func (node *PropertySignatureDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.Modifiers) || visit(v, node.Name) || visit(v, node.PostfixToken) || visit(v, node.TypeNode) || visit(v, node.Initializer)
}

func IsPropertySignatureDeclaration(node *Node) bool {
	return node.Kind == SyntaxKindPropertySignature
}

// PropertyDeclaration

type PropertyDeclaration struct {
	NodeBase
	NamedMemberBase
	ClassElementBase
	TypeNode    *Node // Optional
	Initializer *Node // Optional
}

func (f *NodeFactory) NewPropertyDeclaration(modifiers *Node, name *Node, postfixToken *Node, typeNode *Node, initializer *Node) *Node {
	data := &PropertyDeclaration{}
	data.Modifiers = modifiers
	data.Name = name
	data.PostfixToken = postfixToken
	data.TypeNode = typeNode
	data.Initializer = initializer
	return f.NewNode(SyntaxKindPropertyDeclaration, data)
}

func (node *PropertyDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.Modifiers) || visit(v, node.Name) || visit(v, node.PostfixToken) || visit(v, node.TypeNode) || visit(v, node.Initializer)
}

func IsPropertyDeclaration(node *Node) bool {
	return node.Kind == SyntaxKindPropertyDeclaration
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
	Body *Node
}

func (f *NodeFactory) NewClassStaticBlockDeclaration(modifiers *Node, body *Node) *Node {
	data := &ClassStaticBlockDeclaration{}
	data.Modifiers = modifiers
	data.Body = body
	return f.NewNode(SyntaxKindClassStaticBlockDeclaration, data)
}

func (node *ClassStaticBlockDeclaration) ForEachChild(v Visitor) bool {
	return visit(v, node.Modifiers) || visit(v, node.Body)
}

func IsClassStaticBlockDeclaration(node *Node) bool {
	return node.Kind == SyntaxKindClassStaticBlockDeclaration
}

// TypeParameterList

type TypeParameterList struct {
	NodeBase
	Parameters []*Node
}

func (f *NodeFactory) NewTypeParameterList(parameters []*Node) *Node {
	data := &TypeParameterList{}
	data.Parameters = parameters
	return f.NewNode(SyntaxKindTypeParameterList, data)
}

func (node *TypeParameterList) ForEachChild(v Visitor) bool {
	return visitNodes(v, node.Parameters)
}

func IsTypeParameterList(node *Node) bool {
	return node.Kind == SyntaxKindTypeParameterList
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
	return f.NewNode(SyntaxKindStringLiteral, data)
}

func IsStringLiteral(node *Node) bool {
	return node.Kind == SyntaxKindStringLiteral
}

// NumericLiteral

type NumericLiteral struct {
	ExpressionBase
	LiteralLikeBase
}

func (f *NodeFactory) NewNumericLiteral(text string) *Node {
	data := &NumericLiteral{}
	data.Text = text
	return f.NewNode(SyntaxKindNumericLiteral, data)
}

// BigIntLiteral

type BigIntLiteral struct {
	ExpressionBase
	LiteralLikeBase
}

func (f *NodeFactory) NewBigIntLiteral(text string) *Node {
	data := &BigIntLiteral{}
	data.Text = text
	return f.NewNode(SyntaxKindBigIntLiteral, data)
}

// RegularExpressionLiteral

type RegularExpressionLiteral struct {
	ExpressionBase
	LiteralLikeBase
}

func (f *NodeFactory) NewRegularExpressionLiteral(text string) *Node {
	data := &RegularExpressionLiteral{}
	data.Text = text
	return f.NewNode(SyntaxKindRegularExpressionLiteral, data)
}

// NoSubstitutionTemplateLiteral

type NoSubstitutionTemplateLiteral struct {
	ExpressionBase
	TemplateLiteralLikeBase
}

func (f *NodeFactory) NewNoSubstitutionTemplateLiteral(text string) *Node {
	data := &NoSubstitutionTemplateLiteral{}
	data.Text = text
	return f.NewNode(SyntaxKindNoSubstitutionTemplateLiteral, data)
}

// BinaryExpression

type BinaryExpression struct {
	ExpressionBase
	DeclarationBase
	Left          *Node
	OperatorToken *Node
	Right         *Node
}

func (f *NodeFactory) NewBinaryExpression(left *Node, operatorToken *Node, right *Node) *Node {
	data := &BinaryExpression{}
	data.Left = left
	data.OperatorToken = operatorToken
	data.Right = right
	return f.NewNode(SyntaxKindBinaryExpression, data)
}

func (node *BinaryExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.Left) || visit(v, node.OperatorToken) || visit(v, node.Right)
}

// PrefixUnaryExpression

type PrefixUnaryExpression struct {
	ExpressionBase
	Operator SyntaxKind
	Operand  *Node
}

func (f *NodeFactory) NewPrefixUnaryExpression(operator SyntaxKind, operand *Node) *Node {
	data := &PrefixUnaryExpression{}
	data.Operator = operator
	data.Operand = operand
	return f.NewNode(SyntaxKindPrefixUnaryExpression, data)
}

func (node *PrefixUnaryExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.Operand)
}

func IsPrefixUnaryExpression(node *Node) bool {
	return node.Kind == SyntaxKindPrefixUnaryExpression
}

// PostfixUnaryExpression

type PostfixUnaryExpression struct {
	ExpressionBase
	Operand  *Node
	Operator SyntaxKind
}

func (f *NodeFactory) NewPostfixUnaryExpression(operand *Node, operator SyntaxKind) *Node {
	data := &PostfixUnaryExpression{}
	data.Operand = operand
	data.Operator = operator
	return f.NewNode(SyntaxKindPostfixUnaryExpression, data)
}

func (node *PostfixUnaryExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.Operand)
}

// YieldExpression

type YieldExpression struct {
	ExpressionBase
	AsteriskToken *Node
	Expression    *Node // Optional
}

func (f *NodeFactory) NewYieldExpression(asteriskToken *Node, expression *Node) *Node {
	data := &YieldExpression{}
	data.AsteriskToken = asteriskToken
	data.Expression = expression
	return f.NewNode(SyntaxKindYieldExpression, data)
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
	EqualsGreaterThanToken *Node
}

func (f *NodeFactory) NewArrowFunction(modifiers *Node, typeParameters *Node, parameters []*Node, returnType *Node, equalsGreaterThanToken *Node, body *Node) *Node {
	data := &ArrowFunction{}
	data.Modifiers = modifiers
	data.TypeParameters = typeParameters
	data.Parameters = parameters
	data.ReturnType = returnType
	data.EqualsGreaterThanToken = equalsGreaterThanToken
	data.Body = body
	return f.NewNode(SyntaxKindArrowFunction, data)
}

func (node *ArrowFunction) ForEachChild(v Visitor) bool {
	return visit(v, node.Modifiers) || visit(v, node.TypeParameters) || visitNodes(v, node.Parameters) ||
		visit(v, node.ReturnType) || visit(v, node.EqualsGreaterThanToken) || visit(v, node.Body)
}

func (node *ArrowFunction) GetName() *DeclarationName {
	return nil
}

func IsArrowFunction(node *Node) bool {
	return node.Kind == SyntaxKindArrowFunction
}

// FunctionExpression

type FunctionExpression struct {
	ExpressionBase
	DeclarationBase
	ModifiersBase
	FunctionLikeWithBodyBase
	FlowNodeBase
	Name           *Node // Optional
	ReturnFlowNode *FlowNode
}

func (f *NodeFactory) NewFunctionExpression(modifiers *Node, asteriskToken *Node, name *Node, typeParameters *Node, parameters []*Node, returnType *Node, body *Node) *Node {
	data := &FunctionExpression{}
	data.Modifiers = modifiers
	data.AsteriskToken = asteriskToken
	data.Name = name
	data.TypeParameters = typeParameters
	data.Parameters = parameters
	data.ReturnType = returnType
	data.Body = body
	return f.NewNode(SyntaxKindFunctionExpression, data)
}

func (node *FunctionExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.Modifiers) || visit(v, node.AsteriskToken) || visit(v, node.Name) || visit(v, node.TypeParameters) ||
		visitNodes(v, node.Parameters) || visit(v, node.ReturnType) || visit(v, node.Body)
}

func (node *FunctionExpression) GetName() *DeclarationName {
	return node.Name
}

func IsFunctionExpression(node *Node) bool {
	return node.Kind == SyntaxKindFunctionExpression
}

// AsExpression

type AsExpression struct {
	ExpressionBase
	Expression *Node
	TypeNode   *Node
}

func (f *NodeFactory) NewAsExpression(expression *Node, typeNode *Node) *Node {
	data := &AsExpression{}
	data.Expression = expression
	data.TypeNode = typeNode
	return f.NewNode(SyntaxKindAsExpression, data)
}

func (node *AsExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.Expression) || visit(v, node.TypeNode)
}

// SatisfiesExpression

type SatisfiesExpression struct {
	ExpressionBase
	Expression *Node
	TypeNode   *Node
}

func (f *NodeFactory) NewSatisfiesExpression(expression *Node, typeNode *Node) *Node {
	data := &SatisfiesExpression{}
	data.Expression = expression
	data.TypeNode = typeNode
	return f.NewNode(SyntaxKindSatisfiesExpression, data)
}

func (node *SatisfiesExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.Expression) || visit(v, node.TypeNode)
}

// ConditionalExpression

type ConditionalExpression struct {
	ExpressionBase
	Condition     *Node
	QuestionToken *Node
	WhenTrue      *Node
	ColonToken    *Node
	WhenFalse     *Node
}

func (f *NodeFactory) NewConditionalExpression(condition *Node, questionToken *Node, whenTrue *Node, colonToken *Node, whenFalse *Node) *Node {
	data := &ConditionalExpression{}
	data.Condition = condition
	data.QuestionToken = questionToken
	data.WhenTrue = whenTrue
	data.ColonToken = colonToken
	data.WhenFalse = whenFalse
	return f.NewNode(SyntaxKindConditionalExpression, data)
}

func (node *ConditionalExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.Condition) || visit(v, node.QuestionToken) || visit(v, node.WhenTrue) ||
		visit(v, node.ColonToken) || visit(v, node.WhenFalse)
}

// PropertyAccessExpression

type PropertyAccessExpression struct {
	ExpressionBase
	FlowNodeBase
	Expression       *Node
	QuestionDotToken *Node
	Name             *Node
}

func (f *NodeFactory) NewPropertyAccessExpression(expression *Node, questionDotToken *Node, name *Node, flags NodeFlags) *Node {
	data := &PropertyAccessExpression{}
	data.Expression = expression
	data.QuestionDotToken = questionDotToken
	data.Name = name
	node := f.NewNode(SyntaxKindPropertyAccessExpression, data)
	node.Flags |= flags & NodeFlagsOptionalChain
	return node
}

func (node *PropertyAccessExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.Expression) || visit(v, node.QuestionDotToken) || visit(v, node.Name)
}

func (node *PropertyAccessExpression) GetName() *DeclarationName { return node.Name }

func IsPropertyAccessExpression(node *Node) bool {
	return node.Kind == SyntaxKindPropertyAccessExpression
}

// ElementAccessExpression

type ElementAccessExpression struct {
	ExpressionBase
	FlowNodeBase
	Expression         *Node
	QuestionDotToken   *Node
	ArgumentExpression *Node
}

func (f *NodeFactory) NewElementAccessExpression(expression *Node, questionDotToken *Node, argumentExpression *Node, flags NodeFlags) *Node {
	data := &ElementAccessExpression{}
	data.Expression = expression
	data.QuestionDotToken = questionDotToken
	data.ArgumentExpression = argumentExpression
	node := f.NewNode(SyntaxKindElementAccessExpression, data)
	node.Flags |= flags & NodeFlagsOptionalChain
	return node
}

func (node *ElementAccessExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.Expression) || visit(v, node.QuestionDotToken) || visit(v, node.ArgumentExpression)
}

func IsElementAccessExpression(node *Node) bool {
	return node.Kind == SyntaxKindElementAccessExpression
}

// CallExpression

type CallExpression struct {
	ExpressionBase
	Expression       *Node
	QuestionDotToken *Node
	TypeArguments    *Node
	Arguments        []*Node
}

func (f *NodeFactory) NewCallExpression(expression *Node, questionDotToken *Node, typeArguments *Node, arguments []*Node, flags NodeFlags) *Node {
	data := &CallExpression{}
	data.Expression = expression
	data.QuestionDotToken = questionDotToken
	data.TypeArguments = typeArguments
	data.Arguments = arguments
	node := f.NewNode(SyntaxKindCallExpression, data)
	node.Flags |= flags & NodeFlagsOptionalChain
	return node
}

func (node *CallExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.Expression) || visit(v, node.QuestionDotToken) || visit(v, node.TypeArguments) || visitNodes(v, node.Arguments)
}

func IsCallExpression(node *Node) bool {
	return node.Kind == SyntaxKindCallExpression
}

// NewExpression

type NewExpression struct {
	ExpressionBase
	Expression    *Node
	TypeArguments *Node
	Arguments     []*Node
}

func (f *NodeFactory) NewNewExpression(expression *Node, typeArguments *Node, arguments []*Node) *Node {
	data := &NewExpression{}
	data.Expression = expression
	data.TypeArguments = typeArguments
	data.Arguments = arguments
	return f.NewNode(SyntaxKindNewExpression, data)
}

func (node *NewExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.Expression) || visit(v, node.TypeArguments) || visitNodes(v, node.Arguments)
}

func IsNewExpression(node *Node) bool {
	return node.Kind == SyntaxKindNewExpression
}

// MetaProperty

type MetaProperty struct {
	ExpressionBase
	FlowNodeBase
	KeywordToken SyntaxKind
	Name         *Node
}

func (f *NodeFactory) NewMetaProperty(keywordToken SyntaxKind, name *Node) *Node {
	data := &MetaProperty{}
	data.KeywordToken = keywordToken
	data.Name = name
	return f.NewNode(SyntaxKindMetaProperty, data)
}

func (node *MetaProperty) ForEachChild(v Visitor) bool {
	return visit(v, node.Name)
}

func IsMetaProperty(node *Node) bool {
	return node.Kind == SyntaxKindMetaProperty
}

// NonNullExpression

type NonNullExpression struct {
	ExpressionBase
	Expression *Node
}

func (f *NodeFactory) NewNonNullExpression(expression *Node) *Node {
	data := &NonNullExpression{}
	data.Expression = expression
	return f.NewNode(SyntaxKindNonNullExpression, data)
}

func (node *NonNullExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.Expression)
}

// SpreadElement

type SpreadElement struct {
	ExpressionBase
	Expression *Node
}

func (f *NodeFactory) NewSpreadElement(expression *Node) *Node {
	data := &SpreadElement{}
	data.Expression = expression
	return f.NewNode(SyntaxKindSpreadElement, data)
}

func (node *SpreadElement) ForEachChild(v Visitor) bool {
	return visit(v, node.Expression)
}

func IsSpreadElement(node *Node) bool {
	return node.Kind == SyntaxKindSpreadElement
}

// TemplateExpression

type TemplateExpression struct {
	ExpressionBase
	Head          *Node
	TemplateSpans []*Node
}

func (f *NodeFactory) NewTemplateExpression(head *Node, templateSpans []*Node) *Node {
	data := &TemplateExpression{}
	data.Head = head
	data.TemplateSpans = templateSpans
	return f.NewNode(SyntaxKindTemplateExpression, data)
}

func (node *TemplateExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.Head) || visitNodes(v, node.TemplateSpans)
}

// TemplateLiteralTypeSpan

type TemplateSpan struct {
	NodeBase
	Expression *Node
	Literal    *Node
}

func (f *NodeFactory) NewTemplateSpan(expression *Node, literal *Node) *Node {
	data := &TemplateSpan{}
	data.Expression = expression
	data.Literal = literal
	return f.NewNode(SyntaxKindTemplateSpan, data)
}

func (node *TemplateSpan) ForEachChild(v Visitor) bool {
	return visit(v, node.Expression) || visit(v, node.Literal)
}

func IsTemplateSpan(node *Node) bool {
	return node.Kind == SyntaxKindTemplateSpan
}

// TaggedTemplateExpression

type TaggedTemplateExpression struct {
	ExpressionBase
	Tag              *Node
	QuestionDotToken *Node // For error reporting purposes only
	TypeArguments    *Node
	Template         *Node
}

func (f *NodeFactory) NewTaggedTemplateExpression(tag *Node, questionDotToken *Node, typeArguments *Node, template *Node, flags NodeFlags) *Node {
	data := &TaggedTemplateExpression{}
	data.Tag = tag
	data.QuestionDotToken = questionDotToken
	data.TypeArguments = typeArguments
	data.Template = template
	node := f.NewNode(SyntaxKindTaggedTemplateExpression, data)
	node.Flags |= flags & NodeFlagsOptionalChain
	return node
}

func (node *TaggedTemplateExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.Tag) || visit(v, node.QuestionDotToken) || visit(v, node.TypeArguments) || visit(v, node.Template)
}

// ParenthesizedExpression

type ParenthesizedExpression struct {
	ExpressionBase
	Expression *Node
}

func (f *NodeFactory) NewParenthesizedExpression(expression *Node) *Node {
	data := &ParenthesizedExpression{}
	data.Expression = expression
	return f.NewNode(SyntaxKindParenthesizedExpression, data)
}

func (node *ParenthesizedExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.Expression)
}

func IsParenthesizedExpression(node *Node) bool {
	return node.Kind == SyntaxKindParenthesizedExpression
}

// ArrayLiteralExpression

type ArrayLiteralExpression struct {
	ExpressionBase
	Elements  []*Node
	MultiLine bool
}

func (f *NodeFactory) NewArrayLiteralExpression(elements []*Node, multiLine bool) *Node {
	data := &ArrayLiteralExpression{}
	data.Elements = elements
	data.MultiLine = multiLine
	return f.NewNode(SyntaxKindArrayLiteralExpression, data)
}

func (node *ArrayLiteralExpression) ForEachChild(v Visitor) bool {
	return visitNodes(v, node.Elements)
}

func IsArrayLiteralExpression(node *Node) bool {
	return node.Kind == SyntaxKindArrayLiteralExpression
}

// ObjectLiteralExpression

type ObjectLiteralExpression struct {
	ExpressionBase
	DeclarationBase
	Properties []*Node
	MultiLine  bool
}

func (f *NodeFactory) NewObjectLiteralExpression(properties []*Node, multiLine bool) *Node {
	data := &ObjectLiteralExpression{}
	data.Properties = properties
	data.MultiLine = multiLine
	return f.NewNode(SyntaxKindObjectLiteralExpression, data)

}

func (node *ObjectLiteralExpression) ForEachChild(v Visitor) bool {
	return visitNodes(v, node.Properties)
}

func IsObjectLiteralExpression(node *Node) bool {
	return node.Kind == SyntaxKindObjectLiteralExpression
}

// ObjectLiteralElementBase

type ObjectLiteralElementBase struct{}

// SpreadAssignment

type SpreadAssignment struct {
	NodeBase
	DeclarationBase
	ObjectLiteralElementBase
	Expression *Node
}

func (f *NodeFactory) NewSpreadAssignment(expression *Node) *Node {
	data := &SpreadAssignment{}
	data.Expression = expression
	return f.NewNode(SyntaxKindSpreadAssignment, data)
}

func (node *SpreadAssignment) ForEachChild(v Visitor) bool {
	return visit(v, node.Expression)
}

// PropertyAssignment

type PropertyAssignment struct {
	NodeBase
	NamedMemberBase
	ObjectLiteralElementBase
	Initializer *Node
}

func (f *NodeFactory) NewPropertyAssignment(modifiers *Node, name *Node, postfixToken *Node, initializer *Node) *Node {
	data := &PropertyAssignment{}
	data.Modifiers = modifiers
	data.Name = name
	data.PostfixToken = postfixToken
	data.Initializer = initializer
	return f.NewNode(SyntaxKindPropertyAssignment, data)
}

func (node *PropertyAssignment) ForEachChild(v Visitor) bool {
	return visit(v, node.Modifiers) || visit(v, node.Name) || visit(v, node.PostfixToken) || visit(v, node.Initializer)
}

func IsPropertyAssignment(node *Node) bool {
	return node.Kind == SyntaxKindPropertyAssignment
}

// ShorthandPropertyAssignment

type ShorthandPropertyAssignment struct {
	NodeBase
	NamedMemberBase
	ObjectLiteralElementBase
	ObjectAssignmentInitializer *Node // Optional
}

func (f *NodeFactory) NewShorthandPropertyAssignment(modifiers *Node, name *Node, postfixToken *Node, objectAssignmentInitializer *Node) *Node {
	data := &ShorthandPropertyAssignment{}
	data.Modifiers = modifiers
	data.Name = name
	data.PostfixToken = postfixToken
	data.ObjectAssignmentInitializer = objectAssignmentInitializer
	return f.NewNode(SyntaxKindShorthandPropertyAssignment, data)
}

func (node *ShorthandPropertyAssignment) ForEachChild(v Visitor) bool {
	return visit(v, node.Modifiers) || visit(v, node.Name) || visit(v, node.PostfixToken) || visit(v, node.ObjectAssignmentInitializer)
}

func IsShorthandPropertyAssignment(node *Node) bool {
	return node.Kind == SyntaxKindShorthandPropertyAssignment
}

// DeleteExpression

type DeleteExpression struct {
	ExpressionBase
	Expression *Node
}

func (f *NodeFactory) NewDeleteExpression(expression *Node) *Node {
	data := &DeleteExpression{}
	data.Expression = expression
	return f.NewNode(SyntaxKindDeleteExpression, data)

}

func (node *DeleteExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.Expression)
}

// TypeOfExpression

type TypeOfExpression struct {
	ExpressionBase
	Expression *Node
}

func (f *NodeFactory) NewTypeOfExpression(expression *Node) *Node {
	data := &TypeOfExpression{}
	data.Expression = expression
	return f.NewNode(SyntaxKindTypeOfExpression, data)
}

func (node *TypeOfExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.Expression)
}

func IsTypeOfExpression(node *Node) bool {
	return node.Kind == SyntaxKindTypeOfExpression
}

// VoidExpression

type VoidExpression struct {
	ExpressionBase
	Expression *Node
}

func (f *NodeFactory) NewVoidExpression(expression *Node) *Node {
	data := &VoidExpression{}
	data.Expression = expression
	return f.NewNode(SyntaxKindVoidExpression, data)
}

func (node *VoidExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.Expression)
}

// AwaitExpression

type AwaitExpression struct {
	ExpressionBase
	Expression *Node
}

func (f *NodeFactory) NewAwaitExpression(expression *Node) *Node {
	data := &AwaitExpression{}
	data.Expression = expression
	return f.NewNode(SyntaxKindAwaitExpression, data)
}

func (node *AwaitExpression) ForEachChild(v Visitor) bool {
	return visit(v, node.Expression)
}

// TypeAssertion

type TypeAssertion struct {
	ExpressionBase
	TypeNode   *Node
	Expression *Node
}

func (f *NodeFactory) NewTypeAssertion(typeNode *Node, expression *Node) *Node {
	data := &TypeAssertion{}
	data.TypeNode = typeNode
	data.Expression = expression
	return f.NewNode(SyntaxKindTypeAssertionExpression, data)
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

func (f *NodeFactory) NewKeywordTypeNode(kind SyntaxKind) *Node {
	return f.NewNode(kind, &KeywordTypeNode{})
}

// UnionOrIntersectionTypeBase

type UnionOrIntersectionTypeNodeBase struct {
	TypeNodeBase
	Types []*Node
}

func (node *UnionOrIntersectionTypeNodeBase) ForEachChild(v Visitor) bool {
	return visitNodes(v, node.Types)
}

// UnionTypeNode

type UnionTypeNode struct {
	UnionOrIntersectionTypeNodeBase
}

func (f *NodeFactory) NewUnionTypeNode(types []*Node) *Node {
	data := &UnionTypeNode{}
	data.Types = types
	return f.NewNode(SyntaxKindUnionType, data)
}

// IntersectionTypeNode

type IntersectionTypeNode struct {
	UnionOrIntersectionTypeNodeBase
}

func (f *NodeFactory) NewIntersectionTypeNode(types []*Node) *Node {
	data := &IntersectionTypeNode{}
	data.Types = types
	return f.NewNode(SyntaxKindIntersectionType, data)
}

// ConditionalTypeNode

type ConditionalTypeNode struct {
	TypeNodeBase
	LocalsContainerBase
	CheckType   *Node
	ExtendsType *Node
	TrueType    *Node
	FalseType   *Node
}

func (f *NodeFactory) NewConditionalTypeNode(checkType *Node, extendsType *Node, trueType *Node, falseType *Node) *Node {
	data := &ConditionalTypeNode{}
	data.CheckType = checkType
	data.ExtendsType = extendsType
	data.TrueType = trueType
	data.FalseType = falseType
	return f.NewNode(SyntaxKindConditionalType, data)
}

func (node *ConditionalTypeNode) ForEachChild(v Visitor) bool {
	return visit(v, node.CheckType) || visit(v, node.ExtendsType) || visit(v, node.TrueType) || visit(v, node.FalseType)
}

func IsConditionalTypeNode(node *Node) bool {
	return node.Kind == SyntaxKindConditionalType
}

// TypeOperatorNode

type TypeOperatorNode struct {
	TypeNodeBase
	Operator SyntaxKind // SyntaxKindKeyOfKeyword | SyntaxKindUniqueKeyword | SyntaxKindReadonlyKeyword
	TypeNode *Node
}

func (f *NodeFactory) NewTypeOperatorNode(operator SyntaxKind, typeNode *Node) *Node {
	data := &TypeOperatorNode{}
	data.Operator = operator
	data.TypeNode = typeNode
	return f.NewNode(SyntaxKindTypeOperator, data)
}

func (node *TypeOperatorNode) ForEachChild(v Visitor) bool {
	return visit(v, node.TypeNode)
}

func IsTypeOperatorNode(node *Node) bool {
	return node.Kind == SyntaxKindTypeOperator
}

// InferTypeNode

type InferTypeNode struct {
	TypeNodeBase
	TypeParameter *Node
}

func (f *NodeFactory) NewInferTypeNode(typeParameter *Node) *Node {
	data := &InferTypeNode{}
	data.TypeParameter = typeParameter
	return f.NewNode(SyntaxKindInferType, data)
}

func (node *InferTypeNode) ForEachChild(v Visitor) bool {
	return visit(v, node.TypeParameter)
}

// ArrayTypeNode

type ArrayTypeNode struct {
	TypeNodeBase
	ElementType *Node
}

func (f *NodeFactory) NewArrayTypeNode(elementType *Node) *Node {
	data := &ArrayTypeNode{}
	data.ElementType = elementType
	return f.NewNode(SyntaxKindArrayType, data)
}

func (node *ArrayTypeNode) ForEachChild(v Visitor) bool {
	return visit(v, node.ElementType)
}

// IndexedAccessTypeNode

type IndexedAccessTypeNode struct {
	TypeNodeBase
	ObjectType *Node
	IndexType  *Node
}

func (f *NodeFactory) NewIndexedAccessTypeNode(objectType *Node, indexType *Node) *Node {
	data := &IndexedAccessTypeNode{}
	data.ObjectType = objectType
	data.IndexType = indexType
	return f.NewNode(SyntaxKindIndexedAccessType, data)
}

func (node *IndexedAccessTypeNode) ForEachChild(v Visitor) bool {
	return visit(v, node.ObjectType) || visit(v, node.IndexType)
}

func IsIndexedAccessTypeNode(node *Node) bool {
	return node.Kind == SyntaxKindIndexedAccessType
}

// TypeArgumentList

type TypeArgumentList struct {
	NodeBase
	Arguments []*Node
}

func (f *NodeFactory) NewTypeArgumentList(arguments []*Node) *Node {
	data := &TypeArgumentList{}
	data.Arguments = arguments
	return f.NewNode(SyntaxKindTypeArgumentList, data)
}

func (node *TypeArgumentList) ForEachChild(v Visitor) bool {
	return visitNodes(v, node.Arguments)
}

// TypeReferenceNode

type TypeReferenceNode struct {
	TypeNodeBase
	TypeName      *Node
	TypeArguments *Node // TypeArgumentList
}

func (f *NodeFactory) NewTypeReferenceNode(typeName *Node, typeArguments *Node) *Node {
	data := &TypeReferenceNode{}
	data.TypeName = typeName
	data.TypeArguments = typeArguments
	return f.NewNode(SyntaxKindTypeReference, data)
}

func (node *TypeReferenceNode) ForEachChild(v Visitor) bool {
	return visit(v, node.TypeName) || visit(v, node.TypeArguments)
}

func IsTypeReferenceNode(node *Node) bool {
	return node.Kind == SyntaxKindTypeReference
}

// ExpressionWithTypeArguments

type ExpressionWithTypeArguments struct {
	ExpressionBase
	Expression    *Node
	TypeArguments *Node
}

func (f *NodeFactory) NewExpressionWithTypeArguments(expression *Node, typeArguments *Node) *Node {
	data := &ExpressionWithTypeArguments{}
	data.Expression = expression
	data.TypeArguments = typeArguments
	return f.NewNode(SyntaxKindExpressionWithTypeArguments, data)
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
	return f.NewNode(SyntaxKindLiteralType, data)
}

func (node *LiteralTypeNode) ForEachChild(v Visitor) bool {
	return visit(v, node.Literal)
}

func IsLiteralTypeNode(node *Node) bool {
	return node.Kind == SyntaxKindLiteralType
}

// ThisTypeNode

type ThisTypeNode struct {
	TypeNodeBase
}

func (f *NodeFactory) NewThisTypeNode() *Node {
	return f.NewNode(SyntaxKindThisType, &ThisTypeNode{})
}

// TypePredicateNode

type TypePredicateNode struct {
	TypeNodeBase
	AssertsModifier *Node // Optional
	ParameterName   *Node // Identifier | ThisTypeNode
	TypeNode        *Node // Optional
}

func (f *NodeFactory) NewTypePredicateNode(assertsModifier *Node, parameterName *Node, typeNode *Node) *Node {
	data := &TypePredicateNode{}
	data.AssertsModifier = assertsModifier
	data.ParameterName = parameterName
	data.TypeNode = typeNode
	return f.NewNode(SyntaxKindTypePredicate, data)
}

func (node *TypePredicateNode) ForEachChild(v Visitor) bool {
	return visit(v, node.AssertsModifier) || visit(v, node.ParameterName) || visit(v, node.TypeNode)
}

// ImportTypeNode

type ImportTypeNode struct {
	TypeNodeBase
	IsTypeOf      bool
	Argument      *Node
	Attributes    *Node // Optional
	Qualifier     *Node // Optional
	TypeArguments *Node // Optional
}

func (f *NodeFactory) NewImportTypeNode(isTypeOf bool, argument *Node, attributes *Node, qualifier *Node, typeArguments *Node) *Node {
	data := &ImportTypeNode{}
	data.IsTypeOf = isTypeOf
	data.Argument = argument
	data.Attributes = attributes
	data.Qualifier = qualifier
	data.TypeArguments = typeArguments
	return f.NewNode(SyntaxKindImportType, data)
}

func (node *ImportTypeNode) ForEachChild(v Visitor) bool {
	return visit(v, node.Argument) || visit(v, node.Attributes) || visit(v, node.Qualifier) || visit(v, node.TypeArguments)
}

func isImportTypeNode(node *Node) bool {
	return node.Kind == SyntaxKindImportType
}

// ImportAttribute

type ImportAttribute struct {
	NodeBase
	Name  *Node
	Value *Node
}

func (f *NodeFactory) NewImportAttribute(name *Node, value *Node) *Node {
	data := &ImportAttribute{}
	data.Name = name
	data.Value = value
	return f.NewNode(SyntaxKindImportAttribute, data)
}

func (node *ImportAttribute) ForEachChild(v Visitor) bool {
	return visit(v, node.Name) || visit(v, node.Value)
}

// ImportAttributes

type ImportAttributes struct {
	NodeBase
	Token      SyntaxKind
	Attributes []*Node
	MultiLine  bool
}

func (f *NodeFactory) NewImportAttributes(token SyntaxKind, attributes []*Node, multiLine bool) *Node {
	data := &ImportAttributes{}
	data.Token = token
	data.Attributes = attributes
	data.MultiLine = multiLine
	return f.NewNode(SyntaxKindImportAttributes, data)
}

func (node *ImportAttributes) ForEachChild(v Visitor) bool {
	return visitNodes(v, node.Attributes)
}

// TypeQueryNode

type TypeQueryNode struct {
	TypeNodeBase
	ExprName      *Node
	TypeArguments *Node
}

func (f *NodeFactory) NewTypeQueryNode(exprName *Node, typeArguments *Node) *Node {
	data := &TypeQueryNode{}
	data.ExprName = exprName
	data.TypeArguments = typeArguments
	return f.NewNode(SyntaxKindTypeQuery, data)
}

func (node *TypeQueryNode) ForEachChild(v Visitor) bool {
	return visit(v, node.ExprName) || visit(v, node.TypeArguments)
}

func IsTypeQueryNode(node *Node) bool {
	return node.Kind == SyntaxKindTypeQuery
}

// MappedTypeNode

type MappedTypeNode struct {
	TypeNodeBase
	DeclarationBase
	LocalsContainerBase
	ReadonlyToken *Node // Optional
	TypeParameter *Node
	NameType      *Node   // Optional
	QuestionToken *Node   // Optional
	TypeNode      *Node   // Optional (error if missing)
	Members       []*Node // Used only to produce grammar errors
}

func (f *NodeFactory) NewMappedTypeNode(readonlyToken *Node, typeParameter *Node, nameType *Node, questionToken *Node, typeNode *Node, members []*Node) *Node {
	data := &MappedTypeNode{}
	data.ReadonlyToken = readonlyToken
	data.TypeParameter = typeParameter
	data.NameType = nameType
	data.QuestionToken = questionToken
	data.TypeNode = typeNode
	data.Members = members
	return f.NewNode(SyntaxKindMappedType, data)
}

func (node *MappedTypeNode) ForEachChild(v Visitor) bool {
	return visit(v, node.ReadonlyToken) || visit(v, node.TypeParameter) || visit(v, node.NameType) ||
		visit(v, node.QuestionToken) || visit(v, node.TypeNode) || visitNodes(v, node.Members)
}

func IsMappedTypeNode(node *Node) bool {
	return node.Kind == SyntaxKindMappedType
}

// TypeLiteralNode

type TypeLiteralNode struct {
	TypeNodeBase
	DeclarationBase
	Members []*TypeElement
}

func (f *NodeFactory) NewTypeLiteralNode(members []*TypeElement) *Node {
	data := &TypeLiteralNode{}
	data.Members = members
	return f.NewNode(SyntaxKindTypeLiteral, data)
}

func (node *TypeLiteralNode) ForEachChild(v Visitor) bool {
	return visitNodes(v, node.Members)
}

// TupleTypeNode

type TupleTypeNode struct {
	TypeNodeBase
	Elements []*TypeNode
}

func (f *NodeFactory) NewTupleTypeNode(elements []*TypeNode) *Node {
	data := &TupleTypeNode{}
	data.Elements = elements
	return f.NewNode(SyntaxKindTupleType, data)
}

func (node *TupleTypeNode) Kind() SyntaxKind {
	return SyntaxKindTupleType
}

func (node *TupleTypeNode) ForEachChild(v Visitor) bool {
	return visitNodes(v, node.Elements)
}

// NamedTupleTypeMember

type NamedTupleMember struct {
	TypeNodeBase
	DeclarationBase
	DotDotDotToken *Node
	Name           *Node
	QuestionToken  *Node
	TypeNode       *Node
}

func (f *NodeFactory) NewNamedTupleTypeMember(dotDotDotToken *Node, name *Node, questionToken *Node, typeNode *Node) *Node {
	data := &NamedTupleMember{}
	data.DotDotDotToken = dotDotDotToken
	data.Name = name
	data.QuestionToken = questionToken
	data.TypeNode = typeNode
	return f.NewNode(SyntaxKindNamedTupleMember, data)
}

func (node *NamedTupleMember) ForEachChild(v Visitor) bool {
	return visit(v, node.DotDotDotToken) || visit(v, node.Name) || visit(v, node.QuestionToken) || visit(v, node.TypeNode)
}

func IsNamedTupleMember(node *Node) bool {
	return node.Kind == SyntaxKindNamedTupleMember
}

// OptionalTypeNode

type OptionalTypeNode struct {
	TypeNodeBase
	TypeNode *TypeNode
}

func (f *NodeFactory) NewOptionalTypeNode(typeNode *TypeNode) *Node {
	data := &OptionalTypeNode{}
	data.TypeNode = typeNode
	return f.NewNode(SyntaxKindOptionalType, data)
}

func (node *OptionalTypeNode) ForEachChild(v Visitor) bool {
	return visit(v, node.TypeNode)
}

// RestTypeNode

type RestTypeNode struct {
	TypeNodeBase
	TypeNode *TypeNode
}

func (f *NodeFactory) NewRestTypeNode(typeNode *TypeNode) *Node {
	data := &RestTypeNode{}
	data.TypeNode = typeNode
	return f.NewNode(SyntaxKindRestType, data)
}

func (node *RestTypeNode) ForEachChild(v Visitor) bool {
	return visit(v, node.TypeNode)
}

// ParenthesizedTypeNode

type ParenthesizedTypeNode struct {
	TypeNodeBase
	TypeNode *TypeNode
}

func (f *NodeFactory) NewParenthesizedTypeNode(typeNode *TypeNode) *Node {
	data := &ParenthesizedTypeNode{}
	data.TypeNode = typeNode
	return f.NewNode(SyntaxKindParenthesizedType, data)
}

func (node *ParenthesizedTypeNode) ForEachChild(v Visitor) bool {
	return visit(v, node.TypeNode)
}

func IsParenthesizedTypeNode(node *Node) bool {
	return node.Kind == SyntaxKindParenthesizedType
}

// FunctionOrConstructorTypeNodeBase

type FunctionOrConstructorTypeNodeBase struct {
	TypeNodeBase
	DeclarationBase
	ModifiersBase
	FunctionLikeBase
}

func (node *FunctionOrConstructorTypeNodeBase) ForEachChild(v Visitor) bool {
	return visit(v, node.Modifiers) || visit(v, node.TypeParameters) || visitNodes(v, node.Parameters) || visit(v, node.ReturnType)
}

// FunctionTypeNode

type FunctionTypeNode struct {
	FunctionOrConstructorTypeNodeBase
}

func (f *NodeFactory) NewFunctionTypeNode(typeParameters *Node, parameters []*Node, returnType *Node) *Node {
	data := &FunctionTypeNode{}
	data.TypeParameters = typeParameters
	data.Parameters = parameters
	data.ReturnType = returnType
	return f.NewNode(SyntaxKindFunctionType, data)
}

func IsFunctionTypeNode(node *Node) bool {
	return node.Kind == SyntaxKindFunctionType
}

// ConstructorTypeNode

type ConstructorTypeNode struct {
	FunctionOrConstructorTypeNodeBase
}

func (f *NodeFactory) NewConstructorTypeNode(modifiers *Node, typeParameters *Node, parameters []*Node, returnType *Node) *Node {
	data := &ConstructorTypeNode{}
	data.Modifiers = modifiers
	data.TypeParameters = typeParameters
	data.Parameters = parameters
	data.ReturnType = returnType
	return f.NewNode(SyntaxKindConstructorType, data)
}

func IsConstructorTypeNode(node *Node) bool {
	return node.Kind == SyntaxKindConstructorType
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
	return f.NewNode(SyntaxKindTemplateHead, data)
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
	return f.NewNode(SyntaxKindTemplateMiddle, data)
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
	return f.NewNode(SyntaxKindTemplateTail, data)
}

// TemplateLiteralTypeNode

type TemplateLiteralTypeNode struct {
	TypeNodeBase
	Head          *Node
	TemplateSpans []*Node
}

func (f *NodeFactory) NewTemplateLiteralTypeNode(head *Node, templateSpans []*Node) *Node {
	data := &TemplateLiteralTypeNode{}
	data.Head = head
	data.TemplateSpans = templateSpans
	return f.NewNode(SyntaxKindTemplateLiteralType, data)
}

func (node *TemplateLiteralTypeNode) ForEachChild(v Visitor) bool {
	return visit(v, node.Head) || visitNodes(v, node.TemplateSpans)
}

// TemplateLiteralTypeSpan

type TemplateLiteralTypeSpan struct {
	NodeBase
	TypeNode *Node
	Literal  *Node
}

func (f *NodeFactory) NewTemplateLiteralTypeSpan(typeNode *Node, literal *Node) *Node {
	data := &TemplateLiteralTypeSpan{}
	data.TypeNode = typeNode
	data.Literal = literal
	return f.NewNode(SyntaxKindTemplateLiteralTypeSpan, data)
}

func (node *TemplateLiteralTypeSpan) ForEachChild(v Visitor) bool {
	return visit(v, node.TypeNode) || visit(v, node.Literal)
}

/// A JSX expression of the form <TagName attrs>...</TagName>

type JsxElement struct {
	ExpressionBase
	OpeningElement *Node
	Children       []*Node
	ClosingElement *Node
}

func (f *NodeFactory) NewJsxElement(openingElement *Node, children []*Node, closingElement *Node) *Node {
	data := &JsxElement{}
	data.OpeningElement = openingElement
	data.Children = children
	data.ClosingElement = closingElement
	return f.NewNode(SyntaxKindJsxElement, data)
}

func (node *JsxElement) ForEachChild(v Visitor) bool {
	return visit(v, node.OpeningElement) || visitNodes(v, node.Children) || visit(v, node.ClosingElement)
}

// JsxAttributes

type JsxAttributes struct {
	ExpressionBase
	DeclarationBase
	Properties []*JsxAttributeLike
}

func (f *NodeFactory) NewJsxAttributes(properties []*JsxAttributeLike) *Node {
	data := &JsxAttributes{}
	data.Properties = properties
	return f.NewNode(SyntaxKindJsxAttributes, data)
}

func (node *JsxAttributes) ForEachChild(v Visitor) bool {
	return visitNodes(v, node.Properties)
}

func IsJsxAttributes(node *Node) bool {
	return node.Kind == SyntaxKindJsxAttributes
}

// JsxNamespacedName

type JsxNamespacedName struct {
	ExpressionBase
	Name      *Node
	Namespace *Node
}

func (f *NodeFactory) NewJsxNamespacedName(name *Node, namespace *Node) *Node {
	data := &JsxNamespacedName{}
	data.Name = name
	data.Namespace = namespace
	return f.NewNode(SyntaxKindJsxNamespacedName, data)
}

func (node *JsxNamespacedName) ForEachChild(v Visitor) bool {
	return visit(v, node.Name) || visit(v, node.Namespace)
}

func IsJsxNamespacedName(node *Node) bool {
	return node.Kind == SyntaxKindJsxNamespacedName
}

/// The opening element of a <Tag>...</Tag> JsxElement

type JsxOpeningElement struct {
	ExpressionBase
	TagName       *Node // Identifier | KeywordExpression | PropertyAccessExpression | JsxNamespacedName
	TypeArguments *Node
	Attributes    *Node
}

func (f *NodeFactory) NewJsxOpeningElement(tagName *Node, typeArguments *Node, attributes *Node) *Node {
	data := &JsxOpeningElement{}
	data.TagName = tagName
	data.TypeArguments = typeArguments
	data.Attributes = attributes
	return f.NewNode(SyntaxKindJsxOpeningElement, data)
}

func (node *JsxOpeningElement) ForEachChild(v Visitor) bool {
	return visit(v, node.TagName) || visit(v, node.TypeArguments) || visit(v, node.Attributes)
}

func IsJsxOpeningElement(node *Node) bool {
	return node.Kind == SyntaxKindJsxOpeningElement
}

/// A JSX expression of the form <TagName attrs />

type JsxSelfClosingElement struct {
	ExpressionBase
	TagName       *Node // Identifier | KeywordExpression | PropertyAccessExpression | JsxNamespacedName
	TypeArguments *Node
	Attributes    *Node
}

func (f *NodeFactory) NewJsxSelfClosingElement(tagName *Node, typeArguments *Node, attributes *Node) *Node {
	data := &JsxSelfClosingElement{}
	data.TagName = tagName
	data.TypeArguments = typeArguments
	data.Attributes = attributes
	return f.NewNode(SyntaxKindJsxSelfClosingElement, data)
}

func (node *JsxSelfClosingElement) ForEachChild(v Visitor) bool {
	return visit(v, node.TagName) || visit(v, node.TypeArguments) || visit(v, node.Attributes)
}

func IsJsxSelfClosingElement(node *Node) bool {
	return node.Kind == SyntaxKindJsxSelfClosingElement
}

/// A JSX expression of the form <>...</>

type JsxFragment struct {
	ExpressionBase
	OpeningFragment *Node
	Children        []*Node
	ClosingFragment *Node
}

func (f *NodeFactory) NewJsxFragment(openingFragment *Node, children []*Node, closingFragment *Node) *Node {
	data := &JsxFragment{}
	data.OpeningFragment = openingFragment
	data.Children = children
	data.ClosingFragment = closingFragment
	return f.NewNode(SyntaxKindJsxFragment, data)
}

func (node *JsxFragment) ForEachChild(v Visitor) bool {
	return visit(v, node.OpeningFragment) || visitNodes(v, node.Children) || visit(v, node.ClosingFragment)
}

/// The opening element of a <>...</> JsxFragment

type JsxOpeningFragment struct {
	ExpressionBase
}

func (f *NodeFactory) NewJsxOpeningFragment() *Node {
	return f.NewNode(SyntaxKindJsxOpeningFragment, &JsxOpeningFragment{})
}

func IsJsxOpeningFragment(node *Node) bool {
	return node.Kind == SyntaxKindJsxOpeningFragment
}

/// The closing element of a <>...</> JsxFragment

type JsxClosingFragment struct {
	ExpressionBase
}

func (f *NodeFactory) NewJsxClosingFragment() *Node {
	return f.NewNode(SyntaxKindJsxClosingFragment, &JsxClosingFragment{})
}

// JsxAttribute

type JsxAttribute struct {
	NodeBase
	DeclarationBase
	Name *Node
	/// JSX attribute initializers are optional; <X y /> is sugar for <X y={true} />
	Initializer *Node
}

func (f *NodeFactory) NewJsxAttribute(name *Node, initializer *Node) *Node {
	data := &JsxAttribute{}
	data.Name = name
	data.Initializer = initializer
	return f.NewNode(SyntaxKindJsxAttribute, data)
}

func (node *JsxAttribute) ForEachChild(v Visitor) bool {
	return visit(v, node.Name) || visit(v, node.Initializer)
}

func IsJsxAttribute(node *Node) bool {
	return node.Kind == SyntaxKindJsxAttribute
}

// JsxSpreadAttribute

type JsxSpreadAttribute struct {
	NodeBase
	Expression *Node
}

func (f *NodeFactory) NewJsxSpreadAttribute(expression *Node) *Node {
	data := &JsxSpreadAttribute{}
	data.Expression = expression
	return f.NewNode(SyntaxKindJsxSpreadAttribute, data)
}

func (node *JsxSpreadAttribute) ForEachChild(v Visitor) bool {
	return visit(v, node.Expression)
}

// JsxClosingElement

type JsxClosingElement struct {
	NodeBase
	TagName *Node // Identifier | KeywordExpression | PropertyAccessExpression | JsxNamespacedName
}

func (f *NodeFactory) NewJsxClosingElement(tagName *Node) *Node {
	data := &JsxClosingElement{}
	data.TagName = tagName
	return f.NewNode(SyntaxKindJsxClosingElement, data)
}

func (node *JsxClosingElement) ForEachChild(v Visitor) bool {
	return visit(v, node.TagName)
}

// JsxExpression

type JsxExpression struct {
	ExpressionBase
	DotDotDotToken *Node
	Expression     *Node
}

func (f *NodeFactory) NewJsxExpression(dotDotDotToken *Node, expression *Node) *Node {
	data := &JsxExpression{}
	data.DotDotDotToken = dotDotDotToken
	data.Expression = expression
	return f.NewNode(SyntaxKindJsxExpression, data)
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
	return f.NewNode(SyntaxKindJsxText, data)
}

// JSDocNonNullableType

type JSDocNonNullableType struct {
	TypeNodeBase
	TypeNode *Node
	Postfix  bool
}

func (f *NodeFactory) NewJSDocNonNullableType(typeNode *Node, postfix bool) *Node {
	data := &JSDocNonNullableType{}
	data.TypeNode = typeNode
	data.Postfix = postfix
	return f.NewNode(SyntaxKindJSDocNonNullableType, data)
}

func (node *JSDocNonNullableType) ForEachChild(v Visitor) bool {
	return visit(v, node.TypeNode)
}

// JSDocNullableType

type JSDocNullableType struct {
	TypeNodeBase
	TypeNode *Node
	Postfix  bool
}

func (f *NodeFactory) NewJSDocNullableType(typeNode *Node, postfix bool) *Node {
	data := &JSDocNullableType{}
	data.TypeNode = typeNode
	data.Postfix = postfix
	return f.NewNode(SyntaxKindJSDocNullableType, data)
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
	fileName                    string
	path                        string
	Statements                  []*Statement
	diagnostics                 []*Diagnostic
	bindDiagnostics             []*Diagnostic
	BindSuggestionDiagnostics   []*Diagnostic
	LineMap                     []textpos.TextPos
	LanguageVersion             ScriptTarget
	LanguageVariant             LanguageVariant
	ScriptKind                  ScriptKind
	ExternalModuleIndicator     *Node
	EndFlowNode                 *FlowNode
	JsGlobalAugmentations       SymbolTable
	IsDeclarationFile           bool
	IsBound                     bool
	ModuleReferencesProcessed   bool
	UsesUriStyleNodeCoreModules Tristate
	SymbolCount                 int
	ClassifiableNames           set[string]
	Imports                     []*LiteralLikeNode
	ModuleAugmentations         []*ModuleName
	PatternAmbientModules       []PatternAmbientModule
	AmbientModuleNames          []string
}

func (f *NodeFactory) NewSourceFile(text string, fileName string, statements []*Node) *Node {
	data := &SourceFile{}
	data.Text = text
	data.fileName = fileName
	data.Statements = statements
	data.LanguageVersion = ScriptTargetLatest
	return f.NewNode(SyntaxKindSourceFile, data)
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
	return visitNodes(v, node.Statements)
}

func IsSourceFile(node *Node) bool {
	return node.Kind == SyntaxKindSourceFile
}
