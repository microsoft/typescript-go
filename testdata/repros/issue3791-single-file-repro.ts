// Single-file repro for https://github.com/microsoft/typescript-go/issues/3791
// Expected result when compiled with tsgo: panic
// "Debug failure. False expression: generalized source shouldn't be assignable"
//
// Example:
// ./built/local/tsgo --noEmit --strict false --moduleResolution bundler --module esnext --target es2017 --jsx react --lib es6,dom --skipLibCheck true --esModuleInterop true --allowSyntheticDefaultImports true --resolveJsonModule true --importHelpers true ./testdata/repros/issue3791-single-file-repro.ts

// ===== helpers.d.ts =====
type EmptyAsNever<T extends {}> = T[keyof T] extends never ? never : WithoutNever<T>;

type PickBySubType<T extends {}, TSub> = Pick<T, ({ [key in keyof T]: T[key] extends TSub ? key : never }[keyof T])>;

type ExcludeFieldsWithType<T extends {}, TSub> = Pick<T, ({ [key in keyof T]: T[key] extends TSub ? never : key }[keyof T])>;

type FullFilled<T extends {}> = { [key in keyof T]-?: T[key] };

type IsAny<T> = 0 extends (1 & T) ? true : false;

type IsNever<T> = [T] extends [never] ? true : false;

type ParamWithoutNever<T> = [T] extends [never] ? [] : [T];

type Override<T, R> = Omit<T, keyof R> & R;

type ValueOf<T> = T[keyof T];

type WithoutNever<T extends {}> = ExcludeFieldsWithType<T, never>;

type Writable<T extends {}> = { -readonly [key in keyof T]: T[key] };

type ValueOfArray<T extends any[]> = T extends Array<infer R> ? R : never;

type DeepPartial<T> = (
    T extends Array<infer U> ? Array<DeepPartial<U>> :
    T extends Map<infer K, infer V> ? Map<DeepPartial<K>, DeepPartial<V>> :
    T extends Set<infer M> ? Set<DeepPartial<M>> :
    T extends object ? { [K in keyof T]?: DeepPartial<T[K]> } :
    T
);

type AtLeastOne<T extends Record<string, any>> = keyof T extends infer K ? K extends string ? Pick<T, K & keyof T> & Partial<T> : never : never

type FindKeyByValue<T extends {}, V> = {
    [key in keyof T]: T[key] extends V ? key : never;
}[keyof T];


type Reverse<T extends {}> = {
    [val in ValueOf<T>]: FindKeyByValue<T, val>;
}

type UnionToIntersection<U> =
    (U extends any ? (k: U) => void : never) extends ((k: infer I) => void) ? I : never

type ArrayToIntersection<A> =
    A extends Array<infer U> ? UnionToIntersection<U> : never;

