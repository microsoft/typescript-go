//// [tests/cases/compiler/declarationEmitCastReusesTypeNode3.ts] ////

//// [declarationEmitCastReusesTypeNode3.ts]
type P = { } & { name: string }

export let vLet = <P>null!
export const vConst = <P>null!

export function fn(p = <P>null!) {}

export function fnWithRequiredDefaultParam(p = <P>null!, req: number) {}

export class C {
    field = <P>null!
    optField? = <P>null!
    readonly roFiled = <P>null!;
    method(p = <P>null!) {}
    methodWithRequiredDefault(p = <P>null!, req: number) {}

    constructor(public ctorField = <P>null!) {}

    get x() { return <P>null! }
    set x(v) { }
}

export default <P>null!;

// allows `undefined` on the input side, thanks to the initializer
export function fnWithPartialAnnotationOnDefaultparam(x: P = <P>null!, b: number) {}



//// [declarationEmitCastReusesTypeNode3.d.ts]
type P = {} & {
    name: string;
};
export let vLet: P;
export const vConst: P;
export function fn(p?: P): void;
export function fnWithRequiredDefaultParam(p: P, req: number): void;
export class C {
    ctorField: P;
    field: P;
    optField?: P;
    readonly roFiled: P;
    method(p?: P): void;
    methodWithRequiredDefault(p: P, req: number): void;
    constructor(ctorField?: P);
    get x(): P;
    set x(v: P);
}
const _default: P;
export default _default;
export function fnWithPartialAnnotationOnDefaultparam(x: P, b: number): void;


//// [DtsFileErrors]


declarationEmitCastReusesTypeNode3.d.ts(19,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declarationEmitCastReusesTypeNode3.d.ts (1 errors) ====
    type P = {} & {
        name: string;
    };
    export let vLet: P;
    export const vConst: P;
    export function fn(p?: P): void;
    export function fnWithRequiredDefaultParam(p: P, req: number): void;
    export class C {
        ctorField: P;
        field: P;
        optField?: P;
        readonly roFiled: P;
        method(p?: P): void;
        methodWithRequiredDefault(p: P, req: number): void;
        constructor(ctorField?: P);
        get x(): P;
        set x(v: P);
    }
    const _default: P;
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    export default _default;
    export function fnWithPartialAnnotationOnDefaultparam(x: P, b: number): void;
    