//// [tests/cases/compiler/jsDeclarationsExpandoAlias.ts] ////

//// [main.js]
/**
 * @param {Object} props
 * @param {string} props.name
 */
function FunctionComponent(props) {
    return props.name;
}

FunctionComponent.propTypes = {
    name: "required",
};

export const FunctionComponentAlias = FunctionComponent;

/**
 * @param {Object} props
 * @param {string} props.name
 */
const ArrowComponent = (props) => {
    return props.name;
};

ArrowComponent.propTypes = {
    name: "required",
};

export const ArrowComponentAlias = ArrowComponent;




//// [main.d.ts]
/**
 * @param {Object} props
 * @param {string} props.name
 */
export declare function FunctionComponentAlias(props: {
    name: string;
}): string;
export declare namespace FunctionComponentAlias {
    var propTypes: {
        name: string;
    };
}
export declare const ArrowComponentAlias: {
    (props: {
        name: string;
    }): string;
    propTypes: {
        name: string;
    };
};
