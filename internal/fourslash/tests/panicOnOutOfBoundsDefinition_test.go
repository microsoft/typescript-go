package fourslash_test

import (
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
	"gotest.tools/v3/assert"
)

// TestPanicOnOutOfBoundsDefinition reproduces the panic that occurs when
// the client sends a textDocument/definition request with a line number
// that's beyond the file's line count. This can happen due to
// synchronization issues between client and server.
//
// BUG: The server should handle this gracefully by returning an empty result
// or a proper error, but instead it panics with "bad line number".
// This test currently fails because the server panics instead of handling
// the out-of-bounds position gracefully.
func TestPanicOnOutOfBoundsDefinition(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `export {};
interface Point {
	x: number;
	y: number;
}
declare const p: Point;
p.x;
`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)

	// Send a definition request with a line number that's out of bounds.
	// The file has 8 lines (0-7), but we're requesting line 65.
	// This simulates what happens when the client's view of the file is stale.
	msg, errorCode, errorMsg := f.SendDefinitionRequestAtPosition(t, 65, 24)

	// The server should handle this gracefully, not panic.
	// We expect either:
	// 1. A successful response with an empty result, OR
	// 2. A proper error response (not due to panic)
	//
	// Currently, the server panics and returns InternalError (-32603)
	// which is the bug we're testing for.
	if errorCode == -32603 && strings.Contains(errorMsg, "panic") {
		t.Fatalf("BUG: Server panicked when handling out-of-bounds position.\n"+
			"Error code: %d\n"+
			"Error message: %s\n"+
			"Expected: Server should handle out-of-bounds positions gracefully without panicking.",
			errorCode, errorMsg)
	}

	// If we get here without panic, verify the response is reasonable
	assert.Assert(t, msg != nil, "Expected valid response, got nil")
}
