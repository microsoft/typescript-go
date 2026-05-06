//// [tests/cases/compiler/reactTransitiveImportHasValidDeclaration.ts] ////

//// [index.d.ts]
declare namespace React {
    export interface DetailedHTMLProps<T, U> {}
    export interface HTMLAttributes<T> {}
}
export = React;
export as namespace React;
//// [index.d.ts]
/// <reference types="react" />
declare module 'react' { // augment
    interface HTMLAttributes<T> {
        css?: unknown;
    }
}
export interface StyledOtherComponentList {
    "div": React.DetailedHTMLProps<React.HTMLAttributes<HTMLDivElement>, HTMLDivElement>
}
export interface StyledOtherComponent<A, B, C> {}

//// [index.d.ts]
export * from "./types/react";

//// [index.d.ts]
import {StyledOtherComponent, StyledOtherComponentList} from "create-emotion-styled";
export default function styled(tag: string): (o: object) => StyledOtherComponent<{}, StyledOtherComponentList["div"], any>;

//// [index.ts]
import styled from "react-emotion"

const Form = styled('div')({ color: "red" })

export default Form


//// [index.js]
"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const react_emotion_1 = __importDefault(require("react-emotion"));
const Form = (0, react_emotion_1.default)('div')({ color: "red" });
exports.default = Form;


//// [index.d.ts]
const Form: import("create-emotion-styled").StyledOtherComponent<{}, import("react").DetailedHTMLProps<import("react").HTMLAttributes<HTMLDivElement>, HTMLDivElement>, any>;
export default Form;


//// [DtsFileErrors]


index.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== node_modules/react/index.d.ts (0 errors) ====
    declare namespace React {
        export interface DetailedHTMLProps<T, U> {}
        export interface HTMLAttributes<T> {}
    }
    export = React;
    export as namespace React;
==== node_modules/create-emotion-styled/types/react/index.d.ts (0 errors) ====
    /// <reference types="react" />
    declare module 'react' { // augment
        interface HTMLAttributes<T> {
            css?: unknown;
        }
    }
    export interface StyledOtherComponentList {
        "div": React.DetailedHTMLProps<React.HTMLAttributes<HTMLDivElement>, HTMLDivElement>
    }
    export interface StyledOtherComponent<A, B, C> {}
    
==== node_modules/create-emotion-styled/index.d.ts (0 errors) ====
    export * from "./types/react";
    
==== node_modules/react-emotion/index.d.ts (0 errors) ====
    import {StyledOtherComponent, StyledOtherComponentList} from "create-emotion-styled";
    export default function styled(tag: string): (o: object) => StyledOtherComponent<{}, StyledOtherComponentList["div"], any>;
    
==== index.d.ts (1 errors) ====
    const Form: import("create-emotion-styled").StyledOtherComponent<{}, import("react").DetailedHTMLProps<import("react").HTMLAttributes<HTMLDivElement>, HTMLDivElement>, any>;
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    export default Form;
    