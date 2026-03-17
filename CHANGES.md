CHANGES.md lists intentional changes between the Strada (TypeScript) and Corsa (Go) compilers.

# JavaScript support

At a high level, JavaScript support in Corsa is intended to expose TypeScript features in a .js file, working exactly as they do in TypeScript with different syntax.
This differs from Strada, which has many JavaScript features that do not exist in TypeScript at all, and quite a few differences in features that overlap.
For example, Corsa uses the same rule for checking calls in both TypeScript and JavaScript; Strada lets you skip parameters with type `any`.
And because Corsa uses the same rule for optional parameters, it fixes subtle Strada bugs with `"strict": true` in JavaScript.

We primarily want to support people writing modern JavaScript, using things like ES modules, classes, destructuring, etc.
Not CommonJS modules and constructor functions, although those do still work.
However, we have trimmed a lot of unused or underused features.
This makes the implementation much simpler and more like TypeScript.

The biggest single removed area is support for Closure header files--any Closure-specific features, in fact.
The tables below list removed Closure features along with the other removed features.

Reminder: JavaScript support in TypeScript falls into three main categories:

- JSDoc Tags
- Expando declarations
- CommonJS syntax

An expando declaration is when you declare a property just by assigning to it, on a function, class or empty object literal:

```js
function f() {}
f.called = false;
```

## JSDoc Tags and Types

| Name       | Example                                                                                                                                       | Substitute                                                                                                                                   | Note                                                                                  |
| -------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------- |
| UnknownType                | `?`                                                                                                                                           | `any`                                                                                                                                        |                                                                                       |
| NamepathType               | `Module:file~id`                                                                                                                              | `import("file").id`                                                                                                                          | TS has never had semantics for this                                                   |
| `@class`                   | <pre><code>/** @class */</code><br/><code>function C() {</code><br/>    <code>this.p = 1</code><br/><code>}</code></pre> | Infer from <code>this.p=</code> or <code>C.prototype.m=</code>                                                                               | Only inference from <code>this.p=</code> or <code>C.prototype.m=</code> is supported. |
| `@throws`                  | `/** @throws {E} */`                                                                                                                          | Keep the same                                                                                                                                | TS never had semantics for this                                                       |
| `@enum`                    | <pre><code>/\** @enum {number} */</code><br/><code>const E = { A: 1, B: 2 }</code></pre>                                                      | <pre><code>/** @typedef {number} E \*/</code><br/><code>/** @type {Record<string, E>} */</code><br/><code>const E = { A: 1, B: 2 }</code></pre> | Closure feature.                                                                      |
| `@author`                  | `/** @author Finn <finn@treehouse.com> */`                                                                                                    | Keep the same                                                                                                                                | `@treehouse` parses as a new tag in Corsa.                                            |
| Postfix optional type      | `T?`                                                                                                                                          | `T \| undefined`                                                                                                                             | This was legacy in _Closure_                                                          |
| Postfix definite type      | `T!`                                                                                                                                          | `T`                                                                                                                                          | This was legacy in _Closure_                                                          |
| Uppercase synonyms         | `String`, `Void`, `array`                                                                                                                     | `string`, `void`, `Array`                                                                                                                    |                                                                                       |
| JSDoc index signatures     | `Object.<K,V>`                                                                                                                                | `Record<K, V>`                                                                                                                              |                                                                                       |
| Identifier-named typedefs  | `/** @typedef {T} */ typeName;`                                                                                                               | `/** @typedef {T} typeName */`                                                                                                               | Closure feature.                                                                      |
| Closure function syntax    | `function(string): void`                                                                                                                      | `(s: string) => void`                                                                                                                        |                                                                                       |
| Automatic typeof insertion | <pre><code>const o = { a: 1 }</code><br/><code>/\** @type {o} */</code><br/><code>var o2 = { a: 1 }</code></pre>                                               | <pre><code>const o = { a: 1 }</code><br/><code>/\** @type {typeof o} */</code><br/><code>var o2 = { a: 1 }</code></pre>                                       |                                                                                       |
| `@typedef` nested names    | `/** @typedef {1} NS.T */`                                                                                                                    | Translate to .d.ts                                                                                                                           | Also applies to `@callback`                                                           |

