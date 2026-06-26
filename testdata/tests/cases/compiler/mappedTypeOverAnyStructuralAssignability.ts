// @strict: true
// @moduleResolution: bundler
// @noEmit: true

// Reduced from yup@1.6.1; https://github.com/microsoft/typescript-go/issues/4451
// ObjectSchema<{ foo, bar }> should be assignable to ObjectSchema<any>, but tsgo
// reports the mapped type Shape<any> as missing keys foo, bar; tsc accepts it.
// The contravariant `concat` parameter flips the comparison so Shape<any> is
// checked as the source against Shape<{ foo, bar }>.
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
