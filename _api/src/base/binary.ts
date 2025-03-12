import type {
    GetSymbolAtPositionRequest,
    GetTypeOfSymbolRequest,
    SymbolResponse,
} from "./proto.ts";

export function encodeGetSymbolAtPositionRequest(projectId: number, fileName: string, position: number, encoder: TextEncoder): Uint8Array {
    // assume ASCII filename
    const asciiLength = fileName.length;
    const result = new Uint8Array(2 + 4 + asciiLength);
    const view = new DataView(result.buffer);
    view.setUint16(0, projectId, true);
    view.setUint32(2, position, true);
    const { read } = encoder.encodeInto(fileName, result.subarray(2 + 4));
    // check if ASCII assumption was correct
    if (read !== asciiLength) {
        const encodedFileName = encoder.encode(fileName);
        const newResult = new Uint8Array(2 + 4 + encodedFileName.length);
        newResult.set(result.subarray(0, 2 + 4));
        newResult.set(encodedFileName, 2 + 4);
        return newResult;
    }

    return result;
}

export function encodeGetTypeOfSymbolRequest(projectId: number, symbolId: number): Uint8Array {
    const result = new Uint8Array(6);
    const view = new DataView(result.buffer);
    view.setUint16(0, projectId, true);
    view.setUint32(2, symbolId, true);
    return result;
}