// ===== _big-union.ts =====
// Generated wide union mirroring IPrefixedServerTypes (~150 members)
interface T001 { $a001: string; $b001: number; $c001?: { $d001: string; $e001?: { $f001: number }[] } }
interface T002 { $a002: string; $b002: number; $c002?: { $d002: string; $e002?: { $f002: number }[] } }
interface T003 { $a003: string; $b003: number; $c003?: { $d003: string; $e003?: { $f003: number }[] } }
interface T004 { $a004: string; $b004: number; $c004?: { $d004: string; $e004?: { $f004: number }[] } }
interface T005 { $a005: string; $b005: number; $c005?: { $d005: string; $e005?: { $f005: number }[] } }
interface T006 { $a006: string; $b006: number; $c006?: { $d006: string; $e006?: { $f006: number }[] } }
interface T007 { $a007: string; $b007: number; $c007?: { $d007: string; $e007?: { $f007: number }[] } }
interface T008 { $a008: string; $b008: number; $c008?: { $d008: string; $e008?: { $f008: number }[] } }
interface T009 { $a009: string; $b009: number; $c009?: { $d009: string; $e009?: { $f009: number }[] } }
interface T010 { $a010: string; $b010: number; $c010?: { $d010: string; $e010?: { $f010: number }[] } }
interface T011 { $a011: string; $b011: number; $c011?: { $d011: string; $e011?: { $f011: number }[] } }
interface T012 { $a012: string; $b012: number; $c012?: { $d012: string; $e012?: { $f012: number }[] } }
interface T013 { $a013: string; $b013: number; $c013?: { $d013: string; $e013?: { $f013: number }[] } }
interface T014 { $a014: string; $b014: number; $c014?: { $d014: string; $e014?: { $f014: number }[] } }
interface T015 { $a015: string; $b015: number; $c015?: { $d015: string; $e015?: { $f015: number }[] } }
interface T016 { $a016: string; $b016: number; $c016?: { $d016: string; $e016?: { $f016: number }[] } }
interface T017 { $a017: string; $b017: number; $c017?: { $d017: string; $e017?: { $f017: number }[] } }
interface T018 { $a018: string; $b018: number; $c018?: { $d018: string; $e018?: { $f018: number }[] } }
interface T019 { $a019: string; $b019: number; $c019?: { $d019: string; $e019?: { $f019: number }[] } }
interface T020 { $a020: string; $b020: number; $c020?: { $d020: string; $e020?: { $f020: number }[] } }
interface T021 { $a021: string; $b021: number; $c021?: { $d021: string; $e021?: { $f021: number }[] } }
interface T022 { $a022: string; $b022: number; $c022?: { $d022: string; $e022?: { $f022: number }[] } }
interface T023 { $a023: string; $b023: number; $c023?: { $d023: string; $e023?: { $f023: number }[] } }
interface T024 { $a024: string; $b024: number; $c024?: { $d024: string; $e024?: { $f024: number }[] } }
interface T025 { $a025: string; $b025: number; $c025?: { $d025: string; $e025?: { $f025: number }[] } }
interface T026 { $a026: string; $b026: number; $c026?: { $d026: string; $e026?: { $f026: number }[] } }
interface T027 { $a027: string; $b027: number; $c027?: { $d027: string; $e027?: { $f027: number }[] } }
interface T028 { $a028: string; $b028: number; $c028?: { $d028: string; $e028?: { $f028: number }[] } }
interface T029 { $a029: string; $b029: number; $c029?: { $d029: string; $e029?: { $f029: number }[] } }
interface T030 { $a030: string; $b030: number; $c030?: { $d030: string; $e030?: { $f030: number }[] } }
interface T031 { $a031: string; $b031: number; $c031?: { $d031: string; $e031?: { $f031: number }[] } }
interface T032 { $a032: string; $b032: number; $c032?: { $d032: string; $e032?: { $f032: number }[] } }
interface T033 { $a033: string; $b033: number; $c033?: { $d033: string; $e033?: { $f033: number }[] } }
interface T034 { $a034: string; $b034: number; $c034?: { $d034: string; $e034?: { $f034: number }[] } }
interface T035 { $a035: string; $b035: number; $c035?: { $d035: string; $e035?: { $f035: number }[] } }
interface T036 { $a036: string; $b036: number; $c036?: { $d036: string; $e036?: { $f036: number }[] } }
interface T037 { $a037: string; $b037: number; $c037?: { $d037: string; $e037?: { $f037: number }[] } }
interface T038 { $a038: string; $b038: number; $c038?: { $d038: string; $e038?: { $f038: number }[] } }
interface T039 { $a039: string; $b039: number; $c039?: { $d039: string; $e039?: { $f039: number }[] } }
interface T040 { $a040: string; $b040: number; $c040?: { $d040: string; $e040?: { $f040: number }[] } }
interface T041 { $a041: string; $b041: number; $c041?: { $d041: string; $e041?: { $f041: number }[] } }
interface T042 { $a042: string; $b042: number; $c042?: { $d042: string; $e042?: { $f042: number }[] } }
interface T043 { $a043: string; $b043: number; $c043?: { $d043: string; $e043?: { $f043: number }[] } }
interface T044 { $a044: string; $b044: number; $c044?: { $d044: string; $e044?: { $f044: number }[] } }
interface T045 { $a045: string; $b045: number; $c045?: { $d045: string; $e045?: { $f045: number }[] } }
interface T046 { $a046: string; $b046: number; $c046?: { $d046: string; $e046?: { $f046: number }[] } }
interface T047 { $a047: string; $b047: number; $c047?: { $d047: string; $e047?: { $f047: number }[] } }
interface T048 { $a048: string; $b048: number; $c048?: { $d048: string; $e048?: { $f048: number }[] } }
interface T049 { $a049: string; $b049: number; $c049?: { $d049: string; $e049?: { $f049: number }[] } }
interface T050 { $a050: string; $b050: number; $c050?: { $d050: string; $e050?: { $f050: number }[] } }
interface T051 { $a051: string; $b051: number; $c051?: { $d051: string; $e051?: { $f051: number }[] } }
interface T052 { $a052: string; $b052: number; $c052?: { $d052: string; $e052?: { $f052: number }[] } }
interface T053 { $a053: string; $b053: number; $c053?: { $d053: string; $e053?: { $f053: number }[] } }
interface T054 { $a054: string; $b054: number; $c054?: { $d054: string; $e054?: { $f054: number }[] } }
interface T055 { $a055: string; $b055: number; $c055?: { $d055: string; $e055?: { $f055: number }[] } }
interface T056 { $a056: string; $b056: number; $c056?: { $d056: string; $e056?: { $f056: number }[] } }
interface T057 { $a057: string; $b057: number; $c057?: { $d057: string; $e057?: { $f057: number }[] } }
interface T058 { $a058: string; $b058: number; $c058?: { $d058: string; $e058?: { $f058: number }[] } }
interface T059 { $a059: string; $b059: number; $c059?: { $d059: string; $e059?: { $f059: number }[] } }
interface T060 { $a060: string; $b060: number; $c060?: { $d060: string; $e060?: { $f060: number }[] } }
interface T061 { $a061: string; $b061: number; $c061?: { $d061: string; $e061?: { $f061: number }[] } }
interface T062 { $a062: string; $b062: number; $c062?: { $d062: string; $e062?: { $f062: number }[] } }
interface T063 { $a063: string; $b063: number; $c063?: { $d063: string; $e063?: { $f063: number }[] } }
interface T064 { $a064: string; $b064: number; $c064?: { $d064: string; $e064?: { $f064: number }[] } }
interface T065 { $a065: string; $b065: number; $c065?: { $d065: string; $e065?: { $f065: number }[] } }
interface T066 { $a066: string; $b066: number; $c066?: { $d066: string; $e066?: { $f066: number }[] } }
interface T067 { $a067: string; $b067: number; $c067?: { $d067: string; $e067?: { $f067: number }[] } }
interface T068 { $a068: string; $b068: number; $c068?: { $d068: string; $e068?: { $f068: number }[] } }
interface T069 { $a069: string; $b069: number; $c069?: { $d069: string; $e069?: { $f069: number }[] } }
interface T070 { $a070: string; $b070: number; $c070?: { $d070: string; $e070?: { $f070: number }[] } }
interface T071 { $a071: string; $b071: number; $c071?: { $d071: string; $e071?: { $f071: number }[] } }
interface T072 { $a072: string; $b072: number; $c072?: { $d072: string; $e072?: { $f072: number }[] } }
interface T073 { $a073: string; $b073: number; $c073?: { $d073: string; $e073?: { $f073: number }[] } }
interface T074 { $a074: string; $b074: number; $c074?: { $d074: string; $e074?: { $f074: number }[] } }
interface T075 { $a075: string; $b075: number; $c075?: { $d075: string; $e075?: { $f075: number }[] } }
interface T076 { $a076: string; $b076: number; $c076?: { $d076: string; $e076?: { $f076: number }[] } }
interface T077 { $a077: string; $b077: number; $c077?: { $d077: string; $e077?: { $f077: number }[] } }
interface T078 { $a078: string; $b078: number; $c078?: { $d078: string; $e078?: { $f078: number }[] } }
interface T079 { $a079: string; $b079: number; $c079?: { $d079: string; $e079?: { $f079: number }[] } }
interface T080 { $a080: string; $b080: number; $c080?: { $d080: string; $e080?: { $f080: number }[] } }
interface T081 { $a081: string; $b081: number; $c081?: { $d081: string; $e081?: { $f081: number }[] } }
interface T082 { $a082: string; $b082: number; $c082?: { $d082: string; $e082?: { $f082: number }[] } }
interface T083 { $a083: string; $b083: number; $c083?: { $d083: string; $e083?: { $f083: number }[] } }
interface T084 { $a084: string; $b084: number; $c084?: { $d084: string; $e084?: { $f084: number }[] } }
interface T085 { $a085: string; $b085: number; $c085?: { $d085: string; $e085?: { $f085: number }[] } }
interface T086 { $a086: string; $b086: number; $c086?: { $d086: string; $e086?: { $f086: number }[] } }
interface T087 { $a087: string; $b087: number; $c087?: { $d087: string; $e087?: { $f087: number }[] } }
interface T088 { $a088: string; $b088: number; $c088?: { $d088: string; $e088?: { $f088: number }[] } }
interface T089 { $a089: string; $b089: number; $c089?: { $d089: string; $e089?: { $f089: number }[] } }
interface T090 { $a090: string; $b090: number; $c090?: { $d090: string; $e090?: { $f090: number }[] } }
interface T091 { $a091: string; $b091: number; $c091?: { $d091: string; $e091?: { $f091: number }[] } }
interface T092 { $a092: string; $b092: number; $c092?: { $d092: string; $e092?: { $f092: number }[] } }
interface T093 { $a093: string; $b093: number; $c093?: { $d093: string; $e093?: { $f093: number }[] } }
interface T094 { $a094: string; $b094: number; $c094?: { $d094: string; $e094?: { $f094: number }[] } }
interface T095 { $a095: string; $b095: number; $c095?: { $d095: string; $e095?: { $f095: number }[] } }
interface T096 { $a096: string; $b096: number; $c096?: { $d096: string; $e096?: { $f096: number }[] } }
interface T097 { $a097: string; $b097: number; $c097?: { $d097: string; $e097?: { $f097: number }[] } }
interface T098 { $a098: string; $b098: number; $c098?: { $d098: string; $e098?: { $f098: number }[] } }
interface T099 { $a099: string; $b099: number; $c099?: { $d099: string; $e099?: { $f099: number }[] } }
interface T100 { $a100: string; $b100: number; $c100?: { $d100: string; $e100?: { $f100: number }[] } }
interface T101 { $a101: string; $b101: number; $c101?: { $d101: string; $e101?: { $f101: number }[] } }
interface T102 { $a102: string; $b102: number; $c102?: { $d102: string; $e102?: { $f102: number }[] } }
interface T103 { $a103: string; $b103: number; $c103?: { $d103: string; $e103?: { $f103: number }[] } }
interface T104 { $a104: string; $b104: number; $c104?: { $d104: string; $e104?: { $f104: number }[] } }
interface T105 { $a105: string; $b105: number; $c105?: { $d105: string; $e105?: { $f105: number }[] } }
interface T106 { $a106: string; $b106: number; $c106?: { $d106: string; $e106?: { $f106: number }[] } }
interface T107 { $a107: string; $b107: number; $c107?: { $d107: string; $e107?: { $f107: number }[] } }
interface T108 { $a108: string; $b108: number; $c108?: { $d108: string; $e108?: { $f108: number }[] } }
interface T109 { $a109: string; $b109: number; $c109?: { $d109: string; $e109?: { $f109: number }[] } }
interface T110 { $a110: string; $b110: number; $c110?: { $d110: string; $e110?: { $f110: number }[] } }
interface T111 { $a111: string; $b111: number; $c111?: { $d111: string; $e111?: { $f111: number }[] } }
interface T112 { $a112: string; $b112: number; $c112?: { $d112: string; $e112?: { $f112: number }[] } }
interface T113 { $a113: string; $b113: number; $c113?: { $d113: string; $e113?: { $f113: number }[] } }
interface T114 { $a114: string; $b114: number; $c114?: { $d114: string; $e114?: { $f114: number }[] } }
interface T115 { $a115: string; $b115: number; $c115?: { $d115: string; $e115?: { $f115: number }[] } }
interface T116 { $a116: string; $b116: number; $c116?: { $d116: string; $e116?: { $f116: number }[] } }
interface T117 { $a117: string; $b117: number; $c117?: { $d117: string; $e117?: { $f117: number }[] } }
interface T118 { $a118: string; $b118: number; $c118?: { $d118: string; $e118?: { $f118: number }[] } }
interface T119 { $a119: string; $b119: number; $c119?: { $d119: string; $e119?: { $f119: number }[] } }
interface T120 { $a120: string; $b120: number; $c120?: { $d120: string; $e120?: { $f120: number }[] } }
interface T121 { $a121: string; $b121: number; $c121?: { $d121: string; $e121?: { $f121: number }[] } }
interface T122 { $a122: string; $b122: number; $c122?: { $d122: string; $e122?: { $f122: number }[] } }
interface T123 { $a123: string; $b123: number; $c123?: { $d123: string; $e123?: { $f123: number }[] } }
interface T124 { $a124: string; $b124: number; $c124?: { $d124: string; $e124?: { $f124: number }[] } }
interface T125 { $a125: string; $b125: number; $c125?: { $d125: string; $e125?: { $f125: number }[] } }
interface T126 { $a126: string; $b126: number; $c126?: { $d126: string; $e126?: { $f126: number }[] } }
interface T127 { $a127: string; $b127: number; $c127?: { $d127: string; $e127?: { $f127: number }[] } }
interface T128 { $a128: string; $b128: number; $c128?: { $d128: string; $e128?: { $f128: number }[] } }
interface T129 { $a129: string; $b129: number; $c129?: { $d129: string; $e129?: { $f129: number }[] } }
interface T130 { $a130: string; $b130: number; $c130?: { $d130: string; $e130?: { $f130: number }[] } }
interface T131 { $a131: string; $b131: number; $c131?: { $d131: string; $e131?: { $f131: number }[] } }
interface T132 { $a132: string; $b132: number; $c132?: { $d132: string; $e132?: { $f132: number }[] } }
interface T133 { $a133: string; $b133: number; $c133?: { $d133: string; $e133?: { $f133: number }[] } }
interface T134 { $a134: string; $b134: number; $c134?: { $d134: string; $e134?: { $f134: number }[] } }
interface T135 { $a135: string; $b135: number; $c135?: { $d135: string; $e135?: { $f135: number }[] } }
interface T136 { $a136: string; $b136: number; $c136?: { $d136: string; $e136?: { $f136: number }[] } }
interface T137 { $a137: string; $b137: number; $c137?: { $d137: string; $e137?: { $f137: number }[] } }
interface T138 { $a138: string; $b138: number; $c138?: { $d138: string; $e138?: { $f138: number }[] } }
interface T139 { $a139: string; $b139: number; $c139?: { $d139: string; $e139?: { $f139: number }[] } }
interface T140 { $a140: string; $b140: number; $c140?: { $d140: string; $e140?: { $f140: number }[] } }
interface T141 { $a141: string; $b141: number; $c141?: { $d141: string; $e141?: { $f141: number }[] } }
interface T142 { $a142: string; $b142: number; $c142?: { $d142: string; $e142?: { $f142: number }[] } }
interface T143 { $a143: string; $b143: number; $c143?: { $d143: string; $e143?: { $f143: number }[] } }
interface T144 { $a144: string; $b144: number; $c144?: { $d144: string; $e144?: { $f144: number }[] } }
interface T145 { $a145: string; $b145: number; $c145?: { $d145: string; $e145?: { $f145: number }[] } }
interface T146 { $a146: string; $b146: number; $c146?: { $d146: string; $e146?: { $f146: number }[] } }
interface T147 { $a147: string; $b147: number; $c147?: { $d147: string; $e147?: { $f147: number }[] } }
interface T148 { $a148: string; $b148: number; $c148?: { $d148: string; $e148?: { $f148: number }[] } }
interface T149 { $a149: string; $b149: number; $c149?: { $d149: string; $e149?: { $f149: number }[] } }
interface T150 { $a150: string; $b150: number; $c150?: { $d150: string; $e150?: { $f150: number }[] } }

