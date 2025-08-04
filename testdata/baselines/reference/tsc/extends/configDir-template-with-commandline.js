
currentDirectory::/home/src/projects/myproject
useCaseSensitiveFileNames::true
Input::--explainFiles --outDir ${configDir}/outDir

ExitStatus:: 2

CompilerOptions::{
    "outDir": "/home/src/projects/myproject/${configDir}/outDir",
    "explainFiles": true
}
Output::
[96msrc/secondary.ts[0m:[93m4[0m:[93m20[0m - [91merror[0m[90m TS2307: [0mCannot find module 'other/sometype2' or its corresponding type declarations.

[7m4[0m  import { k } from "other/sometype2";
[7m [0m [91m                   ~~~~~~~~~~~~~~~~~[0m


Found 1 error in src/secondary.ts[90m:4[0m


