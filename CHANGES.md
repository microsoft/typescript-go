CHANGES.md lists intentional changes between the Strada (Typescript) and Corsa (Go) compilers.

# Feature-Level Description

At a high level, Javascript support in Corsa is intended to expose Typescript features in a .js file, working exactly as they do in Typescript with different syntax.
This differs from Strada, which has many Javascript features that do not exist in Typescript at all, and quite a few differences in features that overlap.
For example, Corsa uses the same rule for checking calls in both Typescript and Javascript; Strada lets you skip parameters with type `any`.
And because Corsa uses the same rule for optional parameters, it fixes subtle Strada bugs with `"strict": true` in Javascript.

We primarily want to support people writing modern Javascript, using things like ES modules, classes, destructuring, etc.
Not CommonJS modules and constructor functions, although those do still work.
However, we have trimmed a lot of unused or underused features.

The biggest single removed area is support for Closure header files--any Closure-specific features, in fact.
The tables below list removed Closure features along with the other removed features.

Reminder: Javascript support in Typescript falls into three main categories:

- JSDoc Tags
- Expando declarations
- CommonJS syntax

An expando declaration is when you declare a property just by assigning to it, on a function, class or empty object literal:

```js
function f() {}
f.called = false;
```


## JSDoc Tags and Types

| Name | Example | Substitute | Note  |
| ---- | ------- | ---------- | ----- |
| UnknownType  | `?`              | `any`               |                                     |
| NamepathType | `Module:file~id` | `import("file").id` | TS has never had semantics for this |
| `@class`     | `/** @class \*/`<br/>`function C() {`<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`this.p = 1`<br/>`}` |  | Only inference from `this.p=` or `C.prototype.m=` is supported. |
|`@throws`     |`/** @throws {E} */` | Keep the same    | TS never had semantics for this     |
|`@enum` | `/** @enum {number} */`<br/> `const E = { A: 1, B: 2 }`  |`/** @typedef {number} E */`<br/>`/** @type {Record<string, E>}`<br/>`const E = { A: 1, B: 2 }` | Closure feature.  |
| `@author`    | `/** @author Finn <finn@treehouse.com> */` | Keep the same  | `@treehouse` parses as a new tag in Corsa. |
| Postfix optional type      | `T?` | `T \| undefined`  | This was legacy in *Closure*  |
| Postfix definite type      | `T!` | `T`               | This was legacy in *Closure*  |
| Uppercase synonyms         | `String`, `Void`, `array`| `string`, `void`, `Array`  | |
| JSDoc index signatures     | `Object.<K,V>`           | `{ [x: K]: V }`  |  |
| Identifier-named typedefs  | `/** @typedef {T} */ typeName;` | `/** @typedef {T} typeName */`  | Closure feature.   |
| Closure function syntax    | `function(string): void` | `(s: string) => void` |  |
| Automatic typeof insertion | `const o = { a: 1 }`<br/>`/** @type {o} */ var o2 = { a: 1 }` | `const o = { a: 1 }`<br/>`/** @type {typeof o} */ var o2 = { a: 1 }` |  |
| `@typedef` nested names    | `/** @typedef {1} NS.T */` | Translate to .d.ts  | Also applies to `@callback` |

## Expando declarations

| Name | Example | Substitute | Note  |
| ---- | ------- | ---------- | ----- |
| Fallback initialisers                           | `f.x = f.x \| init` | `if (!f.x) f.x = init` | |
| Nested, undeclared expandos                     | `var N = {};`<br/>`N.X.Y = {}` | `var N = {};`<br/>`N.X = {};`<br/>`N.X.Y = {}` | All intermediate expandos have to be assigned. Closure feature.  |
| Constructor function whole-prototype assignment | `C.prototype = {`<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`m: function() { }`<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`n: function() { }`<br/>`}` | `C.prototype.m = function() { }`<br/>`C.prototype.n = function() { }`                                    | Constructor function feature. See note at end.                     |
| Identifier declarations                         | `class C {`<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`constructor() {`<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`/** @type {T} */`<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`identifier;`<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`}`<br/>`}` | `class C {`<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`/** @type {T} */`<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`identifier;`<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`constructor() { }`<br/>`}`                                 | Closure feature. |
| `this` aliases                                  | `function C() {`<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`var that = this`<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`that.x = 12`<br/>`}` | `function C() {`<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`this.x = 12`<br/>`}` | even better:<br/> `class C {`<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`this.x = 12 `<br/>`}` |                                                                    |
| `this` alias for `globalThis`                   | `this.globby = true` | `globalThis.globby = true` | When used at the top level of a script |

## CommonJS syntax

| Name | Example | Substitute | Note  |
| ---- | ------- | ---------- | ----- |
| Nested, undeclared exports                   | `exports.N.X.p = 1`   | `exports.N = {}`<br/>`exports.N.X = {}`<br/>`exports.N.X.p = 1` | Same as expando rules. |
| Ignored empty module.exports assignment      | `module.exports = {}` | Delete this line   | People used to write in this in case module.exports was not defined. |
| `this` alias for `module.exports`            | `this.p = 1`          | `exports.p = 1`    | When used at the top level of a CommonJS module. |
| Multiple assignments narrow with control flow| `if (isWindows) {`<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`exports.platform = 'win32'`<br/>`} else {`<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`exports.platform = 'posix'`<br/>`}` | Keep the same in most cases | This now unions instead; most uses have the same type in both branches. |
| Single-property access `require`             | `var readFile = require('fs').readFile`    | `var { readFile } = require('fs')` | |
| Aliasing of `module.exports`                 | `var mod = module.exports`<br/>`mod.x = 1` | `module.exports.x = 1` | |

## Features yet to be implemented

Object.defineProperty for CommonJS exports and expandos. The compiler treats this as an alternate to the usual assignment syntax:

```js
function f() { }
Object.defineProperty(f, "p", { value: 1, writable: true })
````

# Component-Level Description

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

#### When `"strict": false`, Corsa no longer allows omitting arguments for parameters with type `undefined`, `unknown`, or `any`:

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

In Strada, `entries: Array<[string, any]>`.
In Corsa it has type `Array<[string, unknown]>`, the same as in TypeScript.

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

This is identical to the Typescript rule.

#### Error messages on async functions that incorrectly return non-Promises now use the same error as TS.

#### `@typedef` and `@callback` in a class body are no longer accessible outside the class.

They must be moved outside the class to use them outside the class.

#### `@class` or `@constructor` does not make a function into a constructor function.

Corsa ignores `@class` and `@constructor`.
This makes a difference on a function without this-property assignments or associated prototype-function assignments.

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

This behaves the same as Typescript's `x?: number` syntax.
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

#### Chained exports no longer work:

```js
exports.x = exports.y = 12;
```

Now only exports `x`, not `y` as well.

#### Exporting `void 0` is no longer ignored as a special case:

```js
exports.x = void 0;
// several lines later...
exports.x = theRealExport;
```

This exports `x: undefined` not `x: typeof theRealExport`.

#### Property access on `require` no longer imports a single property from a module:

```js
const x = require("y").x;
```

If you can't configure your package to use ESM syntax, you can use destructuring instead:

```js
const { x } = require("y");
```
