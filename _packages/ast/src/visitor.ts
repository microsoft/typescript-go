// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
// !!! THIS FILE IS AUTO-GENERATED - DO NOT EDIT !!!
// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
//
// Source: _packages/ast/src/nodes.ts
// Generator: _packages/ast/scripts/generateVisitor.ts
//

import { SyntaxKind } from "#enums/syntaxKind";
import {
    createNodeArray,
    updateArrayBindingPattern,
    updateArrayLiteralExpression,
    updateArrayTypeNode,
    updateArrowFunction,
    updateAsExpression,
    updateAwaitExpression,
    updateBinaryExpression,
    updateBindingElement,
    updateBlock,
    updateBreakStatement,
    updateCallExpression,
    updateCallSignatureDeclaration,
    updateCaseBlock,
    updateCaseClause,
    updateCatchClause,
    updateClassDeclaration,
    updateClassExpression,
    updateClassStaticBlockDeclaration,
    updateCommaListExpression,
    updateComputedPropertyName,
    updateConditionalExpression,
    updateConditionalTypeNode,
    updateConstructorDeclaration,
    updateConstructorTypeNode,
    updateConstructSignatureDeclaration,
    updateContinueStatement,
    updateDecorator,
    updateDefaultClause,
    updateDeleteExpression,
    updateDoStatement,
    updateElementAccessExpression,
    updateEnumDeclaration,
    updateEnumMember,
    updateExportAssignment,
    updateExportDeclaration,
    updateExportSpecifier,
    updateExpressionStatement,
    updateExpressionWithTypeArguments,
    updateExternalModuleReference,
    updateForInStatement,
    updateForOfStatement,
    updateForStatement,
    updateFunctionDeclaration,
    updateFunctionExpression,
    updateFunctionTypeNode,
    updateGetAccessorDeclaration,
    updateHeritageClause,
    updateIfStatement,
    updateImportAttribute,
    updateImportAttributes,
    updateImportClause,
    updateImportDeclaration,
    updateImportEqualsDeclaration,
    updateImportSpecifier,
    updateImportTypeNode,
    updateIndexedAccessTypeNode,
    updateIndexSignatureDeclaration,
    updateInferTypeNode,
    updateInterfaceDeclaration,
    updateIntersectionTypeNode,
    updateJSDoc,
    updateJSDocAugmentsTag,
    updateJSDocCallbackTag,
    updateJSDocDeprecatedTag,
    updateJSDocImplementsTag,
    updateJSDocImportTag,
    updateJSDocLink,
    updateJSDocLinkCode,
    updateJSDocLinkPlain,
    updateJSDocMemberName,
    updateJSDocNameReference,
    updateJSDocNonNullableType,
    updateJSDocNullableType,
    updateJSDocOptionalType,
    updateJSDocOverloadTag,
    updateJSDocOverrideTag,
    updateJSDocParameterTag,
    updateJSDocPrivateTag,
    updateJSDocProtectedTag,
    updateJSDocPublicTag,
    updateJSDocReadonlyTag,
    updateJSDocReturnTag,
    updateJSDocSatisfiesTag,
    updateJSDocSeeTag,
    updateJSDocSignature,
    updateJSDocTemplateTag,
    updateJSDocThisTag,
    updateJSDocTypedefTag,
    updateJSDocTypeExpression,
    updateJSDocTypeTag,
    updateJSDocUnknownTag,
    updateJSDocVariadicType,
    updateJsxAttribute,
    updateJsxAttributes,
    updateJsxClosingElement,
    updateJsxElement,
    updateJsxExpression,
    updateJsxFragment,
    updateJsxNamespacedName,
    updateJsxOpeningElement,
    updateJsxSelfClosingElement,
    updateJsxSpreadAttribute,
    updateLabeledStatement,
    updateLiteralTypeNode,
    updateMappedTypeNode,
    updateMetaProperty,
    updateMethodDeclaration,
    updateMethodSignature,
    updateModuleBlock,
    updateModuleDeclaration,
    updateNamedExports,
    updateNamedImports,
    updateNamedTupleMember,
    updateNamespaceExport,
    updateNamespaceExportDeclaration,
    updateNamespaceImport,
    updateNewExpression,
    updateNonNullExpression,
    updateObjectBindingPattern,
    updateObjectLiteralExpression,
    updateOptionalTypeNode,
    updateParameterDeclaration,
    updateParenthesizedExpression,
    updateParenthesizedTypeNode,
    updatePartiallyEmittedExpression,
    updatePostfixUnaryExpression,
    updatePrefixUnaryExpression,
    updatePropertyAccessExpression,
    updatePropertyAssignment,
    updatePropertyDeclaration,
    updatePropertySignature,
    updateQualifiedName,
    updateRestTypeNode,
    updateReturnStatement,
    updateSatisfiesExpression,
    updateSetAccessorDeclaration,
    updateShorthandPropertyAssignment,
    updateSourceFile,
    updateSpreadAssignment,
    updateSpreadElement,
    updateSwitchStatement,
    updateTaggedTemplateExpression,
    updateTemplateExpression,
    updateTemplateLiteralTypeNode,
    updateTemplateLiteralTypeSpan,
    updateTemplateSpan,
    updateThrowStatement,
    updateTryStatement,
    updateTupleTypeNode,
    updateTypeAliasDeclaration,
    updateTypeAssertion,
    updateTypeLiteralNode,
    updateTypeOfExpression,
    updateTypeOperatorNode,
    updateTypeParameterDeclaration,
    updateTypePredicateNode,
    updateTypeQueryNode,
    updateTypeReferenceNode,
    updateUnionTypeNode,
    updateVariableDeclaration,
    updateVariableDeclarationList,
    updateVariableStatement,
    updateVoidExpression,
    updateWhileStatement,
    updateWithStatement,
    updateYieldExpression,
} from "./factory.ts";
import type {
    ArrayBindingPattern,
    ArrayLiteralExpression,
    ArrayTypeNode,
    ArrowFunction,
    AsExpression,
    AwaitExpression,
    BinaryExpression,
    BindingElement,
    Block,
    BreakStatement,
    CallExpression,
    CallSignatureDeclaration,
    CaseBlock,
    CaseClause,
    CatchClause,
    ClassDeclaration,
    ClassExpression,
    ClassStaticBlockDeclaration,
    CommaListExpression,
    ComputedPropertyName,
    ConditionalExpression,
    ConditionalTypeNode,
    ConstructorDeclaration,
    ConstructorTypeNode,
    ConstructSignatureDeclaration,
    ContinueStatement,
    Decorator,
    DefaultClause,
    DeleteExpression,
    DoStatement,
    ElementAccessExpression,
    EnumDeclaration,
    EnumMember,
    ExportAssignment,
    ExportDeclaration,
    ExportSpecifier,
    ExpressionStatement,
    ExpressionWithTypeArguments,
    ExternalModuleReference,
    ForInStatement,
    ForOfStatement,
    ForStatement,
    FunctionDeclaration,
    FunctionExpression,
    FunctionTypeNode,
    GetAccessorDeclaration,
    HeritageClause,
    IfStatement,
    ImportAttribute,
    ImportAttributes,
    ImportClause,
    ImportDeclaration,
    ImportEqualsDeclaration,
    ImportSpecifier,
    ImportTypeNode,
    IndexedAccessTypeNode,
    IndexSignatureDeclaration,
    InferTypeNode,
    InterfaceDeclaration,
    IntersectionTypeNode,
    JSDoc,
    JSDocAugmentsTag,
    JSDocCallbackTag,
    JSDocDeprecatedTag,
    JSDocImplementsTag,
    JSDocImportTag,
    JSDocLink,
    JSDocLinkCode,
    JSDocLinkPlain,
    JSDocMemberName,
    JSDocNameReference,
    JSDocNonNullableType,
    JSDocNullableType,
    JSDocOptionalType,
    JSDocOverloadTag,
    JSDocOverrideTag,
    JSDocParameterTag,
    JSDocPrivateTag,
    JSDocProtectedTag,
    JSDocPublicTag,
    JSDocReadonlyTag,
    JSDocReturnTag,
    JSDocSatisfiesTag,
    JSDocSeeTag,
    JSDocSignature,
    JSDocTemplateTag,
    JSDocThisTag,
    JSDocTypedefTag,
    JSDocTypeExpression,
    JSDocTypeTag,
    JSDocUnknownTag,
    JSDocVariadicType,
    JsxAttribute,
    JsxAttributes,
    JsxClosingElement,
    JsxElement,
    JsxExpression,
    JsxFragment,
    JsxNamespacedName,
    JsxOpeningElement,
    JsxSelfClosingElement,
    JsxSpreadAttribute,
    LabeledStatement,
    LiteralTypeNode,
    MappedTypeNode,
    MetaProperty,
    MethodDeclaration,
    MethodSignature,
    ModuleBlock,
    ModuleDeclaration,
    NamedExports,
    NamedImports,
    NamedTupleMember,
    NamespaceExport,
    NamespaceExportDeclaration,
    NamespaceImport,
    NewExpression,
    Node,
    NodeArray,
    NonNullExpression,
    ObjectBindingPattern,
    ObjectLiteralExpression,
    OptionalTypeNode,
    ParameterDeclaration,
    ParenthesizedExpression,
    ParenthesizedTypeNode,
    PartiallyEmittedExpression,
    PostfixUnaryExpression,
    PrefixUnaryExpression,
    PropertyAccessExpression,
    PropertyAssignment,
    PropertyDeclaration,
    PropertySignature,
    QualifiedName,
    RestTypeNode,
    ReturnStatement,
    SatisfiesExpression,
    SetAccessorDeclaration,
    ShorthandPropertyAssignment,
    SourceFile,
    SpreadAssignment,
    SpreadElement,
    SwitchStatement,
    TaggedTemplateExpression,
    TemplateExpression,
    TemplateLiteralTypeNode,
    TemplateLiteralTypeSpan,
    TemplateSpan,
    ThrowStatement,
    TryStatement,
    TupleTypeNode,
    TypeAliasDeclaration,
    TypeAssertion,
    TypeLiteralNode,
    TypeOfExpression,
    TypeOperatorNode,
    TypeParameterDeclaration,
    TypePredicateNode,
    TypeQueryNode,
    TypeReferenceNode,
    UnionTypeNode,
    VariableDeclaration,
    VariableDeclarationList,
    VariableStatement,
    VoidExpression,
    WhileStatement,
    WithStatement,
    YieldExpression,
} from "./nodes.ts";

