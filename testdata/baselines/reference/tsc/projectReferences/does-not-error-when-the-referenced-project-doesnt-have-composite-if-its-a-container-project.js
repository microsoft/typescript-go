currentDirectory::/home/src/workspaces/project
useCaseSensitiveFileNames::true
Input::
//// [/home/src/workspaces/project/primary/a.ts] *new* 
export { };
//// [/home/src/workspaces/project/primary/tsconfig.json] *new* 
{
    "compilerOptions": {
        "composite": false,
        "outDir": "bin",
    }
}
//// [/home/src/workspaces/project/reference/b.ts] *new* 
import * as mod_1 from "../primary/a";
//// [/home/src/workspaces/project/reference/tsconfig.json] *new* 
{
    "compilerOptions": {
        "composite": true,
        "outDir": "bin",
    },
    "files": [ ],
    "references": [{
        "path": "../primary"
    }]
}

tsgo --p reference/tsconfig.json
ExitStatus:: DiagnosticsPresent_OutputsGenerated
Output::
[96mreference/tsconfig.json[0m:[93m7[0m:[93m20[0m - [91merror[0m[90m TS6306: [0mReferenced project '/home/src/workspaces/project/primary' must have setting "composite": true.

[7m7[0m     "references": [{
[7m [0m [91m                   ~[0m
[7m8[0m         "path": "../primary"
[7m [0m [91m~~~~~~~~~~~~~~~~~~~~~~~~~~~~[0m
[7m9[0m     }]
[7m [0m [91m~~~~~[0m


Found 1 error in reference/tsconfig.json[90m:7[0m

//// [/home/src/workspaces/project/reference/bin/tsconfig.tsbuildinfo] *new* 
{"version":"FakeTSVersion","errors":true,"fileInfos":[],"options":{"composite":true,"outDir":"./"}}
//// [/home/src/workspaces/project/reference/bin/tsconfig.tsbuildinfo.readable.baseline.txt] *new* 
{
  "version": "FakeTSVersion",
  "errors": true,
  "fileInfos": [],
  "options": {
    "composite": true,
    "outDir": "./"
  },
  "size": 99
}

reference/tsconfig.json::
SemanticDiagnostics::
Signatures::
