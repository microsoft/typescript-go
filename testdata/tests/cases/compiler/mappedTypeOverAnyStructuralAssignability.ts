// @strict: true
// @moduleResolution: bundler
// @noEmit: true

// Reduced from yup@1.6.1; https://github.com/microsoft/typescript-go/issues/4451
// ObjectSchema<MyValues> should be assignable to AnyObjectSchema (= ObjectSchema<any, ...>),
// but tsgo reports the mapped type Shape<any> as missing keys foo, bar; tsc accepts it.

type Maybe<T> = T | null | undefined;
type Optionals<T> = Extract<T, null | undefined>;
type Defined<T> = T extends undefined ? never : T;
type _<T> = T extends {} ? {
    [k in keyof T]: T[k];
} : T;
type Flags = 's' | 'd' | '';
type ResolveFlags<T, F extends Flags, D = T> = Extract<F, 'd'> extends never ? T : D extends undefined ? T : Defined<T>;
declare abstract class Schema<TType = any, TContext = any, TDefault = any, TFlags extends Flags = ''> implements ISchema<TType, TContext, TFlags, TDefault> {
    readonly __outputType: ResolveFlags<TType, TFlags, TDefault>;
}
declare class Reference<TValue = unknown> {
}
type ObjectShape = {
};
type AnyObject = {
};
type ResolveStrip<T extends ISchema<any>> = T extends ISchema<any, any, infer F> ? Extract<F, 's'> extends never ? T['__outputType'] : never : T['__outputType'];
type TypeFromShape<S extends ObjectShape, _C> = {
    [K in keyof S]: S[K] extends ISchema<any> ? ResolveStrip<S[K]> : S[K] extends Reference<infer T> ? T : unknown;
};
type DefaultFromShape<Shape extends ObjectShape> = {
};
type MakePartial<T extends object> = {
};
interface ISchema<T, C = any, F extends Flags = any, D = any> {
    __outputType: T;
}
interface ValidateOptions<TContext = {}> {
    /**
     */
}
interface MessageParams {
}
type Message<Extra extends Record<string, unknown> = any> = string | ((params: Extra & MessageParams) => unknown) | Record<PropertyKey, unknown>;
declare function create$6(): StringSchema;
declare class StringSchema<TType extends Maybe<string> = string | undefined, TContext = AnyObject, TDefault = undefined, TFlags extends Flags = ''> extends Schema<TType, TContext, TDefault, TFlags> {
    length(length: number | Reference<number>, message?: Message<{
    }>): this;
    required(msg?: Message): StringSchema<NonNullable<TType>, TContext, TDefault, TFlags>;
}
declare function create$5<T extends number, TContext extends Maybe<AnyObject> = AnyObject>(): NumberSchema<T | undefined, TContext>;
declare class NumberSchema<TType extends Maybe<number> = number | undefined, TContext = AnyObject, TDefault = undefined, TFlags extends Flags = ''> extends Schema<TType, TContext, TDefault, TFlags> {
    min(min: number | Reference<number>, message?: Message<{
    }>): this;
    required(msg?: Message): NumberSchema<NonNullable<TType>, TContext, TDefault, TFlags>;
}
type MakeKeysOptional<T> = T extends AnyObject ? _<MakePartial<T>> : T;
type Shape<T extends Maybe<AnyObject>, C = any> = {
    [field in keyof T]-?: ISchema<T[field], C> | Reference;
};
declare function create$3<C extends Maybe<AnyObject> = AnyObject, S extends ObjectShape = {}>(spec?: S): ObjectSchema<_<TypeFromShape<S, C>>, C, _<DefaultFromShape<S>>, "">;
interface ObjectSchema<TIn extends Maybe<AnyObject>, TContext = AnyObject, TDefault = any, TFlags extends Flags = ''> extends Schema<MakeKeysOptional<TIn>, TContext, TDefault, TFlags> {
    nullable(msg?: Message): ObjectSchema<TIn | null, TContext, TDefault, TFlags>;
    fields: Shape<NonNullable<TIn>, TContext>;
    concat(schema: this): this;
    partial(): ObjectSchema<Partial<TIn>, TContext, TDefault, TFlags>;
}
declare class ArraySchema<TIn extends any[] | null | undefined, TContext, TDefault = undefined, TFlags extends Flags = ''> extends Schema<TIn, TContext, TDefault, TFlags> {
    length(length: number | Reference<number>, message?: Message<{
    }>): this;
}
type AnyTuple = [unknown, ...unknown[]];
declare function create$1<T extends AnyTuple>(schemas: {
}): TupleSchema<T | undefined, AnyObject, undefined, "">;
declare namespace create$1 {
}
declare class TupleSchema<TType extends Maybe<AnyTuple> = AnyTuple | undefined, TContext = AnyObject, TDefault = undefined, TFlags extends Flags = ''> extends Schema<TType, TContext, TDefault, TFlags> {
}
type AnyObjectSchema = ObjectSchema<any, any, any, any>;
interface MyValues { foo: string; bar: number; }
const specific: ObjectSchema<MyValues> = create$3({ foo: create$6().required(), bar: create$5().required() });
const test: AnyObjectSchema = specific;