/**
 * A callback that receives a node and returns a visited node (or undefined to remove it).
 */
export type Visitor = (node: Node) => Node | undefined;

/**
 * Visits a Node using the supplied visitor, possibly returning a new Node in its place.
 *
 * - If the input node is undefined, then the output is undefined.
 * - If the visitor returns undefined, then the output is undefined.
 */
export function visitNode<T extends Node>(node: T, visitor: Visitor): T | undefined;
export function visitNode<T extends Node>(node: T | undefined, visitor: Visitor): T | undefined;
export function visitNode(node: Node | undefined, visitor: Visitor): Node | undefined {
    if (node === undefined) return undefined;
    return visitor(node);
}

/**
 * Visits a NodeArray using the supplied visitor, possibly returning a new NodeArray in its place.
 *
 * - If the input node array is undefined, the output is undefined.
 * - If the visitor returns undefined for a node, that node is dropped from the result.
 */
export function visitNodes<T extends Node>(nodes: NodeArray<T>, visitor: Visitor): NodeArray<T>;
export function visitNodes<T extends Node>(nodes: NodeArray<T> | undefined, visitor: Visitor): NodeArray<T> | undefined;
export function visitNodes(nodes: NodeArray<Node> | undefined, visitor: Visitor): NodeArray<Node> | undefined {
    if (nodes === undefined) return undefined;
    let updated: Node[] | undefined;
    for (let i = 0; i < nodes.length; i++) {
        const node = nodes[i];
        const visited = visitor(node);
        if (updated) {
            if (visited) updated.push(visited);
        }
        else if (visited !== node) {
            updated = nodes.slice(0, i) as unknown as Node[];
            if (visited) updated.push(visited);
        }
    }
    if (!updated) return nodes;
    return createNodeArray(updated, nodes.pos, nodes.end);
}

/**
 * Visits each child of a Node using the supplied visitor, possibly returning a new Node of the same kind in its place.
 *
 * @param node The Node whose children will be visited.
 * @param visitor The callback used to visit each child.
 * @returns The original node if no children changed, or a new node with visited children.
 */
export function visitEachChild<T extends Node>(node: T, visitor: Visitor): T;
export function visitEachChild<T extends Node>(node: T | undefined, visitor: Visitor): T | undefined;
export function visitEachChild(node: Node | undefined, visitor: Visitor): Node | undefined {
    if (node === undefined) return undefined;
    const fn = visitEachChildTable[node.kind];
    return fn ? fn(node, visitor) : node;
}

type VisitEachChildFunction = (node: any, visitor: Visitor) => Node;

