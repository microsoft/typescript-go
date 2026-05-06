//// [tests/cases/compiler/declarationEmitLocalClassDeclarationMixin.ts] ////

//// [declarationEmitLocalClassDeclarationMixin.ts]
interface Constructor<C> { new (...args: any[]): C; }

function mixin<B extends Constructor<{}>>(Base: B) {
    class PrivateMixed extends Base {
        bar = 2;
    }
    return PrivateMixed;
}

export class Unmixed {
    foo = 1;
}

export const Mixed = mixin(Unmixed);

function Filter<C extends Constructor<{}>>(ctor: C) {
    abstract class FilterMixin extends ctor {
        abstract match(path: string): boolean;
        // other concrete methods, fields, constructor
        thing = 12;
    }
    return FilterMixin;
}

export class FilteredThing extends Filter(Unmixed) {
    match(path: string) {
        return false;
    }
}


//// [declarationEmitLocalClassDeclarationMixin.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.FilteredThing = exports.Mixed = exports.Unmixed = void 0;
function mixin(Base) {
    class PrivateMixed extends Base {
        constructor() {
            super(...arguments);
            this.bar = 2;
        }
    }
    return PrivateMixed;
}
class Unmixed {
    constructor() {
        this.foo = 1;
    }
}
exports.Unmixed = Unmixed;
exports.Mixed = mixin(Unmixed);
function Filter(ctor) {
    class FilterMixin extends ctor {
        constructor() {
            super(...arguments);
            // other concrete methods, fields, constructor
            this.thing = 12;
        }
    }
    return FilterMixin;
}
class FilteredThing extends Filter(Unmixed) {
    match(path) {
        return false;
    }
}
exports.FilteredThing = FilteredThing;


//// [declarationEmitLocalClassDeclarationMixin.d.ts]
export class Unmixed {
    foo: number;
}
export const Mixed: {
    new (...args: any[]): {
        bar: number;
    };
} & typeof Unmixed;
const FilteredThing_base: (abstract new (...args: any[]) => {
    match(path: string): boolean;
    thing: number;
}) & typeof Unmixed;
export class FilteredThing extends FilteredThing_base {
    match(path: string): boolean;
}
export {};


//// [DtsFileErrors]


declarationEmitLocalClassDeclarationMixin.d.ts(9,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declarationEmitLocalClassDeclarationMixin.d.ts (1 errors) ====
    export class Unmixed {
        foo: number;
    }
    export const Mixed: {
        new (...args: any[]): {
            bar: number;
        };
    } & typeof Unmixed;
    const FilteredThing_base: (abstract new (...args: any[]) => {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        match(path: string): boolean;
        thing: number;
    }) & typeof Unmixed;
    export class FilteredThing extends FilteredThing_base {
        match(path: string): boolean;
    }
    export {};
    