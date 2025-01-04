package tsoptions

import (
	"testing"
	"testing/fstest"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/parser"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
	"gotest.tools/v3/assert"
)

type testConfig struct {
	jsonText       string
	configFileName string
	basePath       string
	allFileList    []string
}
type verifyConfig struct {
	fileNames      []string
	configFile     map[string]interface{}
	expectedErrors []string
}

var parseCommandJson = []testConfig{}

type parseConfigHost interface {
	FS() vfs.FS
}

func fixRoot(path string) string {
	rootLength := tspath.GetRootLength(path)
	if rootLength == 0 {
		return path
	}
	if len(path) == rootLength {
		return "."
	}
	return path[rootLength:]
}

func newVFSParseConfigHost(files map[string]string, currentDirectory string) *VfsParseConfigHost {
	fs := fstest.MapFS{}
	for name, content := range files {
		fs[fixRoot(name)] = &fstest.MapFile{
			Data: []byte(content),
		}
	}
	return &VfsParseConfigHost{
		fs:               vfstest.FromMapFS(fs, true /*useCaseSensitiveFileNames*/),
		currentDirectory: currentDirectory,
	}
}

var baselineParseData = []struct {
	title  string
	input  []string
	output []map[string]interface{}
}{
	{
		title: "returns empty config for file with only whitespaces",
		input: []string{
			"",
			" ",
		},
		output: []map[string]interface{}{
			{},
			{},
		},
	},
	{
		title: "returns empty config for file with comments only",
		input: []string{
			"// Comment",
			"/* Comment*/",
		},
		output: []map[string]interface{}{
			{},
			{},
		},
	},
	{
		title: "returns empty config when config is empty object",
		input: []string{
			`{}`,
		},
		output: []map[string]interface{}{
			{},
			{},
		},
	},
	{
		title: "returns config object without comments",
		input: []string{
			`{ // Excluded files
	        "exclude": [
	            // Exclude d.ts
	            "file.d.ts"
	        ]
	    }`,
			`{
	        /* Excluded
	                Files
	        */
	        "exclude": [
	            /* multiline comments can be in the middle of a line */"file.d.ts"
	        ]
	    }`,
		},
		output: []map[string]interface{}{
			{"exclude": []string{"file.d.ts"}},
			{"exclude": []string{"file.d.ts"}},
		},
	},
	{
		title: "keeps string content untouched",
		input: []string{
			`{
				"exclude": [
					"xx//file.d.ts"
				]
			}`,
			`{
				"exclude": [
					"xx/*file.d.ts*/"
				]
			}`,
		},
		output: []map[string]interface{}{
			{"exclude": []string{"xx//file.d.ts"}},
			{"exclude": []string{"xx/*file.d.ts*/"}},
		},
	},
	// {
	// 	title: "handles escaped characters in strings correctly",
	// 	input: []string{
	// 		`{
	// 			"exclude": [
	// 				"xx\\"//files"
	// 			]
	// 		}`,
	// 		`{
	// 			"exclude": [
	// 				"xx\\\\" // end of line comment
	// 			]
	// 		}`,
	// 	},
	// 	output: []map[string]interface{}{
	// 		{"exclude": []string{"xx\"//files"}},
	// 		{"exclude": []string{"xx\\"}},
	// 	},
	// },
	{
		title: "returns object when users correctly specify library",
		input: []string{
			`{
				"compilerOptions": {
					"lib": ["es5"]
				}
			}`,
			`{
				"compilerOptions": {
					"lib": ["es5", "es6"]
				}
			}`,
		},
		output: []map[string]interface{}{
			{"compilerOptions": map[string]interface{}{"lib": []string{"es5"}}},
			{"compilerOptions": map[string]interface{}{"lib": []string{"es5", "es6"}}},
		},
	},
}

func TestBaselineParseResult(t *testing.T) {
	for _, rec := range baselineParseData {
		t.Run(rec.title, func(t *testing.T) {
			t.Parallel()
			var errors []*ast.Diagnostic
			for index, jsonText := range rec.input {
				parsed, _ := ParseConfigFileTextToJson("/apath/tsconfig.json", "/apath", jsonText, errors)
				assert.DeepEqual(t, parsed, rec.output[index])
			}
		})
	}
}

