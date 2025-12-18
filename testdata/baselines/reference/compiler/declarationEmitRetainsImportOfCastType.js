//// [tests/cases/compiler/declarationEmitRetainsImportOfCastType.ts] ////

//// [declarationEmitRetainsImportOfCastType.ts]
import { WritableAtom } from 'jotai'

export function focusAtom() {
  return null as unknown as WritableAtom<any, any, any>
}


//// [declarationEmitRetainsImportOfCastType.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.focusAtom = focusAtom;
function focusAtom() {
    return null;
}


//// [declarationEmitRetainsImportOfCastType.d.ts]
import { WritableAtom } from 'jotai';
export declare function focusAtom(): WritableAtom<any, any, any>;
