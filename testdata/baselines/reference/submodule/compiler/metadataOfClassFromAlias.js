//// [tests/cases/compiler/metadataOfClassFromAlias.ts] ////

//// [auxiliry.ts]
export class SomeClass {
    field: string;
}

//// [test.ts]
import { SomeClass } from './auxiliry';
function annotation(): PropertyDecorator {
    return (target: any): void => { };
}
export class ClassA {
    @annotation() array: SomeClass | null;
}

//// [test.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ClassA = void 0;
function annotation() {
    return (target) => { };
}
class ClassA {
    @annotation()
    array;
}
exports.ClassA = ClassA;
//// [auxiliry.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.SomeClass = void 0;
class SomeClass {
    field;
}
exports.SomeClass = SomeClass;
