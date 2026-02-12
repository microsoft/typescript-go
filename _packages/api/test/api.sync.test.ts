import { createVirtualFileSystem } from "@typescript/api/fs";
import {
    API,
    type Snapshot,
    SymbolFlags,
    TypeFlags,
} from "@typescript/api/sync";
import {
    cast,
    isImportDeclaration,
    isNamedImports,
    isStringLiteral,
    isTemplateHead,
    isTemplateMiddle,
    isTemplateTail,
} from "@typescript/ast";
import assert from "node:assert";
import {
    describe,
    test,
} from "node:test";
import { fileURLToPath } from "node:url";
import { runBenchmarks } from "./api.sync.bench.ts";

const defaultFiles = {
    "/tsconfig.json": "{}",
    "/src/index.ts": `import { foo } from './foo';`,
    "/src/foo.ts": `export const foo = 42;`,
};

describe("API", () => {
    test("parseConfigFile", () => {
        const api = spawnAPI();
        const config = api.parseConfigFile("/tsconfig.json");
        assert.deepEqual(config.fileNames, ["/src/index.ts", "/src/foo.ts"]);
        assert.deepEqual(config.options, { configFilePath: "/tsconfig.json" });
    });
});

describe("Snapshot", () => {
    test("updateSnapshot returns snapshot with projects", () => {
        const api = spawnAPI();
        const snapshot = api.updateSnapshot({ openProject: "/tsconfig.json" });
        assert.ok(snapshot);
        assert.ok(snapshot.id);
        assert.ok(snapshot.projects.length > 0);
        assert.ok(snapshot.projects[0].configFileName);
    });

    test("getSymbolAtPosition", () => {
        const api = spawnAPI();
        const snapshot = api.updateSnapshot({ openProject: "/tsconfig.json" });
        const project = snapshot.projects[0];
        const symbol = project.checker.getSymbolAtPosition("/src/index.ts", 9);
        assert.ok(symbol);
        assert.equal(symbol.name, "foo");
        assert.ok(symbol.flags & SymbolFlags.Alias);
    });

    test("getSymbolAtLocation", () => {
        const api = spawnAPI();
        const snapshot = api.updateSnapshot({ openProject: "/tsconfig.json" });
        const project = snapshot.projects[0];
        const sourceFile = project.program.getSourceFile("/src/index.ts");
        assert.ok(sourceFile);
        const node = cast(
            cast(sourceFile.statements[0], isImportDeclaration).importClause?.namedBindings,
            isNamedImports,
        ).elements[0].name;
        assert.ok(node);
        const symbol = project.checker.getSymbolAtLocation(node);
        assert.ok(symbol);
        assert.equal(symbol.name, "foo");
        assert.ok(symbol.flags & SymbolFlags.Alias);
    });

    test("getTypeOfSymbol", () => {
        const api = spawnAPI();
        const snapshot = api.updateSnapshot({ openProject: "/tsconfig.json" });
        const project = snapshot.projects[0];
        const symbol = project.checker.getSymbolAtPosition("/src/index.ts", 9);
        assert.ok(symbol);
        const type = project.checker.getTypeOfSymbol(symbol);
        assert.ok(type);
        assert.ok(type.flags & TypeFlags.NumberLiteral);
    });
});

