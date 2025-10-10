UseCaseSensitiveFileNames: false
//// [/user/username/projects/myproject/indirect3/main.ts] *new* 
import { foo } from '../target/src/main';
foo()
export function bar() {}
//// [/user/username/projects/myproject/indirect3/tsconfig.json] *new* 
{ }
//// [/user/username/projects/myproject/src/helpers/functions.ts] *new* 
export function foo() { return 1; }
//// [/user/username/projects/myproject/src/main.ts] *new* 
import { foo } from './helpers/functions';
foo()
//// [/user/username/projects/myproject/tsconfig-src.json] *new* 
{
    "compilerOptions": {
        "composite": true,
        "outDir": "./target",
    },
    "include": ["./src/**/*"]
}
//// [/user/username/projects/myproject/tsconfig.json] *new* 
{
    "compilerOptions": {

    },
    "files": [],
    "references": [
        { "path": "./tsconfig-src.json" }
    ]
}
//// [/user/username/workspaces/dummy/dummy.ts] *new* 
const x = 1;
//// [/user/username/workspaces/dummy/tsconfig.json] *new* 
{ }

{
  "method": "textDocument/didOpen",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/myproject/src/main.ts",
      "languageId": "typescript",
      "version": 0,
      "text": "import { foo } from './helpers/functions';\nfoo()"
    }
  }
}
Projects::
  [/user/username/projects/myproject/tsconfig-src.json] *new*
    /user/username/projects/myproject/src/helpers/functions.ts  
    /user/username/projects/myproject/src/main.ts               
Open Files::
  [/user/username/projects/myproject/src/main.ts] *new*
    /user/username/projects/myproject/tsconfig-src.json  (default) 
Config::
  [/user/username/projects/myproject/tsconfig-src.json] *new*
    RetainingProjects:
      /user/username/projects/myproject/tsconfig-src.json  
    RetainingOpenFiles:
      /user/username/projects/myproject/src/main.ts  
  [/user/username/projects/myproject/tsconfig.json] *new*
    RetainingOpenFiles:
      /user/username/projects/myproject/src/main.ts  
Config File Names::
  [/user/username/projects/myproject/src/main.ts] *new*
    NearestConfigFileName: /user/username/projects/myproject/tsconfig.json
    Ancestors:
      /user/username/projects/myproject/tsconfig-src.json   
{
  "method": "textDocument/didOpen",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/workspaces/dummy/dummy.ts",
      "languageId": "typescript",
      "version": 0,
      "text": "const x = 1;"
    }
  }
}
Projects::
  [/user/username/projects/myproject/tsconfig-src.json] 
    /user/username/projects/myproject/src/helpers/functions.ts  
    /user/username/projects/myproject/src/main.ts               
  [/user/username/workspaces/dummy/tsconfig.json] *new*
    /user/username/workspaces/dummy/dummy.ts  
Open Files::
  [/user/username/projects/myproject/src/main.ts] 
    /user/username/projects/myproject/tsconfig-src.json  (default) 
  [/user/username/workspaces/dummy/dummy.ts] *new*
    /user/username/workspaces/dummy/tsconfig.json  (default) 
Config::
  [/user/username/projects/myproject/tsconfig-src.json] 
    RetainingProjects:
      /user/username/projects/myproject/tsconfig-src.json  
    RetainingOpenFiles:
      /user/username/projects/myproject/src/main.ts  
  [/user/username/projects/myproject/tsconfig.json] 
    RetainingOpenFiles:
      /user/username/projects/myproject/src/main.ts  
  [/user/username/workspaces/dummy/tsconfig.json] *new*
    RetainingProjects:
      /user/username/workspaces/dummy/tsconfig.json  
    RetainingOpenFiles:
      /user/username/workspaces/dummy/dummy.ts  
Config File Names::
  [/user/username/projects/myproject/src/main.ts] 
    NearestConfigFileName: /user/username/projects/myproject/tsconfig.json
    Ancestors:
      /user/username/projects/myproject/tsconfig-src.json   
  [/user/username/workspaces/dummy/dummy.ts] *new*
    NearestConfigFileName: /user/username/workspaces/dummy/tsconfig.json
{
  "method": "textDocument/didClose",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/workspaces/dummy/dummy.ts"
    }
  }
}
Open Files::
  [/user/username/projects/myproject/src/main.ts] 
    /user/username/projects/myproject/tsconfig-src.json  (default) 
  [/user/username/workspaces/dummy/dummy.ts] *closed*
{
  "method": "textDocument/didClose",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/myproject/src/main.ts"
    }
  }
}
Open Files::
  [/user/username/projects/myproject/src/main.ts] *closed*
{
  "method": "textDocument/didOpen",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/workspaces/dummy/dummy.ts",
      "languageId": "typescript",
      "version": 0,
      "text": "const x = 1;"
    }
  }
}
Projects::
  [/user/username/projects/myproject/tsconfig-src.json] *deleted*
    /user/username/projects/myproject/src/helpers/functions.ts  
    /user/username/projects/myproject/src/main.ts               
  [/user/username/workspaces/dummy/tsconfig.json] 
    /user/username/workspaces/dummy/dummy.ts  
