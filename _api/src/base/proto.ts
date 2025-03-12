export interface ConfigResponse {
    options: Record<string, unknown>;
    fileNames: string[];
}

export interface ProjectResponse {
    id: number;
    configFileName: string;
    compilerOptions: Record<string, unknown>;
    rootFiles: string[];
}

export interface GetSymbolAtPositionRequest {
    project: number;
    fileName: string;
    position: number;
}

export interface SymbolResponse {
    id: number;
    name: string;
    flags: number;
    checkFlags: number;
}

export interface GetTypeOfSymbolRequest {
    symbol: number;
}

export interface TypeResponse {
    id: number;
    flags: number;
}
