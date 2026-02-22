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

tsgo --generateTrace /home/src/workspaces/project/trace --singleThreaded
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
[{"configFilePath":"/home/src/workspaces/project/tsconfig.json","tracePath":"/home/src/workspaces/project/trace/trace.json","typesPath":"/home/src/workspaces/project/trace/types_0.json"}]
//// [/home/src/workspaces/project/trace/trace.json] *new* 
[
{"pid":1,"tid":1,"ph":"M","cat":"__metadata","ts":1000000000,"name":"process_name","args":{"name":"tsc"}},
{"pid":1,"tid":1,"ph":"M","cat":"__metadata","ts":1000000000,"name":"thread_name","args":{"name":"Main"}},
{"pid":1,"tid":1,"ph":"M","cat":"disabled-by-default-devtools.timeline","ts":1000000000,"name":"TracingStartedInBrowser"},
{"pid":1,"tid":1,"ph":"B","cat":"program","ts":1000000001,"name":"createProgram","args":{"configFilePath":"/home/src/workspaces/project/tsconfig.json"}},
{"pid":1,"tid":1,"ph":"E","ts":1000000002},
{"pid":1,"tid":1,"ph":"B","cat":"bind","ts":1000000003,"name":"bindSourceFiles"},
{"pid":1,"tid":1,"ph":"E","ts":1000000004},
{"pid":1,"tid":1,"ph":"B","cat":"check","ts":1000000005,"name":"checkSourceFiles"},
{"pid":1,"tid":1,"ph":"B","cat":"check","ts":1000000006,"name":"checkSourceFile","args":{"path":"/home/src/workspaces/project/main.ts"}},
{"pid":1,"tid":1,"ph":"B","cat":"check","ts":1000000007,"name":"checkVariableDeclaration","args":{"kind":"261","pos":"52","end":"135"}},
{"pid":1,"tid":1,"ph":"B","cat":"check","ts":1000000008,"name":"checkExpression","args":{"kind":"211","pos":"75","end":"135"}},
{"pid":1,"tid":1,"ph":"B","cat":"check","ts":1000000009,"name":"checkExpression","args":{"kind":"8","pos":"84","end":"87"}},
{"pid":1,"tid":1,"ph":"E","ts":1000000010},
{"pid":1,"tid":1,"ph":"B","cat":"check","ts":1000000011,"name":"checkExpression","args":{"kind":"220","pos":"93","end":"133"}},
{"pid":1,"tid":1,"ph":"B","cat":"check","ts":1000000012,"name":"checkExpression","args":{"kind":"218","pos":"101","end":"133"}},
{"pid":1,"tid":1,"ph":"B","cat":"check","ts":1000000013,"name":"checkExpression","args":{"kind":"211","pos":"103","end":"132"}},
{"pid":1,"tid":1,"ph":"B","cat":"check","ts":1000000014,"name":"checkExpression","args":{"kind":"214","pos":"111","end":"118"}},
{"pid":1,"tid":1,"ph":"B","cat":"check","ts":1000000015,"name":"checkExpression","args":{"kind":"79","pos":"111","end":"114"}},
{"pid":1,"tid":1,"ph":"E","ts":1000000016},
{"pid":1,"tid":1,"ph":"B","cat":"check","ts":1000000017,"name":"checkExpression","args":{"kind":"8","pos":"115","end":"117"}},
{"pid":1,"tid":1,"ph":"E","ts":1000000018},
{"pid":1,"tid":1,"ph":"E","ts":1000000019},
{"pid":1,"tid":1,"ph":"B","cat":"check","ts":1000000020,"name":"checkExpression","args":{"kind":"212","pos":"124","end":"130"}},
{"pid":1,"tid":1,"ph":"B","cat":"check","ts":1000000021,"name":"checkExpression","args":{"kind":"79","pos":"124","end":"126"}},
{"pid":1,"tid":1,"ph":"E","ts":1000000022},
{"pid":1,"tid":1,"ph":"E","ts":1000000023},
{"pid":1,"tid":1,"ph":"E","ts":1000000024},
{"pid":1,"tid":1,"ph":"E","ts":1000000025},
{"pid":1,"tid":1,"ph":"E","ts":1000000026},
{"pid":1,"tid":1,"ph":"E","ts":1000000027},
{"pid":1,"tid":1,"ph":"B","cat":"checkTypes","ts":1000000028,"name":"structuredTypeRelatedTo","args":{"sourceId":"100","targetId":"85"}},
{"pid":1,"tid":1,"ph":"B","cat":"checkTypes","ts":1000000029,"name":"structuredTypeRelatedTo","args":{"sourceId":"92","targetId":"89"}},
{"pid":1,"tid":1,"ph":"B","cat":"checkTypes","ts":1000000030,"name":"structuredTypeRelatedTo","args":{"sourceId":"99","targetId":"96"}},
{"pid":1,"tid":1,"ph":"B","cat":"checkTypes","ts":1000000031,"name":"structuredTypeRelatedTo","args":{"sourceId":"89","targetId":"97"}},
{"pid":1,"tid":1,"ph":"B","cat":"checkTypes","ts":1000000032,"name":"structuredTypeRelatedTo","args":{"sourceId":"15","targetId":"91"}},
{"pid":1,"tid":1,"ph":"E","ts":1000000033},
{"pid":1,"tid":1,"ph":"E","ts":1000000034},
{"pid":1,"tid":1,"ph":"E","ts":1000000035},
{"pid":1,"tid":1,"ph":"E","ts":1000000036},
{"pid":1,"tid":1,"ph":"E","ts":1000000037},
{"pid":1,"tid":1,"ph":"B","cat":"checkTypes","ts":1000000038,"name":"structuredTypeRelatedTo","args":{"sourceId":"99","targetId":"89"}},
{"pid":1,"tid":1,"ph":"E","ts":1000000039},
{"pid":1,"tid":1,"ph":"B","cat":"checkTypes","ts":1000000040,"name":"structuredTypeRelatedTo","args":{"sourceId":"96","targetId":"97"}},
{"pid":1,"tid":1,"ph":"E","ts":1000000041},
{"pid":1,"tid":1,"ph":"B","cat":"check","ts":1000000042,"name":"checkExpression","args":{"kind":"212","pos":"124","end":"130"}},
{"pid":1,"tid":1,"ph":"B","cat":"check","ts":1000000043,"name":"checkExpression","args":{"kind":"79","pos":"124","end":"126"}},
{"pid":1,"tid":1,"ph":"E","ts":1000000044},
{"pid":1,"tid":1,"ph":"E","ts":1000000045},
{"pid":1,"tid":1,"ph":"B","cat":"checkTypes","ts":1000000046,"name":"structuredTypeRelatedTo","args":{"sourceId":"89","targetId":"97"}},
{"pid":1,"tid":1,"ph":"B","cat":"checkTypes","ts":1000000047,"name":"structuredTypeRelatedTo","args":{"sourceId":"15","targetId":"91"}},
{"pid":1,"tid":1,"ph":"E","ts":1000000048},
{"pid":1,"tid":1,"ph":"E","ts":1000000049},
{"pid":1,"tid":1,"ph":"E","ts":1000000050},
{"pid":1,"tid":1,"ph":"B","cat":"check","ts":1000000051,"name":"checkVariableDeclaration","args":{"kind":"261","pos":"142","end":"172"}},
{"pid":1,"tid":1,"ph":"B","cat":"check","ts":1000000052,"name":"checkExpression","args":{"kind":"10","pos":"164","end":"172"}},
{"pid":1,"tid":1,"ph":"E","ts":1000000053},
{"pid":1,"tid":1,"ph":"E","ts":1000000054},
{"pid":1,"tid":1,"ph":"B","cat":"check","ts":1000000055,"name":"checkDeferredNode","args":{"kind":"220","pos":"93","end":"133"}},
{"pid":1,"tid":1,"ph":"B","cat":"check","ts":1000000056,"name":"checkExpression","args":{"kind":"218","pos":"101","end":"133"}},
{"pid":1,"tid":1,"ph":"B","cat":"check","ts":1000000057,"name":"checkExpression","args":{"kind":"211","pos":"103","end":"132"}},
{"pid":1,"tid":1,"ph":"B","cat":"check","ts":1000000058,"name":"checkExpression","args":{"kind":"214","pos":"111","end":"118"}},
{"pid":1,"tid":1,"ph":"E","ts":1000000059},
{"pid":1,"tid":1,"ph":"B","cat":"check","ts":1000000060,"name":"checkExpression","args":{"kind":"212","pos":"124","end":"130"}},
{"pid":1,"tid":1,"ph":"B","cat":"check","ts":1000000061,"name":"checkExpression","args":{"kind":"79","pos":"124","end":"126"}},
{"pid":1,"tid":1,"ph":"E","ts":1000000062},
{"pid":1,"tid":1,"ph":"E","ts":1000000063},
{"pid":1,"tid":1,"ph":"E","ts":1000000064},
{"pid":1,"tid":1,"ph":"E","ts":1000000065},
{"pid":1,"tid":1,"ph":"E","ts":1000000066},
{"pid":1,"tid":1,"ph":"E","ts":1000000067},
{"pid":1,"tid":1,"ph":"B","cat":"check","ts":1000000068,"name":"checkSourceFile","args":{"path":"/home/src/workspaces/project/types.ts"}},
{"pid":1,"tid":1,"ph":"B","cat":"check","ts":1000000069,"name":"checkDeferredNode","args":{"kind":"169","pos":"27","end":"28"}},
{"pid":1,"tid":1,"ph":"E","ts":1000000070},
{"pid":1,"tid":1,"ph":"B","cat":"check","ts":1000000071,"name":"checkDeferredNode","args":{"kind":"169","pos":"54","end":"55"}},
{"pid":1,"tid":1,"ph":"E","ts":1000000072},
{"pid":1,"tid":1,"ph":"B","cat":"check","ts":1000000073,"name":"checkDeferredNode","args":{"kind":"169","pos":"112","end":"113"}},
{"pid":1,"tid":1,"ph":"E","ts":1000000074},
{"pid":1,"tid":1,"ph":"E","ts":1000000075},
{"pid":1,"tid":1,"ph":"B","cat":"check","ts":1000000076,"name":"checkSourceFile","args":{"path":"/home/src/tslibs/TS/Lib/lib.es2025.full.d.ts"}},
{"pid":1,"tid":1,"ph":"B","cat":"check","ts":1000000077,"name":"checkVariableDeclaration","args":{"kind":"261","pos":"511","end":"537"}},
{"pid":1,"tid":1,"ph":"E","ts":1000000078},
{"pid":1,"tid":1,"ph":"B","cat":"check","ts":1000000079,"name":"checkExpression","args":{"kind":"212","pos":"572","end":"590"}},
{"pid":1,"tid":1,"ph":"B","cat":"check","ts":1000000080,"name":"checkExpression","args":{"kind":"79","pos":"572","end":"578"}},
{"pid":1,"tid":1,"ph":"E","ts":1000000081},
{"pid":1,"tid":1,"ph":"E","ts":1000000082},
{"pid":1,"tid":1,"ph":"B","cat":"check","ts":1000000083,"name":"checkVariableDeclaration","args":{"kind":"261","pos":"616","end":"650"}},
{"pid":1,"tid":1,"ph":"E","ts":1000000084},
{"pid":1,"tid":1,"ph":"B","cat":"check","ts":1000000085,"name":"checkDeferredNode","args":{"kind":"169","pos":"296","end":"297"}},
{"pid":1,"tid":1,"ph":"E","ts":1000000086},
{"pid":1,"tid":1,"ph":"B","cat":"check","ts":1000000087,"name":"checkDeferredNode","args":{"kind":"169","pos":"359","end":"360"}},
{"pid":1,"tid":1,"ph":"E","ts":1000000088},
{"pid":1,"tid":1,"ph":"E","ts":1000000089},
{"pid":1,"tid":1,"ph":"E","ts":1000000090},
{"pid":1,"tid":1,"ph":"B","cat":"emit","ts":1000000091,"name":"emit"},
{"pid":1,"tid":1,"ph":"B","cat":"emit","ts":1000000092,"name":"emit","args":{"path":"/home/src/workspaces/project/main.ts"}},
{"pid":1,"tid":1,"ph":"E","ts":1000000093},
{"pid":1,"tid":1,"ph":"B","cat":"emit","ts":1000000094,"name":"emit","args":{"path":"/home/src/workspaces/project/types.ts"}},
{"pid":1,"tid":1,"ph":"E","ts":1000000095},
{"pid":1,"tid":1,"ph":"E","ts":1000000096}
]

