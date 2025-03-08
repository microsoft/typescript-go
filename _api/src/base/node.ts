import { SyntaxKind } from "../syntaxKind.ts";
import {
    childProperties,
    NodeChildMask,
    type NodeDataType,
    NodeDataTypeChildren,
    NodeDataTypeMask,
} from "./decode.ts";

const popcount8 = [0, 1, 1, 2, 1, 2, 2, 3, 1, 2, 2, 3, 2, 3, 3, 4, 1, 2, 2, 3, 2, 3, 3, 4, 2, 3, 3, 4, 3, 4, 4, 5, 1, 2, 2, 3, 2, 3, 3, 4, 2, 3, 3, 4, 3, 4, 4, 5, 2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6, 1, 2, 2, 3, 2, 3, 3, 4, 2, 3, 3, 4, 3, 4, 4, 5, 2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6, 2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6, 3, 4, 4, 5, 4, 5, 5, 6, 4, 5, 5, 6, 5, 6, 6, 7, 1, 2, 2, 3, 2, 3, 3, 4, 2, 3, 3, 4, 3, 4, 4, 5, 2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6, 2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6, 3, 4, 4, 5, 4, 5, 5, 6, 4, 5, 5, 6, 5, 6, 6, 7, 2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6, 3, 4, 4, 5, 4, 5, 5, 6, 4, 5, 5, 6, 5, 6, 6, 7, 3, 4, 4, 5, 4, 5, 5, 6, 4, 5, 5, 6, 5, 6, 6, 7, 4, 5, 5, 6, 5, 6, 6, 7, 5, 6, 6, 7, 6, 7, 7, 8];

const OFFSET_KIND = 0;
const OFFSET_POS = 1;
const OFFSET_END = 2;
const OFFSET_NEXT = 3;
const OFFSET_PARENT = 4;
const OFFSET_DATA = 5;
const NODE_LEN = 6;

const KIND_NODE_LIST = 2 ** 32 - 1;

export class NodeBase {
    protected parent: Node;
    protected view: DataView;
    protected index: number;
    /** Keys are positions */
    protected _children: Map<number, Node | NodeList> | undefined;

    constructor(view: DataView, index: number, parent: Node) {
        this.view = view;
        this.index = index;
        this.parent = parent;
    }

    get kind(): SyntaxKind {
        return this.view.getUint32(this.index * NODE_LEN * 4 + OFFSET_KIND * 4, true);
    }

    get pos(): number {
        return this.view.getUint32(this.index * NODE_LEN * 4 + OFFSET_POS * 4, true);
    }

    get end(): number {
        return this.view.getUint32(this.index * NODE_LEN * 4 + OFFSET_END * 4, true);
    }

    get next(): number {
        return this.view.getUint32(this.index * NODE_LEN * 4 + OFFSET_NEXT * 4, true);
    }

    protected get parentIndex(): number {
        return this.view.getUint32(this.index * NODE_LEN * 4 + OFFSET_PARENT * 4, true);
    }

    protected get data(): number {
        return this.view.getUint32(this.index * NODE_LEN * 4 + OFFSET_DATA * 4, true);
    }

    protected get dataType(): NodeDataType {
        return (this.data & NodeDataTypeMask) as NodeDataType;
    }

    protected get childMask(): number {
        if (this.dataType !== NodeDataTypeChildren) {
            return 0;
        }
        return this.data & NodeChildMask;
    }
}

export class NodeList extends NodeBase implements Iterable<Node> {
    get length(): number {
        return this.view.getUint32(this.index * NODE_LEN * 4 + OFFSET_DATA * 4, true);
    }

    *[Symbol.iterator](): IterableIterator<Node> {
        let next = this.index + 1;
        while (next) {
            const child = this.getOrCreateChildAtNodeIndex(next);
            next = child.next;
            yield child as Node;
        }
    }

    at(index: number): Node {
        let next = this.index + 1;
        for (let i = 0; i < index; i++) {
            const child = this.getOrCreateChildAtNodeIndex(next);
            next = child.next;
        }
        return this.getOrCreateChildAtNodeIndex(next) as Node;
    }

    private getOrCreateChildAtNodeIndex(index: number): Node | NodeList {
        const pos = this.view.getUint32(index * NODE_LEN * 4 + OFFSET_POS * 4, true);
        let child = (this._children ??= new Map()).get(pos);
        if (!child) {
            const kind = this.view.getUint32(index * NODE_LEN * 4 + OFFSET_KIND * 4, true);
            if (kind === KIND_NODE_LIST) {
                throw new Error("NodeList cannot directly contain another NodeList");
            }
            child = new Node(this.view, index, this.parent);
            this._children.set(pos, child);
        }
        return child;
    }
}

export class Node extends NodeBase {
    protected static NODE_LEN: number = NODE_LEN;

