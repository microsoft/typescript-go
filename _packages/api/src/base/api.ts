/**
 * Base interfaces for the TypeScript API client.
 *
 * These interfaces use the `Async` type parameter to support both synchronous
 * and asynchronous implementations:
 * - `Async = false`: Methods return values directly (sync)
 * - `Async = true`: Methods return Promise<T> (async)
 */

import { SymbolFlags } from "#symbolFlags";
import { TypeFlags } from "#typeFlags";
import type {
    Node,
    Path,
    SourceFile,
    SyntaxKind,
} from "@typescript/ast";
import {
    documentURIToFileName,
    fileNameToDocumentURI,
} from "../path.ts";
import type {
    ConfigResponse,
    ProjectChanges,
    ProjectResponse,
} from "../proto.ts";
import type { MaybeAsync } from "./types.ts";

export { SymbolFlags, TypeFlags };

/**
 * A document identifier that can be either a file name (path) or a document URI.
 *
 * @example
 * // Using a file name
 * project.getSourceFile({ fileName: "/path/to/file.ts" });
 *
 * // Using a URI
 * project.getSourceFile({ uri: "file:///path/to/file.ts" });
 */
export type DocumentIdentifier = { fileName: string; } | { uri: string; };

/**
 * A position within a document, combining a document identifier with an offset.
 */
export interface DocumentPosition {
    /** The document containing the position */
    document: DocumentIdentifier | string;
    /** The character offset within the document */
    position: number;
}

/**
 * Resolves a DocumentIdentifier to a file name.
 * If the identifier contains a URI, it is converted to a file name.
 */
export function resolveFileName(identifier: DocumentIdentifier | string): string {
    if (typeof identifier === "string") {
        return identifier;
    }
    if ("fileName" in identifier) {
        return identifier.fileName;
    }
    return documentURIToFileName(identifier.uri);
}

/**
 * Resolves a DocumentIdentifier to a document URI.
 * If the identifier contains a file name, it is converted to a URI.
 */
export function resolveDocumentURI(identifier: DocumentIdentifier | string): string {
    if (typeof identifier === "string") {
        return fileNameToDocumentURI(identifier);
    }
    if ("uri" in identifier) {
        return identifier.uri;
    }
    return fileNameToDocumentURI(identifier.fileName);
}

/**
 * Options for creating an API instance.
 */
export interface APIOptions {
    /** Path to the tsgo executable */
    tsserverPath: string;
    /** Current working directory */
    cwd?: string;
    /** Path to log file for debugging */
    logFile?: string;
}

/**
 * Base interface for the TypeScript API.
 */
export interface API<Async extends boolean> {
    /**
     * Parse a tsconfig.json file.
     */
    parseConfigFile(file: DocumentIdentifier | string): MaybeAsync<Async, ConfigResponse>;

    /**
     * Load a TypeScript project from a tsconfig.json file.
     */
    loadProject(configFile: DocumentIdentifier | string): MaybeAsync<Async, Project<Async>>;

    /**
     * Close the API connection and release all resources.
     */
    close(): MaybeAsync<Async, void>;
}

/**
 * Base interface for a TypeScript project.
 */
export interface Project<Async extends boolean> {
    /** Unique identifier for this project */
    readonly id: string;
    /** Path to the tsconfig.json file */
    configFileName: string;
    /** Compiler options from the config file */
    compilerOptions: Record<string, unknown>;
    /** Root files included in the project */
    rootFiles: readonly string[];

    /**
     * Load project data from a response.
     */
    loadData(data: ProjectResponse): void;

    /**
     * Reload the project, picking up any file changes.
     * @returns File changes if the project was previously loaded, or undefined.
     */
    reload(): MaybeAsync<Async, ProjectChanges | undefined>;

    /**
     * Get a source file from the project by file name or URI.
     */
    getSourceFile(file: DocumentIdentifier | string): MaybeAsync<Async, SourceFile | undefined>;

    /**
     * Get the symbol at a specific location in a source file.
     */
    getSymbolAtLocation(node: Node): MaybeAsync<Async, Symbol<Async> | undefined>;
    getSymbolAtLocation(nodes: readonly Node[]): MaybeAsync<Async, (Symbol<Async> | undefined)[]>;

    /**
     * Get the symbol at a specific position in a file.
     */
    getSymbolAtPosition(file: DocumentIdentifier | string, position: number): MaybeAsync<Async, Symbol<Async> | undefined>;
    getSymbolAtPosition(file: DocumentIdentifier | string, positions: readonly number[]): MaybeAsync<Async, (Symbol<Async> | undefined)[]>;

    /**
     * Get the type of a symbol.
     */
    getTypeOfSymbol(symbol: Symbol<Async>): MaybeAsync<Async, Type<Async> | undefined>;
    getTypeOfSymbol(symbols: readonly Symbol<Async>[]): MaybeAsync<Async, (Type<Async> | undefined)[]>;

    /**
     * Resolve a name to a symbol at a given location.
     * @param name The name to resolve
     * @param meaning Symbol flags indicating what kind of symbol to look for
     * @param location Optional node or document position for location context
     * @param excludeGlobals Whether to exclude global symbols
     */
    resolveName(
        name: string,
        meaning: SymbolFlags,
        location?: Node | DocumentPosition,
        excludeGlobals?: boolean,
    ): MaybeAsync<Async, Symbol<Async> | undefined>;
}

export interface NodeHandle<Async extends boolean> {
    readonly kind: SyntaxKind;
    readonly pos: number;
    readonly end: number;
    readonly path: Path;

    /**
     * Resolve this handle to the actual AST node.
     * @param project The project context to use for fetching the source file
     * @returns The resolved node, or undefined if not found
     */
    resolve(project: Project<Async>): MaybeAsync<Async, Node | undefined>;
}

/**
 * Base interface for a TypeScript symbol.
 */
export interface Symbol<Async extends boolean> {
    /** Unique identifier for this symbol */
    readonly id: string;
    /** Name of the symbol */
    readonly name: string;
    /** Symbol flags */
    readonly flags: SymbolFlags;
    /** Check flags */
    readonly checkFlags: number;
    /** Node handles for declarations of this symbol */
    readonly declarations: readonly NodeHandle<Async>[];
    /** Node handle for the value declaration of this symbol */
    readonly valueDeclaration: NodeHandle<Async> | undefined;
}

/**
 * Base interface for a TypeScript type.
 */
export interface Type<Async extends boolean> {
    /** Unique identifier for this type */
    readonly id: string;
    /** Type flags */
    readonly flags: TypeFlags;
}
