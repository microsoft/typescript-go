// @strict: true
// @noEmit: true

interface A { foo: { bar(x: string): string } }
interface B { foo: { bar(x: string): number } }
declare let a: A;
declare let b: B;
a = b;

interface C { foo: { bar: { baz(x: string): string } } }
interface D { foo: { bar: { baz(x: string): number } } }
declare let c: C;
declare let d: D;
c = d;

interface E { foo: { bar: { baz: string } } }
interface F { foo: { bar: { baz: number } } }
declare let e: E;
declare let f: F;
e = f;

interface G { foo: { bar(): { baz: string } } }
interface H { foo: { bar(): { baz: number } } }
declare let g: G;
declare let h: H;
g = h;
