//// [tests/cases/compiler/staticPropertyUsedBeforeInitInBaseClass.ts] ////

//// [staticPropertyUsedBeforeInitInBaseClass.ts]
class Base {
  p = Derived.S;
}
class Derived extends Base {
  static S = 1;
}


//// [staticPropertyUsedBeforeInitInBaseClass.js]
"use strict";
class Base {
    p = Derived.S;
}
class Derived extends Base {
    static S = 1;
}
