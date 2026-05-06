//// [tests/cases/compiler/es6ExportAll.ts] ////

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
export * from "server";

//// [server.js]
export class c {
}
export var m;
(function (m) {
    m.x = 10;
})(m || (m = {}));
export var x = 10;
//// [client.js]
export * from "server";


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
export * from "server";
