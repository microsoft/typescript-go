//// [tests/cases/conformance/classes/mixinAbstractClassesReturnTypeInference.ts] ////

//// [mixinAbstractClassesReturnTypeInference.ts]
interface Mixin1 {
    mixinMethod(): void;
}

abstract class AbstractBase {
    abstract abstractBaseMethod(): void;
}

function Mixin2<TBase extends abstract new (...args: any[]) => any>(baseClass: TBase) {
    // must be `abstract` because we cannot know *all* of the possible abstract members that need to be
    // implemented for this to be concrete.
    abstract class MixinClass extends baseClass implements Mixin1 {
        mixinMethod(): void {}
        static staticMixinMethod(): void {}
    }
    return MixinClass;
}

class DerivedFromAbstract2 extends Mixin2(AbstractBase) {
    abstractBaseMethod() {}
}


//// [mixinAbstractClassesReturnTypeInference.js]
"use strict";
class AbstractBase {
}
function Mixin2(baseClass) {
    // must be `abstract` because we cannot know *all* of the possible abstract members that need to be
    // implemented for this to be concrete.
    class MixinClass extends baseClass {
        mixinMethod() { }
        static staticMixinMethod() { }
    }
    return MixinClass;
}
class DerivedFromAbstract2 extends Mixin2(AbstractBase) {
    abstractBaseMethod() { }
}


//// [mixinAbstractClassesReturnTypeInference.d.ts]
interface Mixin1 {
    mixinMethod(): void;
}
abstract class AbstractBase {
    abstract abstractBaseMethod(): void;
}
function Mixin2<TBase extends abstract new (...args: any[]) => any>(baseClass: TBase): ((abstract new (...args: any[]) => {
    [x: string]: any;
    mixinMethod(): void;
}) & {
    staticMixinMethod(): void;
}) & TBase;
const DerivedFromAbstract2_base: ((abstract new (...args: any[]) => {
    [x: string]: any;
    mixinMethod(): void;
}) & {
    staticMixinMethod(): void;
}) & typeof AbstractBase;
class DerivedFromAbstract2 extends DerivedFromAbstract2_base {
    abstractBaseMethod(): void;
}


//// [DtsFileErrors]


mixinAbstractClassesReturnTypeInference.d.ts(4,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== mixinAbstractClassesReturnTypeInference.d.ts (1 errors) ====
    interface Mixin1 {
        mixinMethod(): void;
    }
    abstract class AbstractBase {
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        abstract abstractBaseMethod(): void;
    }
    function Mixin2<TBase extends abstract new (...args: any[]) => any>(baseClass: TBase): ((abstract new (...args: any[]) => {
        [x: string]: any;
        mixinMethod(): void;
    }) & {
        staticMixinMethod(): void;
    }) & TBase;
    const DerivedFromAbstract2_base: ((abstract new (...args: any[]) => {
        [x: string]: any;
        mixinMethod(): void;
    }) & {
        staticMixinMethod(): void;
    }) & typeof AbstractBase;
    class DerivedFromAbstract2 extends DerivedFromAbstract2_base {
        abstractBaseMethod(): void;
    }
    