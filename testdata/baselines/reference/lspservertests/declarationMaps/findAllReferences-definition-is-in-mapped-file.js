UseCaseSensitiveFileNames: false
//// [/home/src/projects/project/a/a.ts] *new* 
export function f() {}
//// [/home/src/projects/project/a/tsconfig.json] *new* 
{
    "compilerOptions": {
        "outDir": "bin",
        "declarationMap": true,
        "composite": true
    }
}
//// [/home/src/projects/project/b/b.ts] *new* 
import { f } from "../a/bin/a";
f();
//// [/home/src/projects/project/b/tsconfig.json] *new* 
{
    "references": [
        { "path": "../a" }
    ]
}
//// [/home/src/projects/project/bin/a.d.ts] *new* 
export declare function f(): void;
//# sourceMappingURL=a.d.ts.map
//// [/home/src/projects/project/bin/a.d.ts.map] *new* 
{
    "version":3,
    "file":"a.d.ts",
    "sourceRoot":"",
    "sources":["a.ts"],
    "names":[],
    "mappings":"AAAA,wBAAgB,CAAC,SAAK"
}

{
  "method": "textDocument/didOpen",
  "params": {
    "textDocument": {
      "uri": "file:///home/src/projects/project/b/b.ts",
      "languageId": "typescript",
      "version": 0,
      "text": "import { f } from \"../a/bin/a\";\nf();"
    }
  }
}
Projects::
  [/home/src/projects/project/b/tsconfig.json] *new*
    /home/src/projects/project/a/a.ts  
    /home/src/projects/project/b/b.ts  
Open Files::
  [/home/src/projects/project/b/b.ts] *new*
    /home/src/projects/project/b/tsconfig.json  (default) 
Config::
  [/home/src/projects/project/a/tsconfig.json] *new*
    RetainingProjects:
      /home/src/projects/project/b/tsconfig.json  
  [/home/src/projects/project/b/tsconfig.json] *new*
    RetainingProjects:
      /home/src/projects/project/b/tsconfig.json  
    RetainingOpenFiles:
      /home/src/projects/project/b/b.ts  
Config File Names::
  [/home/src/projects/project/b/b.ts] *new*
    NearestConfigFileName: /home/src/projects/project/b/tsconfig.json
{
  "method": "textDocument/references",
  "params": {
    "textDocument": {
      "uri": "file:///home/src/projects/project/b/b.ts"
    },
    "position": {
      "line": 1,
      "character": 0
    },
    "context": {
      "includeDeclaration": false
    }
  }
}
// === /home/src/projects/project/a/a.ts ===
// export function [|f|]() {}

// === /home/src/projects/project/b/b.ts ===
// import { [|f|] } from "../a/bin/a";
// /*FIND ALL REFS*/[|f|]();
