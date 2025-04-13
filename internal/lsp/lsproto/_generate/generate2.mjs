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

// Define output file paths
const outputDir = path.resolve(__dirname, "../");
const typesFilePath = path.resolve(outputDir, "lsptypes_generated.go");
const enumsFilePath = path.resolve(outputDir, "lspenums_generated.go");
const unionsFilePath = path.resolve(outputDir, "lspunions_generated.go");
const methodsFilePath = path.resolve(outputDir, "lspmethods_generated.go");

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
 * @typedef {Object} TypeInfo
 * @property {Map<string, GoType>} types - Map of type names to types
 * @property {Map<string, string>} literalTypes - Map from literal values to type names
 * @property {Map<string, {name: string, types: Type[]}[]>} unionTypes - Map of union type names to their component types
 * @property {Set<string>} generatedTypes - Set of types that have been generated
 * @property {Map<string, Map<string, {identifier: string, documentation: string, deprecated: string}>>} enumValuesByType - Map of enum type names to their values
 * @property {Map<string, string>} unionTypeAliases - Map from union type name to alias name
 */

/**
 * @type {TypeInfo}
 */
const typeInfo = {
    types: new Map(),
    literalTypes: new Map(),
    unionTypes: new Map(),
    generatedTypes: new Set(),
    enumValuesByType: new Map(),
    unionTypeAliases: new Map(), // Map from union type name to alias name
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
    // Special case for the LSP "any" type structure, which would normally become a complex union
    if (
        type.kind === "or" && type.items.length >= 6 &&
        type.items.some(item => item.kind === "reference" && item.name === "LSPObject") &&
        type.items.some(item => item.kind === "reference" && item.name === "LSPArray") &&
        type.items.some(item => item.kind === "base" && item.name === "string") &&
        type.items.some(item => item.kind === "base" && item.name === "integer") &&
        type.items.some(item => item.kind === "base" && item.name === "boolean")
    ) {
        return { name: "LSPAny", isStruct: false, needsPointer: false };
    }

    switch (type.kind) {
        case "base":
            return mapBaseTypeToGo(type.name);

        case "reference":
            // If it's a reference, we need to check if we know this type
            if (typeInfo.types.has(type.name)) {
                const refType = typeInfo.types.get(type.name);
                if (refType !== undefined) {
                    return refType;
                }
            }

            // By default, assume referenced types are structs that need pointers
            // This will be updated as we process all types
            const refType = { name: type.name, isStruct: true, needsPointer: true };
            typeInfo.types.set(type.name, refType);
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

            if (!typeInfo.unionTypes.has(typeName)) {
                typeInfo.unionTypes.set(typeName, []);
            }

            const union = typeInfo.unionTypes.get(typeName);
            if (union) {
                for (const item of type.items) {
                    union.push({ name: resolveType(item).name, types: [item] });
                }
            }

            return { name: typeName, isStruct: true, needsPointer: true };
        }

        case "stringLiteral": {
            const typeName = `StringLiteral${titleCase(type.value)}`;
            typeInfo.literalTypes.set(String(type.value), typeName);
            return { name: typeName, isStruct: true, needsPointer: false };
        }

        case "integerLiteral": {
            const typeName = `IntegerLiteral${type.value}`;
            typeInfo.literalTypes.set(String(type.value), typeName);
            return { name: typeName, isStruct: true, needsPointer: false };
        }

        case "booleanLiteral": {
            const typeName = `BooleanLiteral${type.value ? "True" : "False"}`;
            typeInfo.literalTypes.set(String(type.value), typeName);
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
            typeInfo.types.set(literalTypeName, literalType);
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
            typeInfo.types.set(typeName, andType);
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
 * Handle OR types specially
 * @param {import('./metaModelSchema.mts').OrType} orType
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

    // Check if all items are references - if so, we can use their names directly
    const allReferences = types.every(type => type.kind === "reference");
    if (allReferences) {
        const memberNames = types.map(type => {
            // Type assertion to help TypeScript understand the type
            const refType = /** @type {import('./metaModelSchema.mts').ReferenceType} */ (type);
            return refType.name;
        });
        const unionTypeName = memberNames.map(titleCase).join("Or");

        if (!typeInfo.unionTypes.has(unionTypeName)) {
            typeInfo.unionTypes.set(unionTypeName, []);
        }

        const union = typeInfo.unionTypes.get(unionTypeName);
        if (union) {
            for (let i = 0; i < types.length; i++) {
                const refType = /** @type {import('./metaModelSchema.mts').ReferenceType} */ (types[i]);
                union.push({
                    name: refType.name,
                    types: [types[i]],
                });
            }
        }

        return {
            name: unionTypeName,
            isStruct: true,
            needsPointer: true,
        };
    }

    // For mixed types, create a union type with more careful naming
    const memberNames = types.map(type => {
        if (type.kind === "reference") {
            const refType = /** @type {import('./metaModelSchema.mts').ReferenceType} */ (type);
            return refType.name;
        }
        else if (type.kind === "base") {
            const baseType = /** @type {import('./metaModelSchema.mts').BaseType} */ (type);
            return titleCase(baseType.name);
        }
        else if (
            type.kind === "array" &&
            (type.element.kind === "reference" || type.element.kind === "base")
        ) {
            const arrayType = /** @type {import('./metaModelSchema.mts').ArrayType} */ (type);
            const elementName = arrayType.element.kind === "reference"
                ? /** @type {import('./metaModelSchema.mts').ReferenceType} */ (arrayType.element).name
                : /** @type {import('./metaModelSchema.mts').BaseType} */ (arrayType.element).name;
            return `${titleCase(elementName)}s`;
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

    if (!typeInfo.unionTypes.has(unionTypeName)) {
        typeInfo.unionTypes.set(unionTypeName, []);
    }

    const union = typeInfo.unionTypes.get(unionTypeName);
    if (union) {
        for (let i = 0; i < types.length; i++) {
            const resolvedType = resolveType(types[i]);
            union.push({
                name: memberNames[i],
                types: [types[i]],
            });
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
function collectTypeDefinitions() {
    // Register built-in types
    typeInfo.types.set("NullType", { name: "NullType", isStruct: true, needsPointer: true });
    typeInfo.types.set("LSPAny", { name: "any", isStruct: false, needsPointer: false });

    // Keep track of used enum identifiers across all enums to avoid conflicts
    const usedEnumIdentifiers = new Set();

    // Process all enumerations first to make them available for struct fields
    for (const enumeration of model.enumerations) {
        // Register the enum type with its own name rather than the base type
        typeInfo.types.set(enumeration.name, {
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
        typeInfo.enumValuesByType.set(enumeration.name, enumValues);
    }

    // Process all structures
    for (const structure of model.structures) {
        typeInfo.types.set(structure.name, {
            name: structure.name,
            isStruct: true,
            needsPointer: true,
        });
    }

    // First pass - process all type aliases to find union types
    for (const typeAlias of model.typeAliases) {
        if (typeAlias.type.kind === "or") {
            // This is a union type - store the alias mapping
            const resolvedType = resolveType(typeAlias.type);
            typeInfo.unionTypeAliases.set(resolvedType.name, typeAlias.name);
        }
    }

    // Second pass - now process all type aliases with the union mappings in place
    for (const typeAlias of model.typeAliases) {
        const resolvedType = resolveType(typeAlias.type);
        typeInfo.types.set(typeAlias.name, {
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

/**
 * Generate types file
 * @returns {string} The generated code
 */
function generateTypesFile() {
    /** @type {string[]} */
    const parts = [];
    /** @type {(s: string) => unknown} */
    const write = s => parts.push(s);
    const writeLine = (s = "") => parts.push(s + "\n");

    // File header
    writeLine("// Code generated by generate2.mjs; DO NOT EDIT.");
    writeLine("");
    writeLine("package lsproto");
    writeLine("");
    // NullType needs the imports
    writeLine(`import "fmt"`);
    writeLine("");
    writeLine("// Meta model version " + model.metaData.version);
    writeLine("");

    // Write basic NullType for null values
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

    // Generate structures
    writeLine("// Structures\n");

    // Keep track of generated types
    const generatedTypes = new Set();

    for (const structure of model.structures) {
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

            writeLine(`\t${titleCase(prop.name)} ${goType} \`json:"${prop.name}${prop.optional ? ",omitzero" : ""}"\``);
            writeLine("");
        }

        writeLine("}");
        writeLine("");

        generatedTypes.add(structure.name);
    }

    // Generate type aliases
    writeLine("// Type aliases\n");

    for (const typeAlias of model.typeAliases) {
        write(formatDocumentation(typeAlias.documentation));
        write(formatDeprecation(typeAlias.deprecated));

        if (typeAlias.name === "LSPAny") {
            writeLine("type LSPAny any");
            writeLine("");
            continue;
        }

        const resolvedType = resolveType(typeAlias.type);
        writeLine(`type ${typeAlias.name} = ${resolvedType.name}`);
        writeLine("");

        typeInfo.generatedTypes.add(typeAlias.name);
    }

    return parts.join("");
}

/**
 * Generate enums file
 * @returns {string} The generated code
 */
function generateEnumsFile() {
    /** @type {string[]} */
    const parts = [];
    /** @type {(s: string) => unknown} */
    const write = s => parts.push(s);
    const writeLine = (s = "") => parts.push(s + "\n");

    // File header
    writeLine("// Code generated by generate2.mjs; DO NOT EDIT.");
    writeLine("");
    writeLine("package lsproto");
    writeLine("");
    writeLine(`import (`);
    writeLine(`\t"encoding/json"`);
    writeLine(`)`);
    writeLine("");
    writeLine("// Meta model version " + model.metaData.version);
    writeLine("");

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
        const enumValues = typeInfo.enumValuesByType.get(enumeration.name);
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

        typeInfo.generatedTypes.add(enumeration.name);
    }

    return parts.join("");
}

/**
 * Generate unions file
 * @returns {string} The generated code
 */
function generateUnionsFile() {
    /** @type {string[]} */
    const parts = [];
    /** @type {(s: string) => unknown} */
    const write = s => parts.push(s);
    const writeLine = (s = "") => parts.push(s + "\n");

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

    // Generate union types
    writeLine("// Union types\n");

    for (const [name, members] of typeInfo.unionTypes.entries()) {
        // Skip if already generated
        if (typeInfo.generatedTypes.has(name)) {
            continue;
        }

        writeLine(`type ${name} struct {`);

        // Use a Map to deduplicate by type name to ensure we don't include multiple fields with the same type
        const uniqueTypeFields = new Map(); // Maps type name -> field name

        for (const member of members) {
            let memberType;
            if (member.types.length === 1) {
                const type = resolveType(member.types[0]);
                memberType = type.name;

                // If this type name already exists in our map, skip it
                if (!uniqueTypeFields.has(memberType)) {
                    const fieldName = titleCase(member.name);
                    uniqueTypeFields.set(memberType, fieldName);
                    writeLine(`\t${fieldName} *${memberType}`);
                }
            }
            else {
                // This shouldn't happen with our current approach, but handle it just in case
                memberType = "any";
                const fieldName = titleCase(member.name);
                uniqueTypeFields.set(memberType, fieldName);
                writeLine(`\t${fieldName} *${memberType}`);
            }
        }

        writeLine(`}`);
        writeLine("");

        // Get the field names and types for marshal/unmarshal methods
        const fieldEntries = Array.from(uniqueTypeFields.entries()).map(([typeName, fieldName]) => ({ fieldName, typeName }));

        // Marshal method
        writeLine(`func (o ${name}) MarshalJSON() ([]byte, error) {`);

        // Create assertion to ensure only one field is set at a time
        write(`\tassertOnlyOne("more than one element of ${name} is set", `);

        // Write the assertion conditions
        for (let i = 0; i < fieldEntries.length; i++) {
            if (i > 0) write(", ");
            write(`o.${fieldEntries[i].fieldName} != nil`);
        }
        writeLine(`)`);
        writeLine("");

        // Write the marshal logic for each field
        for (const entry of fieldEntries) {
            writeLine(`\tif o.${entry.fieldName} != nil {`);
            writeLine(`\t\treturn json.Marshal(*o.${entry.fieldName})`);
            writeLine(`\t}`);
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
        for (let i = 0; i < fieldEntries.length; i++) {
            const entry = fieldEntries[i];
            writeLine(`\t{`);
            writeLine(`\t\tvar v ${entry.typeName}`);
            writeLine(`\t\tif err := json.Unmarshal(data, &v); err == nil {`);
            writeLine(`\t\t\to.${entry.fieldName} = &v`);
            writeLine(`\t\t\treturn nil`);
            writeLine(`\t\t}`);
            writeLine(`\t}`);
        }

        writeLine(`\treturn fmt.Errorf("cannot unmarshal %s into ${name}", string(data))`);
        writeLine(`}`);
        writeLine("");

        typeInfo.generatedTypes.add(name);
    }

    // Generate literal types
    writeLine("// Literal types\n");

    for (const [value, name] of typeInfo.literalTypes.entries()) {
        // Skip if already generated
        if (typeInfo.generatedTypes.has(name)) {
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

        typeInfo.generatedTypes.add(name);
    }

    return parts.join("");
}

/**
 * Generate methods file
 * @returns {string} The generated code
 */
function generateMethodsFile() {
    /** @type {string[]} */
    const parts = [];
    /** @type {(s: string) => unknown} */
    const write = s => parts.push(s);
    const writeLine = (s = "") => parts.push(s + "\n");

    // File header
    writeLine("// Code generated by generate2.mjs; DO NOT EDIT.");
    writeLine("");
    writeLine("package lsproto");
    writeLine("");

    // No imports needed for methods file

    writeLine("// Meta model version " + model.metaData.version);
    writeLine("");

    // Generate Methods for requests and notifications
    writeLine("// Methods\n");

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
            typeName = "any";
        }
        else if (request.params.kind === "reference") {
            typeName = request.params.name;
        }
        else {
            const resolvedType = resolveType(request.params);
            typeName = resolvedType.name;
        }

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
            typeName = "any";
        }
        else if (notification.params.kind === "reference") {
            typeName = notification.params.name;
        }
        else {
            const resolvedType = resolveType(notification.params);
            typeName = resolvedType.name;
        }

        writeLine(`\tMethod${methodName}: unmarshallerFor[${typeName}],`);
    }

    writeLine("}");

    return parts.join("");
}

/**
 * Helper function for asserting only one field is set
 * Since assertOnlyOne already exists in lsp.go, we don't need to generate it
 */
function generateHelperFunctions() {
    return ""; // Return empty string since we don't need to generate anything
}

/**
 * Main function
 */
function main() {
    try {
        collectTypeDefinitions();

        // Remove old generated file if it exists
        const oldGeneratedFilePath = path.resolve(outputDir, "lsp_generated.go");
        if (fs.existsSync(oldGeneratedFilePath)) {
            fs.unlinkSync(oldGeneratedFilePath);
            console.log(`Removed old generated file: ${oldGeneratedFilePath}`);
        }

        // Generate types file
        const generatedTypesCode = generateTypesFile();
        fs.writeFileSync(typesFilePath, generatedTypesCode);

        // Generate enums file
        const generatedEnumsCode = generateEnumsFile();
        fs.writeFileSync(enumsFilePath, generatedEnumsCode);

        // Generate unions file - need to add the helper function
        let generatedUnionsCode = generateUnionsFile();
        // Add the helper function at the beginning of the file after the imports
        const helperCode = generateHelperFunctions();
        const importEndIndex = generatedUnionsCode.indexOf(")\n\n");
        if (importEndIndex !== -1) {
            // Insert after the imports closing parenthesis and the blank line
            generatedUnionsCode = generatedUnionsCode.substring(0, importEndIndex + 3) +
                helperCode +
                generatedUnionsCode.substring(importEndIndex + 3);
        }
        fs.writeFileSync(unionsFilePath, generatedUnionsCode);

        // Generate methods file
        const generatedMethodsCode = generateMethodsFile();
        fs.writeFileSync(methodsFilePath, generatedMethodsCode);

        // Format with gofmt
        const gofmt = which.sync("gofmt");
        cp.execFileSync(gofmt, ["-w", typesFilePath]);
        cp.execFileSync(gofmt, ["-w", enumsFilePath]);
        cp.execFileSync(gofmt, ["-w", unionsFilePath]);
        cp.execFileSync(gofmt, ["-w", methodsFilePath]);

        console.log(`Successfully generated ${typesFilePath}, ${enumsFilePath}, ${unionsFilePath}, and ${methodsFilePath}`);
    }
    catch (error) {
        console.error("Error generating code:", error);
        process.exit(1);
    }
}

main();
