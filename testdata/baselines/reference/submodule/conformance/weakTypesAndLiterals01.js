//// [tests/cases/conformance/types/typeRelationships/comparable/weakTypesAndLiterals01.ts] ////

//// [weakTypesAndLiterals01.ts]
type WeakTypes =
    | { optional?: true; }
    | { toLowerCase?(): string }
    | { toUpperCase?(): string, otherOptionalProp?: number };

type LiteralsOrWeakTypes =
    | "A"
    | "B"
    | WeakTypes;

declare let aOrB: "A" | "B";

const f = (arg: LiteralsOrWeakTypes) => {
    if (arg === "A") {
        return arg;
    }
    else {
        return arg;
    }
}

const g = (arg: WeakTypes) => {
    if (arg === "A") {
        return arg;
    }
    else {
        return arg;
    }
}

const h = (arg: LiteralsOrWeakTypes) => {
    if (arg === aOrB) {
        return arg;
    }
    else {
        return arg;
    }
}

const i = (arg: WeakTypes) => {
    if (arg === aOrB) {
        return arg;
    }
    else {
        return arg;
    }
}




//// [weakTypesAndLiterals01.d.ts]
type WeakTypes = {
    optional?: true;
} | {
    toLowerCase?(): string;
} | {
    toUpperCase?(): string;
    otherOptionalProp?: number;
};
type LiteralsOrWeakTypes = "A" | "B" | WeakTypes;
let aOrB: "A" | "B";
const f: (arg: LiteralsOrWeakTypes) => "A" | "B" | WeakTypes;
const g: (arg: WeakTypes) => WeakTypes;
const h: (arg: LiteralsOrWeakTypes) => LiteralsOrWeakTypes;
const i: (arg: WeakTypes) => WeakTypes;


//// [DtsFileErrors]


weakTypesAndLiterals01.d.ts(10,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== weakTypesAndLiterals01.d.ts (1 errors) ====
    type WeakTypes = {
        optional?: true;
    } | {
        toLowerCase?(): string;
    } | {
        toUpperCase?(): string;
        otherOptionalProp?: number;
    };
    type LiteralsOrWeakTypes = "A" | "B" | WeakTypes;
    let aOrB: "A" | "B";
    ~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    const f: (arg: LiteralsOrWeakTypes) => "A" | "B" | WeakTypes;
    const g: (arg: WeakTypes) => WeakTypes;
    const h: (arg: LiteralsOrWeakTypes) => LiteralsOrWeakTypes;
    const i: (arg: WeakTypes) => WeakTypes;
    