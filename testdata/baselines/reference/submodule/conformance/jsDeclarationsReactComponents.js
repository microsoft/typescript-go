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
/// <reference path="react16.d.ts" preserve="true" />
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
/// <reference path="../..react16.d.ts" preserve="true" />
import PropTypes from "prop-types";
declare const TabbedShowLayout: {
    ({}: {}): JSX.Element;
    propTypes: {
        version: PropTypes.Requireable<number>;
    };
    defaultProps: {
        tabs: undefined;
    };
};
declare namespace TabbedShowLayout {
    const propTypes: {
        version: PropTypes.Requireable<number>;
    };
}
declare namespace TabbedShowLayout {
    const defaultProps: {
        tabs: undefined;
    };
}
export default TabbedShowLayout;
//// [jsDeclarationsReactComponents2.d.ts]
import React from "react";
/**
 * @type {React.SFC}
 */
declare const TabbedShowLayout: React.SFC;
declare namespace TabbedShowLayout {
    const defaultProps: Partial<{}> | undefined;
}
export default TabbedShowLayout;
//// [jsDeclarationsReactComponents3.d.ts]
/**
 * @type {{defaultProps: {tabs: string}} & ((props?: {elem: string}) => JSX.Element)}
 */
declare const TabbedShowLayout: {
    defaultProps: {
        tabs: string;
    };
} & ((props?: {
    elem: string;
}) => JSX.Element);
declare namespace TabbedShowLayout {
    const defaultProps: {
        tabs: string;
    };
}
export default TabbedShowLayout;
//// [jsDeclarationsReactComponents4.d.ts]
declare const TabbedShowLayout: {
    (prop: {
        className: string;
    }): JSX.Element;
    defaultProps: {
        tabs: string;
    };
};
declare namespace TabbedShowLayout {
    const defaultProps: {
        tabs: string;
    };
}
export default TabbedShowLayout;
//// [jsDeclarationsReactComponents5.d.ts]
import PropTypes from 'prop-types';
declare function Tree({ allowDropOnRoot }: {
    allowDropOnRoot: any;
}): JSX.Element;
declare namespace Tree {
    const propTypes: {
        classes: PropTypes.Requireable<object>;
    };
}
declare namespace Tree {
    const defaultProps: {
        classes: {};
        parentSource: string;
    };
}
export default Tree;


//// [DtsFileErrors]


out/jsDeclarationsReactComponents1.d.ts(3,15): error TS2451: Cannot redeclare block-scoped variable 'TabbedShowLayout'.
out/jsDeclarationsReactComponents1.d.ts(12,19): error TS2451: Cannot redeclare block-scoped variable 'TabbedShowLayout'.
out/jsDeclarationsReactComponents1.d.ts(17,19): error TS2451: Cannot redeclare block-scoped variable 'TabbedShowLayout'.
out/jsDeclarationsReactComponents2.d.ts(5,15): error TS2451: Cannot redeclare block-scoped variable 'TabbedShowLayout'.
out/jsDeclarationsReactComponents2.d.ts(6,19): error TS2451: Cannot redeclare block-scoped variable 'TabbedShowLayout'.
out/jsDeclarationsReactComponents3.d.ts(4,15): error TS2451: Cannot redeclare block-scoped variable 'TabbedShowLayout'.
out/jsDeclarationsReactComponents3.d.ts(11,19): error TS2451: Cannot redeclare block-scoped variable 'TabbedShowLayout'.
out/jsDeclarationsReactComponents4.d.ts(1,15): error TS2451: Cannot redeclare block-scoped variable 'TabbedShowLayout'.
out/jsDeclarationsReactComponents4.d.ts(9,19): error TS2451: Cannot redeclare block-scoped variable 'TabbedShowLayout'.


==== out/jsDeclarationsReactComponents1.d.ts (3 errors) ====
    /// <reference path="../../.lib/react16.d.ts" preserve="true" />
    import PropTypes from "prop-types";
    declare const TabbedShowLayout: {
                  ~~~~~~~~~~~~~~~~
!!! error TS2451: Cannot redeclare block-scoped variable 'TabbedShowLayout'.
        ({}: {}): JSX.Element;
        propTypes: {
            version: PropTypes.Requireable<number>;
        };
        defaultProps: {
            tabs: undefined;
        };
    };
    declare namespace TabbedShowLayout {
                      ~~~~~~~~~~~~~~~~
!!! error TS2451: Cannot redeclare block-scoped variable 'TabbedShowLayout'.
        const propTypes: {
            version: PropTypes.Requireable<number>;
        };
    }
    declare namespace TabbedShowLayout {
                      ~~~~~~~~~~~~~~~~
!!! error TS2451: Cannot redeclare block-scoped variable 'TabbedShowLayout'.
        const defaultProps: {
            tabs: undefined;
        };
    }
    export default TabbedShowLayout;
    
==== out/jsDeclarationsReactComponents2.d.ts (2 errors) ====
    import React from "react";
    /**
     * @type {React.SFC}
     */
    declare const TabbedShowLayout: React.SFC;
                  ~~~~~~~~~~~~~~~~
!!! error TS2451: Cannot redeclare block-scoped variable 'TabbedShowLayout'.
    declare namespace TabbedShowLayout {
                      ~~~~~~~~~~~~~~~~
!!! error TS2451: Cannot redeclare block-scoped variable 'TabbedShowLayout'.
        const defaultProps: Partial<{}> | undefined;
    }
    export default TabbedShowLayout;
    
==== out/jsDeclarationsReactComponents3.d.ts (2 errors) ====
    /**
     * @type {{defaultProps: {tabs: string}} & ((props?: {elem: string}) => JSX.Element)}
     */
    declare const TabbedShowLayout: {
                  ~~~~~~~~~~~~~~~~
!!! error TS2451: Cannot redeclare block-scoped variable 'TabbedShowLayout'.
        defaultProps: {
            tabs: string;
        };
    } & ((props?: {
        elem: string;
    }) => JSX.Element);
    declare namespace TabbedShowLayout {
                      ~~~~~~~~~~~~~~~~
!!! error TS2451: Cannot redeclare block-scoped variable 'TabbedShowLayout'.
        const defaultProps: {
            tabs: string;
        };
    }
    export default TabbedShowLayout;
    
==== out/jsDeclarationsReactComponents4.d.ts (2 errors) ====
    declare const TabbedShowLayout: {
                  ~~~~~~~~~~~~~~~~
!!! error TS2451: Cannot redeclare block-scoped variable 'TabbedShowLayout'.
        (prop: {
            className: string;
        }): JSX.Element;
        defaultProps: {
            tabs: string;
        };
    };
    declare namespace TabbedShowLayout {
                      ~~~~~~~~~~~~~~~~
!!! error TS2451: Cannot redeclare block-scoped variable 'TabbedShowLayout'.
        const defaultProps: {
            tabs: string;
        };
    }
    export default TabbedShowLayout;
    
==== out/jsDeclarationsReactComponents5.d.ts (0 errors) ====
    import PropTypes from 'prop-types';
    declare function Tree({ allowDropOnRoot }: {
        allowDropOnRoot: any;
    }): JSX.Element;
    declare namespace Tree {
        const propTypes: {
            classes: PropTypes.Requireable<object>;
        };
    }
    declare namespace Tree {
        const defaultProps: {
            classes: {};
            parentSource: string;
        };
    }
    export default Tree;
    