//// [tests/cases/compiler/jsDocLinkHoverPanic.ts] ////

//// [jsDocLinkHoverPanic.ts]
/**
 * A function with JSDoc links that should not cause hover to panic.
 * See {@link someFunction} for more details.
 * Also check {@linkcode anotherFunction}.
 * And {@linkplain plainLinkFunction}.
 */
function someFunction(): string {
    return "test";
}

/**
 * Another function referenced in links.
 */
function anotherFunction(): number {
    return 42;
}

/**
 * Plain link function.
 */
function plainLinkFunction(): boolean {
    return true;
}

// This should trigger hover functionality which may cause the panic
const result = someFunction();

//// [jsDocLinkHoverPanic.js]
/**
 * A function with JSDoc links that should not cause hover to panic.
 * See {@link someFunction} for more details.
 * Also check {@linkcode anotherFunction}.
 * And {@linkplain plainLinkFunction}.
 */
function someFunction() {
    return "test";
}
/**
 * Another function referenced in links.
 */
function anotherFunction() {
    return 42;
}
/**
 * Plain link function.
 */
function plainLinkFunction() {
    return true;
}
// This should trigger hover functionality which may cause the panic
const result = someFunction();
