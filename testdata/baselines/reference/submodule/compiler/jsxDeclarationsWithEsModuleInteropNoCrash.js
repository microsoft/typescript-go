//// [tests/cases/compiler/jsxDeclarationsWithEsModuleInteropNoCrash.tsx] ////

//// [jsxDeclarationsWithEsModuleInteropNoCrash.jsx]
/// <reference path="/.lib/react16.d.ts" preserve="true" />
import PropTypes from 'prop-types';
import React from 'react';

const propTypes = {
  bar: PropTypes.bool,
};

const defaultProps = {
  bar: false,
};

function Foo({ bar }) {
  return <div>{bar}</div>;
}

Foo.propTypes = propTypes;
Foo.defaultProps = defaultProps;

export default Foo;

//// [jsxDeclarationsWithEsModuleInteropNoCrash.js]
"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const prop_types_1 = __importDefault(require("prop-types"));
const react_1 = __importDefault(require("react"));
const propTypes = {
    bar: prop_types_1.default.bool,
};
const defaultProps = {
    bar: false,
};
function Foo({ bar }) {
    return <div>{bar}</div>;
}
Foo.propTypes = propTypes;
Foo.defaultProps = defaultProps;
exports.default = Foo;
