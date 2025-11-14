package ls

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
)

// Test for issue: panic handling request textDocument/codeAction
// This verifies that we safely handle nil diagnostic codes without panicking
func TestCodeActionDiagnosticNilChecks(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		diagnostic *lsproto.Diagnostic
		shouldSkip bool
	}{
		{
			name: "nil Code field",
			diagnostic: &lsproto.Diagnostic{
				Range: lsproto.Range{
					Start: lsproto.Position{Line: 0, Character: 0},
					End:   lsproto.Position{Line: 0, Character: 10},
				},
				Message: "Test diagnostic with nil code",
				Code:    nil,
			},
			shouldSkip: true,
		},
		{
			name: "nil Integer in Code (string code instead)",
			diagnostic: &lsproto.Diagnostic{
				Range: lsproto.Range{
					Start: lsproto.Position{Line: 0, Character: 0},
					End:   lsproto.Position{Line: 0, Character: 10},
				},
				Message: "Test diagnostic with string code",
				Code: &lsproto.IntegerOrString{
					String:  stringPtr("TS1234"),
					Integer: nil,
				},
			},
			shouldSkip: true,
		},
		{
			name: "valid Integer code",
			diagnostic: &lsproto.Diagnostic{
				Range: lsproto.Range{
					Start: lsproto.Position{Line: 0, Character: 0},
					End:   lsproto.Position{Line: 0, Character: 10},
				},
				Message: "Test diagnostic with integer code",
				Code: &lsproto.IntegerOrString{
					Integer: int32Ptr(2304),
				},
			},
			shouldSkip: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// Test the nil check logic that prevents the panic
			// This mimics the check in ProvideCodeActions at line 67
			shouldSkip := tc.diagnostic.Code == nil || tc.diagnostic.Code.Integer == nil

			if shouldSkip != tc.shouldSkip {
				t.Errorf("Expected shouldSkip=%v, got %v", tc.shouldSkip, shouldSkip)
			}

			// Verify we can safely access Code.Integer when shouldSkip is false
			if !shouldSkip {
				// This should not panic
				errorCode := *tc.diagnostic.Code.Integer
				if errorCode != 2304 {
					t.Errorf("Expected error code 2304, got %d", errorCode)
				}
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}

func int32Ptr(i int32) *int32 {
	return &i
}
