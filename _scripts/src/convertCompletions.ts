import * as fs from 'fs';
import * as path from 'path';
import * as ts from 'typescript';

const folderPath = path.resolve(__dirname, '../', '../', '_submodules', 'TypeScript', 'tests', 'cases', 'fourslash');
const verifyCompletionsListPath = path.join(__dirname, '../', "verify-completions.txt");
const verifyFiles = fs.readFileSync(verifyCompletionsListPath, 'utf-8')
    .split('\n').map(line => line.trim())
    .filter(line => line.length > 0)
    .map(line => {
        line = line.replace('tests/cases/fourslash/server', '');
        return line.replace('tests/cases/fourslash/', '');
    });
const verifyFilesSet = new Set(verifyFiles);

const unparsedFiles: string[] = [];

function parseTypeScriptFiles(folder: string): void {
    const files = fs.readdirSync(folder);

    files.forEach((file) => {
        const filePath = path.join(folder, file);
        const stat = fs.statSync(filePath);
        if (!verifyFilesSet.has(file)) {
            return;
        }

        if (stat.isDirectory()) {
            parseTypeScriptFiles(filePath);
        } else if (file.endsWith('.ts')) {
            const content = fs.readFileSync(filePath, 'utf-8');
            const test = parseFileContent(file, content);
            if (test) {
                const testContent = generateGoTest(test)
                const testPath = path.join(__dirname, '../', '../', 'internal', 'fourslash', 'tests', 'gen', `${test.name}_test.go`);
                fs.writeFileSync(testPath, testContent, 'utf-8');
            }
        }
    });
}

function parseFileContent(filename: string, content: string): GoTest | undefined {
    const sourceFile = ts.createSourceFile('temp.ts', content, ts.ScriptTarget.Latest, true /*setParentNodes*/);
    const statements = sourceFile.statements;
    const goTest: GoTest = {
        name: filename.replace('.ts', ''),
        content: getTestInput(content),
        commands: [],
    }
    for (const statement of statements) {
        const result = parseFourslashStatement(statement);
        if (!result) {
            console.error(`Unrecognized statement in file: ${statement.getText()}`);
            unparsedFiles.push(filename);
            return undefined;
        } else if (typeof result == 'object') {
            goTest.commands.push(...result);
        }
    }
    return goTest;
}

