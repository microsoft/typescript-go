// @allowJs: true
// @checkJs: true
// @outDir: ./out
// @filename: test.js

const foo = {
	/**
	 * @overload
	 * @param {string} termCode
	 * @param {string[]} crnList
	 * @param {string} sis
	 * @returns {Record<string, string>}
	 */
	/**
	 * @overload
	 * @param {string} termCode
	 * @param {string} crn
	 * @param {string} sis
	 * @returns {string}
	 */
	/**
	 * @param {string} termCode
	 * @param {string | string[]} crnList
	 * @param {string} sis
	 * @returns {string | Record<string, string>}
	 */
	getStatus(termCode, crnList, sis) {},
};
