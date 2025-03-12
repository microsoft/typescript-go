export type MaybeAsync<Async extends boolean, T> = Async extends true ? Promise<T> : T;

export interface FileSystemEntries {
    files: string[];
    directories: string[];
}

export interface ParsedCommandLine {
    options: Record<string, unknown>;
    fileNames: string[];
}

export interface ProjectData {
    id: number;
    configFileName: string;
    compilerOptions: Record<string, unknown>;
    rootFiles: string[];
}

export interface SymbolData {
    id: number;
    name: string;
    flags: number;
    checkFlags: number;
}

export interface TypeData {
    id: number;
    flags: number;
}
