//// [tests/cases/compiler/topLevelLambda.ts] ////

//// [topLevelLambda.ts]
namespace M {
	var f = () => {this.window;}
}


//// [topLevelLambda.js]
var M;
(function (M) {
    var f = () => { this.window; };
})(M || (M = {}));
