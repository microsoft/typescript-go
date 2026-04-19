package main_test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

func TestWasmExports(t *testing.T) {
	ctx := context.Background()

	wasmPath := filepath.Join(t.TempDir(), "test.wasm")
	cmd := exec.Command("go", "build", "-buildmode=c-shared", "-o", wasmPath, ".")
	cmd.Dir = "."
	cmd.Env = append(os.Environ(), "GOOS=wasip1", "GOARCH=wasm")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build wasm: %v\nOutput: %s", err, string(out))
	}

	wasmBytes, err := os.ReadFile(wasmPath)
	if err != nil {
		t.Fatalf("Failed to read wasm: %v", err)
	}

	r := wazero.NewRuntimeWithConfig(ctx, wazero.NewRuntimeConfig())
	wasiConfig := wazero.NewModuleConfig().WithStdout(os.Stdout).WithStderr(os.Stderr)
	defer r.Close(ctx)

	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	mod, err := r.InstantiateWithConfig(ctx, wasmBytes, wasiConfig)
	if err != nil {
		t.Fatalf("Failed to instantiate wasm: %v", err)
	}

	
	initFunc := mod.ExportedFunction("_initialize")
	if initFunc != nil {
		_, err := initFunc.Call(ctx)
		if err != nil {
			t.Fatalf("Failed to initialize: %v", err)
		}
	}
	wasmMalloc := mod.ExportedFunction("wasm_malloc")
	parseSource := mod.ExportedFunction("ParseSource")
	getNodeKind := mod.ExportedFunction("GetNodeKind")
	getNodeChildren := mod.ExportedFunction("GetNodeChildren")
	getNodeText := mod.ExportedFunction("GetNodeText")

	if wasmMalloc == nil || parseSource == nil || getNodeKind == nil || getNodeChildren == nil {
		t.Fatal("Missing required exports")
	}

	filename := "/test.ts"
	sourceText := "const x = 1;"

	allocStr := func(s string) uint32 {
		res, err := wasmMalloc.Call(ctx, uint64(len(s)))
		if err != nil {
			t.Fatalf("wasm_malloc failed: %v", err)
		}
		ptr := uint32(res[0])
		if !mod.Memory().Write(ptr, []byte(s)) {
			t.Fatalf("Failed to write to wasm memory")
		}
		return ptr
	}

	filenamePtr := allocStr(filename)
	sourcePtr := allocStr(sourceText)

	res, err := parseSource.Call(ctx, uint64(filenamePtr), uint64(len(filename)), uint64(sourcePtr), uint64(len(sourceText)))
	if err != nil {
		t.Fatalf("ParseSource failed: %v", err)
	}

	rootId := uint32(res[0])
	if rootId == 0 {
		t.Fatal("Expected non-zero rootId")
	}

	res, err = getNodeKind.Call(ctx, uint64(rootId))
	if err != nil {
		t.Fatalf("GetNodeKind failed: %v", err)
	}
	kind := uint32(res[0])
	if kind == 0 {
		t.Fatal("Expected non-zero kind for SourceFile")
	}

	// allocate pointer for children out param
	outArrayPtrRes, err := wasmMalloc.Call(ctx, 4)
	outArrayPtr := uint32(outArrayPtrRes[0])

	res, err = getNodeChildren.Call(ctx, uint64(rootId), uint64(outArrayPtr))
	if err != nil {
		t.Fatalf("GetNodeChildren failed: %v", err)
	}
	numChildren := int32(res[0])
	
	if numChildren <= 0 {
		t.Fatalf("Expected children, got %d", numChildren)
	}

	// Read outArrayPtr to find where the array is
	arrayPtrBytes, ok := mod.Memory().Read(outArrayPtr, 4)
	if !ok { t.Fatalf("Could not read outArrayPtr") }
	
	arrayPtr := uint32(arrayPtrBytes[0]) | (uint32(arrayPtrBytes[1]) << 8) | (uint32(arrayPtrBytes[2]) << 16) | (uint32(arrayPtrBytes[3]) << 24)

	// Read first child
	childIdBytes, ok := mod.Memory().Read(arrayPtr, 4)
	if !ok { t.Fatalf("Could not read child array") }
	
	childId := uint32(childIdBytes[0]) | (uint32(childIdBytes[1]) << 8) | (uint32(childIdBytes[2]) << 16) | (uint32(childIdBytes[3]) << 24)

	res, err = getNodeKind.Call(ctx, uint64(childId))
	childKind := uint32(res[0])
	
	if childKind == 0 {
		t.Fatal("Expected non-zero child kind")
	}

	outTextPtrRes, _ := wasmMalloc.Call(ctx, 4)
	outTextPtr := uint32(outTextPtrRes[0])

	res, err = getNodeText.Call(ctx, uint64(childId), uint64(outTextPtr))
	textLen := int32(res[0])
	if textLen <= 0 {
		t.Fatalf("Expected text length, got %d", textLen)
	}

	textPtrBytes, _ := mod.Memory().Read(outTextPtr, 4)
	textPtr := uint32(textPtrBytes[0]) | (uint32(textPtrBytes[1]) << 8) | (uint32(textPtrBytes[2]) << 16) | (uint32(textPtrBytes[3]) << 24)
	
	textBytes, _ := mod.Memory().Read(textPtr, uint32(textLen))
	if len(textBytes) == 0 {
		t.Fatal("Failed to read text bytes")
	}

	t.Logf("Parsed SourceFile ID: %d, Kind: %d", rootId, kind)
	t.Logf("Child 1 ID: %d, Kind: %d, Text: %s", childId, childKind, string(textBytes))
}
