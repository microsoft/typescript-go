//// [tests/cases/compiler/jsDeclarationsGlobalFileConstFunctionNamed.ts] ////

//// [file.js]
const SomeConstructor = function Named() {
	this.x = 1;
};

const SomeConstructor2 = function Named() {
};
SomeConstructor2.staticMember = "str";

const SomeConstructor3 = function Named() {
	this.x = 1;
};
SomeConstructor3.staticMember = "str";

const SelfReference = function Named() {
    if (!(this instanceof Named)) return new Named();
    this.x = 1;
}
SelfReference.staticMember = "str";




//// [file.d.ts]
const SomeConstructor: () => void;
function SomeConstructor2(): void;
declare namespace SomeConstructor2 {
    var staticMember: string;
}
function SomeConstructor3(): void;
declare namespace SomeConstructor3 {
    var staticMember: string;
}
function SelfReference(): any;
declare namespace SelfReference {
    var staticMember: string;
}
