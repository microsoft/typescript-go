//// [tests/cases/compiler/es6ExportClauseWithoutModuleSpecifier.ts] ////

//// [server.ts]
export class c {
}
export interface i {
}
export namespace m {
    export var x = 10;
}
export var x = 10;
export namespace uninstantiated {
}

//// [client.ts]
export { c } from "./server";
export { c as c2 } from "./server";
export { i, m as instantiatedModule } from "./server";
export { uninstantiated } from "./server";
export { x } from "./server";

//// [server.js]
export class c {
}
export var m;
(function (m) {
    m.x = 10;
})(m || (m = {}));
export var x = 10;
//// [client.js]
export { c } from "./server";
export { c as c2 } from "./server";
export { m as instantiatedModule } from "./server";
export { x } from "./server";


//// [server.d.ts]
export class c {
}
export interface i {
}
export namespace m {
    var x: number;
}
export var x: number;
export namespace uninstantiated {
}
//// [client.d.ts]
export { c } from "./server";
export { c as c2 } from "./server";
export { i, m as instantiatedModule } from "./server";
export { uninstantiated } from "./server";
export { x } from "./server";
