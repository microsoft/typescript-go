#!/usr/bin/env node

import assert from "node:assert";
import cp from "node:child_process";
import fs from "node:fs";
import path from "node:path";
import url from "node:url";
import which from "which";

/**
 * @import { MetaModel, OrType, Type, BaseTypes, BaseType, ReferenceType, ArrayType, MapType } from "./metaModelSchema.mts"
 */
void 0;

const __filename = url.fileURLToPath(new URL(import.meta.url));
const __dirname = path.dirname(__filename);

const out = path.resolve(__dirname, "../lsp_generated.go");
const metaModelPath = path.resolve(__dirname, "metaModel.json");

if (!fs.existsSync(metaModelPath)) {
    console.error("Meta model file not found; did you forget to run fetchModel.mjs?");
    process.exit(1);
}

/** @type {MetaModel} */
const model = JSON.parse(fs.readFileSync(metaModelPath, "utf-8"));

/**
 * Represents a type in our intermediate type system
 * @typedef {Object} GoType
 * @property {string} name - Name of the type in Go
 * @property {boolean} isStruct - Whether this type is a struct
 * @property {boolean} needsPointer - Whether this type should be used with a pointer
 * @property {string} [importPath] - Import path if needed
 * @property {string} [jsonUnmarshaling] - Custom JSON unmarshaling code if required
 */

/**
 * @typedef {Object} TypeRegistry
 * @property {Map<string, GoType>} types - Map of type names to types
 * @property {Map<string, string>} literalTypes - Map from literal values to type names
 * @property {Map<string, {name: string, types: Type[]}[]>} unionTypes - Map of union type names to their component types
 * @property {Set<string>} generatedTypes - Set of types that have been generated
 * @property {Map<string, Map<string, {identifier: string, documentation: string, deprecated: string}>>} enumValuesByType - Map of enum type names to their values
 */

/**
 * @type {TypeRegistry}
 */
const registry = {
    types: new Map(),
    literalTypes: new Map(),
    unionTypes: new Map(),
    generatedTypes: new Set(),
    enumValuesByType: new Map(),
};

/**
 * @param {string} s
 */
function titleCase(s) {
    return s.charAt(0).toUpperCase() + s.slice(1);
}

/**
 * @param {BaseTypes} baseType
 * @returns {GoType}
 */
function mapBaseTypeToGo(baseType) {
    switch (baseType) {
        case "integer":
            return { name: "int32", isStruct: false, needsPointer: false };
        case "uinteger":
            return { name: "uint32", isStruct: false, needsPointer: false };
        case "string":
            return { name: "string", isStruct: false, needsPointer: false };
        case "boolean":
            return { name: "bool", isStruct: false, needsPointer: false };
        case "URI":
            return { name: "URI", isStruct: false, needsPointer: false };
        case "DocumentUri":
            return { name: "DocumentUri", isStruct: false, needsPointer: false };
        case "decimal":
            return { name: "float64", isStruct: false, needsPointer: false };
        case "RegExp":
            return { name: "string", isStruct: false, needsPointer: false }; // Using string for RegExp
        case "null":
            return { name: "NullType", isStruct: true, needsPointer: true }; // Special handling for null
        default:
            console.warn(`Unknown base type: ${baseType}`);
            return { name: `ANY_${baseType}`, isStruct: false, needsPointer: false };
    }
}

/**
 * @param {Type} type
 * @returns {GoType}
 */
