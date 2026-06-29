package lsproto

import (
	"errors"
	"testing"

	"github.com/microsoft/typescript-go/internal/json"
)

// TestUnmarshalParamsRequiresParams verifies that a NoParams method must be
// given no params while every other method must be given params, and that a
// mismatch (including a null value either way) is an InvalidParams error.
func TestUnmarshalParamsRequiresParams(t *testing.T) {
	t.Parallel()

	wantInvalidParams := func(t *testing.T, err error) {
		t.Helper()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, ErrorCodeInvalidParams) {
			t.Fatalf("expected InvalidParams, got %v", err)
		}
	}
	wantNoError := func(t *testing.T, err error) {
		t.Helper()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	// NoParams: only truly-absent/empty params are accepted; null and any
	// present value are rejected.
	t.Run("NoParams/absent", func(t *testing.T) {
		t.Parallel()
		_, err := UnmarshalParams[NoParams](&RequestMessage{Params: nil})
		wantNoError(t, err)
	})
	t.Run("NoParams/empty", func(t *testing.T) {
		t.Parallel()
		_, err := UnmarshalParams[NoParams](&RequestMessage{Params: json.Value(``)})
		wantNoError(t, err)
	})
	t.Run("NoParams/null", func(t *testing.T) {
		t.Parallel()
		_, err := UnmarshalParams[NoParams](&RequestMessage{Params: json.Value(`null`)})
		wantInvalidParams(t, err)
	})
	t.Run("NoParams/object", func(t *testing.T) {
		t.Parallel()
		_, err := UnmarshalParams[NoParams](&RequestMessage{Params: json.Value(`{}`)})
		wantInvalidParams(t, err)
	})

	// Required-params method: only an object or array is accepted; absent,
	// empty, null, and other scalars are rejected.
	for _, tc := range []struct {
		name   string
		params any
	}{
		{"absent", nil},
		{"empty", json.Value(``)},
		{"null", json.Value(`null`)},
		{"number", json.Value(`5`)},
		{"string", json.Value(`"x"`)},
	} {
		t.Run("typed/"+tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := UnmarshalParams[*InitializeParams](&RequestMessage{Params: tc.params})
			wantInvalidParams(t, err)
		})
	}

	// A real params object decodes.
	t.Run("typed/present decodes", func(t *testing.T) {
		t.Parallel()
		got, err := UnmarshalParams[*DidChangeConfigurationParams](&RequestMessage{
			Params: json.Value(`{"settings":{"x":1}}`),
		})
		wantNoError(t, err)
		if got == nil || got.Settings == nil {
			t.Fatalf("params not decoded: %#v", got)
		}
	})
}