## Expando declarations

| Name       | Example                                                                                                                                       | Substitute                                                                                                                                   | Note                                                                                  |
| -------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------- |
| Fallback initialisers                           | `f.x = f.x \|\| init`                                                                                                                                                                                                                                                                                                                                          | `if (!f.x) f.x = init`                                                                                                                                                                                                                            |                                                                 |
| Nested, undeclared expandos                     | <pre><code>var N = {};</code><br/><code>N.X.Y = {}</code></pre>                                                                                                                                                                                                                                                                                              | <pre><code>var N = {};</code><br/><code>N.X = {};</code><br/><code>N.X.Y = {}</code></pre>                                                                                                                                                        | All intermediate expandos have to be assigned. Closure feature. |
| Constructor function whole-prototype assignment | <pre><code>C.prototype = {</code><br/>  <code>m: function() { }</code><br/>  <code>n: function() { }</code><br/><code>}</code></pre>                                                                                                                                                                 | <pre><code>C.prototype.m = function() { }</code><br/><code>C.prototype.n = function() { }</code></pre>                                                                                                                                            | Constructor function feature. See note at end.                  |
| Identifier declarations                         | <pre><code>class C {</code><br/>  <code>constructor() {</code><br/>    <code>/\** @type {T} */</code><br/>    <code>identifier;</code><br/>  <code>}</code><br/><code>}</code></pre> | <pre><code>class C {</code><br/>    <code>/\** @type {T} */</code><br/>    <code>identifier;</code><br/>    <code>constructor() { }</code><br/><code>}</code></pre> | Closure feature.                                                |
| `this` aliases                                  | <pre><code>function C() {</code><br/>  <code>var that = this</code><br/>  <code>that.x = 12</code><br/><code>}</code></pre>                                                                                                                                                                          | <pre><code>function C() {</code><br/>  <code>this.x = 12</code><br/><code>}</code></pre>                                                                                                                              | even better:<br/> <code>class C { this.x = 12 }</code>          |     |
| `this` alias for `globalThis`                   | `this.globby = true`                                                                                                                                                                                                                                                                                                                                         | `globalThis.globby = true`                                                                                                                                                                                                                        | When used at the top level of a script                          |

## CommonJS syntax

| Name                                          | Example                                                                                                                                                                                                                                   | Substitute                                                                                                  | Note                                                                    |
| --------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------- |
| Nested, undeclared exports                    | `exports.N.X.p = 1`                                                                                                                                                                                                                       | <pre><code>exports.N = {}</code><br/><code>exports.N.X = {}</code><br/><code>exports.N.X.p = 1</code></pre> | Same as expando rules.                                                  |
| Ignored empty module.exports assignment       | `module.exports = {}`                                                                                                                                                                                                                     | Delete this line                                                                                            | People used to write in this in case module.exports was not defined.    |
| `this` alias for `module.exports`             | `this.p = 1`                                                                                                                                                                                                                              | `exports.p = 1`                                                                                             | When used at the top level of a CommonJS module.                        |
| Multiple assignments narrow with control flow | <pre><code>if (isWindows) {</code><br/><code>  exports.platform = 'win32'</code><br/><code>} else {</code><br/>  <code>exports.platform = 'posix'</code><br/><code>}</code></pre> | Keep the same in most cases                                                                                 | This now unions instead; most uses have the same type in both branches. |
| Single-property access `require`              | `var readFile = require('fs').readFile`                                                                                                                                                                                                   | `var { readFile } = require('fs')`                                                                          |                                                                         |
| Aliasing of `module.exports`                  | <pre><code>var mod = module.exports</code><br/><code>mod.x = 1</code></pre>                                                                                                                                                               | `module.exports.x = 1`                                                                                      |                                                                         |

