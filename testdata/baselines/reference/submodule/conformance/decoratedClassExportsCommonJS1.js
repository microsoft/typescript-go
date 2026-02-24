//// [tests/cases/conformance/decorators/class/decoratedClassExportsCommonJS1.ts] ////

//// [a.ts]
declare function forwardRef(x: any): any;
declare var Something: any;
@Something({ v: () => Testing123 })
export class Testing123 {
    static prop0: string;
    static prop1 = Testing123.prop0;
}

//// [a.js]
"use strict";
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var _a;
Object.defineProperty(exports, "__esModule", { value: true });
exports.Testing123 = void 0;
let Testing123 = (_a = class Testing123 {
    },
    _a.prop1 = _a.prop0,
    _a);
exports.Testing123 = Testing123;
exports.Testing123 = Testing123 = __decorate([
    Something({ v: () => Testing123 })
], Testing123);
