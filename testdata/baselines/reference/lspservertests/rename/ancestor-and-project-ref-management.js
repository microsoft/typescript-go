UseCaseSensitiveFileNames: false
//// [/user/username/projects/container/compositeExec/index.ts] *new* 
import { myConst } from "../lib";
export function getMyConst() {
    return myConst;
}
//// [/user/username/projects/container/compositeExec/tsconfig.json] *new* 
{
    "compilerOptions": {
        "composite": true,
    },
    "files": ["./index.ts"],
    "references": [
        { "path": "../lib" },
    ],
}
//// [/user/username/projects/container/exec/index.ts] *new* 
import { myConst } from "../lib";
export function getMyConst() {
    return myConst;
}
//// [/user/username/projects/container/exec/tsconfig.json] *new* 
{
    "files": ["./index.ts"],
    "references": [
        { "path": "../lib" },
    ],
}
//// [/user/username/projects/container/lib/index.ts] *new* 
export const myConst = 30;
//// [/user/username/projects/container/lib/tsconfig.json] *new* 
{
    "compilerOptions": {
        "composite": true,
    },
    references: [],
    files: [
        "index.ts",
    ],
}
//// [/user/username/projects/container/tsconfig.json] *new* 
{
    "files": [],
    "include": [],
    "references": [
        { "path": "./exec" },
        { "path": "./compositeExec" },
    ],
}
//// [/user/username/projects/temp/temp.ts] *new* 
let x = 10
//// [/user/username/projects/temp/tsconfig.json] *new* 
{}

{
  "method": "textDocument/didOpen",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/container/compositeExec/index.ts",
      "languageId": "typescript",
      "version": 0,
      "text": "import { myConst } from \"../lib\";\nexport function getMyConst() {\n    return myConst;\n}"
    }
  }
}
Projects::
  [/user/username/projects/container/compositeExec/tsconfig.json] *new*
    /user/username/projects/container/lib/index.ts            
    /user/username/projects/container/compositeExec/index.ts  
  [/user/username/projects/container/tsconfig.json] *new*
Open Files::
  [/user/username/projects/container/compositeExec/index.ts] *new*
    /user/username/projects/container/compositeExec/tsconfig.json  (default) 
Config::
  [/user/username/projects/container/compositeExec/tsconfig.json] *new*
    RetainingProjects:
      /user/username/projects/container/compositeexec/tsconfig.json  
    RetainingOpenFiles:
      /user/username/projects/container/compositeexec/index.ts  
  [/user/username/projects/container/lib/tsconfig.json] *new*
    RetainingProjects:
      /user/username/projects/container/compositeexec/tsconfig.json  
Config File Names::
  [/user/username/projects/container/compositeexec/index.ts] *new*
    NearestConfigFileName: /user/username/projects/container/compositeExec/tsconfig.json
    Ancestors:
      /user/username/projects/container/compositeExec/tsconfig.json  /user/username/projects/container/tsconfig.json 
      /user/username/projects/container/tsconfig.json                 
{
  "method": "textDocument/didOpen",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/temp/temp.ts",
      "languageId": "typescript",
      "version": 0,
      "text": "let x = 10"
    }
  }
}
Projects::
  [/user/username/projects/container/compositeExec/tsconfig.json] 
    /user/username/projects/container/lib/index.ts            
    /user/username/projects/container/compositeExec/index.ts  
  [/user/username/projects/container/tsconfig.json] 
  [/user/username/projects/temp/tsconfig.json] *new*
    /user/username/projects/temp/temp.ts  
Open Files::
  [/user/username/projects/container/compositeExec/index.ts] 
    /user/username/projects/container/compositeExec/tsconfig.json  (default) 
  [/user/username/projects/temp/temp.ts] *new*
    /user/username/projects/temp/tsconfig.json  (default) 
