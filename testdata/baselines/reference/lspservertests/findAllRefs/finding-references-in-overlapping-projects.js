UseCaseSensitiveFileNames: false
//// [/user/username/projects/solution/a/index.ts] *new* 
export interface I {
    M(): void;
}
//// [/user/username/projects/solution/a/tsconfig.json] *new* 
{
    "compilerOptions": {
        "composite": true,
    },
    "files": ["./index.ts"]
}
//// [/user/username/projects/solution/b/index.ts] *new* 
import { I } from "../a";
export class B implements I {
    M() {}
}
//// [/user/username/projects/solution/b/tsconfig.json] *new* 
{
    "compilerOptions": {
        "composite": true
    },
    "files": ["./index.ts"],
    "references": [
        { "path": "../a" },
    ],
}
//// [/user/username/projects/solution/c/index.ts] *new* 
import { I } from "../a";
import { B } from "../b";
export const C: I = new B();
//// [/user/username/projects/solution/c/tsconfig.json] *new* 
{
    "compilerOptions": {
        "composite": true
    },
    "files": ["./index.ts"],
    "references": [
        { "path": "../b" },
    ],
}
//// [/user/username/projects/solution/d/index.ts] *new* 
import { I } from "../a";
import { C } from "../c";
export const D: I = C;
//// [/user/username/projects/solution/d/tsconfig.json] *new* 
{
    "compilerOptions": {
        "composite": true
    },
    "files": ["./index.ts"],
    "references": [
        { "path": "../c" },
    ],
}
//// [/user/username/projects/solution/tsconfig.json] *new* 
{
    "files": [],
    "include": [],
    "references": [
        { "path": "./a" },
        { "path": "./b" },
        { "path": "./c" },
        { "path": "./d" },
    ],
}

{
  "method": "textDocument/didOpen",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/solution/b/index.ts",
      "languageId": "typescript",
      "version": 0,
      "text": "import { I } from \"../a\";\nexport class B implements I {\n    M() {}\n}"
    }
  }
}
Projects::
  [/user/username/projects/solution/b/tsconfig.json] *new*
    /user/username/projects/solution/a/index.ts  
    /user/username/projects/solution/b/index.ts  
Open Files::
  [/user/username/projects/solution/b/index.ts] *new*
    /user/username/projects/solution/b/tsconfig.json  (default) 
Config::
  [/user/username/projects/solution/a/tsconfig.json] *new*
    RetainingProjects:
      /user/username/projects/solution/b/tsconfig.json  
  [/user/username/projects/solution/b/tsconfig.json] *new*
    RetainingProjects:
      /user/username/projects/solution/b/tsconfig.json  
    RetainingOpenFiles:
      /user/username/projects/solution/b/index.ts  
Config File Names::
  [/user/username/projects/solution/b/index.ts] *new*
    NearestConfigFileName: /user/username/projects/solution/b/tsconfig.json
{
  "method": "textDocument/references",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/solution/b/index.ts"
    },
    "position": {
      "line": 1,
      "character": 26
    },
    "context": {
      "includeDeclaration": false
    }
  }
}
Projects::
  [/user/username/projects/solution/a/tsconfig.json] *new*
    /user/username/projects/solution/a/index.ts  
  [/user/username/projects/solution/b/tsconfig.json] 
    /user/username/projects/solution/a/index.ts  
    /user/username/projects/solution/b/index.ts  
Config::
  [/user/username/projects/solution/a/tsconfig.json] *modified*
    RetainingProjects: *modified*
      /user/username/projects/solution/a/tsconfig.json  *new*
      /user/username/projects/solution/b/tsconfig.json  
    RetainingOpenFiles: *modified*
      /user/username/projects/solution/a/index.ts  *new*
  [/user/username/projects/solution/b/tsconfig.json] 
    RetainingProjects:
      /user/username/projects/solution/b/tsconfig.json  
    RetainingOpenFiles:
      /user/username/projects/solution/b/index.ts  
// === /user/username/projects/solution/a/index.ts ===
// export interface [|I|] {
//     M(): void;
// }

// === /user/username/projects/solution/b/index.ts ===
// import { [|I|] } from "../a";
// export class B implements /*FIND ALL REFS*/[|I|] {
//     M() {}
// }
{
  "method": "textDocument/references",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/solution/b/index.ts"
    },
    "position": {
      "line": 1,
      "character": 26
    },
    "context": {
      "includeDeclaration": false
    }
  }
}
// === /user/username/projects/solution/a/index.ts ===
// export interface [|I|] {
//     M(): void;
// }

// === /user/username/projects/solution/b/index.ts ===
// import { [|I|] } from "../a";
// export class B implements /*FIND ALL REFS*/[|I|] {
//     M() {}
// }
