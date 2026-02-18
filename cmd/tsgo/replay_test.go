package main

import (
	"bufio"
	"flag"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/json"
	"github.com/microsoft/typescript-go/internal/jsonrpc"
	"github.com/microsoft/typescript-go/internal/ls/lsconv"
	"github.com/microsoft/typescript-go/internal/lsp"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil/lsptestutil"
	"github.com/microsoft/typescript-go/internal/vfs/osvfs"
)

var replay = flag.String("replay", "", "Path to replay file")
var testDir = flag.String("testDir", "", "Path to project directory")

type initialArguments struct {
	RootDirUriPlaceholder string `json:"rootDirUriPlaceholder"`
	RootDirPlaceholder    string `json:"rootDirPlaceholder"`
}

func TestReplay(t *testing.T) {
	if replay == nil || *replay == "" {
		t.Skip("no replay file specified")
	}
	if testDir == nil || *testDir == "" {
		t.Fatal("testDir must be specified")
	}
	testDirUri := lsconv.FileNameToDocumentURI(*testDir)

	fs := bundled.WrapFS(osvfs.FS())
	defaultLibraryPath := bundled.LibPath()
	typingsLocation := getGlobalTypingsCacheLocation()
	serverOpts := lsp.ServerOptions{
		Err:                os.Stderr,
		Cwd:                core.Must(os.Getwd()),
		FS:                 fs,
		DefaultLibraryPath: defaultLibraryPath,
		TypingsLocation:    typingsLocation,
		NpmInstall: func(cwd string, args []string) ([]byte, error) {
			cmd := exec.Command("npm", args...)
			cmd.Dir = cwd
			return cmd.Output()
		},
	}

	client, closeClient := lsptestutil.NewLSPClient(t, serverOpts, nil)
	defer closeClient()

	f, err := os.Open(*replay)
	if err != nil {
		t.Fatalf("failed to read replay file: %v", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	if !scanner.Scan() {
		t.Fatalf("replay file is empty")
	}

	rootDirPlaceholder := "@PROJECT_ROOT@"
	rootDirUriPlaceholder := "@PROJECT_ROOT_URI@"
	firstLine := scanner.Bytes()
	var initObj initialArguments
	err = json.Unmarshal(firstLine, &initObj)
	if err != nil {
		t.Fatalf("failed to parse initial arguments: %v", err)
	}

	if initObj.RootDirPlaceholder != "" {
		rootDirPlaceholder = initObj.RootDirPlaceholder
	}
	if initObj.RootDirUriPlaceholder != "" {
		rootDirUriPlaceholder = initObj.RootDirUriPlaceholder
	}

	rootDirReplacer := strings.NewReplacer(
		rootDirPlaceholder, *testDir,
		rootDirUriPlaceholder, string(testDirUri),
	)

	var id int32 = 1
	for scanner.Scan() {
		line := scanner.Text()
		line = rootDirReplacer.Replace(line)
		var rawMsg struct {
			Kind   string     `json:"kind"`
			Method string     `json:"method"`
			Params json.Value `json:"params"`
		}
		err := json.Unmarshal([]byte(line), &rawMsg)
		if err != nil {
			t.Fatalf("failed to parse message: %v", err)
		}

		var kind jsonrpc.MessageKind
		var reqID *jsonrpc.ID
		switch rawMsg.Kind {
		case "request":
			kind = jsonrpc.MessageKindRequest
			reqID = lsproto.NewID(lsproto.IntegerOrString{Integer: &id})
		case "notification":
			kind = jsonrpc.MessageKindNotification
		default:
			t.Fatalf("unknown message kind: %s", rawMsg.Kind)
		}

		var rpcMsg struct {
			JSONRPC string      `json:"jsonrpc"`
			ID      *jsonrpc.ID `json:"id"`
			Method  string      `json:"method"`
			Params  json.Value  `json:"params"`
		}
		rpcMsg.JSONRPC = "2.0"
		rpcMsg.ID = reqID
		rpcMsg.Method = rawMsg.Method
		rpcMsg.Params = rawMsg.Params
		rpcData, err := json.Marshal(rpcMsg)
		if err != nil {
			t.Fatalf("failed to marshal rpc message: %v", err)
		}

		var msg lsproto.Message
		err = json.Unmarshal(rpcData, &msg)
		if err != nil {
			t.Fatalf("failed to unmarshal rpc message into lsproto.Message: %v", err)
		}

		switch kind {
		case jsonrpc.MessageKindRequest:
			response, ok := client.SendRequestWorker(t, msg.AsRequest(), reqID)
			if !ok {
				t.Fatalf("failed to send request for method %s", rawMsg.Method)
			}
			if response.Error != nil {
				t.Fatalf("server returned error for method %s: %v", rawMsg.Method, response.Error)
			}
		case jsonrpc.MessageKindNotification:
			client.WriteMsg(t, &msg)
		default:
			t.Fatalf("unknown message kind: %s", rawMsg.Kind)
		}
	}
}