# Component-Level Changes

## Scanner

1. Node positions use UTF8 offsets from the beginning of the file, not UTF16 offsets. Node positions in files with non-ASCII characters will be greater than before.

## Parser

1. Malformed `...T?` at the end of a tuple now fails with a parse error instead of a grammar error.
2. Malformed string ImportSpecifiers (`import x as "OOPS" from "y"`) now contain the string's text instead of an empty identifier.
3. Empty binding elements no longer have a separate kind for OmittedExpression. Instead they have Kind=BindingElement with a nil Initialiser, Name and DotDotDotToken.
4. ShorthandPropertyAssignment no longer includes an EqualsToken as a child when it has an ObjectAssignmentInitializer.
5. JSDoc nodes now include leading whitespace in their location.
6. The parser always parses a JSDocText node for comments in JSDoc. `string` is no longer part of the type of `comment`.
7. In cases where Strada did produce a JSDocText node, Corsa no longer (incorrectly) includes all leading and trailing whitespace/asterisks, as well as initial `/**`.
8. JSDocMemberName is now parsed as QualifiedName. These two nodes previously only differed by type, and now QualifiedName has a much less restrictive type for its left child.

JSDoc types are parsed in normal type annotation position but show a grammar error. Corsa no longer parses the JSDoc types below, giving a parse error instead of a grammar error.

1. No postfix `T?` and `T!` types. Prefix `?T` and `!T` are still parsed and `!T` continues to have no semantics.
2. No Closure `function(string,string): void` types.
3. No JSDoc standalone `?` type.
4. No JSDoc module namepaths: `module:folder/file.C`

Corsa no longer parses the following JSDoc tags with a specific node type. They now parse as generic JSDocTag nodes.

1. `@class`/`@constructor`
2. `@throws`
3. `@author`
4. `@enum`

## Checker

### Miscellaneous

#### Unused type parameters now report TS6196 instead of TS6133.

When `noUnusedParameters` is enabled, unused type parameters are reported with error TS6196 ("'T' is declared but never used") instead of TS6133 ("'T' is declared but its value is never read").
TS6133 is for value-space identifiers; TS6196 is the correct code for type parameters, which have no runtime value.

#### Type assignability errors now use more specific error codes instead of TS2345.

Where Strada reported the generic TS2345 ("Argument of type 'X' is not assignable to parameter of type 'Y'"), Corsa reports a more specific code that identifies the root cause:

- TS2739 — Type is missing properties from another type
- TS2740 — Type is missing some properties from another type
- TS2741 — Property is missing in type but required in another type
- TS2322 — Type is not assignable to type (used in assignment/return contexts)

The diagnostics convey the same semantic error with more actionable detail.

#### Circular import alias errors (TS2303) are now also reported on export statements.

Strada only reported TS2303 on import declarations. Corsa reports the same error on both imports and exports when a circular alias is detected.

#### Duplicate identifier errors are now reported individually (TS2300) rather than batched (TS6200).

Related error spans and wording have also been updated.

#### Isolated declaration error codes have changed.

Some error codes for `isolatedDeclarations` violations have changed (e.g. TS9023 → TS9013). The errors still describe the same isolation constraints.

#### Type assignability error wording has changed for call and construct signatures.

Where Strada said "Call signature return types are incompatible" or wrapped the root cause in "Types of construct signatures are incompatible. Type X is not assignable to type Y.", Corsa reports "Type X is not assignable to type Y." directly. The meaning is the same.

#### With `"strict": false`, Corsa no longer allows omitting arguments for parameters with type `undefined`, `unknown`, or `any`:

```js
/** @param {unknown} x */
function f(x) {
  return x;
}
f(); // Previously allowed, now an error
```

