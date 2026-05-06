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



//// [jsxDeclarationsWithEsModuleInteropNoCrash.d.ts]
/// <reference path="../.lib/react16.d.ts" preserve="true" />
import PropTypes from 'prop-types';
function Foo({ bar }: {
    bar: any;
}): JSX.Element;
declare namespace Foo {
    var propTypes: {
        bar: PropTypes.Requireable<boolean>;
    };
}
declare namespace Foo {
    var defaultProps: {
        bar: boolean;
    };
}
export default Foo;


//// [DtsFileErrors]


jsxDeclarationsWithEsModuleInteropNoCrash.d.ts(3,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== jsxDeclarationsWithEsModuleInteropNoCrash.d.ts (1 errors) ====
    /// <reference path="../.lib/react16.d.ts" preserve="true" />
    import PropTypes from 'prop-types';
    function Foo({ bar }: {
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        bar: any;
    }): JSX.Element;
    declare namespace Foo {
        var propTypes: {
            bar: PropTypes.Requireable<boolean>;
        };
    }
    declare namespace Foo {
        var defaultProps: {
            bar: boolean;
        };
    }
    export default Foo;
    