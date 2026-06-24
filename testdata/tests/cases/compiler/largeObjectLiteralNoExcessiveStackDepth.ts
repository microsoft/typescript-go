// @noEmit: true

// Regression test for false TS2321 ("Excessive stack depth comparing types") when
// comparing large object literals against complex target types.
// The relater's stack depth was previously limited to 100, which could be exceeded
// by object literals with 100+ properties where each property type required a full
// recursive comparison (e.g. function types, union types).

type RuleLevel = "off" | "warn" | "error";

interface RuleConfig {
    level: RuleLevel;
    options?: Record<string, unknown>;
}

// A target type with many properties and some function-typed members,
// simulating a real-world config schema (e.g. Vite UserConfig, ESLint config).
interface LargeConfig {
    mode: string;
    base: string;
    root: string;
    publicDir: string;
    logLevel: string;
    clearScreen: boolean;
    appType: string;
    define: Record<string, string>;
    plugins: Array<{ name: string; apply: string }>;
    resolve: {
        alias: Record<string, string>;
        conditions: string[];
        extensions: string[];
    };
    css: {
        modules: Record<string, unknown>;
        postcss: string | Record<string, unknown>;
        preprocessorOptions: Record<string, Record<string, unknown>>;
    };
    server: {
        host: string | boolean;
        port: number;
        strictPort: boolean;
        open: boolean | string;
        proxy: Record<string, string | { target: string }>;
    };
    build: {
        outDir: string;
        sourcemap: boolean | string;
        minify: boolean | string;
        target: string | string[];
        rollupOptions: Record<string, unknown>;
    };
    // 200 rule-like properties to exercise stack depth
    rule001: RuleConfig;
    rule002: RuleConfig;
    rule003: RuleConfig;
    rule004: RuleConfig;
    rule005: RuleConfig;
    rule006: RuleConfig;
    rule007: RuleConfig;
    rule008: RuleConfig;
    rule009: RuleConfig;
    rule010: RuleConfig;
    rule011: RuleConfig;
    rule012: RuleConfig;
    rule013: RuleConfig;
    rule014: RuleConfig;
    rule015: RuleConfig;
    rule016: RuleConfig;
    rule017: RuleConfig;
    rule018: RuleConfig;
    rule019: RuleConfig;
    rule020: RuleConfig;
    rule021: RuleConfig;
    rule022: RuleConfig;
    rule023: RuleConfig;
    rule024: RuleConfig;
    rule025: RuleConfig;
    rule026: RuleConfig;
    rule027: RuleConfig;
    rule028: RuleConfig;
    rule029: RuleConfig;
    rule030: RuleConfig;
    rule031: RuleConfig;
    rule032: RuleConfig;
    rule033: RuleConfig;
    rule034: RuleConfig;
    rule035: RuleConfig;
    rule036: RuleConfig;
    rule037: RuleConfig;
    rule038: RuleConfig;
    rule039: RuleConfig;
    rule040: RuleConfig;
    rule041: RuleConfig;
    rule042: RuleConfig;
    rule043: RuleConfig;
    rule044: RuleConfig;
    rule045: RuleConfig;
    rule046: RuleConfig;
    rule047: RuleConfig;
    rule048: RuleConfig;
    rule049: RuleConfig;
    rule050: RuleConfig;
    rule051: RuleConfig;
    rule052: RuleConfig;
    rule053: RuleConfig;
    rule054: RuleConfig;
    rule055: RuleConfig;
    rule056: RuleConfig;
    rule057: RuleConfig;
    rule058: RuleConfig;
    rule059: RuleConfig;
    rule060: RuleConfig;
    rule061: RuleConfig;
    rule062: RuleConfig;
    rule063: RuleConfig;
    rule064: RuleConfig;
    rule065: RuleConfig;
    rule066: RuleConfig;
    rule067: RuleConfig;
    rule068: RuleConfig;
    rule069: RuleConfig;
    rule070: RuleConfig;
    rule071: RuleConfig;
    rule072: RuleConfig;
    rule073: RuleConfig;
    rule074: RuleConfig;
    rule075: RuleConfig;
    rule076: RuleConfig;
    rule077: RuleConfig;
    rule078: RuleConfig;
    rule079: RuleConfig;
    rule080: RuleConfig;
    rule081: RuleConfig;
    rule082: RuleConfig;
    rule083: RuleConfig;
    rule084: RuleConfig;
    rule085: RuleConfig;
    rule086: RuleConfig;
    rule087: RuleConfig;
    rule088: RuleConfig;
    rule089: RuleConfig;
    rule090: RuleConfig;
    rule091: RuleConfig;
    rule092: RuleConfig;
    rule093: RuleConfig;
    rule094: RuleConfig;
    rule095: RuleConfig;
    rule096: RuleConfig;
    rule097: RuleConfig;
    rule098: RuleConfig;
    rule099: RuleConfig;
    rule100: RuleConfig;
    rule101: RuleConfig;
    rule102: RuleConfig;
    rule103: RuleConfig;
    rule104: RuleConfig;
    rule105: RuleConfig;
    rule106: RuleConfig;
    rule107: RuleConfig;
    rule108: RuleConfig;
    rule109: RuleConfig;
    rule110: RuleConfig;
    rule111: RuleConfig;
    rule112: RuleConfig;
    rule113: RuleConfig;
    rule114: RuleConfig;
    rule115: RuleConfig;
    rule116: RuleConfig;
    rule117: RuleConfig;
    rule118: RuleConfig;
    rule119: RuleConfig;
    rule120: RuleConfig;
    rule121: RuleConfig;
    rule122: RuleConfig;
    rule123: RuleConfig;
    rule124: RuleConfig;
    rule125: RuleConfig;
    rule126: RuleConfig;
    rule127: RuleConfig;
    rule128: RuleConfig;
    rule129: RuleConfig;
    rule130: RuleConfig;
    rule131: RuleConfig;
    rule132: RuleConfig;
    rule133: RuleConfig;
    rule134: RuleConfig;
    rule135: RuleConfig;
    rule136: RuleConfig;
    rule137: RuleConfig;
    rule138: RuleConfig;
    rule139: RuleConfig;
    rule140: RuleConfig;
    rule141: RuleConfig;
    rule142: RuleConfig;
    rule143: RuleConfig;
    rule144: RuleConfig;
    rule145: RuleConfig;
    rule146: RuleConfig;
    rule147: RuleConfig;
    rule148: RuleConfig;
    rule149: RuleConfig;
    rule150: RuleConfig;
    rule151: RuleConfig;
    rule152: RuleConfig;
    rule153: RuleConfig;
    rule154: RuleConfig;
    rule155: RuleConfig;
    rule156: RuleConfig;
    rule157: RuleConfig;
    rule158: RuleConfig;
    rule159: RuleConfig;
    rule160: RuleConfig;
    rule161: RuleConfig;
    rule162: RuleConfig;
    rule163: RuleConfig;
    rule164: RuleConfig;
    rule165: RuleConfig;
    rule166: RuleConfig;
    rule167: RuleConfig;
    rule168: RuleConfig;
    rule169: RuleConfig;
    rule170: RuleConfig;
    rule171: RuleConfig;
    rule172: RuleConfig;
    rule173: RuleConfig;
    rule174: RuleConfig;
    rule175: RuleConfig;
    rule176: RuleConfig;
    rule177: RuleConfig;
    rule178: RuleConfig;
    rule179: RuleConfig;
    rule180: RuleConfig;
    rule181: RuleConfig;
    rule182: RuleConfig;
    rule183: RuleConfig;
    rule184: RuleConfig;
    rule185: RuleConfig;
    rule186: RuleConfig;
    rule187: RuleConfig;
    rule188: RuleConfig;
    rule189: RuleConfig;
    rule190: RuleConfig;
    rule191: RuleConfig;
    rule192: RuleConfig;
    rule193: RuleConfig;
    rule194: RuleConfig;
    rule195: RuleConfig;
    rule196: RuleConfig;
    rule197: RuleConfig;
    rule198: RuleConfig;
    rule199: RuleConfig;
    rule200: RuleConfig;
}