const visitEachChildTable: Record<number, VisitEachChildFunction> = {
    [SyntaxKind.ArrayBindingPattern]: (node: ArrayBindingPattern, visitor: Visitor): ArrayBindingPattern => {
        const _elements = visitNodes(node.elements, visitor);
        return updateArrayBindingPattern(node, _elements);
    },
    [SyntaxKind.ArrayLiteralExpression]: (node: ArrayLiteralExpression, visitor: Visitor): ArrayLiteralExpression => {
        const _elements = visitNodes(node.elements, visitor);
        return updateArrayLiteralExpression(node, _elements);
    },
    [SyntaxKind.ArrayType]: (node: ArrayTypeNode, visitor: Visitor): ArrayTypeNode => {
        const _elementType = visitNode(node.elementType, visitor)!;
        return updateArrayTypeNode(node, _elementType);
    },
    [SyntaxKind.ArrowFunction]: (node: ArrowFunction, visitor: Visitor): ArrowFunction => {
        const _modifiers = visitNodes(node.modifiers, visitor);
        const _typeParameters = visitNodes(node.typeParameters, visitor);
        const _parameters = visitNodes(node.parameters, visitor);
        const _type = visitNode(node.type, visitor);
        const _equalsGreaterThanToken = visitNode(node.equalsGreaterThanToken, visitor)!;
        const _body = visitNode(node.body, visitor)!;
        return updateArrowFunction(node, _modifiers, _typeParameters, _parameters, _type, _equalsGreaterThanToken, _body);
    },
    [SyntaxKind.AsExpression]: (node: AsExpression, visitor: Visitor): AsExpression => {
        const _expression = visitNode(node.expression, visitor)!;
        const _type = visitNode(node.type, visitor)!;
        return updateAsExpression(node, _expression, _type);
    },
    [SyntaxKind.AwaitExpression]: (node: AwaitExpression, visitor: Visitor): AwaitExpression => {
        const _expression = visitNode(node.expression, visitor)!;
        return updateAwaitExpression(node, _expression);
    },
    [SyntaxKind.BinaryExpression]: (node: BinaryExpression, visitor: Visitor): BinaryExpression => {
        const _left = visitNode(node.left, visitor)!;
        const _operatorToken = visitNode(node.operatorToken, visitor)!;
        const _right = visitNode(node.right, visitor)!;
        return updateBinaryExpression(node, _left, _operatorToken, _right);
    },
    [SyntaxKind.BindingElement]: (node: BindingElement, visitor: Visitor): BindingElement => {
        const _dotDotDotToken = visitNode(node.dotDotDotToken, visitor);
        const _propertyName = visitNode(node.propertyName, visitor);
        const _initializer = visitNode(node.initializer, visitor);
        return updateBindingElement(node, _dotDotDotToken, _propertyName, _initializer);
    },
    [SyntaxKind.Block]: (node: Block, visitor: Visitor): Block => {
        const _statements = visitNodes(node.statements, visitor);
        return updateBlock(node, _statements);
    },
    [SyntaxKind.BreakStatement]: (node: BreakStatement, visitor: Visitor): BreakStatement => {
        const _label = visitNode(node.label, visitor);
        return updateBreakStatement(node, _label);
    },
    [SyntaxKind.CallExpression]: (node: CallExpression, visitor: Visitor): CallExpression => {
        const _expression = visitNode(node.expression, visitor)!;
        const _questionDotToken = visitNode(node.questionDotToken, visitor);
        const _typeArguments = visitNodes(node.typeArguments, visitor);
        const _arguments = visitNodes(node.arguments, visitor);
        return updateCallExpression(node, _expression, _questionDotToken, _typeArguments, _arguments);
    },
    [SyntaxKind.CallSignature]: (node: CallSignatureDeclaration, visitor: Visitor): CallSignatureDeclaration => {
        const _typeParameters = visitNodes(node.typeParameters, visitor);
        const _parameters = visitNodes(node.parameters, visitor);
        const _type = visitNode(node.type, visitor);
        return updateCallSignatureDeclaration(node, _typeParameters, _parameters, _type);
    },
    [SyntaxKind.CaseBlock]: (node: CaseBlock, visitor: Visitor): CaseBlock => {
        const _clauses = visitNodes(node.clauses, visitor);
        return updateCaseBlock(node, _clauses);
    },
    [SyntaxKind.CaseClause]: (node: CaseClause, visitor: Visitor): CaseClause => {
        const _expression = visitNode(node.expression, visitor)!;
        const _statements = visitNodes(node.statements, visitor);
        return updateCaseClause(node, _expression, _statements);
    },
    [SyntaxKind.CatchClause]: (node: CatchClause, visitor: Visitor): CatchClause => {
        const _variableDeclaration = visitNode(node.variableDeclaration, visitor);
        const _block = visitNode(node.block, visitor)!;
        return updateCatchClause(node, _variableDeclaration, _block);
    },
    [SyntaxKind.ClassDeclaration]: (node: ClassDeclaration, visitor: Visitor): ClassDeclaration => {
        const _modifiers = visitNodes(node.modifiers, visitor);
        const _name = visitNode(node.name, visitor);
        const _typeParameters = visitNodes(node.typeParameters, visitor);
        const _heritageClauses = visitNodes(node.heritageClauses, visitor);
        const _members = visitNodes(node.members, visitor);
        return updateClassDeclaration(node, _modifiers, _name, _typeParameters, _heritageClauses, _members);
    },
    [SyntaxKind.ClassExpression]: (node: ClassExpression, visitor: Visitor): ClassExpression => {
        const _modifiers = visitNodes(node.modifiers, visitor);
        const _name = visitNode(node.name, visitor);
        const _typeParameters = visitNodes(node.typeParameters, visitor);
        const _heritageClauses = visitNodes(node.heritageClauses, visitor);
        const _members = visitNodes(node.members, visitor);
        return updateClassExpression(node, _modifiers, _name, _typeParameters, _heritageClauses, _members);
    },
    [SyntaxKind.ClassStaticBlockDeclaration]: (node: ClassStaticBlockDeclaration, visitor: Visitor): ClassStaticBlockDeclaration => {
        const _body = visitNode(node.body, visitor)!;
        return updateClassStaticBlockDeclaration(node, _body);
    },
    [SyntaxKind.CommaListExpression]: (node: CommaListExpression, visitor: Visitor): CommaListExpression => {
        const _elements = visitNodes(node.elements, visitor);
        return updateCommaListExpression(node, _elements);
    },
    [SyntaxKind.ComputedPropertyName]: (node: ComputedPropertyName, visitor: Visitor): ComputedPropertyName => {
        const _expression = visitNode(node.expression, visitor)!;
        return updateComputedPropertyName(node, _expression);
    },
    [SyntaxKind.ConditionalExpression]: (node: ConditionalExpression, visitor: Visitor): ConditionalExpression => {
        const _condition = visitNode(node.condition, visitor)!;
        const _questionToken = visitNode(node.questionToken, visitor)!;
        const _whenTrue = visitNode(node.whenTrue, visitor)!;
        const _colonToken = visitNode(node.colonToken, visitor)!;
        const _whenFalse = visitNode(node.whenFalse, visitor)!;
        return updateConditionalExpression(node, _condition, _questionToken, _whenTrue, _colonToken, _whenFalse);
    },
    [SyntaxKind.ConditionalType]: (node: ConditionalTypeNode, visitor: Visitor): ConditionalTypeNode => {
        const _checkType = visitNode(node.checkType, visitor)!;
        const _extendsType = visitNode(node.extendsType, visitor)!;
        const _trueType = visitNode(node.trueType, visitor)!;
        const _falseType = visitNode(node.falseType, visitor)!;
        return updateConditionalTypeNode(node, _checkType, _extendsType, _trueType, _falseType);
    },
    [SyntaxKind.Constructor]: (node: ConstructorDeclaration, visitor: Visitor): ConstructorDeclaration => {
        const _modifiers = visitNodes(node.modifiers, visitor);
        const _parameters = visitNodes(node.parameters, visitor);
        const _body = visitNode(node.body, visitor);
        return updateConstructorDeclaration(node, _modifiers, _parameters, _body);
    },
    [SyntaxKind.ConstructorType]: (node: ConstructorTypeNode, visitor: Visitor): ConstructorTypeNode => {
        const _modifiers = visitNodes(node.modifiers, visitor);
        const _typeParameters = visitNodes(node.typeParameters, visitor);
        const _parameters = visitNodes(node.parameters, visitor);
        const _type = visitNode(node.type, visitor)!;
        return updateConstructorTypeNode(node, _modifiers, _typeParameters, _parameters, _type);
    },
    [SyntaxKind.ConstructSignature]: (node: ConstructSignatureDeclaration, visitor: Visitor): ConstructSignatureDeclaration => {
        const _typeParameters = visitNodes(node.typeParameters, visitor);
        const _parameters = visitNodes(node.parameters, visitor);
        const _type = visitNode(node.type, visitor);
        return updateConstructSignatureDeclaration(node, _typeParameters, _parameters, _type);
    },
    [SyntaxKind.ContinueStatement]: (node: ContinueStatement, visitor: Visitor): ContinueStatement => {
        const _label = visitNode(node.label, visitor);
        return updateContinueStatement(node, _label);
    },
    [SyntaxKind.Decorator]: (node: Decorator, visitor: Visitor): Decorator => {
        const _expression = visitNode(node.expression, visitor)!;
        return updateDecorator(node, _expression);
    },
    [SyntaxKind.DefaultClause]: (node: DefaultClause, visitor: Visitor): DefaultClause => {
        const _statements = visitNodes(node.statements, visitor);
        return updateDefaultClause(node, _statements);
    },
    [SyntaxKind.DeleteExpression]: (node: DeleteExpression, visitor: Visitor): DeleteExpression => {
        const _expression = visitNode(node.expression, visitor)!;
        return updateDeleteExpression(node, _expression);
    },
    [SyntaxKind.DoStatement]: (node: DoStatement, visitor: Visitor): DoStatement => {
        const _statement = visitNode(node.statement, visitor)!;
        const _expression = visitNode(node.expression, visitor)!;
        return updateDoStatement(node, _statement, _expression);
    },
    [SyntaxKind.ElementAccessExpression]: (node: ElementAccessExpression, visitor: Visitor): ElementAccessExpression => {
        const _expression = visitNode(node.expression, visitor)!;
        const _questionDotToken = visitNode(node.questionDotToken, visitor);
        const _argumentExpression = visitNode(node.argumentExpression, visitor)!;
        return updateElementAccessExpression(node, _expression, _questionDotToken, _argumentExpression);
    },
    [SyntaxKind.EnumDeclaration]: (node: EnumDeclaration, visitor: Visitor): EnumDeclaration => {
        const _modifiers = visitNodes(node.modifiers, visitor);
        const _name = visitNode(node.name, visitor)!;
        const _members = visitNodes(node.members, visitor);
        return updateEnumDeclaration(node, _modifiers, _name, _members);
    },
    [SyntaxKind.EnumMember]: (node: EnumMember, visitor: Visitor): EnumMember => {
        const _name = visitNode(node.name, visitor)!;
        const _initializer = visitNode(node.initializer, visitor);
        return updateEnumMember(node, _name, _initializer);
    },
    [SyntaxKind.ExportAssignment]: (node: ExportAssignment, visitor: Visitor): ExportAssignment => {
        const _modifiers = visitNodes(node.modifiers, visitor);
        const _expression = visitNode(node.expression, visitor)!;
        return updateExportAssignment(node, _modifiers, _expression);
    },
    [SyntaxKind.ExportDeclaration]: (node: ExportDeclaration, visitor: Visitor): ExportDeclaration => {
        const _modifiers = visitNodes(node.modifiers, visitor);
        const _exportClause = visitNode(node.exportClause, visitor);
        const _moduleSpecifier = visitNode(node.moduleSpecifier, visitor);
        const _attributes = visitNode(node.attributes, visitor);
        return updateExportDeclaration(node, _modifiers, _exportClause, _moduleSpecifier, _attributes);
    },
    [SyntaxKind.ExportSpecifier]: (node: ExportSpecifier, visitor: Visitor): ExportSpecifier => {
        const _propertyName = visitNode(node.propertyName, visitor);
        const _name = visitNode(node.name, visitor)!;
        return updateExportSpecifier(node, _propertyName, _name);
    },
    [SyntaxKind.ExpressionStatement]: (node: ExpressionStatement, visitor: Visitor): ExpressionStatement => {
        const _expression = visitNode(node.expression, visitor)!;
        return updateExpressionStatement(node, _expression);
    },
    [SyntaxKind.ExpressionWithTypeArguments]: (node: ExpressionWithTypeArguments, visitor: Visitor): ExpressionWithTypeArguments => {
        const _expression = visitNode(node.expression, visitor)!;
        const _typeArguments = visitNodes(node.typeArguments, visitor);
        return updateExpressionWithTypeArguments(node, _expression, _typeArguments);
    },
    [SyntaxKind.ExternalModuleReference]: (node: ExternalModuleReference, visitor: Visitor): ExternalModuleReference => {
        const _expression = visitNode(node.expression, visitor)!;
        return updateExternalModuleReference(node, _expression);
    },
    [SyntaxKind.ForInStatement]: (node: ForInStatement, visitor: Visitor): ForInStatement => {
        const _initializer = visitNode(node.initializer, visitor)!;
        const _expression = visitNode(node.expression, visitor)!;
        const _statement = visitNode(node.statement, visitor)!;
        return updateForInStatement(node, _initializer, _expression, _statement);
    },
    [SyntaxKind.ForOfStatement]: (node: ForOfStatement, visitor: Visitor): ForOfStatement => {
        const _awaitModifier = visitNode(node.awaitModifier, visitor);
        const _initializer = visitNode(node.initializer, visitor)!;
        const _expression = visitNode(node.expression, visitor)!;
        const _statement = visitNode(node.statement, visitor)!;
        return updateForOfStatement(node, _awaitModifier, _initializer, _expression, _statement);
    },
    [SyntaxKind.ForStatement]: (node: ForStatement, visitor: Visitor): ForStatement => {
        const _initializer = visitNode(node.initializer, visitor);
        const _condition = visitNode(node.condition, visitor);
        const _incrementor = visitNode(node.incrementor, visitor);
        const _statement = visitNode(node.statement, visitor)!;
        return updateForStatement(node, _initializer, _condition, _incrementor, _statement);
    },
    [SyntaxKind.FunctionDeclaration]: (node: FunctionDeclaration, visitor: Visitor): FunctionDeclaration => {
        const _modifiers = visitNodes(node.modifiers, visitor);
        const _asteriskToken = visitNode(node.asteriskToken, visitor);
        const _name = visitNode(node.name, visitor);
        const _typeParameters = visitNodes(node.typeParameters, visitor);
        const _parameters = visitNodes(node.parameters, visitor);
        const _type = visitNode(node.type, visitor);
        const _body = visitNode(node.body, visitor);
        return updateFunctionDeclaration(node, _modifiers, _asteriskToken, _name, _typeParameters, _parameters, _type, _body);
    },
    [SyntaxKind.FunctionExpression]: (node: FunctionExpression, visitor: Visitor): FunctionExpression => {
        const _modifiers = visitNodes(node.modifiers, visitor);
        const _asteriskToken = visitNode(node.asteriskToken, visitor);
        const _name = visitNode(node.name, visitor);
        const _typeParameters = visitNodes(node.typeParameters, visitor);
        const _parameters = visitNodes(node.parameters, visitor);
        const _type = visitNode(node.type, visitor);
        const _body = visitNode(node.body, visitor)!;
        return updateFunctionExpression(node, _modifiers, _asteriskToken, _name, _typeParameters, _parameters, _type, _body);
    },
    [SyntaxKind.FunctionType]: (node: FunctionTypeNode, visitor: Visitor): FunctionTypeNode => {
        const _typeParameters = visitNodes(node.typeParameters, visitor);
        const _parameters = visitNodes(node.parameters, visitor);
        const _type = visitNode(node.type, visitor)!;
        return updateFunctionTypeNode(node, _typeParameters, _parameters, _type);
    },
    [SyntaxKind.GetAccessor]: (node: GetAccessorDeclaration, visitor: Visitor): GetAccessorDeclaration => {
        const _modifiers = visitNodes(node.modifiers, visitor);
        const _name = visitNode(node.name, visitor)!;
        const _parameters = visitNodes(node.parameters, visitor);
        const _type = visitNode(node.type, visitor);
        const _body = visitNode(node.body, visitor);
        return updateGetAccessorDeclaration(node, _modifiers, _name, _parameters, _type, _body);
    },
    [SyntaxKind.HeritageClause]: (node: HeritageClause, visitor: Visitor): HeritageClause => {
        const _types = visitNodes(node.types, visitor);
        return updateHeritageClause(node, _types);
    },
    [SyntaxKind.IfStatement]: (node: IfStatement, visitor: Visitor): IfStatement => {
        const _expression = visitNode(node.expression, visitor)!;
        const _thenStatement = visitNode(node.thenStatement, visitor)!;
        const _elseStatement = visitNode(node.elseStatement, visitor);
        return updateIfStatement(node, _expression, _thenStatement, _elseStatement);
    },
    [SyntaxKind.ImportAttribute]: (node: ImportAttribute, visitor: Visitor): ImportAttribute => {
        const _name = visitNode(node.name, visitor)!;
        const _value = visitNode(node.value, visitor)!;
        return updateImportAttribute(node, _name, _value);
    },
    [SyntaxKind.ImportAttributes]: (node: ImportAttributes, visitor: Visitor): ImportAttributes => {
        const _elements = visitNodes(node.elements, visitor);
        return updateImportAttributes(node, _elements);
    },
    [SyntaxKind.ImportClause]: (node: ImportClause, visitor: Visitor): ImportClause => {
        const _name = visitNode(node.name, visitor);
        const _namedBindings = visitNode(node.namedBindings, visitor);
        return updateImportClause(node, _name, _namedBindings);
    },
    [SyntaxKind.ImportDeclaration]: (node: ImportDeclaration, visitor: Visitor): ImportDeclaration => {
        const _modifiers = visitNodes(node.modifiers, visitor);
        const _importClause = visitNode(node.importClause, visitor);
        const _moduleSpecifier = visitNode(node.moduleSpecifier, visitor)!;
        const _attributes = visitNode(node.attributes, visitor);
        return updateImportDeclaration(node, _modifiers, _importClause, _moduleSpecifier, _attributes);
    },
    [SyntaxKind.ImportEqualsDeclaration]: (node: ImportEqualsDeclaration, visitor: Visitor): ImportEqualsDeclaration => {
        const _modifiers = visitNodes(node.modifiers, visitor);
        const _name = visitNode(node.name, visitor)!;
        const _moduleReference = visitNode(node.moduleReference, visitor)!;
        return updateImportEqualsDeclaration(node, _modifiers, _name, _moduleReference);
    },
    [SyntaxKind.ImportSpecifier]: (node: ImportSpecifier, visitor: Visitor): ImportSpecifier => {
        const _propertyName = visitNode(node.propertyName, visitor);
        const _name = visitNode(node.name, visitor)!;
        return updateImportSpecifier(node, _propertyName, _name);
    },
    [SyntaxKind.ImportType]: (node: ImportTypeNode, visitor: Visitor): ImportTypeNode => {
        const _argument = visitNode(node.argument, visitor)!;
        const _attributes = visitNode(node.attributes, visitor);
        const _qualifier = visitNode(node.qualifier, visitor);
        const _typeArguments = visitNodes(node.typeArguments, visitor);
        return updateImportTypeNode(node, _argument, _attributes, _qualifier, _typeArguments);
    },
    [SyntaxKind.IndexedAccessType]: (node: IndexedAccessTypeNode, visitor: Visitor): IndexedAccessTypeNode => {
        const _objectType = visitNode(node.objectType, visitor)!;
        const _indexType = visitNode(node.indexType, visitor)!;
        return updateIndexedAccessTypeNode(node, _objectType, _indexType);
    },
    [SyntaxKind.IndexSignature]: (node: IndexSignatureDeclaration, visitor: Visitor): IndexSignatureDeclaration => {
        const _modifiers = visitNodes(node.modifiers, visitor);
        const _parameters = visitNodes(node.parameters, visitor);
        const _type = visitNode(node.type, visitor)!;
        return updateIndexSignatureDeclaration(node, _modifiers, _parameters, _type);
    },
    [SyntaxKind.InferType]: (node: InferTypeNode, visitor: Visitor): InferTypeNode => {
        const _typeParameter = visitNode(node.typeParameter, visitor)!;
        return updateInferTypeNode(node, _typeParameter);
    },
    [SyntaxKind.InterfaceDeclaration]: (node: InterfaceDeclaration, visitor: Visitor): InterfaceDeclaration => {
        const _modifiers = visitNodes(node.modifiers, visitor);
        const _name = visitNode(node.name, visitor)!;
        const _typeParameters = visitNodes(node.typeParameters, visitor);
        const _heritageClauses = visitNodes(node.heritageClauses, visitor);
        const _members = visitNodes(node.members, visitor);
        return updateInterfaceDeclaration(node, _modifiers, _name, _typeParameters, _heritageClauses, _members);
    },
    [SyntaxKind.IntersectionType]: (node: IntersectionTypeNode, visitor: Visitor): IntersectionTypeNode => {
        const _types = visitNodes(node.types, visitor);
        return updateIntersectionTypeNode(node, _types);
    },
    [SyntaxKind.JSDoc]: (node: JSDoc, visitor: Visitor): JSDoc => {
        const _tags = visitNodes(node.tags, visitor);
        return updateJSDoc(node, _tags);
    },
    [SyntaxKind.JSDocAugmentsTag]: (node: JSDocAugmentsTag, visitor: Visitor): JSDocAugmentsTag => {
        const _tagName = visitNode(node.tagName, visitor)!;
        const _class = visitNode(node.class, visitor)!;
        return updateJSDocAugmentsTag(node, _tagName, _class);
    },
    [SyntaxKind.JSDocCallbackTag]: (node: JSDocCallbackTag, visitor: Visitor): JSDocCallbackTag => {
        const _tagName = visitNode(node.tagName, visitor)!;
        const _typeExpression = visitNode(node.typeExpression, visitor)!;
        const _fullName = visitNode(node.fullName, visitor);
        return updateJSDocCallbackTag(node, _tagName, _typeExpression, _fullName);
    },
    [SyntaxKind.JSDocDeprecatedTag]: (node: JSDocDeprecatedTag, visitor: Visitor): JSDocDeprecatedTag => {
        const _tagName = visitNode(node.tagName, visitor)!;
        return updateJSDocDeprecatedTag(node, _tagName);
    },
    [SyntaxKind.JSDocImplementsTag]: (node: JSDocImplementsTag, visitor: Visitor): JSDocImplementsTag => {
        const _tagName = visitNode(node.tagName, visitor)!;
        const _class = visitNode(node.class, visitor)!;
        return updateJSDocImplementsTag(node, _tagName, _class);
    },
    [SyntaxKind.JSDocImportTag]: (node: JSDocImportTag, visitor: Visitor): JSDocImportTag => {
        const _tagName = visitNode(node.tagName, visitor)!;
        const _importClause = visitNode(node.importClause, visitor);
        const _moduleSpecifier = visitNode(node.moduleSpecifier, visitor)!;
        const _attributes = visitNode(node.attributes, visitor);
        return updateJSDocImportTag(node, _tagName, _importClause, _moduleSpecifier, _attributes);
    },
    [SyntaxKind.JSDocLink]: (node: JSDocLink, visitor: Visitor): JSDocLink => {
        const _name = visitNode(node.name, visitor);
        return updateJSDocLink(node, _name);
    },
    [SyntaxKind.JSDocLinkCode]: (node: JSDocLinkCode, visitor: Visitor): JSDocLinkCode => {
        const _name = visitNode(node.name, visitor);
        return updateJSDocLinkCode(node, _name);
    },
    [SyntaxKind.JSDocLinkPlain]: (node: JSDocLinkPlain, visitor: Visitor): JSDocLinkPlain => {
        const _name = visitNode(node.name, visitor);
        return updateJSDocLinkPlain(node, _name);
    },
    [SyntaxKind.JSDocMemberName]: (node: JSDocMemberName, visitor: Visitor): JSDocMemberName => {
        const _left = visitNode(node.left, visitor)!;
        const _right = visitNode(node.right, visitor)!;
        return updateJSDocMemberName(node, _left, _right);
    },
    [SyntaxKind.JSDocNameReference]: (node: JSDocNameReference, visitor: Visitor): JSDocNameReference => {
        const _name = visitNode(node.name, visitor)!;
        return updateJSDocNameReference(node, _name);
    },
    [SyntaxKind.JSDocNonNullableType]: (node: JSDocNonNullableType, visitor: Visitor): JSDocNonNullableType => {
        const _type = visitNode(node.type, visitor)!;
        return updateJSDocNonNullableType(node, _type);
    },
    [SyntaxKind.JSDocNullableType]: (node: JSDocNullableType, visitor: Visitor): JSDocNullableType => {
        const _type = visitNode(node.type, visitor)!;
        return updateJSDocNullableType(node, _type);
    },
    [SyntaxKind.JSDocOptionalType]: (node: JSDocOptionalType, visitor: Visitor): JSDocOptionalType => {
        const _type = visitNode(node.type, visitor)!;
        return updateJSDocOptionalType(node, _type);
    },
    [SyntaxKind.JSDocOverloadTag]: (node: JSDocOverloadTag, visitor: Visitor): JSDocOverloadTag => {
        const _tagName = visitNode(node.tagName, visitor)!;
        const _typeExpression = visitNode(node.typeExpression, visitor)!;
        return updateJSDocOverloadTag(node, _tagName, _typeExpression);
    },
    [SyntaxKind.JSDocOverrideTag]: (node: JSDocOverrideTag, visitor: Visitor): JSDocOverrideTag => {
        const _tagName = visitNode(node.tagName, visitor)!;
        return updateJSDocOverrideTag(node, _tagName);
    },
    [SyntaxKind.JSDocParameterTag]: (node: JSDocParameterTag, visitor: Visitor): JSDocParameterTag => {
        const _tagName = visitNode(node.tagName, visitor)!;
        return updateJSDocParameterTag(node, _tagName);
    },
    [SyntaxKind.JSDocPrivateTag]: (node: JSDocPrivateTag, visitor: Visitor): JSDocPrivateTag => {
        const _tagName = visitNode(node.tagName, visitor)!;
        return updateJSDocPrivateTag(node, _tagName);
    },
    [SyntaxKind.JSDocProtectedTag]: (node: JSDocProtectedTag, visitor: Visitor): JSDocProtectedTag => {
        const _tagName = visitNode(node.tagName, visitor)!;
        return updateJSDocProtectedTag(node, _tagName);
    },
    [SyntaxKind.JSDocPublicTag]: (node: JSDocPublicTag, visitor: Visitor): JSDocPublicTag => {
        const _tagName = visitNode(node.tagName, visitor)!;
        return updateJSDocPublicTag(node, _tagName);
    },
    [SyntaxKind.JSDocReadonlyTag]: (node: JSDocReadonlyTag, visitor: Visitor): JSDocReadonlyTag => {
        const _tagName = visitNode(node.tagName, visitor)!;
        return updateJSDocReadonlyTag(node, _tagName);
    },
    [SyntaxKind.JSDocReturnTag]: (node: JSDocReturnTag, visitor: Visitor): JSDocReturnTag => {
        const _tagName = visitNode(node.tagName, visitor)!;
        const _typeExpression = visitNode(node.typeExpression, visitor);
        return updateJSDocReturnTag(node, _tagName, _typeExpression);
    },
    [SyntaxKind.JSDocSatisfiesTag]: (node: JSDocSatisfiesTag, visitor: Visitor): JSDocSatisfiesTag => {
        const _tagName = visitNode(node.tagName, visitor)!;
        const _typeExpression = visitNode(node.typeExpression, visitor)!;
        return updateJSDocSatisfiesTag(node, _tagName, _typeExpression);
    },
    [SyntaxKind.JSDocSeeTag]: (node: JSDocSeeTag, visitor: Visitor): JSDocSeeTag => {
        const _tagName = visitNode(node.tagName, visitor)!;
        const _name = visitNode(node.name, visitor);
        return updateJSDocSeeTag(node, _tagName, _name);
    },
    [SyntaxKind.JSDocSignature]: (node: JSDocSignature, visitor: Visitor): JSDocSignature => {
        const _type = visitNode(node.type, visitor);
        return updateJSDocSignature(node, _type);
    },
    [SyntaxKind.JSDocTemplateTag]: (node: JSDocTemplateTag, visitor: Visitor): JSDocTemplateTag => {
        const _tagName = visitNode(node.tagName, visitor)!;
        const _constraint = visitNode(node.constraint, visitor);
        const _typeParameters = visitNodes(node.typeParameters, visitor);
        return updateJSDocTemplateTag(node, _tagName, _constraint, _typeParameters);
    },
    [SyntaxKind.JSDocThisTag]: (node: JSDocThisTag, visitor: Visitor): JSDocThisTag => {
        const _tagName = visitNode(node.tagName, visitor)!;
        const _typeExpression = visitNode(node.typeExpression, visitor)!;
        return updateJSDocThisTag(node, _tagName, _typeExpression);
    },
    [SyntaxKind.JSDocTypedefTag]: (node: JSDocTypedefTag, visitor: Visitor): JSDocTypedefTag => {
        const _tagName = visitNode(node.tagName, visitor)!;
        const _typeExpression = visitNode(node.typeExpression, visitor);
        const _fullName = visitNode(node.fullName, visitor);
        return updateJSDocTypedefTag(node, _tagName, _typeExpression, _fullName);
    },
    [SyntaxKind.JSDocTypeExpression]: (node: JSDocTypeExpression, visitor: Visitor): JSDocTypeExpression => {
        const _type = visitNode(node.type, visitor)!;
        return updateJSDocTypeExpression(node, _type);
    },
    [SyntaxKind.JSDocTypeTag]: (node: JSDocTypeTag, visitor: Visitor): JSDocTypeTag => {
        const _tagName = visitNode(node.tagName, visitor)!;
        const _typeExpression = visitNode(node.typeExpression, visitor)!;
        return updateJSDocTypeTag(node, _tagName, _typeExpression);
    },
    [SyntaxKind.JSDocTag]: (node: JSDocUnknownTag, visitor: Visitor): JSDocUnknownTag => {
        const _tagName = visitNode(node.tagName, visitor)!;
        return updateJSDocUnknownTag(node, _tagName);
    },
    [SyntaxKind.JSDocVariadicType]: (node: JSDocVariadicType, visitor: Visitor): JSDocVariadicType => {
        const _type = visitNode(node.type, visitor)!;
        return updateJSDocVariadicType(node, _type);
    },
    [SyntaxKind.JsxAttribute]: (node: JsxAttribute, visitor: Visitor): JsxAttribute => {
        const _name = visitNode(node.name, visitor)!;
        const _initializer = visitNode(node.initializer, visitor);
        return updateJsxAttribute(node, _name, _initializer);
    },
    [SyntaxKind.JsxAttributes]: (node: JsxAttributes, visitor: Visitor): JsxAttributes => {
        const _properties = visitNodes(node.properties, visitor);
        return updateJsxAttributes(node, _properties);
    },
    [SyntaxKind.JsxClosingElement]: (node: JsxClosingElement, visitor: Visitor): JsxClosingElement => {
        const _tagName = visitNode(node.tagName, visitor)!;
        return updateJsxClosingElement(node, _tagName);
    },
    [SyntaxKind.JsxElement]: (node: JsxElement, visitor: Visitor): JsxElement => {
        const _openingElement = visitNode(node.openingElement, visitor)!;
        const _children = visitNodes(node.children, visitor);
        const _closingElement = visitNode(node.closingElement, visitor)!;
        return updateJsxElement(node, _openingElement, _children, _closingElement);
    },
    [SyntaxKind.JsxExpression]: (node: JsxExpression, visitor: Visitor): JsxExpression => {
        const _dotDotDotToken = visitNode(node.dotDotDotToken, visitor);
        const _expression = visitNode(node.expression, visitor);
        return updateJsxExpression(node, _dotDotDotToken, _expression);
    },
    [SyntaxKind.JsxFragment]: (node: JsxFragment, visitor: Visitor): JsxFragment => {
        const _openingFragment = visitNode(node.openingFragment, visitor)!;
        const _children = visitNodes(node.children, visitor);
        const _closingFragment = visitNode(node.closingFragment, visitor)!;
        return updateJsxFragment(node, _openingFragment, _children, _closingFragment);
    },
    [SyntaxKind.JsxNamespacedName]: (node: JsxNamespacedName, visitor: Visitor): JsxNamespacedName => {
        const _name = visitNode(node.name, visitor)!;
        const _namespace = visitNode(node.namespace, visitor)!;
        return updateJsxNamespacedName(node, _name, _namespace);
    },
    [SyntaxKind.JsxOpeningElement]: (node: JsxOpeningElement, visitor: Visitor): JsxOpeningElement => {
        const _tagName = visitNode(node.tagName, visitor)!;
        const _typeArguments = visitNodes(node.typeArguments, visitor);
        const _attributes = visitNode(node.attributes, visitor)!;
        return updateJsxOpeningElement(node, _tagName, _typeArguments, _attributes);
    },
    [SyntaxKind.JsxSelfClosingElement]: (node: JsxSelfClosingElement, visitor: Visitor): JsxSelfClosingElement => {
        const _tagName = visitNode(node.tagName, visitor)!;
        const _typeArguments = visitNodes(node.typeArguments, visitor);
        const _attributes = visitNode(node.attributes, visitor)!;
        return updateJsxSelfClosingElement(node, _tagName, _typeArguments, _attributes);
    },
    [SyntaxKind.JsxSpreadAttribute]: (node: JsxSpreadAttribute, visitor: Visitor): JsxSpreadAttribute => {
        const _expression = visitNode(node.expression, visitor)!;
        return updateJsxSpreadAttribute(node, _expression);
    },
    [SyntaxKind.LabeledStatement]: (node: LabeledStatement, visitor: Visitor): LabeledStatement => {
        const _label = visitNode(node.label, visitor)!;
        const _statement = visitNode(node.statement, visitor)!;
        return updateLabeledStatement(node, _label, _statement);
    },
    [SyntaxKind.LiteralType]: (node: LiteralTypeNode, visitor: Visitor): LiteralTypeNode => {
        const _literal = visitNode(node.literal, visitor)!;
        return updateLiteralTypeNode(node, _literal);
    },
    [SyntaxKind.MappedType]: (node: MappedTypeNode, visitor: Visitor): MappedTypeNode => {
        const _readonlyToken = visitNode(node.readonlyToken, visitor);
        const _typeParameter = visitNode(node.typeParameter, visitor)!;
        const _nameType = visitNode(node.nameType, visitor);
        const _questionToken = visitNode(node.questionToken, visitor);
        const _type = visitNode(node.type, visitor);
        const _members = visitNodes(node.members, visitor);
        return updateMappedTypeNode(node, _readonlyToken, _typeParameter, _nameType, _questionToken, _type, _members);
    },
    [SyntaxKind.MetaProperty]: (node: MetaProperty, visitor: Visitor): MetaProperty => {
        const _name = visitNode(node.name, visitor)!;
        return updateMetaProperty(node, _name);
    },
    [SyntaxKind.MethodDeclaration]: (node: MethodDeclaration, visitor: Visitor): MethodDeclaration => {
        const _modifiers = visitNodes(node.modifiers, visitor);
        const _asteriskToken = visitNode(node.asteriskToken, visitor);
        const _name = visitNode(node.name, visitor)!;
        const _postfixToken = visitNode(node.postfixToken, visitor);
        const _typeParameters = visitNodes(node.typeParameters, visitor);
        const _parameters = visitNodes(node.parameters, visitor);
        const _type = visitNode(node.type, visitor);
        const _body = visitNode(node.body, visitor);
        return updateMethodDeclaration(node, _modifiers, _asteriskToken, _name, _postfixToken, _typeParameters, _parameters, _type, _body);
    },
    [SyntaxKind.MethodSignature]: (node: MethodSignature, visitor: Visitor): MethodSignature => {
        const _modifiers = visitNodes(node.modifiers, visitor);
        const _name = visitNode(node.name, visitor)!;
        const _postfixToken = visitNode(node.postfixToken, visitor);
        const _typeParameters = visitNodes(node.typeParameters, visitor);
        const _parameters = visitNodes(node.parameters, visitor);
        const _type = visitNode(node.type, visitor);
        return updateMethodSignature(node, _modifiers, _name, _postfixToken, _typeParameters, _parameters, _type);
    },
    [SyntaxKind.ModuleBlock]: (node: ModuleBlock, visitor: Visitor): ModuleBlock => {
        const _statements = visitNodes(node.statements, visitor);
        return updateModuleBlock(node, _statements);
    },
    [SyntaxKind.ModuleDeclaration]: (node: ModuleDeclaration, visitor: Visitor): ModuleDeclaration => {
        const _modifiers = visitNodes(node.modifiers, visitor);
        const _name = visitNode(node.name, visitor)!;
        return updateModuleDeclaration(node, _modifiers, _name);
    },
    [SyntaxKind.NamedExports]: (node: NamedExports, visitor: Visitor): NamedExports => {
        const _elements = visitNodes(node.elements, visitor);
        return updateNamedExports(node, _elements);
    },
    [SyntaxKind.NamedImports]: (node: NamedImports, visitor: Visitor): NamedImports => {
        const _elements = visitNodes(node.elements, visitor);
        return updateNamedImports(node, _elements);
    },
    [SyntaxKind.NamedTupleMember]: (node: NamedTupleMember, visitor: Visitor): NamedTupleMember => {
        const _dotDotDotToken = visitNode(node.dotDotDotToken, visitor);
        const _name = visitNode(node.name, visitor)!;
        const _questionToken = visitNode(node.questionToken, visitor);
        const _type = visitNode(node.type, visitor)!;
        return updateNamedTupleMember(node, _dotDotDotToken, _name, _questionToken, _type);
    },
    [SyntaxKind.NamespaceExport]: (node: NamespaceExport, visitor: Visitor): NamespaceExport => {
        const _name = visitNode(node.name, visitor)!;
        return updateNamespaceExport(node, _name);
    },
    [SyntaxKind.NamespaceExportDeclaration]: (node: NamespaceExportDeclaration, visitor: Visitor): NamespaceExportDeclaration => {
        const _name = visitNode(node.name, visitor)!;
        return updateNamespaceExportDeclaration(node, _name);
    },
    [SyntaxKind.NamespaceImport]: (node: NamespaceImport, visitor: Visitor): NamespaceImport => {
        const _name = visitNode(node.name, visitor)!;
        return updateNamespaceImport(node, _name);
    },
    [SyntaxKind.NewExpression]: (node: NewExpression, visitor: Visitor): NewExpression => {
        const _expression = visitNode(node.expression, visitor)!;
        const _typeArguments = visitNodes(node.typeArguments, visitor);
        const _arguments = visitNodes(node.arguments, visitor);
        return updateNewExpression(node, _expression, _typeArguments, _arguments);
    },
    [SyntaxKind.NonNullExpression]: (node: NonNullExpression, visitor: Visitor): NonNullExpression => {
        const _expression = visitNode(node.expression, visitor)!;
        return updateNonNullExpression(node, _expression);
    },
    [SyntaxKind.ObjectBindingPattern]: (node: ObjectBindingPattern, visitor: Visitor): ObjectBindingPattern => {
        const _elements = visitNodes(node.elements, visitor);
        return updateObjectBindingPattern(node, _elements);
    },
    [SyntaxKind.ObjectLiteralExpression]: (node: ObjectLiteralExpression, visitor: Visitor): ObjectLiteralExpression => {
        const _properties = visitNodes(node.properties, visitor);
        return updateObjectLiteralExpression(node, _properties);
    },
    [SyntaxKind.OptionalType]: (node: OptionalTypeNode, visitor: Visitor): OptionalTypeNode => {
        const _type = visitNode(node.type, visitor)!;
        return updateOptionalTypeNode(node, _type);
    },
    [SyntaxKind.Parameter]: (node: ParameterDeclaration, visitor: Visitor): ParameterDeclaration => {
        const _modifiers = visitNodes(node.modifiers, visitor);
        const _dotDotDotToken = visitNode(node.dotDotDotToken, visitor);
        const _questionToken = visitNode(node.questionToken, visitor);
        const _type = visitNode(node.type, visitor);
        const _initializer = visitNode(node.initializer, visitor);
        return updateParameterDeclaration(node, _modifiers, _dotDotDotToken, _questionToken, _type, _initializer);
    },
    [SyntaxKind.ParenthesizedExpression]: (node: ParenthesizedExpression, visitor: Visitor): ParenthesizedExpression => {
        const _expression = visitNode(node.expression, visitor)!;
        return updateParenthesizedExpression(node, _expression);
    },
    [SyntaxKind.ParenthesizedType]: (node: ParenthesizedTypeNode, visitor: Visitor): ParenthesizedTypeNode => {
        const _type = visitNode(node.type, visitor)!;
        return updateParenthesizedTypeNode(node, _type);
    },
    [SyntaxKind.PartiallyEmittedExpression]: (node: PartiallyEmittedExpression, visitor: Visitor): PartiallyEmittedExpression => {
        const _expression = visitNode(node.expression, visitor)!;
        return updatePartiallyEmittedExpression(node, _expression);
    },
    [SyntaxKind.PostfixUnaryExpression]: (node: PostfixUnaryExpression, visitor: Visitor): PostfixUnaryExpression => {
        const _operand = visitNode(node.operand, visitor)!;
        return updatePostfixUnaryExpression(node, _operand);
    },
    [SyntaxKind.PrefixUnaryExpression]: (node: PrefixUnaryExpression, visitor: Visitor): PrefixUnaryExpression => {
        const _operand = visitNode(node.operand, visitor)!;
        return updatePrefixUnaryExpression(node, _operand);
    },
    [SyntaxKind.PropertyAccessExpression]: (node: PropertyAccessExpression, visitor: Visitor): PropertyAccessExpression => {
        const _expression = visitNode(node.expression, visitor)!;
        const _questionDotToken = visitNode(node.questionDotToken, visitor);
        const _name = visitNode(node.name, visitor)!;
        return updatePropertyAccessExpression(node, _expression, _questionDotToken, _name);
    },
    [SyntaxKind.PropertyAssignment]: (node: PropertyAssignment, visitor: Visitor): PropertyAssignment => {
        const _name = visitNode(node.name, visitor)!;
        const _postfixToken = visitNode(node.postfixToken, visitor);
        const _initializer = visitNode(node.initializer, visitor)!;
        return updatePropertyAssignment(node, _name, _postfixToken, _initializer);
    },
    [SyntaxKind.PropertyDeclaration]: (node: PropertyDeclaration, visitor: Visitor): PropertyDeclaration => {
        const _modifiers = visitNodes(node.modifiers, visitor);
        const _name = visitNode(node.name, visitor)!;
        const _postfixToken = visitNode(node.postfixToken, visitor);
        const _type = visitNode(node.type, visitor);
        const _initializer = visitNode(node.initializer, visitor);
        return updatePropertyDeclaration(node, _modifiers, _name, _postfixToken, _type, _initializer);
    },
    [SyntaxKind.PropertySignature]: (node: PropertySignature, visitor: Visitor): PropertySignature => {
        const _modifiers = visitNodes(node.modifiers, visitor);
        const _name = visitNode(node.name, visitor)!;
        const _postfixToken = visitNode(node.postfixToken, visitor);
        const _type = visitNode(node.type, visitor);
        return updatePropertySignature(node, _modifiers, _name, _postfixToken, _type);
    },
    [SyntaxKind.QualifiedName]: (node: QualifiedName, visitor: Visitor): QualifiedName => {
        const _left = visitNode(node.left, visitor)!;
        const _right = visitNode(node.right, visitor)!;
        return updateQualifiedName(node, _left, _right);
    },
    [SyntaxKind.RestType]: (node: RestTypeNode, visitor: Visitor): RestTypeNode => {
        const _type = visitNode(node.type, visitor)!;
        return updateRestTypeNode(node, _type);
    },
    [SyntaxKind.ReturnStatement]: (node: ReturnStatement, visitor: Visitor): ReturnStatement => {
        const _expression = visitNode(node.expression, visitor);
        return updateReturnStatement(node, _expression);
    },
    [SyntaxKind.SatisfiesExpression]: (node: SatisfiesExpression, visitor: Visitor): SatisfiesExpression => {
        const _expression = visitNode(node.expression, visitor)!;
        const _type = visitNode(node.type, visitor)!;
        return updateSatisfiesExpression(node, _expression, _type);
    },
    [SyntaxKind.SetAccessor]: (node: SetAccessorDeclaration, visitor: Visitor): SetAccessorDeclaration => {
        const _modifiers = visitNodes(node.modifiers, visitor);
        const _name = visitNode(node.name, visitor)!;
        const _parameters = visitNodes(node.parameters, visitor);
        const _body = visitNode(node.body, visitor);
        return updateSetAccessorDeclaration(node, _modifiers, _name, _parameters, _body);
    },
    [SyntaxKind.ShorthandPropertyAssignment]: (node: ShorthandPropertyAssignment, visitor: Visitor): ShorthandPropertyAssignment => {
        const _name = visitNode(node.name, visitor)!;
        const _postfixToken = visitNode(node.postfixToken, visitor);
        const _equalsToken = visitNode(node.equalsToken, visitor);
        const _objectAssignmentInitializer = visitNode(node.objectAssignmentInitializer, visitor);
        return updateShorthandPropertyAssignment(node, _name, _postfixToken, _equalsToken, _objectAssignmentInitializer);
    },
    [SyntaxKind.SourceFile]: (node: SourceFile, visitor: Visitor): SourceFile => {
        const _statements = visitNodes(node.statements, visitor);
        const _endOfFileToken = visitNode(node.endOfFileToken, visitor)!;
        return updateSourceFile(node, _statements, _endOfFileToken);
    },
    [SyntaxKind.SpreadAssignment]: (node: SpreadAssignment, visitor: Visitor): SpreadAssignment => {
        const _expression = visitNode(node.expression, visitor)!;
        return updateSpreadAssignment(node, _expression);
    },
    [SyntaxKind.SpreadElement]: (node: SpreadElement, visitor: Visitor): SpreadElement => {
        const _expression = visitNode(node.expression, visitor)!;
        return updateSpreadElement(node, _expression);
    },
    [SyntaxKind.SwitchStatement]: (node: SwitchStatement, visitor: Visitor): SwitchStatement => {
        const _expression = visitNode(node.expression, visitor)!;
        const _caseBlock = visitNode(node.caseBlock, visitor)!;
        return updateSwitchStatement(node, _expression, _caseBlock);
    },
    [SyntaxKind.TaggedTemplateExpression]: (node: TaggedTemplateExpression, visitor: Visitor): TaggedTemplateExpression => {
        const _tag = visitNode(node.tag, visitor)!;
        const _typeArguments = visitNodes(node.typeArguments, visitor);
        const _template = visitNode(node.template, visitor)!;
        return updateTaggedTemplateExpression(node, _tag, _typeArguments, _template);
    },
    [SyntaxKind.TemplateExpression]: (node: TemplateExpression, visitor: Visitor): TemplateExpression => {
        const _head = visitNode(node.head, visitor)!;
        const _templateSpans = visitNodes(node.templateSpans, visitor);
        return updateTemplateExpression(node, _head, _templateSpans);
    },
    [SyntaxKind.TemplateLiteralType]: (node: TemplateLiteralTypeNode, visitor: Visitor): TemplateLiteralTypeNode => {
        const _head = visitNode(node.head, visitor)!;
        const _templateSpans = visitNodes(node.templateSpans, visitor);
        return updateTemplateLiteralTypeNode(node, _head, _templateSpans);
    },
    [SyntaxKind.TemplateLiteralTypeSpan]: (node: TemplateLiteralTypeSpan, visitor: Visitor): TemplateLiteralTypeSpan => {
        const _type = visitNode(node.type, visitor)!;
        const _literal = visitNode(node.literal, visitor)!;
        return updateTemplateLiteralTypeSpan(node, _type, _literal);
    },
    [SyntaxKind.TemplateSpan]: (node: TemplateSpan, visitor: Visitor): TemplateSpan => {
        const _expression = visitNode(node.expression, visitor)!;
        const _literal = visitNode(node.literal, visitor)!;
        return updateTemplateSpan(node, _expression, _literal);
    },
    [SyntaxKind.ThrowStatement]: (node: ThrowStatement, visitor: Visitor): ThrowStatement => {
        const _expression = visitNode(node.expression, visitor)!;
        return updateThrowStatement(node, _expression);
    },
    [SyntaxKind.TryStatement]: (node: TryStatement, visitor: Visitor): TryStatement => {
        const _tryBlock = visitNode(node.tryBlock, visitor)!;
        const _catchClause = visitNode(node.catchClause, visitor);
        const _finallyBlock = visitNode(node.finallyBlock, visitor);
        return updateTryStatement(node, _tryBlock, _catchClause, _finallyBlock);
    },
    [SyntaxKind.TupleType]: (node: TupleTypeNode, visitor: Visitor): TupleTypeNode => {
        const _elements = visitNodes(node.elements, visitor);
        return updateTupleTypeNode(node, _elements);
    },
    [SyntaxKind.TypeAliasDeclaration]: (node: TypeAliasDeclaration, visitor: Visitor): TypeAliasDeclaration => {
        const _modifiers = visitNodes(node.modifiers, visitor);
        const _name = visitNode(node.name, visitor)!;
        const _typeParameters = visitNodes(node.typeParameters, visitor);
        const _type = visitNode(node.type, visitor)!;
        return updateTypeAliasDeclaration(node, _modifiers, _name, _typeParameters, _type);
    },
    [SyntaxKind.TypeAssertionExpression]: (node: TypeAssertion, visitor: Visitor): TypeAssertion => {
        const _type = visitNode(node.type, visitor)!;
        const _expression = visitNode(node.expression, visitor)!;
        return updateTypeAssertion(node, _type, _expression);
    },
    [SyntaxKind.TypeLiteral]: (node: TypeLiteralNode, visitor: Visitor): TypeLiteralNode => {
        const _members = visitNodes(node.members, visitor);
        return updateTypeLiteralNode(node, _members);
    },
    [SyntaxKind.TypeOfExpression]: (node: TypeOfExpression, visitor: Visitor): TypeOfExpression => {
        const _expression = visitNode(node.expression, visitor)!;
        return updateTypeOfExpression(node, _expression);
    },
    [SyntaxKind.TypeOperator]: (node: TypeOperatorNode, visitor: Visitor): TypeOperatorNode => {
        const _type = visitNode(node.type, visitor)!;
        return updateTypeOperatorNode(node, _type);
    },
    [SyntaxKind.TypeParameter]: (node: TypeParameterDeclaration, visitor: Visitor): TypeParameterDeclaration => {
        const _modifiers = visitNodes(node.modifiers, visitor);
        const _name = visitNode(node.name, visitor)!;
        const _constraint = visitNode(node.constraint, visitor);
        const _default = visitNode(node.default, visitor);
        return updateTypeParameterDeclaration(node, _modifiers, _name, _constraint, _default);
    },
    [SyntaxKind.TypePredicate]: (node: TypePredicateNode, visitor: Visitor): TypePredicateNode => {
        const _assertsModifier = visitNode(node.assertsModifier, visitor);
        const _parameterName = visitNode(node.parameterName, visitor)!;
        const _type = visitNode(node.type, visitor);
        return updateTypePredicateNode(node, _assertsModifier, _parameterName, _type);
    },
    [SyntaxKind.TypeQuery]: (node: TypeQueryNode, visitor: Visitor): TypeQueryNode => {
        const _exprName = visitNode(node.exprName, visitor)!;
        const _typeArguments = visitNodes(node.typeArguments, visitor);
        return updateTypeQueryNode(node, _exprName, _typeArguments);
    },
    [SyntaxKind.TypeReference]: (node: TypeReferenceNode, visitor: Visitor): TypeReferenceNode => {
        const _typeName = visitNode(node.typeName, visitor)!;
        const _typeArguments = visitNodes(node.typeArguments, visitor);
        return updateTypeReferenceNode(node, _typeName, _typeArguments);
    },
    [SyntaxKind.UnionType]: (node: UnionTypeNode, visitor: Visitor): UnionTypeNode => {
        const _types = visitNodes(node.types, visitor);
        return updateUnionTypeNode(node, _types);
    },
    [SyntaxKind.VariableDeclaration]: (node: VariableDeclaration, visitor: Visitor): VariableDeclaration => {
        const _exclamationToken = visitNode(node.exclamationToken, visitor);
        const _type = visitNode(node.type, visitor);
        const _initializer = visitNode(node.initializer, visitor);
        return updateVariableDeclaration(node, _exclamationToken, _type, _initializer);
    },
    [SyntaxKind.VariableDeclarationList]: (node: VariableDeclarationList, visitor: Visitor): VariableDeclarationList => {
        const _declarations = visitNodes(node.declarations, visitor);
        return updateVariableDeclarationList(node, _declarations);
    },
    [SyntaxKind.VariableStatement]: (node: VariableStatement, visitor: Visitor): VariableStatement => {
        const _modifiers = visitNodes(node.modifiers, visitor);
        const _declarationList = visitNode(node.declarationList, visitor)!;
        return updateVariableStatement(node, _modifiers, _declarationList);
    },
    [SyntaxKind.VoidExpression]: (node: VoidExpression, visitor: Visitor): VoidExpression => {
        const _expression = visitNode(node.expression, visitor)!;
        return updateVoidExpression(node, _expression);
    },
    [SyntaxKind.WhileStatement]: (node: WhileStatement, visitor: Visitor): WhileStatement => {
        const _expression = visitNode(node.expression, visitor)!;
        const _statement = visitNode(node.statement, visitor)!;
        return updateWhileStatement(node, _expression, _statement);
    },
    [SyntaxKind.WithStatement]: (node: WithStatement, visitor: Visitor): WithStatement => {
        const _expression = visitNode(node.expression, visitor)!;
        const _statement = visitNode(node.statement, visitor)!;
        return updateWithStatement(node, _expression, _statement);
    },
    [SyntaxKind.YieldExpression]: (node: YieldExpression, visitor: Visitor): YieldExpression => {
        const _asteriskToken = visitNode(node.asteriskToken, visitor);
        const _expression = visitNode(node.expression, visitor);
        return updateYieldExpression(node, _asteriskToken, _expression);
    },
};
