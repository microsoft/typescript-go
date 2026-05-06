//// [tests/cases/compiler/expandoFunctionSymbolProperty.ts] ////

//// [expandoFunctionSymbolProperty.ts]
// repro from https://github.com/microsoft/TypeScript/issues/54220

const symb = Symbol();

interface TestSymb {
  (): void;
  readonly [symb]: boolean;
}

export function test(): TestSymb {
  function inner() {}
  inner[symb] = true;
  return inner;
}


//// [expandoFunctionSymbolProperty.js]
// repro from https://github.com/microsoft/TypeScript/issues/54220
const symb = Symbol();
export function test() {
    function inner() { }
    inner[symb] = true;
    return inner;
}


//// [expandoFunctionSymbolProperty.d.ts]
const symb: unique symbol;
interface TestSymb {
    (): void;
    readonly [symb]: boolean;
}
export function test(): TestSymb;
export {};


//// [DtsFileErrors]


expandoFunctionSymbolProperty.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== expandoFunctionSymbolProperty.d.ts (1 errors) ====
    const symb: unique symbol;
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    interface TestSymb {
        (): void;
        readonly [symb]: boolean;
    }
    export function test(): TestSymb;
    export {};
    