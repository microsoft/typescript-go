import {
    version,
    versionMajorMinor,
} from "@typescript/native-preview";
import defaultExport from "@typescript/native-preview";
import assert from "node:assert";
import {
    describe,
    test,
} from "node:test";

describe("main entry point", () => {
    test("exposes version and versionMajorMinor as named exports", () => {
        assert.strictEqual(typeof version, "string");
        assert.strictEqual(typeof versionMajorMinor, "string");
        assert.match(version, /^\d+\.\d+\.\d+/);
        assert.strictEqual(versionMajorMinor, version.match(/^\d+\.\d+/)?.[0]);
    });

    test("exposes version and versionMajorMinor on the default export", () => {
        assert.strictEqual(defaultExport.version, version);
        assert.strictEqual(defaultExport.versionMajorMinor, versionMajorMinor);
    });
});
