import { AsyncAPI } from "@typescript/api/async";
import {
    type Node,
    SyntaxKind,
} from "@typescript/ast";
import * as vscode from "vscode";

let outputChannel: vscode.OutputChannel | undefined;

export function activate(context: vscode.ExtensionContext) {
    outputChannel = vscode.window.createOutputChannel("TypeScript AST");

    const disposable = vscode.commands.registerCommand("lsp-example.printAST", async () => {
        const editor = vscode.window.activeTextEditor;
        if (!editor) {
            vscode.window.showErrorMessage("No active editor");
            return;
        }

        const fileName = editor.document.uri.fsPath;
        if (!fileName.endsWith(".ts") && !fileName.endsWith(".tsx")) {
            vscode.window.showErrorMessage("Current file is not a TypeScript file");
            return;
        }

        try {
            // Get the API connection from the TypeScript Native Preview extension
            const pipePath = await vscode.commands.executeCommand<string>(
                "typescript.native-preview.initializeAPIConnection",
            );

            if (!pipePath) {
                vscode.window.showErrorMessage(
                    "Could not get API connection. Is TypeScript Native Preview extension active?",
                );
                return;
            }

            outputChannel!.appendLine(`Connecting to API at: ${pipePath}`);
            outputChannel!.show();

            // Connect to the API using @typescript/api
            const api = AsyncAPI.fromLSPConnection({ pipePath });

            try {
                // Get the project for this file
                const project = await api.getDefaultProjectForFile(fileName);

                if (!project) {
                    outputChannel!.appendLine(`File is not part of any TypeScript project: ${fileName}`);
                    return;
                }

                outputChannel!.appendLine(`Found project: ${project.configFileName}`);

                // Get the source file AST
                const sourceFile = await project.getSourceFile(fileName);

                if (!sourceFile) {
                    outputChannel!.appendLine(`Could not get source file: ${fileName}`);
                    return;
                }

                outputChannel!.appendLine(`\n=== AST for ${fileName} ===\n`);

                // Walk and print the AST
                printNode(sourceFile, 0, outputChannel!);
            }
            finally {
                await api.close();
            }
        }
        catch (error) {
            outputChannel!.appendLine(`Error: ${error}`);
            vscode.window.showErrorMessage(`Failed to print AST: ${error}`);
        }
    });

    context.subscriptions.push(disposable);
    context.subscriptions.push(outputChannel);
}

function printNode(node: Node, depth: number, output: vscode.OutputChannel): void {
    const indent = "  ".repeat(depth);
    output.appendLine(`${indent}${SyntaxKind[node.kind]} [${node.pos}-${node.end}]`);
    node.forEachChild(child => printNode(child, depth + 1, output));
}

export function deactivate() {
    outputChannel?.dispose();
}
