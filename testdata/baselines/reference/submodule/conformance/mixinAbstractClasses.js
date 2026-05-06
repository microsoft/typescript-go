//// [tests/cases/conformance/classes/mixinAbstractClasses.ts] ////

//// [mixinAbstractClasses.ts]
interface Mixin {
    mixinMethod(): void;
}

function Mixin<TBaseClass extends abstract new (...args: any) => any>(baseClass: TBaseClass): TBaseClass & (abstract new (...args: any) => Mixin) {
    abstract class MixinClass extends baseClass implements Mixin {
        mixinMethod() {
        }
    }
    return MixinClass;
}

class ConcreteBase {
    baseMethod() {}
}

abstract class AbstractBase {
    abstract abstractBaseMethod(): void;
}

class DerivedFromConcrete extends Mixin(ConcreteBase) {
}

const wasConcrete = new DerivedFromConcrete();
wasConcrete.baseMethod();
wasConcrete.mixinMethod();

class DerivedFromAbstract extends Mixin(AbstractBase) {
    abstractBaseMethod() {}
}

const wasAbstract = new DerivedFromAbstract();
wasAbstract.abstractBaseMethod();
wasAbstract.mixinMethod();

//// [mixinAbstractClasses.js]
"use strict";
function Mixin(baseClass) {
    class MixinClass extends baseClass {
        mixinMethod() {
        }
    }
    return MixinClass;
}
class ConcreteBase {
    baseMethod() { }
}
class AbstractBase {
}
class DerivedFromConcrete extends Mixin(ConcreteBase) {
}
const wasConcrete = new DerivedFromConcrete();
wasConcrete.baseMethod();
wasConcrete.mixinMethod();
class DerivedFromAbstract extends Mixin(AbstractBase) {
    abstractBaseMethod() { }
}
const wasAbstract = new DerivedFromAbstract();
wasAbstract.abstractBaseMethod();
wasAbstract.mixinMethod();


//// [mixinAbstractClasses.d.ts]
interface Mixin {
    mixinMethod(): void;
}
function Mixin<TBaseClass extends abstract new (...args: any) => any>(baseClass: TBaseClass): TBaseClass & (abstract new (...args: any) => Mixin);
class ConcreteBase {
    baseMethod(): void;
}
abstract class AbstractBase {
    abstract abstractBaseMethod(): void;
}
const DerivedFromConcrete_base: typeof ConcreteBase & (abstract new (...args: any) => Mixin);
class DerivedFromConcrete extends DerivedFromConcrete_base {
}
const wasConcrete: DerivedFromConcrete;
const DerivedFromAbstract_base: typeof AbstractBase & (abstract new (...args: any) => Mixin);
class DerivedFromAbstract extends DerivedFromAbstract_base {
    abstractBaseMethod(): void;
}
const wasAbstract: DerivedFromAbstract;


//// [DtsFileErrors]


mixinAbstractClasses.d.ts(4,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== mixinAbstractClasses.d.ts (1 errors) ====
    interface Mixin {
        mixinMethod(): void;
    }
    function Mixin<TBaseClass extends abstract new (...args: any) => any>(baseClass: TBaseClass): TBaseClass & (abstract new (...args: any) => Mixin);
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    class ConcreteBase {
        baseMethod(): void;
    }
    abstract class AbstractBase {
        abstract abstractBaseMethod(): void;
    }
    const DerivedFromConcrete_base: typeof ConcreteBase & (abstract new (...args: any) => Mixin);
    class DerivedFromConcrete extends DerivedFromConcrete_base {
    }
    const wasConcrete: DerivedFromConcrete;
    const DerivedFromAbstract_base: typeof AbstractBase & (abstract new (...args: any) => Mixin);
    class DerivedFromAbstract extends DerivedFromAbstract_base {
        abstractBaseMethod(): void;
    }
    const wasAbstract: DerivedFromAbstract;
    