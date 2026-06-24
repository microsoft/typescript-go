//// [tests/cases/compiler/plainJsJSDocModifiersInObjectLiteral.ts] ////

//// [test.js]
const obj = {
  /** @override */
  created() {
  },

  /** @private */
  onClose_() {
  },
};




//// [test.d.ts]
declare const obj: {
    /** @override */
    created(): void;
    /** @private */
    onClose_(): void;
};
