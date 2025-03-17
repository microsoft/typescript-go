//// [tests/cases/compiler/jsDeclarationsWithDefaultAsNamespaceLikeMerge.ts] ////

//// [helper.d.ts]
type Computed = () => any;
interface Mapper<R> {
    <Key extends string>(map: Key[]): { [K in Key]: R };
    <Map extends Record<string, string>>(map: Map): { [K in keyof Map]: R };
}
interface NamespacedMappers {
    mapState: Mapper<Computed>;
}
export declare function createNamespacedHelpers(): NamespacedMappers;

//// [index.js]
import { createNamespacedHelpers } from './helper'
const { mapState } = createNamespacedHelpers()
export default {
    computed: {
        ...mapState(['panels'])
    }
}

//// [index.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const helper_1 = require("./helper");
const { mapState } = (0, helper_1.createNamespacedHelpers)();
exports.default = {
    computed: {
        ...mapState(['panels'])
    }
};
