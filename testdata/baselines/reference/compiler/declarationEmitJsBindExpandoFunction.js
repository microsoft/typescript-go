//// [tests/cases/compiler/declarationEmitJsBindExpandoFunction.ts] ////

//// [repro.js]
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




//// [repro.d.ts]
/** @param {{ text: string }} args */
export declare function PublicInternalBinding(args: {
    text: string;
}): string;
export declare namespace PublicInternalBinding {
    var args: {
        text: string;
    };
}
/** @param {{ text: string }} args */
export declare function PublicPlainBinding(args: {
    text: string;
}): string;
export declare namespace PublicPlainBinding {
    var args: {
        text: string;
    };
}
/** @param {{ text: string }} args */
export declare function PublicMixedBinding(args: {
    text: string;
}): string;
export declare namespace PublicMixedBinding {
    var second: {
        count: number;
    };
}
export declare namespace PublicMixedBinding {
    var first: {
        count: number;
    };
}
