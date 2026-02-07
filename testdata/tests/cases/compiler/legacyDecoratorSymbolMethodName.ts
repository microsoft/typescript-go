// @target: es2023
// @module: nodenext
// @strict: true
// @experimentalDecorators: true
// @emitDecoratorMetadata: true

export class A {
  @(fakeDecorator as any)
  [Symbol()]() {}
}

function fakeDecorator() {}