// This large object literal should be assignable to LargeConfig without
// producing TS2321 ("Excessive stack depth comparing types").
const config: LargeConfig = {
    mode: "development",
    base: "/",
    root: "/",
    publicDir: "public",
    logLevel: "info",
    clearScreen: true,
    appType: "spa",
    define: {},
    plugins: [],
    resolve: {
        alias: {},
        conditions: [],
        extensions: [".ts", ".tsx", ".js", ".jsx"],
    },
    css: {
        modules: {},
        postcss: "",
        preprocessorOptions: {},
    },
    server: {
        host: "localhost",
        port: 3000,
        strictPort: false,
        open: true,
        proxy: {},
    },
    build: {
        outDir: "dist",
        sourcemap: true,
        minify: true,
        target: "esnext",
        rollupOptions: {},
    },
    rule001: { level: "off" }, rule002: { level: "warn" }, rule003: { level: "error" },
    rule004: { level: "off" }, rule005: { level: "warn" }, rule006: { level: "error" },
    rule007: { level: "off" }, rule008: { level: "warn" }, rule009: { level: "error" },
    rule010: { level: "off" }, rule011: { level: "warn" }, rule012: { level: "error" },
    rule013: { level: "off" }, rule014: { level: "warn" }, rule015: { level: "error" },
    rule016: { level: "off" }, rule017: { level: "warn" }, rule018: { level: "error" },
    rule019: { level: "off" }, rule020: { level: "warn" }, rule021: { level: "error" },
    rule022: { level: "off" }, rule023: { level: "warn" }, rule024: { level: "error" },
    rule025: { level: "off" }, rule026: { level: "warn" }, rule027: { level: "error" },
    rule028: { level: "off" }, rule029: { level: "warn" }, rule030: { level: "error" },
    rule031: { level: "off" }, rule032: { level: "warn" }, rule033: { level: "error" },
    rule034: { level: "off" }, rule035: { level: "warn" }, rule036: { level: "error" },
    rule037: { level: "off" }, rule038: { level: "warn" }, rule039: { level: "error" },
    rule040: { level: "off" }, rule041: { level: "warn" }, rule042: { level: "error" },
    rule043: { level: "off" }, rule044: { level: "warn" }, rule045: { level: "error" },
    rule046: { level: "off" }, rule047: { level: "warn" }, rule048: { level: "error" },
    rule049: { level: "off" }, rule050: { level: "warn" }, rule051: { level: "error" },
    rule052: { level: "off" }, rule053: { level: "warn" }, rule054: { level: "error" },
    rule055: { level: "off" }, rule056: { level: "warn" }, rule057: { level: "error" },
    rule058: { level: "off" }, rule059: { level: "warn" }, rule060: { level: "error" },
    rule061: { level: "off" }, rule062: { level: "warn" }, rule063: { level: "error" },
    rule064: { level: "off" }, rule065: { level: "warn" }, rule066: { level: "error" },
    rule067: { level: "off" }, rule068: { level: "warn" }, rule069: { level: "error" },
    rule070: { level: "off" }, rule071: { level: "warn" }, rule072: { level: "error" },
    rule073: { level: "off" }, rule074: { level: "warn" }, rule075: { level: "error" },
    rule076: { level: "off" }, rule077: { level: "warn" }, rule078: { level: "error" },
    rule079: { level: "off" }, rule080: { level: "warn" }, rule081: { level: "error" },
    rule082: { level: "off" }, rule083: { level: "warn" }, rule084: { level: "error" },
    rule085: { level: "off" }, rule086: { level: "warn" }, rule087: { level: "error" },
    rule088: { level: "off" }, rule089: { level: "warn" }, rule090: { level: "error" },
    rule091: { level: "off" }, rule092: { level: "warn" }, rule093: { level: "error" },
    rule094: { level: "off" }, rule095: { level: "warn" }, rule096: { level: "error" },
    rule097: { level: "off" }, rule098: { level: "warn" }, rule099: { level: "error" },
    rule100: { level: "off" }, rule101: { level: "warn" }, rule102: { level: "error" },
    rule103: { level: "off" }, rule104: { level: "warn" }, rule105: { level: "error" },
    rule106: { level: "off" }, rule107: { level: "warn" }, rule108: { level: "error" },
    rule109: { level: "off" }, rule110: { level: "warn" }, rule111: { level: "error" },
    rule112: { level: "off" }, rule113: { level: "warn" }, rule114: { level: "error" },
    rule115: { level: "off" }, rule116: { level: "warn" }, rule117: { level: "error" },
    rule118: { level: "off" }, rule119: { level: "warn" }, rule120: { level: "error" },
    rule121: { level: "off" }, rule122: { level: "warn" }, rule123: { level: "error" },
    rule124: { level: "off" }, rule125: { level: "warn" }, rule126: { level: "error" },
    rule127: { level: "off" }, rule128: { level: "warn" }, rule129: { level: "error" },
    rule130: { level: "off" }, rule131: { level: "warn" }, rule132: { level: "error" },
    rule133: { level: "off" }, rule134: { level: "warn" }, rule135: { level: "error" },
    rule136: { level: "off" }, rule137: { level: "warn" }, rule138: { level: "error" },
    rule139: { level: "off" }, rule140: { level: "warn" }, rule141: { level: "error" },
    rule142: { level: "off" }, rule143: { level: "warn" }, rule144: { level: "error" },
    rule145: { level: "off" }, rule146: { level: "warn" }, rule147: { level: "error" },
    rule148: { level: "off" }, rule149: { level: "warn" }, rule150: { level: "error" },
    rule151: { level: "off" }, rule152: { level: "warn" }, rule153: { level: "error" },
    rule154: { level: "off" }, rule155: { level: "warn" }, rule156: { level: "error" },
    rule157: { level: "off" }, rule158: { level: "warn" }, rule159: { level: "error" },
    rule160: { level: "off" }, rule161: { level: "warn" }, rule162: { level: "error" },
    rule163: { level: "off" }, rule164: { level: "warn" }, rule165: { level: "error" },
    rule166: { level: "off" }, rule167: { level: "warn" }, rule168: { level: "error" },
    rule169: { level: "off" }, rule170: { level: "warn" }, rule171: { level: "error" },
    rule172: { level: "off" }, rule173: { level: "warn" }, rule174: { level: "error" },
    rule175: { level: "off" }, rule176: { level: "warn" }, rule177: { level: "error" },
    rule178: { level: "off" }, rule179: { level: "warn" }, rule180: { level: "error" },
    rule181: { level: "off" }, rule182: { level: "warn" }, rule183: { level: "error" },
    rule184: { level: "off" }, rule185: { level: "warn" }, rule186: { level: "error" },
    rule187: { level: "off" }, rule188: { level: "warn" }, rule189: { level: "error" },
    rule190: { level: "off" }, rule191: { level: "warn" }, rule192: { level: "error" },
    rule193: { level: "off" }, rule194: { level: "warn" }, rule195: { level: "error" },
    rule196: { level: "off" }, rule197: { level: "warn" }, rule198: { level: "error" },
    rule199: { level: "off" }, rule200: { level: "warn" },
};