var data = []struct {
	title  string
	input  testConfig
	output verifyConfig
}{
	{
		title: "ignore dotted files and folders",
		input: testConfig{
			jsonText:       `{}`,
			configFileName: "tsconfig.json",
			basePath:       "/apath",
			allFileList:    []string{"/apath/test.ts", "/apath/.git/a.ts", "/apath/.b.ts", "/apath/..c.ts"},
		},
		output: verifyConfig{
			fileNames:      []string{"/apath/test.ts"},
			configFile:     map[string]interface{}{},
			expectedErrors: []string{},
		},
	},
	{
		title: "allow dotted files and folders when explicitly requested",
		input: testConfig{
			jsonText: `{
							"files": ["/apath/.git/a.ts", "/apath/.b.ts", "/apath/..c.ts"]
						}`,
			configFileName: "tsconfig.json",
			basePath:       "/apath",
			allFileList:    []string{"/apath/test.ts", "/apath/.git/a.ts", "/apath/.b.ts", "/apath/..c.ts"},
		},
		output: verifyConfig{
			fileNames: []string{"/apath/.git/a.ts", "/apath/.b.ts", "/apath/..c.ts"},
			configFile: map[string]interface{}{
				"files": []string{"/apath/.git/a.ts", "/apath/.b.ts", "/apath/..c.ts"},
			},
			expectedErrors: []string{},
		},
	},
	{
		title: "implicitly exclude common package folders",
		input: testConfig{
			jsonText:       `{}`,
			configFileName: "tsconfig.json",
			basePath:       "/",
			allFileList:    []string{"/node_modules/a.ts", "/bower_components/b.ts", "/jspm_packages/c.ts", "/d.ts", "/folder/e.ts"},
		},
		output: verifyConfig{
			fileNames:      []string{"/d.ts", "/folder/e.ts"},
			configFile:     map[string]interface{}{},
			expectedErrors: []string{},
		},
	},
	{
		title: "generates errors for empty files list",
		input: testConfig{
			jsonText: `{
					"files": []
				}`,
			configFileName: "/apath/tsconfig.json",
			basePath:       "/apath",
			allFileList:    []string{"/apath/a.ts"},
		},
		output: verifyConfig{
			fileNames: nil,
			configFile: map[string]interface{}{
				"files": []string{},
			},
			expectedErrors: []string{"The 'files' list in config file '/apath/tsconfig.json' is empty."},
		},
	},
	{
		title: "generates errors for empty files list when no references are provided",
		input: testConfig{
			jsonText: `{
					"files": [],
					"references": []
				}`,
			configFileName: "/apath/tsconfig.json",
			basePath:       "/apath",
			allFileList:    []string{"/apath/a.ts"},
		},
		output: verifyConfig{
			fileNames: nil,
			configFile: map[string]interface{}{
				"files":      []string{},
				"references": []map[string]interface{}{},
			},
			expectedErrors: []string{"The 'files' list in config file '/apath/tsconfig.json' is empty."},
		},
	},
	{
		title: "generates errors for directory with no .ts files",
		input: testConfig{
			jsonText:       `{}`,
			configFileName: "/apath/tsconfig.json",
			basePath:       "/apath",
			allFileList:    []string{"/apath/a.js"},
		},
		output: verifyConfig{
			fileNames:      nil,
			configFile:     map[string]interface{}{},
			expectedErrors: []string{"No inputs were found in config file '/apath/tsconfig.json'. Specified 'include' paths were '[**/*]' and 'exclude' paths were '[]'."},
		},
	},
	{
		title: "generates errors for empty include",
		input: testConfig{
			jsonText: `{
			"include": []
		}`,
			configFileName: "/apath/tsconfig.json",
			basePath:       "tests/cases/unittests",
			allFileList:    []string{"/apath/a.ts"},
		},
		output: verifyConfig{
			fileNames: nil,
			configFile: map[string]interface{}{
				"include": []string{},
			},
			expectedErrors: []string{"No inputs were found in config file '/apath/tsconfig.json'. Specified 'include' paths were '[]' and 'exclude' paths were '[]'."},
		},
	},
	{
		title: "parses tsconfig with compilerOptions, files, include, and exclude",
		input: testConfig{
			jsonText: `{
	                "compilerOptions": {
	                    "outDir": "./dist",
						"strict": true,
						"noImplicitAny": true,
						"target": "ES2017",
						"module": "ESNext",
						"moduleResolution": "bundler,
						"moduleDetection": "auto",
						"jsx": "react",
	                },
	                "files": ["/apath/src/index.ts", "/apath/src/app.ts"],
	                "include": ["/apath/src/**/*"],
	                "exclude": ["/apath/node_modules", "/apath/dist"]
	            }`,
			configFileName: "/apath/tsconfig.json",
			basePath:       "/apath",
			allFileList:    []string{"/apath/src/index.ts", "/apath/src/app.ts", "/apath/node_modules/module.ts", "/apath/dist/output.js"},
		},
		output: verifyConfig{
			fileNames: []string{"/apath/src/index.ts", "/apath/src/app.ts"},
			configFile: map[string]interface{}{
				"compilerOptions": core.CompilerOptions{
					OutDir:        "./dist",
					Strict:        core.TSTrue,
					NoImplicitAny: core.TSTrue,
					Target:        core.ScriptTargetES2017,
					ModuleKind:    core.ModuleKindESNext,
					Jsx:           core.JsxEmitReact,
				},
				"files":   []string{"/apath/src/index.ts", "/apath/src/app.ts"},
				"include": []string{"/apath/src/**/*"},
				"exclude": []string{"/apath/node_modules", "/apath/dist"},
			},
			expectedErrors: []string{},
		},
	},
	{
		title: "generates errors when commandline option is in tsconfig",
		input: testConfig{
			jsonText: `{
				"compilerOptions": {
			"help": true,
		},
	}`,
			configFileName: "/apath/tsconfig.json",
			basePath:       "/apath",
			allFileList:    []string{"/apath/a.ts"},
		},
		output: verifyConfig{
			fileNames: []string{"/apath/a.ts"},
			configFile: map[string]interface{}{
				"compilerOptions": core.CompilerOptions{},
			},
			expectedErrors: []string{"Option 'help' can only be specified on command line."},
		},
	},
	{
		title: "does not generate errors for empty files list when one or more references are provided",
		input: testConfig{
			jsonText: `{
	            "files": [],
	            "references": [{ "path": "/apath" }]
	        }`,
			configFileName: "/apath/tsconfig.json",
			basePath:       "/apath",
			allFileList:    []string{"/apath/a.ts"},
		},
		output: verifyConfig{
			fileNames: nil,
			configFile: map[string]interface{}{
				"files":      []string{},
				"references": []map[string]interface{}{{"path": "/apath"}},
			},
			expectedErrors: []string{},
		},
	},
	{
		title: "exclude outDir",
		input: testConfig{
			jsonText: `{
				"compilerOptions": {
					"outDir": "bin"
				},
			}`,
			configFileName: "tsconfig.json",
			basePath:       "/",
			allFileList:    []string{"/bin/a.ts", "/b.ts"},
		},
		output: verifyConfig{
			fileNames: []string{"/b.ts"},
			configFile: map[string]interface{}{
				"compilerOptions": core.CompilerOptions{
					OutDir: "bin",
				},
			},
			expectedErrors: []string{},
		},
	},
	{
		title: "exclude outDir unless overridden",
		input: testConfig{
			jsonText: `{
				"compilerOptions": {
					"outDir": "bin"
				},
				"exclude": ["obj"],
			}`,
			configFileName: "tsconfig.json",
			basePath:       "/",
			allFileList:    []string{"/bin/a.ts", "/b.ts"},
		},
		output: verifyConfig{
			fileNames: []string{"/b.ts", "/bin/a.ts"},
			configFile: map[string]interface{}{
				"compilerOptions": core.CompilerOptions{
					OutDir: "bin",
				},
				"exclude": []string{"obj"},
			},
			expectedErrors: []string{},
		},
	},
	{
		title: "exclude declarationDir",
		input: testConfig{
			jsonText: `{
				"compilerOptions": {
					"declarationDir": "declarations"
				},
			}`,
			configFileName: "tsconfig.json",
			basePath:       "/",
			allFileList:    []string{"/declarations/a.d.ts", "/a.ts"},
		},
		output: verifyConfig{
			fileNames: []string{"/a.ts"},
			configFile: map[string]interface{}{
				"compilerOptions": core.CompilerOptions{
					DeclarationDir: "declarations",
				},
			},
			expectedErrors: []string{},
		},
	},
	{
		title: "exclude declarationDir unless overridden",
		input: testConfig{
			jsonText: `{
				"compilerOptions": {
					"declarationDir": "declarations"
				},
				"exclude": ["types"],
			}`,
			configFileName: "tsconfig.json",
			basePath:       "/",
			allFileList:    []string{"/declarations/a.d.ts", "/a.ts"},
		},
		output: verifyConfig{
			fileNames: []string{"/a.ts", "/declarations/a.d.ts"},
			configFile: map[string]interface{}{
				"compilerOptions": core.CompilerOptions{
					DeclarationDir: "declarations",
				},
				"exclude": []string{"types"},
			},
			expectedErrors: []string{},
		},
	},
	{
		title: "generates errors for empty directory",
		input: testConfig{
			jsonText: `{
				"compilerOptions": {
					"allowJs": true
				},
			}`,
			configFileName: "/apath/tsconfig.json",
			basePath:       "/apath",
			allFileList:    []string{},
		},
		output: verifyConfig{
			fileNames: nil,
			configFile: map[string]interface{}{
				"compilerOptions": core.CompilerOptions{
					AllowJs: core.TSTrue,
				},
			},
			expectedErrors: []string{"No inputs were found in config file '/apath/tsconfig.json'. Specified 'include' paths were '[**/*]' and 'exclude' paths were '[]'."},
		},
	},
	{
		title: "generates errors for includes with outDir",
		input: testConfig{
			jsonText: `{
				"compilerOptions": {
					"outDir": "./"
				},
				"include": ["**/*"]
			}`,
			configFileName: "/apath/tsconfig.json",
			basePath:       "/apath",
			allFileList:    []string{"/apath/a.ts"},
		},
		output: verifyConfig{
			fileNames: nil,
			configFile: map[string]interface{}{
				"compilerOptions": core.CompilerOptions{
					OutDir: "./",
				},
				"include": []string{"**/*"},
			},
			expectedErrors: []string{"No inputs were found in config file '/apath/tsconfig.json'. Specified 'include' paths were '[**/*]' and 'exclude' paths were '[./]'."},
		},
	},
	{
		title: "generates errors when include is not string",
		input: testConfig{
			jsonText: `{
				"include": [./**/*.ts]
				}`,
			configFileName: "/apath/tsconfig.json",
			basePath:       "/apath",
			allFileList:    []string{"/apath/a.ts"},
		},
		output: verifyConfig{
			fileNames:      nil,
			configFile:     map[string]interface{}{"include": []string{}},
			expectedErrors: []string{"No inputs were found in config file '/apath/tsconfig.json'. Specified 'include' paths were '[]' and 'exclude' paths were '[]'."},
		},
	},
	{
		title: "generates errors when files is not string",
		input: testConfig{
			jsonText: `{
				"files": [{"compilerOptions": {}}]
			}`,
			configFileName: "/apath/tsconfig.json",
			basePath:       "/apath",
			allFileList:    []string{"/apath/a.ts"},
		},
		output: verifyConfig{
			fileNames:      nil,
			configFile:     map[string]interface{}{"files": []string{}},
			expectedErrors: []string{"The 'files' list in config file '/apath/tsconfig.json' is empty."},
		},
	},
}

