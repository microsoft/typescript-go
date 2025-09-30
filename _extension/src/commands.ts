import * as vscode from "vscode";
import { Client } from "./client";

export function registerEnablementCommands(context: vscode.ExtensionContext): void {
    context.subscriptions.push(vscode.commands.registerCommand("typescript.native-preview.enable", () => {
        // Fire and forget, because this will restart the extension host and cause an error if we await
        updateUseTsgoSetting(true);
    }));

    context.subscriptions.push(vscode.commands.registerCommand("typescript.native-preview.disable", () => {
        // Fire and forget, because this will restart the extension host and cause an error if we await
        updateUseTsgoSetting(false);
    }));
}

export function registerLanguageCommands(context: vscode.ExtensionContext, client: Client, outputChannel: vscode.OutputChannel, traceOutputChannel: vscode.OutputChannel): vscode.Disposable[] {
    const disposables: vscode.Disposable[] = [];

    disposables.push(vscode.commands.registerCommand("typescript.native-preview.restart", () => {
        return client.restart(context);
    }));

    disposables.push(vscode.commands.registerCommand("typescript.native-preview.output.focus", () => {
        outputChannel.show();
    }));

    disposables.push(vscode.commands.registerCommand("typescript.native-preview.lsp-trace.focus", () => {
        traceOutputChannel.show();
    }));

    disposables.push(vscode.commands.registerCommand("typescript.native-preview.selectVersion", async () => {
    }));

    disposables.push(vscode.commands.registerCommand("typescript.native-preview.showMenu", showCommands));

    disposables.push(vscode.commands.registerCommand("typescript.native-preview.sortImports", async () => {
        return sortImports(client);
    }));

    return disposables;
}

/**
 * Updates the TypeScript Native Preview setting and reloads extension host.
 */
async function updateUseTsgoSetting(enable: boolean): Promise<void> {
    const tsConfig = vscode.workspace.getConfiguration("typescript");
    let target: vscode.ConfigurationTarget | undefined;
    const useTsgo = tsConfig.inspect("experimental.useTsgo");
    if (useTsgo) {
        target = useTsgo.workspaceFolderValue !== undefined ? vscode.ConfigurationTarget.WorkspaceFolder :
            useTsgo.workspaceValue !== undefined ? vscode.ConfigurationTarget.Workspace :
            useTsgo.globalValue !== undefined ? vscode.ConfigurationTarget.Global : undefined;
    }
    // Update the setting and restart the extension host (needed to change the state of the built-in TS extension)
    await tsConfig.update("experimental.useTsgo", enable, target);
    await vscode.commands.executeCommand("workbench.action.restartExtensionHost");
}

async function sortImports(client: Client): Promise<void> {
    const editor = vscode.window.activeTextEditor;
    if (!editor) {
        vscode.window.showErrorMessage("No active editor");
        return;
    }

    const document = editor.document;
    const languageId = document.languageId;

    // Check if the file is TypeScript or JavaScript
    if (!["typescript", "javascript", "typescriptreact", "javascriptreact"].includes(languageId)) {
        vscode.window.showErrorMessage("Sort Imports is only available for TypeScript and JavaScript files");
        return;
    }

    try {
        // Execute the sort imports command on the server via LSP
        await client.executeCommand(
            "typescript-go.organizeImports",
            document.uri.toString(),
        );
        vscode.window.showInformationMessage("Imports sorted successfully");
    }
    catch (error) {
        vscode.window.showErrorMessage(`Failed to sort imports: ${error}`);
    }
}

async function showCommands(): Promise<void> {
    const commands: readonly { label: string; description: string; command: string; }[] = [
        {
            label: "$(symbol-namespace) Sort Imports",
            description: "Sort imports in the current file",
            command: "typescript.native-preview.sortImports",
        },
        {
            label: "$(refresh) Restart Server",
            description: "Restart the TypeScript Native Preview language server",
            command: "typescript.native-preview.restart",
        },
        {
            label: "$(output) Show TS Server Log",
            description: "Show the TypeScript Native Preview server log",
            command: "typescript.native-preview.output.focus",
        },
        {
            label: "$(debug-console) Show LSP Messages",
            description: "Show the LSP communication trace",
            command: "typescript.native-preview.lsp-trace.focus",
        },
        {
            label: "$(stop-circle) Disable TypeScript Native Preview",
            description: "Switch back to the built-in TypeScript extension",
            command: "typescript.native-preview.disable",
        },
    ];

    const selected = await vscode.window.showQuickPick(commands, {
        placeHolder: "TypeScript Native Preview Commands",
    });

    if (selected) {
        await vscode.commands.executeCommand(selected.command);
    }
}
