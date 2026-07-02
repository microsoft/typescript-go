// @strict: true
// @noEmit: true
// @target: es2022
// @module: commonjs
// @skipLibCheck: true

// @filename: /types.d.ts
declare module "estree" {
    export interface Program extends BaseNode { body: Array<Directive | Statement | ModuleDeclaration>; }
    export interface Directive extends BaseNode { body: BlockStatement | Expression; }
    export type Statement = Declaration;
    export type Declaration = FunctionDeclaration | VariableDeclaration | ClassDeclaration;
    export interface FunctionDeclaration extends MaybeNamedFunctionDeclaration { id: Identifier; }
    export interface VariableDeclaration extends BaseDeclaration { consequent: Statement[]; }
    export type Literal = SimpleLiteral | RegExpLiteral | BigIntLiteral;
    export interface SimpleLiteral extends BaseNode, BaseExpression { value: string | boolean | number | null; }
    export interface RegExpLiteral extends BaseNode, BaseExpression { value?: RegExp | null | undefined; regex: {}; }
    export interface BigIntLiteral extends BaseNode, BaseExpression { raw?: string | undefined; }
    export interface ClassDeclaration extends MaybeNamedClassDeclaration { property: Identifier; }
    export type ModuleDeclaration = ExportAllDeclaration;
    export interface ExportAllDeclaration extends BaseModuleDeclaration { source: Literal; }
}

declare module "estree-jsx" {
    export * from "estree";
}

declare module "hast" {
    export type ElementContent = ElementContentMap[keyof ElementContentMap];
    export type RootContent = RootContentMap[keyof RootContentMap];
    export interface RootContentMap { comment: Comment; doctype: Doctype; element: Element; text: Text; }
    export interface Comment extends Literal { type: "comment"; data?: CommentData | undefined; }
    export interface Doctype extends UnistNode { type: "doctype"; data?: DoctypeData | undefined; }
    export interface Element extends Parent { type: "element"; tagName: string; properties: Properties; children: ElementContent[]; content?: Root | undefined; data?: ElementData | undefined; }
    export interface Root extends Parent { type: "root"; children: RootContent[]; data?: RootData | undefined; }
    export interface Text extends Literal { type: "text"; data?: TextData | undefined; }
}

declare module "mdast-util-mdx-jsx" {
    import type { Program } from "estree-jsx";
    export interface MdxJsxExpressionAttribute extends Node { estree?: Program | null | undefined; }
    export interface MdxJsxAttribute extends Node { data?: MdxJsxAttributeData | undefined; }
    export interface MdxJsxTextElementHast extends HastParent { attributes: Array<MdxJsxAttribute | MdxJsxExpressionAttribute>; }
}

declare module "hast" {
    import type { MdxJsxTextElementHast } from "mdast-util-mdx-jsx";
    interface ElementContentMap { mdxJsxTextElement: MdxJsxTextElementHast; }
}

// @filename: /repro.ts
import "mdast-util-mdx-jsx";
import type { Root } from "hast";

type Transform<Base, From, To> = {
    [K in keyof Base]: Exclude<Base[K], undefined | null> extends never ? Base[K]
        : Exclude<Base[K], undefined | null> extends object ? Exclude<Base[K], undefined | null> extends From ? To | Extract<Base[K], null | undefined> : Transform<Base[K], From, To>
        : Base[K];
};
type TransformJson<T> = Transform<T, Date | bigint, string>;

declare function fetchJson<T>(): TransformJson<T>;
declare function parseDate(s: string): Date;

interface IMessage { createdAt: Date; content: Root; }
interface IAggregate { createdAt: Date; messages: IMessage[]; }

export function getAggregate(): IAggregate {
    const response = fetchJson<TransformJson<IAggregate>>();
    return { ...response, createdAt: parseDate(response.createdAt), messages: response.messages.map(message => ({ ...message, createdAt: parseDate(message.createdAt) })) };
}
