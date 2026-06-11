import * as vscode from "vscode";

import {
    registerEnablementCommands,
    updateUseTsgoSetting,
} from "./commands";
import {
    aiConnectionString,
    getUseTsgo,
    getUseTsgoFalseSetting,
    getTypeScriptLanguageFeaturesApi,
    JsTsServerSelection,
    needsExtHostRestartOnChange,
} from "./util";

import { TelemetryReporter as VSCodeTelemetryReporter } from "@vscode/extension-telemetry";
import {
    promptUseWorkspaceVersion,
    SessionManager,
} from "./session";

import { ExperimentationService } from "./experimentationService";
import { createTelemetryReporter } from "./telemetryReporting";

import assert from "node:assert";

export interface ExtensionAPI {
    onLanguageServerInitialized: vscode.Event<void>;
    initializeAPIConnection(pipe?: string): Promise<string>;
}

export async function activate(context: vscode.ExtensionContext): Promise<ExtensionAPI | undefined> {
    await vscode.commands.executeCommand("setContext", "typescript.native-preview.serverRunning", false);

    const telemetryReporter = createTelemetryReporter(new VSCodeTelemetryReporter(aiConnectionString));
    context.subscriptions.push(telemetryReporter);

    const version = context.extension.packageJSON.version;
    assert(typeof version === "string");
    // Constructing the experimentation service actually sets shared properties
    // so that events include context on treatments/flights.
    // If we actually need to read treatment variables we would hold onto this instance,
    // but for now we just construct it to ensure shared properties are set for telemetry.
    void new ExperimentationService(telemetryReporter, context.extension.id, version, context.globalState);

    registerEnablementCommands(context, telemetryReporter);

    const output = vscode.window.createOutputChannel("typescript-native-preview", { log: true });
    context.subscriptions.push(output);

    const languageServerInitializedEventEmitter = new vscode.EventEmitter<void>();
    context.subscriptions.push(languageServerInitializedEventEmitter);

    const sessionManager = new SessionManager(context, output, languageServerInitializedEventEmitter, telemetryReporter);
    context.subscriptions.push(sessionManager);

    const tsApi = await getTypeScriptLanguageFeaturesApi();
    await vscode.commands.executeCommand("setContext", "typescript.native-preview.usingTypeScriptLanguageFeaturesApi", !!tsApi);
    let selection = tsApi?.getServerSelection();
    if (tsApi) {
        context.subscriptions.push(tsApi.onDidChangeServerSelection(async nextSelection => {
            selection = nextSelection;
            await updateSessionForSelection(nextSelection);
        }));
    }
    else {
        let configChangeTimeout: ReturnType<typeof setTimeout> | undefined;
        context.subscriptions.push(vscode.workspace.onDidChangeConfiguration(event => {
            if (event.affectsConfiguration("typescript.experimental.useTsgo") || event.affectsConfiguration("js/ts.experimental.useTsgo")) {
                clearTimeout(configChangeTimeout);
                configChangeTimeout = setTimeout(async () => {
                    if (needsExtHostRestartOnChange()) {
                        const selected = await vscode.window.showInformationMessage(vscode.l10n.t("TypeScript Native Preview setting has changed. Restart extensions to apply changes."), vscode.l10n.t("Restart Extensions"));
                        if (selected) {
                            vscode.commands.executeCommand("workbench.action.restartExtensionHost");
                        }
                    }
                    else {
                        const useTsgo = getUseTsgo();
                        if (useTsgo) {
                            await sessionManager.restart(context);
                        }
                        else {
                            await sessionManager.stop();
                        }
                    }
                }, 100);
            }
        }));
        context.subscriptions.push({ dispose: () => clearTimeout(configChangeTimeout) });
    }

    const hasOnboardedTsgoStateKey = "hasOnboardedTsgo";
    const shouldOnboardTsgo = !context.globalState.get<boolean>(hasOnboardedTsgoStateKey);
    if (!tsApi && shouldOnboardTsgo) {
        await context.globalState.update(hasOnboardedTsgoStateKey, true);
    }

    if (context.extensionMode === vscode.ExtensionMode.Development) {
        const tsExtension = vscode.extensions.getExtension("vscode.typescript-language-features");
        if (!tsExtension) {
            if (tsApi ? selection?.kind !== "lsp" : !getUseTsgo()) {
                vscode.window.showWarningMessage(
                    vscode.l10n.t("The built-in TypeScript extension is disabled. Sync launch.json with launch.template.json to reenable."),
                    vscode.l10n.t("OK"),
                );
                return;
            }
        }
        else if (tsApi ? selection?.kind !== "lsp" : getUseTsgo() === false) {
            const settingName = tsApi ? "js/ts.languageServer.preference" : getUseTsgoFalseSetting() ?? "js/ts.experimental.useTsgo";
            const enableSettingString = vscode.l10n.t("Enable Setting");
            vscode.window.showWarningMessage(
                tsApi
                    ? vscode.l10n.t("TypeScript Native Preview is running in development mode but is not the selected JavaScript and TypeScript language server.")
                    : vscode.l10n.t(`TypeScript Native Preview is running in development mode with "{0}" set to false.`, settingName),
                enableSettingString,
                vscode.l10n.t("Ignore"),
            ).then(selected => {
                if (selected === enableSettingString) {
                    vscode.commands.executeCommand("typescript.native-preview.enable");
                }
            });
            return;
        }
    }
    else if (tsApi) {
        if (selection?.kind !== "lsp") {
            output.appendLine(vscode.l10n.t("TypeScript Native Preview is disabled. Select 'Enable TypeScript Native Preview (Experimental)' in the command palette to enable it."));
            return;
        }
    }
    else {
        const useTsgo = getUseTsgo();
        if (useTsgo === false) {
            output.appendLine(vscode.l10n.t("TypeScript Native Preview is disabled. Select 'Enable TypeScript Native Preview (Experimental)' in the command palette to enable it."));
            return;
        }
        else if (useTsgo === undefined) {
            if (shouldOnboardTsgo) {
                updateUseTsgoSetting(true);
                return;
            }
            output.appendLine(vscode.l10n.t("TypeScript Native Preview is disabled. Select 'Enable TypeScript Native Preview (Experimental)' in the command palette to enable it."));
            return;
        }
    }

    await sessionManager.start(context, selection);

    if (!tsApi) {
        promptUseWorkspaceVersion(context).catch(err => {
            output.appendLine(vscode.l10n.t(`Error prompting to use workspace version: {0}`, String(err)));
        });
    }

    async function updateSessionForSelection(nextSelection: JsTsServerSelection): Promise<void> {
        if (nextSelection.kind === "lsp") {
            await sessionManager.restart(context, nextSelection);
        }
        else {
            await sessionManager.stop();
        }
    }

    function onLanguageServerInitialized(listener: () => void): vscode.Disposable {
        if (sessionManager.currentSession?.client.isInitialized) {
            listener();
        }
        return languageServerInitializedEventEmitter.event(listener);
    }

    return {
        onLanguageServerInitialized: onLanguageServerInitialized,
        async initializeAPIConnection(pipe?: string): Promise<string> {
            return sessionManager.initializeAPIConnection(pipe);
        },
    };
}
