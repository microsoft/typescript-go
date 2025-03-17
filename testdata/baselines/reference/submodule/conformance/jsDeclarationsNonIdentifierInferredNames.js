//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsNonIdentifierInferredNames.ts] ////

//// [jsDeclarationsNonIdentifierInferredNames.jsx]
/// <reference path="/.lib/react16.d.ts" />
import * as React from "react";
const dynPropName = "data-dyn";
export const ExampleFunctionalComponent = ({ "data-testid": dataTestId, [dynPropName]: dynProp }) => (
    <>Hello</>
);

//// [jsDeclarationsNonIdentifierInferredNames.js]
import * as React from "react";
const dynPropName = "data-dyn";
export const ExampleFunctionalComponent = ({ "data-testid": dataTestId, [dynPropName]: dynProp }) => (<>Hello</>);
