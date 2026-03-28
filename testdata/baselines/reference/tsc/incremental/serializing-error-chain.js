currentDirectory::/home/src/workspaces/project
useCaseSensitiveFileNames::true
Input::
//// [/home/src/workspaces/project/index.tsx] *new* 
declare namespace JSX {
    interface ElementChildrenAttribute { children: {}; }
    interface IntrinsicElements { div: {} }
}

declare var React: any;

declare function Component(props: never): any;
declare function Component(props: { children?: number }): any;
(<Component>
    <div />
    <div />
</Component>)
//// [/home/src/workspaces/project/tsconfig.json] *new* 
{
    "compilerOptions": {
        "incremental": true,
        "strict": true,
        "jsx": "react",
        "module": "esnext",
    },
}

tsgo 
ExitStatus:: DiagnosticsPresent_OutputsGenerated
Output::
[96mindex.tsx[0m:[93m10[0m:[93m3[0m - [91merror[0m[90m TS2746: [0mThis JSX tag's 'children' prop expects a single child of type 'never', but multiple children were provided.

[7m10[0m (<Component>
[7m  [0m [91m  ~~~~~~~~~[0m

[96mindex.tsx[0m:[93m10[0m:[93m3[0m - [91merror[0m[90m TS2769: [0mNo overload matches this call.
  Overload 1 of 2, '(props: never): any', gave the following error.
    This JSX tag's 'children' prop expects a single child of type 'never', but multiple children were provided.
  Overload 2 of 2, '(props: { children?: number | undefined; }): any', gave the following error.
    Type '{ children: any[]; }' is not assignable to type '{ children?: number | undefined; }'.
      Types of property 'children' are incompatible.
        Type 'any[]' is not assignable to type 'number'.

[7m10[0m (<Component>
[7m  [0m [91m  ~~~~~~~~~[0m


Found 2 errors in the same file, starting at: index.tsx[90m:10[0m

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
//// [/home/src/workspaces/project/index.js] *new* 
"use strict";
(React.createElement(Component, null,
    React.createElement("div", null),
    React.createElement("div", null)));

//// [/home/src/workspaces/project/tsconfig.tsbuildinfo] *new* 
{"version":"FakeTSVersion","root":[2],"fileNames":["lib.es2025.full.d.ts","./index.tsx"],"fileInfos":[{"version":"8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"8ca521424834f2ae3377cc3ccc9dd3ef-declare namespace JSX {\n    interface ElementChildrenAttribute { children: {}; }\n    interface IntrinsicElements { div: {} }\n}\n\ndeclare var React: any;\n\ndeclare function Component(props: never): any;\ndeclare function Component(props: { children?: number }): any;\n(<Component>\n    <div />\n    <div />\n</Component>)","affectsGlobalScope":true,"impliedNodeFormat":1}],"options":{"jsx":3,"module":99,"strict":true},"semanticDiagnosticsPerFile":[[2,[{"pos":265,"end":274,"code":2746,"category":1,"messageKey":"This_JSX_tag_s_0_prop_expects_a_single_child_of_type_1_but_multiple_children_were_provided_2746","messageArgs":["children","never"]},{"pos":265,"end":274,"code":2769,"category":1,"messageKey":"No_overload_matches_this_call_2769","messageChain":[{"pos":265,"end":274,"code":2772,"category":1,"messageKey":"Overload_0_of_1_2_gave_the_following_error_2772","messageArgs":["1","2","(props: never): any"],"messageChain":[{"pos":265,"end":274,"code":2746,"category":1,"messageKey":"This_JSX_tag_s_0_prop_expects_a_single_child_of_type_1_but_multiple_children_were_provided_2746","messageArgs":["children","never"]}]},{"pos":265,"end":274,"code":2772,"category":1,"messageKey":"Overload_0_of_1_2_gave_the_following_error_2772","messageArgs":["2","2","(props: { children?: number | undefined; }): any"],"messageChain":[{"pos":265,"end":274,"code":2322,"category":1,"messageKey":"Type_0_is_not_assignable_to_type_1_2322","messageArgs":["{ children: any[]; }","{ children?: number | undefined; }"],"messageChain":[{"pos":265,"end":274,"code":2326,"category":1,"messageKey":"Types_of_property_0_are_incompatible_2326","messageArgs":["children"],"messageChain":[{"pos":265,"end":274,"code":2322,"category":1,"messageKey":"Type_0_is_not_assignable_to_type_1_2322","messageArgs":["any[]","number"]}]}]}]}]}]]]}
//// [/home/src/workspaces/project/tsconfig.tsbuildinfo.readable.baseline.txt] *new* 
{
  "version": "FakeTSVersion",
  "root": [
    {
      "files": [
        "./index.tsx"
      ],
      "original": 2
    }
  ],
  "fileNames": [
    "lib.es2025.full.d.ts",
    "./index.tsx"
  ],
  "fileInfos": [
    {
      "fileName": "lib.es2025.full.d.ts",
      "version": "8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
      "signature": "8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./index.tsx",
      "version": "8ca521424834f2ae3377cc3ccc9dd3ef-declare namespace JSX {\n    interface ElementChildrenAttribute { children: {}; }\n    interface IntrinsicElements { div: {} }\n}\n\ndeclare var React: any;\n\ndeclare function Component(props: never): any;\ndeclare function Component(props: { children?: number }): any;\n(<Component>\n    <div />\n    <div />\n</Component>)",
      "signature": "8ca521424834f2ae3377cc3ccc9dd3ef-declare namespace JSX {\n    interface ElementChildrenAttribute { children: {}; }\n    interface IntrinsicElements { div: {} }\n}\n\ndeclare var React: any;\n\ndeclare function Component(props: never): any;\ndeclare function Component(props: { children?: number }): any;\n(<Component>\n    <div />\n    <div />\n</Component>)",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "8ca521424834f2ae3377cc3ccc9dd3ef-declare namespace JSX {\n    interface ElementChildrenAttribute { children: {}; }\n    interface IntrinsicElements { div: {} }\n}\n\ndeclare var React: any;\n\ndeclare function Component(props: never): any;\ndeclare function Component(props: { children?: number }): any;\n(<Component>\n    <div />\n    <div />\n</Component>)",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    }
  ],
  "options": {
    "jsx": 3,
    "module": 99,
    "strict": true
  },
  "semanticDiagnosticsPerFile": [
    [
      "./index.tsx",
      [
        {
          "pos": 265,
          "end": 274,
          "code": 2746,
          "category": 1,
          "messageKey": "This_JSX_tag_s_0_prop_expects_a_single_child_of_type_1_but_multiple_children_were_provided_2746",
          "messageArgs": [
            "children",
            "never"
          ]
        },
        {
          "pos": 265,
          "end": 274,
          "code": 2769,
          "category": 1,
          "messageKey": "No_overload_matches_this_call_2769",
          "messageChain": [
            {
              "pos": 265,
              "end": 274,
              "code": 2772,
              "category": 1,
              "messageKey": "Overload_0_of_1_2_gave_the_following_error_2772",
              "messageArgs": [
                "1",
                "2",
                "(props: never): any"
              ],
              "messageChain": [
                {
                  "pos": 265,
                  "end": 274,
                  "code": 2746,
                  "category": 1,
                  "messageKey": "This_JSX_tag_s_0_prop_expects_a_single_child_of_type_1_but_multiple_children_were_provided_2746",
                  "messageArgs": [
                    "children",
                    "never"
                  ]
                }
              ]
            },
            {
              "pos": 265,
              "end": 274,
              "code": 2772,
              "category": 1,
              "messageKey": "Overload_0_of_1_2_gave_the_following_error_2772",
              "messageArgs": [
                "2",
                "2",
                "(props: { children?: number | undefined; }): any"
              ],
              "messageChain": [
                {
                  "pos": 265,
                  "end": 274,
                  "code": 2322,
                  "category": 1,
                  "messageKey": "Type_0_is_not_assignable_to_type_1_2322",
                  "messageArgs": [
                    "{ children: any[]; }",
                    "{ children?: number | undefined; }"
                  ],
                  "messageChain": [
                    {
                      "pos": 265,
                      "end": 274,
                      "code": 2326,
                      "category": 1,
                      "messageKey": "Types_of_property_0_are_incompatible_2326",
                      "messageArgs": [
                        "children"
                      ],
                      "messageChain": [
                        {
                          "pos": 265,
                          "end": 274,
                          "code": 2322,
                          "category": 1,
                          "messageKey": "Type_0_is_not_assignable_to_type_1_2322",
                          "messageArgs": [
                            "any[]",
                            "number"
                          ]
                        }
                      ]
                    }
                  ]
                }
              ]
            }
          ]
        }
      ]
    ]
  ],
  "size": 2730
}

tsconfig.json::
SemanticDiagnostics::
*refresh*    /home/src/tslibs/TS/Lib/lib.es2025.full.d.ts
*refresh*    /home/src/workspaces/project/index.tsx
Signatures::


Edit [0]:: no change

tsgo 
ExitStatus:: DiagnosticsPresent_OutputsGenerated
Output::
[96mindex.tsx[0m:[93m10[0m:[93m3[0m - [91merror[0m[90m TS2746: [0mThis JSX tag's 'children' prop expects a single child of type 'never', but multiple children were provided.

[7m10[0m (<Component>
[7m  [0m [91m  ~~~~~~~~~[0m

[96mindex.tsx[0m:[93m10[0m:[93m3[0m - [91merror[0m[90m TS2769: [0mNo overload matches this call.
  Overload 1 of 2, '(props: never): any', gave the following error.
    This JSX tag's 'children' prop expects a single child of type 'never', but multiple children were provided.
  Overload 2 of 2, '(props: { children?: number | undefined; }): any', gave the following error.
    Type '{ children: any[]; }' is not assignable to type '{ children?: number | undefined; }'.
      Types of property 'children' are incompatible.
        Type 'any[]' is not assignable to type 'number'.

[7m10[0m (<Component>
[7m  [0m [91m  ~~~~~~~~~[0m


Found 2 errors in the same file, starting at: index.tsx[90m:10[0m


tsconfig.json::
SemanticDiagnostics::
Signatures::
