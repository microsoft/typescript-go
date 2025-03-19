import { createVirtualFileSystem } from "@typescript/api/base/fs";
import { API } from "@typescript/api/sync";
import assert from "node:assert";
import {
    describe,
    test,
} from "node:test";

describe("SourceFile", () => {
    test("file text", () => {
        const files = {
            "/tsconfig.json": "{}",
            "/src/index.ts": `import { foo } from './foo';`,
        };

        const api = new API({
            cwd: new URL("../../../", import.meta.url).pathname,
            tsserverPath: new URL("../../../built/local/tsgo", import.meta.url).pathname,
            fs: createVirtualFileSystem(files),
        });

        const project = api.loadProject("/tsconfig.json");
        const sourceFile = project.getSourceFile("/src/index.ts");

        assert.ok(sourceFile);
        assert.equal(sourceFile.text, files["/src/index.ts"]);
    });
});
