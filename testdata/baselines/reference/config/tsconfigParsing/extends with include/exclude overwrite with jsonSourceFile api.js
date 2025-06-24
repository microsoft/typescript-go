Fs::
//// [/base.json]
{
  "include": [
    "src/**/*",
    "types/**/*"
  ],
  "exclude": [
    "**/*.test.ts"
  ],
  "compilerOptions": {
    "target": "es2017"
  }
}

//// [/lib/util.ts]
export {}

//// [/src/main.ts]
export {}

//// [/src/spec.spec.ts]
export {}

//// [/src/test.test.ts]
export {}

//// [/tsconfig.json]
{
  "extends": "./base.json",
  "include": [
    "lib/**/*"
  ],
  "exclude": [
    "**/*.spec.ts"
  ],
  "compilerOptions": {
    "outDir": "dist"
  }
}

//// [/types/global.d.ts]
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
/lib/util.ts
Errors::

