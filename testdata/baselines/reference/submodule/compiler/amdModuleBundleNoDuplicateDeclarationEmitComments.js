//// [tests/cases/compiler/amdModuleBundleNoDuplicateDeclarationEmitComments.ts] ////

//// [file1.ts]
/// <amd-module name="mynamespace::SomeModuleA" />
export class Foo {}
//// [file2.ts]
/// <amd-module name="mynamespace::SomeModuleB" />
export class Bar {}

//// [file2.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Bar = void 0;
class Bar {
}
exports.Bar = Bar;
//// [file1.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Foo = void 0;
class Foo {
}
exports.Foo = Foo;
