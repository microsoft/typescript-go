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

describe("Multiple snapshots", () => {
    test("two snapshots work independently", () => {
        const api = spawnAPI();
        const snap1 = api.updateSnapshot({ openProject: "/tsconfig.json" });
        const snap2 = api.updateSnapshot({ openProject: "/tsconfig.json" });

        // Both can fetch source files
        const sf1 = snap1.projects[0].program.getSourceFile("/src/index.ts");
        const sf2 = snap2.projects[0].program.getSourceFile("/src/index.ts");
        assert.ok(sf1);
        assert.ok(sf2);

        // Disposing one doesn't break the other
        snap1.dispose();
        assert.ok(snap1.isDisposed());
        assert.ok(!snap2.isDisposed());

        // snap2 still works after snap1 is disposed
        const symbol = snap2.projects[0].checker.getSymbolAtPosition("/src/index.ts", 9);
        assert.ok(symbol);
        assert.equal(symbol.name, "foo");

        snap2.dispose();
        api.close();
    });

    test("each snapshot has its own server-side lifecycle", () => {
        const api = spawnAPI();
        const snap1 = api.updateSnapshot({ openProject: "/tsconfig.json" });
        const snap2 = api.updateSnapshot({ openProject: "/tsconfig.json" });

        snap1.dispose();

        // snap2 still works independently
        const symbol = snap2.projects[0].checker.getSymbolAtPosition("/src/index.ts", 9);
        assert.ok(symbol);

        snap2.dispose();

        // Both are disposed, new snapshot works fine
        const snap3 = api.updateSnapshot({ openProject: "/tsconfig.json" });
        const sf = snap3.projects[0].program.getSourceFile("/src/index.ts");
        assert.ok(sf);

        api.close();
    });
});

describe("Source file caching", () => {
    test("same file from same snapshot returns cached object", () => {
        const api = spawnAPI();
        const snapshot = api.updateSnapshot({ openProject: "/tsconfig.json" });
        const project = snapshot.projects[0];
        const sf1 = project.program.getSourceFile("/src/index.ts");
        const sf2 = project.program.getSourceFile("/src/index.ts");
        assert.ok(sf1);
        assert.strictEqual(sf1, sf2, "Same source file should be returned from cache");
        api.close();
    });

    test("same file from two snapshots (same content) returns cached object", () => {
        const api = spawnAPI();
        const snap1 = api.updateSnapshot({ openProject: "/tsconfig.json" });
        const snap2 = api.updateSnapshot({ openProject: "/tsconfig.json" });
        // Fetch from snap1 first (populates cache), then snap2 (cache hit via hash)
        const sf1 = snap1.projects[0].program.getSourceFile("/src/index.ts");
        const sf2 = snap2.projects[0].program.getSourceFile("/src/index.ts");
        assert.ok(sf1);
        assert.ok(sf2);
        // Same content hash → cache hit → same object
        assert.strictEqual(sf1, sf2, "Same file with same content should share cached object");
        api.close();
    });

    test("cache entries survive when one of two snapshots is disposed", () => {
        const api = spawnAPI();
        const snap1 = api.updateSnapshot({ openProject: "/tsconfig.json" });
        // Fetch from snap1 to populate cache
        const sf1 = snap1.projects[0].program.getSourceFile("/src/index.ts");
        assert.ok(sf1);

        // snap2 retains snap1's cache refs for unchanged files via snapshot changes
        const snap2 = api.updateSnapshot({ openProject: "/tsconfig.json" });

        // Dispose snap1 — snap2 still holds a ref, so the entry survives
        snap1.dispose();

        // Fetching from snap2 should still return the cached object
        const sf2 = snap2.projects[0].program.getSourceFile("/src/index.ts");
        assert.ok(sf2);
        assert.strictEqual(sf1, sf2, "Cache entry should survive when retained by the next snapshot");
        api.close();
    });
});

describe("Snapshot disposal", () => {
    test("dispose is idempotent", () => {
        const api = spawnAPI();
        const snapshot = api.updateSnapshot({ openProject: "/tsconfig.json" });
        snapshot.dispose();
        assert.ok(snapshot.isDisposed());
        // Second dispose should not throw
        snapshot.dispose();
        assert.ok(snapshot.isDisposed());
        api.close();
    });

    test("api.close disposes all active snapshots", () => {
        const api = spawnAPI();
        const snap1 = api.updateSnapshot({ openProject: "/tsconfig.json" });
        const snap2 = api.updateSnapshot({ openProject: "/tsconfig.json" });
        assert.ok(!snap1.isDisposed());
        assert.ok(!snap2.isDisposed());
        api.close();
        assert.ok(snap1.isDisposed());
        assert.ok(snap2.isDisposed());
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
