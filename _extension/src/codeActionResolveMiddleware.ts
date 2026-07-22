import * as vscode from "vscode";

import type {
    CodeActionMiddleware,
    FormattingOptions,
} from "vscode-languageclient/node";

type CodeActionDataWithFormattingOptions = {
    uri?: string;
    formattingOptions?: FormattingOptions;
};

type ProtocolCodeActionWithWritableData = vscode.CodeAction & {
    data?: CodeActionDataWithFormattingOptions;
};

const sourceRemoveUnusedImportsKind = vscode.CodeActionKind.Source.append("removeUnusedImports");
const sourceSortImportsKind = vscode.CodeActionKind.Source.append("sortImports");

function isOrganizeImportsKind(kind: vscode.CodeActionKind): boolean {
    return vscode.CodeActionKind.SourceOrganizeImports.contains(kind)
        || sourceRemoveUnusedImportsKind.contains(kind)
        || sourceSortImportsKind.contains(kind);
}

function getEditor(uri: string): vscode.TextEditor | undefined {
    const activeEditor = vscode.window.activeTextEditor;
    return activeEditor?.document.uri.toString() === uri
        ? activeEditor
        : vscode.window.visibleTextEditors.find(candidate => candidate.document.uri.toString() === uri);
}

function getFormattingOptions(editor: vscode.TextEditor): FormattingOptions | undefined {
    if (typeof editor.options.tabSize !== "number" || typeof editor.options.insertSpaces !== "boolean") {
        return undefined;
    }

    return {
        tabSize: editor.options.tabSize,
        insertSpaces: editor.options.insertSpaces,
    };
}

export const codeActionResolveMiddleware: CodeActionMiddleware = {
    resolveCodeAction(item, token, next) {
        if (item.kind === undefined || !isOrganizeImportsKind(item.kind)) {
            return next(item, token);
        }

        const codeAction = item as ProtocolCodeActionWithWritableData;
        const uri = codeAction.data?.uri;
        if (uri === undefined) {
            return next(item, token);
        }

        const editor = getEditor(uri);
        if (editor === undefined) {
            return next(item, token);
        }

        const formattingOptions = getFormattingOptions(editor);
        if (formattingOptions === undefined) {
            return next(item, token);
        }

        codeAction.data = {
            ...codeAction.data,
            formattingOptions,
        };
        return next(item, token);
    },
};