// Also test function-typed properties which require deeper recursion
type FnType = (a: string, b: number) => boolean;

interface ConfigWithFunctions {
    fn001: FnType; fn002: FnType; fn003: FnType; fn004: FnType; fn005: FnType;
    fn006: FnType; fn007: FnType; fn008: FnType; fn009: FnType; fn010: FnType;
    fn011: FnType; fn012: FnType; fn013: FnType; fn014: FnType; fn015: FnType;
    fn016: FnType; fn017: FnType; fn018: FnType; fn019: FnType; fn020: FnType;
    fn021: FnType; fn022: FnType; fn023: FnType; fn024: FnType; fn025: FnType;
    fn026: FnType; fn027: FnType; fn028: FnType; fn029: FnType; fn030: FnType;
    fn031: FnType; fn032: FnType; fn033: FnType; fn034: FnType; fn035: FnType;
    fn036: FnType; fn037: FnType; fn038: FnType; fn039: FnType; fn040: FnType;
    fn041: FnType; fn042: FnType; fn043: FnType; fn044: FnType; fn045: FnType;
    fn046: FnType; fn047: FnType; fn048: FnType; fn049: FnType; fn050: FnType;
    fn051: FnType; fn052: FnType; fn053: FnType; fn054: FnType; fn055: FnType;
    fn056: FnType; fn057: FnType; fn058: FnType; fn059: FnType; fn060: FnType;
    fn061: FnType; fn062: FnType; fn063: FnType; fn064: FnType; fn065: FnType;
    fn066: FnType; fn067: FnType; fn068: FnType; fn069: FnType; fn070: FnType;
    fn071: FnType; fn072: FnType; fn073: FnType; fn074: FnType; fn075: FnType;
    fn076: FnType; fn077: FnType; fn078: FnType; fn079: FnType; fn080: FnType;
    fn081: FnType; fn082: FnType; fn083: FnType; fn084: FnType; fn085: FnType;
    fn086: FnType; fn087: FnType; fn088: FnType; fn089: FnType; fn090: FnType;
    fn091: FnType; fn092: FnType; fn093: FnType; fn094: FnType; fn095: FnType;
    fn096: FnType; fn097: FnType; fn098: FnType; fn099: FnType; fn100: FnType;
    fn101: FnType; fn102: FnType; fn103: FnType; fn104: FnType; fn105: FnType;
    fn106: FnType; fn107: FnType; fn108: FnType; fn109: FnType; fn110: FnType;
    fn111: FnType; fn112: FnType; fn113: FnType; fn114: FnType; fn115: FnType;
    fn116: FnType; fn117: FnType; fn118: FnType; fn119: FnType; fn120: FnType;
    fn121: FnType; fn122: FnType; fn123: FnType; fn124: FnType; fn125: FnType;
    fn126: FnType; fn127: FnType; fn128: FnType; fn129: FnType; fn130: FnType;
    fn131: FnType; fn132: FnType; fn133: FnType; fn134: FnType; fn135: FnType;
    fn136: FnType; fn137: FnType; fn138: FnType; fn139: FnType; fn140: FnType;
    fn141: FnType; fn142: FnType; fn143: FnType; fn144: FnType; fn145: FnType;
    fn146: FnType; fn147: FnType; fn148: FnType; fn149: FnType; fn150: FnType;
}

