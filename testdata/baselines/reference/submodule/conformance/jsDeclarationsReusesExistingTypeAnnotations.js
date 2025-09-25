//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsReusesExistingTypeAnnotations.ts] ////

//// [index.js]
class С1 {
    /** @type {string=} */
    p1 = undefined;

    /** @type {string | undefined} */
    p2 = undefined;

    /** @type {?string} */
    p3 = null;

    /** @type {string | null} */
    p4 = null;
}

class С2 {
    /** @type {string=} */
    get p1() {
        return undefined;
    }

    /** @type {string | undefined} */
    get p2() {
        return undefined;
    }

    /** @type {?string} */
    get p3() {
        return null;
    }

    /** @type {string | null} */
    get p4() {
        return null;
    }
}


class С3 {
    /** @type {string=} */
    get p1() {
        return undefined;
    }

    /** @param {string=} value */
    set p1(value) {
        this.p1 = value;
    }

    /** @type {string | undefined} */
    get p2() {
        return undefined;
    }

    /** @param {string | undefined} value */
    set p2(value) {
        this.p2 = value;
    }

    /** @type {?string} */
    get p3() {
        return null;
    }

    /** @param {?string} value */
    set p3(value) {
        this.p3 = value;
    }

    /** @type {string | null} */
    get p4() {
        return null;
    }

    /** @param {string | null} value */
    set p4(value) {
        this.p4 = value;
    }
}


class С4 {
    /** @param {string=} value */
    set p1(value) {
        this.p1 = value;
    }

    /** @param {string | undefined} value */
    set p2(value) {
        this.p2 = value;
    }

    /** @param {?string} value */
    set p3(value) {
        this.p3 = value;
    }

    /** @param {string | null} value */
    set p4(value) {
        this.p4 = value;
    }
}


//// [index.js]
class С1 {
    /** @type {string=} */
    p1 = undefined;
    /** @type {string | undefined} */
    p2 = undefined;
    /** @type {?string} */
    p3 = null;
    /** @type {string | null} */
    p4 = null;
}
class С2 {
    /** @type {string=} */
    get p1() {
        return undefined;
    }
    /** @type {string | undefined} */
    get p2() {
        return undefined;
    }
    /** @type {?string} */
    get p3() {
        return null;
    }
    /** @type {string | null} */
    get p4() {
        return null;
    }
}
class С3 {
    /** @type {string=} */
    get p1() {
        return undefined;
    }
    /** @param {string=} value */
    set p1(value) {
        this.p1 = value;
    }
    /** @type {string | undefined} */
    get p2() {
        return undefined;
    }
    /** @param {string | undefined} value */
    set p2(value) {
        this.p2 = value;
    }
    /** @type {?string} */
    get p3() {
        return null;
    }
    /** @param {?string} value */
    set p3(value) {
        this.p3 = value;
    }
    /** @type {string | null} */
    get p4() {
        return null;
    }
    /** @param {string | null} value */
    set p4(value) {
        this.p4 = value;
    }
}
class С4 {
    /** @param {string=} value */
    set p1(value) {
        this.p1 = value;
    }
    /** @param {string | undefined} value */
    set p2(value) {
        this.p2 = value;
    }
    /** @param {?string} value */
    set p3(value) {
        this.p3 = value;
    }
    /** @param {string | null} value */
    set p4(value) {
        this.p4 = value;
    }
}


//// [index.d.ts]
declare class С1 {
    p1: string | undefined;
    p2: string | undefined;
    p3: string | null;
    p4: string | null;
}
declare class С2 {
    get p1(): undefined;
    get p2(): undefined;
    get p3(): null;
    get p4(): null;
}
declare class С3 {
    get p1(): string | undefined;
    set p1(value?: string | undefined);
    get p2(): string | undefined;
    set p2(value: string | undefined);
    get p3(): string | null;
    set p3(value: string | null);
    get p4(): string | null;
    set p4(value: string | null);
}
declare class С4 {
    set p1(value?: string | undefined);
    set p2(value: string | undefined);
    set p3(value: string | null);
    set p4(value: string | null);
}
