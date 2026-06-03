import * as vscode from "vscode";

import {
    ClientCapabilities,
    CloseAction,
    CloseHandlerResult,
    ErrorAction,
    ErrorHandler,
    ErrorHandlerResult,
    LanguageClient,
    LanguageClientOptions,
    Message,
    NotebookDocumentFilter,
    ServerOptions,
    StaticFeature,
    TextDocumentFilter,
    TransportKind,
} from "vscode-languageclient/node";

import { codeLensShowLocationsCommandName } from "./commands";
import {
    configurationMiddleware,
    sendNotificationMiddleware,
} from "./configurationMiddleware";
import { registerMultiDocumentHighlightFeature } from "./languageFeatures/documentHighlight";
import { registerHoverFeature } from "./languageFeatures/hover";
import { registerOnAutoInsertFeature } from "./languageFeatures/onAutoInsert";
import { registerSourceDefinitionFeature } from "./languageFeatures/sourceDefinition";
import * as tr from "./telemetryReporting";
import {
    ExeInfo,
    getExe,
    jsTsLanguageModes,
} from "./util";
import { getLanguageForUri } from "./util";

export class Client implements vscode.Disposable {
    private outputChannel: vscode.LogOutputChannel;
    private traceOutputChannel: vscode.LogOutputChannel;
    private initializedEventEmitter: vscode.EventEmitter<void>;
    private telemetryReporter: tr.TelemetryReporter;

    private documentSelector: Array<{ scheme: string; language: string; }>;
    private client?: LanguageClient;

    private isDisposed = false;
    private disposables: vscode.Disposable[] = [];
    isInitialized = false;

    private exe: ExeInfo | undefined;
    private errorHandler: ReportingErrorHandler | undefined;

    private reporterCommonProperties: tr.LSCommonProperties | undefined;

    constructor(
        outputChannel: vscode.LogOutputChannel,
        traceOutputChannel: vscode.LogOutputChannel,
        initializedEventEmitter: vscode.EventEmitter<void>,
        telemetryReporter: tr.TelemetryReporter,
    ) {
        this.outputChannel = outputChannel;
        this.traceOutputChannel = traceOutputChannel;
        this.initializedEventEmitter = initializedEventEmitter;
        this.telemetryReporter = telemetryReporter;

        this.documentSelector = [
            ...jsTsLanguageModes.map(language => ({ scheme: "file", language })),
            ...jsTsLanguageModes.map(language => ({ scheme: "untitled", language })),
        ];
    }

