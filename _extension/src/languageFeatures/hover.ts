import * as vscode from "vscode";
import {
    Hover,
    HoverRequest,
    LanguageClient,
    MarkupContent,
    TextDocumentPositionParams,
} from "vscode-languageclient/node";

interface HoverResult extends Hover {
    canIncreaseVerbosity?: boolean;
    canDecreaseVerbosity?: boolean;
}

interface HoverParamsWithVerbosity extends TextDocumentPositionParams {
    verbosityLevel?: number;
}

class VerboseHoverProvider implements vscode.HoverProvider {
    private lastHoverAndLevel: [vscode.Hover, number] | undefined;

    constructor(private readonly client: LanguageClient) {}

    async provideHover(
        document: vscode.TextDocument,
        position: vscode.Position,
        token: vscode.CancellationToken,
        context?: vscode.HoverContext,
    ): Promise<vscode.VerboseHover | undefined> {
        const verbosityLevel = Math.max(0, this.getPreviousLevel(context?.previousHover) + (context?.verbosityDelta ?? 0));

        const params: HoverParamsWithVerbosity = {
            ...this.client.code2ProtocolConverter.asTextDocumentPositionParams(document, position),
            verbosityLevel,
        };

        let response: HoverResult | null;
        try {
            response = await this.client.sendRequest(HoverRequest.type, params, token);
        }
        catch (error) {
            return this.client.handleFailedRequest(HoverRequest.type, token, error, null) ?? undefined;
        }

        if (!response || token.isCancellationRequested) {
            return undefined;
        }

        const markupContent = response.contents as MarkupContent;
        const contents = new vscode.MarkdownString(markupContent.value);
        contents.isTrusted = true;
        contents.supportHtml = true;

        let range: vscode.Range | undefined;
        if (response.range) {
            range = this.client.protocol2CodeConverter.asRange(response.range);
        }

        const hover = new vscode.VerboseHover(
            [contents],
            range,
            response.canIncreaseVerbosity,
            verbosityLevel > 0,
        );

        this.lastHoverAndLevel = [hover, verbosityLevel];
        return hover;
    }

    private getPreviousLevel(previousHover: vscode.Hover | undefined): number {
        if (previousHover && this.lastHoverAndLevel && this.lastHoverAndLevel[0] === previousHover) {
            return this.lastHoverAndLevel[1];
        }
        return 0;
    }
}

export function registerHoverFeature(
    selector: vscode.DocumentSelector,
    client: LanguageClient,
): vscode.Disposable {
    return vscode.languages.registerHoverProvider(selector, new VerboseHoverProvider(client));
}
