import * as vscode from "vscode";
import { jsTsLanguageModes } from "./util";

export function setupStatusBar(version: string, pid?: number): vscode.Disposable {
    const statusItem = vscode.languages.createLanguageStatusItem("typescript.native-preview.status", jsTsLanguageModes);
    statusItem.name = "TypeScript Native Preview";

    function updateText() {
        const showPID = vscode.workspace.getConfiguration("typescript.native-preview").get<boolean>("showPID", false);
        statusItem.text = showPID && pid
            ? `$(beaker) tsgo ${version} (PID: ${pid})`
            : `$(beaker) tsgo ${version}`;
    }
    updateText();

    statusItem.detail = "TypeScript Native Preview Language Server";
    statusItem.command = {
        title: "Show Menu",
        command: "typescript.native-preview.showMenu",
    };

    const configListener = vscode.workspace.onDidChangeConfiguration(e => {
        if (e.affectsConfiguration("typescript.native-preview.showPID")) {
            updateText();
        }
    });

    return {
        dispose() {
            statusItem.dispose();
            configListener.dispose();
        },
    };
}