describe("SourceFile", () => {
    test("file properties", () => {
        const api = spawnAPI();
        const snapshot = api.updateSnapshot({ openProject: "/tsconfig.json" });
        const project = snapshot.projects[0];
        const sourceFile = project.program.getSourceFile("/src/index.ts");

        assert.ok(sourceFile);
        assert.equal(sourceFile.text, defaultFiles["/src/index.ts"]);
        assert.equal(sourceFile.fileName, "/src/index.ts");
    });

    test("extended data", () => {
        const api = spawnAPI();
        const snapshot = api.updateSnapshot({ openProject: "/tsconfig.json" });
        const project = snapshot.projects[0];
        const sourceFile = project.program.getSourceFile("/src/index.ts");

        assert.ok(sourceFile);
        let nodeCount = 1;
        sourceFile.forEachChild(function visit(node) {
            if (isTemplateHead(node)) {
                assert.equal(node.text, "head ");
                assert.equal(node.rawText, "head ");
                assert.equal(node.templateFlags, 0);
            }
            else if (isTemplateMiddle(node)) {
                assert.equal(node.text, "middle");
                assert.equal(node.rawText, "middle");
                assert.equal(node.templateFlags, 0);
            }
            else if (isTemplateTail(node)) {
                assert.equal(node.text, " tail");
                assert.equal(node.rawText, " tail");
                assert.equal(node.templateFlags, 0);
            }
            nodeCount++;
            node.forEachChild(visit);
        });
        assert.equal(nodeCount, 8);
    });
});

test("unicode escapes", () => {
    const srcFiles = {
        "/src/1.ts": `"😃"`,
        "/src/2.ts": `"\\ud83d\\ude03"`, // this is "😃"
    };

    const api = spawnAPI({
        "/tsconfig.json": "{}",
        ...srcFiles,
    });
    const snapshot = api.updateSnapshot({ openProject: "/tsconfig.json" });
    const project = snapshot.projects[0];

    Object.keys(srcFiles).forEach(file => {
        const sourceFile = project.program.getSourceFile(file);
        assert.ok(sourceFile);

        sourceFile.forEachChild(function visit(node) {
            if (isStringLiteral(node)) {
                assert.equal(node.text, "😃");
            }
            node.forEachChild(visit);
        });
    });
});

test("Object equality", () => {
    const api = spawnAPI();
    const snapshot = api.updateSnapshot({ openProject: "/tsconfig.json" });
    const project = snapshot.projects[0];
    // Same symbol returned from same snapshot's checker
    assert.strictEqual(
        project.checker.getSymbolAtPosition("/src/index.ts", 9),
        project.checker.getSymbolAtPosition("/src/index.ts", 10),
    );
});

test("Snapshot dispose", () => {
    const api = spawnAPI();
    const snapshot = api.updateSnapshot({ openProject: "/tsconfig.json" });
    const project = snapshot.projects[0];
    const symbol = project.checker.getSymbolAtPosition("/src/index.ts", 9);
    assert.ok(symbol);

    // Snapshot dispose should release server-side resources
    assert.ok(snapshot.isDisposed() === false);
    snapshot.dispose();
    assert.ok(snapshot.isDisposed() === true);

    // After dispose, snapshot methods should throw
    assert.throws(() => {
        snapshot.getProject(project.id);
    }, {
        name: "Error",
        message: "Snapshot is disposed",
    });
});

test("Server-side release", () => {
    const api = spawnAPI();
    const snapshot = api.updateSnapshot({ openProject: "/tsconfig.json" });
    const project = snapshot.projects[0];
    const symbol = project.checker.getSymbolAtPosition("/src/index.ts", 9);
    assert.ok(symbol);

    // Manually release the snapshot on the server
    // @ts-ignore private API
    api.client.request("release", { handle: snapshot.id });

    // Symbol handle should no longer be resolvable on the server
    assert.throws(() => {
        project.checker.getTypeOfSymbol(symbol);
    }, {
        name: "Error",
        message: `api: client error: snapshot ${snapshot.id} not found`,
    });
});

test("Benchmarks", async () => {
    await runBenchmarks(/*singleIteration*/ true);
});

function spawnAPI(files: Record<string, string> = defaultFiles) {
    return new API({
        cwd: fileURLToPath(new URL("../../../", import.meta.url).toString()),
        tsserverPath: fileURLToPath(new URL(`../../../built/local/tsgo${process.platform === "win32" ? ".exe" : ""}`, import.meta.url).toString()),
        fs: createVirtualFileSystem(files),
    });
}
