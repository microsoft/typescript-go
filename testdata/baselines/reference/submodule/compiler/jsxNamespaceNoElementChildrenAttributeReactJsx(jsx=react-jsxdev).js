//// [tests/cases/compiler/jsxNamespaceNoElementChildrenAttributeReactJsx.tsx] ////

//// [jsx.d.ts]
declare namespace JSX {
  interface IntrinsicElements {
    h1: { children: string }
  }

  type Element = string;
}

//// [test.tsx]
const Title = (props: { children: string }) => <h1>{props.children}</h1>;

const element = <Title>Hello, world!</Title>;

//// [jsx-runtime.ts]
export {};
//// [jsx-dev-runtime.ts]
export {};

//// [jsx-dev-runtime.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
//// [test.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const jsx_dev_runtime_1 = require("/jsx/jsx-dev-runtime");
const _jsxFileName = "/test.tsx";
const Title = (props) => jsx_dev_runtime_1.jsxDEV("h1", { children: props.children }, void 0, false, { fileName: _jsxFileName, lineNumber: 1, columnNumber: 47 }, this);
const element = jsx_dev_runtime_1.jsxDEV(Title, { children: "Hello, world!" }, void 0, false, { fileName: _jsxFileName, lineNumber: 3, columnNumber: 16 }, this);
//// [jsx-runtime.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
