//// [tests/cases/compiler/declarationEmitExportValueSymbolWithTsIgnore.ts] ////

//// [index.d.ts]
export declare const MySymbol: unique symbol;
export declare function createService<T>(): {
    new (): {
        [MySymbol](): T | undefined;
    };
};

//// [client.ts]
// @ts-ignore Import needed for type visibility but appears unused
import { MySymbol } from "lib";
import { createService } from "lib";

// The extends clause references the factory result which uses MySymbol
// This should trigger symbol accessibility check for MySymbol
export class Client extends createService<string>() {
    doSomething(): string {
        return "hello";
    }
}


//// [client.js]
import { createService } from "lib";
// The extends clause references the factory result which uses MySymbol
// This should trigger symbol accessibility check for MySymbol
export class Client extends createService() {
    doSomething() {
        return "hello";
    }
}


//// [client.d.ts]
import { MySymbol } from "lib";
const Client_base: new () => {
    [MySymbol](): string | undefined;
};
export class Client extends Client_base {
    doSomething(): string;
}
export {};


//// [DtsFileErrors]


client.d.ts(2,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== node_modules/lib/index.d.ts (0 errors) ====
    export declare const MySymbol: unique symbol;
    export declare function createService<T>(): {
        new (): {
            [MySymbol](): T | undefined;
        };
    };
    
==== client.d.ts (1 errors) ====
    import { MySymbol } from "lib";
    const Client_base: new () => {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        [MySymbol](): string | undefined;
    };
    export class Client extends Client_base {
        doSomething(): string;
    }
    export {};
    