type Type = { values: string };

type Wrapped<T extends Type> = { values: T["values"] };

type FromObject<
    Context extends Record<string, Type>,
    Props extends Context,
    Shape extends Record<keyof Props, Type> = {
        [Key in keyof Props]: Wrapped<FromSchema<Context, Props[Key]>>;
    },
> = Shape;

type FromSchema<Context extends Record<string, Type>, S> =
    S extends Type ? FromSchema<Context, Context[keyof Context]> : never;
