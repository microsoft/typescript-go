// is.ts — Hand-written type guard functions
// Generated guards are in is.generated.ts

import { NodeFlags } from "#enums/nodeFlags";
import { OuterExpressionKinds } from "#enums/outerExpressionKinds";
import { ScriptKind } from "#enums/scriptKind";
import { SyntaxKind } from "#enums/syntaxKind";
import type {
    AsExpression,
    AssertsKeyword,
    AsteriskToken,
    AwaitKeyword,
    BinaryOperatorToken,
    BindingName,
    BindingPattern,
    BooleanLiteral,
    ColonToken,
    ConciseBody,
    DeclarationName,
    DotDotDotToken,
    EndOfFile,
    EntityName,
    EqualsGreaterThanToken,
    EqualsToken,
    ExclamationToken,
    Expression,
    ExpressionWithTypeArguments,
    ForInitializer,
    FunctionBody,
    Identifier,
    JSDocTypeExpression,
    JSDocTypeLiteral,
    LeftHandSideExpression,
    LiteralExpression,
    MinusToken,
    Modifier,
    ModifierLike,
    ModuleBody,
    ModuleDeclaration,
    Node,
    NonNullExpression,
    NullLiteral,
    ParenthesizedExpression,
    PartiallyEmittedExpression,
    PlusToken,
    PrefixUnaryExpression,
    QuestionDotToken,
    QuestionToken,
    ReadonlyKeyword,
    SatisfiesExpression,
    Statement,
    TemplateMiddle,
    TemplateTail,
    ThisTypeNode,
    Token,
    TokenSyntaxKind,
    TypeAssertion,
    TypeNode,
    UnaryExpressionBase,
} from "./ast.ts";

export * from "./is.generated.ts";
import { isDecorator } from "./is.generated.ts";

type JSDocNamespaceDeclaration = ModuleDeclaration;

type WrappedExpression<T extends Expression> =
    | ParenthesizedExpression
    | TypeAssertion
    | AsExpression
    | SatisfiesExpression
    | ExpressionWithTypeArguments
    | NonNullExpression
    | PartiallyEmittedExpression;

type OuterExpression = WrappedExpression<Expression>;

export function isAssertsKeyword(node: Node): node is AssertsKeyword {
    return node.kind === SyntaxKind.AssertsKeyword;
}

export function isAsteriskToken(node: Node): node is AsteriskToken {
    return node.kind === SyntaxKind.AsteriskToken;
}

export function isAwaitKeyword(node: Node): node is AwaitKeyword {
    return node.kind === SyntaxKind.AwaitKeyword;
}

export function isColonToken(node: Node): node is ColonToken {
    return node.kind === SyntaxKind.ColonToken;
}

export function isDotDotDotToken(node: Node): node is DotDotDotToken {
    return node.kind === SyntaxKind.DotDotDotToken;
}

export function isEqualsGreaterThanToken(node: Node): node is EqualsGreaterThanToken {
    return node.kind === SyntaxKind.EqualsGreaterThanToken;
}

export function isExclamationToken(node: Node): node is ExclamationToken {
    return node.kind === SyntaxKind.ExclamationToken;
}

export function isQuestionDotToken(node: Node): node is QuestionDotToken {
    return node.kind === SyntaxKind.QuestionDotToken;
}

export function isQuestionToken(node: Node): node is QuestionToken {
    return node.kind === SyntaxKind.QuestionToken;
}

export function isTokenKind(kind: SyntaxKind): boolean {
    return kind >= SyntaxKind.FirstToken && kind <= SyntaxKind.LastToken;
}

export function isKeywordKind(kind: SyntaxKind): boolean {
    return kind >= SyntaxKind.FirstKeyword && kind <= SyntaxKind.LastKeyword;
}

export function isJSDocKind(kind: SyntaxKind): boolean {
    return kind >= SyntaxKind.FirstJSDocNode && kind <= SyntaxKind.LastJSDocNode;
}

export function isPropertyNameLiteral(node: Node): boolean {
    const kind = node.kind;
    return kind === SyntaxKind.Identifier
        || kind === SyntaxKind.StringLiteral
        || kind === SyntaxKind.NoSubstitutionTemplateLiteral
        || kind === SyntaxKind.NumericLiteral;
}

