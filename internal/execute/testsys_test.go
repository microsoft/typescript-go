package execute_test

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"maps"
	"slices"
	"strings"
	"time"

	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/testutil/incrementaltestutil"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
)

type FileMap map[string]any

var (
	tscLibPath           = "/home/src/tslibs/TS/Lib"
	tscDefaultLibContent = `/// <reference no-default-lib="true"/>
interface Boolean {}
interface Function {}
interface CallableFunction {}
interface NewableFunction {}
interface IArguments {}
interface Number { toExponential: any; }
interface Object {}
interface RegExp {}
interface String { charAt: any; }
interface Array<T> { length: number; [n: number]: T; }
interface ReadonlyArray<T> {}
interface SymbolConstructor {
    (desc?: string | number): symbol;
    for(name: string): symbol;
    readonly toStringTag: symbol;
}
declare var Symbol: SymbolConstructor;
interface Symbol {
    readonly [Symbol.toStringTag]: string;
}
declare const console: { log(msg: any): void; };`
)

func newTestSys(fileOrFolderList FileMap, cwd string) *testSys {
	if cwd == "" {
		cwd = "/home/src/workspaces/project"
	}
	sys := &testSys{
		fs:                 NewFSTrackingLibs(incrementaltestutil.NewFsHandlingBuildInfo(vfstest.FromMap(fileOrFolderList, true /*useCaseSensitiveFileNames*/))),
		defaultLibraryPath: tscLibPath,
		cwd:                cwd,
		files:              slices.Collect(maps.Keys(fileOrFolderList)),
		output:             []string{},
		currentWrite:       &strings.Builder{},
		start:              time.Now(),
	}

	// Ensure the default library file is present
	sys.ensureLibPathExists("lib.d.ts")
	for _, libFile := range tsoptions.TargetToLibMap() {
		sys.ensureLibPathExists(libFile)
	}
	for libFile := range tsoptions.LibFilesSet.Keys() {
		sys.ensureLibPathExists(libFile)
	}
	return sys
}

type diffEntry struct {
	content  string
	fileInfo fs.FileInfo
}

type snapshot struct {
	snap        map[string]*diffEntry
	defaultLibs *collections.Set[string]
}

type testSys struct {
	// todo: original has write to output as a string[] because the separations are needed for baselining
	output         []string
	currentWrite   *strings.Builder
	serializedDiff *snapshot

	fs                 *testFsTrackingLibs
	defaultLibraryPath string
	cwd                string
	files              []string

	start time.Time
}

func (s *testSys) Now() time.Time {
	// todo: make a "test time" structure
	return time.Now()
}

func (s *testSys) SinceStart() time.Duration {
	return time.Since(s.start)
}

func (s *testSys) FS() vfs.FS {
	return s.fs
}

func (s *testSys) TestFS() *incrementaltestutil.FsHandlingBuildInfo {
	return s.fs.fs.(*incrementaltestutil.FsHandlingBuildInfo)
}

func (s *testSys) ensureLibPathExists(path string) {
	path = tscLibPath + "/" + path
	if _, ok := s.TestFS().ReadFile(path); !ok {
		if s.fs.defaultLibs == nil {
			s.fs.defaultLibs = collections.NewSetWithSizeHint[string](tsoptions.LibFilesSet.Len() + len(tsoptions.TargetToLibMap()) + 1)
		}
		s.fs.defaultLibs.Add(path)
		err := s.TestFS().WriteFile(path, tscDefaultLibContent, false)
		if err != nil {
			panic("Failed to write default library file: " + err.Error())
		}
	}
}

func (s *testSys) DefaultLibraryPath() string {
	return s.defaultLibraryPath
}

func (s *testSys) GetCurrentDirectory() string {
	return s.cwd
}

func (s *testSys) NewLine() string {
	return "\n"
}

func (s *testSys) Writer() io.Writer {
	return s.currentWrite
}

