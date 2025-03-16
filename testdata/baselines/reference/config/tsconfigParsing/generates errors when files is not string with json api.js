Fs::
//// [/apath/a.ts]


//// [/apath/tsconfig.json]
{
  "files": [
    {
      "compilerOptions": {
        "experimentalDecorators": true,
        "allowJs": true
      }
    }
  ]
}


configFileName:: /apath/tsconfig.json
CompilerOptions::
{
  "configFilePath": "/apath/tsconfig.json",
  "help": null,
  "all": null
}

FileNames::

Errors::
[91merror[0m[90m TS5024: [0mCompiler option 'files' requires a value of type string.
