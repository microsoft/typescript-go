//// [tests/cases/conformance/classes/members/privateNames/privateNamesIncompatibleModifiersJs.ts] ////

//// [privateNamesIncompatibleModifiersJs.js]
class A {
    /**
     * @public
     */
    #a = 1;

    /**
     * @private
     */
    #b = 1;

    /**
     * @protected
     */
    #c = 1;

    /**
     * @public
     */
    #aMethod() { return 1; }

    /**
     * @private
     */
    #bMethod() { return 1; }

    /**
     * @protected
     */
    #cMethod() { return 1; }

    /**
     * @public
     */
    get #aProp() { return 1; }
    /**
     * @public
     */
    set #aProp(value) { }

    /**
     * @private
     */
    get #bProp() { return 1; }
    /**
     * @private
     */
    set #bProp(value) { }

    /**
    * @protected
    */
    get #cProp() { return 1; }
    /**
     * @protected
     */
    set #cProp(value) { }
}


//// [privateNamesIncompatibleModifiersJs.js]
"use strict";
var _A_a, _A_b, _A_c;
class A {
    constructor() {
        /**
         * @public
         */
        _A_a.set(this, 1);
        /**
         * @private
         */
        _A_b.set(this, 1);
        /**
         * @protected
         */
        _A_c.set(this, 1);
    }
    /**
     * @public
     */
    #aMethod() { return 1; }
    /**
     * @private
     */
    #bMethod() { return 1; }
    /**
     * @protected
     */
    #cMethod() { return 1; }
    /**
     * @public
     */
    get #aProp() { return 1; }
    /**
     * @public
     */
    set #aProp(value) { }
    /**
     * @private
     */
    get #bProp() { return 1; }
    /**
     * @private
     */
    set #bProp(value) { }
    /**
    * @protected
    */
    get #cProp() { return 1; }
    /**
     * @protected
     */
    set #cProp(value) { }
}
_A_a = new WeakMap(), _A_b = new WeakMap(), _A_c = new WeakMap();
