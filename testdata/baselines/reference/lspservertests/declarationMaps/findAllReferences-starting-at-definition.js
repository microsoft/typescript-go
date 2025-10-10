UseCaseSensitiveFileNames: false
//// [/home/src/projects/project/a/a.ts] *new* 
export function fnA() {}
export interface IfaceA {}
export const instanceA: IfaceA = {};
//// [/home/src/projects/project/a/bin/a.d.ts] *new* 
export declare function fnA(): void;
export interface IfaceA {
}
export declare const instanceA: IfaceA;
//# sourceMappingURL=a.d.ts.map
//// [/home/src/projects/project/a/bin/a.d.ts.map] *new* 
{
    "version": 3,
    "file": "a.d.ts",
    "sourceRoot": "",
    "sources": ["../a.ts"],
    "names": [],
    "mappings": "AAAA,wBAAgB,GAAG,SAAK;AACxB,MAAM,WAAW,MAAM;CAAG;AAC1B,eAAO,MAAM,SAAS,EAAE,MAAW,CAAC"
}
//// [/home/src/projects/project/a/tsconfig.json] *new* 
{
    "compilerOptions": {
        "outDir": "bin",
        "declarationMap": true,
        "composite": true
    }
}
//// [/home/src/projects/project/b/bin/b.d.ts] *new* 
export declare function fnB(): void;
//# sourceMappingURL=b.d.ts.map
//// [/home/src/projects/project/b/bin/b.d.ts.map] *new* 
{
    "version": 3,
    "file": "b.d.ts",
    "sourceRoot": "",
    "sources": ["../b.ts"],
    "names": [],
    "mappings": "AAAA,wBAAgB,GAAG,SAAK"
}
//// [/home/src/projects/project/b/tsconfig.json] *new* 
{
    "compilerOptions": {
        "outDir": "bin",
        "declarationMap": true,
        "composite": true
    }
}
//// [/home/src/projects/project/dummy/dummy.ts] *new* 
export const a = 10;
//// [/home/src/projects/project/dummy/tsconfig.json] *new* 
{}
//// [/home/src/projects/project/user/user.ts] *new* 
import * as a from "../a/bin/a";
import * as b from "../b/bin/b";
export function fnUser() { a.fnA(); b.fnB(); a.instanceA; }

{
  "method": "textDocument/didOpen",
  "params": {
    "textDocument": {
      "uri": "file:///home/src/projects/project/user/user.ts",
      "languageId": "typescript",
      "version": 0,
      "text": "import * as a from \"../a/bin/a\";\nimport * as b from \"../b/bin/b\";\nexport function fnUser() { a.fnA(); b.fnB(); a.instanceA; }"
    }
  }
}
Projects::
  [/dev/null/inferred] *new*
    /home/src/projects/project/a/bin/a.d.ts  
    /home/src/projects/project/b/bin/b.d.ts  
    /home/src/projects/project/user/user.ts  
Open Files::
  [/home/src/projects/project/user/user.ts] *new*
    /dev/null/inferred  (default) 
Config File Names::
  [/home/src/projects/project/user/user.ts] *new*
    NearestConfigFileName: 
{
  "method": "textDocument/didOpen",
  "params": {
    "textDocument": {
      "uri": "file:///home/src/projects/project/a/a.ts",
      "languageId": "typescript",
      "version": 0,
      "text": "export function fnA() {}\nexport interface IfaceA {}\nexport const instanceA: IfaceA = {};"
    }
  }
}
Projects::
  [/dev/null/inferred] 
    /home/src/projects/project/a/bin/a.d.ts  
    /home/src/projects/project/b/bin/b.d.ts  
    /home/src/projects/project/user/user.ts  
  [/home/src/projects/project/a/tsconfig.json] *new*
    /home/src/projects/project/a/a.ts  
Open Files::
  [/home/src/projects/project/a/a.ts] *new*
    /home/src/projects/project/a/tsconfig.json  (default) 
  [/home/src/projects/project/user/user.ts] 
    /dev/null/inferred  (default) 
Config::
  [/home/src/projects/project/a/tsconfig.json] *new*
    RetainingProjects:
      /home/src/projects/project/a/tsconfig.json  
    RetainingOpenFiles:
      /home/src/projects/project/a/a.ts  
Config File Names::
  [/home/src/projects/project/a/a.ts] *new*
    NearestConfigFileName: /home/src/projects/project/a/tsconfig.json
    Ancestors:
      /home/src/projects/project/a/tsconfig.json   
  [/home/src/projects/project/user/user.ts] 
    NearestConfigFileName: 
{
  "method": "textDocument/references",
  "params": {
    "textDocument": {
      "uri": "file:///home/src/projects/project/a/a.ts"
    },
    "position": {
      "line": 0,
      "character": 16
    },
    "context": {
      "includeDeclaration": false
    }
  }
}
// === /home/src/projects/project/a/a.ts ===
// export function /*FIND ALL REFS*/[|fnA|]() {}
// export interface IfaceA {}
// export const instanceA: IfaceA = {};
{
  "method": "textDocument/didClose",
  "params": {
    "textDocument": {
      "uri": "file:///home/src/projects/project/user/user.ts"
    }
  }
}
Open Files::
  [/home/src/projects/project/a/a.ts] 
    /home/src/projects/project/a/tsconfig.json  (default) 
  [/home/src/projects/project/user/user.ts] *closed*
{
  "method": "textDocument/didOpen",
  "params": {
    "textDocument": {
      "uri": "file:///home/src/projects/project/dummy/dummy.ts",
      "languageId": "typescript",
      "version": 0,
      "text": "export const a = 10;"
    }
  }
}
Projects::
  [/dev/null/inferred] *deleted*
    /home/src/projects/project/a/bin/a.d.ts  
    /home/src/projects/project/b/bin/b.d.ts  
    /home/src/projects/project/user/user.ts  
  [/home/src/projects/project/a/tsconfig.json] 
    /home/src/projects/project/a/a.ts  
  [/home/src/projects/project/dummy/tsconfig.json] *new*
    /home/src/projects/project/dummy/dummy.ts  
Open Files::
  [/home/src/projects/project/a/a.ts] 
    /home/src/projects/project/a/tsconfig.json  (default) 
  [/home/src/projects/project/dummy/dummy.ts] *new*
    /home/src/projects/project/dummy/tsconfig.json  (default) 
Config::
  [/home/src/projects/project/a/tsconfig.json] 
    RetainingProjects:
      /home/src/projects/project/a/tsconfig.json  
    RetainingOpenFiles:
      /home/src/projects/project/a/a.ts  
  [/home/src/projects/project/dummy/tsconfig.json] *new*
    RetainingProjects:
      /home/src/projects/project/dummy/tsconfig.json  
    RetainingOpenFiles:
      /home/src/projects/project/dummy/dummy.ts  
Config File Names::
  [/home/src/projects/project/a/a.ts] 
    NearestConfigFileName: /home/src/projects/project/a/tsconfig.json
    Ancestors:
      /home/src/projects/project/a/tsconfig.json   
  [/home/src/projects/project/dummy/dummy.ts] *new*
    NearestConfigFileName: /home/src/projects/project/dummy/tsconfig.json
  [/home/src/projects/project/user/user.ts] *deleted*
