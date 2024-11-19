package parsingoptions

import (
	// "runtime"

	"testing"

	// "encoding/json"
	json2 "github.com/go-json-experiment/json"
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
	// 	`,
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
}

func TestBaselineParseResult(t *testing.T) {
	//var baseline []string = []string{}

	for _, jsonText := range jsonTexts {
		//baseline = append(baseline, "Input::", jsonText)
		parsed := ParseConfigFileTextToJson("/apath/tsconfig.json", jsonText)
		s, ok := (*parsed.config).(string)
		if ok {
			json2.Unmarshal([]byte(s), &parsed)
		}
		//fmt.Println(s)
	}
}
