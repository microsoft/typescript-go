package tsoptions

import (
	// "runtime"

	"fmt"
	"testing"
	"testing/fstest"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
	"gotest.tools/v3/assert"
	// "encoding/json"
)

//type jsonTexts func() []string
// var jsonTexts = func() []string {
// 	text := []string{[
//         "// Comment",
//         "/* Comment*/",]
// 	}
// 	return text
// }

var jsonTexts = []string{
	// returns empty config for file with only whitespaces
	// `"",
	// 	" ",
	//  	`,
	// 	// returns empty config for file with comments only
	// 	`"// Comment",
	// "/* Comment*/",`,
	// 	// return empty config when file is empty object
	// 	`{}`,
	// returns config object without comments
	// `{ // Excluded files
	// 	"exclude": [
	// 		// Exclude d.ts
	// 		"file.d.ts"
	// 	]
	// }`,
	// `{
	// 	/* Excluded
	// 			Files
	// 	*/
	// 	"exclude": [
	// 		/* multiline comments can be in the middle of a line */"file.d.ts"
	// 	]
	// }`,
	// keeps string content untouched
	// `{
	// 	"exclude": [
	// 		"xx//file.d.ts"
	// 	]
	// }`,
	// `{
	// 	"exclude": [
	// 		"xx/*file.d.ts*/"
	// 	]
	// }`,
	// handles escaped characters in strings correctly
	// `{ //doesn't work
	// 	"exclude": [
	// 		"xx\\"//files"
	// 	]
	// }`,
	// `{
	// 	"exclude": [
	// 		"xx\\\\" // end of line comment
	// 	]
	// }`,
	// returns object when users correctly specify library
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
}

// func TestBaselineParseResult(t *testing.T) {
// 	//var baseline []string = []string{}

// 	for _, jsonText := range jsonTexts {
// 		//baseline = append(baseline, "Input::", jsonText)
// 		parsed := ParseConfigFileTextToJson("/apath/tsconfig.json", jsonText)
// 		config := json.Unmarshal(parsed.config.([]byte), "Config::")
// 		fmt.Println(config)
// 		// s, ok := (parsed.config).([]byte)
// 		// if ok {
// 		// 	json.Unmarshal([]byte(s), &parsed)
// 		// }
// 		// fmt.Println(s)
// 	}
// }

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

// func TestGetParsedCommandJson(t *testing.T) {
// 	for _, test := range parseCommandJson {
// 		//host := newVFSParseConfigHost(test.allFileList, "")
// 		parsed := ParseConfigFileTextToJson(test.configFileName, test.jsonText)
// 		parseConfigFileContent := ParseJsonConfigFileContent(
// 			parsed.config.(map[string]interface{}),
// 			*host,
// 			test.basePath,
// 			//basePath ?? ts.getNormalizedAbsolutePath(ts.getDirectoryPath(configFileName), host.sys.getCurrentDirectory()),
// 			nil,
// 			test.configFileName,
// 			/*resolutionStack*/ nil,
// 			/*extraFileExtensions*/ nil,
// 			/*extendedConfigCache*/ nil,
// 		)
// 		configJson, err := json.Marshal(parseConfigFileContent.Options)
// 		if err != nil {
// 			t.Errorf("Failed to marshal parseConfigFileContent: %v", err)
// 		}
// 		fmt.Println("****************************************************")
// 		fmt.Println(string(configJson))
// 	}
// }

// func TestGetParsedCommandJsonSourceFile(t *testing.T) {

// 	for _, test := range parseCommandJson {
// 		var currentTestFile = make(map[string]string, len(test.allFileList))
// 		for _, file := range test.allFileList {
// 			currentTestFile[file] = ""
// 		}
// 		host := newVFSParseConfigHost(currentTestFile, test.basePath)
// 		parsed := compiler.ParseJSONText(test.configFileName, test.jsonText)
// 		var basePath string
// 		if test.basePath != "" {
// 			basePath = test.basePath
// 		} else {
// 			basePath = tspath.GetNormalizedAbsolutePath(tspath.GetDirectoryPath(test.configFileName), "")
// 		}
// 		var tsConfigSourceFile *tsConfigSourceFile = &tsConfigSourceFile{
// 			sourceFile: parsed,
// 		}
// 		parseConfigFileContent := ParseJsonSourceFileConfigFileContent(
// 			tsConfigSourceFile,
// 			*host,
// 			host.currentDirectory,
// 			nil,
// 			tspath.GetNormalizedAbsolutePath(test.configFileName, basePath), //&test.configFileName,
// 			/*resolutionStack*/ nil,
// 			/*extraFileExtensions*/ nil,
// 			/*extendedConfigCache*/ nil,
// 		)
// 		// k := ParseRawConfig(parseConfigFileContent.Raw)
// 		// l := ParseRawConfig(test.expectedResult)
// 		// assert.DeepEqual(t, k.prop, l)

// 		configJson, err := json.Marshal(parseConfigFileContent.Raw)
// 		if err != nil {
// 			t.Errorf("Failed to marshal parseConfigFileContent: %v", err)
// 		}
// 		// expectedResultJson, err := json.Marshal(test.expectedResult)
// 		// if err != nil {
// 		// 	t.Errorf("Failed to marshal expectedResult: %v", err)
// 		// }
// 		//k := string(expectedResultJson)
// 		// assert.DeepEqual(t, string(configJson), strings.ReplaceAll(k, " ", ""))
// 		// assert.Equal(t, parseConfigFileContent.Errors[0].Message(), test.expectedErrors)
// 		fmt.Println("****************************************************")
// 		fmt.Println(string(configJson))
// 		if parseConfigFileContent.Errors != nil {
// 			fmt.Println("errors: ", parseConfigFileContent.Errors[0].Message())
// 		}
// 		fmt.Println("fileNames: ", parseConfigFileContent.FileNames)
// 		fmt.Println("configFileName: ", tspath.GetNormalizedAbsolutePath(test.configFileName, ""))
// 		fmt.Println("****************************************************")
// 		fmt.Println("")
// 	}
// }

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
			fmt.Printf(rec.title)
			fmt.Println(parseConfigFileContent.FileNames)
			fmt.Println(rec.output.fileNames)
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
			parsed := compiler.ParseJSONText(rec.input.configFileName, rec.input.jsonText)
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
