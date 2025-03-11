import { SyntaxKind } from "#ast/syntax";
import {
    bench,
    run,
} from "mitata";
import { readFileSync } from "node:fs";
import { dirname } from "node:path";
import ts from "typescript";
import { isIdentifier } from "./ast/nodeTests.ts";
import { SymbolFlags } from "./base/api.ts";
import { API } from "./sync/api.ts";
{
    const api = new API({
        tsserverPath: new URL("../../built/local/tsgo", import.meta.url).pathname,
        cwd: dirname(new URL(import.meta.url).pathname),
        logFile: "tsgo.log",
    });
    const project = api.loadProject("../../../TypeScript/src/compiler/tsconfig.json");
    const file = project.getSourceFile("debug.ts")!;
    // bench("native - batched", () => {
    const symbolRequests: { fileName: string; position: number; }[] = [];
    file.forEachChild(function visitNode(node) {
        if (isIdentifier(node)) {
            symbolRequests.push({ fileName: "debug.ts", position: node.pos });
        }
        node.forEachChild(child => visitNode(child));
    });

    project.getSymbolAtPosition(symbolRequests);
    // });

    // bench("native - many calls", () => {
    // const symbolRequests: { fileName: string; position: number; }[] = [];
    // file.forEachChild(function visitNode(node) {
    //     if (isIdentifier(node)) {
    //         project.getSymbolAtPosition("debug.ts", node.pos);
    //     }
    //     node.forEachChild(child => visitNode(child));
    // });
    // project.getSymbolAtPosition(symbolRequests);
    // });
}

// {
//     const configFilePath = new URL("../../../TypeScript/src/compiler/tsconfig.json", import.meta.url).pathname;
//     const configFileText = readFileSync(configFilePath, "utf-8");
//     const jsonSourceFile = ts.parseJsonText(configFilePath, configFileText);

//     const parseConfigHost: ts.ParseConfigHost = {
//         fileExists: ts.sys.fileExists,
//         readFile: ts.sys.readFile,
//         readDirectory: ts.sys.readDirectory,
//         useCaseSensitiveFileNames: ts.sys.useCaseSensitiveFileNames,
//     };

//     const parsedConfig = ts.parseJsonSourceFileConfigFileContent(
//         jsonSourceFile,
//         parseConfigHost,
//         dirname(configFilePath),
//     );

//     const program = ts.createProgram({
//         rootNames: parsedConfig.fileNames,
//         options: parsedConfig.options,
//     });

//     const checker = program.getTypeChecker();
//     const file = program.getSourceFile("/Users/andrew/Developer/microsoft/TypeScript/src/compiler/checker.ts")!;
//     bench("js", () => {
//         file.forEachChild(function visitNode(node) {
//             if (ts.isIdentifier(node)) {
//                 const symbol = checker.getSymbolAtLocation(node);
//                 // if (symbol?.flags! & SymbolFlags.Value) {
//                 //     counts.types++;
//                 //     checker.getTypeOfSymbolAtLocation(symbol!, node);
//                 // }
//             }
//             node.forEachChild(child => visitNode(child));
//         });
//     });
// }
// await run();