    async start(exe: { path: string; version: string; }): Promise<void> {
        this.exe = exe;
        this.reporterCommonProperties = {
            "tscommon.version": exe.version,
            "tscommon.serverSessionId": `${Date.now()}`,
        };

        this.errorHandler = new ReportingErrorHandler(this.telemetryReporter, 5, this.reporterCommonProperties);

        // Monkey-patch the output channel's error method to capture recent stderr lines.
        // When the server crashes, vscode-languageclient pipes stderr to outputChannel.error(),
        // so the error handler can include the last N lines in crash telemetry.
        const unboundOriginalError = this.outputChannel.error;
        const originalError = unboundOriginalError.bind(this.outputChannel);
        this.outputChannel.error = (...args: Parameters<typeof this.outputChannel.error>) => {
            originalError(...args);
            this.errorHandler!.pushStderrLine(String(args[0]));
        };
        // Remember to dispose here so that we don't keep stacking references to the original Client instance.
        this.disposables.push({
            dispose: () => {
                this.outputChannel.error = unboundOriginalError;
            },
        });

        this.outputChannel.appendLine(`Resolved to ${this.exe.path}`);
        this.telemetryReporter.sendTelemetryEvent("languageServer.start", {
            version: exe.version,
            ...this.reporterCommonProperties,
        });

        // Get pprofDir
        const config = vscode.workspace.getConfiguration("typescript.native-preview");
        const pprofDir = config.get<string>("pprofDir");
        const pprofArgs = pprofDir ? ["--pprofDir", pprofDir] : [];

        const goMemLimit = config.get<string>("goMemLimit");
        const env = { ...process.env };
        if (goMemLimit) {
            // Keep this regex aligned with the pattern in package.json.
            if (/^[0-9]+(([KMGT]i)?B)?$/.test(goMemLimit)) {
                this.outputChannel.appendLine(`Setting GOMEMLIMIT=${goMemLimit}`);
                env.GOMEMLIMIT = goMemLimit;
            }
            else {
                this.outputChannel.error(`Invalid goMemLimit: ${goMemLimit}. Must be a valid memory limit (e.g., '2048MiB', '4GiB'). Not overriding GOMEMLIMIT.`);
            }
        }

        const serverOptions: ServerOptions = {
            run: {
                command: this.exe.path,
                args: ["--lsp", ...pprofArgs],
                transport: TransportKind.stdio,
                options: { env },
            },
            debug: {
                command: this.exe.path,
                args: ["--lsp", ...pprofArgs],
                transport: TransportKind.stdio,
                options: { env },
            },
        };

        this.client = new LanguageClient(
            "typescript.native-preview",
            "typescript.native-preview-lsp",
            serverOptions,
            this.makeClientOptions(this.reporterCommonProperties),
        );
        this.disposables.push(this.client);

        // Register a static feature to advertise verbosityLevel support in hover capabilities.
        this.client.registerFeature(
            {
                fillClientCapabilities(capabilities: ClientCapabilities): void {
                    capabilities.textDocument = capabilities.textDocument ?? {};
                    capabilities.textDocument.hover = capabilities.textDocument.hover ?? {};
                    (capabilities.textDocument.hover as { verbosityLevel?: boolean; }).verbosityLevel = true;
                },
                initialize(): void {},
                getState() {
                    return { kind: "static" as const };
                },
                clear(): void {},
            } satisfies StaticFeature,
        );

        this.outputChannel.appendLine(`Starting language server...`);
        await this.client.start();
        this.isInitialized = true;
        this.initializedEventEmitter.fire();

        if (this.traceOutputChannel.logLevel !== vscode.LogLevel.Trace) {
            this.traceOutputChannel.appendLine(`To see LSP trace output, set this output's log level to "Trace" (gear icon next to the dropdown).`);
        }

        type TelemetryData = {
            eventName: string;
            telemetryPurpose: "usage" | "error";
            properties?: Record<string, string>;
            measurements?: Record<string, number>;
        };

        const serverTelemetryListener = this.client.onTelemetry((d: TelemetryData) => {
            switch (d.telemetryPurpose) {
                case "usage":
                    this.telemetryReporter.sendTelemetryEventUntyped(
                        d.eventName,
                        { ...this.reporterCommonProperties, ...d.properties },
                        d.measurements,
                    );
                    break;
                case "error":
                    this.telemetryReporter.sendTelemetryErrorEventUntyped(
                        d.eventName,
                        { ...this.reporterCommonProperties, ...d.properties },
                        d.measurements,
                    );
                    break;
                default:
                    const _: never = d.telemetryPurpose;
                    this.telemetryReporter.sendTelemetryErrorEvent("languageServer.unexpectedTelemetryPurpose", {
                        ...this.reporterCommonProperties,
                        telemetryPurpose: String(d.telemetryPurpose),
                    });
                    break;
            }
        });

        this.disposables.push(
            serverTelemetryListener,
            registerMultiDocumentHighlightFeature(this.documentSelector, this.client),
            registerSourceDefinitionFeature(this.client),
            registerHoverFeature(this.documentSelector, this.client),
            registerOnAutoInsertFeature(this.documentSelector, this.client),
        );
    }

    dispose(): void {
        if (this.isDisposed) {
            return;
        }
        this.isDisposed = true;
        while (this.disposables.length > 0) {
            const disposable = this.disposables.pop()!;
            disposable.dispose();
        }
    }

    getCurrentExe(): ExeInfo | undefined {
        return this.exe;
    }

    private makeClientOptions(reporterCommonProperties: tr.LSCommonProperties): LanguageClientOptions {
        return {
            documentSelector: this.documentSelector,
            outputChannel: this.outputChannel,
            traceOutputChannel: this.traceOutputChannel,
            initializationOptions: {
                codeLensShowLocationsCommandName,
                enableTelemetry: true,
            },
            errorHandler: new ReportingErrorHandler(this.telemetryReporter, 5, reporterCommonProperties),
            middleware: {
                workspace: {
                    ...configurationMiddleware,
                },
                sendNotification: sendNotificationMiddleware,
                provideHover: () => undefined,
            },
            diagnosticCollectionName: "typescript",
            diagnosticPullOptions: {
                onChange: true,
                onSave: true,
                onTabs: true,
                match(documentSelector, resource) {
                    // This function is called when diagnostics are requested but
                    // only the URI itself is known (e.g. open but not yet focused tabs),
                    // so will not be present in vscode.workspace.textDocuments.
                    // See if this file matches without consulting vscode.languages.match
                    // (which requires a TextDocument).

                    const language = getLanguageForUri(resource);

                    for (const selector of documentSelector) {
                        if (typeof selector === "string") {
                            if (selector === language) {
                                return true;
                            }
                            continue;
                        }
                        if (NotebookDocumentFilter.is(selector)) {
                            continue;
                        }
                        if (TextDocumentFilter.is(selector)) {
                            if (selector.language !== undefined && selector.language !== language) {
                                continue;
                            }

                            if (selector.scheme !== undefined && selector.scheme !== resource.scheme) {
                                continue;
                            }

                            if (selector.pattern !== undefined) {
                                // VS Code's glob matcher is not available via the API;
                                // see: https://github.com/microsoft/vscode/issues/237304
                                // But, we're only called on selectors passed above, so just ignore this for now.
                                throw new Error("Not implemented");
                            }

                            return true;
                        }
                    }

                    return false;
                },
            },
        };
    }

