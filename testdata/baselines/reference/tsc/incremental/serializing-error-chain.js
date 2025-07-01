
currentDirectory::/home/src/workspaces/project
useCaseSensitiveFileNames::true
Input::
//// [/home/src/workspaces/project/index.tsx] new file

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
//// [/home/src/workspaces/project/tsconfig.json] new file
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
[96mindex.tsx[0m:[93m11[0m:[93m23[0m - [91merror[0m[90m TS2746: [0mThis JSX tag's 'children' prop expects a single child of type 'number | undefined', but multiple children were provided.

[7m11[0m                     (<Component>
[7m  [0m [91m                      ~~~~~~~~~[0m

[96mindex.tsx[0m:[93m11[0m:[93m23[0m - [91merror[0m[90m TS2769: [0mNo overload matches this call.
  The last overload gave the following error.
    This JSX tag's 'children' prop expects a single child of type 'number | undefined', but multiple children were provided.

[7m11[0m                     (<Component>
[7m  [0m [91m                      ~~~~~~~~~[0m

  [96mindex.tsx[0m:[93m10[0m:[93m38[0m - The last overload is declared here.
    [7m10[0m                     declare function Component(props: { children?: number }): any;
    [7m  [0m [96m                                     ~~~~~~~~~[0m


Found 2 errors in the same file, starting at: index.tsx[90m:11[0m

//// [/home/src/workspaces/project/index.js] new file
(React.createElement(Component, null, React.createElement("div", null), React.createElement("div", null)));

//// [/home/src/workspaces/project/index.tsx] no change
//// [/home/src/workspaces/project/tsconfig.json] no change
//// [/home/src/workspaces/project/tsconfig.tsbuildinfo] new file
{"version":"FakeTSVersion","fileNames":["bundled:///libs/lib.d.ts","bundled:///libs/lib.es5.d.ts","bundled:///libs/lib.dom.d.ts","bundled:///libs/lib.webworker.importscripts.d.ts","bundled:///libs/lib.scripthost.d.ts","bundled:///libs/lib.decorators.d.ts","bundled:///libs/lib.decorators.legacy.d.ts","./index.tsx"],"fileInfos":["a7297ff837fcdf174a9524925966429eb8e5feecc2cc55cc06574e6b092c1eaa",{"version":"69684132aeb9b5642cbcd9e22dff7818ff0ee1aa831728af0ecf97d3364d5546","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"092c2bfe125ce69dbb1223c85d68d4d2397d7d8411867b5cc03cec902c233763","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"80e18897e5884b6723488d4f5652167e7bb5024f946743134ecc4aa4ee731f89","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"cd034f499c6cdca722b60c04b5b1b78e058487a7085a8e0d6fb50809947ee573","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"8e7f8264d0fb4c5339605a15daadb037bf238c10b654bb3eee14208f860a32ea","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"782dec38049b92d4e85c1585fbea5474a219c6984a35b004963b00beb1aab538","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"c7980af975245f04431574a9c187c9abd1c0ba29d83a127ad2af4b952296f935","affectsGlobalScope":true,"impliedNodeFormat":1}],"options":{"jsx":3,"module":99,"strict":true},"semanticDiagnosticsPerFile":[[8,[{"pos":426,"end":435,"code":2746,"category":1,"message":"This JSX tag's 'children' prop expects a single child of type 'number | undefined', but multiple children were provided."},{"pos":426,"end":435,"code":2769,"category":1,"message":"No overload matches this call.","messageChain":[{"pos":426,"end":435,"code":2770,"category":1,"message":"The last overload gave the following error.","messageChain":[{"pos":426,"end":435,"code":2746,"category":1,"message":"This JSX tag's 'children' prop expects a single child of type 'number | undefined', but multiple children were provided."}]}],"relatedInformation":[{"pos":358,"end":367,"code":2771,"category":1,"message":"The last overload is declared here."}]}]]]}
//// [/home/src/workspaces/project/tsconfig.tsbuildinfo.readable.baseline.txt] new file
{
  "version": "FakeTSVersion",
  "fileNames": [
    "bundled:///libs/lib.d.ts",
    "bundled:///libs/lib.es5.d.ts",
    "bundled:///libs/lib.dom.d.ts",
    "bundled:///libs/lib.webworker.importscripts.d.ts",
    "bundled:///libs/lib.scripthost.d.ts",
    "bundled:///libs/lib.decorators.d.ts",
    "bundled:///libs/lib.decorators.legacy.d.ts",
    "./index.tsx"
  ],
  "fileInfos": [
    {
      "fileName": "bundled:///libs/lib.d.ts",
      "version": "a7297ff837fcdf174a9524925966429eb8e5feecc2cc55cc06574e6b092c1eaa",
      "signature": "a7297ff837fcdf174a9524925966429eb8e5feecc2cc55cc06574e6b092c1eaa",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "bundled:///libs/lib.es5.d.ts",
      "version": "69684132aeb9b5642cbcd9e22dff7818ff0ee1aa831728af0ecf97d3364d5546",
      "signature": "69684132aeb9b5642cbcd9e22dff7818ff0ee1aa831728af0ecf97d3364d5546",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "69684132aeb9b5642cbcd9e22dff7818ff0ee1aa831728af0ecf97d3364d5546",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.dom.d.ts",
      "version": "092c2bfe125ce69dbb1223c85d68d4d2397d7d8411867b5cc03cec902c233763",
      "signature": "092c2bfe125ce69dbb1223c85d68d4d2397d7d8411867b5cc03cec902c233763",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "092c2bfe125ce69dbb1223c85d68d4d2397d7d8411867b5cc03cec902c233763",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.webworker.importscripts.d.ts",
      "version": "80e18897e5884b6723488d4f5652167e7bb5024f946743134ecc4aa4ee731f89",
      "signature": "80e18897e5884b6723488d4f5652167e7bb5024f946743134ecc4aa4ee731f89",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "80e18897e5884b6723488d4f5652167e7bb5024f946743134ecc4aa4ee731f89",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.scripthost.d.ts",
      "version": "cd034f499c6cdca722b60c04b5b1b78e058487a7085a8e0d6fb50809947ee573",
      "signature": "cd034f499c6cdca722b60c04b5b1b78e058487a7085a8e0d6fb50809947ee573",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "cd034f499c6cdca722b60c04b5b1b78e058487a7085a8e0d6fb50809947ee573",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.decorators.d.ts",
      "version": "8e7f8264d0fb4c5339605a15daadb037bf238c10b654bb3eee14208f860a32ea",
      "signature": "8e7f8264d0fb4c5339605a15daadb037bf238c10b654bb3eee14208f860a32ea",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "8e7f8264d0fb4c5339605a15daadb037bf238c10b654bb3eee14208f860a32ea",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.decorators.legacy.d.ts",
      "version": "782dec38049b92d4e85c1585fbea5474a219c6984a35b004963b00beb1aab538",
      "signature": "782dec38049b92d4e85c1585fbea5474a219c6984a35b004963b00beb1aab538",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "782dec38049b92d4e85c1585fbea5474a219c6984a35b004963b00beb1aab538",
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
          "code": 2746,
          "category": 1,
          "message": "This JSX tag's 'children' prop expects a single child of type 'number | undefined', but multiple children were provided."
        },
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
                  "code": 2746,
                  "category": 1,
                  "message": "This JSX tag's 'children' prop expects a single child of type 'number | undefined', but multiple children were provided."
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
  "size": 2074
}

