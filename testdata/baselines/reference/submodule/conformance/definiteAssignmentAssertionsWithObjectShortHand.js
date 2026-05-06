//// [tests/cases/conformance/controlFlow/definiteAssignmentAssertionsWithObjectShortHand.ts] ////

//// [definiteAssignmentAssertionsWithObjectShortHand.ts]
const a: string | undefined = 'ff';
const foo = { a! }

const bar = {
    a ? () { }
}

//// [definiteAssignmentAssertionsWithObjectShortHand.js]
"use strict";
const a = 'ff';
const foo = { a };
const bar = {
    a() { }
};


//// [definiteAssignmentAssertionsWithObjectShortHand.d.ts]
const a: string | undefined;
const foo: {
    a: string;
};
const bar: {
    a(): void;
};
