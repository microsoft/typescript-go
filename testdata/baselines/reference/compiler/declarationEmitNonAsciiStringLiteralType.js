//// [tests/cases/compiler/declarationEmitNonAsciiStringLiteralType.ts] ////

//// [declarationEmitNonAsciiStringLiteralType.ts]
export const IconEmojis = {
    alert_low: "⚠️",
} as const;

export const singleEmoji = "⚠️" as const;

export const tuple = ["⚠️", "日本語"] as const;

export function returnsEmoji(): "⚠️" {
    return "⚠️";
}




//// [declarationEmitNonAsciiStringLiteralType.d.ts]
export declare const IconEmojis: {
    readonly alert_low: "⚠️";
};
export declare const singleEmoji: "⚠️";
export declare const tuple: readonly ["⚠️", "日本語"];
export declare function returnsEmoji(): "⚠️";
