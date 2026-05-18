// @declaration: true
// @emitDeclarationOnly: true

export const IconEmojis = {
    alert_low: "⚠️",
} as const;

export const singleEmoji = "⚠️" as const;

export const tuple = ["⚠️", "日本語"] as const;

export function returnsEmoji(): "⚠️" {
    return "⚠️";
}
