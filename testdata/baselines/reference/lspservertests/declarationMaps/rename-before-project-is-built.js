UseCaseSensitiveFileNames: false
//// [/user/username/projects/myproject/dependency/FnS.ts] *new* 
export function fn1() { }
export function fn2() { }
export function fn3() { }
export function fn4() { }
export function fn5() { }

//// [/user/username/projects/myproject/dependency/tsconfig.json] *new* 
{
    "compilerOptions": {
        "composite": true,
        "declarationMap": true,
        "declarationDir": "../decls"
    }
}
//// [/user/username/projects/myproject/main/main.ts] *new* 
import {
    fn1,
    fn2,
    fn3,
    fn4,
    fn5
} from "../decls/FnS";

fn1();
fn2();
fn3();
fn4();
fn5();
//// [/user/username/projects/myproject/main/tsconfig.json] *new* 
{
    "compilerOptions": {
        "composite": true,
        "declarationMap": true,
    },
    "references": [
        { "path": "../dependency" }
    ]
}
//// [/user/username/projects/myproject/tsconfig.json] *new* 
{
    "references": [
        { "path": "main" }
    ]
}
//// [/user/username/projects/random/random.ts] *new* 
export const a = 10;
//// [/user/username/projects/random/tsconfig.json] *new* 
{}

{
  "method": "textDocument/didOpen",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/myproject/dependency/FnS.ts",
      "languageId": "typescript",
      "version": 0,
      "text": "export function fn1() { }\nexport function fn2() { }\nexport function fn3() { }\nexport function fn4() { }\nexport function fn5() { }\n"
    }
  }
}
Projects::
  [/user/username/projects/myproject/dependency/tsconfig.json] *new*
    /user/username/projects/myproject/dependency/FnS.ts  
  [/user/username/projects/myproject/tsconfig.json] *new*
Open Files::
  [/user/username/projects/myproject/dependency/FnS.ts] *new*
    /user/username/projects/myproject/dependency/tsconfig.json  (default) 
Config::
  [/user/username/projects/myproject/dependency/tsconfig.json] *new*
    RetainingProjects:
      /user/username/projects/myproject/dependency/tsconfig.json  
    RetainingOpenFiles:
      /user/username/projects/myproject/dependency/fns.ts  
Config File Names::
  [/user/username/projects/myproject/dependency/fns.ts] *new*
    NearestConfigFileName: /user/username/projects/myproject/dependency/tsconfig.json
    Ancestors:
      /user/username/projects/myproject/dependency/tsconfig.json  /user/username/projects/myproject/tsconfig.json 
      /user/username/projects/myproject/tsconfig.json              
{
  "method": "textDocument/didOpen",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/random/random.ts",
      "languageId": "typescript",
      "version": 0,
      "text": "export const a = 10;"
    }
  }
}
Projects::
  [/user/username/projects/myproject/dependency/tsconfig.json] 
    /user/username/projects/myproject/dependency/FnS.ts  
  [/user/username/projects/myproject/tsconfig.json] 
  [/user/username/projects/random/tsconfig.json] *new*
    /user/username/projects/random/random.ts  
Open Files::
  [/user/username/projects/myproject/dependency/FnS.ts] 
    /user/username/projects/myproject/dependency/tsconfig.json  (default) 
  [/user/username/projects/random/random.ts] *new*
    /user/username/projects/random/tsconfig.json  (default) 
Config::
  [/user/username/projects/myproject/dependency/tsconfig.json] 
    RetainingProjects:
      /user/username/projects/myproject/dependency/tsconfig.json  
    RetainingOpenFiles:
      /user/username/projects/myproject/dependency/fns.ts  
  [/user/username/projects/random/tsconfig.json] *new*
    RetainingProjects:
      /user/username/projects/random/tsconfig.json  
    RetainingOpenFiles:
      /user/username/projects/random/random.ts  
Config File Names::
  [/user/username/projects/myproject/dependency/fns.ts] 
    NearestConfigFileName: /user/username/projects/myproject/dependency/tsconfig.json
    Ancestors:
      /user/username/projects/myproject/dependency/tsconfig.json  /user/username/projects/myproject/tsconfig.json 
      /user/username/projects/myproject/tsconfig.json              
  [/user/username/projects/random/random.ts] *new*
    NearestConfigFileName: /user/username/projects/random/tsconfig.json
{
  "method": "textDocument/rename",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/myproject/dependency/FnS.ts"
    },
    "position": {
      "line": 2,
      "character": 16
    },
    "newName": "?"
  }
}
Projects::
  [/user/username/projects/myproject/dependency/tsconfig.json] 
    /user/username/projects/myproject/dependency/FnS.ts  
  [/user/username/projects/myproject/main/tsconfig.json] *new*
    /user/username/projects/myproject/dependency/FnS.ts  
    /user/username/projects/myproject/main/main.ts       
  [/user/username/projects/myproject/tsconfig.json] *modified*
    /user/username/projects/myproject/dependency/FnS.ts  *new*
    /user/username/projects/myproject/main/main.ts       *new*
  [/user/username/projects/random/tsconfig.json] 
    /user/username/projects/random/random.ts  
