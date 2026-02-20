#!/usr/bin/env node

/**
 * Code generator for the monomorphic NodeObject class and factory functions.
 *
 * Usage:
 *   node _packages/ast/scripts/generateFactory.ts
 *
 * Reads:  _packages/ast/src/nodes.ts
 * Writes: _packages/ast/src/factory.ts
 */

import { execaSync } from "execa";
import * as fs from "node:fs";
import * as path from "node:path";
import { fileURLToPath } from "node:url";
import ts from "typescript";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const nodesPath = path.resolve(__dirname, "../src/nodes.ts");
const outputPath = path.resolve(__dirname, "../src/factory.ts");

const errors: string[] = [];

function reportError(msg: string): void {
    errors.push(msg);
}

function fail(msg: string): never {
    throw new Error(msg);
}

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

interface PropertyInfo {
    name: string;
    type: string;
    optional: boolean;
}

interface ExtendsInfo {
    name: string;
    typeArguments: string[];
}

interface InterfaceInfo {
    name: string;
    syntaxKind: string;
    properties: PropertyInfo[];
    extends: ExtendsInfo[];
    typeParameters: string[];
}

interface FactoryDef {
    interfaceName: string;
    syntaxKind: string;
    factoryName: string;
    params: PropertyInfo[];
}

// ---------------------------------------------------------------------------
// Step 1: Parse nodes.ts and build interface map
// ---------------------------------------------------------------------------

const nodesSource = fs.readFileSync(nodesPath, "utf-8");
const sourceFile = ts.createSourceFile("nodes.ts", nodesSource, ts.ScriptTarget.Latest, true);

const interfaces = new Map<string, InterfaceInfo>();

for (const stmt of sourceFile.statements) {
    if (!ts.isInterfaceDeclaration(stmt)) continue;

    const name = stmt.name.text;
    let syntaxKind = "";
    const properties: PropertyInfo[] = [];
    const extendsInfos: ExtendsInfo[] = [];
    const typeParameters: string[] = [];

    if (stmt.typeParameters) {
        for (const tp of stmt.typeParameters) {
            typeParameters.push(tp.name.text);
        }
    }

    if (stmt.heritageClauses) {
        for (const clause of stmt.heritageClauses) {
            if (clause.token === ts.SyntaxKind.ExtendsKeyword) {
                for (const type of clause.types) {
                    if (!ts.isIdentifier(type.expression)) {
                        fail(`${name}: extends clause has non-identifier expression: ${type.expression.getText(sourceFile)}`);
                    }
                    const typeArgs = type.typeArguments
                        ? type.typeArguments.map(a => a.getText(sourceFile))
                        : [];
                    extendsInfos.push({ name: type.expression.text, typeArguments: typeArgs });
                }
            }
        }
    }

    for (const member of stmt.members) {
        // Skip index signatures, method signatures, call signatures, etc.
        if (!ts.isPropertySignature(member)) continue;

        if (!member.name) {
            fail(`${name}: property signature has no name at pos ${member.pos}`);
        }

        let propName: string;
        if (ts.isIdentifier(member.name)) {
            propName = member.name.text;
        }
        else if (ts.isStringLiteral(member.name)) {
            propName = member.name.text;
        }
        else {
            fail(`${name}: unexpected property name kind ${ts.SyntaxKind[member.name.kind]} at pos ${member.name.pos}`);
        }

        if (!member.type) {
            fail(`${name}.${propName}: property has no type annotation`);
        }
        const propType = member.type.getText(sourceFile);
        const isOptional = !!member.questionToken;

        if (propName === "kind" && propType.startsWith("SyntaxKind.")) {
            syntaxKind = propType;
        }

        properties.push({ name: propName, type: propType, optional: isOptional });
    }

    interfaces.set(name, { name, syntaxKind, properties, extends: extendsInfos, typeParameters });
}

// ---------------------------------------------------------------------------
// Step 2: Collect exported type names for import filtering
// ---------------------------------------------------------------------------

const exportedTypeNames = new Set<string>();
for (const stmt of sourceFile.statements) {
    if (ts.isTypeAliasDeclaration(stmt)) {
        exportedTypeNames.add(stmt.name.text);
    }
    else if (ts.isInterfaceDeclaration(stmt)) {
        exportedTypeNames.add(stmt.name.text);
    }
}

// ---------------------------------------------------------------------------
// Step 3: Resolve all properties for an interface (including inherited)
// ---------------------------------------------------------------------------

const EXCLUDED_PROPS = new Set(["kind", "parent", "pos", "end"]);

function isBrandField(name: string): boolean {
    return name.startsWith("_") && (name.endsWith("Brand") || name.endsWith("brand"));
}

function substituteTypeParams(type: string, substitutions: Map<string, string>): string {
    if (substitutions.size === 0) return type;
    // Replace standalone type parameter references with their concrete types
    let result = type;
    for (const [param, concrete] of substitutions) {
        result = result.replace(new RegExp(`\\b${param}\\b`, "g"), concrete);
    }
    return result;
}

