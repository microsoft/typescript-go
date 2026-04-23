// @allowJs: true
// @checkJs: true
// @noEmit: true

// @filename: EnvironmentPlugin.js
/** @typedef {string | number | boolean} CodeValue */

class EnvironmentPlugin {
	/**
	 * @param {(string | string[] | Record<string, CodeValue>)[]} keys
	 */
	constructor(...keys) {
		if (keys.length === 1 && Array.isArray(keys[0])) {
			/** @type {string[]} */
			this.keys = keys[0];
			this.defaultValues = {};
		} else if (keys.length === 1 && keys[0] && typeof keys[0] === "object") {
			this.keys = Object.keys(keys[0]);
			this.defaultValues =
				/** @type {Record<string, CodeValue>} */
				(keys[0]);
		} else {
			this.keys = /** @type {string[]} */ (keys);
			this.defaultValues = {};
		}
	}

	apply() {
		/** @type {Record<string, CodeValue>} */
		const definitions = {};
		for (const key of this.keys) {
			this.defaultValues // current: Record<string, CodeValue>, broken: {}
			const value = this.defaultValues[key];
			definitions[key] = value;
		}
	}
}