`void` can still be omitted, regardless of strict mode:

```js
/** @param {void} x */
function f(x) {
  return x;
}
f(); // Still allowed
```

#### Strada's JS-specific rules for inferring type arguments no longer apply in Corsa.

Inferred type arguments may change. For example:

```js
/** @type {any} */
var x = { a: 1, b: 2 };
var entries = Object.entries(x);
```

In Strada, `entries: Array<[string, any]>`. In Corsa it has type `Array<[string, unknown]>`, the same as in TypeScript.

#### Values are no longer resolved as types in JSDoc type positions.

```js
/** @typedef {FORWARD | BACKWARD} Direction */
const FORWARD = 1,
  BACKWARD = 2;
```

Must now use `typeof` the same way TS does:

```js
/** @typedef {typeof FORWARD | typeof BACKWARD} Direction */
const FORWARD = 1,
  BACKWARD = 2;
```

### JSDoc Types

#### JSDoc variadic types are now only synonyms for array types.

```js
/** @param {...number} ns */
function sum(...ns) {}
```

is equivalent to

```js
/** @param {number[]} ns */
function sum(...ns) {}
```

They have no other semantics.

#### A variadic type on a parameter no longer makes it a rest parameter. The parameter must use standard rest syntax.

```js
/** @param {...number} ns */
function sum(ns) {}
```

Must now be written as

```js
/** @param {...number} ns */
function sum(...ns) {}
```

#### The postfix `=` type no longer adds `undefined` even when `strictNullChecks` is off

This is a bug in Strada: it adds `undefined` to the type even when `strictNullChecks` is off.
This bug is fixed in Corsa.

```js
/** @param {number=} x */
function f(x) {
  return x;
}
```

will now have `x?: number` not `x?: number | undefined` with `strictNullChecks` off.
Regardless of strictness, it still makes parameters optional when used in a `@param` tag.

### JSDoc Tags

#### `asserts` annotation for an arrow function must be on the declaring variable, not on the arrow itself. This no longer works:

```js
/**
 * @param {A} a
 * @returns { asserts a is B }
 */
const foo = (a) => {
  if (/** @type { B } */ (a).y !== 0) throw TypeError();
  return undefined;
};
```

And must be written like this:

```js
/**
 * @type {(a: A) => asserts a is B}
 */
const foo = (a) => {
  if (/** @type { B } */ (a).y !== 0) throw TypeError();
  return undefined;
};
```

This is identical to the TypeScript rule.

#### Error messages on async functions that incorrectly return non-Promises now use the same error as TS.

#### `@typedef` and `@callback` in a class body are hoisted outside the class.

This means the declarations are accessible outside the class and may conflict with similarly named declarations in the outer scope.

#### `@class` or `@constructor` does not make a function into a constructor function.

Corsa ignores `@class` and `@constructor`. This makes a difference on a function without this-property assignments or associated prototype-function assignments.

#### `@param` tags now apply to at most one function.

If they're in a place where they could apply to multiple functions, they apply only to the first one.
If you have `"strict": true`, you will see a noImplicitAny error on the now-untyped parameters.

```js
/** @param {number} x */
var f = (x) => x,
  g = (x) => x;
```

#### Optional marking on parameter names now makes the parameter both optional and undefined:

```js
/** @param {number} [x] */
function f(x) {
  return x;
}
```

This behaves the same as TypeScript's `x?: number` syntax.
Strada makes the parameter optional but does not add `undefined` to the type.

#### Type assertions with `@type` tags now prevent narrowing of the type.

```js
/** @param {C | undefined} cu */
function f(cu) {
  if (/** @type {any} */ (cu).undeclaredProperty) {
    cu; // still has type C | undefined
  }
}
```

In Strada, `cu` incorrectly narrows to `C` inside the `if` block, unlike with TS assertion syntax.
In Corsa, the behaviour is the same between TS and JS.

### Expandos

