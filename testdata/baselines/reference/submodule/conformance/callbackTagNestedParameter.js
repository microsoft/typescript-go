//// [tests/cases/conformance/jsdoc/callbackTagNestedParameter.ts] ////

//// [cb_nested.js]
/**
 * @callback WorksWithPeopleCallback
 * @param {Object} person
 * @param {string} person.name
 * @param {number} [person.age]
 * @returns {void}
 */

/**
 * For each person, calls your callback.
 * @param {WorksWithPeopleCallback} callback
 * @returns {void}
 */
function eachPerson(callback) {
    callback({ name: "Empty" });
}




//// [cb_nested.d.ts]
/**
 * @callback WorksWithPeopleCallback
 * @param {Object} person
 * @param {string} person.name
 * @param {number} [person.age]
 * @returns {void}
 */
type WorksWithPeopleCallback = (person: {
    name: string;
    age?: number;
}) => void;
/**
 * For each person, calls your callback.
 * @param {WorksWithPeopleCallback} callback
 * @returns {void}
 */
function eachPerson(callback: WorksWithPeopleCallback): void;


//// [DtsFileErrors]


cb_nested.d.ts(17,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== cb_nested.d.ts (1 errors) ====
    /**
     * @callback WorksWithPeopleCallback
     * @param {Object} person
     * @param {string} person.name
     * @param {number} [person.age]
     * @returns {void}
     */
    type WorksWithPeopleCallback = (person: {
        name: string;
        age?: number;
    }) => void;
    /**
     * For each person, calls your callback.
     * @param {WorksWithPeopleCallback} callback
     * @returns {void}
     */
    function eachPerson(callback: WorksWithPeopleCallback): void;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    