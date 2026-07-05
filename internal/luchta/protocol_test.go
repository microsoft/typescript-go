package luchta

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestDecodeRun(t *testing.T) {
	msg, err := DecodeMessage([]byte(`{"type":"run","id":"pkg#task","command":"","cwd":"packages/pkg","env":{"A":"b"}}`))
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	run, ok := msg.(*Run)
	if !ok {
		t.Fatalf("want *Run, got %T", msg)
	}
	if run.ID != "pkg#task" || run.Cwd != "packages/pkg" || run.Env["A"] != "b" {
		t.Fatalf("bad run: %+v", run)
	}
}

func TestDecodeResolveTask(t *testing.T) {
	msg, err := DecodeMessage([]byte(`{"type":"resolveTask","id":"j","name":"build","command":"","package":"@r/a","cwd":"packages/a","scripts":["build"],"mode":"run"}`))
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if _, ok := msg.(*ResolveTask); !ok {
		t.Fatalf("want *ResolveTask, got %T", msg)
	}
}

func TestWriterEmitsTaggedCamelCase(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)
	w.Log("id1", "stdout", "hello")
	w.Done("id1", 0, []string{"dist/a.js"})
	w.ResolvedModify("id2", []string{"src/**", "tsconfig.json"})
	out := buf.String()
	for _, want := range []string{
		`{"type":"log","id":"id1","stream":"stdout","line":"hello"}`,
		`"type":"done"`, `"exitCode":0`, `"outputs":["dist/a.js"]`,
		`{"type":"resolved","id":"id2","result":{"decision":"modify","inputs":["src/**","tsconfig.json"]}}`,
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("output missing %q:\n%s", want, out)
		}
	}
	// each message on its own line
	if lines := strings.Count(strings.TrimSpace(out), "\n"); lines != 2 {
		t.Fatalf("want 3 lines (2 newlines), got %d:\n%s", lines, out)
	}
}

func TestDoneHasNoInputs(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)
	w.Done("id1", 0, []string{"dist/a.js"})
	out := buf.String()
	// Ensure no "inputs" key is present
	if strings.Contains(out, `"inputs"`) {
		t.Fatalf("Done must not contain inputs key:\n%s", out)
	}
}

func TestResolvedModifyHasInputs(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)
	w.ResolvedModify("id1", []string{"src/**", "tsconfig.json"})
	out := buf.String()
	if !strings.Contains(out, `"decision":"modify"`) {
		t.Fatalf("missing modify decision:\n%s", out)
	}
	if !strings.Contains(out, `"inputs":["src/**","tsconfig.json"]`) {
		t.Fatalf("missing inputs:\n%s", out)
	}
}

func TestResolvedAcceptNoInputs(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)
	w.ResolvedAccept("id1")
	out := buf.String()
	if !strings.Contains(out, `"decision":"accept"`) {
		t.Fatalf("missing accept decision:\n%s", out)
	}
	// Accept should not have inputs key (omitempty)
	if strings.Contains(out, `"inputs"`) {
		t.Fatalf("accept must not contain inputs key:\n%s", out)
	}
}

func jsonString(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}
