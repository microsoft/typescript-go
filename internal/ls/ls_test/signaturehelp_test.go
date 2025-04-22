package ls

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls"
	"github.com/microsoft/typescript-go/internal/parser"
	"github.com/microsoft/typescript-go/internal/testutil/projecttestutil"
	"github.com/microsoft/typescript-go/internal/tspath"
	"gotest.tools/v3/assert"

	"github.com/microsoft/typescript-go/internal/scanner"
)

type verifySignatureHelpOptions struct {
	marker                    string
	overloadsCount            int
	docComment                string
	text                      string
	parameterName             string
	parameterSpan             string
	parameterDocComment       string
	parameterCount            int
	isVariadic                bool
	triggerReason             ls.SignatureHelpTriggerReason
	overrideSelectedItemIndex int
	//tags?: ReadonlyArray<JSDocTagInfo>;
}

var data = []struct {
	title  string
	input  string
	output []verifySignatureHelpOptions
}{
	// 	{
	// 		title: "SignatureHelpCallExpressions",
	// 		input: `function fnTest(str: string, num: number) { }
	// fnTest(/*1*/'', /*2*/5);`,
	// 		output: []verifySignatureHelpOptions{
	// 			{
	// 				marker:         "1",
	// 				text:           `fnTest(str: string, num: number): void`,
	// 				parameterCount: 2,
	// 				parameterSpan:  "str: string",
	// 			},
	// 			{
	// 				marker:         "2",
	// 				text:           `fnTest(str: string, num: number): void`,
	// 				parameterCount: 2,
	// 				parameterSpan:  "num: number",
	// 			},
	// 		},
	// 	},
	{
		title: "SignatureHelpCallExpressionTuples",
		input: `function fnTest(str: string, num: number) { }
declare function wrap<A extends any[], R>(fn: (...a: A) => R) : (...a: A) => R;
var fnWrapped = wrap(fnTest);
fnWrapped(/*1*/'', /*2*/5);
function fnTestVariadic (str: string, ...num: number[]) { }
var fnVariadicWrapped = wrap(fnTestVariadic);
fnVariadicWrapped(/*3*/'', /*4*/5);
function fnNoParams () { }
var fnNoParamsWrapped = wrap(fnNoParams);
fnNoParamsWrapped(/*5*/);`,
		output: []verifySignatureHelpOptions{
			{
				marker:         "1",
				text:           `fnWrapped(str: string, num: number): void`,
				parameterCount: 2,
				parameterSpan:  "str: string",
			},
			{
				marker:         "2",
				text:           `fnWrapped(str: string, num: number): void`,
				parameterCount: 2,
				parameterSpan:  "num: number",
			},
			{
				marker:         "3",
				text:           `fnVariadicWrapped(str: string, ...num: number[]): void`,
				parameterCount: 2,
				parameterSpan:  "str: string",
				isVariadic:     true,
			},
			{
				marker:         "4",
				text:           `fnVariadicWrapped(str: string, ...num: number[]): void`,
				parameterCount: 2,
				parameterSpan:  "...num: number[]",
				isVariadic:     true,
			},
			{
				marker:         "5",
				text:           `fnNoParamsWrapped(): void`,
				parameterCount: 0,
			},
		},
	},
}

func TestSignature(t *testing.T) {
	t.Parallel()

	for _, rec := range data {
		testData := parseTestdata("/file1.ts", rec.input, "/file1.ts")
		// Creating a program
		// fs := vfstest.FromMap(map[string]string{
		// 	testData.files[0].filename: testData.files[0].content,
		// 	"/tsconfig.json": `
		// 				  {
		// 					  "compilerOptions": {}
		// 				  }
		// 			  `,
		// }, false /*useCaseSensitiveFileNames*/)
		// fs = bundled.WrapFS(fs)
		// host := compiler.NewCompilerHost(nil, "/", fs, bundled.LibPath())
		// opts := compiler.ProgramOptions{
		// 	Host:           host,
		// 	ConfigFileName: "/tsconfig.json",
		// }
		// p := compiler.NewProgram(opts)
		files := map[string]string{
			testData.files[0].filename: testData.files[0].content,
		}
		service := createLanguageService(testData.files[0].filename, files)
		file := parser.ParseSourceFile(testData.files[0].filename, tspath.Path(testData.files[0].filename), testData.files[0].content, core.ScriptTargetLatest, scanner.JSDocParsingModeParseAll)

		markerNumber := 0
		for i, marker := range testData.markerPositions {
			result := service.GetSignatureHelpItems(file.FileName(), marker.position, nil)
			if result == nil {
				t.Fatal("expected result to be non-nil")
			}
			assert.Equal(t, rec.output[markerNumber].marker, i, "marker")
			assert.Equal(t, rec.output[markerNumber].text, result.Signatures[result.ActiveSignature].Label, "text")
			assert.Equal(t, rec.output[markerNumber].parameterCount, len(*result.Signatures[result.ActiveSignature].Parameters), "parameterCount")
			assert.Equal(t, rec.output[markerNumber].parameterSpan, (*result.Signatures[result.ActiveSignature].Parameters)[result.ActiveParameter].Label, "parameterSpan")
			markerNumber++
		}
	}
}

func createLanguageService(fileName string, files map[string]string) *ls.LanguageService {
	projectService, _ := projecttestutil.Setup(files)
	projectService.OpenFile(fileName, files[fileName], core.ScriptKindTS, "")
	project := projectService.Projects()[0]
	return project.LanguageService()
}
