currentDirectory::/home/src/workspaces/project
useCaseSensitiveFileNames::true
Input::
//// [/home/src/workspaces/project/tsconfig.first.json] *new* 
{
    "extends": "./foobar.json",
    "compilerOptions": {
        "composite": true
    }
}
//// [/home/src/workspaces/project/tsconfig.json] *new* 
{
    "compilerOptions": {
        "composite": true
    },
    "references": [
        { "path": "./tsconfig.first.json" },
        { "path": "./tsconfig.second.json" }
    ]
}
//// [/home/src/workspaces/project/tsconfig.second.json] *new* 
{
    "extends": "./foobar.json",
    "compilerOptions": {
        "composite": true
    }
}

tsgo --b
ExitStatus:: Success
Output::

