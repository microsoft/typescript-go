currentDirectory::/home/src/workspaces/project
useCaseSensitiveFileNames::true
Input::
//// [/home/src/workspaces/project/main.ts] *new* 
import { Container, Nullable } from "./types";
const c: Container<number> = { value: 42, map: (fn) => ({ value: fn(42), map: c.map }) };
const n: Nullable<string> = "hello";
//// [/home/src/workspaces/project/tsconfig.json] *new* 
{
    "compilerOptions": {
        "strict": true,
        "noEmit": true
    }
}
//// [/home/src/workspaces/project/types.ts] *new* 
export interface Container<T> {
    value: T;
    map<U>(fn: (x: T) => U): Container<U>;
}
export type Nullable<T> = T | null | undefined;

tsgo --generateTrace /home/src/workspaces/project/trace
ExitStatus:: DiagnosticsPresent_OutputsSkipped
Output::
[96mmain.ts[0m:[93m2[0m:[93m74[0m - [91merror[0m[90m TS2322: [0mType '<U>(fn: (x: number) => U) => Container<U>' is not assignable to type '<U>(fn: (x: U) => U) => Container<U>'.
  Types of parameters 'fn' and 'fn' are incompatible.
    Types of parameters 'x' and 'x' are incompatible.
      Type 'number' is not assignable to type 'U'.
        'U' could be instantiated with an arbitrary type which could be unrelated to 'number'.

[7m2[0m const c: Container<number> = { value: 42, map: (fn) => ({ value: fn(42), map: c.map }) };
[7m [0m [91m                                                                         ~~~[0m

  [96mtypes.ts[0m:[93m3[0m:[93m5[0m - The expected type comes from property 'map' which is declared here on type 'Container<U>'
    [7m3[0m     map<U>(fn: (x: T) => U): Container<U>;
    [7m [0m [96m    ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~[0m


Found 1 error in main.ts[90m:2[0m

//// [/home/src/tslibs/TS/Lib/lib.es2025.full.d.ts] *Lib*
/// <reference no-default-lib="true"/>
interface Boolean {}
interface Function {}
interface CallableFunction {}
interface NewableFunction {}
interface IArguments {}
interface Number { toExponential: any; }
interface Object {}
interface RegExp {}
interface String { charAt: any; }
interface Array<T> { length: number; [n: number]: T; }
interface ReadonlyArray<T> {}
interface SymbolConstructor {
    (desc?: string | number): symbol;
    for(name: string): symbol;
    readonly toStringTag: symbol;
}
declare var Symbol: SymbolConstructor;
interface Symbol {
    readonly [Symbol.toStringTag]: string;
}
declare const console: { log(msg: any): void; };
//// [/home/src/workspaces/project/trace/legend.json] *new* 
[{"configFilePath":"/home/src/workspaces/project/tsconfig.json","typesPath":"/home/src/workspaces/project/trace/types_0.json"},{"configFilePath":"/home/src/workspaces/project/tsconfig.json","typesPath":"/home/src/workspaces/project/trace/types_1.json"},{"configFilePath":"/home/src/workspaces/project/tsconfig.json","typesPath":"/home/src/workspaces/project/trace/types_2.json"}]
//// [/home/src/workspaces/project/trace/types_0.json] *new* 
[{"id":1,"intrinsicName":"any","isTuple":false,"flags":["Any"]},
{"id":2,"intrinsicName":"any","isTuple":false,"flags":["Any"]},
{"id":3,"intrinsicName":"any","isTuple":false,"flags":["Any"]},
{"id":4,"intrinsicName":"any","isTuple":false,"flags":["Any"]},
{"id":5,"intrinsicName":"any","isTuple":false,"flags":["Any"]},
{"id":6,"intrinsicName":"unresolved","isTuple":false,"flags":["Any"]},
{"id":7,"intrinsicName":"any","isTuple":false,"flags":["Any"]},
{"id":8,"intrinsicName":"intrinsic","isTuple":false,"flags":["Any"]},
{"id":9,"intrinsicName":"unknown","isTuple":false,"flags":["Unknown"]},
{"id":10,"intrinsicName":"undefined","isTuple":false,"flags":["Undefined"]},
{"id":11,"intrinsicName":"undefined","isTuple":false,"flags":["Undefined"]},
{"id":12,"intrinsicName":"undefined","isTuple":false,"flags":["Undefined"]},
{"id":13,"intrinsicName":"null","isTuple":false,"flags":["Null"]},
{"id":14,"intrinsicName":"string","isTuple":false,"flags":["String"]},
{"id":15,"intrinsicName":"number","isTuple":false,"flags":["Number"]},
{"id":16,"intrinsicName":"bigint","isTuple":false,"flags":["BigInt"]},
{"id":17,"isTuple":false,"flags":["BooleanLiteral"],"display":"false"},
{"id":18,"isTuple":false,"flags":["BooleanLiteral"],"display":"false"},
{"id":19,"isTuple":false,"flags":["BooleanLiteral"],"display":"true"},
{"id":20,"isTuple":false,"flags":["BooleanLiteral"],"display":"true"},
{"id":21,"isTuple":false,"unionTypes":[17,19],"flags":["Boolean","Union"]},
{"id":22,"intrinsicName":"symbol","isTuple":false,"flags":["ESSymbol"]},
{"id":23,"intrinsicName":"void","isTuple":false,"flags":["Void"]},
{"id":24,"intrinsicName":"never","isTuple":false,"flags":["Never"]},
{"id":25,"intrinsicName":"never","isTuple":false,"flags":["Never"]},
{"id":26,"intrinsicName":"never","isTuple":false,"flags":["Never"]},
{"id":27,"intrinsicName":"never","isTuple":false,"flags":["Never"]},
{"id":28,"intrinsicName":"object","isTuple":false,"flags":["NonPrimitive"]},
{"id":29,"isTuple":false,"unionTypes":[14,15],"flags":["Union"]},
{"id":30,"isTuple":false,"unionTypes":[14,15,22],"flags":["Union"]},
{"id":31,"isTuple":false,"unionTypes":[15,16],"flags":["Union"]},
{"id":32,"isTuple":false,"flags":["TemplateLiteral"]},
{"id":33,"isTuple":false,"unionTypes":[10,13,14,15,16,17,19],"flags":["Union"]},
{"id":34,"intrinsicName":"never","isTuple":false,"flags":["Never"]},
{"id":35,"isTuple":false,"flags":["Object"],"display":"{}"},
{"id":36,"isTuple":false,"flags":["Object"],"display":"{}"},
{"id":37,"isTuple":false,"flags":["Object"],"display":"{}"},
{"id":38,"symbolName":"__type","isTuple":false,"flags":["Object"],"display":"{}"},
{"id":39,"isTuple":false,"flags":["Object"],"display":"{}"},
{"id":40,"isTuple":false,"unionTypes":[10,13,39],"flags":["Union"]},
{"id":41,"isTuple":false,"flags":["Object"],"display":"{}"},
{"id":42,"isTuple":false,"flags":["Object"],"display":"{}"},
{"id":43,"isTuple":false,"flags":["Object"],"display":"{}"},
{"id":44,"isTuple":false,"flags":["Object"],"display":"{}"},
{"id":45,"isTuple":false,"flags":["Object"],"display":"{}"},
{"id":46,"isTuple":false,"flags":["TypeParameter"]},
{"id":47,"isTuple":false,"flags":["TypeParameter"]},
{"id":48,"isTuple":false,"flags":["TypeParameter"]},
{"id":49,"isTuple":false,"flags":["TypeParameter"]},
{"id":50,"isTuple":false,"flags":["TypeParameter"]},
{"id":51,"isTuple":false,"flags":["StringLiteral"],"display":"\"\""},
{"id":52,"isTuple":false,"flags":["NumberLiteral"],"display":"0"},
{"id":53,"isTuple":false,"flags":["BigIntLiteral"],"display":"0n"},
{"id":54,"isTuple":false,"flags":["StringLiteral"],"display":"\"bigint\""},
{"id":55,"isTuple":false,"flags":["StringLiteral"],"display":"\"boolean\""},
{"id":56,"isTuple":false,"flags":["StringLiteral"],"display":"\"function\""},
{"id":57,"isTuple":false,"flags":["StringLiteral"],"display":"\"number\""},
{"id":58,"isTuple":false,"flags":["StringLiteral"],"display":"\"object\""},
{"id":59,"isTuple":false,"flags":["StringLiteral"],"display":"\"string\""},
{"id":60,"isTuple":false,"flags":["StringLiteral"],"display":"\"symbol\""},
{"id":61,"isTuple":false,"flags":["StringLiteral"],"display":"\"undefined\""},
{"id":62,"isTuple":false,"unionTypes":[54,55,56,57,58,59,60,61],"flags":["Union"]},
{"id":63,"symbolName":"IArguments","isTuple":false,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":5,"character":29},"end":{"line":6,"character":24}},"flags":["Object"]},
{"id":64,"symbolName":"globalThis","isTuple":false,"flags":["Object"],"display":"typeof globalThis"},
{"id":65,"symbolName":"Array","isTuple":false,"instantiatedType":65,"typeArguments":[66],"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":10,"character":34},"end":{"line":11,"character":55}},"flags":["Object"]},
{"id":66,"symbolName":"T","isTuple":false,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":11,"character":17},"end":{"line":11,"character":18}},"flags":["TypeParameter"]},
{"id":67,"symbolName":"Array","isTuple":false,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":10,"character":34},"end":{"line":11,"character":55}},"flags":["TypeParameter"]},
{"id":68,"symbolName":"Object","isTuple":false,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":7,"character":41},"end":{"line":8,"character":20}},"flags":["Object"]},
{"id":69,"symbolName":"Function","isTuple":false,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":2,"character":21},"end":{"line":3,"character":22}},"flags":["Object"]},
{"id":70,"symbolName":"CallableFunction","isTuple":false,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":3,"character":22},"end":{"line":4,"character":30}},"flags":["Object"]},
{"id":71,"symbolName":"NewableFunction","isTuple":false,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":4,"character":30},"end":{"line":5,"character":29}},"flags":["Object"]},
{"id":72,"symbolName":"String","isTuple":false,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":9,"character":20},"end":{"line":10,"character":34}},"flags":["Object"]},
{"id":73,"symbolName":"Number","isTuple":false,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":6,"character":24},"end":{"line":7,"character":41}},"flags":["Object"]},
{"id":74,"symbolName":"Boolean","isTuple":false,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":1,"character":1},"end":{"line":2,"character":21}},"flags":["Object"]},
{"id":75,"symbolName":"RegExp","isTuple":false,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":8,"character":20},"end":{"line":9,"character":20}},"flags":["Object"]},
{"id":76,"symbolName":"Array","isTuple":false,"instantiatedType":65,"typeArguments":[1],"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":10,"character":34},"end":{"line":11,"character":55}},"flags":["Object"]},
{"id":77,"symbolName":"Array","isTuple":false,"instantiatedType":65,"typeArguments":[2],"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":10,"character":34},"end":{"line":11,"character":55}},"flags":["Object"]},
{"id":78,"symbolName":"ReadonlyArray","isTuple":false,"instantiatedType":78,"typeArguments":[79],"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":11,"character":55},"end":{"line":12,"character":30}},"flags":["Object"]},
{"id":79,"symbolName":"T","isTuple":false,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":12,"character":25},"end":{"line":12,"character":26}},"flags":["TypeParameter"]},
{"id":80,"symbolName":"ReadonlyArray","isTuple":false,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":11,"character":55},"end":{"line":12,"character":30}},"flags":["TypeParameter"]},
{"id":81,"symbolName":"ReadonlyArray","isTuple":false,"instantiatedType":78,"typeArguments":[1],"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":11,"character":55},"end":{"line":12,"character":30}},"flags":["Object"]},
{"id":82,"symbolName":"Array","isTuple":false,"instantiatedType":65,"typeArguments":[66,67],"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":10,"character":34},"end":{"line":11,"character":55}},"flags":["Object"]},
{"id":83,"isTuple":false,"flags":["StringLiteral"],"display":"\"length\""},
{"id":84,"symbolName":"ReadonlyArray","isTuple":false,"instantiatedType":78,"typeArguments":[79,80],"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":11,"character":55},"end":{"line":12,"character":30}},"flags":["Object"]},
{"id":85,"symbolName":"SymbolConstructor","isTuple":false,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":12,"character":30},"end":{"line":17,"character":2}},"flags":["Object"]},
{"id":86,"isTuple":false,"unionTypes":[10,14,15],"flags":["Union"]},
{"id":87,"symbolName":"toStringTag","isTuple":false,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":15,"character":31},"end":{"line":16,"character":34}},"flags":["UniqueESSymbol"]},
{"id":88,"symbolName":"Symbol","isTuple":false,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":18,"character":12},"end":{"line":18,"character":38}},"flags":["Object"]},
{"id":89,"symbolName":"__type","isTuple":false,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":22,"character":23},"end":{"line":22,"character":48}},"flags":["Object"],"display":"{ log(msg: any): void; }"}]

//// [/home/src/workspaces/project/trace/types_1.json] *new* 
[{"id":1,"intrinsicName":"any","isTuple":false,"flags":["Any"]},
{"id":2,"intrinsicName":"any","isTuple":false,"flags":["Any"]},
{"id":3,"intrinsicName":"any","isTuple":false,"flags":["Any"]},
{"id":4,"intrinsicName":"any","isTuple":false,"flags":["Any"]},
{"id":5,"intrinsicName":"any","isTuple":false,"flags":["Any"]},
{"id":6,"intrinsicName":"unresolved","isTuple":false,"flags":["Any"]},
{"id":7,"intrinsicName":"any","isTuple":false,"flags":["Any"]},
{"id":8,"intrinsicName":"intrinsic","isTuple":false,"flags":["Any"]},
{"id":9,"intrinsicName":"unknown","isTuple":false,"flags":["Unknown"]},
{"id":10,"intrinsicName":"undefined","isTuple":false,"flags":["Undefined"]},
{"id":11,"intrinsicName":"undefined","isTuple":false,"flags":["Undefined"]},
{"id":12,"intrinsicName":"undefined","isTuple":false,"flags":["Undefined"]},
{"id":13,"intrinsicName":"null","isTuple":false,"flags":["Null"]},
{"id":14,"intrinsicName":"string","isTuple":false,"flags":["String"]},
{"id":15,"intrinsicName":"number","isTuple":false,"flags":["Number"]},
{"id":16,"intrinsicName":"bigint","isTuple":false,"flags":["BigInt"]},
{"id":17,"isTuple":false,"flags":["BooleanLiteral"],"display":"false"},
{"id":18,"isTuple":false,"flags":["BooleanLiteral"],"display":"false"},
{"id":19,"isTuple":false,"flags":["BooleanLiteral"],"display":"true"},
{"id":20,"isTuple":false,"flags":["BooleanLiteral"],"display":"true"},
{"id":21,"isTuple":false,"unionTypes":[17,19],"flags":["Boolean","Union"]},
{"id":22,"intrinsicName":"symbol","isTuple":false,"flags":["ESSymbol"]},
{"id":23,"intrinsicName":"void","isTuple":false,"flags":["Void"]},
{"id":24,"intrinsicName":"never","isTuple":false,"flags":["Never"]},
{"id":25,"intrinsicName":"never","isTuple":false,"flags":["Never"]},
{"id":26,"intrinsicName":"never","isTuple":false,"flags":["Never"]},
{"id":27,"intrinsicName":"never","isTuple":false,"flags":["Never"]},
{"id":28,"intrinsicName":"object","isTuple":false,"flags":["NonPrimitive"]},
{"id":29,"isTuple":false,"unionTypes":[14,15],"flags":["Union"]},
{"id":30,"isTuple":false,"unionTypes":[14,15,22],"flags":["Union"]},
{"id":31,"isTuple":false,"unionTypes":[15,16],"flags":["Union"]},
{"id":32,"isTuple":false,"flags":["TemplateLiteral"]},
{"id":33,"isTuple":false,"unionTypes":[10,13,14,15,16,17,19],"flags":["Union"]},
{"id":34,"intrinsicName":"never","isTuple":false,"flags":["Never"]},
{"id":35,"isTuple":false,"flags":["Object"],"display":"{}"},
{"id":36,"isTuple":false,"flags":["Object"],"display":"{}"},
{"id":37,"isTuple":false,"flags":["Object"],"display":"{}"},
{"id":38,"symbolName":"__type","isTuple":false,"flags":["Object"],"display":"{}"},
{"id":39,"isTuple":false,"flags":["Object"],"display":"{}"},
{"id":40,"isTuple":false,"unionTypes":[10,13,39],"flags":["Union"]},
{"id":41,"isTuple":false,"flags":["Object"],"display":"{}"},
{"id":42,"isTuple":false,"flags":["Object"],"display":"{}"},
{"id":43,"isTuple":false,"flags":["Object"],"display":"{}"},
{"id":44,"isTuple":false,"flags":["Object"],"display":"{}"},
{"id":45,"isTuple":false,"flags":["Object"],"display":"{}"},
{"id":46,"isTuple":false,"flags":["TypeParameter"]},
{"id":47,"isTuple":false,"flags":["TypeParameter"]},
{"id":48,"isTuple":false,"flags":["TypeParameter"]},
{"id":49,"isTuple":false,"flags":["TypeParameter"]},
{"id":50,"isTuple":false,"flags":["TypeParameter"]},
{"id":51,"isTuple":false,"flags":["StringLiteral"],"display":"\"\""},
{"id":52,"isTuple":false,"flags":["NumberLiteral"],"display":"0"},
{"id":53,"isTuple":false,"flags":["BigIntLiteral"],"display":"0n"},
{"id":54,"isTuple":false,"flags":["StringLiteral"],"display":"\"bigint\""},
{"id":55,"isTuple":false,"flags":["StringLiteral"],"display":"\"boolean\""},
{"id":56,"isTuple":false,"flags":["StringLiteral"],"display":"\"function\""},
{"id":57,"isTuple":false,"flags":["StringLiteral"],"display":"\"number\""},
{"id":58,"isTuple":false,"flags":["StringLiteral"],"display":"\"object\""},
{"id":59,"isTuple":false,"flags":["StringLiteral"],"display":"\"string\""},
{"id":60,"isTuple":false,"flags":["StringLiteral"],"display":"\"symbol\""},
{"id":61,"isTuple":false,"flags":["StringLiteral"],"display":"\"undefined\""},
{"id":62,"isTuple":false,"unionTypes":[54,55,56,57,58,59,60,61],"flags":["Union"]},
{"id":63,"symbolName":"IArguments","isTuple":false,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":5,"character":29},"end":{"line":6,"character":24}},"flags":["Object"]},
{"id":64,"symbolName":"globalThis","isTuple":false,"flags":["Object"],"display":"typeof globalThis"},
{"id":65,"symbolName":"Array","isTuple":false,"instantiatedType":65,"typeArguments":[66],"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":10,"character":34},"end":{"line":11,"character":55}},"flags":["Object"]},
{"id":66,"symbolName":"T","isTuple":false,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":11,"character":17},"end":{"line":11,"character":18}},"flags":["TypeParameter"]},
{"id":67,"symbolName":"Array","isTuple":false,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":10,"character":34},"end":{"line":11,"character":55}},"flags":["TypeParameter"]},
{"id":68,"symbolName":"Object","isTuple":false,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":7,"character":41},"end":{"line":8,"character":20}},"flags":["Object"]},
{"id":69,"symbolName":"Function","isTuple":false,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":2,"character":21},"end":{"line":3,"character":22}},"flags":["Object"]},
{"id":70,"symbolName":"CallableFunction","isTuple":false,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":3,"character":22},"end":{"line":4,"character":30}},"flags":["Object"]},
{"id":71,"symbolName":"NewableFunction","isTuple":false,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":4,"character":30},"end":{"line":5,"character":29}},"flags":["Object"]},
{"id":72,"symbolName":"String","isTuple":false,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":9,"character":20},"end":{"line":10,"character":34}},"flags":["Object"]},
{"id":73,"symbolName":"Number","isTuple":false,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":6,"character":24},"end":{"line":7,"character":41}},"flags":["Object"]},
{"id":74,"symbolName":"Boolean","isTuple":false,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":1,"character":1},"end":{"line":2,"character":21}},"flags":["Object"]},
{"id":75,"symbolName":"RegExp","isTuple":false,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":8,"character":20},"end":{"line":9,"character":20}},"flags":["Object"]},
{"id":76,"symbolName":"Array","isTuple":false,"instantiatedType":65,"typeArguments":[1],"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":10,"character":34},"end":{"line":11,"character":55}},"flags":["Object"]},
{"id":77,"symbolName":"Array","isTuple":false,"instantiatedType":65,"typeArguments":[2],"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":10,"character":34},"end":{"line":11,"character":55}},"flags":["Object"]},
{"id":78,"symbolName":"ReadonlyArray","isTuple":false,"instantiatedType":78,"typeArguments":[79],"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":11,"character":55},"end":{"line":12,"character":30}},"flags":["Object"]},
{"id":79,"symbolName":"T","isTuple":false,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":12,"character":25},"end":{"line":12,"character":26}},"flags":["TypeParameter"]},
{"id":80,"symbolName":"ReadonlyArray","isTuple":false,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":11,"character":55},"end":{"line":12,"character":30}},"flags":["TypeParameter"]},
{"id":81,"symbolName":"ReadonlyArray","isTuple":false,"instantiatedType":78,"typeArguments":[1],"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":11,"character":55},"end":{"line":12,"character":30}},"flags":["Object"]},
{"id":82,"symbolName":"T","isTuple":false,"firstDeclaration":{"path":"/home/src/workspaces/project/types.ts","start":{"line":1,"character":28},"end":{"line":1,"character":29}},"flags":["TypeParameter"]},
{"id":83,"symbolName":"Container","isTuple":false,"instantiatedType":83,"typeArguments":[82],"firstDeclaration":{"path":"/home/src/workspaces/project/types.ts","start":{"line":1,"character":1},"end":{"line":4,"character":2}},"flags":["Object"]},
{"id":84,"symbolName":"Container","isTuple":false,"firstDeclaration":{"path":"/home/src/workspaces/project/types.ts","start":{"line":1,"character":1},"end":{"line":4,"character":2}},"flags":["TypeParameter"]},
{"id":85,"symbolName":"Container","isTuple":false,"instantiatedType":83,"typeArguments":[82,84],"firstDeclaration":{"path":"/home/src/workspaces/project/types.ts","start":{"line":1,"character":1},"end":{"line":4,"character":2}},"flags":["Object"]},
{"id":86,"symbolName":"U","isTuple":false,"firstDeclaration":{"path":"/home/src/workspaces/project/types.ts","start":{"line":3,"character":9},"end":{"line":3,"character":10}},"flags":["TypeParameter"]},
{"id":87,"symbolName":"__type","isTuple":false,"firstDeclaration":{"path":"/home/src/workspaces/project/types.ts","start":{"line":3,"character":15},"end":{"line":3,"character":27}},"flags":["Object"],"display":"(x: T) => U"},
{"id":88,"symbolName":"Container","isTuple":false,"instantiatedType":83,"typeArguments":[86],"firstDeclaration":{"path":"/home/src/workspaces/project/types.ts","start":{"line":1,"character":1},"end":{"line":4,"character":2}},"flags":["Object"]},
{"id":89,"symbolName":"T","isTuple":false,"firstDeclaration":{"path":"/home/src/workspaces/project/types.ts","start":{"line":5,"character":22},"end":{"line":5,"character":23}},"flags":["TypeParameter"]},
{"id":90,"symbolName":"Nullable","isTuple":false,"unionTypes":[10,13,89],"aliasTypeArguments":[89],"flags":["Union"]}]

//// [/home/src/workspaces/project/trace/types_2.json] *new* 
[{"id":1,"intrinsicName":"any","isTuple":false,"flags":["Any"]},
{"id":2,"intrinsicName":"any","isTuple":false,"flags":["Any"]},
{"id":3,"intrinsicName":"any","isTuple":false,"flags":["Any"]},
{"id":4,"intrinsicName":"any","isTuple":false,"flags":["Any"]},
{"id":5,"intrinsicName":"any","isTuple":false,"flags":["Any"]},
{"id":6,"intrinsicName":"unresolved","isTuple":false,"flags":["Any"]},
{"id":7,"intrinsicName":"any","isTuple":false,"flags":["Any"]},
{"id":8,"intrinsicName":"intrinsic","isTuple":false,"flags":["Any"]},
{"id":9,"intrinsicName":"unknown","isTuple":false,"flags":["Unknown"]},
{"id":10,"intrinsicName":"undefined","isTuple":false,"flags":["Undefined"]},
{"id":11,"intrinsicName":"undefined","isTuple":false,"flags":["Undefined"]},
{"id":12,"intrinsicName":"undefined","isTuple":false,"flags":["Undefined"]},
{"id":13,"intrinsicName":"null","isTuple":false,"flags":["Null"]},
{"id":14,"intrinsicName":"string","isTuple":false,"flags":["String"]},
{"id":15,"intrinsicName":"number","isTuple":false,"flags":["Number"]},
{"id":16,"intrinsicName":"bigint","isTuple":false,"flags":["BigInt"]},
{"id":17,"isTuple":false,"flags":["BooleanLiteral"],"display":"false"},
{"id":18,"isTuple":false,"flags":["BooleanLiteral"],"display":"false"},
{"id":19,"isTuple":false,"flags":["BooleanLiteral"],"display":"true"},
{"id":20,"isTuple":false,"flags":["BooleanLiteral"],"display":"true"},
{"id":21,"isTuple":false,"unionTypes":[17,19],"flags":["Boolean","Union"]},
{"id":22,"intrinsicName":"symbol","isTuple":false,"flags":["ESSymbol"]},
{"id":23,"intrinsicName":"void","isTuple":false,"flags":["Void"]},
{"id":24,"intrinsicName":"never","isTuple":false,"flags":["Never"]},
{"id":25,"intrinsicName":"never","isTuple":false,"flags":["Never"]},
{"id":26,"intrinsicName":"never","isTuple":false,"flags":["Never"]},
{"id":27,"intrinsicName":"never","isTuple":false,"flags":["Never"]},
{"id":28,"intrinsicName":"object","isTuple":false,"flags":["NonPrimitive"]},
{"id":29,"isTuple":false,"unionTypes":[14,15],"flags":["Union"]},
{"id":30,"isTuple":false,"unionTypes":[14,15,22],"flags":["Union"]},
{"id":31,"isTuple":false,"unionTypes":[15,16],"flags":["Union"]},
{"id":32,"isTuple":false,"flags":["TemplateLiteral"]},
{"id":33,"isTuple":false,"unionTypes":[10,13,14,15,16,17,19],"flags":["Union"]},
{"id":34,"intrinsicName":"never","isTuple":false,"flags":["Never"]},
{"id":35,"isTuple":false,"flags":["Object"],"display":"{}"},
{"id":36,"isTuple":false,"flags":["Object"],"display":"{}"},
{"id":37,"isTuple":false,"flags":["Object"],"display":"{}"},
{"id":38,"symbolName":"__type","isTuple":false,"flags":["Object"],"display":"{}"},
{"id":39,"isTuple":false,"flags":["Object"],"display":"{}"},
{"id":40,"isTuple":false,"unionTypes":[10,13,39],"flags":["Union"]},
{"id":41,"isTuple":false,"flags":["Object"],"display":"{}"},
{"id":42,"isTuple":false,"flags":["Object"],"display":"{}"},
{"id":43,"isTuple":false,"flags":["Object"],"display":"{}"},
{"id":44,"isTuple":false,"flags":["Object"],"display":"{}"},
{"id":45,"isTuple":false,"flags":["Object"],"display":"{}"},
{"id":46,"isTuple":false,"flags":["TypeParameter"]},
{"id":47,"isTuple":false,"flags":["TypeParameter"]},
{"id":48,"isTuple":false,"flags":["TypeParameter"]},
{"id":49,"isTuple":false,"flags":["TypeParameter"]},
{"id":50,"isTuple":false,"flags":["TypeParameter"]},
{"id":51,"isTuple":false,"flags":["StringLiteral"],"display":"\"\""},
{"id":52,"isTuple":false,"flags":["NumberLiteral"],"display":"0"},
{"id":53,"isTuple":false,"flags":["BigIntLiteral"],"display":"0n"},
{"id":54,"isTuple":false,"flags":["StringLiteral"],"display":"\"bigint\""},
{"id":55,"isTuple":false,"flags":["StringLiteral"],"display":"\"boolean\""},
{"id":56,"isTuple":false,"flags":["StringLiteral"],"display":"\"function\""},
{"id":57,"isTuple":false,"flags":["StringLiteral"],"display":"\"number\""},
{"id":58,"isTuple":false,"flags":["StringLiteral"],"display":"\"object\""},
{"id":59,"isTuple":false,"flags":["StringLiteral"],"display":"\"string\""},
{"id":60,"isTuple":false,"flags":["StringLiteral"],"display":"\"symbol\""},
{"id":61,"isTuple":false,"flags":["StringLiteral"],"display":"\"undefined\""},
{"id":62,"isTuple":false,"unionTypes":[54,55,56,57,58,59,60,61],"flags":["Union"]},
{"id":63,"symbolName":"IArguments","isTuple":false,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":5,"character":29},"end":{"line":6,"character":24}},"flags":["Object"]},
{"id":64,"symbolName":"globalThis","isTuple":false,"flags":["Object"],"display":"typeof globalThis"},
{"id":65,"symbolName":"Array","isTuple":false,"instantiatedType":65,"typeArguments":[66],"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":10,"character":34},"end":{"line":11,"character":55}},"flags":["Object"]},
{"id":66,"symbolName":"T","isTuple":false,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":11,"character":17},"end":{"line":11,"character":18}},"flags":["TypeParameter"]},
{"id":67,"symbolName":"Array","isTuple":false,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":10,"character":34},"end":{"line":11,"character":55}},"flags":["TypeParameter"]},
{"id":68,"symbolName":"Object","isTuple":false,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":7,"character":41},"end":{"line":8,"character":20}},"flags":["Object"]},
{"id":69,"symbolName":"Function","isTuple":false,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":2,"character":21},"end":{"line":3,"character":22}},"flags":["Object"]},
{"id":70,"symbolName":"CallableFunction","isTuple":false,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":3,"character":22},"end":{"line":4,"character":30}},"flags":["Object"]},
{"id":71,"symbolName":"NewableFunction","isTuple":false,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":4,"character":30},"end":{"line":5,"character":29}},"flags":["Object"]},
{"id":72,"symbolName":"String","isTuple":false,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":9,"character":20},"end":{"line":10,"character":34}},"flags":["Object"]},
{"id":73,"symbolName":"Number","isTuple":false,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":6,"character":24},"end":{"line":7,"character":41}},"flags":["Object"]},
{"id":74,"symbolName":"Boolean","isTuple":false,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":1,"character":1},"end":{"line":2,"character":21}},"flags":["Object"]},
{"id":75,"symbolName":"RegExp","isTuple":false,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":8,"character":20},"end":{"line":9,"character":20}},"flags":["Object"]},
{"id":76,"symbolName":"Array","isTuple":false,"instantiatedType":65,"typeArguments":[1],"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":10,"character":34},"end":{"line":11,"character":55}},"flags":["Object"]},
{"id":77,"symbolName":"Array","isTuple":false,"instantiatedType":65,"typeArguments":[2],"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":10,"character":34},"end":{"line":11,"character":55}},"flags":["Object"]},
{"id":78,"symbolName":"ReadonlyArray","isTuple":false,"instantiatedType":78,"typeArguments":[79],"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":11,"character":55},"end":{"line":12,"character":30}},"flags":["Object"]},
{"id":79,"symbolName":"T","isTuple":false,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":12,"character":25},"end":{"line":12,"character":26}},"flags":["TypeParameter"]},
{"id":80,"symbolName":"ReadonlyArray","isTuple":false,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":11,"character":55},"end":{"line":12,"character":30}},"flags":["TypeParameter"]},
{"id":81,"symbolName":"ReadonlyArray","isTuple":false,"instantiatedType":78,"typeArguments":[1],"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":11,"character":55},"end":{"line":12,"character":30}},"flags":["Object"]},
{"id":82,"symbolName":"Container","isTuple":false,"instantiatedType":82,"typeArguments":[83],"firstDeclaration":{"path":"/home/src/workspaces/project/types.ts","start":{"line":1,"character":1},"end":{"line":4,"character":2}},"flags":["Object"]},
{"id":83,"symbolName":"T","isTuple":false,"firstDeclaration":{"path":"/home/src/workspaces/project/types.ts","start":{"line":1,"character":28},"end":{"line":1,"character":29}},"flags":["TypeParameter"]},
{"id":84,"symbolName":"Container","isTuple":false,"firstDeclaration":{"path":"/home/src/workspaces/project/types.ts","start":{"line":1,"character":1},"end":{"line":4,"character":2}},"flags":["TypeParameter"]},
{"id":85,"symbolName":"Container","isTuple":false,"instantiatedType":82,"typeArguments":[15],"firstDeclaration":{"path":"/home/src/workspaces/project/types.ts","start":{"line":1,"character":1},"end":{"line":4,"character":2}},"flags":["Object"]},
{"id":86,"isTuple":false,"flags":["NumberLiteral"],"display":"42"},
{"id":87,"isTuple":false,"flags":["NumberLiteral"],"display":"42"},
{"id":88,"symbolName":"map","isTuple":false,"firstDeclaration":{"path":"/home/src/workspaces/project/types.ts","start":{"line":2,"character":14},"end":{"line":3,"character":43}},"flags":["Object"],"display":"<U>(fn: (x: T) => U) => Container<U>"},
{"id":89,"symbolName":"map","isTuple":false,"firstDeclaration":{"path":"/home/src/workspaces/project/types.ts","start":{"line":2,"character":14},"end":{"line":3,"character":43}},"flags":["Object"],"display":"<U>(fn: (x: number) => U) => Container<U>"},
{"id":90,"symbolName":"U","isTuple":false,"firstDeclaration":{"path":"/home/src/workspaces/project/types.ts","start":{"line":3,"character":9},"end":{"line":3,"character":10}},"flags":["TypeParameter"]},
{"id":91,"symbolName":"U","isTuple":false,"firstDeclaration":{"path":"/home/src/workspaces/project/types.ts","start":{"line":3,"character":9},"end":{"line":3,"character":10}},"flags":["TypeParameter"]},
{"id":92,"symbolName":"__function","isTuple":false,"firstDeclaration":{"path":"/home/src/workspaces/project/main.ts","start":{"line":2,"character":47},"end":{"line":2,"character":87}},"flags":["Object"],"display":"<U>(fn: (x: number) => U) => { value: U; map: <U>(fn: (x: number) => U) => Container<U>; }"},
{"id":93,"symbolName":"__type","isTuple":false,"firstDeclaration":{"path":"/home/src/workspaces/project/types.ts","start":{"line":3,"character":15},"end":{"line":3,"character":27}},"flags":["Object"],"display":"(x: T) => U"},
{"id":94,"symbolName":"__type","isTuple":false,"firstDeclaration":{"path":"/home/src/workspaces/project/types.ts","start":{"line":3,"character":15},"end":{"line":3,"character":27}},"flags":["Object"],"display":"(x: number) => U"},
{"id":95,"symbolName":"Container","isTuple":false,"instantiatedType":82,"typeArguments":[90],"firstDeclaration":{"path":"/home/src/workspaces/project/types.ts","start":{"line":1,"character":1},"end":{"line":4,"character":2}},"flags":["Object"]},
{"id":96,"symbolName":"Container","isTuple":false,"instantiatedType":82,"typeArguments":[91],"firstDeclaration":{"path":"/home/src/workspaces/project/types.ts","start":{"line":1,"character":1},"end":{"line":4,"character":2}},"flags":["Object"]},
{"id":97,"symbolName":"map","isTuple":false,"firstDeclaration":{"path":"/home/src/workspaces/project/types.ts","start":{"line":2,"character":14},"end":{"line":3,"character":43}},"flags":["Object"],"display":"<U>(fn: (x: U) => U) => Container<U>"},
{"id":98,"symbolName":"__object","isTuple":false,"firstDeclaration":{"path":"/home/src/workspaces/project/main.ts","start":{"line":2,"character":57},"end":{"line":2,"character":86}},"flags":["Object"],"display":"{ value: U; map: <U>(fn: (x: number) => U) => Container<U>; }"},
{"id":99,"symbolName":"__object","isTuple":false,"firstDeclaration":{"path":"/home/src/workspaces/project/main.ts","start":{"line":2,"character":57},"end":{"line":2,"character":86}},"flags":["Object"],"display":"{ value: U; map: <U>(fn: (x: number) => U) => Container<U>; }"},
{"id":100,"symbolName":"__object","isTuple":false,"firstDeclaration":{"path":"/home/src/workspaces/project/main.ts","start":{"line":2,"character":29},"end":{"line":2,"character":89}},"flags":["Object"],"display":"{ value: number; map: <U>(fn: (x: number) => U) => { value: U; map: <U>(fn: (x: number) => U) => Container<U>; }; }"},
{"id":101,"symbolName":"U","isTuple":false,"firstDeclaration":{"path":"/home/src/workspaces/project/types.ts","start":{"line":3,"character":9},"end":{"line":3,"character":10}},"flags":["TypeParameter"]},
{"id":102,"symbolName":"__type","isTuple":false,"firstDeclaration":{"path":"/home/src/workspaces/project/types.ts","start":{"line":3,"character":15},"end":{"line":3,"character":27}},"flags":["Object"],"display":"(x: number) => any"},
{"id":103,"symbolName":"__type","isTuple":false,"firstDeclaration":{"path":"/home/src/workspaces/project/types.ts","start":{"line":3,"character":15},"end":{"line":3,"character":27}},"flags":["Object"],"display":"(x: U) => any"},
{"id":104,"symbolName":"__type","isTuple":false,"firstDeclaration":{"path":"/home/src/workspaces/project/types.ts","start":{"line":3,"character":15},"end":{"line":3,"character":27}},"flags":["Object"],"display":"(x: U) => U"},
{"id":105,"isTuple":false,"flags":["StringLiteral"],"display":"\"value\""},
{"id":106,"isTuple":false,"flags":["StringLiteral"],"display":"\"map\""},
{"id":107,"symbolName":"Container","isTuple":false,"instantiatedType":82,"typeArguments":[101],"firstDeclaration":{"path":"/home/src/workspaces/project/types.ts","start":{"line":1,"character":1},"end":{"line":4,"character":2}},"flags":["Object"]},
{"id":108,"symbolName":"T","isTuple":false,"firstDeclaration":{"path":"/home/src/workspaces/project/types.ts","start":{"line":5,"character":22},"end":{"line":5,"character":23}},"flags":["TypeParameter"]},
{"id":109,"symbolName":"Nullable","isTuple":false,"unionTypes":[10,13,108],"aliasTypeArguments":[108],"flags":["Union"]},
{"id":110,"symbolName":"Nullable","isTuple":false,"unionTypes":[10,13,14],"aliasTypeArguments":[14],"flags":["Union"]},
{"id":111,"isTuple":false,"flags":["StringLiteral"],"display":"\"hello\""},
{"id":112,"isTuple":false,"flags":["StringLiteral"],"display":"\"hello\""},
{"id":113,"symbolName":"__object","isTuple":false,"firstDeclaration":{"path":"/home/src/workspaces/project/main.ts","start":{"line":2,"character":57},"end":{"line":2,"character":86}},"flags":["Object"],"display":"{ value: U; map: <U>(fn: (x: number) => U) => Container<U>; }"}]


