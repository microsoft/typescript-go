UseCaseSensitiveFileNames: false
//// [/home/src/projects/project/packages/babel-loader/src/index.ts] *new* 
import type { Foo } from "../../core/src/index.js";
//// [/home/src/projects/project/packages/babel-loader/tsconfig.json] *new* 
{
    "compilerOptions": {
        "target": "ES2018",
        "module": "commonjs",
        "strict": true,
        "esModuleInterop": true,
        "composite": true,
        "rootDir": "src",
        "outDir": "dist"
    },
    "include": ["src"],
    "references": [{"path": "../core"}]
}
//// [/home/src/projects/project/packages/core/src/index.ts] *new* 
import { Bar } from "./loading-indicator.js";
export type Foo = {};
const bar: Bar = {
    prop: 0
}
//// [/home/src/projects/project/packages/core/src/loading-indicator.ts] *new* 
export interface Bar {
    prop: number;
}
const bar: Bar = {
    prop: 1
}
//// [/home/src/projects/project/packages/core/tsconfig.json] *new* 
{
    "compilerOptions": {
        "target": "ES2018",
        "module": "commonjs",
        "strict": true,
        "esModuleInterop": true,
        "composite": true,
        "rootDir": "./src",
        "outDir": "./dist",
    },
    "include": ["./src"]
}

{
  "method": "textDocument/didOpen",
  "params": {
    "textDocument": {
      "uri": "file:///home/src/projects/project/packages/babel-loader/src/index.ts",
      "languageId": "typescript",
      "version": 0,
      "text": "import type { Foo } from \"../../core/src/index.js\";"
    }
  }
}
Projects::
  [/home/src/projects/project/packages/babel-loader/tsconfig.json] *new*
    /home/src/projects/project/packages/core/src/loading-indicator.ts  
    /home/src/projects/project/packages/core/src/index.ts              
    /home/src/projects/project/packages/babel-loader/src/index.ts      
Open Files::
  [/home/src/projects/project/packages/babel-loader/src/index.ts] *new*
    /home/src/projects/project/packages/babel-loader/tsconfig.json  (default) 
Config::
  [/home/src/projects/project/packages/babel-loader/tsconfig.json] *new*
    RetainingProjects:
      /home/src/projects/project/packages/babel-loader/tsconfig.json  
    RetainingOpenFiles:
      /home/src/projects/project/packages/babel-loader/src/index.ts  
  [/home/src/projects/project/packages/core/tsconfig.json] *new*
    RetainingProjects:
      /home/src/projects/project/packages/babel-loader/tsconfig.json  
Config File Names::
  [/home/src/projects/project/packages/babel-loader/src/index.ts] *new*
    NearestConfigFileName: /home/src/projects/project/packages/babel-loader/tsconfig.json
{
  "method": "textDocument/didOpen",
  "params": {
    "textDocument": {
      "uri": "file:///home/src/projects/project/packages/core/src/index.ts",
      "languageId": "typescript",
      "version": 0,
      "text": "import { Bar } from \"./loading-indicator.js\";\nexport type Foo = {};\nconst bar: Bar = {\n    prop: 0\n}"
    }
  }
}
Projects::
  [/home/src/projects/project/packages/babel-loader/tsconfig.json] 
    /home/src/projects/project/packages/core/src/loading-indicator.ts  
    /home/src/projects/project/packages/core/src/index.ts              
    /home/src/projects/project/packages/babel-loader/src/index.ts      
  [/home/src/projects/project/packages/core/tsconfig.json] *new*
    /home/src/projects/project/packages/core/src/loading-indicator.ts  
    /home/src/projects/project/packages/core/src/index.ts              
Open Files::
  [/home/src/projects/project/packages/babel-loader/src/index.ts] 
    /home/src/projects/project/packages/babel-loader/tsconfig.json  (default) 
  [/home/src/projects/project/packages/core/src/index.ts] *new*
    /home/src/projects/project/packages/babel-loader/tsconfig.json  
    /home/src/projects/project/packages/core/tsconfig.json          (default) 
Config::
  [/home/src/projects/project/packages/babel-loader/tsconfig.json] 
    RetainingProjects:
      /home/src/projects/project/packages/babel-loader/tsconfig.json  
    RetainingOpenFiles:
      /home/src/projects/project/packages/babel-loader/src/index.ts  
  [/home/src/projects/project/packages/core/tsconfig.json] *modified*
    RetainingProjects: *modified*
      /home/src/projects/project/packages/babel-loader/tsconfig.json  
      /home/src/projects/project/packages/core/tsconfig.json          *new*
    RetainingOpenFiles: *modified*
      /home/src/projects/project/packages/core/src/index.ts  *new*
Config File Names::
  [/home/src/projects/project/packages/babel-loader/src/index.ts] 
    NearestConfigFileName: /home/src/projects/project/packages/babel-loader/tsconfig.json
  [/home/src/projects/project/packages/core/src/index.ts] *new*
    NearestConfigFileName: /home/src/projects/project/packages/core/tsconfig.json
{
  "method": "textDocument/didChange",
  "params": {
    "textDocument": {
      "uri": "file:///home/src/projects/project/packages/babel-loader/src/index.ts",
      "version": 2
    },
    "contentChanges": [
      {
        "range": {
          "start": {
            "line": 0,
            "character": 0
          },
          "end": {
            "line": 0,
            "character": 0
          }
        },
        "text": "// comment"
      }
    ]
  }
}
{
  "method": "textDocument/references",
  "params": {
    "textDocument": {
      "uri": "file:///home/src/projects/project/packages/core/src/index.ts"
    },
    "position": {
      "line": 3,
      "character": 4
    },
    "context": {
      "includeDeclaration": false
    }
  }
}
Projects::
  [/home/src/projects/project/packages/babel-loader/tsconfig.json] *modified*
    /home/src/projects/project/packages/babel-loader/src/index.ts      *modified*
    /home/src/projects/project/packages/core/src/loading-indicator.ts  *deleted*
    /home/src/projects/project/packages/core/src/index.ts              *deleted*
  [/home/src/projects/project/packages/core/tsconfig.json] 
    /home/src/projects/project/packages/core/src/loading-indicator.ts  
    /home/src/projects/project/packages/core/src/index.ts              
Open Files::
  [/home/src/projects/project/packages/babel-loader/src/index.ts] 
    /home/src/projects/project/packages/babel-loader/tsconfig.json  (default) 
  [/home/src/projects/project/packages/core/src/index.ts] *modified*
    /home/src/projects/project/packages/babel-loader/tsconfig.json  *deleted*
    /home/src/projects/project/packages/core/tsconfig.json          (default) 
// === /home/src/projects/project/packages/core/src/index.ts ===
// import { Bar } from "./loading-indicator.js";
// export type Foo = {};
// const bar: Bar = {
//     /*FIND ALL REFS*/[|prop|]: 0
// }

// === /home/src/projects/project/packages/core/src/loading-indicator.ts ===
// export interface Bar {
//     [|prop|]: number;
// }
// const bar: Bar = {
//     [|prop|]: 1
// }
