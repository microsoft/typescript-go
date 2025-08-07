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
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/diagnosticwriter"
	"github.com/microsoft/typescript-go/internal/execute"
	"github.com/microsoft/typescript-go/internal/incremental"
	"github.com/microsoft/typescript-go/internal/testutil/incrementaltestutil"
	"github.com/microsoft/typescript-go/internal/testutil/stringtestutil"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
)

type FileMap map[string]any

var tscLibPath = "/home/src/tslibs/TS/Lib"

var tscDefaultLibContent = stringtestutil.Dedent(`
/// <reference no-default-lib="true"/>
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
declare const console: { log(msg: any): void; };
`)

func newTestSys(fileOrFolderList FileMap, cwd string) *testSys {
	if cwd == "" {
		cwd = "/home/src/workspaces/project"
	}
	fs, timeImpl := vfstest.FromMapWithTime(fileOrFolderList, true /*useCaseSensitiveFileNames*/)
	sys := &testSys{
		fs: &incrementaltestutil.FsHandlingBuildInfo{
			FS: &testFs{
				FS: fs,
			},
		},
		defaultLibraryPath: tscLibPath,
		cwd:                cwd,
		files:              slices.Collect(maps.Keys(fileOrFolderList)),
		currentWrite:       &strings.Builder{},
		timeImpl:           timeImpl,
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
	content   string
	mTime     time.Time
	isWritten bool
}

type snapshot struct {
	snap        map[string]*diffEntry
	defaultLibs *collections.SyncSet[string]
}

type testSys struct {
	// todo: original has write to output as a string[] because the separations are needed for baselining
	currentWrite   *strings.Builder
	serializedDiff *snapshot

	fs                 *incrementaltestutil.FsHandlingBuildInfo
	defaultLibraryPath string
	cwd                string
	files              []string

	timeImpl *vfstest.Time
}

var (
	_ execute.System             = (*testSys)(nil)
	_ execute.CommandLineTesting = (*testSys)(nil)
)

func (s *testSys) Now() time.Time {
	return s.timeImpl.Now()
}

func (s *testSys) SinceStart() time.Duration {
	return s.timeImpl.SinceStart()
}

func (s *testSys) FS() vfs.FS {
	return s.fs
}

func (s *testSys) testFs() *testFs {
	return s.fs.FS.(*testFs)
}

func (s *testSys) fsFromFileMap() vfs.FS {
	return s.testFs().FS
}

func (s *testSys) ensureLibPathExists(path string) {
	path = tscLibPath + "/" + path
	if _, ok := s.fsFromFileMap().ReadFile(path); !ok {
		if s.testFs().defaultLibs == nil {
			s.testFs().defaultLibs = &collections.SyncSet[string]{}
		}
		s.testFs().defaultLibs.Add(path)
		err := s.fsFromFileMap().WriteFile(path, tscDefaultLibContent, false)
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

func (s *testSys) Writer() io.Writer {
	return s.currentWrite
}

func (s *testSys) OnEmittedFiles(result *compiler.EmitResult) {
	if result != nil {
		for _, file := range result.EmittedFiles {
			// Ensure that the timestamp for emitted files is in the order
			now := s.timeImpl.Now()
			s.fsFromFileMap().Chtimes(file, time.Time{}, now)
		}
	}
}

func (s *testSys) baselinePrograms(baseline *strings.Builder, programs []*incremental.Program, watcher *execute.Watcher) {
	if watcher != nil {
		programs = []*incremental.Program{watcher.GetProgram()}
	}
	for index, program := range programs {
		if index > 0 {
			baseline.WriteString("\n")
		}
		s.baselineProgram(baseline, program)
	}
}

func (s *testSys) baselineProgram(baseline *strings.Builder, program *incremental.Program) {
	if program == nil {
		return
	}

	testingData := program.GetTestingData(program.GetProgram())
	if testingData.ConfigFilePath != "" {
		baseline.WriteString(testingData.ConfigFilePath + "::\n")
	}
	baseline.WriteString("SemanticDiagnostics::\n")
	for _, file := range program.GetProgram().GetSourceFiles() {
		if diagnostics, ok := testingData.SemanticDiagnosticsPerFile.Load(file.Path()); ok {
			if oldDiagnostics, ok := testingData.OldProgramSemanticDiagnosticsPerFile.Load(file.Path()); !ok || oldDiagnostics != diagnostics {
				baseline.WriteString("*refresh*    " + file.FileName() + "\n")
			}
		} else {
			baseline.WriteString("*not cached* " + file.FileName() + "\n")
		}
	}

	// Write signature updates
	baseline.WriteString("Signatures::\n")
	for _, file := range program.GetProgram().GetSourceFiles() {
		if kind, ok := testingData.UpdatedSignatureKinds[file.Path()]; ok {
			switch kind {
			case incremental.SignatureUpdateKindComputedDts:
				baseline.WriteString("(computed .d.ts) " + file.FileName() + "\n")
			case incremental.SignatureUpdateKindStoredAtEmit:
				baseline.WriteString("(stored at emit) " + file.FileName() + "\n")
			case incremental.SignatureUpdateKindUsedVersion:
				baseline.WriteString("(used version)   " + file.FileName() + "\n")
			}
		}
	}
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
	output := s.getOutput(false)
	fmt.Fprint(baseline, output)
}

var (
	fakeTimeStamp = "HH:MM:SS AM"
	fakeDuration  = "d.ddds"

	buildStartingAt            = "build starting at "
	buildFinishedIn            = "build finished in "
	prettyStatusTimeStampStart = "[" + diagnosticwriter.ForegroundColorEscapeGrey
	listFileStart              = "TSFILE:  "
)

func (s *testSys) getOutput(forComparing bool) string {
	lines := strings.Split(s.currentWrite.String(), "\n")
	outputLines := make([]string, 0, len(lines))
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if change := strings.Replace(line, "Version "+core.Version(), "Version "+incrementaltestutil.FakeTsVersion, 1); change != line {
			outputLines = append(outputLines, change)
			continue
		}
		if strings.HasPrefix(line, buildStartingAt) {
			if !forComparing {
				outputLines = append(outputLines, buildStartingAt+fakeTimeStamp)
			}
			continue
		}
		if strings.HasPrefix(line, buildFinishedIn) {
			if !forComparing {
				outputLines = append(outputLines, buildFinishedIn+fakeDuration)
			}
			continue
		}
		if change := strings.TrimPrefix(line, prettyStatusTimeStampStart); change != line {
			if !forComparing {
				outputLines = append(outputLines, prettyStatusTimeStampStart+fakeTimeStamp+change[len(fakeTimeStamp):])
			}
		} else if strings.Index(line, " -") == len(fakeTimeStamp) && strings.Index(line, ":") == 2 && strings.Index(line[3:], ":") == 2 {
			// Fuzzy check for hh:mm:ss AM/PM - string
			outputLines = append(outputLines, fakeTimeStamp+line[11:])
		} else if !forComparing || !strings.HasPrefix(line, listFileStart) {
			outputLines = append(outputLines, line)
			continue
		}

		// Consume all the buildStatus reported
		for j := i + 1; j < len(lines); j++ {
			if !forComparing {
				outputLines = append(outputLines, lines[j])
			}
			if lines[j] == "" {
				i = j
				break
			}
		}
	}
	return strings.Join(outputLines, "\n")
}

func (s *testSys) clearOutput() {
	s.currentWrite.Reset()
}

func (s *testSys) baselineFSwithDiff(baseline io.Writer) {
	// todo: baselines the entire fs, possibly doesn't correctly diff all cases of emitted files, since emit isn't fully implemented and doesn't always emit the same way as strada
	snap := map[string]*diffEntry{}

	testFs := s.testFs()
	diffs := map[string]string{}
	err := s.fsFromFileMap().WalkDir("/", func(path string, d vfs.DirEntry, e error) error {
		if e != nil {
			return e
		}

		if !d.Type().IsRegular() {
			return nil
		}

		newContents, ok := s.fsFromFileMap().ReadFile(path)
		if !ok {
			return nil
		}
		stat := s.fsFromFileMap().Stat(path)
		if stat == nil {
			panic("stat is nil: " + path)
		}
		newEntry := &diffEntry{content: newContents, mTime: stat.ModTime(), isWritten: testFs.writtenFiles.Has(path)}
		snap[path] = newEntry
		s.addFsEntryDiff(diffs, newEntry, path)

		return nil
	})
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		panic("walkdir error during diff: " + err.Error())
	}
	if s.serializedDiff != nil {
		for path := range s.serializedDiff.snap {
			_, ok := s.fsFromFileMap().ReadFile(path)
			if !ok {
				// report deleted
				s.addFsEntryDiff(diffs, nil, path)
			}
		}
	}
	var defaultLibs collections.SyncSet[string]
	if s.testFs().defaultLibs != nil {
		s.testFs().defaultLibs.Range(func(libPath string) bool {
			defaultLibs.Add(libPath)
			return true
		})
	}
	s.serializedDiff = &snapshot{
		snap:        snap,
		defaultLibs: &defaultLibs,
	}
	diffKeys := slices.Collect(maps.Keys(diffs))
	slices.Sort(diffKeys)
	for _, path := range diffKeys {
		fmt.Fprint(baseline, "//// ["+path+"] ", diffs[path], "\n")
	}
	fmt.Fprintln(baseline)
	testFs.writtenFiles = collections.SyncSet[string]{} // Reset written files after baseline
}

func (s *testSys) addFsEntryDiff(diffs map[string]string, newDirContent *diffEntry, path string) {
	var oldDirContent *diffEntry
	var defaultLibs *collections.SyncSet[string]
	if s.serializedDiff != nil {
		oldDirContent = s.serializedDiff.snap[path]
		defaultLibs = s.serializedDiff.defaultLibs
	}
	// todo handle more cases of fs changes
	if oldDirContent == nil {
		if s.testFs().defaultLibs == nil || !s.testFs().defaultLibs.Has(path) {
			diffs[path] = "*new* \n" + newDirContent.content
		}
	} else if newDirContent == nil {
		diffs[path] = "*deleted*"
	} else if newDirContent.content != oldDirContent.content {
		diffs[path] = "*modified* \n" + newDirContent.content
	} else if newDirContent.isWritten {
		diffs[path] = "*rewrite with same content*"
	} else if newDirContent.mTime != oldDirContent.mTime {
		diffs[path] = "*mTime changed*"
	} else if defaultLibs != nil && defaultLibs.Has(path) && s.testFs().defaultLibs != nil && !s.testFs().defaultLibs.Has(path) {
		// Lib file that was read
		diffs[path] = "*Lib*\n" + newDirContent.content
	}
}

func (s *testSys) writeFileNoError(path string, content string, writeByteOrderMark bool) {
	if err := s.fsFromFileMap().WriteFile(path, content, writeByteOrderMark); err != nil {
		panic(err)
	}
}

func (s *testSys) replaceFileText(path string, oldText string, newText string) {
	content, ok := s.fsFromFileMap().ReadFile(path)
	if !ok {
		panic("File not found: " + path)
	}
	content = strings.Replace(content, oldText, newText, 1)
	s.writeFileNoError(path, content, false)
}

func (s *testSys) appendFile(path string, text string) {
	content, ok := s.fsFromFileMap().ReadFile(path)
	if !ok {
		panic("File not found: " + path)
	}
	s.writeFileNoError(path, content+text, false)
}

func (s *testSys) prependFile(path string, text string) {
	content, ok := s.fsFromFileMap().ReadFile(path)
	if !ok {
		panic("File not found: " + path)
	}
	s.writeFileNoError(path, text+content, false)
}
