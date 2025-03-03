import * as fs from "node:fs";
import { dirname } from "node:path";
import * as path from "node:path";
import ts from "typescript";
import { API } from "./api.ts";

console.log("=== LSP-based tsgo ===");
console.time("Total execution time");

console.time("Server startup time");
const api = new API({
    tsserverPath: new URL("../../built/local/tsgo", import.meta.url).pathname,
    cwd: dirname(new URL(import.meta.url).pathname),
    // logServer: msg => logs.push(msg),
});
console.timeEnd("Server startup time");

// Start project loading timer
console.time("Project loading time");
const project = await api.loadProject("../../../eg/ts/tsconfig.json");
console.timeEnd("Project loading time");

console.time("Symbol lookup time");
const symbol = await project.getSymbolAtPosition("a.ts", 4);
console.timeEnd("Symbol lookup time");

console.log(symbol?.name, symbol?.flags);
await api.close();

console.timeEnd("Total execution time");

// Implement equivalent functionality with TypeScript API directly
console.log("\n=== Direct TypeScript API implementation ===");
console.time("Total execution time");

try {
    console.time("Config parsing");
    // Parse tsconfig.json
    const configFilePath = new URL("../../../eg/ts/tsconfig.json", import.meta.url).pathname;
    const configFileText = fs.readFileSync(configFilePath, "utf-8");
    const jsonSourceFile = ts.parseJsonText(configFilePath, configFileText);

    const parseConfigHost: ts.ParseConfigHost = {
        fileExists: ts.sys.fileExists,
        readFile: ts.sys.readFile,
        readDirectory: ts.sys.readDirectory,
        useCaseSensitiveFileNames: ts.sys.useCaseSensitiveFileNames,
    };

    const parsedConfig = ts.parseJsonSourceFileConfigFileContent(
        jsonSourceFile,
        parseConfigHost,
        path.dirname(configFilePath),
    );
    console.timeEnd("Config parsing");

    console.time("Program creation");
    // Create a program
    const program = ts.createProgram({
        rootNames: parsedConfig.fileNames,
        options: parsedConfig.options,
    });
    console.timeEnd("Program creation");

    console.time("Symbol lookup");
    // Get the checker and find symbol at position
    const checker = program.getTypeChecker();
    const aFilePath = path.resolve(path.dirname(configFilePath), "a.ts");
    const sourceFile = program.getSourceFile(aFilePath);

    if (sourceFile) {
        // Position 4 is in the first line
        const position = 4;

        // @ts-expect-error
        const token = ts.getTokenAtPosition(sourceFile, position);

        // Get symbol at position
        const symbol = checker.getSymbolAtLocation(token);
        console.log("Direct TS API symbol:", symbol?.getName(), symbol?.getFlags());
    }
    console.timeEnd("Symbol lookup");
}
catch (err) {
    console.error("Error in direct TS API implementation:", err);
}

console.timeEnd("Total execution time");
