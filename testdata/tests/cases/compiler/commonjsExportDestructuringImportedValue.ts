// @target: es2022
// @module: commonjs
// @noTypesAndSymbols: true

// @filename: /enum.ts
export class CodePriceType {
    static A = "a";
    static B = "b";
}

// @filename: /repro.ts
import { CodePriceType } from "./enum";
export const { A, B } = CodePriceType;
