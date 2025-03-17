//// [tests/cases/conformance/types/uniqueSymbol/uniqueSymbolsDeclarationsInJsErrors.ts] ////

//// [uniqueSymbolsDeclarationsInJsErrors.js]
class C {
    /**
     * @type {unique symbol}
     */
    static readwriteStaticType;
    /**
     * @type {unique symbol}
     * @readonly
     */
    static readonlyType;
    /**
     * @type {unique symbol}
     */
    static readwriteType;
}


//// [uniqueSymbolsDeclarationsInJsErrors.js]
class C {
    static readwriteStaticType;
    static readonlyType;
    static readwriteType;
}