function resolveType(type) {
    switch (type.kind) {
        case "base":
            return mapBaseTypeToGo(type.name);

        case "reference":
            // If it's a reference, we need to check if we know this type
            if (registry.types.has(type.name)) {
                const refType = registry.types.get(type.name);
                if (refType !== undefined) {
                    return refType;
                }
            }

            // By default, assume referenced types are structs that need pointers
            // This will be updated as we process all types
            const refType = { name: type.name, isStruct: true, needsPointer: true };
            registry.types.set(type.name, refType);
            return refType;

        case "array": {
            const elementType = resolveType(type.element);
            // Arrays of structs should be arrays of pointers to structs
            const arrayTypeName = elementType.needsPointer
                ? `[]*${elementType.name}`
                : `[]${elementType.name}`;
            return {
                name: arrayTypeName,
                isStruct: false,
                needsPointer: false,
            };
        }

        case "map": {
            const keyType = type.key.kind === "base"
                ? mapBaseTypeToGo(type.key.name).name
                : resolveType(type.key).name;

            const valueType = resolveType(type.value);
            const valueTypeName = valueType.needsPointer && valueType.isStruct
                ? `*${valueType.name}`
                : valueType.name;

            return {
                name: `map[${keyType}]${valueTypeName}`,
                isStruct: false,
                needsPointer: false,
            };
        }

        case "tuple": {
            if (
                type.items.length === 2 &&
                type.items[0].kind === "base" && type.items[0].name === "uinteger" &&
                type.items[1].kind === "base" && type.items[1].name === "uinteger"
            ) {
                return { name: "[2]uint32", isStruct: false, needsPointer: false };
            }

            // For other tuples, create a custom type name
            const typeName = `Tuple${
                type.items.map(item => {
                    const resolvedType = resolveType(item);
                    return titleCase(resolvedType.name.replace(/[\[\]{}*]/g, ""));
                }).join("And")
            }`;

            if (!registry.unionTypes.has(typeName)) {
                registry.unionTypes.set(typeName, []);
            }

            const union = registry.unionTypes.get(typeName);
            if (union) {
                for (const item of type.items) {
                    union.push({ name: resolveType(item).name, types: [item] });
                }
            }

            return { name: typeName, isStruct: true, needsPointer: true };
        }

        case "stringLiteral": {
            const typeName = `StringLiteral${titleCase(type.value)}`;
            registry.literalTypes.set(String(type.value), typeName);
            return { name: typeName, isStruct: true, needsPointer: false };
        }

        case "integerLiteral": {
            const typeName = `IntegerLiteral${type.value}`;
            registry.literalTypes.set(String(type.value), typeName);
            return { name: typeName, isStruct: true, needsPointer: false };
        }

        case "booleanLiteral": {
            const typeName = `BooleanLiteral${type.value ? "True" : "False"}`;
            registry.literalTypes.set(String(type.value), typeName);
            return { name: typeName, isStruct: true, needsPointer: false };
        }

        case "literal":
            // Empty object literal
            if (type.value.properties.length === 0) {
                return { name: "struct{}", isStruct: true, needsPointer: false };
            }

            // Handle literal structs (this is a simplification, may need enhancement)
            const literalTypeName = `AnonymousStruct${Math.floor(Math.random() * 10000)}`;
            const literalType = { name: literalTypeName, isStruct: true, needsPointer: true };
            registry.types.set(literalTypeName, literalType);
            return literalType;

        case "or": {
            return handleOrType(type);
        }

        case "and": {
            // For AND types, we'll create a struct that embeds all component types
            const typeName = `And${
                type.items.map(item => {
                    if (item.kind === "reference") {
                        return item.name;
                    }
                    return "Anonymous";
                }).join("")
            }`;

            const andType = { name: typeName, isStruct: true, needsPointer: true };
            registry.types.set(typeName, andType);
            return andType;
        }

        default: {
            // This is a safeguard for the TypeScript compiler, should not happen in practice
            // Handle unknown type kind safely by using type assertion with the 'any' type
            // @ts-ignore - We know this will have a kind property, but TypeScript thinks it's 'never'
            const unknownKind = String(type["kind"] || "unknown");
            console.warn(`Unhandled type kind: ${unknownKind}`);
            return { name: `ANY_${unknownKind}`, isStruct: false, needsPointer: false };
        }
    }
}

/**
 * @param {OrType} orType
 * @returns {GoType}
 */
