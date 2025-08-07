package execute_test

import (
    "testing"

    "github.com/microsoft/typescript-go/internal/execute"
)

// This test intentionally reproduces a panic in source-map emission during declaration emit.
// It mirrors the blackbox invocation that crashes when declaration maps are enabled.
// The test is expected to FAIL until the underlying bug is fixed.
func TestDeclarationMap_DestructuredParam_BlackboxPanic(t *testing.T) {
    t.Parallel()

    files := FileMap{
        "/home/src/workspaces/project/tsconfig.json": `{
            "compilerOptions": {
                "declaration": true,
                "declarationMap": true,
                "emitDeclarationOnly": true,
                "strict": true
            }
        }`,
        "/home/src/workspaces/project/index.ts": `export const fn = ({ a, b }: { a: string; b: number }) => { console.log(a, b); };`,
    }

    sys := newTestSys(files, "")
    // Intentionally do not recover: we want the panic to crash the test to reproduce the bug.
    _ = execute.CommandLine(sys, []string{}, true)
}


