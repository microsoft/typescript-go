
currentDirectory::/home/src/workspaces/project
useCaseSensitiveFileNames::true
Input::
//// [/home/src/workspaces/project/tsconfig.json] new file
{
	"compilerOptions": {
		"composite": true,
		"noEmit": true,
	},
	"typeAcquisition": {
		"enable": true,
		"include": ["0.d.ts", "1.d.ts"],
		"exclude": ["0.js", "1.js"],
		"disableFilenameBasedTypeAcquisition": true,
	},
}

ExitStatus:: 2

CompilerOptions::{}
Output::
error TS18003: No inputs were found in config file '/home/src/workspaces/project/tsconfig.json'. Specified 'include' paths were '["**/*"]' and 'exclude' paths were '[]'.


Found 1 error.

//// [/home/src/workspaces/project/tsconfig.json] no change