function handleOrType(orType) {
    // Check for nullable types (OR with null)
    const nullIndex = orType.items.findIndex(item => item.kind === "base" && item.name === "null");

    // If it's nullable and only has one other type
    if (nullIndex !== -1 && orType.items.length === 2) {
        const otherType = orType.items[1 - nullIndex];
        const resolvedType = resolveType(otherType);

        // Make it a pointer type
        return {
            name: resolvedType.name,
            isStruct: resolvedType.isStruct,
            needsPointer: true,
        };
    }

    // Filter out null if present
    const types = nullIndex !== -1
        ? orType.items.filter((_, i) => i !== nullIndex)
        : orType.items;

    // If only one type remains after filtering null
    if (types.length === 1) {
        const resolvedType = resolveType(types[0]);
        return {
            name: resolvedType.name,
            isStruct: resolvedType.isStruct,
            needsPointer: nullIndex !== -1 ? true : resolvedType.needsPointer,
        };
    }

    // For multiple types, create a union type
    const memberNames = types.map(type => {
        if (type.kind === "reference") {
            return type.name;
        }
        else if (type.kind === "base") {
            return titleCase(type.name);
        }
        else if (
            type.kind === "array" &&
            (type.element.kind === "reference" || type.element.kind === "base")
        ) {
            return `${
                titleCase(
                    type.element.kind === "reference"
                        ? type.element.name
                        : type.element.name,
                )
            }s`;
        }
        else if (type.kind === "literal" && type.value.properties.length === 0) {
            return "EmptyObject";
        }
        else if (type.kind === "tuple") {
            return "Tuple";
        }
        else {
            return `Type${Math.floor(Math.random() * 10000)}`;
        }
    });

    const unionTypeName = memberNames.map(titleCase).join("Or");

    if (!registry.unionTypes.has(unionTypeName)) {
        registry.unionTypes.set(unionTypeName, []);
    }

    const union = registry.unionTypes.get(unionTypeName);
    if (union) {
        for (let i = 0; i < types.length; i++) {
            union.push({ name: memberNames[i], types: [types[i]] });
        }
    }

    return {
        name: unionTypeName,
        isStruct: true,
        needsPointer: true,
    };
}

/**
 * First pass: Resolve all type information
 */
function buildTypeSystem() {
    // Register built-in types
    registry.types.set("LSPAny", { name: "any", isStruct: false, needsPointer: false });
    registry.types.set("NullType", { name: "NullType", isStruct: true, needsPointer: true });

    // Keep track of used enum identifiers across all enums to avoid conflicts
    const usedEnumIdentifiers = new Set();

    // Process all enumerations first to make them available for struct fields
    for (const enumeration of model.enumerations) {
        // Register the enum type with its own name rather than the base type
        registry.types.set(enumeration.name, {
            name: enumeration.name, // Use the enum type name, not the base type
            isStruct: false,
            needsPointer: false,
        });

        // Create a map for this enum's values (not an array)
        const enumValues = new Map();

        // Process values for this enum
        for (const value of enumeration.values) {
            // Generate a unique identifier for this enum constant
            let identifier = `${enumeration.name}${value.name}`;

            // If this identifier is already used, create a more unique one
            if (usedEnumIdentifiers.has(identifier)) {
                // Try with underscores
                identifier = `${enumeration.name}_${value.name}`;

                // If still not unique, add a numeric suffix
                let counter = 1;
                while (usedEnumIdentifiers.has(identifier)) {
                    identifier = `${enumeration.name}_${value.name}_${counter++}`;
                }
            }

            // Mark this identifier as used
            usedEnumIdentifiers.add(identifier);

            // Store the entry in the map with the value literal as the key
            // and an object with all needed information as the value
            enumValues.set(String(value.value), {
                identifier,
                documentation: value.documentation,
                deprecated: value.deprecated,
            });
        }

        // Store the map of values for this enum
        registry.enumValuesByType.set(enumeration.name, enumValues);
    }

    // Process all structures
    for (const structure of model.structures) {
        registry.types.set(structure.name, {
            name: structure.name,
            isStruct: true,
            needsPointer: true,
        });
    }

    // Process all type aliases
    for (const typeAlias of model.typeAliases) {
        const resolvedType = resolveType(typeAlias.type);
        registry.types.set(typeAlias.name, {
            name: resolvedType.name,
            isStruct: resolvedType.isStruct,
            needsPointer: resolvedType.needsPointer,
        });
    }
}