export function isTypeNode(node: Node): node is TypeNode {
    return isTypeNodeKind(node.kind);
}

function isTypeNodeKind(kind: SyntaxKind): boolean {
    return kind >= SyntaxKind.FirstTypeNode && kind <= SyntaxKind.LastTypeNode
        || kind === SyntaxKind.AnyKeyword
        || kind === SyntaxKind.UnknownKeyword
        || kind === SyntaxKind.NumberKeyword
        || kind === SyntaxKind.BigIntKeyword
        || kind === SyntaxKind.ObjectKeyword
        || kind === SyntaxKind.BooleanKeyword
        || kind === SyntaxKind.StringKeyword
        || kind === SyntaxKind.SymbolKeyword
        || kind === SyntaxKind.VoidKeyword
        || kind === SyntaxKind.UndefinedKeyword
        || kind === SyntaxKind.NeverKeyword
        || kind === SyntaxKind.IntrinsicKeyword
        || kind === SyntaxKind.ExpressionWithTypeArguments
        || kind === SyntaxKind.JSDocAllType
        || kind === SyntaxKind.JSDocNullableType
        || kind === SyntaxKind.JSDocNonNullableType
        || kind === SyntaxKind.JSDocOptionalType
        || kind === SyntaxKind.JSDocVariadicType
        || kind === SyntaxKind.JSDocTypeExpression
        || kind === SyntaxKind.JSDocTypeLiteral
        || kind === SyntaxKind.JSDocSignature;
}

export function isStatement(node: Node): node is Statement {
    const kind = node.kind;
    return kind === SyntaxKind.VariableStatement || kind === SyntaxKind.EmptyStatement
        || kind === SyntaxKind.ExpressionStatement || kind === SyntaxKind.IfStatement
        || kind === SyntaxKind.DoStatement || kind === SyntaxKind.WhileStatement
        || kind === SyntaxKind.ForStatement || kind === SyntaxKind.ForInStatement
        || kind === SyntaxKind.ForOfStatement || kind === SyntaxKind.ContinueStatement
        || kind === SyntaxKind.BreakStatement || kind === SyntaxKind.ReturnStatement
        || kind === SyntaxKind.WithStatement || kind === SyntaxKind.SwitchStatement
        || kind === SyntaxKind.LabeledStatement || kind === SyntaxKind.ThrowStatement
        || kind === SyntaxKind.TryStatement || kind === SyntaxKind.DebuggerStatement
        || kind === SyntaxKind.InterfaceDeclaration || kind === SyntaxKind.TypeAliasDeclaration
        || kind === SyntaxKind.EnumDeclaration || kind === SyntaxKind.ModuleDeclaration
        || kind === SyntaxKind.ImportDeclaration || kind === SyntaxKind.ImportEqualsDeclaration
        || kind === SyntaxKind.ExportDeclaration || kind === SyntaxKind.ExportAssignment
        || kind === SyntaxKind.NamespaceExportDeclaration || kind === SyntaxKind.FunctionDeclaration
        || kind === SyntaxKind.ClassDeclaration || kind === SyntaxKind.MissingDeclaration
        || kind === SyntaxKind.NotEmittedStatement || kind === SyntaxKind.Block;
}

