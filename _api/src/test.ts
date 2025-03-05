import * as fs from "node:fs";
import { dirname } from "node:path";
import * as path from "node:path";
import ts from "typescript";
import { API as AsyncAPI } from "./async/api.ts";
import { API as SyncAPI } from "./sync/api.ts";

{
    console.log("=== LSP-based tsgo ===");
    console.time("Total execution time");

    console.time("Server startup time");
    const api = new AsyncAPI({
        tsserverPath: new URL("../../built/local/tsgo", import.meta.url).pathname,
        cwd: dirname(new URL(import.meta.url).pathname),
        // logServer: msg => logs.push(msg),
    });
    console.timeEnd("Server startup time");

    // Start project loading timer
    console.time("Project loading time");
    const project = await api.loadProject("../../../eg/ts/tsconfig.json");
    console.timeEnd("Project loading time");

    console.time("Symbol lookup time 1");
    const symbol1 = await project.getSymbolAtPosition("a.ts", 4);
    console.timeEnd("Symbol lookup time 1");

    console.time("Symbol lookup time 2");
    const symbol2 = await project.getSymbolAtPosition("a.ts", 4);
    console.timeEnd("Symbol lookup time 2");

    console.log(symbol1?.name, symbol1?.flags);
    await api.close();

    console.timeEnd("Total execution time");
}

{
    console.log("\n=== libsyncrpc-based tsgo ===");
    console.time("Total execution time");

    console.time("Server startup time");
    const api = new SyncAPI({
        tsserverPath: new URL("../../built/local/tsgo", import.meta.url).pathname,
        cwd: dirname(new URL(import.meta.url).pathname),
        // logServer: msg => logs.push(msg),
    });
    console.timeEnd("Server startup time");

    console.time("Project loading time");
    const project = api.loadProject("../../../eg/ts/tsconfig.json");
    console.timeEnd("Project loading time");

    console.time("Symbol lookup time 1");
    const symbol1 = project.getSymbolAtPosition("a.ts", 4);
    console.timeEnd("Symbol lookup time 1");

    console.time("Symbol lookup time 2");
    const symbol2 = project.getSymbolAtPosition("a.ts", 4);
    console.timeEnd("Symbol lookup time 2");

    console.log(symbol1?.name, symbol1?.flags);
    api.close();
    console.timeEnd("Total execution time");
}

{
    console.log("\n=== Direct TypeScript API implementation ===");
    console.time("Total execution time");

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

    const aFilePath = path.resolve(path.dirname(configFilePath), "a.ts");
    const sourceFile = program.getSourceFile(aFilePath)!;
    // @ts-expect-error
    const token = ts.getTokenAtPosition(sourceFile, 4);

    console.time("Symbol lookup 1");
    // Get the checker and find symbol at position
    const checker = program.getTypeChecker();
    // Get symbol at position
    const symbol = checker.getSymbolAtLocation(token);
    console.timeEnd("Symbol lookup 1");

    console.time("Symbol lookup 2");
    const _ = program.getTypeChecker().getSymbolAtLocation(token);
    console.timeEnd("Symbol lookup 2");

    console.log("Direct TS API symbol:", symbol?.getName(), symbol?.getFlags());
    console.timeEnd("Total execution time");
}
