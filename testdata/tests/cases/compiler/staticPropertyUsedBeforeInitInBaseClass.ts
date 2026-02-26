// @useDefineForClassFields: true
// @target: es2022

class Base {
  p = Derived.S;
}
class Derived extends Base {
  static S = 1;
}
