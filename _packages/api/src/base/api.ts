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
    SourceFile,
} from "@typescript/ast";
import type {
    ConfigResponse,
    ProjectResponse,
} from "../proto.ts";
import type { MaybeAsync } from "./types.ts";

export { SymbolFlags, TypeFlags };

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
    parseConfigFile(fileName: string): MaybeAsync<Async, ConfigResponse>;

    /**
     * Load a TypeScript project from a tsconfig.json file.
     */
    loadProject(configFileName: string): MaybeAsync<Async, Project<Async>>;

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
     */
    reload(): MaybeAsync<Async, void>;

    /**
     * Get a source file from the project by file name.
     */
    getSourceFile(fileName: string): MaybeAsync<Async, SourceFile | undefined>;

    /**
     * Get the symbol at a specific location in a source file.
     */
    getSymbolAtLocation(node: Node): MaybeAsync<Async, Symbol<Async> | undefined>;

    /**
     * Get symbols at multiple locations.
     */
    getSymbolsAtLocations(nodes: readonly Node[]): MaybeAsync<Async, (Symbol<Async> | undefined)[]>;

    /**
     * Get the symbol at a specific position in a file.
     */
    getSymbolAtPosition(fileName: string, position: number): MaybeAsync<Async, Symbol<Async> | undefined>;

    /**
     * Get symbols at multiple positions in a file.
     */
    getSymbolsAtPositions(fileName: string, positions: readonly number[]): MaybeAsync<Async, (Symbol<Async> | undefined)[]>;

    /**
     * Get the type of a symbol.
     */
    getTypeOfSymbol(symbol: Symbol<Async>): MaybeAsync<Async, Type<Async> | undefined>;

    /**
     * Get types of multiple symbols.
     */
    getTypesOfSymbols(symbols: readonly Symbol<Async>[]): MaybeAsync<Async, (Type<Async> | undefined)[]>;
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
