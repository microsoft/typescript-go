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
import PropTypes from 'prop-types';
declare function Foo({ bar }: {
    bar: any;
}): JSX.Element;
declare namespace Foo {
    var propTypes: {
        bar: PropTypes.Requireable<boolean>;
    };
    var defaultProps: {
        bar: boolean;
    };
}
export default Foo;


//// [DtsFileErrors]


jsxDeclarationsWithEsModuleInteropNoCrash.d.ts(1,23): error TS2307: Cannot find module 'prop-types' or its corresponding type declarations.
jsxDeclarationsWithEsModuleInteropNoCrash.d.ts(4,5): error TS2503: Cannot find namespace 'JSX'.


==== jsxDeclarationsWithEsModuleInteropNoCrash.d.ts (2 errors) ====
    import PropTypes from 'prop-types';
                          ~~~~~~~~~~~~
!!! error TS2307: Cannot find module 'prop-types' or its corresponding type declarations.
    declare function Foo({ bar }: {
        bar: any;
    }): JSX.Element;
        ~~~
!!! error TS2503: Cannot find namespace 'JSX'.
    declare namespace Foo {
        var propTypes: {
            bar: PropTypes.Requireable<boolean>;
        };
        var defaultProps: {
            bar: boolean;
        };
    }
    export default Foo;
    