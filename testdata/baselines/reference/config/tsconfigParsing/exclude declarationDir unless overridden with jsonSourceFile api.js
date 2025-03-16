Fs::
//// [/a.ts]


//// [/declarations/a.d.ts]


//// [/tsconfig.json]
{
                "compilerOptions": {
                    "declarationDir": "declarations"
                }
            }


configFileName:: tsconfig.json
CompilerOptions::
{
  "declarationDir": "/declarations",
  "configFilePath": "/tsconfig.json",
  "help": null,
  "all": null
}

FileNames::
/a.ts
Errors::


Fs::
//// [/a.ts]


//// [/declarations/a.d.ts]


//// [/tsconfig.json]
{
                "compilerOptions": {
                    "declarationDir": "declarations"
                },
                "exclude": [ "types" ]
            }


configFileName:: tsconfig.json
CompilerOptions::
{
  "declarationDir": "/declarations",
  "configFilePath": "/tsconfig.json",
  "help": null,
  "all": null
}

FileNames::
/a.ts,/declarations/a.d.ts
Errors::

