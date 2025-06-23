Fs::
//// [/base.json]
{
  "files": [
    "types/luxon.d.ts"
  ],
  "compilerOptions": {
    "target": "es2017"
  }
}

//// [/tsconfig.json]
{
  "extends": "./base.json",
  "compilerOptions": {
    "outDir": "dist"
  }
}

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
/types/luxon.d.ts
Errors::

