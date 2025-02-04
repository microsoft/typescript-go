import assert from "node:assert";
import type {
    MetaModelSchema,
    OrType,
    Type,
} from "./metaModelSchema.mts";

import fs from "node:fs";
import path from "node:path";
import url from "node:url";

const __filename = url.fileURLToPath(new URL(import.meta.url));
const __dirname = path.dirname(__filename);

const out = process.argv[2];
if (!out) {
    console.error("Usage: node generate.mts <output file>");
    process.exit(1);
}

const metaModelPath = path.resolve(__dirname, "metaModel.json");

const model: MetaModelSchema = JSON.parse(fs.readFileSync(metaModelPath, "utf-8"));

let parts: string[] = [];
let indentLevel = 0;

function indent() {
    indentLevel++;
}

function dedent() {
    indentLevel--;
}

function write(s: string) {
    parts.push(s);
}

function writeLine(s: string) {
    startLine(s);
    write("\n");
}

function startLine(s: string) {
    parts.push("\t".repeat(indentLevel));
    write(s);
}

function finishLine(s: string) {
    write(s);
    write("\n");
}

function writeDocumentation(doc: string | undefined) {
    if (doc) {
        const lines = doc.split("\n");
        for (const line of lines) {
            startLine("// ");
            finishLine(line);
        }
    }
}

function writeDeprecation(deprecated: string | undefined) {
    if (deprecated) {
        writeLine("//");
        startLine("// Deprecated: ");
        finishLine(deprecated);
    }
}

function titleCase(s: string): string {
    return s.charAt(0).toUpperCase() + s.slice(1);
}

const unionTypes = new Map<string, Type[]>();

function writeOr(t: OrType, wasOptional = false) {
    let nullable = false;
    const types = t.items.filter(item => {
        if (item.kind === "base" && item.name === "null") {
            nullable = true;
            return false;
        }
        return true;
    });
    if (nullable) {
        if (wasOptional) {
            write("Nullable[");
        }
        else {
            write("*");
        }
    }
    if (types.length === 1) {
        writeTypeElement(types[0]);
    }
    else {
        const names = [];
        for (const t of types) {
            if (t.kind === "reference") {
                names.push(t.name);
            }
            else if (t.kind === "base") {
                names.push(t.name);
            }
            else if (t.kind === "array" && t.element.kind === "reference") {
                names.push("ArrayOf" + titleCase(t.element.name));
            }
            else {
                write("TODO_or_" + t.kind);
            }
        }

        const name = names.map(titleCase).join("Or");
        unionTypes.set(name, types);
        write(name);
    }
    if (nullable && wasOptional) {
        write("]");
    }
}

function compareTypes(a: Type, b: Type): number {
    if (a.kind === "base" && b.kind === "base") {
        return a.name.localeCompare(b.name);
    }
    if (a.kind === "reference" && b.kind === "reference") {
        return a.name.localeCompare(b.name);
    }
    if (a.kind === "array" && b.kind === "array") {
        return compareTypes(a.element, b.element);
    }
    if (a.kind === "map" && b.kind === "map") {
        return compareTypes(a.key, b.key) || compareTypes(a.value, b.value);
    }
    if (a.kind === "or" && b.kind === "or") {
        const cmp = a.items.length - b.items.length;
        if (cmp !== 0) {
            return cmp;
        }
        const aItems = a.items.slice().sort(compareTypes);
        const bItems = b.items.slice().sort(compareTypes);

        for (let i = 0; i < aItems.length; i++) {
            const cmp = compareTypes(aItems[i], bItems[i]);
            if (cmp !== 0) {
                return cmp;
            }
        }

        return 0;
    }
    if (a.kind === "tuple" && b.kind === "tuple") {
        const cmp = a.items.length - b.items.length;
        if (cmp !== 0) {
            return cmp;
        }
        for (let i = 0; i < a.items.length; i++) {
            const cmp = compareTypes(a.items[i], b.items[i]);
            if (cmp !== 0) {
                return cmp;
            }
        }
        return 0;
    }
    if (a.kind === "literal" && b.kind === "literal") {
        // For now, the spec only uses this for empty arrays
        assert(a.value.properties.length === 0);
        assert(b.value.properties.length === 0);
        return 0;
    }
    if (a.kind === "stringLiteral" && b.kind === "stringLiteral") {
        return a.value.localeCompare(b.value);
    }
    if (a.kind === "integerLiteral" && b.kind === "integerLiteral") {
        return a.value - b.value;
    }
    if (a.kind === "booleanLiteral" && b.kind === "booleanLiteral") {
        return a.value === b.value ? 0 : a.value ? 1 : -1;
    }
    return a.kind.localeCompare(b.kind);
}

