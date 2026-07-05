currentDirectory::/home/src/workspaces/project
useCaseSensitiveFileNames::true
Input::
//// [/home/src/workspaces/project/mainJs.js] *new* 
/**
 * @template T
 * @param {T} component
 * @returns {T}
 */
function wrap(component) { return component; }

function FunctionComponent() { return null; }
FunctionComponent.propTypes = { num: 0 };

export const WrappedFunction = wrap(FunctionComponent);
//// [/home/src/workspaces/project/mainTs.ts] *new* 
declare function wrap<T>(component: T): T;

function FunctionComponent() { return null; }
FunctionComponent.propTypes = { num: 0 };

const ArrowComponent = () => null;
ArrowComponent.propTypes = { num: 0 };

function UnusedComponent() { return null; }
UnusedComponent.propTypes = { num: 0 };

export const WrappedFunction = wrap(FunctionComponent);
export const WrappedArrow = wrap(ArrowComponent);
//// [/home/src/workspaces/project/tsconfig.json] *new* 
{
    "compilerOptions": {
        "strict": true,
        "declaration": true,
        "emitDeclarationOnly": true,
        "allowJs": true,
        "checkJs": true,
        "outDir": "dist",
    },
}

tsgo -p .
ExitStatus:: Success
Output::
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
//// [/home/src/workspaces/project/dist/mainJs.d.ts] *new* 
declare function FunctionComponent(): null;
declare namespace FunctionComponent {
    var propTypes: {
        num: number;
    };
}
export declare const WrappedFunction: typeof FunctionComponent;
export {};

//// [/home/src/workspaces/project/dist/mainTs.d.ts] *new* 
declare function FunctionComponent(): null;
declare namespace FunctionComponent {
    var propTypes: {
        num: number;
    };
}
export declare const WrappedFunction: typeof FunctionComponent;
export declare const WrappedArrow: {
    (): null;
    propTypes: {
        num: number;
    };
};
export {};


