export interface ParsedCommandLine {
    options: Record<string, unknown>;
    fileNames: string[];
}

export interface ProjectData {
    configFileName: string;
    commandLine: ParsedCommandLine;
}

export interface SymbolData {
    name: string;
    flags: number;
    checkFlags: number;
    projectVersion: number;
}
