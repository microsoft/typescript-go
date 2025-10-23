UseCaseSensitiveFileNames: false
//// [/user/username/projects/myproject/playground/tests.ts] *new* 
export function foo() {}
//// [/user/username/projects/myproject/playground/tsconfig-json/src/src.ts] *new* 
export function foobar() { }
//// [/user/username/projects/myproject/playground/tsconfig-json/tests/spec.ts] *new* 
export function bar() { }
//// [/user/username/projects/myproject/playground/tsconfig-json/tsconfig.json] *new* 
{
    "include": ["./src"],
}
//// [/user/username/projects/myproject/playground/tsconfig.json] *new* 
{}

{
  "method": "textDocument/didOpen",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/myproject/playground/tests.ts",
      "languageId": "typescript",
      "version": 0,
      "text": "export function foo() {}"
    }
  }
}
Projects::
  [/user/username/projects/myproject/playground/tsconfig.json] *new*
    /user/username/projects/myproject/playground/tests.ts                     
    /user/username/projects/myproject/playground/tsconfig-json/src/src.ts     
    /user/username/projects/myproject/playground/tsconfig-json/tests/spec.ts  
Open Files::
  [/user/username/projects/myproject/playground/tests.ts] *new*
    /user/username/projects/myproject/playground/tsconfig.json  (default) 
Config::
  [/user/username/projects/myproject/playground/tsconfig.json] *new*
    RetainingProjects:
      /user/username/projects/myproject/playground/tsconfig.json  
    RetainingOpenFiles:
      /user/username/projects/myproject/playground/tests.ts  
Config File Names::
  [/user/username/projects/myproject/playground/tests.ts] *new*
    NearestConfigFileName: /user/username/projects/myproject/playground/tsconfig.json
{
  "method": "textDocument/didClose",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/myproject/playground/tests.ts"
    }
  }
}
Open Files::
  [/user/username/projects/myproject/playground/tests.ts] *closed*
{
  "method": "textDocument/didOpen",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/myproject/playground/tsconfig-json/tests/spec.ts",
      "languageId": "typescript",
      "version": 0,
      "text": "export function bar() { }"
    }
  }
}
Open Files::
  [/user/username/projects/myproject/playground/tsconfig-json/tests/spec.ts] *new*
    /user/username/projects/myproject/playground/tsconfig.json  (default) 
Config::
  [/user/username/projects/myproject/playground/tsconfig-json/tsconfig.json] *new*
    RetainingOpenFiles:
      /user/username/projects/myproject/playground/tsconfig-json/tests/spec.ts  
  [/user/username/projects/myproject/playground/tsconfig.json] *modified*
    RetainingProjects:
      /user/username/projects/myproject/playground/tsconfig.json  
    RetainingOpenFiles: *modified*
      /user/username/projects/myproject/playground/tests.ts                     *deleted*
      /user/username/projects/myproject/playground/tsconfig-json/tests/spec.ts  *new*
Config File Names::
  [/user/username/projects/myproject/playground/tests.ts] *deleted*
  [/user/username/projects/myproject/playground/tsconfig-json/tests/spec.ts] *new*
    NearestConfigFileName: /user/username/projects/myproject/playground/tsconfig-json/tsconfig.json
    Ancestors:
      /user/username/projects/myproject/playground/tsconfig-json/tsconfig.json  /user/username/projects/myproject/playground/tsconfig.json 
{
  "method": "textDocument/references",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/myproject/playground/tsconfig-json/tests/spec.ts"
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
// === /user/username/projects/myproject/playground/tsconfig-json/tests/spec.ts ===
// export function /*FIND ALL REFS*/[|bar|]() { }
