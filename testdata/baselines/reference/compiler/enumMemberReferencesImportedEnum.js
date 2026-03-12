//// [tests/cases/compiler/enumMemberReferencesImportedEnum.ts] ////

//// [imported.ts]
export enum Imported {
	Option1 = 1,
	Option2 = "hello",
}

//// [usage.ts]
import { Imported } from "./imported.js";
export enum Usage {
	Option1 = Imported.Option1,
	Option2 = 2,
	Option3 = Imported.Option2,
}


//// [imported.js]
export var Imported;
(function (Imported) {
    Imported[Imported["Option1"] = 1] = "Option1";
    Imported["Option2"] = "hello";
})(Imported || (Imported = {}));
//// [usage.js]
export var Usage;
(function (Usage) {
    Usage[Usage["Option1"] = 1] = "Option1";
    Usage[Usage["Option2"] = 2] = "Option2";
    Usage["Option3"] = "hello";
})(Usage || (Usage = {}));
