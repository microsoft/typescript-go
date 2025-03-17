//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsClassStatic.ts] ////

//// [source.js]
class Handler {
	static get OPTIONS() {
		return 1;
	}

	process() {
	}
}
Handler.statische = function() { }
const Strings = {
    a: "A",
    b: "B"
}

module.exports = Handler;
module.exports.Strings = Strings

/**
 * @typedef {Object} HandlerOptions
 * @property {String} name
 * Should be able to export a type alias at the same time.
 */


//// [source.js]
class Handler {
    static get OPTIONS() {
        return 1;
    }
    process() {
    }
}
Handler.statische = function () { };
const Strings = {
    a: "A",
    b: "B"
};
module.exports = Handler;
module.exports.Strings = Strings;
