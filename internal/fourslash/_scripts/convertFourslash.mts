import * as cp from "child_process";
import * as fs from "fs";
import * as path from "path";
import * as ts from "typescript";
import which from "which";

const stradaFourslashPath = path.resolve(import.meta.dirname, "../", "../", "../", "_submodules", "TypeScript", "tests", "cases", "fourslash");

let inputFileSet: Set<string> | undefined;

const failingTestsPath = path.join(import.meta.dirname, "failingTests.txt");
const failingTestsList = fs.readFileSync(failingTestsPath, "utf-8").split("\n").map(line => line.trim().substring(4)).filter(line => line.length > 0);
const failingTests = new Set(failingTestsList);
const helperFilePath = path.join(import.meta.dirname, "../", "tests", "util_test.go");

const outputDir = path.join(import.meta.dirname, "../", "tests", "gen");

const unparsedFiles: string[] = [];

function main() {
    const args = process.argv.slice(2);
    const inputFilesPath = args[0];
    if (inputFilesPath) {
        const inputFiles = fs.readFileSync(inputFilesPath, "utf-8")
            .split("\n").map(line => line.trim())
            .filter(line => line.length > 0)
            .map(line => path.basename(line));
        inputFileSet = new Set(inputFiles);
    }

    if (!fs.existsSync(outputDir)) {
        fs.mkdirSync(outputDir, { recursive: true });
    }

    generateHelperFile();
    parseTypeScriptFiles(stradaFourslashPath);
    console.log(unparsedFiles.join("\n"));
    const gofmt = which.sync("go");
    cp.execFileSync(gofmt, ["tool", "mvdan.cc/gofumpt", "-lang=go1.24", "-w", outputDir]);
}

function parseTypeScriptFiles(folder: string): void {
    const files = fs.readdirSync(folder);

    files.forEach(file => {
        const filePath = path.join(folder, file);
        const stat = fs.statSync(filePath);
        if (inputFileSet && !inputFileSet.has(file)) {
            return;
        }

        if (stat.isDirectory()) {
            parseTypeScriptFiles(filePath);
        }
        else if (file.endsWith(".ts")) {
            const content = fs.readFileSync(filePath, "utf-8");
            const test = parseFileContent(file, content);
            if (test) {
                const testContent = generateGoTest(test);
                const testPath = path.join(outputDir, `${test.name}_test.go`);
                fs.writeFileSync(testPath, testContent, "utf-8");
            }
        }
    });
}

function parseFileContent(filename: string, content: string): GoTest | undefined {
    console.error(`Parsing file: ${filename}`);
    const sourceFile = ts.createSourceFile("temp.ts", content, ts.ScriptTarget.Latest, true /*setParentNodes*/);
    const statements = sourceFile.statements;
    const goTest: GoTest = {
        name: filename.replace(".ts", ""),
        content: getTestInput(content),
        commands: [],
    };
    for (const statement of statements) {
        const result = parseFourslashStatement(statement);
        if (!result) {
            unparsedFiles.push(filename);
            return undefined;
        }
        else {
            goTest.commands.push(...result);
        }
    }
    return goTest;
}

function getTestInput(content: string): string {
    const lines = content.split("\n");
    let testInput: string[] = [];
    for (const line of lines) {
        let newLine = "";
        if (line.startsWith("////")) {
            const parts = line.substring(4).split("`");
            for (let i = 0; i < parts.length; i++) {
                if (i > 0) {
                    newLine += `\` + "\`" + \``;
                }
                newLine += parts[i];
            }
            testInput.push(newLine);
        }
        else if (line.startsWith("// @") || line.startsWith("//@")) {
            testInput.push(line);
        }
        // !!! preserve non-input comments?
    }

    // chomp leading spaces
    if (!testInput.some(line => line.length != 0 && !line.startsWith(" ") && !line.startsWith("// "))) {
        testInput = testInput.map(line => {
            if (line.startsWith(" ")) return line.substring(1);
            return line;
        });
    }
    return `\`${testInput.join("\n")}\``;
}

