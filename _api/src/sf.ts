import { dirname } from "node:path";
import { API } from "./sync/api.ts";
import { SyntaxKind } from "./syntaxKind.ts";

const api = new API({
    tsserverPath: new URL("../../built/local/tsgo", import.meta.url).pathname,
    cwd: dirname(new URL(import.meta.url).pathname),
});

const project = api.loadProject("../../../TypeScript/src/compiler/tsconfig.json");

console.time("getSourceFile");
const file = project.getSourceFile("corePublic.ts")!;
console.timeEnd("getSourceFile");
console.log("");

console.log("node count:", file.nodeCount());
console.log();

console.log(file.statements!.at(0).declarationList);

// console.log(SyntaxKind[file.kind], file.pos, file.end);
// file.forEachChild(function visitNode(node, depth = 0) {
//     console.log(" ".repeat(depth), SyntaxKind[node.kind] ?? "(NodeList)", node.pos, node.end);
//     node.forEachChild(child => visitNode(child, depth + 1));
// });
