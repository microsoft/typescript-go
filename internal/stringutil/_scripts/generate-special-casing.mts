#!/usr/bin/env -S node --experimental-strip-types --no-warnings

import * as fs from "fs";
import * as path from "path";

const OUTPUT_PATH = path.join(import.meta.dirname, "..", "js_case_generated.go");
const SPECIAL_CASING_URL = "https://www.unicode.org/Public/UCD/latest/ucd/SpecialCasing.txt";
const DERIVED_CORE_PROPERTIES_URL = "https://www.unicode.org/Public/UCD/latest/ucd/DerivedCoreProperties.txt";

const knownContextConditions = new Set([
    "Final_Sigma",
    "After_Soft_Dotted",
    "More_Above",
    "After_I",
    "Not_Before_Dot",
]);

const knownLocaleConditions = new Set([
    "az",
    "lt",
    "tr",
]);

type SpecialCasingEntry = {
    codePoint: number;
    lower: number[];
    upper: number[];
    condition: string;
    comment: string;
};

type Range = {
    start: number;
    end: number;
};

function assert(condition: unknown, message: string): asserts condition {
    if (!condition) {
        throw new Error(message);
    }
}

async function fetchText(url: string): Promise<string> {
    const response = await fetch(url);
    if (!response.ok) {
        throw new Error(`Failed to fetch ${url}: ${response.status} ${response.statusText}`);
    }
    return await response.text();
}

function parseCodePointList(field: string): number[] {
    const trimmed = field.trim();
    if (!trimmed) return [];
    return trimmed.split(/\s+/).map(codePoint => parseInt(codePoint, 16));
}

function parseRange(field: string): Range {
    const [startHex, endHex] = field.split("..");
    const start = parseInt(startHex, 16);
    const end = endHex ? parseInt(endHex, 16) : start;
    return { start, end };
}

function goRuneLiteral(codePoint: number): string {
    return `0x${codePoint.toString(16).toUpperCase()}`;
}

function goStringLiteral(codePoints: number[]): string {
    let text = '"';
    for (const codePoint of codePoints) {
        if (codePoint <= 0xFFFF) {
            text += `\\u${codePoint.toString(16).toUpperCase().padStart(4, "0")}`;
        }
        else {
            text += `\\U${codePoint.toString(16).toUpperCase().padStart(8, "0")}`;
        }
    }
    text += '"';
    return text;
}

function parseSpecialCasing(text: string): { unicodeVersion: string; entries: SpecialCasingEntry[]; } {
    const entries: SpecialCasingEntry[] = [];
    let unicodeVersion = "unknown";

    for (const line of text.split(/\r?\n/)) {
        const versionMatch = line.match(/^# SpecialCasing-(.+)\.txt$/);
        if (versionMatch) {
            unicodeVersion = versionMatch[1];
            continue;
        }

        const trimmed = line.trim();
        if (!trimmed || trimmed.startsWith("#")) {
            continue;
        }

        const [data, comment = ""] = line.split("#", 2);
        const parts = data.split(";").map(part => part.trim());
        assert(parts.length >= 4, `Malformed SpecialCasing row: ${line}`);

        const [codeField, lowerField, _titleField, upperField, conditionField = ""] = parts;
        const code = parseCodePointList(codeField);
        assert(code.length === 1, `Expected single code point in SpecialCasing row: ${line}`);

        let hasLocaleCondition = false;
        let condition = "specialCasingConditionNone";
        let sawContextCondition = false;

        for (const token of conditionField.split(/\s+/).filter(Boolean)) {
            if (knownContextConditions.has(token)) {
                sawContextCondition = true;
                if (token === "Final_Sigma") {
                    condition = "specialCasingConditionFinalSigma";
                }
                continue;
            }
            if (knownLocaleConditions.has(token.toLowerCase())) {
                hasLocaleCondition = true;
                continue;
            }
            throw new Error(`Unknown SpecialCasing condition token: ${token}`);
        }

        if (hasLocaleCondition) {
            continue;
        }
        if (sawContextCondition && condition === "specialCasingConditionNone") {
            throw new Error(`Unsupported locale-insensitive context-only SpecialCasing row: ${line}`);
        }

        entries.push({
            codePoint: code[0],
            lower: parseCodePointList(lowerField),
            upper: parseCodePointList(upperField),
            condition,
            comment: comment.trim(),
        });
    }

    return { unicodeVersion, entries };
}

function parseDerivedCorePropertyRanges(text: string, propertyName: string): Range[] {
    const ranges: Range[] = [];

    for (const line of text.split(/\r?\n/)) {
        const trimmed = line.trim();
        if (!trimmed || trimmed.startsWith("#")) {
            continue;
        }

        const [data] = line.split("#", 1);
        const parts = data.split(";").map(part => part.trim());
        if (parts.length < 2 || parts[1] !== propertyName) {
            continue;
        }

        ranges.push(parseRange(parts[0]));
    }

    return ranges;
}

function renderRanges(name: string, ranges: Range[]): string {
    const values = ranges.flatMap(range => [goRuneLiteral(range.start), goRuneLiteral(range.end)]).join(", ");
    return `var ${name} = []rune{${values}}\n`;
}

function render(unicodeVersion: string, entries: SpecialCasingEntry[], casedRanges: Range[], caseIgnorableRanges: Range[]): string {
    const mappings = entries.map(entry => `\t${goRuneLiteral(entry.codePoint)}: {lower: ${goStringLiteral(entry.lower)}, upper: ${goStringLiteral(entry.upper)}, condition: ${entry.condition}}, // ${entry.comment}`).join("\n");

    return `// Code generated by internal/stringutil/_scripts/generate-special-casing.mts. DO NOT EDIT.
// Based on Unicode SpecialCasing.txt and DerivedCoreProperties.txt (${unicodeVersion}).
// Includes only the locale-insensitive mappings needed for ECMAScript default casing.
// Go's unicode package handles simple one-rune mappings, but not these multi-rune
// mappings or the Cased/Case_Ignorable properties needed for Final_Sigma handling.

package stringutil

type specialCasingCondition uint8

const (
\tspecialCasingConditionNone specialCasingCondition = iota
\tspecialCasingConditionFinalSigma
)

type specialCasingMapping struct {
\tlower     string
\tupper     string
\tcondition specialCasingCondition
}

var specialCasingMappings = map[rune]specialCasingMapping{
${mappings}
}

${renderRanges("unicodeCasedRanges", casedRanges)}
${renderRanges("unicodeCaseIgnorableRanges", caseIgnorableRanges)}
`;
}

async function main() {
    const [specialCasingText, derivedCorePropertiesText] = await Promise.all([
        fetchText(SPECIAL_CASING_URL),
        fetchText(DERIVED_CORE_PROPERTIES_URL),
    ]);

    const { unicodeVersion, entries } = parseSpecialCasing(specialCasingText);
    const casedRanges = parseDerivedCorePropertyRanges(derivedCorePropertiesText, "Cased");
    const caseIgnorableRanges = parseDerivedCorePropertyRanges(derivedCorePropertiesText, "Case_Ignorable");
    fs.writeFileSync(OUTPUT_PATH, render(unicodeVersion, entries, casedRanges, caseIgnorableRanges));
}

await main();
