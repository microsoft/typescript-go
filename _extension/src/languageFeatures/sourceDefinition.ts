import * as vscode from "vscode";
import { LanguageClient } from "vscode-languageclient/node";

const sourceDefinitionMethod = "custom/textDocument/sourceDefinition";
const sourceDefinitionCommand = "typescript.native-preview.goToSourceDefinition";
const sourceDefinitionContext = "tsSupportsSourceDefinition";

interface LspPosition {
    line: number;
    character: number;
}

interface LspRange {
    start: LspPosition;
    end: LspPosition;
}

interface LspLocation {
    uri: string;
    range: LspRange;
}

interface LspLocationLink {
    targetUri: string;
    targetSelectionRange: LspRange;
}

type SourceDefinitionResponse = LspLocation | LspLocation[] | LspLocationLink[] | null;

function isLocationLink(value: LspLocation | LspLocationLink): value is LspLocationLink {
    return "targetUri" in value;
}

function toVsRange(range: LspRange): vscode.Range {
    return new vscode.Range(
        new vscode.Position(range.start.line, range.start.character),
        new vscode.Position(range.end.line, range.end.character),
    );
}

function toVsLocations(response: SourceDefinitionResponse): vscode.Location[] {
    if (!response) {
        return [];
    }

    const items = Array.isArray(response) ? response : [response];
    return items.map(item => {
        if (isLocationLink(item)) {
            return new vscode.Location(vscode.Uri.parse(item.targetUri), toVsRange(item.targetSelectionRange));
        }
        return new vscode.Location(vscode.Uri.parse(item.uri), toVsRange(item.range));
    });
}

export function registerSourceDefinitionFeature(client: LanguageClient): vscode.Disposable {
    const capabilities = client.initializeResult?.capabilities as { customSourceDefinitionProvider?: boolean; } | undefined;
    const enabled = !!capabilities?.customSourceDefinitionProvider;
    void vscode.commands.executeCommand("setContext", sourceDefinitionContext, enabled);

    if (!enabled) {
        return new vscode.Disposable(() => {
            void vscode.commands.executeCommand("setContext", sourceDefinitionContext, false);
        });
    }

    const disposable = vscode.commands.registerCommand(sourceDefinitionCommand, async () => {
        const activeEditor = vscode.window.activeTextEditor;
        if (!activeEditor) {
            vscode.window.showErrorMessage("Go to Source Definition failed. No editor is active.");
            return;
        }

        const { document } = activeEditor;
        if (!["javascript", "javascriptreact", "typescript", "typescriptreact"].includes(document.languageId)) {
            vscode.window.showErrorMessage("Go to Source Definition failed. Unsupported file type.");
            return;
        }

        const position = activeEditor.selection.active;
        await vscode.window.withProgress({
            location: vscode.ProgressLocation.Window,
            title: "Finding source definitions",
        }, async (_, token) => {
            let response: SourceDefinitionResponse;
            try {
                response = await client.sendRequest<SourceDefinitionResponse>(
                    sourceDefinitionMethod,
                    client.code2ProtocolConverter.asTextDocumentPositionParams(document, position),
                    token,
                );
            }
            catch {
                return;
            }

            if (token.isCancellationRequested) {
                return;
            }

            const locations = toVsLocations(response);
            if (locations.length === 0) {
                vscode.window.showErrorMessage("No source definitions found.");
                return;
            }

            if (locations.length === 1) {
                const location = locations[0];
                await vscode.commands.executeCommand(
                    "vscode.open",
                    location.uri.with({
                        fragment: `L${location.range.start.line + 1},${location.range.start.character + 1}`,
                    }),
                );
                return;
            }

            await vscode.commands.executeCommand("editor.action.showReferences", document.uri, position, locations);
        });
    });

    return vscode.Disposable.from(
        disposable,
        new vscode.Disposable(() => {
            void vscode.commands.executeCommand("setContext", sourceDefinitionContext, false);
        }),
    );
}
