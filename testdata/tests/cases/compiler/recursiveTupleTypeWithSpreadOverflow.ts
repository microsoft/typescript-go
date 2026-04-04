// @target: es2015
// @strict: true
// @noEmit: true

// Repro from #63270 - should not crash with stack overflow

type Recur<T> = (
    T extends (unknown[]) ? {} : { [K in keyof T]?: Recur<T> }
) | [...Recur<T>[number][]];

function join<T>(l: Recur<T>[]): Recur<T> {
    return ['marker', ...l];
}

function a<T>(l: Recur<T>[]): void {
    const x: Recur<T> | undefined = join(l);
}
