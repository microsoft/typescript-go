UseCaseSensitiveFileNames: false
//// [/user/username/projects/myproject/core/src/file1.ts] *new* 
export const coreConst = 10;
//// [/user/username/projects/myproject/core/tsconfig.json] *new* 
{
    "compilerOptions": {
        "composite": true,

    },

}
//// [/user/username/projects/myproject/coreRef1/src/file1.ts] *new* 
export const coreRef1Const = 10;
//// [/user/username/projects/myproject/coreRef1/tsconfig.json] *new* 
{
    "compilerOptions": {
        "composite": true,

    },
    { "path": "../core" },\n
}
//// [/user/username/projects/myproject/coreRef2/src/file1.ts] *new* 
export const coreRef2Const = 10;
//// [/user/username/projects/myproject/coreRef2/tsconfig.json] *new* 
{
    "compilerOptions": {
        "composite": true,

    },
    { "path": "../core" },\n
}
//// [/user/username/projects/myproject/coreRef3/src/file1.ts] *new* 
export const coreRef3Const = 10;
//// [/user/username/projects/myproject/coreRef3/tsconfig.json] *new* 
{
    "compilerOptions": {
        "composite": true,

    },
    { "path": "../core" },\n
}
//// [/user/username/projects/myproject/indirect/src/file1.ts] *new* 
export const indirectConst = 10;
//// [/user/username/projects/myproject/indirect/tsconfig.json] *new* 
{
    "compilerOptions": {
        "composite": true,

    },
    { "path": "../coreRef1" },\n
}
//// [/user/username/projects/myproject/indirectDisabledChildLoad1/src/file1.ts] *new* 
export const indirectDisabledChildLoad1Const = 10;
//// [/user/username/projects/myproject/indirectDisabledChildLoad1/tsconfig.json] *new* 
{
    "compilerOptions": {
        "composite": true,
        "disableReferencedProjectLoad": true,
    },
    { "path": "../coreRef2" },\n
}
//// [/user/username/projects/myproject/indirectDisabledChildLoad2/src/file1.ts] *new* 
export const indirectDisabledChildLoad2Const = 10;
//// [/user/username/projects/myproject/indirectDisabledChildLoad2/tsconfig.json] *new* 
{
    "compilerOptions": {
        "composite": true,
        "disableReferencedProjectLoad": true,
    },
    { "path": "../coreRef3" },\n
}
//// [/user/username/projects/myproject/indirectNoCoreRef/src/file1.ts] *new* 
export const indirectNoCoreRefConst = 10;
//// [/user/username/projects/myproject/indirectNoCoreRef/tsconfig.json] *new* 
{
    "compilerOptions": {
        "composite": true,

    },
    { "path": "../noCoreRef2" },\n
}
//// [/user/username/projects/myproject/main/src/file1.ts] *new* 
export const mainConst = 10;
//// [/user/username/projects/myproject/main/tsconfig.json] *new* 
{
    "compilerOptions": {
        "composite": true,

    },
    { "path": "../core" },\n{ "path": "../indirect" },\n{ "path": "../noCoreRef1" },\n{ "path": "../indirectDisabledChildLoad1" },\n{ "path": "../indirectDisabledChildLoad2" },\n{ "path": "../refToCoreRef3" },\n{ "path": "../indirectNoCoreRef" },\n
}
//// [/user/username/projects/myproject/noCoreRef1/src/file1.ts] *new* 
export const noCoreRef1Const = 10;
//// [/user/username/projects/myproject/noCoreRef1/tsconfig.json] *new* 
{
    "compilerOptions": {
        "composite": true,

    },

}
//// [/user/username/projects/myproject/noCoreRef2/src/file1.ts] *new* 
export const noCoreRef2Const = 10;
//// [/user/username/projects/myproject/noCoreRef2/tsconfig.json] *new* 
{
    "compilerOptions": {
        "composite": true,

    },

}
//// [/user/username/projects/myproject/refToCoreRef3/src/file1.ts] *new* 
export const refToCoreRef3Const = 10;
//// [/user/username/projects/myproject/refToCoreRef3/tsconfig.json] *new* 
{
    "compilerOptions": {
        "composite": true,

    },
    { "path": "../coreRef3" },\n
}

{
  "method": "textDocument/didOpen",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/myproject/main/src/file1.ts",
      "languageId": "typescript",
      "version": 0,
      "text": "export const mainConst = 10;"
    }
  }
}
Projects::
  [/user/username/projects/myproject/main/tsconfig.json] *new*
    /user/username/projects/myproject/main/src/file1.ts  
Open Files::
  [/user/username/projects/myproject/main/src/file1.ts] *new*
    /user/username/projects/myproject/main/tsconfig.json  (default) 
Config::
  [/user/username/projects/myproject/main/tsconfig.json] *new*
    RetainingProjects:
      /user/username/projects/myproject/main/tsconfig.json  
    RetainingOpenFiles:
      /user/username/projects/myproject/main/src/file1.ts  
Config File Names::
  [/user/username/projects/myproject/main/src/file1.ts] *new*
    NearestConfigFileName: /user/username/projects/myproject/main/tsconfig.json
{
  "method": "textDocument/didOpen",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/myproject/core/src/file1.ts",
      "languageId": "typescript",
      "version": 0,
      "text": "export const coreConst = 10;"
    }
  }
}
Projects::
  [/user/username/projects/myproject/core/tsconfig.json] *new*
    /user/username/projects/myproject/core/src/file1.ts  
  [/user/username/projects/myproject/main/tsconfig.json] 
    /user/username/projects/myproject/main/src/file1.ts  
Open Files::
  [/user/username/projects/myproject/core/src/file1.ts] *new*
    /user/username/projects/myproject/core/tsconfig.json  (default) 
  [/user/username/projects/myproject/main/src/file1.ts] 
    /user/username/projects/myproject/main/tsconfig.json  (default) 
Config::
  [/user/username/projects/myproject/core/tsconfig.json] *new*
    RetainingProjects:
      /user/username/projects/myproject/core/tsconfig.json  
    RetainingOpenFiles:
      /user/username/projects/myproject/core/src/file1.ts  
  [/user/username/projects/myproject/main/tsconfig.json] 
    RetainingProjects:
      /user/username/projects/myproject/main/tsconfig.json  
    RetainingOpenFiles:
      /user/username/projects/myproject/main/src/file1.ts  
Config File Names::
  [/user/username/projects/myproject/core/src/file1.ts] *new*
    NearestConfigFileName: /user/username/projects/myproject/core/tsconfig.json
  [/user/username/projects/myproject/main/src/file1.ts] 
    NearestConfigFileName: /user/username/projects/myproject/main/tsconfig.json
{
  "method": "textDocument/references",
  "params": {
    "textDocument": {
      "uri": "file:///user/username/projects/myproject/core/src/file1.ts"
    },
    "position": {
      "line": 0,
      "character": 13
    },
    "context": {
      "includeDeclaration": false
    }
  }
}
// === /user/username/projects/myproject/core/src/file1.ts ===
// export const /*FIND ALL REFS*/[|coreConst|] = 10;