#### Expando assignments of `void 0` are no longer ignored as a special case:

```js
var o = {};
o.y = void 0;
```

creates a property `y: undefined` on `o` (which will widen to `y: any` if strictNullChecks is off).

#### A this-property expression with a type annotation in the constructor no longer creates a property:

```js
class SharedClass {
  constructor() {
    /** @type {SharedId} */
    this.id;
  }
}
```

Provide an initializer or use a property declaration in the class body:

```js
class SharedClass1 {
  /** @type {SharedId} */
  id;
}
class SharedClass2 {
  constructor() {
    /** @type {SharedId} */
    this.id = 1;
  }
}
```

#### Assigning an object literal to the `prototype` property of a function no longer makes it a constructor function:

```js
function Foo() {}
Foo.prototype = {
  /** @param {number} x */
  bar(x) {
    return x;
  },
};
```

If you still need to use constructor functions instead of classes, you should declare methods individually on the prototype:

```js
function Foo() {}
/** @param {number} x */
Foo.prototype.bar = function (x) {
  return x;
};
```

Although classes are a much better way to write this code.

### CommonJS

#### Initializing exports to `undefined`:

To accommodate the pattern of initializing CommonJS exports to `undefined` (sometimes written as `void 0`) and then subsequently assigning their intended values, when CommonJS exports have multiple assignments and an initial assignment of `undefined`, the `undefined` is ignored when determining the type of the export.

```js
exports.foo = exports.bar = void 0;
// Later in the same file...
exports.foo = 123  // Exported type is `123`
exports.bar = "abc"  // Exported type is `"abc"`
```

## Extreme Trivia

These changes are listed only for the purposes of completely accounting for all error diffs.

#### Multi-overload "no overload matches" errors now show only the last overload's failure.

Where Strada listed every overload's failure ("Overload 1 of N, '(sig)', gave the following error. ..."), Corsa shows only "The last overload gave the following error." with a TS2771 related info span pointing to the last overload's declaration. The error position may also shift to highlight the specific argument that failed.

#### Optional parameter types no longer include `| undefined` in most error messages.

Corsa omits the redundant `| undefined` from optional parameter types in error messages. For example, `(p?: string | undefined) => T` is printed as `(p?: string) => T`.

#### Type names may be displayed differently in error messages.

Some type names are represented differently:

- `DateConstructor` is now displayed as `typeof Date`.
- Unqualified names may be fully qualified (e.g. `ReactNode` → `React.ReactNode`, `connectExport` → `server.connectExport`).
- Mixin anonymous class types omit their generic type arguments in display.
- `undefined` optional properties may appear as `never` in certain intersection contexts.

#### String literal types use single quotes in type displays.

Where Strada printed string literal types with double quotes (`"value"`), Corsa uses single quotes (`'value'`).

#### When all destructured elements are unused, a single error replaces individual TS6133 errors.

Where Strada reported individual TS6133 ("'x' is declared but its value is never read") for each unused element in a destructuring pattern, Corsa reports a single TS6198 ("All destructured elements are unused") covering the whole pattern. Similarly, when all type parameters are unused, TS6205 ("All type parameters are unused") replaces individual TS6133 errors.

#### Unused identifier error spans are narrower.

For unused imports, the error squiggle covers only the imported name rather than the entire import statement. For unused type parameters, the span covers only the parameter name (not the surrounding angle brackets).

#### Related info message text is displayed inline with its source location.

Where Strada printed related info text on a separate line below the file:line:col location, Corsa prints it on the same line: `file.ts:1:5 - message text`.

#### Duplicate identifier errors are now also reported on the first declaration.

Strada reported TS2300 only on the second (and later) occurrences of a duplicate identifier. Corsa also reports it on the first occurrence.

#### The 'Object' class name restriction (TS2725) has been removed.

Strada reported TS2725 when a class inside a namespace was named `Object` when targeting ES5 with CommonJS. Corsa removes this check.