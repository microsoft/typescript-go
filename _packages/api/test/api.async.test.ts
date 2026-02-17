import {
    API,
    SymbolFlags,
    TypeFlags,
} from "@typescript/api/async";
import { createVirtualFileSystem } from "@typescript/api/fs";
import {
    cast,
    isImportDeclaration,
    isNamedImports,
    isStringLiteral,
} from "@typescript/ast";
import assert from "node:assert";
import {
    describe,
    test,
} from "node:test";
import { fileURLToPath } from "node:url";
import { runBenchmarks } from "./api.async.bench.ts";

const defaultFiles = {
    "/tsconfig.json": "{}",
    "/src/index.ts": `import { foo } from './foo';`,
    "/src/foo.ts": `export const foo = 42;`,
};

describe("API", () => {
    test("parseConfigFile", async () => {
        const api = spawnAPI();
        try {
            const config = await api.parseConfigFile("/tsconfig.json");
            assert.deepEqual(config.fileNames, ["/src/index.ts", "/src/foo.ts"]);
            assert.deepEqual(config.options, { configFilePath: "/tsconfig.json" });
        }
        finally {
            await api.close();
        }
    });
});

describe("Snapshot", () => {
    test("updateSnapshot returns snapshot with projects", async () => {
        const api = spawnAPI();
        try {
            const snapshot = await api.updateSnapshot({ openProject: "/tsconfig.json" });
            assert.ok(snapshot);
            assert.ok(snapshot.id);
            assert.ok(snapshot.projects.length > 0);
            assert.ok(snapshot.projects[0].configFileName);
        }
        finally {
            await api.close();
        }
    });

    test("getSymbolAtPosition", async () => {
        const api = spawnAPI();
        try {
            const snapshot = await api.updateSnapshot({ openProject: "/tsconfig.json" });
            const project = snapshot.projects[0];
            const symbol = await project.checker.getSymbolAtPosition("/src/index.ts", 9);
            assert.ok(symbol);
            assert.equal(symbol.name, "foo");
            assert.ok(symbol.flags & SymbolFlags.Alias);
        }
        finally {
            await api.close();
        }
    });

    test("getSymbolAtLocation", async () => {
        const api = spawnAPI();
        try {
            const snapshot = await api.updateSnapshot({ openProject: "/tsconfig.json" });
            const project = snapshot.projects[0];
            const sourceFile = await project.program.getSourceFile("/src/index.ts");
            assert.ok(sourceFile);
            const node = cast(
                cast(sourceFile.statements[0], isImportDeclaration).importClause?.namedBindings,
                isNamedImports,
            ).elements[0].name;
            assert.ok(node);
            const symbol = await project.checker.getSymbolAtLocation(node);
            assert.ok(symbol);
            assert.equal(symbol.name, "foo");
            assert.ok(symbol.flags & SymbolFlags.Alias);
        }
        finally {
            await api.close();
        }
    });

    test("getTypeOfSymbol", async () => {
        const api = spawnAPI();
        try {
            const snapshot = await api.updateSnapshot({ openProject: "/tsconfig.json" });
            const project = snapshot.projects[0];
            const symbol = await project.checker.getSymbolAtPosition("/src/index.ts", 9);
            assert.ok(symbol);
            const type = await project.checker.getTypeOfSymbol(symbol);
            assert.ok(type);
            assert.ok(type.flags & TypeFlags.NumberLiteral);
        }
        finally {
            await api.close();
        }
    });
});

describe("SourceFile", () => {
    test("file properties", async () => {
        const api = spawnAPI();
        try {
            const snapshot = await api.updateSnapshot({ openProject: "/tsconfig.json" });
            const project = snapshot.projects[0];
            const sourceFile = await project.program.getSourceFile("/src/index.ts");

            assert.ok(sourceFile);
            assert.equal(sourceFile.text, defaultFiles["/src/index.ts"]);
            assert.equal(sourceFile.fileName, "/src/index.ts");
        }
        finally {
            await api.close();
        }
    });

    test("extended data", async () => {
        const api = spawnAPI();
        try {
            const snapshot = await api.updateSnapshot({ openProject: "/tsconfig.json" });
            const project = snapshot.projects[0];
            const sourceFile = await project.program.getSourceFile("/src/index.ts");

            assert.ok(sourceFile);
            let nodeCount = 1;
            sourceFile.forEachChild(function visit(node) {
                nodeCount++;
                node.forEachChild(visit);
            });
            assert.equal(nodeCount, 8);
        }
        finally {
            await api.close();
        }
    });
});

