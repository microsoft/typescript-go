//// [tests/cases/conformance/classes/members/instanceAndStaticMembers/typeOfThisInStaticMembers11.ts] ////

//// [typeOfThisInStaticMembers11.ts]
declare const foo: any;

@foo
class C {
    static a = 1;
    static b = this.a + 1;
}

@foo
class D extends C {
    static c = 2;
    static d = this.c + 1;
    static e = super.a + this.c + 1;
    static f = () => this.c + 1;
    static ff = function () { this.c + 1 }
    static foo () {
        return this.c + 1;
    }
    static get fa () {
        return this.c + 1;
    }
    static set fa (v: number) {
        this.c = v + 1;
    }
}

class CC {
    static a = 1;
    static b = this.a + 1;
}

class DD extends CC {
    static c = 2;
    static d = this.c + 1;
    static e = super.a + this.c + 1;
    static f = () => this.c + 1;
    static ff = function () { this.c + 1 }
    static foo () {
        return this.c + 1;
    }
    static get fa () {
        return this.c + 1;
    }
    static set fa (v: number) {
        this.c = v + 1;
    }
}


//// [typeOfThisInStaticMembers11.js]
"use strict";
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var _a, _b, _c, _d, _e, _f;
let C = (_a = class C {
    },
    Object.defineProperty(_a, "a", {
        enumerable: true,
        configurable: true,
        writable: true,
        value: 1
    }),
    Object.defineProperty(_a, "b", {
        enumerable: true,
        configurable: true,
        writable: true,
        value: _a.a + 1
    }),
    _a);
C = __decorate([
    foo
], C);
let D = (_b = class D extends (_c = C) {
        static foo() {
            return this.c + 1;
        }
        static get fa() {
            return this.c + 1;
        }
        static set fa(v) {
            this.c = v + 1;
        }
    },
    Object.defineProperty(_b, "c", {
        enumerable: true,
        configurable: true,
        writable: true,
        value: 2
    }),
    Object.defineProperty(_b, "d", {
        enumerable: true,
        configurable: true,
        writable: true,
        value: _b.c + 1
    }),
    Object.defineProperty(_b, "e", {
        enumerable: true,
        configurable: true,
        writable: true,
        value: Reflect.get(_c, "a", _b) + _b.c + 1
    }),
    Object.defineProperty(_b, "f", {
        enumerable: true,
        configurable: true,
        writable: true,
        value: () => _b.c + 1
    }),
    Object.defineProperty(_b, "ff", {
        enumerable: true,
        configurable: true,
        writable: true,
        value: function () { this.c + 1; }
    }),
    _b);
D = __decorate([
    foo
], D);
class CC {
}
_d = CC;
Object.defineProperty(CC, "a", {
    enumerable: true,
    configurable: true,
    writable: true,
    value: 1
});
Object.defineProperty(CC, "b", {
    enumerable: true,
    configurable: true,
    writable: true,
    value: _d.a + 1
});
class DD extends (_f = CC) {
    static foo() {
        return this.c + 1;
    }
    static get fa() {
        return this.c + 1;
    }
    static set fa(v) {
        this.c = v + 1;
    }
}
_e = DD;
Object.defineProperty(DD, "c", {
    enumerable: true,
    configurable: true,
    writable: true,
    value: 2
});
Object.defineProperty(DD, "d", {
    enumerable: true,
    configurable: true,
    writable: true,
    value: _e.c + 1
});
Object.defineProperty(DD, "e", {
    enumerable: true,
    configurable: true,
    writable: true,
    value: Reflect.get(_f, "a", _e) + _e.c + 1
});
Object.defineProperty(DD, "f", {
    enumerable: true,
    configurable: true,
    writable: true,
    value: () => _e.c + 1
});
Object.defineProperty(DD, "ff", {
    enumerable: true,
    configurable: true,
    writable: true,
    value: function () { this.c + 1; }
});
