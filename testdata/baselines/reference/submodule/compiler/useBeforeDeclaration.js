//// [tests/cases/compiler/useBeforeDeclaration.ts] ////

//// [A.ts]
namespace ts {
    export function printVersion():void {
        log("Version: " + sys.version);  // the call of sys.version is deferred, should not report an error.
    }

    export function log(info:string):void {

    }
}

//// [B.ts]
namespace ts {

    export let sys:{version:string} = {version: "2.0.5"};

}



//// [A.js]
var ts;
(function (ts) {
    function printVersion() {
        log("Version: " + sys.version);
    }
    ts.printVersion = printVersion;
    function log(info) {
    }
    ts.log = log;
})(ts || (ts = {}));
//// [B.js]
var ts;
(function (ts) {
    ts.sys = { version: "2.0.5" };
})(ts || (ts = {}));
