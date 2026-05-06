//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsReactComponents.ts] ////

//// [jsDeclarationsReactComponents1.jsx]
/// <reference path="/.lib/react16.d.ts" preserve="true" />
import React from "react";
import PropTypes from "prop-types"

const TabbedShowLayout = ({
}) => {
    return (
        <div />
    );
};

TabbedShowLayout.propTypes = {
    version: PropTypes.number,

};

TabbedShowLayout.defaultProps = {
    tabs: undefined
};

export default TabbedShowLayout;

//// [jsDeclarationsReactComponents2.jsx]
import React from "react";
/**
 * @type {React.SFC}
 */
const TabbedShowLayout = () => {
    return (
        <div className="" key="">
            ok
        </div>
    );
};

TabbedShowLayout.defaultProps = {
    tabs: "default value"
};

export default TabbedShowLayout;

//// [jsDeclarationsReactComponents3.jsx]
import React from "react";
/**
 * @type {{defaultProps: {tabs: string}} & ((props?: {elem: string}) => JSX.Element)}
 */
const TabbedShowLayout = () => {
    return (
        <div className="" key="">
            ok
        </div>
    );
};

TabbedShowLayout.defaultProps = {
    tabs: "default value"
};

export default TabbedShowLayout;

//// [jsDeclarationsReactComponents4.jsx]
import React from "react";
const TabbedShowLayout = (/** @type {{className: string}}*/prop) => {
    return (
        <div className={prop.className} key="">
            ok
        </div>
    );
};

TabbedShowLayout.defaultProps = {
    tabs: "default value"
};

export default TabbedShowLayout;
//// [jsDeclarationsReactComponents5.jsx]
import React from 'react';
import PropTypes from 'prop-types';

function Tree({ allowDropOnRoot }) {
  return <div />
}

Tree.propTypes = {
    classes: PropTypes.object,
};

Tree.defaultProps = {
    classes: {},
    parentSource: 'parent_id',
};

export default Tree;

//// [jsDeclarationsReactComponents1.js]
"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
/// <reference path="/.lib/react16.d.ts" preserve="true" />
const react_1 = __importDefault(require("react"));
const prop_types_1 = __importDefault(require("prop-types"));
const TabbedShowLayout = ({}) => {
    return (react_1.default.createElement("div", null));
};
TabbedShowLayout.propTypes = {
    version: prop_types_1.default.number,
};
TabbedShowLayout.defaultProps = {
    tabs: undefined
};
exports.default = TabbedShowLayout;
//// [jsDeclarationsReactComponents2.js]
"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const react_1 = __importDefault(require("react"));
/**
 * @type {React.SFC}
 */
const TabbedShowLayout = () => {
    return (react_1.default.createElement("div", { className: "", key: "" }, "ok"));
};
TabbedShowLayout.defaultProps = {
    tabs: "default value"
};
exports.default = TabbedShowLayout;
//// [jsDeclarationsReactComponents3.js]
"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const react_1 = __importDefault(require("react"));
/**
 * @type {{defaultProps: {tabs: string}} & ((props?: {elem: string}) => JSX.Element)}
 */
const TabbedShowLayout = () => {
    return (react_1.default.createElement("div", { className: "", key: "" }, "ok"));
};
TabbedShowLayout.defaultProps = {
    tabs: "default value"
};
exports.default = TabbedShowLayout;
//// [jsDeclarationsReactComponents4.js]
"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const react_1 = __importDefault(require("react"));
const TabbedShowLayout = (/** @type {{className: string}}*/ prop) => {
    return (react_1.default.createElement("div", { className: prop.className, key: "" }, "ok"));
};
TabbedShowLayout.defaultProps = {
    tabs: "default value"
};
exports.default = TabbedShowLayout;
//// [jsDeclarationsReactComponents5.js]
"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const react_1 = __importDefault(require("react"));
const prop_types_1 = __importDefault(require("prop-types"));
function Tree({ allowDropOnRoot }) {
    return react_1.default.createElement("div", null);
}
Tree.propTypes = {
    classes: prop_types_1.default.object,
};
Tree.defaultProps = {
    classes: {},
    parentSource: 'parent_id',
};
exports.default = Tree;


