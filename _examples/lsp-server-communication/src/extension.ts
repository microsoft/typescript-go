import {
    AsyncAPI,
    SymbolFlags,
    TypeFlags,
} from "@typescript/api/async";
import {
    type Node,
    SyntaxKind,
} from "@typescript/ast";
import * as vscode from "vscode";

let outputChannel: vscode.OutputChannel | undefined;

export function activate(context: vscode.ExtensionContext) {
    outputChannel = vscode.window.createOutputChannel("TypeScript AST");

    // Command: Print AST of current file
    const printASTCommand = vscode.commands.registerCommand("lsp-example.printAST", async () => {
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

    // Command: Print symbol info at cursor position
    const printSymbolInfoCommand = vscode.commands.registerCommand("lsp-example.printSymbolInfo", async () => {
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

        const position = editor.document.offsetAt(editor.selection.active);

        try {
            const pipePath = await vscode.commands.executeCommand<string>(
                "typescript.native-preview.initializeAPIConnection",
            );

            if (!pipePath) {
                vscode.window.showErrorMessage(
                    "Could not get API connection. Is TypeScript Native Preview extension active?",
                );
                return;
            }

            outputChannel!.show();
            outputChannel!.appendLine(`\n=== Symbol Info at position ${position} ===\n`);

            const api = AsyncAPI.fromLSPConnection({ pipePath });

            try {
                const project = await api.getDefaultProjectForFile(fileName);
                if (!project) {
                    outputChannel!.appendLine(`File is not part of any TypeScript project: ${fileName}`);
                    return;
                }

                // Get the symbol at the cursor position
                const symbol = await project.getSymbolAtPosition(fileName, position);

                if (!symbol) {
                    outputChannel!.appendLine(`No symbol found at position ${position}`);
                    return;
                }

                outputChannel!.appendLine(`Symbol: ${symbol.name}`);
                outputChannel!.appendLine(`Symbol Flags: ${formatSymbolFlags(symbol.flags)}`);

                // Get the type of the symbol
                const type = await project.getTypeOfSymbol(symbol);

                if (type) {
                    outputChannel!.appendLine(`Type Flags: ${formatTypeFlags(type.flags)}`);
                }
                else {
                    outputChannel!.appendLine(`No type found for symbol`);
                }
            }
            finally {
                await api.close();
            }
        }
        catch (error) {
            outputChannel!.appendLine(`Error: ${error}`);
            vscode.window.showErrorMessage(`Failed to get symbol info: ${error}`);
        }
    });

    context.subscriptions.push(printASTCommand);
    context.subscriptions.push(printSymbolInfoCommand);
    context.subscriptions.push(outputChannel);
}

function printNode(node: Node, depth: number, output: vscode.OutputChannel): void {
    const indent = "  ".repeat(depth);
    output.appendLine(`${indent}${SyntaxKind[node.kind]} [${node.pos}-${node.end}]`);
    node.forEachChild(child => printNode(child, depth + 1, output));
}

function formatFlags<T extends number>(flags: T, enumObj: object): string {
    const names: string[] = [];
    for (const [name, value] of Object.entries(enumObj)) {
        // Skip reverse mappings (numeric keys) and composite flags
        if (typeof value === "number" && value !== 0 && (flags & value) === value) {
            // Check it's a power of 2 (single bit flag)
            if ((value & (value - 1)) === 0) {
                names.push(name);
            }
        }
    }
    return names.length > 0 ? names.join(" | ") : `0x${flags.toString(16)}`;
}

function formatSymbolFlags(flags: SymbolFlags): string {
    return formatFlags(flags, SymbolFlags);
}

function formatTypeFlags(flags: TypeFlags): string {
    return formatFlags(flags, TypeFlags);
}

export function deactivate() {
    outputChannel?.dispose();
}
