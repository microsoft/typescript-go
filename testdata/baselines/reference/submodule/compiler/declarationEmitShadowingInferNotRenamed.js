//// [tests/cases/compiler/declarationEmitShadowingInferNotRenamed.ts] ////

//// [declarationEmitShadowingInferNotRenamed.ts]
// Any instance type
type Client = string

// Modified instance
type UpdatedClient<C> = C & {foo: number}

export const createClient = <
  D extends
    | (new (...args: any[]) => Client) // accept class
    | Record<string, new (...args: any[]) => Client> // or map of classes
>(
  clientDef: D
): D extends new (...args: any[]) => infer C
  ? UpdatedClient<C> // return instance
  : {
      [K in keyof D]: D[K] extends new (...args: any[]) => infer C // or map of instances respectively
        ? UpdatedClient<C>
        : never
    } => {
  return null as any
}

//// [declarationEmitShadowingInferNotRenamed.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.createClient = void 0;
const createClient = (clientDef) => {
    return null;
};
exports.createClient = createClient;


//// [declarationEmitShadowingInferNotRenamed.d.ts]
// Modified instance
type UpdatedClient<C> = C & {
    foo: number;
};
export declare const createClient: <D extends Record<string, new (...args: any[]) => string> | (new (...args: any[]) => string)>(clientDef: D) => D extends new (...args: any[]) => infer C ? UpdatedClient<C> : { [K in keyof D]: D[K] extends new (...args: any[]) => infer C // or map of instances respectively
 ? UpdatedClient<C> : never; };
export {};
