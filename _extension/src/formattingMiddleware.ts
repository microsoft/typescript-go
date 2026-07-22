import * as vscode from "vscode";

import type {
    CodeActionParams,
    FormattingOptions,
    MessageSignature,
} from "vscode-languageclient/node";

type CodeActionParamsWithFormattingOptions = Partial<CodeActionParams> & { formattingOptions?: FormattingOptions; };

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

    const codeActionParams = param as CodeActionParamsWithFormattingOptions;
    const uri = codeActionParams.textDocument?.uri;
    if (uri === undefined) {
        return next(type, param, token);
    }

    const activeEditor = vscode.window.activeTextEditor;
    const editor = activeEditor?.document.uri.toString() === uri
        ? activeEditor
        : vscode.window.visibleTextEditors.find(candidate => candidate.document.uri.toString() === uri);
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
