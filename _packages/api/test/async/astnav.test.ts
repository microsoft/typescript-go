import { API } from "@typescript/api/async"; // @sync-skip
// @sync-only-start
// import { API } from "@typescript/api/sync";
// @sync-only-end
import { createVirtualFileSystem } from "@typescript/api/fs";
import {
    findNextToken,
    findPrecedingToken,
    formatSyntaxKind,
    getTokenAtPosition,
    getTouchingPropertyName,
} from "@typescript/ast";
import { SyntaxKind } from "@typescript/ast";
import type {
    Node,
    SourceFile,
} from "@typescript/ast";
import assert from "node:assert";
import { readFileSync } from "node:fs";
import { resolve } from "node:path";
import {
    after,
    before,
    describe,
    test,
} from "node:test";
import { fileURLToPath } from "node:url";

// ---------------------------------------------------------------------------
// Go JSON baseline format
// ---------------------------------------------------------------------------

interface TokenRun {
    startPos: number;
    endPos: number;
    kind: string;
    nodePos: number;
    nodeEnd: number;
}

/**
 * Expand run-length encoded Go baseline into a per-position map.
 */
function expandBaseline(runs: TokenRun[]): Map<number, { kind: string; pos: number; end: number; }> {
    const map = new Map<number, { kind: string; pos: number; end: number; }>();
    for (const run of runs) {
        const entry = { kind: run.kind, pos: run.nodePos, end: run.nodeEnd };
        for (let p = run.startPos; p <= run.endPos; p++) {
            map.set(p, entry);
        }
    }
    return map;
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

/** Convert an astnav result node to a string-based token info for comparison. */
function toTokenInfo(node: Node): { kind: string; pos: number; end: number; } {
    return {
        kind: formatSyntaxKind(node.kind),
        pos: node.pos,
        end: node.end,
    };
}

// ---------------------------------------------------------------------------
// Test configuration
// ---------------------------------------------------------------------------

const repoRoot = resolve(import.meta.dirname!, "..", "..", "..", "..");
const testFile = resolve(repoRoot, "_submodules/TypeScript/src/services/mapCode.ts");
const baselineDir = resolve(repoRoot, "testdata/baselines/reference/astnav");

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

describe("astnav", () => {
    let fileText: string;

    try {
        fileText = readFileSync(testFile, "utf-8");
    }
    catch {
        console.log("Skipping astnav tests: submodule not available");
        fileText = "";
    }

    if (!fileText) return;

    // Use the Go API to parse the file — the resulting SourceFile is already
    // in our SyntaxKind/NodeFlags enum space with correct JSDoc structure.
    let api: API;
    let sourceFile: SourceFile;

    before(async () => {
        api = new API({
            cwd: fileURLToPath(new URL("../../../../", import.meta.url).toString()),
            tsserverPath: fileURLToPath(new URL(`../../../../built/local/tsgo${process.platform === "win32" ? ".exe" : ""}`, import.meta.url).toString()),
            fs: createVirtualFileSystem({
                "/tsconfig.json": JSON.stringify({ files: ["/src/testFile.ts"] }),
                "/src/testFile.ts": fileText,
            }),
        });

        const snapshot = await api.updateSnapshot({ openProject: "/tsconfig.json" });
        const project = snapshot.getProject("/tsconfig.json")!;
        const sf = await project.program.getSourceFile("/src/testFile.ts");
        assert.ok(sf, "Failed to get source file from API");
        sourceFile = sf;
    });

    after(async () => {
        await api.close();
    });

    const testCases = [
        {
            name: "getTokenAtPosition",
            baselineFile: "GetTokenAtPosition.mapCode.ts.baseline.json",
            fn: getTokenAtPosition,
        },
        {
            name: "getTouchingPropertyName",
            baselineFile: "GetTouchingPropertyName.mapCode.ts.baseline.json",
            fn: getTouchingPropertyName,
        },
    ];

    for (const tc of testCases) {
        test(tc.name, () => {
            const baselinePath = resolve(baselineDir, tc.baselineFile);
            const runs: TokenRun[] = JSON.parse(readFileSync(baselinePath, "utf-8"));
            const expected = expandBaseline(runs);

            const failures: string[] = [];

            for (let pos = 0; pos < fileText.length; pos++) {
                const result = toTokenInfo(tc.fn(sourceFile, pos));
                const goExpected = expected.get(pos);

                if (!goExpected) continue;

                if (result.kind !== goExpected.kind || result.pos !== goExpected.pos || result.end !== goExpected.end) {
                    failures.push(
                        `  pos ${pos}: expected ${goExpected.kind} [${goExpected.pos}, ${goExpected.end}), ` +
                            `got ${result.kind} [${result.pos}, ${result.end})`,
                    );
                    if (failures.length >= 50) {
                        failures.push("  ... (truncated, too many failures)");
                        break;
                    }
                }
            }

            console.log(`  ${tc.name}: checked ${fileText.length} positions`);

            if (failures.length > 0) {
                assert.fail(
                    `${tc.name}: ${failures.length} position(s) differ from Go baseline:\n` +
                        failures.join("\n"),
                );
            }
        });
    }

    // ---------------------------------------------------------------------------
    // findNextToken tests
    // ---------------------------------------------------------------------------

    test("findNextToken: result always starts at the end of the previous token", () => {
        // For each unique token in the first 2000 characters, verify that
        // findNextToken returns a node whose pos equals the previous token's end.
        const limit = Math.min(fileText.length, 2000);
        const failures: string[] = [];
        const seen = new Set<number>();

        for (let pos = 0; pos < limit; pos++) {
            const token = getTokenAtPosition(sourceFile, pos);
            if (seen.has(token.pos)) continue;
            seen.add(token.pos);

            const result = findNextToken(token, sourceFile, sourceFile);
            if (result === undefined) continue;

            if (result.pos !== token.end) {
                failures.push(
                    `  token ${formatSyntaxKind(token.kind)} [${token.pos}, ${token.end}): ` +
                        `next token ${formatSyntaxKind(result.kind)} [${result.pos}, ${result.end}) ` +
                        `does not start at ${token.end}`,
                );
                if (failures.length >= 10) break;
            }
        }

        if (failures.length > 0) {
            assert.fail(`findNextToken: next token pos !== previous token end:\n${failures.join("\n")}`);
        }
    });

    test("findNextToken: returns undefined when parent is the token itself", () => {
        // When a token is its own parent, there is no next token within it.
        const token = getTokenAtPosition(sourceFile, 0);
        assert.notEqual(token.kind, SyntaxKind.EndOfFile);

        const result = findNextToken(token, token, sourceFile);
        assert.equal(result, undefined, `Expected undefined, got ${result !== undefined ? formatSyntaxKind(result.kind) : "undefined"}`);
    });

    test("findNextToken: returns undefined for EndOfFile token", () => {
        // The EndOfFile token has no successor.
        const eof = getTokenAtPosition(sourceFile, fileText.length);
        assert.equal(eof.kind, SyntaxKind.EndOfFile);

        const result = findNextToken(eof, sourceFile, sourceFile);
        assert.equal(result, undefined);
    });

    test("findNextToken: finds punctuation tokens not directly stored in the AST", () => {
        // mapCode.ts starts with `import {`.
        // The `{` after `import` is a syntactic token found via the scanner.
        const importToken = getTokenAtPosition(sourceFile, 0);
        assert.equal(importToken.kind, SyntaxKind.ImportKeyword, "expected 'import' at pos 0");

        const openBrace = findNextToken(importToken, sourceFile, sourceFile);
        assert.ok(openBrace !== undefined, "Expected a token after 'import'");
        assert.equal(openBrace.kind, SyntaxKind.OpenBraceToken, `Expected '{', got ${formatSyntaxKind(openBrace.kind)}`);
        assert.equal(openBrace.pos, importToken.end, `Expected next token pos === importToken.end (${importToken.end}), got ${openBrace.pos}`);
    });

    // ---------------------------------------------------------------------------
    // findPrecedingToken tests
    // ---------------------------------------------------------------------------

    test("findPrecedingToken: is consistent with findNextToken (roundtrip)", () => {
        // For each unique non-EOF, non-JSDoc token, the preceding token of (token.end) should be
        // the token itself (or one that ends at the same position).
        // JSDoc nodes are skipped: findPrecedingToken visits them differently from
        // getTokenAtPosition (a pre-existing limitation of the TypeScript port).
        const limit = Math.min(fileText.length, 2000);
        const failures: string[] = [];
        const seen = new Set<number>();

        for (let pos = 0; pos < limit; pos++) {
            const token = getTokenAtPosition(sourceFile, pos);
            if (token.kind === SyntaxKind.EndOfFile) continue;
            // Skip JSDoc nodes — findPrecedingToken does not yet visit JSDoc nodes.
            if (token.kind >= SyntaxKind.FirstJSDocNode && token.kind <= SyntaxKind.LastJSDocNode) continue;
            if (seen.has(token.pos)) continue;
            seen.add(token.pos);

            const preceding = findPrecedingToken(sourceFile, token.end);
            if (preceding === undefined) {
                failures.push(`  token ${formatSyntaxKind(token.kind)} [${token.pos}, ${token.end}): findPrecedingToken(${token.end}) returned undefined`);
            }
            else if (preceding.end !== token.end) {
                failures.push(
                    `  token ${formatSyntaxKind(token.kind)} [${token.pos}, ${token.end}): ` +
                        `findPrecedingToken returned ${formatSyntaxKind(preceding.kind)} [${preceding.pos}, ${preceding.end}) with different end`,
                );
            }
            if (failures.length >= 10) break;
        }

        if (failures.length > 0) {
            assert.fail(`findPrecedingToken roundtrip failures:\n${failures.join("\n")}`);
        }
    });

    test("findPrecedingToken: returns undefined at position 0", () => {
        // There is no token before the very start of the file.
        const result = findPrecedingToken(sourceFile, 0);
        assert.equal(result, undefined);
    });

    test("findPrecedingToken: returns a token at end of file", () => {
        // At the end of file there should be a valid preceding token.
        const result = findPrecedingToken(sourceFile, fileText.length);
        assert.ok(result !== undefined, "Expected a preceding token at end of file");
        assert.notEqual(result.kind, SyntaxKind.EndOfFile);
    });
});
