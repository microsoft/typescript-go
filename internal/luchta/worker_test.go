package luchta

import (
	"bytes"
	"context"
	"path/filepath"
	"strings"
	"testing"
)

func TestServeRunEmitsLogAndDone(t *testing.T) {
	cwd := t.TempDir()
	writeTsPackage(t, cwd, `{"compilerOptions":{"outDir":"dist","rootDir":"src","module":"nodenext","moduleResolution":"nodenext"}}`,
		"index.ts", "export const x = 1;\n")

	// JSON-encode cwd safely into the Run message.
	in := strings.NewReader(`{"type":"run","id":"t1","command":"","cwd":` + jsonString(cwd) + `,"env":{}}` + "\n")
	var out, errw bytes.Buffer
	if err := Serve(context.Background(), in, &out, &errw); err != nil {
		t.Fatalf("Serve: %v", err)
	}
	s := out.String()
	if !strings.Contains(s, `"type":"done"`) || !strings.Contains(s, `"id":"t1"`) || !strings.Contains(s, `"exitCode":0`) {
		t.Fatalf("missing done: %s", s)
	}
	if !fileExists(filepath.Join(cwd, "dist", "index.js")) {
		t.Fatalf("expected emit")
	}
}

func TestServeResolveTaskAccepts(t *testing.T) {
	in := strings.NewReader(`{"type":"resolveTask","id":"r1","name":"build","command":"","package":"p","cwd":"x","scripts":[],"mode":"run"}` + "\n")
	var out, errw bytes.Buffer
	if err := Serve(context.Background(), in, &out, &errw); err != nil {
		t.Fatalf("Serve: %v", err)
	}
	if !strings.Contains(out.String(), `"decision":"accept"`) {
		t.Fatalf("expected accept: %s", out.String())
	}
}

func TestServeMalformedLineGoesToStderr(t *testing.T) {
	in := strings.NewReader("not json\n")
	var out, errw bytes.Buffer
	if err := Serve(context.Background(), in, &out, &errw); err != nil {
		t.Fatalf("Serve: %v", err)
	}
	if out.Len() != 0 {
		t.Fatalf("malformed input must not write to protocol stdout: %s", out.String())
	}
	if errw.Len() == 0 {
		t.Fatalf("expected stderr diagnostic")
	}
}
