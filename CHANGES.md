# JSDoc and JS inference support
JSDoc and JS inference in Typescript supports a huge variety of users, some of which no longer exist. 
I propose dropping some features in order to simplify and speed up the compiler and language service.

When we added Javascript support to the Typescript compiler 7-8 years ago, we tried to support 4 kinds of users:

1. JS-only -- people opening a random JS file in VSCode, often by somebody who's never heard of TS.
2. TS-early-adopter -- somebody who's adding `// @ts-check` to a JSDoc-but-unchecked code base in order to convince their team of the benefits of Typescript.
3. Closure/other JSDoc processors -- an entire team writing for Google's Closure compiler, which used JSDoc type syntax, or perhaps the original JSDoc documentation generator.
4. TS-via-JSDoc -- people who wanted to use Typescript but
	- didn't want a build step.
	- didn't want to rely on a non-standard JS variant with a relatively short history.
	- didn't want to restrict their contributor pool to people who knew Typescript.
	- didn't want to rely on a Microsoft-provided compiler.

That is a wide variety of users, from people who don't even know they're using Typescript to those who are writing complete Typescript-first code, just in a different syntax. Today, however, only the two users on the extremes are left. The wild variation in JSDoc is long gone, in part because JS files are now mostly written in an editor running tsserver. And all the Closure codebases that I know of have migrated to Typescript (although it's still used as an optimiser). Even TS-via-JSDoc projects are getting rarer as Typescript syntax is supported everywhere in the JS ecosystem.

Also, Javascript itself has changed a lot since 2017. `class` has replaced constructor functions. ES modules are finally replacing CommonJS, and the syntax for both has long since been replaced for people with a build. In a word, Typescript's JS support needs cuts. Patterns with no Typescript equivalent are less important than in 2017, and there are many fewer to support. Those that remain are easier to analyse statically because of ES standards and because of Typescript's own place as the de facto interpreter of JSDoc semantics.

So what do those two kinds of customers need when editing modern JS?

JS-only users need goto-definition, find-all-references and quickinfo/signature help. They are not sensitive to semantic errors because they only see syntax errors.

TS-via-JSDoc users need compatibility with Typescript. They want to write code as strictly as in a .ts file. That means:
1. All of Typescript needs to be accessible, including new features like `satisfies`.
2. JS inference shouldn't add extra semantics.

By now, TS-via-JSDoc users also need backward compatibility with Typescript's JS support itself. I doubt many new projects are adopting pure JSDoc, so it's important for existing ones to be able to switch to tsgo when it's ready.


## Overview
Here are the highlights of what I plan to keep and plan to cut. Please comment on this proposal with your ideas, especially if you see a cut feature you'd like to see kept.

JS support falls broadly into:
1. Basic type and tag syntax. (`@param`, `@typedef`, etc) -- keep almost all of it.
2. CommonJS -- keep almost all of it.
2. Expando assignments (called "assignment declarations") -- keep most of them.
3. Constructor function support -- keep most of it.
4. Closure backward compatibility -- cut almost all of it.

Now let's get into the details.

## Keep
### Basic type and tag syntax
JSDoc specifies types with `@param`, `@return`, `@type` and `@typedef` tags, plus some others. This stays.

```js
/** @typedef {{ x: number, y: number }} Pos */
/**
 * @param {number} n
 * @param {number} m
 * @return {Pos}
 */
function polar2Cartesian(n,m) { ... }
```

In addition, most JSDoc-specific type syntax stays:

```js
/** @typedef {{ anything: *, indexSignature: Object.<string,string> }} Example */
```

### CommonJS 
CommonJS keeps `module.exports` and `require` support; most new handwritten code no longer uses it, but it's still useful so Typescript can understand semantics of a JS file inside node_modules:
 ```js
const fs = require('fs')
function f() { }
function g() { }
module.exports = { f, g }
 ```

### Expando assignments

Simple expando assignments stay, mainly to help typing of undeclared or untyped variables:
If you assign to an undeclared property of an object, it'll be treated like a declaration in a namespace of the same name:

```ts
function f(o) {
  o.p = {}
  o.q = 12
}
// is equivalent to
function f(o: any) {
  namespace o {
    export const p: {}
    export const q: number
  }
  o.p = {}
  o.q = 12
}
```

If you assign to an undeclared property of a class, it'll be treated like a property declaration. This works even outside the constructor.
```js
class C {
  constructor() {
    this.ctorUndeclared = 1
  }
  m() {
    this.methodUndeclared = this.ctorUndeclared + 1
  }
}
```

### Constructor functions
Constructor functions are rarely handwritten but still show up in generated code that targets ES5. However, generated code has much less variety than hand-written code.

```ts
function C(x) {
  this.x = x
  this.y = 1
}
C.prototype.add = function(z) {
  return this.x + this.y + z
}

// is equivalent to 
class C {
  x: any
  y: number
  constructor(x) {
    this.x = x
    this.y = 1
  }
  add(z) {
    return this.x + this.y + z
  }
}
```

## Cut

### Tags
Note that removing support for a tag will still parse it, but as `JSDocUnknownTag`.

#### `@class`
`@class` was only useful for constructor functions, and Typescript rarely needs help identifying constructor functions. Basically unused in handwritten code. 

**NOTE**: Typescript will continue to *emit* `@class` for ES5 in order to aid other tools.
```js
/** @class - is redundant because TS uses `this.x = x` as a marker. */
function C(x) {
  this.x = x
}
```
#### `@throws`
`@throws` has no semantics in TS, no special syntax and is rarely used.
```js
/** @throws {Error} - it's always Error in JS */
function div(m,n) { return m / n }
```

#### `@author`
`@author Name <name@email.com>` is supported to avoid parsing `@email` as a fresh tag and hurting readability of the resulting hover. But it's so rarely used that I'm not concerned about the loss.
```js
/** 
 * @author Finn <finn@treehouse.com> 
 * Without the author tag, TS will parse a second tag, @treehouse, by mistake
 */
function mathematical() { }
```
#### `@enum`
`@enum` is a Closure tag. See below.

### Types

#### Standalone `?`
`?` is a synonym of `any` -- but it's barely used:

```js
/** @type {?} */
var unknown
```

#### Namepaths
`module:path~Type` - Barely used, no semantics. Displaced entirely by `import("path").Type`
#### Postfix types
`T!` and `T?`, but keep `!T` and `?T` - The postfix variants came from Closure. See below. Prefix `?T` is used in Flow, which is useful for TS to understand.

#### `String`, `Number` etc synonyms
Using `String` for `string` is rare in modern code, so I decided that I'd rather have slightly wrong types in a few JS files, but less confusion among people who write TS-via-JSDoc.

```ts
/** @type {Number} */
var n = 1
// is equivalent to 
var n: number = 1

/** @type {Object} */
var o = { n }
// is equivalent to
var o: any = { n }
```
#### automatic insertion of missing `typeof` 
That is, resolving values as types. This was an aid intended for people in a `// @ts-check` world. Today it's largely unused by experts using TS-via-JSDoc, and too complex for random JS files to contain.

```js
const value = { property: 1 }
/** @type {value.property} */
var incorrect
/** @type {typeof value["property"]} */
var correct
```

#### `@type` on functions
This is a nice feature from the original JSDoc document processor, but around 80% of uses are in places where an arrow would already be contextually typed:
```js
/** @type {(x: number) => number} */
var sub1 = x => x - 1
```
And the remaining uses are rare -- about 0.2% of the most-used tag, `@param`:
```js
/** @type {FancySignature} */
function f(x,y,z) {
}
```
The workaround is easy -- switch to a variable declaration with an arrow function.
#### `function(new, string, string): T` 
This is Closure syntax. See below.

### Constructor functions
#### Multiple assignments to the same `C.prototype.property` 
Currently the types are unioned, but this code is quite rare.
The main use of this pattern is to provide polyfills or platform-specific shims, which all have the same type, so instead typescript-go will take the type of the first declaration.

```js
function C() { }
if (isWindows) {
  C.prototype.m = function() { return -1 }
} else {
  C.prototype.m = function() { return 0 }
}
```
#### Assignment of an object as the entire prototype
```js
function C(x) {
  this.x = x
}
C.prototype = {
  successor() {
    return this.x + 1
  }
  predecessor() {
    return this.x - 1
  }
  // more methods here...
}
```
### CommonJS
#### toplevel `this.p = {}` assignments
```js
this.p = 12
// inside a module is equivalent to
module.exports.p = 12
// outside a module is equivalent to
global {
  var p: number
}
```
#### Multiple assignments to the same export
Currently the type is determined using control flow.
The main use of this pattern is to provide polyfills or platform-specific shims, which all have the same type, so instead typescript-go will take the type of the first declaration.

```js
if (isWindows) {
  module.exports.normalise = function(path) {
    // replace slashes
  }
} else {
  module.exports.normalise = function(path) {
    return path
  }
}
```

The same applies to single-object exports, but I didn't observe any usage like this:
```js
if (isWindows) {
  module.exports = { ... }
} else {
  module.exports = { ... }
}
```

#### single-property access `require`
```js
var readFileSync = require('fs').readFileSync
// modern code uses destructuring:
const { readFileSync } = require('fs')
```
#### aliasing of `module.exports`:
```js
var mod = module.exports
// f is exported because `mod` aliases module.exports
mod.f = function() { ... }
```
#### ignoring preventative empty assignment of `module.exports`:
```js
module.exports = {}
module.exports.f = function() { ... }
```
Modern code assumes that `module.exports` is always defined.

### Global code
#### expando assignments skip  *x* `||` in assignments of the form *x* `=` *x* `||` *expression*
This was intended to allow TS to understand global polyfills, but nobody uses this form these days.
```js
glob = glob || function polyfill() { ... }
```
### Closure backward-compatibility support.
#### Anonymous typedefs and enums get their names from following ExpressionStatement.
```ts
/** @typedef {{ o: number }} */
var Name;
// is equivalent to
type Name = { o: number }
var Name;
```
#### ExpressionStatement property declarations 
This feature is used in Closure extern files.
```ts
var Ns = {}
/** @type {number} */
Ns.Num;

// is equivalent to
namespace Ns {
  export var Num: number
}

// also works in classes
class C {
  constructor() {
    /** @type {number} */
    this.property;
  }
}
/** @type {() => void} */
C.prototype.p;
// is equivalent to
class C {
  property: number
  declare p(): void
  constructor() { }
}
```
#### Undeclared-namespace support.
Also from Closure extern files.
```ts
var Ns = {}
Ns.Nested1.Nested2.Num = 12;

// is equivalent to
namespace Ns {
  namespace Nested1 {
	namespace Nested2 {
	   export var Num: number
	}
  }
}
// also supported in expando assignments elsewhere, eg
Ns.Nested1.Nested2.C.prototype = {
  method() { }
}
// is equivalent to
namespace Ns {
  namespace Nested1 {
	namespace Nested2 {
	   export class C {
	     method(): void 
	   }
	}
  }
}
// also works with other features:
Object.defineProperty(Ns.Nested1.Nested2.o, "r")
/** @typedef {number} */
Ns.Nested1.Nested2.T;
Ns.Nested1.Nested2.o.expando = 34
```
These features are only used in Closure extern files, which aren't useful for Typescript to understand, even if somebody does happen to open one in VSCode.
#### Postfix `T!` and `T?` 
(which were only in *Closure* for backward compatibility with an even older system!)

- `T?` is still allowed in tuples, where it already had Typescript semantics: `[number, string?]`
- `?T` (and `!T` are still allowed), since `?T` is Flow syntax as well.
#### `function(string,string):void` syntax for signatures
This is a little-used synonym for `(s: string, t: string) => void`
#### `@enum`
This is an Closure enum with Closure semantics, not Typescript semantics.
```ts
/** @enum {number} */
const Tristate = {
	True: 1,
	False: 0,
	Maybe: -1,
}
// is equivalent to
type Tristate = number
const Tristate: { [s: string]: number } = { ... }
```
