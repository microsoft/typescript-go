"use strict";

const vscode = require("vscode");

const newExtensionId = "TypeScriptTeam.vscode-typescript";

async function activate() {
    if (vscode.extensions.getExtension(newExtensionId)) {
        return;
    }

    const install = "Install";
    const later = "Later";
    const result = await vscode.window.showInformationMessage(
        "TypeScript 7 Native Preview has moved to a new extension. Install it now?",
        install,
        later,
    );
    if (result !== install) {
        return;
    }

    try {
        await vscode.commands.executeCommand("workbench.extensions.installExtension", newExtensionId);
    }
    catch (error) {
        const message = error instanceof Error ? error.message : String(error);
        await vscode.window.showErrorMessage(`Failed to install ${newExtensionId}: ${message}`);
    }
}

function deactivate() {}

module.exports = {
    activate,
    deactivate,
};
