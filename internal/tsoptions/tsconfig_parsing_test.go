package tsoptions

import (
	// "runtime"

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
type typeParseConfig struct {
	host parseConfigHost
}

func newParseConfigHost(host parseConfigHost) *typeParseConfig {
	return &typeParseConfig{
		host: host,
	}
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
	// {
	// 	title: "returns object with error when json is invalid",
	// 	input: []string{
	// 		"invalid",
	// 	},
	// 	output: []map[string]interface{}{
	// 		{},
	// 		{},
	// 	},
	// },
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
				"files": nil,
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
				"files":      nil,
				"references": nil,
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
			expectedErrors: []string{"No inputs were found in config file '/apath/tsconfig.json'. Specified 'include' paths were '[**/*]' and 'exclude' paths were '[no-prop]'."},
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
				"include": nil,
			},
			expectedErrors: []string{"No inputs were found in config file '/apath/tsconfig.json'. Specified 'include' paths were '[]' and 'exclude' paths were '[no-prop]'."},
		},
	},
	// {
	// 	title: "generates errors for includes with outDir",
	// 	input: testConfig{
	// 		jsonText: `{
	// 	"compilerOptions": {
	// 		"outDir": "./"
	// 	},
	// 	"include": ["**/*"]
	// }`,
	// 		configFileName: "/apath/tsconfig.json",
	// 		basePath:       "/apath",
	// 		allFileList:    []string{"/apath/a.ts"},
	// 	},
	// 	output: verifyConfig{
	// 		fileNames: nil,
	// 		configFile: map[string]interface{}{
	// 			"compilerOptions": map[string]interface{}{
	// 				"outDir": "./",
	// 			},
	// 			"include": []string{"**/*"},
	// 		},
	// 		expectedErrors: []string{"No inputs were found in config file '/apath/tsconfig.json'. Specified 'include' paths were '[**/*]' and 'exclude' paths were '[/apath]'."},
	// 	},
	// },
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

			assert.DeepEqual(t, parseConfigFileContent.FileNames, rec.output.fileNames)
			compareTsConfigOptions(t, ParseRawConfig(parseConfigFileContent.Raw, basePath, parseConfigFileContent.Errors, "tsconfig.json"), ParseRawConfig(rec.output.configFile, basePath, nil, "tsconfig.json"), rec.output.configFile)
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
			assert.DeepEqual(t, parseConfigFileContent.FileNames, rec.output.fileNames)

			compareTsConfigOptions(t, ParseRawConfig(parseConfigFileContent.Raw, basePath, parseConfigFileContent.Errors, "tsconfig.json"), ParseRawConfig(rec.output.configFile, basePath, nil, "tsconfig.json"), rec.output.configFile)
			var actualErrorMessages = []string{}
			for _, error := range parseConfigFileContent.Errors {
				actualErrorMessages = append(actualErrorMessages, error.Message())
			}
			compareDiagnosticMessages(t, actualErrorMessages, rec.output.expectedErrors)
		})
	}
}

func compareTsConfigOptions(t *testing.T, tsConfigOptionsInput tsConfigOptions, tsConfigOptionsResults tsConfigOptions, configFileOutput map[string]interface{}) {
	assert.DeepEqual(t, tsConfigOptionsInput.prop, tsConfigOptionsResults.prop)

	var parsedOutputCompilerOptions *core.CompilerOptions = &core.CompilerOptions{}
	if configFileOutput["compilerOptions"] != nil {
		for key, value := range configFileOutput["compilerOptions"].(map[string]interface{}) {
			parseCompilerOptions(key, value, parsedOutputCompilerOptions)
		}
	}

	assert.DeepEqual(t, tsConfigOptionsInput.compilerOptionsProp, *parsedOutputCompilerOptions)
	// assert.DeepEqual(t, tsConfigOptionsInput.references, tsConfigOptionsResults.references)
}

func compareDiagnosticMessages(t *testing.T, actual []string, expected []string) {
	assert.DeepEqual(t, actual, expected)
}
