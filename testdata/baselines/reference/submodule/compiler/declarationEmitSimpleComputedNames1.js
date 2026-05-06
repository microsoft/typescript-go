//// [tests/cases/compiler/declarationEmitSimpleComputedNames1.ts] ////

//// [declarationEmitSimpleComputedNames1.ts]
export const fieldName = Math.random() > 0.5 ? "f1" : "f2";
export const conatainer = {
    [fieldName]() {
        return "result";
    }
};

const classFieldName = Math.random() > 0.5 ? "g1" : "g2";
const otherField = classFieldName === "g1" ? "g2" : "g1";
const staticField = Math.random() > 0.5 ? "s1" : "s2";
export class Holder {
    [classFieldName]() {
        return "value";
    }
    [otherField]() {
        return 42;
    }
    static [staticField]() {
        return { static: true };
    }
    static [staticField]() {
        return { static: "sometimes" };
    }
}

/**
 * Could be `"prototype"`, so all static string indexers include the instance type
 */
export const staticLookup = Holder["some" + "thing"];
export const instanceLookup = (new Holder())["some" + "thing"];


//// [declarationEmitSimpleComputedNames1.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.instanceLookup = exports.staticLookup = exports.Holder = exports.conatainer = exports.fieldName = void 0;
exports.fieldName = Math.random() > 0.5 ? "f1" : "f2";
exports.conatainer = {
    [exports.fieldName]() {
        return "result";
    }
};
const classFieldName = Math.random() > 0.5 ? "g1" : "g2";
const otherField = classFieldName === "g1" ? "g2" : "g1";
const staticField = Math.random() > 0.5 ? "s1" : "s2";
class Holder {
    [classFieldName]() {
        return "value";
    }
    [otherField]() {
        return 42;
    }
    static [staticField]() {
        return { static: true };
    }
    static [staticField]() {
        return { static: "sometimes" };
    }
}
exports.Holder = Holder;
/**
 * Could be `"prototype"`, so all static string indexers include the instance type
 */
exports.staticLookup = Holder["some" + "thing"];
exports.instanceLookup = (new Holder())["some" + "thing"];


//// [declarationEmitSimpleComputedNames1.d.ts]
export const fieldName: string;
export const conatainer: {
    [x: string]: () => string;
};
const classFieldName: string;
const otherField: string;
const staticField: string;
export class Holder {
    static [staticField]: () => {
        static: boolean;
    };
    static [staticField]: () => {
        static: string;
    };
    [classFieldName]: () => string;
    [otherField]: () => number;
}
/**
 * Could be `"prototype"`, so all static string indexers include the instance type
 */
export const staticLookup: Holder | (() => {
    static: boolean;
}) | (() => {
    static: string;
});
export const instanceLookup: (() => string) | (() => number);
export {};


//// [DtsFileErrors]


declarationEmitSimpleComputedNames1.d.ts(5,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declarationEmitSimpleComputedNames1.d.ts (1 errors) ====
    export const fieldName: string;
    export const conatainer: {
        [x: string]: () => string;
    };
    const classFieldName: string;
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    const otherField: string;
    const staticField: string;
    export class Holder {
        static [staticField]: () => {
            static: boolean;
        };
        static [staticField]: () => {
            static: string;
        };
        [classFieldName]: () => string;
        [otherField]: () => number;
    }
    /**
     * Could be `"prototype"`, so all static string indexers include the instance type
     */
    export const staticLookup: Holder | (() => {
        static: boolean;
    }) | (() => {
        static: string;
    });
    export const instanceLookup: (() => string) | (() => number);
    export {};
    