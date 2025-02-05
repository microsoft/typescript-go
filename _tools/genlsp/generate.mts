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

function compareValues(a: string | number, b: string | number): number {
    if (typeof a === "string" && typeof b === "string") {
        return a < b ? -1 : a > b ? 1 : 0;
    }
    if (typeof a === "number" && typeof b === "number") {
        return a - b;
    }
    throw new Error("Cannot compare values of different types");
}

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

interface UnionMember {
    type: Type;
    name: string;
}

const unionTypes = new Map<string, UnionMember[]>();

function writeOr(t: OrType, wasOptional = false): boolean {
    let nullable = false;
    let omitEmpty = true;
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
            omitEmpty = false;
        }
    }
    if (types.length === 1) {
        writeTypeElement(types[0]);
    }
    else {
        const members: UnionMember[] = [];
        for (const t of types) {
            let name;
            if (t.kind === "reference") {
                name = t.name;
            }
            else if (t.kind === "base") {
                name = t.name;
            }
            else if (t.kind === "array" && (t.element.kind === "reference" || t.element.kind === "base")) {
                name = "ArrayOf" + titleCase(t.element.name);
            }
            else if (t.kind === "tuple") {
                assert(t.items.length === 2);
                assert(t.items[0].kind === "base" && t.items[0].name === "uinteger");
                assert(t.items[1].kind === "base" && t.items[1].name === "uinteger");
                name = "UintegerPair";
            }
            else if (t.kind === "or") {
                throw new Error("Nested or types are not supported");
            }
            else if (t.kind === "booleanLiteral") {
                name = t.value ? "True" : "False";
            }
            else if (t.kind === "literal") {
                assert(t.value.properties.length === 0);
                name = "EmptyObject";
            }
            else {
                name = "_TODO_or_" + t.kind + "_";
            }
            members.push({ type: t, name });
        }

        const name = members.map(m => titleCase(m.name)).join("Or");
        unionTypes.set(name, members);
        write(name);
    }
    if (nullable && wasOptional) {
        write("]");
    }
    return omitEmpty;
}

function writeTypeElement(t: Type, wasOptional = false): boolean {
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
        case "booleanLiteral":
            write("bool");
            break;
        case "integerLiteral":
            write("int32");
            break;
        case "literal":
            assert(t.value.properties.length === 0);
            write("struct{}");
            break;
        case "tuple":
            assert(t.items.length === 2);
            assert(t.items[0].kind === "base" && t.items[0].name === "uinteger");
            assert(t.items[1].kind === "base" && t.items[1].name === "uinteger");
            write("[2]uint32");
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
            return writeOr(t, wasOptional);
        default:
            write("TODO_" + t.kind);
            break;
    }
    return wasOptional;
}

// Generation

writeLine("// Code generated by genlsp; DO NOT EDIT.");
writeLine("");
writeLine("package lsproto2");
writeLine("");
writeLine(`import "encoding/json"`);
writeLine(`import "fmt"`);
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

writeLine("// Structures\n");

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

        const omitEmpty = writeTypeElement(p.type, !!p.optional);
        write(' `json:"');
        write(p.name);
        if (omitEmpty) {
            write(",omitempty");
        }
        finishLine('"`');
    }

    dedent();
    writeLine("}");
    writeLine("\n");
}

writeLine("// Enumerations\n");

