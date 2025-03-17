//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsClassAccessor.ts] ////

//// [supplement.d.ts]
export { };
declare module "./argument.js" {
    interface Argument {
        idlType: any;
        default: null;
    }
}
//// [base.js]
export class Base {
    constructor() { }

    toJSON() {
        const json = { type: undefined, name: undefined, inheritance: undefined };
        return json;
    }
}
//// [argument.js]
import { Base } from "./base.js";
export class Argument extends Base {
    /**
     * @param {*} tokeniser
     */
    static parse(tokeniser) {
        return;
    }

    get type() {
        return "argument";
    }

    /**
     * @param {*} defs
     */
    *validate(defs) { }
}

//// [argument.js]
import { Base } from "./base.js";
export class Argument extends Base {
    static parse(tokeniser) {
        return;
    }
    get type() {
        return "argument";
    }
    *validate(defs) { }
}
//// [base.js]
export class Base {
    constructor() { }
    toJSON() {
        const json = { type: undefined, name: undefined, inheritance: undefined };
        return json;
    }
}
