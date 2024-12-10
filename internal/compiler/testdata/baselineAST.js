// @ts-check
const ts = require("../../../_submodules/TypeScript/built/local/typescript");
const fs = require("fs");
const path = require("path");
/** @param node {ts.Node} */
function printNode(node, indentLevel = 0) {
    let s = "";
    if (ts.isIdentifier(node)) {
        s = `${" ".repeat(indentLevel * 2)}${unaliasKind(node.kind)}: '${node.getFullText()}'\n`;
    }
    else if (node.kind === ts.SyntaxKind.EndOfFileToken) {
        return "";
    }
    else {
        s = `${" ".repeat(indentLevel * 2)}${unaliasKind(node.kind)}\n`;
    }
    node.forEachChild(child => {
        s += printNode(child, indentLevel + 1);
        return false;
    });
    return s;
}
function unaliasKind(kind) {
    switch (kind) {
        // Special case (T? now parses as an optional type), not an alias
        case ts.SyntaxKind.JSDocNullableType:
            return "OptionalType";
        case ts.SyntaxKind.FirstAssignment:
            return "EqualsToken";
        case ts.SyntaxKind.LastAssignment:
            return "CaretEqualsToken";
        case ts.SyntaxKind.FirstCompoundAssignment:
            return "PlusEqualsToken";
        case ts.SyntaxKind.LastCompoundAssignment:
            return "CaretEqualsToken";
        case ts.SyntaxKind.FirstReservedWord:
            return "BreakKeyword";
        case ts.SyntaxKind.LastReservedWord:
            return "WithKeyword";
        case ts.SyntaxKind.FirstKeyword:
            return "BreakKeyword";
        case ts.SyntaxKind.LastKeyword:
            return "OfKeyword";
        case ts.SyntaxKind.FirstFutureReservedWord:
            return "ImplementsKeyword";
        case ts.SyntaxKind.LastFutureReservedWord:
            return "YieldKeyword";
        case ts.SyntaxKind.FirstTypeNode:
            return "TypePredicate";
        case ts.SyntaxKind.LastTypeNode:
            return "ImportType";
        case ts.SyntaxKind.FirstPunctuation:
            return "OpenBraceToken";
        case ts.SyntaxKind.LastPunctuation:
            return "CaretEqualsToken";
        case ts.SyntaxKind.FirstToken:
            return "Unknown";
        case ts.SyntaxKind.LastToken:
            return "OfKeyword";
        case ts.SyntaxKind.FirstTriviaToken:
            return "SingleLineCommentTrivia";
        case ts.SyntaxKind.LastTriviaToken:
            return "ConflictMarkerTrivia";
        case ts.SyntaxKind.FirstLiteralToken:
            return "NumericLiteral";
        case ts.SyntaxKind.LastLiteralToken:
            return "NoSubstitutionTemplateLiteral";
        case ts.SyntaxKind.FirstTemplateToken:
            return "NoSubstitutionTemplateLiteral";
        case ts.SyntaxKind.LastTemplateToken:
            return "TemplateTail";
        case ts.SyntaxKind.FirstBinaryOperator:
            return "LessThanToken";
        case ts.SyntaxKind.LastBinaryOperator:
            return "CaretEqualsToken";
        case ts.SyntaxKind.FirstStatement:
            return "VariableStatement";
        case ts.SyntaxKind.LastStatement:
            return "DebuggerStatement";
        case ts.SyntaxKind.FirstNode:
            return "QualifiedName";
        case ts.SyntaxKind.FirstJSDocNode:
            return "JSDocTypeExpression";
        case ts.SyntaxKind.LastJSDocNode:
            return "JSDocImportTag";
        case ts.SyntaxKind.FirstJSDocTagNode:
            return "JSDocTag";
        case ts.SyntaxKind.LastJSDocTagNode:
            return "JSDocImportTag";
        case ts.SyntaxKind.AssertClause:
            return "ImportAttributes";
        case ts.SyntaxKind.AssertEntry:
            return "ImportAttribute";
        default:
            kind = ts.SyntaxKind[kind];
            if (kind === "FirstContextualKeyword") {
                return "AbstractKeyword";
            }
            else if (kind === "LastContextualKeyword") {
                return "OfKeyword";
            }
            else {
                return kind;
            }
    }
}
/** @param filePath {string} */
function printAST(filePath) {
    const fileContent = fs.readFileSync(filePath, "utf8");
    const sourceFile = ts.createSourceFile(
        path.basename(filePath),
        fileContent,
        ts.ScriptTarget.ESNext,
        true,
    );
    return printNode(sourceFile);
}
/**
 * @param inputRoot {string}
 * @param targetRoot {string}
 */
function processDirectory(inputRoot, targetRoot) {
    worker(inputRoot);
    function worker(dir) {
        if (!inputRoot.endsWith("/")) {
            inputRoot += "/";
        }
       for (const dirent of fs.readdirSync(dir, { withFileTypes: true })) {
            let fullPath = path.join(dir, dirent.name);
            const ext = path.extname(dirent.name)
            if (dirent.isDirectory()) {
                worker(fullPath);
            }
            else if (dirent.isFile() && (ext === '.ts' || ext === '.tsx' || ext === '.js' || ext === '.jsx')) {
                if (
                    // Too deep for a simple tree walker
                    dirent.name.endsWith("binderBinaryExpressionStress.ts") ||
                    dirent.name.endsWith("binderBinaryExpressionStress.js") ||
                    dirent.name.endsWith("binderBinaryExpressionStressJs.ts") ||
                    dirent.name.endsWith("binderBinaryExpressionStressJs.js") ||
                    // Very large minified code
                    dirent.name.includes("codeMirrorModule") ||
                    // Not actually .js
                    dirent.name.includes("reference/config/") ||
                    dirent.name.includes("reference/tsc") ||
                    dirent.name.includes("reference/tsserver") ||
                    dirent.name.includes("reference/tsbuild")
                ) {
                    continue;
                }
                const astContent = printAST(fullPath);
                if (fullPath.startsWith(inputRoot)) {
                    fullPath = fullPath.slice(inputRoot.length);
                }
                else {
                    console.error("Unexpected file path: " + fullPath);
                }
                const outputFileName =  fullPath.replace(/[\/\\]/g, "_") + ".ast";
                const outputFilePath = path.join(targetRoot, outputFileName);
                fs.writeFileSync(outputFilePath, astContent);
            }
        }
    }
}
const outputFlagOrFile = process.argv[2];
if (outputFlagOrFile === "-r") {
    const inputDir = process.argv[3];
    const outputDir = process.argv[4];
    if (process.argv.length !== 5 || !inputDir || !outputDir) {
        console.error(
            "node internal/compiler/testdata/baselineAST.js -r _submodules/TypeScript testdata/baselines/gold",
        );
        process.exit(1);
    }
    if (!fs.existsSync(outputDir)) {
        fs.mkdirSync(outputDir, { recursive: true });
    }
    processDirectory(inputDir, outputDir);
} else {
    const inputFile = process.argv[2];
    if (process.argv.length !== 3) {
        console.error(
            "node internal/compiler/testdata/baselineAST.js _submodules/TypeScript/src/compiler/checker.ts",
        );
        process.exit(1);
    }
    console.log(printAST(inputFile))
}