type _BigUnion =
    | T001
    | T002
    | T003
    | T004
    | T005
    | T006
    | T007
    | T008
    | T009
    | T010
    | T011
    | T012
    | T013
    | T014
    | T015
    | T016
    | T017
    | T018
    | T019
    | T020
    | T021
    | T022
    | T023
    | T024
    | T025
    | T026
    | T027
    | T028
    | T029
    | T030
    | T031
    | T032
    | T033
    | T034
    | T035
    | T036
    | T037
    | T038
    | T039
    | T040
    | T041
    | T042
    | T043
    | T044
    | T045
    | T046
    | T047
    | T048
    | T049
    | T050
    | T051
    | T052
    | T053
    | T054
    | T055
    | T056
    | T057
    | T058
    | T059
    | T060
    | T061
    | T062
    | T063
    | T064
    | T065
    | T066
    | T067
    | T068
    | T069
    | T070
    | T071
    | T072
    | T073
    | T074
    | T075
    | T076
    | T077
    | T078
    | T079
    | T080
    | T081
    | T082
    | T083
    | T084
    | T085
    | T086
    | T087
    | T088
    | T089
    | T090
    | T091
    | T092
    | T093
    | T094
    | T095
    | T096
    | T097
    | T098
    | T099
    | T100
    | T101
    | T102
    | T103
    | T104
    | T105
    | T106
    | T107
    | T108
    | T109
    | T110
    | T111
    | T112
    | T113
    | T114
    | T115
    | T116
    | T117
    | T118
    | T119
    | T120
    | T121
    | T122
    | T123
    | T124
    | T125
    | T126
    | T127
    | T128
    | T129
    | T130
    | T131
    | T132
    | T133
    | T134
    | T135
    | T136
    | T137
    | T138
    | T139
    | T140
    | T141
    | T142
    | T143
    | T144
    | T145
    | T146
    | T147
    | T148
    | T149
    | T150
    | { $value: never };

