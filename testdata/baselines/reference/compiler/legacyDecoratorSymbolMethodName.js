//// [tests/cases/compiler/legacyDecoratorSymbolMethodName.ts] ////

//// [legacyDecoratorSymbolMethodName.ts]
export class A {
  @(fakeDecorator as any)
  [Symbol()]() {}
}

function fakeDecorator() {}


//// [legacyDecoratorSymbolMethodName.js]
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
var _a;
Object.defineProperty(exports, "__esModule", { value: true });
exports.A = void 0;
class A {
    [_a = Symbol()]() { }
}
exports.A = A;
__decorate([
    fakeDecorator,
    __metadata("design:type", Function),
    __metadata("design:paramtypes", []),
    __metadata("design:returntype", void 0)
], A.prototype, _a, null);
function fakeDecorator() { }
