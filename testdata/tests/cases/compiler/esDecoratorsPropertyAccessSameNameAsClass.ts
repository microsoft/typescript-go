// @target: es2020

export enum MyEnum {
    Foo = "FooValue",
    Bar = "BarValue",
}

function myDecorator(target: any, context: ClassDecoratorContext) {
    return target;
}

// ES decorators with enum member access sharing a class name
@myDecorator
export class Foo {
    type: MyEnum = MyEnum.Foo;

    getType(): MyEnum {
        return this.type || MyEnum.Foo;
    }
}

// Property access on objects
declare const obj: { Bar: string };

@myDecorator
export class Bar {
    prop = obj.Bar;

    method() {
        return obj.Bar;
    }
}

// Static member of another class
class Other {
    static Baz = 42;
}

@myDecorator
export class Baz {
    prop = Other.Baz;
}