for (const t of model.enumerations) {
    writeDocumentation(t.documentation);
    writeDeprecation(t.deprecated);

    let underlyingType: string;
    switch (t.type.name) {
        case "string":
            underlyingType = "string";
            break;
        case "integer":
            underlyingType = "int32";
            break;
        case "uinteger":
            underlyingType = "uint32";
            break;
    }

    writeLine("type " + t.name + " " + underlyingType);
    writeLine("\n");

    function valueToLiteral(v: string | number): string {
        return typeof v === "string" ? '"' + v + '"' : `${v}`;
    }

    writeLine("const (");
    indent();
    for (const v of t.values) {
        writeDocumentation(v.documentation);
        writeDeprecation(v.deprecated);

        startLine(t.name);
        write(v.name);
        write(" ");
        write(t.name);
        write(" = ");
        finishLine(valueToLiteral(v.value));
    }
    dedent();
    writeLine(")");

    writeLine("");

    writeLine("func (e *" + t.name + ") UnmarshalJSON(data []byte) error {");
    indent();
    writeLine("var v " + underlyingType);
    writeLine("if err := json.Unmarshal(data, &v); err != nil {");
    indent();
    writeLine("return err");
    dedent();
    writeLine("}");
    writeLine("switch v {");
    indent();
    const values = [...new Set<string | number>(t.values.map(v => v.value))].sort(compareValues);
    for (let i = 0; i < values.length; i++) {
        const v = values[i];
        if (i === 0) {
            startLine("case ");
        }
        write(valueToLiteral(v));
        if (i === values.length - 1) {
            writeLine(":");
        }
        else {
            if (i % 3 === 2) {
                writeLine(",");
            }
            else {
                write(", ");
            }
        }
    }
    indent();
    writeLine("*e = " + t.name + "(v)");
    writeLine("return nil");
    writeLine("default:");
    indent();
    writeLine(`return fmt.Errorf("unknown ${t.name} value: %v", v)`);
    dedent();
    writeLine("}");
    dedent();
    writeLine("}");
    writeLine("");
}

writeLine("// Type aliases\n");

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

writeLine("// Requests\n");

for (const t of model.requests) {
    writeDocumentation(t.documentation);
    writeDeprecation(t.deprecated);
    writeLine("const MethodRequest" + methodNameToIdentifier(t.method) + ' Method = "' + t.method + '"\n');
}

writeLine("// Notifications\n");

for (const t of model.notifications) {
    writeDocumentation(t.documentation);
    writeDeprecation(t.deprecated);
    writeLine("const MethodNotification" + methodNameToIdentifier(t.method) + ' Method = "' + t.method + '"\n');
}

writeLine("// Union types\n");

writeLine("func assertOnlyOneTrue(message string, values ...bool) {");
indent();
writeLine("count := 0");
writeLine("for _, v := range values {");
indent();
writeLine("if v {");
indent();
writeLine("count++");
dedent();
writeLine("}");
dedent();
writeLine("}");
writeLine("if count != 1 {");
indent();
writeLine("panic(message)");
dedent();
writeLine("}");
dedent();
writeLine("}");
writeLine("");

for (const [name, members] of unionTypes) {
    writeLine("type " + name + " struct {");
    indent();

    for (const member of members) {
        startLine(titleCase(member.name) + " *");
        writeTypeElement(member.type, false);
        finishLine("");
    }

    dedent();
    writeLine("}");
    writeLine("");

    writeLine("func (o " + name + ") MarshalJSON() ([]byte, error) {");
    indent();
    startLine('assertOnlyOneTrue("invalid union type", ');
    for (let i = 0; i < members.length; i++) {
        if (i > 0) {
            write(", ");
        }
        write("o." + titleCase(members[i].name) + " != nil");
    }
    finishLine(")");

    for (const member of members) {
        const name = titleCase(member.name);
        writeLine("if o." + name + " != nil {");
        indent();
        writeLine("return json.Marshal(o." + name + ")");
        dedent();
        writeLine("}");
    }
    writeLine('panic("unreachable")');
    dedent();
    writeLine("}");
    writeLine("");

    // TODO: do this way more efficiently
    // TODO: this doesn't work when union members overlap
    writeLine("func (o *" + name + ") UnmarshalJSON(data []byte) error {");
    indent();
    writeLine("*o = " + name + "{}");
    for (const member of members) {
        const name = titleCase(member.name);
        const local = "v" + name;
        startLine("var " + local + " ");
        writeTypeElement(member.type);
        finishLine("");
        writeLine("if err := json.Unmarshal(data, &" + local + "); err == nil {");
        indent();
        writeLine("o." + name + " = &" + local);
        writeLine("return nil");
        dedent();
        writeLine("}");
    }
    writeLine(`return fmt.Errorf("invalid ${name}: %s", data)`);
    dedent();
    writeLine("}");
}

fs.writeFileSync(out, parts.join(""));
