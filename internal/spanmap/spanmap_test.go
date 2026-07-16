package spanmap_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/spanmap"
	"gotest.tools/v3/assert"
)

func TestMapSpanVerbatim(t *testing.T) {
	t.Parallel()

	// Generated [0,10) is a verbatim copy of original [100,110).
	m := spanmap.New([]spanmap.Segment{
		{GenStart: 0, GenEnd: 10, OrigStart: 100, OrigEnd: 110, Kind: spanmap.KindVerbatim},
	})

	got, fidelity := m.MapSpan(core.NewTextRange(3, 7))
	assert.Equal(t, got.Pos(), 103)
	assert.Equal(t, got.End(), 107)
	assert.Equal(t, fidelity, spanmap.FidelityExact)
}

func TestMapSpanAtom(t *testing.T) {
	t.Parallel()

	// Generated [0,3) is a synthesized gap; [3,14) ("MyComponent") is an atom of the original [60,71).
	m := spanmap.New([]spanmap.Segment{
		{GenStart: 3, GenEnd: 14, OrigStart: 60, OrigEnd: 71, Kind: spanmap.KindAtom},
	})

	// A span inside the atom maps to the whole atom span.
	got, fidelity := m.MapSpan(core.NewTextRange(5, 9))
	assert.Equal(t, got.Pos(), 60)
	assert.Equal(t, got.End(), 71)
	assert.Equal(t, fidelity, spanmap.FidelityAtom)
}

func TestMapSpanSynthesizedGap(t *testing.T) {
	t.Parallel()

	// A gap between two verbatim segments is synthesized: it maps to the insertion point (the preceding
	// segment's original end) with no fidelity.
	m := spanmap.New([]spanmap.Segment{
		{GenStart: 0, GenEnd: 10, OrigStart: 100, OrigEnd: 110, Kind: spanmap.KindVerbatim},
		{GenStart: 20, GenEnd: 30, OrigStart: 200, OrigEnd: 210, Kind: spanmap.KindVerbatim},
	})

	got, fidelity := m.MapSpan(core.NewTextRange(12, 15))
	assert.Equal(t, got.Pos(), 110)
	assert.Equal(t, got.End(), 110)
	assert.Equal(t, fidelity, spanmap.FidelityNone)
}

func TestMapSpanEmptyIsSynthesized(t *testing.T) {
	t.Parallel()

	// An empty map describes fully synthesized output: everything maps to the start with no fidelity.
	m := spanmap.New(nil)
	got, fidelity := m.MapSpan(core.NewTextRange(5, 10))
	assert.Equal(t, got.Pos(), 0)
	assert.Equal(t, got.End(), 0)
	assert.Equal(t, fidelity, spanmap.FidelityNone)
}

func TestMapSpanCrossingSegments(t *testing.T) {
	t.Parallel()

	m := spanmap.New([]spanmap.Segment{
		{GenStart: 0, GenEnd: 10, OrigStart: 100, OrigEnd: 110, Kind: spanmap.KindVerbatim},
		{GenStart: 10, GenEnd: 20, OrigStart: 200, OrigEnd: 210, Kind: spanmap.KindVerbatim},
	})

	got, fidelity := m.MapSpan(core.NewTextRange(5, 15))
	assert.Equal(t, got.Pos(), 105)
	assert.Equal(t, got.End(), 205)
	assert.Equal(t, fidelity, spanmap.FidelityApproximate)
}

func TestMapSpanNilIdentity(t *testing.T) {
	t.Parallel()

	var m *spanmap.SpanMap
	got, fidelity := m.MapSpan(core.NewTextRange(3, 7))
	assert.Equal(t, got.Pos(), 3)
	assert.Equal(t, got.End(), 7)
	assert.Equal(t, fidelity, spanmap.FidelityExact)
}

func TestMarshalRoundTrip(t *testing.T) {
	t.Parallel()

	original := spanmap.New([]spanmap.Segment{
		{GenStart: 3, GenEnd: 14, OrigStart: 60, OrigEnd: 71, Kind: spanmap.KindAtom},
		{GenStart: 14, GenEnd: 24, OrigStart: 71, OrigEnd: 81, Kind: spanmap.KindVerbatim},
	})

	data, err := original.Marshal()
	assert.NilError(t, err)
	decoded, err := spanmap.Unmarshal(data)
	assert.NilError(t, err)

	for _, r := range []core.TextRange{core.NewTextRange(1, 2), core.NewTextRange(4, 10), core.NewTextRange(16, 20)} {
		wantRange, wantFidelity := original.MapSpan(r)
		gotRange, gotFidelity := decoded.MapSpan(r)
		assert.Equal(t, gotRange, wantRange)
		assert.Equal(t, gotFidelity, wantFidelity)
	}
}

func TestValidate(t *testing.T) {
	t.Parallel()

	const transformed = "const greeting = 1;\n"
	const original = "<x>const greeting = 1;\n</x>"
	scriptStart := 3 // index of "const" in original

	testCases := []struct {
		name     string
		segs     []spanmap.Segment
		wantKind spanmap.MappingErrorKind
		wantOK   bool
	}{
		{
			name:   "valid verbatim",
			segs:   []spanmap.Segment{{GenStart: 0, GenEnd: core.TextPos(len(transformed)), OrigStart: core.TextPos(scriptStart), OrigEnd: core.TextPos(scriptStart + len(transformed)), Kind: spanmap.KindVerbatim}},
			wantOK: true,
		},
		{
			name:   "empty is valid",
			segs:   nil,
			wantOK: true,
		},
		{
			name:   "gap is allowed",
			segs:   []spanmap.Segment{{GenStart: 3, GenEnd: core.TextPos(len(transformed)), OrigStart: 0, OrigEnd: 0, Kind: spanmap.KindAtom}},
			wantOK: true,
		},
		{
			name: "overlap",
			segs: []spanmap.Segment{
				{GenStart: 0, GenEnd: 10, OrigStart: 0, OrigEnd: 0, Kind: spanmap.KindAtom},
				{GenStart: 5, GenEnd: core.TextPos(len(transformed)), OrigStart: 0, OrigEnd: 0, Kind: spanmap.KindAtom},
			},
			wantKind: spanmap.MappingErrorKindOverlap,
		},
		{
			name:     "original out of bounds",
			segs:     []spanmap.Segment{{GenStart: 0, GenEnd: core.TextPos(len(transformed)), OrigStart: 0, OrigEnd: core.TextPos(len(original) + 10), Kind: spanmap.KindAtom}},
			wantKind: spanmap.MappingErrorKindOutOfBounds,
		},
		{
			name:     "verbatim text mismatch",
			segs:     []spanmap.Segment{{GenStart: 0, GenEnd: core.TextPos(len(transformed)), OrigStart: 0, OrigEnd: core.TextPos(len(transformed)), Kind: spanmap.KindVerbatim}},
			wantKind: spanmap.MappingErrorKindVerbatimMismatch,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			problem := spanmap.New(tc.segs).Validate(transformed, original)
			if tc.wantOK {
				assert.Assert(t, problem == nil, "expected valid, got %+v", problem)
				return
			}
			assert.Assert(t, problem != nil, "expected a problem")
			assert.Equal(t, problem.Kind, tc.wantKind)
		})
	}
}

func TestValidateNilIsValid(t *testing.T) {
	t.Parallel()
	var m *spanmap.SpanMap
	assert.Assert(t, m.Validate("abc", "abc") == nil)
}