Config::
  [/user/username/projects/container/compositeExec/tsconfig.json] 
    RetainingProjects:
      /user/username/projects/container/compositeexec/tsconfig.json  
    RetainingOpenFiles:
      /user/username/projects/container/compositeexec/index.ts  
  [/user/username/projects/container/lib/tsconfig.json] 
    RetainingProjects:
      /user/username/projects/container/compositeexec/tsconfig.json  
  [/user/username/projects/temp/tsconfig.json] *new*
    RetainingProjects:
      /user/username/projects/temp/tsconfig.json  
    RetainingOpenFiles:
      /user/username/projects/temp/temp.ts  
Config File Names::
  [/user/username/projects/container/compositeexec/index.ts] 
    NearestConfigFileName: /user/username/projects/container/compositeExec/tsconfig.json
    Ancestors:
      /user/username/projects/container/compositeExec/tsconfig.json  /user/username/projects/container/tsconfig.json 
      /user/username/projects/container/tsconfig.json                 
  [/user/username/projects/temp/temp.ts] *new*
    NearestConfigFileName: /user/username/projects/temp/tsconfig.json
{
  "method": "textDocument/rename",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/container/compositeExec/index.ts"
    },
    "position": {
      "line": 0,
      "character": 9
    },
    "newName": "?"
  }
}
// === /user/username/projects/container/compositeExec/index.ts ===
// import { /*START PREFIX*/myConst as /*RENAME*/[|myConstRENAME|] } from "../lib";
// export function getMyConst() {
//     return [|myConstRENAME|];
// }
{
  "method": "textDocument/didClose",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/temp/temp.ts"
    }
  }
}
Open Files::
  [/user/username/projects/container/compositeExec/index.ts] 
    /user/username/projects/container/compositeExec/tsconfig.json  (default) 
  [/user/username/projects/temp/temp.ts] *closed*
{
  "method": "textDocument/didOpen",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/temp/temp.ts",
      "languageId": "typescript",
      "version": 0,
      "text": "let x = 10"
    }
  }
}
Open Files::
  [/user/username/projects/container/compositeExec/index.ts] 
    /user/username/projects/container/compositeExec/tsconfig.json  (default) 
  [/user/username/projects/temp/temp.ts] *new*
    /user/username/projects/temp/tsconfig.json  (default) 
{
  "method": "textDocument/didClose",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/container/compositeExec/index.ts"
    }
  }
}
Open Files::
  [/user/username/projects/container/compositeExec/index.ts] *closed*
  [/user/username/projects/temp/temp.ts] 
    /user/username/projects/temp/tsconfig.json  (default) 
{
  "method": "textDocument/didClose",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/temp/temp.ts"
    }
  }
}
Open Files::
  [/user/username/projects/temp/temp.ts] *closed*
{
  "method": "textDocument/didOpen",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/temp/temp.ts",
      "languageId": "typescript",
      "version": 0,
      "text": "let x = 10"
    }
  }
}
Projects::
  [/user/username/projects/container/compositeExec/tsconfig.json] *deleted*
    /user/username/projects/container/lib/index.ts            
    /user/username/projects/container/compositeExec/index.ts  
  [/user/username/projects/container/tsconfig.json] *deleted*
  [/user/username/projects/temp/tsconfig.json] 
    /user/username/projects/temp/temp.ts  
Open Files::
  [/user/username/projects/temp/temp.ts] *new*
    /user/username/projects/temp/tsconfig.json  (default) 
Config::
  [/user/username/projects/container/compositeExec/tsconfig.json] *deleted*
  [/user/username/projects/container/lib/tsconfig.json] *deleted*
  [/user/username/projects/temp/tsconfig.json] 
    RetainingProjects:
      /user/username/projects/temp/tsconfig.json  
    RetainingOpenFiles:
      /user/username/projects/temp/temp.ts  
Config File Names::
  [/user/username/projects/container/compositeexec/index.ts] *deleted*
  [/user/username/projects/temp/temp.ts] 
    NearestConfigFileName: /user/username/projects/temp/tsconfig.json