const fns: ConfigWithFunctions = {
    fn001: (a, b) => a.length > b, fn002: (a, b) => a.length > b, fn003: (a, b) => a.length > b,
    fn004: (a, b) => a.length > b, fn005: (a, b) => a.length > b, fn006: (a, b) => a.length > b,
    fn007: (a, b) => a.length > b, fn008: (a, b) => a.length > b, fn009: (a, b) => a.length > b,
    fn010: (a, b) => a.length > b, fn011: (a, b) => a.length > b, fn012: (a, b) => a.length > b,
    fn013: (a, b) => a.length > b, fn014: (a, b) => a.length > b, fn015: (a, b) => a.length > b,
    fn016: (a, b) => a.length > b, fn017: (a, b) => a.length > b, fn018: (a, b) => a.length > b,
    fn019: (a, b) => a.length > b, fn020: (a, b) => a.length > b, fn021: (a, b) => a.length > b,
    fn022: (a, b) => a.length > b, fn023: (a, b) => a.length > b, fn024: (a, b) => a.length > b,
    fn025: (a, b) => a.length > b, fn026: (a, b) => a.length > b, fn027: (a, b) => a.length > b,
    fn028: (a, b) => a.length > b, fn029: (a, b) => a.length > b, fn030: (a, b) => a.length > b,
    fn031: (a, b) => a.length > b, fn032: (a, b) => a.length > b, fn033: (a, b) => a.length > b,
    fn034: (a, b) => a.length > b, fn035: (a, b) => a.length > b, fn036: (a, b) => a.length > b,
    fn037: (a, b) => a.length > b, fn038: (a, b) => a.length > b, fn039: (a, b) => a.length > b,
    fn040: (a, b) => a.length > b, fn041: (a, b) => a.length > b, fn042: (a, b) => a.length > b,
    fn043: (a, b) => a.length > b, fn044: (a, b) => a.length > b, fn045: (a, b) => a.length > b,
    fn046: (a, b) => a.length > b, fn047: (a, b) => a.length > b, fn048: (a, b) => a.length > b,
    fn049: (a, b) => a.length > b, fn050: (a, b) => a.length > b, fn051: (a, b) => a.length > b,
    fn052: (a, b) => a.length > b, fn053: (a, b) => a.length > b, fn054: (a, b) => a.length > b,
    fn055: (a, b) => a.length > b, fn056: (a, b) => a.length > b, fn057: (a, b) => a.length > b,
    fn058: (a, b) => a.length > b, fn059: (a, b) => a.length > b, fn060: (a, b) => a.length > b,
    fn061: (a, b) => a.length > b, fn062: (a, b) => a.length > b, fn063: (a, b) => a.length > b,
    fn064: (a, b) => a.length > b, fn065: (a, b) => a.length > b, fn066: (a, b) => a.length > b,
    fn067: (a, b) => a.length > b, fn068: (a, b) => a.length > b, fn069: (a, b) => a.length > b,
    fn070: (a, b) => a.length > b, fn071: (a, b) => a.length > b, fn072: (a, b) => a.length > b,
    fn073: (a, b) => a.length > b, fn074: (a, b) => a.length > b, fn075: (a, b) => a.length > b,
    fn076: (a, b) => a.length > b, fn077: (a, b) => a.length > b, fn078: (a, b) => a.length > b,
    fn079: (a, b) => a.length > b, fn080: (a, b) => a.length > b, fn081: (a, b) => a.length > b,
    fn082: (a, b) => a.length > b, fn083: (a, b) => a.length > b, fn084: (a, b) => a.length > b,
    fn085: (a, b) => a.length > b, fn086: (a, b) => a.length > b, fn087: (a, b) => a.length > b,
    fn088: (a, b) => a.length > b, fn089: (a, b) => a.length > b, fn090: (a, b) => a.length > b,
    fn091: (a, b) => a.length > b, fn092: (a, b) => a.length > b, fn093: (a, b) => a.length > b,
    fn094: (a, b) => a.length > b, fn095: (a, b) => a.length > b, fn096: (a, b) => a.length > b,
    fn097: (a, b) => a.length > b, fn098: (a, b) => a.length > b, fn099: (a, b) => a.length > b,
    fn100: (a, b) => a.length > b, fn101: (a, b) => a.length > b, fn102: (a, b) => a.length > b,
    fn103: (a, b) => a.length > b, fn104: (a, b) => a.length > b, fn105: (a, b) => a.length > b,
    fn106: (a, b) => a.length > b, fn107: (a, b) => a.length > b, fn108: (a, b) => a.length > b,
    fn109: (a, b) => a.length > b, fn110: (a, b) => a.length > b, fn111: (a, b) => a.length > b,
    fn112: (a, b) => a.length > b, fn113: (a, b) => a.length > b, fn114: (a, b) => a.length > b,
    fn115: (a, b) => a.length > b, fn116: (a, b) => a.length > b, fn117: (a, b) => a.length > b,
    fn118: (a, b) => a.length > b, fn119: (a, b) => a.length > b, fn120: (a, b) => a.length > b,
    fn121: (a, b) => a.length > b, fn122: (a, b) => a.length > b, fn123: (a, b) => a.length > b,
    fn124: (a, b) => a.length > b, fn125: (a, b) => a.length > b, fn126: (a, b) => a.length > b,
    fn127: (a, b) => a.length > b, fn128: (a, b) => a.length > b, fn129: (a, b) => a.length > b,
    fn130: (a, b) => a.length > b, fn131: (a, b) => a.length > b, fn132: (a, b) => a.length > b,
    fn133: (a, b) => a.length > b, fn134: (a, b) => a.length > b, fn135: (a, b) => a.length > b,
    fn136: (a, b) => a.length > b, fn137: (a, b) => a.length > b, fn138: (a, b) => a.length > b,
    fn139: (a, b) => a.length > b, fn140: (a, b) => a.length > b, fn141: (a, b) => a.length > b,
    fn142: (a, b) => a.length > b, fn143: (a, b) => a.length > b, fn144: (a, b) => a.length > b,
    fn145: (a, b) => a.length > b, fn146: (a, b) => a.length > b, fn147: (a, b) => a.length > b,
    fn148: (a, b) => a.length > b, fn149: (a, b) => a.length > b, fn150: (a, b) => a.length > b,
};
