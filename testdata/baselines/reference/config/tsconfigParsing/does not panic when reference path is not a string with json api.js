Fs::
//// [/apath/a.ts]


//// [/apath/tsconfig.json]
{
                "files": [],
                "references": [{ "path": 3 }]
            }


configFileName:: /apath/tsconfig.json
CompilerOptions::
{
  "configFilePath": "/apath/tsconfig.json"
}

TypeAcquisition::
{}

FileNames::

Errors::
[91merror[0m[90m TS5024: [0mCompiler option 'reference.path' requires a value of type string.
