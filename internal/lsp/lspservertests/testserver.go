package lspservertests

import (
	"fmt"
	"maps"
	"slices"
	"strings"
	"testing"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/ls"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/project"
	"github.com/microsoft/typescript-go/internal/testutil/fsbaselineutil"
	"github.com/microsoft/typescript-go/internal/testutil/lsptestutil"
	"github.com/microsoft/typescript-go/internal/testutil/projecttestutil"
	"github.com/microsoft/typescript-go/internal/vfs/iovfs"
	"gotest.tools/v3/assert"
)

type projectInfo struct {
	program *compiler.Program
}

type testServer struct {
	t            *testing.T
	files        map[string]any
	server       *lsptestutil.TestLspServer
	utils        *projecttestutil.SessionUtils
	baseline     strings.Builder
	fsDiffer     *fsbaselineutil.FSDiffer
	writtenFiles collections.SyncSet[string]

	serializedProjects map[string]*projectInfo
	openFiles          map[string]string
}

func newTestServer(t *testing.T, files map[string]any) *testServer {
	t.Helper()
	server, utils := lsptestutil.Setup(t, files)
	testServer := &testServer{
		t:                  t,
		files:              files,
		server:             server,
		utils:              utils,
		serializedProjects: make(map[string]*projectInfo),
		openFiles:          make(map[string]string),
	}
	wrappedFs := testServer.server.FS.(bundled.WrappedFS)
	testServer.fsDiffer = &fsbaselineutil.FSDiffer{
		FS:           wrappedFs.InternalFS().(iovfs.FsWithSys),
		WrittenFiles: &testServer.writtenFiles,
	}
	fmt.Fprintf(&testServer.baseline, "UseCaseSensitiveFileNames: %v\n", wrappedFs.UseCaseSensitiveFileNames())
	testServer.fsDiffer.BaselineFSwithDiff(&testServer.baseline)
	return testServer
}

func (s *testServer) content(fileName string) string {
	return s.files[fileName].(string)
}

func (s *testServer) hoverToWriteProjectStatus(fileName string) {
	// Do hover so we have snapshot to check things on!!
	_, _, resultOk := lsptestutil.SendRequest(s.t, s.server, lsproto.TextDocumentHoverInfo, &lsproto.HoverParams{
		TextDocument: lsproto.TextDocumentIdentifier{
			Uri: lsproto.DocumentUri("file://" + fileName),
		},
		Position: lsproto.Position{
			Line:      uint32(0),
			Character: uint32(0),
		},
	})
	assert.Assert(s.t, resultOk)
}

func (s *testServer) baselineProjectsAfterNotification(fileName string) {
	s.t.Helper()
	s.hoverToWriteProjectStatus(fileName)
	s.baselineState(false)
}

func (s *testServer) baselineState(before bool) {
	s.t.Helper()

	serialized := s.serializedState()
	if serialized != "" {
		s.baseline.WriteString(serialized)
	}
}

func (s *testServer) serializedState() string {
	var builder strings.Builder
	s.fsDiffer.BaselineFSwithDiff(&builder)
	if strings.TrimSpace(builder.String()) == "" {
		builder.Reset()
	}

	session := s.server.Session()
	if session == nil {
		return builder.String()
	}

	snapshot, release := session.Snapshot()
	defer release()

	hasChange := false
	currentProjects := make(map[string]*projectInfo)
	diffs := make(map[string]string)
	subDiffs := make(map[string][]string)

	for _, project := range snapshot.ProjectCollection.ConfiguredProjects() {
		hasChange = s.addProjectDiff(currentProjects, diffs, subDiffs, project) || hasChange
	}
	if inferredProject := snapshot.ProjectCollection.InferredProject(); inferredProject != nil {
		hasChange = s.addProjectDiff(currentProjects, diffs, subDiffs, inferredProject) || hasChange
	}
	for projectName, info := range s.serializedProjects {
		if currentProjects[projectName] == nil {
			diffs[projectName] = "*deleted*"
			var subDiff []string
			if info.program != nil {
				for _, file := range info.program.GetSourceFiles() {
					if fileName := file.FileName(); !s.isLibFile(fileName) {
						subDiff = append(subDiff, fileName)
					}
				}
			}
			subDiffs[projectName] = subDiff
			hasChange = true
		}
	}
	s.serializedProjects = currentProjects

	if hasChange {
		diffKeys := slices.Collect(maps.Keys(diffs))
		slices.Sort(diffKeys)
		builder.WriteString("Projects::\n")
		for _, projectName := range diffKeys {
			diff := diffs[projectName]
			builder.WriteString(fmt.Sprintf("  [%s] %s\n", projectName, diff))
			for _, subDiff := range subDiffs[projectName] {
				builder.WriteString(fmt.Sprintf("    %s\n", subDiff))
			}
		}
	}
	return builder.String()
}

