//// [tests/cases/conformance/types/stringLiteral/stringLiteralTypesAsTags02.ts] ////

//// [stringLiteralTypesAsTags02.ts]
type Kind = "A" | "B"

interface Entity {
    kind: Kind;
}

interface A extends Entity {
    kind: "A";
    a: number;
}

interface B extends Entity {
    kind: "B";
    b: string;
}

function hasKind(entity: Entity, kind: "A"): entity is A;
function hasKind(entity: Entity, kind: "B"): entity is B;
function hasKind(entity: Entity, kind: Kind): entity is (A | B) {
    return entity.kind === kind;
}

let x: A = {
    kind: "A",
    a: 100,
}

if (hasKind(x, "A")) {
    let a = x;
}
else {
    let b = x;
}

if (!hasKind(x, "B")) {
    let c = x;
}
else {
    let d = x;
}

//// [stringLiteralTypesAsTags02.js]
"use strict";
function hasKind(entity, kind) {
    return entity.kind === kind;
}
let x = {
    kind: "A",
    a: 100,
};
if (hasKind(x, "A")) {
    let a = x;
}
else {
    let b = x;
}
if (!hasKind(x, "B")) {
    let c = x;
}
else {
    let d = x;
}


//// [stringLiteralTypesAsTags02.d.ts]
type Kind = "A" | "B";
interface Entity {
    kind: Kind;
}
interface A extends Entity {
    kind: "A";
    a: number;
}
interface B extends Entity {
    kind: "B";
    b: string;
}
function hasKind(entity: Entity, kind: "A"): entity is A;
function hasKind(entity: Entity, kind: "B"): entity is B;
let x: A;


//// [DtsFileErrors]


stringLiteralTypesAsTags02.d.ts(13,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== stringLiteralTypesAsTags02.d.ts (1 errors) ====
    type Kind = "A" | "B";
    interface Entity {
        kind: Kind;
    }
    interface A extends Entity {
        kind: "A";
        a: number;
    }
    interface B extends Entity {
        kind: "B";
        b: string;
    }
    function hasKind(entity: Entity, kind: "A"): entity is A;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    function hasKind(entity: Entity, kind: "B"): entity is B;
    let x: A;
    