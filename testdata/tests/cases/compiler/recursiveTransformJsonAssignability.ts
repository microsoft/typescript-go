// @noEmit: true
// @strict: true
// @target: es2020

// https://github.com/microsoft/typescript-go/issues/4408

type Transform<Base, From, To> = {
    [K in keyof Base]: Exclude<Base[K], undefined | null> extends never ? Base[K]
        : Exclude<Base[K], undefined | null> extends object
            ? Exclude<Base[K], undefined | null> extends From ? To | Extract<Base[K], null | undefined> : Transform<Base[K], From, To>
        : Base[K];
};

type TransformJson<T> = Transform<T, Date | bigint, string>;

type Shared = {
    base: {
        id: string;
        meta: {
            tags: string[];
        };
    };
};

type Variant = Shared & (
    | { type: 1; value: { a: string } }
    | { type: 2; value: { b: string } }
    | { type: 3; value: { c: string } }
    | { type: 4; value: { d: string } }
    | { type: 5; value: { e: string } }
    | { type: 6; value: { f: string } }
    | { type: 7; value: { g: string } }
    | { type: 8; value: { h: string } }
);

type X = {
    root: {
        levels: Array<{
            items: Array<{
                variants: Variant[];
            }>;
        }>;
    };
};

declare const transformed: TransformJson<X>;
const x: X = transformed;
