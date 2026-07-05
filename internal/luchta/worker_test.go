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
	// Ensure no "inputs" key in done
	if strings.Contains(s, `"inputs"`) {
		t.Fatalf("done must not have inputs key: %s", s)
	}
	if !fileExists(filepath.Join(cwd, "dist", "index.js")) {
		t.Fatalf("expected emit")
	}
}

func TestServeResolveTaskModifies(t *testing.T) {
	cwd := t.TempDir()
	writeTsPackage(t, cwd, `{"include":["src/**"],"compilerOptions":{"outDir":"dist"}}`,
		"index.ts", "export const x = 1;\n")

	// JSON-encode cwd safely into the ResolveTask message.
	in := strings.NewReader(`{"type":"resolveTask","id":"r1","name":"build","command":"","package":"p","cwd":` + jsonString(cwd) + `,"scripts":[],"mode":"run"}` + "\n")
	var out, errw bytes.Buffer
	if err := Serve(context.Background(), in, &out, &errw); err != nil {
		t.Fatalf("Serve: %v", err)
	}
	outStr := out.String()
	// Should emit a modify decision with inputs
	if !strings.Contains(outStr, `"decision":"modify"`) {
		t.Fatalf("expected modify decision: %s", outStr)
	}
	if !strings.Contains(outStr, `"inputs"`) {
		t.Fatalf("expected inputs in modify: %s", outStr)
	}
	// Should contain tsconfig.json and src/**
	if !strings.Contains(outStr, "tsconfig.json") {
		t.Fatalf("expected tsconfig.json in inputs: %s", outStr)
	}
	if !strings.Contains(outStr, "src/**") {
		t.Fatalf("expected src/** in inputs: %s", outStr)
	}
}

func TestServeResolveTaskMergesDeclaredInputs(t *testing.T) {
	cwd := t.TempDir()
	writeTsPackage(t, cwd, `{"include":["src/**"],"compilerOptions":{"outDir":"dist"}}`,
		"index.ts", "export const x = 1;\n")

	// ResolveTask JSON with declared inputs (custom/**) that should be merged
	in := strings.NewReader(`{"type":"resolveTask","id":"r1","name":"build","command":"","package":"p","cwd":` + jsonString(cwd) + `,"scripts":[],"mode":"run","inputs":["custom/**","extra.txt"]}` + "\n")
	var out, errw bytes.Buffer
	if err := Serve(context.Background(), in, &out, &errw); err != nil {
		t.Fatalf("Serve: %v", err)
	}
	outStr := out.String()

	// Should emit a modify decision with merged inputs
	if !strings.Contains(outStr, `"decision":"modify"`) {
		t.Fatalf("expected modify decision: %s", outStr)
	}

	// Should contain declared inputs
	if !strings.Contains(outStr, "custom/**") {
		t.Fatalf("expected custom/** (declared) in inputs: %s", outStr)
	}
	if !strings.Contains(outStr, "extra.txt") {
		t.Fatalf("expected extra.txt (declared) in inputs: %s", outStr)
	}

	// Should also contain worker-derived inputs
	if !strings.Contains(outStr, "tsconfig.json") {
		t.Fatalf("expected tsconfig.json (worker-derived) in inputs: %s", outStr)
	}
	if !strings.Contains(outStr, "src/**") {
		t.Fatalf("expected src/** (worker-derived) in inputs: %s", outStr)
	}
}

func TestMergeInputs(t *testing.T) {
	tests := []struct {
		name        string
		declared    []string
		worker      []string
		want        []string
	}{
		{
			name:     "both empty",
			declared: nil,
			worker:   nil,
			want:     nil,
		},
		{
			name:     "declared only",
			declared: []string{"a.txt", "b.txt"},
			worker:   nil,
			want:     []string{"a.txt", "b.txt"},
		},
		{
			name:     "worker only",
			declared: nil,
			worker:   []string{"x/**", "y/**"},
			want:     []string{"x/**", "y/**"},
		},
		{
			name:     "both with overlap",
			declared: []string{"a.txt", "shared/**"},
			worker:   []string{"shared/**", "b.txt"},
			want:     []string{"a.txt", "b.txt", "shared/**"},
		},
		{
			name:     "full overlap",
			declared: []string{"a", "b", "c"},
			worker:   []string{"a", "b", "c"},
			want:     []string{"a", "b", "c"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mergeInputs(tt.declared, tt.worker)
			if len(got) != len(tt.want) {
				t.Fatalf("mergeInputs(%v, %v) = %v, want %v", tt.declared, tt.worker, got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Fatalf("mergeInputs(%v, %v) = %v, want %v", tt.declared, tt.worker, got, tt.want)
				}
			}
		})
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
