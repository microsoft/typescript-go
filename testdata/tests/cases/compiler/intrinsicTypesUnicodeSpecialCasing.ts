// @strict: true

// Test Unicode special case mappings for intrinsic string types
// These characters have 1:many case mappings that Go's strings.ToUpper/ToLower
// don't handle, but JavaScript's toUpperCase()/toLowerCase() do.

// ß (U+00DF) uppercases to "SS" in JavaScript
type T1 = Uppercase<"ß">;
const t1: T1 = "SS";

// İ (U+0130) lowercases to "i̇" (i + combining dot above U+0307) in JavaScript
type T2 = Lowercase<"İ">;
const t2: T2 = "i\u0307";

// Ligatures: ﬁ (U+FB01) uppercases to "FI"
type T3 = Uppercase<"ﬁ">;
const t3: T3 = "FI";

// ﬂ (U+FB02) uppercases to "FL"
type T4 = Uppercase<"ﬂ">;
const t4: T4 = "FL";

// ﬀ (U+FB00) uppercases to "FF"
type T5 = Uppercase<"ﬀ">;
const t5: T5 = "FF";

// Capitalize should only affect first character
type T6 = Capitalize<"ßtest">;
const t6: T6 = "SStest";

// Uncapitalize with İ
type T7 = Uncapitalize<"İSPANYOL">;
const t7: T7 = "i\u0307SPANYOL";

// Mixed string with special characters
type T8 = Uppercase<"straße">;
const t8: T8 = "STRASSE";
