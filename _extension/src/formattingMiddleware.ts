import * as vscode from "vscode";

import type { MessageSignature } from "vscode-languageclient/node";

interface CodeActionParams {
    textDocument?: {
        uri?: string;
    };
    formattingOptions?: {
        tabSize: number;
        insertSpaces: boolean;
    };
}

export function sendRequestMiddleware<P, R>(
    type: string | MessageSignature,
    param: P | undefined,
    token: vscode.CancellationToken | undefined,
    next: (type: string | MessageSignature, param?: P, token?: vscode.CancellationToken) => Promise<R>,
): Promise<R> {
    const method = typeof type === "string" ? type : type.method;
    if (method !== "textDocument/codeAction" || param === undefined || typeof param !== "object") {
        return next(type, param, token);
    }

    const codeActionParams = param as CodeActionParams;
    const uri = codeActionParams.textDocument?.uri;
    if (uri === undefined) {
        return next(type, param, token);
    }

    const editor = vscode.window.visibleTextEditors.find(
        candidate => candidate.document.uri.toString() === uri,
    );
    if (
        editor === undefined
        || typeof editor.options.tabSize !== "number"
        || typeof editor.options.insertSpaces !== "boolean"
    ) {
        return next(type, param, token);
    }

    return next(type, {
        ...codeActionParams,
        formattingOptions: {
            tabSize: editor.options.tabSize,
            insertSpaces: editor.options.insertSpaces,
        },
    } as P, token);
}