func (s *testServer) isLibFile(fileName string) bool {
	return strings.HasPrefix(fileName, bundled.LibPath()+"/")
}

func (s *testServer) addProjectDiff(currentProjects map[string]*projectInfo, diffs map[string]string, subDiffs map[string][]string, project *project.Project) bool {
	s.t.Helper()
	program := project.GetProgram()
	var oldProgram *compiler.Program
	entry := &projectInfo{program}
	currentProjects[project.Name()] = entry
	changed := false
	diff := ""
	if existing, ok := s.serializedProjects[project.Name()]; ok {
		oldProgram = existing.program
		if oldProgram != program {
			diff = "*modified*"
			changed = true
		} else {
			diff = ""
		}
	} else {
		diff = "*new*"
		changed = true
	}

	diffs[project.Name()] = diff
	var subDiff []string
	if program != nil {
		for _, file := range program.GetSourceFiles() {
			fileDiff := ""
			// No need to write "*new*" for files as its obvious
			if diff == "*modified*" {
				if oldProgram == nil {
					fileDiff = " *new*"
				} else if oldFile := oldProgram.GetSourceFileByPath(file.Path()); oldFile == nil {
					fileDiff = " *new*"
				} else if oldFile != file {
					fileDiff = " *modified*"
				}
			}
			fileName := file.FileName()
			if fileDiff != "" || !s.isLibFile(fileName) {
				subDiff = append(subDiff, fileName+fileDiff)
			}
		}
	}
	if oldProgram != program && oldProgram != nil {
		for _, file := range oldProgram.GetSourceFiles() {
			if program == nil || program.GetSourceFileByPath(file.Path()) == nil {
				subDiff = append(subDiff, file.FileName()+" *deleted*")
			}
		}
	}
	subDiffs[project.Name()] = subDiff
	return changed
}

type requestOrMessage struct {
	Method lsproto.Method `json:"method"`
	Params any            `json:"params,omitzero"`
}

func baselineRequestOrNotification(t *testing.T, s *testServer, method lsproto.Method, params any) {
	s.t.Helper()
	s.baselineState(true)
	res, _ := json.Marshal(requestOrMessage{
		Method: method,
		Params: params,
	}, jsontext.WithIndent("  "))
	s.baseline.WriteString(fmt.Sprintln(string(res)))
}

func sendNotification[Params any](t *testing.T, s *testServer, info lsproto.NotificationInfo[Params], params Params) {
	s.t.Helper()
	baselineRequestOrNotification(t, s, info.Method, params)
	lsptestutil.SendNotification(s.t, s.server, info, params)
}

func sendRequest[Params, Resp any](t *testing.T, s *testServer, info lsproto.RequestInfo[Params, Resp], params Params) Resp {
	s.t.Helper()
	baselineRequestOrNotification(t, s, info.Method, params)
	resMsg, result, resultOk := lsptestutil.SendRequest(t, s.server, info, params)
	s.baselineState(false)
	if resMsg == nil {
		s.t.Fatalf("Nil response received for %s", info.Method)
	}
	if !resultOk {
		s.t.Fatalf("Unexpected response type for %s: %T", info.Method, resMsg.AsResponse().Result)
	}
	return result
}