/**
 * @param {string | undefined} s
 * @returns {string}
 */
function formatDocumentation(s) {
    if (!s) return "";

    let formatted = s.split("\n")
        .map(line => {
            line = line.replace(/(\w ) +/g, "$1");
            line = line.replace(/\{@link(?:code)?.*?([^} ]+)\}/g, "$1");
            line = line.replace(/@since (.*)/g, "Since: $1");
            if (line.startsWith("@deprecated")) {
                return null;
            }
            if (line.startsWith("@proposed")) {
                return "// Proposed.";
            }
            return "// " + line;
        })
        .filter(Boolean)
        .join("\n");

    return formatted ? formatted + "\n" : "";
}

/**
 * @param {string | undefined} deprecated
 * @returns {string}
 */
function formatDeprecation(deprecated) {
    if (!deprecated) return "";
    return "//\n// Deprecated: " + deprecated + "\n";
}

/** @type {string[]} */
const parts = [];

/**
 * @param {string} s
 */
function write(s) {
    parts.push(s);
}

/**
 * @param {string} s
 */
function writeLine(s = "") {
    parts.push(s + "\n");
}

/**
 * Generate the Go code
 */
function generateCode() {
    // File header
    writeLine("// Code generated by generate2.mjs; DO NOT EDIT.");
    writeLine("");
    writeLine("package lsproto");
    writeLine("");
    writeLine(`import (`);
    writeLine(`\t"encoding/json"`);
    writeLine(`\t"fmt"`);
    writeLine(`)`);
    writeLine("");
    writeLine("// Meta model version " + model.metaData.version);
    writeLine("");

    // Write basic NullType for null values - skip assertOnlyOne which exists in lsp.go
    writeLine("// NullType represents a JSON null value");
    writeLine("type NullType struct{}");
    writeLine("");
    writeLine("func (n NullType) MarshalJSON() ([]byte, error) {");
    writeLine('\treturn []byte("null"), nil');
    writeLine("}");
    writeLine("");
    writeLine("func (n *NullType) UnmarshalJSON(data []byte) error {");
    writeLine('\tif string(data) != "null" {');
    writeLine('\t\treturn fmt.Errorf("invalid null value: %s", data)');
    writeLine("\t}");
    writeLine("\treturn nil");
    writeLine("}");
    writeLine("");

    // Skip helper function for union types - already exists in lsp.go

    // Generate structures
    writeLine("// Structures\n");

    // Keep track of generated types to avoid duplicates
    const generatedTypes = new Set();

    for (const structure of model.structures) {
        // Skip URI and DocumentUri as they are already defined in lsp.go
        if (structure.name === "URI" || structure.name === "DocumentUri") {
            continue;
        }

        write(formatDocumentation(structure.documentation));
        write(formatDeprecation(structure.deprecated));

        writeLine(`type ${structure.name} struct {`);

        // First embed extended types
        for (const e of structure.extends || []) {
            if (e.kind !== "reference") {
                throw new Error(`Unexpected extends kind: ${e.kind}`);
            }
            writeLine(`\t${e.name}`);
        }

        // Then embed mixin types
        for (const m of structure.mixins || []) {
            if (m.kind !== "reference") {
                throw new Error(`Unexpected mixin kind: ${m.kind}`);
            }
            writeLine(`\t${m.name}`);
        }

        // Insert a blank line after embeds if there were any
        if (
            (structure.extends && structure.extends.length > 0) ||
            (structure.mixins && structure.mixins.length > 0)
        ) {
            writeLine("");
        }

        // Then properties
        for (const prop of structure.properties) {
            write(formatDocumentation(prop.documentation));
            write(formatDeprecation(prop.deprecated));

            const type = resolveType(prop.type);
            const goType = prop.optional || type.needsPointer ? `*${type.name}` : type.name;

            writeLine(`\t${titleCase(prop.name)} ${goType} \`json:"${prop.name}${prop.optional ? ",omitempty" : ""}"\``);
            writeLine("");
        }

        writeLine("}");
        writeLine("");

        generatedTypes.add(structure.name);
    }

    // Generate enumerations
    writeLine("// Enumerations\n");

    for (const enumeration of model.enumerations) {
        write(formatDocumentation(enumeration.documentation));
        write(formatDeprecation(enumeration.deprecated));

        let baseType;
        switch (enumeration.type.name) {
            case "string":
                baseType = "string";
                break;
            case "integer":
                baseType = "int32";
                break;
            case "uinteger":
                baseType = "uint32";
                break;
            default:
                baseType = "string";
        }

        writeLine(`type ${enumeration.name} ${baseType}`);
        writeLine("");

        // Get the pre-processed enum entries map that avoids duplicates
        const enumValues = registry.enumValuesByType.get(enumeration.name);
        if (!enumValues || !enumValues.size) {
            continue; // Skip if no entries (shouldn't happen)
        }

        writeLine("const (");

        // Process entries with unique identifiers
        for (const [value, entry] of enumValues.entries()) {
            write(formatDocumentation(entry.documentation));
            write(formatDeprecation(entry.deprecated));

            let valueLiteral;
            // Handle string values
            if (enumeration.type.name === "string") {
                valueLiteral = `"${String(value).replace(/^"|"$/g, "")}"`;
            }
            else {
                valueLiteral = String(value);
            }

            writeLine(`\t${entry.identifier} ${enumeration.name} = ${valueLiteral}`);
        }

        writeLine(")");
        writeLine("");

        // Add custom JSON unmarshaling
        writeLine(`func (e *${enumeration.name}) UnmarshalJSON(data []byte) error {`);
        writeLine(`\tvar v ${baseType}`);
        writeLine(`\tif err := json.Unmarshal(data, &v); err != nil {`);
        writeLine(`\t\treturn err`);
        writeLine(`\t}`);
        writeLine(`\t*e = ${enumeration.name}(v)`);
        writeLine(`\treturn nil`);
        writeLine(`}`);
        writeLine("");

        generatedTypes.add(enumeration.name);
    }

    // Generate type aliases
    writeLine("// Type aliases\n");

    for (const typeAlias of model.typeAliases) {
        // Skip URI and DocumentUri as they are already defined in lsp.go
        if (typeAlias.name === "URI" || typeAlias.name === "DocumentUri") {
            continue;
        }

        write(formatDocumentation(typeAlias.documentation));
        write(formatDeprecation(typeAlias.deprecated));

        if (typeAlias.name === "LSPAny") {
            writeLine("type LSPAny = any");
            writeLine("");
            continue;
        }

        const resolvedType = resolveType(typeAlias.type);
        writeLine(`type ${typeAlias.name} = ${resolvedType.name}`);
        writeLine("");

        generatedTypes.add(typeAlias.name);
    }

    generateUnionTypes();

    // Generate literal types
    writeLine("// Literal types\n");

    for (const [value, name] of registry.literalTypes.entries()) {
        // Skip if already generated
        if (generatedTypes.has(name)) {
            continue;
        }

        const jsonValue = JSON.stringify(value);

        writeLine(`// ${name} is a literal type for ${jsonValue}`);
        writeLine(`type ${name} struct{}`);
        writeLine("");

        writeLine(`func (o ${name}) MarshalJSON() ([]byte, error) {`);
        writeLine(`\treturn []byte(${jsonValue}), nil`);
        writeLine(`}`);
        writeLine("");

        writeLine(`func (o *${name}) UnmarshalJSON(data []byte) error {`);
        writeLine(`\tif string(data) != ${jsonValue} {`);
        writeLine(`\t\treturn fmt.Errorf("invalid ${name}: %s", string(data))`);
        writeLine(`\t}`);
        writeLine(`\treturn nil`);
        writeLine(`}`);
        writeLine("");

        generatedTypes.add(name);
    }

    // Generate Methods for requests and notifications
    writeLine("// Methods\n");

    // Method type exists in lsp.go, so skip declaring it

    writeLine("// Request Methods");
    writeLine("const (");
    for (const request of model.requests) {
        write(formatDocumentation(request.documentation));
        write(formatDeprecation(request.deprecated));

        const methodName = request.method.split("/")
            .map(v => v === "$" ? "" : titleCase(v))
            .join("");

        writeLine(`\tMethod${methodName} Method = "${request.method}"`);
    }
    writeLine(")");
    writeLine("");

    writeLine("// Notification Methods");
    writeLine("const (");
    for (const notification of model.notifications) {
        write(formatDocumentation(notification.documentation));
        write(formatDeprecation(notification.deprecated));

        const methodName = notification.method.split("/")
            .map(v => v === "$" ? "" : titleCase(v))
            .join("");

        writeLine(`\tMethod${methodName} Method = "${notification.method}"`);
    }
    writeLine(")");
    writeLine("");

    // Unmarshallers
    writeLine("// Unmarshallers\n");

    // Note: The unmarshallerFor function already exists in lsp.go, so we don't generate it

    // The unmarshallers map is expected by jsonrpc.go
    writeLine("var unmarshallers = map[Method]func([]byte) (any, error){");

    // Client-to-server requests
    for (const request of model.requests) {
        if (request.messageDirection !== "clientToServer" && request.messageDirection !== "both") {
            continue;
        }

        const methodName = request.method.split("/")
            .map(v => v === "$" ? "" : titleCase(v))
            .join("");

        if (!request.params) {
            continue;
        }

        let typeName;
        if (Array.isArray(request.params)) {
            // This shouldn't typically happen in the LSP spec
            typeName = "any";
        }
        else if (request.params.kind === "reference") {
            typeName = request.params.name;
        }
        else {
            const resolvedType = resolveType(request.params);
            typeName = resolvedType.name;
        }

        // Make sure to use the function properly
        writeLine(`\tMethod${methodName}: unmarshallerFor[${typeName}],`);
    }

    // Client-to-server notifications
    for (const notification of model.notifications) {
        if (notification.messageDirection !== "clientToServer" && notification.messageDirection !== "both") {
            continue;
        }

        const methodName = notification.method.split("/")
            .map(v => v === "$" ? "" : titleCase(v))
            .join("");

        if (!notification.params) {
            continue;
        }

        let typeName;
        if (Array.isArray(notification.params)) {
            // This shouldn't typically happen in the LSP spec
            typeName = "any";
        }
        else if (notification.params.kind === "reference") {
            typeName = notification.params.name;
        }
        else {
            const resolvedType = resolveType(notification.params);
            typeName = resolvedType.name;
        }

        // Make sure to use the function properly
        writeLine(`\tMethod${methodName}: unmarshallerFor[${typeName}],`);
    }

    writeLine("}");

    return parts.join("");
}

