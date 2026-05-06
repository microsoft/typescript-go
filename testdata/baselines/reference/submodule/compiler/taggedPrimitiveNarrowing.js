//// [tests/cases/compiler/taggedPrimitiveNarrowing.ts] ////

//// [taggedPrimitiveNarrowing.ts]
type Hash = string & { __hash: true };

function getHashLength(hash: Hash): number {
    if (typeof hash !== "string") {
        throw new Error("This doesn't look like a hash");
    }
    return hash.length;
}

function getHashLength2<T extends { __tag__: unknown}>(hash: string & T): number {
    if (typeof hash !== "string") {
        throw new Error("This doesn't look like a hash");
    }
    return hash.length;
}


//// [taggedPrimitiveNarrowing.js]
"use strict";
function getHashLength(hash) {
    if (typeof hash !== "string") {
        throw new Error("This doesn't look like a hash");
    }
    return hash.length;
}
function getHashLength2(hash) {
    if (typeof hash !== "string") {
        throw new Error("This doesn't look like a hash");
    }
    return hash.length;
}


//// [taggedPrimitiveNarrowing.d.ts]
type Hash = string & {
    __hash: true;
};
function getHashLength(hash: Hash): number;
function getHashLength2<T extends {
    __tag__: unknown;
}>(hash: string & T): number;


//// [DtsFileErrors]


taggedPrimitiveNarrowing.d.ts(4,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== taggedPrimitiveNarrowing.d.ts (1 errors) ====
    type Hash = string & {
        __hash: true;
    };
    function getHashLength(hash: Hash): number;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    function getHashLength2<T extends {
        __tag__: unknown;
    }>(hash: string & T): number;
    