// @ts-check
const ts = require("../../_submodules/TypeScript/built/local/typescript");
const fs = require("fs");
const path = require("path");
/** @param node {ts.Node} */
function printNode(node, indentLevel = 0) {
  let s = "";
  if (ts.isIdentifier(node)) {
    s = `${" ".repeat(indentLevel * 2)}${
      ts.SyntaxKind[node.kind]
    }: '${node.getFullText()}'\n`;
  } else {
    s = `${" ".repeat(indentLevel * 2)}${ts.SyntaxKind[node.kind]}\n`;
  }
  node.forEachChild((child) => {
    s += printNode(child, indentLevel + 1);
    return false;
  });
  return s;
}
/** @param filePath {string} */
function printAST(filePath) {
  const fileContent = fs.readFileSync(filePath, "utf8");
  const sourceFile = ts.createSourceFile(
    path.basename(filePath),
    fileContent,
    ts.ScriptTarget.ESNext,
    true
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
    fs.readdirSync(dir, { withFileTypes: true }).forEach((dirent) => {
      let fullPath = path.join(dir, dirent.name);
      if (dirent.isDirectory()) {
        worker(fullPath);
      } else if (dirent.isFile() && path.extname(dirent.name) === ".ts") {
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
        } else {
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
    "Please provide an input directory and an output directory as arguments."
  );
  process.exit(1);
}
if (!fs.existsSync(outputDir)) {
  fs.mkdirSync(outputDir, { recursive: true });
}
processDirectory(inputDir, outputDir);
