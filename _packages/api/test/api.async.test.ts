import {
    API,
    SymbolFlags,
    TypeFlags,
} from "@typescript/api/async";
import {
    cast,
    isImportDeclaration,
    isNamedImports,
    isStringLiteral,
} from "@typescript/ast";
import assert from "node:assert";
import fs from "node:fs";
import path from "node:path";
import {
    after,
    before,
    describe,
    test,
} from "node:test";
import { fileURLToPath } from "node:url";
import { runBenchmarks } from "./api.async.bench.ts";

// Create a temp directory with test files
const repoRoot = fileURLToPath(new URL("../../../", import.meta.url).toString());
const fixtureDir = path.join(repoRoot, "testdata/fixtures/async-api-test");
const indexContent = `import { foo } from './foo';`;
const fooContent = `export const foo = 42;`;

before(() => {
    fs.mkdirSync(path.join(fixtureDir, "src"), { recursive: true });
    fs.writeFileSync(path.join(fixtureDir, "tsconfig.json"), "{}");
    fs.writeFileSync(path.join(fixtureDir, "src/index.ts"), indexContent);
    fs.writeFileSync(path.join(fixtureDir, "src/foo.ts"), fooContent);
});

after(() => {
    fs.rmSync(fixtureDir, { recursive: true, force: true });
});

describe("API", () => {
    test("parseConfigFile", async () => {
        const api = spawnAPI();
        try {
            const config = await api.parseConfigFile(path.join(fixtureDir, "tsconfig.json"));
            assert.ok(config.fileNames.some(f => f.endsWith("src/index.ts")));
            assert.ok(config.fileNames.some(f => f.endsWith("src/foo.ts")));
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
            const snapshot = await api.updateSnapshot({ openProject: path.join(fixtureDir, "tsconfig.json") });
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
            const snapshot = await api.updateSnapshot({ openProject: path.join(fixtureDir, "tsconfig.json") });
            const project = snapshot.projects[0];
            const symbol = await project.checker.getSymbolAtPosition(path.join(fixtureDir, "src/index.ts"), 9);
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
            const snapshot = await api.updateSnapshot({ openProject: path.join(fixtureDir, "tsconfig.json") });
            const project = snapshot.projects[0];
            const sourceFile = await project.program.getSourceFile(path.join(fixtureDir, "src/index.ts"));
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
            const snapshot = await api.updateSnapshot({ openProject: path.join(fixtureDir, "tsconfig.json") });
            const project = snapshot.projects[0];
            const symbol = await project.checker.getSymbolAtPosition(path.join(fixtureDir, "src/index.ts"), 9);
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
            const snapshot = await api.updateSnapshot({ openProject: path.join(fixtureDir, "tsconfig.json") });
            const project = snapshot.projects[0];
            const sourceFile = await project.program.getSourceFile(path.join(fixtureDir, "src/index.ts"));

            assert.ok(sourceFile);
            assert.equal(sourceFile.text, indexContent);
            assert.ok(sourceFile.fileName.endsWith("src/index.ts"));
        }
        finally {
            await api.close();
        }
    });

    test("extended data", async () => {
        const api = spawnAPI();
        try {
            const snapshot = await api.updateSnapshot({ openProject: path.join(fixtureDir, "tsconfig.json") });
            const project = snapshot.projects[0];
            const sourceFile = await project.program.getSourceFile(path.join(fixtureDir, "src/index.ts"));

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
    const unicodeDir = path.join(fixtureDir, "unicode");
    fs.mkdirSync(unicodeDir, { recursive: true });
    fs.writeFileSync(path.join(unicodeDir, "tsconfig.json"), "{}");
    fs.writeFileSync(path.join(unicodeDir, "1.ts"), `"😃"`);
    fs.writeFileSync(path.join(unicodeDir, "2.ts"), `"\\ud83d\\ude03"`);

    const api = spawnAPI();
    try {
        const snapshot = await api.updateSnapshot({ openProject: path.join(unicodeDir, "tsconfig.json") });
        const project = snapshot.projects[0];

        for (const file of ["1.ts", "2.ts"]) {
            const sourceFile = await project.program.getSourceFile(path.join(unicodeDir, file));
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
        fs.rmSync(unicodeDir, { recursive: true, force: true });
    }
});

test("async Object equality", async () => {
    const api = spawnAPI();
    try {
        const snapshot = await api.updateSnapshot({ openProject: path.join(fixtureDir, "tsconfig.json") });
        const project = snapshot.projects[0];
        // Same symbol returned from same snapshot's checker
        assert.strictEqual(
            await project.checker.getSymbolAtPosition(path.join(fixtureDir, "src/index.ts"), 9),
            await project.checker.getSymbolAtPosition(path.join(fixtureDir, "src/index.ts"), 10),
        );
    }
    finally {
        await api.close();
    }
});

test("async Snapshot dispose", async () => {
    const api = spawnAPI();
    try {
        const snapshot = await api.updateSnapshot({ openProject: path.join(fixtureDir, "tsconfig.json") });
        const project = snapshot.projects[0];
        const symbol = await project.checker.getSymbolAtPosition(path.join(fixtureDir, "src/index.ts"), 9);
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

test("async Benchmarks", async () => {
    await runBenchmarks(/*singleIteration*/ true);
});

function spawnAPI() {
    return new API({
        cwd: repoRoot,
        tsserverPath: fileURLToPath(new URL(`../../../built/local/tsgo${process.platform === "win32" ? ".exe" : ""}`, import.meta.url).toString()),
    });
}
