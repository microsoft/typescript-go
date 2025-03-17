//// [tests/cases/compiler/jsdocTypeCast.ts] ////

//// [jsdocTypeCast.js]
/**
 * @param {string} x
 */
 function f(x) {
    /** @type {'a' | 'b'} */
    let a = (x); // Error
    a;

    /** @type {'a' | 'b'} */
    let b = (((x))); // Error
    b;

    /** @type {'a' | 'b'} */
    let c = /** @type {'a' | 'b'} */ (x); // Ok
    c;
}


//// [jsdocTypeCast.js]
function f(x) {
    let a = (x);
    a;
    let b = (((x)));
    b;
    let c = (x);
    c;
}