test("async unicode escapes", async () => {
    const api = spawnAPI({
        "/tsconfig.json": "{}",
        "/src/1.ts": `"😃"`,
        "/src/2.ts": `"\\ud83d\\ude03"`,
    });
    try {
        const snapshot = await api.updateSnapshot({ openProject: "/tsconfig.json" });
        const project = snapshot.projects[0];

        for (const file of ["/src/1.ts", "/src/2.ts"]) {
            const sourceFile = await project.program.getSourceFile(file);
            assert.ok(sourceFile);

            sourceFile.forEachChild(function visit(node) {
                if (isStringLiteral(node)) {
                    assert.equal(node.text, "😃");
                }
                node.forEachChild(visit);
            });
        }
    }
    finally {
        await api.close();
    }
});

test("async Object equality", async () => {
    const api = spawnAPI();
    try {
        const snapshot = await api.updateSnapshot({ openProject: "/tsconfig.json" });
        const project = snapshot.projects[0];
        // Same symbol returned from same snapshot's checker
        assert.strictEqual(
            await project.checker.getSymbolAtPosition("/src/index.ts", 9),
            await project.checker.getSymbolAtPosition("/src/index.ts", 10),
        );
    }
    finally {
        await api.close();
    }
});

test("async Snapshot dispose", async () => {
    const api = spawnAPI();
    try {
        const snapshot = await api.updateSnapshot({ openProject: "/tsconfig.json" });
        const project = snapshot.projects[0];
        const symbol = await project.checker.getSymbolAtPosition("/src/index.ts", 9);
        assert.ok(symbol);

        // Snapshot dispose should release server-side resources
        assert.ok(snapshot.isDisposed() === false);
        await snapshot.dispose();
        assert.ok(snapshot.isDisposed() === true);

        // After dispose, snapshot methods should throw
        assert.throws(() => {
            snapshot.getProject(project.id);
        }, {
            name: "Error",
            message: "Snapshot is disposed",
        });
    }
    finally {
        await api.close();
    }
});

describe("async Multiple snapshots", () => {
    test("two snapshots work independently", async () => {
        const api = spawnAPI();
        try {
            const snap1 = await api.updateSnapshot({ openProject: "/tsconfig.json" });
            const snap2 = await api.updateSnapshot({ openProject: "/tsconfig.json" });

            // Both can fetch source files
            const sf1 = await snap1.projects[0].program.getSourceFile("/src/index.ts");
            const sf2 = await snap2.projects[0].program.getSourceFile("/src/index.ts");
            assert.ok(sf1);
            assert.ok(sf2);

            // Disposing one doesn't break the other
            await snap1.dispose();
            assert.ok(snap1.isDisposed());
            assert.ok(!snap2.isDisposed());

            // snap2 still works after snap1 is disposed
            const symbol = await snap2.projects[0].checker.getSymbolAtPosition("/src/index.ts", 9);
            assert.ok(symbol);
            assert.equal(symbol.name, "foo");
        }
        finally {
            await api.close();
        }
    });

    test("each snapshot has its own server-side lifecycle", async () => {
        const api = spawnAPI();
        try {
            const snap1 = await api.updateSnapshot({ openProject: "/tsconfig.json" });
            const snap2 = await api.updateSnapshot({ openProject: "/tsconfig.json" });

            await snap1.dispose();

            // snap2 still works independently
            const symbol = await snap2.projects[0].checker.getSymbolAtPosition("/src/index.ts", 9);
            assert.ok(symbol);

            await snap2.dispose();

            // Both are disposed, new snapshot works fine
            const snap3 = await api.updateSnapshot({ openProject: "/tsconfig.json" });
            const sf = await snap3.projects[0].program.getSourceFile("/src/index.ts");
            assert.ok(sf);
        }
        finally {
            await api.close();
        }
    });
});

