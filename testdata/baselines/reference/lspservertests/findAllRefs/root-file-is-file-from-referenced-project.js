UseCaseSensitiveFileNames: false
//// [/user/username/projects/project/src/common/input/keyboard.test.ts] *new* 
import { evaluateKeyboardEvent } from 'common/input/keyboard';
function testEvaluateKeyboardEvent() {
    return evaluateKeyboardEvent();
}
//// [/user/username/projects/project/src/common/input/keyboard.ts] *new* 
function bar() { return "just a random function so .d.ts location doesnt match"; }
export function evaluateKeyboardEvent() { }
//// [/user/username/projects/project/src/common/tsconfig.json] *new* 
{
    "compilerOptions": {
        "composite": true,
        "declarationMap": true,
        "outDir": "../../out",
        "disableSourceOfProjectReferenceRedirect": false,
        "paths": {
            "*": ["../*"],
        },
    },
    "include": ["./**/*"]
}
//// [/user/username/projects/project/src/terminal.ts] *new* 
import { evaluateKeyboardEvent } from 'common/input/keyboard';
function foo() {
    return evaluateKeyboardEvent();
}
//// [/user/username/projects/project/src/tsconfig.json] *new* 
{
    "compilerOptions": {
        "composite": true,
        "declarationMap": true,
        "outDir": "../out",
        "disableSourceOfProjectReferenceRedirect": false,
        "paths": {
            "common/*": ["./common/*"],
        },
        "tsBuildInfoFile": "../out/src.tsconfig.tsbuildinfo"
    },
    "include": ["./**/*"],
    "references": [
        { "path": "./common" },
    ],
}

{
  "method": "textDocument/didOpen",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/project/src/common/input/keyboard.ts",
      "languageId": "typescript",
      "version": 0,
      "text": "function bar() { return \"just a random function so .d.ts location doesnt match\"; }\nexport function evaluateKeyboardEvent() { }"
    }
  }
}
Projects::
  [/user/username/projects/project/src/common/tsconfig.json] *new*
    /user/username/projects/project/src/common/input/keyboard.ts       
    /user/username/projects/project/src/common/input/keyboard.test.ts  
Open Files::
  [/user/username/projects/project/src/common/input/keyboard.ts] *new*
    /user/username/projects/project/src/common/tsconfig.json  (default) 
Config::
  [/user/username/projects/project/src/common/tsconfig.json] *new*
    RetainingProjects:
      /user/username/projects/project/src/common/tsconfig.json  
    RetainingOpenFiles:
      /user/username/projects/project/src/common/input/keyboard.ts  
Config File Names::
  [/user/username/projects/project/src/common/input/keyboard.ts] *new*
    NearestConfigFileName: /user/username/projects/project/src/common/tsconfig.json
{
  "method": "textDocument/didOpen",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/project/src/terminal.ts",
      "languageId": "typescript",
      "version": 0,
      "text": "import { evaluateKeyboardEvent } from 'common/input/keyboard';\nfunction foo() {\n    return evaluateKeyboardEvent();\n}"
    }
  }
}
Projects::
  [/user/username/projects/project/src/common/tsconfig.json] 
    /user/username/projects/project/src/common/input/keyboard.ts       
    /user/username/projects/project/src/common/input/keyboard.test.ts  
  [/user/username/projects/project/src/tsconfig.json] *new*
    /user/username/projects/project/src/common/input/keyboard.ts       
    /user/username/projects/project/src/terminal.ts                    
    /user/username/projects/project/src/common/input/keyboard.test.ts  
Open Files::
  [/user/username/projects/project/src/common/input/keyboard.ts] *modified*
    /user/username/projects/project/src/common/tsconfig.json  (default) 
    /user/username/projects/project/src/tsconfig.json         *new*
  [/user/username/projects/project/src/terminal.ts] *new*
    /user/username/projects/project/src/tsconfig.json  (default) 
Config::
  [/user/username/projects/project/src/common/tsconfig.json] *modified*
    RetainingProjects: *modified*
      /user/username/projects/project/src/common/tsconfig.json  
      /user/username/projects/project/src/tsconfig.json         *new*
    RetainingOpenFiles:
      /user/username/projects/project/src/common/input/keyboard.ts  
  [/user/username/projects/project/src/tsconfig.json] *new*
    RetainingProjects:
      /user/username/projects/project/src/tsconfig.json  
    RetainingOpenFiles:
      /user/username/projects/project/src/terminal.ts  
Config File Names::
  [/user/username/projects/project/src/common/input/keyboard.ts] 
    NearestConfigFileName: /user/username/projects/project/src/common/tsconfig.json
  [/user/username/projects/project/src/terminal.ts] *new*
    NearestConfigFileName: /user/username/projects/project/src/tsconfig.json
{
  "method": "textDocument/references",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/project/src/common/input/keyboard.ts"
    },
    "position": {
      "line": 1,
      "character": 16
    },
    "context": {
      "includeDeclaration": false
    }
  }
}
// === /user/username/projects/project/src/common/input/keyboard.test.ts ===
// import { [|evaluateKeyboardEvent|] } from 'common/input/keyboard';
// function testEvaluateKeyboardEvent() {
//     return [|evaluateKeyboardEvent|]();
// }

// === /user/username/projects/project/src/common/input/keyboard.ts ===
// function bar() { return "just a random function so .d.ts location doesnt match"; }
// export function /*FIND ALL REFS*/[|evaluateKeyboardEvent|]() { }

// === /user/username/projects/project/src/terminal.ts ===
// import { [|evaluateKeyboardEvent|] } from 'common/input/keyboard';
// function foo() {
//     return [|evaluateKeyboardEvent|]();
// }