func (s *testServer) openFile(fileName string, languageID lsproto.LanguageKind) {
	s.t.Helper()
	s.openFileWithContent(fileName, s.content(fileName), languageID)
}

func (s *testServer) openFileWithContent(fileName string, content string, languageID lsproto.LanguageKind) {
	s.t.Helper()
	s.openFiles[fileName] = content
	sendNotification(s.t, s, lsproto.TextDocumentDidOpenInfo, &lsproto.DidOpenTextDocumentParams{
		TextDocument: &lsproto.TextDocumentItem{
			Uri:        lsproto.DocumentUri("file://" + fileName),
			LanguageId: languageID,
			Text:       content,
		},
	})
	s.baselineProjectsAfterNotification(fileName)
}

func (s *testServer) closeFile(fileName string) {
	s.t.Helper()
	delete(s.openFiles, fileName)
	sendNotification(s.t, s, lsproto.TextDocumentDidCloseInfo, &lsproto.DidCloseTextDocumentParams{
		TextDocument: lsproto.TextDocumentIdentifier{
			Uri: lsproto.DocumentUri("file://" + fileName),
		},
	})
	// Skip baselining projects here since updated snapshot is not generated right away after this
}

func (s *testServer) changeFile(params *lsproto.DidChangeTextDocumentParams) {
	s.t.Helper()
	fileName := params.TextDocument.Uri.FileName()
	text := s.openFiles[fileName]
	converters := ls.NewConverters(lsproto.PositionEncodingKindUTF8, func(fileName string) *ls.LSPLineMap {
		return ls.ComputeLSPLineStarts(text)
	})
	// Update the contents in openFiles
	for _, textChange := range params.ContentChanges {
		if partialChange := textChange.Partial; partialChange != nil {
			text = converters.FromLSPTextChange(lsptestutil.NewLsScript(fileName, text), partialChange).ApplyTo(text)
		} else if wholeChange := textChange.WholeDocument; wholeChange != nil {
			text = wholeChange.Text
		}
	}
	s.openFiles[fileName] = text
	sendNotification(s.t, s, lsproto.TextDocumentDidChangeInfo, params)
	// Skip baselining projects here since updated snapshot is not generated right away after this
}

func (s *testServer) baselineReferences(fileName string, position lsproto.Position) {
	s.t.Helper()
	result := sendRequest(s.t, s, lsproto.TextDocumentReferencesInfo, &lsproto.ReferenceParams{
		TextDocument: lsproto.TextDocumentIdentifier{
			Uri: lsproto.DocumentUri("file://" + fileName),
		},
		Position: position,
		Context:  &lsproto.ReferenceContext{},
	})
	s.baseline.WriteString(lsptestutil.GetBaselineForLocationsWithFileContents(s.server.FS, *result.Locations, lsptestutil.BaselineLocationsOptions{
		Marker:     &marker{fileName, position},
		MarkerName: "/*FIND ALL REFS*/",
		OpenFiles:  s.openFiles,
	}) + "\n")
}

func (s *testServer) baselineRename(fileName string, position lsproto.Position) {
	s.t.Helper()
	result := sendRequest(s.t, s, lsproto.TextDocumentRenameInfo, &lsproto.RenameParams{
		TextDocument: lsproto.TextDocumentIdentifier{
			Uri: lsproto.DocumentUri("file://" + fileName),
		},
		Position: position,
		NewName:  "?",
	})
	s.baseline.WriteString(lsptestutil.GetBaselineForRename(s.server.FS, result, lsptestutil.BaselineLocationsOptions{
		Marker:    &marker{fileName, position},
		OpenFiles: s.openFiles,
	}) + "\n")
}

type marker struct {
	fileName string
	position lsproto.Position
}

var _ lsptestutil.LocationMarker = (*marker)(nil)

func (m *marker) FileName() string {
	return m.fileName
}

func (m *marker) LSPos() lsproto.Position {
	return m.position
}