Open Files::
  [/user/username/workspaces/dummy/dummy.ts] *new*
    /user/username/workspaces/dummy/tsconfig.json  (default) 
Config::
  [/user/username/projects/myproject/tsconfig-src.json] *deleted*
  [/user/username/projects/myproject/tsconfig.json] *deleted*
  [/user/username/workspaces/dummy/tsconfig.json] 
    RetainingProjects:
      /user/username/workspaces/dummy/tsconfig.json  
    RetainingOpenFiles:
      /user/username/workspaces/dummy/dummy.ts  
Config File Names::
  [/user/username/projects/myproject/src/main.ts] *deleted*
  [/user/username/workspaces/dummy/dummy.ts] 
    NearestConfigFileName: /user/username/workspaces/dummy/tsconfig.json
{
  "method": "textDocument/didClose",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/workspaces/dummy/dummy.ts"
    }
  }
}
Open Files::
  [/user/username/workspaces/dummy/dummy.ts] *closed*
{
  "method": "textDocument/didOpen",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/myproject/src/main.ts",
      "languageId": "typescript",
      "version": 0,
      "text": "import { foo } from './helpers/functions';\nfoo()"
    }
  }
}
Projects::
  [/user/username/projects/myproject/tsconfig-src.json] *new*
    /user/username/projects/myproject/src/helpers/functions.ts  
    /user/username/projects/myproject/src/main.ts               
  [/user/username/workspaces/dummy/tsconfig.json] *deleted*
    /user/username/workspaces/dummy/dummy.ts  
Open Files::
  [/user/username/projects/myproject/src/main.ts] *new*
    /user/username/projects/myproject/tsconfig-src.json  (default) 
Config::
  [/user/username/projects/myproject/tsconfig-src.json] *new*
    RetainingProjects:
      /user/username/projects/myproject/tsconfig-src.json  
    RetainingOpenFiles:
      /user/username/projects/myproject/src/main.ts  
  [/user/username/projects/myproject/tsconfig.json] *new*
    RetainingOpenFiles:
      /user/username/projects/myproject/src/main.ts  
  [/user/username/workspaces/dummy/tsconfig.json] *deleted*
Config File Names::
  [/user/username/projects/myproject/src/main.ts] *new*
    NearestConfigFileName: /user/username/projects/myproject/tsconfig.json
    Ancestors:
      /user/username/projects/myproject/tsconfig-src.json   
  [/user/username/workspaces/dummy/dummy.ts] *deleted*
{
  "method": "textDocument/references",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/myproject/src/main.ts"
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
// === /user/username/projects/myproject/src/helpers/functions.ts ===
// export function [|foo|]() { return 1; }

// === /user/username/projects/myproject/src/main.ts ===
// import { [|foo|] } from './helpers/functions';
// /*FIND ALL REFS*/[|foo|]()
{
  "method": "textDocument/didClose",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/myproject/src/main.ts"
    }
  }
}
Open Files::
  [/user/username/projects/myproject/src/main.ts] *closed*
{
  "method": "textDocument/didOpen",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/myproject/indirect3/main.ts",
      "languageId": "typescript",
      "version": 0,
      "text": "import { foo } from '../target/src/main';\nfoo()\nexport function bar() {}"
    }
  }
}
Projects::
  [/user/username/projects/myproject/indirect3/tsconfig.json] *new*
    /user/username/projects/myproject/indirect3/main.ts  
  [/user/username/projects/myproject/tsconfig-src.json] *deleted*
    /user/username/projects/myproject/src/helpers/functions.ts  
    /user/username/projects/myproject/src/main.ts               
Open Files::
  [/user/username/projects/myproject/indirect3/main.ts] *new*
    /user/username/projects/myproject/indirect3/tsconfig.json  (default) 
Config::
  [/user/username/projects/myproject/indirect3/tsconfig.json] *new*
    RetainingProjects:
      /user/username/projects/myproject/indirect3/tsconfig.json  
    RetainingOpenFiles:
      /user/username/projects/myproject/indirect3/main.ts  
  [/user/username/projects/myproject/tsconfig-src.json] *deleted*
  [/user/username/projects/myproject/tsconfig.json] *deleted*
Config File Names::
  [/user/username/projects/myproject/indirect3/main.ts] *new*
    NearestConfigFileName: /user/username/projects/myproject/indirect3/tsconfig.json
  [/user/username/projects/myproject/src/main.ts] *deleted*
{
  "method": "textDocument/references",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/myproject/indirect3/main.ts"
    },
    "position": {
      "line": 0,
      "character": 9
    },
    "context": {
      "includeDeclaration": false
    }
  }
}
// === /user/username/projects/myproject/indirect3/main.ts ===
// import { /*FIND ALL REFS*/[|foo|] } from '../target/src/main';
// [|foo|]()
// export function bar() {}
