//// [tests/cases/conformance/override/overrideLateBindableIndexSignature1.ts] ////

//// [overrideLateBindableIndexSignature1.ts]
const sym: symbol = Symbol();

class Base1 {
  [sym]() {}
}

class Derived1 extends Base1 {
  override [sym]() {}
}

class Base2 {
  [sym]() {}
}

class Derived2 extends Base2 {
  [sym]() {}
}

class Base3 {}

class Derived3 extends Base3 {
  override [sym]() {}
}




//// [overrideLateBindableIndexSignature1.d.ts]
const sym: symbol;
class Base1 {
    [sym]: () => void;
}
class Derived1 extends Base1 {
    [sym]: () => void;
}
class Base2 {
    [sym]: () => void;
}
class Derived2 extends Base2 {
    [sym]: () => void;
}
class Base3 {
}
class Derived3 extends Base3 {
    [sym]: () => void;
}
