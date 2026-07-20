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

func TestMapPosition(t *testing.T) {
	t.Parallel()

	// Generated [0,10) is a verbatim copy of original [100,110); [10,20) is an atom of original [200,210).
	m := spanmap.New([]spanmap.Segment{
		{GenStart: 0, GenEnd: 10, OrigStart: 100, OrigEnd: 110, Kind: spanmap.KindVerbatim},
		{GenStart: 20, GenEnd: 30, OrigStart: 200, OrigEnd: 210, Kind: spanmap.KindAtom},
	})

	testCases := []struct {
		name     string
		pos      core.TextPos
		want     core.TextPos
		fidelity spanmap.Fidelity
	}{
		{"verbatim interpolates", 3, 103, spanmap.FidelityExact},
		{"atom maps to its start", 25, 200, spanmap.FidelityAtom},
		{"gap maps to insertion point", 15, 110, spanmap.FidelityNone},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, fidelity := m.MapPosition(tc.pos)
			assert.Equal(t, got, tc.want)
			assert.Equal(t, fidelity, tc.fidelity)
			// MapPosition must agree with MapSpan on a zero-length range.
			span, spanFidelity := m.MapSpan(core.NewTextRange(int(tc.pos), int(tc.pos)))
			assert.Equal(t, got, core.TextPos(span.Pos()))
			assert.Equal(t, fidelity, spanFidelity)
		})
	}
}

func TestMapPositionNilIdentity(t *testing.T) {
	t.Parallel()

	var m *spanmap.SpanMap
	got, fidelity := m.MapPosition(7)
	assert.Equal(t, got, core.TextPos(7))
	assert.Equal(t, fidelity, spanmap.FidelityExact)
}

func TestMapRangeToGeneratedVerbatim(t *testing.T) {
	t.Parallel()

	// Generated [0,10) is a verbatim copy of original [100,110).
	m := spanmap.New([]spanmap.Segment{
		{GenStart: 0, GenEnd: 10, OrigStart: 100, OrigEnd: 110, Kind: spanmap.KindVerbatim},
	})

	got, fidelity := m.MapRangeToGenerated(core.NewTextRange(103, 107))
	assert.Equal(t, got.Pos(), 3)
	assert.Equal(t, got.End(), 7)
	assert.Equal(t, fidelity, spanmap.FidelityExact)
}

func TestMapRangeToGeneratedAtom(t *testing.T) {
	t.Parallel()

	// Generated [3,14) is an atom of the original [60,71).
	m := spanmap.New([]spanmap.Segment{
		{GenStart: 3, GenEnd: 14, OrigStart: 60, OrigEnd: 71, Kind: spanmap.KindAtom},
	})

	// A span inside the original atom maps to the whole generated span.
	got, fidelity := m.MapRangeToGenerated(core.NewTextRange(63, 67))
	assert.Equal(t, got.Pos(), 3)
	assert.Equal(t, got.End(), 14)
	assert.Equal(t, fidelity, spanmap.FidelityAtom)
}

func TestMapRangeToGeneratedGap(t *testing.T) {
	t.Parallel()

	// An original range with no covering segment has no generated counterpart: it maps to the insertion
	// point (the preceding segment's generated end) with no fidelity.
	m := spanmap.New([]spanmap.Segment{
		{GenStart: 0, GenEnd: 10, OrigStart: 100, OrigEnd: 110, Kind: spanmap.KindVerbatim},
		{GenStart: 20, GenEnd: 30, OrigStart: 200, OrigEnd: 210, Kind: spanmap.KindVerbatim},
	})

	got, fidelity := m.MapRangeToGenerated(core.NewTextRange(150, 160))
	assert.Equal(t, got.Pos(), 10)
	assert.Equal(t, got.End(), 10)
	assert.Equal(t, fidelity, spanmap.FidelityNone)
}

func TestMapRangeToGeneratedNilIdentity(t *testing.T) {
	t.Parallel()

	var m *spanmap.SpanMap
	got, fidelity := m.MapRangeToGenerated(core.NewTextRange(3, 7))
	assert.Equal(t, got.Pos(), 3)
	assert.Equal(t, got.End(), 7)
	assert.Equal(t, fidelity, spanmap.FidelityExact)
}

func TestMapPositionToGenerated(t *testing.T) {
	t.Parallel()

	// Original [100,110) is a verbatim copy of generated [0,10); [200,210) is an atom of generated [20,30).
	m := spanmap.New([]spanmap.Segment{
		{GenStart: 0, GenEnd: 10, OrigStart: 100, OrigEnd: 110, Kind: spanmap.KindVerbatim},
		{GenStart: 20, GenEnd: 30, OrigStart: 200, OrigEnd: 210, Kind: spanmap.KindAtom},
	})

	testCases := []struct {
		name     string
		pos      core.TextPos
		want     core.TextPos
		fidelity spanmap.Fidelity
	}{
		{"verbatim interpolates", 103, 3, spanmap.FidelityExact},
		{"atom maps to its start", 205, 20, spanmap.FidelityAtom},
		{"gap maps to insertion point", 150, 10, spanmap.FidelityNone},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, fidelity := m.MapPositionToGenerated(tc.pos)
			assert.Equal(t, got, tc.want)
			assert.Equal(t, fidelity, tc.fidelity)
			// MapPositionToGenerated must agree with MapRangeToGenerated on a zero-length range.
			span, spanFidelity := m.MapRangeToGenerated(core.NewTextRange(int(tc.pos), int(tc.pos)))
			assert.Equal(t, got, core.TextPos(span.Pos()))
			assert.Equal(t, fidelity, spanFidelity)
		})
	}
}

func TestMapRangeToGeneratedRoundTrip(t *testing.T) {
	t.Parallel()

	// Original spans are out of order relative to generated spans, exercising the reverse index.
	m := spanmap.New([]spanmap.Segment{
		{GenStart: 0, GenEnd: 10, OrigStart: 200, OrigEnd: 210, Kind: spanmap.KindVerbatim},
		{GenStart: 10, GenEnd: 20, OrigStart: 100, OrigEnd: 110, Kind: spanmap.KindVerbatim},
	})

	for _, r := range []core.TextRange{core.NewTextRange(2, 8), core.NewTextRange(12, 18)} {
		orig, fidelity := m.MapSpan(r)
		assert.Equal(t, fidelity, spanmap.FidelityExact)
		back, backFidelity := m.MapRangeToGenerated(orig)
		assert.Equal(t, backFidelity, spanmap.FidelityExact)
		assert.Equal(t, back.Pos(), r.Pos())
		assert.Equal(t, back.End(), r.End())
	}
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
