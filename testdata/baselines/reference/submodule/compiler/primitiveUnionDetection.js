//// [tests/cases/compiler/primitiveUnionDetection.ts] ////

//// [primitiveUnionDetection.ts]
// Repro from #46624

type Kind = "one" | "two" | "three";

declare function getInterfaceFromString<T extends Kind>(options?: { type?: T } & { type?: Kind }): T;

const result = getInterfaceFromString({ type: 'two' });


//// [primitiveUnionDetection.js]
"use strict";
// Repro from #46624
const result = getInterfaceFromString({ type: 'two' });


//// [primitiveUnionDetection.d.ts]
type Kind = "one" | "two" | "three";
function getInterfaceFromString<T extends Kind>(options?: {
    type?: T;
} & {
    type?: Kind;
}): T;
const result: "two";


//// [DtsFileErrors]


primitiveUnionDetection.d.ts(2,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== primitiveUnionDetection.d.ts (1 errors) ====
    type Kind = "one" | "two" | "three";
    function getInterfaceFromString<T extends Kind>(options?: {
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        type?: T;
    } & {
        type?: Kind;
    }): T;
    const result: "two";
    