//// [tests/cases/compiler/declarationEmitInferredObjectLiteralWithThis.ts] ////

//// [declarationEmitInferredObjectLiteralWithThis.ts]
export class C {
    foo() {
        return {
            self: this,
        };
    }

    prop = {
        self: this,
    };
}


//// [declarationEmitInferredObjectLiteralWithThis.js]
export class C {
    foo() {
        return {
            self: this,
        };
    }
    prop = {
        self: this,
    };
}


//// [declarationEmitInferredObjectLiteralWithThis.d.ts]
export declare class C {
    foo(): any;
    prop: any;
}
