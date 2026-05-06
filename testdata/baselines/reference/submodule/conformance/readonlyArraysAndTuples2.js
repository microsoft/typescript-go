//// [tests/cases/conformance/types/tuple/readonlyArraysAndTuples2.ts] ////

//// [readonlyArraysAndTuples2.ts]
type T10 = string[];
type T11 = Array<string>;
type T12 = readonly string[];
type T13 = ReadonlyArray<string>;

type T20 = [number, number];
type T21 = readonly [number, number];

declare function f1(ma: string[], ra: readonly string[], mt: [string, string], rt: readonly [string, string]): readonly [string, string];

declare const someDec: any;

class A {
  @someDec
  j: readonly string[] = [];
  @someDec
  k: readonly [string, number] = ['foo', 42];
}


//// [readonlyArraysAndTuples2.js]
"use strict";
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};
class A {
    constructor() {
        this.j = [];
        this.k = ['foo', 42];
    }
}
__decorate([
    someDec,
    __metadata("design:type", Array)
], A.prototype, "j", void 0);
__decorate([
    someDec,
    __metadata("design:type", Array)
], A.prototype, "k", void 0);


//// [readonlyArraysAndTuples2.d.ts]
type T10 = string[];
type T11 = Array<string>;
type T12 = readonly string[];
type T13 = ReadonlyArray<string>;
type T20 = [number, number];
type T21 = readonly [number, number];
function f1(ma: string[], ra: readonly string[], mt: [string, string], rt: readonly [string, string]): readonly [string, string];
const someDec: any;
class A {
    j: readonly string[];
    k: readonly [string, number];
}


//// [DtsFileErrors]


readonlyArraysAndTuples2.d.ts(7,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== readonlyArraysAndTuples2.d.ts (1 errors) ====
    type T10 = string[];
    type T11 = Array<string>;
    type T12 = readonly string[];
    type T13 = ReadonlyArray<string>;
    type T20 = [number, number];
    type T21 = readonly [number, number];
    function f1(ma: string[], ra: readonly string[], mt: [string, string], rt: readonly [string, string]): readonly [string, string];
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    const someDec: any;
    class A {
        j: readonly string[];
        k: readonly [string, number];
    }
    