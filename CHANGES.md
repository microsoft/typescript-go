CHANGES.md lists intentional changes between the Strada (Typescript) and Corsa (Go) compilers.

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

1. When `"strict": false`, Corsa no longer allows omitting arguments for parameters with type `undefined`, `unknown`, or `any`:


```js
/** @param {unknown} x */
function f(x) { return x; }
f(); // Previously allowed, now an error
```

`void` can still be omitted, regardless of strict mode:

```js
/** @param {void} x */
function f(x) { return x; }
f(); // Still allowed
```

2. Strada's JS-specific rules for inferring type arguments no longer apply in Corsa.

Inferred type arguments may change. For example:

```js
/** @type {any} */
var x = { a: 1, b: 2 };
var entries = Object.entries(x);
```
In Strada, `entries: Array<[string, any]>`.
In Corsa it has type `Array<[string, unknown]>`, the same as in TypeScript.

### JSDoc Types

1. JSDoc variadic types are now only synonyms for array types.

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


2. A variadic type on a parameter no longer makes it a rest parameter. The parameter must use standard rest syntax.

```js
/** @param {...number} ns */
function sum(ns) {}
```

Must now be written as

```js
/** @param {...number} ns */
function sum(...ns) {}
```

3. The postfix `=` type no longer adds `undefined` even when `strictNullChecks` is off:

```js
/** @param {number=} x */
function f(x) {
    return x;
}
```

will now have `x?: number` not `x?: number | undefined` with `strictNullChecks` off.
Regardless of strictness, it still makes parameters optional when used in a `@param` tag.


### JSDoc Tags

1. `@type` tags no longer apply to function declarations, and now contextually type functionn expressions instead of applying directly. So this annotation no longer does anything:

```js

/** @type {(x: unknown) => asserts x is string } */
function assertIsString(x) {
    if (!(typeof x === "string")) throw new Error();
}
```

Although this one still works via contextual typing:

```js
/** @typedef {(check: boolean) => asserts check} AssertFunc */

/** @type {AssertFunc} */
const assert = check => {
    if (!check) throw new Error();
}
```

A number of things change slightly because of differences between type annotation and contextual typing.

2. `asserts` annotation for an arrow function must be on the declaring variable, not on the arrow itself. This no longer works:

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

3. Error messages on async functions that incorrectly return non-Promises now use the same error as TS.

4. `@typedef` and `@callback` in a class body are no longer accessible outside the class. 
They must be moved outside the class to use them outside the class.

5. `@class` or `@constructor` does not make a function into a constructor function.

Corsa ignores `@class` and `@constructor`.
This makes a difference on a function without this-property assignments or associated prototype-function assignments.

6. `@param` tags now apply to at most one function.

If they're in a place where they could apply to multiple functions, they apply only to the first one.
If you have `"strict": true`, you will see a noImplicitAny error on the now-untyped parameters.

```js
/** @param {number} x */
var f = x => x, g = x => x;
```

7. Optional marking on parameter names now makes the parameter both optional and undefined:

```js
/** @param {number} [x] */
function f(x) {
    return x;
}
```

This behaves the same as Typescript's `x?: number` syntax. 
Strada makes the parameter optional but does not add `undefined` to the type.

8. Type assertions with `@type` tags now prevent narrowing of the type.

```js
/** @param {C | undefined} cu */
function f(cu) {
    if (/** @type {any} */ (cu).undeclaredProperty) {
        cu // still has type C | undefined
    }
}
```

In Strada, `cu` incorrectly narrows to `C` inside the `if` block, unlike with TS assertion syntax.
In Corsa, the behaviour is the same between TS and JS.

### Expandos

1. Expando assignments of `void 0` are no longer ignored as a special case:

```js
var o = {}
o.y = void 0
```

creates a property `y: undefined` on `o` (which will widen to `y: any` if strictNullChecks is off).

2. A this-property expression with a type annotation in the constructor no longer creates a property:

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

3. Assigning an object literal to the `prototype` property of a function no longer makes it a constructor function:

```js
function Foo() {}
Foo.prototype = {
    /** @param {number} x */
    bar(x) {
        return x;
    }
};
```

If you still need to use constructor functions instead of classes, you should declare methods individually on the prototype:

```js
function Foo() {}
/** @param {number} x */ 
Foo.prototype.bar = function(x) {
    return x;
};
```

Although classes are a much better way to write this code.

### CommonJS

1. Chained exports no longer work:

```js
exports.x = exports.y = 12
```

Now only exports `x`, not `y` as well.

2. Exporting `void 0` is no longer ignored as a special case:

```js
exports.x = void 0
// several lines later...
exports.x = theRealExport
```

This exports `x: undefined` not `x: typeof theRealExport`.

3. Type info for `module` shows a property with name of an instead of `exports`:

```js
module.exports = singleIdentifier
```

results in `module: { singleIdentifier: any }`

4. Property access on `require` no longer imports a single property from a module:

```js
const x = require("y").x
```

If you can't configure your package to use ESM syntax, you can use destructuring instead:

```js
const { x } = require("y")
```

5. `Object.defineProperty` on `exports` no longer creates an export:

```js
Object.defineProperty(exports, "x", { value: 12 })
```

This applies to `module.exports` as well.
Use `exports.x = 12` instead.