//// [jsDeclarationsReactComponents1.d.ts]
/// <reference path="../../.lib/react16.d.ts" preserve="true" />
import PropTypes from "prop-types";
function TabbedShowLayout({}: {}): JSX.Element;
declare namespace TabbedShowLayout {
    var propTypes: {
        version: PropTypes.Requireable<number>;
    };
}
declare namespace TabbedShowLayout {
    var defaultProps: {
        tabs: undefined;
    };
}
export default TabbedShowLayout;
//// [jsDeclarationsReactComponents2.d.ts]
import React from "react";
/**
 * @type {React.SFC}
 */
const TabbedShowLayout: React.SFC;
export default TabbedShowLayout;
//// [jsDeclarationsReactComponents3.d.ts]
/**
 * @type {{defaultProps: {tabs: string}} & ((props?: {elem: string}) => JSX.Element)}
 */
const TabbedShowLayout: {
    defaultProps: {
        tabs: string;
    };
} & ((props?: {
    elem: string;
}) => JSX.Element);
export default TabbedShowLayout;
//// [jsDeclarationsReactComponents4.d.ts]
function TabbedShowLayout(/** @type {{className: string}}*/ prop: {
    className: string;
}): JSX.Element;
declare namespace TabbedShowLayout {
    var defaultProps: {
        tabs: string;
    };
}
export default TabbedShowLayout;
//// [jsDeclarationsReactComponents5.d.ts]
import PropTypes from 'prop-types';
function Tree({ allowDropOnRoot }: {
    allowDropOnRoot: any;
}): JSX.Element;
declare namespace Tree {
    var propTypes: {
        classes: PropTypes.Requireable<object>;
    };
}
declare namespace Tree {
    var defaultProps: {
        classes: {};
        parentSource: string;
    };
}
export default Tree;


//// [DtsFileErrors]


out/jsDeclarationsReactComponents1.d.ts(3,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
out/jsDeclarationsReactComponents2.d.ts(5,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
out/jsDeclarationsReactComponents3.d.ts(4,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
out/jsDeclarationsReactComponents4.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
out/jsDeclarationsReactComponents5.d.ts(2,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== out/jsDeclarationsReactComponents1.d.ts (1 errors) ====
    /// <reference path="../../.lib/react16.d.ts" preserve="true" />
    import PropTypes from "prop-types";
    function TabbedShowLayout({}: {}): JSX.Element;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    declare namespace TabbedShowLayout {
        var propTypes: {
            version: PropTypes.Requireable<number>;
        };
    }
    declare namespace TabbedShowLayout {
        var defaultProps: {
            tabs: undefined;
        };
    }
    export default TabbedShowLayout;
    
==== out/jsDeclarationsReactComponents2.d.ts (1 errors) ====
    import React from "react";
    /**
     * @type {React.SFC}
     */
    const TabbedShowLayout: React.SFC;
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    export default TabbedShowLayout;
    
==== out/jsDeclarationsReactComponents3.d.ts (1 errors) ====
    /**
     * @type {{defaultProps: {tabs: string}} & ((props?: {elem: string}) => JSX.Element)}
     */
    const TabbedShowLayout: {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        defaultProps: {
            tabs: string;
        };
    } & ((props?: {
        elem: string;
    }) => JSX.Element);
    export default TabbedShowLayout;
    
==== out/jsDeclarationsReactComponents4.d.ts (1 errors) ====
    function TabbedShowLayout(/** @type {{className: string}}*/ prop: {
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        className: string;
    }): JSX.Element;
    declare namespace TabbedShowLayout {
        var defaultProps: {
            tabs: string;
        };
    }
    export default TabbedShowLayout;
    
==== out/jsDeclarationsReactComponents5.d.ts (1 errors) ====
    import PropTypes from 'prop-types';
    function Tree({ allowDropOnRoot }: {
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        allowDropOnRoot: any;
    }): JSX.Element;
    declare namespace Tree {
        var propTypes: {
            classes: PropTypes.Requireable<object>;
        };
    }
    declare namespace Tree {
        var defaultProps: {
            classes: {};
            parentSource: string;
        };
    }
    export default Tree;
    