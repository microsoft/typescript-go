// @allowJs: true
// @checkJs: true
// @declaration: true
// @emitDeclarationOnly: true

// @filename: repro.js
/** @param {{ text: string }} args */
function Internal(args) {
    return args.text;
}
Internal.args = { text: "source" };

export const PublicInternalBinding = Internal.bind({});

/** @param {{ text: string }} args */
function Plain(args) {
    return args.text;
}
export const PublicPlainBinding = Plain.bind({});
// @ts-ignore
PublicPlainBinding.args = { text: "bound" };

/** @param {{ text: string }} args */
function Mixed(args) {
    return args.text;
}
Mixed.first = { text: "source" };
Mixed.second = { count: 0 };

export const PublicMixedBinding = Mixed.bind({});
// @ts-ignore
PublicMixedBinding.first = { count: 1 };
