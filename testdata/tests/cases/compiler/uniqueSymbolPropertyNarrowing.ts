// @strict: true
// @noEmit: true

// Unique symbol defined as attribute should be narrowed correctly
const Sym: unique symbol = Symbol("test");
const Objs = () => {};
Objs.Sym = Sym;

type Val = { attr: 1 } | typeof Objs.Sym;
const vals: Val[] = [Objs.Sym];
console.log(vals[0] != Objs.Sym && vals[0].attr);

// Simpler case that should also work (and already does)
const Sym2: unique symbol = Symbol("test2");
type Val2 = { attr: 1 } | typeof Sym2;
const vals2: Val2[] = [Sym2];
console.log(vals2[0] != Sym2 && vals2[0].attr);
