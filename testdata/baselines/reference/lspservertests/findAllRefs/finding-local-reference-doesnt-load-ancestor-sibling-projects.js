UseCaseSensitiveFileNames: false
//// [/user/username/projects/solution/compiler/program.ts] *new* 
namespace ts {
    export const program: Program = {
        getSourceFiles: () => [getSourceFile()]
    };
    function getSourceFile() { return "something"; }
}
//// [/user/username/projects/solution/compiler/tsconfig.json] *new* 
{
    "compilerOptions": {
        "composite": true,

    },
    "files": ["./types.ts", "./program.ts"]
}
//// [/user/username/projects/solution/compiler/types.ts] *new* 
namespace ts {
    export interface Program {
        getSourceFiles(): string[];
    }
}
//// [/user/username/projects/solution/services/services.ts] *new* 
/// <reference path="../compiler/types.ts" />
/// <reference path="../compiler/program.ts" />
namespace ts {
    const result = program.getSourceFiles();
}
//// [/user/username/projects/solution/services/tsconfig.json] *new* 
{
    "compilerOptions": {
        "composite": true
    },
    "files": ["./services.ts"],
    "references": [
        { "path": "../compiler" },
    ],
}
//// [/user/username/projects/solution/tsconfig.json] *new* 
{
    "files": [],
    "include": [],
    "references": [
        { "path": "./compiler" },
        { "path": "./services" },
    ],
}

{
  "method": "textDocument/didOpen",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/solution/compiler/program.ts",
      "languageId": "typescript",
      "version": 0,
      "text": "namespace ts {\n    export const program: Program = {\n        getSourceFiles: () => [getSourceFile()]\n    };\n    function getSourceFile() { return \"something\"; }\n}"
    }
  }
}
Projects::
  [/user/username/projects/solution/compiler/tsconfig.json] *new*
    /user/username/projects/solution/compiler/types.ts    
    /user/username/projects/solution/compiler/program.ts  
  [/user/username/projects/solution/tsconfig.json] *new*
Open Files::
  [/user/username/projects/solution/compiler/program.ts] *new*
    /user/username/projects/solution/compiler/tsconfig.json  (default) 
Config::
  [/user/username/projects/solution/compiler/tsconfig.json] *new*
    RetainingProjects:
      /user/username/projects/solution/compiler/tsconfig.json  
    RetainingOpenFiles:
      /user/username/projects/solution/compiler/program.ts  
Config File Names::
  [/user/username/projects/solution/compiler/program.ts] *new*
    NearestConfigFileName: /user/username/projects/solution/compiler/tsconfig.json
    Ancestors:
      /user/username/projects/solution/compiler/tsconfig.json  /user/username/projects/solution/tsconfig.json 
      /user/username/projects/solution/tsconfig.json            
{
  "method": "textDocument/references",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/solution/compiler/program.ts"
    },
    "position": {
      "line": 2,
      "character": 31
    },
    "context": {
      "includeDeclaration": false
    }
  }
}
// === /user/username/projects/solution/compiler/program.ts ===
// namespace ts {
//     export const program: Program = {
//         getSourceFiles: () => [/*FIND ALL REFS*/[|getSourceFile|]()]
//     };
//     function [|getSourceFile|]() { return "something"; }
// }
{
  "method": "textDocument/references",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/solution/compiler/program.ts"
    },
    "position": {
      "line": 2,
      "character": 8
    },
    "context": {
      "includeDeclaration": false
    }
  }
}
Projects::
  [/user/username/projects/solution/compiler/tsconfig.json] 
    /user/username/projects/solution/compiler/types.ts    
    /user/username/projects/solution/compiler/program.ts  
  [/user/username/projects/solution/services/tsconfig.json] *new*
    /user/username/projects/solution/compiler/types.ts     
    /user/username/projects/solution/compiler/program.ts   
    /user/username/projects/solution/services/services.ts  
  [/user/username/projects/solution/tsconfig.json] *modified*
Open Files::
  [/user/username/projects/solution/compiler/program.ts] *modified*
    /user/username/projects/solution/compiler/tsconfig.json  (default) 
    /user/username/projects/solution/services/tsconfig.json  *new*
Config::
  [/user/username/projects/solution/compiler/tsconfig.json] *modified*
    RetainingProjects: *modified*
      /user/username/projects/solution/compiler/tsconfig.json  
      /user/username/projects/solution/services/tsconfig.json  *new*
      /user/username/projects/solution/tsconfig.json           *new*
    RetainingOpenFiles:
      /user/username/projects/solution/compiler/program.ts  
  [/user/username/projects/solution/services/tsconfig.json] *new*
    RetainingProjects:
      /user/username/projects/solution/services/tsconfig.json  
      /user/username/projects/solution/tsconfig.json           
  [/user/username/projects/solution/tsconfig.json] *new*
    RetainingProjects:
      /user/username/projects/solution/tsconfig.json  
// === /user/username/projects/solution/compiler/program.ts ===
// namespace ts {
//     export const program: Program = {
//         /*FIND ALL REFS*/[|getSourceFiles|]: () => [getSourceFile()]
//     };
//     function getSourceFile() { return "something"; }
// }

// === /user/username/projects/solution/compiler/types.ts ===
// namespace ts {
//     export interface Program {
//         [|getSourceFiles|](): string[];
//     }
// }

// === /user/username/projects/solution/services/services.ts ===
// /// <reference path="../compiler/types.ts" />
// /// <reference path="../compiler/program.ts" />
// namespace ts {
//     const result = program.[|getSourceFiles|]();
// }
