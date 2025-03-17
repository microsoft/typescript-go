//// [tests/cases/compiler/declarationEmitBundleWithAmbientReferences.ts] ////

//// [lib.d.ts]
declare module "lib/result" {
    export type Result<E extends Error, T> = (E & Failure<E>) | (T & Success<T>);
    export interface Failure<E extends Error> { }
    export interface Success<T> { }
}

//// [datastore_result.ts]
import { Result } from "lib/result";

export type T<T> = Result<Error, T>;

//// [conditional_directive_field.ts]
import * as DatastoreResult from "src/datastore_result";

export const build = (): DatastoreResult.T<string> => {
	return null;
};


//// [conditional_directive_field.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.build = void 0;
const build = () => {
    return null;
};
exports.build = build;
//// [datastore_result.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
