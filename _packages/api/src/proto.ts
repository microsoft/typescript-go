/**
 * Response from the initialize method.
 */
export interface InitializeResponse {
    /** Whether the host file system is case-sensitive */
    useCaseSensitiveFileNames: boolean;
    /** The server's current working directory */
    currentDirectory: string;
}

export interface ConfigResponse {
    options: Record<string, unknown>;
    fileNames: string[];
}

export interface ProjectResponse {
    id: string;
    configFileName: string;
    compilerOptions: Record<string, unknown>;
    rootFiles: string[];
    /** Encodes the source-file-independent parse options for this project (for cache keying) */
    parseOptionsKey: string;
}

/**
 * Response from loadProject, includes project info plus optional changes.
 */
export interface LoadProjectResponse extends ProjectResponse {
    /** File changes if the project was previously loaded */
    changes?: ProjectChanges;
}

export interface SourceFileResponse {
    /** Base64-encoded binary AST data */
    data: string;
}

export interface SymbolResponse {
    id: string;
    name: string;
    flags: number;
    checkFlags: number;
    declarations?: string[];
    valueDeclaration?: string;
}

export interface TypeResponse {
    id: string;
    flags: number;
}

/**
 * Response from adoptLSPState with cache invalidation information.
 */
export interface AdoptLSPStateResponse {
    /** List of project handles that no longer exist */
    removedProjects?: string[];
    /** Map of project handles to their file changes */
    projectChanges?: Record<string, ProjectChanges>;
}

/**
 * Describes file changes within a project.
 */
export interface ProjectChanges {
    /** List of file paths whose content has changed */
    changedFiles?: string[];
    /** List of file paths that no longer exist in the project */
    removedFiles?: string[];
}
