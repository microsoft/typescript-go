//// [tests/cases/compiler/infiniteRecursionDestructuringLoop.ts] ////

//// [infiniteRecursionDestructuringLoop.ts]
// Repro from https://github.com/microsoft/TypeScript/issues/63192

interface Node {
    children?: readonly Node[];
    index?: number;
}

function IterateNodes(data: { node: Node }) {
    let node: Node | undefined = data.node;
    while (node) {
        const { children, index = -1 } = node;
        const activeNode: Node | undefined = index != -1 && children ? children[index] : undefined;

        node = activeNode;
    }
}

// Simplified repro
interface MyNode {
    children: MyNode[];
    index?: number;
}

function f(init: MyNode) {
    let node: MyNode | undefined = init;
    while (node) {
        const { children, index = 0 } = node;
        node = children[index];
    }
}


//// [infiniteRecursionDestructuringLoop.js]
"use strict";
// Repro from https://github.com/microsoft/TypeScript/issues/63192
function IterateNodes(data) {
    let node = data.node;
    while (node) {
        const { children, index = -1 } = node;
        const activeNode = index != -1 && children ? children[index] : undefined;
        node = activeNode;
    }
}
function f(init) {
    let node = init;
    while (node) {
        const { children, index = 0 } = node;
        node = children[index];
    }
}
