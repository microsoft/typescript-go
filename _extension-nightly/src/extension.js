"use strict";

const vscode = require("vscode");

const newExtensionId = "TypeScriptTeam.vscode-typescript";

function activate() {
    if (vscode.extensions.getExtension(newExtensionId)) {
        return;
    }

    void installReplacement();
}

async function installReplacement() {
    try {
        await vscode.commands.executeCommand("workbench.extensions.installExtension", newExtensionId);
    }
    catch (error) {
        const message = error instanceof Error ? error.message : String(error);
        await vscode.window.showErrorMessage(vscode.l10n.t("Failed to install {0}: {1}", newExtensionId, message));
    }
}

function deactivate() {}

module.exports = {
    activate,
    deactivate,
};