export function isExpression(node: Node): node is Expression {
    const kind = node.kind;
    return kind === SyntaxKind.ConditionalExpression || kind === SyntaxKind.YieldExpression
        || kind === SyntaxKind.ArrowFunction || kind === SyntaxKind.BinaryExpression
        || kind === SyntaxKind.SpreadElement || kind === SyntaxKind.AsExpression
        || kind === SyntaxKind.OmittedExpression
        || kind === SyntaxKind.SatisfiesExpression
        || kind === SyntaxKind.PrefixUnaryExpression || kind === SyntaxKind.PostfixUnaryExpression
        || kind === SyntaxKind.DeleteExpression || kind === SyntaxKind.TypeOfExpression
        || kind === SyntaxKind.VoidExpression || kind === SyntaxKind.AwaitExpression
        || kind === SyntaxKind.TypeAssertionExpression
        || kind === SyntaxKind.CallExpression || kind === SyntaxKind.NewExpression
        || kind === SyntaxKind.TaggedTemplateExpression || kind === SyntaxKind.NonNullExpression
        || kind === SyntaxKind.MetaProperty || kind === SyntaxKind.JsxExpression
        || kind === SyntaxKind.PropertyAccessExpression || kind === SyntaxKind.ElementAccessExpression
        || kind === SyntaxKind.FunctionExpression || kind === SyntaxKind.ClassExpression
        || kind === SyntaxKind.ParenthesizedExpression || kind === SyntaxKind.ArrayLiteralExpression
        || kind === SyntaxKind.ObjectLiteralExpression || kind === SyntaxKind.TemplateExpression
        || kind === SyntaxKind.Identifier
        || kind === SyntaxKind.PrivateIdentifier
        || kind === SyntaxKind.NumericLiteral || kind === SyntaxKind.BigIntLiteral
        || kind === SyntaxKind.StringLiteral || kind === SyntaxKind.RegularExpressionLiteral
        || kind === SyntaxKind.NoSubstitutionTemplateLiteral || kind === SyntaxKind.JsxElement
        || kind === SyntaxKind.JsxSelfClosingElement || kind === SyntaxKind.JsxFragment
        || kind === SyntaxKind.NullKeyword || kind === SyntaxKind.TrueKeyword
        || kind === SyntaxKind.FalseKeyword || kind === SyntaxKind.ThisKeyword
        || kind === SyntaxKind.SuperKeyword || kind === SyntaxKind.ImportKeyword
        || kind === SyntaxKind.ExpressionWithTypeArguments;
}

export function isLeftHandSideExpression(node: Node): node is LeftHandSideExpression {
    return isLeftHandSideExpressionKind(skipPartiallyEmittedExpressions(node).kind);
}

export function skipPartiallyEmittedExpressions(node: Expression): Expression;
export function skipPartiallyEmittedExpressions(node: Node): Node;
export function skipPartiallyEmittedExpressions(node: Node) {
    return skipOuterExpressions(node, OuterExpressionKinds.PartiallyEmittedExpressions);
}

function isLeftHandSideExpressionKind(kind: SyntaxKind): boolean {
    switch (kind) {
        case SyntaxKind.PropertyAccessExpression:
        case SyntaxKind.ElementAccessExpression:
        case SyntaxKind.NewExpression:
        case SyntaxKind.CallExpression:
        case SyntaxKind.JsxElement:
        case SyntaxKind.JsxSelfClosingElement:
        case SyntaxKind.JsxFragment:
        case SyntaxKind.TaggedTemplateExpression:
        case SyntaxKind.ArrayLiteralExpression:
        case SyntaxKind.ParenthesizedExpression:
        case SyntaxKind.ObjectLiteralExpression:
        case SyntaxKind.ClassExpression:
        case SyntaxKind.FunctionExpression:
        case SyntaxKind.Identifier:
        case SyntaxKind.PrivateIdentifier: // technically this is only an Expression if it's in a `#field in expr` BinaryExpression
        case SyntaxKind.RegularExpressionLiteral:
        case SyntaxKind.NumericLiteral:
        case SyntaxKind.BigIntLiteral:
        case SyntaxKind.StringLiteral:
        case SyntaxKind.NoSubstitutionTemplateLiteral:
        case SyntaxKind.TemplateExpression:
        case SyntaxKind.FalseKeyword:
        case SyntaxKind.NullKeyword:
        case SyntaxKind.ThisKeyword:
        case SyntaxKind.TrueKeyword:
        case SyntaxKind.SuperKeyword:
        case SyntaxKind.NonNullExpression:
        case SyntaxKind.ExpressionWithTypeArguments:
        case SyntaxKind.MetaProperty:
        case SyntaxKind.ImportKeyword: // technically this is only an Expression if it's in a CallExpression
        case SyntaxKind.MissingDeclaration:
            return true;
        default:
            return false;
    }
}

export function isUnaryExpression(node: Node): node is UnaryExpressionBase {
    return isUnaryExpressionKind(skipPartiallyEmittedExpressions(node).kind);
}

function isUnaryExpressionKind(kind: SyntaxKind): boolean {
    switch (kind) {
        case SyntaxKind.PrefixUnaryExpression:
        case SyntaxKind.PostfixUnaryExpression:
        case SyntaxKind.DeleteExpression:
        case SyntaxKind.TypeOfExpression:
        case SyntaxKind.VoidExpression:
        case SyntaxKind.AwaitExpression:
        case SyntaxKind.TypeAssertionExpression:
            return true;
        default:
            return isLeftHandSideExpressionKind(kind);
    }
}