describe("async Source file caching", () => {
    test("same file from same snapshot returns cached object", async () => {
        const api = spawnAPI();
        try {
            const snapshot = await api.updateSnapshot({ openProject: "/tsconfig.json" });
            const project = snapshot.projects[0];
            const sf1 = await project.program.getSourceFile("/src/index.ts");
            const sf2 = await project.program.getSourceFile("/src/index.ts");
            assert.ok(sf1);
            assert.strictEqual(sf1, sf2, "Same source file should be returned from cache");
        }
        finally {
            await api.close();
        }
    });

    test("same file from two snapshots (same content) returns cached object", async () => {
        const api = spawnAPI();
        try {
            const snap1 = await api.updateSnapshot({ openProject: "/tsconfig.json" });
            const snap2 = await api.updateSnapshot({ openProject: "/tsconfig.json" });
            // Fetch from snap1 first (populates cache), then snap2 (cache hit via hash)
            const sf1 = await snap1.projects[0].program.getSourceFile("/src/index.ts");
            const sf2 = await snap2.projects[0].program.getSourceFile("/src/index.ts");
            assert.ok(sf1);
            assert.ok(sf2);
            // Same content hash → cache hit → same object
            assert.strictEqual(sf1, sf2, "Same file with same content should share cached object");
        }
        finally {
            await api.close();
        }
    });

    test("cache entries survive when one of two snapshots is disposed", async () => {
        const api = spawnAPI();
        try {
            const snap1 = await api.updateSnapshot({ openProject: "/tsconfig.json" });
            // Fetch from snap1 to populate cache
            const sf1 = await snap1.projects[0].program.getSourceFile("/src/index.ts");
            assert.ok(sf1);

            // snap2 retains snap1's cache refs for unchanged files via snapshot changes
            const snap2 = await api.updateSnapshot({ openProject: "/tsconfig.json" });

            // Dispose snap1 — snap2 still holds a ref, so the entry survives
            await snap1.dispose();

            // Fetching from snap2 should still return the cached object
            const sf2 = await snap2.projects[0].program.getSourceFile("/src/index.ts");
            assert.ok(sf2);
            assert.strictEqual(sf1, sf2, "Cache entry should survive when retained by the next snapshot");
        }
        finally {
            await api.close();
        }
    });
});

describe("async Snapshot disposal", () => {
    test("dispose is idempotent", async () => {
        const api = spawnAPI();
        try {
            const snapshot = await api.updateSnapshot({ openProject: "/tsconfig.json" });
            await snapshot.dispose();
            assert.ok(snapshot.isDisposed());
            // Second dispose should not throw
            await snapshot.dispose();
            assert.ok(snapshot.isDisposed());
        }
        finally {
            await api.close();
        }
    });

    test("api.close disposes all active snapshots", async () => {
        const api = spawnAPI();
        const snap1 = await api.updateSnapshot({ openProject: "/tsconfig.json" });
        const snap2 = await api.updateSnapshot({ openProject: "/tsconfig.json" });
        assert.ok(!snap1.isDisposed());
        assert.ok(!snap2.isDisposed());
        await api.close();
        assert.ok(snap1.isDisposed());
        assert.ok(snap2.isDisposed());
    });
});

test("async Benchmarks", async () => {
    await runBenchmarks(/*singleIteration*/ true);
});

function spawnAPI(files: Record<string, string> = defaultFiles) {
    return new API({
        cwd: fileURLToPath(new URL("../../../", import.meta.url).toString()),
        tsserverPath: fileURLToPath(new URL(`../../../built/local/tsgo${process.platform === "win32" ? ".exe" : ""}`, import.meta.url).toString()),
        fs: createVirtualFileSystem(files),
    });
}
