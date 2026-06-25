const { version } = require("../package.json");
exports.version = version;
exports.versionMajorMinor = version.split(".").slice(0, 2).join(".");
