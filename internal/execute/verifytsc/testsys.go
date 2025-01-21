package verifytsc

import (
	"fmt"
	"io"
	"strings"
	"testing/fstest"

	"github.com/microsoft/typescript-go/internal/bundled"
	ts "github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	dw "github.com/microsoft/typescript-go/internal/diagnosticwriter"
	"github.com/microsoft/typescript-go/internal/execute"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
)

type FileMap map[string]string

func NewTestSys(fileOrFolderList FileMap, args ...string) *testSys {
	// todo: rest of TestServerHost constructor
	mapFS := fstest.MapFS{}
	fileList := []string{}
	for name, content := range fileOrFolderList {
		mapFS[name] = &fstest.MapFile{
			Data: []byte(content),
		}
		fileList = append(fileList, name)
	}
	fs := bundled.WrapFS(vfstest.FromMapFS(mapFS, true /*useCaseSensitiveFileNames*/))
	newHost := ts.NewCompilerHost(&core.CompilerOptions{}, "/home/src/workspaces/project", fs)
	return &testSys{
		host:         newHost,
		files:        fileList,
		output:       []string{},
		currentWrite: &strings.Builder{},
		formatOpts: &dw.FormattingOptions{
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

	exitVal    execute.ExitStatus
	host       ts.CompilerHost
	formatOpts *dw.FormattingOptions
	files      []string
}

func (s *testSys) FS() vfs.FS {
	return s.Host().FS()
}

func (s *testSys) Host() ts.CompilerHost {
	return s.host
}

func (s *testSys) Exit(e execute.ExitStatus) execute.ExitStatus {
	s.exitVal = e
	return s.exitVal
}

func (s *testSys) GetFormatOpts() *dw.FormattingOptions {
	return s.formatOpts
}

func (s *testSys) Write(p []byte) (n int, err error) {
	// todo: check accuracy with original
	return fmt.Fprint(s.currentWrite, string(p))
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

func (s *testSys) serializeOutput(baseline io.Writer) {
	fmt.Fprintln(baseline, "\nOutput::")
	// todo screen clears
	s.baselineOutputs(baseline, 0, len(s.output))
}

func (s *testSys) diff(baseline io.Writer) {
	snap := map[string]string{}

	err := s.FS().WalkDir(s.Host().GetCurrentDirectory(), func(path string, d vfs.DirEntry, e error) error {
		if d == nil || d.IsDir() {
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