    get serverPid(): number | undefined {
        return (this.client as any)?._serverProcess?.pid;
    }

    /**
     * Initialize an API session and return the socket path for connecting.
     * This allows other extensions to get a direct connection to the API server.
     */
    async initializeAPISession(pipe?: string): Promise<{ sessionId: string; pipe: string; }> {
        if (!this.client) {
            throw new Error("Language client is not initialized");
        }
        return this.client.sendRequest<{ sessionId: string; pipe: string; }>("custom/initializeAPISession", { pipe });
    }

    /**
     * Restart the language server if the executable path has not changed.
     * Returns true if a restart was performed.
     */
    async tryRestart(context: vscode.ExtensionContext): Promise<boolean> {
        if (!this.client) {
            return Promise.reject(new Error("Language client is not initialized"));
        }
        const exe = await getExe(context);
        if (exe.path !== this.exe?.path) {
            return false;
        }

        this.isInitialized = false;
        this.outputChannel.appendLine(`Restarting language server...`);
        try {
            await this.client.restart();
        }
        catch (err) {
            this.outputChannel.appendLine(`Graceful shutdown failed, forcing restart: ${err}`);
            await this.client.start();
        }
        this.isInitialized = true;
        this.initializedEventEmitter.fire();
        return true;
    }

    // Developer/debugging methods

    async runGC(): Promise<void> {
        if (!this.client) {
            throw new Error("Language client is not initialized");
        }
        await this.client.sendRequest("custom/runGC");
    }

    async saveHeapProfile(dir: string): Promise<string> {
        if (!this.client) {
            throw new Error("Language client is not initialized");
        }
        const result = await this.client.sendRequest<{ file: string; }>("custom/saveHeapProfile", { dir });
        return result.file;
    }

    async saveAllocProfile(dir: string): Promise<string> {
        if (!this.client) {
            throw new Error("Language client is not initialized");
        }
        const result = await this.client.sendRequest<{ file: string; }>("custom/saveAllocProfile", { dir });
        return result.file;
    }

    async startCPUProfile(dir: string): Promise<void> {
        if (!this.client) {
            throw new Error("Language client is not initialized");
        }
        await this.client.sendRequest("custom/startCPUProfile", { dir });
    }

    async stopCPUProfile(): Promise<string> {
        if (!this.client) {
            throw new Error("Language client is not initialized");
        }
        const result = await this.client.sendRequest<{ file: string; }>("custom/stopCPUProfile");
        return result.file;
    }

    async getProjectInfo(uri: string, token?: vscode.CancellationToken): Promise<{ configFilePath: string; }> {
        if (!this.client) {
            throw new Error("Language client is not initialized");
        }
        return this.client.sendRequest<{ configFilePath: string; }>("custom/projectInfo", {
            textDocument: { uri },
        }, token);
    }
}

// Adapted from the default error handler in vscode-languageclient.
class ReportingErrorHandler implements ErrorHandler {
    private telemetryReporter: tr.TelemetryReporter;
    private maxRestartCount: number;
    private restarts: number[];
    private stderrBuffer: string[] = [];
    private capturingPanic = false;
    private static readonly maxStderrLines = 40;
    private static readonly maxStderrLength = 8192;

    private reporterCommonProperties: tr.LSCommonProperties;

    constructor(telemetryReporter: tr.TelemetryReporter, maxRestartCount: number, reporterCommonProperties: tr.LSCommonProperties) {
        this.telemetryReporter = telemetryReporter;
        this.maxRestartCount = maxRestartCount;
        this.restarts = [];

        this.reporterCommonProperties = reporterCommonProperties;
    }

    pushStderrLine(line: string): void {
        for (const l of line.split("\n")) {
            if (!this.capturingPanic) {
                if (/^panic:/.test(l.trimStart())) {
                    // Clear any stale data from a previous session/panic.
                    this.stderrBuffer = [];
                    this.capturingPanic = true;
                }
                else {
                    continue;
                }
            }
            if (this.stderrBuffer.length < ReportingErrorHandler.maxStderrLines) {
                this.stderrBuffer.push(l);
            }
            else {
                this.capturingPanic = false;
            }
        }
    }

