import * as vscode from "vscode";
import type { CancellationToken } from "vscode";
import type { MessageSignature } from "vscode-languageserver-protocol";
import {
    disabledSchemes,
    isSupportedLanguageMode,
    readUnifiedConfig,
} from "./util";

function getDocument(): vscode.TextDocument | undefined {
    const activeDocument = vscode.window.activeTextEditor?.document;
    if (activeDocument && isSupportedLanguageMode(activeDocument) && !disabledSchemes.has(activeDocument.uri.scheme)) {
        return activeDocument;
    }

    return vscode.workspace.textDocuments.find(
        document => isSupportedLanguageMode(document) && !disabledSchemes.has(document.uri.scheme),
    );
}

export function workspaceSymbolSendRequestMiddleware<P, R>(
    type: string | MessageSignature,
    params: P | undefined,
    token: CancellationToken | undefined,
    next: (type: string | MessageSignature, params?: P, token?: CancellationToken) => Promise<R>,
): Promise<R> {
    const method = typeof type === "string" ? type : type.method;
    if (method !== "workspace/symbol") {
        return next(type, params, token);
    }

    const scope = readUnifiedConfig<"allOpenProjects" | "currentProject">(
        "workspaceSymbols.scope",
        "typescript",
        "workspaceSymbols.scope",
        undefined,
        "allOpenProjects",
    );
    if (scope !== "currentProject") {
        return next(type, params, token);
    }

    const document = getDocument();
    if (!document) {
        return next(type, params, token);
    }

    return next(type, {
        ...params,
        textDocument: { uri: document.uri.toString() },
    } as P, token);
}
