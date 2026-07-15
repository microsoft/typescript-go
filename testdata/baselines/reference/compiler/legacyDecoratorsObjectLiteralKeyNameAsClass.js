//// [tests/cases/compiler/legacyDecoratorsObjectLiteralKeyNameAsClass.ts] ////

//// [legacyDecoratorsObjectLiteralKeyNameAsClass.ts]
// Repro from issue: class name used as object literal key should not be renamed to the alias

const dec = (t: any) => t;

@dec
class SessionAuth {
    static requirement() {
        return { SessionAuth: [] };
    }
}

// Method shorthand name in object literal should not be renamed
@dec
class Foo {
    static methods() {
        return {
            Foo() { return 1; }
        };
    }
}

// Class field name should not be renamed
@dec
class Bar {
    Bar = 42;
}

// Self-reference in expression positions SHOULD still be renamed
@dec
class SelfRef {
    static instance = new SelfRef();
    method() {
        return SelfRef;
    }
}


//// [legacyDecoratorsObjectLiteralKeyNameAsClass.js]
"use strict";
// Repro from issue: class name used as object literal key should not be renamed to the alias
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var SelfRef_1;
const dec = (t) => t;
let SessionAuth = class SessionAuth {
    static requirement() {
        return { SessionAuth: [] };
    }
};
SessionAuth = __decorate([
    dec
], SessionAuth);
// Method shorthand name in object literal should not be renamed
let Foo = class Foo {
    static methods() {
        return {
            Foo() { return 1; }
        };
    }
};
Foo = __decorate([
    dec
], Foo);
// Class field name should not be renamed
let Bar = class Bar {
    Bar = 42;
};
Bar = __decorate([
    dec
], Bar);
// Self-reference in expression positions SHOULD still be renamed
let SelfRef = class SelfRef {
    static { SelfRef_1 = this; }
    static instance = new SelfRef_1();
    method() {
        return SelfRef_1;
    }
};
SelfRef = SelfRef_1 = __decorate([
    dec
], SelfRef);