Open Files::
  [/user/username/projects/myproject/dependency/FnS.ts] *modified*
    /user/username/projects/myproject/dependency/tsconfig.json  (default) 
    /user/username/projects/myproject/main/tsconfig.json        *new*
    /user/username/projects/myproject/tsconfig.json             *new*
  [/user/username/projects/random/random.ts] 
    /user/username/projects/random/tsconfig.json  (default) 
Config::
  [/user/username/projects/myproject/dependency/tsconfig.json] *modified*
    RetainingProjects: *modified*
      /user/username/projects/myproject/dependency/tsconfig.json  
      /user/username/projects/myproject/main/tsconfig.json        *new*
      /user/username/projects/myproject/tsconfig.json             *new*
    RetainingOpenFiles:
      /user/username/projects/myproject/dependency/fns.ts  
  [/user/username/projects/myproject/main/tsconfig.json] *new*
    RetainingProjects:
      /user/username/projects/myproject/main/tsconfig.json  
      /user/username/projects/myproject/tsconfig.json       
    RetainingOpenFiles:
      /user/username/projects/myproject/main/main.ts  
  [/user/username/projects/myproject/tsconfig.json] *new*
    RetainingProjects:
      /user/username/projects/myproject/tsconfig.json  
  [/user/username/projects/random/tsconfig.json] 
    RetainingProjects:
      /user/username/projects/random/tsconfig.json  
    RetainingOpenFiles:
      /user/username/projects/random/random.ts  
// === /user/username/projects/myproject/dependency/FnS.ts ===
// export function fn1() { }
// export function fn2() { }
// export function /*RENAME*/[|fn3RENAME|]() { }
// export function fn4() { }
// export function fn5() { }
// 

// === /user/username/projects/myproject/main/main.ts ===
// import {
//     fn1,
//     fn2,
//     [|fn3RENAME|],
//     fn4,
//     fn5
// } from "../decls/FnS";
// 
// fn1();
// fn2();
// [|fn3RENAME|]();
// fn4();
// fn5();
{
  "method": "textDocument/didClose",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/random/random.ts"
    }
  }
}
Open Files::
  [/user/username/projects/myproject/dependency/FnS.ts] 
    /user/username/projects/myproject/dependency/tsconfig.json  (default) 
    /user/username/projects/myproject/main/tsconfig.json        
    /user/username/projects/myproject/tsconfig.json             
  [/user/username/projects/random/random.ts] *closed*
{
  "method": "textDocument/didOpen",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/random/random.ts",
      "languageId": "typescript",
      "version": 0,
      "text": "export const a = 10;"
    }
  }
}
Open Files::
  [/user/username/projects/myproject/dependency/FnS.ts] 
    /user/username/projects/myproject/dependency/tsconfig.json  (default) 
    /user/username/projects/myproject/main/tsconfig.json        
    /user/username/projects/myproject/tsconfig.json             
  [/user/username/projects/random/random.ts] *new*
    /user/username/projects/random/tsconfig.json  (default) 
{
  "method": "textDocument/didClose",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/myproject/dependency/FnS.ts"
    }
  }
}
Open Files::
  [/user/username/projects/myproject/dependency/FnS.ts] *closed*
  [/user/username/projects/random/random.ts] 
    /user/username/projects/random/tsconfig.json  (default) 
{
  "method": "textDocument/didClose",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/random/random.ts"
    }
  }
}
Open Files::
  [/user/username/projects/random/random.ts] *closed*
{
  "method": "textDocument/didOpen",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/random/random.ts",
      "languageId": "typescript",
      "version": 0,
      "text": "export const a = 10;"
    }
  }
}
Projects::
  [/user/username/projects/myproject/dependency/tsconfig.json] *deleted*
    /user/username/projects/myproject/dependency/FnS.ts  
  [/user/username/projects/myproject/main/tsconfig.json] *deleted*
    /user/username/projects/myproject/dependency/FnS.ts  
    /user/username/projects/myproject/main/main.ts       
  [/user/username/projects/myproject/tsconfig.json] *deleted*
    /user/username/projects/myproject/dependency/FnS.ts  
    /user/username/projects/myproject/main/main.ts       
  [/user/username/projects/random/tsconfig.json] 
    /user/username/projects/random/random.ts  
Open Files::
  [/user/username/projects/random/random.ts] *new*
    /user/username/projects/random/tsconfig.json  (default) 
Config::
  [/user/username/projects/myproject/dependency/tsconfig.json] *deleted*
  [/user/username/projects/myproject/main/tsconfig.json] *modified*
    RetainingProjects: *modified*
      /user/username/projects/myproject/main/tsconfig.json  *deleted*
      /user/username/projects/myproject/tsconfig.json       *deleted*
    RetainingOpenFiles:
      /user/username/projects/myproject/main/main.ts  
  [/user/username/projects/myproject/tsconfig.json] *deleted*
  [/user/username/projects/random/tsconfig.json] 
    RetainingProjects:
      /user/username/projects/random/tsconfig.json  
    RetainingOpenFiles:
      /user/username/projects/random/random.ts  
Config File Names::
  [/user/username/projects/myproject/dependency/fns.ts] *deleted*
  [/user/username/projects/random/random.ts] 
    NearestConfigFileName: /user/username/projects/random/tsconfig.json
