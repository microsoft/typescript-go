package tsoptions

import (
	// "runtime"

	"encoding/json"
	"fmt"
	"testing"
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
}

func TestGetParsedCommandJson(t *testing.T) {
	for _, test := range parseCommandJson {
		parsed := ParseConfigFileTextToJson(test.configFileName, test.jsonText)
		parseConfigFileContent := ParseJsonConfigFileContent(
			parsed.config.(map[string]interface{}),
			host,
			test.basePath,
			//basePath ?? ts.getNormalizedAbsolutePath(ts.getDirectoryPath(configFileName), host.sys.getCurrentDirectory()),
			nil,
			&test.configFileName,
			/*resolutionStack*/ nil,
			/*extraFileExtensions*/ nil,
			/*extendedConfigCache*/ nil,
		)
		configJson, err := json.Marshal(parseConfigFileContent)
		if err != nil {
			t.Errorf("Failed to marshal parseConfigFileContent: %v", err)
		}
		fmt.Println("****************************************************")
		fmt.Println(string(configJson))
	}
}

var parseCommandJson = []verifyConfig{
	{
		jsonText: `{
		"compilerOptions": {
			"lib": ["es5"]
		},
	}`,
		configFileName: "tsconfig.json",
		basePath:       "/apath",
		allFileList:    []string{"/apath/test.ts", "/apath/foge.ts"},
	},
	{
		jsonText:       `{}`,
		configFileName: "tsconfig.json",
		basePath:       "/apath",
		allFileList:    []string{"/apath/test.ts", "/apath/.git/a.ts", "/apath/.b.ts", "/apath/..c.ts"},
	},
}

var host ParseConfigHost = ParseConfigHost{
	useCaseSensitiveFileNames: true,
	readDirectory: func(rootDir string, extensions []string, excludes []string, includes []string, depth int) []string {
		return []string{}
	},
	fileExists: func(path string) bool {
		return true
	},
	readFile: func(path string) string {
		return parseCommandJson[0].jsonText
	},
}

// type ParseConfigHost struct {
// 	useCaseSensitiveFileNames bool
// 	readDirectory             func(rootDir string, extensions []string, excludes []string, includes []string, depth int) []string
// 	/**
// 	 * Gets a value indicating whether the specified path exists and is a file.
// 	 * @param path The path to test.
// 	 */
// 	fileExists func(path string) bool
// 	readFile   func(path string) string
// 	trace      func(s string)
// }
