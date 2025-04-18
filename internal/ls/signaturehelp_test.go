package ls

import (
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/parser"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
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
	triggerReason             SignatureHelpTriggerReason
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
		fs := vfstest.FromMap(map[string]string{
			testData.files[0].filename: testData.files[0].content,
			"/tsconfig.json": `
						  {
							  "compilerOptions": {}
						  }
					  `,
		}, false /*useCaseSensitiveFileNames*/)
		fs = bundled.WrapFS(fs)
		host := compiler.NewCompilerHost(nil, "/", fs, bundled.LibPath())
		opts := compiler.ProgramOptions{
			Host:           host,
			ConfigFileName: "/tsconfig.json",
		}
		p := compiler.NewProgram(opts)
		files := map[string]string{
			testData.files[0].filename: testData.files[0].content,
		}
		service := NewLanguageService(newLanguageServiceHost(files, p))
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

// setting up language service
type languageServiceHost struct {
	fs                 vfs.FS
	defaultLibraryPath string
	program            *compiler.Program
}

func newLanguageServiceHost(files map[string]string, p *compiler.Program) *languageServiceHost {
	fs := bundled.WrapFS(vfstest.FromMap(files, false /*useCaseSensitiveFileNames*/))
	host := &languageServiceHost{
		fs:                 fs,
		defaultLibraryPath: bundled.LibPath(),
		program:            p,
	}
	return host
}

func (l *languageServiceHost) DefaultLibraryPath() string {
	return l.defaultLibraryPath
}

func (l *languageServiceHost) FS() vfs.FS {
	return l.fs
}

func (l *languageServiceHost) GetProgram() *compiler.Program {
	return l.program
}

func (l *languageServiceHost) GetCurrentDirectory() string {
	return "/"
}

func (l *languageServiceHost) NewLine() string {
	return "\n"
}

func (l *languageServiceHost) GetProjectVersion() int {
	return 1
}

func (l *languageServiceHost) GetCompilerOptions() *core.CompilerOptions {
	return l.program.GetCompilerOptions()
}

func (l *languageServiceHost) GetRootFileNames() []string {
	return []string{"/file1.ts"}
}

func (l *languageServiceHost) GetSourceFile(fileName string, path tspath.Path, languageVersion core.ScriptTarget) *ast.SourceFile {
	return l.program.GetSourceFile(fileName)
}

func (l *languageServiceHost) Trace(msg string) {
	fmt.Println("Trace:", msg)
}

func (l *languageServiceHost) GetDefaultLibraryPath() string {
	return l.defaultLibraryPath
}

// setting up project service

type projectServiceHost struct {
	fs                 vfs.FS
	mu                 sync.Mutex
	defaultLibraryPath string
	output             strings.Builder
}

func newProjectServiceHost(files map[string]string) *projectServiceHost {
	fs := bundled.WrapFS(vfstest.FromMap(files, false /*useCaseSensitiveFileNames*/))
	host := &projectServiceHost{
		fs:                 fs,
		defaultLibraryPath: bundled.LibPath(),
	}
	return host
}

// DefaultLibraryPath implements project.ProjectServiceHost.
func (p *projectServiceHost) DefaultLibraryPath() string {
	return p.defaultLibraryPath
}

// FS implements project.ProjectServiceHost.
func (p *projectServiceHost) FS() vfs.FS {
	return p.fs
}

// GetCurrentDirectory implements project.ProjectServiceHost.
func (p *projectServiceHost) GetCurrentDirectory() string {
	return "/"
}

// Log implements project.ProjectServiceHost.
func (p *projectServiceHost) Log(msg ...any) {
	p.mu.Lock()
	defer p.mu.Unlock()
	fmt.Fprintln(&p.output, msg...)
}

// NewLine implements project.ProjectServiceHost.
func (p *projectServiceHost) NewLine() string {
	return "\n"
}

func (p *projectServiceHost) replaceFS(files map[string]string) {
	p.fs = bundled.WrapFS(vfstest.FromMap(files, false /*useCaseSensitiveFileNames*/))
}
