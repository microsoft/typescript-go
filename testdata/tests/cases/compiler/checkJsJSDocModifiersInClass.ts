// @allowJs: true
// @checkJs: true
// @noEmit: true
// @filename: test.js

// Unlike object literal members, JSDoc modifiers on class members are still
// semantically validated under checkJs: @override without a base class errors.

class NoBase {
  /** @override */
  foo() {
  }
}

class Base {
  foo() {
  }
}

class Derived extends Base {
  /** @override */
  foo() {
  }
}