// ===== prefix.ts =====
// Mirrors Forms' Common/Prefix/Prefix.ts machinery.
// Uses globals from helpers.d.ts: EmptyAsNever, IsAny, IsNever, ParamWithoutNever.

type BasicType = null | string | number | boolean | Date | Function | Blob;

type UnprefixedKey<TKey extends string> = TKey extends `$${infer K}` ? K : never;
type PrefixedKey<TKey extends string> = `$${TKey}`;

type Unprefixed<T> = T extends BasicType
    ? T
    : T extends Array<infer TItem>
      ? Array<Unprefixed<TItem>>
      : T extends {}
        ? string extends keyof T
            ? { [key: string]: Unprefixed<T[string]> }
            : {
                  // @ts-ignore
                  [key in keyof T as UnprefixedKey<key>]: Unprefixed<T[key]>;
              }
        : never;

type IPrefixEscapeSpec<T> = T extends BasicType
    ? never
    : T extends Array<infer TItem>
      ? IPrefixEscapeSpec<TItem> extends never
          ? never
          : [IPrefixEscapeSpec<TItem>]
      : T extends {}
        ? keyof T extends never
            ? never
            : string extends keyof T
              ? T extends Record<string, infer TValue>
                  ? EmptyAsNever<{ $isMap: true; $value: IPrefixEscapeSpec<TValue> }>
                  : never
              : EmptyAsNever<{
                    $fields: EmptyAsNever<{
                        [key in keyof T & string]: IPrefixEscapeSpec<T[key]>;
                    }>;
                }>
        : never;

