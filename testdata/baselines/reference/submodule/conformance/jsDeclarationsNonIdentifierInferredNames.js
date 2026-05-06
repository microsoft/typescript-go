//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsNonIdentifierInferredNames.ts] ////

//// [jsDeclarationsNonIdentifierInferredNames.jsx]
/// <reference path="/.lib/react16.d.ts" />
import * as React from "react";
const dynPropName = "data-dyn";
export const ExampleFunctionalComponent = ({ "data-testid": dataTestId, [dynPropName]: dynProp }) => (
    <>Hello</>
);

//// [jsDeclarationsNonIdentifierInferredNames.js]
/// <reference path="/.lib/react16.d.ts" />
import * as React from "react";
const dynPropName = "data-dyn";
export const ExampleFunctionalComponent = ({ "data-testid": dataTestId, [dynPropName]: dynProp }) => (React.createElement(React.Fragment, null, "Hello"));


//// [jsDeclarationsNonIdentifierInferredNames.d.ts]
const dynPropName = "data-dyn";
export const ExampleFunctionalComponent: ({ "data-testid": dataTestId, [dynPropName]: dynProp }: {
    "data-dyn": any;
    "data-testid": any;
}) => JSX.Element;
export {};


//// [DtsFileErrors]


out/jsDeclarationsNonIdentifierInferredNames.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
out/jsDeclarationsNonIdentifierInferredNames.d.ts(5,7): error TS2503: Cannot find namespace 'JSX'.


==== out/jsDeclarationsNonIdentifierInferredNames.d.ts (2 errors) ====
    const dynPropName = "data-dyn";
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    export const ExampleFunctionalComponent: ({ "data-testid": dataTestId, [dynPropName]: dynProp }: {
        "data-dyn": any;
        "data-testid": any;
    }) => JSX.Element;
          ~~~
!!! error TS2503: Cannot find namespace 'JSX'.
    export {};
    