    forEachChild<T>(visitor: (node: Node) => T): T | undefined {
        if (this.hasChildren()) {
            let next = this.index + 1;
            do {
                const child = this.getOrCreateChildAtNodeIndex(next);
                if (child instanceof NodeList) {
                    for (const node of child) {
                        const result = visitor(node);
                        if (result) {
                            return result;
                        }
                    }
                }
                else {
                    const result = visitor(child);
                    if (result) {
                        return result;
                    }
                }
                next = child.next;
            }
            while (next);
        }
    }

    private getOrCreateChildAtNodeIndex(index: number): Node | NodeList {
        const pos = this.view.getUint32(index * NODE_LEN * 4 + OFFSET_POS * 4, true);
        let child = (this._children ??= new Map()).get(pos);
        if (!child) {
            const kind = this.view.getUint32(index * NODE_LEN * 4 + OFFSET_KIND * 4, true);
            child = kind === KIND_NODE_LIST ? new NodeList(this.view, index, this) : new Node(this.view, index, this);
            this._children.set(pos, child);
        }
        return child;
    }

    private hasChildren(): boolean {
        if (this._children) {
            return true;
        }
        if (this.index === this.view.byteLength / (NODE_LEN * 4) - 1) {
            return false;
        }
        const nextNodeParent = this.view.getUint32((this.index + 1) * NODE_LEN * 4 + OFFSET_PARENT * 4, true);
        return nextNodeParent === this.index;
    }

    private getNamedChild(propertyName: string): Node | NodeList | undefined {
        const propertyNames = childProperties[this.kind];
        if (!propertyNames) {
            // `childProperties` is only defined for nodes with more than one child property.
            // Get the only child if it exists.
            const child = this.getOrCreateChildAtNodeIndex(this.index + 1);
            if (child.next !== 0) {
                throw new Error("Expected only one child");
            }
            return child;
        }

        let order = propertyNames.indexOf(propertyName);
        if (order === -1) {
            // JSDocPropertyTag and JSDocParameterTag need special handling
            // because they have a conditional property order
            const kind = this.kind;
            if (kind === SyntaxKind.JSDocPropertyTag) {
                switch (propertyName) {
                    case "name":
                        order = this.isNameFirst ? 0 : 1;
                        break;
                    case "typeExpression":
                        order = this.isNameFirst ? 1 : 0;
                        break;
                }
            }
            else if (kind === SyntaxKind.JSDocParameterTag) {
                switch (propertyName) {
                    case "name":
                        order = this.isNameFirst ? 1 : 2;
                    case "typeExpression":
                        order = this.isNameFirst ? 2 : 1;
                }
            }
            // Node kind does not have this property
            return undefined;
        }
        const mask = this.childMask;
        if (!(mask & (1 << order))) {
            // Property is not present
            return undefined;
        }

        const propertyIndex = order - popcount8[mask & ((1 << order) - 1)];
        return this.getOrCreateChildAtNodeIndex(this.index + 1 + propertyIndex);
    }

    // Boolean properties
    get isArrayType(): boolean | undefined {
        switch (this.kind) {
            case SyntaxKind.JSDocTypeLiteral:
                return (this.data & 1 << 24) !== 0;
        }
    }

    get isTypeOnly(): boolean | undefined {
        switch (this.kind) {
            case SyntaxKind.ImportSpecifier:
            case SyntaxKind.ImportClause:
            case SyntaxKind.ExportSpecifier:
            case SyntaxKind.ImportEqualsDeclaration:
            case SyntaxKind.ExportDeclaration:
                return (this.data & 1 << 24) !== 0;
        }
    }

    get isTypeOf(): boolean | undefined {
        switch (this.kind) {
            case SyntaxKind.ImportType:
                return (this.data & 1 << 24) !== 0;
        }
    }

    get multiline(): boolean | undefined {
        switch (this.kind) {
            case SyntaxKind.Block:
            case SyntaxKind.ArrayLiteralExpression:
            case SyntaxKind.ObjectLiteralExpression:
            case SyntaxKind.ImportAttributes:
                return (this.data & 1 << 24) !== 0;
        }
    }

    get isExportEquals(): boolean | undefined {
        switch (this.kind) {
            case SyntaxKind.ExportAssignment:
                return (this.data & 1 << 24) !== 0;
        }
    }

    get isBracketed(): boolean | undefined {
        switch (this.kind) {
            case SyntaxKind.JSDocPropertyTag:
            case SyntaxKind.JSDocParameterTag:
                return (this.data & 1 << 24) !== 0;
        }
    }

    get containsOnlyTriviaWhiteSpaces(): boolean | undefined {
        switch (this.kind) {
            case SyntaxKind.JsxText:
                return (this.data & 1 << 24) !== 0;
        }
    }

