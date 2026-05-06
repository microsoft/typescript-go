//// [tests/cases/compiler/lateBoundAssignmentCandidateJS2.ts] ////

//// [index.js]
const prop = 'prop';

export class foo1 {
    constructor() {
        this[prop] = 'bar'
    }

    /**
     * @protected
     * @type {string}
     */
    [prop] = 'baz';
}


//// [index.js]
const prop = 'prop';
export class foo1 {
    constructor() {
        this[prop] = 'bar';
    }
    /**
     * @protected
     * @type {string}
     */
    [prop] = 'baz';
}


//// [index.d.ts]
const prop = "prop";
export class foo1 {
    constructor();
    /**
     * @protected
     * @type {string}
     */
    protected [prop]: string;
}
export {};


//// [DtsFileErrors]


dist/index.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== dist/index.d.ts (1 errors) ====
    const prop = "prop";
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    export class foo1 {
        constructor();
        /**
         * @protected
         * @type {string}
         */
        protected [prop]: string;
    }
    export {};
    