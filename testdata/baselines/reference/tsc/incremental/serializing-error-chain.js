
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

ExitStatus:: 2

CompilerOptions::{}
Output::
[96mindex.tsx[0m:[93m11[0m:[93m23[0m - [91merror[0m[90m TS2769: [0mNo overload matches this call.
  The last overload gave the following error.
    Type '{ children: any[]; }' is not assignable to type '{ children?: number | undefined; }'.
      Types of property 'children' are incompatible.
        Type 'any[]' is not assignable to type 'number'.

[7m11[0m                     (<Component>
[7m  [0m [91m                      ~~~~~~~~~[0m

  [96mindex.tsx[0m:[93m10[0m:[93m38[0m - The last overload is declared here.
    [7m10[0m                     declare function Component(props: { children?: number }): any;
    [7m  [0m [96m                                     ~~~~~~~~~[0m


Found 1 error in index.tsx[90m:11[0m

//// [/home/src/tslibs/TS/Lib/lib.d.ts] *Lib*
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
(React.createElement(Component, null, React.createElement("div", null), React.createElement("div", null)));

//// [/home/src/workspaces/project/tsconfig.tsbuildinfo] *new* 
{"version":"FakeTSVersion","fileNames":["../../tslibs/TS/Lib/lib.d.ts","./index.tsx"],"fileInfos":[{"version":"7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"c7980af975245f04431574a9c187c9abd1c0ba29d83a127ad2af4b952296f935","affectsGlobalScope":true,"impliedNodeFormat":1}],"options":{"jsx":3,"module":99,"strict":true},"semanticDiagnosticsPerFile":[[2,[{"pos":426,"end":435,"code":2769,"category":1,"message":"No overload matches this call.","messageChain":[{"pos":426,"end":435,"code":2770,"category":1,"message":"The last overload gave the following error.","messageChain":[{"pos":426,"end":435,"code":2322,"category":1,"message":"Type '{ children: any[]; }' is not assignable to type '{ children?: number | undefined; }'.","messageChain":[{"pos":426,"end":435,"code":2326,"category":1,"message":"Types of property 'children' are incompatible.","messageChain":[{"pos":426,"end":435,"code":2322,"category":1,"message":"Type 'any[]' is not assignable to type 'number'."}]}]}]}],"relatedInformation":[{"pos":358,"end":367,"code":2771,"category":1,"message":"The last overload is declared here."}]}]]]}
//// [/home/src/workspaces/project/tsconfig.tsbuildinfo.readable.baseline.txt] *new* 
{
  "version": "FakeTSVersion",
  "fileNames": [
    "../../tslibs/TS/Lib/lib.d.ts",
    "./index.tsx"
  ],
  "fileInfos": [
    {
      "fileName": "../../tslibs/TS/Lib/lib.d.ts",
      "version": "7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e",
      "signature": "7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./index.tsx",
      "version": "c7980af975245f04431574a9c187c9abd1c0ba29d83a127ad2af4b952296f935",
      "signature": "c7980af975245f04431574a9c187c9abd1c0ba29d83a127ad2af4b952296f935",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "c7980af975245f04431574a9c187c9abd1c0ba29d83a127ad2af4b952296f935",
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
          "pos": 426,
          "end": 435,
          "code": 2769,
          "category": 1,
          "message": "No overload matches this call.",
          "messageChain": [
            {
              "pos": 426,
              "end": 435,
              "code": 2770,
              "category": 1,
              "message": "The last overload gave the following error.",
              "messageChain": [
                {
                  "pos": 426,
                  "end": 435,
                  "code": 2322,
                  "category": 1,
                  "message": "Type '{ children: any[]; }' is not assignable to type '{ children?: number | undefined; }'.",
                  "messageChain": [
                    {
                      "pos": 426,
                      "end": 435,
                      "code": 2326,
                      "category": 1,
                      "message": "Types of property 'children' are incompatible.",
                      "messageChain": [
                        {
                          "pos": 426,
                          "end": 435,
                          "code": 2322,
                          "category": 1,
                          "message": "Type 'any[]' is not assignable to type 'number'."
                        }
                      ]
                    }
                  ]
                }
              ]
            }
          ],
          "relatedInformation": [
            {
              "pos": 358,
              "end": 367,
              "code": 2771,
              "category": 1,
              "message": "The last overload is declared here."
            }
          ]
        }
      ]
    ]
  ],
  "size": 1181
}

