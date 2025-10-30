package lspservertests

import (
	"fmt"
	"strings"
	"testing"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/ls/lsconv"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/project"
	"github.com/microsoft/typescript-go/internal/testutil/fsbaselineutil"
	"github.com/microsoft/typescript-go/internal/testutil/lsptestutil"
	"github.com/microsoft/typescript-go/internal/testutil/projecttestutil"
	"gotest.tools/v3/assert"
)

type testServer struct {
	t             *testing.T
	files         map[string]any
	server        *lsptestutil.TestLspServer
	utils         *projecttestutil.SessionUtils
	baseline      strings.Builder
	fsDiffer      *fsbaselineutil.FSDiffer
	writtenFiles  collections.SyncSet[string]
	isInitialized bool

	serializedProjects           map[string]projectInfo
	serializedOpenFiles          map[string]*openFileInfo
	serializedConfigFileRegistry *project.ConfigFileRegistry
	openFiles                    map[string]string
}

func newTestServer(t *testing.T, files map[string]any) *testServer {
	t.Helper()
	server, utils := lsptestutil.Setup(t, files)
	testServer := &testServer{
		t:         t,
		files:     files,
		server:    server,
		utils:     utils,
		openFiles: make(map[string]string),
	}
	testServer.fsDiffer = &fsbaselineutil.FSDiffer{
		FS:           utils.FsFromFileMap(),
		WrittenFiles: &testServer.writtenFiles,
	}
	fmt.Fprintf(&testServer.baseline, "UseCaseSensitiveFileNames: %v\n", utils.FsFromFileMap().UseCaseSensitiveFileNames())
	testServer.fsDiffer.BaselineFSwithDiff(&testServer.baseline)
	return testServer
}

func (s *testServer) content(fileName string) string {
	if text, ok := s.openFiles[fileName]; ok {
		return text
	}
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

	printStateDiff(s, &builder)
	return builder.String()
}

func (s *testServer) isLibFile(fileName string) bool {
	return strings.HasPrefix(fileName, bundled.LibPath()+"/")
}

type requestOrMessage struct {
	Method lsproto.Method `json:"method"`
	Params any            `json:"params,omitzero"`
}

func baselineRequestOrNotification(t *testing.T, server *testServer, method lsproto.Method, params any) {
	server.t.Helper()
	server.baselineState(true)
	res, _ := json.Marshal(requestOrMessage{
		Method: method,
		Params: params,
	}, jsontext.WithIndent("  "))
	server.baseline.WriteString(fmt.Sprintln(string(res)))
	server.isInitialized = true
}

func sendNotification[Params any](t *testing.T, server *testServer, info lsproto.NotificationInfo[Params], params Params) {
	server.t.Helper()
	baselineRequestOrNotification(t, server, info.Method, params)
	switch info.Method {
	case lsproto.TextDocumentDidOpenInfo.Method:
		openFileParams := any(params).(*lsproto.DidOpenTextDocumentParams)
		server.openFiles[openFileParams.TextDocument.Uri.FileName()] = openFileParams.TextDocument.Text
	case lsproto.TextDocumentDidCloseInfo.Method:
		closeFileParams := any(params).(*lsproto.DidCloseTextDocumentParams)
		delete(server.openFiles, closeFileParams.TextDocument.Uri.FileName())
	case lsproto.TextDocumentDidChangeInfo.Method:
		changeFileParams := any(params).(*lsproto.DidChangeTextDocumentParams)
		fileName := changeFileParams.TextDocument.Uri.FileName()
		text := server.openFiles[fileName]
		converters := lsconv.NewConverters(lsproto.PositionEncodingKindUTF8, func(fileName string) *lsconv.LSPLineMap {
			return lsconv.ComputeLSPLineStarts(text)
		})
		// Update the contents in openFiles
		for _, textChange := range changeFileParams.ContentChanges {
			if partialChange := textChange.Partial; partialChange != nil {
				text = converters.FromLSPTextChange(lsptestutil.NewLsScript(fileName, text), partialChange).ApplyTo(text)
			} else if wholeChange := textChange.WholeDocument; wholeChange != nil {
				text = wholeChange.Text
			}
		}
		server.openFiles[fileName] = text
	}
	lsptestutil.SendNotification(server.t, server.server, info, params)
}

func sendRequest[Params, Resp any](t *testing.T, server *testServer, info lsproto.RequestInfo[Params, Resp], params Params) Resp {
	server.t.Helper()
	baselineRequestOrNotification(t, server, info.Method, params)
	resMsg, result, resultOk := lsptestutil.SendRequest(t, server.server, info, params)
	server.baselineState(false)
	if resMsg == nil {
		server.t.Fatalf("Nil response received for %s", info.Method)
	}
	if !resultOk {
		server.t.Fatalf("Unexpected response type for %s: %T", info.Method, resMsg.AsResponse().Result)
	}
	return result
}

func (s *testServer) openFile(fileName string, languageID lsproto.LanguageKind) {
	s.t.Helper()
	s.openFileWithContent(fileName, s.content(fileName), languageID)
}

func (s *testServer) openFileWithContent(fileName string, content string, languageID lsproto.LanguageKind) {
	s.t.Helper()
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
	sendNotification(s.t, s, lsproto.TextDocumentDidCloseInfo, &lsproto.DidCloseTextDocumentParams{
		TextDocument: lsproto.TextDocumentIdentifier{
			Uri: lsproto.DocumentUri("file://" + fileName),
		},
	})
	// Skip baselining projects here since updated snapshot is not generated right away after this
}

func (s *testServer) changeFile(params *lsproto.DidChangeTextDocumentParams) {
	s.t.Helper()
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

func (s *testServer) baselineWorkspaceSymbol(query string) {
	s.t.Helper()
	result := sendRequest(s.t, s, lsproto.WorkspaceSymbolInfo, &lsproto.WorkspaceSymbolParams{
		Query: query,
	})
	s.baseline.WriteString(lsptestutil.GetBaselineForWorkspaceSymbol(s.server.FS, result, lsptestutil.BaselineLocationsOptions{
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
