//// [tests/cases/compiler/amdDeclarationEmitNoExtraDeclare.ts] ////

//// [Class.ts]
import { Configurable } from "./Configurable"

export class HiddenClass {}

export class ActualClass extends Configurable(HiddenClass) {}
//// [Configurable.ts]
export type Constructor<T> = {
    new(...args: any[]): T;
}
export function Configurable<T extends Constructor<{}>>(base: T): T {
    return class extends base {

        constructor(...args: any[]) {
            super(...args);
        }

    };
}


//// [Class.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ActualClass = exports.HiddenClass = void 0;
const Configurable_1 = require("./Configurable");
class HiddenClass {
}
exports.HiddenClass = HiddenClass;
class ActualClass extends (0, Configurable_1.Configurable)(HiddenClass) {
}
exports.ActualClass = ActualClass;
//// [Configurable.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Configurable = Configurable;
function Configurable(base) {
    return class extends base {
        constructor(...args) {
            super(...args);
        }
    };
}