func TestParsedCommandJson(t *testing.T) {
	for _, rec := range data {
		t.Run(rec.title, func(t *testing.T) {
			t.Parallel()

			var allFileLists = make(map[string]string, len(rec.input.allFileList))
			for _, file := range rec.input.allFileList {
				allFileLists[file] = ""
			}
			host := newVFSParseConfigHost(allFileLists, rec.input.basePath)
			var errors []*ast.Diagnostic
			parsed, _ := ParseConfigFileTextToJson(rec.input.configFileName, rec.input.basePath, rec.input.jsonText, errors)
			var basePath string
			if rec.input.basePath != "" {
				basePath = rec.input.basePath
			} else {
				basePath = tspath.GetNormalizedAbsolutePath(tspath.GetDirectoryPath(rec.input.configFileName), "")
			}
			parseConfigFileContent := ParseJsonConfigFileContent(
				parsed,
				*host,
				basePath,
				nil,
				tspath.GetNormalizedAbsolutePath(rec.input.configFileName, basePath),
				/*resolutionStack*/ nil,
				/*extraFileExtensions*/ nil,
				/*extendedConfigCache*/ nil,
			)

			rawConfigResult := ParseRawConfig(parseConfigFileContent.Raw, basePath, parseConfigFileContent.Errors, "tsconfig.json")
			rawConfigExpected := ParseRawConfig(rec.output.configFile, basePath, nil, "tsconfig.json")

			// Check for file names
			assert.DeepEqual(t, parseConfigFileContent.Options.FileNames, rec.output.fileNames)
			// Check for compiler options
			if rec.output.configFile["compilerOptions"] != nil {
				assert.DeepEqual(t, rawConfigResult.compilerOptionsProp, rec.output.configFile["compilerOptions"])
			}
			// Check for all other options
			assert.DeepEqual(t, rawConfigResult.prop, rawConfigExpected.prop)
			// Check for diagnostics
			var actualErrorMessages = []string{}
			for _, error := range parseConfigFileContent.Errors {
				actualErrorMessages = append(actualErrorMessages, error.Message())
			}
			compareDiagnosticMessages(t, actualErrorMessages, rec.output.expectedErrors)
		})
	}
}

