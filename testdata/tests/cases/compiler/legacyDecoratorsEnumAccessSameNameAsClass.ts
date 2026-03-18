// @target: esnext
// @experimentalDecorators: true

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