func (s *testSys) EndWrite() {
	// todo: revisit if improving tsc/build/watch unittest baselines
	s.output = append(s.output, s.currentWrite.String())
	s.currentWrite.Reset()
}

func (s *testSys) serializeState(baseline *strings.Builder) {
	s.baselineOutput(baseline)
	s.baselineFSwithDiff(baseline)
	// todo watch
	// this.serializeWatches(baseline);
	// this.timeoutCallbacks.serialize(baseline);
	// this.immediateCallbacks.serialize(baseline);
	// this.pendingInstalls.serialize(baseline);
	// this.service?.baseline();
}

func (s *testSys) baselineOutput(baseline io.Writer) {
	fmt.Fprint(baseline, "\nOutput::\n")
	if len(s.output) == 0 {
		fmt.Fprint(baseline, "No output\n")
		return
	}
	// todo screen clears
	s.printOutputs(baseline)
	s.output = []string{}
}

func (s *testSys) baselineFSwithDiff(baseline io.Writer) {
	// todo: baselines the entire fs, possibly doesn't correctly diff all cases of emitted files, since emit isn't fully implemented and doesn't always emit the same way as strada
	snap := map[string]*diffEntry{}

	err := s.FS().WalkDir("/", func(path string, d vfs.DirEntry, e error) error {
		if e != nil {
			return e
		}

		if !d.Type().IsRegular() {
			return nil
		}

		newContents, ok := s.TestFS().FS().ReadFile(path)
		if !ok {
			return nil
		}
		fileInfo, err := d.Info()
		if err != nil {
			return nil
		}
		newEntry := &diffEntry{content: newContents, fileInfo: fileInfo}
		snap[path] = newEntry
		s.reportFSEntryDiff(baseline, newEntry, path)

		return nil
	})
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		panic("walkdir error during diff: " + err.Error())
	}
	if s.serializedDiff != nil {
		for path := range s.serializedDiff.snap {
			if s.FS().FileExists(path) {
				_, ok := s.TestFS().FS().ReadFile(path)
				if !ok {
					// report deleted
					s.reportFSEntryDiff(baseline, nil, path)
				}
			}
		}
	}
	var defaultLibs *collections.Set[string]
	if s.fs.defaultLibs != nil {
		defaultLibs = s.fs.defaultLibs.Clone()
	}
	s.serializedDiff = &snapshot{
		snap:        snap,
		defaultLibs: defaultLibs,
	}
	fmt.Fprintln(baseline)
}

func (s *testSys) reportFSEntryDiff(baseline io.Writer, newDirContent *diffEntry, path string) {
	var oldDirContent *diffEntry
	var defaultLibs *collections.Set[string]
	if s.serializedDiff != nil {
		oldDirContent = s.serializedDiff.snap[path]
		defaultLibs = s.serializedDiff.defaultLibs
	}
	// todo handle more cases of fs changes
	if oldDirContent == nil {
		if s.fs.defaultLibs == nil || !s.fs.defaultLibs.Has(path) {
			fmt.Fprint(baseline, "//// [", path, "] *new* \n", newDirContent.content, "\n")
		}
	} else if newDirContent == nil {
		fmt.Fprint(baseline, "//// [", path, "] *deleted*\n")
	} else if newDirContent.content != oldDirContent.content {
		fmt.Fprint(baseline, "//// [", path, "] *modified* \n", newDirContent.content, "\n")
	} else if newDirContent.fileInfo.ModTime() != oldDirContent.fileInfo.ModTime() {
		fmt.Fprint(baseline, "//// [", path, "] *modified time*\n")
	} else if defaultLibs != nil && defaultLibs.Has(path) && s.fs.defaultLibs != nil && !s.fs.defaultLibs.Has(path) {
		// Lib file that was read
		fmt.Fprint(baseline, "//// [", path, "] *Lib*\n", newDirContent.content, "\n")
	}
}

func (s *testSys) printOutputs(baseline io.Writer) {
	// todo sanitize sys output
	fmt.Fprint(baseline, strings.Join(s.output, "\n"))
}