type CompileTypeError<TError extends string, TErrorPrefix extends string> = TErrorPrefix extends ""
    ? [TError]
    : [`${TErrorPrefix}: ${TError}`];

type GetObjectSubFieldNames<T extends {}, TKey extends keyof T, TErrorPrefix extends string> = TKey extends
    | `$${any}`
    | `@${any}`
    ? GetFieldNamesRecursively<T[TKey], `${TErrorPrefix}.${TKey}`>
    : CompileTypeError<"only $-prefixed", TErrorPrefix>;

type GetFieldNamesRecursively<T, TErrorPrefix extends string> = IsAny<T> extends true
    ? CompileTypeError<"any not allowed", TErrorPrefix>
    : IsNever<T> extends true
      ? never
      : T extends T
        ? T extends BasicType
            ? never
            : T extends any[]
              ? GetFieldNamesRecursively<T[number], `${TErrorPrefix}[]`>
              : T extends {}
                ? IsNever<keyof T> extends true
                    ? CompileTypeError<"empty", TErrorPrefix>
                    : string extends keyof T
                      ? GetFieldNamesRecursively<T[keyof T], `${TErrorPrefix}.*`>
                      : keyof T | GetObjectSubFieldNames<T, keyof T, TErrorPrefix>
                : CompileTypeError<"not supported", TErrorPrefix>
        : never;