function writeTypeElement(t: Type, wasOptional = false) {
    switch (t.kind) {
        case "reference":
            write(t.name);
            break;
        case "base":
            switch (t.name) {
                case "integer":
                    write("int32");
                    break;
                case "uinteger":
                    write("uint32");
                    break;
                case "string":
                    write("string");
                    break;
                case "boolean":
                    write("bool");
                    break;
                case "URI":
                    write("URI");
                    break;
                case "DocumentUri":
                    write("DocumentUri");
                    break;
                case "decimal":
                    write("float64");
                    break;
                default:
                    write("TODO_base_" + t.name);
                    break;
            }
            break;
        case "array":
            write("[]");
            writeTypeElement(t.element);
            break;
        case "stringLiteral":
            write("string");
            break;
        case "map":
            write("map[");
            write(t.key.name);
            write("]");

            const vt = t.value;
            switch (vt.kind) {
                case "reference":
                    write(vt.name);
                    break;
                case "array":
                    write("[]");
                    writeTypeElement(vt.element);
                    break;
                case "or":
                    writeOr(vt);
                    break;
                default:
                    write("TODO_map_value_" + vt.kind);
                    break;
            }
            break;
        case "or":
            writeOr(t, wasOptional);
            break;
        default:
            write("TODO_" + t.kind);
            break;
    }
}

// Generation

writeLine("// Code generated by genlsp; DO NOT EDIT.");
writeLine("");
writeLine("package lsproto2");
writeLine("");
writeLine(`import "encoding/json"`);
writeLine("");
writeLine("// Meta model version " + model.metaData.version);
writeLine("");

writeLine("type URI string\n");
writeLine("type DocumentUri string\n");
writeLine("type Method string\n");

writeLine("");

writeLine("type Nullable[T any] struct {");
indent();
writeLine("Value T");
writeLine("Null bool");
dedent();
writeLine("}");
writeLine("");

writeLine("func (n Nullable[T]) MarshalJSON() ([]byte, error) {");
indent();
writeLine("if n.Null {");
indent();
writeLine("return []byte(`null`), nil");
dedent();
writeLine("}");
writeLine("return json.Marshal(n.Value)");
dedent();
writeLine("}");
writeLine("");

writeLine("func (n *Nullable[T]) UnmarshalJSON(data []byte) error {");
indent();
writeLine("*n = Nullable[T]{}");
writeLine("if string(data) == `null` {");
indent();
writeLine("n.Null = true");
writeLine("return nil");
dedent();
writeLine("}");
writeLine("return json.Unmarshal(data, &n.Value)");
dedent();
writeLine("}");
writeLine("");

for (const t of model.structures) {
    writeDocumentation(t.documentation);
    writeDeprecation(t.deprecated);

    writeLine("type " + t.name + " struct {");
    indent();

    for (const e of t.extends ?? []) {
        if (e.kind !== "reference") {
            throw new Error("Unexpected extends kind: " + e.kind);
        }
        writeLine(e.name);
    }
    for (const m of t.mixins ?? []) {
        if (m.kind !== "reference") {
            throw new Error("Unexpected mixin kind: " + m.kind);
        }
        writeLine(m.name);
    }

    for (const p of t.properties) {
        writeDocumentation(p.documentation);
        writeDeprecation(p.deprecated);

        startLine(titleCase(p.name) + " ");

        if (p.optional) {
            write("*");
        }

        writeTypeElement(p.type, !!p.optional);

        finishLine(' `json:"' + p.name + '"`');
    }

    dedent();
    writeLine("}");
    writeLine("\n");
}

for (const t of model.enumerations) {
    writeDocumentation(t.documentation);
    writeDeprecation(t.deprecated);

    writeLine("type " + t.name + " int");
    writeLine("\n");
}

for (const t of model.typeAliases) {
    writeDocumentation(t.documentation);
    writeDeprecation(t.deprecated);

    if (t.name === "LSPAny") {
        writeLine("type LSPAny = any\n");
        continue;
    }

    startLine("type " + t.name + " = ");
    writeTypeElement(t.type);
    writeLine("\n");
}

function methodNameToIdentifier(method: string): string {
    return method.split("/").map(v => v === "$" ? "" : titleCase(v)).join("");
}

for (const t of model.requests) {
    writeDocumentation(t.documentation);
    writeDeprecation(t.deprecated);
    writeLine("const MethodRequest" + methodNameToIdentifier(t.method) + ' Method = "' + t.method + '"\n');
}

for (const t of model.notifications) {
    writeDocumentation(t.documentation);
    writeDeprecation(t.deprecated);
    writeLine("const MethodNotification" + methodNameToIdentifier(t.method) + ' Method = "' + t.method + '"\n');
}

fs.writeFileSync(out, parts.join(""));
