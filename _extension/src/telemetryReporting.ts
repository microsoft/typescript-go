import type { TelemetryReporter as VSCodeTelemetryReporter } from "@vscode/extension-telemetry";

export interface TelemetryReporter {
    publicLog2: LogTelemetrySignature;
    publicLogError2: LogTelemetrySignature;
}

export function createTelemetryReporter(vscReporter: VSCodeTelemetryReporter): TelemetryReporter {
    return {
        publicLog2(eventName, data) {
            vscReporter.sendTelemetryEvent(eventName, data);
        },
        publicLogError2(eventName, data) {
            vscReporter.sendTelemetryErrorEvent(eventName, data);
        },
    };
}

type LogTelemetrySignature = <
    E extends ClassifiedEvent<OmitMetadata<T>> = never,
    T extends IGDPRProperty = never,
>(
    eventName: string,
    data?: StrictPropertyCheck<T, E>,
) => void;

// Note that all of the following types are defined as type aliases rather than interfaces
// because of how object type literals are related with respect to index signatures in TypeScript.

export type LanguageServerStart = {
    version: string;
};

export type LanguageServerStartClassification = {
    owner: "joj";
    comment: "Event emitted when the TypeScript language server starts";
    version: string;
};

export type LanguageServerConnectionError = {
    causedServerShutdown: boolean;
};

export type LanguageServerConnectionErrorClassification = {
    owner: "joj";
    comment: "Event emitted when the TypeScript language server encounters a connection error";
    causedServerShutdown: { classification: "SystemMetaData"; purpose: "PerformanceAndHealth"; comment: "Whether the error caused the language server to shut down"; };
};

export type LanguageServerConnectionClosed = {
    exceededMaxRestarts: boolean;
};

export type LanguageServerConnectionClosedClassification = {
    owner: "joj";
    comment: "Event emitted when the TypeScript language server encounters a connection error";
    exceededMaxRestarts: { classification: "SystemMetaData"; purpose: "PerformanceAndHealth"; comment: "Whether the language server closed enough times such that it restarted"; };
};

export type LanguageServerErrorResponse = {
    errorCode: string;
    requestMethod: string;
    stack: string;
};

export type LanguageServerErrorResponseClassification = {
    owner: "joj";
    comment: "Event emitted when the TypeScript language server returns an error response";
    errorCode: { classification: "CallstackOrException"; purpose: "PerformanceAndHealth"; comment: "The error code returned by the language server"; };
    requestMethod: { classification: "SystemMetaData"; purpose: "PerformanceAndHealth"; comment: "The method of the request that caused the error"; };
    stack: { classification: "CallstackOrException"; purpose: "PerformanceAndHealth"; comment: "The callstack of the error"; };
};

export type EnableNativePreview = {};

export type EnableNativePreviewClassification = {
    owner: "joj";
    comment: "Event emitted when the user enables TypeScript Native Preview";
};

export type DisableNativePreview = {};

export type DisableNativePreviewClassification = {
    owner: "joj";
    comment: "Event emitted when the user disables TypeScript Native Preview";
};

export type RestartLanguageServer = {};

export type RestartLanguageServerClassification = {
    owner: "joj";
    comment: "Event emitted when the user restarts the TypeScript Native Preview language server";
};

export type ReportIssue = {};

export type ReportIssueClassification = {
    owner: "joj";
    comment: "Event emitted when the user decides to report an issue on TypeScript Native Preview through the command palette.";
};

// The following types are from
// https://github.com/microsoft/vscode-telemetry-extractor/blob/eb6c1452b7c4b78141ec090105be7127efafb572/documentation/typescript-code-annotations.md

export interface IPropertyData {
    classification: "SystemMetaData" | "CallstackOrException" | "CustomerContent" | "PublicNonPersonalData" | "EndUserPseudonymizedInformation";
    purpose: "PerformanceAndHealth" | "FeatureInsight" | "BusinessInsight";
    comment: string;
    expiration?: string;
    endpoint?: string;
    isMeasurement?: boolean;
}

export interface IGDPRProperty {
    owner: string;
    comment: string;
    expiration?: string;
    readonly [name: string]: IPropertyData | undefined | IGDPRProperty | string;
}

type IGDPRPropertyWithoutMetadata = Omit<IGDPRProperty, "owner" | "comment" | "expiration">;
export type OmitMetadata<T> = Omit<T, "owner" | "comment" | "expiration">;

export type ClassifiedEvent<T extends IGDPRPropertyWithoutMetadata> = {
    [k in keyof T]: any;
};

export type StrictPropertyChecker<TEvent, TClassification, TError> = keyof TEvent extends keyof OmitMetadata<TClassification> ? keyof OmitMetadata<TClassification> extends keyof TEvent ? TEvent : TError : TError;

export type StrictPropertyCheckError = { error: "Type of classified event does not match event properties"; };

export type StrictPropertyCheck<T extends IGDPRProperty, E> = StrictPropertyChecker<E, ClassifiedEvent<OmitMetadata<T>>, StrictPropertyCheckError>;
