//// [tests/cases/compiler/jsFileCompilationAwaitModifier.ts] ////

//// [a.js]
class Foo {
    async a() {
        await Promise.resolve(1);
    }

    b = async () => {
        await Promise.resolve(1);
    }
}


//// [a.js]
"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
class Foo {
    constructor() {
        this.b = () => __awaiter(this, void 0, void 0, function* () {
            yield Promise.resolve(1);
        });
    }
    a() {
        return __awaiter(this, void 0, void 0, function* () {
            yield Promise.resolve(1);
        });
    }
}


//// [a.d.ts]
class Foo {
    a(): Promise<void>;
    b: () => Promise<void>;
}


//// [DtsFileErrors]


out/a.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== out/a.d.ts (1 errors) ====
    class Foo {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        a(): Promise<void>;
        b: () => Promise<void>;
    }
    