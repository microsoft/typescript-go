//// [tests/cases/conformance/classes/mixinAccessors5.ts] ////

//// [mixinAccessors5.ts]
// https://github.com/microsoft/TypeScript/issues/61967

declare function basicMixin<T extends object, U extends object>(
  t: T,
  u: U,
): T & U;
  
declare class GetterA {
  constructor(...args: any[]);

  get inCompendium(): boolean;
}
  
declare class GetterB {
  constructor(...args: any[]);

  get inCompendium(): boolean;
}
  
declare class TestB extends basicMixin(GetterA, GetterB) {
  override get inCompendium(): boolean;
}
  

//// [mixinAccessors5.js]
"use strict";
// https://github.com/microsoft/TypeScript/issues/61967


//// [mixinAccessors5.d.ts]
function basicMixin<T extends object, U extends object>(t: T, u: U): T & U;
class GetterA {
    constructor(...args: any[]);
    get inCompendium(): boolean;
}
class GetterB {
    constructor(...args: any[]);
    get inCompendium(): boolean;
}
const TestB_base: typeof GetterA & typeof GetterB;
class TestB extends TestB_base {
    get inCompendium(): boolean;
}


//// [DtsFileErrors]


mixinAccessors5.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== mixinAccessors5.d.ts (1 errors) ====
    function basicMixin<T extends object, U extends object>(t: T, u: U): T & U;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    class GetterA {
        constructor(...args: any[]);
        get inCompendium(): boolean;
    }
    class GetterB {
        constructor(...args: any[]);
        get inCompendium(): boolean;
    }
    const TestB_base: typeof GetterA & typeof GetterB;
    class TestB extends TestB_base {
        get inCompendium(): boolean;
    }
    