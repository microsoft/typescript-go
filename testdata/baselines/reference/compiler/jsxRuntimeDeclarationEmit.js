//// [tests/cases/compiler/jsxRuntimeDeclarationEmit.tsx] ////

//// [jsxRuntimeDeclarationEmit.tsx]
// Reproduces issue #1011: jsx-runtime declaration emit requires type annotation
// https://github.com/microsoft/typescript-go/issues/1011


/// <reference path="/.lib/react16.d.ts" />

// This should produce clean jsx-runtime imports without "unsafe import" errors
export const MyComponent = () => {
  return <div>Hello World</div>
}



//// [jsxRuntimeDeclarationEmit.d.ts]
// Reproduces issue #1011: jsx-runtime declaration emit requires type annotation
// https://github.com/microsoft/typescript-go/issues/1011
// This should produce clean jsx-runtime imports without "unsafe import" errors
export declare const MyComponent: () => JSX.Element;


//// [DtsFileErrors]


jsxRuntimeDeclarationEmit.d.ts(4,41): error TS2503: Cannot find namespace 'JSX'.


==== jsxRuntimeDeclarationEmit.d.ts (1 errors) ====
    // Reproduces issue #1011: jsx-runtime declaration emit requires type annotation
    // https://github.com/microsoft/typescript-go/issues/1011
    // This should produce clean jsx-runtime imports without "unsafe import" errors
    export declare const MyComponent: () => JSX.Element;
                                            ~~~
!!! error TS2503: Cannot find namespace 'JSX'.
    