type GetPrefixTypeErrors<T> = Exclude<GetFieldNamesRecursively<T, "">, string>;

// Stand-in for IPrefixedServerTypes. Mirrors the real shape: Partial<wide union>.
type IPrefixedServerTypes = Partial<_BigUnion>;

type IsValidPrefixedType<
    T extends IPrefixedParam,
    TTarget = T,
    TError = CompileTypeError<"unknown", "">,
> = GetFieldNamesRecursively<T, ""> extends GetFieldNamesRecursively<IPrefixedServerTypes, ""> ? TTarget : TError;

type PrefixApiReturnType<T extends IPrefixedParam, TReturn = T> =
    IsNever<GetPrefixTypeErrors<T>> extends true
        ? IsValidPrefixedType<T, TReturn>
        : GetPrefixTypeErrors<T>;

type IEscaperParamOrError<T> = IsNever<GetPrefixTypeErrors<T>> extends true
    ? IsValidPrefixedType<T, ParamWithoutNever<IPrefixEscapeSpec<T>>>
    : [GetPrefixTypeErrors<T>];

type IPrefixedParam =
    | IPrefixedServerTypes
    | IPrefixedServerTypes[]
    | { $value: IPrefixedServerTypes }
    | { $value: IPrefixedServerTypes[] }
    | BasicType
    | BasicType[];

// ===== prefixer.ts =====
// Mirrors Forms' Common/Prefix/Prefixer.ts — `unprefix` const is the trigger surface.
// The assignment of (data: any, escapeSpecifier?: any) => any to the generic
// signature with PrefixApiReturnType<T, Unprefixed<T>> is what tsgo panics on.


declare const IS_DEBUG: boolean;

function getRecursiveConvertor(mapKey: (key: string) => string, prefixing: boolean) {
    return function convertor(data: any, escapeSpecifier?: any): any {
        if (IS_DEBUG && escapeSpecifier === null) {
            throw new Error("escape specifier should not be null.");
        }
        if (Array.isArray(data)) {
            return data.map(subData => convertor(subData, escapeSpecifier && escapeSpecifier[0]));
        }
        if (!(data instanceof Object) || data instanceof Blob || data instanceof Date || data instanceof Function) {
            return data;
        }
        const escSpec = (escapeSpecifier || {}) as { $isMap?: boolean; $value?: unknown; $fields?: Record<string, unknown> };
        const { $isMap, $value: valueEscaper, $fields } = escSpec;
        return Object.keys(data).reduce<Record<string, unknown>>((res, key) => {
            const mappedKey = $isMap ? key : mapKey(key);
            if (mappedKey) {
                const subEscaper = $isMap ? valueEscaper : $fields && $fields[prefixing ? mappedKey : key];
                res[mappedKey] = convertor(data[key], subEscaper);
            }
            return res;
        }, {});
    };
}

