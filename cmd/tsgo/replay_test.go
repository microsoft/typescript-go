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
	"github.com/microsoft/typescript-go/internal/lsp"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil/lsptestutil"
	"github.com/microsoft/typescript-go/internal/vfs/osvfs"
)

var replay = flag.String("replay", "", "Path to replay file")

// !!! move & reuse from fourslash; refactor into some general "test server"
// that is like a harness for running the server and sending it messages and receiving responses etc

type initialArguments struct {
	RootDirPlaceholder string `json:"rootDirPlaceholder"`
}

func TestReplay(t *testing.T) {
	if replay == nil || *replay == "" {
		t.Skip("no replay file specified")
	}

	var testDir string // !!!

	fs := bundled.WrapFS(osvfs.FS())
	defaultLibraryPath := bundled.LibPath()
	typingsLocation := getGlobalTypingsCacheLocation()
	serverOpts := lsp.ServerOptions{
		Err:                os.Stderr,
		Cwd:                core.Must(os.Getwd()), // !!!
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
	firstLine := scanner.Bytes()
	var initObj initialArguments
	err = json.Unmarshal(firstLine, initObj)
	if err != nil {
		t.Fatalf("failed to parse initial arguments: %v", err)
	}

	if initObj.RootDirPlaceholder != "" {
		rootDirPlaceholder = initObj.RootDirPlaceholder
	}

	var id int32 = 1
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.ReplaceAll(line, rootDirPlaceholder, core.Must(os.Getwd()))
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

		// !!! this seems redundant, maybe use byte stream for input.
		var msg lsproto.Message
		err = json.Unmarshal(rpcData, &msg)
		if err != nil {
			t.Fatalf("failed to unmarshal rpc message into lsproto.Message: %v", err)
		}

		if kind == jsonrpc.MessageKindRequest {
			response, ok := client.SendRequestWorker()
		}

	}

	// !!! how to wait for everything to finish running on the server?? wait for response??
}
