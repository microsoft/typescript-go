package tsoptions

import (
	// "runtime"

	"encoding/json"
	"fmt"
	"testing"
	"testing/fstest"

	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
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

func TestBaselineParseResult(t *testing.T) {
	//var baseline []string = []string{}

	for _, jsonText := range jsonTexts {
		//baseline = append(baseline, "Input::", jsonText)
		parsed := ParseConfigFileTextToJson("/apath/tsconfig.json", jsonText)
		config := json.Unmarshal(parsed.config.([]byte), "Config::")
		fmt.Println(config)
		// s, ok := (parsed.config).([]byte)
		// if ok {
		// 	json.Unmarshal([]byte(s), &parsed)
		// }
		// fmt.Println(s)
	}
}

type verifyConfig struct {
	jsonText       string
	configFileName string
	basePath       string
	allFileList    []string
	expectedResult any
	expectedErrors string
}

func TestGetParsedCommandJson(t *testing.T) {
	for _, test := range parseCommandJson {
		host := newVFSParseConfigHost(test.allFileList, "")
		parsed := ParseConfigFileTextToJson(test.configFileName, test.jsonText)
		parseConfigFileContent := ParseJsonConfigFileContent(
			parsed.config.(map[string]interface{}),
			*host,
			test.basePath,
			//basePath ?? ts.getNormalizedAbsolutePath(ts.getDirectoryPath(configFileName), host.sys.getCurrentDirectory()),
			nil,
			test.configFileName,
			/*resolutionStack*/ nil,
			/*extraFileExtensions*/ nil,
			/*extendedConfigCache*/ nil,
		)
		configJson, err := json.Marshal(parseConfigFileContent.Options)
		if err != nil {
			t.Errorf("Failed to marshal parseConfigFileContent: %v", err)
		}
		fmt.Println("****************************************************")
		fmt.Println(string(configJson))
	}
}

func TestGetParsedCommandJsonSourceFile(t *testing.T) {

	//parseConfig := newParseConfigHost(host)
	for _, test := range parseCommandJson {
		host := newVFSParseConfigHost(test.allFileList, "")
		parsed := compiler.ParseJSONText(test.configFileName, test.jsonText)
		var basePath string
		if test.basePath != "" {
			basePath = test.basePath
		} else {
			basePath = tspath.GetNormalizedAbsolutePath(tspath.GetDirectoryPath(test.configFileName), "")
		}
		parseConfigFileContent := ParseJsonSourceFileConfigFileContent(
			parsed,
			*host,
			basePath,
			nil,
			tspath.GetNormalizedAbsolutePath(test.configFileName, basePath), //&test.configFileName,
			/*resolutionStack*/ nil,
			/*extraFileExtensions*/ nil,
			/*extendedConfigCache*/ nil,
		)
		// k := ParseRawConfig(parseConfigFileContent.Raw)
		// l := ParseRawConfig(test.expectedResult)
		// assert.DeepEqual(t, k.prop, l)

		configJson, err := json.Marshal(parseConfigFileContent.Raw)
		if err != nil {
			t.Errorf("Failed to marshal parseConfigFileContent: %v", err)
		}
		// expectedResultJson, err := json.Marshal(test.expectedResult)
		// if err != nil {
		// 	t.Errorf("Failed to marshal expectedResult: %v", err)
		// }
		//k := string(expectedResultJson)
		// assert.DeepEqual(t, string(configJson), strings.ReplaceAll(k, " ", ""))
		// assert.Equal(t, parseConfigFileContent.Errors[0].Message(), test.expectedErrors)
		// fmt.Println("****************************************************")
		fmt.Println(string(configJson))
		if parseConfigFileContent.Errors != nil {
			fmt.Println("errors: ", parseConfigFileContent.Errors[0].Message())
		}
		fmt.Println("fileNames: ", parseConfigFileContent.FileNames)
		fmt.Println("configFileName: ", tspath.GetNormalizedAbsolutePath(test.configFileName, ""))
		fmt.Println("****************************************************")
		fmt.Println("")
	}
}

var parseCommandJson = []verifyConfig{
	//"returns error when tsconfig have excludes"
	// {
	// 	jsonText: `{
	// 	"compilerOptions": {
	// 		"lib": ["es5"]
	// 	},
	// }`,
	// },
	// 		expectedResult: `{
	// "compilerOptions": {
	// 	"lib": ["es5"]
	// },
	// "excludes": [
	// 	"foge.ts"
	// ]
	// }`,
	// 	expectedErrors: "Unknown option 'excludes'. Did you mean 'exclude'?",
	// },
	// {
	// 	jsonText:       `{}`,
	// 	configFileName: "tsconfig.json",
	// 	basePath:       "/apath",
	// 	allFileList:    []string{"/apath/test.ts", "/apath/.git/a.ts", "/apath/.b.ts", "/apath/..c.ts"},
	// },
	// 	configFileName: "tsconfig.json",
	// 	basePath:       "/apath",
	// 	allFileList:    []string{"/apath/test.ts", "/apath/foge.ts"},
	// },
	// {
	// 	jsonText: `{
	// 		"files": []
	// 	}`,
	// 	configFileName: "/apath/tsconfig.json",
	// 	basePath:       "/apath",
	// 	allFileList:    []string{"/apath/a.ts"},
	// },
	// "exclude outDir unless overridden"
	{
		jsonText: `{
			"compilerOptions": {
				"outDir": "bin"
			}
		}`,
		configFileName: "tsconfig.json",
		basePath:       "/",
		allFileList:    []string{"bin/a.ts", "b.ts"},
	},
	// {
	// 	jsonText: `{
	// 		"compilerOptions": {
	// 			"outDir": "bin"
	// 		},
	// 		"exclude": ["obj"]
	// 	}`,
	// 	configFileName: "tsconfig.json",
	// 	basePath:       "/",
	// 	allFileList:    []string{"/bin/a.ts", "/b.ts"},
	// },
	// {
	// 	jsonText: `{
	// 		"files": [],
	// 		"references": [{ "path": "/apath" }]
	// 	}`,
	// 	configFileName: "/apath/tsconfig.json",
	// 	basePath:       "/apath",
	// 	allFileList:    []string{"/apath/a.ts"},
	// },
	// {
	// 	jsonText: `{
	// 		"compilerOptions": {
	// 			"target": "es5",
	// 			"module": "commonjs",
	// 			"lib": ["es2015", "dom"],
	// 			"strict": true,
	// 			"esModuleInterop": true,
	// 			"skipLibCheck": true,
	// 			"forceConsistentCasingInFileNames": true
	// 		},
	// 		"include": ["src/**/*"]
	// 	}`,
	// 	configFileName: "tsconfig.json",
	// 	basePath:       "/apath",
	// 	allFileList:    []string{"/apath/test.ts", "/apath/foge.ts"},
	// },
	// {
	// 	jsonText: `{
	//     "exclude": ["node_modules", "dist"]
	// }`,
	// 	configFileName: "tsconfig.json",
	// 	basePath:       "/apath",
	// 	allFileList:    []string{"/apath/test.ts", "/apath/.git/a.ts", "/apath/.b.ts", "/apath/..c.ts"},
	// },
}

// var host ParseConfigHost = ParseConfigHost{
// 	useCaseSensitiveFileNames: true,
// 	readDirectory: func(rootDir string, extensions []string, excludes []string, includes []string, depth int) []string {
// 		return []string{}
// 	},
// 	fileExists: func(path string) bool {
// 		return true
// 	},
// 	readFile: func(path string) string {
// 		return parseCommandJson[0].jsonText
// 	},
// }

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

func newVFSParseConfigHost(file []string, currentDirectory string) *VfsParseConfigHost {
	fs := fstest.MapFS{}
	for _, f := range file {
		fs[f] = &fstest.MapFile{
			Data: []byte(""),
		}
	}
	return &VfsParseConfigHost{
		vfstest.FromMapFS(fs, true /*useCaseSensitiveFileNames*/),
	}
}