/**
 * Parses a Strada fourslash statement and returns the corresponding Corsa commands.
 * @returns an array of commands if the statement is a valid fourslash command, or `false` if the statement could not be parsed.
 */
function parseFourslashStatement(statement: ts.Statement): Cmd[] | undefined {
    if (ts.isVariableStatement(statement)) {
        // variable declarations (for ranges and markers), e.g. `const range = test.ranges()[0];`
        return [];
    }
    else if (ts.isExpressionStatement(statement) && ts.isCallExpression(statement.expression)) {
        const callExpression = statement.expression;
        if (!ts.isPropertyAccessExpression(callExpression.expression)) {
            console.error(`Expected property access expression, got ${callExpression.expression.getText()}`);
            return undefined;
        }
        const namespace = callExpression.expression.expression;
        const func = callExpression.expression.name;
        if (!ts.isIdentifier(namespace) || !ts.isIdentifier(func)) {
            console.error(`Expected identifiers for namespace and function, got ${namespace.getText()} and ${func.getText()}`);
            return undefined;
        }
        // `verify.completions(...)`
        if (namespace.text === "verify" && func.text === "completions") {
            return parseVerifyCompletionsArgs(callExpression.arguments);
        }
        // `goTo....`
        if (namespace.text === "goTo") {
            return parseGoToArgs(callExpression.arguments, func.text);
        }
        // `edit....`
        if (namespace.text === "edit") {
            const result = parseEditStatement(func.text, callExpression.arguments);
            if (!result) {
                return undefined;
            }
            return [result];
        }
        // !!! other fourslash commands
    }
    console.error(`Unrecognized fourslash statement: ${statement.getText()}`);
    return undefined;
}

function parseEditStatement(funcName: string, args: readonly ts.Expression[]): EditCmd | undefined {
    switch (funcName) {
        case "insert":
        case "paste":
        case "insertLine":
            if (args.length !== 1 || !ts.isStringLiteralLike(args[0])) {
                console.error(`Expected a single string literal argument in edit.${funcName}, got ${args.map(arg => arg.getText()).join(", ")}`);
                return undefined;
            }
            return {
                kind: "edit",
                goStatement: `f.${funcName.charAt(0).toUpperCase() + funcName.slice(1)}(t, ${getGoStringLiteral(args[0].text)})`,
            };
        case "replaceLine":
            if (args.length !== 2 || !ts.isNumericLiteral(args[0]) || !ts.isStringLiteral(args[1])) {
                console.error(`Expected a single string literal argument in edit.insert, got ${args.map(arg => arg.getText()).join(", ")}`);
                return undefined;
            }
            return {
                kind: "edit",
                goStatement: `f.ReplaceLine(t, ${args[0].text}, ${getGoStringLiteral(args[1].text)})`,
            };
        case "backspace":
            const arg = args[0];
            if (arg) {
                if (!ts.isNumericLiteral(arg)) {
                    console.error(`Expected numeric literal argument in edit.backspace, got ${arg.getText()}`);
                    return undefined;
                }
                return {
                    kind: "edit",
                    goStatement: `f.Backspace(t, ${arg.text})`,
                };
            }
            return {
                kind: "edit",
                goStatement: `f.Backspace(t, 1)`,
            };
        default:
            console.error(`Unrecognized edit function: ${funcName}`);
            return undefined;
    }
}

function getGoStringLiteral(text: string): string {
    return `${JSON.stringify(text)}`;
}

