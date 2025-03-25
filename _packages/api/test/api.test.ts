import { API } from "@typescript/api";
import { createVirtualFileSystem } from "@typescript/api/fs";
import {
    isTemplateHead,
    isTemplateMiddle,
    isTemplateTail,
} from "@typescript/ast";
import assert from "node:assert";
import {
    describe,
    test,
} from "node:test";

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

describe("SourceFile", () => {
    test("file text", () => {
        const api = spawnAPI();
        const project = api.loadProject("/tsconfig.json");
        const sourceFile = project.getSourceFile("/src/index.ts");

        assert.ok(sourceFile);
        assert.equal(sourceFile.text, defaultFiles["/src/index.ts"]);
    });

    test("extended data", () => {
        const api = spawnAPI();
        const project = api.loadProject("/tsconfig.json");
        const sourceFile = project.getSourceFile("/src/index.ts");

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
        assert.equal(nodeCount, 7);
    });
});

test("Object equality", () => {
    const api = spawnAPI();
    const project = api.loadProject("/tsconfig.json");
    assert.strictEqual(project, api.loadProject("/tsconfig.json"));
    assert.strictEqual(
        project.getSymbolAtPosition("/src/index.ts", 9),
        project.getSymbolAtPosition("/src/index.ts", 10),
    );
});

function spawnAPI(files: Record<string, string> = defaultFiles) {
    return new API({
        cwd: new URL("../../../", import.meta.url).pathname,
        tsserverPath: new URL("../../../built/local/tsgo", import.meta.url).pathname,
        fs: createVirtualFileSystem(files),
    });
}
