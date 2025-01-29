package execute_test

import (
	"fmt"
	"io"
	"strings"
	"testing/fstest"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/diagnosticwriter"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
)

type FileMap map[string]string

func NewTestSys(fileOrFolderList FileMap, cwd string, args ...string) *testSys {
	if cwd == "" {
		cwd = "/home/src/workspaces/project"
	}
	// todo: rest of TestServerHost constructor
	mapFS := fstest.MapFS{}
	fileList := []string{}
	for name, content := range fileOrFolderList {
		mapFS[strings.TrimPrefix(name, "/")] = &fstest.MapFile{
			Data: []byte(content),
		}
		fileList = append(fileList, name)
	}
	fs := bundled.WrapFS(vfstest.FromMapFS(mapFS, true /*useCaseSensitiveFileNames*/))
	newHost := compiler.NewCompilerHost(&core.CompilerOptions{}, cwd, fs)
	return &testSys{
		host:         newHost,
		files:        fileList,
		output:       []string{},
		currentWrite: &strings.Builder{},
		formatOpts: &diagnosticwriter.FormattingOptions{
			NewLine: newHost.NewLine(),
			ComparePathsOptions: tspath.ComparePathsOptions{
				CurrentDirectory:          newHost.GetCurrentDirectory(),
				UseCaseSensitiveFileNames: newHost.FS().UseCaseSensitiveFileNames(),
			},
		},
	}
}

type testSys struct {
	// todo: original has write to output as a string[] because the separations are needed for baselining
	output         []string
	currentWrite   *strings.Builder
	serializedDiff map[string]string
	host           compiler.CompilerHost
	formatOpts     *diagnosticwriter.FormattingOptions
	files          []string
}

func (s *testSys) FS() vfs.FS {
	return s.Host().FS()
}

func (s *testSys) Host() compiler.CompilerHost {
	return s.host
}

func (s *testSys) GetFormatOpts() *diagnosticwriter.FormattingOptions {
	return s.formatOpts
}

func (s *testSys) Writer() io.Writer {
	// todo: check accuracy with original
	return s.currentWrite
}

func (s *testSys) EndWrite() {
	// todo: revisit if improving tsc/build/watch unittest baselines
	s.output = append(s.output, s.currentWrite.String())
	s.currentWrite.Reset()
}

func (s *testSys) serializeState(baseline io.Writer, order serializeOutputOrder) {
	if order == serializeOutputOrderBefore {
		s.serializeOutput(baseline)
	}
	s.diff(baseline)
	if order == serializeOutputOrderAfter {
		s.serializeOutput(baseline)
	}
	// todo watch
	// this.serializeWatches(baseline);
	// this.timeoutCallbacks.serialize(baseline);
	// this.immediateCallbacks.serialize(baseline);
	// this.pendingInstalls.serialize(baseline);
	// this.service?.baseline();
}

func (s *testSys) baselineFS(baseline io.Writer) {
	fmt.Fprint(baseline, "\n\nCurrentFiles::")
	err := s.FS().WalkDir(s.Host().GetCurrentDirectory(), func(path string, d vfs.DirEntry, e error) error {
		if d == nil {
			return nil
		}
		if !d.IsDir() {
			contents, ok := s.FS().ReadFile(path)
			if !ok {
				return e
			}
			fmt.Fprint(baseline, "\n//// ["+path+"]\n"+contents+"\n")
		}
		return nil
	})
	if err != nil {
		panic("walkdir error during fs baseline")
	}
}

func (s *testSys) serializeOutput(baseline io.Writer) {
	fmt.Fprintln(baseline, "\nOutput::")
	// todo screen clears
	s.baselineOutputs(baseline, 0, len(s.output))
}

func (s *testSys) diff(baseline io.Writer) {
	// todo: watch isnt implemented
	// todo: doesn't actually do anything rn, but don't really care atm because we aren't passing edits into the test, so we don't care abt diffs
	snap := map[string]string{}

	err := s.FS().WalkDir(s.Host().GetCurrentDirectory(), func(path string, d vfs.DirEntry, e error) error {
		if d == nil {
			return nil
		}

		newContents, ok := s.FS().ReadFile(path)
		if !ok {
			return e
		}
		snap[path] = newContents
		diffFSEntry(baseline, s.serializedDiff[path], newContents, path)

		return nil
	})
	if err != nil {
		panic("walkdir error during diff")
	}
	for path, oldDirContents := range s.serializedDiff {
		if s.FS().FileExists(path) {
			_, ok := s.FS().ReadFile(path)
			if !ok {
				// report deleted
				diffFSEntry(baseline, oldDirContents, "", path)
			}
		}
	}
	s.serializedDiff = snap
	fmt.Fprint(baseline, s.host.NewLine())
}

func diffFSEntry(baseline io.Writer, oldDirContent string, newDirContent string, path string) {
	// todo handle more cases of fs changes
	if newDirContent == "" {
		fmt.Fprint(baseline, `//// [`, path, `] deleted`, "\n")
	} else if newDirContent == oldDirContent {
		return
	} else {
		fmt.Fprint(baseline, `//// [`, path, `]\n`, newDirContent, "\n")
	}
}

func (s *testSys) baselineOutputs(baseline io.Writer, start int, end int) {
	// todo sanitize sys output
	fmt.Fprint(baseline, strings.Join(s.output[start:end], "\n"))
}

type serializeOutputOrder int

const (
	serializeOutputOrderNone   serializeOutputOrder = iota
	serializeOutputOrderBefore serializeOutputOrder = 1
	serializeOutputOrderAfter  serializeOutputOrder = 2
)
