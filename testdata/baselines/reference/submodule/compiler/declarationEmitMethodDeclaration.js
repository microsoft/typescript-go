//// [tests/cases/compiler/declarationEmitMethodDeclaration.ts] ////

//// [a.js]
export default {
    methods: {
        foo() { }
    }
}


//// [a.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.default = {
    methods: {
        foo() { }
    }
};
