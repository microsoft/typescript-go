currentDirectory::/home/src/workspaces/project
useCaseSensitiveFileNames::true
Input::
//// [/home/src/workspaces/project/index.tsx] *new* 
declare namespace JSX {
    interface ElementChildrenAttribute { children: {}; }
    interface IntrinsicElements { div: {} }
}

declare var React: any;

declare function Component(props: never): any;
declare function Component(props: { children?: number }): any;
(<Component>
    <div />
    <div />
</Component>)

tsgo --outputFormat minimal --noEmit --strict --jsx react index.tsx
ExitStatus:: DiagnosticsPresent_OutputsSkipped
Output::
index.tsx:10:3: TS2769: No overload matches this call. | The last overload gave the following error. | Type '{ children: any[]; }' is not assignable to type '{ children?: number | undefined; }'. | Types of property 'children' are incompatible. | Type 'any[]' is not assignable to type 'number'.
//// [/home/src/tslibs/TS/Lib/lib.es2025.full.d.ts] *Lib*
/// <reference no-default-lib="true"/>
interface Boolean {}
interface Function {}
interface CallableFunction {}
interface NewableFunction {}
interface IArguments {}
interface Number { toExponential: any; }
interface Object {}
interface RegExp {}
interface String { charAt: any; }
interface Array<T> { length: number; [n: number]: T; }
interface ReadonlyArray<T> {}
interface SymbolConstructor {
    (desc?: string | number): symbol;
    for(name: string): symbol;
    readonly toStringTag: symbol;
}
declare var Symbol: SymbolConstructor;
interface Symbol {
    readonly [Symbol.toStringTag]: string;
}
declare const console: { log(msg: any): void; };