function parseGoToArgs(args: readonly ts.Expression[], funcName: string): GoToCmd[] | undefined {
    switch (funcName) {
        case "marker":
            const arg = args[0];
            if (arg === undefined) {
                return [{
                    kind: "goTo",
                    funcName: "marker",
                    args: [`""`],
                }];
            }
            if (!ts.isStringLiteral(arg)) {
                console.error(`Unrecognized argument in goTo.marker: ${arg.getText()}`);
                return undefined;
            }
            return [{
                kind: "goTo",
                funcName: "marker",
                args: [getGoStringLiteral(arg.text)],
            }];
        case "file":
            if (args.length !== 1) {
                console.error(`Expected a single argument in goTo.file, got ${args.map(arg => arg.getText()).join(", ")}`);
                return undefined;
            }
            if (ts.isStringLiteral(args[0])) {
                return [{
                    kind: "goTo",
                    funcName: "file",
                    args: [getGoStringLiteral(args[0].text)],
                }];
            }
            else if (ts.isNumericLiteral(args[0])) {
                return [{
                    kind: "goTo",
                    funcName: "fileNumber",
                    args: [args[0].text],
                }];
            }
            console.error(`Expected string or number literal argument in goTo.file, got ${args[0].getText()}`);
            return undefined;
        case "position":
            if (args.length !== 1 || !ts.isNumericLiteral(args[0])) {
                console.error(`Expected a single numeric literal argument in goTo.position, got ${args.map(arg => arg.getText()).join(", ")}`);
                return undefined;
            }
            return [{
                kind: "goTo",
                funcName: "position",
                args: [`${args[0].text}`],
            }];
        case "eof":
            return [{
                kind: "goTo",
                funcName: "EOF",
                args: [],
            }];
        case "bof":
            return [{
                kind: "goTo",
                funcName: "BOF",
                args: [],
            }];
        case "select":
            if (args.length !== 2 || !ts.isStringLiteral(args[0]) || !ts.isStringLiteral(args[1])) {
                console.error(`Expected two string literal arguments in goTo.select, got ${args.map(arg => arg.getText()).join(", ")}`);
                return undefined;
            }
            return [{
                kind: "goTo",
                funcName: "select",
                args: [getGoStringLiteral(args[0].text), getGoStringLiteral(args[1].text)],
            }];
        default:
            console.error(`Unrecognized goTo function: ${funcName}`);
            return undefined;
    }
}

function parseVerifyCompletionsArgs(args: readonly ts.Expression[]): VerifyCompletionsCmd[] | undefined {
    const cmds = [];
    for (const arg of args) {
        const result = parseVerifyCompletionArg(arg);
        if (!result) {
            return undefined;
        }
        cmds.push(result);
    }
    return cmds;
}

const completionConstants = new Map([
    ["completion.globals", "completionGlobals"],
    ["completion.globalTypes", "completionGlobalTypes"],
    ["completion.classElementKeywords", "completionClassElementKeywords"],
    ["completion.classElementInJsKeywords", "completionClassElementInJSKeywords"],
    ["completion.constructorParameterKeywords", "completionConstructorParameterKeywords"],
    ["completion.functionMembersWithPrototype", "completionFunctionMembersWithPrototype"],
    ["completion.functionMembers", "completionFunctionMembers"],
    ["completion.typeKeywords", "completionTypeKeywords"],
    ["completion.undefinedVarEntry", "completionUndefinedVarItem"],
    ["completion.typeAssertionKeywords", "completionTypeAssertionKeywords"],
]);

const completionPlus = new Map([
    ["completion.globalsPlus", "completionGlobalsPlus"],
    ["completion.globalTypesPlus", "completionGlobalTypesPlus"],
    ["completion.functionMembersPlus", "completionFunctionMembersPlus"],
    ["completion.functionMembersWithPrototypePlus", "completionFunctionMembersWithPrototypePlus"],
    ["completion.globalsInJsPlus", "completionGlobalsInJSPlus"],
    ["completion.typeKeywordsPlus", "completionTypeKeywordsPlus"],
]);

