import * as vscode from "vscode";
import { LanguageClient } from "vscode-languageclient/node";

const multiDocumentHighlightMethod = "custom/textDocument/multiDocumentHighlight";

interface MultiDocumentHighlightParams {
    textDocument: { uri: string; };
    position: { line: number; character: number; };
    filesToSearch: string[];
}

interface MultiDocumentHighlightItem {
    uri: string;
    highlights: { range: { start: { line: number; character: number; }; end: { line: number; character: number; }; }; kind?: number; }[];
}

class MultiDocumentHighlightProvider implements vscode.MultiDocumentHighlightProvider {
    constructor(private readonly client: LanguageClient) {}

    async provideMultiDocumentHighlights(
        document: vscode.TextDocument,
        position: vscode.Position,
        otherDocuments: vscode.TextDocument[],
        token: vscode.CancellationToken,
    ): Promise<vscode.MultiDocumentHighlight[]> {
        const allFiles = [document, ...otherDocuments]
            .map(doc => this.client.code2ProtocolConverter.asUri(doc.uri))
            .filter(file => !!file);

        if (allFiles.length === 0) {
            return [];
        }

        const params: MultiDocumentHighlightParams = {
            textDocument: this.client.code2ProtocolConverter.asTextDocumentIdentifier(document),
            position: this.client.code2ProtocolConverter.asPosition(position),
            filesToSearch: allFiles,
        };

        let response: MultiDocumentHighlightItem[] | null;
        try {
            response = await this.client.sendRequest<MultiDocumentHighlightItem[] | null>(multiDocumentHighlightMethod, params, token);
        }
        catch (error) {
            return [];
        }

        if (!response || token.isCancellationRequested) {
            return [];
        }

        return response.map(item =>
            new vscode.MultiDocumentHighlight(
                vscode.Uri.parse(item.uri),
                item.highlights.map(h =>
                    new vscode.DocumentHighlight(
                        new vscode.Range(
                            new vscode.Position(h.range.start.line, h.range.start.character),
                            new vscode.Position(h.range.end.line, h.range.end.character),
                        ),
                        h.kind === 3 ? vscode.DocumentHighlightKind.Write : vscode.DocumentHighlightKind.Read,
                    )
                ),
            )
        );
    }
}

export function registerMultiDocumentHighlightFeature(
    selector: vscode.DocumentSelector,
    client: LanguageClient,
): vscode.Disposable {
    const capabilities = client.initializeResult?.capabilities as { customMultiDocumentHighlightProvider?: boolean; } | undefined;
    if (!capabilities?.customMultiDocumentHighlightProvider) {
        return { dispose() {} };
    }
    return vscode.languages.registerMultiDocumentHighlightProvider(selector, new MultiDocumentHighlightProvider(client));
}
