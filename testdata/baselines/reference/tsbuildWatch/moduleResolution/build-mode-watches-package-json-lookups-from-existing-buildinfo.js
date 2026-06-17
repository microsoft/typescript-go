currentDirectory::/user/username/projects/myproject
useCaseSensitiveFileNames::true
Input::
//// [/user/username/projects/myproject/node_modules/pkg2/index.d.ts] *new* 
export type TheNum = 42;
//// [/user/username/projects/myproject/node_modules/pkg2/package.json] *new* 
{
    "name": "pkg2",
    "version": "1.0.0",
    "types": "index.d.ts"
}
//// [/user/username/projects/myproject/packages/pkg1/index.ts] *new* 
import type { TheNum } from 'pkg2'
export const theNum: TheNum = 42;
//// [/user/username/projects/myproject/packages/pkg1/tsconfig.json] *new* 
{
    "compilerOptions": {
        "outDir": "zzbuild",
    },
}
//// [/user/username/projects/myproject/packages/pkg1/zzbuild/index.js] *new* 
export const theNum = 42;

//// [/user/username/projects/myproject/packages/pkg1/zzbuild/tsconfig.tsbuildinfo] *new* 
{"version":"FakeTSVersion","root":["../index.ts"],"packageJsons":["../../../node_modules/pkg2/package.json"]}
//// [/user/username/projects/myproject/packages/pkg1/zzbuild/tsconfig.tsbuildinfo.readable.baseline.txt] *new* 
{
  "version": "FakeTSVersion",
  "root": [
    {
      "files": [
        "../index.ts"
      ],
      "original": "../index.ts"
    }
  ],
  "packageJsons": [
    "../../../node_modules/pkg2/package.json"
  ],
  "size": 109
}

tsgo -b packages/pkg1 -w --verbose --traceResolution
ExitStatus:: Success
Output::
[2J[3J[H[[90mHH:MM:SS AM[0m] Starting compilation in watch mode...

[[90mHH:MM:SS AM[0m] Projects in this build: 
    * packages/pkg1/tsconfig.json

[[90mHH:MM:SS AM[0m] Project 'packages/pkg1/tsconfig.json' is up to date because newest input 'packages/pkg1/index.ts' is older than output 'packages/pkg1/zzbuild/index.js'

[[90mHH:MM:SS AM[0m] Found 0 errors. Watching for file changes.


Watch Registrations::
Directory watches::
  /user/username/projects/myproject/node_modules/pkg2
  /user/username/projects/myproject/packages/pkg1 (recursive)


Edit [0]:: reports import errors after package is removed
//// [/user/username/projects/myproject/node_modules/pkg2/index.d.ts] *deleted*
//// [/user/username/projects/myproject/node_modules/pkg2/package.json] *deleted*


Output::
[2J[3J[H[[90mHH:MM:SS AM[0m] File change detected. Starting incremental compilation...

[[90mHH:MM:SS AM[0m] Projects in this build: 
    * packages/pkg1/tsconfig.json

[[90mHH:MM:SS AM[0m] Project 'packages/pkg1/tsconfig.json' is up to date because newest input 'packages/pkg1/index.ts' is older than output 'packages/pkg1/zzbuild/index.js'

[[90mHH:MM:SS AM[0m] Found 0 errors. Watching for file changes.


Watch Registrations::
Directory watches::
  /user/username/projects/myproject/node_modules/pkg2
  /user/username/projects/myproject/packages/pkg1 (recursive)
