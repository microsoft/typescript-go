// @strict: true
// @strictArrayVariance: true, false
// @noEmit: true

// Mutable Array: covariant element typing is unsound when the callee mutates.
function widen(xs: (string | number)[]) {
    xs.push(3);
}

const ys: string[] = ["a"];
widen(ys);

// ReadonlyArray stays covariant in T; passing readonly string[] should remain ok.
function read(xs: ReadonlyArray<string | number>) {
    return xs[0];
}

const rs: readonly string[] = ["b"];
read(rs);

// `Array<T>` syntax uses the same `globalArrayType` variance as `(T)[]`.
function idArray(xs: Array<string | number>) {
    xs.push(3);
}

const zs: Array<string> = ["c"];
idArray(zs);

// Type alias to `Array<T>` should use the same variance as `Array<T>` / `(T)[]`.
type ArrayAlias<T> = Array<T>;
function idAlias(xs: ArrayAlias<string | number>) {
    xs.push(3);
}

const ws: ArrayAlias<string> = ["d"];
idAlias(ws);

// Mutable tuples: element positions are invariant under strictArrayVariance.
function widenTuple(xs: [string | number, string]) {
    xs[0] = 1;
}

const tup: [string, string] = ["x", "y"];
widenTuple(tup);

// Nested `Array<Array<U>>`: element type `U[]` is still `globalArrayType`.
function nest(xs: (string | number)[][]) {
    xs[0].push(3);
}

const row: string[] = ["p"];
const grid: string[][] = [row];
nest(grid);

// Readonly tuples use tuple variance (covariant), like ReadonlyArray.
function rtup(xs: readonly [string | number, string]) {
    return xs[0];
}

const rt: readonly [string, string] = ["u", "v"] as const;
rtup(rt);

// Structural assignability: object with `string[]` field to type expecting `(string | number)[]`.
interface Hold {
    xs: (string | number)[];
}

function mutateCell(h: Hold) {
    h.xs.push(1);
}

const cell: { xs: string[] } = { xs: ["q"] };
mutateCell(cell);

// Inference: generic element extraction via `infer U` still works; the instance of
// Array<string> has its type argument inferred as `string` regardless of the flag.
type ElemOf<T> = T extends Array<infer U> ? U : never;
type E1 = ElemOf<string[]>; // string
const e1: E1 = "ok";

// Inference: a generic identity function should infer T from the argument type.
// Under strictArrayVariance this still works because T is inferred, not constrained
// by a widening target.
function idGeneric<T>(xs: T[]): T[] {
    return xs;
}

const inferred = idGeneric(ys); // inferred: string[]
const checkInferred: string[] = inferred;

// Inference: when the target is wider than the source and both are mutable Array,
// the strict flag rejects the call (same as the widen cases above) rather than
// silently widening the inferred type argument.
function takeWider(xs: Array<string | number>) {}
takeWider(ys);

// Spread alone does NOT widen under this flag: `[...ys]` has type `string[]`,
// so this still errors. The diagnostic message must not suggest otherwise.
widen([...ys]);

// The primary migration target: `readonly` at the parameter boundary stays
// covariant under the flag, so passing `string[]` to a consumer typed for
// `readonly (string | number)[]` is still accepted.
function readWider(xs: readonly (string | number)[]) {
    return xs.length;
}
readWider(ys);

// Callback contravariance is orthogonal to array variance and must keep
// working: a callback accepting the wider element type is assignable to a
// parameter typed for the narrower one.
function forEachString(xs: readonly string[], cb: (x: string) => void) {
    for (const x of xs) cb(x);
}
const wideCb: (x: string | number) => void = () => {};
forEachString(ys, wideCb);
