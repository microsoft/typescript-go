//// [tests/cases/compiler/cssTreeTypeInference.ts] ////

//// [cssTreeTypeInference.ts]
// Simplified reproduction of css-tree type inference issue
// https://github.com/microsoft/typescript-go/issues/1727

interface Declaration {
    type: 'Declaration';
    property: string;
    value: string;
}

interface Rule {
    type: 'Rule';
    selector: string;
    children: Declaration[];
}

type ASTNode = Declaration | Rule;

interface WalkOptions<T extends ASTNode> {
    visit: T['type'];
    enter(node: T): void;
}

declare function walk<T extends ASTNode>(ast: ASTNode, options: WalkOptions<T>): void;

// Test case 1: Simple type inference
const ast: ASTNode = {
    type: 'Declaration',
    property: 'color',
    value: 'red'
};

// This should infer node as Declaration type
walk(ast, {
    visit: 'Declaration',
    enter(node) {
        console.log(node.property); // Should not error - node should be inferred as Declaration
    },
});

// Test case 2: More complex scenario
declare const complexAst: Rule;

walk(complexAst, {
    visit: 'Declaration',
    enter(node) {
        console.log(node.value); // Should infer node as Declaration
    },
});

//// [cssTreeTypeInference.js]
// Test case 1: Simple type inference
const ast = {
    type: 'Declaration',
    property: 'color',
    value: 'red'
};
// This should infer node as Declaration type
walk(ast, {
    visit: 'Declaration',
    enter(node) {
        console.log(node.property); // Should not error - node should be inferred as Declaration
    },
});
walk(complexAst, {
    visit: 'Declaration',
    enter(node) {
        console.log(node.value); // Should infer node as Declaration
    },
});
