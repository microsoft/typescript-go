//// [tests/cases/compiler/jsxRuntimeDeclarationEmit.tsx] ////

//// [jsxRuntimeDeclarationEmit.tsx]
// Test case to reproduce jsx-runtime declaration emit issue


/// <reference path="/.lib/react16.d.ts" />

// This should trigger the jsx-runtime import without type annotation error

export const FunctionComponent = () => {
  return <div>Hello World</div>
}

export const AnotherComponent = () => {
  return <FunctionComponent />
}



//// [jsxRuntimeDeclarationEmit.d.ts]
// Test case to reproduce jsx-runtime declaration emit issue
// This should trigger the jsx-runtime import without type annotation error
export declare const FunctionComponent: () => JSX.Element;
export declare const AnotherComponent: () => JSX.Element;


//// [DtsFileErrors]


jsxRuntimeDeclarationEmit.d.ts(3,47): error TS2503: Cannot find namespace 'JSX'.
jsxRuntimeDeclarationEmit.d.ts(4,46): error TS2503: Cannot find namespace 'JSX'.


==== jsxRuntimeDeclarationEmit.d.ts (2 errors) ====
    // Test case to reproduce jsx-runtime declaration emit issue
    // This should trigger the jsx-runtime import without type annotation error
    export declare const FunctionComponent: () => JSX.Element;
                                                  ~~~
!!! error TS2503: Cannot find namespace 'JSX'.
    export declare const AnotherComponent: () => JSX.Element;
                                                 ~~~
!!! error TS2503: Cannot find namespace 'JSX'.
    