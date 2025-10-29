UseCaseSensitiveFileNames: false
//// [/user/username/projects/a/a.d.ts] *new* 
export declare class A {
}
//# sourceMappingURL=a.d.ts.map
//// [/user/username/projects/a/a.d.ts.map] *new* 
{
    "version": 3,
    "file": "a.d.ts",
    "sourceRoot": "",
    "sources": ["./a.ts"],
    "names": [],
    "mappings": "AAAA,qBAAa,CAAC;CAAI"
}
//// [/user/username/projects/a/a.ts] *new* 
export class A { }
//// [/user/username/projects/a/tsconfig.json] *new* 
{}
//// [/user/username/projects/b/b.ts] *new* 
import {A} from "../a/a";
new A();
//// [/user/username/projects/b/tsconfig.json] *new* 
{
    "compilerOptions": {
        "disableSourceOfProjectReferenceRedirect": false
    },
    "references": [
        { "path": "../a" }
    ]
}

{
  "method": "textDocument/didOpen",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/b/b.ts",
      "languageId": "typescript",
      "version": 0,
      "text": "import {A} from \"../a/a\";\nnew A();"
    }
  }
}
Projects::
  [/user/username/projects/b/tsconfig.json] *new*
    /user/username/projects/a/a.ts  
    /user/username/projects/b/b.ts  
Open Files::
  [/user/username/projects/b/b.ts] *new*
    /user/username/projects/b/tsconfig.json  (default) 
Config::
  [/user/username/projects/a/tsconfig.json] *new*
    RetainingProjects:
      /user/username/projects/b/tsconfig.json  
  [/user/username/projects/b/tsconfig.json] *new*
    RetainingProjects:
      /user/username/projects/b/tsconfig.json  
    RetainingOpenFiles:
      /user/username/projects/b/b.ts  
Config File Names::
  [/user/username/projects/b/b.ts] *new*
    NearestConfigFileName: /user/username/projects/b/tsconfig.json
{
  "method": "textDocument/references",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/b/b.ts"
    },
    "position": {
      "line": 1,
      "character": 4
    },
    "context": {
      "includeDeclaration": false
    }
  }
}
Projects::
  [/user/username/projects/a/tsconfig.json] *new*
    /user/username/projects/a/a.ts  
  [/user/username/projects/b/tsconfig.json] 
    /user/username/projects/a/a.ts  
    /user/username/projects/b/b.ts  
Config::
  [/user/username/projects/a/tsconfig.json] *modified*
    RetainingProjects: *modified*
      /user/username/projects/a/tsconfig.json  *new*
      /user/username/projects/b/tsconfig.json  
  [/user/username/projects/b/tsconfig.json] 
    RetainingProjects:
      /user/username/projects/b/tsconfig.json  
    RetainingOpenFiles:
      /user/username/projects/b/b.ts  
// === /user/username/projects/a/a.ts ===
// export class [|A|] { }

// === /user/username/projects/b/b.ts ===
// import {[|A|]} from "../a/a";
// new /*FIND ALL REFS*/[|A|]();