    private consumeStderrBuffer(): string {
        const raw = this.stderrBuffer.join("\n");
        this.stderrBuffer = [];
        this.capturingPanic = false;
        return sanitizeStderr(raw).slice(0, ReportingErrorHandler.maxStderrLength);
    }

    error(_error: Error, _message: Message | undefined, count: number | undefined): ErrorHandlerResult | Promise<ErrorHandlerResult> {
        let errorAction = ErrorAction.Shutdown;
        if (count && count <= 3) {
            errorAction = ErrorAction.Continue;
        }

        let actionString = "";
        switch (errorAction) {
            case ErrorAction.Continue:
                actionString = "continue";
                break;
            case ErrorAction.Shutdown:
                actionString = "shutdown";
                break;
            default:
                const _: never = errorAction;
        }
        this.telemetryReporter.sendTelemetryErrorEvent("languageServer.connectionError", {
            ...this.reporterCommonProperties,
            resultingAction: actionString,
        });

        return { action: errorAction };
    }

    closed(): CloseHandlerResult | Promise<CloseHandlerResult> {
        let resultingAction: CloseAction;

        this.restarts.push(Date.now());
        if (this.restarts.length <= this.maxRestartCount) {
            resultingAction = CloseAction.Restart;
        }
        else {
            const diff = this.restarts[this.restarts.length - 1] - this.restarts[0];
            if (diff <= 3 * 60 * 1000) {
                resultingAction = CloseAction.DoNotRestart;
            }
            else {
                this.restarts.shift();
                resultingAction = CloseAction.Restart;
            }
        }

        let actionString = "";
        switch (resultingAction) {
            case CloseAction.DoNotRestart:
                actionString = "doNotRestart";
                break;
            case CloseAction.Restart:
                actionString = "restart";
                break;
            default:
                const _: never = resultingAction;
        }
        const lastStderr = this.consumeStderrBuffer();
        this.telemetryReporter.sendTelemetryErrorEvent("languageServer.connectionClosed", {
            ...this.reporterCommonProperties,
            resultingAction: actionString,
            lastStderr,
        });

        if (resultingAction === CloseAction.DoNotRestart) {
            return {
                action: resultingAction,
                message: `The typescript.native-preview-lsp server crashed ${this.maxRestartCount + 1} times in the last 3 minutes. The server will not be restarted. See the output for more information.`,
            };
        }

        return { action: resultingAction };
    }
}

// Matches the server-side sanitizeStackTrace in internal/lsp/stack_sanitizer.go.
// Strips file path prefixes that may contain PII and redacts frames outside of our module.
const genericSecretKeywordRegex = /\b(key|token|signature|sig|pwd)([(\[.|])/gi;

function sanitizeStderr(stderr: string): string {
    if (!stderr) {
        return "";
    }
    return stderr.split("\n").map(sanitizeStderrLine).join("\n");
}

function sanitizeStderrLine(line: string): string {
    // Keep "goroutine N [status]:" headers as-is.
    if (/^goroutine \d+/.test(line)) {
        return line;
    }
    // Redact the panic message itself — assert messages may contain user data.
    // Keep only "panic:" as a marker.
    if (/^panic:/.test(line.trimStart())) {
        return "panic: (REDACTED)";
    }
    // Keep "Server process exited" messages from vscode-languageclient.
    if (line.includes("Server process exited")) {
        return line;
    }

    const leadingWhitespace = line.match(/^(\s*)/)?.[1] ?? "";

    // Stack frame file path lines look like: \t/full/path/to/file.go:123 +0x40
    // Function lines look like: github.com/microsoft/typescript-go/internal/foo.Bar(...)
    const ourModuleMarker = "typescript-go/internal";
    const idx = line.indexOf(ourModuleMarker);
    if (idx >= 0) {
        let relevantPart = line.slice(idx);
        // Strip hex offset suffixes like " +0x40"
        relevantPart = relevantPart.replace(/ \+0x[0-9a-fA-F]+$/, "");
        // Strip " in goroutine N" suffixes
        relevantPart = relevantPart.replace(/ in goroutine \d+$/, "");
        // Strip function arguments (keep parens empty)
        relevantPart = relevantPart.replace(/\([^)]*\)$/, "()");
        // Replace / with |> to defeat path-based secret detection
        relevantPart = relevantPart.replace(/\//g, "|>");
        // Defeat generic secret keyword regex
        relevantPart = relevantPart.replace(genericSecretKeywordRegex, "$1X_X$2");
        return leadingWhitespace + relevantPart;
    }

    // Preserve completely blank lines.
    if (line.trim() === "") {
        return "";
    }

    // Non-internal frames get fully redacted.
    return leadingWhitespace + "(REDACTED)";
}
