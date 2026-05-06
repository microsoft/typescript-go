//// [tests/cases/compiler/declarationNoDanglingGenerics.ts] ////

//// [declarationNoDanglingGenerics.ts]
const kindCache: { [kind: string]: boolean } = {};

function register(kind: string): void | never {
  if (kindCache[kind]) {
    throw new Error(`Class with kind "${kind}" is already registered.`);
  }
  kindCache[kind] = true;
}

function ClassFactory<TKind extends string>(kind: TKind) {
  register(kind);

  return class {
    static readonly THE_KIND: TKind = kind;
    readonly kind: TKind = kind;
  };
}

class Kinds {
  static readonly A = "A";
  static readonly B = "B";
  static readonly C = "C";
}

export class AKind extends ClassFactory(Kinds.A) {
}

export class BKind extends ClassFactory(Kinds.B) {
}

export class CKind extends ClassFactory(Kinds.C) {
}

//// [declarationNoDanglingGenerics.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.CKind = exports.BKind = exports.AKind = void 0;
const kindCache = {};
function register(kind) {
    if (kindCache[kind]) {
        throw new Error(`Class with kind "${kind}" is already registered.`);
    }
    kindCache[kind] = true;
}
function ClassFactory(kind) {
    var _a;
    register(kind);
    return _a = class {
            constructor() {
                this.kind = kind;
            }
        },
        _a.THE_KIND = kind,
        _a;
}
class Kinds {
}
Kinds.A = "A";
Kinds.B = "B";
Kinds.C = "C";
class AKind extends ClassFactory(Kinds.A) {
}
exports.AKind = AKind;
class BKind extends ClassFactory(Kinds.B) {
}
exports.BKind = BKind;
class CKind extends ClassFactory(Kinds.C) {
}
exports.CKind = CKind;


//// [declarationNoDanglingGenerics.d.ts]
const AKind_base: {
    new (): {
        readonly kind: "A";
    };
    readonly THE_KIND: "A";
};
export class AKind extends AKind_base {
}
const BKind_base: {
    new (): {
        readonly kind: "B";
    };
    readonly THE_KIND: "B";
};
export class BKind extends BKind_base {
}
const CKind_base: {
    new (): {
        readonly kind: "C";
    };
    readonly THE_KIND: "C";
};
export class CKind extends CKind_base {
}
export {};


//// [DtsFileErrors]


declarationNoDanglingGenerics.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declarationNoDanglingGenerics.d.ts (1 errors) ====
    const AKind_base: {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        new (): {
            readonly kind: "A";
        };
        readonly THE_KIND: "A";
    };
    export class AKind extends AKind_base {
    }
    const BKind_base: {
        new (): {
            readonly kind: "B";
        };
        readonly THE_KIND: "B";
    };
    export class BKind extends BKind_base {
    }
    const CKind_base: {
        new (): {
            readonly kind: "C";
        };
        readonly THE_KIND: "C";
    };
    export class CKind extends CKind_base {
    }
    export {};
    