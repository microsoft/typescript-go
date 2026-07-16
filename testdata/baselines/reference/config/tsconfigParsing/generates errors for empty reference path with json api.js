Fs::
//// [/apath/a.ts]


//// [/apath/tsconfig.json]
{
                "references": [{ "path": "" }],
                "files": ["a.ts"]
            }


configFileName:: /apath/tsconfig.json
CompilerOptions::
{
  "configFilePath": "/apath/tsconfig.json"
}

TypeAcquisition::
{}

FileNames::
/apath/a.ts
Errors::
[91merror[0m[90m TS18052: [0mA project reference path cannot be an empty string. Did you mean '.'?
