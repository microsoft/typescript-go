//// [tests/cases/conformance/types/typeRelationships/comparable/optionalProperties01.ts] ////

//// [optionalProperties01.ts]
interface Foo {
  required1: string;
  required2: string;
  optional?: string;
}

const foo1 = { required1: "hello" } as Foo;
const foo2 = { required1: "hello", optional: "bar" } as Foo;


//// [optionalProperties01.js]
"use strict";
const foo1 = { required1: "hello" };
const foo2 = { required1: "hello", optional: "bar" };


//// [optionalProperties01.d.ts]
interface Foo {
    required1: string;
    required2: string;
    optional?: string;
}
const foo1: Foo;
const foo2: Foo;


//// [DtsFileErrors]


optionalProperties01.d.ts(6,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== optionalProperties01.d.ts (1 errors) ====
    interface Foo {
        required1: string;
        required2: string;
        optional?: string;
    }
    const foo1: Foo;
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    const foo2: Foo;
    