function parseVerifyCompletionArg(arg: ts.Expression): VerifyCompletionsCmd | undefined {
    let marker: string | undefined;
    let goArgs: VerifyCompletionsArgs | undefined;
    if (!ts.isObjectLiteralExpression(arg)) {
        console.error(`Expected object literal expression in verify.completions, got ${arg.getText()}`);
        return undefined;
    }
    let isNewIdentifierLocation: true | undefined;
    for (const prop of arg.properties) {
        if (!ts.isPropertyAssignment(prop) || !ts.isIdentifier(prop.name)) {
            console.error(`Expected property assignment with identifier name, got ${prop.getText()}`);
            return undefined;
        }
        const propName = prop.name.text;
        const init = prop.initializer;
        switch (propName) {
            case "marker":
                if (ts.isStringLiteral(init)) {
                    marker = getGoStringLiteral(init.text);
                }
                else if (ts.isArrayLiteralExpression(init)) {
                    marker = "[]string{";
                    for (const elem of init.elements) {
                        if (!ts.isStringLiteral(elem)) {
                            console.error(`Expected string literal in marker array, got ${elem.getText()}`);
                            return undefined; // !!! parse marker arrays?
                        }
                        marker += `${getGoStringLiteral(elem.text)}, `;
                    }
                    marker += "}";
                }
                else if (ts.isObjectLiteralExpression(init)) {
                    // !!! parse marker objects?
                    console.error(`Unrecognized marker initializer: ${init.getText()}`);
                    return undefined;
                }
                else if (init.getText() === "test.markers()") {
                    marker = "f.Markers()";
                }
                else {
                    console.error(`Unrecognized marker initializer: ${init.getText()}`);
                    return undefined;
                }
                break;
            case "exact":
            case "includes":
            case "unsorted":
                if (init.getText() === "undefined") {
                    return {
                        kind: "verifyCompletions",
                        marker: marker ? marker : "nil",
                        args: undefined,
                    };
                }
                let expected: string;
                const initText = init.getText();
                if (completionConstants.has(initText)) {
                    expected = completionConstants.get(initText)!;
                }
                else if (completionPlus.keys().some(funcName => initText.startsWith(funcName))) {
                    const tsFunc = completionPlus.keys().find(funcName => initText.startsWith(funcName));
                    const funcName = completionPlus.get(tsFunc!)!;
                    const items = (init as ts.CallExpression).arguments[0];
                    const opts = (init as ts.CallExpression).arguments[1];
                    if (!ts.isArrayLiteralExpression(items)) {
                        console.error(`Expected array literal expression for completion.globalsPlus items, got ${items.getText()}`);
                        return undefined;
                    }
                    expected = `${funcName}([]fourslash.CompletionsExpectedItem{`;
                    for (const elem of items.elements) {
                        const result = parseExpectedCompletionItem(elem);
                        if (!result) {
                            return undefined;
                        }
                        expected += result + ", ";
                    }
                    expected += "}";
                    if (opts) {
                        if (!ts.isObjectLiteralExpression(opts)) {
                            console.error(`Expected object literal expression for completion.globalsPlus options, got ${opts.getText()}`);
                            return undefined;
                        }
                        const noLib = opts.properties[0];
                        if (noLib && ts.isPropertyAssignment(noLib) && noLib.name.getText() === "noLib") {
                            if (noLib.initializer.kind === ts.SyntaxKind.TrueKeyword) {
                                expected += ", true";
                            }
                            else if (noLib.initializer.kind === ts.SyntaxKind.FalseKeyword) {
                                expected += ", false";
                            }
                            else {
                                console.error(`Expected boolean literal for noLib, got ${noLib.initializer.getText()}`);
                                return undefined;
                            }
                        }
                        else {
                            console.error(`Expected noLib property in completion.globalsPlus options, got ${opts.getText()}`);
                            return undefined;
                        }
                    }
                    else if (tsFunc === "completion.globalsPlus" || tsFunc === "completion.globalsInJsPlus") {
                        expected += ", false"; // Default for noLib
                    }
                    expected += ")";
                }
                else {
                    expected = "[]fourslash.CompletionsExpectedItem{";
                    if (ts.isArrayLiteralExpression(init)) {
                        for (const elem of init.elements) {
                            const result = parseExpectedCompletionItem(elem);
                            if (!result) {
                                return undefined;
                            }
                            expected += result + ", ";
                        }
                    }
                    else {
                        const result = parseExpectedCompletionItem(init);
                        if (!result) {
                            return undefined;
                        }
                        expected += result;
                    }
                    expected += "}";
                }
                if (propName === "includes") {
                    (goArgs ??= {}).includes = expected;
                }
                else if (propName === "exact") {
                    (goArgs ??= {}).exact = expected;
                }
                else {
                    (goArgs ??= {}).unsorted = expected;
                }
                break;
            case "excludes":
                let excludes = "[]string{";
                if (ts.isStringLiteral(init)) {
                    excludes += `${getGoStringLiteral(init.text)}, `;
                }
                else if (ts.isArrayLiteralExpression(init)) {
                    for (const elem of init.elements) {
                        if (!ts.isStringLiteral(elem)) {
                            return undefined; // Shouldn't happen
                        }
                        excludes += `${getGoStringLiteral(elem.text)}, `;
                    }
                }
                excludes += "}";
                (goArgs ??= {}).excludes = excludes;
                break;
            case "isNewIdentifierLocation":
                if (init.kind === ts.SyntaxKind.TrueKeyword) {
                    isNewIdentifierLocation = true;
                }
                break;
            case "preferences":
            case "triggerCharacter":
            case "defaultCommitCharacters":
                break; // !!! parse once they're supported in fourslash
            case "optionalReplacementSpan": // the only two tests that use this will require manual conversion
            case "isGlobalCompletion":
                break; // Ignored, unused
            default:
                console.error(`Unrecognized expected completion item: ${init.parent.getText()}`);
                return undefined;
        }
    }
    return {
        kind: "verifyCompletions",
        marker: marker ? marker : "nil",
        args: goArgs,
        isNewIdentifierLocation: isNewIdentifierLocation,
    };
}