const stubMapKey = (key: string) => key;

// @ts-ignore
const prefix: <T extends IPrefixedParam>(
    serverData: Unprefixed<T>,
    ...escapeSpecifier: IEscaperParamOrError<T>
) => PrefixApiReturnType<T> = getRecursiveConvertor(stubMapKey, true);

// === The trigger surface ===
// Assigning the (data: any, escapeSpecifier?: any) => any return value to this
// generic signature provokes the panic in conjunction with the function expression
// in requestWithPrefix.ts that calls unprefix<TInput>(...) inside an object literal.
const unprefix: <T extends IPrefixedParam>(
    prefixedData: T,
    ...escapeSpecifier: IEscaperParamOrError<T>
) => PrefixApiReturnType<T, Unprefixed<T>> = getRecursiveConvertor(stubMapKey, false);

// ===== request.ts =====
// Mirrors <Common>/Utilities/Request.ts (ISendRequest) and <NeoCommon>/Utils/Request.ts (IRequest, ISendRequestEscapeSpecs, RequestInputTypeChecker)

interface ISendRequest<TResponse = never> {
    accessToken?: string;
    correlationId?: string;
    data?: {};
    defaultValueForTryoutMode?: TResponse;
    disableRetry?: boolean;
    headers?: { [headerKey: string]: string };
    method: "POST" | "GET" | "PUT" | "PATCH" | "DELETE";
    params?: string;
    qosEventName?: string;
    timeout?: number;
    url: string;
    withoutToken?: boolean;
    retryConfig?: { retryCount: number };
}

// === Uses GLOBAL Override<T, R> from helpers.d.ts with inline object-literal type ===
// Same shape as the Forms-internal NeoCommon/Utils/Request.ts.
type IRequest<TData, TResponse = never> = Override<
    ISendRequest<TResponse>,
    {
        data?: TData;
        qosEventName: string;
    }
>;

interface ISendRequestEscapeSpecs<TInput, TOutput> {
    $input: IPrefixEscapeSpec<TInput>;
    $output: IPrefixEscapeSpec<TOutput>;
}

type RequestInputTypeChecker<T> = [...EmptyObjectTypeCheck<T>, ...IsValidPrefixedTypeCheck<T>];
type EmptyObjectTypeCheck<T> = keyof T extends never ? [{ TInput_ShouldNotBeEmptyObject: true }] : [];
type IsValidPrefixedTypeCheck<T> = IsValidPrefixedType<T, [], [{ TInput_ShouldBeListedInPrefixedTypeList: true }]>;

// ===== common.ts =====
// Mirror of Common/FormData/Types/ServerTypes/Common.ts
type OverrideExistingFields<T extends {}, R> = Override<T, Pick<R, keyof T & keyof R>>;

// ===== questioninfo.ts =====
// Mirrors Forms' Common/FormData/Types/ServerTypes/Form/QuestionInfo.ts pattern:
// IUnionQuestionInfo is a wide intersection of Partial<...> types extending a base.
// In the real codebase the imports here pull in <Forms>/Pages/LightResponsePage/types/questions
// and <Common>/Enums/* — both essential to the panic. Kept inline here to be self-contained.

interface IBranchInfo { $TargetQuestionId?: string; $ToTheEnd?: boolean }
interface IVideo { $Url: string; $ShowWarningBar: boolean }
interface IImageSettings { $Layout: string; $Size: string; $NaturalWidth: number; $NaturalHeight: number }

interface IBaseQuestionInfo {
    $Point?: number;
    $ShowSubtitle?: boolean;
    $IsMathQuiz?: boolean;
    $BranchInfo?: IBranchInfo;
    $Multiline?: boolean;
    $Hidden?: boolean;
    $VideoInfo?: IVideo;
    $ImageSettings?: IImageSettings;
}

interface ITextValidation { $rule: number; $value: string }

interface ITextQuestionInfo extends IBaseQuestionInfo {
    $Validation?: ITextValidation;
    $Multiline: boolean;
}

interface IChoiceItem { $Description: string; $Index?: number; $BranchInfo?: IBranchInfo }

interface IChoiceQuestionInfo extends IBaseQuestionInfo {
    $AllowOtherAnswer: boolean;
    $BranchInfoForOtherChoice: IBranchInfo;
    $Choices: IChoiceItem[];
    $ChoiceType: number;
    $OptionDisplayStyle: string;
    $ShuffleOptions?: boolean;
}