function getTestInput(content: string): string {
    const lines = content.split('\n');
    let testInput: string[] = [];
    for (const line of lines) {
        let newLine = "";
        if (line.startsWith('////')) {
            const parts = line.substring(4).split('`');
            for (let i = 0; i < parts.length; i++) {
                if (i > 0) {
                    newLine += `\` + "\`" + \``;
                }
                newLine += parts[i];
            }
            testInput.push(newLine);
        } else if (line.startsWith('// @')) {
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
    return `\`${testInput.join('\n')}\``;
}

// !!! return a reason for the failure
function parseFourslashStatement(statement: ts.Statement): Cmd[] | boolean {
    if (ts.isVariableStatement(statement)) {
        // !!! variable declarations (for ranges and markers) or function calls
        return true;
    } else if (ts.isExpressionStatement(statement) && ts.isCallExpression(statement.expression)) {
        const callExpression = statement.expression;
        if (!ts.isPropertyAccessExpression(callExpression.expression)) {
            return false;
        }
        const namespace = callExpression.expression.expression;
        const func = callExpression.expression.name;
        if (!ts.isIdentifier(namespace) || !ts.isIdentifier(func)) {
            return false;
        }
        if (namespace.text === 'verify' && func.text === 'completions') {
            return parseVerifyCompletionsArgs(callExpression.arguments);
        }
        if (namespace.text === 'goTo' && func.text === 'marker') {
            // !!! parse args
            return true;
        }
        // !!! baseline completions
    }
    return false;
}

function getGoStringLiteral(text: string): string {
    // if (text.includes('`')) {
    //     return `"${text}"`;
    // }
    // return `\`${text}\``;
    return `${JSON.stringify(text)}`;
}

function parseVerifyCompletionsArgs(args: readonly ts.Expression[]): VerifyCompletionsCmd[] | false {
    const cmds = [];
    for (const arg of args) {
        const result = parseVerifyCompletionArg(arg);
        if (!result) {
            return false;
        }
        cmds.push(result);
    }
    return cmds;
}
function parseVerifyCompletionArg(arg: ts.Expression): VerifyCompletionsCmd | false {
    let marker: string | undefined;
    let goArgs: VerifyCompletionsArgs | undefined;
    if (!ts.isObjectLiteralExpression(arg)) {
        return false;
    }
    for (const prop of arg.properties) {
        if (!ts.isPropertyAssignment(prop) || !ts.isIdentifier(prop.name)) {
            return false;
        }
        const propName = prop.name.text;
        const init = prop.initializer;
        switch (propName) {
            case 'marker':
                if (ts.isStringLiteral(init)) {
                    marker = getGoStringLiteral(init.text);
                }
                else if (ts.isArrayLiteralExpression(init)) {
                    marker = "[]string{"
                    for (const elem of init.elements) {
                        if (!ts.isStringLiteral(elem)) {
                            return false; // !!! parse marker arrays?
                        }
                        marker += `${getGoStringLiteral(elem.text)}, `;
                    }
                    marker += '}';
                }
                else if (ts.isObjectLiteralExpression(init)) {
                    // !!! parse marker objects?
                    return false;
                }
                else if (init.getText() === 'test.markers()') {
                    marker = "f.Markers()";
                }
                else {
                    return false;
                }
                break;
            case 'exact':
            case 'includes':
                if (init.getText() === 'undefined') {
                    return {
                        kind: CommandKind.verifyCompletions,
                        marker: marker ? marker : "nil",
                        args: undefined,
                    };
                }
                let expected = "[]fourslash.ExpectedCompletionItem{";
                if (ts.isArrayLiteralExpression(init)) {
                    for (const elem of init.elements) {
                        const result = parseExpectedCompletionItem(elem);
                        if (!result) {
                            return false;
                        }
                        expected += result + ', ';
                    }
                } else {
                    const result = parseExpectedCompletionItem(init);
                    if (!result) {
                        return false;
                    }
                    expected += result;
                }
                expected += '}';
                if (propName === 'includes') {
                    (goArgs ??= {}).includes = expected;
                } else {
                    (goArgs ??= {}).exact = expected;
                }
                break; // !!! parse these args
            case 'excludes':
                let excludes = "[]string{";
                if (ts.isStringLiteral(init)) {
                    excludes += `${getGoStringLiteral(init.text)}, `;
                }
                else if (ts.isArrayLiteralExpression(init)) {
                    for (const elem of init.elements) {
                        if (!ts.isStringLiteral(elem)) {
                            return false; // Shouldn't happen
                        }
                        excludes += `${getGoStringLiteral(elem.text)}, `;
                    }
                }
                excludes += '}';
                (goArgs ??= {}).excludes = excludes;
                break;
            case 'isNewIdentifierLocation':
                break; // !!! parse this into item defaults/commit characters
            case 'preferences':
            case 'triggerCharacter':
            case 'defaultCommitCharacters':
                break; // !!! parse once they're supported
            case 'optionalReplacementSpan':
            case 'isGlobalCompletion':
                break; // Ignored
            default:
                return false;
        }
    }
    return {
        kind: CommandKind.verifyCompletions,
        marker: marker ? marker : "nil",
        args: goArgs,
    };
}

function parseExpectedCompletionItem(expr: ts.Expression): string | false {
    if (ts.isStringLiteral(expr)) {
        return getGoStringLiteral(expr.text);
    }
    if (ts.isObjectLiteralExpression(expr)) {
        let isDeprecated = false;
        let isOptional = false;
        let extensions: string[] = [];
        let item = "&lsproto.CompletionItem{";
        let name: string | undefined;
        let insertText: string | undefined;
        let filterText: string | undefined;
        for (const prop of expr.properties) {
            if (!ts.isPropertyAssignment(prop) || !ts.isIdentifier(prop.name)) {
                return false;
            }
            const propName = prop.name.text;
            const init = prop.initializer;
            switch (propName) {
                case 'name':
                    if (ts.isStringLiteral(init)) {
                        name = init.text;
                    } else {
                        return false;
                    }
                    break;
                case 'sortText':
                    const result = parseSortText(init);
                    if (!result) {
                        return false;
                    }
                    item += `SortText: fourslash.PtrTo(string(${result})), `;
                    break;
                case 'insertText':
                    if (ts.isStringLiteral(init)) {
                        insertText = init.text;
                    } else {
                        return false;
                    }
                    break;
                case 'filterText':
                    if (ts.isStringLiteral(init)) {
                        filterText = init.text;
                    } else {
                        return false;
                    }
                    break;
                case 'isRecommended':
                    if (init.kind === ts.SyntaxKind.TrueKeyword) {
                        item += `Preselect: fourslash.PtrTo(true), `;
                    }
                    break;
                case 'kind': 
                    const kind = parseKind(init);
                    if (!kind) {
                        return false;
                    }
                    item += `Kind: fourslash.PtrTo(${kind}), `;
                    break;
                case 'kindModifiers':
                    const modifiers = parseKindModifiers(init);
                    if (!modifiers) {
                        return false;
                    }
                    ({ isDeprecated, isOptional, extensions } = modifiers);
                    break;
                case 'commitCharacters':
                case 'replacementSpan':
                    return false; // !!! support these later
                    break;
                default:
                    return false; // Unsupported property
            }
        }
        if (!name) {
            return false; // Shouldn't happen
        }
        if (isOptional) {
            insertText ??= name;
            filterText ??= name;
            name += '?';
        }
        item += `Label: ${getGoStringLiteral(name!)}, `;
        if (insertText) item += `InsertText: fourslash.PtrTo(${getGoStringLiteral(insertText)}), `;
        if (filterText) item += `FilterText: fourslash.PtrTo(${getGoStringLiteral(filterText)}), `;
        item += "}";
        return item;
    }
    return false; // Unsupported expression type
}

function parseKind(expr: ts.Expression): string | false {
    if (!ts.isStringLiteral(expr)) {
        return false;
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

const fileKindModifiers = new Set([".d.ts",
	".ts",
    ".tsx",
    ".js",
    ".jsx",
    ".json",
]);

function parseKindModifiers(expr: ts.Expression): { isOptional: boolean, isDeprecated: boolean, extensions: string[] } | false {
    if (!ts.isStringLiteral(expr)) {
        return false;
    }
    let isOptional = false;
    let isDeprecated = false;
    const extensions: string[] = [];
    const modifiers = expr.text.split(',');
    for (const modifier of modifiers) {
        switch (modifier) {
            case 'optional':
                isOptional = true;
                break;
            case 'deprecated':
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

function parseSortText(expr: ts.Expression): string | false {
    const text = expr.getText();
    switch (text) {
        case 'completion.SortText.LocalDeclarationPriority':
            return 'ls.SortTextLocalDeclarationPriority';
        case 'completion.SortText.LocationPriority':
            return 'ls.SortTextLocationPriority';
        case 'completion.SortText.OptionalMember':
            return 'ls.SortTextOptionalMember';
        case 'completion.SortText.MemberDeclaredBySpreadAssignment':
            return 'ls.SortTextMemberDeclaredBySpreadAssignment';
        case 'completion.SortText.SuggestedClassMember':
            return 'ls.SortTextSuggestedClassMember';
        case 'completion.SortText.GlobalsOrKeywords':
            return 'ls.SortTextGlobalsOrKeywords';
        case 'completion.SortText.AutoImportSuggestions':
            return 'ls.SortTextAutoImportSuggestions';
        case 'completion.SortText.ClassMemberSnippets':
            return 'ls.SortTextClassMemberSnippets';
        case 'completion.SortText.JavaScriptIdentifiers':
            return 'ls.SortTextJavaScriptIdentifiers';
        default:
            return false; // !!! support deprecated/obj literal prop/etc
    }
}

enum CommandKind {
    verifyCompletions,
    goToMarker,
}


interface VerifyCompletionsCmd {
    kind: CommandKind.verifyCompletions;
    marker: string;
    args?: VerifyCompletionsArgs;
}

interface VerifyCompletionsArgs {
    includes?: string;
    excludes?: string;
    exact?: string;
}

interface GoToMarkerCmd {
    kind: CommandKind.goToMarker;
    marker: string;
}

type Cmd = VerifyCompletionsCmd | GoToMarkerCmd;

function generateVerifyCompletions({ marker, args }: VerifyCompletionsCmd): string {
    let expectedList = "nil";
    if (args) {
        const includes = args.includes ? `Includes: ${args.includes},` : '';
        const excludes = args.excludes ? `Excludes: ${args.excludes},` : '';
        const exact = args.exact ? `Exact: ${args.exact},` : '';
        // !!! isIncomplete
        // !!! itemDefaults/commitCharacters
        expectedList = `&fourslash.VerifyCompletionsExpectedList{
    IsIncomplete: false,
    ItemDefaults: &lsproto.CompletionItemDefaults{
        CommitCharacters: &fourslash.DefaultCommitCharacters,
    },
    Items: &fourslash.VerifyCompletionsExpectedItems{
        ${exact}
        ${includes}
        ${excludes}
    },
}`;
    }
    return `f.VerifyCompletions(t, ${marker}, ${expectedList})`;
}

function generateGoToMarker({ marker }: GoToMarkerCmd): string {
    return `f.GoToMarker(t, ${marker})`;
}

function generateCmd(cmd: Cmd): string {
    switch (cmd.kind) {
        case CommandKind.verifyCompletions:
            return generateVerifyCompletions(cmd as VerifyCompletionsCmd);
        case CommandKind.goToMarker:
            return generateGoToMarker(cmd as GoToMarkerCmd);
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
    const commands = test.commands.map(cmd => generateCmd(cmd)).join('\n');
    const imports = [`"github.com/microsoft/typescript-go/internal/fourslash"`];
    if (commands.includes('ls.')) {
        imports.push(`"github.com/microsoft/typescript-go/internal/ls"`);
    }
    if (commands.includes('lsproto.')) {
        imports.push(`"github.com/microsoft/typescript-go/internal/lsp/lsproto"`);
    }
    imports.push(`"github.com/microsoft/typescript-go/internal/testutil"`)
    const template = `package ls_test
import (
	"testing"

    ${imports.join('\n\t')}
)

func Test${testName}(t *testing.T) {
    t.Parallel()
    defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = ${content}
    f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
    defer done()
    ${commands}
}`
    return template;
}


parseTypeScriptFiles(folderPath);
console.log(unparsedFiles.join('\n'));