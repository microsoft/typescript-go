UseCaseSensitiveFileNames: false
//// [/user/username/projects/solution/api/src/server.ts] *new* 
import * as shared from "../../shared/dist"
shared.dog();
//// [/user/username/projects/solution/api/tsconfig.json] *new* 
{
    "compilerOptions": {
        "composite": true,
        "outDir": "dist",
        "rootDir": "src"
    },
    "include": ["src"],
    "references": [{ "path": "../shared" }],
}
//// [/user/username/projects/solution/app/src/app.ts] *new* 
import * as shared from "../../shared/dist"
shared.dog();
//// [/user/username/projects/solution/app/tsconfig.json] *new* 
{
    "compilerOptions": {
        "composite": true,
        "outDir": "dist",
        "rootDir": "src"
    },
    "include": ["src"],
    "references": [{ "path": "../shared" }],
}
//// [/user/username/projects/solution/shared/src/index.ts] *new* 
export const dog = () => { };
//// [/user/username/projects/solution/shared/tsconfig.json] *new* 
{
    "compilerOptions": {
        "composite": true,
        "outDir": "dist",
        "rootDir": "src"
    },
    "include": ["src"],
}
//// [/user/username/projects/solution/tsconfig.json] *new* 
{
    "files": [],
    "references": [
        { "path": "./api" },
        { "path": "./app" },
    ],
}

{
  "method": "textDocument/didOpen",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/solution/api/src/server.ts",
      "languageId": "typescript",
      "version": 0,
      "text": "import * as shared from \"../../shared/dist\"\nshared.dog();"
    }
  }
}
Projects::
  [/user/username/projects/solution/api/tsconfig.json] *new*
    /user/username/projects/solution/shared/src/index.ts  
    /user/username/projects/solution/api/src/server.ts    
Open Files::
  [/user/username/projects/solution/api/src/server.ts] *new*
    /user/username/projects/solution/api/tsconfig.json  (default) 
Config::
  [/user/username/projects/solution/api/tsconfig.json] *new*
    RetainingProjects:
      /user/username/projects/solution/api/tsconfig.json  
    RetainingOpenFiles:
      /user/username/projects/solution/api/src/server.ts  
  [/user/username/projects/solution/shared/tsconfig.json] *new*
    RetainingProjects:
      /user/username/projects/solution/api/tsconfig.json  
Config File Names::
  [/user/username/projects/solution/api/src/server.ts] *new*
    NearestConfigFileName: /user/username/projects/solution/api/tsconfig.json
{
  "method": "textDocument/references",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/solution/api/src/server.ts"
    },
    "position": {
      "line": 1,
      "character": 7
    },
    "context": {
      "includeDeclaration": false
    }
  }
}
Projects::
  [/user/username/projects/solution/api/tsconfig.json] 
    /user/username/projects/solution/shared/src/index.ts  
    /user/username/projects/solution/api/src/server.ts    
  [/user/username/projects/solution/shared/tsconfig.json] *new*
    /user/username/projects/solution/shared/src/index.ts  
Config::
  [/user/username/projects/solution/api/tsconfig.json] 
    RetainingProjects:
      /user/username/projects/solution/api/tsconfig.json  
    RetainingOpenFiles:
      /user/username/projects/solution/api/src/server.ts  
  [/user/username/projects/solution/shared/tsconfig.json] *modified*
    RetainingProjects: *modified*
      /user/username/projects/solution/api/tsconfig.json     
      /user/username/projects/solution/shared/tsconfig.json  *new*
    RetainingOpenFiles: *modified*
      /user/username/projects/solution/shared/src/index.ts  *new*
// === /user/username/projects/solution/api/src/server.ts ===
// import * as shared from "../../shared/dist"
// shared./*FIND ALL REFS*/[|dog|]();

// === /user/username/projects/solution/shared/src/index.ts ===
// export const [|dog|] = () => { };
