//// [tests/cases/conformance/salsa/expandoOnAlias.ts] ////

//// [vue.js]
export class Vue {}
export const config = { x: 0 };

//// [test.js]
import { Vue, config } from "./vue";

// Expando declarations aren't allowed on aliases.
Vue.config = {};
new Vue();

// This is not an expando declaration; it's just a plain property assignment.
config.x = 1;

// This is not an expando declaration; it works because non-strict JS allows
// loosey goosey assignment on objects.
config.y = {};
config.x;
config.y;


//// [vue.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.config = exports.Vue = void 0;
class Vue {
}
exports.Vue = Vue;
exports.config = { x: 0 };
//// [test.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const vue_1 = require("./vue");
vue_1.Vue.config = {};
new vue_1.Vue();
vue_1.config.x = 1;
vue_1.config.y = {};
vue_1.config.x;
vue_1.config.y;
