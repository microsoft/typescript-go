//// [tests/cases/compiler/plainJsJSDocModifiersInObjectLiteral.ts] ////

//// [test.js]
const obj = {
  /** @override */
  created() {
  },

  /** @private */
  onClose_() {
  },

  /** @readonly */
  get value() {
    return 1;
  },

  /** @protected */
  set value(v) {
  },
};




//// [test.d.ts]
declare const obj: {
    /** @override */
    created(): void;
    /** @private */
    onClose_(): void;
    /** @readonly */
    value: number;
};
