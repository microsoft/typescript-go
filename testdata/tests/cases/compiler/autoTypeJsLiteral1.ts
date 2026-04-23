// @allowJs: true
// @checkJs: true
// @noEmit: true

// @filename: EnvironmentPlugin.js
/** @typedef {string | number | boolean} CodeValue */

class EnvironmentPlugin {
	constructor() {
        this.keys = /** @type {string[]} */ ([]);
		this.defaultValues = {};
	}

	apply() {
		/** @type {Record<string, CodeValue>} */
		const definitions = {};
		for (const key of this.keys) {
			const value = this.defaultValues[key];
			definitions[key] = value;
		}
	}
}
