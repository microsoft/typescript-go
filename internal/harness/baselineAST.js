// @ts-check
const ts = require("../../_submodules/TypeScript/built/local/typescript");
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
        fs.readdirSync(dir, { withFileTypes: true }).forEach(dirent => {
            let fullPath = path.join(dir, dirent.name);
            if (dirent.isDirectory()) {
                worker(fullPath);
            }
            else if (dirent.isFile() && path.extname(dirent.name) === ".ts") {
                // Too deep for a simple tree walker
                if (
                    dirent.name.endsWith("binderBinaryExpressionStress.ts") ||
                    dirent.name.endsWith("binderBinaryExpressionStressJs.ts")
                ) {
                    return;
                }
                const astContent = printAST(fullPath);
                if (fullPath.startsWith(inputRoot)) {
                    fullPath = fullPath.slice(inputRoot.length);
                }
                else {
                    console.error("Unexpected file path: " + fullPath);
                }
                const outputFileName = fullPath.replace(/[\/\\]/g, "_") + ".ast.txt";
                const outputFilePath = path.join(targetRoot, outputFileName);
                fs.writeFileSync(outputFilePath, astContent);
            }
        });
    }
}
const inputDir = process.argv[2];
const outputDir = process.argv[3];
if (!inputDir || !outputDir) {
    console.error(
        "node internal/harness/baselineAST.js _submodules/TypeScript testdata/baselines/gold",
    );
    process.exit(1);
}
if (!fs.existsSync(outputDir)) {
    fs.mkdirSync(outputDir, { recursive: true });
}
processDirectory(inputDir, outputDir);