function parseExpectedCompletionItem(expr: ts.Expression): string | undefined {
    if (completionConstants.has(expr.getText())) {
        return completionConstants.get(expr.getText())!;
    }
    if (ts.isStringLiteral(expr)) {
        return getGoStringLiteral(expr.text);
    }
    if (ts.isObjectLiteralExpression(expr)) {
        let isDeprecated = false; // !!!
        let isOptional = false;
        let extensions: string[] = []; // !!!
        let item = "&lsproto.CompletionItem{";
        let name: string | undefined;
        let insertText: string | undefined;
        let filterText: string | undefined;
        for (const prop of expr.properties) {
            if (!ts.isPropertyAssignment(prop) || !ts.isIdentifier(prop.name)) {
                console.error(`Expected property assignment with identifier name for completion item, got ${prop.getText()}`);
                return undefined;
            }
            const propName = prop.name.text;
            const init = prop.initializer;
            switch (propName) {
                case "name":
                    if (ts.isStringLiteral(init)) {
                        name = init.text;
                    }
                    else {
                        console.error(`Expected string literal for completion item name, got ${init.getText()}`);
                        return undefined;
                    }
                    break;
                case "sortText":
                    const result = parseSortText(init);
                    if (!result) {
                        return undefined;
                    }
                    item += `SortText: ptrTo(string(${result})), `;
                    if (result === "ls.SortTextOptionalMember") {
                        isOptional = true;
                    }
                    break;
                case "insertText":
                    if (ts.isStringLiteral(init)) {
                        insertText = init.text;
                    }
                    else {
                        console.error(`Expected string literal for insertText, got ${init.getText()}`);
                        return undefined;
                    }
                    break;
                case "filterText":
                    if (ts.isStringLiteral(init)) {
                        filterText = init.text;
                    }
                    else {
                        console.error(`Expected string literal for filterText, got ${init.getText()}`);
                        return undefined;
                    }
                    break;
                case "isRecommended":
                    if (init.kind === ts.SyntaxKind.TrueKeyword) {
                        item += `Preselect: ptrTo(true), `;
                    }
                    break;
                case "kind":
                    const kind = parseKind(init);
                    if (!kind) {
                        return undefined;
                    }
                    item += `Kind: ptrTo(${kind}), `;
                    break;
                case "kindModifiers":
                    const modifiers = parseKindModifiers(init);
                    if (!modifiers) {
                        return undefined;
                    }
                    ({ isDeprecated, isOptional, extensions } = modifiers);
                    break;
                case "commitCharacters":
                case "replacementSpan":
                    // !!! support these later
                    break;
                default:
                    console.error(`Unrecognized property in expected completion item: ${propName}`);
                    return undefined; // Unsupported property
            }
        }
        if (!name) {
            return undefined; // Shouldn't happen
        }
        if (isOptional) {
            insertText ??= name;
            filterText ??= name;
            name += "?";
        }
        item += `Label: ${getGoStringLiteral(name!)}, `;
        if (insertText) item += `InsertText: ptrTo(${getGoStringLiteral(insertText)}), `;
        if (filterText) item += `FilterText: ptrTo(${getGoStringLiteral(filterText)}), `;
        item += "}";
        return item;
    }
    console.error(`Expected string literal or object literal for expected completion item, got ${expr.getText()}`);
    return undefined; // Unsupported expression type
}

