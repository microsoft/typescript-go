UseCaseSensitiveFileNames: false
//// [/users/username/projects/a/a.ts] *new* 
import {C} from "./c/fc";
console.log(C)
//// [/users/username/projects/a/c] -> /users/username/projects/c *new*
//// [/users/username/projects/a/tsconfig.json] *new* 
{}
//// [/users/username/projects/b/b.ts] *new* 
import {C} from "../c/fc";
console.log(C)
//// [/users/username/projects/b/c] -> /users/username/projects/c *new*
//// [/users/username/projects/b/tsconfig.json] *new* 
{}
//// [/users/username/projects/c/fc.ts] *new* 
export const C = 42;

{
  "method": "textDocument/didOpen",
  "params": {
    "textDocument": {
      "uri": "file:///users/username/projects/a/a.ts",
      "languageId": "typescript",
      "version": 0,
      "text": "import {C} from \"./c/fc\";\nconsole.log(C)"
    }
  }
}
Projects::
  [/users/username/projects/a/tsconfig.json] *new*
    /users/username/projects/a/c/fc.ts  
    /users/username/projects/a/a.ts     
Open Files::
  [/users/username/projects/a/a.ts] *new*
    /users/username/projects/a/tsconfig.json  (default) 
Config::
  [/users/username/projects/a/tsconfig.json] *new*
    RetainingProjects:
      /users/username/projects/a/tsconfig.json  
    RetainingOpenFiles:
      /users/username/projects/a/a.ts  
Config File Names::
  [/users/username/projects/a/a.ts] *new*
    NearestConfigFileName: /users/username/projects/a/tsconfig.json
{
  "method": "textDocument/didOpen",
  "params": {
    "textDocument": {
      "uri": "file:///users/username/projects/b/b.ts",
      "languageId": "typescript",
      "version": 0,
      "text": "import {C} from \"../c/fc\";\nconsole.log(C)"
    }
  }
}
Projects::
  [/users/username/projects/a/tsconfig.json] 
    /users/username/projects/a/c/fc.ts  
    /users/username/projects/a/a.ts     
  [/users/username/projects/b/tsconfig.json] *new*
    /users/username/projects/c/fc.ts    
    /users/username/projects/b/b.ts     
    /users/username/projects/b/c/fc.ts  
Open Files::
  [/users/username/projects/a/a.ts] 
    /users/username/projects/a/tsconfig.json  (default) 
  [/users/username/projects/b/b.ts] *new*
    /users/username/projects/b/tsconfig.json  (default) 
Config::
  [/users/username/projects/a/tsconfig.json] 
    RetainingProjects:
      /users/username/projects/a/tsconfig.json  
    RetainingOpenFiles:
      /users/username/projects/a/a.ts  
  [/users/username/projects/b/tsconfig.json] *new*
    RetainingProjects:
      /users/username/projects/b/tsconfig.json  
    RetainingOpenFiles:
      /users/username/projects/b/b.ts  
Config File Names::
  [/users/username/projects/a/a.ts] 
    NearestConfigFileName: /users/username/projects/a/tsconfig.json
  [/users/username/projects/b/b.ts] *new*
    NearestConfigFileName: /users/username/projects/b/tsconfig.json
{
  "method": "textDocument/didOpen",
  "params": {
    "textDocument": {
      "uri": "file:///users/username/projects/a/c/fc.ts",
      "languageId": "typescript",
      "version": 0,
      "text": "export const C = 42;"
    }
  }
}
Open Files::
  [/users/username/projects/a/a.ts] 
    /users/username/projects/a/tsconfig.json  (default) 
  [/users/username/projects/a/c/fc.ts] *new*
    /users/username/projects/a/tsconfig.json  (default) 
  [/users/username/projects/b/b.ts] 
    /users/username/projects/b/tsconfig.json  (default) 
Config::
  [/users/username/projects/a/tsconfig.json] *modified*
    RetainingProjects:
      /users/username/projects/a/tsconfig.json  
    RetainingOpenFiles: *modified*
      /users/username/projects/a/a.ts     
      /users/username/projects/a/c/fc.ts  *new*
  [/users/username/projects/b/tsconfig.json] 
    RetainingProjects:
      /users/username/projects/b/tsconfig.json  
    RetainingOpenFiles:
      /users/username/projects/b/b.ts  
Config File Names::
  [/users/username/projects/a/a.ts] 
    NearestConfigFileName: /users/username/projects/a/tsconfig.json
  [/users/username/projects/a/c/fc.ts] *new*
    NearestConfigFileName: /users/username/projects/a/tsconfig.json
  [/users/username/projects/b/b.ts] 
    NearestConfigFileName: /users/username/projects/b/tsconfig.json
{
  "method": "textDocument/didOpen",
  "params": {
    "textDocument": {
      "uri": "file:///users/username/projects/b/c/fc.ts",
      "languageId": "typescript",
      "version": 0,
      "text": "export const C = 42;"
    }
  }
}
Open Files::
  [/users/username/projects/a/a.ts] 
    /users/username/projects/a/tsconfig.json  (default) 
  [/users/username/projects/a/c/fc.ts] 
    /users/username/projects/a/tsconfig.json  (default) 
  [/users/username/projects/b/b.ts] 
    /users/username/projects/b/tsconfig.json  (default) 
  [/users/username/projects/b/c/fc.ts] *new*
    /users/username/projects/b/tsconfig.json  (default) 
Config::
  [/users/username/projects/a/tsconfig.json] 
    RetainingProjects:
      /users/username/projects/a/tsconfig.json  
    RetainingOpenFiles:
      /users/username/projects/a/a.ts     
      /users/username/projects/a/c/fc.ts  
  [/users/username/projects/b/tsconfig.json] *modified*
    RetainingProjects:
      /users/username/projects/b/tsconfig.json  
    RetainingOpenFiles: *modified*
      /users/username/projects/b/b.ts     
      /users/username/projects/b/c/fc.ts  *new*
Config File Names::
  [/users/username/projects/a/a.ts] 
    NearestConfigFileName: /users/username/projects/a/tsconfig.json
  [/users/username/projects/a/c/fc.ts] 
    NearestConfigFileName: /users/username/projects/a/tsconfig.json
  [/users/username/projects/b/b.ts] 
    NearestConfigFileName: /users/username/projects/b/tsconfig.json
  [/users/username/projects/b/c/fc.ts] *new*
    NearestConfigFileName: /users/username/projects/b/tsconfig.json
{
  "method": "textDocument/rename",
  "params": {
    "textDocument": {
      "uri": "file:///users/username/projects/a/c/fc.ts"
    },
    "position": {
      "line": 0,
      "character": 13
    },
    "newName": "?"
  }
}
// === /users/username/projects/a/a.ts ===
// import {[|CRENAME|]} from "./c/fc";
// console.log([|CRENAME|])

// === /users/username/projects/a/c/fc.ts ===
// export const /*RENAME*/C = 42;
