//// [tests/cases/conformance/jsdoc/jsdocImplements_interface_multiple.ts] ////

//// [defs.d.ts]
interface Drawable {
    draw(): number;
}
interface Sizable {
    size(): number;
}
//// [a.js]
/** 
 * @implements {Drawable} 
 * @implements Sizable 
 **/
class Square {
    draw() {
        return 0;
    }
    size() {
        return 0;
    }
}
/**
 * @implements Drawable
 * @implements {Sizable}
 **/
class BadSquare {
    size() {
        return 0;
    }
}

//// [a.js]
class Square {
    draw() {
        return 0;
    }
    size() {
        return 0;
    }
}
class BadSquare {
    size() {
        return 0;
    }
}
