//// [tests/cases/conformance/externalModules/moduleResolutionWithoutExtension3.ts] ////

//// [foo.tsx]
export function foo() {
    return "";
}

//// [bar.mts]
// Extensionless relative path ES import in an ES module
import { foo } from "./foo"; // should error, suggest adding ".jsx"


//// [bar.mjs]
export {};
//// [foo.jsx]
export function foo() {
    return "";
}
