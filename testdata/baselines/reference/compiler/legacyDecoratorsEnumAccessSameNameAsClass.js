//// [tests/cases/compiler/legacyDecoratorsEnumAccessSameNameAsClass.ts] ////

//// [legacyDecoratorsEnumAccessSameNameAsClass.ts]
export enum MyEnum {
    Foo = "FooValue",
    Bar = "BarValue",
}

function myDecorator(target: any) {
    return target;
}

@myDecorator
export class Foo {
    type: MyEnum = MyEnum.Foo;

    getType(): MyEnum {
        return this.type || MyEnum.Foo;
    }
}


//// [legacyDecoratorsEnumAccessSameNameAsClass.js]
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var Foo_1;
export var MyEnum;
(function (MyEnum) {
    MyEnum["Foo"] = "FooValue";
    MyEnum["Bar"] = "BarValue";
})(MyEnum || (MyEnum = {}));
function myDecorator(target) {
    return target;
}
let Foo = Foo_1 = class Foo {
    type = MyEnum.Foo_1;
    getType() {
        return this.type || MyEnum.Foo_1;
    }
};
Foo = Foo_1 = __decorate([
    myDecorator
], Foo);
export { Foo };