// Known external base interfaces that we don't need to resolve properties from.
const EXTERNAL_BASES = new Set([
    "ReadonlyArray",
    "ReadonlyTextRange",
]);

function getAllProperties(name: string, visited = new Set<string>(), substitutions = new Map<string, string>()): PropertyInfo[] {
    if (visited.has(name)) return [];
    visited.add(name);

    const iface = interfaces.get(name);
    if (!iface) {
        if (!EXTERNAL_BASES.has(name)) {
            reportError(`getAllProperties: interface "${name}" not found in nodes.ts (not in EXTERNAL_BASES either)`);
        }
        return [];
    }

    const result: PropertyInfo[] = [];

    for (const ext of iface.extends) {
        // Build substitutions for the parent's type parameters
        const parentIface = interfaces.get(ext.name);
        const parentSubs = new Map(substitutions);
        if (parentIface && ext.typeArguments.length > 0) {
            for (let i = 0; i < Math.min(parentIface.typeParameters.length, ext.typeArguments.length); i++) {
                // Also apply current substitutions to the type argument itself
                parentSubs.set(parentIface.typeParameters[i], substituteTypeParams(ext.typeArguments[i], substitutions));
            }
        }
        result.push(...getAllProperties(ext.name, visited, parentSubs));
    }

    for (const prop of iface.properties) {
        if (substitutions.size > 0) {
            result.push({ ...prop, type: substituteTypeParams(prop.type, substitutions) });
        }
        else {
            result.push(prop);
        }
    }

    return result;
}

// ---------------------------------------------------------------------------
// Step 4: Build factory definitions for concrete interfaces
// ---------------------------------------------------------------------------

// Rename reserved words used as parameter names
const RESERVED_WORDS = new Set(["arguments", "class", "default", "delete", "export", "extends", "import", "in", "new", "return", "super", "switch", "this", "throw", "typeof", "var", "void", "with", "yield"]);
function safeParamName(name: string): string {
    return RESERVED_WORDS.has(name) ? `${name}_` : name;
}

const factoryDefs: FactoryDef[] = [];
const allPropertyNames = new Set<string>();

for (const [name, iface] of interfaces) {
    if (!iface.syntaxKind) continue;
    // Skip union kinds (generic Token types)
    if (iface.syntaxKind.includes(" | ")) continue;

    const allProps = getAllProperties(name);

    // Deduplicate by name, last definition wins
    const propMap = new Map<string, PropertyInfo>();
    for (const prop of allProps) {
        if (EXCLUDED_PROPS.has(prop.name) || isBrandField(prop.name)) continue;
        propMap.set(prop.name, prop);
    }

    const params = [...propMap.values()];

    // Validate no unresolved type parameters leaked through
    for (const p of params) {
        // Single uppercase letter like T, K, U indicates an unresolved generic
        if (/\b[A-Z]\b/.test(p.type)) {
            // Check it's not just a normal single-letter type used in the codebase
            const singleLetters = p.type.match(/\b[A-Z]\b/g) ?? [];
            for (const letter of singleLetters) {
                if (!exportedTypeNames.has(letter)) {
                    fail(`${name}.${p.name}: unresolved type parameter "${letter}" in type "${p.type}"`);
                }
            }
        }
    }

    const factoryName = `create${name}`;

    for (const prop of params) {
        allPropertyNames.add(prop.name);
    }

    factoryDefs.push({
        interfaceName: name,
        syntaxKind: iface.syntaxKind,
        factoryName,
        params,
    });
}

factoryDefs.sort((a, b) => a.interfaceName.localeCompare(b.interfaceName));

// ---------------------------------------------------------------------------
// Step 5: Collect type references for imports
// ---------------------------------------------------------------------------

function extractTypeReferences(typeStr: string): string[] {
    // Remove SyntaxKind.XYZ and TokenFlags.XYZ references before extracting type names
    const cleaned = typeStr.replace(/\b(?:SyntaxKind|TokenFlags)\.\w+/g, "");
    const matches = cleaned.match(/\b[A-Z][A-Za-z0-9]*\b/g) ?? [];
    const builtins = new Set([
        "Array",
        "ReadonlyArray",
        "Record",
        "Map",
        "Set",
        "Promise",
        "Partial",
        "Required",
        "Readonly",
        "Pick",
        "Omit",
        "Exclude",
        "Extract",
        "NonNullable",
        "ReturnType",
        "InstanceType",
    ]);
    return matches.filter(m => !builtins.has(m));
}

const referencedTypes = new Set<string>(["Node"]);
for (const def of factoryDefs) {
    referencedTypes.add(def.interfaceName);
    for (const param of def.params) {
        for (const t of extractTypeReferences(param.type)) {
            referencedTypes.add(t);
        }
    }
}
referencedTypes.delete("SyntaxKind");
const needsTokenFlags = referencedTypes.has("TokenFlags");
referencedTypes.delete("TokenFlags");

