//// [tests/cases/compiler/jsxRuntimeDeclarationEmit.tsx] ////

//// [jsxRuntimeDeclarationEmit.tsx]
// Test case to reproduce jsx-runtime declaration emit issue


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
export declare const FunctionComponent: () => any;
export declare const AnotherComponent: () => any;
