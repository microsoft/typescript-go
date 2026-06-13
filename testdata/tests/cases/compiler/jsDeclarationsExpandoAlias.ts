// @allowJs: true
// @checkJs: true
// @declaration: true
// @emitDeclarationOnly: true
// @stableTypeOrdering: true
// @filename: src/main.js

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