    get isNameFirst(): boolean | undefined {
        switch (this.kind) {
            case SyntaxKind.JSDocPropertyTag:
            case SyntaxKind.JSDocParameterTag:
                return (this.data & 1 << 25) !== 0;
        }
    }

    // Children properties
    get argument(): Node | undefined {
        return this.getNamedChild("argument") as Node;
    }
    get argumentExpression(): Node | undefined {
        return this.getNamedChild("argumentExpression") as Node;
    }
    get arguments(): NodeList | undefined {
        return this.getNamedChild("arguments") as NodeList;
    }
    get assertsModifier(): Node | undefined {
        return this.getNamedChild("assertsModifier") as Node;
    }
    get asteriskToken(): Node | undefined {
        return this.getNamedChild("asteriskToken") as Node;
    }
    get attributes(): Node | undefined {
        return this.getNamedChild("attributes") as Node;
    }
    get awaitModifier(): Node | undefined {
        return this.getNamedChild("awaitModifier") as Node;
    }
    get block(): Node | undefined {
        return this.getNamedChild("block") as Node;
    }
    get body(): Node | undefined {
        return this.getNamedChild("body") as Node;
    }
    get caseBlock(): Node | undefined {
        return this.getNamedChild("caseBlock") as Node;
    }
    get catchClause(): Node | undefined {
        return this.getNamedChild("catchClause") as Node;
    }
    get checkType(): Node | undefined {
        return this.getNamedChild("checkType") as Node;
    }
    get children(): NodeList | undefined {
        return this.getNamedChild("children") as NodeList;
    }
    get className(): Node | undefined {
        return this.getNamedChild("className") as Node;
    }
    get closingElement(): Node | undefined {
        return this.getNamedChild("closingElement") as Node;
    }
    get closingFragment(): Node | undefined {
        return this.getNamedChild("closingFragment") as Node;
    }
    get colonToken(): Node | undefined {
        return this.getNamedChild("colonToken") as Node;
    }
    get comment(): Node | undefined {
        return this.getNamedChild("comment") as Node;
    }
    get condition(): Node | undefined {
        return this.getNamedChild("condition") as Node;
    }
    get constraint(): Node | undefined {
        return this.getNamedChild("constraint") as Node;
    }
    get declarationList(): Node | undefined {
        return this.getNamedChild("declarationList") as Node;
    }
    get defaultType(): Node | undefined {
        return this.getNamedChild("defaultType") as Node;
    }
    get dotDotDotToken(): Node | undefined {
        return this.getNamedChild("dotDotDotToken") as Node;
    }
    get elseStatement(): Node | undefined {
        return this.getNamedChild("elseStatement") as Node;
    }
    get equalsGreaterThanToken(): Node | undefined {
        return this.getNamedChild("equalsGreaterThanToken") as Node;
    }
    get equalsToken(): Node | undefined {
        return this.getNamedChild("equalsToken") as Node;
    }
    get exclamationToken(): Node | undefined {
        return this.getNamedChild("exclamationToken") as Node;
    }
    get exportClause(): Node | undefined {
        return this.getNamedChild("exportClause") as Node;
    }
    get expression(): Node | undefined {
        return this.getNamedChild("expression") as Node;
    }
    get exprName(): Node | undefined {
        return this.getNamedChild("exprName") as Node;
    }
    get extendsType(): Node | undefined {
        return this.getNamedChild("extendsType") as Node;
    }
    get falseType(): Node | undefined {
        return this.getNamedChild("falseType") as Node;
    }
    get finallyBlock(): Node | undefined {
        return this.getNamedChild("finallyBlock") as Node;
    }
    get fullName(): Node | undefined {
        return this.getNamedChild("fullName") as Node;
    }
    get head(): Node | undefined {
        return this.getNamedChild("head") as Node;
    }
    get heritageClauses(): NodeList | undefined {
        return this.getNamedChild("heritageClauses") as NodeList;
    }
    get importClause(): Node | undefined {
        return this.getNamedChild("importClause") as Node;
    }
    get incrementor(): Node | undefined {
        return this.getNamedChild("incrementor") as Node;
    }
    get indexType(): Node | undefined {
        return this.getNamedChild("indexType") as Node;
    }
    get initializer(): Node | undefined {
        return this.getNamedChild("initializer") as Node;
    }
    get label(): Node | undefined {
        return this.getNamedChild("label") as Node;
    }
    get left(): Node | undefined {
        return this.getNamedChild("left") as Node;
    }
    get literal(): Node | undefined {
        return this.getNamedChild("literal") as Node;
    }
    get members(): NodeList | undefined {
        return this.getNamedChild("members") as NodeList;
    }
    get modifiers(): NodeList | undefined {
        return this.getNamedChild("modifiers") as NodeList;
    }
    get moduleReference(): Node | undefined {
        return this.getNamedChild("moduleReference") as Node;
    }
    get moduleSpecifier(): Node | undefined {
        return this.getNamedChild("moduleSpecifier") as Node;
    }
    get name(): Node | undefined {
        return this.getNamedChild("name") as Node;
    }
    get namedBindings(): Node | undefined {
        return this.getNamedChild("namedBindings") as Node;
    }
    get nameExpression(): Node | undefined {
        return this.getNamedChild("nameExpression") as Node;
    }
    get namespace(): Node | undefined {
        return this.getNamedChild("namespace") as Node;
    }
    get nameType(): Node | undefined {
        return this.getNamedChild("nameType") as Node;
    }
    get objectAssignmentInitializer(): Node | undefined {
        return this.getNamedChild("objectAssignmentInitializer") as Node;
    }
    get objectType(): Node | undefined {
        return this.getNamedChild("objectType") as Node;
    }
    get openingElement(): Node | undefined {
        return this.getNamedChild("openingElement") as Node;
    }
    get openingFragment(): Node | undefined {
        return this.getNamedChild("openingFragment") as Node;
    }
    get operatorToken(): Node | undefined {
        return this.getNamedChild("operatorToken") as Node;
    }
    get parameterName(): Node | undefined {
        return this.getNamedChild("parameterName") as Node;
    }
    get parameters(): NodeList | undefined {
        return this.getNamedChild("parameters") as NodeList;
    }
    get postfixToken(): Node | undefined {
        return this.getNamedChild("postfixToken") as Node;
    }
    get propertyName(): Node | undefined {
        return this.getNamedChild("propertyName") as Node;
    }
    get qualifier(): Node | undefined {
        return this.getNamedChild("qualifier") as Node;
    }
    get questionDotToken(): Node | undefined {
        return this.getNamedChild("questionDotToken") as Node;
    }
    get questionToken(): Node | undefined {
        return this.getNamedChild("questionToken") as Node;
    }
    get readonlyToken(): Node | undefined {
        return this.getNamedChild("readonlyToken") as Node;
    }
    get right(): Node | undefined {
        return this.getNamedChild("right") as Node;
    }
    get statement(): Node | undefined {
        return this.getNamedChild("statement") as Node;
    }
    get statements(): NodeList | undefined {
        return this.getNamedChild("statements") as NodeList;
    }
    get tag(): Node | undefined {
        return this.getNamedChild("tag") as Node;
    }
    get tagName(): Node | undefined {
        return this.getNamedChild("tagName") as Node;
    }
    get tags(): NodeList | undefined {
        return this.getNamedChild("tags") as NodeList;
    }
    get template(): Node | undefined {
        return this.getNamedChild("template") as Node;
    }
    get templateSpans(): NodeList | undefined {
        return this.getNamedChild("templateSpans") as NodeList;
    }
    get thenStatement(): Node | undefined {
        return this.getNamedChild("thenStatement") as Node;
    }
    get trueType(): Node | undefined {
        return this.getNamedChild("trueType") as Node;
    }
    get tryBlock(): Node | undefined {
        return this.getNamedChild("tryBlock") as Node;
    }
    get type(): Node | undefined {
        return this.getNamedChild("type") as Node;
    }
    get typeArguments(): Node | undefined {
        return this.getNamedChild("typeArguments") as Node;
    }
    get typeExpression(): Node | undefined {
        return this.getNamedChild("typeExpression") as Node;
    }
    get typeName(): Node | undefined {
        return this.getNamedChild("typeName") as Node;
    }
    get typeParameter(): Node | undefined {
        return this.getNamedChild("typeParameter") as Node;
    }
    get typeParameters(): NodeList | undefined {
        return this.getNamedChild("typeParameters") as NodeList;
    }
    get value(): Node | undefined {
        return this.getNamedChild("value") as Node;
    }
    get variableDeclaration(): Node | undefined {
        return this.getNamedChild("variableDeclaration") as Node;
    }
    get whenFalse(): Node | undefined {
        return this.getNamedChild("whenFalse") as Node;
    }
    get whenTrue(): Node | undefined {
        return this.getNamedChild("whenTrue") as Node;
    }

    // Other properties
    get flags(): number {
        switch (this.kind) {
            case SyntaxKind.VariableDeclarationList:
                return this.data & (1 << 24 | 1 << 25) >> 24;
            default:
                return 0;
        }
    }

    get token(): SyntaxKind | undefined {
        switch (this.kind) {
            case SyntaxKind.ImportAttributes:
                if ((this.data & 1 << 25) !== 0) {
                    return SyntaxKind.AssertKeyword;
                }
                return SyntaxKind.WithKeyword;
        }
    }
}
