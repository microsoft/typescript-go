//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsCommonjsRelativePath.ts] ////

//// [thing.js]
'use strict';
class Thing {}
module.exports = { Thing }

//// [reexport.js]
'use strict';
const Thing = require('./thing').Thing
module.exports = { Thing }


//// [reexport.js]
'use strict';
const Thing = require('./thing').Thing;
module.exports = { Thing };
