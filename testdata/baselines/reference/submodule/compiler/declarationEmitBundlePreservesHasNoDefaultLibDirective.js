//// [tests/cases/compiler/declarationEmitBundlePreservesHasNoDefaultLibDirective.ts] ////

//// [extensions.ts]
/// <reference no-default-lib="true"/>
class Foo {
    public: string;
}
//// [core.ts]
interface Array<T> {}
interface Boolean {}
interface Function {}
interface IArguments {}
interface Number {}
interface Object {}
interface RegExp {}
interface String {}


//// [core.js]
//// [extensions.js]
class Foo {
    public;
}
