// @strict: true
// @noEmit: true
// @stableTypeOrdering: true

// https://github.com/microsoft/typescript-go/issues/4817

type MaybeRef<T> = T | { value: T };

declare function f<D = string>(options: {
    select?: MaybeRef<(data: string) => D>;
}): D | undefined;

declare const arg: {
    select?: MaybeRef<(data: string) => string>;
} & {
    select: (data: string) => number;
};

const n: number | undefined = f(arg);

declare const reversedArg: {
    select: (data: string) => number;
} & {
    select?: MaybeRef<(data: string) => string>;
};

const s: string | undefined = f(reversedArg);
