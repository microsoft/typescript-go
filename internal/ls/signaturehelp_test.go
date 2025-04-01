package ls

import (
	"fmt"
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/parser"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"

	"github.com/microsoft/typescript-go/internal/scanner"
)

// var fileText = `function fnTest(str: string, num: number) { }
// fnTest(”, 5);`
var fileText = `declare function f(s: string);
declare function f(n: number);
declare function f(s: string, b: boolean);
declare function f(n: number, b: boolean);

f(1)`
var files = map[string]string{
	"/file1.ts": fileText,
}

func TestSignatureHelp(t *testing.T) {
	t.Parallel()

	file := parser.ParseSourceFile("/file.ts", "/file.ts", fileText, core.ScriptTargetLatest, scanner.JSDocParsingModeParseAll)
	position := strings.Index(fileText, "1")

	// creating a program
	fs := vfstest.FromMap(map[string]string{
		"/file.ts": fileText,
		"/tsconfig.json": `
					  {
						  "compilerOptions": {}
					  }
				  `,
	}, false /*useCaseSensitiveFileNames*/)
	fs = bundled.WrapFS(fs)

	cd := "/"
	host := compiler.NewCompilerHost(nil, cd, fs, bundled.LibPath())
	opts := compiler.ProgramOptions{
		Host:           host,
		ConfigFileName: "/tsconfig.json",
	}
	p := compiler.NewProgram(opts)
	service := NewLanguageService(newLanguageServiceHost(files, p))
	result := service.GetSignatureHelpItems(file.FileName(), position, nil)
	fmt.Sprintf("result: %v\n", result)
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

// type projectServiceHost struct {
// 	fs                 vfs.FS
// 	mu                 sync.Mutex
// 	defaultLibraryPath string
// 	output             strings.Builder
// }

// func newProjectServiceHost(files map[string]string) *projectServiceHost {
// 	fs := bundled.WrapFS(vfstest.FromMap(files, false /*useCaseSensitiveFileNames*/))
// 	host := &projectServiceHost{
// 		fs:                 fs,
// 		defaultLibraryPath: bundled.LibPath(),
// 	}
// 	return host
// }

// // DefaultLibraryPath implements project.ProjectServiceHost.
// func (p *projectServiceHost) DefaultLibraryPath() string {
// 	return p.defaultLibraryPath
// }

// // FS implements project.ProjectServiceHost.
// func (p *projectServiceHost) FS() vfs.FS {
// 	return p.fs
// }

// // GetCurrentDirectory implements project.ProjectServiceHost.
// func (p *projectServiceHost) GetCurrentDirectory() string {
// 	return "/"
// }

// // Log implements project.ProjectServiceHost.
// func (p *projectServiceHost) Log(msg ...any) {
// 	p.mu.Lock()
// 	defer p.mu.Unlock()
// 	fmt.Fprintln(&p.output, msg...)
// }

// // NewLine implements project.ProjectServiceHost.
// func (p *projectServiceHost) NewLine() string {
// 	return "\n"
// }

// func (p *projectServiceHost) replaceFS(files map[string]string) {
// 	p.fs = bundled.WrapFS(vfstest.FromMap(files, false /*useCaseSensitiveFileNames*/))
// }
