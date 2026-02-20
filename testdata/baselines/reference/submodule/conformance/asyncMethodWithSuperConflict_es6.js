//// [tests/cases/conformance/async/es2017/asyncMethodWithSuperConflict_es6.ts] ////

//// [asyncMethodWithSuperConflict_es6.ts]
class A {
    x() {
    }
    y() {
    }
}

class B extends A {
    // async method with only call/get on 'super' does not require a binding
    async simple() {
        const _super = null;
        const _superIndex = null;
        // call with property access
        super.x();
        // call additional property.
        super.y();

        // call with element access
        super["x"]();

        // property access (read)
        const a = super.x;

        // element access (read)
        const b = super["x"];
    }

    // async method with assignment/destructuring on 'super' requires a binding
    async advanced() {
        const _super = null;
        const _superIndex = null;
        const f = () => {};

        // call with property access
        super.x();

        // call with element access
        super["x"]();

        // property access (read)
        const a = super.x;

        // element access (read)
        const b = super["x"];

        // property access (assign)
        super.x = f;

        // element access (assign)
        super["x"] = f;

        // destructuring assign with property access
        ({ f: super.x } = { f });

        // destructuring assign with element access
        ({ f: super["x"] } = { f });
    }
}


//// [asyncMethodWithSuperConflict_es6.js]
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
class A {
    x() {
    }
    y() {
    }
}
class B extends A {
    // async method with only call/get on 'super' does not require a binding
    simple() {
        return __awaiter(this, void 0, void 0, function* () {
            const _super = null;
            const _superIndex = null;
            // call with property access
            super.x();
            // call additional property.
            super.y();
            // call with element access
            super["x"]();
            // property access (read)
            const a = super.x;
            // element access (read)
            const b = super["x"];
        });
    }
    // async method with assignment/destructuring on 'super' requires a binding
    advanced() {
        return __awaiter(this, void 0, void 0, function* () {
            const _super = null;
            const _superIndex = null;
            const f = () => { };
            // call with property access
            super.x();
            // call with element access
            super["x"]();
            // property access (read)
            const a = super.x;
            // element access (read)
            const b = super["x"];
            // property access (assign)
            super.x = f;
            // element access (assign)
            super["x"] = f;
            // destructuring assign with property access
            ({ f: super.x } = { f });
            // destructuring assign with element access
            ({ f: super["x"] } = { f });
        });
    }
}
