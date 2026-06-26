// @strict: true
// @moduleResolution: bundler
// @noEmit: true

// Reduced from yup@1.6.1; https://github.com/microsoft/typescript-go/issues/4451
// ObjectSchema<{ foo, bar }> should be assignable to ObjectSchema<any>. tsgo
// previously rejected this because getRecursionIdentity gave from-type-node
// references their own object identity instead of their symbol, so the deeply
// nested self-referential ObjectSchema instantiations were never recognized as
// recursive and depth-limited, and the comparison ran to a spurious failure.
type Shape<T> = { [field in keyof T]-?: T[field]; };
interface ObjectSchema<TIn> {
    __outputType: TIn extends {} ? {} : TIn;
    nullable(): ObjectSchema<TIn | null>;
    fields: Shape<NonNullable<TIn>>;
    concat(schema: this): this;
    partial(): ObjectSchema<Partial<TIn>>;
}
declare const specific: ObjectSchema<{ foo: string; bar: number }>;
const test: ObjectSchema<any> = specific;
