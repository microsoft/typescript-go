export type LSPErrorType = "MessageFormat" | "Parse" | "Server";

export class LSPError extends Error {
    code?: number;
    raw?: string;
    constructor(
        message: string,
        type: LSPErrorType,
        details?: { code?: number; raw?: string; },
    ) {
        super(message);
        this.name = `LSP${type}Error`;
        this.code = details?.code;
        this.raw = details?.raw;
    }
}

export class LSPErrorCollection extends Error {
    errors: LSPError[];
    constructor(
        errors: LSPError[],
    ) {
        super("LSP Errors:\n" + errors.map(e => e.message).join("\n"));
        this.name = "LSPErrorCollection";
        this.errors = errors;
    }
}