// Validate all referenced types are accounted for
const KNOWN_EXTERNAL_TYPES = new Set([
    "SyntaxKind",
    "TokenFlags",
]);
for (const t of referencedTypes) {
    if (!exportedTypeNames.has(t) && !KNOWN_EXTERNAL_TYPES.has(t)) {
        reportError(`Referenced type "${t}" is not exported from nodes.ts and is not a known external type`);
    }
}

const importTypes = [...referencedTypes].filter(t => exportedTypeNames.has(t)).sort();

// ---------------------------------------------------------------------------
// Step 6: Emit output
// ---------------------------------------------------------------------------

const lines: string[] = [];

function emit(line: string) {
    lines.push(line);
}

// Header
emit("// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!");
emit("// !!! THIS FILE IS AUTO-GENERATED - DO NOT EDIT !!!");
emit("// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!");
emit("//");
emit("// Source: _packages/ast/src/nodes.ts");
emit("// Generator: _packages/ast/scripts/generateFactory.ts");
emit("//");
emit("");
emit(`import { SyntaxKind } from "#syntaxKind";`);

if (needsTokenFlags) {
    emit(`import { TokenFlags } from "#tokenFlags";`);
}

emit(`import type {`);
for (const t of importTypes) {
    emit(`    ${t},`);
}
emit(`} from "./nodes.ts";`);
emit("");

const sortedPropertyNames = [...allPropertyNames].sort();

// NodeObject class
emit("/**");
emit(" * Monomorphic AST node implementation.");
emit(" * All synthetic nodes share the same V8 hidden class for optimal property access.");
emit(" *");
emit(" * Common fields live directly on the object; kind-specific fields are stored");
emit(" * in the `_data` bag and accessed via generated property accessors.");
emit(" */");
emit("export class NodeObject {");
emit("    readonly kind!: SyntaxKind;");
emit("    readonly pos!: number;");
emit("    readonly end!: number;");
emit("    readonly parent!: Node;");
emit("    /** @internal */");
emit("    _data: any;");
emit("");
emit("    constructor(kind: SyntaxKind, data: any) {");
emit("        this.kind = kind;");
emit("        this.pos = -1;");
emit("        this.end = -1;");
emit("        this.parent = undefined!;");
emit("        this._data = data;");
emit("    }");
emit("");

for (const propName of sortedPropertyNames) {
    emit(`    get ${propName}(): any { return this._data?.${propName}; }`);
}

emit("}");
emit("");

// createToken helper
emit("/**");
emit(" * Create a simple token node with only a `kind`.");
emit(" */");
emit("export function createToken<TKind extends SyntaxKind>(kind: TKind): Node & { readonly kind: TKind } {");
emit("    return new NodeObject(kind, undefined) as any;");
emit("}");
emit("");

// Factory functions
for (const def of factoryDefs) {
    const { interfaceName, syntaxKind, factoryName, params } = def;

    const requiredParams = params.filter(p => !p.optional);
    const optionalParams = params.filter(p => p.optional);
    const orderedParams = [...requiredParams, ...optionalParams];

    const paramList = orderedParams.map(p => {
        const opt = p.optional ? "?" : "";
        return `${safeParamName(p.name)}${opt}: ${p.type}`;
    });

    emit(`export function ${factoryName}(${paramList.join(", ")}): ${interfaceName} {`);

    if (params.length === 0) {
        emit(`    return new NodeObject(${syntaxKind}, undefined) as unknown as ${interfaceName};`);
    }
    else {
        emit(`    return new NodeObject(${syntaxKind}, {`);
        for (const p of orderedParams) {
            const safe = safeParamName(p.name);
            if (safe !== p.name) {
                emit(`        ${p.name}: ${safe},`);
            }
            else {
                emit(`        ${p.name},`);
            }
        }
        emit(`    }) as unknown as ${interfaceName};`);
    }
    emit("}");
    emit("");
}

// ---------------------------------------------------------------------------
// Step 7: Write output
// ---------------------------------------------------------------------------

const output = lines.join("\n") + "\n";
fs.writeFileSync(outputPath, output);

console.log("Formatting...");
execaSync("dprint", ["fmt", outputPath]);

console.log(`Generated ${outputPath}`);
console.log(`  ${factoryDefs.length} factory functions`);
console.log(`  ${sortedPropertyNames.length} accessors on NodeObject`);

if (factoryDefs.length === 0) {
    reportError("No factory definitions generated — something is very wrong");
}

if (sortedPropertyNames.length === 0) {
    reportError("No properties collected — something is very wrong");
}

if (errors.length > 0) {
    console.error(`\n${errors.length} error(s):`);
    for (const e of errors) {
        console.error(`  - ${e}`);
    }
    process.exit(1);
}