interface IRatingQuestionInfo extends IBaseQuestionInfo { $Length: number; $RatingShape: string }
interface IRankingQuestionInfo extends IBaseQuestionInfo { $ShuffleOptions?: boolean }
interface INPSQuestionInfo extends IBaseQuestionInfo { $LeftDescription: string; $RightDescription: string }
interface IFileUploadQuestionInfo extends IBaseQuestionInfo {
    $MaxFileCount: number;
    $MaxFileSize: number;
    $FileTypes: { [t: string]: boolean };
}
interface IDateTimeQuestionInfo extends IBaseQuestionInfo { $Date: boolean; $Time: boolean }

// === Wide intersection — same shape as IUnionQuestionInfo in Forms ===
type IUnionQuestionInfo = IBaseQuestionInfo &
    Partial<ITextQuestionInfo> &
    Partial<IChoiceQuestionInfo> &
    Partial<IRatingQuestionInfo> &
    Partial<IRankingQuestionInfo> &
    Partial<INPSQuestionInfo> &
    Partial<IFileUploadQuestionInfo> &
    Partial<IDateTimeQuestionInfo>;

type IQuestionInfo = IUnionQuestionInfo;

// ===== choice.ts =====
interface IChoice { $id: string; $value: string; }
interface IRawChoice { $id: string; $value: string; }

// ===== image.ts =====
// Self-contained — no Forms refs. Real codebase Image.ts has more enums and union types.
interface IRawImage {
    $altText: string;
    $contentType: string;
    $fileIdentifier: string;
    $originalFileName: string;
    $resourceId: string;
    $resourceUrl: string;
    $height: number;
    $width: number;
    $size: "Small" | "Large";
}

interface IImage extends IRawImage {}

interface IOpenTypeImage {
    $width: number;
    $height: number;
    $left: number;
    $top: number;
    $altText?: string;
}

interface IOpenTypeImageDictionary {
    [quid: string]: IOpenTypeImage;
}

// ===== question.ts =====
// Mirrors Forms' Common/FormData/Types/ServerTypes/Form/Question.ts —
// the file that bisection identified as the trigger.
//
// The construct of interest is `OverrideExistingFields<IRawQuestion, IQuestionOverriddenProperties<T>>`
// where T defaults to `IUnionQuestionInfo` (a wide intersection of Partial<...> types).


interface IRawQuestion {
    $allowMultipleValues?: boolean;
    $choices?: IRawChoice[];
    $defaultValue?: string;
    $groupId?: string;
    $id: string;
    $image?: IRawImage;
    $imageDictionary?: string;
    $insightsInfo?: string;
    $isFromSuggestion?: boolean;
    $isQuiz?: boolean;
    $modifiedDate?: string;
    $order?: number;
    $questionInfo: string;
    $required?: boolean;
    $shuffleOrder?: number;
    $subtitle?: string;
    $title?: string;
    $trackingId?: string;
    $type?: number;
}

interface IQuestionOverriddenProperties<T> {
    $choices?: IChoice[];
    $image?: IImage;
    $imageDictionary?: IOpenTypeImageDictionary;
    $questionInfo: T;
}

// === The construct that the panic centers on ===
type IQuestion<T = IQuestionInfo> = OverrideExistingFields<IRawQuestion, IQuestionOverriddenProperties<T>>;

// ===== requestWithPrefix.ts =====
// Mirrors <NeoCommon>/Utils/RequestWithPrefix.ts.
// The arrow-function returned by sendRequestWithPrefix is the panic site:
// inside its body, `request<Unprefixed<TOutput>>({ ...requestParams, data: unprefix<TInput>(...) }).then(...)`
// is the property-access-on-call-with-object-literal pattern that the panic stack ends at.

function sendRequestWithPrefix(request: <T>(param: ISendRequest<T>) => Promise<T>) {
    return <TInput extends IPrefixedParam = never, TOutput extends IPrefixedParam = never>(
        requestParams: IRequest<TInput, Unprefixed<TOutput>>,
        escapeSpecs: WithoutNever<ISendRequestEscapeSpecs<TInput, TOutput>>,
        ...compileTimeChecker: RequestInputTypeChecker<TInput>
    ): Promise<PrefixApiReturnType<TOutput>> => {
        const { $input, $output } = escapeSpecs as ISendRequestEscapeSpecs<TInput, TOutput>;
        const inputEscapeSpec = ($input ? [$input] : []) as IEscaperParamOrError<TInput>;

        // @ts-ignore
        return request<Unprefixed<TOutput>>({
            ...requestParams,
            data: unprefix<TInput>(requestParams.data, ...inputEscapeSpec),
        }).then((data: Unprefixed<TOutput>) => {
            const outputEscapeSpec = ($output ? [$output] : []) as IEscaperParamOrError<TOutput>;
            return prefix<TOutput>(data, ...outputEscapeSpec);
        });
    };
}