function parseKind(expr: ts.Expression): string | undefined {
    if (!ts.isStringLiteral(expr)) {
        console.error(`Expected string literal for kind, got ${expr.getText()}`);
        return undefined;
    }
    switch (expr.text) {
        case "primitive type":
        case "keyword":
            return "lsproto.CompletionItemKindKeyword";
        case "const":
        case "let":
        case "var":
        case "local var":
        case "alias":
        case "parameter":
            return "lsproto.CompletionItemKindVariable";
        case "property":
        case "getter":
        case "setter":
            return "lsproto.CompletionItemKindField";
        case "function":
        case "local function":
            return "lsproto.CompletionItemKindFunction";
        case "method":
        case "construct":
        case "call":
        case "index":
            return "lsproto.CompletionItemKindMethod";
        case "enum":
            return "lsproto.CompletionItemKindEnum";
        case "enum member":
            return "lsproto.CompletionItemKindEnumMember";
        case "module":
        case "external module name":
            return "lsproto.CompletionItemKindModule";
        case "class":
        case "type":
            return "lsproto.CompletionItemKindClass";
        case "interface":
            return "lsproto.CompletionItemKindInterface";
        case "warning":
            return "lsproto.CompletionItemKindText";
        case "script":
            return "lsproto.CompletionItemKindFile";
        case "directory":
            return "lsproto.CompletionItemKindFolder";
        case "string":
            return "lsproto.CompletionItemKindConstant";
        default:
            return "lsproto.CompletionItemKindProperty";
    }
}

const fileKindModifiers = new Set([".d.ts", ".ts", ".tsx", ".js", ".jsx", ".json"]);

function parseKindModifiers(expr: ts.Expression): { isOptional: boolean; isDeprecated: boolean; extensions: string[]; } | undefined {
    if (!ts.isStringLiteral(expr)) {
        console.error(`Expected string literal for kind modifiers, got ${expr.getText()}`);
        return undefined;
    }
    let isOptional = false;
    let isDeprecated = false;
    const extensions: string[] = [];
    const modifiers = expr.text.split(",");
    for (const modifier of modifiers) {
        switch (modifier) {
            case "optional":
                isOptional = true;
                break;
            case "deprecated":
                isDeprecated = true;
                break;
            default:
                if (fileKindModifiers.has(modifier)) {
                    extensions.push(modifier);
                }
        }
    }
    return {
        isOptional,
        isDeprecated,
        extensions,
    };
}

function parseSortText(expr: ts.Expression): string | undefined {
    const text = expr.getText();
    switch (text) {
        case "completion.SortText.LocalDeclarationPriority":
            return "ls.SortTextLocalDeclarationPriority";
        case "completion.SortText.LocationPriority":
            return "ls.SortTextLocationPriority";
        case "completion.SortText.OptionalMember":
            return "ls.SortTextOptionalMember";
        case "completion.SortText.MemberDeclaredBySpreadAssignment":
            return "ls.SortTextMemberDeclaredBySpreadAssignment";
        case "completion.SortText.SuggestedClassMembers":
            return "ls.SortTextSuggestedClassMembers";
        case "completion.SortText.GlobalsOrKeywords":
            return "ls.SortTextGlobalsOrKeywords";
        case "completion.SortText.AutoImportSuggestions":
            return "ls.SortTextAutoImportSuggestions";
        case "completion.SortText.ClassMemberSnippets":
            return "ls.SortTextClassMemberSnippets";
        case "completion.SortText.JavascriptIdentifiers":
            return "ls.SortTextJavascriptIdentifiers";
        default:
            console.error(`Unrecognized sort text: ${text}`);
            return undefined; // !!! support deprecated/obj literal prop/etc
    }
}