export function isModifier(node: Node): node is Modifier {
    return isModifierKind(node.kind);
}

function isModifierKind(kind: SyntaxKind): boolean {
    return kind === SyntaxKind.AbstractKeyword || kind === SyntaxKind.AccessorKeyword
        || kind === SyntaxKind.AsyncKeyword || kind === SyntaxKind.ConstKeyword
        || kind === SyntaxKind.DeclareKeyword || kind === SyntaxKind.DefaultKeyword
        || kind === SyntaxKind.ExportKeyword || kind === SyntaxKind.InKeyword
        || kind === SyntaxKind.PrivateKeyword || kind === SyntaxKind.ProtectedKeyword
        || kind === SyntaxKind.PublicKeyword || kind === SyntaxKind.ReadonlyKeyword
        || kind === SyntaxKind.OutKeyword || kind === SyntaxKind.OverrideKeyword
        || kind === SyntaxKind.StaticKeyword;
}

/** @internal */
export function isOuterExpression(node: Node, kinds: OuterExpressionKinds = OuterExpressionKinds.All): node is OuterExpression {
    switch (node.kind) {
        case SyntaxKind.ParenthesizedExpression:
            if (kinds & OuterExpressionKinds.ExcludeJSDocTypeAssertion && isJSDocTypeAssertion(node as ParenthesizedExpression)) {
                return false;
            }
            return (kinds & OuterExpressionKinds.Parentheses) !== 0;
        case SyntaxKind.TypeAssertionExpression:
        case SyntaxKind.AsExpression:
            return (kinds & OuterExpressionKinds.TypeAssertions) !== 0;
        case SyntaxKind.SatisfiesExpression:
            return (kinds & (OuterExpressionKinds.ExpressionsWithTypeArguments | OuterExpressionKinds.Satisfies)) !== 0;
        case SyntaxKind.ExpressionWithTypeArguments:
            return (kinds & OuterExpressionKinds.ExpressionsWithTypeArguments) !== 0;
        case SyntaxKind.NonNullExpression:
            return (kinds & OuterExpressionKinds.NonNullAssertions) !== 0;
        case SyntaxKind.PartiallyEmittedExpression:
            return (kinds & OuterExpressionKinds.PartiallyEmittedExpressions) !== 0;
    }
    return false;
}

/** @internal */
export function skipOuterExpressions<T extends Expression>(node: WrappedExpression<T>): T;
/** @internal */
export function skipOuterExpressions(node: Expression, kinds?: OuterExpressionKinds): Expression;
/** @internal */
export function skipOuterExpressions(node: Node, kinds?: OuterExpressionKinds): Node;
/** @internal */
export function skipOuterExpressions(node: Node, kinds = OuterExpressionKinds.All) {
    while (isOuterExpression(node, kinds)) {
        node = node.expression;
    }
    return node;
}

function isJSDocTypeAssertion(node: ParenthesizedExpression): boolean {
    const sourceFile = node.getSourceFile();
    if (sourceFile.scriptKind !== ScriptKind.JS && sourceFile.scriptKind !== ScriptKind.JSX) {
        return false;
    }
    const expression = node.expression;
    if (expression.kind !== SyntaxKind.AsExpression) {
        return false;
    }
    const asExpression = expression as AsExpression;
    return !!asExpression.type
        && (asExpression.type.flags & NodeFlags.Reparsed) !== 0;
}

export function isModifierLike(node: Node): node is ModifierLike {
    return isModifier(node) || isDecorator(node);
}

export function isBindingName(node: Node): node is BindingName {
    return node.kind === SyntaxKind.Identifier || node.kind === SyntaxKind.ObjectBindingPattern || node.kind === SyntaxKind.ArrayBindingPattern;
}

export function isBindingPattern(node: Node): node is BindingPattern {
    return node.kind === SyntaxKind.ObjectBindingPattern || node.kind === SyntaxKind.ArrayBindingPattern;
}

export function isConciseBody(node: Node): node is ConciseBody {
    return node.kind === SyntaxKind.Block || isExpression(node);
}

export function isToken(node: Node): node is Token<TokenSyntaxKind> {
    return isTokenKind(node.kind);
}

