currentDirectory::/home/src/workspaces/project
useCaseSensitiveFileNames::true
Input::
//// [/home/src/workspaces/project/tsconfig.json] *new* 
{"files": []}

tsgo --watch
ExitStatus:: Success
Output::
build starting at HH:MM:SS AM
[96mtsconfig.json[0m:[93m1[0m:[93m11[0m - [91merror[0m[90m TS18002: [0mThe 'files' list in config file '/home/src/workspaces/project/tsconfig.json' is empty.

[7m1[0m {"files": []}
[7m [0m [91m          ~~[0m


Found 1 error in tsconfig.json[90m:1[0m

build finished in d.ddds

tsconfig.json::
SemanticDiagnostics::
Signatures::
