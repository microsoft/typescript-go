/**
 * Temporary script to compare generated factory functions against the original nodeFactory.ts.
 * Extracts all `create*` function signatures from both and reports differences.
 */

import { readFileSync } from "node:fs";
import { resolve, dirname } from "node:path";
import { fileURLToPath } from "node:url";

const __dirname = dirname(fileURLToPath(import.meta.url));

const generatedPath = resolve(__dirname, "../src/factory.ts");
const originalPath = resolve(__dirname, "../../../_submodules/TypeScript/src/compiler/factory/nodeFactory.ts");

const generated = readFileSync(generatedPath, "utf-8");
const original = readFileSync(originalPath, "utf-8");

// ---- Extract create functions from the generated factory ----

interface FuncInfo {
    name: string;
    params: string[];
    returnType: string;
    fullSig: string;
}

function extractGeneratedCreateFunctions(src: string): Map<string, FuncInfo> {
    const map = new Map<string, FuncInfo>();
    // Pattern: export function createXxx(params): ReturnType {
    const re = /^export function (create\w+)\(([^)]*)\):\s*([^\s{]+)\s*\{/gm;
    let m;
    while ((m = re.exec(src)) !== null) {
        const name = m[1];
        const rawParams = m[2];
        const returnType = m[3];
        const params = rawParams.split(",").map(p => p.trim()).filter(Boolean);
        map.set(name, { name, params, returnType, fullSig: m[0] });
    }
    return map;
}

// ---- Extract create functions from the original nodeFactory.ts ----
// The original defines them as `function createXxx(...)` inside createNodeFactory().
// They're also on the factory object. Let's find the function declarations.

function extractOriginalCreateFunctions(src: string): Map<string, FuncInfo> {
    const map = new Map<string, FuncInfo>();
    // Match `function createXxx(params): ReturnType {` or `function createXxx(params) {`
    // Also multiline params - we'll grab until closing paren
    const lines = src.split("\n");
    for (let i = 0; i < lines.length; i++) {
        const line = lines[i];
        // Look for `function create` at start of line (possibly indented)
        const funcMatch = line.match(/^\s*function (create\w+)\(/);
        if (!funcMatch) continue;
        const name = funcMatch[1];

        // Skip internal worker functions
        if (name.includes("Worker") || name.includes("worker")) continue;

        // Collect full signature (may span multiple lines)
        let sig = line;
        let depth = 0;
        for (const ch of line) {
            if (ch === "(") depth++;
            if (ch === ")") depth--;
        }
        let j = i;
        while (depth > 0 && j < lines.length - 1) {
            j++;
            sig += " " + lines[j].trim();
            for (const ch of lines[j]) {
                if (ch === "(") depth++;
                if (ch === ")") depth--;
            }
        }

        // Extract params and return type from the signature
        const sigMatch = sig.match(/function \w+\(([^)]*)\)(?::\s*(\S+))?\s*\{/);
        if (sigMatch) {
            const rawParams = sigMatch[1];
            const returnType = sigMatch[2] || "unknown";
            const params = rawParams.split(",").map(p => p.trim()).filter(Boolean);
            // Only keep the first definition if duplicated
            if (!map.has(name)) {
                map.set(name, { name, params, returnType, fullSig: sig.trim() });
            }
        }
    }
    return map;
}

const genFuncs = extractGeneratedCreateFunctions(generated);
const origFuncs = extractOriginalCreateFunctions(original);

console.log(`Generated factory: ${genFuncs.size} create* functions`);
console.log(`Original factory: ${origFuncs.size} create* functions`);
console.log();

// ---- Compare ----

// 1. Functions in original but NOT in generated
const missingFromGenerated: string[] = [];
for (const [name] of origFuncs) {
    if (!genFuncs.has(name)) {
        missingFromGenerated.push(name);
    }
}

// 2. Functions in generated but NOT in original
const extraInGenerated: string[] = [];
for (const [name] of genFuncs) {
    if (!origFuncs.has(name)) {
        extraInGenerated.push(name);
    }
}

// 3. Functions in both but with different param counts
interface ParamDiff {
    name: string;
    genParams: string[];
    origParams: string[];
    genReturn: string;
    origReturn: string;
}
const paramDiffs: ParamDiff[] = [];
for (const [name, genInfo] of genFuncs) {
    const origInfo = origFuncs.get(name);
    if (!origInfo) continue;

    // Compare param names (strip types for simpler comparison)
    const genNames = genInfo.params.map(p => p.replace(/[?:].*$/, "").trim());
    const origNames = origInfo.params.map(p => p.replace(/[?:].*$/, "").trim());

    const different = genNames.length !== origNames.length ||
        genNames.some((n, i) => n !== origNames[i]);

    if (different) {
        paramDiffs.push({
            name,
            genParams: genInfo.params,
            origParams: origInfo.params,
            genReturn: genInfo.returnType,
            origReturn: origInfo.returnType,
        });
    }
}

// ---- Report ----

if (missingFromGenerated.length) {
    console.log(`=== MISSING from generated (${missingFromGenerated.length}) ===`);
    for (const name of missingFromGenerated.sort()) {
        const orig = origFuncs.get(name)!;
        console.log(`  ${name}(${orig.params.length} params) -> ${orig.returnType}`);
    }
    console.log();
}

if (extraInGenerated.length) {
    console.log(`=== EXTRA in generated (not in original) (${extraInGenerated.length}) ===`);
    for (const name of extraInGenerated.sort()) {
        console.log(`  ${name}`);
    }
    console.log();
}

if (paramDiffs.length) {
    console.log(`=== PARAMETER DIFFERENCES (${paramDiffs.length}) ===`);
    for (const diff of paramDiffs.sort((a, b) => a.name.localeCompare(b.name))) {
        console.log(`  ${diff.name}:`);
        console.log(`    gen:  (${diff.genParams.join(", ")})`);
        console.log(`    orig: (${diff.origParams.join(", ")})`);
    }
    console.log();
}

if (!missingFromGenerated.length && !extraInGenerated.length && !paramDiffs.length) {
    console.log("No differences found!");
}
