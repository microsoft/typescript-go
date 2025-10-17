UseCaseSensitiveFileNames: false
//// [/user/username/projects/myproject/a/index.ts] *new* 
import { B } from "../b/lib";
const b: B = new B();
//// [/user/username/projects/myproject/a/tsconfig.json] *new* 
{
    "disableReferencedProjectLoad": true,
    "disableSourceOfProjectReferenceRedirect": true,
    "composite": true
}
//// [/user/username/projects/myproject/b/helper.ts] *new* 
import { B } from ".";
const b: B = new B();
//// [/user/username/projects/myproject/b/index.ts] *new* 
export class B {
    M() {}
}
//// [/user/username/projects/myproject/b/lib/index.d.ts] *new* 
export declare class B {
    M(): void;
}
//# sourceMappingURL=index.d.ts.map
//// [/user/username/projects/myproject/b/lib/index.d.ts.map] *new* 
{
    "version": 3,
    "file": "index.d.ts",
    "sourceRoot": "",
    "sources": ["../index.ts"],
    "names": [],
    "mappings": "AAAA,qBAAa,CAAC;IACV,CAAC;CACJ"
}
//// [/user/username/projects/myproject/b/tsconfig.json] *new* 
{
    "declarationMap": true,
    "outDir": "lib",
    "composite": true,
}

{
  "method": "textDocument/didOpen",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/myproject/a/index.ts",
      "languageId": "typescript",
      "version": 0,
      "text": "import { B } from \"../b/lib\";\nconst b: B = new B();"
    }
  }
}
Projects::
  [/user/username/projects/myproject/a/tsconfig.json] *new*
    /user/username/projects/myproject/b/lib/index.d.ts  
    /user/username/projects/myproject/a/index.ts        
Open Files::
  [/user/username/projects/myproject/a/index.ts] *new*
    /user/username/projects/myproject/a/tsconfig.json  (default) 
Config::
  [/user/username/projects/myproject/a/tsconfig.json] *new*
    RetainingProjects:
      /user/username/projects/myproject/a/tsconfig.json  
    RetainingOpenFiles:
      /user/username/projects/myproject/a/index.ts  
Config File Names::
  [/user/username/projects/myproject/a/index.ts] *new*
    NearestConfigFileName: /user/username/projects/myproject/a/tsconfig.json
{
  "method": "textDocument/didOpen",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/myproject/b/helper.ts",
      "languageId": "typescript",
      "version": 0,
      "text": "import { B } from \".\";\nconst b: B = new B();"
    }
  }
}
Projects::
  [/user/username/projects/myproject/a/tsconfig.json] 
    /user/username/projects/myproject/b/lib/index.d.ts  
    /user/username/projects/myproject/a/index.ts        
  [/user/username/projects/myproject/b/tsconfig.json] *new*
    /user/username/projects/myproject/b/index.ts        
    /user/username/projects/myproject/b/helper.ts       
    /user/username/projects/myproject/b/lib/index.d.ts  
Open Files::
  [/user/username/projects/myproject/a/index.ts] 
    /user/username/projects/myproject/a/tsconfig.json  (default) 
  [/user/username/projects/myproject/b/helper.ts] *new*
    /user/username/projects/myproject/b/tsconfig.json  (default) 
Config::
  [/user/username/projects/myproject/a/tsconfig.json] 
    RetainingProjects:
      /user/username/projects/myproject/a/tsconfig.json  
    RetainingOpenFiles:
      /user/username/projects/myproject/a/index.ts  
  [/user/username/projects/myproject/b/tsconfig.json] *new*
    RetainingProjects:
      /user/username/projects/myproject/b/tsconfig.json  
    RetainingOpenFiles:
      /user/username/projects/myproject/b/helper.ts  
Config File Names::
  [/user/username/projects/myproject/a/index.ts] 
    NearestConfigFileName: /user/username/projects/myproject/a/tsconfig.json
  [/user/username/projects/myproject/b/helper.ts] *new*
    NearestConfigFileName: /user/username/projects/myproject/b/tsconfig.json
{
  "method": "textDocument/references",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/myproject/a/index.ts"
    },
    "position": {
      "line": 1,
      "character": 9
    },
    "context": {
      "includeDeclaration": false
    }
  }
}
// === /user/username/projects/myproject/a/index.ts ===
// import { [|B|] } from "../b/lib";
// const b: /*FIND ALL REFS*/[|B|] = new [|B|]();

// === /user/username/projects/myproject/b/helper.ts ===
// import { [|B|] } from ".";
// const b: [|B|] = new [|B|]();

// === /user/username/projects/myproject/b/index.ts ===
// export class [|B|] {
//     M() {}
// }

// === /user/username/projects/myproject/b/lib/index.d.ts ===
// export declare class [|B|] {
//     M(): void;
// }
// //# sourceMappingURL=index.d.ts.map
