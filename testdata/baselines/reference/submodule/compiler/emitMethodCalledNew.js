//// [tests/cases/compiler/emitMethodCalledNew.ts] ////

//// [emitMethodCalledNew.ts]
// https://github.com/microsoft/TypeScript/issues/55075

export const a = {
  new(x: number) { return x + 1 }
}
export const b = {
  "new"(x: number) { return x + 1 }
}
export const c = {
  ["new"](x: number) { return x + 1 }
}


//// [emitMethodCalledNew.js]
// https://github.com/microsoft/TypeScript/issues/55075
export const a = {
    new(x) { return x + 1; }
};
export const b = {
    "new"(x) { return x + 1; }
};
export const c = {
    ["new"](x) { return x + 1; }
};


//// [emitMethodCalledNew.d.ts]
export const a: {
    "new"(x: number): number;
};
export const b: {
    "new"(x: number): number;
};
export const c: {
    ["new"](x: number): number;
};
