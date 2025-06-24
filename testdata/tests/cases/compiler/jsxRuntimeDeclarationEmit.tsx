// Reproduces issue #1011: jsx-runtime declaration emit requires type annotation
// https://github.com/microsoft/typescript-go/issues/1011

// @jsx: react-jsx
// @declaration: true  
// @emitDeclarationOnly: true
// @strict: true
// @target: es6
// @module: esnext
// @moduleResolution: node

/// <reference path="/.lib/react16.d.ts" />

// This should produce clean jsx-runtime imports without "unsafe import" errors
export const MyComponent = () => {
  return <div>Hello World</div>
}