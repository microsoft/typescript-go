import assert from "node:assert/strict";
import path from "node:path";
import test from "node:test";
import { getNodeModulesCandidates } from "../src/tsdkPackage";

test("uses the physical node_modules root for a scoped npm alias", () => {
    const nodeModules = path.resolve("workspace", "node_modules");
    const packagePath = path.join(nodeModules, "@typescript", "native");

    assert.deepEqual(getNodeModulesCandidates(packagePath, "typescript"), [
        nodeModules,
        path.join(nodeModules, "@typescript"),
    ]);
});

test("does not duplicate the metadata candidate for a regular package", () => {
    const nodeModules = path.resolve("workspace", "node_modules");
    const packagePath = path.join(nodeModules, "typescript");

    assert.deepEqual(getNodeModulesCandidates(packagePath, "typescript"), [nodeModules]);
});

test("uses the nearest node_modules directory for pnpm installations", () => {
    const nodeModules = path.resolve("workspace", "node_modules");
    const packagePath = path.join(nodeModules, ".pnpm", "typescript@7.1.0", "node_modules", "typescript");

    assert.deepEqual(getNodeModulesCandidates(packagePath, "typescript"), [
        path.join(nodeModules, ".pnpm", "typescript@7.1.0", "node_modules"),
    ]);
});
