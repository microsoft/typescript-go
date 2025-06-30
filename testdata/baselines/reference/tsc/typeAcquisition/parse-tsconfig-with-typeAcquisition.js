
currentDirectory::/home/src/workspaces/project
useCaseSensitiveFileNames::true
Input::
//// [/home/src/workspaces/project/tsconfig.json] new file
{
	"compilerOptions": {
		"composite": true,
		"noEmit": true,
	},
	"typeAcquisition": {
		"enable": true,
		"include": ["0.d.ts", "1.d.ts"],
		"exclude": ["0.js", "1.js"],
		"disableFilenameBasedTypeAcquisition": true,
	},
}

ExitStatus:: 1

CompilerOptions::{}
Output::
[91merror[0m[90m TS18003: [0mNo inputs were found in config file '/home/src/workspaces/project/tsconfig.json'. Specified 'include' paths were '["**/*"]' and 'exclude' paths were '[]'.

Found 1 error.

//// [/home/src/workspaces/project/tsconfig.json] no change
//// [/home/src/workspaces/project/tsconfig.tsbuildinfo] new file
{"Version":"FakeTSVersion","fileNames":["bundled:///libs/lib.d.ts","bundled:///libs/lib.es5.d.ts","bundled:///libs/lib.dom.d.ts","bundled:///libs/lib.webworker.importscripts.d.ts","bundled:///libs/lib.scripthost.d.ts","bundled:///libs/lib.decorators.d.ts","bundled:///libs/lib.decorators.legacy.d.ts"],"fileInfos":["a7297ff837fcdf174a9524925966429eb8e5feecc2cc55cc06574e6b092c1eaa",{"version":"69684132aeb9b5642cbcd9e22dff7818ff0ee1aa831728af0ecf97d3364d5546","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"092c2bfe125ce69dbb1223c85d68d4d2397d7d8411867b5cc03cec902c233763","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"80e18897e5884b6723488d4f5652167e7bb5024f946743134ecc4aa4ee731f89","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"cd034f499c6cdca722b60c04b5b1b78e058487a7085a8e0d6fb50809947ee573","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"8e7f8264d0fb4c5339605a15daadb037bf238c10b654bb3eee14208f860a32ea","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"782dec38049b92d4e85c1585fbea5474a219c6984a35b004963b00beb1aab538","affectsGlobalScope":true,"impliedNodeFormat":1}],"options":{"composite":true},"semanticDiagnosticsPerFile":[1,2,3,4,5,6,7]}
//// [/home/src/workspaces/project/tsconfig.tsbuildinfo.readable.baseline.txt] new file
{
  "Version": "FakeTSVersion",
  "fileNames": [
    "bundled:///libs/lib.d.ts",
    "bundled:///libs/lib.es5.d.ts",
    "bundled:///libs/lib.dom.d.ts",
    "bundled:///libs/lib.webworker.importscripts.d.ts",
    "bundled:///libs/lib.scripthost.d.ts",
    "bundled:///libs/lib.decorators.d.ts",
    "bundled:///libs/lib.decorators.legacy.d.ts"
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
    }
  ],
  "options": {
    "composite": true
  },
  "semanticDiagnosticsPerFile": [
    "bundled:///libs/lib.d.ts",
    "bundled:///libs/lib.es5.d.ts",
    "bundled:///libs/lib.dom.d.ts",
    "bundled:///libs/lib.webworker.importscripts.d.ts",
    "bundled:///libs/lib.scripthost.d.ts",
    "bundled:///libs/lib.decorators.d.ts",
    "bundled:///libs/lib.decorators.legacy.d.ts"
  ],
  "size": 1219
}

