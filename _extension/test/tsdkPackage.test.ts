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

function linkPackage(root: string, relativePath: string, targetPackageJson: string): void {
    const packagePath = path.join(root, relativePath);
    fs.mkdirSync(path.dirname(packagePath), { recursive: true });
    fs.symlinkSync(
        path.dirname(targetPackageJson),
        packagePath,
        process.platform === "win32" ? "junction" : "dir",
    );
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
    const platformPackageJson = createPackage(
        root,
        `node_modules/.pnpm/@typescript+${platformPackage}@7.0.2/node_modules/@typescript/${platformPackage}`,
    );
    linkPackage(root, `${virtualStore}/@typescript/${platformPackage}`, platformPackageJson);

    assert.equal(
        resolvePackageExecutable(packageJsonPath, platformPackage, exeName),
        path.join(path.dirname(platformPackageJson), "lib", exeName),
    );
});

test("resolves a pnpm native-preview virtual-store package", t => {
    const root = createFixture(t);
    const version = "7.0.0-dev.20260707.2";
    const nativePlatformPackage = "native-preview-linux-x64";
    const virtualStore = `node_modules/.pnpm/@typescript+native-preview@${version}/node_modules`;
    const packageJsonPath = createPackage(root, `${virtualStore}/@typescript/native-preview`);
    const platformPackageJson = createPackage(
        root,
        `node_modules/.pnpm/@typescript+${nativePlatformPackage}@${version}/node_modules/@typescript/${nativePlatformPackage}`,
    );
    linkPackage(root, `${virtualStore}/@typescript/${nativePlatformPackage}`, platformPackageJson);

    assert.equal(
        resolvePackageExecutable(packageJsonPath, nativePlatformPackage, "tsgo"),
        path.join(path.dirname(platformPackageJson), "lib", "tsgo"),
    );
});

test("resolves an npm linked-store package", t => {
    const root = createFixture(t);
    const store = "node_modules/.store/typescript@7.0.2-hash/node_modules";
    const packageJsonPath = createPackage(root, `${store}/typescript`);
    const platformPackageJson = createPackage(
        root,
        `node_modules/.store/@typescript+${platformPackage}@7.0.2-hash/node_modules/@typescript/${platformPackage}`,
    );
    linkPackage(root, `${store}/@typescript/${platformPackage}`, platformPackageJson);

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