interface VerifyCompletionsCmd {
    kind: "verifyCompletions";
    marker: string;
    isNewIdentifierLocation?: true;
    args?: VerifyCompletionsArgs;
}

interface VerifyCompletionsArgs {
    includes?: string;
    excludes?: string;
    exact?: string;
    unsorted?: string;
}

interface GoToCmd {
    kind: "goTo";
    // !!! `selectRange` and `rangeStart` require parsing variables and `test.ranges()[n]`
    funcName: "marker" | "file" | "fileNumber" | "EOF" | "BOF" | "position" | "select";
    args: string[];
}

interface EditCmd {
    kind: "edit";
    goStatement: string;
}

type Cmd = VerifyCompletionsCmd | GoToCmd | EditCmd;

function generateVerifyCompletions({ marker, args, isNewIdentifierLocation }: VerifyCompletionsCmd): string {
    let expectedList = "nil";
    if (args) {
        const expected = [];
        if (args.includes) expected.push(`Includes: ${args.includes},`);
        if (args.excludes) expected.push(`Excludes: ${args.excludes},`);
        if (args.exact) expected.push(`Exact: ${args.exact},`);
        if (args.unsorted) expected.push(`Unsorted: ${args.unsorted},`);
        // !!! isIncomplete
        const commitCharacters = isNewIdentifierLocation ? "[]string{}" : "defaultCommitCharacters";
        expectedList = `&fourslash.CompletionsExpectedList{
    IsIncomplete: false,
    ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
        CommitCharacters: &${commitCharacters},
        EditRange: ignored,
    },
    Items: &fourslash.CompletionsExpectedItems{
        ${expected.join("\n")}
    },
}`;
    }
    return `f.VerifyCompletions(t, ${marker}, ${expectedList})`;
}

function generateGoToCommand({ funcName, args }: GoToCmd): string {
    const funcNameCapitalized = funcName.charAt(0).toUpperCase() + funcName.slice(1);
    return `f.GoTo${funcNameCapitalized}(t, ${args.join(", ")})`;
}

function generateCmd(cmd: Cmd): string {
    switch (cmd.kind) {
        case "verifyCompletions":
            return generateVerifyCompletions(cmd as VerifyCompletionsCmd);
        case "goTo":
            return generateGoToCommand(cmd as GoToCmd);
        case "edit":
            return cmd.goStatement;
        default:
            throw new Error(`Unknown command kind: ${cmd}`);
    }
}

interface GoTest {
    name: string;
    content: string;
    commands: Cmd[];
}

function generateGoTest(test: GoTest): string {
    const testName = test.name[0].toUpperCase() + test.name.substring(1);
    const content = test.content;
    const commands = test.commands.map(cmd => generateCmd(cmd)).join("\n");
    const imports = [`"github.com/microsoft/typescript-go/internal/fourslash"`];
    // Only include these imports if the commands use them to avoid unused import errors.
    if (commands.includes("ls.")) {
        imports.push(`"github.com/microsoft/typescript-go/internal/ls"`);
    }
    if (commands.includes("lsproto.")) {
        imports.push(`"github.com/microsoft/typescript-go/internal/lsp/lsproto"`);
    }
    imports.push(`"github.com/microsoft/typescript-go/internal/testutil"`);
    const template = `package fourslash_test

import (
	"testing"

    ${imports.join("\n\t")}
)

func Test${testName}(t *testing.T) {
    t.Parallel()
    ${failingTests.has(testName) ? "t.Skip()" : ""}
    defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = ${content}
    f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
    ${commands}
}`;
    return template;
}

function generateHelperFile() {
    fs.copyFileSync(helperFilePath, path.join(outputDir, "util_test.go"));
}

main();
