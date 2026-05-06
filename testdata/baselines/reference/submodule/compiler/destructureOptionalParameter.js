//// [tests/cases/compiler/destructureOptionalParameter.ts] ////

//// [destructureOptionalParameter.ts]
declare function f1({ a, b }?: { a: number, b: string }): void;

function f2({ a, b }: { a: number, b: number } = { a: 0, b: 0 }) {
    a;
    b;
}

// Repro from #8681

interface Type { t: void }
interface QueryMetadata { q: void }

interface QueryMetadataFactory {
    (selector: Type | string, {descendants, read}?: {
        descendants?: boolean;
        read?: any;
    }): ParameterDecorator;
    new (selector: Type | string, {descendants, read}?: {
        descendants?: boolean;
        read?: any;
    }): QueryMetadata;
}


//// [destructureOptionalParameter.js]
"use strict";
function f2({ a, b } = { a: 0, b: 0 }) {
    a;
    b;
}


//// [destructureOptionalParameter.d.ts]
function f1({ a, b }?: {
    a: number;
    b: string;
}): void;
function f2({ a, b }?: {
    a: number;
    b: number;
}): void;
interface Type {
    t: void;
}
interface QueryMetadata {
    q: void;
}
interface QueryMetadataFactory {
    (selector: Type | string, { descendants, read }?: {
        descendants?: boolean;
        read?: any;
    }): ParameterDecorator;
    new (selector: Type | string, { descendants, read }?: {
        descendants?: boolean;
        read?: any;
    }): QueryMetadata;
}


//// [DtsFileErrors]


destructureOptionalParameter.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== destructureOptionalParameter.d.ts (1 errors) ====
    function f1({ a, b }?: {
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        a: number;
        b: string;
    }): void;
    function f2({ a, b }?: {
        a: number;
        b: number;
    }): void;
    interface Type {
        t: void;
    }
    interface QueryMetadata {
        q: void;
    }
    interface QueryMetadataFactory {
        (selector: Type | string, { descendants, read }?: {
            descendants?: boolean;
            read?: any;
        }): ParameterDecorator;
        new (selector: Type | string, { descendants, read }?: {
            descendants?: boolean;
            read?: any;
        }): QueryMetadata;
    }
    