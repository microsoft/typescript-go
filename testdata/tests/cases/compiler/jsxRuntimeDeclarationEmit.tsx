// Test case to reproduce jsx-runtime declaration emit issue

// @jsx: react-jsx
// @declaration: true  
// @emitDeclarationOnly: true
// @strict: true
// @target: esnext
// @module: esnext

// This should trigger the jsx-runtime import without type annotation error

export const FunctionComponent = () => {
  return <div>Hello World</div>
}

export const AnotherComponent = () => {
  return <FunctionComponent />
}