export function isEndOfFile(node: Node): node is EndOfFile {
    return node.kind === SyntaxKind.EndOfFile;
}

export function isEqualsToken(node: Node): node is EqualsToken {
    return node.kind === SyntaxKind.EqualsToken;
}

export function isFunctionBody(node: Node): node is FunctionBody {
    return node.kind === SyntaxKind.Block;
}

export function isForInitializer(node: Node): node is ForInitializer {
    return node.kind === SyntaxKind.VariableDeclarationList || isExpression(node);
}

export function isQuestionOrExclamationToken(node: Node): node is QuestionToken | ExclamationToken {
    return node.kind === SyntaxKind.QuestionToken || node.kind === SyntaxKind.ExclamationToken;
}

export function isIdentifierOrThisTypeNode(node: Node): node is Identifier | ThisTypeNode {
    return node.kind === SyntaxKind.Identifier || node.kind === SyntaxKind.ThisType;
}

export function isReadonlyKeywordOrPlusOrMinusToken(node: Node): node is ReadonlyKeyword | PlusToken | MinusToken {
    return node.kind === SyntaxKind.ReadonlyKeyword || node.kind === SyntaxKind.PlusToken || node.kind === SyntaxKind.MinusToken;
}

export function isQuestionOrPlusOrMinusToken(node: Node): node is QuestionToken | PlusToken | MinusToken {
    return node.kind === SyntaxKind.QuestionToken || node.kind === SyntaxKind.PlusToken || node.kind === SyntaxKind.MinusToken;
}

export function isModuleBody(node: Node): node is ModuleBody {
    return node.kind === SyntaxKind.ModuleBlock || node.kind === SyntaxKind.ModuleDeclaration;
}

export function isTemplateMiddleOrTemplateTail(node: Node): node is TemplateMiddle | TemplateTail {
    return node.kind === SyntaxKind.TemplateMiddle || node.kind === SyntaxKind.TemplateTail;
}

export function isLiteralExpression(node: Node): node is LiteralExpression {
    return node.kind >= SyntaxKind.FirstLiteralToken && node.kind <= SyntaxKind.LastLiteralToken;
}

export function isLiteralTypeLiteral(node: Node): node is NullLiteral | BooleanLiteral | LiteralExpression | PrefixUnaryExpression {
    const kind = node.kind;
    return kind === SyntaxKind.NullKeyword || kind === SyntaxKind.TrueKeyword
        || kind === SyntaxKind.FalseKeyword || kind === SyntaxKind.PrefixUnaryExpression
        || isLiteralExpression(node);
}

export function isBinaryOperatorToken(node: Node): node is BinaryOperatorToken {
    return isBinaryOperatorKind(node.kind);
}

function isBinaryOperatorKind(kind: SyntaxKind): boolean {
    return kind >= SyntaxKind.FirstBinaryOperator && kind <= SyntaxKind.LastBinaryOperator;
}

export function isEntityName(node: Node): node is EntityName {
    return node.kind === SyntaxKind.Identifier || node.kind === SyntaxKind.QualifiedName;
}

export function isIdentifierOrJSDocNamespaceDeclaration(node: Node): node is Identifier | JSDocNamespaceDeclaration {
    return node.kind === SyntaxKind.Identifier || node.kind === SyntaxKind.ModuleDeclaration;
}

export function isJSDocTypeExpressionOrJSDocTypeLiteral(node: Node): node is JSDocTypeExpression | JSDocTypeLiteral {
    return node.kind === SyntaxKind.JSDocTypeExpression || node.kind === SyntaxKind.JSDocTypeLiteral;
}

export function isDeclarationName(node: Node): node is DeclarationName {
    const kind = node.kind;
    return kind === SyntaxKind.Identifier || kind === SyntaxKind.PrivateIdentifier
        || kind === SyntaxKind.StringLiteral || kind === SyntaxKind.NumericLiteral
        || kind === SyntaxKind.ComputedPropertyName || kind === SyntaxKind.ElementAccessExpression
        || kind === SyntaxKind.BindingElement || kind === SyntaxKind.NoSubstitutionTemplateLiteral
        || kind === SyntaxKind.ObjectBindingPattern || kind === SyntaxKind.ArrayBindingPattern
        || kind === SyntaxKind.JsxNamespacedName;
}
