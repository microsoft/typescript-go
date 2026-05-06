//// [tests/cases/compiler/inferenceOptionalProperties.ts] ////

//// [inferenceOptionalProperties.ts]
declare function test<T>(x: { [key: string]: T }): T;

declare let x1: { a?: string, b?: number };
declare let x2: { a?: string, b?: number | undefined };

const y1 = test(x1);
const y2 = test(x2);

var v1: Required<{ a?: string, b?: number }>;
var v1: { a: string, b: number };

var v2: Required<{ a?: string, b?: number | undefined }>;
var v2: { a: string, b: number | undefined };

var v3: Partial<{ a: string, b: string }>;
var v3: { a?: string, b?: string };

var v4: Partial<{ a: string, b: string | undefined }>;
var v4: { a?: string, b?: string | undefined };

var v5: Required<Partial<{ a: string, b: string }>>;
var v5: { a: string, b: string };

var v6: Required<Partial<{ a: string, b: string | undefined }>>;
var v6: { a: string, b: string | undefined };


//// [inferenceOptionalProperties.js]
"use strict";
const y1 = test(x1);
const y2 = test(x2);
var v1;
var v1;
var v2;
var v2;
var v3;
var v3;
var v4;
var v4;
var v5;
var v5;
var v6;
var v6;


//// [inferenceOptionalProperties.d.ts]
function test<T>(x: {
    [key: string]: T;
}): T;
let x1: {
    a?: string;
    b?: number;
};
let x2: {
    a?: string;
    b?: number | undefined;
};
const y1: string | number;
const y2: string | number | undefined;
var v1: Required<{
    a?: string;
    b?: number;
}>;
var v1: {
    a: string;
    b: number;
};
var v2: Required<{
    a?: string;
    b?: number | undefined;
}>;
var v2: {
    a: string;
    b: number | undefined;
};
var v3: Partial<{
    a: string;
    b: string;
}>;
var v3: {
    a?: string;
    b?: string;
};
var v4: Partial<{
    a: string;
    b: string | undefined;
}>;
var v4: {
    a?: string;
    b?: string | undefined;
};
var v5: Required<Partial<{
    a: string;
    b: string;
}>>;
var v5: {
    a: string;
    b: string;
};
var v6: Required<Partial<{
    a: string;
    b: string | undefined;
}>>;
var v6: {
    a: string;
    b: string | undefined;
};


//// [DtsFileErrors]


inferenceOptionalProperties.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== inferenceOptionalProperties.d.ts (1 errors) ====
    function test<T>(x: {
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        [key: string]: T;
    }): T;
    let x1: {
        a?: string;
        b?: number;
    };
    let x2: {
        a?: string;
        b?: number | undefined;
    };
    const y1: string | number;
    const y2: string | number | undefined;
    var v1: Required<{
        a?: string;
        b?: number;
    }>;
    var v1: {
        a: string;
        b: number;
    };
    var v2: Required<{
        a?: string;
        b?: number | undefined;
    }>;
    var v2: {
        a: string;
        b: number | undefined;
    };
    var v3: Partial<{
        a: string;
        b: string;
    }>;
    var v3: {
        a?: string;
        b?: string;
    };
    var v4: Partial<{
        a: string;
        b: string | undefined;
    }>;
    var v4: {
        a?: string;
        b?: string | undefined;
    };
    var v5: Required<Partial<{
        a: string;
        b: string;
    }>>;
    var v5: {
        a: string;
        b: string;
    };
    var v6: Required<Partial<{
        a: string;
        b: string | undefined;
    }>>;
    var v6: {
        a: string;
        b: string | undefined;
    };
    