func TestParsedCommandJsonSourceFile(t *testing.T) {
	for _, rec := range data {
		t.Run(rec.title, func(t *testing.T) {
			t.Parallel()

			var allFileLists = make(map[string]string, len(rec.input.allFileList))
			for _, file := range rec.input.allFileList {
				allFileLists[file] = ""
			}
			host := newVFSParseConfigHost(allFileLists, rec.input.basePath)
			parsed := parser.ParseJSONText(rec.input.configFileName, rec.input.jsonText)
			var basePath string
			if rec.input.basePath != "" {
				basePath = rec.input.basePath
			} else {
				basePath = tspath.GetNormalizedAbsolutePath(tspath.GetDirectoryPath(rec.input.configFileName), "")
			}
			var tsConfigSourceFile *tsConfigSourceFile = &tsConfigSourceFile{
				sourceFile: parsed,
			}
			parseConfigFileContent := ParseJsonSourceFileConfigFileContent(
				tsConfigSourceFile,
				*host,
				host.currentDirectory,
				nil,
				tspath.GetNormalizedAbsolutePath(rec.input.configFileName, basePath),
				/*resolutionStack*/ nil,
				/*extraFileExtensions*/ nil,
				/*extendedConfigCache*/ nil,
			)

			rawConfigResult := ParseRawConfig(parseConfigFileContent.Raw, basePath, parseConfigFileContent.Errors, "tsconfig.json")
			rawConfigExpected := ParseRawConfig(rec.output.configFile, basePath, nil, "tsconfig.json")

			// Check for file names
			assert.DeepEqual(t, parseConfigFileContent.Options.FileNames, rec.output.fileNames)
			// Check for compiler options
			if rec.output.configFile["compilerOptions"] != nil {
				assert.DeepEqual(t, rawConfigResult.compilerOptionsProp, rec.output.configFile["compilerOptions"])
			}
			// Check for all other options
			assert.DeepEqual(t, rawConfigResult.prop, rawConfigExpected.prop)
			// Check for diagnostics
			var actualErrorMessages = []string{}
			for _, error := range parseConfigFileContent.Errors {
				actualErrorMessages = append(actualErrorMessages, error.Message())
			}
			compareDiagnosticMessages(t, actualErrorMessages, rec.output.expectedErrors)
		})
	}
}

func compareDiagnosticMessages(t *testing.T, actual []string, expected []string) {
	assert.DeepEqual(t, actual, expected)
}