/**
 * Main function
 */
function main() {
    try {
        buildTypeSystem();
        const generatedCode = generateCode();
        fs.writeFileSync(out, generatedCode);

        // Format with gofmt
        const gofmt = which.sync("gofmt");
        cp.execFileSync(gofmt, ["-w", out]);

        console.log(`Successfully generated ${out}`);
    }
    catch (error) {
        console.error("Error generating code:", error);
        process.exit(1);
    }
}

main();

/**
 * Generate union types
 */
function generateUnionTypes() {
    writeLine("// Union types\n");

    for (const [name, members] of registry.unionTypes.entries()) {
        // Skip if already generated
        if (registry.generatedTypes.has(name)) {
            continue;
        }

        writeLine(`type ${name} struct {`);

        // Track field names to avoid duplicates in the Go struct
        const usedFieldNames = new Set();

        for (const member of members) {
            let memberType;
            if (member.types.length === 1) {
                const type = resolveType(member.types[0]);
                memberType = type.name;
            }
            else {
                // This case shouldn't really happen with our current approach
                memberType = "any";
            }

            // Use a unique field name by adding index if needed
            let fieldName = titleCase(member.name);
            let counter = 1;

            while (usedFieldNames.has(fieldName)) {
                fieldName = `${titleCase(member.name)}${counter++}`;
            }

            usedFieldNames.add(fieldName);
            writeLine(`\t${fieldName} *${memberType}`);
        }

        writeLine(`}`);
        writeLine("");

        // Marshal method
        writeLine(`func (o ${name}) MarshalJSON() ([]byte, error) {`);

        // Create assertion to ensure only one field is set at a time
        write(`\tassertOnlyOne("more than one element of ${name} is set", `);

        // Get field names again for the assertion
        const fieldNames = [];
        const seenFieldNames = new Set();

        for (const member of members) {
            let fieldName = titleCase(member.name);
            let counter = 1;

            while (seenFieldNames.has(fieldName)) {
                fieldName = `${titleCase(member.name)}${counter++}`;
            }

            seenFieldNames.add(fieldName);
            fieldNames.push(fieldName);
        }

        // Write the assertion conditions
        for (let i = 0; i < fieldNames.length; i++) {
            if (i > 0) write(", ");
            write(`o.${fieldNames[i]} != nil`);
        }
        writeLine(`)`);
        writeLine("");

        // Write the marshal logic for each field
        let fieldIndex = 0;
        seenFieldNames.clear();

        for (const member of members) {
            let fieldName = titleCase(member.name);
            let counter = 1;

            while (seenFieldNames.has(fieldName)) {
                fieldName = `${titleCase(member.name)}${counter++}`;
            }

            seenFieldNames.add(fieldName);

            writeLine(`\tif o.${fieldName} != nil {`);
            writeLine(`\t\treturn json.Marshal(*o.${fieldName})`);
            writeLine(`\t}`);
            fieldIndex++;
        }

        writeLine(`\treturn []byte("null"), nil`);
        writeLine(`}`);
        writeLine("");

        // Unmarshal method
        writeLine(`func (o *${name}) UnmarshalJSON(data []byte) error {`);
        writeLine(`\t*o = ${name}{}`);
        writeLine(`\tif string(data) == "null" {`);
        writeLine(`\t\treturn nil`);
        writeLine(`\t}`);
        writeLine("");

        // Write the unmarshal logic for each field
        seenFieldNames.clear();

        for (const member of members) {
            let fieldName = titleCase(member.name);
            let counter = 1;

            while (seenFieldNames.has(fieldName)) {
                fieldName = `${titleCase(member.name)}${counter++}`;
            }

            seenFieldNames.add(fieldName);

            let memberType;
            if (member.types.length === 1) {
                const type = resolveType(member.types[0]);
                memberType = type.name;
            }
            else {
                memberType = "any";
            }

            writeLine(`\t{`);
            writeLine(`\t\tvar v ${memberType}`);
            writeLine(`\t\tif err := json.Unmarshal(data, &v); err == nil {`);
            writeLine(`\t\t\to.${fieldName} = &v`);
            writeLine(`\t\t\treturn nil`);
            writeLine(`\t\t}`);
            writeLine(`\t}`);
        }

        writeLine(`\treturn fmt.Errorf("cannot unmarshal %s into ${name}", string(data))`);
        writeLine(`}`);
        writeLine("");

        registry.generatedTypes.add(name);
    }
}