//// [/home/src/workspaces/project/trace/types_0.json] *new* 
[{"id":1,"intrinsicName":"any","recursionId":0,"flags":["Any"]},
{"id":2,"intrinsicName":"any","recursionId":1,"flags":["Any"]},
{"id":3,"intrinsicName":"any","recursionId":2,"flags":["Any"]},
{"id":4,"intrinsicName":"any","recursionId":3,"flags":["Any"]},
{"id":5,"intrinsicName":"any","recursionId":4,"flags":["Any"]},
{"id":6,"intrinsicName":"unresolved","recursionId":5,"flags":["Any"]},
{"id":7,"intrinsicName":"any","recursionId":6,"flags":["Any"]},
{"id":8,"intrinsicName":"intrinsic","recursionId":7,"flags":["Any"]},
{"id":9,"intrinsicName":"unknown","recursionId":8,"flags":["Unknown"]},
{"id":10,"intrinsicName":"undefined","recursionId":9,"flags":["Undefined"]},
{"id":11,"intrinsicName":"undefined","recursionId":10,"flags":["Undefined"]},
{"id":12,"intrinsicName":"undefined","recursionId":11,"flags":["Undefined"]},
{"id":13,"intrinsicName":"null","recursionId":12,"flags":["Null"]},
{"id":14,"intrinsicName":"string","recursionId":13,"flags":["String"]},
{"id":15,"intrinsicName":"number","recursionId":14,"flags":["Number"]},
{"id":16,"intrinsicName":"bigint","recursionId":15,"flags":["BigInt"]},
{"id":17,"recursionId":16,"flags":["BooleanLiteral"],"display":"false"},
{"id":18,"recursionId":17,"flags":["BooleanLiteral"],"display":"false"},
{"id":19,"recursionId":18,"flags":["BooleanLiteral"],"display":"true"},
{"id":20,"recursionId":19,"flags":["BooleanLiteral"],"display":"true"},
{"id":21,"recursionId":20,"unionTypes":[17,19],"flags":["Boolean","Union"]},
{"id":22,"intrinsicName":"symbol","recursionId":21,"flags":["ESSymbol"]},
{"id":23,"intrinsicName":"void","recursionId":22,"flags":["Void"]},
{"id":24,"intrinsicName":"never","recursionId":23,"flags":["Never"]},
{"id":25,"intrinsicName":"never","recursionId":24,"flags":["Never"]},
{"id":26,"intrinsicName":"never","recursionId":25,"flags":["Never"]},
{"id":27,"intrinsicName":"never","recursionId":26,"flags":["Never"]},
{"id":28,"intrinsicName":"object","recursionId":27,"flags":["NonPrimitive"]},
{"id":29,"recursionId":28,"unionTypes":[14,15],"flags":["Union"]},
{"id":30,"recursionId":29,"unionTypes":[14,15,22],"flags":["Union"]},
{"id":31,"recursionId":30,"unionTypes":[15,16],"flags":["Union"]},
{"id":32,"recursionId":31,"flags":["TemplateLiteral"]},
{"id":33,"recursionId":32,"unionTypes":[10,13,14,15,16,17,19],"flags":["Union"]},
{"id":34,"intrinsicName":"never","recursionId":33,"flags":["Never"]},
{"id":35,"recursionId":34,"flags":["Object"],"display":"{}"},
{"id":36,"recursionId":35,"flags":["Object"],"display":"{}"},
{"id":37,"recursionId":36,"flags":["Object"],"display":"{}"},
{"id":38,"symbolName":"__type","recursionId":37,"flags":["Object"],"display":"{}"},
{"id":39,"recursionId":38,"flags":["Object"],"display":"{}"},
{"id":40,"recursionId":39,"unionTypes":[10,13,39],"flags":["Union"]},
{"id":41,"recursionId":40,"flags":["Object"],"display":"{}"},
{"id":42,"recursionId":41,"flags":["Object"],"display":"{}"},
{"id":43,"recursionId":42,"flags":["Object"],"display":"{}"},
{"id":44,"recursionId":43,"flags":["Object"],"display":"{}"},
{"id":45,"recursionId":44,"flags":["Object"],"display":"{}"},
{"id":46,"recursionId":45,"flags":["TypeParameter"]},
{"id":47,"recursionId":46,"flags":["TypeParameter"]},
{"id":48,"recursionId":47,"flags":["TypeParameter"]},
{"id":49,"recursionId":48,"flags":["TypeParameter"]},
{"id":50,"recursionId":49,"flags":["TypeParameter"]},
{"id":51,"recursionId":50,"flags":["StringLiteral"],"display":"\"\""},
{"id":52,"recursionId":51,"flags":["NumberLiteral"],"display":"0"},
{"id":53,"recursionId":52,"flags":["BigIntLiteral"],"display":"0n"},
{"id":54,"recursionId":53,"flags":["StringLiteral"],"display":"\"bigint\""},
{"id":55,"recursionId":54,"flags":["StringLiteral"],"display":"\"boolean\""},
{"id":56,"recursionId":55,"flags":["StringLiteral"],"display":"\"function\""},
{"id":57,"recursionId":56,"flags":["StringLiteral"],"display":"\"number\""},
{"id":58,"recursionId":57,"flags":["StringLiteral"],"display":"\"object\""},
{"id":59,"recursionId":58,"flags":["StringLiteral"],"display":"\"string\""},
{"id":60,"recursionId":59,"flags":["StringLiteral"],"display":"\"symbol\""},
{"id":61,"recursionId":60,"flags":["StringLiteral"],"display":"\"undefined\""},
{"id":62,"recursionId":61,"unionTypes":[54,55,56,57,58,59,60,61],"flags":["Union"]},
{"id":63,"symbolName":"IArguments","recursionId":62,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":5,"character":29},"end":{"line":6,"character":24}},"flags":["Object"]},
{"id":64,"symbolName":"globalThis","recursionId":63,"flags":["Object"],"display":"typeof globalThis"},
{"id":65,"symbolName":"Array","recursionId":64,"instantiatedType":65,"typeArguments":[66],"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":10,"character":34},"end":{"line":11,"character":55}},"flags":["Object"]},
{"id":66,"symbolName":"T","recursionId":65,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":11,"character":17},"end":{"line":11,"character":18}},"flags":["TypeParameter"]},
{"id":67,"symbolName":"Array","recursionId":64,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":10,"character":34},"end":{"line":11,"character":55}},"flags":["TypeParameter"]},
{"id":68,"symbolName":"Object","recursionId":66,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":7,"character":41},"end":{"line":8,"character":20}},"flags":["Object"]},
{"id":69,"symbolName":"Function","recursionId":67,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":2,"character":21},"end":{"line":3,"character":22}},"flags":["Object"]},
{"id":70,"symbolName":"CallableFunction","recursionId":68,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":3,"character":22},"end":{"line":4,"character":30}},"flags":["Object"]},
{"id":71,"symbolName":"NewableFunction","recursionId":69,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":4,"character":30},"end":{"line":5,"character":29}},"flags":["Object"]},
{"id":72,"symbolName":"String","recursionId":70,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":9,"character":20},"end":{"line":10,"character":34}},"flags":["Object"]},
{"id":73,"symbolName":"Number","recursionId":71,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":6,"character":24},"end":{"line":7,"character":41}},"flags":["Object"]},
{"id":74,"symbolName":"Boolean","recursionId":72,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":1,"character":1},"end":{"line":2,"character":21}},"flags":["Object"]},
{"id":75,"symbolName":"RegExp","recursionId":73,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":8,"character":20},"end":{"line":9,"character":20}},"flags":["Object"]},
{"id":76,"symbolName":"Array","recursionId":64,"instantiatedType":65,"typeArguments":[1],"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":10,"character":34},"end":{"line":11,"character":55}},"flags":["Object"]},
{"id":77,"symbolName":"Array","recursionId":64,"instantiatedType":65,"typeArguments":[2],"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":10,"character":34},"end":{"line":11,"character":55}},"flags":["Object"]},
{"id":78,"symbolName":"ReadonlyArray","recursionId":74,"instantiatedType":78,"typeArguments":[79],"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":11,"character":55},"end":{"line":12,"character":30}},"flags":["Object"]},
{"id":79,"symbolName":"T","recursionId":75,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":12,"character":25},"end":{"line":12,"character":26}},"flags":["TypeParameter"]},
{"id":80,"symbolName":"ReadonlyArray","recursionId":74,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":11,"character":55},"end":{"line":12,"character":30}},"flags":["TypeParameter"]},
{"id":81,"symbolName":"ReadonlyArray","recursionId":74,"instantiatedType":78,"typeArguments":[1],"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":11,"character":55},"end":{"line":12,"character":30}},"flags":["Object"]},
{"id":82,"symbolName":"Container","recursionId":76,"instantiatedType":82,"typeArguments":[83],"firstDeclaration":{"path":"/home/src/workspaces/project/types.ts","start":{"line":1,"character":1},"end":{"line":4,"character":2}},"flags":["Object"]},
{"id":83,"symbolName":"T","recursionId":77,"firstDeclaration":{"path":"/home/src/workspaces/project/types.ts","start":{"line":1,"character":28},"end":{"line":1,"character":29}},"flags":["TypeParameter"]},
{"id":84,"symbolName":"Container","recursionId":76,"firstDeclaration":{"path":"/home/src/workspaces/project/types.ts","start":{"line":1,"character":1},"end":{"line":4,"character":2}},"flags":["TypeParameter"]},
{"id":85,"symbolName":"Container","recursionId":76,"instantiatedType":82,"typeArguments":[15],"firstDeclaration":{"path":"/home/src/workspaces/project/types.ts","start":{"line":1,"character":1},"end":{"line":4,"character":2}},"flags":["Object"]},
{"id":86,"recursionId":78,"flags":["NumberLiteral"],"display":"42"},
{"id":87,"recursionId":79,"flags":["NumberLiteral"],"display":"42"},
{"id":88,"symbolName":"map","recursionId":80,"firstDeclaration":{"path":"/home/src/workspaces/project/types.ts","start":{"line":2,"character":14},"end":{"line":3,"character":43}},"flags":["Object"],"display":"<U>(fn: (x: T) => U) => Container<U>"},
{"id":89,"symbolName":"map","recursionId":80,"firstDeclaration":{"path":"/home/src/workspaces/project/types.ts","start":{"line":2,"character":14},"end":{"line":3,"character":43}},"flags":["Object"],"display":"<U>(fn: (x: number) => U) => Container<U>"},
{"id":90,"symbolName":"U","recursionId":81,"firstDeclaration":{"path":"/home/src/workspaces/project/types.ts","start":{"line":3,"character":9},"end":{"line":3,"character":10}},"flags":["TypeParameter"]},
{"id":91,"symbolName":"U","recursionId":81,"firstDeclaration":{"path":"/home/src/workspaces/project/types.ts","start":{"line":3,"character":9},"end":{"line":3,"character":10}},"flags":["TypeParameter"]},
{"id":92,"symbolName":"__function","recursionId":82,"firstDeclaration":{"path":"/home/src/workspaces/project/main.ts","start":{"line":2,"character":47},"end":{"line":2,"character":87}},"flags":["Object"],"display":"<U>(fn: (x: number) => U) => { value: U; map: <U>(fn: (x: number) => U) => Container<U>; }"},
{"id":93,"symbolName":"__type","recursionId":83,"firstDeclaration":{"path":"/home/src/workspaces/project/types.ts","start":{"line":3,"character":15},"end":{"line":3,"character":27}},"flags":["Object"],"display":"(x: T) => U"},
{"id":94,"symbolName":"__type","recursionId":83,"firstDeclaration":{"path":"/home/src/workspaces/project/types.ts","start":{"line":3,"character":15},"end":{"line":3,"character":27}},"flags":["Object"],"display":"(x: number) => U"},
{"id":95,"symbolName":"Container","recursionId":76,"instantiatedType":82,"typeArguments":[90],"firstDeclaration":{"path":"/home/src/workspaces/project/types.ts","start":{"line":1,"character":1},"end":{"line":4,"character":2}},"flags":["Object"]},
{"id":96,"symbolName":"Container","recursionId":76,"instantiatedType":82,"typeArguments":[91],"firstDeclaration":{"path":"/home/src/workspaces/project/types.ts","start":{"line":1,"character":1},"end":{"line":4,"character":2}},"flags":["Object"]},
{"id":97,"symbolName":"map","recursionId":80,"firstDeclaration":{"path":"/home/src/workspaces/project/types.ts","start":{"line":2,"character":14},"end":{"line":3,"character":43}},"flags":["Object"],"display":"<U>(fn: (x: U) => U) => Container<U>"},
{"id":98,"symbolName":"__object","recursionId":84,"firstDeclaration":{"path":"/home/src/workspaces/project/main.ts","start":{"line":2,"character":57},"end":{"line":2,"character":86}},"flags":["Object"],"display":"{ value: U; map: <U>(fn: (x: number) => U) => Container<U>; }"},
{"id":99,"symbolName":"__object","recursionId":85,"firstDeclaration":{"path":"/home/src/workspaces/project/main.ts","start":{"line":2,"character":57},"end":{"line":2,"character":86}},"flags":["Object"],"display":"{ value: U; map: <U>(fn: (x: number) => U) => Container<U>; }"},
{"id":100,"symbolName":"__object","recursionId":86,"firstDeclaration":{"path":"/home/src/workspaces/project/main.ts","start":{"line":2,"character":29},"end":{"line":2,"character":89}},"flags":["Object"],"display":"{ value: number; map: <U>(fn: (x: number) => U) => { value: U; map: <U>(fn: (x: number) => U) => Container<U>; }; }"},
{"id":101,"symbolName":"U","recursionId":81,"firstDeclaration":{"path":"/home/src/workspaces/project/types.ts","start":{"line":3,"character":9},"end":{"line":3,"character":10}},"flags":["TypeParameter"]},
{"id":102,"symbolName":"__type","recursionId":83,"firstDeclaration":{"path":"/home/src/workspaces/project/types.ts","start":{"line":3,"character":15},"end":{"line":3,"character":27}},"flags":["Object"],"display":"(x: number) => any"},
{"id":103,"symbolName":"__type","recursionId":83,"firstDeclaration":{"path":"/home/src/workspaces/project/types.ts","start":{"line":3,"character":15},"end":{"line":3,"character":27}},"flags":["Object"],"display":"(x: U) => any"},
{"id":104,"symbolName":"__type","recursionId":83,"firstDeclaration":{"path":"/home/src/workspaces/project/types.ts","start":{"line":3,"character":15},"end":{"line":3,"character":27}},"flags":["Object"],"display":"(x: U) => U"},
{"id":105,"recursionId":87,"flags":["StringLiteral"],"display":"\"value\""},
{"id":106,"recursionId":88,"flags":["StringLiteral"],"display":"\"map\""},
{"id":107,"symbolName":"Container","recursionId":76,"instantiatedType":82,"typeArguments":[101],"firstDeclaration":{"path":"/home/src/workspaces/project/types.ts","start":{"line":1,"character":1},"end":{"line":4,"character":2}},"flags":["Object"]},
{"id":108,"symbolName":"T","recursionId":89,"firstDeclaration":{"path":"/home/src/workspaces/project/types.ts","start":{"line":5,"character":22},"end":{"line":5,"character":23}},"flags":["TypeParameter"]},
{"id":109,"symbolName":"Nullable","recursionId":90,"unionTypes":[10,13,108],"aliasTypeArguments":[108],"flags":["Union"]},
{"id":110,"symbolName":"Nullable","recursionId":91,"unionTypes":[10,13,14],"aliasTypeArguments":[14],"flags":["Union"]},
{"id":111,"recursionId":92,"flags":["StringLiteral"],"display":"\"hello\""},
{"id":112,"recursionId":93,"flags":["StringLiteral"],"display":"\"hello\""},
{"id":113,"symbolName":"__object","recursionId":94,"firstDeclaration":{"path":"/home/src/workspaces/project/main.ts","start":{"line":2,"character":57},"end":{"line":2,"character":86}},"flags":["Object"],"display":"{ value: U; map: <U>(fn: (x: number) => U) => Container<U>; }"},
{"id":114,"symbolName":"Container","recursionId":76,"instantiatedType":82,"typeArguments":[83,84],"firstDeclaration":{"path":"/home/src/workspaces/project/types.ts","start":{"line":1,"character":1},"end":{"line":4,"character":2}},"flags":["Object"]},
{"id":115,"symbolName":"Array","recursionId":64,"instantiatedType":65,"typeArguments":[66,67],"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":10,"character":34},"end":{"line":11,"character":55}},"flags":["Object"]},
{"id":116,"recursionId":95,"flags":["StringLiteral"],"display":"\"length\""},
{"id":117,"symbolName":"ReadonlyArray","recursionId":74,"instantiatedType":78,"typeArguments":[79,80],"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":11,"character":55},"end":{"line":12,"character":30}},"flags":["Object"]},
{"id":118,"symbolName":"SymbolConstructor","recursionId":96,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":12,"character":30},"end":{"line":17,"character":2}},"flags":["Object"]},
{"id":119,"recursionId":97,"unionTypes":[10,14,15],"flags":["Union"]},
{"id":120,"symbolName":"toStringTag","recursionId":98,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":15,"character":31},"end":{"line":16,"character":34}},"flags":["UniqueESSymbol"]},
{"id":121,"symbolName":"Symbol","recursionId":99,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":18,"character":12},"end":{"line":18,"character":38}},"flags":["Object"]},
{"id":122,"symbolName":"__type","recursionId":100,"firstDeclaration":{"path":"/home/src/tslibs/ts/lib/lib.es2025.full.d.ts","start":{"line":22,"character":23},"end":{"line":22,"character":48}},"flags":["Object"],"display":"{ log(msg: any): void; }"}]


