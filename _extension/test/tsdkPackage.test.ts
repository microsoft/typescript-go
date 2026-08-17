import assert from "node:assert/strict";
import fs from "node:fs";
import os from "node:os";
import path from "node:path";
import test from "node:test";
import { resolvePackageExecutable } from "../src/tsdkPackage";

const platformPackage = "typescript-linux-x64";
const exeName = "tsc";

function createPackage(root: string, relativePath: string): string {
    const packagePath = path.join(root, relativePath);
    fs.mkdirSync(packagePath, { recursive: true });
    const packageJsonPath = path.join(packagePath, "package.json");
    fs.writeFileSync(packageJsonPath, "{}");
    return packageJsonPath;
}

function createFixture(t: test.TestContext): string {
    const root = fs.mkdtempSync(path.join(os.tmpdir(), "tsdk-package-"));
    t.after(() => fs.rmSync(root, { recursive: true, force: true }));
    return root;
}

test("resolves the TypeScript 7 package", t => {
    const root = createFixture(t);
    const packageJsonPath = createPackage(root, "node_modules/typescript");
    const platformPackageJson = createPackage(root, `node_modules/@typescript/${platformPackage}`);

    assert.equal(
        resolvePackageExecutable(packageJsonPath, platformPackage, exeName),
        path.join(path.dirname(platformPackageJson), "lib", exeName),
    );
});

test("resolves the TypeScript 7 scoped alias recommended in the 7.0 release post", t => {
    const root = createFixture(t);
    const packageJsonPath = createPackage(root, "node_modules/@typescript/native");
    const platformPackageJson = createPackage(root, `node_modules/@typescript/${platformPackage}`);

    assert.equal(
        resolvePackageExecutable(packageJsonPath, platformPackage, exeName),
        path.join(path.dirname(platformPackageJson), "lib", exeName),
    );
});

test("resolves an unscoped TypeScript 7 alias", t => {
    const root = createFixture(t);
    const packageJsonPath = createPackage(root, "node_modules/typescript-next");
    const platformPackageJson = createPackage(root, `node_modules/@typescript/${platformPackage}`);

    assert.equal(
        resolvePackageExecutable(packageJsonPath, platformPackage, exeName),
        path.join(path.dirname(platformPackageJson), "lib", exeName),
    );
});

test("resolves the native-preview package", t => {
    const root = createFixture(t);
    const packageJsonPath = createPackage(root, "node_modules/@typescript/native-preview");
    const nativePlatformPackage = "native-preview-linux-x64";
    const platformPackageJson = createPackage(root, `node_modules/@typescript/${nativePlatformPackage}`);

    assert.equal(
        resolvePackageExecutable(packageJsonPath, nativePlatformPackage, "tsgo"),
        path.join(path.dirname(platformPackageJson), "lib", "tsgo"),
    );
});

test("resolves a non-hoisted platform package", t => {
    const root = createFixture(t);
    const packageJsonPath = createPackage(root, "node_modules/@typescript/native");
    const platformPackageJson = createPackage(root, `node_modules/@typescript/native/node_modules/@typescript/${platformPackage}`);

    assert.equal(
        resolvePackageExecutable(packageJsonPath, platformPackage, exeName),
        path.join(path.dirname(platformPackageJson), "lib", exeName),
    );
});

test("resolves a platform package hoisted above a workspace", t => {
    const root = createFixture(t);
    const packageJsonPath = createPackage(root, "packages/app/node_modules/@typescript/native");
    const platformPackageJson = createPackage(root, `node_modules/@typescript/${platformPackage}`);

    assert.equal(
        resolvePackageExecutable(packageJsonPath, platformPackage, exeName),
        path.join(path.dirname(platformPackageJson), "lib", exeName),
    );
});

test("resolves a pnpm virtual-store package", t => {
    const root = createFixture(t);
    const virtualStore = "node_modules/.pnpm/typescript@7.0.2/node_modules";
    const packageJsonPath = createPackage(root, `${virtualStore}/typescript`);
    const platformPackageJson = createPackage(root, `${virtualStore}/@typescript/${platformPackage}`);

    assert.equal(
        resolvePackageExecutable(packageJsonPath, platformPackage, exeName),
        path.join(path.dirname(platformPackageJson), "lib", exeName),
    );
});

test("throws when the platform package is missing", t => {
    const root = createFixture(t);
    const packageJsonPath = createPackage(root, "node_modules/@typescript/native");

    assert.throws(
        () => resolvePackageExecutable(packageJsonPath, platformPackage, exeName),
        { code: "MODULE_NOT_FOUND" },
    );
});
