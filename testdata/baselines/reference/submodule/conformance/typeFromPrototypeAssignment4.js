//// [tests/cases/conformance/salsa/typeFromPrototypeAssignment4.ts] ////

//// [a.js]
function Multimap4() {
  this._map = {};
};

Multimap4["prototype"] = {
  /**
   * @param {string} key
   * @returns {number} the value ok
   */
  get(key) {
    return this._map[key + ''];
  }
};

Multimap4["prototype"]["add-on"] = function() {};
Multimap4["prototype"]["addon"] = function() {};
Multimap4["prototype"]["__underscores__"] = function() {};

const map4 = new Multimap4();
map4.get("");
map4["add-on"]();
map4.addon();
map4.__underscores__();




//// [a.d.ts]
function Multimap4(): void;
declare namespace Multimap4 {
    var prototype: {
        /**
         * @param {string} key
         * @returns {number} the value ok
         */
        get(key: string): number;
    };
}
const map4: any;


//// [DtsFileErrors]


out/a.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== out/a.d.ts (1 errors) ====
    function Multimap4(): void;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    declare namespace Multimap4 {
        var prototype: {
            /**
             * @param {string} key
             * @returns {number} the value ok
             */
            get(key: string): number;
        };
    }
    const map4: any;
    