Fs::
//// [/base.json]
{
  "files": [
    "types/luxon.d.ts",
    "types/express.d.ts"
  ],
  "compilerOptions": {
    "target": "es2017"
  }
}

//// [/src/main.ts]
export {}

//// [/tsconfig.json]
{
  "extends": "./base.json",
  "files": [
    "src/main.ts"
  ],
  "compilerOptions": {
    "outDir": "dist"
  }
}

//// [/types/express.d.ts]
export {}

//// [/types/luxon.d.ts]
export {}


configFileName:: tsconfig.json
CompilerOptions::
{
  "outDir": "/dist",
  "target": 4,
  "configFilePath": "/tsconfig.json"
}

TypeAcquisition::
{}

FileNames::
/types/luxon.d.ts,/types/express.d.ts,/src/main.ts
Errors::

