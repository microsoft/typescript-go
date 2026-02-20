// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
// !!! THIS FILE IS AUTO-GENERATED — DO NOT EDIT !!!
// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
//
// Source: _packages/ast/src/nodes.ts
// Generator: _packages/ast/scripts/generateFactory.ts
//

import { SyntaxKind } from "#syntaxKind";
import { TokenFlags } from "#tokenFlags";
import type {
    ArrayBindingElement,
    ArrayBindingPattern,
    ArrayLiteralExpression,
    ArrayTypeNode,
    ArrowFunction,
    AsExpression,
    AssertClause,
    AssertsKeyword,
    AsteriskToken,
    AwaitExpression,
    AwaitKeyword,
    BigIntLiteral,
    BinaryExpression,
    BinaryOperatorToken,
    BindingElement,
    BindingName,
    Block,
    BooleanLiteral,
    BreakStatement,
    CallExpression,
    CallSignatureDeclaration,
    CaseBlock,
    CaseClause,
    CaseOrDefaultClause,
    CatchClause,
    ClassDeclaration,
    ClassElement,
    ClassExpression,
    ClassStaticBlockDeclaration,
    ColonToken,
    CommaListExpression,
    ComputedPropertyName,
    ConciseBody,
    ConditionalExpression,
    ConditionalTypeNode,
    ConstructorDeclaration,
    ConstructorTypeNode,
    ConstructSignatureDeclaration,
    ContinueStatement,
    DebuggerStatement,
    Decorator,
    DefaultClause,
    DeleteExpression,
    DoStatement,
    DotDotDotToken,
    ElementAccessExpression,
    EmptyStatement,
    EndOfFile,
    EntityName,
    EnumDeclaration,
    EnumMember,
    EqualsGreaterThanToken,
    EqualsToken,
    ExclamationToken,
    ExportAssignment,
    ExportDeclaration,
    ExportSpecifier,
    Expression,
    ExpressionStatement,
    ExpressionWithTypeArguments,
    ExternalModuleReference,
    FalseLiteral,
    ForInitializer,
    ForInStatement,
    ForOfStatement,
    ForStatement,
    FunctionBody,
    FunctionDeclaration,
    FunctionExpression,
    FunctionTypeNode,
    GetAccessorDeclaration,
    HeritageClause,
    Identifier,
    IfStatement,
    ImportAttribute,
    ImportAttributeName,
    ImportAttributes,
    ImportClause,
    ImportDeclaration,
    ImportEqualsDeclaration,
    ImportExpression,
    ImportPhaseModifierSyntaxKind,
    ImportSpecifier,
    ImportTypeNode,
    IndexedAccessTypeNode,
    IndexSignatureDeclaration,
    InferTypeNode,
    InterfaceDeclaration,
    IntersectionTypeNode,
    JSDoc,
    JSDocAllType,
    JSDocAugmentsTag,
    JSDocCallbackTag,
    JSDocComment,
    JSDocDeprecatedTag,
    JSDocImplementsTag,
    JSDocImportTag,
    JSDocLink,
    JSDocLinkCode,
    JSDocLinkPlain,
    JSDocMemberName,
    JSDocNameReference,
    JSDocNamespaceDeclaration,
    JSDocNonNullableType,
    JSDocNullableType,
    JSDocOptionalType,
    JSDocOverloadTag,
    JSDocOverrideTag,
    JSDocParameterTag,
    JSDocPrivateTag,
    JSDocPropertyLikeTag,
    JSDocPropertyTag,
    JSDocProtectedTag,
    JSDocPublicTag,
    JSDocReadonlyTag,
    JSDocReturnTag,
    JSDocSatisfiesTag,
    JSDocSeeTag,
    JSDocSignature,
    JSDocTag,
    JSDocTemplateTag,
    JSDocText,
    JSDocThisTag,
    JSDocTypedefTag,
    JSDocTypeExpression,
    JSDocTypeLiteral,
    JSDocTypeTag,
    JSDocUnknownTag,
    JSDocVariadicType,
    JsxAttribute,
    JsxAttributeLike,
    JsxAttributeName,
    JsxAttributes,
    JsxAttributeValue,
    JsxChild,
    JsxClosingElement,
    JsxClosingFragment,
    JsxElement,
    JsxExpression,
    JsxFragment,
    JsxNamespacedName,
    JsxOpeningElement,
    JsxOpeningFragment,
    JsxSelfClosingElement,
    JsxSpreadAttribute,
    JsxTagNameExpression,
    JsxText,
    LabeledStatement,
    LeftHandSideExpression,
    LiteralExpression,
    LiteralTypeNode,
    MappedTypeNode,
    MemberName,
    MetaProperty,
    MethodDeclaration,
    MethodSignature,
    MinusToken,
    Modifier,
    ModifierLike,
    ModuleBlock,
    ModuleBody,
    ModuleDeclaration,
    ModuleExportName,
    ModuleName,
    ModuleReference,
    NamedExportBindings,
    NamedExports,
    NamedImportBindings,
    NamedImports,
    NamedTupleMember,
    NamespaceExport,
    NamespaceExportDeclaration,
    NamespaceImport,
    NewExpression,
    Node,
    NodeArray,
    NonNullExpression,
    NoSubstitutionTemplateLiteral,
    NullLiteral,
    NumericLiteral,
    ObjectBindingPattern,
    ObjectLiteralElementLike,
    ObjectLiteralExpression,
    OmittedExpression,
    OptionalTypeNode,
    ParameterDeclaration,
    ParenthesizedExpression,
    ParenthesizedTypeNode,
    PartiallyEmittedExpression,
    Path,
    PlusToken,
    PostfixUnaryExpression,
    PostfixUnaryOperator,
    PrefixUnaryExpression,
    PrefixUnaryOperator,
    PrivateIdentifier,
    PropertyAccessEntityNameExpression,
    PropertyAccessExpression,
    PropertyAssignment,
    PropertyDeclaration,
    PropertyName,
    PropertySignature,
    QualifiedName,
    QuestionDotToken,
    QuestionToken,
    ReadonlyKeyword,
    RegularExpressionLiteral,
    RestTypeNode,
    ReturnStatement,
    SatisfiesExpression,
    SemicolonClassElement,
    SetAccessorDeclaration,
    ShorthandPropertyAssignment,
    SourceFile,
    SpreadAssignment,
    SpreadElement,
    Statement,
    StringLiteral,
    SuperExpression,
    SwitchStatement,
    TaggedTemplateExpression,
    TemplateExpression,
    TemplateHead,
    TemplateLiteral,
    TemplateLiteralTypeNode,
    TemplateLiteralTypeSpan,
    TemplateMiddle,
    TemplateSpan,
    TemplateTail,
    ThisExpression,
    ThisTypeNode,
    ThrowStatement,
    Token,
    TrueLiteral,
    TryStatement,
    TupleTypeNode,
    TypeAliasDeclaration,
    TypeAssertion,
    TypeElement,
    TypeLiteralNode,
    TypeNode,
    TypeOfExpression,
    TypeOperatorNode,
    TypeParameterDeclaration,
    TypePredicateNode,
    TypeQueryNode,
    TypeReferenceNode,
    UnaryExpression,
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
 * Monomorphic AST node implementation.
 * All synthetic nodes share the same V8 hidden class for optimal property access.
 *
 * Common fields live directly on the object; kind-specific fields are stored
 * in the `_data` bag and accessed via generated property accessors.
 */
export class NodeObject {
    readonly kind!: SyntaxKind;
    readonly pos!: number;
    readonly end!: number;
    readonly parent!: Node;
    /** @internal */
    _data: any;

    constructor(kind: SyntaxKind, data: any) {
        this.kind = kind;
        this.pos = -1;
        this.end = -1;
        this.parent = undefined!;
        this._data = data;
    }

    get argument(): any {
        return this._data?.argument;
    }
    get argumentExpression(): any {
        return this._data?.argumentExpression;
    }
    get arguments(): any {
        return this._data?.arguments;
    }
    get assertClause(): any {
        return this._data?.assertClause;
    }
    get assertsModifier(): any {
        return this._data?.assertsModifier;
    }
    get asteriskToken(): any {
        return this._data?.asteriskToken;
    }
    get attributes(): any {
        return this._data?.attributes;
    }
    get awaitModifier(): any {
        return this._data?.awaitModifier;
    }
    get block(): any {
        return this._data?.block;
    }
    get body(): any {
        return this._data?.body;
    }
    get caseBlock(): any {
        return this._data?.caseBlock;
    }
    get catchClause(): any {
        return this._data?.catchClause;
    }
    get checkType(): any {
        return this._data?.checkType;
    }
    get children(): any {
        return this._data?.children;
    }
    get class(): any {
        return this._data?.class;
    }
    get clauses(): any {
        return this._data?.clauses;
    }
    get closingElement(): any {
        return this._data?.closingElement;
    }
    get closingFragment(): any {
        return this._data?.closingFragment;
    }
    get colonToken(): any {
        return this._data?.colonToken;
    }
    get comment(): any {
        return this._data?.comment;
    }
    get condition(): any {
        return this._data?.condition;
    }
    get constraint(): any {
        return this._data?.constraint;
    }
    get containsOnlyTriviaWhiteSpaces(): any {
        return this._data?.containsOnlyTriviaWhiteSpaces;
    }
    get declarationList(): any {
        return this._data?.declarationList;
    }
    get declarations(): any {
        return this._data?.declarations;
    }
    get default(): any {
        return this._data?.default;
    }
    get dotDotDotToken(): any {
        return this._data?.dotDotDotToken;
    }
    get elementType(): any {
        return this._data?.elementType;
    }
    get elements(): any {
        return this._data?.elements;
    }
    get elseStatement(): any {
        return this._data?.elseStatement;
    }
    get endOfFileToken(): any {
        return this._data?.endOfFileToken;
    }
    get equalsGreaterThanToken(): any {
        return this._data?.equalsGreaterThanToken;
    }
    get equalsToken(): any {
        return this._data?.equalsToken;
    }
    get escapedText(): any {
        return this._data?.escapedText;
    }
    get exclamationToken(): any {
        return this._data?.exclamationToken;
    }
    get exportClause(): any {
        return this._data?.exportClause;
    }
    get exprName(): any {
        return this._data?.exprName;
    }
    get expression(): any {
        return this._data?.expression;
    }
    get extendsType(): any {
        return this._data?.extendsType;
    }
    get falseType(): any {
        return this._data?.falseType;
    }
    get fileName(): any {
        return this._data?.fileName;
    }
    get finallyBlock(): any {
        return this._data?.finallyBlock;
    }
    get fullName(): any {
        return this._data?.fullName;
    }
    get hasExtendedUnicodeEscape(): any {
        return this._data?.hasExtendedUnicodeEscape;
    }
    get head(): any {
        return this._data?.head;
    }
    get heritageClauses(): any {
        return this._data?.heritageClauses;
    }
    get importClause(): any {
        return this._data?.importClause;
    }
    get incrementor(): any {
        return this._data?.incrementor;
    }
    get indexType(): any {
        return this._data?.indexType;
    }
    get initializer(): any {
        return this._data?.initializer;
    }
    get isArrayType(): any {
        return this._data?.isArrayType;
    }
    get isBracketed(): any {
        return this._data?.isBracketed;
    }
    get isExportEquals(): any {
        return this._data?.isExportEquals;
    }
    get isNameFirst(): any {
        return this._data?.isNameFirst;
    }
    get isTypeOf(): any {
        return this._data?.isTypeOf;
    }
    get isTypeOnly(): any {
        return this._data?.isTypeOnly;
    }
    get isUnterminated(): any {
        return this._data?.isUnterminated;
    }
    get jsDocPropertyTags(): any {
        return this._data?.jsDocPropertyTags;
    }
    get keywordToken(): any {
        return this._data?.keywordToken;
    }
    get label(): any {
        return this._data?.label;
    }
    get left(): any {
        return this._data?.left;
    }
    get literal(): any {
        return this._data?.literal;
    }
    get members(): any {
        return this._data?.members;
    }
    get modifiers(): any {
        return this._data?.modifiers;
    }
    get moduleReference(): any {
        return this._data?.moduleReference;
    }
    get moduleSpecifier(): any {
        return this._data?.moduleSpecifier;
    }
    get multiLine(): any {
        return this._data?.multiLine;
    }
    get name(): any {
        return this._data?.name;
    }
    get nameType(): any {
        return this._data?.nameType;
    }
    get namedBindings(): any {
        return this._data?.namedBindings;
    }
    get namespace(): any {
        return this._data?.namespace;
    }
    get numericLiteralFlags(): any {
        return this._data?.numericLiteralFlags;
    }
    get objectAssignmentInitializer(): any {
        return this._data?.objectAssignmentInitializer;
    }
    get objectType(): any {
        return this._data?.objectType;
    }
    get openingElement(): any {
        return this._data?.openingElement;
    }
    get openingFragment(): any {
        return this._data?.openingFragment;
    }
    get operand(): any {
        return this._data?.operand;
    }
    get operator(): any {
        return this._data?.operator;
    }
    get operatorToken(): any {
        return this._data?.operatorToken;
    }
    get parameterName(): any {
        return this._data?.parameterName;
    }
    get parameters(): any {
        return this._data?.parameters;
    }
    get path(): any {
        return this._data?.path;
    }
    get phaseModifier(): any {
        return this._data?.phaseModifier;
    }
    get possiblyExhaustive(): any {
        return this._data?.possiblyExhaustive;
    }
    get postfix(): any {
        return this._data?.postfix;
    }
    get properties(): any {
        return this._data?.properties;
    }
    get propertyName(): any {
        return this._data?.propertyName;
    }
    get qualifier(): any {
        return this._data?.qualifier;
    }
    get questionDotToken(): any {
        return this._data?.questionDotToken;
    }
    get questionToken(): any {
        return this._data?.questionToken;
    }
    get rawText(): any {
        return this._data?.rawText;
    }
    get readonlyToken(): any {
        return this._data?.readonlyToken;
    }
    get right(): any {
        return this._data?.right;
    }
    get statement(): any {
        return this._data?.statement;
    }
    get statements(): any {
        return this._data?.statements;
    }
    get tag(): any {
        return this._data?.tag;
    }
    get tagName(): any {
        return this._data?.tagName;
    }
    get tags(): any {
        return this._data?.tags;
    }
    get template(): any {
        return this._data?.template;
    }
    get templateFlags(): any {
        return this._data?.templateFlags;
    }
    get templateSpans(): any {
        return this._data?.templateSpans;
    }
    get text(): any {
        return this._data?.text;
    }
    get thenStatement(): any {
        return this._data?.thenStatement;
    }
    get token(): any {
        return this._data?.token;
    }
    get trueType(): any {
        return this._data?.trueType;
    }
    get tryBlock(): any {
        return this._data?.tryBlock;
    }
    get type(): any {
        return this._data?.type;
    }
    get typeArguments(): any {
        return this._data?.typeArguments;
    }
    get typeExpression(): any {
        return this._data?.typeExpression;
    }
    get typeName(): any {
        return this._data?.typeName;
    }
    get typeParameter(): any {
        return this._data?.typeParameter;
    }
    get typeParameters(): any {
        return this._data?.typeParameters;
    }
    get types(): any {
        return this._data?.types;
    }
    get value(): any {
        return this._data?.value;
    }
    get variableDeclaration(): any {
        return this._data?.variableDeclaration;
    }
    get whenFalse(): any {
        return this._data?.whenFalse;
    }
    get whenTrue(): any {
        return this._data?.whenTrue;
    }
}

/**
 * Create a simple token node with only a `kind`.
 */
export function createToken<TKind extends SyntaxKind>(kind: TKind): Node & { readonly kind: TKind; } {
    return new NodeObject(kind, undefined) as any;
}

export function createArrayBindingPattern(elements: NodeArray<ArrayBindingElement>): ArrayBindingPattern {
    return new NodeObject(SyntaxKind.ArrayBindingPattern, {
        elements,
    }) as unknown as ArrayBindingPattern;
}

export function createArrayLiteralExpression(elements: NodeArray<Expression>, multiLine?: boolean): ArrayLiteralExpression {
    return new NodeObject(SyntaxKind.ArrayLiteralExpression, {
        elements,
        multiLine,
    }) as unknown as ArrayLiteralExpression;
}

export function createArrayTypeNode(elementType: TypeNode): ArrayTypeNode {
    return new NodeObject(SyntaxKind.ArrayType, {
        elementType,
    }) as unknown as ArrayTypeNode;
}

export function createArrowFunction(name: never, parameters: NodeArray<ParameterDeclaration>, body: ConciseBody, equalsGreaterThanToken: EqualsGreaterThanToken, typeParameters?: NodeArray<TypeParameterDeclaration> | undefined, type?: TypeNode | undefined, asteriskToken?: AsteriskToken | undefined, questionToken?: QuestionToken | undefined, exclamationToken?: ExclamationToken | undefined, modifiers?: NodeArray<Modifier>): ArrowFunction {
    return new NodeObject(SyntaxKind.ArrowFunction, {
        name,
        parameters,
        body,
        equalsGreaterThanToken,
        typeParameters,
        type,
        asteriskToken,
        questionToken,
        exclamationToken,
        modifiers,
    }) as unknown as ArrowFunction;
}

export function createAsExpression(expression: Expression, type: TypeNode): AsExpression {
    return new NodeObject(SyntaxKind.AsExpression, {
        expression,
        type,
    }) as unknown as AsExpression;
}

export function createAwaitExpression(expression: UnaryExpression): AwaitExpression {
    return new NodeObject(SyntaxKind.AwaitExpression, {
        expression,
    }) as unknown as AwaitExpression;
}

export function createBigIntLiteral(text: string, isUnterminated?: boolean, hasExtendedUnicodeEscape?: boolean): BigIntLiteral {
    return new NodeObject(SyntaxKind.BigIntLiteral, {
        text,
        isUnterminated,
        hasExtendedUnicodeEscape,
    }) as unknown as BigIntLiteral;
}

export function createBinaryExpression(left: Expression, operatorToken: BinaryOperatorToken, right: Expression): BinaryExpression {
    return new NodeObject(SyntaxKind.BinaryExpression, {
        left,
        operatorToken,
        right,
    }) as unknown as BinaryExpression;
}

export function createBindingElement(name: BindingName, propertyName?: PropertyName, dotDotDotToken?: DotDotDotToken, initializer?: Expression): BindingElement {
    return new NodeObject(SyntaxKind.BindingElement, {
        name,
        propertyName,
        dotDotDotToken,
        initializer,
    }) as unknown as BindingElement;
}

export function createBlock(statements: NodeArray<Statement>, multiLine?: boolean): Block {
    return new NodeObject(SyntaxKind.Block, {
        statements,
        multiLine,
    }) as unknown as Block;
}

export function createBreakStatement(label?: Identifier): BreakStatement {
    return new NodeObject(SyntaxKind.BreakStatement, {
        label,
    }) as unknown as BreakStatement;
}

export function createCallExpression(expression: LeftHandSideExpression, arguments_: NodeArray<Expression>, questionDotToken?: QuestionDotToken, typeArguments?: NodeArray<TypeNode>): CallExpression {
    return new NodeObject(SyntaxKind.CallExpression, {
        expression,
        arguments: arguments_,
        questionDotToken,
        typeArguments,
    }) as unknown as CallExpression;
}

export function createCallSignatureDeclaration(parameters: NodeArray<ParameterDeclaration>, name?: PropertyName, typeParameters?: NodeArray<TypeParameterDeclaration> | undefined, type?: TypeNode | undefined, questionToken?: QuestionToken | undefined): CallSignatureDeclaration {
    return new NodeObject(SyntaxKind.CallSignature, {
        parameters,
        name,
        typeParameters,
        type,
        questionToken,
    }) as unknown as CallSignatureDeclaration;
}

export function createCaseBlock(clauses: NodeArray<CaseOrDefaultClause>): CaseBlock {
    return new NodeObject(SyntaxKind.CaseBlock, {
        clauses,
    }) as unknown as CaseBlock;
}

export function createCaseClause(expression: Expression, statements: NodeArray<Statement>): CaseClause {
    return new NodeObject(SyntaxKind.CaseClause, {
        expression,
        statements,
    }) as unknown as CaseClause;
}

export function createCatchClause(block: Block, variableDeclaration?: VariableDeclaration): CatchClause {
    return new NodeObject(SyntaxKind.CatchClause, {
        block,
        variableDeclaration,
    }) as unknown as CatchClause;
}

export function createClassDeclaration(members: NodeArray<ClassElement>, name?: Identifier, typeParameters?: NodeArray<TypeParameterDeclaration>, heritageClauses?: NodeArray<HeritageClause>, modifiers?: NodeArray<ModifierLike>): ClassDeclaration {
    return new NodeObject(SyntaxKind.ClassDeclaration, {
        members,
        name,
        typeParameters,
        heritageClauses,
        modifiers,
    }) as unknown as ClassDeclaration;
}

export function createClassExpression(members: NodeArray<ClassElement>, name?: Identifier, typeParameters?: NodeArray<TypeParameterDeclaration>, heritageClauses?: NodeArray<HeritageClause>, modifiers?: NodeArray<ModifierLike>): ClassExpression {
    return new NodeObject(SyntaxKind.ClassExpression, {
        members,
        name,
        typeParameters,
        heritageClauses,
        modifiers,
    }) as unknown as ClassExpression;
}

export function createClassStaticBlockDeclaration(body: Block, name?: PropertyName): ClassStaticBlockDeclaration {
    return new NodeObject(SyntaxKind.ClassStaticBlockDeclaration, {
        body,
        name,
    }) as unknown as ClassStaticBlockDeclaration;
}

export function createCommaListExpression(elements: NodeArray<Expression>): CommaListExpression {
    return new NodeObject(SyntaxKind.CommaListExpression, {
        elements,
    }) as unknown as CommaListExpression;
}

export function createComputedPropertyName(expression: Expression): ComputedPropertyName {
    return new NodeObject(SyntaxKind.ComputedPropertyName, {
        expression,
    }) as unknown as ComputedPropertyName;
}

export function createConditionalExpression(condition: Expression, questionToken: QuestionToken, whenTrue: Expression, colonToken: ColonToken, whenFalse: Expression): ConditionalExpression {
    return new NodeObject(SyntaxKind.ConditionalExpression, {
        condition,
        questionToken,
        whenTrue,
        colonToken,
        whenFalse,
    }) as unknown as ConditionalExpression;
}

export function createConditionalTypeNode(checkType: TypeNode, extendsType: TypeNode, trueType: TypeNode, falseType: TypeNode): ConditionalTypeNode {
    return new NodeObject(SyntaxKind.ConditionalType, {
        checkType,
        extendsType,
        trueType,
        falseType,
    }) as unknown as ConditionalTypeNode;
}

export function createConstructorDeclaration(parameters: NodeArray<ParameterDeclaration>, name?: PropertyName, typeParameters?: NodeArray<TypeParameterDeclaration> | undefined, type?: TypeNode | undefined, asteriskToken?: AsteriskToken | undefined, questionToken?: QuestionToken | undefined, exclamationToken?: ExclamationToken | undefined, body?: FunctionBody | undefined, modifiers?: NodeArray<ModifierLike> | undefined): ConstructorDeclaration {
    return new NodeObject(SyntaxKind.Constructor, {
        parameters,
        name,
        typeParameters,
        type,
        asteriskToken,
        questionToken,
        exclamationToken,
        body,
        modifiers,
    }) as unknown as ConstructorDeclaration;
}

export function createConstructorTypeNode(parameters: NodeArray<ParameterDeclaration>, type: TypeNode, name?: PropertyName, typeParameters?: NodeArray<TypeParameterDeclaration> | undefined, modifiers?: NodeArray<Modifier>): ConstructorTypeNode {
    return new NodeObject(SyntaxKind.ConstructorType, {
        parameters,
        type,
        name,
        typeParameters,
        modifiers,
    }) as unknown as ConstructorTypeNode;
}

export function createConstructSignatureDeclaration(parameters: NodeArray<ParameterDeclaration>, name?: PropertyName, typeParameters?: NodeArray<TypeParameterDeclaration> | undefined, type?: TypeNode | undefined, questionToken?: QuestionToken | undefined): ConstructSignatureDeclaration {
    return new NodeObject(SyntaxKind.ConstructSignature, {
        parameters,
        name,
        typeParameters,
        type,
        questionToken,
    }) as unknown as ConstructSignatureDeclaration;
}

export function createContinueStatement(label?: Identifier): ContinueStatement {
    return new NodeObject(SyntaxKind.ContinueStatement, {
        label,
    }) as unknown as ContinueStatement;
}

export function createDebuggerStatement(): DebuggerStatement {
    return new NodeObject(SyntaxKind.DebuggerStatement, undefined) as unknown as DebuggerStatement;
}

export function createDecorator(expression: LeftHandSideExpression): Decorator {
    return new NodeObject(SyntaxKind.Decorator, {
        expression,
    }) as unknown as Decorator;
}

export function createDefaultClause(statements: NodeArray<Statement>): DefaultClause {
    return new NodeObject(SyntaxKind.DefaultClause, {
        statements,
    }) as unknown as DefaultClause;
}

export function createDeleteExpression(expression: UnaryExpression): DeleteExpression {
    return new NodeObject(SyntaxKind.DeleteExpression, {
        expression,
    }) as unknown as DeleteExpression;
}

export function createDoStatement(statement: Statement, expression: Expression): DoStatement {
    return new NodeObject(SyntaxKind.DoStatement, {
        statement,
        expression,
    }) as unknown as DoStatement;
}

export function createElementAccessExpression(expression: LeftHandSideExpression, argumentExpression: Expression, questionDotToken?: QuestionDotToken): ElementAccessExpression {
    return new NodeObject(SyntaxKind.ElementAccessExpression, {
        expression,
        argumentExpression,
        questionDotToken,
    }) as unknown as ElementAccessExpression;
}

export function createEmptyStatement(): EmptyStatement {
    return new NodeObject(SyntaxKind.EmptyStatement, undefined) as unknown as EmptyStatement;
}

export function createEnumDeclaration(name: Identifier, members: NodeArray<EnumMember>, modifiers?: NodeArray<ModifierLike>): EnumDeclaration {
    return new NodeObject(SyntaxKind.EnumDeclaration, {
        name,
        members,
        modifiers,
    }) as unknown as EnumDeclaration;
}

export function createEnumMember(name: PropertyName, initializer?: Expression): EnumMember {
    return new NodeObject(SyntaxKind.EnumMember, {
        name,
        initializer,
    }) as unknown as EnumMember;
}

export function createExportAssignment(expression: Expression, name?: Identifier | StringLiteral | NumericLiteral, modifiers?: NodeArray<ModifierLike>, isExportEquals?: boolean): ExportAssignment {
    return new NodeObject(SyntaxKind.ExportAssignment, {
        expression,
        name,
        modifiers,
        isExportEquals,
    }) as unknown as ExportAssignment;
}

export function createExportDeclaration(isTypeOnly: boolean, name?: Identifier | StringLiteral | NumericLiteral, modifiers?: NodeArray<ModifierLike>, exportClause?: NamedExportBindings, moduleSpecifier?: Expression, assertClause?: AssertClause, attributes?: ImportAttributes): ExportDeclaration {
    return new NodeObject(SyntaxKind.ExportDeclaration, {
        isTypeOnly,
        name,
        modifiers,
        exportClause,
        moduleSpecifier,
        assertClause,
        attributes,
    }) as unknown as ExportDeclaration;
}

export function createExportSpecifier(name: ModuleExportName, isTypeOnly: boolean, propertyName?: ModuleExportName): ExportSpecifier {
    return new NodeObject(SyntaxKind.ExportSpecifier, {
        name,
        isTypeOnly,
        propertyName,
    }) as unknown as ExportSpecifier;
}

export function createExpressionStatement(expression: Expression): ExpressionStatement {
    return new NodeObject(SyntaxKind.ExpressionStatement, {
        expression,
    }) as unknown as ExpressionStatement;
}

export function createExpressionWithTypeArguments(expression: LeftHandSideExpression, typeArguments?: NodeArray<TypeNode>): ExpressionWithTypeArguments {
    return new NodeObject(SyntaxKind.ExpressionWithTypeArguments, {
        expression,
        typeArguments,
    }) as unknown as ExpressionWithTypeArguments;
}

export function createExternalModuleReference(expression: Expression): ExternalModuleReference {
    return new NodeObject(SyntaxKind.ExternalModuleReference, {
        expression,
    }) as unknown as ExternalModuleReference;
}

export function createFalseLiteral(): FalseLiteral {
    return new NodeObject(SyntaxKind.FalseKeyword, undefined) as unknown as FalseLiteral;
}

export function createForInStatement(statement: Statement, initializer: ForInitializer, expression: Expression): ForInStatement {
    return new NodeObject(SyntaxKind.ForInStatement, {
        statement,
        initializer,
        expression,
    }) as unknown as ForInStatement;
}

export function createForOfStatement(statement: Statement, initializer: ForInitializer, expression: Expression, awaitModifier?: AwaitKeyword): ForOfStatement {
    return new NodeObject(SyntaxKind.ForOfStatement, {
        statement,
        initializer,
        expression,
        awaitModifier,
    }) as unknown as ForOfStatement;
}

export function createForStatement(statement: Statement, initializer?: ForInitializer, condition?: Expression, incrementor?: Expression): ForStatement {
    return new NodeObject(SyntaxKind.ForStatement, {
        statement,
        initializer,
        condition,
        incrementor,
    }) as unknown as ForStatement;
}

export function createFunctionDeclaration(parameters: NodeArray<ParameterDeclaration>, name?: Identifier, typeParameters?: NodeArray<TypeParameterDeclaration> | undefined, type?: TypeNode | undefined, asteriskToken?: AsteriskToken | undefined, questionToken?: QuestionToken | undefined, exclamationToken?: ExclamationToken | undefined, body?: FunctionBody, modifiers?: NodeArray<ModifierLike>): FunctionDeclaration {
    return new NodeObject(SyntaxKind.FunctionDeclaration, {
        parameters,
        name,
        typeParameters,
        type,
        asteriskToken,
        questionToken,
        exclamationToken,
        body,
        modifiers,
    }) as unknown as FunctionDeclaration;
}

export function createFunctionExpression(parameters: NodeArray<ParameterDeclaration>, body: FunctionBody, name?: Identifier, typeParameters?: NodeArray<TypeParameterDeclaration> | undefined, type?: TypeNode | undefined, asteriskToken?: AsteriskToken | undefined, questionToken?: QuestionToken | undefined, exclamationToken?: ExclamationToken | undefined, modifiers?: NodeArray<Modifier>): FunctionExpression {
    return new NodeObject(SyntaxKind.FunctionExpression, {
        parameters,
        body,
        name,
        typeParameters,
        type,
        asteriskToken,
        questionToken,
        exclamationToken,
        modifiers,
    }) as unknown as FunctionExpression;
}

export function createFunctionTypeNode(parameters: NodeArray<ParameterDeclaration>, type: TypeNode, name?: PropertyName, typeParameters?: NodeArray<TypeParameterDeclaration> | undefined): FunctionTypeNode {
    return new NodeObject(SyntaxKind.FunctionType, {
        parameters,
        type,
        name,
        typeParameters,
    }) as unknown as FunctionTypeNode;
}

export function createGetAccessorDeclaration(name: PropertyName, parameters: NodeArray<ParameterDeclaration>, typeParameters?: NodeArray<TypeParameterDeclaration> | undefined, type?: TypeNode | undefined, asteriskToken?: AsteriskToken | undefined, questionToken?: QuestionToken | undefined, exclamationToken?: ExclamationToken | undefined, body?: FunctionBody, modifiers?: NodeArray<ModifierLike>): GetAccessorDeclaration {
    return new NodeObject(SyntaxKind.GetAccessor, {
        name,
        parameters,
        typeParameters,
        type,
        asteriskToken,
        questionToken,
        exclamationToken,
        body,
        modifiers,
    }) as unknown as GetAccessorDeclaration;
}

export function createHeritageClause(token: SyntaxKind.ExtendsKeyword | SyntaxKind.ImplementsKeyword, types: NodeArray<ExpressionWithTypeArguments>): HeritageClause {
    return new NodeObject(SyntaxKind.HeritageClause, {
        token,
        types,
    }) as unknown as HeritageClause;
}

export function createIdentifier(text: string): Identifier {
    return new NodeObject(SyntaxKind.Identifier, {
        text,
    }) as unknown as Identifier;
}

export function createIfStatement(expression: Expression, thenStatement: Statement, elseStatement?: Statement): IfStatement {
    return new NodeObject(SyntaxKind.IfStatement, {
        expression,
        thenStatement,
        elseStatement,
    }) as unknown as IfStatement;
}

export function createImportAttribute(name: ImportAttributeName, value: Expression): ImportAttribute {
    return new NodeObject(SyntaxKind.ImportAttribute, {
        name,
        value,
    }) as unknown as ImportAttribute;
}

export function createImportAttributes(token: SyntaxKind.WithKeyword | SyntaxKind.AssertKeyword, elements: NodeArray<ImportAttribute>, multiLine?: boolean): ImportAttributes {
    return new NodeObject(SyntaxKind.ImportAttributes, {
        token,
        elements,
        multiLine,
    }) as unknown as ImportAttributes;
}

export function createImportClause(phaseModifier: ImportPhaseModifierSyntaxKind, name?: Identifier, namedBindings?: NamedImportBindings): ImportClause {
    return new NodeObject(SyntaxKind.ImportClause, {
        phaseModifier,
        name,
        namedBindings,
    }) as unknown as ImportClause;
}

export function createImportDeclaration(moduleSpecifier: Expression, modifiers?: NodeArray<ModifierLike>, importClause?: ImportClause, assertClause?: AssertClause, attributes?: ImportAttributes): ImportDeclaration {
    return new NodeObject(SyntaxKind.ImportDeclaration, {
        moduleSpecifier,
        modifiers,
        importClause,
        assertClause,
        attributes,
    }) as unknown as ImportDeclaration;
}

export function createImportEqualsDeclaration(name: Identifier, isTypeOnly: boolean, moduleReference: ModuleReference, modifiers?: NodeArray<ModifierLike>): ImportEqualsDeclaration {
    return new NodeObject(SyntaxKind.ImportEqualsDeclaration, {
        name,
        isTypeOnly,
        moduleReference,
        modifiers,
    }) as unknown as ImportEqualsDeclaration;
}

export function createImportExpression(): ImportExpression {
    return new NodeObject(SyntaxKind.ImportKeyword, undefined) as unknown as ImportExpression;
}

export function createImportSpecifier(name: Identifier, isTypeOnly: boolean, propertyName?: ModuleExportName): ImportSpecifier {
    return new NodeObject(SyntaxKind.ImportSpecifier, {
        name,
        isTypeOnly,
        propertyName,
    }) as unknown as ImportSpecifier;
}

export function createImportTypeNode(isTypeOf: boolean, argument: TypeNode, typeArguments?: NodeArray<TypeNode>, attributes?: ImportAttributes, qualifier?: EntityName): ImportTypeNode {
    return new NodeObject(SyntaxKind.ImportType, {
        isTypeOf,
        argument,
        typeArguments,
        attributes,
        qualifier,
    }) as unknown as ImportTypeNode;
}

export function createIndexedAccessTypeNode(objectType: TypeNode, indexType: TypeNode): IndexedAccessTypeNode {
    return new NodeObject(SyntaxKind.IndexedAccessType, {
        objectType,
        indexType,
    }) as unknown as IndexedAccessTypeNode;
}

export function createIndexSignatureDeclaration(parameters: NodeArray<ParameterDeclaration>, type: TypeNode, name?: PropertyName, typeParameters?: NodeArray<TypeParameterDeclaration> | undefined, questionToken?: QuestionToken | undefined, modifiers?: NodeArray<ModifierLike>): IndexSignatureDeclaration {
    return new NodeObject(SyntaxKind.IndexSignature, {
        parameters,
        type,
        name,
        typeParameters,
        questionToken,
        modifiers,
    }) as unknown as IndexSignatureDeclaration;
}

export function createInferTypeNode(typeParameter: TypeParameterDeclaration): InferTypeNode {
    return new NodeObject(SyntaxKind.InferType, {
        typeParameter,
    }) as unknown as InferTypeNode;
}

export function createInterfaceDeclaration(name: Identifier, members: NodeArray<TypeElement>, modifiers?: NodeArray<ModifierLike>, typeParameters?: NodeArray<TypeParameterDeclaration>, heritageClauses?: NodeArray<HeritageClause>): InterfaceDeclaration {
    return new NodeObject(SyntaxKind.InterfaceDeclaration, {
        name,
        members,
        modifiers,
        typeParameters,
        heritageClauses,
    }) as unknown as InterfaceDeclaration;
}

export function createIntersectionTypeNode(types: NodeArray<TypeNode>): IntersectionTypeNode {
    return new NodeObject(SyntaxKind.IntersectionType, {
        types,
    }) as unknown as IntersectionTypeNode;
}

export function createJSDoc(tags?: NodeArray<JSDocTag>, comment?: string | NodeArray<JSDocComment>): JSDoc {
    return new NodeObject(SyntaxKind.JSDoc, {
        tags,
        comment,
    }) as unknown as JSDoc;
}

export function createJSDocAllType(): JSDocAllType {
    return new NodeObject(SyntaxKind.JSDocAllType, undefined) as unknown as JSDocAllType;
}

export function createJSDocAugmentsTag(tagName: Identifier, class_: ExpressionWithTypeArguments & { readonly expression: Identifier | PropertyAccessEntityNameExpression; }, comment?: string | NodeArray<JSDocComment>): JSDocAugmentsTag {
    return new NodeObject(SyntaxKind.JSDocAugmentsTag, {
        tagName,
        class: class_,
        comment,
    }) as unknown as JSDocAugmentsTag;
}

export function createJSDocCallbackTag(tagName: Identifier, typeExpression: JSDocSignature, comment?: string | NodeArray<JSDocComment>, name?: Identifier, fullName?: JSDocNamespaceDeclaration | Identifier): JSDocCallbackTag {
    return new NodeObject(SyntaxKind.JSDocCallbackTag, {
        tagName,
        typeExpression,
        comment,
        name,
        fullName,
    }) as unknown as JSDocCallbackTag;
}

export function createJSDocDeprecatedTag(tagName: Identifier, comment?: string | NodeArray<JSDocComment>): JSDocDeprecatedTag {
    return new NodeObject(SyntaxKind.JSDocDeprecatedTag, {
        tagName,
        comment,
    }) as unknown as JSDocDeprecatedTag;
}

export function createJSDocImplementsTag(tagName: Identifier, class_: ExpressionWithTypeArguments & { readonly expression: Identifier | PropertyAccessEntityNameExpression; }, comment?: string | NodeArray<JSDocComment>): JSDocImplementsTag {
    return new NodeObject(SyntaxKind.JSDocImplementsTag, {
        tagName,
        class: class_,
        comment,
    }) as unknown as JSDocImplementsTag;
}

export function createJSDocImportTag(tagName: Identifier, moduleSpecifier: Expression, comment?: string | NodeArray<JSDocComment>, importClause?: ImportClause, attributes?: ImportAttributes): JSDocImportTag {
    return new NodeObject(SyntaxKind.JSDocImportTag, {
        tagName,
        moduleSpecifier,
        comment,
        importClause,
        attributes,
    }) as unknown as JSDocImportTag;
}

export function createJSDocLink(text: string, name?: EntityName | JSDocMemberName): JSDocLink {
    return new NodeObject(SyntaxKind.JSDocLink, {
        text,
        name,
    }) as unknown as JSDocLink;
}

export function createJSDocLinkCode(text: string, name?: EntityName | JSDocMemberName): JSDocLinkCode {
    return new NodeObject(SyntaxKind.JSDocLinkCode, {
        text,
        name,
    }) as unknown as JSDocLinkCode;
}

export function createJSDocLinkPlain(text: string, name?: EntityName | JSDocMemberName): JSDocLinkPlain {
    return new NodeObject(SyntaxKind.JSDocLinkPlain, {
        text,
        name,
    }) as unknown as JSDocLinkPlain;
}

export function createJSDocMemberName(left: EntityName | JSDocMemberName, right: Identifier): JSDocMemberName {
    return new NodeObject(SyntaxKind.JSDocMemberName, {
        left,
        right,
    }) as unknown as JSDocMemberName;
}

export function createJSDocNameReference(name: EntityName | JSDocMemberName): JSDocNameReference {
    return new NodeObject(SyntaxKind.JSDocNameReference, {
        name,
    }) as unknown as JSDocNameReference;
}

export function createJSDocNonNullableType(type: TypeNode, postfix: boolean): JSDocNonNullableType {
    return new NodeObject(SyntaxKind.JSDocNonNullableType, {
        type,
        postfix,
    }) as unknown as JSDocNonNullableType;
}

export function createJSDocNullableType(type: TypeNode, postfix: boolean): JSDocNullableType {
    return new NodeObject(SyntaxKind.JSDocNullableType, {
        type,
        postfix,
    }) as unknown as JSDocNullableType;
}

export function createJSDocOptionalType(type: TypeNode): JSDocOptionalType {
    return new NodeObject(SyntaxKind.JSDocOptionalType, {
        type,
    }) as unknown as JSDocOptionalType;
}

export function createJSDocOverloadTag(tagName: Identifier, typeExpression: JSDocSignature, comment?: string | NodeArray<JSDocComment>): JSDocOverloadTag {
    return new NodeObject(SyntaxKind.JSDocOverloadTag, {
        tagName,
        typeExpression,
        comment,
    }) as unknown as JSDocOverloadTag;
}

export function createJSDocOverrideTag(tagName: Identifier, comment?: string | NodeArray<JSDocComment>): JSDocOverrideTag {
    return new NodeObject(SyntaxKind.JSDocOverrideTag, {
        tagName,
        comment,
    }) as unknown as JSDocOverrideTag;
}

export function createJSDocParameterTag(tagName: Identifier, name: EntityName, isNameFirst: boolean, isBracketed: boolean, comment?: string | NodeArray<JSDocComment>, typeExpression?: JSDocTypeExpression): JSDocParameterTag {
    return new NodeObject(SyntaxKind.JSDocParameterTag, {
        tagName,
        name,
        isNameFirst,
        isBracketed,
        comment,
        typeExpression,
    }) as unknown as JSDocParameterTag;
}

export function createJSDocPrivateTag(tagName: Identifier, comment?: string | NodeArray<JSDocComment>): JSDocPrivateTag {
    return new NodeObject(SyntaxKind.JSDocPrivateTag, {
        tagName,
        comment,
    }) as unknown as JSDocPrivateTag;
}

export function createJSDocPropertyTag(tagName: Identifier, name: EntityName, isNameFirst: boolean, isBracketed: boolean, comment?: string | NodeArray<JSDocComment>, typeExpression?: JSDocTypeExpression): JSDocPropertyTag {
    return new NodeObject(SyntaxKind.JSDocPropertyTag, {
        tagName,
        name,
        isNameFirst,
        isBracketed,
        comment,
        typeExpression,
    }) as unknown as JSDocPropertyTag;
}

export function createJSDocProtectedTag(tagName: Identifier, comment?: string | NodeArray<JSDocComment>): JSDocProtectedTag {
    return new NodeObject(SyntaxKind.JSDocProtectedTag, {
        tagName,
        comment,
    }) as unknown as JSDocProtectedTag;
}

export function createJSDocPublicTag(tagName: Identifier, comment?: string | NodeArray<JSDocComment>): JSDocPublicTag {
    return new NodeObject(SyntaxKind.JSDocPublicTag, {
        tagName,
        comment,
    }) as unknown as JSDocPublicTag;
}

export function createJSDocReadonlyTag(tagName: Identifier, comment?: string | NodeArray<JSDocComment>): JSDocReadonlyTag {
    return new NodeObject(SyntaxKind.JSDocReadonlyTag, {
        tagName,
        comment,
    }) as unknown as JSDocReadonlyTag;
}

export function createJSDocReturnTag(tagName: Identifier, comment?: string | NodeArray<JSDocComment>, typeExpression?: JSDocTypeExpression): JSDocReturnTag {
    return new NodeObject(SyntaxKind.JSDocReturnTag, {
        tagName,
        comment,
        typeExpression,
    }) as unknown as JSDocReturnTag;
}

export function createJSDocSatisfiesTag(tagName: Identifier, typeExpression: JSDocTypeExpression, comment?: string | NodeArray<JSDocComment>): JSDocSatisfiesTag {
    return new NodeObject(SyntaxKind.JSDocSatisfiesTag, {
        tagName,
        typeExpression,
        comment,
    }) as unknown as JSDocSatisfiesTag;
}

export function createJSDocSeeTag(tagName: Identifier, comment?: string | NodeArray<JSDocComment>, name?: JSDocNameReference): JSDocSeeTag {
    return new NodeObject(SyntaxKind.JSDocSeeTag, {
        tagName,
        comment,
        name,
    }) as unknown as JSDocSeeTag;
}

export function createJSDocSignature(parameters: readonly JSDocParameterTag[], type: JSDocReturnTag | undefined, typeParameters?: readonly JSDocTemplateTag[]): JSDocSignature {
    return new NodeObject(SyntaxKind.JSDocSignature, {
        parameters,
        type,
        typeParameters,
    }) as unknown as JSDocSignature;
}

export function createJSDocTemplateTag(tagName: Identifier, constraint: JSDocTypeExpression | undefined, typeParameters: NodeArray<TypeParameterDeclaration>, comment?: string | NodeArray<JSDocComment>): JSDocTemplateTag {
    return new NodeObject(SyntaxKind.JSDocTemplateTag, {
        tagName,
        constraint,
        typeParameters,
        comment,
    }) as unknown as JSDocTemplateTag;
}

export function createJSDocText(text: string): JSDocText {
    return new NodeObject(SyntaxKind.JSDocText, {
        text,
    }) as unknown as JSDocText;
}

export function createJSDocThisTag(tagName: Identifier, typeExpression: JSDocTypeExpression, comment?: string | NodeArray<JSDocComment>): JSDocThisTag {
    return new NodeObject(SyntaxKind.JSDocThisTag, {
        tagName,
        typeExpression,
        comment,
    }) as unknown as JSDocThisTag;
}

export function createJSDocTypedefTag(tagName: Identifier, comment?: string | NodeArray<JSDocComment>, name?: Identifier, fullName?: JSDocNamespaceDeclaration | Identifier, typeExpression?: JSDocTypeExpression | JSDocTypeLiteral): JSDocTypedefTag {
    return new NodeObject(SyntaxKind.JSDocTypedefTag, {
        tagName,
        comment,
        name,
        fullName,
        typeExpression,
    }) as unknown as JSDocTypedefTag;
}

export function createJSDocTypeExpression(type: TypeNode): JSDocTypeExpression {
    return new NodeObject(SyntaxKind.JSDocTypeExpression, {
        type,
    }) as unknown as JSDocTypeExpression;
}

export function createJSDocTypeLiteral(isArrayType: boolean, jsDocPropertyTags?: readonly JSDocPropertyLikeTag[]): JSDocTypeLiteral {
    return new NodeObject(SyntaxKind.JSDocTypeLiteral, {
        isArrayType,
        jsDocPropertyTags,
    }) as unknown as JSDocTypeLiteral;
}

export function createJSDocTypeTag(tagName: Identifier, typeExpression: JSDocTypeExpression, comment?: string | NodeArray<JSDocComment>): JSDocTypeTag {
    return new NodeObject(SyntaxKind.JSDocTypeTag, {
        tagName,
        typeExpression,
        comment,
    }) as unknown as JSDocTypeTag;
}

export function createJSDocUnknownTag(tagName: Identifier, comment?: string | NodeArray<JSDocComment>): JSDocUnknownTag {
    return new NodeObject(SyntaxKind.JSDocTag, {
        tagName,
        comment,
    }) as unknown as JSDocUnknownTag;
}

export function createJSDocVariadicType(type: TypeNode): JSDocVariadicType {
    return new NodeObject(SyntaxKind.JSDocVariadicType, {
        type,
    }) as unknown as JSDocVariadicType;
}

export function createJsxAttribute(name: JsxAttributeName, initializer?: JsxAttributeValue): JsxAttribute {
    return new NodeObject(SyntaxKind.JsxAttribute, {
        name,
        initializer,
    }) as unknown as JsxAttribute;
}

export function createJsxAttributes(properties: NodeArray<JsxAttributeLike>): JsxAttributes {
    return new NodeObject(SyntaxKind.JsxAttributes, {
        properties,
    }) as unknown as JsxAttributes;
}

export function createJsxClosingElement(tagName: JsxTagNameExpression): JsxClosingElement {
    return new NodeObject(SyntaxKind.JsxClosingElement, {
        tagName,
    }) as unknown as JsxClosingElement;
}

export function createJsxClosingFragment(): JsxClosingFragment {
    return new NodeObject(SyntaxKind.JsxClosingFragment, undefined) as unknown as JsxClosingFragment;
}

export function createJsxElement(openingElement: JsxOpeningElement, children: NodeArray<JsxChild>, closingElement: JsxClosingElement): JsxElement {
    return new NodeObject(SyntaxKind.JsxElement, {
        openingElement,
        children,
        closingElement,
    }) as unknown as JsxElement;
}

export function createJsxExpression(dotDotDotToken?: Token<SyntaxKind.DotDotDotToken>, expression?: Expression): JsxExpression {
    return new NodeObject(SyntaxKind.JsxExpression, {
        dotDotDotToken,
        expression,
    }) as unknown as JsxExpression;
}

export function createJsxFragment(openingFragment: JsxOpeningFragment, children: NodeArray<JsxChild>, closingFragment: JsxClosingFragment): JsxFragment {
    return new NodeObject(SyntaxKind.JsxFragment, {
        openingFragment,
        children,
        closingFragment,
    }) as unknown as JsxFragment;
}

export function createJsxNamespacedName(name: Identifier, namespace: Identifier): JsxNamespacedName {
    return new NodeObject(SyntaxKind.JsxNamespacedName, {
        name,
        namespace,
    }) as unknown as JsxNamespacedName;
}

export function createJsxOpeningElement(tagName: JsxTagNameExpression, attributes: JsxAttributes, typeArguments?: NodeArray<TypeNode>): JsxOpeningElement {
    return new NodeObject(SyntaxKind.JsxOpeningElement, {
        tagName,
        attributes,
        typeArguments,
    }) as unknown as JsxOpeningElement;
}

export function createJsxOpeningFragment(): JsxOpeningFragment {
    return new NodeObject(SyntaxKind.JsxOpeningFragment, undefined) as unknown as JsxOpeningFragment;
}

export function createJsxSelfClosingElement(tagName: JsxTagNameExpression, attributes: JsxAttributes, typeArguments?: NodeArray<TypeNode>): JsxSelfClosingElement {
    return new NodeObject(SyntaxKind.JsxSelfClosingElement, {
        tagName,
        attributes,
        typeArguments,
    }) as unknown as JsxSelfClosingElement;
}

export function createJsxSpreadAttribute(expression: Expression, name?: PropertyName): JsxSpreadAttribute {
    return new NodeObject(SyntaxKind.JsxSpreadAttribute, {
        expression,
        name,
    }) as unknown as JsxSpreadAttribute;
}

export function createJsxText(text: string, containsOnlyTriviaWhiteSpaces: boolean, isUnterminated?: boolean, hasExtendedUnicodeEscape?: boolean): JsxText {
    return new NodeObject(SyntaxKind.JsxText, {
        text,
        containsOnlyTriviaWhiteSpaces,
        isUnterminated,
        hasExtendedUnicodeEscape,
    }) as unknown as JsxText;
}

export function createLabeledStatement(label: Identifier, statement: Statement): LabeledStatement {
    return new NodeObject(SyntaxKind.LabeledStatement, {
        label,
        statement,
    }) as unknown as LabeledStatement;
}

export function createLiteralTypeNode(literal: NullLiteral | BooleanLiteral | LiteralExpression | PrefixUnaryExpression): LiteralTypeNode {
    return new NodeObject(SyntaxKind.LiteralType, {
        literal,
    }) as unknown as LiteralTypeNode;
}

export function createMappedTypeNode(typeParameter: TypeParameterDeclaration, readonlyToken?: ReadonlyKeyword | PlusToken | MinusToken, nameType?: TypeNode, questionToken?: QuestionToken | PlusToken | MinusToken, type?: TypeNode, members?: NodeArray<TypeElement>): MappedTypeNode {
    return new NodeObject(SyntaxKind.MappedType, {
        typeParameter,
        readonlyToken,
        nameType,
        questionToken,
        type,
        members,
    }) as unknown as MappedTypeNode;
}

export function createMetaProperty(keywordToken: SyntaxKind.NewKeyword | SyntaxKind.ImportKeyword, name: Identifier): MetaProperty {
    return new NodeObject(SyntaxKind.MetaProperty, {
        keywordToken,
        name,
    }) as unknown as MetaProperty;
}

export function createMethodDeclaration(name: PropertyName, parameters: NodeArray<ParameterDeclaration>, typeParameters?: NodeArray<TypeParameterDeclaration> | undefined, type?: TypeNode | undefined, asteriskToken?: AsteriskToken | undefined, questionToken?: QuestionToken | undefined, exclamationToken?: ExclamationToken | undefined, body?: FunctionBody | undefined, modifiers?: NodeArray<ModifierLike> | undefined): MethodDeclaration {
    return new NodeObject(SyntaxKind.MethodDeclaration, {
        name,
        parameters,
        typeParameters,
        type,
        asteriskToken,
        questionToken,
        exclamationToken,
        body,
        modifiers,
    }) as unknown as MethodDeclaration;
}

export function createMethodSignature(name: PropertyName, parameters: NodeArray<ParameterDeclaration>, typeParameters?: NodeArray<TypeParameterDeclaration> | undefined, type?: TypeNode | undefined, questionToken?: QuestionToken | undefined, modifiers?: NodeArray<Modifier>): MethodSignature {
    return new NodeObject(SyntaxKind.MethodSignature, {
        name,
        parameters,
        typeParameters,
        type,
        questionToken,
        modifiers,
    }) as unknown as MethodSignature;
}

export function createModuleBlock(statements: NodeArray<Statement>): ModuleBlock {
    return new NodeObject(SyntaxKind.ModuleBlock, {
        statements,
    }) as unknown as ModuleBlock;
}

export function createModuleDeclaration(name: ModuleName, modifiers?: NodeArray<ModifierLike>, body?: ModuleBody | JSDocNamespaceDeclaration): ModuleDeclaration {
    return new NodeObject(SyntaxKind.ModuleDeclaration, {
        name,
        modifiers,
        body,
    }) as unknown as ModuleDeclaration;
}

export function createNamedExports(elements: NodeArray<ExportSpecifier>): NamedExports {
    return new NodeObject(SyntaxKind.NamedExports, {
        elements,
    }) as unknown as NamedExports;
}

export function createNamedImports(elements: NodeArray<ImportSpecifier>): NamedImports {
    return new NodeObject(SyntaxKind.NamedImports, {
        elements,
    }) as unknown as NamedImports;
}

export function createNamedTupleMember(name: Identifier, type: TypeNode, dotDotDotToken?: Token<SyntaxKind.DotDotDotToken>, questionToken?: Token<SyntaxKind.QuestionToken>): NamedTupleMember {
    return new NodeObject(SyntaxKind.NamedTupleMember, {
        name,
        type,
        dotDotDotToken,
        questionToken,
    }) as unknown as NamedTupleMember;
}

export function createNamespaceExport(name: ModuleExportName): NamespaceExport {
    return new NodeObject(SyntaxKind.NamespaceExport, {
        name,
    }) as unknown as NamespaceExport;
}

export function createNamespaceExportDeclaration(name: Identifier): NamespaceExportDeclaration {
    return new NodeObject(SyntaxKind.NamespaceExportDeclaration, {
        name,
    }) as unknown as NamespaceExportDeclaration;
}

export function createNamespaceImport(name: Identifier): NamespaceImport {
    return new NodeObject(SyntaxKind.NamespaceImport, {
        name,
    }) as unknown as NamespaceImport;
}

export function createNewExpression(expression: LeftHandSideExpression, typeArguments?: NodeArray<TypeNode>, arguments_?: NodeArray<Expression>): NewExpression {
    return new NodeObject(SyntaxKind.NewExpression, {
        expression,
        typeArguments,
        arguments: arguments_,
    }) as unknown as NewExpression;
}

export function createNonNullExpression(expression: Expression): NonNullExpression {
    return new NodeObject(SyntaxKind.NonNullExpression, {
        expression,
    }) as unknown as NonNullExpression;
}

export function createNoSubstitutionTemplateLiteral(text: string, isUnterminated?: boolean, hasExtendedUnicodeEscape?: boolean, rawText?: string, templateFlags?: TokenFlags): NoSubstitutionTemplateLiteral {
    return new NodeObject(SyntaxKind.NoSubstitutionTemplateLiteral, {
        text,
        isUnterminated,
        hasExtendedUnicodeEscape,
        rawText,
        templateFlags,
    }) as unknown as NoSubstitutionTemplateLiteral;
}

export function createNullLiteral(): NullLiteral {
    return new NodeObject(SyntaxKind.NullKeyword, undefined) as unknown as NullLiteral;
}

export function createNumericLiteral(text: string, numericLiteralFlags: TokenFlags, isUnterminated?: boolean, hasExtendedUnicodeEscape?: boolean): NumericLiteral {
    return new NodeObject(SyntaxKind.NumericLiteral, {
        text,
        numericLiteralFlags,
        isUnterminated,
        hasExtendedUnicodeEscape,
    }) as unknown as NumericLiteral;
}

export function createObjectBindingPattern(elements: NodeArray<BindingElement>): ObjectBindingPattern {
    return new NodeObject(SyntaxKind.ObjectBindingPattern, {
        elements,
    }) as unknown as ObjectBindingPattern;
}

export function createObjectLiteralExpression(properties: NodeArray<ObjectLiteralElementLike>, multiLine?: boolean): ObjectLiteralExpression {
    return new NodeObject(SyntaxKind.ObjectLiteralExpression, {
        properties,
        multiLine,
    }) as unknown as ObjectLiteralExpression;
}

export function createOmittedExpression(): OmittedExpression {
    return new NodeObject(SyntaxKind.OmittedExpression, undefined) as unknown as OmittedExpression;
}

export function createOptionalTypeNode(type: TypeNode): OptionalTypeNode {
    return new NodeObject(SyntaxKind.OptionalType, {
        type,
    }) as unknown as OptionalTypeNode;
}

export function createParameterDeclaration(name: BindingName, modifiers?: NodeArray<ModifierLike>, dotDotDotToken?: DotDotDotToken, questionToken?: QuestionToken, type?: TypeNode, initializer?: Expression): ParameterDeclaration {
    return new NodeObject(SyntaxKind.Parameter, {
        name,
        modifiers,
        dotDotDotToken,
        questionToken,
        type,
        initializer,
    }) as unknown as ParameterDeclaration;
}

export function createParenthesizedExpression(expression: Expression): ParenthesizedExpression {
    return new NodeObject(SyntaxKind.ParenthesizedExpression, {
        expression,
    }) as unknown as ParenthesizedExpression;
}

export function createParenthesizedTypeNode(type: TypeNode): ParenthesizedTypeNode {
    return new NodeObject(SyntaxKind.ParenthesizedType, {
        type,
    }) as unknown as ParenthesizedTypeNode;
}

export function createPartiallyEmittedExpression(expression: Expression): PartiallyEmittedExpression {
    return new NodeObject(SyntaxKind.PartiallyEmittedExpression, {
        expression,
    }) as unknown as PartiallyEmittedExpression;
}

export function createPostfixUnaryExpression(operand: LeftHandSideExpression, operator: PostfixUnaryOperator): PostfixUnaryExpression {
    return new NodeObject(SyntaxKind.PostfixUnaryExpression, {
        operand,
        operator,
    }) as unknown as PostfixUnaryExpression;
}

export function createPrefixUnaryExpression(operator: PrefixUnaryOperator, operand: UnaryExpression): PrefixUnaryExpression {
    return new NodeObject(SyntaxKind.PrefixUnaryExpression, {
        operator,
        operand,
    }) as unknown as PrefixUnaryExpression;
}

export function createPrivateIdentifier(escapedText: string): PrivateIdentifier {
    return new NodeObject(SyntaxKind.PrivateIdentifier, {
        escapedText,
    }) as unknown as PrivateIdentifier;
}

export function createPropertyAccessExpression(name: MemberName, expression: LeftHandSideExpression, questionDotToken?: QuestionDotToken): PropertyAccessExpression {
    return new NodeObject(SyntaxKind.PropertyAccessExpression, {
        name,
        expression,
        questionDotToken,
    }) as unknown as PropertyAccessExpression;
}

export function createPropertyAssignment(name: PropertyName, initializer: Expression): PropertyAssignment {
    return new NodeObject(SyntaxKind.PropertyAssignment, {
        name,
        initializer,
    }) as unknown as PropertyAssignment;
}

export function createPropertyDeclaration(name: PropertyName, modifiers?: NodeArray<ModifierLike>, questionToken?: QuestionToken, exclamationToken?: ExclamationToken, type?: TypeNode, initializer?: Expression): PropertyDeclaration {
    return new NodeObject(SyntaxKind.PropertyDeclaration, {
        name,
        modifiers,
        questionToken,
        exclamationToken,
        type,
        initializer,
    }) as unknown as PropertyDeclaration;
}

export function createPropertySignature(name: PropertyName, questionToken?: QuestionToken, modifiers?: NodeArray<Modifier>, type?: TypeNode): PropertySignature {
    return new NodeObject(SyntaxKind.PropertySignature, {
        name,
        questionToken,
        modifiers,
        type,
    }) as unknown as PropertySignature;
}

export function createQualifiedName(left: EntityName, right: Identifier): QualifiedName {
    return new NodeObject(SyntaxKind.QualifiedName, {
        left,
        right,
    }) as unknown as QualifiedName;
}

export function createRegularExpressionLiteral(text: string, isUnterminated?: boolean, hasExtendedUnicodeEscape?: boolean): RegularExpressionLiteral {
    return new NodeObject(SyntaxKind.RegularExpressionLiteral, {
        text,
        isUnterminated,
        hasExtendedUnicodeEscape,
    }) as unknown as RegularExpressionLiteral;
}

export function createRestTypeNode(type: TypeNode): RestTypeNode {
    return new NodeObject(SyntaxKind.RestType, {
        type,
    }) as unknown as RestTypeNode;
}

export function createReturnStatement(expression?: Expression): ReturnStatement {
    return new NodeObject(SyntaxKind.ReturnStatement, {
        expression,
    }) as unknown as ReturnStatement;
}

export function createSatisfiesExpression(expression: Expression, type: TypeNode): SatisfiesExpression {
    return new NodeObject(SyntaxKind.SatisfiesExpression, {
        expression,
        type,
    }) as unknown as SatisfiesExpression;
}

export function createSemicolonClassElement(name?: PropertyName): SemicolonClassElement {
    return new NodeObject(SyntaxKind.SemicolonClassElement, {
        name,
    }) as unknown as SemicolonClassElement;
}

export function createSetAccessorDeclaration(name: PropertyName, parameters: NodeArray<ParameterDeclaration>, typeParameters?: NodeArray<TypeParameterDeclaration> | undefined, type?: TypeNode | undefined, asteriskToken?: AsteriskToken | undefined, questionToken?: QuestionToken | undefined, exclamationToken?: ExclamationToken | undefined, body?: FunctionBody, modifiers?: NodeArray<ModifierLike>): SetAccessorDeclaration {
    return new NodeObject(SyntaxKind.SetAccessor, {
        name,
        parameters,
        typeParameters,
        type,
        asteriskToken,
        questionToken,
        exclamationToken,
        body,
        modifiers,
    }) as unknown as SetAccessorDeclaration;
}

export function createShorthandPropertyAssignment(name: Identifier, equalsToken?: EqualsToken, objectAssignmentInitializer?: Expression): ShorthandPropertyAssignment {
    return new NodeObject(SyntaxKind.ShorthandPropertyAssignment, {
        name,
        equalsToken,
        objectAssignmentInitializer,
    }) as unknown as ShorthandPropertyAssignment;
}

export function createSourceFile(statements: NodeArray<Statement>, endOfFileToken: EndOfFile, text: string, fileName: string, path: Path): SourceFile {
    return new NodeObject(SyntaxKind.SourceFile, {
        statements,
        endOfFileToken,
        text,
        fileName,
        path,
    }) as unknown as SourceFile;
}

export function createSpreadAssignment(expression: Expression, name?: PropertyName): SpreadAssignment {
    return new NodeObject(SyntaxKind.SpreadAssignment, {
        expression,
        name,
    }) as unknown as SpreadAssignment;
}

export function createSpreadElement(expression: Expression): SpreadElement {
    return new NodeObject(SyntaxKind.SpreadElement, {
        expression,
    }) as unknown as SpreadElement;
}

export function createStringLiteral(text: string, isUnterminated?: boolean, hasExtendedUnicodeEscape?: boolean): StringLiteral {
    return new NodeObject(SyntaxKind.StringLiteral, {
        text,
        isUnterminated,
        hasExtendedUnicodeEscape,
    }) as unknown as StringLiteral;
}

export function createSuperExpression(): SuperExpression {
    return new NodeObject(SyntaxKind.SuperKeyword, undefined) as unknown as SuperExpression;
}

export function createSwitchStatement(expression: Expression, caseBlock: CaseBlock, possiblyExhaustive?: boolean): SwitchStatement {
    return new NodeObject(SyntaxKind.SwitchStatement, {
        expression,
        caseBlock,
        possiblyExhaustive,
    }) as unknown as SwitchStatement;
}

export function createTaggedTemplateExpression(tag: LeftHandSideExpression, template: TemplateLiteral, typeArguments?: NodeArray<TypeNode>): TaggedTemplateExpression {
    return new NodeObject(SyntaxKind.TaggedTemplateExpression, {
        tag,
        template,
        typeArguments,
    }) as unknown as TaggedTemplateExpression;
}

export function createTemplateExpression(head: TemplateHead, templateSpans: NodeArray<TemplateSpan>): TemplateExpression {
    return new NodeObject(SyntaxKind.TemplateExpression, {
        head,
        templateSpans,
    }) as unknown as TemplateExpression;
}

export function createTemplateHead(text: string, templateFlags: TokenFlags, isUnterminated?: boolean, hasExtendedUnicodeEscape?: boolean, rawText?: string): TemplateHead {
    return new NodeObject(SyntaxKind.TemplateHead, {
        text,
        templateFlags,
        isUnterminated,
        hasExtendedUnicodeEscape,
        rawText,
    }) as unknown as TemplateHead;
}

export function createTemplateLiteralTypeNode(head: TemplateHead, templateSpans: NodeArray<TemplateLiteralTypeSpan>): TemplateLiteralTypeNode {
    return new NodeObject(SyntaxKind.TemplateLiteralType, {
        head,
        templateSpans,
    }) as unknown as TemplateLiteralTypeNode;
}

export function createTemplateLiteralTypeSpan(type: TypeNode, literal: TemplateMiddle | TemplateTail): TemplateLiteralTypeSpan {
    return new NodeObject(SyntaxKind.TemplateLiteralTypeSpan, {
        type,
        literal,
    }) as unknown as TemplateLiteralTypeSpan;
}

export function createTemplateMiddle(text: string, templateFlags: TokenFlags, isUnterminated?: boolean, hasExtendedUnicodeEscape?: boolean, rawText?: string): TemplateMiddle {
    return new NodeObject(SyntaxKind.TemplateMiddle, {
        text,
        templateFlags,
        isUnterminated,
        hasExtendedUnicodeEscape,
        rawText,
    }) as unknown as TemplateMiddle;
}

export function createTemplateSpan(expression: Expression, literal: TemplateMiddle | TemplateTail): TemplateSpan {
    return new NodeObject(SyntaxKind.TemplateSpan, {
        expression,
        literal,
    }) as unknown as TemplateSpan;
}

export function createTemplateTail(text: string, templateFlags: TokenFlags, isUnterminated?: boolean, hasExtendedUnicodeEscape?: boolean, rawText?: string): TemplateTail {
    return new NodeObject(SyntaxKind.TemplateTail, {
        text,
        templateFlags,
        isUnterminated,
        hasExtendedUnicodeEscape,
        rawText,
    }) as unknown as TemplateTail;
}

export function createThisExpression(): ThisExpression {
    return new NodeObject(SyntaxKind.ThisKeyword, undefined) as unknown as ThisExpression;
}

export function createThisTypeNode(): ThisTypeNode {
    return new NodeObject(SyntaxKind.ThisType, undefined) as unknown as ThisTypeNode;
}

export function createThrowStatement(expression: Expression): ThrowStatement {
    return new NodeObject(SyntaxKind.ThrowStatement, {
        expression,
    }) as unknown as ThrowStatement;
}

export function createTrueLiteral(): TrueLiteral {
    return new NodeObject(SyntaxKind.TrueKeyword, undefined) as unknown as TrueLiteral;
}

export function createTryStatement(tryBlock: Block, catchClause?: CatchClause, finallyBlock?: Block): TryStatement {
    return new NodeObject(SyntaxKind.TryStatement, {
        tryBlock,
        catchClause,
        finallyBlock,
    }) as unknown as TryStatement;
}

export function createTupleTypeNode(elements: NodeArray<TypeNode | NamedTupleMember>): TupleTypeNode {
    return new NodeObject(SyntaxKind.TupleType, {
        elements,
    }) as unknown as TupleTypeNode;
}

export function createTypeAliasDeclaration(name: Identifier, type: TypeNode, modifiers?: NodeArray<ModifierLike>, typeParameters?: NodeArray<TypeParameterDeclaration>): TypeAliasDeclaration {
    return new NodeObject(SyntaxKind.TypeAliasDeclaration, {
        name,
        type,
        modifiers,
        typeParameters,
    }) as unknown as TypeAliasDeclaration;
}

export function createTypeAssertion(type: TypeNode, expression: UnaryExpression): TypeAssertion {
    return new NodeObject(SyntaxKind.TypeAssertionExpression, {
        type,
        expression,
    }) as unknown as TypeAssertion;
}

export function createTypeLiteralNode(members: NodeArray<TypeElement>): TypeLiteralNode {
    return new NodeObject(SyntaxKind.TypeLiteral, {
        members,
    }) as unknown as TypeLiteralNode;
}

export function createTypeOfExpression(expression: UnaryExpression): TypeOfExpression {
    return new NodeObject(SyntaxKind.TypeOfExpression, {
        expression,
    }) as unknown as TypeOfExpression;
}

export function createTypeOperatorNode(operator: SyntaxKind.KeyOfKeyword | SyntaxKind.UniqueKeyword | SyntaxKind.ReadonlyKeyword, type: TypeNode): TypeOperatorNode {
    return new NodeObject(SyntaxKind.TypeOperator, {
        operator,
        type,
    }) as unknown as TypeOperatorNode;
}

export function createTypeParameterDeclaration(name: Identifier, modifiers?: NodeArray<Modifier>, constraint?: TypeNode, default_?: TypeNode, expression?: Expression): TypeParameterDeclaration {
    return new NodeObject(SyntaxKind.TypeParameter, {
        name,
        modifiers,
        constraint,
        default: default_,
        expression,
    }) as unknown as TypeParameterDeclaration;
}

export function createTypePredicateNode(parameterName: Identifier | ThisTypeNode, assertsModifier?: AssertsKeyword, type?: TypeNode): TypePredicateNode {
    return new NodeObject(SyntaxKind.TypePredicate, {
        parameterName,
        assertsModifier,
        type,
    }) as unknown as TypePredicateNode;
}

export function createTypeQueryNode(exprName: EntityName, typeArguments?: NodeArray<TypeNode>): TypeQueryNode {
    return new NodeObject(SyntaxKind.TypeQuery, {
        exprName,
        typeArguments,
    }) as unknown as TypeQueryNode;
}

export function createTypeReferenceNode(typeName: EntityName, typeArguments?: NodeArray<TypeNode>): TypeReferenceNode {
    return new NodeObject(SyntaxKind.TypeReference, {
        typeName,
        typeArguments,
    }) as unknown as TypeReferenceNode;
}

export function createUnionTypeNode(types: NodeArray<TypeNode>): UnionTypeNode {
    return new NodeObject(SyntaxKind.UnionType, {
        types,
    }) as unknown as UnionTypeNode;
}

export function createVariableDeclaration(name: BindingName, exclamationToken?: ExclamationToken, type?: TypeNode, initializer?: Expression): VariableDeclaration {
    return new NodeObject(SyntaxKind.VariableDeclaration, {
        name,
        exclamationToken,
        type,
        initializer,
    }) as unknown as VariableDeclaration;
}

export function createVariableDeclarationList(declarations: NodeArray<VariableDeclaration>): VariableDeclarationList {
    return new NodeObject(SyntaxKind.VariableDeclarationList, {
        declarations,
    }) as unknown as VariableDeclarationList;
}

export function createVariableStatement(declarationList: VariableDeclarationList, modifiers?: NodeArray<ModifierLike>): VariableStatement {
    return new NodeObject(SyntaxKind.VariableStatement, {
        declarationList,
        modifiers,
    }) as unknown as VariableStatement;
}

export function createVoidExpression(expression: UnaryExpression): VoidExpression {
    return new NodeObject(SyntaxKind.VoidExpression, {
        expression,
    }) as unknown as VoidExpression;
}

export function createWhileStatement(statement: Statement, expression: Expression): WhileStatement {
    return new NodeObject(SyntaxKind.WhileStatement, {
        statement,
        expression,
    }) as unknown as WhileStatement;
}

export function createWithStatement(expression: Expression, statement: Statement): WithStatement {
    return new NodeObject(SyntaxKind.WithStatement, {
        expression,
        statement,
    }) as unknown as WithStatement;
}

export function createYieldExpression(asteriskToken?: AsteriskToken, expression?: Expression): YieldExpression {
    return new NodeObject(SyntaxKind.YieldExpression, {
        asteriskToken,
        expression,
    }) as unknown as YieldExpression;
}
