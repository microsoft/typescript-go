//// [tests/cases/compiler/declarationEmitIndexedAccessPreservation.ts] ////

//// [declarationEmitIndexedAccessPreservation.ts]
export interface Transaction {
  readonly Service: {
    retry: boolean
    readonly journal: Map<unknown, { readonly version: number; value: unknown }>
  }
}

export interface Config {
  readonly Options: {
    readonly timeout: number
    readonly retries: number
  }
}

// Indexed access should be preserved in declaration emit
export const withTransaction = <A>(
  f: (state: Transaction["Service"]) => A
): A => f({ retry: false, journal: new Map() })

export const getOptions = (): Config["Options"] => ({ timeout: 1000, retries: 3 })

// Nested indexed access
export interface Nested {
  readonly Level1: {
    readonly Level2: {
      readonly value: string
    }
  }
}

export const getDeep = (): Nested["Level1"]["Level2"] => ({ value: "test" })


//// [declarationEmitIndexedAccessPreservation.js]
// Indexed access should be preserved in declaration emit
export const withTransaction = (f) => f({ retry: false, journal: new Map() });
export const getOptions = () => ({ timeout: 1000, retries: 3 });
export const getDeep = () => ({ value: "test" });


//// [declarationEmitIndexedAccessPreservation.d.ts]
export interface Transaction {
    readonly Service: {
        retry: boolean;
        readonly journal: Map<unknown, {
            readonly version: number;
            value: unknown;
        }>;
    };
}
export interface Config {
    readonly Options: {
        readonly timeout: number;
        readonly retries: number;
    };
}
export declare const withTransaction: <A>(f: (state: Transaction["Service"]) => A) => A;
export declare const getOptions: () => Config["Options"];
export interface Nested {
    readonly Level1: {
        readonly Level2: {
            readonly value: string;
        };
    };
}
export declare const getDeep: () => Nested["Level1"]["Level2"];
