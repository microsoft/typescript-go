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
{
  "method": "textDocument/didClose",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/workspaces/dummy/dummy.ts"
    }
  }
}
{
  "method": "textDocument/didClose",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/myproject/src/main.ts"
    }
  }
}
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
{
  "method": "textDocument/didClose",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/workspaces/dummy/dummy.ts"
    }
  }
}
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
